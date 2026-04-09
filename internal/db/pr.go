package db

import (
	"database/sql"

	"dev.azure.com/xbox/xb-tasks/domain"
)

func (db *DB) UpsertPullRequest(pr domain.PullRequest) error {
	_, err := db.Exec(`
		INSERT INTO pull_requests (title, pr_url, pr_number, repo, task_id, ado_id, status, reviewers, source_branch, target_branch, votes, created_by, created_at, updated_at, merged_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(pr_number, repo) DO UPDATE SET
			title         = excluded.title,
			pr_url        = excluded.pr_url,
			status        = excluded.status,
			reviewers     = excluded.reviewers,
			source_branch = excluded.source_branch,
			target_branch = excluded.target_branch,
			votes         = excluded.votes,
			created_by    = excluded.created_by,
			updated_at    = CURRENT_TIMESTAMP,
			merged_at     = excluded.merged_at`,
		pr.Title, pr.PRURL, pr.PRNumber, pr.Repo, pr.TaskID, pr.AdoID,
		pr.Status, pr.Reviewers, pr.SourceBranch, pr.TargetBranch, pr.Votes,
		pr.CreatedBy, pr.CreatedAt, pr.UpdatedAt, pr.MergedAt,
	)
	return err
}

func (db *DB) ListPullRequests() ([]domain.PullRequest, error) {
	return db.queryPRs(`
		SELECT id, title, pr_url, pr_number, repo, task_id, ado_id, status, reviewers,
		       source_branch, target_branch, votes, created_by, created_at, updated_at, merged_at,
		       dismissed_at, watched
		FROM pull_requests
		ORDER BY created_at DESC`)
}

func (db *DB) ListPullRequestsByStatus(status string) ([]domain.PullRequest, error) {
	return db.queryPRs(`
		SELECT id, title, pr_url, pr_number, repo, task_id, ado_id, status, reviewers,
		       source_branch, target_branch, votes, created_by, created_at, updated_at, merged_at,
		       dismissed_at, watched
		FROM pull_requests
		WHERE status = ?
		ORDER BY created_at DESC`, status)
}

func (db *DB) DismissPullRequest(prNumber int, repo string) error {
	_, err := db.Exec(`UPDATE pull_requests SET dismissed_at = CURRENT_TIMESTAMP WHERE pr_number = ? AND repo = ?`, prNumber, repo)
	return err
}

func (db *DB) UndismissPullRequest(prNumber int, repo string) error {
	_, err := db.Exec(`UPDATE pull_requests SET dismissed_at = NULL WHERE pr_number = ? AND repo = ?`, prNumber, repo)
	return err
}

func (db *DB) SetPullRequestWatched(prNumber int, repo string, watched bool) error {
	val := 0
	if watched {
		val = 1
	}
	_, err := db.Exec(`UPDATE pull_requests SET watched = ? WHERE pr_number = ? AND repo = ?`, val, prNumber, repo)
	return err
}

func (db *DB) queryPRs(query string, args ...any) ([]domain.PullRequest, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		var mergedAt, dismissedAt sql.NullTime
		var taskID sql.NullInt64
		var watched sql.NullInt64
		if err := rows.Scan(
			&pr.ID, &pr.Title, &pr.PRURL, &pr.PRNumber, &pr.Repo,
			&taskID, &pr.AdoID, &pr.Status, &pr.Reviewers,
			&pr.SourceBranch, &pr.TargetBranch, &pr.Votes,
			&pr.CreatedBy, &pr.CreatedAt, &pr.UpdatedAt, &mergedAt,
			&dismissedAt, &watched,
		); err != nil {
			return nil, err
		}
		if taskID.Valid {
			id := int(taskID.Int64)
			pr.TaskID = &id
		}
		if mergedAt.Valid {
			pr.MergedAt = &mergedAt.Time
		}
		if dismissedAt.Valid {
			pr.DismissedAt = &dismissedAt.Time
		}
		pr.Watched = watched.Valid && watched.Int64 == 1
		prs = append(prs, pr)
	}
	return prs, rows.Err()
}
