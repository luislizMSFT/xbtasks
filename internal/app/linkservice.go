package app

import (
	"fmt"
	"log"
	"strconv"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/auth"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/pkg/ado"
)

// LinkService manages the link/promote/import/unlink lifecycle between
// local tasks and Azure DevOps work items.
type LinkService struct {
	db        *db.DB
	tokenProv auth.TokenProvider
	cfg       *config.ConfigService
}

// NewLinkService creates a LinkService with database, token provider, and config.
func NewLinkService(database *db.DB, tokenProv auth.TokenProvider, cfg *config.ConfigService) *LinkService {
	return &LinkService{db: database, tokenProv: tokenProv, cfg: cfg}
}

// getClient returns a REST client for the specified org/project, or the first configured pair.
func (s *LinkService) getClient(org, project string) (*ado.Client, error) {
	token, err := s.tokenProv.GetToken()
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	if org != "" && project != "" {
		return ado.NewClient(org, project, token), nil
	}
	return ado.NewDefaultClient(token, config.GetOrgProjects())
}

// getClients returns REST clients for all configured org/project pairs.
func (s *LinkService) getClients() ([]*ado.Client, error) {
	token, err := s.tokenProv.GetToken()
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	return ado.NewClients(token, config.GetOrgProjects())
}

// LinkTask connects an existing local task to an existing ADO work item (D-17).
// Creates task_ado_links entry with direction='linked', updates tasks.ado_id,
// fetches the ADO item into cache, and creates initial sync_state snapshot.
func (s *LinkService) LinkTask(taskID int, adoID string) (domain.TaskADOLink, error) {
	// Verify task exists
	task, err := s.getTask(taskID)
	if err != nil {
		return domain.TaskADOLink{}, fmt.Errorf("task %d not found: %w", taskID, err)
	}

	// Insert link
	_, err = s.db.Exec(
		`INSERT OR IGNORE INTO task_ado_links (task_id, ado_id, direction) VALUES (?, ?, 'linked')`,
		taskID, adoID)
	if err != nil {
		return domain.TaskADOLink{}, fmt.Errorf("link task %d to ADO %s: %w", taskID, adoID, err)
	}

	// Update tasks.ado_id for backward compatibility
	_, err = s.db.Exec(`UPDATE tasks SET ado_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, adoID, taskID)
	if err != nil {
		log.Printf("[link] update tasks.ado_id for %d: %v", taskID, err)
	}

	// Fetch the ADO item and upsert into cache
	adoItem, err := s.fetchADOItem(adoID)
	if err != nil {
		log.Printf("[link] fetch ADO item %s for cache: %v", adoID, err)
	} else {
		if err := s.db.UpsertADOWorkItem(adoItem); err != nil {
			log.Printf("[link] cache ADO item %s: %v", adoID, err)
		}

		// Create initial sync_state snapshot
		localStatus := ado.MapStatusToADO(task.Status, adoItem.Type)
		if err := s.db.UpsertSyncState(taskID, adoID,
			task.Title, task.Status, task.Description,
			adoItem.Title, adoItem.State, adoItem.Description,
		); err != nil {
			log.Printf("[link] upsert sync state for task %d / ADO %s: %v", taskID, adoID, err)
		}
		_ = localStatus // used above in UpsertSyncState call context
	}

	return s.getTaskADOLink(taskID, adoID)
}

// PromoteTask creates a new ADO work item from a local task (D-18).
// Only pushes title/status/description — NOT subtasks, personal priority, or notes (D-09).
func (s *LinkService) PromoteTask(taskID int, wiType string) (domain.TaskADOLink, error) {
	// Fetch local task
	task, err := s.getTask(taskID)
	if err != nil {
		return domain.TaskADOLink{}, fmt.Errorf("task %d not found: %w", taskID, err)
	}

	if wiType == "" {
		wiType = "Task"
	}

	// Get client for first configured org/project
	client, err := s.getClient("", "")
	if err != nil {
		return domain.TaskADOLink{}, fmt.Errorf("get ADO client: %w", err)
	}

	// Create work item in ADO — only title + description (D-09)
	wi, err := ado.CreateWorkItem(client, wiType, task.Title, task.Description)
	if err != nil {
		return domain.TaskADOLink{}, fmt.Errorf("create ADO work item: %w", err)
	}

	adoID := fmt.Sprintf("%d", wi.ID)

	// Insert link
	_, err = s.db.Exec(
		`INSERT INTO task_ado_links (task_id, ado_id, direction) VALUES (?, ?, 'promoted')`,
		taskID, adoID)
	if err != nil {
		return domain.TaskADOLink{}, fmt.Errorf("create promote link: %w", err)
	}

	// Update tasks.ado_id
	_, err = s.db.Exec(`UPDATE tasks SET ado_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, adoID, taskID)
	if err != nil {
		log.Printf("[promote] update tasks.ado_id for %d: %v", taskID, err)
	}

	// Cache the new ADO item
	domainItem := adoWorkItemToDomain(*wi)
	if err := s.db.UpsertADOWorkItem(domainItem); err != nil {
		log.Printf("[promote] cache ADO item %s: %v", adoID, err)
	}

	// Create initial sync_state snapshot
	if err := s.db.UpsertSyncState(taskID, adoID,
		task.Title, task.Status, task.Description,
		wi.Title, wi.State, wi.Description,
	); err != nil {
		log.Printf("[promote] upsert sync state for task %d / ADO %s: %v", taskID, adoID, err)
	}

	return s.getTaskADOLink(taskID, adoID)
}

// ImportWorkItem fetches an ADO work item and creates a local task with a link (D-19).
func (s *LinkService) ImportWorkItem(adoID string) (domain.Task, error) {
	// Fetch the ADO item
	adoItem, err := s.fetchADOItem(adoID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("fetch ADO item %s: %w", adoID, err)
	}

	// Map ADO state to local status
	localStatus := ado.MapADOToStatus(adoItem.State)

	// Create local task
	res, err := s.db.Exec(
		`INSERT INTO tasks (title, description, status, priority, ado_id) VALUES (?, ?, ?, ?, ?)`,
		adoItem.Title, adoItem.Description, localStatus, ado.MapPriorityToLocal(adoItem.Priority), adoID,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create local task from ADO %s: %w", adoID, err)
	}
	taskID64, _ := res.LastInsertId()
	taskID := int(taskID64)

	// Insert link
	_, err = s.db.Exec(
		`INSERT INTO task_ado_links (task_id, ado_id, direction) VALUES (?, ?, 'imported')`,
		taskID, adoID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create import link: %w", err)
	}

	// Cache the ADO item
	if err := s.db.UpsertADOWorkItem(adoItem); err != nil {
		log.Printf("[import] cache ADO item %s: %v", adoID, err)
	}

	// Create sync_state snapshot
	task, err := s.getTask(taskID)
	if err != nil {
		log.Printf("[import] re-fetch task %d for sync state: %v", taskID, err)
	} else if err := s.db.UpsertSyncState(taskID, adoID,
		task.Title, task.Status, task.Description,
		adoItem.Title, adoItem.State, adoItem.Description,
	); err != nil {
		log.Printf("[import] upsert sync state for task %d / ADO %s: %v", taskID, adoID, err)
	}

	return task, nil
}

// ImportWorkItemAsProject fetches an ADO work item and creates a local project
// with a project-ADO link, then imports child work items as tasks on the board.
func (s *LinkService) ImportWorkItemAsProject(adoID string) (domain.Project, error) {
	adoItem, err := s.fetchADOItem(adoID)
	if err != nil {
		return domain.Project{}, fmt.Errorf("fetch ADO item %s: %w", adoID, err)
	}

	// Create local project
	res, err := s.db.Exec(
		`INSERT INTO projects (name, description, status) VALUES (?, ?, 'active')`,
		adoItem.Title, adoItem.Description,
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf("create project from ADO %s: %w", adoID, err)
	}
	projectID64, _ := res.LastInsertId()
	projectID := int(projectID64)

	// Create project-ADO link
	if err := s.db.CreateProjectADOLink(projectID, adoID, "imported"); err != nil {
		log.Printf("[import-project] create project-ADO link: %v", err)
	}

	// Cache the parent ADO item
	if err := s.db.UpsertADOWorkItem(adoItem); err != nil {
		log.Printf("[import-project] cache ADO item %s: %v", adoID, err)
	}

	// Fetch ADO children and import each as a task under this project
	parentInt, _ := strconv.Atoi(adoID)
	clients, clientErr := s.getClients()
	if clientErr != nil {
		log.Printf("[import-project] get clients for children: %v", clientErr)
	} else {
		for _, c := range clients {
			children, err := ado.GetWorkItemChildren(c, parentInt)
			if err != nil {
				continue
			}
			for _, child := range children {
				childAdoID := fmt.Sprintf("%d", child.ID)
				localStatus := ado.MapADOToStatus(child.State)
				taskRes, err := s.db.Exec(
					`INSERT INTO tasks (title, description, status, priority, ado_id, project_id) VALUES (?, ?, ?, ?, ?, ?)`,
					child.Title, child.Description, localStatus, ado.MapPriorityToLocal(child.Priority), childAdoID, projectID,
				)
				if err != nil {
					log.Printf("[import-project] create child task for ADO %s: %v", childAdoID, err)
					continue
				}
				taskID64, _ := taskRes.LastInsertId()
				taskID := int(taskID64)
				s.db.Exec(`INSERT OR IGNORE INTO task_ado_links (task_id, ado_id, direction) VALUES (?, ?, 'imported')`, taskID, childAdoID)

				// Cache each child ADO item
				domainChild := adoWorkItemToDomain(child)
				if err := s.db.UpsertADOWorkItem(domainChild); err != nil {
					log.Printf("[import-project] cache child ADO %s: %v", childAdoID, err)
				}
			}
			break // found children from one client, done
		}
	}

	var project domain.Project
	row := s.db.QueryRow(`SELECT id, name, description, status, is_pinned, created_at, updated_at FROM projects WHERE id = ?`, projectID)
	if err := row.Scan(&project.ID, &project.Name, &project.Description, &project.Status, &project.IsPinned, &project.CreatedAt, &project.UpdatedAt); err != nil {
		return domain.Project{}, fmt.Errorf("read created project: %w", err)
	}
	return project, nil
}

// UnlinkTask disconnects a task from an ADO work item (D-19a).
// If deleteLocal is true, the local task is also deleted.
func (s *LinkService) UnlinkTask(taskID int, adoID string, deleteLocal bool) error {
	// Remove the link
	_, err := s.db.Exec(`DELETE FROM task_ado_links WHERE task_id = ? AND ado_id = ?`, taskID, adoID)
	if err != nil {
		return fmt.Errorf("unlink task %d from ADO %s: %w", taskID, adoID, err)
	}

	// Clear tasks.ado_id if no links remain
	var remaining int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM task_ado_links WHERE task_id = ?`, taskID).Scan(&remaining)
	if err != nil {
		return fmt.Errorf("counting remaining links for task %d: %w", taskID, err)
	}
	if remaining == 0 {
		if _, err := s.db.Exec(`UPDATE tasks SET ado_id = '', updated_at = CURRENT_TIMESTAMP WHERE id = ?`, taskID); err != nil {
			log.Printf("[unlink] clear ado_id for task %d: %v", taskID, err)
		}
	}

	// Remove sync_state
	if _, err := s.db.Exec(`DELETE FROM sync_state WHERE task_id = ? AND ado_id = ?`, taskID, adoID); err != nil {
		log.Printf("[unlink] delete sync_state for task %d / ADO %s: %v", taskID, adoID, err)
	}

	// Optionally delete the local task
	if deleteLocal {
		_, err = s.db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
		if err != nil {
			return fmt.Errorf("delete local task %d: %w", taskID, err)
		}
	}

	return nil
}

// IsPublic returns true if a task has any ADO link (D-05).
// Computed from task_ado_links table presence.
func (s *LinkService) IsPublic(taskID int) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM task_ado_links WHERE task_id = ?`, taskID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check public status for task %d: %w", taskID, err)
	}
	return count > 0, nil
}

// GetTaskLinks returns all ADO links for a task.
func (s *LinkService) GetTaskLinks(taskID int) ([]domain.TaskADOLink, error) {
	rows, err := s.db.Query(
		`SELECT task_id, ado_id, direction, created_at FROM task_ado_links WHERE task_id = ?`, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task links for %d: %w", taskID, err)
	}
	defer rows.Close()

	var links []domain.TaskADOLink
	for rows.Next() {
		var l domain.TaskADOLink
		if err := rows.Scan(&l.TaskID, &l.AdoID, &l.Direction, &l.CreatedAt); err != nil {
			log.Printf("[link] scan link row: %v", err)
			continue
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

// ListPublicTaskIDs returns all task IDs that have at least one ADO link.
func (s *LinkService) ListPublicTaskIDs() ([]int, error) {
	rows, err := s.db.Query(`SELECT DISTINCT task_id FROM task_ado_links`)
	if err != nil {
		return nil, fmt.Errorf("list public task IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ListLinkedAdoIDs returns all ADO IDs that have at least one local task link.
func (s *LinkService) ListLinkedAdoIDs() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT ado_id FROM task_ado_links`)
	if err != nil {
		return nil, fmt.Errorf("list linked ADO IDs: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// --- helpers ---

// getTask fetches a task from the database by ID.
func (s *LinkService) getTask(id int) (domain.Task, error) {
	var t domain.Task
	err := s.db.QueryRow(
		`SELECT id, title, description, status, priority, category, project_id, area, due_date,
		        ado_id, tags, blocked_reason, blocked_by, parent_id, personal_priority,
		        sort_order, created_at, updated_at, completed_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
		&t.Category, &t.ProjectID, &t.Area, &t.DueDate,
		&t.AdoID, &t.Tags, &t.BlockedReason, &t.BlockedBy,
		&t.ParentID, &t.PersonalPriority,
		&t.SortOrder, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
	)
	return t, err
}

// getTaskADOLink returns a specific task-ADO link.
func (s *LinkService) getTaskADOLink(taskID int, adoID string) (domain.TaskADOLink, error) {
	var l domain.TaskADOLink
	err := s.db.QueryRow(
		`SELECT task_id, ado_id, direction, created_at FROM task_ado_links WHERE task_id = ? AND ado_id = ?`,
		taskID, adoID,
	).Scan(&l.TaskID, &l.AdoID, &l.Direction, &l.CreatedAt)
	return l, err
}

// fetchADOItem fetches a single ADO item by string ID, trying all configured org/project clients.
func (s *LinkService) fetchADOItem(adoID string) (domain.WorkItem, error) {
	id, err := strconv.Atoi(adoID)
	if err != nil {
		return domain.WorkItem{}, fmt.Errorf("invalid ADO ID %q: %w", adoID, err)
	}
	clients, err := s.getClients()
	if err != nil {
		return domain.WorkItem{}, err
	}
	for _, c := range clients {
		wi, err := ado.GetWorkItem(c, id)
		if err != nil {
			continue
		}
		return adoWorkItemToDomain(*wi), nil
	}
	return domain.WorkItem{}, fmt.Errorf("ADO item %s not found in any configured org/project", adoID)
}

