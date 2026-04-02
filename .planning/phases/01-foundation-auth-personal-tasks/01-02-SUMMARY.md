---
phase: 01
plan: 02
subsystem: backend-task-services
tags: [go, sqlite, dependencies, crud, tags, cycle-detection]
dependency_graph:
  requires: [internal/db/db.go, domain/types.go]
  provides: [internal/app/deps.go, extended internal/app/tasks.go]
  affects: [cmd/main.go (needs DependencyService registration)]
tech_stack:
  added: []
  patterns: [DFS cycle detection, CSV tag search with SQL LIKE, dynamic SQL filter building]
key_files:
  created: [internal/app/deps.go]
  modified: [internal/app/tasks.go]
key_decisions:
  - "Used DFS (iterative stack) for cycle detection — simple, correct, no recursion depth issues"
  - "Tags searched via SQL (',' || tags || ',') LIKE pattern for comma-separated field"
  - "scanTasks helper extracted to deps.go for shared row scanning pattern"
  - "Used domain package (domain.Task) instead of plan's pkg/models reference"
metrics:
  duration: "2m 37s"
  tasks_completed: 2
  tasks_total: 2
  files_created: 1
  files_modified: 1
  completed: "2026-04-02T19:48:23Z"
---

# Phase 01 Plan 02: Go Backend — Task Dependencies, Tags & Extended CRUD Summary

DependencyService with DFS circular detection and TaskService extended with personal priority, subtasks, filtered list, and tag aggregation — all using domain.Task types against SQLite task_deps table.

## What Was Done

### Task 02.1: Create DependencyService with circular detection
- **Created** `internal/app/deps.go` with 5 exported methods + 1 private helper
- `AddDependency`: validates self-dep, both-exist, DFS cycle check, INSERT OR IGNORE
- `hasCircularDep`: iterative DFS from dependsOn following task_deps edges
- `RemoveDependency`: simple DELETE from task_deps
- `GetDependencies`: returns full Task objects that taskID depends on (JOIN query)
- `GetBlockedBy`: returns full Task objects blocked by taskID (reverse direction JOIN)
- `scanTasks` helper for shared row scanning across GetDependencies/GetBlockedBy
- **Commit:** `856cee2`

### Task 02.2: Extend TaskService with tags, personal priority, filtered list, subtask creation
- **Added** 4 new methods to `internal/app/tasks.go` (11 total methods now)
- `SetPersonalPriority`: UPDATE personal_priority field independent of ADO priority (TASK-05)
- `CreateSubtask`: validates parent exists, delegates to Create with parentID (TASK-04)
- `ListFiltered`: dynamic SQL with status, projectID, parentID, tag filters (TASK-01/03/07)
- `GetAllTags`: aggregates unique tags from comma-separated field, sorted (TASK-07)
- Added `sort` and `strings` imports
- **Commit:** `d5e774b`

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Corrected domain import path from pkg/models to domain**
- **Found during:** Task 02.1 & 02.2
- **Issue:** Plan referenced `pkg/models.Task` and import `dev.azure.com/xbox/xb-tasks/pkg/models` but actual package is `domain` at `dev.azure.com/xbox/xb-tasks/domain`
- **Fix:** Used `domain.Task` and correct import path throughout both files
- **Files modified:** internal/app/deps.go, internal/app/tasks.go

**2. [Rule 2 - Missing functionality] Added scanTasks helper to reduce duplication**
- **Found during:** Task 02.1
- **Issue:** GetDependencies and GetBlockedBy would duplicate the 18-field Scan logic
- **Fix:** Extracted `scanTasks()` helper function in deps.go
- **Files modified:** internal/app/deps.go

## Verification Results

- `go build ./internal/app/ ./internal/db/ ./domain/` — exits 0 ✓
- DependencyService has 5 exported methods ✓
- TaskService has 11 total methods (7 original + 4 new) ✓
- `hasCircularDep` defined and called ✓
- `INSERT OR IGNORE INTO task_deps` present ✓
- All 4 new TaskService methods match acceptance criteria ✓

## Known Stubs

None — all methods are fully implemented with real SQL queries against the existing schema.

## Commits

| Task | Commit | Message |
|------|--------|---------|
| 02.1 | 856cee2 | feat(01-02): create DependencyService with circular detection |
| 02.2 | d5e774b | feat(01-02): extend TaskService with tags, priority, filtered list, subtasks |

## Self-Check: PASSED

- [x] internal/app/deps.go exists
- [x] internal/app/tasks.go exists
- [x] 01-02-SUMMARY.md exists
- [x] Commit 856cee2 found
- [x] Commit d5e774b found
