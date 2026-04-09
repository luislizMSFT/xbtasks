---
phase: 11-integrate-new-task-list-into-real-app
plan: 06
subsystem: ui
tags: [vue, subtasks, pipelines, sync-indicators, ado-integration, dirty-fields]

# Dependency graph
requires:
  - phase: 11-integrate-new-task-list-into-real-app
    provides: TreeTaskRow, FilterCycleButton, QuickAddInput components, ADO meta cache
provides:
  - Playground-parity subtask section in TaskDetail (filter cycle, ADO type icons, badges)
  - Pipeline status icons per PR in TaskDetail
  - Honest dirty-field sync indicators replacing hardcoded "Synced"
  - Amber dirty-field dots on status/priority badges for ADO-linked tasks
affects: [11-07, task-detail, ado-sync]

# Tech tracking
tech-stack:
  added: []
  patterns: [pendingFieldCount computed for field-level dirty tracking, pipelinesForPr branch matching]

key-files:
  created: []
  modified:
    - frontend/src/components/tasks/TaskDetail.vue

key-decisions:
  - "Subtask filter wired with future-proof conditionals — ADO fields use optional chaining since Task type doesn't have syncStatus/source/assignedTo yet"
  - "pendingFieldCount compares edit refs against stored task values for honest sync indicator"
  - "Pipeline matching uses normalized branch names (strips refs/heads/ prefix)"

patterns-established:
  - "Pattern: pendingFieldCount computed for local dirty-field tracking on ADO-linked tasks"
  - "Pattern: pipelinesForPr branch matching for per-PR pipeline status display"
  - "Pattern: Filter-aware empty state messages per subtask filter value"

requirements-completed: [P11-DETAIL-01]

# Metrics
duration: 6min
completed: 2026-04-09
---

# Phase 11 Plan 06: TaskDetail Playground Parity Summary

**Playground-parity subtask section with filter cycle + ADO icons, per-PR pipeline status icons, and honest dirty-field sync indicators**

## Performance

- **Duration:** 6 min
- **Started:** 2026-04-09T12:49:32Z
- **Completed:** 2026-04-09T12:55:29Z
- **Tasks:** 2
- **Files modified:** 1

## Accomplishments
- Subtask section now has filter cycle button (All/Mine/ADO/Personal) matching playground, with type-aware icons and badges
- ADO integration bar shows honest sync status — "{N} pending" with amber indicator when fields are locally modified, "Synced" with green checkmark when clean
- Status and priority badges show amber dirty-field dots when locally modified on ADO-linked tasks
- PR section shows pipeline status icons per PR (green check/red X/blue play/amber clock) with pipeline names
- Push button only appears and shows field count when there are pending changes

## Task Commits

Each task was committed atomically:

1. **Task 1: Port subtask section to playground parity** - `c4bad30` (feat)
2. **Task 2: Add pipeline status icons per PR** - `fad468b` (feat)

## Files Created/Modified
- `frontend/src/components/tasks/TaskDetail.vue` - Subtask filter cycle, ADO type differentiation, sync indicators, PR pipeline status, dirty-field dots

## Decisions Made
- Subtask filter uses optional chaining for future ADO fields (syncStatus, source, assignedTo) not yet on the Task type — forward-proofs template for when ADO subtask sync is implemented
- pendingFieldCount computed compares editTitle/editDescription/editStatus/editPriority against stored task values for real-time dirty tracking
- Pipeline branch matching normalizes by stripping `refs/heads/` prefix for reliable matching

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed stray closing tag in subtask empty state**
- **Found during:** Task 2
- **Issue:** A duplicate `</p>` tag remained from the subtask section replacement in Task 1
- **Fix:** Removed the stray `</p>` tag
- **Files modified:** frontend/src/components/tasks/TaskDetail.vue
- **Verification:** Template renders correctly, no vue-tsc errors
- **Committed in:** fad468b (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** Trivial fix from edit artifact. No scope creep.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- TaskDetail now at playground feature parity for subtasks, PRs, and sync indicators
- Ready for Plan 07 (remaining integration work)
- Pre-existing FilterBar.vue TS errors (unrelated to this plan) remain in codebase

---
*Phase: 11-integrate-new-task-list-into-real-app*
*Completed: 2026-04-09*
