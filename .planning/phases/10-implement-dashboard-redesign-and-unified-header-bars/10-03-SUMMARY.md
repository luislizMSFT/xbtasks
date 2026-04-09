---
phase: 10-implement-dashboard-redesign-and-unified-header-bars
plan: 03
title: "UAT Gap Closure: SyncCluster Loading State + Attention Bar Date Fix"
subsystem: frontend-sync-dashboard
tags: [sync, dashboard, uat, bugfix, date-normalization]
dependency_graph:
  requires: [10-01, 10-02]
  provides: [composite-sync-loading, date-normalized-due-soon]
  affects: [SyncCluster, DashboardView, AdoView, sync-store]
tech_stack:
  added: []
  patterns: [composite-loading-state, start-of-day-normalization, guarded-fetch]
key_files:
  created: []
  modified:
    - frontend/src/stores/sync.ts
    - frontend/src/components/SyncCluster.vue
    - frontend/src/views/AdoView.vue
    - frontend/src/views/DashboardView.vue
decisions:
  - "isFullyLoaded computed aggregates syncing + taskStore.loading + adoStore.loading + adoStore.pipelinesLoading + prStore.loading"
  - "SyncCluster label uses isFullyLoaded for status text; refresh button keeps using syncing (user-triggered action)"
  - "AdoView guards each fetch with emptiness check (same pattern as DashboardView)"
  - "startOfDay helper replicated locally in DashboardView (not exported from tasks.ts)"
metrics:
  duration: "3m"
  completed: "2026-04-09T16:37:00Z"
  tasks_completed: 2
  tasks_total: 2
  files_modified: 4
---

# Phase 10 Plan 03: UAT Gap Closure — SyncCluster Loading State + Attention Bar Date Fix

Composite isFullyLoaded computed in sync store aggregates all frontend loading states; SyncCluster label accurately reflects data readiness; dueSoonTasks uses startOfDay normalization so today's tasks always appear in Attention Bar.

## Tasks Completed

### Task 1: Composite isFullyLoaded + SyncCluster label + AdoView guards
**Commit:** `e71ca96`

**Changes:**
- **sync.ts**: Added `useADOStore` and `usePRStore` imports; created `isFullyLoaded` computed that returns true only when `!syncing && !taskStore.loading && !adoStore.loading && !adoStore.pipelinesLoading && !prStore.loading`; exported in return object. Cleaned up duplicate `useTaskStore()` calls in event handlers.
- **SyncCluster.vue**: Label logic changed from `syncStore.syncing` to `!syncStore.isFullyLoaded` for status text and text color class. Relative time display condition changed to `syncStore.isFullyLoaded`. Refresh button disabled/spinner kept using `syncStore.syncing` (intentional — user-triggered action scope).
- **AdoView.vue**: Replaced unconditional `Promise.all` of 4 fetches with guarded pattern: `workItemTree.length`, `linkedAdoIds.size`, `savedQueries?.length`, `myPRs.length + reviewPRs.length` checks before pushing to fetches array.

### Task 2: Date-normalized dueSoonTasks and upcomingTasks
**Commit:** `cbfafdd`

**Changes:**
- **DashboardView.vue**: Added local `startOfDay(d: Date)` helper normalizing dates to midnight local time. `dueSoonTasks` computed now compares `startOfDay(dueDate)` vs `startOfDay(now)` so today's tasks always have `diffDays=0`. `upcomingTasks` computed uses same normalization for consistency.

## Deviations from Plan

None — plan executed exactly as written.

## Verification Results

1. `vue-tsc --noEmit` — zero TypeScript errors (both tasks)
2. Grep verification:
   - `isFullyLoaded` present in sync.ts (computed + export) and SyncCluster.vue (3 references)
   - AdoView.vue has emptiness guards (`if (!...`) in onMounted
   - DashboardView.vue has `startOfDay` in both dueSoonTasks and upcomingTasks
3. `syncStore.syncing` in SyncCluster.vue only appears in refresh button lines (42-44) — correct

## Self-Check: PASSED
