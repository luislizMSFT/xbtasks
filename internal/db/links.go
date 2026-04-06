package db

import (
	"dev.azure.com/xbox/xb-tasks/domain"
)

func (db *DB) CreateLink(taskID int, url, label, linkType string) (domain.TaskLink, error) {
	res, err := db.Exec(`
		INSERT INTO task_links (task_id, url, label, type)
		VALUES (?, ?, ?, ?)`, taskID, url, label, linkType)
	if err != nil {
		return domain.TaskLink{}, err
	}

	id, _ := res.LastInsertId()
	return db.getLink(int(id))
}

func (db *DB) ListLinks(taskID int) ([]domain.TaskLink, error) {
	rows, err := db.Query(`
		SELECT id, task_id, url, label, type, created_at
		FROM task_links
		WHERE task_id = ?
		ORDER BY created_at DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []domain.TaskLink
	for rows.Next() {
		var l domain.TaskLink
		if err := rows.Scan(&l.ID, &l.TaskID, &l.URL, &l.Label, &l.Type, &l.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

func (db *DB) DeleteLink(id int) error {
	_, err := db.Exec(`DELETE FROM task_links WHERE id = ?`, id)
	return err
}

func (db *DB) getLink(id int) (domain.TaskLink, error) {
	var l domain.TaskLink
	err := db.QueryRow(`
		SELECT id, task_id, url, label, type, created_at
		FROM task_links
		WHERE id = ?`, id,
	).Scan(&l.ID, &l.TaskID, &l.URL, &l.Label, &l.Type, &l.CreatedAt)
	return l, err
}
