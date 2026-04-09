---
phase: 10-implement-dashboard-redesign-and-unified-header-bars
plan: 01
title: "SyncCluster + 3-Zone Header Bar"
subsystem: frontend/ui
tags: [header, sync, teleport, layout]
dependency_graph:
  requires: []
  provides: [SyncCluster.vue, topbar-center, priorityBgColor]
  affects: [AppShell.vue, DashboardView, TasksView, AdoView, ProjectsView, DependencyGraphView, SettingsView]
tech_stack:
  added: []
  patterns: [3-zone-header, shared-sync-cluster, teleport-center]
key_files:
  created:
    - frontend/src/components/SyncCluster.vue
  modified:
    - frontend/src/lib/styles.ts
    - frontend/src/layouts/AppShell.vue
    - frontend/src/views/DashboardView.vue
    - frontend/src/views/TasksView.vue
    - frontend/src/views/AdoView.vue
    - frontend/src/views/ProjectsView.vue
    - frontend/src/views/DependencyGraphView.vue
    - frontend/src/views/SettingsView.vue
decisions:
  - "3-zone top bar: left (page name + divider + SyncCluster), center (#topbar-center), right (search+new+activity)"
  - "SyncCluster replaces per-page sync buttons — single shared component in AppShell left zone"
  - "All 6 views teleport to #topbar-center — unified pattern"
  - "AdoView teleport trimmed to Tabs only — sync/connection UI moved to SyncCluster"
  - "TasksView teleport trimmed to status filter chips only — sync button moved to SyncCluster"
metrics:
  duration: "27 min"
  completed: "2026-04-08T23:47:42Z"
  tasks_completed: 2
  tasks_total: 2
  files_created: 1
  files_modified: 8
---

# Phase 10 Plan 01: SyncCluster + 3-Zone Header Bar Summary

SyncCluster shared component with 3-state sync status (connected/syncing/offline) + AppShell refactored to 3-zone top bar + all 6 views migrated to #topbar-center teleport

## What Was Done

### Task 1: Create SyncCluster component + add priorityBgColor helper
- **Created** `frontend/src/components/SyncCluster.vue` — shared sync status component
  - Shows green/red dot based on ADO connection state
  - 3 states: "Synced" (green), "Syncing…" (muted with spinner), "Offline" (red)
  - Displays relative time since last sync
  - Pending conflicts badge (amber, shown when > 0)
  - Refresh button with `aria-label="Refresh sync"` for accessibility
  - Calls `syncStore.manualSync()` on click
- **Added** `priorityBgColor()` to `frontend/src/lib/styles.ts`
  - Returns Tailwind bg classes: P0=red, P1=orange, P2=amber, P3=zinc
- **Commit:** `922cdd6`

### Task 2: Refactor AppShell to 3-zone top bar + update all 6 view teleports
- **Refactored** `frontend/src/layouts/AppShell.vue`:
  - Replaced 2-zone layout with 3-zone: left (page name + divider + SyncCluster), center (#topbar-center), right (search+new+activity)
  - Removed `#topbar-actions` teleport target entirely
  - Added `/dependencies` to breadcrumb map
  - Changed page name from `font-medium` to `font-semibold`
  - Added `shrink-0` to right zone wrapper
- **Updated** `DashboardView.vue`: teleport → `#topbar-center`
- **Updated** `TasksView.vue`: teleport → `#topbar-center`, removed sync RefreshCw button from teleport content (now in SyncCluster)
- **Updated** `AdoView.vue`: teleport → `#topbar-center`, removed Connected/Offline pill, Sync button, lastSyncLabel, loading spinner from teleport (kept only Tabs)
- **Added teleport** to `ProjectsView.vue`: project count + active count badge
- **Added teleport** to `DependencyGraphView.vue`: node/edge stats + cycle detection badge
- **Added teleport** to `SettingsView.vue`: empty teleport (no center content needed)
- **Commit:** `08f5676`

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing] Removed unused RefreshCw import from TasksView**
- **Found during:** Task 2, step C
- **Issue:** After removing the sync button from TasksView teleport, `RefreshCw` import was unused
- **Fix:** Removed `RefreshCw` from lucide import
- **Files modified:** frontend/src/views/TasksView.vue

**2. [Rule 2 - Missing] Removed unused Wifi/WifiOff/RefreshCw imports from AdoView**
- **Found during:** Task 2, step D
- **Issue:** After removing connection pill and sync button from AdoView teleport, `Wifi`, `WifiOff`, `RefreshCw` imports were unused
- **Fix:** Removed from lucide import (kept `Clock` and `cn` which are still used elsewhere)
- **Files modified:** frontend/src/views/AdoView.vue

## Verification Results

- ✅ Zero `topbar-actions` references in frontend/src/
- ✅ `topbar-center` found in AppShell + all 6 views
- ✅ `SyncCluster` found in AppShell.vue (import + template)
- ✅ `aria-label="Refresh sync"` in SyncCluster.vue
- ✅ `priorityBgColor` exported from styles.ts
- ✅ All 3 SyncCluster states handled (connected/syncing/offline)
