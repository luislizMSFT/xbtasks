# Feature Parity Matrix: Playground vs Integrated App

> Generated from exhaustive comparison of `PlaygroundIntegrated.vue` (source of truth)
> against `TasksView.vue`, `TaskDetail.vue`, `TreeTaskRow.vue`, `FilterBar.vue`,
> `FilterCycleButton.vue`, `DashboardView.vue`, and `DashboardTaskRow.vue`.

## Legend

- ✅ = Fully implemented and matching playground
- ⚠️ = Partially implemented / differs from playground
- ❌ = Missing entirely
- N/A = Not applicable (playground-only concept or intentionally different)

---

## 1. Task Row Styling (TreeTaskRow)

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 1.1 | 2-row layout (Row 1: icons/title/badges, Row 2: metadata) | ✅ Present | ✅ TreeTaskRow.vue matches | None |
| 1.2 | Status icon (colored circle via `statusIcon()`) | ✅ Clickable button | ✅ TreeTaskRow has clickable status button | None |
| 1.3 | ADO type icon (via `adoTypeIcon()` when meta has type) | ✅ Shows per-task ADO type | ✅ TreeTaskRow renders `adoTypeIcon` from `adoMeta` prop | None |
| 1.4 | Personal indicator (User icon when no ADO/project) | ✅ `isPersonalTask(task)` check | ✅ TreeTaskRow uses `isPublic` prop (inverted) | None |
| 1.5 | Done tasks: strikethrough + muted text | ✅ `line-through text-muted-foreground` | ✅ TreeTaskRow has same styling via `isDone` | None |
| 1.6 | Parent tasks: font-medium title | ✅ `hasChildren ? 'font-medium' : ''` | ✅ TreeTaskRow has same logic | None |
| 1.7 | State badge (ADO state or task status) | ✅ Shows ADO state or fallback to status | ✅ TreeTaskRow uses `adoStateClasses`/`statusClasses` | None |
| 1.8 | Subtask progress bar (done/total + green bar) | ✅ Shows `subtaskProgress()` with bar | ✅ TreeTaskRow accepts `subtaskProgress` prop | None |
| 1.9 | Row 2: ADO ID number badge | ✅ Shows `#50010` style ADO number | ⚠️ TreeTaskRow shows raw `task.adoId` (not `#`-prefixed) — uses `isPublic && task.adoId` | Playground displays `adoNumber(meta.adoId)` with `#` prefix. TreeTaskRow displays raw `task.adoId` without `#` formatting. |
| 1.10 | Row 2: Personal badge (dashed outline) | ✅ Shows when `isPersonalTask()` | ✅ TreeTaskRow shows when `isPersonal` | None |
| 1.11 | Row 2: Project name | ✅ Shows `projectNames[task.projectId]` | ✅ TreeTaskRow accepts `projectName` prop | None |
| 1.12 | Row 2: Area text | ✅ Shows `task.area` | ✅ TreeTaskRow shows `task.area` | None |
| 1.13 | Row 2: Due date with overdue detection (red) | ✅ Inline overdue logic, raw date display | ✅ TreeTaskRow uses `formatDate()` + overdue/today logic | None |
| 1.14 | Row 2: PriorityBadge component | ✅ `<PriorityBadge :priority>` | ✅ TreeTaskRow includes PriorityBadge | None |
| 1.15 | Row 3: Blocked reason banner (⚠ red text) | ✅ Shows when `blocked && blockedReason` | ✅ TreeTaskRow renders blocked banner | None |
| 1.16 | Selected state: `bg-primary/5 border-l-2 border-l-primary` | ✅ Applied on `selectedId === task.id` | ✅ TreeTaskRow uses same classes via `selected` prop | None |
| 1.17 | Status toggle on click (emit `toggleStatus`) | ✅ Playground: status button exists but no toggle emit | ✅ TreeTaskRow emits `toggleStatus` on status click | TreeTaskRow actually improved over playground here |

---

## 2. Detail Pane

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 2.1 | Permanent panel (no slide-out animation) | ✅ Direct render, no Transition | ✅ No `translate-x-full`, renders directly | None |
| 2.2 | Header Row 1: ADO type icon + editable title + close (X) | ✅ Type icon + `<h2>` title + X button | ⚠️ Type icon + `<Input>` editable title + X button | Integrated uses inline editable Input (improvement), but matches intent |
| 2.3 | Header: Status badge with dirty-field amber dot | ✅ Badge with `statusBgColor` dot + amber dot for dirty | ⚠️ Status Select dropdown (editable), no dirty-field indicator | **GAP**: Playground shows amber dot when field modified locally (dirty tracking). Integrated has no dirty-field indicators at all. |
| 2.4 | Header: Priority badge with dirty-field amber dot | ✅ Badge with priority dot + amber dot for dirty | ⚠️ Priority Select dropdown (editable), no dirty-field indicator | **GAP**: Same — no dirty-field tracking in integrated. |
| 2.5 | Header: Due date display | ✅ CalendarDays icon + date text | ✅ CalendarDays + date input (editable) | Integrated improved: editable date picker |
| 2.6 | Header: Delete button | ✅ In footer | ✅ In header Row 2 (moved for quicker access) | Intentional improvement |
| 2.7 | Project context bar (Folder icon + project name + Switch select + Unlink) | ✅ Full bar with Switch dropdown + Link2Off unlink | ⚠️ Shows Folder + project name only, no Switch dropdown or Unlink | **GAP**: Integrated lacks project Switch and project Unlink functionality from playground. |
| 2.8 | Personal task indicator (no project) | ✅ Dashed border + "Personal task — not linked to a project" | ❌ Not rendered when no project | **GAP**: Playground shows explicit "Personal task" indicator when no project. Integrated shows nothing. |
| 2.9 | ADO integration row (blue bg, ADO icon, #ID, sync status, Push/Open buttons) | ✅ Full row: ADO icon, #ID, Synced/pending status, Push N button, Open button | ✅ ADO row with icon, #ID, Synced text, Push button, Open button, Unlink button | Integrated actually adds Unlink button (improvement) |
| 2.10 | ADO sync: dirty-field count + conditional Push button | ✅ Shows "N pending" with amber dot when dirty, Push button only appears when dirty | ⚠️ Always shows "Synced" text + Push button always visible | **GAP**: Integrated has no concept of dirty fields / pending changes count. Always shows "Synced" regardless. Push button always shown instead of conditionally. |
| 2.11 | Personal task indicator (no ADO) | ✅ Dashed border + "Personal task — not linked to ADO" | ✅ Same dashed border + identical text | None |
| 2.12 | Subtask progress bar in header | ✅ Blue progress bar when subtasks exist | ✅ Same implementation | None |
| 2.13 | **Subtasks section layout** | ✅ Continuous section (not in tab) | ✅ Continuous section (not in tab) | None — both use sections, not tabs |
| 2.14 | Subtask filter (All/Mine/ADO/Personal cycle button) | ✅ 4-way cycle filter in subtask section header | ❌ No subtask filter at all | **GAP**: Playground has rich subtask filtering by source (ADO vs personal vs mine). Integrated has no filter. |
| 2.15 | Subtask: ADO type icon for ADO subtasks | ✅ Shows `adoTypeIcon(st.adoType)` with colored icon | ❌ Subtasks only show `statusIcon()` | **GAP**: Playground distinguishes ADO subtasks (type icon) from personal subtasks (checkbox). Integrated shows status icon for all. |
| 2.16 | Subtask: Personal checkbox toggle | ✅ Clickable checkbox (green when done, border when todo) | ❌ All subtasks use status icon toggle | **GAP**: Playground uses a visual checkbox for personal subtasks. Integrated uses same status icon for all. |
| 2.17 | Subtask: ADO state badge | ✅ Shows `adoStateClasses(st.adoState)` badge | ❌ No ADO state badge on subtasks | **GAP**: Playground shows ADO state (New/Active/Closed) per subtask. Integrated shows only priority dot. |
| 2.18 | Subtask: sync status indicators (synced/pending/not-pulled) | ✅ Amber dot for pending, "not pulled" text, opacity for not-pulled | ❌ No sync status on subtasks | **GAP**: Playground visualizes per-subtask sync status. Integrated has none. |
| 2.19 | Subtask: "Push to ADO" / "Pull from ADO" hover actions | ✅ Shows on hover for pending/not-pulled subtasks | ❌ Not present | **GAP**: Playground shows sync action buttons on hover. |
| 2.20 | Subtask: Assigned-to pill (name, violet for others) | ✅ Shows `st.assignedTo` with violet color for other people | ❌ Not shown | **GAP**: Playground shows who each subtask is assigned to. Integrated doesn't. |
| 2.21 | Subtask: Priority dot | ✅ Colored dot via `priorityDot` map | ✅ Uses `priorityDotBgColor()` | None |
| 2.22 | Subtask: ADO ID number (#xxxx) | ✅ Shows `#st.adoId` for ADO subtasks | ❌ Not shown | **GAP**: Playground shows ADO work item ID on subtasks. |
| 2.23 | Subtask: Delete button (personal only, on hover) | ✅ Trash2 icon, hover-reveal, personal only | ❌ Not present | **GAP**: No way to remove subtasks in integrated. |
| 2.24 | Subtask: "personal" tag badge | ✅ Small tag badge for personal subtasks | ❌ Not shown | **GAP**: No visual distinction between personal and ADO subtasks. |
| 2.25 | Subtask: Add subtask inline input | ✅ Input with Add/Cancel buttons | ✅ Input with Add/Cancel buttons | None |
| 2.26 | **Description section** | ✅ Editable with Edit/Save/Cancel, supports HTML rendering | ✅ Editable with Edit/Save/Cancel, plain text | Minor: playground renders HTML (`v-html`), integrated uses plain text |
| 2.27 | Description: dirty-field amber dot | ✅ Shows amber dot when description modified | ❌ Not present | **GAP**: No dirty tracking for description. |
| 2.28 | **Pull Requests section** | ✅ PR list with icon, title, #number, status badge | ✅ PR list with same structure | None |
| 2.29 | PR: Pipeline status indicators per PR | ✅ Shows CI Build / Integration Tests status icons per PR branch | ❌ Not present | **GAP**: Playground shows pipeline run status (succeeded/failed/in-progress) under each PR. Integrated does not. |
| 2.30 | PR: Branch name display | ✅ Not shown in playground PR list | ⚠️ Integrated uses `branchName()` but only in the store | N/A |
| 2.31 | **Notes section** | ✅ Add input + note cards with timestamp + delete on hover | ⚠️ Notes/ADO tab toggle + CommentsSection / ADODiscussion components | Integrated has richer notes: dual-mode (local notes vs ADO Discussion), shared input, Cmd+Enter submit. **Improvement over playground.** |
| 2.32 | Notes: ADO Discussion link button | ✅ Small "ADO Discussion" link with ExternalLink icon | ✅ Full ADO Discussion tab with reply capability | Integrated improved: full ADO discussion integration |
| 2.33 | **Footer: timestamps** | ✅ "Created ... · Updated ..." | ✅ Same format | None |
| 2.34 | Footer: Delete button | ✅ In footer with Trash2 icon | ❌ Moved to header Row 2 | Intentional design change |
| 2.35 | **Tags section** | ❌ Not in playground | ✅ Editable tags in integrated | Integrated added tags (improvement) |
| 2.36 | **External Links section** | ❌ Not in playground | ✅ ExternalLinks component between description and PRs | Integrated added external links (improvement) |
| 2.37 | **SyncConfirmDialog** | ❌ Not in playground | ✅ Rendered at bottom of detail | Integrated added sync confirm flow |
| 2.38 | **ConflictResolver** | ❌ Not in playground | ✅ Rendered at bottom of detail | Integrated added conflict resolution |

---

## 3. Toolbar / Header Bar

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 3.1 | Tree/Flat view toggle buttons | ✅ Two buttons: Tree (GitBranch) / Flat (List) | ✅ Tree/flat toggle in FilterBar (Network/List icon) | Presentation differs but functionality matches |
| 3.2 | Show/Hide Detail toggle | ✅ Button toggles `showDetail` | ❌ Detail panel always visible | **GAP**: Playground allows hiding the detail pane entirely. Integrated always shows detail panel. |
| 3.3 | Expand All / Collapse All buttons (tree mode) | ✅ Two buttons in header, tree mode only | ❌ Not present | **GAP**: No expand/collapse all in integrated. |
| 3.4 | Task count display | ✅ Shows "N roots · M tasks" (tree) or "N tasks" (flat) | ❌ Not displayed | **GAP**: No task count shown in integrated toolbar. |
| 3.5 | Filter cycle button (All/ADO/Personal) | ✅ Inline cycle button in list header | ✅ FilterCycleButton in topbar Teleport | None — same logic, different placement |
| 3.6 | Quick-add "Add" button in list header | ✅ Plus icon + "Add" in list header | ✅ QuickAddInput always visible at top of list | QuickAddInput is always visible (improvement over toggle) |

---

## 4. Filters (FilterBar)

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 4.1 | Status chip tabs (All/Todo/In Progress/etc.) | ❌ Not in playground (playground has simple cycle filter) | ✅ 6 status chip Badges in Teleport topbar | Integrated improvement |
| 4.2 | Priority filter dropdown | ❌ Not in playground | ✅ FilterBar: Priority select (All/P0-P3) | Integrated improvement |
| 4.3 | Project filter dropdown | ❌ Not in playground | ✅ FilterBar: Project select | Integrated improvement |
| 4.4 | Due date filter dropdown | ❌ Not in playground | ✅ FilterBar: Due select (Overdue/Today/Week/No Date) | Integrated improvement |
| 4.5 | Scope (ADO link) filter dropdown | ✅ Cycle button only | ✅ Both cycle button + FilterBar Scope select | Integrated improvement |
| 4.6 | Sort-by dropdown | ❌ Not in playground | ✅ FilterBar: Sort (Priority/Due/Title/Status) | Integrated improvement |
| 4.7 | Group-by dropdown | ❌ Not in playground | ✅ FilterBar: Group (None/Status/Priority/Project) | Integrated improvement |
| 4.8 | Active filter count badge | ❌ Not in playground | ✅ FilterBar shows count when filters active | Integrated improvement |

---

## 5. List Rendering Modes

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 5.1 | Tree view: Root → Child → Grandchild (3 levels) | ✅ Inline rendering with 3 indent levels | ✅ TreeTaskRow with `indentLevel` 0/1/2 | None |
| 5.2 | Tree view: Expand/collapse per node | ✅ `expandedNodes` Set, toggle on chevron click | ✅ `expandedTreeNodes` Set, same behavior | None |
| 5.3 | Flat view: All tasks in flat list | ✅ `filteredAllTasks` shown without hierarchy | ✅ Uses `enhancedFilteredTasks` from store | None |
| 5.4 | Flat view: ADO badges inline | ✅ ADO number badge + Personal badge in flat rows | ⚠️ Flat mode uses same TreeTaskRow (2-row layout, not flat-style) | Minor: Playground has a simpler single-row flat view with inline badges. Integrated reuses the 2-row TreeTaskRow layout in flat mode too. This is acceptable. |
| 5.5 | Grouped view: Collapsible group headers | ❌ Not in playground | ✅ Group headers with chevron, type icon, ADO badge, count | Integrated improvement |

---

## 6. Drag-Drop Reordering

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 6.1 | Drag handle (dot grid SVG) | ❌ Not in playground | ✅ In flat mode via vuedraggable with ghost class | Integrated improvement |
| 6.2 | Drag-drop reorder in flat mode | ❌ Not in playground | ✅ vuedraggable with `@start/@end`, persists to backend | Integrated improvement |
| 6.3 | Drag disabled in tree/group modes | N/A | ✅ Only renders draggable in non-grouped, non-tree mode | Correct behavior |

---

## 7. Status Toggle from Row

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 7.1 | Click status icon to toggle done/todo | ⚠️ Status button exists but no toggle behavior | ✅ TreeTaskRow emits `toggleStatus`, TasksView calls `toggleDone()` | Integrated improved: actually toggles status |

---

## 8. ADO Integration Indicators

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 8.1 | Per-task ADO type icon in list | ✅ Via `adoMeta` reactive object | ✅ Via `useAdoMeta` composable + batch cache | None |
| 8.2 | Per-task ADO state badge in list | ✅ Via `adoMeta` reactive object | ✅ Via `useAdoMeta` composable | None |
| 8.3 | ADO ID number in list row 2 | ✅ `#50010` prefixed display | ⚠️ Shows raw `task.adoId` without `#` prefix | **GAP**: TreeTaskRow displays `task.adoId` raw. Should use `#` prefix formatting. |
| 8.4 | Detail pane: ADO sync status (synced vs pending count) | ✅ Shows "Synced" or "N pending" with amber styling | ⚠️ Always shows "Synced" | **GAP**: No dirty-field awareness in integrated detail. |
| 8.5 | Detail pane: conditional Push button (only when dirty) | ✅ Push button only appears when `dirtyFields.length > 0` | ⚠️ Push button always visible | **GAP**: Push button always shown regardless of sync state. |
| 8.6 | Detail pane: Open in ADO link | ✅ Opens external link | ✅ Opens via `openAdoLink()` with URL built from cached work item | None |
| 8.7 | Detail pane: Unlink from ADO | ❌ Not in playground | ✅ Unlink button with `unlinkFromADO()` | Integrated improvement |

---

## 9. Dashboard

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 9.1 | Dashboard uses TreeTaskRow for task rows | N/A (no dashboard in playground) | ❌ Uses `DashboardTaskRow` (old component) | **GAP**: DashboardView still imports `DashboardTaskRow.vue` — NOT `TreeTaskRow`. Summary said it was replaced (commit `23ee5f5`), but current code still uses the old component. |
| 9.2 | Dashboard task rows show ADO metadata | N/A | ❌ `DashboardTaskRow` uses `task.category` for ADO type (not `useAdoMeta`) | **GAP**: DashboardTaskRow doesn't use the `useAdoMeta` composable or batch cache. Falls back to `task.category` which may be empty. |
| 9.3 | Dashboard task rows: status icon | N/A | ❌ Uses priority dot (colored circle) instead of status icon | **GAP**: DashboardTaskRow shows a priority dot, not a status icon. Doesn't match TreeTaskRow styling. |
| 9.4 | Dashboard: Blocked section styling | N/A | ⚠️ Inline rendering (not DashboardTaskRow), shows priority dot + title + blocked reason | Different from TreeTaskRow but acceptable for dashboard context |

---

## 10. Detail Pane: Tab vs Section Structure

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 10.1 | Overall layout | ✅ Continuous scrollable sections (Subtasks → Description → PRs → Notes) | ✅ Continuous scrollable sections (same order) | None — both use sections, not tabs. UAT #2 reported "3 tabs" but both implementations use continuous sections. |
| 10.2 | Section order | ✅ Subtasks → Description → PRs → Notes | ✅ Subtasks → Description → External Links → PRs → Notes | Integrated adds External Links section (improvement) |

---

## 11. Quick-Add Task

| # | Feature | Playground | Integrated | Gap Description |
|---|---------|-----------|------------|-----------------|
| 11.1 | Quick-add input with Enter/Escape | ✅ Toggle-reveal input at bottom of list | ✅ Always-visible QuickAddInput at top of list | Integrated improved: always visible, uses extracted component |
| 11.2 | Auto-focus on add | ✅ `nextTick(() => ref.focus())` | ✅ QuickAddInput handles auto-focus | None |
| 11.3 | Add/Cancel buttons | ✅ Both buttons present | ✅ Both buttons in QuickAddInput | None |
| 11.4 | After add: auto-select new task + show detail | ✅ Sets `selectedId` and `showDetail = true` | ⚠️ Calls `taskStore.quickAdd()` — selection depends on store impl | Minor: playground auto-selects the new task; integrated delegates to store |

---

## Critical Gaps Summary (Fixes Needed)

| Priority | Gap | Impact | Files to Fix |
|----------|-----|--------|-------------|
| **P0** | Dashboard still uses old `DashboardTaskRow` instead of `TreeTaskRow` | UAT #9 fails — "old views" | `DashboardView.vue` |
| **P0** | Detail pane: No dirty-field tracking / sync status indicators | UAT #2 partially — playground shows pending sync status | `TaskDetail.vue` |
| **P0** | Subtask section: Missing ADO type icons, ADO state badges, sync status, assigned-to, personal badge, source filter | UAT #2 — "does not match playground" | `TaskDetail.vue` |
| **P1** | Detail: No subtask filter cycle (All/Mine/ADO/Personal) | Functionality gap vs playground | `TaskDetail.vue` |
| **P1** | Detail: No subtask delete button | Functionality gap | `TaskDetail.vue` |
| **P1** | Detail: No pipeline status under PRs | Visual gap vs playground | `TaskDetail.vue` |
| **P1** | Detail: No project Switch dropdown or project Unlink in project context bar | Functionality gap | `TaskDetail.vue` |
| **P1** | Detail: No personal-task-no-project indicator | Visual gap | `TaskDetail.vue` |
| **P2** | TreeTaskRow: ADO ID not `#`-prefixed in Row 2 | Minor visual gap | `TreeTaskRow.vue` |
| **P2** | No Show/Hide detail pane toggle | Functionality gap | `TasksView.vue` |
| **P2** | No Expand All / Collapse All buttons for tree mode | Functionality gap | `TasksView.vue` |
| **P2** | No task count display in toolbar | Visual gap | `TasksView.vue` |

---

## Features Where Integrated IMPROVES on Playground

| Feature | Description |
|---------|-------------|
| Editable title in detail header | Input field instead of static h2 |
| Full FilterBar with dropdowns | Priority, Project, Due, Scope, Sort, Group filters |
| Status chip tabs | 6-status quick filter in topbar |
| Drag-drop reorder | vuedraggable in flat mode |
| Group-by with collapsible headers | Status/Priority/Project grouping |
| Notes: dual-mode (Notes + ADO Discussion) | Two-tab notes with ADO reply capability |
| External Links section | Dedicated section in detail |
| SyncConfirmDialog + ConflictResolver | Real sync confirmation and conflict resolution |
| Tags editing | Add/remove tags in detail |
| Status toggle from row | Actually works (playground button had no toggle logic) |
| Unlink from ADO | Button in ADO integration row |
| Real backend persistence | Store-backed CRUD vs playground mock data |
