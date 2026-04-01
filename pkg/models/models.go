package models

import "time"

type Task struct {
	ID                 int        `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	Status             string     `json:"status"`   // todo, in_progress, in_review, done, blocked, cancelled
	Priority           string     `json:"priority"` // P0, P1, P2, P3
	Category           string     `json:"category"`
	ProjectID          *int       `json:"projectId"`
	Area               string     `json:"area"`
	DueDate            string     `json:"dueDate"`
	AdoID              string     `json:"adoId"`
	Tags               string     `json:"tags"`
	BlockedReason      string     `json:"blockedReason"`
	BlockedBy          string     `json:"blockedBy"`
	ParentID           *int       `json:"parentId"`
	PersonalPriority   string     `json:"personalPriority"` // override priority for linked ADO items
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	CompletedAt        *time.Time `json:"completedAt"`
}

type TaskDep struct {
	TaskID    int `json:"taskId"`
	DependsOn int `json:"dependsOn"`
}

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // active, paused, completed, archived
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PullRequest struct {
	ID               int        `json:"id"`
	Title            string     `json:"title"`
	PRURL            string     `json:"prUrl"`
	PRNumber         int        `json:"prNumber"`
	Repo             string     `json:"repo"`
	TaskID           *int       `json:"taskId"`
	AdoID            string     `json:"adoId"`
	Status           string     `json:"status"` // draft, active, completed, abandoned
	Reviewers        string     `json:"reviewers"`
	SourceBranch     string     `json:"sourceBranch"`
	TargetBranch     string     `json:"targetBranch"`
	Votes            int        `json:"votes"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	MergedAt         *time.Time `json:"mergedAt"`
}

type ADOWorkItem struct {
	ID          int    `json:"id"`
	AdoID       string `json:"adoId"`
	Title       string `json:"title"`
	State       string `json:"state"`
	Type        string `json:"type"` // Bug, Task, User Story, etc.
	AssignedTo  string `json:"assignedTo"`
	Priority    int    `json:"priority"`
	AreaPath    string `json:"areaPath"`
	Description string `json:"description"`
	URL         string `json:"url"`
	SyncedAt    time.Time `json:"syncedAt"`
}

type TaskADOLink struct {
	TaskID    int    `json:"taskId"`
	AdoID     string `json:"adoId"`
	Direction string `json:"direction"` // promoted, imported, linked
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatarUrl"`
}
