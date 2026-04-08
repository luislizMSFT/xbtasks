---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
status: Deep UI iteration — playgrounds built, real views being refined
stopped_at: Completed 02-07-PLAN.md
last_updated: "2026-04-06T19:36:25.387Z"
last_activity: 2026-04-03
progress:
  total_phases: 3
  completed_phases: 2
  total_plans: 16
  completed_plans: 16
  percent: 50
---

# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-03-31)

**Core value:** One pane of glass for all your work — personal tasks linked to ADO items, PRs, and comments — so nobody has to context-switch between tools.
**Current focus:** Phase 1.1 — UI Overhaul & Cleanup (then Phase 02 — ADO Integration & Sync Workflow)

## Current Position

Phase: 02 (ado-integration-prs-unified-dashboard) — EXECUTING
Status: Deep UI iteration — playgrounds built, real views being refined
Last activity: 2026-04-03

Progress: [█████░░░░░] 50%

## What's Been Done

### Phase 01 — Foundation (COMPLETE)
- Go backend: auth service, task CRUD, dependencies, projects, SQLite
- Viper config with OS-appropriate paths
- Frontend scaffold with Wails v3 + Vue 3 + vue-router

### Phase 02 — UI Build-Out (IN PROGRESS → SCOPE RESTRUCTURED)
**Completed this session:**
- shadcn-vue installed (79 components), Zinc theme, Tailwind v4
- AppShell: 40px top bar (breadcrumb + stats pills + search/new/activity), icon sidebar, global activity panel
- TasksView: richer rows (checkbox, subtask expand, progress bar, ADO badge with AzureDevOpsIcon, PR count)
- TaskDetail: rebuilt — Work/Discussion tabs, slim collapsible PRs, sticky config footer
- DashboardView: 2-column (tasks+inbox left, PRs right), quick-capture inbox
- ADO view (/ado): Management tab (tree browser, queue & confirm), DevOps tab (2-col PRs+Pipelines)
- PageHeader component for per-page toolbars
- AzureDevOpsIcon component (official Azure DevOps SVG)
- Codebase mapped: 7 docs in .planning/codebase/ (1519 lines)

**Scope change (2026-04-03):**
- Phase 2 restructured: PRs deferred to Phase 3
- Phase 2 now focused on: auth (az cli token), personal→public task model, ADO browser, sync with confirmation, conflict resolution
- New requirements added: TASK-08/09, ADO-06/07/08, SYNC-01-04, updated DASH-01/02
- Phase 3 added: PR Monitoring & Team Views

**Playgrounds built:**
- /playground/tasks — Todoist hybrid layout
- /playground/detail — 2-col with full shell, compact variant
- /playground/chain — Task→PR→Build→Deploy flow chains + blockers dashboard
- /playground/shell — 3 shell layout variants (Compact Rail, Expanded Nav, Floating Panels)
- /playground/ado — ADO tree + GParted queue (now promoted to /ado)

**Still needed:**
- Shell header fix (padding, traffic light overlap, cross-platform macOS/Windows)
- Wire stores to real Go backend (flip useMock to false)
- ADO service implementation (Go backend for ADO API)
- PR service implementation
- Subtasks backend (currently mock data)
- Work chain visualization integration into real views
- Drag-drop task reordering

## Accumulated Context

### Decisions

- Stack: Wails v3 (Go) + Vue 3 + SQLite
- Auth: Abstracted token provider — az cli first, swappable for PAT/OAuth
- ADO client: Direct REST API from Go (no shelling out per query)
- Personal→public task model: tasks start local, become ADO-synced when linked/promoted
- Outbound sync: Always requires preview diff + user confirmation — never auto-push
- Inbound sync: Auto-pull silently on timer + manual refresh
- Conflict resolution: Per-field, user picks local or ADO value
- Subtasks stay personal unless individually linked to ADO
- No emojis — Lucide icons + AzureDevOpsIcon only
- Todoist hybrid task rows: checkbox → title → subtask progress → ADO badge
- ADO browser view: browse assigned items, hide linked, import/link from here
- PRs deferred to Phase 3 — focus on task lifecycle + ADO sync first
- List view filters: status, priority, project, due date, ADO link status (not tags)
- Quick-add: title-only capture, expand later
- Work tasks only (dev + non-dev) — not personal/life tasks
- Desktop app primary, design for future VS Code/MCP integration
- [Phase 01]: Semantic status/priority CSS custom properties added to style.css for UI-SPEC color contract
- [Phase 01]: Accent color uses hex values from UI-SPEC (blue-600 light, blue-500 dark)
- [Phase 01]: Badge components use CSS custom properties via inline style binding for UI-SPEC color contract (not Tailwind utility classes)
- [Phase 01]: color-mix(in srgb, ...) for opacity variants on semantic tokens
- [Phase 01]: PAT button uses ghost variant per UI-SPEC
- [Phase 01]: Kept AppShell in layouts/ directory (Vue convention) rather than components/
- [Phase 01]: Auth guard in App.vue with tryRestore() on mount for session persistence
- [Phase 01]: shadcn-vue Command (cmdk-vue) preferred over raw Headless UI Dialog+Combobox for command palette
- [Phase 01]: DashboardView rebuilt with stat cards layout per UI-SPEC (replacing 2-column playground layout)
- [Phase 01]: TasksView uses useMagicKeys for keyboard shortcuts (Cmd+N/Ctrl+N, Escape)
- [Phase 01]: Expanded status tabs to 6 (All, Todo, In Progress, In Review, Done, Blocked) per UI-SPEC
- [Phase 02]: Multi-org config falls back to legacy single-org when ado.orgs is empty
- [Phase 02]: SyncState uses composite PK (task_id, ado_id) for per-link conflict tracking
- [Phase 02]: TokenProvider interface uses GetToken()+Name(); CachedTokenProvider wraps with TTL mutex caching
- [Phase 02]: ADO client accepts token string (not provider)  caller manages token lifecycle
- [Phase 02]: ADOService refactored from az-cli to pkg/ado REST client with multi-org iteration
- [Phase 02]: Token provider chain: AzCli -> CachedWrapper created once in main.go, shared across services
- [Phase 02]: LinkService implements personal-to-public model with 4 flows: link/promote/import/unlink
- [Phase 02]: IsPublic computed from task_ado_links presence (not a column)
- [Phase 02]: Sync store uses dynamic imports for Wails bindings (same pattern as ado.ts)
- [Phase 02]: ConflictResolver walks conflicts sequentially - shows first, resolves, advances to next
- [Phase 02]: SyncService uses ticker goroutine with configurable interval; outbound never auto-pushes
- [Phase 02]: SyncService uses syncMu to serialize background and manual sync; emits lifecycle events (sync:started/completed/failed)
- [Phase 02]: CachedTokenProvider uses RWMutex + refreshMu double-check pattern — reads non-blocking, refresh serialized
- [Phase 02]: All az CLI calls use exec.CommandContext with 15s timeout (AzCliTimeout constant)
- [Phase 02]: ADO multi-org fan-out has 20s timeout; returns partial results on slow orgs
- [Phase 02]: Frontend subscribes to backend Wails events via Events.On() from @wailsio/runtime
- [Phase 02]: Sync store has initEvents() called once on auth in App.vue for backend→frontend push
- [Phase 02]: ADO tree search uses useDebounceFn (200ms) from @vueuse/core
- [Phase 02]: Domain SyncDiff/FieldDiff types separate from ado package for frontend decoupling
- [Phase 02]: Comments always private/local by default; PushCommentToADO is explicit opt-in
- [Phase 02]: ProjectService constructor extended with tokenProv+cfg for ADO operations
- [Phase 02]: ExternalLinks uses window.open for real browser opening (UX-02)
- [Phase 02]: CommentsSection shows Push to ADO only for private comments on linked tasks
- [Phase 02]: Az CLI auth fetches ADO profile for real user info, falls back to generic user
- [Phase 02]: Settings page replaced legacy single-org inputs with multi-org list management
- [Phase 02]: Card grid replaces list+detail split panel for Projects page
- [Phase 02]: Dual progress bars (local+ADO) kept separate, not merged, per PROJ-06
- [Phase 02]: FilterBar uses AcceptableValue from reka-ui for type-safe Select handlers
- [Phase 02]: Quick-add always-visible Input (not toggle) for faster task capture
- [Phase 02]: enhancedFilteredTasks replaces grouped-by-status as primary rendering pipeline
- [Phase 02]: Recursive render function component for tree nodes (defineComponent with h()) for same-file recursion
- [Phase 02]: ListLinkedAdoIDs method added to LinkService for frontend linked status tracking (ADO IDs not task IDs)
- [Phase 02]: Saved query picker uses sentinel value __my_assignments__ to revert to default tree fetch

### Roadmap Evolution

- Phase 4 added: Work Item Lifecycle Tracking — PRs, Pipelines & Task Traceability (unified progress view so you never have to hunt across ADO)
- Phase 5 added: Fix TaskDetail null crash and UI freeze
- Phase 6 added: Dashboard header summary display
- Phase 7 added: Task hierarchy tree rendering
- Phase 8 added: Task row styling parity with ADO
- Phase 9 added: Async Architecture Polish — context propagation, debounce, circuit breakers, SQLite retry, lazy-load (deferred items from async audit)
- Phase 10 added: Implement dashboard redesign and unified header bars

### Pending Todos

- Shell header: fix padding, traffic light overlap, Windows compat
- Wire mock stores to real Go backend
- Implement abstracted token provider (az cli → token → REST)
- Implement ADO REST client package (`pkg/ado/`)
- Refactor ADOService to use REST client instead of az cli subprocess
- Build ADO browser view (browse items, linked status, import/link)
- Implement personal→public task model in frontend (badges, transitions)
- Build sync confirmation dialog (preview diff for outbound changes)
- Build conflict resolution UI (per-field picker)
- Add list view filters (status, priority, project, due date, ADO link status)
- Add quick-add task input
- Subtask backend implementation
- Drag-drop task reordering
- Splash screen (Xbox logo)

### Blockers/Concerns

- Wails v3 is alpha — may need v2 fallback if blockers arise
- Detached HEAD state — commits need to be merged to a branch
- az cli must be installed and authenticated on user's machine for initial auth flow

## Session Continuity

Last session: 2026-04-06T19:26:50.644Z
Stopped at: Completed 02-07-PLAN.md
Resume: Continue refining shell header, then wire backend services
