package app

import (
	"context"
	"fmt"

	"dev.azure.com/xbox/xb-tasks/internal/db"
)

// AdoMeta holds the cached ADO work item type and state for a linked task.
type AdoMeta struct {
	Type  string `json:"type"`
	State string `json:"state"`
}

// ADOMetaCacheService provides batch access to cached ADO metadata.
// The cache is populated by SyncService after each sync cycle.
type ADOMetaCacheService struct {
	db *db.DB
}

// NewADOMetaCacheService creates a new cache service.
func NewADOMetaCacheService(database *db.DB) *ADOMetaCacheService {
	return &ADOMetaCacheService{db: database}
}

// GetAll returns a map of task_id → AdoMeta for all linked tasks.
// Called by the frontend to batch-load metadata for task list rendering.
func (s *ADOMetaCacheService) GetAll() (map[int]AdoMeta, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT task_id, ado_type, ado_state
		FROM ado_meta_cache
	`)
	if err != nil {
		return nil, fmt.Errorf("query ado_meta_cache: %w", err)
	}
	defer rows.Close()

	result := make(map[int]AdoMeta)
	for rows.Next() {
		var taskID int
		var meta AdoMeta
		if err := rows.Scan(&taskID, &meta.Type, &meta.State); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result[taskID] = meta
	}
	return result, rows.Err()
}

// Refresh rebuilds the cache by joining task_ado_links with ado_work_items.
// Called by SyncService after pullChanges completes.
func (s *ADOMetaCacheService) Refresh() error {
	ctx := context.Background()
	if _, err := s.db.ExecContext(ctx, `DELETE FROM ado_meta_cache`); err != nil {
		return fmt.Errorf("clear ado_meta_cache: %w", err)
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO ado_meta_cache (task_id, ado_type, ado_state, synced_at)
		SELECT
			l.task_id,
			w.type,
			w.state,
			CURRENT_TIMESTAMP
		FROM task_ado_links l
		INNER JOIN ado_work_items w ON l.ado_id = w.ado_id
	`)
	if err != nil {
		return fmt.Errorf("rebuild ado_meta_cache: %w", err)
	}
	return nil
}
