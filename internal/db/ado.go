package db

import (
	"dev.azure.com/xbox/xb-tasks/domain"
)

func (db *DB) UpsertADOWorkItem(item domain.ADOWorkItem) error {
	_, err := db.Exec(`
		INSERT INTO ado_work_items (ado_id, title, state, type, assigned_to, priority, area_path, description, url, synced_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(ado_id) DO UPDATE SET
			title       = excluded.title,
			state       = excluded.state,
			type        = excluded.type,
			assigned_to = excluded.assigned_to,
			priority    = excluded.priority,
			area_path   = excluded.area_path,
			description = excluded.description,
			url         = excluded.url,
			synced_at   = CURRENT_TIMESTAMP`,
		item.AdoID, item.Title, item.State, item.Type, item.AssignedTo,
		item.Priority, item.AreaPath, item.Description, item.URL,
	)
	return err
}

func (db *DB) ListADOWorkItems() ([]domain.ADOWorkItem, error) {
	rows, err := db.Query(`
		SELECT id, ado_id, title, state, type, assigned_to, priority, area_path, description, url, synced_at
		FROM ado_work_items
		ORDER BY priority ASC, title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ADOWorkItem
	for rows.Next() {
		var item domain.ADOWorkItem
		if err := rows.Scan(
			&item.ID, &item.AdoID, &item.Title, &item.State, &item.Type,
			&item.AssignedTo, &item.Priority, &item.AreaPath, &item.Description,
			&item.URL, &item.SyncedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (db *DB) GetADOWorkItem(adoID string) (domain.ADOWorkItem, error) {
	var item domain.ADOWorkItem
	err := db.QueryRow(`
		SELECT id, ado_id, title, state, type, assigned_to, priority, area_path, description, url, synced_at
		FROM ado_work_items
		WHERE ado_id = ?`, adoID,
	).Scan(
		&item.ID, &item.AdoID, &item.Title, &item.State, &item.Type,
		&item.AssignedTo, &item.Priority, &item.AreaPath, &item.Description,
		&item.URL, &item.SyncedAt,
	)
	return item, err
}
