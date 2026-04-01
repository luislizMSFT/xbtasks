# Requirements: Team ADO Tool

**Defined:** 2026-03-31
**Core Value:** One pane of glass for all your work — personal tasks linked to ADO items, PRs, and comments — so nobody has to context-switch between tools.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Authentication

- [ ] **AUTH-01**: User can sign in with Microsoft Entra ID (SSO with existing work account)
- [ ] **AUTH-02**: User session persists across browser refresh (token auto-refresh)
- [ ] **AUTH-03**: User profile populated from Entra ID (display name, email, avatar)

### Personal Tasks

- [ ] **TASK-01**: User can create a personal task with title, description, priority (P0-P3), and category
- [ ] **TASK-02**: User can edit and delete their own tasks
- [ ] **TASK-03**: User can set task status (todo, in_progress, in_review, done, blocked, cancelled)
- [ ] **TASK-04**: User can create subtasks under a parent task or ADO-linked item
- [ ] **TASK-05**: User can set a personal priority overlay independent of ADO priority
- [ ] **TASK-06**: User can define dependencies between tasks (task A blocks task B)
- [ ] **TASK-07**: User can add tags to tasks for organization

### ADO Integration

- [ ] **ADO-01**: User can view all ADO work items assigned to them in the dashboard
- [ ] **ADO-02**: User can view ADO work item details (state, priority, type, assigned_to, description, area path)
- [ ] **ADO-03**: User can link a personal task to an ADO work item (bidirectional)
- [ ] **ADO-04**: User can promote a personal task to a new ADO work item
- [ ] **ADO-05**: User can import an ADO work item as a personal task with local tracking

### Pull Requests

- [ ] **PR-01**: User can see all their own PRs with current status (draft, active, completed, abandoned)
- [ ] **PR-02**: User can see all PRs assigned to them for review with status
- [ ] **PR-03**: User can view PR details (title, repo, source/target branch, reviewers, votes)

### Dashboard

- [ ] **DASH-01**: User sees a unified dashboard on login showing tasks, ADO items, and PRs in one view
- [ ] **DASH-02**: Dashboard shows clear visual distinction between personal tasks and ADO-linked items
- [ ] **DASH-03**: Linked items show the connection (e.g., task → ADO work item with status badge)

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### ADO Sync

- **SYNC-01**: Task status changes automatically push to linked ADO work item
- **SYNC-02**: ADO work item changes automatically pull to linked local task
- **SYNC-03**: User can bulk import all assigned ADO items in one action
- **SYNC-04**: User can import all child work items under an ADO parent

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
| CI/CD pipeline management | ADO Pipelines handles this; show status only via PR |
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

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| AUTH-01 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| AUTH-02 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| AUTH-03 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| TASK-01 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| TASK-02 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| TASK-03 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| TASK-04 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| TASK-05 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| TASK-06 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| TASK-07 | Phase 1: Foundation, Auth & Personal Tasks | Pending |
| ADO-01 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| ADO-02 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| ADO-03 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| ADO-04 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| ADO-05 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| PR-01 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| PR-02 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| PR-03 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| DASH-01 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| DASH-02 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |
| DASH-03 | Phase 2: ADO Integration, PRs & Unified Dashboard | Pending |

**Coverage:**
- v1 requirements: 21 total
- Mapped to phases: 21 ✓
- Unmapped: 0

---
*Requirements defined: 2026-03-31*
*Last updated: 2026-03-31 after initial definition*
