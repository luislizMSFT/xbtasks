---
status: diagnosed
trigger: "Projects show old view in the task list instead of new TreeTaskRow styling"
created: 2025-01-01T00:00:00Z
updated: 2025-01-01T00:00:00Z
---

## Current Focus

hypothesis: CONFIRMED — groupBy defaults to 'project', groups start collapsed, user sees only project headers (no TreeTaskRow)
test: Read task store default + TasksView grouped rendering + expandedGroups init
expecting: groupBy='project' default + expandedGroups=empty causes collapsed groups
next_action: Return root cause diagnosis

## Symptoms

expected: Each task row shows: status icon (colored circle), ADO type icon (if linked), title, state badge, priority badge, and due date. Blocked tasks show a red blocked-reason banner. Tasks with subtasks show a progress indicator.
actual: "wrong no dot projects show old view in the list"
errors: None reported
reproduction: Test 3 in UAT - look at task rows in the task list view
started: Discovered during UAT

## Eliminated

- hypothesis: TasksView still imports/uses old TaskRow instead of TreeTaskRow
  evidence: TasksView.vue imports TreeTaskRow (line 14) and uses it in ALL three rendering modes (grouped line 404/421, tree line 440/458/475, flat line 516). Old TaskRow is NOT imported.
  timestamp: 2025-01-01T00:01:00Z

## Evidence

- timestamp: 2025-01-01T00:01:00Z
  checked: TasksView.vue imports and component usage
  found: TreeTaskRow is imported and used in all rendering paths (grouped, tree, flat/draggable). Old TaskRow.vue is NOT imported.
  implication: The component usage is correct — the problem is elsewhere

- timestamp: 2025-01-01T00:02:00Z
  checked: Task store default groupBy value (stores/tasks.ts line 44)
  found: groupBy defaults to 'project' — `const groupBy = ref<string | null>('project')`
  implication: TasksView always opens in grouped-by-project mode by default

- timestamp: 2025-01-01T00:03:00Z
  checked: TasksView expandedGroups initialization (line 42)
  found: `const expandedGroups = ref<Set<string>>(new Set())` — empty set, all groups start collapsed
  implication: In grouped-by-project mode, ALL project groups are collapsed on load

- timestamp: 2025-01-01T00:04:00Z
  checked: TasksView grouped rendering template (lines 355-434)
  found: Collapsed groups show ONLY project header divs (chevron + type icon + label + ADO badge + count). TreeTaskRow only renders inside `v-if="expandedGroups.has(key)"`. Group headers have NO colored status circle dot.
  implication: Users see project headers without TreeTaskRow styling — "no dot, old view"

- timestamp: 2025-01-01T00:05:00Z
  checked: 'No Project' tasks rendering (TasksView lines 420-433)
  found: Tasks without a project render directly as TreeTaskRow (no group header, no collapse gate). These DO have status dots and new styling.
  implication: Visual split — project tasks hidden behind collapsed headers, non-project tasks shown with new styling. User sees contrast = "projects show old view"

## Resolution

root_cause: |
  groupBy in the task store defaults to 'project' (stores/tasks.ts line 44), so TasksView renders in grouped-by-project mode on initial load. Combined with expandedGroups starting as an empty Set (TasksView.vue line 42), ALL project groups are collapsed by default. Users see only project group header divs (chevron + label + count) which lack TreeTaskRow styling — no colored status circle ("no dot"), no ADO type icon, no state/priority badges, no due dates. Tasks without a project appear below as TreeTaskRow with full new styling, creating a visual split where "projects show old view" while other tasks show new styling with dots.
fix:
verification:
files_changed: []
