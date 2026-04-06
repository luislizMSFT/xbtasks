package app

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/db"
)

// TaskService is bound to the frontend via Wails.
// All task CRUD operations go through here.
type TaskService struct {
	db *db.DB
}

func NewTaskService(database *db.DB) *TaskService {
	return &TaskService{db: database}
}

func (s *TaskService) Create(title, description, priority, category string, projectID *int, parentID *int) (domain.Task, error) {
	if title == "" {
		return domain.Task{}, fmt.Errorf("title is required")
	}
	if priority == "" {
		priority = "P2"
	}

	res, err := s.db.Exec(
		`INSERT INTO tasks (title, description, priority, category, project_id, parent_id) VALUES (?, ?, ?, ?, ?, ?)`,
		title, description, priority, category, projectID, parentID,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	id, _ := res.LastInsertId()
	return s.GetByID(int(id))
}

func (s *TaskService) GetByID(id int) (domain.Task, error) {
	var t domain.Task
	err := s.db.QueryRow(
		`SELECT id, title, description, status, priority, category, project_id, area, due_date,
		        ado_id, tags, blocked_reason, blocked_by, parent_id, personal_priority,
		        created_at, updated_at, completed_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
		&t.Category, &t.ProjectID, &t.Area, &t.DueDate,
		&t.AdoID, &t.Tags, &t.BlockedReason, &t.BlockedBy,
		&t.ParentID, &t.PersonalPriority,
		&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
	)
	if err != nil {
		return t, fmt.Errorf("get task %d: %w", id, err)
	}
	return t, nil
}

func (s *TaskService) List(status string) ([]domain.Task, error) {
	query := `SELECT id, title, description, status, priority, category, project_id, area, due_date,
	                  ado_id, tags, blocked_reason, blocked_by, parent_id, personal_priority,
	                  created_at, updated_at, completed_at
	           FROM tasks`
	var args []any
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY CASE priority WHEN 'P0' THEN 0 WHEN 'P1' THEN 1 WHEN 'P2' THEN 2 WHEN 'P3' THEN 3 END, updated_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.Category, &t.ProjectID, &t.Area, &t.DueDate,
			&t.AdoID, &t.Tags, &t.BlockedReason, &t.BlockedBy,
			&t.ParentID, &t.PersonalPriority,
			&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
		); err != nil {
			log.Printf("scan task row: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *TaskService) Update(id int, title, description, status, priority, category, area, dueDate, tags string) (domain.Task, error) {
	_, err := s.db.Exec(
		`UPDATE tasks SET title=?, description=?, status=?, priority=?, category=?, area=?, due_date=?, tags=?,
		 updated_at=CURRENT_TIMESTAMP,
		 completed_at = CASE WHEN ? = 'done' THEN CURRENT_TIMESTAMP ELSE completed_at END
		 WHERE id=?`,
		title, description, status, priority, category, area, dueDate, tags, status, id,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task %d: %w", id, err)
	}
	return s.GetByID(id)
}

func (s *TaskService) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (s *TaskService) SetStatus(id int, status string) (domain.Task, error) {
	_, err := s.db.Exec(
		`UPDATE tasks SET status=?, updated_at=CURRENT_TIMESTAMP,
		 completed_at = CASE WHEN ? = 'done' THEN CURRENT_TIMESTAMP ELSE completed_at END
		 WHERE id=?`,
		status, status, id,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("set status task %d: %w", id, err)
	}
	return s.GetByID(id)
}

func (s *TaskService) GetSubtasks(parentID int) ([]domain.Task, error) {
	rows, err := s.db.Query(
		`SELECT id, title, description, status, priority, category, project_id, area, due_date,
		        ado_id, tags, blocked_reason, blocked_by, parent_id, personal_priority,
		        created_at, updated_at, completed_at
		 FROM tasks WHERE parent_id = ?
		 ORDER BY CASE priority WHEN 'P0' THEN 0 WHEN 'P1' THEN 1 WHEN 'P2' THEN 2 WHEN 'P3' THEN 3 END`, parentID,
	)
	if err != nil {
		return nil, fmt.Errorf("get subtasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.Category, &t.ProjectID, &t.Area, &t.DueDate,
			&t.AdoID, &t.Tags, &t.BlockedReason, &t.BlockedBy,
			&t.ParentID, &t.PersonalPriority,
			&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
		); err != nil {
			log.Printf("scan subtask row: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// SetPersonalPriority updates the personal priority overlay independent of ADO priority.
func (s *TaskService) SetPersonalPriority(id int, personalPriority string) (domain.Task, error) {
	_, err := s.db.Exec(
		`UPDATE tasks SET personal_priority=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		personalPriority, id,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("set personal priority task %d: %w", id, err)
	}
	return s.GetByID(id)
}

// CreateSubtask creates a child task under the given parent.
func (s *TaskService) CreateSubtask(parentID int, title, description, priority string) (domain.Task, error) {
	if title == "" {
		return domain.Task{}, fmt.Errorf("title is required")
	}
	if priority == "" {
		priority = "P2"
	}
	// Verify parent exists
	_, err := s.GetByID(parentID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parent task %d not found: %w", parentID, err)
	}
	return s.Create(title, description, priority, "", nil, &parentID)
}

// ListFiltered returns tasks matching the provided filters.
// Empty filter values are ignored (all-pass).
func (s *TaskService) ListFiltered(status, projectID, parentID, tag string) ([]domain.Task, error) {
	query := `SELECT id, title, description, status, priority, category, project_id, area, due_date,
	                  ado_id, tags, blocked_reason, blocked_by, parent_id, personal_priority,
	                  created_at, updated_at, completed_at
	           FROM tasks WHERE 1=1`
	var args []any
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if projectID != "" {
		query += " AND project_id = ?"
		args = append(args, projectID)
	}
	if parentID != "" {
		query += " AND parent_id = ?"
		args = append(args, parentID)
	}
	if tag != "" {
		query += " AND (',' || tags || ',') LIKE '%,' || ? || ',%'"
		args = append(args, tag)
	}
	query += " ORDER BY CASE priority WHEN 'P0' THEN 0 WHEN 'P1' THEN 1 WHEN 'P2' THEN 2 WHEN 'P3' THEN 3 END, updated_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list filtered tasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.Category, &t.ProjectID, &t.Area, &t.DueDate,
			&t.AdoID, &t.Tags, &t.BlockedReason, &t.BlockedBy,
			&t.ParentID, &t.PersonalPriority,
			&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
		); err != nil {
			log.Printf("scan filtered task row: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// ---------------------------------------------------------------------------
// ADO link read operations — LOCAL ONLY
// Link/unlink mutations live in LinkService (linkservice.go), which handles
// sync state, ADO item caching, and full lifecycle management.
// ---------------------------------------------------------------------------

// GetADOLinks returns all ADO links for a task (local read only).
func (s *TaskService) GetADOLinks(taskID int) ([]domain.TaskADOLink, error) {
	rows, err := s.db.Query(
		`SELECT task_id, ado_id, direction, created_at FROM task_ado_links WHERE task_id = ?`, taskID)
	if err != nil {
		return nil, fmt.Errorf("get ADO links for task %d: %w", taskID, err)
	}
	defer rows.Close()

	var links []domain.TaskADOLink
	for rows.Next() {
		var l domain.TaskADOLink
		if err := rows.Scan(&l.TaskID, &l.AdoID, &l.Direction, &l.CreatedAt); err != nil {
			log.Printf("scan ADO link row: %v", err)
			continue
		}
		links = append(links, l)
	}
	return links, nil
}

// GetAllTags returns all unique tags used across tasks, sorted alphabetically.
func (s *TaskService) GetAllTags() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT tags FROM tasks WHERE tags != ''`)
	if err != nil {
		return nil, fmt.Errorf("get all tags: %w", err)
	}
	defer rows.Close()

	tagSet := make(map[string]bool)
	for rows.Next() {
		var tags string
		if err := rows.Scan(&tags); err != nil {
			continue
		}
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tagSet[t] = true
			}
		}
	}
	result := make([]string, 0, len(tagSet))
	for t := range tagSet {
		result = append(result, t)
	}
	sort.Strings(result)
	return result, nil
}
