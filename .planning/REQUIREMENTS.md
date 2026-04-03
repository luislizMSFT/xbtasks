# Requirements: Team ADO Tool

**Defined:** 2026-03-31
**Core Value:** One pane of glass for all your work — personal tasks linked to ADO items, PRs, and comments — so nobody has to context-switch between tools.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Authentication

- [ ] **AUTH-01**: App authenticates to ADO via abstracted token provider — az cli `get-access-token` initially, swappable for PAT or OAuth later
- [ ] **AUTH-02**: Token auto-refreshes transparently; user doesn't need to re-auth unless az cli session expires
- [ ] **AUTH-03**: Token provider abstraction allows future swap to Entra ID OAuth2 PKCE or PAT without changing consuming code

### Personal Tasks

- [x] **TASK-01**: User can create a personal task with title, description, priority (P0-P3), and category
- [x] **TASK-02**: User can edit and delete their own tasks
- [x] **TASK-03**: User can set task status (todo, in_progress, in_review, done, blocked, cancelled)
- [x] **TASK-04**: User can create subtasks under a parent task or ADO-linked item
- [x] **TASK-05**: User can set a personal priority overlay independent of ADO priority
- [x] **TASK-06**: User can define dependencies between tasks (task A blocks task B)
- [x] **TASK-07**: User can add tags to tasks for organization
- [ ] **TASK-08**: Tasks start as personal (local-only). Linking/promoting to ADO makes them "public" (synced). Subtasks of public tasks stay personal unless individually linked.
- [ ] **TASK-09**: Quick-add task with just a title for fast capture; expand with full details (description, priority, project, due date) later

### ADO Integration

- [ ] **ADO-01**: User can view all ADO work items assigned to them in the dashboard
- [ ] **ADO-02**: User can view ADO work item details (state, priority, type, assigned_to, description, area path)
- [ ] **ADO-03**: User can link a personal task to an ADO work item (bidirectional)
- [ ] **ADO-04**: User can promote a personal task to a new ADO work item (with confirmation)
- [ ] **ADO-05**: User can import an ADO work item as a personal task with local tracking
- [ ] **ADO-06**: ADO browser view — browse all assigned ADO items, show which are already linked to local tasks, toggle to hide already-linked items, select items to import/link
- [ ] **ADO-07**: All outbound changes to ADO (push, promote, update) require a preview diff showing what will change, with user confirmation before applying
- [ ] **ADO-08**: Direct ADO REST API calls from Go using token from abstracted provider — no shelling out to az cli per query
- [ ] **ADO-09**: User can configure multiple ADO orgs, then pick specific projects within each org to follow
- [ ] **ADO-10**: Items from all configured org/project pairs appear in a unified list with org/project label; toggleable group-by-project view

### External Links & Context

- [ ] **LINK-01**: User can attach external URLs to any task (ICMs, dashboards, wiki pages, etc.) as structured links with label
- [ ] **LINK-02**: Auto-detect known URL patterns (ICM → incident icon, Grafana → dashboard icon, ADO → work item icon, Wiki) with fallback to manual label
- [ ] **LINK-03**: Links displayed in a dedicated section on task detail — clickable, labeled, with type icon

### Comments & Sync

- [ ] **CMT-01**: User can add local comments to any task (private, never pushed to ADO by default)
- [ ] **CMT-02**: User can selectively push a comment to linked ADO work item (with confirmation); synced comments clearly marked as public
- [ ] **CMT-03**: User can update task description locally, then confirm-push to ADO (preview diff)

### ADO Sync

- [ ] **SYNC-01**: Background auto-sync pulls ADO changes to linked items silently on configurable timer (default 15 min) + manual refresh trigger
- [ ] **SYNC-02**: All outbound pushes to ADO require user confirmation via preview diff — never auto-push
- [ ] **SYNC-03**: Per-field conflict resolution when both local and ADO sides changed the same linked item — user picks per field
- [ ] **SYNC-04**: Linked task syncs title, status, and description to ADO; subtasks and personal breakdowns never pushed unless individually linked

### Dashboard & Views

- [ ] **DASH-01**: Unified list view showing all tasks (personal + public/ADO-linked) with visual indicator (icon/badge/color) distinguishing personal from public
- [ ] **DASH-02**: List view filterable by status, priority, project, due date, and ADO link status; sortable by these dimensions
- [ ] **DASH-03**: Linked items show the connection (e.g., task → ADO work item with status badge)

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
| AUTH-01 | Phase 2: ADO Integration & Sync Workflow | Pending |
| AUTH-02 | Phase 2: ADO Integration & Sync Workflow | Pending |
| AUTH-03 | Phase 2: ADO Integration & Sync Workflow | Pending |
| TASK-01 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-02 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-03 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-04 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-05 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-06 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-07 | Phase 1: Foundation, Auth & Personal Tasks | Complete |
| TASK-08 | Phase 2: ADO Integration & Sync Workflow | Pending |
| TASK-09 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-01 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-02 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-03 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-04 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-05 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-06 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-07 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-08 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-09 | Phase 2: ADO Integration & Sync Workflow | Pending |
| ADO-10 | Phase 2: ADO Integration & Sync Workflow | Pending |
| SYNC-01 | Phase 2: ADO Integration & Sync Workflow | Pending |
| SYNC-02 | Phase 2: ADO Integration & Sync Workflow | Pending |
| SYNC-03 | Phase 2: ADO Integration & Sync Workflow | Pending |
| SYNC-04 | Phase 2: ADO Integration & Sync Workflow | Pending |
| LINK-01 | Phase 2: ADO Integration & Sync Workflow | Pending |
| LINK-02 | Phase 2: ADO Integration & Sync Workflow | Pending |
| LINK-03 | Phase 2: ADO Integration & Sync Workflow | Pending |
| CMT-01 | Phase 2: ADO Integration & Sync Workflow | Pending |
| CMT-02 | Phase 2: ADO Integration & Sync Workflow | Pending |
| CMT-03 | Phase 2: ADO Integration & Sync Workflow | Pending |
| DASH-01 | Phase 2: ADO Integration & Sync Workflow | Pending |
| DASH-02 | Phase 2: ADO Integration & Sync Workflow | Pending |
| DASH-03 | Phase 2: ADO Integration & Sync Workflow | Pending |

**Coverage:**
- v1 requirements: 35 total
- Mapped to phases: 35 ✓
- Unmapped: 0

---
*Requirements defined: 2026-03-31*
*Last updated: 2026-04-03 — added multi-org (ADO-09/10), external links (LINK-01/02/03), comments sync (CMT-01/02/03)*
