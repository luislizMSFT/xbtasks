package db

// migrate runs the schema creation followed by incremental column additions.
// The schema uses CREATE TABLE IF NOT EXISTS so it's safe on fresh and existing DBs.
// ALTER TABLE additions are idempotent — SQLite returns an error if the column
// already exists, which we intentionally ignore.
func (db *DB) migrate() error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	for _, alt := range alterations {
		db.Exec(alt) // ignore "duplicate column" errors
	}

	return nil
}

const schema = `
-- ─── Core tables ───

CREATE TABLE IF NOT EXISTS projects (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	name        TEXT NOT NULL,
	description TEXT DEFAULT '',
	status      TEXT DEFAULT 'active' CHECK(status IN ('active','paused','completed','archived')),
	is_pinned   INTEGER DEFAULT 0,
	created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
	id                INTEGER PRIMARY KEY AUTOINCREMENT,
	title             TEXT NOT NULL,
	description       TEXT DEFAULT '',
	status            TEXT DEFAULT 'todo' CHECK(status IN ('todo','in_progress','in_review','done','blocked','cancelled')),
	priority          TEXT DEFAULT 'P2'  CHECK(priority IN ('P0','P1','P2','P3')),
	category          TEXT DEFAULT '',
	project_id        INTEGER REFERENCES projects(id),
	area              TEXT DEFAULT '',
	due_date          TEXT DEFAULT '',
	ado_id            TEXT DEFAULT '',
	tags              TEXT DEFAULT '',
	blocked_reason    TEXT DEFAULT '',
	blocked_by        TEXT DEFAULT '',
	parent_id         INTEGER REFERENCES tasks(id),
	personal_priority TEXT DEFAULT '',
	sort_order        INTEGER DEFAULT 0,
	created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
	completed_at      DATETIME
);

CREATE TABLE IF NOT EXISTS task_deps (
	task_id    INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	depends_on INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	PRIMARY KEY (task_id, depends_on)
);

-- ─── Pull requests ───

CREATE TABLE IF NOT EXISTS pull_requests (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	title         TEXT NOT NULL,
	pr_url        TEXT DEFAULT '',
	pr_number     INTEGER DEFAULT 0,
	repo          TEXT DEFAULT '',
	task_id       INTEGER REFERENCES tasks(id),
	ado_id        TEXT DEFAULT '',
	status        TEXT DEFAULT 'active' CHECK(status IN ('draft','active','completed','abandoned')),
	reviewers     TEXT DEFAULT '',
	source_branch TEXT DEFAULT '',
	target_branch TEXT DEFAULT '',
	votes         INTEGER DEFAULT 0,
	created_by    TEXT DEFAULT '',
	created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
	merged_at     DATETIME
);

-- ─── ADO integration ───

CREATE TABLE IF NOT EXISTS ado_work_items (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	ado_id      TEXT NOT NULL UNIQUE,
	title       TEXT NOT NULL,
	state       TEXT DEFAULT '',
	type        TEXT DEFAULT '',
	assigned_to TEXT DEFAULT '',
	priority    INTEGER DEFAULT 2,
	area_path   TEXT DEFAULT '',
	description TEXT DEFAULT '',
	url         TEXT DEFAULT '',
	org         TEXT DEFAULT '',
	project     TEXT DEFAULT '',
	parent_id   INTEGER DEFAULT 0,
	changed_date DATETIME DEFAULT CURRENT_TIMESTAMP,
	synced_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task_ado_links (
	task_id    INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	ado_id     TEXT NOT NULL,
	direction  TEXT DEFAULT 'linked' CHECK(direction IN ('promoted','imported','linked')),
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (task_id, ado_id)
);

CREATE TABLE IF NOT EXISTS project_ado_links (
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	ado_id     TEXT NOT NULL,
	direction  TEXT DEFAULT 'linked' CHECK(direction IN ('promoted','imported','linked')),
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (project_id, ado_id)
);

CREATE TABLE IF NOT EXISTS sync_state (
	task_id        INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	ado_id         TEXT NOT NULL,
	last_synced_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	local_title    TEXT DEFAULT '',
	local_status   TEXT DEFAULT '',
	local_desc     TEXT DEFAULT '',
	remote_title   TEXT DEFAULT '',
	remote_status  TEXT DEFAULT '',
	remote_desc    TEXT DEFAULT '',
	PRIMARY KEY (task_id, ado_id)
);

-- ─── Supplementary tables ───

CREATE TABLE IF NOT EXISTS users (
	id           TEXT PRIMARY KEY,
	display_name TEXT NOT NULL,
	email        TEXT DEFAULT '',
	avatar_url   TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS task_links (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	task_id    INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	url        TEXT NOT NULL,
	label      TEXT DEFAULT '',
	type       TEXT DEFAULT 'url',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task_comments (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	task_id        INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	content        TEXT NOT NULL,
	is_public      INTEGER DEFAULT 0,
	ado_comment_id TEXT DEFAULT '',
	created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ─── ADO metadata cache ───

CREATE TABLE IF NOT EXISTS ado_meta_cache (
	task_id   INTEGER PRIMARY KEY REFERENCES tasks(id) ON DELETE CASCADE,
	ado_type  TEXT DEFAULT '',
	ado_state TEXT DEFAULT '',
	synced_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ─── Indexes (after all tables) ───

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_ado_id ON tasks(ado_id);
CREATE INDEX IF NOT EXISTS idx_tasks_parent ON tasks(parent_id);
CREATE INDEX IF NOT EXISTS idx_tasks_project_status ON tasks(project_id, status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_pull_requests_pr_repo ON pull_requests(pr_number, repo);
CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pull_requests_task ON pull_requests(task_id);
CREATE INDEX IF NOT EXISTS idx_ado_work_items_ado_id ON ado_work_items(ado_id);
CREATE INDEX IF NOT EXISTS idx_ado_work_items_parent_id ON ado_work_items(parent_id);
CREATE INDEX IF NOT EXISTS idx_task_ado_links_ado_id ON task_ado_links(ado_id);
CREATE INDEX IF NOT EXISTS idx_sync_state_task_id ON sync_state(task_id);
CREATE INDEX IF NOT EXISTS idx_task_links_task ON task_links(task_id);
CREATE INDEX IF NOT EXISTS idx_task_comments_task ON task_comments(task_id);
CREATE INDEX IF NOT EXISTS idx_ado_meta_cache_task ON ado_meta_cache(task_id);
`

// alterations are idempotent column additions for databases created before
// those columns existed. On a fresh DB these are no-ops (column already in
// the CREATE TABLE). Keep this list append-only.
var alterations = []string{
	// v0.1 → v0.2
	`ALTER TABLE pull_requests ADD COLUMN created_by TEXT DEFAULT ''`,
	`ALTER TABLE projects ADD COLUMN is_pinned INTEGER DEFAULT 0`,
	`ALTER TABLE ado_work_items ADD COLUMN org TEXT DEFAULT ''`,
	`ALTER TABLE ado_work_items ADD COLUMN project TEXT DEFAULT ''`,
	`ALTER TABLE ado_work_items ADD COLUMN parent_id INTEGER DEFAULT 0`,
	`ALTER TABLE ado_work_items ADD COLUMN changed_date DATETIME DEFAULT CURRENT_TIMESTAMP`,
	// v0.2 → v0.3
	`ALTER TABLE tasks ADD COLUMN sort_order INTEGER DEFAULT 0`,
}
