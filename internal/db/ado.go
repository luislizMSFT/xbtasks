package db

import (
	"dev.azure.com/xbox/xb-tasks/domain"
)

func (db *DB) UpsertADOWorkItem(item domain.ADOWorkItem) error {
	_, err := db.Exec(`
		INSERT INTO ado_work_items (ado_id, title, state, type, assigned_to, priority, area_path, description, url, org, project, parent_id, changed_date, synced_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(ado_id) DO UPDATE SET
			title        = excluded.title,
			state        = excluded.state,
			type         = excluded.type,
			assigned_to  = excluded.assigned_to,
			priority     = excluded.priority,
			area_path    = excluded.area_path,
			description  = excluded.description,
			url          = excluded.url,
			org          = excluded.org,
			project      = excluded.project,
			parent_id    = excluded.parent_id,
			changed_date = excluded.changed_date,
			synced_at    = CURRENT_TIMESTAMP`,
		item.AdoID, item.Title, item.State, item.Type, item.AssignedTo,
		item.Priority, item.AreaPath, item.Description, item.URL,
		item.Org, item.Project, item.ParentID, item.ChangedDate,
	)
	return err
}

func (db *DB) ListADOWorkItems() ([]domain.ADOWorkItem, error) {
	rows, err := db.Query(`
		SELECT id, ado_id, title, state, type, assigned_to, priority, area_path, description, url, org, project, parent_id, changed_date, synced_at
		FROM ado_work_items
		ORDER BY priority ASC, title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanADOWorkItems(rows)
}

func (db *DB) ListADOWorkItemsByOrg(org, project string) ([]domain.ADOWorkItem, error) {
	rows, err := db.Query(`
		SELECT id, ado_id, title, state, type, assigned_to, priority, area_path, description, url, org, project, parent_id, changed_date, synced_at
		FROM ado_work_items
		WHERE org = ? AND project = ?
		ORDER BY priority ASC, title ASC`, org, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanADOWorkItems(rows)
}

func (db *DB) GetADOWorkItem(adoID string) (domain.ADOWorkItem, error) {
	var item domain.ADOWorkItem
	err := db.QueryRow(`
		SELECT id, ado_id, title, state, type, assigned_to, priority, area_path, description, url, org, project, parent_id, changed_date, synced_at
		FROM ado_work_items
		WHERE ado_id = ?`, adoID,
	).Scan(
		&item.ID, &item.AdoID, &item.Title, &item.State, &item.Type,
		&item.AssignedTo, &item.Priority, &item.AreaPath, &item.Description,
		&item.URL, &item.Org, &item.Project, &item.ParentID, &item.ChangedDate,
		&item.SyncedAt,
	)
	return item, err
}

// --- Sync State ---

func (db *DB) UpsertSyncState(taskID int, adoID string, localTitle, localStatus, localDesc, remoteTitle, remoteStatus, remoteDesc string) error {
	_, err := db.Exec(`
		INSERT INTO sync_state (task_id, ado_id, last_synced_at, local_title, local_status, local_desc, remote_title, remote_status, remote_desc)
		VALUES (?, ?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(task_id, ado_id) DO UPDATE SET
			last_synced_at = CURRENT_TIMESTAMP,
			local_title    = excluded.local_title,
			local_status   = excluded.local_status,
			local_desc     = excluded.local_desc,
			remote_title   = excluded.remote_title,
			remote_status  = excluded.remote_status,
			remote_desc    = excluded.remote_desc`,
		taskID, adoID, localTitle, localStatus, localDesc, remoteTitle, remoteStatus, remoteDesc,
	)
	return err
}

func (db *DB) GetSyncState(taskID int, adoID string) (domain.SyncState, error) {
	var s domain.SyncState
	err := db.QueryRow(`
		SELECT task_id, ado_id, last_synced_at, local_title, local_status, local_desc, remote_title, remote_status, remote_desc
		FROM sync_state
		WHERE task_id = ? AND ado_id = ?`, taskID, adoID,
	).Scan(
		&s.TaskID, &s.AdoID, &s.LastSyncedAt,
		&s.LocalTitle, &s.LocalStatus, &s.LocalDesc,
		&s.RemoteTitle, &s.RemoteStatus, &s.RemoteDesc,
	)
	return s, err
}

func (db *DB) ListSyncStates() ([]domain.SyncState, error) {
	rows, err := db.Query(`
		SELECT task_id, ado_id, last_synced_at, local_title, local_status, local_desc, remote_title, remote_status, remote_desc
		FROM sync_state
		ORDER BY last_synced_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []domain.SyncState
	for rows.Next() {
		var s domain.SyncState
		if err := rows.Scan(
			&s.TaskID, &s.AdoID, &s.LastSyncedAt,
			&s.LocalTitle, &s.LocalStatus, &s.LocalDesc,
			&s.RemoteTitle, &s.RemoteStatus, &s.RemoteDesc,
		); err != nil {
			return nil, err
		}
		states = append(states, s)
	}
	return states, rows.Err()
}

// --- Project ADO Links ---

func (db *DB) CreateProjectADOLink(projectID int, adoID, direction string) error {
	_, err := db.Exec(`
		INSERT INTO project_ado_links (project_id, ado_id, direction)
		VALUES (?, ?, ?)`, projectID, adoID, direction)
	return err
}

func (db *DB) DeleteProjectADOLink(projectID int, adoID string) error {
	_, err := db.Exec(`
		DELETE FROM project_ado_links
		WHERE project_id = ? AND ado_id = ?`, projectID, adoID)
	return err
}

func (db *DB) GetProjectADOLink(projectID int) (domain.ProjectADOLink, error) {
	var l domain.ProjectADOLink
	err := db.QueryRow(`
		SELECT project_id, ado_id, direction, created_at
		FROM project_ado_links
		WHERE project_id = ?`, projectID,
	).Scan(&l.ProjectID, &l.AdoID, &l.Direction, &l.CreatedAt)
	return l, err
}

// scanADOWorkItems scans rows into a slice of ADOWorkItem.
func scanADOWorkItems(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]domain.ADOWorkItem, error) {
	var items []domain.ADOWorkItem
	for rows.Next() {
		var item domain.ADOWorkItem
		if err := rows.Scan(
			&item.ID, &item.AdoID, &item.Title, &item.State, &item.Type,
			&item.AssignedTo, &item.Priority, &item.AreaPath, &item.Description,
			&item.URL, &item.Org, &item.Project, &item.ParentID, &item.ChangedDate,
			&item.SyncedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
