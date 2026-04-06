package ado

import "time"

// OrgProject represents an Azure DevOps organization/project pair.
type OrgProject struct {
	Org     string `json:"org"`
	Project string `json:"project"`
}

// WorkItem represents an Azure DevOps work item with all relevant fields.
type WorkItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
	Type        string    `json:"type"`
	AssignedTo  string    `json:"assignedTo"`
	Priority    int       `json:"priority"`
	AreaPath    string    `json:"areaPath"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	ParentID    int       `json:"parentId"`
	ChangedDate time.Time `json:"changedDate"`
	Org         string    `json:"org"`
	Project     string    `json:"project"`
}

// PatchOperation represents a single JSON Patch operation for ADO work item updates.
type PatchOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}

// WIQLResponse is the response from an ADO WIQL query.
type WIQLResponse struct {
	WorkItems []WIQLWorkItemRef `json:"workItems"`
}

// WIQLWorkItemRef is a work item reference returned by WIQL (ID + URL only).
type WIQLWorkItemRef struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// Identity represents an ADO user identity.
type Identity struct {
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
	ID          string `json:"id"`
}

// Comment represents a work item comment.
type Comment struct {
	ID          int       `json:"id"`
	Text        string    `json:"text"`
	CreatedBy   Identity  `json:"createdBy"`
	CreatedDate time.Time `json:"createdDate"`
}

// FieldDiff represents a difference in a single field between local and remote.
type FieldDiff struct {
	Field    string `json:"field"`
	Local    string `json:"local"`
	Remote   string `json:"remote"`
	Proposed string `json:"proposed"`
}

// SyncDiff captures the changes between a local task and its ADO counterpart.
type SyncDiff struct {
	TaskID    int         `json:"taskId"`
	AdoID     string      `json:"adoId"`
	Changes   []FieldDiff `json:"changes"`
	Direction string      `json:"direction"`
}
