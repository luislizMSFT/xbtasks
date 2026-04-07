package domain

import "time"

// ProjectStatus represents the lifecycle state of a project.
type ProjectStatus string

const (
	ProjectStatusActive    = "active"
	ProjectStatusPaused    = "paused"
	ProjectStatusCompleted = "completed"
	ProjectStatusArchived  = "archived"
	ProjectStatusBlocked   = "blocked"
)

// Project represents a local project.
type Project struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"` // active, paused, completed, archived
	IsPinned    bool          `json:"isPinned"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// ProjectProviderLink represents a project linked to a provider work item (PROJ-04).
type ProjectProviderLink struct {
	ProjectID      int       `json:"projectId"`
	ProviderItemID string    `json:"adoId"`     // json kept as "adoId" for frontend backward compat
	Direction      string    `json:"direction"` // promoted, imported, linked
	CreatedAt      time.Time `json:"createdAt"`
}
