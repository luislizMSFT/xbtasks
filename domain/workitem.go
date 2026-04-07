package domain

import "time"

// WorkItem represents a work item fetched from a provider (e.g. Azure DevOps).
type WorkItem struct {
	ID          int       `json:"id"`
	AdoID       string    `json:"adoId"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
	Type        string    `json:"type"`
	AssignedTo  string    `json:"assignedTo"`
	Priority    int       `json:"priority"`
	AreaPath    string    `json:"areaPath"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Org         string    `json:"org"`
	Project     string    `json:"project"`
	ParentID    int       `json:"parentId"`
	ChangedDate time.Time `json:"changedDate"`
	SyncedAt    time.Time `json:"syncedAt"`
}

// ADOComment represents a comment fetched from Azure DevOps.
type ADOComment struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
}
