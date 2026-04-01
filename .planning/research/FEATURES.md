# Feature Research

**Domain:** Team Developer Productivity Dashboard (ADO-integrated)
**Researched:** 2025-07-14
**Confidence:** HIGH (well-understood domain; based on analysis of Linear, Shortcut, GitHub Dashboard, Backstage, Port, Cortex, Sleuth, ADO native features, and the proven xl TUI data model)

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist from day one. Missing any of these and the tool feels broken or not worth switching to.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Unified "My Work" view** | Every tool (Linear, GitHub Dashboard, ADO) has a personal work queue. Users expect to see all assigned items in one list. | MEDIUM | Core value prop. Must aggregate ADO work items, authored PRs, pending reviews, and personal tasks into a single sorted/filterable view. |
| **ADO work items display** | The tool layers on ADO; if it can't show ADO items clearly, it's useless. | MEDIUM | Title, state, priority, assignment, area path, iteration, type (Bug/Task/UserStory). Must handle ADO's data model faithfully. |
| **PR list with review status** | Every dev dashboard shows your PRs. GitHub Dashboard, Linear integrations, ADO itself all do this. | MEDIUM | Authored PRs + review-requested PRs. Show title, status (active/completed/abandoned), reviewers, vote status (approved/rejected/waiting), comment count. |
| **Personal task CRUD** | Core differentiator requires it; basic task management is also table stakes for any productivity tool. | LOW | Create, edit, delete, mark complete. P0-P3 priority. Statuses: todo, in_progress, done, blocked (proven from xl). Title, description, category. |
| **Microsoft/Entra ID auth** | Team is on Microsoft ecosystem. No other auth option is acceptable. | MEDIUM | SSO via Entra ID. Use MSAL.js for browser auth. Must obtain ADO API tokens transparently. No separate PAT management for end users. |
| **Search** | Every tool has search. Linear is famous for fast search. GitHub has powerful search. Users expect to find anything. | MEDIUM | Full-text across all entity types: tasks, work items, PRs. Must be fast (<200ms perceived). FTS on local data; WIQL for ADO server-side. |
| **Filtering and sorting** | Standard in Linear (custom views), Shortcut (filters), ADO (queries). Users expect to slice data. | MEDIUM | By status, priority, type, assignee, iteration/sprint, date range. Combinable filters. Saveable filter presets are nice-to-have but not v1. |
| **Deep links to ADO** | Users will need to jump to ADO for advanced operations (editing fields, running queries, sprint planning). | LOW | Every ADO-sourced item clickable → opens in ADO web. URL format is deterministic from org/project/id. |
| **Responsive web UI** | It's a web app. Must not feel sluggish or broken. | MEDIUM | Fast initial load, smooth interactions. Keyboard navigable. Not "enterprise software" feel — more Linear than ADO. |
| **Data freshness** | Stale data destroys trust. If I complete a work item in ADO and it still shows "Active" here, I stop using the tool. | MEDIUM | Poll-based sync on reasonable interval (2-5 min) + manual refresh. Webhooks ideal but harder in internal ADO. Visual "last synced" indicator. |
| **Work item status visualization** | Kanban/status columns are universal (ADO, Linear, Shortcut, Trello). Users expect visual status indication. | LOW | Color-coded status badges/pills at minimum. Optional board view later. List view with clear status indicators is sufficient for v1. |
| **Keyboard shortcuts** | Linear set the bar. Power users (your target audience: engineers) expect keyboard-driven navigation. | LOW | j/k navigation, Enter to open, Esc to close, / to search, quick-action shortcuts. Can be incremental but some must exist at launch. |

### Differentiators (Competitive Advantage)

Features that no single tool in the ecosystem offers in this combination. These are why the team would use this over raw ADO + personal tools.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Bidirectional task ↔ ADO linking** | THE key differentiator. No tool lets you promote a local task to an ADO work item OR pull an ADO item as a personal task with personal metadata. Linear/Shortcut are separate systems; ADO has no "personal staging" concept. | HIGH | Promote: create ADO work item from local task, set area/iteration, link back. Pull: import ADO item as local task with personal priority/notes. Sync: keep status aligned bidirectionally. xl has proven patterns for this. |
| **Personal priority/notes overlay** | Set your own P0-P3 priority and free-text notes on ANY item (ADO work items, PRs) without modifying the source. Your personal lens on organizational data. | LOW | Store as local metadata keyed by ADO ID. Display alongside ADO's native priority. Decouples personal workflow from organizational data hygiene. |
| **"Staging area" workflow** | Tasks start local (zero-friction capture), mature, then get promoted to ADO when ready. Reduces ADO noise from half-baked items. This is the anti-pattern to "everything must be in the tracker immediately." | MEDIUM | Status progression: local_only → ready_to_push → synced. Clear visual distinction between local-only and ADO-linked items. Batch promote capability. |
| **Cross-entity unified timeline** | "What happened today?" across all entity types — work item state changes, PR updates, review completions, task completions — in a single chronological feed. No tool does this across personal + ADO data. | MEDIUM | Aggregate activity from ADO API (work item updates, PR events) + local task changes. Chronological or grouped by entity. Powers standup prep and weekly reviews. |
| **PR comment thread preview** | View ADO PR comment threads inline without leaving the dashboard. ADO's PR UI is slow and context-switching is painful. | MEDIUM | Fetch PR threads via ADO API. Display threaded comments with author, date, status (active/resolved/won't fix). Don't need full reply capability v1 — read-only is the value. |
| **Sprint context awareness** | Auto-detect current iteration from ADO. Group and highlight items in current sprint vs backlog. Show sprint progress. | LOW | Query current iteration from ADO API. Tag items accordingly. Dashboard default view = current sprint. Simple progress bar. |
| **Team workload visibility** | Manager/lead view: see each team member's active items, PR review load, blocked items. No individual tool surface shows this across personal + organizational data. | MEDIUM | Query ADO by team members. Aggregate per-person. Show: items in progress, PRs awaiting review, blocked items. Privacy-respecting — personal tasks only visible if user opts in. |
| **Daily log / standup prep** | Auto-generate "yesterday/today" from activity data. xl already proves this is high-value (53 daily logs). Eliminates the "what did I do yesterday?" problem. | MEDIUM | Pull from: completed items, PR merges, review completions, status changes — all timestamped. Pre-fill standup template. Editable before sharing. |
| **Deliverable/scenario tree view** | ADO organizes work as Scenario → Deliverable → Task hierarchy. Show this tree with your personal items mapped in. See where your work fits in the big picture. | MEDIUM | Fetch ADO parent-child relationships. Render as collapsible tree. Map local tasks to deliverables via linking. Visual progress at each level. |
| **PR deployment tracking** | Track PR lifecycle beyond merge: deployed to dev? staging? prod? xl tracks this (proven value). Sleuth does this at company scale; this does it at personal scale. | MEDIUM | Track deployment status per PR. Integration with ADO release pipelines or manual status updates. Shows "is my change live?" without checking multiple systems. |
| **Subtask decomposition** | Break ADO work items into personal subtasks locally. "I need to do X, Y, Z for this work item" without creating ADO child items for every small step. | LOW | Parent-child relationship in local DB. Group under linked ADO item. Check off subtasks locally. Track completion percentage. |

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem obvious but would either bloat scope, fragment the value prop, or compete with tools that do them better.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| **Full ADO boards/backlog replication** | "If we're showing work items, let's show the board too" | ADO's board is the source of truth. Replicating it creates a split-brain problem and enormous scope. You'd be rebuilding ADO badly. | Deep link to ADO boards. Show items in list/filtered views. Don't do drag-and-drop kanban on ADO items — let ADO handle that. |
| **Sprint planning capabilities** | "We plan sprints, so the tool should help" | Sprint planning is a team ceremony with complex requirements (capacity, velocity, parent linking). ADO and Linear do this well. Building it means competing with mature tools. | Show sprint contents and progress. Let ADO handle sprint planning. Surface "unplanned items" from personal tasks as candidates to add. |
| **Notification/alert system** | "Tell me when something changes" | Notifications are a product in themselves. Email, in-app, push — each is complex. Teams/Outlook already deliver ADO notifications. Building another one adds noise. | Show "what changed since last visit" on dashboard load. Badge counts for new items. Don't build a notification infrastructure. |
| **Real-time collaboration** | "Let the team see each other's updates live" | WebSocket infrastructure, conflict resolution, presence indicators — enormous complexity for a team of 6-15. They're in the same room/Teams channel. | Refresh on navigate. "Last synced" timestamps. Team view refreshes on demand. Real-time is unnecessary at this team size. |
| **Custom workflow engine** | "Let each person define their own statuses/transitions" | Configuration complexity explodes. Every custom workflow is a maintenance burden. xl's fixed status model works. | Fixed status model (todo/in_progress/blocked/done + ADO state mapping). If someone needs custom workflows, ADO already supports them. |
| **Time tracking** | "Track how long tasks take" | Rabbit hole. Accurate time tracking requires discipline nobody has. Estimation vs actual is a whole product (Toggl, Harvest). Low value for team adoption. | Don't track time. If needed later, integrate with existing time tracking tools. |
| **Wiki/documentation system** | "We need docs near our tasks" | Team already has ADO Wiki. Building another doc system fragments knowledge. | Link to ADO Wiki pages from tasks. Don't store docs. |
| **Chat/messaging** | "Discuss work items in-app" | Teams exists. Building chat means competing with Teams + maintaining real-time infra. Zero chance of adoption when Teams is already open. | Link to Teams conversations. PR comments (read-only) covers the code discussion need. |
| **Mobile app** | "Check my tasks on my phone" | Tiny team, dev tool. Engineers work at their desks. Mobile adds two platforms of maintenance for near-zero usage. | Responsive web app works on mobile browsers in a pinch. Don't build native. |
| **GitHub integration** | "Some repos might move to GitHub" | ADO is the system. Adding GitHub doubles every integration surface. If the team migrates, that's a separate project. | Explicitly out of scope per PROJECT.md. Revisit only if team actually migrates. |
| **AI/LLM-powered features** | "Use AI to summarize PRs, suggest priorities" | Scope explosion. Every AI feature requires prompt engineering, testing, cost management. Team already has M365 Copilot + Agency CLI MCP. | Let the Agency CLI MCP handle natural language queries. The dashboard is a structured data surface, not an AI interface. Keep them separate. |
| **DORA metrics / engineering analytics** | "Track team velocity, lead time, deploy frequency" | Different product (Sleuth, LinearB, Swarmia). Analytics requires historical data pipeline, visualizations, statistical reasoning. Enormous scope, different audience. | If metrics are wanted, integrate with Sleuth or similar. The dashboard shows current state, not historical trends. |
| **Custom dashboards / plugin system** | "Let teams build their own widgets like Backstage" | Premature abstraction for a team of 6-15. Plugin architecture is a platform investment; you'd be building a dashboard builder instead of a dashboard. | Build the right fixed dashboard. Add views when real needs emerge. Don't abstract until you have 3+ concrete use cases for extensibility. |

## Feature Dependencies

```
[Entra ID Auth]
    └──(foundational, no dependencies)

[ADO API Integration Layer]
    └──requires──> [Entra ID Auth]

[Database/Storage Layer]
    └──(foundational, no dependencies)

[ADO Work Items Display]
    └──requires──> [ADO API Integration Layer]

[PR List + Review Status]
    └──requires──> [ADO API Integration Layer]

[PR Comment Thread Preview]
    └──requires──> [PR List + Review Status]

[Personal Task CRUD]
    └──requires──> [Database/Storage Layer]

[Subtask Decomposition]
    └──requires──> [Personal Task CRUD]

[Bidirectional ADO Linking]
    └──requires──> [Personal Task CRUD]
    └──requires──> [ADO API Integration Layer]

[Personal Priority/Notes Overlay]
    └──requires──> [Database/Storage Layer]
    └──requires──> [ADO Work Items Display]

[Unified "My Work" View]
    └──requires──> [ADO Work Items Display]
    └──requires──> [PR List + Review Status]
    └──requires──> [Personal Task CRUD]

[Search]
    └──requires──> [Unified "My Work" View]

[Filtering/Sorting]
    └──requires──> [Unified "My Work" View]

["Staging Area" Workflow]
    └──requires──> [Bidirectional ADO Linking]

[Sprint Context Awareness]
    └──requires──> [ADO API Integration Layer]
    └──enhances──> [Unified "My Work" View]

[Cross-Entity Timeline]
    └──requires──> [ADO Work Items Display]
    └──requires──> [PR List + Review Status]
    └──requires──> [Personal Task CRUD]

[Daily Log / Standup Prep]
    └──requires──> [Cross-Entity Timeline]

[Team Workload Visibility]
    └──requires──> [ADO API Integration Layer]
    └──requires──> [Unified "My Work" View]

[Deliverable/Scenario Tree]
    └──requires──> [ADO Work Items Display]

[PR Deployment Tracking]
    └──requires──> [PR List + Review Status]
    └──requires──> [ADO API Integration Layer]
```

### Dependency Notes

- **Bidirectional ADO Linking requires both Personal Task CRUD and ADO API Integration:** You can't link things that don't exist on both sides. The linking logic is the bridge between two already-working systems.
- **Unified "My Work" View requires three data sources:** This is the integration layer. It can't exist until work items, PRs, and personal tasks each work independently.
- **Cross-Entity Timeline requires all three display features:** The timeline aggregates events from all entity types. Build the individual views first, then unify them temporally.
- **Daily Log requires Cross-Entity Timeline:** The standup prep is essentially a filtered timeline (yesterday's activity + today's planned). Don't build it as a separate data pipeline.
- **Team Workload Visibility builds on individual views:** It's the same data, queried across team members instead of just the current user. The per-user infrastructure must exist first.
- **Sprint Context Awareness enhances the unified view:** It's not a separate feature — it's a filter/grouping dimension applied to the existing view. Low incremental cost once the view exists.
- **"Staging Area" Workflow requires Bidirectional Linking:** The staging concept only makes sense if tasks can actually be promoted to ADO. Without linking, it's just local task management.

## MVP Definition

### Launch With (v1)

Minimum viable product — what one person (Luis, dogfooding) needs to validate the concept and demonstrate value to the team.

- [ ] **Entra ID authentication** — Non-negotiable for ADO API access. MSAL.js browser auth flow.
- [ ] **ADO work items display** — Show my assigned items with state, priority, type, iteration. This alone proves ADO integration works.
- [ ] **PR list with review status** — My authored PRs + review-requested PRs. Vote status, comment count. Second highest-traffic view.
- [ ] **Personal task CRUD** — Create/edit/complete/delete local tasks. P0-P3 priority, 4 statuses (todo/in_progress/done/blocked).
- [ ] **Unified "My Work" view** — Single page combining the above three. Sortable, filterable by type and status. THE product.
- [ ] **Bidirectional ADO linking (promote flow)** — Promote a local task to an ADO work item. This is the key differentiator; must be in v1 to test the hypothesis.
- [ ] **Deep links to ADO** — Every ADO-sourced item links back to ADO web UI.
- [ ] **Keyboard shortcuts (basic)** — j/k navigation, Enter to open, / to search. Sets the tone for a dev-centric tool.
- [ ] **Manual refresh + auto-poll** — Refresh button + 5-minute auto-poll for ADO data freshness.

### Add After Validation (v1.x)

Features to add once the core loop (view work → manage tasks → link to ADO) is proven with real usage.

- [ ] **Bidirectional ADO linking (pull flow)** — Pull ADO items as personal tasks with personal metadata overlay. Trigger: v1 users want to annotate ADO items personally.
- [ ] **Personal priority/notes overlay** — Add personal P0-P3 + notes to any ADO item. Trigger: users want to re-prioritize ADO items for their own workflow.
- [ ] **PR comment thread preview** — Read comment threads inline. Trigger: users keep clicking through to ADO to read PR comments.
- [ ] **Search (full-text)** — Search across all entity types. Trigger: more than ~30 items makes scrolling/filtering insufficient.
- [ ] **Sprint context awareness** — Auto-detect current iteration, group items. Trigger: users ask "what's in my sprint?" frequently.
- [ ] **Filtering presets / saved views** — Save common filter combinations. Trigger: users repeatedly apply the same filters.
- [ ] **"Staging area" visual workflow** — Clear local-only vs synced vs needs-sync visual states. Trigger: users have enough local tasks that sync state becomes confusing.
- [ ] **Subtask decomposition** — Break ADO work items into personal subtasks. Trigger: users need granular local tracking under a single ADO item.

### Future Consideration (v2+)

Features to defer until the tool has team adoption and proven individual value.

- [ ] **Team workload visibility** — Requires team adoption (>3 users) to be meaningful. Don't build for one user.
- [ ] **Daily log / standup prep** — High value but depends on rich activity data. Needs weeks of usage data to auto-generate meaningfully.
- [ ] **Cross-entity unified timeline** — Activity feed across all types. Needs mature data pipeline. Build after the entity views are solid.
- [ ] **Deliverable/scenario tree view** — ADO hierarchy visualization. Complex ADO API queries. Value increases with team adoption.
- [ ] **PR deployment tracking** — Requires ADO release pipeline integration. Complex infrastructure. High value for "is my change live?" but not core to the dashboard.
- [ ] **VS Code extension** — Different surface, different codebase, different distribution. Build after web proves value.
- [ ] **Bulk operations** — Multi-select and batch promote/update. Needed at scale, not at 10 items.
- [ ] **Incident dashboard** — Active incidents with severity. Valuable for on-call engineers. xl tracks sev1-4. Defer until team adoption provides the data density.

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Entra ID auth | HIGH | MEDIUM | P1 |
| ADO work items display | HIGH | MEDIUM | P1 |
| PR list + review status | HIGH | MEDIUM | P1 |
| Personal task CRUD | HIGH | LOW | P1 |
| Unified "My Work" view | HIGH | MEDIUM | P1 |
| Bidirectional linking (promote) | HIGH | HIGH | P1 |
| Deep links to ADO | MEDIUM | LOW | P1 |
| Keyboard shortcuts (basic) | MEDIUM | LOW | P1 |
| Data freshness (poll + manual) | HIGH | LOW | P1 |
| Bidirectional linking (pull) | HIGH | MEDIUM | P2 |
| Personal priority/notes overlay | MEDIUM | LOW | P2 |
| PR comment thread preview | MEDIUM | MEDIUM | P2 |
| Search (full-text) | MEDIUM | MEDIUM | P2 |
| Sprint context awareness | MEDIUM | LOW | P2 |
| Subtask decomposition | MEDIUM | LOW | P2 |
| Filtering presets / saved views | LOW | LOW | P2 |
| Staging area workflow visuals | MEDIUM | LOW | P2 |
| Team workload visibility | HIGH | MEDIUM | P3 |
| Daily log / standup prep | MEDIUM | MEDIUM | P3 |
| Cross-entity timeline | MEDIUM | HIGH | P3 |
| Deliverable/scenario tree | MEDIUM | HIGH | P3 |
| PR deployment tracking | MEDIUM | HIGH | P3 |
| Incident dashboard | MEDIUM | MEDIUM | P3 |

**Priority key:**
- P1: Must have for launch (dogfooding with Luis)
- P2: Should have, add once core loop validated
- P3: Nice to have, requires team adoption to justify

## Competitor Feature Analysis

| Feature | Linear | GitHub Dashboard | ADO Native | Backstage/Port/Cortex | Sleuth | Our Approach |
|---------|--------|-----------------|------------|----------------------|--------|--------------|
| Personal work queue | ✅ My Issues | ✅ Dashboard feed | ✅ "Assigned to me" query | ❌ Not their focus | ❌ Not their focus | ✅ Unified across all entity types — this is the core |
| Task creation (low friction) | ✅ Cmd+K, fast | ❌ Issues only | ⚠️ Work item forms are heavy | ❌ | ❌ | ✅ Personal tasks with zero ADO overhead, promote when ready |
| PR visibility | ⚠️ Via Git integration | ✅ Native | ✅ Native but separate from boards | ❌ | ✅ Deploy tracking | ✅ PRs alongside work items and tasks in one view |
| Bidirectional external linking | ❌ One-way integrations | ❌ | ❌ | ❌ | ❌ | ✅ KEY DIFFERENTIATOR. Promote local → ADO, pull ADO → local |
| Personal metadata overlay | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ Personal priority + notes on any item without touching source |
| Comment threads inline | ✅ Native | ✅ Native | ✅ Native (but slow UI) | ❌ | ❌ | ✅ Fast read of ADO PR threads without ADO's heavy UI |
| Sprint awareness | ✅ Cycles | ❌ Milestones only | ✅ Native | ❌ | ❌ | ✅ Auto-detect current iteration, group items |
| Team visibility | ✅ Team views | ⚠️ Org-level only | ✅ Team queries | ✅ Service catalog | ❌ | ✅ Per-person workload across all entity types |
| Keyboard-first UX | ✅ Best in class | ⚠️ Some shortcuts | ❌ Mouse-heavy | ❌ | ❌ | ✅ Must match Linear's bar for dev-centric feel |
| Standup/daily log | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ Auto-generated from activity data. Nobody does this well. |
| Service catalog | ❌ | ❌ | ❌ | ✅ Core feature | ✅ | ❌ Not our domain. Out of scope. |
| DORA metrics | ❌ | ❌ | ⚠️ Via Analytics | ✅ Scorecards | ✅ Core feature | ❌ Not our domain. Out of scope. |
| Deployment tracking | ❌ | ✅ Environments | ✅ Release pipelines | ✅ | ✅ Core feature | ⚠️ Future (v2). Track PR → deploy lifecycle. |

### Key Competitive Insight

**No tool in the ecosystem bridges the gap between personal task management and organizational work tracking.** Linear/Shortcut are standalone project management. ADO is organizational. GitHub Dashboard is repository-centric. Backstage/Port/Cortex are service catalogs. Sleuth is metrics.

The unique position: **a personal productivity layer that speaks ADO natively.** The closest analogue is a developer building their own dashboard (which is what xl already is — and now it becomes a team tool).

The two features that nobody else has:
1. **Bidirectional promote/pull** — tasks flow between personal and organizational
2. **Personal metadata overlay** — your priority on top of the org's priority

Everything else is table stakes done well in a unified view. The combination is the moat.

## Sources

- **Linear** — Known from extensive usage and documentation. Keyboard-first project management, cycles, custom views, Git integrations. Training data, HIGH confidence in feature set.
- **Shortcut (formerly Clubhouse)** — Stories/epics model, flexible workflows. Training data, HIGH confidence.
- **GitHub Dashboard** — Personal feed of PRs, issues, review requests across repos. Training data, HIGH confidence.
- **Azure DevOps** — Boards, backlogs, queries, sprint planning, PR UI. Direct team usage, HIGH confidence.
- **Backstage (Spotify)** — Open-source IDP, service catalog, TechDocs, scaffolding. Training data, HIGH confidence.
- **Port / Cortex** — Commercial IDPs, service catalogs, scorecards, self-service actions. Training data, MEDIUM confidence (feature sets may have evolved).
- **Sleuth** — DORA metrics, deployment tracking, change lead time. Training data, MEDIUM confidence.
- **xl TUI** — Existing personal tool with 17 entities, proven ADO integration, 29 tasks / 10 PRs / 53 daily logs / 158 memory entries. Direct project context, HIGH confidence.

---
*Feature research for: Team Developer Productivity Dashboard (ADO-integrated)*
*Researched: 2025-07-14*
