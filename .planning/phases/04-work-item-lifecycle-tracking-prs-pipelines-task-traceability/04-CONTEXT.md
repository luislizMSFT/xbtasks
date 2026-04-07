# Phase 4: Work Item Lifecycle Tracking — PRs, Pipelines & Task Traceability - Context

**Gathered:** 2026-04-07
**Status:** Ready for planning

<domain>
## Phase Boundary

Surface the full connection chain — tasks ↔ ADO work items ↔ PRs ↔ pipeline runs — so users never have to hunt through ADO to understand work item progress. When a PR completes, related items and pipelines appear right in the dashboard. A repurposed graph view page shows all relationships interactively, plus a timeline of all work done.

**Core problem this solves:** ADO fragments context across work items, PRs, pipelines, and boards. The user currently has to click through multiple tools to see "what happened with this task." This phase makes all connections visible in one place.

</domain>

<decisions>
## Implementation Decisions

### Completed PR Presentation (Dashboard)
- **D-01:** When a PR completes, its dashboard row becomes compact (smaller, less info) with connected badges/chips showing related tasks, work items, and pipeline runs underneath
- **D-02:** Related items are auto-linked by branch name match (PR source branch = pipeline source branch) and existing `task_id`/`ado_id` fields on the PR record — no manual linking needed
- **D-03:** Completed PRs are stored permanently in local SQLite (no re-fetching from ADO). Dashboard shows recent completions; full history lives in the Timeline tab

### PR Snooze
- **D-04:** PRs can be snoozed with quick preset durations: 1h / 4h / 1d / 1w
- **D-05:** No permanent hide — snooze only, PR reappears after duration expires
- **D-06:** Snooze state stored locally (pr_id + snooze_until timestamp)

### PR → Task Status Nudge
- **D-07:** When a linked PR merges, show a toast/nudge ("PR merged — mark task done?") but do NOT auto-change task status
- **D-08:** User explicitly confirms or dismisses — consistent with the "never auto-mutate" philosophy from Phase 2

### PR → ADO Work Item Detection
- **D-09:** Parse PR title/description for `AB#12345` patterns referencing ADO work items
- **D-10:** Surface a linking suggestion ("This PR references ADO #12345 — link it?") — user confirms before any connection is created
- **D-11:** No auto-linking — detect and suggest only

### Graph View Page — Relationships Tab
- **D-12:** Repurpose existing force graph to show ALL entity types as nodes: tasks, ADO work items, PRs, and pipeline runs with typed edges between them
- **D-13:** Replace current force-directed layout with a hierarchical/layered layout (top-to-bottom) — much cleaner than tangled force graph
- **D-14:** Polish existing graph: thicker links, straighter edges, larger nodes — current links are too small and thin
- **D-15:** Click a node → highlights its full relationship chain and opens detail panel
- **D-16:** Graph currently uses mock data — wire to real data from stores (tasks, PRs, pipelines, ADO items)

### Graph View Page — Timeline Tab
- **D-17:** New tab alongside the Relationships graph: a horizontal left-to-right timeline with time axis
- **D-18:** All entity events on one chronological stream — tasks, PRs, pipelines mixed together with type indicators
- **D-19:** Items placed by creation/completion dates — shows the full history of all work done
- **D-20:** Timeline is the long-term archive — all completed PRs, pipeline runs, and task completions visible here even after they age out of the dashboard

### Implementation Approach
- **D-21:** Playground-first — build prototype views/components before touching any production code
- **D-22:** PR/Pipeline services currently use az cli subprocess (not REST API like ADOService) — evaluate migrating to REST during research

### Agent's Discretion
- Exact hierarchical graph layout algorithm (dagre, elk, custom)
- Dashboard section ordering for completed PRs relative to active PRs
- Timeline visual design (dot markers, cards, etc.)
- Snooze UI placement (icon button, context menu, etc.)
- How many days of completed PRs to show on dashboard before they're timeline-only

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Existing PR Infrastructure
- `domain/pullrequest.go` — PullRequest domain type with TaskID, AdoID, Status fields
- `domain/pipeline.go` — Pipeline domain type with status, sourceBranch, result
- `internal/app/prservice.go` — PR fetching service (currently az cli based)
- `internal/app/pipelineservice.go` — Pipeline fetching service (currently az cli based)
- `internal/db/pr.go` — PR database operations (upsert, list, query by status)
- `internal/db/migrate.go` §pull_requests — PR table schema with task_id, ado_id foreign keys

### Existing Frontend
- `frontend/src/views/DashboardView.vue` — Dashboard with PR sections (activeMyPRs, activeReviewPRs, groupedPipelines)
- `frontend/src/views/DependencyGraphView.vue` — Force graph view (mock data, ForceGraph component)
- `frontend/src/components/ForceGraph.vue` — Force-directed graph component (needs polish)
- `frontend/src/stores/prs.ts` — PR store (myPRs, reviewPRs, fetchAll)
- `frontend/src/stores/ado.ts` — ADO store (pipelines, fetchPipelines)
- `frontend/src/api/prs.ts` — PR API wrapper
- `frontend/src/api/pipelines.ts` — Pipeline API wrapper
- `frontend/src/components/tasks/TaskDetail.vue` — Task detail showing taskPRs computed

### Phase 2 Decisions
- `.planning/phases/02-ado-integration-prs-unified-dashboard/02-CONTEXT.md` — Personal→public model, sync philosophy, no auto-push/mutate pattern

### Requirements
- `.planning/REQUIREMENTS.md` §v2 — PR-01 through PR-04, PIPE-01, PIPE-02 (v2 requirements now being pulled forward)
- `.planning/REQUIREMENTS.md` §v1 — UX-06 (PR scope: authored + reviewing, no abandoned), UX-07 (pipeline names not IDs)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `DashboardView.vue` already has `activeMyPRs`, `activeReviewPRs`, `groupedPipelines` computeds — extend for completed PRs
- `myPRBranches` computed already collects branch names from PRs — reuse for pipeline matching
- `ForceGraph.vue` component exists — refactor from force-directed to hierarchical layout
- `TaskDetail.vue` already computes `taskPRs` — pattern for linking PRs to tasks
- `domain.PullRequest.TaskID` and `domain.PullRequest.AdoID` fields exist for linking

### Established Patterns
- PR/pipeline services use az cli subprocess — different from ADOService which uses REST. Migration candidate.
- Dashboard fetches on mount if stores empty — same pattern for new data
- `openPR()`/`openPipeline()` use `window.open` or Wails `openURL` — consistent external link pattern
- Toast notifications via vue-sonner — use for PR merge nudges

### Integration Points
- Dashboard PR sections — add completed PR section with related items
- DependencyGraphView — replace mock data, add tabs, refactor layout
- PR store — add completed PR tracking, snooze state
- SQLite migrations — add snooze table/columns

</code_context>

<specifics>
## Specific Ideas

- "I want to be aware of all connections" — the core UX principle is full visibility of task↔PR↔pipeline↔work item relationships
- "Never have to go hunting" — the pain point is ADO's fragmented UX where you click through multiple tools to understand progress
- "The original intention of the graph view is to interactively see relation" — graph view was always meant for this, just needs real data and polish
- "Instead of messy graph show a hierarchy with top and bottom" — clean layered layout, not tangled force-directed
- "Links are too small and thin and should be straighter" — specific graph polish feedback
- "Keep forever in local db just don't query over and over" — cache completed PRs permanently, don't re-fetch

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 04-work-item-lifecycle-tracking-prs-pipelines-task-traceability*
*Context gathered: 2026-04-07*
