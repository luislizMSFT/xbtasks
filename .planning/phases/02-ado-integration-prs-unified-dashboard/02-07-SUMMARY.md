---
phase: 02-ado-integration-prs-unified-dashboard
plan: 07
subsystem: ui
tags: [vue, tree-view, ado-browser, dialog, pinia, lucide, saved-queries]

requires:
  - phase: 02-03
    provides: ADOService REST methods (ListMyWorkItems, GetWorkItemTree), LinkService (LinkTask, PromoteTask, ImportWorkItem)
  - phase: 02-04
    provides: SyncService, CommentService, backend wiring in main.go
provides:
  - ADO Browser tree view with hierarchy (Scenario -> Deliverable -> Task)
  - LinkDialog for connecting local tasks to ADO items
  - PromoteDialog for creating ADO items from local tasks
  - Extended ADO store with tree data, filtering, and saved query support
  - GetSavedQueries and RunSavedQuery REST API methods
  - ListLinkedAdoIDs backend method for linked status tracking
affects: [02-08, 02-09, frontend-views]

tech-stack:
  added: []
  patterns: [recursive-render-function-component, pinia-computed-tree-filtering, dynamic-wails-binding-imports]

key-files:
  created:
    - frontend/src/components/AdoTreeBrowser.vue
    - frontend/src/components/LinkDialog.vue
    - frontend/src/components/PromoteDialog.vue
  modified:
    - frontend/src/stores/ado.ts
    - frontend/src/views/AdoView.vue
    - internal/app/adoservice.go
    - internal/app/linkservice.go
    - pkg/ado/query.go

key-decisions:
  - "Recursive tree node via defineComponent render function (inline in AdoTreeBrowser) instead of separate file"
  - "ADO store getChildren applies same filters as treeRoots for consistent filtered tree"
  - "ListLinkedAdoIDs added to LinkService (plan had ListPublicTaskIDs which returns task IDs, not ADO IDs)"
  - "Saved Queries use __my_assignments__ sentinel value to revert to default tree fetch"
  - "LinkDialog searches both workItems and workItemTree for maximum coverage"

patterns-established:
  - "Render function components: Use defineComponent with h() for recursive tree nodes that need same-file recursion"
  - "Store-driven filtering: Filters live in Pinia store, computed properties derive filtered results"
  - "Dialog pattern: v-model:open for two-way binding, emit events on success, reset state on close via watcher"

requirements-completed: [ADO-01, ADO-02, ADO-06, ADO-10, UX-01, UX-04, UX-05]

duration: 8min
completed: 2026-04-06
---

# Phase 02 Plan 07: ADO Browser View Summary

**ADO tree browser with hierarchy navigation, type-filtered search, saved query picker, linked status tracking, and link/import/promote dialogs**

## Performance

- **Duration:** ~8 min
- **Started:** 2026-04-06T19:17:12Z
- **Completed:** 2026-04-06T19:25:14Z
- **Tasks:** 2
- **Files modified:** 7

## Accomplishments

- Extended ADO store with full tree state: workItemTree, filters (type/state/area), hideLinked toggle, savedQueries, treeRoots computed
- Built recursive AdoTreeBrowser component with type-colored icons (Bug=red, Task=blue, Story=purple, Feature=green, Epic=orange), state badges, linked indicators, and import/link hover actions
- Created LinkDialog: search ADO items by ID or title, select and confirm to connect local task
- Created PromoteDialog: editable title, work item type selector (Task/Bug/User Story), subtask warning, create-in-ADO flow
- Overhauled AdoView: replaced flat Work Items tab with tree-based ADO Browser tab, added filter toolbar with type/state/area/saved-query pickers and hide-linked toggle, two-column layout with tree (60%) and detail panel (40%)
- Added Go backend: GetSavedQueries/RunSavedQuery in ADOService, ListLinkedAdoIDs in LinkService, pkg/ado saved query REST calls

## Task Commits

Each task was committed atomically:

1. **Task 1: ADO Store Extension + ADO Tree Browser Component** - `b98ec58` (feat)
2. **Task 2: LinkDialog + PromoteDialog + AdoView Overhaul** - `1ebf854` (feat)

## Files Created/Modified

- `frontend/src/stores/ado.ts` — Extended with tree state, filters, savedQueries, treeRoots computed, getChildren, fetchWorkItemTree, fetchLinkedAdoIds, fetchSavedQueries, runSavedQuery
- `frontend/src/components/AdoTreeBrowser.vue` — Recursive tree browser with type icons, state badges, linked indicators, import/link/open actions
- `frontend/src/components/LinkDialog.vue` — Dialog to link existing local task to ADO item via search
- `frontend/src/components/PromoteDialog.vue` — Dialog to create new ADO work item from local task
- `frontend/src/views/AdoView.vue` — Overhauled with tree browser tab, filter toolbar, two-column detail panel
- `internal/app/adoservice.go` — Added GetSavedQueries, RunSavedQuery methods
- `internal/app/linkservice.go` — Added ListLinkedAdoIDs method
- `pkg/ado/query.go` — Added GetSavedQueries, RunSavedQuery REST API functions

## Decisions Made

- **Recursive render function**: Used `defineComponent` with `h()` for the recursive tree node inside AdoTreeBrowser.vue, avoiding need for a separate file while keeping TypeScript support
- **ListLinkedAdoIDs**: Plan referenced ListPublicTaskIDs (returns task IDs), but frontend needs ADO IDs for linked status. Added new ListLinkedAdoIDs method — Rule 2 auto-fix for missing critical functionality
- **Saved query sentinel value**: Used `__my_assignments__` as a special value in the saved query picker to revert to the default fetchWorkItemTree() call
- **Dual-source search in LinkDialog**: Searches both cached workItems and workItemTree for maximum coverage when linking

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Added ListLinkedAdoIDs method to LinkService**
- **Found during:** Task 1 (ADO Store Extension)
- **Issue:** Plan's fetchLinkedAdoIds referenced ListPublicTaskIDs which returns task IDs, but frontend needs ADO IDs for the isLinked() check
- **Fix:** Added ListLinkedAdoIDs to LinkService that queries `SELECT DISTINCT ado_id FROM task_ado_links`
- **Files modified:** internal/app/linkservice.go
- **Verification:** Go compilation succeeds
- **Committed in:** b98ec58 (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (1 missing critical)
**Impact on plan:** Essential for correctness of linked status tracking. No scope creep.

## Issues Encountered

None — plan executed cleanly.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- ADO browser view complete with tree hierarchy, filtering, and import/link actions
- LinkDialog and PromoteDialog ready for use from TaskDetail and other views
- PromoteDialog currently uses first configured org/project — may need org/project selector in future
- Sync confirmation (Plan 08) and conflict resolution can now reference the linked status infrastructure

---
*Phase: 02-ado-integration-prs-unified-dashboard*
*Plan: 07*
*Completed: 2026-04-06*
