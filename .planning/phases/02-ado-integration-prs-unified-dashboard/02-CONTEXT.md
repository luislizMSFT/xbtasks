# Phase 2: ADO Integration & Sync Workflow - Context

**Gathered:** 2026-04-03 (updated from 2026-04-02 discussion)
**Status:** Ready for planning

<domain>
## Phase Boundary

Authenticate to ADO, implement the personal→public task model, build the ADO browser view, and deliver safe bidirectional sync with user confirmation for all outbound changes. PRs are deferred to Phase 3.

**Core mental model:**
- **Personal tasks** — Local-only, fully private, never touch ADO. Can have subtasks and dependencies.
- **Public tasks** — Linked to ADO (via link, promote, or import). Title, status, and description sync bidirectionally. Show in same list but visually distinguished (badge/icon).
- **Transition** — A personal task becomes public when the user explicitly links/promotes it to ADO (with confirmation). Once public, it stays synced.
- **Subtasks of public tasks** — Remain personal unless individually linked. User can break down an ADO item locally without cluttering the ADO board.

Our domain: **Projects** and **Tasks** (with subtasks). ADO has its own hierarchy (Scenario → Deliverable → Task/Bug/Story) but we don't replicate that in our domain. Users **link** our projects/tasks to any ADO item they choose.

</domain>

<decisions>
## Implementation Decisions

### Authentication & ADO Client
- **D-01:** Abstracted token provider interface in Go — implementations: az cli (`az account get-access-token`), PAT (from keyring), future OAuth
- **D-02:** Token auto-refreshes transparently; cached with TTL, re-fetched when expired
- **D-03:** Direct ADO REST API calls from Go using token from provider — no shelling out to az cli per query
- **D-04:** ADO client package (`pkg/ado/`) — stateless HTTP client accepting token, handles pagination, rate limits, WIQL queries, JSON Patch
- **D-04a:** Single token assumed to work across all configured orgs (same Azure tenant). Multi-tenant support deferred.

### Multi-Org Configuration
- **D-04b:** Config supports a list of org/project pairs (replaces single `ado.organization` / `ado.project`)
- **D-04c:** User configures orgs first, then picks specific projects within each org to follow
- **D-04d:** All configured org/project pairs are synced — work items and ADO browser aggregate across all of them
- **D-04e:** Each item carries its org/project origin as metadata (shown as badge/label in UI)
- **D-04f:** Unified list is default; user can toggle group-by-project view (collapsible sections per org/project)

### Personal → Public Task Model
- **D-05:** Tasks have a computed "public" status derived from whether they have an entry in `task_ado_links` table
- **D-06:** Unified list view shows personal + public tasks together; public tasks have a filled ADO badge, personal have hollow/empty badge
- **D-07:** Quick-add: create task with just a title (all other fields optional, filled in later)
- **D-08:** Subtasks of a public (linked) task remain personal unless the user explicitly links them individually
- **D-09:** When promoting a personal task to ADO, only title/status/description are pushed — subtasks, personal priority, local notes stay local

### ADO Browser View
- **D-10:** Dedicated view showing all ADO work items assigned to the user (fetched via ADO REST API)
- **D-11:** Each item shows whether it's already linked to a local task (visual indicator)
- **D-12:** Toggle to hide already-linked items (so user only sees what's not imported yet)
- **D-13:** From this view, user can select items to import (create local task) or link (to existing local task)
- **D-14:** ADO tree hierarchy — expandable Scenario → Deliverable → Task/Bug/Story with type icons
- **D-15:** Tree scope: sync items assigned to me + their parent chain for context
- **D-16:** Tree search: instant search of local SQLite cache + "Search all ADO" button for live API query

### ADO Linking (3 Flows)
- **D-17:** **Link** — Connect local task to existing ADO item (search by ID or title)
- **D-18:** **Promote** — Create new ADO work item from local task (preview diff of what will be created, confirm)
- **D-19:** **Import** — Pull ADO work item into local task list (creates new local task with ADO link)
- **D-20:** Task row ADO icon: hollow/empty if unlinked, filled with ADO type icon if linked. Click hollow → opens link dialog

### Sync Behavior
- **D-21:** Background auto-sync pulls ADO changes to linked items silently on configurable timer (default 15 min via Viper) + manual refresh button
- **D-22:** All outbound pushes to ADO require user confirmation via preview diff — shows exactly what fields will change (title, status, description), user confirms or cancels
- **D-23:** Never auto-push to ADO. Every outbound change is explicit and confirmed.
- **D-24:** Linked fields that sync: title, status (with ADO state mapping), description
- **D-25:** Fields that never sync: subtasks, personal priority, local notes/comments, tags, external links

### External Links
- **D-25a:** Tasks have a `links` section — structured list of URLs with label and type
- **D-25b:** Auto-detect known URL patterns: `portal.microsofticm.com` → ICM incident, Grafana URLs → dashboard, ADO URLs → work item, Wiki URLs → wiki page. Fallback to user-provided label.
- **D-25c:** Links stored locally in SQLite — never pushed to ADO. They're personal context.
- **D-25d:** Task detail shows links section with type icon, label, and clickable URL

### Comments
- **D-25e:** Comments are local by default (private, never auto-pushed)
- **D-25f:** User can selectively push a comment to the linked ADO work item (requires confirmation)
- **D-25g:** Synced/public comments are visually marked (e.g., ADO icon, "synced" badge) so user always knows what's public vs private
- **D-25h:** Description field can be edited locally, then confirm-pushed to ADO (preview diff shows old vs new)

### Investigation Workflow
- **D-25i:** Task flow is bidirectional: user can start with a local investigation task, then link/promote to ADO when a bug is created. Or import an existing ADO bug and enrich locally.
- **D-25j:** A linked task becomes an "investigation hub" — ADO bug context + ICM links + dashboard links + local subtasks + private notes, all in one place

### Conflict Resolution
- **D-26:** When both local and ADO changed the same linked item, show conflict to user
- **D-27:** Per-field conflict resolution — user picks local or ADO value for each conflicting field
- **D-28:** Non-conflicting fields merge silently (e.g., only title changed locally, only status changed in ADO)

### Dashboard & Filters
- **D-29:** Unified list view (not kanban) as primary task view
- **D-30:** Filter dimensions: status, priority, project, due date, ADO link status (personal/public)
- **D-31:** Sortable by all filter dimensions
- **D-32:** No tag-based filtering for v1 (tags exist on tasks but not a filter dimension)

### Agent's Discretion
- ADO REST API client library design (package structure, error handling)
- Token caching strategy and TTL
- Preview diff UI component design
- Conflict resolution UI design
- ADO state mapping (todo→Proposed, in_progress→Active, done→Completed, blocked→Blocked)
- Bulk import UX from ADO browser view
- Search implementation (local cache vs live API)

### Deferred Ideas
- PR monitoring under tasks — Phase 3
- Pipeline status — Phase 3
- Customizable dashboard widget grid — v2
- Rich text WYSIWYG for comments — v2 (plain text for v1)
- Activity timeline per task (chronological events) — v2
- @mention syntax for quick linking — v2
- Bulk linking with smart auto-match — v2
- In-app PR diff viewer — if needed later
- Multi-tenant auth (different tokens per org) — if needed later
- Syncing external links to ADO — links stay local for now

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Context
- `.planning/PROJECT.md` — Full project vision, personal→public model, sync safety constraints, key decisions
- `.planning/REQUIREMENTS.md` — AUTH-01 through AUTH-03, TASK-08/09, ADO-01 through ADO-08, SYNC-01 through SYNC-04, DASH-01 through DASH-03
- `.planning/ROADMAP.md` — Phase 2 success criteria and scope boundary

### Phase 1 Context
- `.planning/phases/01-foundation-auth-personal-tasks/01-CONTEXT.md` — Phase 1 decisions (sidebar, command palette, system tray, theme)

### Research (Port Reference)
- `.planning/research/XL-CODE-REFERENCE.md` — xl's proven ADO client patterns: PAT + Azure CLI auth, WIQL queries, JSON Patch, state mapping, bulk import, per-repo config

### Existing Code
- `internal/db/db.go` — SQLite schema with tables: ado_work_items, task_ado_links, pull_requests (already created in Phase 1)
- `internal/app/adoservice.go` — Current ADO service (shells out to az cli — needs refactoring to use token + REST)
- `internal/config/config.go` — Viper config with ado.organization, ado.project, sync.interval_minutes
- `internal/config/service.go` — ConfigService exposed to frontend via Wails bindings
- `frontend/src/stores/tasks.ts` — Task store with mock data, binding fallback pattern
- `frontend/src/views/TasksView.vue` — Current task list (needs personal/public badges, filters)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/db/db.go` schema already has: `ado_work_items` (ado_id, title, state, type, assigned_to, priority, area_path, url), `task_ado_links` (task_id, ado_id, direction), `pull_requests` (deferred to Phase 3)
- `internal/config/` — Viper config with ADO org/project/sync interval, ConfigService for frontend
- `frontend/src/components/ui/` — Full shadcn-vue library (79 components): button, badge, card, select, input, textarea, separator, tooltip, scroll-area, dialog, dropdown-menu, command, tabs
- `frontend/src/lib/utils.ts` — `cn()` utility for class merging
- `pkg/ado/` — Empty placeholder, ready for ADO REST client implementation

### Established Patterns
- **Wails service binding:** Go structs → frontend via `application.NewService()`. Refactor ADOService to use REST.
- **Store pattern:** Pinia stores with `useMock` flag, try bindings → catch fallback to mock
- **shadcn-vue tokens:** bg-primary/text-primary-foreground for blue, bg-muted for subtle, text-foreground/text-muted-foreground

### Integration Points
- `main.go` — Register refactored ADOService with token provider
- `pkg/ado/` — New ADO REST client package (token provider interface, HTTP client, WIQL, JSON Patch)
- `internal/auth/` — Refactor to implement token provider interface (az cli, PAT, future OAuth)
- `frontend/src/views/TasksView.vue` — Add personal/public badges, filters, quick-add
- `frontend/src/views/AdoView.vue` — ADO browser (browse items, linked status, import/link)
- `frontend/src/components/SyncConfirmDialog.vue` — Preview diff + confirm for outbound changes
- `frontend/src/components/ConflictResolver.vue` — Per-field conflict resolution UI

</code_context>

<specifics>
## Specific Ideas

- ADO bugs can be tasks — they keep their bug icon throughout
- "Search all ADO" button in tree search for items outside your assignments
- Task row ADO badge: hollow = personal, filled with type icon = public/linked
- Quick-add input at top of task list — just type and press Enter
- Preview diff shows side-by-side: "Local value" vs "Will push to ADO"
- Conflict resolution shows: "Local value" / "ADO value" / pick buttons per field
- Token provider logs which source it's using (az cli / PAT / OAuth) for debugging

</specifics>

---

*Phase: 02-ado-integration-sync-workflow*
*Context gathered: 2026-04-03 (restructured from original 2026-04-02 discussion)*
