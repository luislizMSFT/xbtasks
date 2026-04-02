package app

import (
	"fmt"
	"log"

	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/domain"
)

// ProjectService is bound to the frontend via Wails.
type ProjectService struct {
	db *db.DB
}

func NewProjectService(database *db.DB) *ProjectService {
	return &ProjectService{db: database}
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
