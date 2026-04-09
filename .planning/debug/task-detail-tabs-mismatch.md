---
status: diagnosed
trigger: "Task detail pane 3-tab structure does not match the playground UI"
created: 2025-01-01T00:00:00Z
updated: 2025-01-01T00:00:00Z
---

## Current Focus

hypothesis: Real TaskDetail sections (Subtasks, PRs, Notes) exist but have significantly reduced content richness compared to PlaygroundIntegrated.vue — missing subtask filter cycling, personal/ADO differentiation, sync indicators, pipeline status on PRs, and inline notes with count badges
test: Side-by-side comparison of section templates in both files
expecting: Confirm specific missing features in each section
next_action: Return diagnosis

## Symptoms

expected: Clicking a task shows detail pane with 3 tabs: Subtasks, PRs, Notes. Subtasks tab lists child tasks with status icons and add-subtask input. PRs tab shows linked pull requests. Notes tab renders comments section.
actual: User reported: "this does not match the playground UI"
errors: None reported
reproduction: Test 2 in UAT - click any task and inspect the detail pane tabs
started: Discovered during UAT

## Eliminated

(none — root cause found on first hypothesis)

## Evidence

- timestamp: 2025-01-01T00:00:00Z
  checked: Overall layout structure of both TaskDetail.vue and PlaygroundIntegrated.vue
  found: Both use the SAME structural pattern — continuous scrollable sections (NOT actual tabs) with header→sections→footer. The overall layout is structurally aligned.
  implication: The layout skeleton is correct; issue is content richness within sections.

- timestamp: 2025-01-01T00:00:00Z
  checked: Subtasks section (TaskDetail.vue lines 488-540 vs PlaygroundIntegrated.vue lines 1105-1222)
  found: |
    Playground has:
    1. Filter cycling button (All/Mine/ADO/Personal) via cycleSubtaskFilter()
    2. Personal subtasks shown with checkboxes, ADO subtasks with type icons via adoTypeIcon()
    3. Sync status indicators per subtask (pending amber dot, "not pulled" label) with Push/Pull action buttons on hover
    4. Assigned-to badges per subtask with violet styling for other-person
    5. ADO state badges per subtask (e.g. "Active", "New")
    6. ADO IDs displayed per subtask (#50010 etc)
    7. Delete button for personal subtasks only
    Real TaskDetail has:
    1. NO filter cycling — shows all subtasks always
    2. Uses statusIcon() for ALL subtasks — no personal/ADO differentiation
    3. NO sync status indicators
    4. NO assigned-to badges
    5. NO ADO state badges per subtask
    6. NO ADO IDs per subtask
    7. Priority dots only (no delete, no sync actions)
  implication: Subtask section is the BIGGEST visual mismatch — playground shows rich ADO-integrated subtask list while real version shows basic list.

- timestamp: 2025-01-01T00:00:00Z
  checked: Pull Requests section (TaskDetail.vue lines 580-604 vs PlaygroundIntegrated.vue lines 1257-1283)
  found: |
    Playground has pipeline status per PR — shows build/deploy icons with pipeline name under each PR via pipelinesForPr() and pipelineIcon()/pipelineColor().
    Real TaskDetail has NO pipeline status rendering.
    PR title, number, and status badge are present in both.
  implication: PRs section is close but missing pipeline status display.

- timestamp: 2025-01-01T00:00:00Z
  checked: Notes section (TaskDetail.vue lines 606-657 vs PlaygroundIntegrated.vue lines 1285-1321)
  found: |
    Playground has:
    1. Notes count badge
    2. ADO Discussion link button (with AzureDevOpsIcon + ExternalLink icon)
    3. Inline note cards with timestamps, delete-on-hover
    4. Simple Input + Add button for adding notes
    Real TaskDetail has:
    1. NO notes count badge
    2. Notes/ADO chip toggle (different approach to ADO discussion)
    3. Delegates to CommentsSection and ADODiscussion components (not inline)
    4. Uses Textarea for input with meta+enter / ctrl+enter submit
  implication: Notes section uses different component architecture — real version delegates to sub-components while playground renders inline.

- timestamp: 2025-01-01T00:00:00Z
  checked: Description section and header sync indicators
  found: |
    Playground has dirty field amber dots on status/priority badges and "N pending" count in ADO bar.
    Real TaskDetail always shows "Synced" with no dirty tracking.
    Playground uses v-html for rich text description; real uses plain text.
  implication: Sync status feedback is absent from the real detail pane header.

## Resolution

root_cause: |
  The real TaskDetail component (components/tasks/TaskDetail.vue) was not updated to match the richer PlaygroundIntegrated.vue detail pane. While the overall section layout is structurally correct (continuous scrollable sections: Subtasks → Description → PRs → Notes), the content within each section is significantly less featured than the playground version:

  1. SUBTASKS (biggest gap): Missing filter cycling (All/Mine/ADO/Personal), personal vs ADO subtask differentiation (checkboxes vs type icons), sync status indicators, assigned-to badges, ADO state/ID per subtask, and hover actions (Push/Pull)
  2. PRs: Missing pipeline status display per PR
  3. NOTES: Uses CommentsSection/ADODiscussion delegation instead of playground's inline notes with count badge and ADO Discussion link
  4. HEADER: Missing dirty-field amber dots on badges and pending count in ADO bar

fix: (not applied — diagnosis only)
verification: (not applicable)
files_changed: []
