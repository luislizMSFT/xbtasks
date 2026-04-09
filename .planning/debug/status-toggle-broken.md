---
status: diagnosed
trigger: "Status toggle from clicking the status icon on a task row is not working"
created: 2026-04-09T19:00:00Z
updated: 2026-04-09T19:15:00Z
---

## Current Focus

hypothesis: During UAT, TasksView still used old TaskRow (TreeTaskRow swap not yet done) AND groups were collapsed on load, making task rows invisible — user couldn't access status icon to click
test: Compared git diff between UAT commit (3476300) and fix commit (bea5a5a)
expecting: Diff shows component swap and display fix happened AFTER UAT
next_action: Report root cause — already fixed in commit bea5a5a

## Symptoms

expected: Clicking the status icon circle on a task row toggles its status (e.g., todo → done). The icon and styling update immediately.
actual: User reported "false" — feature does not work
errors: None reported
reproduction: Test 10 in UAT — click the status icon circle on any task row
started: Discovered during phase 11 UAT

## Eliminated

- hypothesis: "Status toggle click handler not wired in TreeTaskRow"
  evidence: "TreeTaskRow.vue has onStatusClick() → emit('toggleStatus', task.id) on the button @click. Fully wired."
  timestamp: 2026-04-09T19:02:00Z

- hypothesis: "@toggle-status not handled in TasksView parent"
  evidence: "All 3 rendering paths (grouped, tree, flat/draggable) have @toggle-status handlers calling toggleDone(t)."
  timestamp: 2026-04-09T19:03:00Z

- hypothesis: "taskStore.setStatus or backend SetStatus is broken"
  evidence: "Store setStatus calls tasksApi.setStatus which calls Wails binding SetStatus. Backend Go method exists with correct SQL. Binding IDs match. No TypeScript errors (vue-tsc --noEmit passes clean)."
  timestamp: 2026-04-09T19:05:00Z

- hypothesis: "CSS pointer-events blocking the button click"
  evidence: "Searched all CSS for pointer-events rules. No global rule affects raw <button> elements. shadcn Button's [&_svg]:pointer-events-none only applies to Button component, not raw <button>."
  timestamp: 2026-04-09T19:06:00Z

- hypothesis: "Type mismatch on task.id causing find() to miss"
  evidence: "Task interface has id: number. TreeTaskRow emits props.task.id (number). Handler does taskStore.tasks.find(x => x.id === id) with number comparison. No mismatch."
  timestamp: 2026-04-09T19:07:00Z

## Evidence

- timestamp: 2026-04-09T19:08:00Z
  checked: "Git diff between UAT commit (3476300) and fix commit (bea5a5a)"
  found: "During UAT, TasksView imported OLD TaskRow, not TreeTaskRow. The fix commit swapped TaskRow → TreeTaskRow."
  implication: "Phase 11 integration was incomplete at UAT time — the component swap hadn't been done yet."

- timestamp: 2026-04-09T19:09:00Z
  checked: "groupBy default value at UAT time vs after fix"
  found: "At UAT: groupBy = ref('project'). After fix: groupBy = ref(null). Also, expandedGroups was new Set() with no auto-expand watcher."
  implication: "During UAT, all project groups were collapsed on page load. Task rows were hidden behind collapsed group headers. Users couldn't see task rows to click status icons."

- timestamp: 2026-04-09T19:10:00Z
  checked: "Fix commit bea5a5a contents"
  found: "Fix includes: (1) auto-expand watcher for groups, (2) groupBy default changed to null, (3) TaskRow → TreeTaskRow swap, (4) manual sort mode, (5) filter label unification, (6) TaskDetail reactivity fix, (7) DashboardView adoMeta init."
  implication: "The display bug (collapsed groups) and the component swap were both addressed in the same fix commit. Status toggle should work in current code."

- timestamp: 2026-04-09T19:11:00Z
  checked: "Old TaskRow component's status toggle wiring"
  found: "Old TaskRow also had onCheckClick → emit('toggleStatus', task.id) and the old TasksView had @toggle-status handlers. The toggle mechanism was correct in the old code too."
  implication: "The toggle code itself was never broken — the issue was that task rows were invisible due to collapsed groups."

- timestamp: 2026-04-09T19:12:00Z
  checked: "Current code (post-fix) status toggle wiring"
  found: "TreeTaskRow.vue: button @click=onStatusClick → emit('toggleStatus', task.id). TasksView: @toggle-status handler on all 3 rendering paths → toggleDone → taskStore.setStatus → Wails binding → Go SetStatus. TypeScript compiles clean."
  implication: "Status toggle is correctly implemented in current code and should work."

## Resolution

root_cause: "Two issues combined during UAT: (1) PRIMARY — groupBy defaulted to 'project' with expandedGroups starting empty and no auto-expand watcher, causing all group headers to render collapsed on page load. Task rows were invisible, so the user couldn't find/click the status icon circle. (2) SECONDARY — TasksView was still importing old TaskRow component (phase 11 TreeTaskRow swap not yet done at UAT time). Both issues were fixed in commit bea5a5a."
fix: "Already applied in commit bea5a5a: changed groupBy default to null (flat list visible on load), added auto-expand watcher for groups, swapped TaskRow → TreeTaskRow in all rendering paths."
verification: "Code analysis confirms fix is in place. Status toggle wiring is correct across all rendering modes. TypeScript compiles clean. Needs runtime re-test to confirm."
files_changed: [frontend/src/views/TasksView.vue, frontend/src/stores/tasks.ts]
