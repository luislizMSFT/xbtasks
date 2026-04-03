# Roadmap: Team ADO Tool

## Overview

Deliver a native desktop productivity dashboard (Wails v3 + Vue 3 + Go + SQLite) that gives the Xbox Services team one pane of glass for all their work. Three phases: solid foundation with personal task management, then ADO integration with the personal→public workflow and safe bidirectional sync, then PR monitoring and team views. Porting proven patterns from xl (Go ADO client, SQLite schema, data model) accelerates development.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (1.1, 1.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Foundation, Auth & Personal Tasks** - Native desktop app with SQLite database and full personal task management (CRUD, statuses, subtasks, dependencies, tags)
- [ ] **Phase 1.1: UI Overhaul & Cleanup** *(INSERTED)* - Clean out playground pages, redesign TaskDetail sidebar, redo Projects page, add dependency graph visualization, integrate real APIs (SQLite + az cli)
- [ ] **Phase 2: ADO Integration & Sync Workflow** - Abstracted token provider, ADO REST client, personal→public task model, ADO browser view, bidirectional sync with confirmation, conflict resolution, unified list view with filters
- [ ] **Phase 3: PR Monitoring & Team Views** - Surface PRs under tasks, pipeline status, team-wide views, PR comment threads

## Phase Details

### Phase 1: Foundation, Auth & Personal Tasks
**Goal**: User can launch the native desktop app and manage all personal work with full CRUD, organization, hierarchy, and workflow states
**Depends on**: Nothing (first phase)
**Requirements**: TASK-01, TASK-02, TASK-03, TASK-04, TASK-05, TASK-06, TASK-07
**Success Criteria** (what must be TRUE):
  1. User can launch the native desktop app with SQLite persistence
  2. User can create, edit, and delete personal tasks with title, description, priority (P0-P3), and category
  3. User can move tasks through statuses (todo → in_progress → in_review → done, plus blocked and cancelled) and organize tasks with tags
  4. User can create subtasks under any parent task, define task dependencies (A blocks B), and set a personal priority overlay on any task
**Plans**: TBD
**UI hint**: yes

### Phase 2: ADO Integration & Sync Workflow
**Goal**: User can authenticate via az cli, browse ADO items, link/promote/import tasks with the personal→public model, and bidirectionally sync with safe confirmation before any outbound changes
**Depends on**: Phase 1
**Requirements**: AUTH-01, AUTH-02, AUTH-03, TASK-08, TASK-09, ADO-01, ADO-02, ADO-03, ADO-04, ADO-05, ADO-06, ADO-07, ADO-08, SYNC-01, SYNC-02, SYNC-03, SYNC-04, DASH-01, DASH-02, DASH-03
**Success Criteria** (what must be TRUE):
  1. App authenticates to ADO via abstracted token provider (az cli initially); token auto-refreshes
  2. Tasks start as personal (local-only). Linking or promoting makes them "public" (ADO-synced). Unified list view shows both with visual distinction (badge/icon)
  3. User can quick-add a task with just a title, then expand with full details later
  4. ADO browser view shows all assigned ADO items, indicates which are already linked, allows toggling to hide linked items, and supports selecting items to import/link
  5. User can link a personal task to an existing ADO item, promote a task to a new ADO work item, or import an ADO item as a local task — all three flows work
  6. Background auto-sync pulls ADO changes silently on a timer + manual refresh. All outbound pushes show a preview diff of what will change and require user confirmation — never auto-push
  7. When both local and ADO changed the same linked item, user sees per-field conflict resolution
  8. Linked tasks sync title, status, and description to ADO. Subtasks stay personal unless individually linked
  9. List view filterable by status, priority, project, due date, and ADO link status
**Plans**: TBD
**UI hint**: yes

### Phase 3: PR Monitoring & Team Views
**Goal**: Surface PRs under tasks, show pipeline status, enable team-wide visibility into work and blockers
**Depends on**: Phase 2
**Requirements**: PR-01, PR-02, PR-03, PR-04, PIPE-01, PIPE-02 (v2 requirements)
**Success Criteria** (what must be TRUE):
  1. User can see all their own PRs and PRs assigned to them for review with status, title, repo, branch info, reviewers, and votes
  2. PRs appear under linked tasks in the detail panel
  3. Pipeline run status visible on linked PRs/tasks
  4. Team-wide views for PRs and work items available
**Plans**: TBD
**UI hint**: yes

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 1.1 → 2 → 3

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & Personal Tasks | 6/6 | Complete | ✓ |
| 1.1. UI Overhaul & Cleanup | 0/0 | Not started | - |
| 2. ADO Integration & Sync Workflow | 0/0 | Not started | - |
| 3. PR Monitoring & Team Views | 0/0 | Not started | - |
