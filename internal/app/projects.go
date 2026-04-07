package app

import (
	"fmt"
	"log"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/auth"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
)

// ProjectService is bound to the frontend via Wails.
type ProjectService struct {
	db        *db.DB
	tokenProv auth.TokenProvider
	cfg       *config.ConfigService
}

func NewProjectService(database *db.DB, tokenProv auth.TokenProvider, cfg *config.ConfigService) *ProjectService {
	return &ProjectService{db: database, tokenProv: tokenProv, cfg: cfg}
}

func (s *ProjectService) Create(name, description string) (domain.Project, error) {
	if name == "" {
		return domain.Project{}, fmt.Errorf("name is required")
	}

	res, err := s.db.Exec(
		`INSERT INTO projects (name, description) VALUES (?, ?)`,
		name, description,
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf("create project: %w", err)
	}

	id, _ := res.LastInsertId()
	return s.GetByID(int(id))
}

func (s *ProjectService) GetByID(id int) (domain.Project, error) {
	var p domain.Project
	err := s.db.QueryRow(
		`SELECT id, name, description, status, created_at, updated_at FROM projects WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, fmt.Errorf("get project %d: %w", id, err)
	}
	return p, nil
}

func (s *ProjectService) List() ([]domain.Project, error) {
	rows, err := s.db.Query(
		`SELECT id, name, description, status, created_at, updated_at FROM projects ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			log.Printf("scan project row: %v", err)
			continue
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (s *ProjectService) Update(id int, name, description, status string) (domain.Project, error) {
	_, err := s.db.Exec(
		`UPDATE projects SET name=?, description=?, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		name, description, status, id,
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf("update project %d: %w", id, err)
	}
	return s.GetByID(id)
}

func (s *ProjectService) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}

// --- ADO Linking (D-20a) ---

// LinkProjectToADO connects a project to an ADO scenario/deliverable.
func (s *ProjectService) LinkProjectToADO(projectID int, adoID, direction string) error {
	if direction == "" {
		direction = "linked"
	}
	return s.db.CreateProjectADOLink(projectID, adoID, direction)
}

// UnlinkProject disconnects a project from an ADO item.
// If deleteLocal is true, the local project is also deleted.
func (s *ProjectService) UnlinkProject(projectID int, adoID string, deleteLocal bool) error {
	if err := s.db.DeleteProjectADOLink(projectID, adoID); err != nil {
		return fmt.Errorf("unlink project %d from ADO %s: %w", projectID, adoID, err)
	}
	if deleteLocal {
		return s.Delete(projectID)
	}
	return nil
}

// GetProjectADOLink returns the ADO link for a project.
func (s *ProjectService) GetProjectADOLink(projectID int) (domain.ProjectADOLink, error) {
	return s.db.GetProjectADOLink(projectID)
}

// PinProject sets or clears the pinned/starred flag on a project.
func (s *ProjectService) PinProject(projectID int, pinned bool) error {
	val := 0
	if pinned {
		val = 1
	}
	_, err := s.db.Exec(`UPDATE projects SET is_pinned = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, val, projectID)
	if err != nil {
		return fmt.Errorf("pin project %d: %w", projectID, err)
	}
	return nil
}

// GetProjectProgress returns task completion counts for a project (D-PROJ-06).
// Includes both local task counts and (if linked) ADO children counts from cache.
func (s *ProjectService) GetProjectProgress(projectID int) (map[string]any, error) {
	// Local task progress (single query with conditional count)
	var localDone, localTotal int
	err := s.db.QueryRow(
		`SELECT COUNT(*) AS total,
		        COUNT(CASE WHEN status = 'done' THEN 1 END) AS done
		 FROM tasks WHERE project_id = ?`, projectID,
	).Scan(&localTotal, &localDone)
	if err != nil {
		return nil, fmt.Errorf("count project tasks: %w", err)
	}

	result := map[string]any{
		"localDone":  localDone,
		"localTotal": localTotal,
		"adoDone":    0,
		"adoTotal":   0,
	}

	// Check if project is linked to ADO
	link, err := s.db.GetProjectADOLink(projectID)
	if err != nil {
		// Real DB failure — return local-only progress
		return result, nil
	}
	if link.AdoID == "" {
		// No ADO link exists for this project
		return result, nil
	}

	// Count ADO children from cache (single query with conditional count)
	var adoTotal, adoDone int
	err = s.db.QueryRow(
		`SELECT COUNT(*) AS total,
		        COUNT(CASE WHEN state IN ('Closed','Completed') THEN 1 END) AS done
		 FROM ado_work_items WHERE parent_id = CAST(? AS INTEGER)`, link.AdoID,
	).Scan(&adoTotal, &adoDone)
	if err != nil {
		log.Printf("[project] count ADO children for %s: %v", link.AdoID, err)
		return result, nil
	}

	result["adoDone"] = adoDone
	result["adoTotal"] = adoTotal
	return result, nil
}
