---
phase: 11-integrate-new-task-list-into-real-app
plan: 02
title: "Extract TreeTaskRow, QuickAddInput, FilterCycleButton Components"
subsystem: frontend
tags: [vue-components, tree-task-row, quick-add, filter-cycle, extraction]
dependency_graph:
  requires: [ado_meta_cache_table, useAdoMeta, statusIcon, priorityDotBgColor]
  provides: [TreeTaskRow, QuickAddInput, FilterCycleButton]
  affects: [TasksView, DashboardView, PlaygroundIntegrated]
tech_stack:
  added: []
  patterns: [props-down-events-up, v-model-pattern, component-extraction]
key_files:
  created:
    - frontend/src/components/tasks/TreeTaskRow.vue
    - frontend/src/components/tasks/QuickAddInput.vue
    - frontend/src/components/tasks/FilterCycleButton.vue
  modified: []
decisions:
  - "TreeTaskRow uses formatDate from lib/date.ts for due date display (consistent with existing TaskRow)"
  - "FilterCycleButton uses 'linked' value (not 'ado') to match taskStore.filterAdoLink semantics"
  - "QuickAddInput uses Button component from shadcn-vue for Add/Cancel (consistent with design system)"
metrics:
  duration: "2m 31s"
  completed: "2026-04-09T00:35:44Z"
  tasks_completed: 2
  tasks_total: 2
  files_created: 3
  files_modified: 0
---

# Phase 11 Plan 02: Extract TreeTaskRow, QuickAddInput, FilterCycleButton Components Summary

Three standalone Vue components extracted from PlaygroundIntegrated.vue — TreeTaskRow (195-line unified task row with indent/expand/status/ADO support), QuickAddInput (Enter/Esc/auto-focus quick-add), FilterCycleButton (all→ADO→personal v-model cycle).

## What Was Done

### Task 1: Create TreeTaskRow.vue — unified task row component
**Commit:** `c48a35d`

- Created `TreeTaskRow.vue` (195 lines) with 2-row layout extracted from PlaygroundIntegrated.vue
- Row 1: expand chevron, status icon (via `statusIcon()`), ADO type icon (via `adoTypeIcon()`), title, state badge (ADO state or task status), subtask progress bar
- Row 2: ADO ID badge, personal badge, project name, area, due date with overdue detection, PriorityBadge component
- Row 3 (conditional): blocked reason banner for blocked tasks
- Props: `task` (Task type from `@/types`), `indentLevel`, `selected`, `expanded`, `hasChildren`, `adoMeta`, `projectName`, `subtaskProgress`, `isPublic`
- Emits: `click`, `toggleExpand`, `toggleStatus(id)`
- Imports all style helpers from `@/lib/styles` — no inline switch statements
- Due date display ported from TaskRow.vue with overdue/today logic using `formatDate` from `@/lib/date`

### Task 2: Create QuickAddInput.vue and FilterCycleButton.vue
**Commit:** `fe72e6d`

- Created `QuickAddInput.vue` (43 lines):
  - Auto-focuses input on mount via `onMounted` + `nextTick`
  - Emits `add(title)` on Enter key or Add button click
  - Emits `cancel` on Escape key or Cancel button click
  - Re-focuses input after successful add for rapid entry
  - Add button disabled when title is empty
  - Uses Circle icon + shadcn Button components
- Created `FilterCycleButton.vue` (49 lines):
  - v-model pattern with `modelValue` prop and `update:modelValue` emit
  - Cycles through `all` → `linked` → `personal` on click
  - Dynamic icon: ListChecks (all), Lock (linked), User (personal)
  - Visual active state: `bg-primary/10 text-primary font-medium` when not 'all'

## Deviations from Plan

None — plan executed exactly as written.

## Decisions Made

1. **Date formatting:** TreeTaskRow uses `formatDate` from `@/lib/date` for due date display (same utility as existing TaskRow)
2. **Filter values:** FilterCycleButton uses 'linked' (not 'ado') as the value for ADO-linked filter, matching the existing taskStore.filterAdoLink semantics
3. **Button component:** QuickAddInput uses shadcn-vue Button for Add/Cancel to stay consistent with the design system

## Verification

- All three component files exist in `frontend/src/components/tasks/`
- `vue-tsc --noEmit --skipLibCheck` produces zero errors for new files
- TreeTaskRow imports `statusIcon`, `adoTypeIcon`, and all style helpers from `@/lib/styles`
- TreeTaskRow uses `Task` type from `@/types` (not custom interface)
- QuickAddInput handles Enter, Escape, and auto-focuses on mount
- FilterCycleButton uses v-model pattern (`modelValue` + `update:modelValue`)

## Self-Check: PASSED

- All 3 created files exist
- SUMMARY.md exists
- Commit c48a35d found (TreeTaskRow)
- Commit fe72e6d found (QuickAddInput + FilterCycleButton)
