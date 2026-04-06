# Deferred Items — Phase 02

## Pre-existing Build Failures

- `internal/db/comments.go` references `domain.TaskComment` — type not yet defined in `domain/`
- `internal/db/links.go` references `domain.TaskLink` — type not yet defined
- `internal/db/ado.go` references `domain.SyncState`, `domain.ProjectADOLink` — types not yet defined
- These prevent transitive compilation of packages that import `internal/db` (e.g., `internal/auth`)
- **Not caused by Phase 02 changes** — pre-existing from prior phase work
