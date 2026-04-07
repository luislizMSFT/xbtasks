package domain

import "time"

// PullRequest represents a pull request tracked locally.
type PullRequest struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	PRURL        string     `json:"prUrl"`
	PRNumber     int        `json:"prNumber"`
	Repo         string     `json:"repo"`
	TaskID       *int       `json:"taskId"`
	AdoID        string     `json:"adoId"`
	Status       string     `json:"status"` // draft, active, completed, abandoned
	Reviewers    string     `json:"reviewers"` // JSON array of reviewer objects
	SourceBranch string     `json:"sourceBranch"`
	TargetBranch string     `json:"targetBranch"`
	Votes        int        `json:"votes"`
	CreatedBy    string     `json:"createdBy"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	MergedAt     *time.Time `json:"mergedAt"`
}
