# Team ADO Tool

## What This Is

A unified productivity dashboard for the Xbox Services team — a native desktop app (Wails v3 + Vue + Go + SQLite) that brings all work — ADO work items, PRs, comments, and personal tasks — into a single view. It replaces the fragmented landscape where team members use different tools (ADO, Wiki, Loop, CLIs, M365 Copilot) with one app everyone gravitates to. Local-first with SQLite, system tray presence, native notifications.

## Core Value

One pane of glass for all your work — personal tasks linked to ADO items, PRs, and comments — so nobody has to context-switch between tools to know what they need to do.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Unified dashboard showing all work in one view (ADO items, PRs, personal tasks)
- [ ] Personal task management with create, edit, prioritize, subtask support
- [ ] Bidirectional ADO linking — promote local task → ADO work item, pull ADO item → local task
- [ ] ADO work items surfaced with status, priority, assignment
- [ ] ADO PRs surfaced with review status and comment threads
- [ ] PR comment threads viewable in-app (ADO-style)
- [ ] Tasks linkable to ADO deliverables/scenarios with clear visual indication
- [ ] User can set personal priority and notes on any linked item
- [ ] Auth via Microsoft identity (ADO/Entra ID)
- [ ] Multi-role support — ICs, infra engineers, managers all find value

### Out of Scope

- GitHub integration — ADO is the primary system; GitHub support deferred to later milestone
- VS Code extension — planned for code review flow, but not v1
- Local code review (reviewme-style) — depends on VS Code extension, deferred
- Notifications feed — valuable but not core to the unified view; deferred
- Linkwarden integration — personal power feature, not team dependency
- Real-time collaboration — not competing with Teams/Slack
- Mobile app — web-first

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

- **Auth**: Must use Microsoft identity / Entra ID — team is on Microsoft ecosystem. PAT fallback for simplicity.
- **ADO API**: Azure DevOps REST API is the primary integration surface; must handle rate limits and pagination
- **Audience**: Start with Luis (dogfooding), grow to team — must be useful for one person before it scales
- **Stack**: Wails v3 (Go) + Vue 3 (thin shell) + SQLite — native desktop app, not a web app
- **No XAML, no React, no C#**: Team prefers Go backend with lightweight web frontend in native shell
- **Vue is thin shell**: All logic lives in Go. Vue is display + interaction only.
- **Local-first**: SQLite per user, no server dependency for core functionality

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| ADO as source of truth, tool as productivity layer | Team already lives in ADO; tool adds personal workflow on top | — Pending |
| Wails v3 + Vue + Go + SQLite | Go backend (port xl), Vue thin shell, native desktop, local-first SQLite | — Pending |
| Native desktop app over web app | System tray, notifications, offline-capable, no browser dependency | — Pending |
| Bidirectional task linking (not just read-only ADO view) | Users need to promote/pull tasks, not just see them | — Pending |

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
*Last updated: 2026-03-31 after initialization*
