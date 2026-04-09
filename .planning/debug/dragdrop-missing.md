---
status: diagnosed
trigger: "Drag and drop reordering is not working"
created: 2025-01-09T12:00:00Z
updated: 2025-01-09T12:05:00Z
---

## Current Focus

hypothesis: CONFIRMED - Drag-drop is wired up but functionally broken: enhancedFilteredTasks force-sorts by sortBy (priority), ignoring sortOrder; default groupBy='project' prevents flat-mode rendering; drag handle is invisible
test: n/a - root cause confirmed via code analysis
expecting: n/a
next_action: Return diagnosis

## Symptoms

expected: In flat mode (no grouping, no tree), tasks can be dragged to reorder. The new order persists.
actual: User reported "no drag and drop"
errors: None reported
reproduction: Test 8 in UAT - try to drag tasks in flat mode
started: Discovered during UAT

## Eliminated

- hypothesis: Drag-and-drop was never integrated (not in TasksView)
  evidence: vuedraggable IS imported and used in TasksView.vue lines 5, 495-529 with drag handle at line 508
  timestamp: 2025-01-09T12:01:00Z

- hypothesis: Backend reorderTasks not implemented
  evidence: ReorderTasks exists in internal/app/tasks.go:296 with proper tx + sort_order UPDATE; frontend API binding at tasks.ts:21; backend fetches ORDER BY sort_order ASC
  timestamp: 2025-01-09T12:03:00Z

## Evidence

- timestamp: 2025-01-09T12:01:00Z
  checked: TasksView.vue for vuedraggable usage
  found: draggable component IS rendered at lines 495-529 with v-else-if="hasAnyTasks && !taskStore.loading && !treeView" - only in flat non-grouped mode
  implication: Drag-drop code exists but only renders when groupBy is null AND treeView is false

- timestamp: 2025-01-09T12:01:30Z
  checked: Task store groupBy default value
  found: groupBy defaults to 'project' (tasks.ts line 44) - on page load, the grouped rendering branch (line 355) is active, NOT the draggable branch
  implication: User never sees the draggable component unless they manually switch group-by to "None" in FilterBar

- timestamp: 2025-01-09T12:02:00Z
  checked: enhancedFilteredTasks sorting logic (tasks.ts lines 120-146)
  found: Always force-sorts by sortBy value (default 'priority'). No sort-by-sortOrder/manual path exists. After drag reorder updates sortOrder, the computed re-sorts by priority on next tick, snapping the list back.
  implication: Even if user reaches flat mode and drags successfully, the reorder is immediately undone by the computed sort

- timestamp: 2025-01-09T12:02:30Z
  checked: FilterBar sort options (FilterBar.vue lines 173-177)
  found: Only options are Priority, Due Date, Title, Status. No "Manual"/"Custom" option that would sort by sortOrder.
  implication: User has no way to select a sort mode that respects drag-drop ordering

- timestamp: 2025-01-09T12:03:00Z
  checked: Drag handle visibility (TasksView.vue line 508)
  found: drag-handle uses text-muted-foreground/0 (fully transparent) with hover:text-muted-foreground/30 (30% opacity)
  implication: Drag handle is invisible by default, barely visible on hover - users can't discover drag affordance

- timestamp: 2025-01-09T12:04:00Z
  checked: Backend sort_order support
  found: Backend has sort_order column (migrate.go:48), ReorderTasks method (tasks.go:296-311) with proper tx, and fetches tasks ORDER BY sort_order ASC (tasks.go:73). Backend is fully functional.
  implication: Backend is ready; all issues are frontend-only

## Resolution

root_cause: |
  Drag-drop reordering is wired up in TasksView.vue but functionally broken due to 3 compounding frontend issues:
  
  1. PRIMARY: enhancedFilteredTasks (tasks.ts:120-146) always force-sorts by the active sortBy criterion (default: 'priority'), completely ignoring the sortOrder field. When the user drags to reorder, the draggableList setter updates sortOrder and calls reorderTasks, but the computed immediately re-evaluates and re-sorts by priority, snapping the list back. There is no 'manual' sort mode.
  
  2. SECONDARY: groupBy defaults to 'project' (tasks.ts:44). The draggable component only renders in the flat-mode branch (v-else-if on line 495). On initial page load, the grouped rendering branch is active, so the draggable component never mounts. Users must manually switch group-by to "None" first.
  
  3. TERTIARY: The drag handle (TasksView.vue:508) uses text-muted-foreground/0 (fully transparent) and only shows at 30% opacity on hover. Even if the user reaches flat mode, they can't discover the drag affordance.
  
  The backend (sort_order column, ReorderTasks method, ORDER BY sort_order ASC) is fully implemented and correct.
fix:
verification:
files_changed: []
