---
phase: 11-integrate-new-task-list-into-real-app
plan: 03
title: "Rebuild TaskDetail 3-Tab Structure & DashboardView TreeTaskRow"
subsystem: frontend
tags: [task-detail, tabs, tree-task-row, dashboard, component-reuse]
dependency_graph:
  requires: [TreeTaskRow, useAdoMeta, statusIcon, priorityDotBgColor, CommentsSection, ExternalLinks, SyncConfirmDialog, ConflictResolver]
  provides: [TaskDetail-3-tabs, DashboardView-TreeTaskRow]
  affects: [TasksView, DashboardView]
tech_stack:
  added: []
  patterns: [permanent-panel, 3-tab-layout, composable-wiring]
key_files:
  created: []
  modified:
    - frontend/src/components/tasks/TaskDetail.vue
    - frontend/src/views/DashboardView.vue
decisions:
  - "Removed outer Transition wrapper entirely — TaskDetail now renders as permanent panel controlled by parent sizing"
  - "Moved project select from header badges into a dedicated project context bar (Row 4) for cleaner layout"
  - "Push to ADO moved into ADO dropdown menu alongside Open in ADO and Unlink for consolidated ADO actions"
  - "Subtask icons use statusIcon() from styles.ts instead of inline checkbox SVGs for consistency with TreeTaskRow"
  - "Delete button moved to header Row 2 (inline with status/priority) instead of footer for quicker access"
  - "statusLabel renamed to statusLabelMap to avoid shadowing template variable conflicts"
metrics:
  duration: "6m 16s"
  completed: "2026-04-09T00:44:23Z"
  tasks_completed: 2
  tasks_total: 2
  files_created: 0
  files_modified: 2
---

# Phase 11 Plan 03: Rebuild TaskDetail 3-Tab Structure & DashboardView TreeTaskRow Summary

TaskDetail.vue rebuilt with permanent-panel layout and 3-tab structure (Subtasks/PRs/Notes replacing Work/Discussion), DashboardView.vue swapped from TaskRow to TreeTaskRow with ADO metadata badges.

## What Was Done

### Task 1: Rebuild TaskDetail.vue with 3-tab structure
**Commit:** `e798388`

- **Removed slide-out Transition wrapper** — no more `translate-x-full` animation; `<aside>` renders directly as permanent panel (`v-if="isOpen || lastTask"`)
- **Replaced Work/Discussion tab layout with 3 tabs:**
  - **Subtasks tab** — subtask list with `statusIcon()` components and `priorityDotBgColor()` dots, progress bar, add subtask input
  - **PRs tab** — PR list filtered by task ID/ADO ID with status badges, branch names via `branchName()`, and "Open" click handler
  - **Notes tab** — `<CommentsSection>` component rendered directly
- **Restructured header** into 4 rows:
  - Row 1: Editable title + close (X) button
  - Row 2: Status select + Priority select + due date + delete button
  - Row 3 (if ADO linked): ADO badge dropdown with Open/Push/Unlink actions
  - Row 4 (if project): Folder icon + project name context bar
- **Description section** moved above tabs as collapsible section with border-b separator
- **External Links** rendered as own section between description and tabs
- **Preserved all functional logic:** `lastTask` ref pattern, `save()` with dirty-check and race-guard, `loadSubtasks()` watcher, `taskPRs` computed, ADO URL building via `getCachedWorkItem()`, `unlinkFromADO()`, sync integration, tag add/remove
- **SyncConfirmDialog + ConflictResolver** rendered outside tabs at bottom of aside
- **Footer** simplified to timestamps only (delete moved to header)

### Task 2: Update DashboardView.vue to use TreeTaskRow
**Commit:** `23ee5f5`

- Replaced `import TaskRow` with `import TreeTaskRow`
- Added `import { useAdoMeta } from '@/composables/useAdoMeta'` and `const adoMeta = useAdoMeta()`
- Replaced `<TaskRow>` in Today's Focus section with `<TreeTaskRow>` + `:is-public`, `:ado-meta` props
- Replaced `<TaskRow>` in Blocked section with `<TreeTaskRow>` + `:is-public`, `:ado-meta` props
- Mapped `@select` emit to `@click` emit pattern (TreeTaskRow emits `click` not `select`)
- Added `@toggle-status` handler for blocked tasks with status toggle logic

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] statusLabel variable shadowing**
- **Found during:** Task 1
- **Issue:** Template variable `statusLabel` (Record) conflicted with potential reactive/template naming
- **Fix:** Renamed to `statusLabelMap` to avoid any template shadowing
- **Files modified:** TaskDetail.vue

**2. [Rule 2 - Missing functionality] Project name display**
- **Found during:** Task 1
- **Issue:** Plan specified Row 4 with project name but the original only had a project Select dropdown, not a display bar
- **Fix:** Added `projectName` computed that looks up project name from projectStore, rendered in Row 4 with Folder icon
- **Files modified:** TaskDetail.vue

## Decisions Made

1. **Permanent panel:** Removed all Transition wrappers — the parent TasksView now controls 2-panel split sizing
2. **ADO actions consolidated:** Push to ADO moved from separate button into ADO dropdown menu for cleaner UI
3. **Style helpers:** Used `statusIcon()` and `priorityDotBgColor()` from styles.ts instead of inline checkbox SVGs and hardcoded maps
4. **Delete placement:** Moved from footer to header row 2 for faster access, matching common task-app patterns
5. **statusLabelMap:** Renamed from `statusLabel` to avoid any shadowing with template usage

## Verification

- TaskDetail.vue has `<TabsTrigger value="subtasks">`, `<TabsTrigger value="prs">`, `<TabsTrigger value="notes">`
- TaskDetail.vue has no `translate-x-full` (no slide-out animation)
- TaskDetail.vue imports all 4 sub-components: CommentsSection, ExternalLinks, SyncConfirmDialog, ConflictResolver
- TaskDetail.vue preserves `lastTask`, `save()`, `loadSubtasks`, `taskPRs`
- DashboardView.vue has zero references to old `TaskRow`
- DashboardView.vue uses `TreeTaskRow` with `adoMeta.getAdoMeta()` prop
- `vue-tsc --noEmit --skipLibCheck` produces zero errors for modified files

## Self-Check: PASSED

- All 2 modified files exist
- SUMMARY.md exists
- Commit e798388 found (TaskDetail rebuild)
- Commit 23ee5f5 found (DashboardView TreeTaskRow)
