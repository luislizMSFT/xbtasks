package app

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"dev.azure.com/xbox/xb-tasks/domain"
	"dev.azure.com/xbox/xb-tasks/internal/auth"
	"dev.azure.com/xbox/xb-tasks/internal/config"
	"dev.azure.com/xbox/xb-tasks/internal/db"
	"dev.azure.com/xbox/xb-tasks/pkg/ado"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// SyncService manages bidirectional sync between local tasks and Azure DevOps work items.
// Inbound: auto-pulls ADO changes on a configurable timer (silent merge).
// Outbound: generates preview diffs and pushes only after user confirmation.
type SyncService struct {
	db        *db.DB
	tokenProv auth.TokenProvider
	cfg       *config.ConfigService
	app       *application.App
	stopCh    chan struct{}
}

// NewSyncService creates a SyncService with all required dependencies.
func NewSyncService(database *db.DB, tokenProv auth.TokenProvider, cfg *config.ConfigService, app *application.App) *SyncService {
	return &SyncService{db: database, tokenProv: tokenProv, cfg: cfg, app: app}
}

// StartBackgroundSync launches a goroutine that pulls ADO changes on a configurable interval.
// Default interval is 15 minutes (from sync.interval_minutes config).
func (s *SyncService) StartBackgroundSync() {
	interval := config.SyncIntervalMinutes()
	if interval <= 0 {
		interval = 15
	}

	s.stopCh = make(chan struct{})
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)

	log.Printf("[sync] background sync started (every %d min)", interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				diffs, err := s.pullChanges()
				if err != nil {
					log.Printf("[sync] background pull error: %v", err)
				} else {
					conflicts := 0
					for _, d := range diffs {
						if d.Direction == "conflict" {
							conflicts++
						}
					}
					log.Printf("[sync] background pull complete: %d updated, %d conflicts", len(diffs)-conflicts, conflicts)
				}
			case <-s.stopCh:
				ticker.Stop()
				log.Printf("[sync] background sync stopped")
				return
			}
		}
	}()
}

// StopSync stops the background sync goroutine.
func (s *SyncService) StopSync() {
	if s.stopCh != nil {
		close(s.stopCh)
		s.stopCh = nil
	}
}

// ManualSync triggers an immediate sync and returns any diffs/conflicts found.
// Called when user clicks the sync button.
func (s *SyncService) ManualSync() ([]domain.SyncDiff, error) {
	diffs, err := s.pullChanges()
	if err != nil {
		return nil, fmt.Errorf("manual sync: %w", err)
	}

	conflicts := 0
	updated := 0
	for _, d := range diffs {
		if d.Direction == "conflict" {
			conflicts++
		} else {
			updated++
		}
	}

	// Emit sync:completed event to frontend
	if s.app != nil {
		s.app.Event.Emit("sync:completed", map[string]int{
			"conflicts": conflicts,
			"updated":   updated,
		})
	}

	return diffs, nil
}

// pullChanges performs the inbound sync: fetches current ADO state for all linked items,
// auto-merges remote-only changes, and flags conflicts where both local and remote changed.
func (s *SyncService) pullChanges() ([]domain.SyncDiff, error) {
	// Get all sync states (linked items)
	states, err := s.db.ListSyncStates()
	if err != nil {
		return nil, fmt.Errorf("list sync states: %w", err)
	}
	if len(states) == 0 {
		return nil, nil
	}

	// Get token
	token, err := s.tokenProv.GetToken()
	if err != nil {
		return nil, fmt.Errorf("get token for sync: %w", err)
	}

	// Build map of adoID → sync state for quick lookup
	stateMap := make(map[string]domain.SyncState)
	adoIDs := make([]int, 0, len(states))
	for _, st := range states {
		stateMap[st.AdoID] = st
		id, err := strconv.Atoi(st.AdoID)
		if err != nil {
			log.Printf("[sync] invalid ADO ID %q, skipping", st.AdoID)
			continue
		}
		adoIDs = append(adoIDs, id)
	}

	if len(adoIDs) == 0 {
		return nil, nil
	}

	// Fetch current ADO state for all linked items across all org/project pairs
	orgProjects := config.GetOrgProjects()
	var remoteItems []ado.WorkItem
	for _, op := range orgProjects {
		for _, proj := range op.Projects {
			c := ado.NewClient(op.Org, proj, token)
			items, err := ado.GetWorkItemsByIDs(c, adoIDs)
			if err != nil {
				log.Printf("[sync] fetch from %s/%s failed: %v", op.Org, proj, err)
				continue
			}
			remoteItems = append(remoteItems, items...)
		}
	}

	// Build remoteMap: adoID string → WorkItem
	remoteMap := make(map[string]ado.WorkItem)
	for _, wi := range remoteItems {
		remoteMap[strconv.Itoa(wi.ID)] = wi
	}

	var diffs []domain.SyncDiff

	for adoID, syncState := range stateMap {
		remote, ok := remoteMap[adoID]
		if !ok {
			continue // item not found in any org — skip
		}

		// Get current local task
		localTask, err := s.getTask(syncState.TaskID)
		if err != nil {
			log.Printf("[sync] task %d not found: %v", syncState.TaskID, err)
			continue
		}

		// Detect changes since last sync
		localChanged := s.hasLocalChanged(localTask, syncState)
		remoteChanged := s.hasRemoteChanged(remote, syncState)

		if !localChanged && !remoteChanged {
			continue // nothing changed
		}

		if remoteChanged && !localChanged {
			// Remote-only changes → auto-merge silently (D-21)
			diff := s.autoMergeInbound(localTask, remote, syncState)
			if len(diff.Changes) > 0 {
				diffs = append(diffs, diff)
			}
		} else if localChanged && remoteChanged {
			// Both changed → conflict (D-26)
			diff := s.detectConflictFields(localTask, remote, syncState)
			if len(diff.Changes) > 0 {
				diffs = append(diffs, diff)
				// Emit conflict event
				if s.app != nil {
					s.app.Event.Emit("sync:conflict", map[string]any{
						"taskId": syncState.TaskID,
						"adoId":  adoID,
						"fields": len(diff.Changes),
					})
				}
			}
		}
		// localChanged && !remoteChanged → no action (outbound is manual)
	}

	return diffs, nil
}

// GenerateOutboundDiff previews what will be pushed to ADO before user confirms (D-22).
// Returns a SyncDiff with direction="outbound" showing Local vs Remote vs Proposed for each field.
func (s *SyncService) GenerateOutboundDiff(taskID int) (domain.SyncDiff, error) {
	// Get local task
	task, err := s.getTask(taskID)
	if err != nil {
		return domain.SyncDiff{}, fmt.Errorf("task %d not found: %w", taskID, err)
	}

	// Get ADO link
	adoID, err := s.getLinkedAdoID(taskID)
	if err != nil {
		return domain.SyncDiff{}, fmt.Errorf("no ADO link for task %d: %w", taskID, err)
	}

	// Fetch current remote state
	remote, client, err := s.fetchRemoteItem(adoID)
	if err != nil {
		return domain.SyncDiff{}, fmt.Errorf("fetch ADO item %s: %w", adoID, err)
	}
	_ = client // needed for push, not for diff

	diff := domain.SyncDiff{
		TaskID:    taskID,
		AdoID:     adoID,
		Direction: "outbound",
	}

	// Compare title
	if task.Title != remote.Title {
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:    "title",
			Local:    task.Title,
			Remote:   remote.Title,
			Proposed: task.Title,
		})
	}

	// Compare status (mapped)
	mappedStatus := ado.MapStatusToADO(task.Status, remote.Type)
	if mappedStatus != remote.State {
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:    "state",
			Local:    task.Status,
			Remote:   remote.State,
			Proposed: mappedStatus,
		})
	}

	// Compare description
	if task.Description != remote.Description {
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:    "description",
			Local:    task.Description,
			Remote:   remote.Description,
			Proposed: task.Description,
		})
	}

	return diff, nil
}

// PushChanges pushes local task changes to ADO after user confirms the diff (D-22/D-23).
func (s *SyncService) PushChanges(taskID int) error {
	// Get local task
	task, err := s.getTask(taskID)
	if err != nil {
		return fmt.Errorf("task %d not found: %w", taskID, err)
	}

	// Get ADO link
	adoID, err := s.getLinkedAdoID(taskID)
	if err != nil {
		return fmt.Errorf("no ADO link for task %d: %w", taskID, err)
	}

	// Fetch remote to get the work item type for status mapping
	remote, client, err := s.fetchRemoteItem(adoID)
	if err != nil {
		return fmt.Errorf("fetch ADO item %s: %w", adoID, err)
	}

	// Build fields map
	fields := map[string]string{
		"title":       task.Title,
		"description": task.Description,
		"state":       ado.MapStatusToADO(task.Status, remote.Type),
	}

	adoIDInt, _ := strconv.Atoi(adoID)
	_, err = ado.UpdateWorkItemFields(client, adoIDInt, fields)
	if err != nil {
		return fmt.Errorf("push changes for task %d to ADO %s: %w", taskID, adoID, err)
	}

	// Update sync_state with new snapshot
	mappedState := ado.MapStatusToADO(task.Status, remote.Type)
	err = s.db.UpsertSyncState(taskID, adoID,
		task.Title, task.Status, task.Description,
		task.Title, mappedState, task.Description,
	)
	if err != nil {
		log.Printf("[sync] update sync state after push for task %d: %v", taskID, err)
	}

	log.Printf("[sync] pushed changes for task %d to ADO %s", taskID, adoID)
	return nil
}

// ResolveConflict applies per-field conflict resolution (D-26/D-27).
// resolutions maps field name → "local" or "remote".
func (s *SyncService) ResolveConflict(taskID int, resolutions map[string]string) error {
	task, err := s.getTask(taskID)
	if err != nil {
		return fmt.Errorf("task %d not found: %w", taskID, err)
	}

	adoID, err := s.getLinkedAdoID(taskID)
	if err != nil {
		return fmt.Errorf("no ADO link for task %d: %w", taskID, err)
	}

	remote, _, err := s.fetchRemoteItem(adoID)
	if err != nil {
		return fmt.Errorf("fetch ADO item %s: %w", adoID, err)
	}

	// Apply resolutions
	newTitle := task.Title
	newStatus := task.Status
	newDesc := task.Description

	for field, choice := range resolutions {
		switch field {
		case "title":
			if choice == "remote" {
				newTitle = remote.Title
			}
		case "state":
			if choice == "remote" {
				newStatus = ado.MapADOToStatus(remote.State)
			}
		case "description":
			if choice == "remote" {
				newDesc = remote.Description
			}
		}
	}

	// Update local task with resolved values
	_, err = s.db.Exec(
		`UPDATE tasks SET title=?, description=?, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		newTitle, newDesc, newStatus, taskID,
	)
	if err != nil {
		return fmt.Errorf("update task %d after conflict resolution: %w", taskID, err)
	}

	// Update sync_state with new snapshot
	mappedState := ado.MapStatusToADO(newStatus, remote.Type)
	err = s.db.UpsertSyncState(taskID, adoID,
		newTitle, newStatus, newDesc,
		remote.Title, remote.State, remote.Description,
	)
	if err != nil {
		log.Printf("[sync] update sync state after resolution for task %d: %v", taskID, err)
	}

	// If any fields resolved to "local", user may want to push those changes
	_ = mappedState
	log.Printf("[sync] resolved conflict for task %d (ADO %s): %v", taskID, adoID, resolutions)
	return nil
}

// DetectConflicts checks a specific task for field-level conflicts between local and remote.
func (s *SyncService) DetectConflicts(taskID int) ([]domain.FieldDiff, error) {
	task, err := s.getTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("task %d not found: %w", taskID, err)
	}

	adoID, err := s.getLinkedAdoID(taskID)
	if err != nil {
		return nil, fmt.Errorf("no ADO link for task %d: %w", taskID, err)
	}

	syncState, err := s.db.GetSyncState(taskID, adoID)
	if err != nil {
		return nil, fmt.Errorf("get sync state for task %d: %w", taskID, err)
	}

	remote, _, err := s.fetchRemoteItem(adoID)
	if err != nil {
		return nil, fmt.Errorf("fetch ADO item %s: %w", adoID, err)
	}

	diff := s.detectConflictFields(task, *remote, syncState)
	return diff.Changes, nil
}

// --- internal helpers ---

// hasLocalChanged checks if the local task has changed since the last sync snapshot.
func (s *SyncService) hasLocalChanged(task domain.Task, syncState domain.SyncState) bool {
	return task.Title != syncState.LocalTitle ||
		task.Status != syncState.LocalStatus ||
		task.Description != syncState.LocalDesc
}

// hasRemoteChanged checks if the remote ADO item has changed since the last sync snapshot.
func (s *SyncService) hasRemoteChanged(remote ado.WorkItem, syncState domain.SyncState) bool {
	return remote.Title != syncState.RemoteTitle ||
		remote.State != syncState.RemoteStatus ||
		remote.Description != syncState.RemoteDesc
}

// autoMergeInbound silently applies remote-only changes to the local task.
func (s *SyncService) autoMergeInbound(task domain.Task, remote ado.WorkItem, syncState domain.SyncState) domain.SyncDiff {
	diff := domain.SyncDiff{
		TaskID:    task.ID,
		AdoID:     syncState.AdoID,
		Direction: "inbound",
	}

	newTitle := task.Title
	newStatus := task.Status
	newDesc := task.Description

	if remote.Title != syncState.RemoteTitle {
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:    "title",
			Local:    task.Title,
			Remote:   remote.Title,
			Proposed: remote.Title,
		})
		newTitle = remote.Title
	}

	if remote.State != syncState.RemoteStatus {
		localStatus := ado.MapADOToStatus(remote.State)
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:    "state",
			Local:    task.Status,
			Remote:   remote.State,
			Proposed: localStatus,
		})
		newStatus = localStatus
	}

	if remote.Description != syncState.RemoteDesc {
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:    "description",
			Local:    task.Description,
			Remote:   remote.Description,
			Proposed: remote.Description,
		})
		newDesc = remote.Description
	}

	if len(diff.Changes) > 0 {
		// Update local task
		_, err := s.db.Exec(
			`UPDATE tasks SET title=?, description=?, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
			newTitle, newDesc, newStatus, task.ID,
		)
		if err != nil {
			log.Printf("[sync] auto-merge update task %d failed: %v", task.ID, err)
		}

		// Update sync_state snapshot
		err = s.db.UpsertSyncState(task.ID, syncState.AdoID,
			newTitle, newStatus, newDesc,
			remote.Title, remote.State, remote.Description,
		)
		if err != nil {
			log.Printf("[sync] update sync state for task %d: %v", task.ID, err)
		}
	}

	return diff
}

// detectConflictFields identifies fields where both local and remote have diverged.
func (s *SyncService) detectConflictFields(task domain.Task, remote ado.WorkItem, syncState domain.SyncState) domain.SyncDiff {
	diff := domain.SyncDiff{
		TaskID:    task.ID,
		AdoID:     syncState.AdoID,
		Direction: "conflict",
	}

	// Title: both changed and values differ
	localTitleChanged := task.Title != syncState.LocalTitle
	remoteTitleChanged := remote.Title != syncState.RemoteTitle
	if localTitleChanged && remoteTitleChanged && task.Title != remote.Title {
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:  "title",
			Local:  task.Title,
			Remote: remote.Title,
		})
	}

	// State: both changed and values differ
	localStatusChanged := task.Status != syncState.LocalStatus
	remoteStatusChanged := remote.State != syncState.RemoteStatus
	if localStatusChanged && remoteStatusChanged {
		mappedLocal := ado.MapStatusToADO(task.Status, remote.Type)
		if mappedLocal != remote.State {
			diff.Changes = append(diff.Changes, domain.FieldDiff{
				Field:  "state",
				Local:  task.Status,
				Remote: remote.State,
			})
		}
	}

	// Description: both changed and values differ
	localDescChanged := task.Description != syncState.LocalDesc
	remoteDescChanged := remote.Description != syncState.RemoteDesc
	if localDescChanged && remoteDescChanged && task.Description != remote.Description {
		diff.Changes = append(diff.Changes, domain.FieldDiff{
			Field:  "description",
			Local:  task.Description,
			Remote: remote.Description,
		})
	}

	return diff
}

// getTask fetches a task from the database by ID.
func (s *SyncService) getTask(id int) (domain.Task, error) {
	var t domain.Task
	err := s.db.QueryRow(
		`SELECT id, title, description, status, priority, category, project_id, area, due_date,
		        ado_id, tags, blocked_reason, blocked_by, parent_id, personal_priority,
		        created_at, updated_at, completed_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
		&t.Category, &t.ProjectID, &t.Area, &t.DueDate,
		&t.AdoID, &t.Tags, &t.BlockedReason, &t.BlockedBy,
		&t.ParentID, &t.PersonalPriority,
		&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
	)
	return t, err
}

// getLinkedAdoID returns the first ADO ID linked to a task.
func (s *SyncService) getLinkedAdoID(taskID int) (string, error) {
	var adoID string
	err := s.db.QueryRow(
		`SELECT ado_id FROM task_ado_links WHERE task_id = ? LIMIT 1`, taskID,
	).Scan(&adoID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("task %d has no ADO link", taskID)
	}
	return adoID, err
}

// fetchRemoteItem fetches a single ADO work item by string ID, trying all configured orgs.
// Returns the work item and the client that found it (for subsequent push operations).
func (s *SyncService) fetchRemoteItem(adoID string) (*ado.WorkItem, *ado.Client, error) {
	id, err := strconv.Atoi(adoID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid ADO ID %q: %w", adoID, err)
	}

	token, err := s.tokenProv.GetToken()
	if err != nil {
		return nil, nil, fmt.Errorf("get token: %w", err)
	}

	orgProjects := config.GetOrgProjects()
	for _, op := range orgProjects {
		for _, proj := range op.Projects {
			c := ado.NewClient(op.Org, proj, token)
			wi, err := ado.GetWorkItem(c, id)
			if err != nil {
				continue
			}
			return wi, c, nil
		}
	}
	return nil, nil, fmt.Errorf("ADO item %s not found in any configured org/project", adoID)
}
