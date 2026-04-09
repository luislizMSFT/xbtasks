---
status: diagnosed
trigger: "Group-by mode missing section headers/dividers"
created: 2025-01-01T00:00:00Z
updated: 2025-01-01T00:00:00Z
---

## Current Focus

hypothesis: Group-by rendering has two bugs — (1) all groups start collapsed so tasks are invisible, (2) non-project grouping modes reuse project-centric template that filters out 'No Project' key, hiding all headers when tasks lack projectId
test: Code analysis of TasksView.vue grouped rendering section (lines 354-434) vs task store groupBy logic
expecting: Confirmed by reading code
next_action: Return root cause diagnosis

## Symptoms

expected: Selecting group-by (status, priority, or project) shows tasks in collapsible sections with headers showing group name and task count.
actual: User reported: "no headers divided"
errors: None reported
reproduction: Test 7 in UAT - select group-by mode and check if headers appear between groups
started: Discovered during UAT

## Eliminated

(none)

## Evidence

- timestamp: 2025-01-01T00:00:05Z
  checked: TasksView.vue lines 40-45 — expandedGroups initialization
  found: expandedGroups = ref<Set<string>>(new Set()) — starts EMPTY, no watcher to auto-expand when groupBy changes
  implication: All groups start collapsed; task rows only render when expandedGroups.has(key) (line 403), so NO tasks appear under any group header

- timestamp: 2025-01-01T00:00:06Z
  checked: TasksView.vue line 357 — grouped rendering loop
  found: v-for="key in groupKeys.filter(k => k !== 'No Project')" — filters out 'No Project' key
  implication: For project grouping when tasks have no projectId, ALL keys are 'No Project', so the filter removes ALL keys and NO headers render at all

- timestamp: 2025-01-01T00:00:07Z
  checked: TasksView.vue lines 419-433 — 'No Project' rendering
  found: 'No Project' tasks render flat with NO header — just TreeTaskRow components with no group header above them
  implication: Even for mixed tasks, ungrouped tasks have no visual separation

- timestamp: 2025-01-01T00:00:08Z
  checked: tasks.ts line 44 — groupBy default
  found: groupBy defaults to 'project' — const groupBy = ref<string | null>('project')
  implication: On initial load, view is already in group-by-project mode. If tasks lack projectId, user immediately sees "no headers"

- timestamp: 2025-01-01T00:00:09Z
  checked: TasksView.vue lines 354-434 — overall grouped template structure
  found: Entire template comment says "Project groups (with headers)" — confirms this was designed for project grouping only, not status/priority. Header contains project-specific elements (projectAdoLookup, projectStore.isLinked, selectProject)
  implication: Status and priority group-by modes reuse project-oriented header template. Headers render but with project-specific elements (ADO type icon, linked badge) that show nothing for non-project keys

- timestamp: 2025-01-01T00:00:10Z
  checked: PlaygroundIntegrated.vue — source playground
  found: No group-by rendering exists in playground — only tree view (hierarchical parent/child). Group-by was implemented directly in TasksView without a playground prototype
  implication: Feature was added without iterating on design; explains why it's incomplete

- timestamp: 2025-01-01T00:00:11Z
  checked: tasks.ts lines 152-172 — groupedEnhanced computed
  found: groupedEnhanced correctly groups tasks by status/priority/project. For status: key=t.status, for priority: key=t.priority||'None', for project: key=String(t.projectId)||'No Project'
  implication: Store logic is correct — the issue is purely in the template rendering

## Resolution

root_cause: |
  Two bugs in TasksView.vue grouped rendering (lines 354-434):
  
  1. **Groups always start collapsed**: `expandedGroups` is initialized as empty Set (line 42) with no watcher to auto-populate when `groupBy` changes. Task rows only render when `expandedGroups.has(key)` is true (line 403). Result: user sees headers (maybe) but NO tasks underneath.
  
  2. **'No Project' filtering hides all headers for project grouping**: Line 357 uses `groupKeys.filter(k => k !== 'No Project')` to skip ungrouped tasks. If all tasks lack a projectId (likely during UAT), ALL keys are 'No Project', the filter eliminates everything, and ZERO headers render. Tasks fall through to the flat 'No Project' block (lines 420-433) which has no header.
  
  Combined effect: With groupBy defaulting to 'project' (tasks.ts line 44) and tasks having no projects assigned, the user sees a flat task list with no group headers at all — exactly matching "no headers divided."
  
  Secondary issue: The group header template is project-centric (ADO type icons, project link badges, selectProject click handler) and doesn't adapt for status/priority grouping modes.

fix:
verification:
files_changed: []
