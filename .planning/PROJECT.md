# Team ADO Tool

## What This Is

A unified productivity dashboard for the Xbox Services team — a native desktop app (Wails v3 + Vue + Go + SQLite) that brings all work — ADO work items, PRs, comments, and personal tasks — into a single view. It replaces the fragmented landscape where team members use different tools (ADO, Wiki, Loop, CLIs, M365 Copilot) with one app everyone gravitates to. Local-first with SQLite, system tray presence, native notifications.

## Core Value

One pane of glass for all your work — personal tasks linked to ADO items, PRs, and comments — so nobody has to context-switch between tools to know what they need to do.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Unified list view showing all tasks (personal + public/ADO-linked) with visual badge distinguishing them
- [ ] Personal task management with create (quick-add + expand), edit, prioritize, subtask support, dependencies
- [ ] Personal→public task model: tasks start local-only, become ADO-synced when linked/promoted
- [ ] Bidirectional ADO linking — promote local task → ADO work item, pull ADO item → local task, link to existing
- [ ] ADO browser view — browse assigned items, see linked status, toggle hide linked, import from here
- [ ] ADO work items surfaced with status, priority, assignment
- [ ] Multi-org support — configure multiple ADO orgs/projects, unified view with group-by-project toggle
- [ ] External links on tasks — attach ICMs, Grafana dashboards, wiki pages with auto-detected type icons
- [ ] Local comments (private) with selective push to ADO (clearly marked as public)
- [ ] Investigation hub workflow — start local or import ADO bug, enrich with links/subtasks/notes, link to ADO anytime
- [ ] Bidirectional sync — auto-pull silently, outbound requires preview diff + user confirmation
- [ ] Per-field conflict resolution when both sides changed a linked item
- [ ] Tasks linkable to ADO deliverables/scenarios with clear visual indication
- [ ] User can set personal priority and notes on any linked item
- [ ] Auth via abstracted token provider (az cli initially, swappable for PAT/OAuth)
- [ ] List view filterable by status, priority, project, due date, ADO link status
- [ ] Subtasks of public tasks stay personal unless individually linked to ADO

### Out of Scope

- GitHub integration — ADO is the primary system; GitHub support deferred to later milestone
- VS Code extension — planned for code review flow, but not v1
- Local code review (reviewme-style) — depends on VS Code extension, deferred
- Notifications feed — valuable but not core to the unified view; deferred
- Linkwarden integration — personal power feature, not team dependency
- Real-time collaboration — not competing with Teams/Slack
- Mobile app — web-first
- PR monitoring & review — deferred to Phase 3; focus on task lifecycle + ADO sync first
- Pipeline status — deferred to Phase 3
- Personal/life tasks — work tasks only (dev + non-dev like meetings, docs, reviews)
- Auto-push to ADO — all outbound changes require user confirmation

## Context

- **Team environment:** Xbox Services engineering team. Everyone uses ADO for organizational work tracking (work items, bugs, deliverables, scenarios, PRs). Individual productivity tools are fragmented — some use Wiki, others ADO boards, others Loop, Luis uses a CLI. No shared surface.
- **Prior work — xl TUI:** Luis built `xl` (Go/Bubble Tea), a personal TUI managing 17 entities in SQLite: tasks, projects, PRs, reviews, daily_logs, incidents, goals, areas, milestones, notes, memory, meetings, relationships, sessions, weekly_reviews. It has a full ADO integration (push/pull/sync/link with bidirectional status mapping), tree view (Projects→Tasks→PRs), knowledge browser (memory/notes/meetings), FTS5 search, filter panel, and graph renderer. Currently tracks 29 tasks, 10 PRs, 53 daily logs, 158 memory entries, 2 incidents. The data model and ADO integration patterns are proven and should inform this tool's design.
- **Prior work — reviewme:** A working VSCode extension POC for local code review — compare branches with native diff editor without creating a PR. Phase 1 complete (TypeScript + simple-git). Designed for reviewing agent-generated code before merging. Phases 2-5 planned (diff browsing, annotations, review workflow, hosting integration).
- **xl's ADO integration (reusable patterns):** PAT + Azure CLI auth, WIQL queries, JSON Patch operations for updates, state mapping (todo→Proposed, in_progress→Active, done→Completed, blocked→Blocked), parent-child relationship creation, bulk import (pull-mine, pull-children), per-repo config stored in memory table.
- **xl's data model (proven entities):** Tasks (P0-P3 priority, 7 statuses, categories, ADO linking, time tracking), PRs (full lifecycle + deployment tracking), Reviews (others' PRs), Incidents (sev1-4), Goals/OKRs, Daily logs, Notes, Memory (persistent key-value), Relationships (knowledge graph edges), Task dependencies.
- **The fragmentation problem:** Team members each have their own system. M365 Copilot can't reliably be used as a task manager. No way to connect personal task management to organizational work items.
- **Agent-era code review:** With agents writing code on top of developer work, there's a gap in reviewing those changes locally before they become a PR. The reviewme POC addresses this; it'll integrate via VS Code extension in a future milestone.
- **ADO is central:** ADO manages everything for the team including PRs (Azure DevOps Repos, not GitHub). The tool layers on top of ADO, not beside it.

## Constraints

- **Auth**: Abstracted token provider — az cli `get-access-token` initially, swappable for PAT or Entra ID OAuth later. Team is on Microsoft ecosystem.
- **ADO API**: Direct ADO REST API calls from Go using token from provider — no shelling out to az cli per query. Must handle rate limits and pagination. Supports multiple org/project pairs.
- **Sync safety**: All outbound changes to ADO require preview diff + user confirmation. Never auto-push. Subtasks/personal breakdowns never pushed unless individually linked.
- **Audience**: Start with Luis (dogfooding), grow to team — must be useful for one person before it scales
- **Stack**: Wails v3 (Go) + Vue 3 (thin shell) + SQLite — native desktop app, not a web app. Design for future VS Code/MCP extension integration.
- **No XAML, no React, no C#**: Team prefers Go backend with lightweight web frontend in native shell
- **Vue is thin shell**: All logic lives in Go. Vue is display + interaction only.
- **Local-first**: SQLite per user, no server dependency for core functionality
- **Work tasks only**: Dev and non-dev work tasks (meetings, docs, reviews) — not personal/life tasks

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| ADO as source of truth, tool as productivity layer | Team already lives in ADO; tool adds personal workflow on top | — Pending |
| Wails v3 + Vue + Go + SQLite | Go backend (port xl), Vue thin shell, native desktop, local-first SQLite | — Pending |
| Native desktop app over web app | System tray, notifications, offline-capable, no browser dependency | — Pending |
| Bidirectional task linking (not just read-only ADO view) | Users need to promote/pull tasks, not just see them | — Pending |
| Personal→public task model | Tasks start local-only, become ADO-synced when explicitly linked/promoted. Subtasks stay personal unless individually linked. | — Decided 2026-04-03 |
| Outbound sync requires confirmation | Preview diff before any push to ADO — prevents accidental updates, protects personal breakdowns | — Decided 2026-04-03 |
| Abstracted token provider (az cli first) | az cli is easiest auth path; abstract so PAT/OAuth can swap in later | — Decided 2026-04-03 |
| Direct ADO REST API from Go | Grab token from provider, call REST directly — no shelling out per query | — Decided 2026-04-03 |
| PRs deferred to Phase 3 | Focus on task lifecycle + ADO sync first; PRs are additive, not core | — Decided 2026-04-03 |
| Multi-org support | Team members work across multiple ADO orgs/projects; configure orgs → pick projects, unified view | — Decided 2026-04-03 |
| Design for future VS Code/MCP integration | Desktop app primary, but architecture should support other surfaces | — Decided 2026-04-03 |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-04-03 — restructured per workflow discussion (personal/public model, sync confirmation, ADO browser, auth via az cli token, PRs deferred)*
