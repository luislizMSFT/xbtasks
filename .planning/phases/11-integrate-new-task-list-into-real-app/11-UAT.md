---
status: diagnosed
phase: 11-integrate-new-task-list-into-real-app
source: [11-01-SUMMARY.md, 11-02-SUMMARY.md, 11-03-SUMMARY.md, 11-04-SUMMARY.md]
started: 2026-04-09T01:06:20Z
updated: 2026-04-09T19:08:24Z
---

## Current Test

[testing complete]

## Tests

### 1. 2-Panel Split Layout Sizing
expected: TasksView shows a balanced 2-panel split. Task list left, detail pane right. Detail pane is large enough to comfortably view and edit task content. No slide-out animation.
result: issue
reported: "the task list is bigger than task detail task detail too small"
severity: major
fix: "Removed flex-1 from left panel, wrapped right panel in w-[45%] shrink-0 container. Commit 0f3c69a"

### 2. Task Detail 3-Tab Structure
expected: Clicking a task shows detail pane with 3 tabs: Subtasks, PRs, Notes. Subtasks tab lists child tasks with status icons and add-subtask input. PRs tab shows linked pull requests. Notes tab renders comments section.
result: issue
reported: "this does not match the playground UI"
severity: major

### 3. TreeTaskRow Display
expected: Each task row shows: status icon (colored circle), ADO type icon (if linked), title, state badge, priority badge, and due date. Blocked tasks show a red blocked-reason banner. Tasks with subtasks show a progress indicator.
result: issue
reported: "wrong no dot projects show old view in the list"
severity: major

### 4. Filter Cycle Button
expected: A small filter button in the toolbar cycles through All → ADO → Personal on each click. The task list filters accordingly — ADO shows only ADO-linked tasks, Personal shows only personal tasks. Active filter is visually highlighted.
result: issue
reported: "There are conflicts"
severity: major

### 5. Quick Add Task
expected: A quick-add input is always visible. Typing a title and pressing Enter creates a new task. Pressing Escape cancels. Input auto-focuses after adding for rapid entry.
result: pass

### 6. Tree View Mode
expected: Toggling tree view shows tasks in parent/child hierarchy with visual indentation. Parent tasks have expand/collapse chevrons. Clicking a chevron toggles children visibility.
result: pass

### 7. Group-by Mode
expected: Selecting group-by (status, priority, or project) shows tasks in collapsible sections with headers showing group name and task count.
result: issue
reported: "no headers divided"
severity: major

### 8. Drag-Drop Reordering
expected: In flat mode (no grouping, no tree), tasks can be dragged to reorder. The new order persists.
result: issue
reported: "no drag and drop"
severity: major

### 9. Dashboard TreeTaskRow
expected: Dashboard view (/dashboard) renders task rows with the new TreeTaskRow styling — status icons, ADO type icons, and ADO metadata badges visible.
result: issue
reported: "seems like old views"
severity: major

### 10. Status Toggle from Row
expected: Clicking the status icon circle on a task row toggles its status (e.g., todo → done). The icon and styling update immediately.
result: issue
reported: "false"
severity: major

## Summary

total: 10
passed: 2
issues: 8
pending: 0
skipped: 0

## Gaps

- truth: "TasksView shows a balanced 2-panel split with detail pane large enough to comfortably view content"
  status: fixed
  reason: "User reported: the task list is bigger than task detail task detail too small. Root cause: left panel had flex-1 which overrode w-[55%], and TaskDetail/ProjectDetail had no width wrapper. Fixed in commit 0f3c69a."
  severity: major
  test: 1
  artifacts: [frontend/src/views/TasksView.vue]
  missing: []

- truth: "Clicking a task shows detail pane with 3 tabs: Subtasks, PRs, Notes"
  status: failed
  reason: "User reported: this does not match the playground UI"
  severity: major
  test: 2
  root_cause: "TaskDetail.vue sections exist but content is far less featured than PlaygroundIntegrated.vue. Missing: subtask filter cycling (All/Mine/ADO/Personal), personal vs ADO subtask differentiation, sync status indicators per subtask, assigned-to badges, ADO state badges on subtasks, pipeline status icons per PR, inline notes with count badge."
  artifacts:
    - path: "frontend/src/components/tasks/TaskDetail.vue"
      issue: "Subtask section is bare-bones vs playground (lines 488-540 vs playground 1105-1222)"
    - path: "frontend/src/views/playground/PlaygroundIntegrated.vue"
      issue: "Gold standard reference (lines 1000-1335)"
  missing:
    - "Port subtask filter cycle button (All/Mine/ADO/Personal)"
    - "Add personal vs ADO subtask differentiation (checkboxes vs ADO type icons)"
    - "Add sync status indicators per subtask (pending/not-pulled)"
    - "Add pipeline status icons per PR"
    - "Add dirty-field amber dots on status/priority badges"
  debug_session: .planning/debug/task-detail-tabs-mismatch.md

- truth: "Each task row shows status icon, ADO type icon, title, state badge, priority badge, due date with new TreeTaskRow styling"
  status: partially-fixed
  reason: "User reported: wrong no dot projects show old view in the list. User clarified: projects in the task list use an old view component entirely, and the old task row for projects has a bug where it just keeps updating."
  severity: major
  test: 3
  root_cause: "groupBy defaulted to 'project' with expandedGroups empty — all groups collapsed showing plain div headers without TreeTaskRow styling. PARTIALLY FIXED in bea5a5a (groupBy default changed to null, auto-expand watcher added). Remaining issue: project group headers still lack TreeTaskRow-level styling, and old project row component may have reactivity loop."
  artifacts:
    - path: "frontend/src/views/TasksView.vue"
      issue: "Project group headers are plain divs, not TreeTaskRow components"
    - path: "frontend/src/stores/tasks.ts"
      issue: "groupBy default fixed to null in bea5a5a"
  missing:
    - "Ensure project group headers use TreeTaskRow-level styling"
    - "Fix old project task row reactivity loop (keeps updating)"
  debug_session: .planning/debug/treetaskrow-old-view.md

- truth: "Filter cycle button cycles All → ADO → Personal, filtering task list accordingly with visual highlight"
  status: partially-fixed
  reason: "User reported: There are conflicts. User clarified: the preview task list header that allows you to toggle personal/ado tasks should be in the new version."
  severity: major
  test: 4
  root_cause: "FilterCycleButton (topbar) and FilterBar Scope dropdown both control taskStore.filterAdoLink with inconsistent labels (ADO/Personal vs Public/Private). Labels unified in bea5a5a but duplicate controls remain."
  artifacts:
    - path: "frontend/src/views/TasksView.vue"
      issue: "Both FilterCycleButton (line 299) and FilterBar (lines 311,321) bind to same filterAdoLink"
    - path: "frontend/src/components/FilterBar.vue"
      issue: "Scope dropdown duplicates FilterCycleButton function"
    - path: "frontend/src/components/tasks/FilterCycleButton.vue"
      issue: "Correct implementation but creates duplicate control"
  missing:
    - "Remove Scope dropdown from FilterBar (FilterCycleButton replaces it)"
    - "Add task list header with task count + filter toggle from playground"
  debug_session: .planning/debug/filter-cycle-conflicts.md

- truth: "Group-by mode shows tasks in collapsible sections with headers showing group name and task count"
  status: partially-fixed
  reason: "User reported: no headers divided"
  severity: major
  test: 7
  root_cause: "Two bugs: (1) expandedGroups started empty — FIXED in bea5a5a. (2) 'No Project' filter removes all headers when tasks lack projectId — still broken. (3) Group headers are project-centric and don't adapt for status/priority grouping modes."
  artifacts:
    - path: "frontend/src/views/TasksView.vue"
      issue: "Line 357: groupKeys.filter(k => k !== 'No Project') eliminates all keys when tasks have no project"
  missing:
    - "Remove 'No Project' filter — render proper header for ungrouped tasks"
    - "Make group header template adaptive for status/priority modes"
  debug_session: .planning/debug/groupby-no-headers.md

- truth: "In flat mode, tasks can be dragged to reorder and the new order persists"
  status: partially-fixed
  reason: "User reported: no drag and drop"
  severity: major
  test: 8
  root_cause: "3 compounding issues: (1) enhancedFilteredTasks ignored sortOrder — FIXED in bea5a5a (manual sort mode added). (2) groupBy default prevented draggable from mounting — FIXED in bea5a5a. (3) Drag handle is invisible (opacity-0) — NOT FIXED."
  artifacts:
    - path: "frontend/src/views/TasksView.vue"
      issue: "Line 508: drag handle uses text-muted-foreground/0 (fully transparent)"
    - path: "frontend/src/components/FilterBar.vue"
      issue: "Missing 'Manual' sort option in dropdown"
  missing:
    - "Make drag handle visible (opacity-40 default, opacity-70 on hover)"
    - "Add 'Manual' sort option to FilterBar sort dropdown"
  debug_session: .planning/debug/dragdrop-missing.md

- truth: "Dashboard view renders task rows with new TreeTaskRow styling — status icons, ADO type icons, ADO metadata badges"
  status: failed
  reason: "User reported: seems like old views"
  severity: major
  test: 9
  root_cause: "DashboardView.vue imports and renders DashboardTaskRow (Phase 10 compact component) instead of TreeTaskRow. Phase 11 updated TasksView but never updated DashboardView. useAdoMeta() call added in bea5a5a but component swap not done."
  artifacts:
    - path: "frontend/src/views/DashboardView.vue"
      issue: "Imports DashboardTaskRow, not TreeTaskRow (line 12); renders it in Focus/Upcoming/Blocked sections"
    - path: "frontend/src/components/tasks/DashboardTaskRow.vue"
      issue: "Old component missing status toggle, progress bars, ADO state badges"
  missing:
    - "Replace DashboardTaskRow with TreeTaskRow in DashboardView (or enhance DashboardTaskRow with missing features)"
    - "Wire TreeTaskRow props (adoMeta, isPublic, subtaskProgress) via useAdoMeta composable"
    - "Replace inline blocked-task markup (lines 370-386) with component usage"
  debug_session: .planning/debug/dashboard-old-taskrow.md

- truth: "Clicking status icon circle on task row toggles status and updates icon/styling immediately"
  status: partially-fixed
  reason: "User reported: false"
  severity: major
  test: 10
  root_cause: "At UAT time, TasksView imported old TaskRow instead of TreeTaskRow. Combined with collapsed groups, user had nothing to click. TaskRow→TreeTaskRow swap done in bea5a5a. Toggle mechanism was always correctly wired. Needs re-test."
  artifacts:
    - path: "frontend/src/views/TasksView.vue"
      issue: "Fixed in bea5a5a — TaskRow replaced with TreeTaskRow"
  missing:
    - "Re-test to confirm status toggle works in current build"
  debug_session: .planning/debug/status-toggle-broken.md
