# Roadmap: Team ADO Tool

## Overview

Deliver a native desktop productivity dashboard (Wails v3 + Vue 3 + Go + SQLite) that gives the Xbox Services team one pane of glass for all their work. Two coarse phases: first deliver a solid foundation with auth and full personal task management, then layer on ADO integration, PRs, and the unified dashboard that ties it all together. Porting proven patterns from xl (Go ADO client, SQLite schema, data model) accelerates both phases.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (1.1, 1.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Foundation, Auth & Personal Tasks** - Native desktop app with Entra ID sign-in, SQLite database, and full personal task management (CRUD, statuses, subtasks, dependencies, tags)
- [ ] **Phase 2: ADO Integration, PRs & Unified Dashboard** - Fetch ADO work items, bridge tasks ↔ ADO, surface PRs, and tie everything together in one unified view

## Phase Details

### Phase 1: Foundation, Auth & Personal Tasks
**Goal**: User can launch the native desktop app, sign in with their Microsoft work account, and manage all personal work with full CRUD, organization, hierarchy, and workflow states
**Depends on**: Nothing (first phase)
**Requirements**: AUTH-01, AUTH-02, AUTH-03, TASK-01, TASK-02, TASK-03, TASK-04, TASK-05, TASK-06, TASK-07
**Success Criteria** (what must be TRUE):
  1. User can launch the native desktop app and authenticate with their Microsoft work account (Entra ID OAuth2 PKCE)
  2. User session persists across app restarts without re-authenticating, and user profile (display name, email, avatar) is visible in the app
  3. User can create, edit, and delete personal tasks with title, description, priority (P0-P3), and category
  4. User can move tasks through statuses (todo → in_progress → in_review → done, plus blocked and cancelled) and organize tasks with tags
  5. User can create subtasks under any parent task, define task dependencies (A blocks B), and set a personal priority overlay on any task
**Plans**: TBD
**UI hint**: yes

### Phase 2: ADO Integration, PRs & Unified Dashboard
**Goal**: User can see all assigned ADO work, bridge tasks ↔ ADO bidirectionally, view PRs, and access everything from one unified dashboard
**Depends on**: Phase 1
**Requirements**: ADO-01, ADO-02, ADO-03, ADO-04, ADO-05, PR-01, PR-02, PR-03, DASH-01, DASH-02, DASH-03
**Success Criteria** (what must be TRUE):
  1. User can see all ADO work items assigned to them with state, priority, type, assigned_to, description, and area path
  2. User can link a personal task to an ADO work item, promote a task to a new ADO work item, and import an ADO work item as a local task
  3. User can see all their own PRs and PRs assigned to them for review with status (draft, active, completed, abandoned), title, repo, branch info, reviewers, and votes
  4. User sees a unified dashboard on login showing tasks, ADO items, and PRs together in one view
  5. Dashboard visually distinguishes personal tasks from ADO-linked items and shows link connections with status badges
**Plans**: TBD
**UI hint**: yes

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation, Auth & Personal Tasks | 0/0 | Not started | - |
| 2. ADO Integration, PRs & Unified Dashboard | 0/0 | Not started | - |
