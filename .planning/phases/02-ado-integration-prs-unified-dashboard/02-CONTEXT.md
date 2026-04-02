# Phase 2: ADO Integration, PRs & Unified Dashboard - Context

**Gathered:** 2026-04-02
**Status:** Ready for planning

<domain>
## Phase Boundary

Bridge local tasks ↔ ADO items, surface PRs under tasks, add comments/timeline to tasks, and deliver a unified dashboard. ADO is a management/linking layer — the primary experience is the user's own tasks, projects, and notes. ADO answers "is this tracked upstream?" not "what should I work on?"

Our domain: **Projects** and **Tasks** (with subtasks). ADO has its own hierarchy (Scenario → Deliverable → Task/Bug/Story) but we don't replicate that in our domain. Users **link** our projects/tasks to any ADO item they choose.

</domain>

<decisions>
## Implementation Decisions

### ADO Tree & Bulk Linking
- **D-01:** ADO Items tab shows full tree hierarchy — expandable Scenario → Deliverable → Task/Bug/Story with type icons (🐛 Bug, ✅ Task, 📖 Story, 📦 Deliverable, 🎯 Scenario)
- **D-02:** Tree scope: sync items assigned to me + their parent chain (Deliverable → Scenario context). Full team items visible in ADO Items tab.
- **D-03:** Auto-sync on configurable interval (default 15min via Viper) + manual refresh button
- **D-04:** Bulk linking via checkboxes on tree nodes → confirm view showing selected items
- **D-05:** Bulk confirm is a table with smart auto-match — suggest links based on title similarity, user confirms/overrides. Per-item actions: Track (import as new local task), Link (to existing task/project), Skip
- **D-06:** Tree search: hybrid — instant search of local SQLite cache + "Search all ADO" button for live API query. Results show in-tree with highlights OR flat list with breadcrumb paths (toggle between views). Recursive search through full hierarchy.
- **D-07:** ADO bugs can be linked as tasks — they keep their bug icon (🐛) throughout the app

### PRs Under Tasks
- **D-08:** PRs belong under tasks, not a separate view. Task detail has a collapsible PR section + PRs appear in the activity timeline
- **D-09:** Tasks can have multiple PRs. Each PR row is expandable: compact by default (title, status badge, vote summary), click to expand (repo, branches, all reviewers, comment count)
- **D-10:** PR association: auto-link from ADO (if PR references a linked work item) + manual add for unlinked PRs
- **D-11:** Click PR title → deep link opens ADO PR page in browser (no in-app diff viewer)
- **D-12:** Sync all team PRs from configured repos. Dashboard shows personal PRs (authored + assigned for review + followed). ADO Items tab shows all team PRs.
- **D-13:** PR mark-as-viewed: unread dot on PR rows, click to mark viewed. "X unread" counter on dashboard

### Comments & Timeline
- **D-14:** No standalone "notes" entity. Two concepts: **Comments** on tasks (thread style) and **Description** field
- **D-15:** Comments have two tabs: **Personal** (local only, private) and **ADO** (synced to linked ADO work item)
- **D-16:** Task description can optionally sync to ADO linked work item
- **D-17:** Rich text WYSIWYG for comments and descriptions (bold, italic, links, code blocks, lists)
- **D-18:** Full activity timeline per task — comments, status changes, PR events, ADO sync events, all chronological in detail panel
- **D-19:** @mention syntax for quick linking (type @task-123 or #ADO-4829) + toolbar button for search-based linking within comments
- **D-20:** Project-level timeline: aggregates all child task activity into one chronological stream

### Dashboard
- **D-21:** Fixed widget/card grid layout (no customization for v1). PRs on top.
- **D-22:** Dashboard is personal: my PRs (authored + reviewing + followed), my focus tasks, stats, blocked items
- **D-23:** ADO Items tab is the team-wide view — all team PRs + all ADO items

### Task Detail Panel Reorg
- **D-24:** Top: Title (editable) → Subtasks with progress bar → PR section (collapsible)
- **D-25:** Middle: Description (rich text, optional ADO sync) → Comments (Personal | ADO tabs) → Activity Timeline
- **D-26:** Bottom: Configuration (status, priority, tags, due date, project, ADO link) → Delete
- **D-27:** Progress bar on task rows and in detail panel (X of Y subtasks complete)

### ADO Linking UX
- **D-28:** Task row ADO icon: hollow/empty if unlinked, filled with ADO type icon if linked. Click hollow → opens link dialog
- **D-29:** Projects can be linked to ADO Scenarios or Deliverables — user's choice, we don't enforce hierarchy

### Agent's Discretion
- Rich text editor library choice (Tiptap, ProseMirror, etc.)
- Activity timeline event rendering design
- Smart auto-match algorithm for bulk linking
- PR sync implementation (ADO REST API, polling strategy)
- Dashboard card layout/sizing

### Deferred Ideas
- Customizable dashboard widget grid (drag to rearrange) — v2
- Upcoming/reminder widget on dashboard — v2
- Xbox splash screen with zoom animation — v2
- In-app PR diff viewer — if needed later

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Context
- `.planning/PROJECT.md` — Full project vision, xl porting context (ADO integration patterns, WIQL queries, JSON Patch, state mapping), key decisions
- `.planning/REQUIREMENTS.md` — ADO-01 through ADO-05, PR-01 through PR-03, DASH-01 through DASH-03
- `.planning/ROADMAP.md` — Phase 2 success criteria and scope boundary

### Phase 1 Context
- `.planning/phases/01-foundation-auth-personal-tasks/01-CONTEXT.md` — Phase 1 decisions (sidebar, command palette, system tray, theme)

### Research (Port Reference)
- `.planning/research/XL-CODE-REFERENCE.md` — xl's proven ADO client patterns: PAT + Azure CLI auth, WIQL queries, JSON Patch, state mapping, bulk import, per-repo config

### Existing Code
- `internal/db/db.go` — SQLite schema with tables: ado_work_items, task_ado_links, pull_requests (already created in Phase 1)
- `internal/config/config.go` — Viper config with ado.organization, ado.project, sync.interval_minutes
- `internal/config/service.go` — ConfigService exposed to frontend via Wails bindings
- `frontend/src/stores/tasks.ts` — Task store with mock data, binding fallback pattern
- `frontend/src/views/TasksView.vue` — Current task list with status sections (needs ADO Items tab)
- `frontend/src/components/TaskDetail.vue` — Current detail panel (needs PR section, comments, timeline)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/db/db.go` schema already has: `ado_work_items` (ado_id, title, state, type, assigned_to, priority, area_path, url), `task_ado_links` (task_id, ado_id, direction), `pull_requests` (title, pr_url, repo, task_id, ado_id, status, reviewers, votes, branches)
- `internal/config/` — Viper config with ADO org/project/sync interval, ConfigService for frontend
- `frontend/src/components/ui/` — Full shadcn-vue library (79 components): button, badge, card, select, input, textarea, separator, tooltip, scroll-area, dialog, dropdown-menu, command, tabs
- `frontend/src/lib/utils.ts` — `cn()` utility for class merging
- `frontend/src/components/ui/AdoBadge.vue` — Existing ADO badge component (needs type icon variants)

### Established Patterns
- **Wails service binding:** Go structs → frontend via `application.NewService()`. Add ADOService, PRService
- **Store pattern:** Pinia stores with `useMock` flag, try bindings → catch fallback to mock
- **shadcn-vue tokens:** bg-primary/text-primary-foreground for blue, bg-muted for subtle, text-foreground/text-muted-foreground

### Integration Points
- `main.go` — Register new services (ADOService, PRService)
- `frontend/src/router/index.ts` — No new routes needed (ADO Items is a tab within Tasks page)
- `frontend/src/views/TasksView.vue` — Add ADO Items, Linked tabs via shadcn Tabs
- `frontend/src/components/TaskDetail.vue` — Add PR section, comments, timeline
- `frontend/src/views/DashboardView.vue` — Add PR widget, redesign as card grid

</code_context>

<specifics>
## Specific Ideas

- ADO bugs can be tasks — they keep their bug icon throughout
- Smart auto-match for bulk linking uses title similarity to suggest which local task to link to
- "Search all ADO" button in tree search for items outside your assignments
- Comments are distinguished by tab (Personal vs ADO), not by a separate notes entity
- Dashboard is the personal view; ADO Items tab is the team-wide view
- Breadcrumb paths in flat search results (Scenario › Deliverable › Task)
- Progress bar: thin (2-3px), colored by status, on both task rows and detail panel

</specifics>

---

*Phase: 02-ado-integration-prs-unified-dashboard*
*Context gathered: 2026-04-02*
