# Deferred Items — Phase 02

## Pre-existing Build Failures

- Previously noted missing domain types for `internal/db/comments.go`, `internal/db/links.go`, and `internal/db/ado.go` have been added in `domain/types.go` in this PR.
- Do not treat `domain.TaskComment`, `domain.TaskLink`, `domain.SyncState`, or `domain.ProjectADOLink` as current build blockers in Phase 02 planning.
- If compilation issues remain in packages importing `internal/db` (for example `internal/auth`), they should be re-evaluated based on the current code rather than this older missing-type note.
