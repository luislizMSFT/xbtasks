---
phase: 02-ado-integration-prs-unified-dashboard
plan: 04
subsystem: backend-services
tags: [sync, ado, comments, links, goroutine, conflict-detection, url-detection]

# Dependency graph
requires:
  - phase: 02-01
    provides: "ADO REST client (pkg/ado), TokenProvider, auth chain"
  - phase: 02-02
    provides: "DB schema + CRUD for sync_state, comments, links, project_ado_links"
  - phase: 02-03
    provides: "LinkService personal-to-public model, ADOService REST integration"
provides:
  - "SyncService: background pull, outbound diff preview, push with confirmation, per-field conflict resolution"
  - "CommentService: local CRUD + selective ADO push"
  - "ExternalLinksService: URL pattern detection (ICM, Grafana, ADO, Wiki)"
  - "ProjectService: ADO linking, pin/star, dual progress tracking"
  - "domain.SyncDiff and domain.FieldDiff types for frontend consumption"
affects: [02-05, 02-06, 02-07, 02-08]

# Tech tracking
tech-stack:
  added: []
  patterns: [background-goroutine-with-ticker, event-driven-sync, per-field-conflict-resolution, url-pattern-detection]

key-files:
  created:
    - internal/app/syncservice.go
    - internal/app/commentservice.go
    - internal/app/externallinks.go
  modified:
    - internal/app/projects.go
    - main.go
    - domain/types.go

key-decisions:
  - "SyncService uses time.NewTicker goroutine with configurable interval from config"
  - "Outbound diff never auto-pushes — always requires user confirmation via GenerateOutboundDiff + PushChanges"
  - "Inbound auto-merge silently updates local task when only remote changed"
  - "Conflict detection compares both local and remote against last sync snapshot"
  - "Comments are always private/local by default; PushCommentToADO is explicit opt-in"
  - "DetectLinkType is a package-level function for reuse outside ExternalLinksService"
  - "ProjectService now accepts tokenProv+cfg for ADO operations (breaking constructor change)"
  - "Event API uses app.Event.Emit() per Wails v3 pattern"

patterns-established:
  - "Background goroutine pattern: ticker + stopCh for graceful shutdown"
  - "fetchRemoteItem returns both WorkItem and Client for subsequent operations"
  - "Domain types (SyncDiff, FieldDiff) separate from ado package types for frontend"
  - "URL pattern detection via compiled regexp table with fallback to 'url' type"

requirements-completed: [SYNC-01, SYNC-02, SYNC-03, SYNC-04, CMT-02, CMT-03, LINK-02, PROJ-04, PROJ-06, PROJ-07]

# Metrics
duration: 6min
completed: 2026-04-06
---

# Phase 02 Plan 04: Remaining Backend Services Summary

**SyncService with background pull + diff preview + conflict resolution, CommentService with selective ADO push, ExternalLinksService with URL pattern detection, ProjectService extended for ADO linking**

## Performance

- **Duration:** 6 min
- **Started:** 2026-04-06T19:06:22Z
- **Completed:** 2026-04-06T19:12:30Z
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments
- SyncService implements full bidirectional sync: auto-pull on timer, outbound diff preview, push with confirmation, per-field conflict resolution
- CommentService manages local comments with selective push to ADO (private by default)
- ExternalLinksService auto-detects ICM, Grafana, ADO, and Wiki URLs from regex patterns
- ProjectService extended with ADO linking, pin/star functionality, and dual progress tracking (local + ADO)
- All new services wired into main.go with Wails v3 service registration
- Background sync starts automatically after auth restore

## Task Commits

Each task was committed atomically:

1. **Task 1: SyncService — Background Pull + Diff Generation + Conflict Detection** - `55ede99` (feat)
2. **Task 2: CommentService + ExternalLinksService + Project Extension + Main Wiring** - `b4b0594` (feat)

## Files Created/Modified
- `internal/app/syncservice.go` - Background sync, outbound diff, push, conflict resolution (380+ lines)
- `internal/app/commentservice.go` - Local comment CRUD + selective ADO push
- `internal/app/externallinks.go` - URL pattern detection + link CRUD
- `internal/app/projects.go` - Extended with ADO linking, pin, progress tracking
- `main.go` - Service registration, constructor update, background sync start
- `domain/types.go` - Added SyncDiff and FieldDiff domain types

## Decisions Made
- Used `app.Event.Emit()` for Wails v3 event emission (not `EmitEvent`)
- Created domain-level `SyncDiff`/`FieldDiff` types separate from `ado.SyncDiff`/`ado.FieldDiff` to decouple frontend from ADO package
- Moved `projectService` creation after `tokenProvider` initialization in main.go (was using undefined variable)
- `fetchRemoteItem` returns both `*ado.WorkItem` and `*ado.Client` — enables subsequent push operations on the same client
- Conflict detection uses snapshot comparison (not timestamp-based) for reliability

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed Wails v3 event emission API**
- **Found during:** Task 1 (SyncService)
- **Issue:** Plan used `s.app.EmitEvent()` but Wails v3 uses `s.app.Event.Emit()`
- **Fix:** Changed to correct API pattern matching existing auth service usage
- **Files modified:** internal/app/syncservice.go
- **Verification:** `go build ./internal/app/...` passes
- **Committed in:** 55ede99 (Task 1 commit)

**2. [Rule 1 - Bug] Fixed variable declaration order in main.go**
- **Found during:** Task 2 (main.go wiring)
- **Issue:** `NewProjectService` called with `tokenProvider` before it was declared
- **Fix:** Moved `projectService` creation after token provider initialization
- **Files modified:** main.go
- **Verification:** `go build .` passes
- **Committed in:** b4b0594 (Task 2 commit)

---

**Total deviations:** 2 auto-fixed (2 bugs)
**Impact on plan:** Both auto-fixes necessary for compilation. No scope creep.

## Issues Encountered
- `build/ios` package has pre-existing compilation error (empty main) — not related to our changes, excluded from verification

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- All backend services complete — ready for frontend stores and UI components (Plans 05-08)
- SyncService, CommentService, ExternalLinksService all registered in Wails for frontend binding
- Background sync infrastructure in place for real-time updates

---
*Phase: 02-ado-integration-prs-unified-dashboard*
*Completed: 2026-04-06*
