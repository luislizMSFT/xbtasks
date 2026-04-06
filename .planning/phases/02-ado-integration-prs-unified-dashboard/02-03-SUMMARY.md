---
phase: 02-ado-integration-prs-unified-dashboard
plan: 03
subsystem: api
tags: [ado, rest-client, token-provider, link-service, multi-org]

# Dependency graph
requires:
  - phase: 02-01
    provides: "TokenProvider interface, CachedTokenProvider, AzCliTokenProvider, pkg/ado REST client"
  - phase: 02-02
    provides: "Multi-org config, domain types (TaskADOLink, SyncState, OrgProject), DB migrations"
provides:
  - "ADOService using REST client with multi-org iteration"
  - "LinkService with link/promote/import/unlink operations"
  - "Token provider chain wired in main.go"
  - "GetWorkItemTree for ADO browser hierarchy"
  - "IsPublic computed from task_ado_links presence"
affects: [02-04, 02-05, 02-06, 02-07, 02-08]

# Tech tracking
tech-stack:
  added: []
  patterns: [token-provider-chain, multi-org-client-iteration, personal-to-public-model]

key-files:
  created:
    - internal/app/linkservice.go
  modified:
    - internal/app/adoservice.go
    - main.go

key-decisions:
  - "ADOService getClients() flattens OrgProject pairs into individual ado.Client instances for iteration"
  - "GetWorkItemTree fetches parent hierarchy up to 3 levels deep for ADO browser tree rendering"
  - "LinkService reuses adoWorkItemToDomain converter from ADOService (same package)"
  - "priorityToLocal maps ADO 1-4 integer priority to P0-P3 string for imported tasks"
  - "Token provider chain: AzCli -> CachedWrapper with 5 minute buffer, created once in main.go"

patterns-established:
  - "Token provider chain: create once in main.go, pass to services that need ADO auth"
  - "Multi-org iteration: getClients() returns all org/project clients, methods iterate with continue-on-error"
  - "Link directions: 'linked' (manual connect), 'promoted' (local->ADO create), 'imported' (ADO->local pull)"

requirements-completed: [ADO-01, ADO-02, ADO-03, ADO-04, ADO-05, ADO-06, ADO-07, ADO-10, TASK-08]

# Metrics
duration: 5min
completed: 2026-04-06
---

# Phase 02 Plan 03: ADOService REST Refactor + LinkService Summary

**ADOService rewritten from az-cli subprocess to pkg/ado REST client with multi-org support; LinkService implements link/promote/import/unlink for personal→public task model**

## Performance

- **Duration:** 5 min
- **Started:** 2026-04-06T18:59:03Z
- **Completed:** 2026-04-06T19:03:52Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments
- ADOService completely refactored: removed all exec.Command/az-cli calls, now uses pkg/ado REST client with Bearer tokens
- Multi-org support: all ADO methods iterate configured org/project pairs via getClients()
- GetWorkItemTree added for ADO browser — fetches assigned items + parent hierarchy up to 3 levels
- LinkService created with all 4 linking flows: LinkTask, PromoteTask, ImportWorkItem, UnlinkTask
- IsPublic/GetTaskLinks/ListPublicTaskIDs for frontend personal↔public task model
- Token provider chain (AzCli → Cached) wired once in main.go and shared across ADOService + LinkService

## Task Commits

Each task was committed atomically:

1. **Task 1: Refactor ADOService to Use REST Client + Multi-Org** - `c73d5a6` (feat)
2. **Task 2: LinkService — Link, Promote, Import, Unlink + Main Wiring** - `e691eb9` (feat)

## Files Created/Modified
- `internal/app/adoservice.go` - Refactored from az-cli to REST client with TokenProvider, multi-org, and tree fetching
- `internal/app/linkservice.go` - New: link/promote/import/unlink + IsPublic + GetTaskLinks + ListPublicTaskIDs
- `main.go` - Token provider chain creation, ADOService constructor updated, LinkService registered

## Decisions Made
- ADOService.getClients() flattens OrgProject (org → []projects) into individual ado.Client instances since each client targets one org/project pair
- GetWorkItemTree fetches 3 levels of parent hierarchy by batch-fetching parent IDs across all clients, deduplicating by AdoID
- LinkService has its own getTask() helper (same SQL as TaskService.GetByID) to avoid circular service dependency
- priorityToLocal maps ADO integer priority 1→P0, 2→P1, 3→P2, 4→P3 for imported tasks
- Token provider chain created once in main.go (not per-service) — both ADOService and LinkService share the same cached provider

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
- `go build ./...` shows pre-existing error in `build/ios` (Wails scaffold stub with no main function) — not caused by our changes, core packages all build cleanly

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- ADOService and LinkService are ready for frontend integration (Plan 04+)
- Sync workflow (Plan 05) can use SyncWorkItems and sync_state snapshots created by LinkService
- Conflict resolution (Plan 06) has sync_state baseline from link/promote/import operations
- ADO browser frontend (Plan 07) can use GetWorkItemTree for hierarchy display

---
*Phase: 02-ado-integration-prs-unified-dashboard*
*Completed: 2026-04-06*
