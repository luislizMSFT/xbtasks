# Requirements: Team ADO Tool

**Defined:** 2026-03-31
**Core Value:** One pane of glass for all your work — personal tasks linked to ADO items, PRs, and comments — so nobody has to context-switch between tools.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Authentication

- [x] **AUTH-01**: App authenticates to ADO via abstracted token provider — az cli `get-access-token` initially, swappable for PAT or OAuth later
- [x] **AUTH-02**: Token auto-refreshes transparently; user doesn't need to re-auth unless az cli session expires
- [x] **AUTH-03**: Token provider abstraction allows future swap to Entra ID OAuth2 PKCE or PAT without changing consuming code

### Personal Tasks

- [x] **TASK-01**: User can create a personal task with title, description, priority (P0-P3), and type (bug, feature, meeting, doc, investigation)
- [x] **TASK-02**: User can edit and delete their own tasks
- [x] **TASK-03**: User can set task status (todo, in_progress, in_review, done, blocked, cancelled)
- [x] **TASK-04**: User can create subtasks under a parent task or ADO-linked item
- [x] **TASK-05**: User can set a personal priority overlay independent of ADO priority
- [x] **TASK-06**: User can define dependencies between tasks (task A blocks task B)
- [x] **TASK-07**: User can add tags to tasks for organization
- [x] **TASK-08**: Tasks start as personal (local-only). Linking/promoting to ADO makes them "public" (synced). Subtasks of public tasks stay personal unless individually linked.
- [x] **TASK-09**: Quick-add task with just a title for fast capture; expand with full details (description, priority, project, due date) later

### ADO Integration

- [x] **ADO-01**: User can view all ADO work items assigned to them in the dashboard
- [x] **ADO-02**: User can view ADO work item details (state, priority, type, assigned_to, description, area path)
- [x] **ADO-03**: User can link a personal task to an ADO work item (bidirectional)
- [x] **ADO-04**: User can promote a personal task to a new ADO work item (with confirmation)
- [x] **ADO-05**: User can import an ADO work item as a personal task with local tracking
- [x] **ADO-06**: ADO browser view — browse all assigned ADO items, show which are already linked to local tasks, toggle to hide already-linked items, select items to import/link
- [x] **ADO-07**: User can unlink a task from ADO — with option to keep or delete the local task
- [x] **ADO-08**: Direct ADO REST API calls from Go using token from abstracted provider — no shelling out to az cli per query
- [x] **ADO-09**: User can configure multiple ADO orgs, then pick specific projects within each org to follow
- [x] **ADO-10**: Items from all configured org/project pairs appear in a unified list with org/project label; toggleable group-by-project view

### External Links & Context

- [x] **LINK-01**: User can attach external URLs to any task (ICMs, dashboards, wiki pages, etc.) as structured links with label
- [x] **LINK-02**: Auto-detect known URL patterns (ICM → incident icon, Grafana → dashboard icon, ADO → work item icon, Wiki) with fallback to manual label
- [x] **LINK-03**: Links displayed in a dedicated section on task detail — clickable, labeled, with type icon

### Comments & Sync

- [x] **CMT-01**: User can add local comments to any task (private, never pushed to ADO by default)
- [x] **CMT-02**: User can selectively push a comment to linked ADO work item (with confirmation); synced comments clearly marked as public
- [x] **CMT-03**: User can update task description locally, then confirm-push to ADO (preview diff)

### ADO Sync

- [x] **SYNC-01**: Background auto-sync pulls ADO changes to linked items silently on configurable timer (default 15 min) + manual refresh trigger
- [x] **SYNC-02**: All outbound pushes to ADO require user confirmation via preview diff — never auto-push
- [x] **SYNC-03**: Per-field conflict resolution when both local and ADO sides changed the same linked item — user picks per field
- [x] **SYNC-04**: Linked task syncs title, status, and description to ADO; subtasks and personal breakdowns never pushed unless individually linked

### Dashboard & Views

- [x] **DASH-01**: Unified list view showing all tasks (personal + public/ADO-linked) with visual indicator (icon/badge/color) distinguishing personal from public
- [x] **DASH-02**: List view filterable by status, priority, project, due date, and ADO link status; sortable by these dimensions
- [x] **DASH-03**: Linked items show the connection (e.g., task → ADO work item with status badge)

### Task List

- [x] **TL-01**: Single global task list across all projects; filter by project to scope down
- [x] **TL-02**: Medium-density rows: checkbox + title + priority badge + project tag + ADO badge + due date + description preview (1 line) + subtask progress bar
- [x] **TL-03**: Flat sorted by default (priority then due date). Optional group-by toggle: status, priority, or project
- [x] **TL-04**: Click task → slide-out right panel (TaskDetail sidebar)
- [x] **TL-05**: Tasks can exist without a project (orphan/inbox — unassigned tasks)

### Projects

- [x] **PROJ-01**: Projects page — card grid layout. Each card: name, ADO badge (filled if linked to scenario/deliverable, hollow if local), dual progress (ADO children % + local tasks %), task count
- [x] **PROJ-02**: Pin/star favorite projects to top of the grid
- [x] **PROJ-03**: Projects are flat (no sub-projects). Hierarchy lives at the task level (subtasks/dependencies)
- [x] **PROJ-04**: Projects can be local-only groupings or linked to ADO scenarios/deliverables (same personal/public model as tasks)
- [x] **PROJ-05**: Click project card → project dashboard: stats, ADO progress, linked ADO context, then tasks filtered to that project
- [x] **PROJ-06**: Project progress shows two indicators: ADO children completion (from linked ADO item) and local task completion (personal tasks under project) — kept separate, never mixed
- [x] **PROJ-07**: User can link a project to an ADO scenario or deliverable; user can unlink with same keep/delete option as tasks

### UX Patterns

- [ ] **UX-01**: Clicking an ADO work item opens an in-app detail panel (relevant fields + link/import action). "Open in ADO" button deep-links to real browser.
- [x] **UX-02**: Clicking external links (ICMs, Grafana, wikis) opens in real browser — never in-app webview
- [x] **UX-03**: Tabs + sync button + filter chips all in same toolbar row (compact, no wasted space)
- [ ] **UX-04**: ADO browser shows tree view (Scenario → Deliverable → Task/Bug/Story) — not a flat list
- [ ] **UX-05**: ADO browser supports filter chips (status, type, area) + text search + saved ADO query picker (queries created in ADO, browsed/executed in app)
- [ ] **UX-06**: PR view scoped to: user's authored PRs + PRs user is required/reviewing. Excludes abandoned. No "all team PRs" flooding.
- [ ] **UX-07**: Pipeline runs show proper pipeline names (not just IDs), readable status

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Pull Requests

- **PR-01**: User can see all their own PRs with current status (draft, active, completed, abandoned)
- **PR-02**: User can see all PRs assigned to them for review with status
- **PR-03**: User can view PR details (title, repo, source/target branch, reviewers, votes)
- **PR-04**: PRs displayed under linked tasks in detail panel

### Pipeline Monitoring

- **PIPE-01**: User can see recent pipeline runs with status (succeeded, failed, running, queued, cancelled)
- **PIPE-02**: Pipeline status shown on linked PR or task

### PR Comments & Deployment

- **PRC-01**: User can view PR comment threads in-app
- **PRC-02**: User can respond to PR comments from the dashboard
- **DEP-01**: User can track PR deployment lifecycle (pending → staging → prod)

### Search & Filtering

- **SRCH-01**: User can full-text search across tasks, work items, and PRs
- **SRCH-02**: User can filter by status, priority, area, and assignment
- **SRCH-03**: Search results grouped by type (faceted)
- **SRCH-04**: User can save filter configurations as named views

### Notifications

- **NOTF-01**: User receives in-app notifications for PR review requests
- **NOTF-02**: User receives notifications for ADO work item assignment changes
- **NOTF-03**: User can configure notification preferences

### Team Views

- **TEAM-01**: User can see team-wide task and PR metrics
- **TEAM-02**: User can see blockers across team members
- **TEAM-03**: Manager view showing team workload distribution

### Daily Log

- **LOG-01**: User can add daily log entries (standup, progress, blocker, decision)
- **LOG-02**: User can view their log history

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| Kanban boards / sprint planning | ADO does this well — we layer on top, not replace |
| CI/CD pipeline management | ADO Pipelines handles this; show status only (v2) |
| Chat / messaging | Teams/Slack exists; don't compete |
| Wiki / document editing | Link to existing wikis instead |
| Time tracking | Adds friction; not core to unified view |
| Custom dashboard widgets / plugin system | Premature abstraction; build the right dashboard first |
| AI-powered insights | Complex; defer to future milestone |
| GitHub integration | ADO is primary; GitHub deferred to later milestone |
| VS Code extension | Planned for code review flow; not v1 |
| Local code review (reviewme) | Depends on VS Code extension; deferred |
| Linkwarden integration | Personal power feature; not team dependency |
| Mobile app | Web-first; mobile deferred |
| Personal/life tasks | Work tasks only (dev + non-dev like meetings, docs, reviews) |
| PR monitoring & review | Deferred to v2 — focus on task lifecycle + ADO sync first |
| Tag-based filtering | Tags exist on tasks but not a filter dimension for v1 |
| Real-time push on local change | Outbound sync is always manual/confirmed, never auto-push |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| AUTH-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| AUTH-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| AUTH-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| TASK-01 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-02 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-03 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-04 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-05 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-06 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-07 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-08 | Phase 2: ADO Integration & Sync Workflow | Complete |
| TASK-09 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-04 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-05 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-06 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-07 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-08 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-09 | Phase 2: ADO Integration & Sync Workflow | Complete |
| ADO-10 | Phase 2: ADO Integration & Sync Workflow | Complete |
| SYNC-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| SYNC-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| SYNC-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| SYNC-04 | Phase 2: ADO Integration & Sync Workflow | Complete |
| LINK-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| LINK-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| LINK-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| CMT-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| CMT-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| CMT-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| DASH-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| DASH-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| DASH-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| TL-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| TL-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| TL-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| TL-04 | Phase 2: ADO Integration & Sync Workflow | Complete |
| TL-05 | Phase 2: ADO Integration & Sync Workflow | Complete |
| PROJ-01 | Phase 2: ADO Integration & Sync Workflow | Complete |
| PROJ-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| PROJ-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| PROJ-04 | Phase 2: ADO Integration & Sync Workflow | Complete |
| PROJ-05 | Phase 2: ADO Integration & Sync Workflow | Complete |
| PROJ-06 | Phase 2: ADO Integration & Sync Workflow | Complete |
| PROJ-07 | Phase 2: ADO Integration & Sync Workflow | Complete |
| UX-01 | Phase 2: ADO Integration & Sync Workflow | Pending |
| UX-02 | Phase 2: ADO Integration & Sync Workflow | Complete |
| UX-03 | Phase 2: ADO Integration & Sync Workflow | Complete |
| UX-04 | Phase 2: ADO Integration & Sync Workflow | Pending |
| UX-05 | Phase 2: ADO Integration & Sync Workflow | Pending |
| UX-06 | Phase 3: PR Monitoring & Team Views | Pending |
| UX-07 | Phase 3: PR Monitoring & Team Views | Pending |

**Coverage:**
- v1 requirements: 54 total
- Mapped to phases: 54 ✓
- Unmapped: 0

---
*Requirements defined: 2026-03-31*
*Last updated: 2026-04-03 — added multi-org (ADO-09/10), external links (LINK-01/02/03), comments sync (CMT-01/02/03)*
