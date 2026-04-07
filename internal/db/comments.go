package db

import (
	"dev.azure.com/xbox/xb-tasks/domain"
)

func (db *DB) CreateComment(taskID int, content string) (domain.TaskComment, error) {
	res, err := db.Exec(`
		INSERT INTO task_comments (task_id, content)
		VALUES (?, ?)`, taskID, content)
	if err != nil {
		return domain.TaskComment{}, err
	}

	id, _ := res.LastInsertId()
	return db.getComment(int(id))
}

func (db *DB) ListComments(taskID int) ([]domain.TaskComment, error) {
	rows, err := db.Query(`
		SELECT id, task_id, content, is_public, ado_comment_id, created_at, updated_at
		FROM task_comments
		WHERE task_id = ?
		ORDER BY created_at ASC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.TaskComment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (db *DB) UpdateComment(id int, content string) error {
	_, err := db.Exec(`
		UPDATE task_comments
		SET content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, content, id)
	return err
}

func (db *DB) DeleteComment(id int) error {
	_, err := db.Exec(`DELETE FROM task_comments WHERE id = ?`, id)
	return err
}

func (db *DB) MarkCommentPublic(id int, adoCommentID string) error {
	_, err := db.Exec(`
		UPDATE task_comments
		SET is_public = 1, ado_comment_id = ?
		WHERE id = ?`, adoCommentID, id)
	return err
}

func (db *DB) getComment(id int) (domain.TaskComment, error) {
	var c domain.TaskComment
	var isPublic int
	err := db.QueryRow(`
		SELECT id, task_id, content, is_public, ado_comment_id, created_at, updated_at
		FROM task_comments
		WHERE id = ?`, id,
	).Scan(&c.ID, &c.TaskID, &c.Content, &isPublic, &c.AdoCommentID, &c.CreatedAt, &c.UpdatedAt)
	c.IsPublic = isPublic == 1
	return c, err
}

type commentScanner interface {
	Scan(dest ...any) error
}

func scanComment(s commentScanner) (domain.TaskComment, error) {
	var c domain.TaskComment
	var isPublic int
	err := s.Scan(&c.ID, &c.TaskID, &c.Content, &isPublic, &c.AdoCommentID, &c.CreatedAt, &c.UpdatedAt)
	c.IsPublic = isPublic == 1
	return c, err
}
