package domain

import "time"

// Task represents a local task.
type Task struct {
	ID               int        `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`   // todo, in_progress, in_review, done, blocked, cancelled
	Priority         string     `json:"priority"` // P0, P1, P2, P3
	Category         string     `json:"category"`
	ProjectID        *int       `json:"projectId"`
	Area             string     `json:"area"`
	DueDate          string     `json:"dueDate"`
	AdoID            string     `json:"adoId"`
	Tags             string     `json:"tags"`
	BlockedReason    string     `json:"blockedReason"`
	BlockedBy        string     `json:"blockedBy"`
	ParentID         *int       `json:"parentId"`
	PersonalPriority string     `json:"personalPriority"` // override priority for linked ADO items
	SortOrder        int        `json:"sortOrder"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CompletedAt      *time.Time `json:"completedAt"`
}

// TaskADOLink represents a link between a task and an ADO work item.
type TaskADOLink struct {
	TaskID    int       `json:"taskId"`
	AdoID     string    `json:"adoId"`
	Direction string    `json:"direction"` // promoted, imported, linked
	CreatedAt time.Time `json:"createdAt"`
}

// TaskDep represents a dependency between two tasks.
type TaskDep struct {
	TaskID    int `json:"taskId"`
	DependsOn int `json:"dependsOn"`
}

// TaskLink represents an external URL attached to a task (LINK-01).
type TaskLink struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"taskId"`
	URL       string    `json:"url"`
	Label     string    `json:"label"`
	Type      string    `json:"type"` // url, icm, grafana, ado, wiki
	CreatedAt time.Time `json:"createdAt"`
}

// TaskComment represents a local comment on a task (CMT-01).
type TaskComment struct {
	ID           int       `json:"id"`
	TaskID       int       `json:"taskId"`
	Content      string    `json:"content"`
	IsPublic     bool      `json:"isPublic"`     // true if pushed to ADO
	AdoCommentID string    `json:"adoCommentId"` // ADO comment ID if pushed
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
