package domain

import "time"

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatarUrl"`
}

type ProjectStatus string

const (
	ProjectStatusActive    = "active"
	ProjectStatusPaused    = "paused"
	ProjectStatusCompleted = "completed"
	ProjectStatusArchived  = "archived"
	ProjectStatusBlocked   = "blocked"
)

type Project struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"` // active, paused, completed, archived
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

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
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CompletedAt      *time.Time `json:"completedAt"`
}

type TaskADOLink struct {
	TaskID    int       `json:"taskId"`
	AdoID     string    `json:"adoId"`
	Direction string    `json:"direction"` // promoted, imported, linked
	CreatedAt time.Time `json:"createdAt"`
}

type TaskDep struct {
	TaskID    int `json:"taskId"`
	DependsOn int `json:"dependsOn"`
}
