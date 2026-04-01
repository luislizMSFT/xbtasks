# Phase 1: Foundation, Auth & Personal Tasks - Context

**Gathered:** 2026-04-01
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver a native desktop app where the user can launch the app, sign in with their Microsoft work account (Entra ID OAuth2 PKCE), and manage all personal work with full CRUD, organization (projects, tags, categories), hierarchy (subtasks, dependencies), and workflow states (todo → in_progress → in_review → done, plus blocked and cancelled). This is the foundation layer — no ADO integration, no PR views, no unified dashboard yet.

</domain>

<decisions>
## Implementation Decisions

### Navigation & App Shell
- **D-01:** Slim icon sidebar (Linear/Notion style) — icons for each section, content area takes most of the screen. Modern, desktop-native feel.
- **D-02:** Command palette (⌘K / Ctrl+K) for power-user search and management — navigate anywhere, create tasks, change status, all from keyboard.
- **D-03:** Minimize to system tray on window close — app stays running in background, always available (Slack/1Password pattern). Click tray icon to restore.
- **D-04:** Color scheme follows system preference (auto dark/light). Respect macOS/Windows appearance setting.
- **D-05:** Landing page is a dashboard overview — today's work, recent tasks, quick stats. Not a raw task list.

### Agent's Discretion
- Auth flow details (OAuth2 PKCE implementation, PAT fallback, token storage approach) — planner/researcher decides best approach for Wails v3 desktop context
- Task list layout (table vs cards, grouping, sort/filter controls) — agent picks based on codebase patterns and desktop UX conventions
- Task detail experience (inline editing, side panel, full page) — agent picks based on what works best with the sidebar + content area layout
- Tags UX (free-form vs predefined, autocomplete) — agent decides based on the existing CSV tags column
- Dependency visualization (text, tree, graph) — agent picks appropriate for v1 complexity

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Context
- `.planning/PROJECT.md` — Full project vision, constraints, xl porting context, key decisions
- `.planning/REQUIREMENTS.md` — AUTH-01 through AUTH-03, TASK-01 through TASK-07 acceptance criteria
- `.planning/ROADMAP.md` — Phase 1 success criteria and scope boundary

### Existing Code (Port Reference)
- `.planning/research/XL-CODE-REFERENCE.md` — xl's proven ADO client, SQLite schema, and data model patterns to port
- `.planning/research/ARCHITECTURE.md` — System architecture (needs updating for desktop — still describes web BFF pattern)
- `.planning/research/PITFALLS.md` — 7 critical pitfalls with prevention strategies

### Stack Research
- `.planning/research/STACK.md` — Final stack decision (Wails v3 + Vue + Go + SQLite)
- `.planning/research/FEATURES.md` — Feature landscape (table stakes, differentiators, anti-features)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/db/db.go` — SQLite with WAL mode, foreign keys, full schema for tasks, projects, PRs, ADO work items, users, task_deps, task_ado_links. Migration on startup.
- `internal/app/tasks.go` — TaskService with Create, GetByID, List (filterable by status), Update, Delete, SetStatus, GetSubtasks. Bound to Wails frontend.
- `internal/app/projects.go` — ProjectService with full CRUD. Bound to Wails frontend.
- `pkg/models/models.go` — Go structs for Task, Project, PullRequest, ADOWorkItem, TaskADOLink, User with JSON tags.
- `frontend/src/style.css` — Tailwind CSS configured and importing.

### Established Patterns
- **Wails service binding:** Go structs with exported methods are bound to frontend via `application.NewService()`. Frontend calls Go methods directly (no REST API).
- **SQLite schema-in-code:** Schema defined as const string in `db.go`, auto-migrated on startup.
- **Priority ordering:** Tasks sorted by priority (P0→P3) then updated_at DESC.
- **Status flow:** todo → in_progress → in_review → done + blocked/cancelled. `completed_at` auto-set when status becomes "done".

### Integration Points
- `main.go` — Register new services here (e.g., AuthService, DependencyService)
- `internal/auth/` — Empty package, ready for Entra ID implementation
- `frontend/src/components/` — Only HelloWorld.vue exists, entire UI to be built
- `frontend/src/App.vue` — Root component, needs app shell/layout structure

</code_context>

<specifics>
## Specific Ideas

- Command palette should be a first-class citizen, not an afterthought — the user explicitly wants it alongside sidebar nav
- System tray presence means the app should feel "always there" — quick access to create tasks, check status without full window restore
- Dashboard overview as landing implies we need a lightweight stats/summary component even in Phase 1

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 01-foundation-auth-personal-tasks*
*Context gathered: 2026-04-01*
