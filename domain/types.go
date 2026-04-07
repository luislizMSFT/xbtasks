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
	IsPinned    bool          `json:"isPinned"`
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
	SortOrder        int        `json:"sortOrder"`
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

type ADOPipeline struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Status       string     `json:"status"` // succeeded, failed, running, queued, cancelled
	Result       string     `json:"result"`
	URL          string     `json:"url"`
	SourceBranch string     `json:"sourceBranch"`
	QueueTime    time.Time  `json:"queueTime"`
	FinishTime   *time.Time `json:"finishTime"`
}

type ADOWorkItem struct {
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

// TaskLink represents an external URL attached to a task (LINK-01)
type TaskLink struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"taskId"`
	URL       string    `json:"url"`
	Label     string    `json:"label"`
	Type      string    `json:"type"` // url, icm, grafana, ado, wiki
	CreatedAt time.Time `json:"createdAt"`
}

// TaskComment represents a local comment on a task (CMT-01)
type TaskComment struct {
	ID           int       `json:"id"`
	TaskID       int       `json:"taskId"`
	Content      string    `json:"content"`
	IsPublic     bool      `json:"isPublic"`    // true if pushed to ADO
	AdoCommentID string    `json:"adoCommentId"` // ADO comment ID if pushed
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ProjectADOLink represents a project linked to an ADO scenario/deliverable (PROJ-04)
type ProjectADOLink struct {
	ProjectID int       `json:"projectId"`
	AdoID     string    `json:"adoId"`
	Direction string    `json:"direction"` // promoted, imported, linked
	CreatedAt time.Time `json:"createdAt"`
}

// SyncState tracks the last known state for conflict detection (SYNC-03)
type SyncState struct {
	TaskID       int       `json:"taskId"`
	AdoID        string    `json:"adoId"`
	LastSyncedAt time.Time `json:"lastSyncedAt"`
	LocalTitle   string    `json:"localTitle"`
	LocalStatus  string    `json:"localStatus"`
	LocalDesc    string    `json:"localDesc"`
	RemoteTitle  string    `json:"remoteTitle"`
	RemoteStatus string    `json:"remoteStatus"`
	RemoteDesc   string    `json:"remoteDesc"`
}

// SyncDiff captures the changes between a local task and its ADO counterpart for frontend display.
type SyncDiff struct {
	TaskID    int         `json:"taskId"`
	AdoID     string      `json:"adoId"`
	Changes   []FieldDiff `json:"changes"`
	Direction string      `json:"direction"` // inbound, outbound, conflict
}

// FieldDiff represents a difference in a single field between local and remote.
type FieldDiff struct {
	Field    string `json:"field"`
	Local    string `json:"local"`
	Remote   string `json:"remote"`
	Proposed string `json:"proposed"`
}

// ADOComment represents a comment fetched from Azure DevOps.
type ADOComment struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
}

// OrgProject represents a configured ADO org/project pair (ADO-09)
type OrgProject struct {
	Org      string   `json:"org" mapstructure:"org"`
	Projects []string `json:"projects" mapstructure:"projects"`
}
