# Phase 4: Work Item Lifecycle Tracking — PRs, Pipelines & Task Traceability - Research

**Researched:** 2026-04-08
**Domain:** Azure DevOps PR/Pipeline REST API, hierarchical graph visualization, timeline UI, SQLite state management
**Confidence:** HIGH

## Summary

This phase transforms the existing dashboard and graph view into a full lifecycle traceability surface. The core work is: (1) migrating PR and pipeline services from `az cli` subprocess calls to the existing `pkg/ado` REST client pattern, (2) adding completed PR presentation with linked badges on the dashboard, (3) PR snooze functionality, (4) PR→task merge nudges, (5) PR→ADO work item detection via `AB#` pattern matching, (6) replacing the force-directed graph with a hierarchical/layered layout, (7) wiring the graph to real data with all entity types (tasks, PRs, pipelines, ADO items), and (8) building a new timeline tab as a chronological archive.

The codebase is well-structured for this work. The `pkg/ado.Client` already implements REST API patterns (Bearer auth, JSON decoding, batch fetching) that PRService and PipelineService should adopt. The frontend has established patterns for stores, API wrappers (dynamic Wails binding imports), and dashboard sections. The `ForceGraph.vue` component (664 lines) uses d3-force with SVG rendering — it needs to be replaced/refactored to use dagre for hierarchical layout while keeping the SVG rendering approach. Vue-sonner is already installed for toast notifications (PR merge nudge). SQLite schema already has `pull_requests` table with `task_id` and `ado_id` foreign keys.

**Primary recommendation:** Migrate PR/Pipeline services to REST API first (unblocks multi-org + removes az cli dependency for these calls), then layer dashboard completed PR section + snooze, then graph view overhaul, then timeline tab. Playground-first per D-21.

<user_constraints>

## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** When a PR completes, its dashboard row becomes compact (smaller, less info) with connected badges/chips showing related tasks, work items, and pipeline runs underneath
- **D-02:** Related items are auto-linked by branch name match (PR source branch = pipeline source branch) and existing `task_id`/`ado_id` fields on the PR record — no manual linking needed
- **D-03:** Completed PRs are stored permanently in local SQLite (no re-fetching from ADO). Dashboard shows recent completions; full history lives in the Timeline tab
- **D-04:** PRs can be snoozed with quick preset durations: 1h / 4h / 1d / 1w
- **D-05:** No permanent hide — snooze only, PR reappears after duration expires
- **D-06:** Snooze state stored locally (pr_id + snooze_until timestamp)
- **D-07:** When a linked PR merges, show a toast/nudge ("PR merged — mark task done?") but do NOT auto-change task status
- **D-08:** User explicitly confirms or dismisses — consistent with the "never auto-mutate" philosophy from Phase 2
- **D-09:** Parse PR title/description for `AB#12345` patterns referencing ADO work items
- **D-10:** Surface a linking suggestion ("This PR references ADO #12345 — link it?") — user confirms before any connection is created
- **D-11:** No auto-linking — detect and suggest only
- **D-12:** Repurpose existing force graph to show ALL entity types as nodes: tasks, ADO work items, PRs, and pipeline runs with typed edges between them
- **D-13:** Replace current force-directed layout with a hierarchical/layered layout (top-to-bottom) — much cleaner than tangled force graph
- **D-14:** Polish existing graph: thicker links, straighter edges, larger nodes — current links are too small and thin
- **D-15:** Click a node → highlights its full relationship chain and opens detail panel
- **D-16:** Graph currently uses mock data — wire to real data from stores (tasks, PRs, pipelines, ADO items)
- **D-17:** New tab alongside the Relationships graph: a horizontal left-to-right timeline with time axis
- **D-18:** All entity events on one chronological stream — tasks, PRs, pipelines mixed together with type indicators
- **D-19:** Items placed by creation/completion dates — shows the full history of all work done
- **D-20:** Timeline is the long-term archive — all completed PRs, pipeline runs, and task completions visible here even after they age out of the dashboard
- **D-21:** Playground-first — build prototype views/components before touching any production code
- **D-22:** PR/Pipeline services currently use az cli subprocess (not REST API like ADOService) — evaluate migrating to REST during research

### Agent's Discretion
- Exact hierarchical graph layout algorithm (dagre, elk, custom)
- Dashboard section ordering for completed PRs relative to active PRs
- Timeline visual design (dot markers, cards, etc.)
- Snooze UI placement (icon button, context menu, etc.)
- How many days of completed PRs to show on dashboard before they're timeline-only

### Deferred Ideas
None — discussion stayed within phase scope

</user_constraints>

<phase_requirements>

## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| PR-01 | User can see all their own PRs with current status (draft, active, completed, abandoned) | Already partially implemented in `PRService.ListMyPRs()` — needs REST migration + completed PR presentation on dashboard |
| PR-02 | User can see all PRs assigned to them for review with status | Already implemented in `PRService.ListReviewPRs()` — needs REST migration |
| PR-03 | User can view PR details (title, repo, source/target branch, reviewers, votes) | PR domain type already has all fields; TaskDetail already shows `taskPRs` computed |
| PR-04 | PRs displayed under linked tasks in detail panel | Already working via `taskPRs` computed in `TaskDetail.vue` — enhance with completed PR badges |
| PIPE-01 | User can see recent pipeline runs with status (succeeded, failed, running, queued, cancelled) | Already implemented in `PipelineService.ListRecentRuns()` — needs REST migration |
| PIPE-02 | Pipeline status shown on linked PR or task | Dashboard already matches pipelines to PR branches via `myPRBranches` computed; extend to completed PR badges |
| UX-06 | PR view scoped to: user's authored PRs + PRs user is reviewing. Excludes abandoned. No "all team PRs" flooding. | Filter logic needed — currently `ListMyPRs` fetches `--status all` including abandoned |
| UX-07 | Pipeline runs show proper pipeline names (not just IDs), readable status | Already implemented — `mapPipelineRun` resolves `pipeline.Name` with fallback to `name` then `Pipeline #ID` |

</phase_requirements>

## Standard Stack

### Core (Already Installed)
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| d3-force | 3.0.0 | Physics simulation (existing, to be reduced in scope) | Already in package.json, still used for any physics needs |
| d3-selection | 3.0.0 | DOM manipulation for SVG | Already installed, used by ForceGraph |
| d3-zoom | 3.0.0 | Pan/zoom behavior | Already installed, used by ForceGraph |
| d3-drag | 3.0.0 | Node dragging | Already installed, used by ForceGraph |
| vue-sonner | 2.0.9 | Toast notifications | Already installed, used for PR merge nudge (D-07) |
| pinia | 3.0.4 | State management | Already installed, all stores use it |

### New Dependencies
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @dagrejs/dagre | 3.0.0 | Hierarchical/layered graph layout | Replace d3-force layout for graph view (D-13). Computes x,y positions for nodes in a DAG layout |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| @dagrejs/dagre | elkjs 0.11.1 | ELK is more powerful (handles compound graphs, port constraints) but 10x larger (~2MB vs ~200KB), async-only (WASM), overkill for this use case |
| @dagrejs/dagre | Custom d3-hierarchy | d3-hierarchy only does trees, not DAGs with multiple parents — dagre handles cross-edges correctly |
| @dagrejs/dagre | Original `dagre` 0.8.5 | `@dagrejs/dagre` is the maintained fork (v3.0.0), original is abandoned since 2019 |

**Recommendation: Use `@dagrejs/dagre` 3.0.0.** It's the right fit: lightweight (~200KB), synchronous layout computation, handles DAGs with cross-edges, produces clean top-to-bottom hierarchical output. No WASM overhead. The existing SVG rendering approach with d3-zoom/d3-drag can be kept — only the layout engine changes.

**Installation:**
```bash
cd frontend && npm install @dagrejs/dagre @types/dagre
```

**Version verification:** `@dagrejs/dagre` 3.0.0 is the latest on npm. `@types/dagre` provides TypeScript declarations compatible with both the original dagre and the fork (same API surface).

## Architecture Patterns

### Backend: PR/Pipeline REST API Migration

**Current state:** `PRService` and `PipelineService` shell out to `az cli` via `ado.RunAzCli()`. This is slow (process spawn per call), doesn't support multi-org, and is the only remaining az cli dependency for data fetching (ADOService already migrated to REST in Phase 2).

**Target state:** Both services should use `pkg/ado.Client` REST API directly, following the same pattern as `pkg/ado/query.go`.

#### ADO REST API Endpoints for PRs

```
GET https://dev.azure.com/{org}/{project}/_apis/git/pullrequests?searchCriteria.creatorId={userId}&searchCriteria.status=all&$top=50&api-version=7.0
GET https://dev.azure.com/{org}/{project}/_apis/git/pullrequests?searchCriteria.reviewerId={userId}&searchCriteria.status=active&$top=50&api-version=7.0
```

- API version: `7.0` (consistent with existing `query.go` usage)
- Response shape: `{ "value": [{ "pullRequestId": int, "title": string, ... }] }`
- `searchCriteria.creatorId` / `searchCriteria.reviewerId` accept user GUIDs
- Current user GUID available via `https://dev.azure.com/{org}/_apis/connectionData` or the profile API

#### ADO REST API Endpoints for Pipeline Runs

```
GET https://dev.azure.com/{org}/{project}/_apis/pipelines/runs?$top=20&api-version=7.0
```

Alternative (Build API, more fields):
```
GET https://dev.azure.com/{org}/{project}/_apis/build/builds?$top=20&api-version=7.0
```

- Build API returns `definition.name` (pipeline name), `sourceBranch`, `status`, `result`, `queueTime`, `finishTime`
- The Build API is more reliable for pipeline names than the Pipelines API

#### Migration Pattern

```go
// pkg/ado/pullrequests.go — new file
func ListPullRequests(c *Client, creatorID string, status string, top int) ([]PRResponse, error) {
    url := c.apiURL(fmt.Sprintf(
        "git/pullrequests?searchCriteria.creatorId=%s&searchCriteria.status=%s&$top=%d&api-version=7.0",
        creatorID, status, top,
    ))
    resp, err := c.doRequest("GET", url, nil, "")
    if err != nil {
        return nil, fmt.Errorf("list PRs: %w", err)
    }
    var result struct {
        Value []PRResponse `json:"value"`
    }
    if err := decodeResponse(resp, &result); err != nil {
        return nil, fmt.Errorf("parse PR response: %w", err)
    }
    return result.Value, nil
}
```

PRService would then:
1. Accept `tokenProvider` and `configService` (like ADOService)
2. Create `ado.Client` instances per org/project pair
3. Get current user GUID from connection data API (cache it)
4. Call REST endpoints instead of `ado.RunAzCli`
5. Multi-org support: iterate all clients, merge results

### Frontend: Completed PR Presentation on Dashboard

**Pattern:** Add a `completedPRs` computed to `DashboardView.vue` that filters from the PR store:

```typescript
const completedPRs = computed(() =>
  prStore.myPRs
    .filter(pr => pr.status === 'completed' && !isSnoozed(pr.id))
    .slice(0, maxCompletedDays) // agent's discretion: recommend 7 days
)
```

**Badge linking pattern (D-02):** For each completed PR, find related items by:
1. `pr.taskId` → directly linked task
2. `pr.adoId` → directly linked ADO work item  
3. Branch name match: pipelines where `pipeline.sourceBranch === pr.sourceBranch`

This pattern already exists in `DashboardView.vue` via `myPRBranches` computed.

### Frontend: Snooze State

**Pattern:** Add a snooze table to SQLite + a backend service method:

```sql
CREATE TABLE IF NOT EXISTS pr_snooze (
    pr_id   INTEGER NOT NULL,
    repo    TEXT NOT NULL,
    snooze_until DATETIME NOT NULL,
    PRIMARY KEY (pr_id, repo)
);
```

Backend exposes: `SnoozePR(prID int, repo string, duration string)` and `ListSnoozedPRs()`.

Frontend PR store adds a `snoozedPRs` map and filters snoozed PRs from active views.

### Frontend: Graph View Overhaul

**Current architecture:** `DependencyGraphView.vue` (819 lines) + `ForceGraph.vue` (664 lines).
- `DependencyGraphView` has 3 view modes: graph, tree, flat
- `ForceGraph.vue` renders SVG with d3-force simulation, d3-zoom, d3-drag
- Graph already loads real task data + ADO parent nodes (mixed local + ADO)

**Target architecture:**
1. Rename route/view: `DependencyGraphView` → `RelationshipsView` (or keep path, update content)
2. Add tabs: **Relationships** (graph) | **Timeline** (new)
3. Replace `ForceGraph.vue` with `HierarchicalGraph.vue` that uses dagre for layout
4. Add PR and pipeline nodes to the graph data
5. Keep d3-zoom, d3-drag, SVG rendering — only replace the layout engine

**Dagre integration pattern:**

```typescript
import dagre from '@dagrejs/dagre'

function layoutGraph(nodes: GraphNode[], edges: GraphEdge[]) {
  const g = new dagre.graphlib.Graph()
  g.setGraph({ rankdir: 'TB', nodesep: 60, ranksep: 80, edgesep: 20 })
  g.setDefaultEdgeLabel(() => ({}))
  
  for (const node of nodes) {
    g.setNode(String(node.id), { width: node.width || 180, height: node.height || 50 })
  }
  for (const edge of edges) {
    g.setEdge(String(edge.source), String(edge.target))
  }
  
  dagre.layout(g)
  
  // Read computed positions
  return nodes.map(node => {
    const pos = g.node(String(node.id))
    return { ...node, x: pos.x, y: pos.y }
  })
}
```

### Frontend: Timeline Tab

**Pattern:** Horizontal left-to-right timeline rendering with CSS/SVG.

Collect all events from stores:
- Tasks: `createdAt`, `completedAt`
- PRs: `createdAt`, `mergedAt`
- Pipelines: `queueTime`, `finishTime`

Merge into a single sorted array by date, render as a horizontal scrollable timeline with type-colored indicators.

### Anti-Patterns to Avoid
- **Don't build a second graph library:** Use dagre for layout computation only. Render with the existing SVG + d3-zoom approach. Don't pull in a full graph rendering framework like cytoscape.js or vis-network.
- **Don't auto-mutate task status on PR merge:** Per D-07/D-08, show toast only. This is a locked design principle from Phase 2.
- **Don't re-fetch completed PRs from ADO:** Per D-03, completed PRs live permanently in local SQLite. Only active PRs should be re-fetched on sync.
- **Don't bundle the full ELK runtime:** elkjs includes a WASM runtime that's 2MB+. Dagre at ~200KB is sufficient.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Hierarchical graph layout | Custom topological sort + position algorithm | `@dagrejs/dagre` | Layout algorithms are deceptively complex — node overlap avoidance, edge crossing minimization, rank assignment. Dagre handles all of this |
| Toast notifications | Custom notification system | `vue-sonner` (already installed) | Already used via `useNotify()` composable. Supports action buttons for "Mark task done?" flow |
| `AB#12345` pattern detection | Custom regex | `RegExp(/AB#(\d+)/gi)` | This is simple enough for regex, but don't over-engineer — a simple regex is correct here |
| Timeline date formatting | Custom date math | `relativeTime()` from `@/lib/date` | Already exists, already used throughout dashboard |
| SVG pan/zoom | Custom mouse handlers | `d3-zoom` (already installed) | Already working in ForceGraph.vue, reuse in HierarchicalGraph |

## Common Pitfalls

### Pitfall 1: ADO REST API User ID Resolution
**What goes wrong:** The PR REST API uses user GUIDs (`searchCriteria.creatorId`), not display names or UPNs. The current az cli approach uses display names.
**Why it happens:** ADO REST API and az cli use different identity formats.
**How to avoid:** Use the ADO Connection Data API (`GET https://dev.azure.com/{org}/_apis/connectionData`) or the Profile API to get the current user's GUID. Cache it (like `PRService.meEmail` currently does).
**Warning signs:** Empty PR results even though user has PRs — usually means wrong user identifier format.

### Pitfall 2: Multi-Org PR Deduplication
**What goes wrong:** When fetching PRs across multiple org/project pairs, the same PR can appear in different projects if repos are shared.
**Why it happens:** ADO repos can be shared across projects. The PR search scope is per-project.
**How to avoid:** Deduplicate by `(prNumber, repo)` key — this is already done in `SyncPRs()`. Maintain this pattern in the REST migration.
**Warning signs:** Duplicate PR cards on dashboard.

### Pitfall 3: Dagre Layout with Disconnected Components
**What goes wrong:** Dagre only lays out connected components. Disconnected nodes pile up at (0,0).
**Why it happens:** Not all tasks have dependencies — many are standalone.
**How to avoid:** Either (a) compute separate layouts for disconnected components and arrange them in a grid, or (b) add invisible root nodes to connect components. Option (a) is cleaner. D3-force handled this automatically because it uses physics simulation.
**Warning signs:** Nodes overlapping in the top-left corner.

### Pitfall 4: Snooze Timer Precision
**What goes wrong:** Snoozed PRs don't reappear at the right time because the frontend only checks on navigation/mount.
**Why it happens:** No periodic re-evaluation of snooze expiry.
**How to avoid:** Use a 1-minute `setInterval` in the PR store to check snooze expiry, or check on every dashboard mount/activate. Since the dashboard already uses `onActivated`, checking there is sufficient for a desktop app.
**Warning signs:** User snoozes for 1h but PR doesn't reappear until next app restart.

### Pitfall 5: AB# Pattern False Positives
**What goes wrong:** Regex matches `AB#` in contexts where it's not an ADO work item reference (e.g., in code snippets, markdown headers).
**Why it happens:** `AB#12345` is a convention, not a guaranteed link.
**How to avoid:** Only scan PR title and first paragraph of description. Validate that the extracted ID actually exists in the local ADO work items cache before suggesting a link. Present as suggestion (D-10), never auto-link (D-11).
**Warning signs:** Suggestions to link non-existent work items.

### Pitfall 6: Graph Performance with Many Nodes
**What goes wrong:** Rendering hundreds of nodes with SVG causes janky pan/zoom.
**Why it happens:** SVG DOM updates are expensive at scale. Each node is a group of DOM elements (circle + text + badges).
**How to avoid:** For the current dataset (personal tasks + linked ADO items), SVG is fine — likely <100 nodes. If scaling becomes an issue, switch to Canvas rendering. But don't prematurely optimize.
**Warning signs:** Laggy pan/zoom with >200 nodes. Current ForceGraph already handles ~50 nodes smoothly.

## Code Examples

### 1. ADO REST API: List Pull Requests (Go)

```go
// Source: Azure DevOps REST API v7.0 - Git Pull Requests
// GET {org}/{project}/_apis/git/pullrequests
func ListPullRequests(c *Client, creatorID string, status string, top int) ([]PRResult, error) {
    url := c.apiURL(fmt.Sprintf(
        "git/pullrequests?searchCriteria.creatorId=%s&searchCriteria.status=%s&$top=%d&api-version=7.0",
        creatorID, status, top,
    ))
    resp, err := c.doRequest("GET", url, nil, "")
    if err != nil {
        return nil, err
    }
    var result struct {
        Value []PRResult `json:"value"`
    }
    return result.Value, decodeResponse(resp, &result)
}

// PRResult matches the ADO REST API response for a pull request.
type PRResult struct {
    PullRequestID int    `json:"pullRequestId"`
    Title         string `json:"title"`
    Status        string `json:"status"`
    IsDraft       bool   `json:"isDraft"`
    Repository    struct {
        Name string `json:"name"`
    } `json:"repository"`
    SourceRefName string `json:"sourceRefName"`
    TargetRefName string `json:"targetRefName"`
    CreatedBy     Identity `json:"createdBy"`
    Reviewers     []struct {
        DisplayName string `json:"displayName"`
        UniqueName  string `json:"uniqueName"`
        Vote        int    `json:"vote"`
    } `json:"reviewers"`
    CreationDate string  `json:"creationDate"`
    ClosedDate   *string `json:"closedDate"`
    URL          string  `json:"url"`
}
```

### 2. ADO REST API: Current User GUID (Go)

```go
// Source: Azure DevOps REST API - Connection Data
// GET {org}/_apis/connectionData
func GetCurrentUserID(c *Client) (string, error) {
    url := fmt.Sprintf("https://dev.azure.com/%s/_apis/connectionData", c.Org())
    resp, err := c.doRequest("GET", url, nil, "")
    if err != nil {
        return "", err
    }
    var result struct {
        AuthenticatedUser struct {
            ID string `json:"id"`
        } `json:"authenticatedUser"`
    }
    if err := decodeResponse(resp, &result); err != nil {
        return "", err
    }
    return result.AuthenticatedUser.ID, nil
}
```

### 3. ADO REST API: List Pipeline Runs / Builds (Go)

```go
// Source: Azure DevOps REST API v7.0 - Build
// GET {org}/{project}/_apis/build/builds
func ListRecentBuilds(c *Client, top int) ([]BuildResult, error) {
    url := c.apiURL(fmt.Sprintf("build/builds?$top=%d&api-version=7.0", top))
    resp, err := c.doRequest("GET", url, nil, "")
    if err != nil {
        return nil, err
    }
    var result struct {
        Value []BuildResult `json:"value"`
    }
    return result.Value, decodeResponse(resp, &result)
}

type BuildResult struct {
    ID         int    `json:"id"`
    BuildNumber string `json:"buildNumber"`
    Status     string `json:"status"`     // notStarted, inProgress, completed
    Result     string `json:"result"`     // succeeded, failed, canceled, partiallySucceeded
    Definition struct {
        Name string `json:"name"`
    } `json:"definition"`
    SourceBranch string  `json:"sourceBranch"`
    QueueTime    string  `json:"queueTime"`
    FinishTime   *string `json:"finishTime"`
    URL          string  `json:"url"`
}
```

### 4. Dagre Hierarchical Layout (TypeScript)

```typescript
// Source: @dagrejs/dagre documentation
import dagre from '@dagrejs/dagre'

interface LayoutNode {
  id: string
  x: number
  y: number
  width: number
  height: number
}

function computeHierarchicalLayout(
  nodes: { id: string; width: number; height: number }[],
  edges: { source: string; target: string }[],
): LayoutNode[] {
  const g = new dagre.graphlib.Graph()
  g.setGraph({
    rankdir: 'TB',     // Top-to-bottom (D-13)
    nodesep: 60,        // Horizontal spacing between nodes
    ranksep: 80,        // Vertical spacing between ranks
    edgesep: 20,        // Spacing between edges
    marginx: 20,
    marginy: 20,
  })
  g.setDefaultEdgeLabel(() => ({}))

  for (const node of nodes) {
    g.setNode(node.id, { width: node.width, height: node.height })
  }
  for (const edge of edges) {
    g.setEdge(edge.source, edge.target)
  }

  dagre.layout(g)

  return nodes.map(node => {
    const pos = g.node(node.id)
    return { ...node, x: pos.x, y: pos.y, width: pos.width, height: pos.height }
  })
}
```

### 5. PR Merge Nudge Toast (TypeScript/Vue)

```typescript
// Source: Existing useNotify composable + vue-sonner action support
import { toast } from 'vue-sonner'

function showMergeNudge(pr: PullRequest, taskTitle: string) {
  toast.info(`PR merged — mark "${taskTitle}" done?`, {
    duration: 10000,
    action: {
      label: 'Mark Done',
      onClick: () => {
        taskStore.setStatus(pr.taskId!, 'done')
      },
    },
    cancel: {
      label: 'Dismiss',
      onClick: () => { /* no-op */ },
    },
  })
}
```

### 6. AB# Pattern Detection (TypeScript)

```typescript
// Source: ADO auto-linking convention
const AB_PATTERN = /AB#(\d+)/gi

function detectWorkItemRefs(title: string, description: string): number[] {
  const text = `${title}\n${description.split('\n').slice(0, 3).join('\n')}` // first 3 lines
  const matches = [...text.matchAll(AB_PATTERN)]
  return [...new Set(matches.map(m => parseInt(m[1], 10)))]
}
```

### 7. Snooze Duration Calculation (TypeScript)

```typescript
// Source: D-04 preset durations
const SNOOZE_PRESETS = {
  '1h': 60 * 60 * 1000,
  '4h': 4 * 60 * 60 * 1000,
  '1d': 24 * 60 * 60 * 1000,
  '1w': 7 * 24 * 60 * 60 * 1000,
} as const

function snoozeUntil(preset: keyof typeof SNOOZE_PRESETS): Date {
  return new Date(Date.now() + SNOOZE_PRESETS[preset])
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `az repos pr list` subprocess | ADO REST API `_apis/git/pullrequests` | Phase 2 established REST pattern for work items | PRService should follow same migration path |
| `az pipelines runs list` subprocess | ADO REST API `_apis/build/builds` | Phase 2 established REST pattern | PipelineService should follow same migration path |
| Force-directed graph (d3-force) | Hierarchical DAG layout (dagre) | User decision D-13 | Cleaner visualization for parent-child + dependency relationships |
| `dagre` 0.8.5 (abandoned) | `@dagrejs/dagre` 3.0.0 (maintained) | 2023+ | Use the maintained fork, not the abandoned original |

**Deprecated/outdated:**
- `dagre` npm package: Last publish 2019, unmaintained. Use `@dagrejs/dagre` instead.
- `az repos pr list` for PR fetching: Works but is slow (subprocess spawn), single-org only, and inconsistent with the REST pattern used by ADOService.

## Open Questions

1. **User GUID caching strategy**
   - What we know: The ADO REST API needs user GUIDs for PR queries. The Connection Data API returns the authenticated user's GUID.
   - What's unclear: Should the GUID be cached in memory (like `meEmail` currently) or persisted to SQLite?
   - Recommendation: Cache in memory on the service struct, same as current `meEmail` pattern. It's per-session and doesn't change.

2. **Completed PR retention on dashboard**
   - What we know: D-03 says completed PRs stored permanently, dashboard shows "recent completions"
   - What's unclear: How many days of completed PRs to show on dashboard
   - Recommendation: 7 days is a reasonable default. This is in the agent's discretion area. Could be a config setting.

3. **Graph node sizing for different entity types**
   - What we know: Tasks, PRs, pipelines, and ADO work items all become graph nodes. They have different information density.
   - What's unclear: Should all node types be the same size or vary by type?
   - Recommendation: Use consistent width (180px) but allow different heights. Tasks and PRs get taller nodes with more detail. Pipeline nodes can be smaller/simpler.

4. **PR description access for AB# parsing**
   - What we know: The current PR domain type doesn't include a `description` field. AB# references are often in the PR description.
   - What's unclear: Do we need to add a description field or just parse the title?
   - Recommendation: Add a `description` field to the PR response parsing. The ADO REST API returns it. For storage, either add it to the SQLite schema or just scan at fetch time and store detected work item IDs.

## Project Constraints (from copilot-instructions.md)

- **Stack:** Wails v3 (Go) + Vue 3 + SQLite — native desktop app
- **Vue is thin shell:** All logic lives in Go. Vue is display + interaction only.
- **Local-first:** SQLite per user, no server dependency
- **No emojis:** Lucide icons + AzureDevOpsIcon only
- **Auth:** Abstracted token provider (az cli first, swappable)
- **ADO client:** Direct REST API from Go (no shelling out per query) — this reinforces the migration decision
- **Outbound sync:** Always requires user confirmation — never auto-push (relevant for D-07/D-08)
- **GSD Workflow:** Use `/gsd-execute-phase` for planned work

## Sources

### Primary (HIGH confidence)
- Codebase analysis: `internal/app/prservice.go`, `internal/app/pipelineservice.go`, `internal/db/pr.go`, `internal/db/migrate.go` — current implementation fully read
- Codebase analysis: `pkg/ado/client.go`, `pkg/ado/query.go` — REST API client pattern established
- Codebase analysis: `frontend/src/views/DashboardView.vue`, `frontend/src/views/DependencyGraphView.vue`, `frontend/src/components/ForceGraph.vue` — full current state understood
- Codebase analysis: `frontend/src/stores/prs.ts`, `frontend/src/stores/ado.ts` — store patterns documented
- npm registry: `@dagrejs/dagre` 3.0.0 verified, `elkjs` 0.11.1 verified

### Secondary (MEDIUM confidence)
- Azure DevOps REST API documentation: PR search criteria, Build API, Connection Data API — API endpoint formats based on established ADO v7.0 API patterns used in existing codebase
- dagre API surface: `graphlib.Graph`, `setGraph({ rankdir })`, `layout()` — well-established API, stable since dagre 0.8.x

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all core libraries already in use, only dagre is new and it's a well-known library
- Architecture: HIGH — clear migration path from az cli to REST, established patterns in codebase
- Pitfalls: HIGH — based on direct codebase analysis and understanding of ADO API behavior
- Graph layout: MEDIUM — dagre recommendation is strong but disconnected component handling needs implementation attention

**Research date:** 2026-04-08
**Valid until:** 2026-05-08 (stable stack, no fast-moving dependencies)
