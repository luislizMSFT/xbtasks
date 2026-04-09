---
phase: 11-integrate-new-task-list-into-real-app
plan: 04
title: "Rewrite TasksView with 2-Panel Split, TreeTaskRow, FilterCycleButton & ADO Metadata"
subsystem: frontend
tags: [tasks-view, tree-task-row, filter-cycle, ado-meta, 2-panel-split]
dependency_graph:
  requires: [TreeTaskRow, QuickAddInput, FilterCycleButton, useAdoMeta, TaskDetail-3-tabs]
  provides: [TasksView-production]
  affects: [TasksView]
tech_stack:
  added: []
  patterns: [permanent-split-panel, composable-wiring, component-extraction-reuse]
key_files:
  created: []
  modified:
    - frontend/src/views/TasksView.vue
decisions:
  - "QuickAddInput replaces raw Input for quick-add — always visible with Add/Cancel buttons"
  - "FilterCycleButton placed in topbar Teleport next to status chips for easy access"
  - "TreeTaskRow used in all three rendering modes: flat with drag-drop, grouped, tree"
  - "adoMeta.getAdoMeta() passed as prop to every TreeTaskRow instance for batch ADO metadata"
metrics:
  duration: "4m 23s"
  completed: "2026-04-08T17:51:41Z"
  tasks_completed: 1
  tasks_total: 2
  files_created: 0
  files_modified: 1
---

# Phase 11 Plan 04: Rewrite TasksView with 2-Panel Split, TreeTaskRow & ADO Metadata Summary

TasksView.vue fully rewritten — TaskRow replaced with TreeTaskRow everywhere (flat/grouped/tree), FilterCycleButton added to topbar, QuickAddInput replaces raw Input, useAdoMeta wired to every row for batch ADO type/state display.

## What Was Done

### Task 1: Rewrite TasksView.vue with permanent 2-panel split layout
**Commit:** `69ff237`

- **Replaced all `TaskRow` usages with `TreeTaskRow`** in flat (drag-drop), grouped (collapsible group headers), and tree (parent→child→grandchild) rendering modes
- **Added `FilterCycleButton`** to topbar Teleport alongside status chips — cycles All/ADO/Personal filter via `v-model="taskStore.filterAdoLink"`
- **Replaced raw `<Input>` quick-add** with `QuickAddInput` component — accepts `@add` emit with title string instead of managing `v-model` locally
- **Wired `useAdoMeta` composable** — `adoMeta.getAdoMeta(task.id)` passed as `:ado-meta` prop to every TreeTaskRow instance for batch-loaded ADO metadata
- **Added new TreeTaskRow props**: `:subtask-progress`, `:is-public`, `:project-name`, `:ado-meta` on all rendering paths
- **Preserved all 13 existing features:**
  1. Status chip tabs in topbar via Teleport (6 chips)
  2. Sync button in topbar (RefreshCw, animate-spin)
  3. FilterBar with all dropdowns (priority, project, dueDate, adoLink, sortBy, groupBy, treeView)
  4. Three rendering modes: flat with drag-drop, grouped by status/priority/project, tree view
  5. Drag-drop reordering in flat mode (vuedraggable with handle)
  6. Group-by mode with collapsible headers (project ADO type icons, badges, counts)
  7. Tree view with expand/collapse (3 depth levels)
  8. Quick-add task input (always visible at top)
  9. Loading skeleton state (6 rows)
  10. Empty state (ClipboardList icon)
  11. Keyboard shortcuts (Meta+N, Ctrl+N, Escape)
  12. Project detail panel (click group header → ProjectDetail)
  13. Right panel: TaskDetail / ProjectDetail / empty placeholder
- **Removed unused imports**: `Input` (replaced by QuickAddInput), `adoStateClasses`/`adoPriorityClasses` (not needed — TreeTaskRow handles internally)
- **Removed unused function**: `toggleSubtaskDone` (subtask toggling now handled inside TaskDetail)
- **`handleQuickAdd` simplified**: accepts `title: string` from QuickAddInput emit instead of reading `quickAddTitle.value`

### Task 2: Human Verification (checkpoint)
**Status:** Awaiting user verification

## Deviations from Plan

None — plan executed exactly as written.

## Decisions Made

1. **QuickAddInput always visible:** Kept the always-visible pattern (matching current UX) rather than toggle-based `showQuickAdd`
2. **FilterCycleButton placement:** Added inside the Teleport alongside status Badge chips for natural topbar integration
3. **TreeTaskRow in all modes:** All three rendering paths (flat, grouped, tree) now use TreeTaskRow with consistent props
4. **adoMeta prop passing:** Every TreeTaskRow instance receives `adoMeta.getAdoMeta(task.id)` for batch-loaded ADO metadata

## Verification

- All 20 acceptance criteria pass (checked via automated script)
- `vue-tsc --noEmit --skipLibCheck` — zero errors
- TasksView.vue contains no references to old `TaskRow` component
- TasksView.vue is 557 lines (well above 200 minimum)
- All imports, composables, and store connections verified

## Self-Check: PASSED

- Modified file `frontend/src/views/TasksView.vue` exists
- Commit `69ff237` found
- SUMMARY.md exists
