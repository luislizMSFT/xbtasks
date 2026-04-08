# Phase 11: Integrate new task list into real app - Context

**Gathered:** 2026-04-08
**Status:** Ready for planning

<domain>
## Phase Boundary

Replace the current TasksView with the playground's integrated tree+detail layout. Extract shared components from PlaygroundIntegrated.vue, wire them to real stores (taskStore, adoStore, syncStore), and replace the existing TaskRow and TaskDetail components with new unified versions. Add a backend ADO metadata cache so task rows can display ADO type icons and state badges without N+1 lookups. This is a frontend-heavy phase with one backend addition (ado_meta_cache table).

</domain>

<decisions>
## Implementation Decisions

### Replacement Strategy
- Replace TasksView.vue entirely — rewrite to use the playground's tree+detail 2-panel layout
- Preserve ALL existing features: drag-drop reordering, group-by (status/priority/project), full FilterBar (status, priority, project, dueDate, ADO link), alongside new playground features (tree/flat toggle, All/ADO/Personal filter)
- Permanent 2-panel split layout: task list always on left, detail pane always on right when a task is selected (not a slide-out — permanent split)
- Keep existing store wiring (`enhancedFilteredTasks`, `groupedEnhanced`, `filterAdoLink`) — these already work

### Task Row Component
- Create new `TreeTaskRow.vue` from playground pattern — supports indent levels, ADO type icons, state badges, expand/collapse chevrons built-in
- TreeTaskRow replaces TaskRow everywhere: TasksView, DashboardView (Today's Focus, Upcoming, Blocked), and any other list context — single source of truth
- Simplified row content: status icon + title + ADO type icon + priority badge + due date
- Drop from current TaskRow: tags, description preview, time-ago timestamp (too noisy)
- Keep: status circle (Things 3 style), ADO badge (hollow/filled), blocked reason subtitle, project name caption

### Detail Pane
- Build new TaskDetail.vue from playground pattern — replaces current TaskDetail.vue completely
- 3-tab structure: Subtasks, PRs, Notes (simpler than current Work/Discussion tabs)
- Include sync features: sync confirmation dialog, conflict resolver, push-to-ADO actions
- Include external links section (reuse ExternalLinks.vue component)
- Reuse CommentsSection.vue for Notes tab (comments with push-to-ADO capability)
- Detail pane is the right panel in the permanent 2-panel split

### ADO Metadata Enrichment
- Create `ado_meta_cache` table in SQLite: `(task_id INTEGER, ado_type TEXT, ado_state TEXT)` — populated during sync, keeps tasks table clean
- Create `useAdoMeta()` composable on frontend — provides `getAdoMeta(taskId): { type, state } | null`
- Batch load all ADO metadata on app startup + refresh on `sync:completed` Wails event
- Backend: new `ADOMetaCacheService` with `GetAll() → map[int]{type, state}` and `Refresh()` called during sync
- Frontend composable first checks the batch-loaded cache (fast), no per-task API calls

### Shared Components to Extract
- `TreeTaskRow.vue` — task row with indent, ADO type icon, state badge, expand/collapse
- `QuickAddInput.vue` — inline quick-add pattern (Enter to add, Esc to cancel) used for tasks + subtasks
- `FilterCycleButton.vue` — filter cycle UI (click to cycle through All → ADO → Personal)

### Style Helpers to Add
- `statusIcon()` in `styles.ts` — maps status → Lucide icon component (currently inline in playground and TaskRow)
- `priorityDotBgColor()` in `styles.ts` — maps P0-P3 → Tailwind bg class

### Claude's Discretion
- Exact component prop API for TreeTaskRow (which props are required vs optional)
- How to handle the transition animation when selecting/deselecting tasks in the split pane
- Whether QuickAddInput needs focus management (auto-focus on mount vs manual)
- Internal state management for expand/collapse in tree view (local ref vs store)
- How to wire drag-drop in tree view (existing Sortable.js or new approach)
- Whether to keep the existing FilterBar component or replace with a simpler toolbar

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Playground Prototype (approved design)
- `frontend/src/views/playground/PlaygroundIntegrated.vue` — The source-of-truth prototype: tree+detail 2-panel layout, mock data shapes, flat/tree toggle, filter cycle, quick-add, status icons, ADO type icons, subtask list with source tracking. Extract components and patterns from here.

### Current Production Code (to be replaced)
- `frontend/src/views/TasksView.vue` — Current task view (23KB). Has FilterBar, drag-drop, group-by, TaskRow usage, TaskDetail slide-out. Must preserve all filter/sort/group features.
- `frontend/src/components/tasks/TaskRow.vue` — Current task row (183 lines). Being replaced by TreeTaskRow but review for features to carry forward.
- `frontend/src/components/tasks/TaskDetail.vue` — Current detail pane. Being replaced but has sync/conflict/comments/links logic to understand.

### Components to Reuse (not replace)
- `frontend/src/components/tasks/CommentsSection.vue` — Comments with push-to-ADO. Reuse in Notes tab.
- `frontend/src/components/tasks/ExternalLinks.vue` — External link management. Reuse in detail pane.
- `frontend/src/components/ado/SyncConfirmDialog.vue` — Sync confirmation with diff preview. Reuse.
- `frontend/src/components/ado/ConflictResolver.vue` — Per-field conflict resolution. Reuse.

### Stores (data sources)
- `frontend/src/stores/tasks.ts` — Task store: `enhancedFilteredTasks`, `groupedEnhanced`, `filterAdoLink`, `getSubtasks()`, `quickAdd()`, `isPublic()`
- `frontend/src/stores/ado.ts` — ADO store: `workItemTree`, ADO work item data for metadata cross-reference
- `frontend/src/stores/sync.ts` — Sync store: events, `generateOutboundDiff()`, conflict tracking

### Style & Helpers
- `frontend/src/lib/styles.ts` — `statusColor`, `adoTypeIcon`, `adoTypeColor`, `adoStateClasses`, `prStatusClasses`. Add `statusIcon()` and `priorityDotBgColor()` here.
- `frontend/src/components/ui/PriorityBadge.vue` — Priority badge component (keep using)
- `frontend/src/components/icons/AzureDevOpsIcon.vue` — ADO icon (keep using)

### Backend (for ado_meta_cache)
- `internal/db/db.go` — SQLite schema. Add `ado_meta_cache` table here.
- `internal/app/syncservice.go` — Sync service. Populate `ado_meta_cache` during sync.
- `internal/app/adoservice.go` — ADO service. Source of work item type/state data.

### Prior Phase Context
- `.planning/phases/02-ado-integration-prs-unified-dashboard/02-CONTEXT.md` — Personal→public model, sync behavior, ADO linking flows
- `.planning/phases/10-implement-dashboard-redesign-and-unified-header-bars/10-CONTEXT.md` — Dashboard layout, unified header, richer task rows

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `enhancedFilteredTasks` in taskStore — already excludes subtasks, applies all filters (status/priority/project/dueDate/ADO link), sorts by selected dimension. Wire directly to new TasksView.
- `filterAdoLink` in taskStore — exists but unused in UI. Values: 'all', 'linked', 'personal'. Wire to FilterCycleButton.
- `groupedEnhanced` in taskStore — groups by status/priority/project. Wire to group-by toggle.
- `isPublic(taskId)` in taskStore — checks `publicTaskIds` set. Use for ADO badge hollow/filled.
- `getSubtasks(parentID)` in taskStore — per-parent subtask loading. Use for detail pane subtasks tab.
- `statusColor()`, `adoTypeIcon()`, `adoTypeColor()` in styles.ts — existing mappings for styling.
- `CommentsSection.vue`, `ExternalLinks.vue` — full-featured components with real API wiring.
- `SyncConfirmDialog.vue`, `ConflictResolver.vue` — sync UX components.

### Established Patterns
- **Teleport for header content**: Views teleport page-specific content to AppShell's `#topbar-actions`. New TasksView should do the same.
- **Pinia store + Wails bindings**: Stores call `@/api/*` wrappers which call Wails bindings.
- **shadcn-vue components**: Badge, Button, Input, ScrollArea, Tabs, Select all available.
- **Lucide icons**: All icons from lucide-vue-next.
- **`lastTask` ref pattern**: TaskDetail uses `lastTask` to survive Vue Transition leave phase.

### Integration Points
- `frontend/src/router/index.ts` — `/tasks` route points to TasksView. Keep same route, replace component.
- `frontend/src/layouts/AppShell.vue` — Shell wraps all views. New TasksView works within AppShell.
- `frontend/src/components/CommandPalette.vue` — Quick-add shortcut (Cmd+N) creates tasks. Should work with new view.
- `main.go` — Register new ADOMetaCacheService with Wails.
- `internal/db/db.go` — Add ado_meta_cache table migration.

</code_context>

<specifics>
## Specific Ideas

- Base the TreeTaskRow on the playground's tree row markup (lines ~731-950 of PlaygroundIntegrated.vue) — the indent levels, expand chevrons, type icons, and state badges are all prototyped there
- The playground's `adoMeta` record pattern maps directly to what `ado_meta_cache` will provide — same shape (task_id → {type, state})
- Use the playground's `flattenedTasks` computed for flat view mode — it recursively walks the tree and outputs indented items
- The `filterAdoLink` store field maps directly to the playground's task list filter (All/ADO/Personal) — just wire it to the FilterCycleButton
- Dashboard's richer task rows (Phase 10) should also use TreeTaskRow once extracted — this creates a single row component used everywhere

</specifics>

<deferred>
## Deferred Ideas

- Batch sync state endpoint (GetAllSyncStates) — Phase 12 scope for dirty-field tracking UI
- Per-subtask sync status tracking — Phase 12 scope
- PR-task linking backend — Phase 12 scope
- N+1 subtask query optimization — Phase 12 scope (current per-parent loading is fine for Phase 11)

</deferred>

---

*Phase: 11-integrate-new-task-list-into-real-app*
*Context gathered: 2026-04-08*
