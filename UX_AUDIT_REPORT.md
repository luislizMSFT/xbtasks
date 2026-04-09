# UX Integration Audit Report
## Playground → Production Inconsistencies

**Date:** 2025-07-16  
**Auditor:** Senior UX Engineer  
**Reference:** `PlaygroundIntegrated.vue` (1354 lines)  
**Targets:** `TasksView.vue`, `TaskDetail.vue`, `DashboardView.vue`, `FilterBar.vue`, `FilterCycleButton.vue`, task store

---

## Critical (P0) — Broken Interactions

### 1. Status Toggle From Row Click Not Working
**Location:** `TreeTaskRow.vue:103-108`, `TasksView.vue:423`  
**Issue:** The status icon button in TreeTaskRow correctly emits `toggleStatus` with `e.stopPropagation()`, and TasksView handles the event — but the row's outer `@click="emit('click')"` fires first because the status button doesn't prevent the row click from selecting the task. The real problem: `toggleDone()` calls `taskStore.setStatus()` which replaces the task object in the store, causing a re-render race with the selection change happening simultaneously.  
**Playground:** Status icon click directly mutates `task.status` inline (no async, no store round-trip).  
**Fix:** In `toggleDone()`, use optimistic local update before the API call. Also verify `@click.stop` is on the status `<button>` in TreeTaskRow (it is — via `onStatusClick` which calls `e.stopPropagation()`). The likely root cause is the `async setStatus()` re-fetch that replaces the reactive task array mid-render. Apply optimistic update:
```ts
async function toggleDone(task: Task) {
  const newStatus = task.status === 'done' ? 'todo' : 'done'
  const original = task.status
  task.status = newStatus // optimistic
  try { await taskStore.setStatus(task.id, newStatus) }
  catch { task.status = original } // rollback
}
```

### 2. Drag-Drop Non-Functional (3 compounding bugs)
**Location:** `TasksView.vue:263-279, 504-538`, task store `sortBy`/`groupBy`  

**Bug 2a — Default `groupBy` prevents draggable rendering:**  
`groupBy` defaults to `null` (line 44 of store), so flat draggable *does* render when `!treeView && !groupBy`. This is actually correct now. But if a user ever sets groupBy and then clears it, the draggable path activates. ✅ Not a blocker.

**Bug 2b — Sort snaps back:**  
`draggableList` is a computed writable that sets `sortBy = 'manual'` in the setter. But `enhancedFilteredTasks` re-sorts on every reactive tick. When the setter mutates `taskStore.tasks[idx].sortOrder`, Vue reactivity triggers `enhancedFilteredTasks` to recompute *before* the manual sort mode is fully applied, causing a flash-snap.  
**Fix:** Set `sortBy = 'manual'` *before* mutating sortOrder, and use `nextTick` to batch:
```ts
set: (val: Task[]) => {
  taskStore.sortBy = 'manual'
  nextTick(() => {
    const ids = val.map(t => t.id)
    for (let i = 0; i < val.length; i++) {
      const idx = taskStore.tasks.findIndex(t => t.id === val[i].id)
      if (idx !== -1) taskStore.tasks[idx].sortOrder = i
    }
    taskStore.reorderTasks(ids)
  })
}
```

**Bug 2c — Invisible drag handle:**  
`TasksView.vue:517` — the drag handle has `text-muted-foreground/0` (fully transparent) and only shows on `hover:text-muted-foreground/30` (30% opacity). This makes the grip dots nearly invisible even on hover.  
**Fix:** Change to `text-muted-foreground/20 hover:text-muted-foreground/60` for discoverability.

### 3. Group-By Headers Start Collapsed (groups invisible on first load)
**Location:** `TasksView.vue:42, 166-170`  
**Issue:** `expandedGroups` starts as `new Set()`. There is a watcher at line 166-170 that auto-expands all groups when `groupBy` changes or `groupKeys` updates — this was added as a fix. However, the `{ immediate: true }` fires before `groupKeys` is populated (it depends on `groupedEnhanced` which is async), so on initial load with `groupBy` set, groups may still appear collapsed until the next reactive tick.  
**Playground:** No group-by mode exists — playground only has tree/flat. This is a feature gap.  
**Fix:** The watcher fix is correct in principle. Ensure `expandedGroups` is initialized to include all keys from the first `groupKeys` computation. Add a fallback in the template: if `groupBy` is set but `expandedGroups` is empty, render tasks un-collapsed:
```vue
<template v-if="expandedGroups.has(key) || expandedGroups.size === 0">
```

---

## High (P1) — Feature Gaps vs Playground

### 4. Task Detail Missing: Subtask Filters
**Location:** `TaskDetail.vue:488-540` vs Playground lines `614-633, 1106-1222`  
**Issue:** TaskDetail shows subtasks as a flat unfiltered list. Playground has a subtask filter cycle button (All → Mine → ADO → Personal) with `cycleSubtaskFilter()`, filtered display, and source-aware rendering (personal checkbox vs ADO type icon).  
**What's missing:**
- No `subtaskFilter` ref or `filteredSubtasks` computed
- No cycle button in subtask section header
- No visual distinction between personal/ADO subtasks
- No sync status indicators (pending/not-pulled badges)
- No assigned-to pill on subtasks
- No inline "Push to ADO" / "Pull from ADO" hover actions

**Fix:** Port the subtask filter system from playground:
```ts
const subtaskFilter = ref<'all'|'ado'|'personal'|'mine'>('all')
// Add filteredSubtasks computed, cycleSubtaskFilter()
// Render cycle button next to "Add" in subtask header
// Add source-aware icons/badges per subtask row
```

### 5. Task Detail Missing: Pipeline Status on PRs
**Location:** `TaskDetail.vue:580-604` vs Playground lines `1257-1283`  
**Issue:** TaskDetail shows PRs with icon + title + number + status badge, but no pipeline status underneath. Playground shows nested pipeline runs per PR (CI Build ✓, Integration Tests ✗) with color-coded status icons.  
**Fix:** Look up pipelines by PR source branch (from `adoStore.pipelines`) and render nested under each PR row:
```vue
<div v-if="pipelinesForPr(pr.sourceBranch).length" class="flex items-center gap-3 pl-6 pb-1.5">
  <div v-for="pipe in pipelinesForPr(pr.sourceBranch)" class="flex items-center gap-1 text-[9px]">
    <component :is="pipelineIcon(pipe)" :size="10" :class="pipelineColor(pipe)" />
    <span class="text-muted-foreground">{{ pipe.name }}</span>
  </div>
</div>
```

### 6. Task Detail Missing: Inline Notes (playground style)
**Location:** `TaskDetail.vue:606-657` vs Playground lines `1285-1321`  
**Issue:** TaskDetail uses `CommentsSection`/`ADODiscussion` components with tab toggle — this is *more* functional than the playground. However, the playground's inline note cards (with timestamp, hover-delete) are visually richer. The integrated version delegates entirely to sub-components, losing the card-style note presentation.  
**Recommendation:** This is acceptable divergence — the integrated version is functionally superior. Low-priority cosmetic alignment only. Keep the current tab-based approach but ensure CommentsSection renders notes with the card styling from playground (bg-muted/40, border, hover-to-delete).

### 7. Task Detail Missing: Sync Indicators
**Location:** `TaskDetail.vue:451-474` vs Playground lines `1064-1091`  
**Issue:** The integrated ADO bar shows a static "Synced" label regardless of actual dirty state. Playground shows dynamic sync status:
- `dirtyFields.length === 0` → "✓ Synced" (green)
- `dirtyFields.length > 0` → "● N pending" (amber) + "Push N" button

**Fix:** Track dirty fields by comparing `editTitle/editStatus/etc` against `task.value` originals:
```ts
const dirtyFields = computed(() => {
  if (!task.value) return []
  const dirty: string[] = []
  if (editTitle.value !== task.value.title) dirty.push('title')
  if (editStatus.value !== task.value.status) dirty.push('status')
  // ...etc
  return dirty
})
```
Then conditionally render "Synced" vs "N pending" + Push button.

### 8. Dashboard Uses DashboardTaskRow Instead of TreeTaskRow
**Location:** `DashboardView.vue:326-333, 349-356` — uses `DashboardTaskRow`  
**Issue:** Dashboard renders task rows with the old `DashboardTaskRow` component (simpler: priority dot + icon + title + status). The playground reference and TasksView both use `TreeTaskRow` which has richer rendering (expand chevron, ADO type icon, state badge, subtask progress bar, area, due date, blocked banner).  
**Impact:** Visual inconsistency between dashboard task cards and main task list. Clicking a dashboard task → navigating to task list shows a completely different row design.  
**Fix:** Replace `DashboardTaskRow` with `TreeTaskRow` in dashboard sections:
```vue
<TreeTaskRow
  :task="task"
  :is-public="taskStore.isPublic(task.id)"
  :ado-meta="adoMeta.getAdoMeta(task.id)"
  :project-name="projectNameMap[task.projectId]"
  :subtask-progress="subtaskProgress(task.id)"
  @click="selectTask(task.id)"
  @toggle-status="(id) => toggleDone(id)"
/>
```
Note: TreeTaskRow may need a `compact` prop for dashboard contexts (no expand chevron, tighter padding).

---

## Medium (P2) — Duplicate/Conflicting Controls

### 9. FilterCycleButton Duplicates FilterBar Scope Dropdown
**Location:** `TasksView.vue:308` (topbar) + `FilterBar.vue:139-152` (scope dropdown)  
**Issue:** Two controls modify `taskStore.filterAdoLink`:
1. **FilterCycleButton** in the topbar — cycles All → ADO → Personal (labels: All/ADO/Personal)
2. **FilterBar "Scope" dropdown** — same values (All/ADO/Personal) via `<Select>`

Both are wired to the same store value, so they stay in sync. But having two controls is confusing — users don't know which one is the "real" filter.  
**Playground:** Has only ONE filter control — the cycle button in the list header (lines 715-723). No scope dropdown.  
**Fix:** Remove the "Scope" dropdown from FilterBar since the topbar FilterCycleButton already handles it. This eliminates the duplication:
```diff
- <span class="text-[10px] text-muted-foreground shrink-0">Scope</span>
- <Select :model-value="filterAdoLink" ...>
-   ...
- </Select>
```
Also remove `filterAdoLink` from FilterBar's props/emits since it's controlled elsewhere.

### 10. Missing List Header (task count + filter + add button)
**Location:** `TasksView.vue` vs Playground lines `713-727`  
**Issue:** Playground has a dedicated list header row showing: task count label ("7 roots · 28 tasks"), the filter cycle button, and an "Add" button. TasksView skips this — the filter cycle button is teleported to the topbar, and the count is nowhere.  
**User feedback:** "the preview task list header that allows you to toggle personal/ado tasks should be in the new version"  
**Fix:** Add a list header row above the ScrollArea in TasksView:
```vue
<div class="px-3 py-2 border-b border-border/50 flex items-center gap-2">
  <span class="text-xs font-semibold text-muted-foreground flex-1">
    {{ taskStore.enhancedFilteredTasks.length }} tasks
  </span>
  <FilterCycleButton v-model="taskStore.filterAdoLink" />
  <button class="text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1 px-1.5 py-0.5 rounded hover:bg-muted"
    @click="showQuickAdd = true">
    <Plus :size="11" /> Add
  </button>
</div>
```

### 11. Projects in Task List Use Old View
**Location:** `TasksView.vue:544-549` — right panel uses `ProjectDetail` component  
**User feedback:** "for projects in the task list it's using an old view but should have the new full view"  
**Issue:** When a project group header is clicked in TasksView, it opens `ProjectDetail` in the right panel. This is a separate component that likely uses the old task row rendering internally. Users expect the project view to show TreeTaskRow-based task rendering with the full new visual treatment.  
**Fix:** Ensure `ProjectDetail` internally uses `TreeTaskRow` for its task list. Audit `ProjectDetail.vue` to confirm it renders tasks consistently. If it uses old row components, swap them for TreeTaskRow.

---

## Low (P3) — Cosmetic & Polish

### 12. FilterCycleButton Label Mismatch  
**Location:** `FilterCycleButton.vue:14-18` uses "ADO" for linked state  
**FilterBar.vue:148-149** uses "ADO" for linked, "Personal" for personal  
**Playground:622** subtask filter uses "ADO"/"Personal"/"Mine"  
**Status:** Labels are consistent now. ✅ No action needed.

### 13. Tree View Not Default (Playground defaults to tree)
**Location:** `TasksView.vue:44` — `treeView = ref(false)` (flat by default)  
**Playground:** `listViewMode = ref<'tree' | 'flat'>('tree')` (tree by default)  
**Fix:** Change default to `true` to match playground behavior:
```ts
const treeView = ref(true)
```

### 14. Tree Nodes Not Pre-Expanded
**Location:** `TasksView.vue:45` — `expandedTreeNodes = new Set()` (all collapsed)  
**Playground:** `expandedNodes = new Set([100, 101, 200, 400])` (key nodes pre-expanded)  
**Impact:** Users see a flat list of root tasks with no children visible on first load.  
**Fix:** Auto-expand root tasks that have children on initial data load:
```ts
watch(() => taskStore.tasks, (tasks) => {
  if (expandedTreeNodes.value.size === 0 && tasks.length > 0) {
    const roots = tasks.filter(t => !t.parentId)
    for (const r of roots) {
      if (tasks.some(t => t.parentId === r.id)) expandedTreeNodes.value.add(r.id)
    }
  }
}, { immediate: true })
```

### 15. No Expand All / Collapse All Controls
**Location:** TasksView has no equivalent of playground's "Expand All" / "Collapse All" buttons (playground line 704-706).  
**Fix:** Add toggle buttons to the FilterBar or inline above the tree when `treeView` is active.

### 16. DashboardTaskRow Missing Information Density
**Location:** `DashboardTaskRow.vue` vs `TreeTaskRow.vue`  
**Issue:** DashboardTaskRow shows: priority dot, type icon, title, personal badge, status badge, due date. It omits: area, project name, ADO ID, subtask progress, blocked banner, expand chevron. Even if TreeTaskRow is adopted (P1 #8), dashboard may need a compact variant.  
**Fix:** Add a `compact` prop to TreeTaskRow that hides Row 2 metadata and Row 3 blocked banner, tightening vertical space for dashboard contexts.

### 17. Detail Pane Description: Prose vs Plain Text
**Location:** `TaskDetail.vue:566-572` renders `editDescription` as plain text `<p>`. Playground (line 1248) renders `v-html="detail.description"` with prose styling.  
**Issue:** Playground supports HTML description (from ADO). Integrated version strips HTML on load (line 104) and renders plain text.  
**Impact:** Minor — rich text from ADO loses formatting in the detail pane.  
**Fix:** Preserve original HTML in a separate ref, render with `v-html` and prose classes, edit as plain text.

---

## Summary Priority Matrix

| # | Issue | Severity | Effort | User Impact |
|---|-------|----------|--------|-------------|
| 1 | Status toggle re-render race | P0 | S | Broken interaction |
| 2 | Drag-drop (3 bugs) | P0 | M | Core feature broken |
| 3 | Group headers collapsed | P0 | S | Groups appear empty |
| 4 | Missing subtask filters | P1 | M | Lost feature parity |
| 5 | Missing pipeline status on PRs | P1 | M | Missing context |
| 7 | Static sync indicators | P1 | S | Misleading status |
| 8 | Dashboard old row component | P1 | M | Visual inconsistency |
| 9 | Duplicate scope controls | P2 | S | UX confusion |
| 10 | Missing list header | P2 | S | Lost discoverability |
| 11 | Projects old view | P2 | M | Stale rendering |
| 13 | Tree not default | P3 | XS | Preference mismatch |
| 14 | Tree not pre-expanded | P3 | S | Poor first impression |
| 15 | No expand/collapse all | P3 | S | Missing convenience |
| 16 | Dashboard row info density | P3 | M | Visual gap |
| 17 | Description HTML stripping | P3 | S | Formatting loss |

**Effort:** XS = <30min, S = 1-2hr, M = 2-4hr

---

## Recommended Fix Order

1. **Sprint 1 (P0):** Issues 1, 2, 3 — unblock core interactions
2. **Sprint 2 (P1):** Issues 4, 5, 7, 8 — restore feature parity with playground
3. **Sprint 3 (P2):** Issues 9, 10, 11 — clean up control duplication
4. **Backlog (P3):** Issues 13-17 — polish and defaults
