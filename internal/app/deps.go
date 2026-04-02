package app

import (
	"fmt"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/db"
)

// DependencyService manages task-to-task dependencies with cycle detection.
type DependencyService struct {
	db *db.DB
}

func NewDependencyService(database *db.DB) *DependencyService {
	return &DependencyService{db: database}
}

// AddDependency creates a dependency: taskID depends on dependsOn.
// Returns an error if this would create a circular dependency chain.
func (s *DependencyService) AddDependency(taskID, dependsOn int) error {
	if taskID == dependsOn {
		return fmt.Errorf("a task cannot depend on itself")
	}

	// Verify both tasks exist
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM tasks WHERE id IN (?, ?)`, taskID, dependsOn,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("check tasks exist: %w", err)
	}
	if count != 2 {
		return fmt.Errorf("one or both tasks not found (task_id=%d, depends_on=%d)", taskID, dependsOn)
	}

	// Check for circular dependency before inserting
	if s.hasCircularDep(taskID, dependsOn) {
		return fmt.Errorf("adding this dependency would create a circular chain: task %d already depends on task %d", dependsOn, taskID)
	}

	_, err = s.db.Exec(
		`INSERT OR IGNORE INTO task_deps (task_id, depends_on) VALUES (?, ?)`,
		taskID, dependsOn,
	)
	if err != nil {
		return fmt.Errorf("add dependency: %w", err)
	}
	return nil
}

// hasCircularDep performs DFS from dependsOn following task_deps edges.
// If we can reach taskID, adding this edge would create a cycle.
func (s *DependencyService) hasCircularDep(taskID, dependsOn int) bool {
	visited := map[int]bool{}
	stack := []int{dependsOn}

	for len(stack) > 0 {
		// Pop
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[current] {
			continue
		}
		visited[current] = true

		rows, err := s.db.Query(
			`SELECT depends_on FROM task_deps WHERE task_id = ?`, current,
		)
		if err != nil {
			continue
		}

		for rows.Next() {
			var dep int
			if err := rows.Scan(&dep); err != nil {
				continue
			}
			if dep == taskID {
				rows.Close()
				return true
			}
			if !visited[dep] {
				stack = append(stack, dep)
			}
		}
		rows.Close()
	}

	return false
}

// RemoveDependency removes a dependency relationship.
func (s *DependencyService) RemoveDependency(taskID, dependsOn int) error {
	_, err := s.db.Exec(
		`DELETE FROM task_deps WHERE task_id = ? AND depends_on = ?`,
		taskID, dependsOn,
	)
	if err != nil {
		return fmt.Errorf("remove dependency: %w", err)
	}
	return nil
}

// GetDependencies returns all tasks that taskID depends on.
func (s *DependencyService) GetDependencies(taskID int) ([]domain.Task, error) {
	rows, err := s.db.Query(
		`SELECT t.id, t.title, t.description, t.status, t.priority, t.category,
		        t.project_id, t.area, t.due_date, t.ado_id, t.tags, t.blocked_reason,
		        t.blocked_by, t.parent_id, t.personal_priority,
		        t.created_at, t.updated_at, t.completed_at
		 FROM tasks t
		 JOIN task_deps td ON t.id = td.depends_on
		 WHERE td.task_id = ?`, taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("get dependencies for task %d: %w", taskID, err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// GetBlockedBy returns all tasks that are blocked by taskID (reverse direction).
func (s *DependencyService) GetBlockedBy(taskID int) ([]domain.Task, error) {
	rows, err := s.db.Query(
		`SELECT t.id, t.title, t.description, t.status, t.priority, t.category,
		        t.project_id, t.area, t.due_date, t.ado_id, t.tags, t.blocked_reason,
		        t.blocked_by, t.parent_id, t.personal_priority,
		        t.created_at, t.updated_at, t.completed_at
		 FROM tasks t
		 JOIN task_deps td ON t.id = td.task_id
		 WHERE td.depends_on = ?`, taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("get blocked by task %d: %w", taskID, err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// scanTasks is a helper to scan rows into a slice of domain.Task.
func scanTasks(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]domain.Task, error) {
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
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
