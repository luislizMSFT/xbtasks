package domain

import "time"

// SyncState tracks the last known state for conflict detection (SYNC-03).
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

// SyncDiff captures the changes between a local task and its provider counterpart for frontend display.
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

// OrgProject represents a configured ADO org/project pair (ADO-09).
type OrgProject struct {
	Org      string   `json:"org" mapstructure:"org"`
	Projects []string `json:"projects" mapstructure:"projects"`
}
