package domain

import "time"

// Pipeline represents a CI/CD pipeline run.
type Pipeline struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Status       string     `json:"status"` // succeeded, failed, running, queued, cancelled
	Result       string     `json:"result"`
	URL          string     `json:"url"`
	SourceBranch string     `json:"sourceBranch"`
	QueueTime    time.Time  `json:"queueTime"`
	FinishTime   *time.Time `json:"finishTime"`
}
