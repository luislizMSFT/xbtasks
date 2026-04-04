---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
status: Deep UI iteration — playgrounds built, real views being refined
stopped_at: Completed 01-03-PLAN.md
last_updated: "2026-04-04T00:04:22.371Z"
last_activity: 2026-04-03
progress:
  total_phases: 3
  completed_phases: 0
  total_plans: 6
  completed_plans: 3
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

Last session: 2026-04-04T00:04:22.367Z
Stopped at: Completed 01-03-PLAN.md
Resume: Continue refining shell header, then wire backend services
