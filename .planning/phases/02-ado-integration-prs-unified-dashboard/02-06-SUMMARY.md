---
phase: 02-ado-integration-prs-unified-dashboard
plan: 06
subsystem: ui
tags: [vue, pinia, filter-bar, task-list, personal-public-model, quick-add, shadcn-vue]

# Dependency graph
requires:
  - phase: 02-03
    provides: "LinkService with ListPublicTaskIDs, IsPublic for personal→public model"
  - phase: 02-04
    provides: "SyncService with manualSync for sync button integration"
provides:
  - "Enhanced TaskRow with personal/public ADO badge, priority badge, project tag, due date, description preview"
  - "FilterBar component with 5 filter dimensions + sort + group-by + sync button"
  - "Task store extended with publicTaskIds, enhanced filtering, sorting, grouping, quick-add"
  - "TasksView fully integrated with FilterBar, quick-add, sync store, project store"
affects: [02-07, 02-09, 02-10]

# Tech tracking
tech-stack:
  added: []
  patterns: [multi-dimensional-filtering, group-by-toggle, personal-public-badge, quick-add-capture]

key-files:
  created:
    - frontend/src/components/FilterBar.vue
  modified:
    - frontend/src/components/TaskRow.vue
    - frontend/src/stores/tasks.ts
    - frontend/src/views/TasksView.vue

key-decisions:
  - "AcceptableValue type from reka-ui used for Select handler typing"
  - "FilterBar uses Badge chips for status (clickable, highlighted when active) and Select dropdowns for other filters"
  - "enhancedFilteredTasks replaces grouped-by-status as primary task rendering in TasksView"
  - "Quick-add is always-visible Input at top of task list (not toggle-based) for faster capture"
  - "Personal/public ADO badge: hollow dashed circle for personal, filled blue circle with AzureDevOpsIcon for linked"
  - "Group-by rendering resolves project IDs to names via projectNameMap computed"

patterns-established:
  - "Filter chips pattern: Badge variant=default for active, variant=outline for inactive"
  - "Select handlers use AcceptableValue type from reka-ui for type safety with shadcn-vue Select"
  - "publicTaskIds fetched alongside tasks in fetchTasks() for atomic data loading"
  - "Quick-add delegates to store.quickAdd(title) which defaults P2 priority"

requirements-completed: [TASK-08, TASK-09, DASH-01, DASH-02, DASH-03, TL-01, TL-02, TL-03, TL-04, TL-05]

# Metrics
duration: 8min
completed: 2026-04-06
---

# Phase 02 Plan 06: Task List Overhaul Summary

**Enhanced task list with personal/public ADO badges, 5-dimension filter bar, quick-add input, sortable/groupable task rendering integrated with sync and project stores**

## Performance

- **Duration:** 8 min
- **Started:** 2026-04-06T12:16:46Z
- **Completed:** 2026-04-06T12:25:00Z
- **Tasks:** 2
- **Files modified:** 4

## Accomplishments
- TaskRow enhanced with personal/public ADO badge (hollow=personal, filled=linked), priority badge, project tag, due date display (with overdue highlighting), and description preview
- FilterBar component built with 5 filter dimensions (status chips, priority, project, due date, ADO link status), sort dropdown, group-by toggle, and sync button with spinning animation
- Task store extended with publicTaskIds (from LinkService), multi-dimensional filtering, multi-key sorting (priority then due date), groupBy computed, and quickAdd function
- TasksView fully integrated: FilterBar replaces old status tabs, quick-add input always visible at top, grouped rendering when groupBy active, passes isPublic and projectName to each row

## Task Commits

Each task was committed atomically:

1. **Task 1: TaskRow Enhancement + Task Store Extension** - `1f4e246` (feat)
2. **Task 2: FilterBar Component + TasksView Integration** - `965c53c` (feat)

## Files Created/Modified
- `frontend/src/components/FilterBar.vue` - New reusable filter bar with status chips, dropdown filters, sort/group-by controls, sync button
- `frontend/src/components/TaskRow.vue` - Enhanced with isPublic prop, ADO badge (filled/hollow), project tag, due date, description preview, link-task emit
- `frontend/src/stores/tasks.ts` - Extended with publicTaskIds, fetchPublicTaskIds, isPublic, enhancedFilteredTasks, groupedEnhanced, quickAdd, filter/sort/groupBy state
- `frontend/src/views/TasksView.vue` - Integrated FilterBar, quick-add input, syncStore, projectStore, grouped rendering, enhanced filtered tasks

## Decisions Made
- Used `AcceptableValue` type import from `reka-ui` for type-safe Select handlers (auto-fixed TS2322 errors)
- Quick-add is always-visible Input field (not toggle-based like before) — faster capture aligns with D-07 "quick-add: title-only capture"
- FilterBar uses Badge chips for status (most frequently toggled) and Select dropdowns for other filters (compact, less frequent)
- Group-by rendering resolves project IDs to readable names via projectNameMap computed from projectStore
- enhancedFilteredTasks provides unified filter+sort pipeline; original filteredTasks/grouped preserved for backward compatibility

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed Select handler type mismatch with AcceptableValue**
- **Found during:** Task 2 (FilterBar component)
- **Issue:** Select component's `@update:model-value` expects `AcceptableValue` type (from reka-ui), not `string`
- **Fix:** Imported `AcceptableValue` from `reka-ui` and typed all handler functions accordingly
- **Files modified:** frontend/src/components/FilterBar.vue
- **Verification:** `vue-tsc --noEmit` passes with no FilterBar errors
- **Committed in:** 965c53c (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** Type fix necessary for TypeScript compilation. No scope creep.

## Issues Encountered
- Pre-existing TS errors in other files (binding modules not generated yet) — not caused by our changes, ignored

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Task list UI complete — ready for ADO browser view (Plan 07) and dashboard integration (Plan 09/10)
- FilterBar reusable component available for other views
- publicTaskIds infrastructure ready for any view that needs personal/public distinction
- syncStore integration provides manual refresh capability throughout the app

---
*Phase: 02-ado-integration-prs-unified-dashboard*
*Completed: 2026-04-06*
