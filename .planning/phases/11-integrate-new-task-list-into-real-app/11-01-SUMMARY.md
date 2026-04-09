---
phase: 11-integrate-new-task-list-into-real-app
plan: 01
title: "Backend ADO Metadata Cache & Frontend Foundation"
subsystem: backend, frontend
tags: [ado-meta, cache, sqlite, composable, styles]
dependency_graph:
  requires: []
  provides: [ado_meta_cache_table, ADOMetaCacheService, useAdoMeta, statusIcon, priorityDotBgColor]
  affects: [syncservice, main, styles]
tech_stack:
  added: []
  patterns: [module-level-shared-state, batch-cache, dynamic-wails-imports]
key_files:
  created:
    - internal/app/adometa.go
    - frontend/src/api/adometa.ts
    - frontend/src/composables/useAdoMeta.ts
  modified:
    - internal/db/migrate.go
    - internal/app/syncservice.go
    - main.go
    - frontend/src/lib/styles.ts
decisions:
  - "ADOMetaCacheService uses DELETE + INSERT for cache rebuild (simple, no upsert complexity)"
  - "useAdoMeta uses module-level refs for shared singleton state across all consumers"
  - "API wrapper uses relative import path matching existing tasks.ts pattern"
  - "Dynamic import with try/catch for Wails bindings (not yet generated at dev time)"
metrics:
  duration: "6m 15s"
  completed: "2026-04-08T17:30:34Z"
  tasks_completed: 2
  tasks_total: 2
  files_created: 3
  files_modified: 4
---

# Phase 11 Plan 01: Backend ADO Metadata Cache & Frontend Foundation Summary

SQLite ado_meta_cache table + Go ADOMetaCacheService + Wails binding + useAdoMeta composable + statusIcon/priorityDotBgColor style helpers — batch ADO metadata for N+1-free task list rendering.

## What Was Done

### Task 1: Backend — ADO metadata cache table, service, and wiring
**Commit:** `11cc118`

- Added `ado_meta_cache` table to SQLite schema in `internal/db/migrate.go` with `task_id` (PK, FK to tasks), `ado_type`, `ado_state`, `synced_at`
- Added `idx_ado_meta_cache_task` index
- Created `internal/app/adometa.go` with `ADOMetaCacheService`:
  - `GetAll()` — returns `map[int]AdoMeta` from cache for batch frontend loading
  - `Refresh()` — rebuilds cache via `DELETE` + `INSERT ... SELECT` from `task_ado_links JOIN ado_work_items`
- Registered `ADOMetaCacheService` with Wails in `main.go`
- Added `adoMetaCache *ADOMetaCacheService` field to `SyncService` struct
- Updated `NewSyncService` to accept cache service parameter
- Added `s.adoMetaCache.Refresh()` call at end of `pullChanges()`

### Task 2: Frontend — style helpers, API wrapper, and useAdoMeta composable
**Commit:** `43c7464`

- Added `statusIcon(status)` to `styles.ts` — maps 6 task statuses to Lucide components (Circle, CircleDot, Eye, CheckCircle2, Octagon, XCircle)
- Added `priorityDotBgColor(priority)` to `styles.ts` — maps P0-P3 to Tailwind bg classes (red, orange, amber, zinc)
- Created `frontend/src/api/adometa.ts` with `getAllADOMeta()` and `refreshADOMeta()` wrappers using dynamic imports matching existing pattern
- Created `frontend/src/composables/useAdoMeta.ts`:
  - Module-level shared `metaCache` ref (singleton across consumers)
  - `getAdoMeta(taskId)` returns `AdoMeta | null` from batch-loaded cache
  - Auto-loads on first component mount
  - Auto-refreshes on `sync:completed` Wails event

## Deviations from Plan

None — plan executed exactly as written.

## Decisions Made

1. **Cache rebuild strategy:** DELETE + INSERT approach (not UPSERT) for simplicity — cache is small and fully rebuilt each sync
2. **Module-level shared state:** useAdoMeta uses module-level refs so all components share one cache instance without Pinia
3. **Import path pattern:** Used relative `../../bindings/...` path matching existing `tasks.ts` pattern (not `@/bindings/`)
4. **Dynamic imports:** Wails bindings are auto-generated at build time; dynamic imports with try/catch handle dev-time absence gracefully

## Verification

- `go build ./internal/... ./main.go` — exits 0
- `tsc --noEmit --skipLibCheck` — only pre-existing TS2307 errors for Vue SFCs and expected Wails binding path
- All acceptance criteria met per plan
