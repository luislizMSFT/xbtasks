package domain

import "time"

type ADOPullRequestStatus string

const (
	ADOPullRequestStatusDraft     = "draft"
	ADOPullRequestStatusActive    = "active"
	ADOPullRequestStatusCompleted = "completed"
	ADOPullRequestStatusAbandoned = "abandoned"
)

type ADOPullRequest struct {
	ID           int                  `json:"id"`
	Title        string               `json:"title"`
	PRURL        string               `json:"prUrl"`
	PRNumber     int                  `json:"prNumber"`
	Repo         string               `json:"repo"`
	TaskID       *int                 `json:"taskId"`
	AdoID        string               `json:"adoId"`
	Status       ADOPullRequestStatus `json:"status"`
	Reviewers    string               `json:"reviewers"`
	SourceBranch string               `json:"sourceBranch"`
	TargetBranch string               `json:"targetBranch"`
	Votes        int                  `json:"votes"`
	CreatedAt    time.Time            `json:"createdAt"`
	UpdatedAt    time.Time            `json:"updatedAt"`
	MergedAt     *time.Time           `json:"mergedAt"`
}

type ADOWorkitemType string

const (
	ADOWorkItemTypeBug         = "bug"
	ADOWorkItemTypeTask        = "task"
	ADOWorkItemTypeUserStory   = "user-story"
	ADOWorkItemTypeDeliverable = "deliverable"
	ADOWorkItemTypeScenario    = "scenario"
)

type ADOWorkItem struct {
	ID          int             `json:"id"`
	AdoID       string          `json:"adoId"`
	Title       string          `json:"title"`
	State       string          `json:"state"`
	Type        ADOWorkitemType `json:"type"`
	AssignedTo  string          `json:"assignedTo"`
	Priority    int             `json:"priority"`
	AreaPath    string          `json:"areaPath"`
	Description string          `json:"description"`
	URL         string          `json:"url"`
	SyncedAt    time.Time       `json:"syncedAt"`
}
