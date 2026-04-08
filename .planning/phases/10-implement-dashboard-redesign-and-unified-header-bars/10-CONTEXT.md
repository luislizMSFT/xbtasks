# Phase 10: Implement dashboard redesign and unified header bars - Context

**Gathered:** 2026-04-08
**Status:** Ready for planning

<domain>
## Phase Boundary

Redesign the dashboard landing page and implement a unified header bar pattern across all pages. The dashboard gets an Attention Bar, richer task rows, and an Upcoming section. The top bar gets a consistent 3-zone layout (sync left, page content center, actions right) applied to every page. No new backend functionality — this is a frontend-only visual/layout overhaul.

</domain>

<decisions>
## Implementation Decisions

### Dashboard Layout
- Keep 2-column layout: tasks left (3/5 width), PRs+Pipelines right (2/5 width)
- Keep greeting + summary line at top ("Good morning, Luis" + "3 tasks in progress · 1 blocked")
- Keep inline compact stats line below greeting ("3 in progress · 1 blocked · 5 done of 47")
- Activity sidebar stays as the global AppShell panel (240px, toggleable via clock icon) — not dashboard-specific

### Dashboard Sections (Left Column)
- **Today's Focus** — tasks with status `in_progress` or `in_review` (kept from current)
- **Upcoming** — tasks due within 3 days (replaces "Recent Activity" section). Overdue tasks (past due date) highlighted with red indicator
- **Blocked** — tasks with status `blocked`, shown with red left border accent + blocked reason text (from playground)

### Attention Bar
- Horizontal scrollable row of urgency nudges between greeting and stats line
- Three nudge types:
  - **Due soon**: amber background, CalendarDays icon, "3 due within 3 days" with first 2 task title previews
  - **Pipeline failure**: red background, XCircle icon, pipeline name + branch name
  - **PR approval ready**: green background, GitMerge icon, "PR #N has N approvals — ready to merge"
- Each nudge is a pill-shaped card with icon + text, horizontally scrollable on overflow
- Only shown when relevant (conditional rendering per nudge type)

### Dashboard Task Rows (Richer)
- Enhanced rows in Today's Focus and Upcoming sections (not reusing existing TaskRow component):
  - Priority dot (color: P0=red, P1=orange, P2=amber, P3=zinc)
  - ADO type icon (from `adoTypeIcon()`) or SquareCheckBig for personal tasks
  - Task title (truncated)
  - "personal" badge for personal tasks (small text, primary/8 bg)
  - Pending sync amber dot indicator
  - Status badge (blue for active, violet for in_review)
  - Due date (amber highlight if within 2 days)
- Blocked section shows blocked reason text below task title (italic, red-tinted)

### Right Column (PRs + Pipelines)
- Keep current 3-section structure: Needs Your Review, Your Pull Requests, Pipelines
- No structural changes to right column — matches current production design

### Unified Header Bar Pattern
- Apply to all 6 pages: Dashboard, Tasks, ADO, Projects, Dependencies, Settings
- Modify AppShell top bar to support the 3-zone pattern:
  - **Left zone**: Page name + vertical divider + sync cluster
    - Sync cluster: colored dot (green=connected, red=offline) + "Synced"/"Offline" label + relative time ("3m ago") + pending changes badge (amber, "N pending") + refresh button
    - Sync cluster is consistent and identical on every page
  - **Center zone**: Page-specific content (teleported from each page)
    - Dashboard: stat badges ("5 active", "2 blocked", "31/47 done")
    - Tasks: status filter chips (All, Active, Blocked, Done)
    - ADO: tabs (Browser, PRs, Pipelines) with count badges
    - Projects: project count + active badge ("4 projects", "2 active")
    - Dependencies: graph stats ("8 nodes · 5 edges") + cycle badge
    - Settings: empty (no center content)
  - **Right zone**: Search ⌘K + New dropdown + Activity toggle (always the same, unchanged)

### Claude's Discretion
- Exact implementation of sync cluster component (shared component vs inline)
- Whether to extract the Attention Bar nudges into separate sub-components
- Exact responsive breakpoints for Attention Bar horizontal scroll
- Whether dashboard task rows should be a new DashboardTaskRow component or inline in DashboardView
- Animation/transition details for Attention Bar nudges appearing/disappearing

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Playground Prototypes (approved designs)
- `frontend/src/views/playground/PlaygroundDashboardHeader.vue` — Proposed unified header pattern for all 6 pages (current vs proposed side-by-side). Contains exact sync cluster markup, center zone content per page, and design notes.
- `frontend/src/views/playground/PlaygroundDashboard.vue` — Proposed dashboard layout with Attention Bar, richer task rows, sync cluster, and 2-column layout. Contains mock data and exact markup to implement.

### Current Production Code (to be modified)
- `frontend/src/layouts/AppShell.vue` — Current top bar implementation (46px, breadcrumb + teleport target + search + new + activity). Must be refactored for 3-zone pattern.
- `frontend/src/views/DashboardView.vue` — Current dashboard implementation to be redesigned. Has Teleport to #topbar-actions.
- `frontend/src/components/PageHeader.vue` — Sub-bar component with left/right slots. May need updating or removal if header unification absorbs its role.
- `frontend/src/components/Sidebar.vue` — Sidebar nav (56px icon rail). Not changing structure, but playground nav items need review.

### Style & Component References
- `frontend/src/lib/styles.ts` — adoTypeIcon, adoTypeColor, adoStateClasses, prStatusClasses mappings
- `frontend/src/components/ui/StatusBadge.vue` — Current status badge component
- `frontend/src/components/tasks/TaskRow.vue` — Current task row component (dashboard will use richer rows instead)

### Stores (data sources for dashboard)
- `frontend/src/stores/tasks.ts` — Task store with stats computed, fetchTasks
- `frontend/src/stores/prs.ts` — PR store with myPRs, reviewPRs
- `frontend/src/stores/ado.ts` — ADO store with pipelines, connected state
- `frontend/src/stores/sync.ts` — Sync store with syncing state, lastSyncedAt, manualSync()

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `AppShell.vue` top bar: Has `#topbar-actions` teleport target already used by DashboardView — teleport pattern can be extended for center zone
- `PageHeader.vue`: Simple left/right slot bar — may become redundant if unified header absorbs page-specific content into center zone
- `StatusBadge.vue`, `PriorityBadge.vue`: Existing badge components for status/priority display
- `adoTypeIcon()`, `adoTypeColor()` from `styles.ts`: Already used in playground dashboard rows
- `relativeTime()` from `lib/date.ts`: Used in DashboardView for timestamps
- `useSyncStore`: Already provides syncing, lastSyncedAt, manualSync, connected state
- `useADOStore`: Provides pipelines and connected state

### Established Patterns
- **Teleport for page-specific header content**: DashboardView already teleports status indicators to `#topbar-actions` — this pattern extends naturally to center zone content
- **Conditional section rendering**: Dashboard already conditionally shows blocked section — same pattern for Attention Bar nudges
- **shadcn-vue Badge component**: Used throughout for counts, status, labels
- **Lucide icons**: All icons from lucide-vue-next, no emojis

### Integration Points
- `AppShell.vue` top bar needs structural refactoring (3-zone layout)
- Each view needs to teleport its center zone content to the new center teleport target
- DashboardView needs major template rewrite (Attention Bar, Upcoming section, richer rows)
- Sync cluster may be a new shared component used by AppShell

</code_context>

<specifics>
## Specific Ideas

- Base the implementation on the approved PlaygroundDashboardHeader.vue and PlaygroundDashboard.vue designs — they contain exact markup, spacing, and styling to replicate
- The "Attention Bar" is the key new dashboard feature — horizontal urgency nudges that surface pipeline failures, due-soon tasks, and merge-ready PRs at a glance
- Dashboard task rows should show more info than current TaskRow: priority dot + ADO type icon + personal badge + sync dot + status badge + due date
- The playground's design notes section describes the pattern: "Left zone = Page name + sync cluster, Center zone = Per-page content, Right zone = Search/New/Activity"
- Blocked section should show the reason text (italic, red-tinted) below the task title — this is visible in the playground

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 10-implement-dashboard-redesign-and-unified-header-bars*
*Context gathered: 2026-04-08*
