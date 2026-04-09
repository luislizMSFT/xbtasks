---
status: diagnosed
trigger: "Dashboard still showing old task row views instead of new TreeTaskRow styling"
created: 2025-01-20T12:00:00Z
updated: 2025-01-20T12:00:00Z
---

## Current Focus

hypothesis: CONFIRMED — DashboardView uses DashboardTaskRow (Phase 10 component) instead of TreeTaskRow (Phase 11 component)
test: Confirmed by reading DashboardView.vue imports and template
expecting: n/a — root cause found
next_action: Return diagnosis

## Symptoms

expected: Dashboard view (/dashboard) renders task rows with the new TreeTaskRow styling — status icons, ADO type icons, and ADO metadata badges visible.
actual: User reported "seems like old views"
errors: None reported
reproduction: Test 9 in UAT - navigate to /dashboard and inspect task rows
started: Discovered during UAT

## Eliminated

## Evidence

- timestamp: 2025-01-20T12:01:00Z
  checked: DashboardView.vue line 12 import
  found: "import DashboardTaskRow from '@/components/tasks/DashboardTaskRow.vue'" — uses DashboardTaskRow, NOT TreeTaskRow
  implication: Dashboard renders the Phase 10 DashboardTaskRow component (priority dot, SquareCheckBig, simple status badge) instead of the Phase 11 TreeTaskRow (status icons, ADO type icons, ADO metadata badges, expand/collapse, blocked banners)

- timestamp: 2025-01-20T12:01:30Z
  checked: DashboardView.vue template lines 324 and 347
  found: DashboardTaskRow used in "Today's Focus" (line 324) and "Upcoming" (line 347) sections. Blocked section (lines 370-385) uses inline rendering with just priority dot and title — also lacks TreeTaskRow styling.
  implication: All three task-rendering areas in the dashboard use old-style row rendering

- timestamp: 2025-01-20T12:02:00Z
  checked: DashboardTaskRow.vue vs TreeTaskRow.vue props and rendering
  found: DashboardTaskRow has NO adoMeta prop, no status icon circle, no expand/collapse, no blocked banner — only has priority dot, simple icon, title, personal badge, status badge, optional due date. TreeTaskRow has adoMeta, statusIcon(), adoTypeIcon(), state badge, subtask progress, blocked banner, expand/collapse.
  implication: DashboardTaskRow is a fundamentally different (simpler) component; swapping to TreeTaskRow is needed

- timestamp: 2025-01-20T12:02:30Z
  checked: TasksView.vue imports (for comparison)
  found: TasksView correctly imports and uses TreeTaskRow (6 usages across flat/grouped/tree modes). Phase 11 integration updated TasksView but NOT DashboardView.
  implication: Phase 11 integration was incomplete — it updated TasksView but missed DashboardView

- timestamp: 2025-01-20T12:03:00Z
  checked: STATE.md decisions
  found: "[Phase 10]: DashboardTaskRow is a separate component from TaskRow" and "[Phase 11]: TasksView rewritten: TreeTaskRow replaces TaskRow in all rendering modes". Phase 11 scope explicitly covered TasksView only.
  implication: DashboardView was never in scope for TreeTaskRow integration — the gap is an omission, not a regression

## Resolution

root_cause: DashboardView.vue imports and renders DashboardTaskRow (a Phase 10 component with simplified layout — priority dot, SquareCheckBig icon, basic status badge) instead of the Phase 11 TreeTaskRow component (which has status icon circles, ADO type icons, ADO metadata badges, expand/collapse, subtask progress, blocked banners). Phase 11 integration updated TasksView to use TreeTaskRow but never updated DashboardView. The blocked tasks section in DashboardView also uses inline rendering with old styling. Three areas need updating: (1) import TreeTaskRow instead of DashboardTaskRow, (2) replace DashboardTaskRow usages in "Today's Focus" and "Upcoming" sections with TreeTaskRow, passing adoMeta and other required props, (3) replace inline blocked task rendering with TreeTaskRow.
fix:
verification:
files_changed: []
