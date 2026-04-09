---
status: complete
phase: 11-integrate-new-task-list-into-real-app
source: [11-01-SUMMARY.md, 11-02-SUMMARY.md, 11-03-SUMMARY.md, 11-04-SUMMARY.md]
started: 2026-04-09T01:06:20Z
updated: 2026-04-09T18:57:44Z
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
  artifacts: []
  missing: []

- truth: "Each task row shows status icon, ADO type icon, title, state badge, priority badge, due date with new TreeTaskRow styling"
  status: failed
  reason: "User reported: wrong no dot projects show old view in the list"
  severity: major
  test: 3
  artifacts: []
  missing: []

- truth: "Filter cycle button cycles All → ADO → Personal, filtering task list accordingly with visual highlight"
  status: failed
  reason: "User reported: There are conflicts"
  severity: major
  test: 4
  artifacts: []
  missing: []

- truth: "Group-by mode shows tasks in collapsible sections with headers showing group name and task count"
  status: failed
  reason: "User reported: no headers divided"
  severity: major
  test: 7
  artifacts: []
  missing: []

- truth: "In flat mode, tasks can be dragged to reorder and the new order persists"
  status: failed
  reason: "User reported: no drag and drop"
  severity: major
  test: 8
  artifacts: []
  missing: []

- truth: "Dashboard view renders task rows with new TreeTaskRow styling — status icons, ADO type icons, ADO metadata badges"
  status: failed
  reason: "User reported: seems like old views"
  severity: major
  test: 9
  artifacts: []
  missing: []

- truth: "Clicking status icon circle on task row toggles status and updates icon/styling immediately"
  status: failed
  reason: "User reported: false"
  severity: major
  test: 10
  artifacts: []
  missing: []
