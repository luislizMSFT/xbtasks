---
phase: 11-integrate-new-task-list-into-real-app
plan: 05
title: "Fix broken interactions and missing UI elements"
one_liner: "Visible drag handles, adaptive group headers, task list header with expand/collapse, Scope dropdown removed from FilterBar"
subsystem: frontend-tasksview
tags: [drag-drop, group-headers, filter-cleanup, task-list-header, expand-collapse]
dependency_graph:
  requires: []
  provides: [adaptive-group-headers, task-list-header, expand-collapse-controls, filterbar-cleanup]
  affects: [TasksView, FilterBar]
tech_stack:
  added: []
  patterns: [adaptive-template-switching, statusBgColor, priorityDotBgColor]
key_files:
  created: []
  modified:
    - frontend/src/views/TasksView.vue
    - frontend/src/components/FilterBar.vue
decisions:
  - "Group headers adapt based on groupBy mode: status dots, priority dots, or project ADO icons"
  - "No Project group included in v-for loop with proper header (not filtered out or rendered headerless)"
  - "FilterCycleButton dual-placed in topbar and task list header (intentional dual placement)"
  - "Scope dropdown removed from FilterBar — ADO/Personal filter exclusively via FilterCycleButton"
metrics:
  duration: "3m 32s"
  completed: "2026-04-09T19:52:52Z"
  tasks_completed: 2
  tasks_total: 2
  files_modified: 2
requirements: [P11-TASKLIST-01]
---

# Phase 11 Plan 05: Fix Broken Interactions and Missing UI Elements Summary

Visible drag handles, adaptive group headers (status dots / priority dots / project ADO icons), task list header with count + FilterCycleButton + expand/collapse buttons, and FilterBar cleaned of duplicate Scope dropdown.

## Changes Made

### Task 1: Fix TasksView — drag handle, adaptive group headers, task list header, expand/collapse
**Commit:** `b8bbf08`

- **Gap A (Drag handle):** Changed opacity from `/0` to `/40` (hover `/70`) — drag grip dots are now visible on hover
- **Gap B (No Project filter):** Removed `.filter(k => k !== 'No Project')` so all group keys render, including 'No Project' with a proper header
- **Gap C (Adaptive group headers):** Refactored group header template:
  - `groupBy === 'status'`: colored status dot via `statusBgColor(key)`
  - `groupBy === 'priority'`: colored priority dot via `priorityDotBgColor(key)`
  - `groupBy === 'project'`: ADO type icon + project badge (unchanged behavior)
  - Click behavior: `selectProject` only for project mode; status/priority click toggles group expand
- **Gap I (Task list header):** Added header bar between FilterBar and ScrollArea showing `taskListCount` computed and `FilterCycleButton`
- **Gap K (Expand/Collapse):** Added Expand All / Collapse All buttons in task list header, visible when `groupBy` or `treeView` is active
- Removed `filterAdoLink` prop and `@update:filter-ado-link` emit from FilterBar usage

### Task 2: Remove Scope dropdown from FilterBar
**Commit:** `2699dea`

- Removed `filterAdoLink` from props interface
- Removed `update:filterAdoLink` from emits
- Removed `handleAdoLinkChange` function
- Removed Scope label and Select dropdown from template
- Removed `filterAdoLink` from `activeFilterCount` computation
- FilterBar now has: Priority, Project, Due, Sort, Group, Tree toggle (no Scope)

## Verification Results

- ✅ TypeScript compiles cleanly (`vue-tsc --noEmit --skipLibCheck` — no errors)
- ✅ No `text-muted-foreground/0` in TasksView (drag handle fixed)
- ✅ No `filter.*No Project` in TasksView (No Project filter removed)
- ✅ No `Scope` in FilterBar (Scope dropdown removed)
- ✅ `taskListCount` present in TasksView (task list header added)
- ✅ `expandAll`/`collapseAll` present in TasksView (expand/collapse added)

## Deviations from Plan

None — plan executed exactly as written.

## Self-Check: PASSED

- ✅ frontend/src/views/TasksView.vue exists
- ✅ frontend/src/components/FilterBar.vue exists
- ✅ 11-05-SUMMARY.md exists
- ✅ Commit b8bbf08 found
- ✅ Commit 2699dea found
