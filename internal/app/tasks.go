package app

import (
	"fmt"
	"log"

	"dev.azure.com/microsoft/Xbox/xb-tasks/internal/db"
	"dev.azure.com/microsoft/Xbox/xb-tasks/pkg/models"
)

// TaskService is bound to the frontend via Wails.
// All task CRUD operations go through here.
type TaskService struct {
	db *db.DB
}

func NewTaskService(database *db.DB) *TaskService {
	return &TaskService{db: database}
}

func (s *TaskService) Create(title, description, priority, category string, projectID *int, parentID *int) (models.Task, error) {
	if title == "" {
		return models.Task{}, fmt.Errorf("title is required")
	}
	if priority == "" {
		priority = "P2"
	}

	res, err := s.db.Exec(
		`INSERT INTO tasks (title, description, priority, category, project_id, parent_id) VALUES (?, ?, ?, ?, ?, ?)`,
		title, description, priority, category, projectID, parentID,
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("create task: %w", err)
	}

	id, _ := res.LastInsertId()
	return s.GetByID(int(id))
}

func (s *TaskService) GetByID(id int) (models.Task, error) {
	var t models.Task
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

func (s *TaskService) List(status string) ([]models.Task, error) {
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

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
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

func (s *TaskService) Update(id int, title, description, status, priority, category, area, dueDate, tags string) (models.Task, error) {
	_, err := s.db.Exec(
		`UPDATE tasks SET title=?, description=?, status=?, priority=?, category=?, area=?, due_date=?, tags=?,
		 updated_at=CURRENT_TIMESTAMP,
		 completed_at = CASE WHEN ? = 'done' THEN CURRENT_TIMESTAMP ELSE completed_at END
		 WHERE id=?`,
		title, description, status, priority, category, area, dueDate, tags, status, id,
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("update task %d: %w", id, err)
	}
	return s.GetByID(id)
}

func (s *TaskService) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (s *TaskService) SetStatus(id int, status string) (models.Task, error) {
	_, err := s.db.Exec(
		`UPDATE tasks SET status=?, updated_at=CURRENT_TIMESTAMP,
		 completed_at = CASE WHEN ? = 'done' THEN CURRENT_TIMESTAMP ELSE completed_at END
		 WHERE id=?`,
		status, status, id,
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("set status task %d: %w", id, err)
	}
	return s.GetByID(id)
}

func (s *TaskService) GetSubtasks(parentID int) ([]models.Task, error) {
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

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
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
