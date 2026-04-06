# Phase 2: ADO Integration & Sync Workflow - Research

**Researched:** 2026-04-04
**Domain:** Azure DevOps REST API integration, bidirectional sync, token-based auth, Go backend + Vue 3 frontend
**Confidence:** HIGH

## Summary

Phase 2 transforms xb-tasks from a local task manager into a bidirectional ADO integration hub. The core domains are: (1) abstracted token provider for ADO auth (az cli first), (2) a proper ADO REST client in Go replacing the current az-cli-subprocess approach, (3) the personal→public task model with link/promote/import flows, (4) safe bidirectional sync with preview-diff confirmation for outbound changes, and (5) significant frontend work — unified task list with filters, ADO browser tree view, conflict resolution UI, external links, and comments.

The project has a strong foundation from Phase 1: SQLite schema already has `ado_work_items`, `task_ado_links`, and `tasks` tables. The XL codebase provides proven ADO client patterns (WIQL queries, JSON Patch, state mapping, batch fetching). The current `ADOService` shells out to `az cli` per query and needs full refactoring to use token + REST directly.

**Primary recommendation:** Build the ADO REST client (`pkg/ado/`) first as the foundational layer, then token provider (`internal/auth/token.go`), then refactor `ADOService` to use them. Frontend work flows naturally from there: personal/public badges → ADO browser → sync confirmation → conflict resolution → filters/projects/links/comments.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Abstracted token provider interface in Go — implementations: az cli (`az account get-access-token`), PAT (from keyring), future OAuth
- **D-02:** Token auto-refreshes transparently; cached with TTL, re-fetched when expired
- **D-03:** Direct ADO REST API calls from Go using token from provider — no shelling out to az cli per query
- **D-04:** ADO client package (`pkg/ado/`) — stateless HTTP client accepting token, handles pagination, rate limits, WIQL queries, JSON Patch
- **D-04a:** Single token assumed to work across all configured orgs (same Azure tenant). Multi-tenant support deferred.
- **D-04b:** Config supports a list of org/project pairs (replaces single `ado.organization` / `ado.project`)
- **D-04c:** User configures orgs first, then picks specific projects within each org to follow
- **D-04d:** All configured org/project pairs are synced — work items and ADO browser aggregate across all of them
- **D-04e:** Each item carries its org/project origin as metadata (shown as badge/label in UI)
- **D-04f:** Unified list is default; user can toggle group-by-project view (collapsible sections per org/project)
- **D-05:** Tasks have a computed "public" status derived from whether they have an entry in `task_ado_links` table
- **D-06:** Unified list view shows personal + public tasks together; public tasks have a filled ADO badge, personal have hollow/empty badge
- **D-07:** Quick-add: create task with just a title (all other fields optional, filled in later)
- **D-08:** Subtasks of a public (linked) task remain personal unless the user explicitly links them individually
- **D-09:** When promoting a personal task to ADO, only title/status/description are pushed — subtasks, personal priority, local notes stay local
- **D-10 through D-16:** ADO browser view with tree hierarchy, linked status indicators, hide-linked toggle, import/link actions, search
- **D-17 through D-20a:** Link, Promote, Import, Unlink flows + badge UX + project linking
- **D-21 through D-25:** Sync behavior — auto-pull on timer, outbound requires confirmation, linked fields: title/status/description
- **D-25a through D-25d:** External links with auto-detect patterns, stored locally
- **D-25e through D-25h:** Comments — local by default, selective push to ADO, visual distinction
- **D-25i through D-25j:** Investigation workflow — task as hub
- **D-26 through D-28:** Per-field conflict resolution
- **D-29 through D-32:** Dashboard filters and sorting

### Claude's Discretion
- ADO REST API client library design (package structure, error handling)
- Token caching strategy and TTL
- Preview diff UI component design
- Conflict resolution UI design
- ADO state mapping (todo→Proposed, in_progress→Active, done→Completed, blocked→Blocked)
- Bulk import UX from ADO browser view
- Search implementation (local cache vs live API)

### Deferred Ideas (OUT OF SCOPE)
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
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| AUTH-01 | App authenticates to ADO via abstracted token provider (az cli initially) | Token provider interface pattern from XL; az cli `get-access-token --resource 499b84ac-...` |
| AUTH-02 | Token auto-refreshes transparently | Cached token with TTL, re-fetch on expiry via token provider |
| AUTH-03 | Token provider abstraction allows future swap to OAuth/PAT | Go interface `TokenProvider` with multiple implementations |
| TASK-08 | Personal→public model: tasks start local, linking makes them public | Computed from `task_ado_links` table presence; schema already exists |
| TASK-09 | Quick-add with just title | Frontend quick-add already partially built in TasksView; needs backend simplification |
| ADO-01 | View all ADO work items assigned to user | WIQL `@Me` query + batch fetch pattern from XL |
| ADO-02 | View ADO work item details | ADO REST GET `/wit/workitems/{id}` with field expansion |
| ADO-03 | Link personal task to ADO item | `task_ado_links` INSERT with direction='linked' + ado_id on task |
| ADO-04 | Promote personal task to new ADO work item | JSON Patch POST to create work item + link back |
| ADO-05 | Import ADO item as local task | Fetch ADO item, create local task, insert link with direction='imported' |
| ADO-06 | ADO browser view with linked status | Dedicated Vue view + Go service aggregating across org/project pairs |
| ADO-07 | Unlink task from ADO | Delete from `task_ado_links`, clear `ado_id`, user chooses keep/delete |
| ADO-08 | Direct ADO REST API from Go | `pkg/ado/` client with Bearer auth, proven XL patterns |
| ADO-09 | Configure multiple ADO orgs/projects | Viper config list structure, settings UI |
| ADO-10 | Unified list with org/project labels, group-by-project toggle | Frontend computed grouping + UI toggle |
| SYNC-01 | Background auto-sync pulls changes on timer + manual refresh | Go goroutine with `time.Ticker`, configurable interval via Viper |
| SYNC-02 | Outbound pushes require preview diff + confirmation | SyncService generates diff, frontend dialog shows changes |
| SYNC-03 | Per-field conflict resolution | Compare local vs ADO timestamps/values per field |
| SYNC-04 | Linked task syncs title/status/description; subtasks stay personal | Selective field sync in SyncService |
| LINK-01 | Attach external URLs with label | New `task_links` table in SQLite |
| LINK-02 | Auto-detect URL patterns (ICM, Grafana, ADO, Wiki) | Regex-based pattern matching in Go |
| LINK-03 | Links displayed with type icon and clickable URL | Frontend links section in TaskDetail |
| CMT-01 | Add local comments (private) | New `task_comments` table in SQLite |
| CMT-02 | Selectively push comment to ADO (with confirmation) | ADO REST API POST comment + mark as synced |
| CMT-03 | Update description locally, confirm-push to ADO | Preview diff dialog for description changes |
| DASH-01 | Unified list with personal/public visual indicator | Filled vs hollow ADO badge on task rows |
| DASH-02 | Filterable by status, priority, project, due date, ADO link status | Enhanced ListFiltered in Go + frontend filter chips |
| DASH-03 | Linked items show connection | ADO badge + status display |
| TL-01 | Single global task list, filter by project | Existing tasks query + project filter |
| TL-02 | Medium-density rows with all indicators | TaskRow component enhancement |
| TL-03 | Flat sorted, optional group-by toggle | Frontend computed grouping |
| TL-04 | Click task → slide-out right panel | Already exists (TaskDetail sidebar) |
| TL-05 | Tasks can exist without project (orphan/inbox) | Already supported — `project_id` is nullable |
| PROJ-01 | Projects page — card grid with ADO badge, progress | ProjectsView enhancement |
| PROJ-02 | Pin/star favorite projects | New `is_pinned` column on projects table |
| PROJ-03 | Projects are flat (no sub-projects) | Already the case |
| PROJ-04 | Projects can be local-only or linked to ADO | Same link model as tasks — `project_ado_links` table |
| PROJ-05 | Project dashboard with stats and filtered tasks | ProjectDetailView with ADO context |
| PROJ-06 | Project progress: ADO children % + local task % | Computed from linked ADO item's children + local tasks |
| PROJ-07 | Link/unlink project to ADO scenario/deliverable | Same UX pattern as task linking |
| UX-01 | ADO work item opens in-app detail panel | ADO detail panel component |
| UX-02 | External links open in real browser | Wails `browser.OpenURL()` |
| UX-03 | Tabs + sync button + filter chips in same toolbar | PageHeader component enhancement |
| UX-04 | ADO browser shows tree view (Scenario→Deliverable→Task) | TreeNodeItem component + parent chain fetching |
| UX-05 | ADO browser supports filter chips + text search + saved query picker | ADO browser toolbar with search + ADO query execution |
</phase_requirements>

## Standard Stack

### Core (Already in Project)
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25.0 | Backend language | Project standard, proven in xl |
| Wails v3 | alpha.74 | Desktop framework | Project decision — Go→frontend bindings |
| Vue 3 | 3.x | Frontend framework | Project decision — thin shell |
| Pinia | 2.x | State management | Already used for all stores |
| modernc.org/sqlite | 1.48.0 | Pure-Go SQLite | Already in go.mod, no CGO needed |
| github.com/spf13/viper | 1.21.0 | Configuration | Already used for all config |
| github.com/zalando/go-keyring | 0.2.8 | OS credential storage | Already used for token storage |
| shadcn-vue | latest | UI components | 13+ components already installed |
| Tailwind CSS | v4 | Styling | Already configured |
| lucide-vue-next | latest | Icons | Already used throughout |

### New for Phase 2
| Library | Purpose | When to Use |
|---------|---------|-------------|
| `net/http` (stdlib) | ADO REST API client | All ADO API calls — no external HTTP lib needed |
| `encoding/json` (stdlib) | JSON marshaling for ADO API | Request/response serialization |
| `time` (stdlib) | Token TTL, sync timers | Token caching, background sync goroutine |
| `sync` (stdlib) | Mutex for token cache | Thread-safe token provider |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| net/http | go-retryablehttp | Adds retry logic but extra dependency; handle retries manually for now |
| Manual JSON Patch | go-jsonpatch lib | ADO uses specific JSON Patch format; simpler to build ops manually |

**No new Go dependencies needed.** The stdlib `net/http` + `encoding/json` are sufficient for the ADO REST client, matching the proven xl pattern.

## Architecture Patterns

### Recommended Package Structure
```
pkg/ado/
├── client.go        # HTTP client, auth header injection, base request methods
├── models.go        # WorkItem, PatchOperation, WIQLResponse, Identity types
├── query.go         # WIQL queries, batch fetch, work item retrieval
├── sync.go          # State mapping (xl↔ADO), sync logic
├── push.go          # Create/update work items in ADO (JSON Patch)
├── pull.go          # Import work items from ADO
├── comments.go      # ADO work item comments API
└── urls.go          # URL builders for ADO web/API endpoints

internal/auth/
├── auth.go          # Existing OAuth flow (keep as-is)
├── token.go         # NEW: TokenProvider interface + implementations
└── azcli.go         # NEW: az cli token provider implementation

internal/app/
├── adoservice.go    # REFACTOR: Use pkg/ado client instead of az cli subprocess
├── syncservice.go   # NEW: Background sync, conflict detection, diff generation
├── linkservice.go   # NEW: Link/promote/import/unlink operations
├── commentservice.go # NEW: Local comments + selective ADO push
└── tasks.go         # EXTEND: Add isPublic computed, enhanced filtering

internal/db/
├── db.go            # EXTEND: New tables (task_links, task_comments, project_ado_links)
├── ado.go           # EXTEND: Multi-org support, parent chain queries
├── comments.go      # NEW: Comment CRUD
└── links.go         # NEW: External links CRUD

frontend/src/
├── stores/
│   ├── ado.ts       # REFACTOR: Use new service bindings, multi-org
│   ├── sync.ts      # NEW: Sync state, conflict resolution
│   ├── tasks.ts     # EXTEND: isPublic computed, enhanced filters
│   └── projects.ts  # EXTEND: ADO linking, progress
├── views/
│   ├── TasksView.vue    # EXTEND: Personal/public badges, filters, quick-add
│   ├── AdoView.vue      # REFACTOR: Tree browser, linked status, import/link
│   ├── ProjectsView.vue # EXTEND: Card grid, ADO badges, pin/star
│   └── SettingsView.vue # EXTEND: Multi-org config UI
├── components/
│   ├── SyncConfirmDialog.vue   # NEW: Preview diff for outbound changes
│   ├── ConflictResolver.vue    # NEW: Per-field conflict resolution
│   ├── LinkDialog.vue          # NEW: Search/select ADO item to link
│   ├── PromoteDialog.vue       # NEW: Preview what will be created in ADO
│   ├── ExternalLinks.vue       # NEW: Links section with type icons
│   ├── CommentsSection.vue     # NEW: Local/public comments
│   ├── AdoTreeBrowser.vue      # NEW: Tree view component
│   ├── FilterBar.vue           # NEW: Reusable filter chips
│   └── TaskRow.vue             # EXTEND: ADO badge filled/hollow
│   └── ProjectCard.vue         # NEW: Project grid card
```

### Pattern 1: Token Provider Interface
**What:** Abstract token acquisition behind a Go interface so auth method is swappable.
**When to use:** All ADO API calls go through this.
**Example:**
```go
// Source: Derived from xl's AutoClient pattern + D-01/D-02 decisions
package auth

import (
    "sync"
    "time"
)

// TokenProvider abstracts ADO token acquisition.
type TokenProvider interface {
    // GetToken returns a valid ADO access token, refreshing if needed.
    GetToken() (string, error)
    // Name returns the provider name for logging (e.g., "az-cli", "pat").
    Name() string
}

// CachedTokenProvider wraps a TokenProvider with TTL-based caching.
type CachedTokenProvider struct {
    inner    TokenProvider
    token    string
    expiry   time.Time
    mu       sync.Mutex
    ttl      time.Duration // re-fetch buffer before actual expiry
}

func NewCachedTokenProvider(inner TokenProvider, ttl time.Duration) *CachedTokenProvider {
    return &CachedTokenProvider{inner: inner, ttl: ttl}
}

func (c *CachedTokenProvider) GetToken() (string, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.token != "" && time.Now().Before(c.expiry.Add(-c.ttl)) {
        return c.token, nil
    }
    token, err := c.inner.GetToken()
    if err != nil {
        return "", err
    }
    c.token = token
    c.expiry = time.Now().Add(50 * time.Minute) // az cli tokens last ~60min
    return token, nil
}
```

### Pattern 2: Stateless ADO Client Accepting Token
**What:** ADO client takes token per-request, no internal auth state.
**When to use:** All ADO REST API operations.
**Example:**
```go
// Source: xl pkg/ado/client.go adapted for token provider pattern
package ado

type Client struct {
    org     string
    project string
    http    *http.Client
    token   string // set per-call from TokenProvider
}

func NewClient(org, project, token string) *Client {
    return &Client{
        org:     org,
        project: project,
        token:   token,
        http:    &http.Client{Timeout: 15 * time.Second},
    }
}

// For multi-org: create one Client per org/project pair, all using same token
func NewMultiOrgClients(orgProjects []OrgProject, token string) []*Client {
    clients := make([]*Client, len(orgProjects))
    for i, op := range orgProjects {
        clients[i] = NewClient(op.Org, op.Project, token)
    }
    return clients
}
```

### Pattern 3: ADO State Mapping (Proven from XL)
**What:** Bidirectional mapping between local task statuses and ADO work item states, keyed by work item type.
**When to use:** Sync operations — both inbound (ADO→local) and outbound (local→ADO).
**Example:**
```go
// Source: xl pkg/ado/sync.go — proven in production
var statusByType = map[string]map[string]string{
    "default": {
        "todo":        "New",
        "in_progress": "Active",
        "in_review":   "Resolved",
        "done":        "Closed",
        "blocked":     "Active",  // ADO has no "Blocked" state
        "cancelled":   "Removed",
    },
    "Task": {
        "done": "Completed",  // Tasks use "Completed" not "Closed"
        // ... rest same as default
    },
}

var ReverseStatusMap = map[string]string{
    "New":       "todo",
    "Active":    "in_progress",
    "Resolved":  "in_review",
    "Closed":    "done",
    "Completed": "done",
    "Removed":   "cancelled",
    "Proposed":  "todo",
}
```

### Pattern 4: JSON Patch for ADO Updates
**What:** ADO REST API uses JSON Patch (`application/json-patch+json`) for work item modifications.
**When to use:** All outbound writes — update state, create work items.
**Example:**
```go
// Source: xl pkg/ado/client.go
type PatchOperation struct {
    Op    string `json:"op"`
    Path  string `json:"path"`
    Value any    `json:"value"`
}

// Update title + state
ops := []PatchOperation{
    {Op: "replace", Path: "/fields/System.Title", Value: "New Title"},
    {Op: "replace", Path: "/fields/System.State", Value: "Active"},
}
// POST to: /wit/workitems/$Task?api-version=7.0 (create)
// PATCH to: /wit/workitems/{id}?api-version=7.0 (update)
// Content-Type: application/json-patch+json  ← CRITICAL
```

### Pattern 5: WIQL Query + Batch Fetch
**What:** ADO WIQL returns only IDs. Must batch-fetch full work items separately.
**When to use:** Listing assigned work items, querying by parent.
**Example:**
```go
// Source: xl pkg/ado/query.go
// Step 1: WIQL query returns IDs
wiql := `SELECT [System.Id] FROM WorkItems
    WHERE [System.AssignedTo] = @Me
    AND [System.State] NOT IN ('Closed','Removed','Completed','Cut')
    ORDER BY [System.ChangedDate] DESC`
// POST to: /_apis/wit/wiql?api-version=7.0

// Step 2: Batch fetch by IDs (max 200 per batch)
// GET: /_apis/wit/workitems?ids=1,2,3&$expand=none&api-version=7.0
```

### Pattern 6: Preview Diff Before Outbound Push
**What:** Generate a structured diff of what will change before pushing to ADO.
**When to use:** D-22 — every outbound change must show preview.
**Example:**
```go
// Sync service generates diff
type SyncDiff struct {
    TaskID    int         `json:"taskId"`
    AdoID     string      `json:"adoId"`
    Changes   []FieldDiff `json:"changes"`
    Direction string      `json:"direction"` // "push" or "pull"
}

type FieldDiff struct {
    Field    string `json:"field"`    // "title", "status", "description"
    Local    string `json:"local"`    // current local value
    Remote   string `json:"remote"`   // current ADO value
    Proposed string `json:"proposed"` // what will be written
}
```

### Pattern 7: Multi-Org Aggregation
**What:** Config stores list of org/project pairs; aggregate results across all.
**When to use:** D-04b through D-04f — unified view across orgs.
**Example:**
```yaml
# config.yaml — new structure
ado:
  orgs:
    - org: "xbox"
      projects:
        - "XES"
        - "XboxLive"
    - org: "microsoft"
      projects:
        - "OneCoreUAP"
```

### Anti-Patterns to Avoid
- **Shelling out to az cli per query:** Current `ADOService.runAzCli()` pattern. Replace with token + REST.
- **Auto-pushing to ADO:** Never. Every outbound change requires explicit user confirmation via preview diff.
- **Syncing subtasks automatically:** Subtasks stay personal unless individually linked.
- **Storing token in config file:** Use keyring for tokens, not Viper config.
- **Single-org assumption in code paths:** Always iterate over org/project pairs.
- **Blocking UI on ADO API calls:** All ADO operations should be async with loading states.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| OS credential storage | Custom file-based token store | `github.com/zalando/go-keyring` | Already in project; handles Windows Credential Manager, macOS Keychain, Linux Secret Service |
| JSON Patch ops | Custom diff/patch library | Manual `PatchOperation` structs | ADO's JSON Patch is simple (3-4 fields); xl's pattern is proven and minimal |
| WIQL query building | SQL query builder abstraction | String templates with parameters | WIQL is limited; templates are clearer than a builder pattern |
| Config management | Custom config file parser | `github.com/spf13/viper` | Already used; supports YAML, env vars, defaults |
| Component library | Custom buttons/dialogs/inputs | shadcn-vue (already installed) | 13+ components already in use; use Dialog, Command, ScrollArea, Tabs, etc. |
| Tree view component | Custom recursive DOM tree | Adapt existing `TreeNodeItem.vue` | Already exists in codebase; extend for ADO hierarchy |
| URL pattern detection | Complex NLP-based categorization | Simple regex map | Known URL patterns are finite and predictable |

**Key insight:** The XL codebase has already solved the hard ADO integration problems. Port the patterns (token acquisition, WIQL, JSON Patch, state mapping, batch fetch) rather than redesigning them.

## Common Pitfalls

### Pitfall 1: ADO API Rate Limiting
**What goes wrong:** Rapid API calls during initial sync or bulk operations trigger 429 responses.
**Why it happens:** ADO enforces per-user rate limits (~200 requests/5 min for REST API). Batch fetching 200 items = 1 request, but WIQL + fetch = 2+ requests per query.
**How to avoid:** Batch work item fetches (max 200 IDs per call). Cache aggressively in SQLite. Respect `Retry-After` headers. Implement exponential backoff.
**Warning signs:** 429 status codes, increasing latency, intermittent failures.

### Pitfall 2: Token Expiry Mid-Operation
**What goes wrong:** Long-running sync operations fail partway because the token expired.
**Why it happens:** az cli tokens last ~60 minutes. A large sync could exceed this.
**How to avoid:** Check token expiry before each API call. Re-fetch token if within 5-minute buffer of expiry. The `CachedTokenProvider` pattern handles this.
**Warning signs:** 401 responses after successful initial auth.

### Pitfall 3: JSON Patch Content-Type Header
**What goes wrong:** ADO returns 400 Bad Request when creating/updating work items.
**Why it happens:** ADO requires `Content-Type: application/json-patch+json` for patch operations, not `application/json`.
**How to avoid:** Always set this header for PATCH and POST to work item endpoints.
**Warning signs:** "The request body is not a valid JSON Patch document" error.

### Pitfall 4: ADO State Transitions Are Not Free-Form
**What goes wrong:** Trying to set a work item directly to "Closed" from "New" fails.
**Why it happens:** ADO enforces state transition rules per process template (Agile, CMMI, Scrum). Not all states are reachable from all other states.
**How to avoid:** Use the xl pattern: create in default state ("New"), then transition. For complex transitions, may need intermediate states. The `statusByType` mapping handles this for common cases.
**Warning signs:** "The field 'System.State' contains an invalid value" or "VS402323: The transition from ... to ... is not valid."

### Pitfall 5: WIQL Returns IDs Only
**What goes wrong:** Assuming WIQL returns full work item data and displaying empty items.
**Why it happens:** WIQL queries return only `{ id, url }` references. You must batch-fetch the actual work items separately.
**How to avoid:** Always follow the 2-step pattern: WIQL query → extract IDs → batch GET by IDs.
**Warning signs:** Work items display with no title/state/type.

### Pitfall 6: Wails Binding Limitations with Goroutines
**What goes wrong:** Background sync goroutine can't emit events to the frontend.
**Why it happens:** Wails v3 requires the `application.App` reference to emit events.
**How to avoid:** Pass `*application.App` to SyncService. Use `app.Event.Emit()` for sync status updates. Frontend listens via Wails event API.
**Warning signs:** Frontend never updates during background sync.

### Pitfall 7: Conflict Detection Timing
**What goes wrong:** Changes appear as conflicts even when they shouldn't be.
**Why it happens:** Comparing timestamps without accounting for the last sync point. A local change made before the last sync is not a conflict.
**How to avoid:** Track `last_synced_at` per linked task. Only flag conflicts when both local `updated_at > last_synced_at` AND remote `ChangedDate > last_synced_at`.
**Warning signs:** Every linked task shows as conflicted after each sync.

### Pitfall 8: SQLite Schema Migration Without Breaking Existing Data
**What goes wrong:** Adding new tables/columns breaks the existing Phase 1 database.
**Why it happens:** The current migration runs the full schema via `CREATE TABLE IF NOT EXISTS`. New columns on existing tables need `ALTER TABLE`.
**How to avoid:** Use the existing pattern (see `db.go` line 46): `db.Exec('ALTER TABLE ... ADD COLUMN ...')` with error suppression for idempotency. Add new tables as `CREATE TABLE IF NOT EXISTS`.
**Warning signs:** "duplicate column name" errors or missing columns at runtime.

## Code Examples

### ADO REST API: Get Access Token via az cli
```go
// Source: xl pkg/ado/auth.go — proven pattern
const adoResource = "499b84ac-1321-427f-aa17-267ca6975798"

func getAzCliToken() (string, time.Time, error) {
    cmd := exec.Command("az", "account", "get-access-token",
        "--resource", adoResource,
        "--output", "json")
    out, err := cmd.Output()
    if err != nil {
        return "", time.Time{}, fmt.Errorf("az CLI failed: run `az login` first")
    }
    var resp struct {
        AccessToken string `json:"accessToken"`
        ExpiresOn   string `json:"expiresOn"`
    }
    json.Unmarshal(out, &resp)
    expiry, _ := time.Parse("2006-01-02 15:04:05.000000", resp.ExpiresOn)
    return resp.AccessToken, expiry, nil
}
```

### ADO REST API: Create Work Item (Promote)
```go
// Source: xl pkg/ado/push.go adapted for xb-tasks
func (c *Client) CreateWorkItem(wiType string, title, description, state string) (*WorkItem, error) {
    ops := []PatchOperation{
        {Op: "add", Path: "/fields/System.Title", Value: title},
        {Op: "add", Path: "/fields/System.Description", Value: description},
    }
    // Create in default state, then transition
    url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/wit/workitems/$%s?api-version=7.0",
        c.org, c.project, wiType)
    // POST with Content-Type: application/json-patch+json
    // ... then transition to target state if not "New"
}
```

### ADO REST API: Add Comment to Work Item
```go
// ADO Comments API — for CMT-02
func (c *Client) AddComment(workItemID, text string) error {
    url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/wit/workitems/%s/comments?api-version=7.0-preview.3",
        c.org, c.project, workItemID)
    body, _ := json.Marshal(map[string]string{"text": text})
    req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.token)
    resp, err := c.http.Do(req)
    // ... handle response
}
```

### New SQLite Tables for Phase 2
```sql
-- External links on tasks (LINK-01 through LINK-03)
CREATE TABLE IF NOT EXISTS task_links (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id  INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    url      TEXT NOT NULL,
    label    TEXT DEFAULT '',
    type     TEXT DEFAULT 'url', -- icm, grafana, ado, wiki, url
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_task_links_task ON task_links(task_id);

-- Local comments on tasks (CMT-01 through CMT-03)
CREATE TABLE IF NOT EXISTS task_comments (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id    INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    is_public  INTEGER DEFAULT 0,  -- 0=local/private, 1=pushed to ADO
    ado_comment_id TEXT DEFAULT '', -- ADO comment ID if pushed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_task_comments_task ON task_comments(task_id);

-- Project ADO links (PROJ-04, PROJ-07)
CREATE TABLE IF NOT EXISTS project_ado_links (
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    ado_id     TEXT NOT NULL,
    direction  TEXT DEFAULT 'linked' CHECK(direction IN ('promoted','imported','linked')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, ado_id)
);

-- Sync tracking per linked task (for conflict detection)
CREATE TABLE IF NOT EXISTS sync_state (
    task_id       INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    ado_id        TEXT NOT NULL,
    last_synced_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    local_title    TEXT DEFAULT '',
    local_status   TEXT DEFAULT '',
    local_desc     TEXT DEFAULT '',
    remote_title   TEXT DEFAULT '',
    remote_status  TEXT DEFAULT '',
    remote_desc    TEXT DEFAULT '',
    PRIMARY KEY (task_id, ado_id)
);
```

### Multi-Org Config Structure
```go
// internal/config/config.go — extend for multi-org
type OrgProject struct {
    Org      string   `mapstructure:"org"`
    Projects []string `mapstructure:"projects"`
}

// In setDefaults():
viper.SetDefault("ado.orgs", []OrgProject{})
// Keep legacy single-org for backward compat:
// viper.SetDefault("ado.organization", "")
// viper.SetDefault("ado.project", "")

// GetOrgProjects returns all configured org/project pairs
func GetOrgProjects() []OrgProject {
    var orgs []OrgProject
    viper.UnmarshalKey("ado.orgs", &orgs)
    // Backward compat: if empty but legacy single-org is set, use that
    if len(orgs) == 0 {
        org := viper.GetString("ado.organization")
        proj := viper.GetString("ado.project")
        if org != "" && proj != "" {
            orgs = []OrgProject{{Org: org, Projects: []string{proj}}}
        }
    }
    return orgs
}
```

### URL Pattern Detection (LINK-02)
```go
// internal/app/linkservice.go
var urlPatterns = []struct {
    Pattern *regexp.Regexp
    Type    string
}{
    {regexp.MustCompile(`portal\.microsofticm\.com`), "icm"},
    {regexp.MustCompile(`grafana\.|\.grafana\.`), "grafana"},
    {regexp.MustCompile(`dev\.azure\.com/.+/_workitems`), "ado"},
    {regexp.MustCompile(`dev\.azure\.com/.+/_wiki`), "wiki"},
    {regexp.MustCompile(`\.visualstudio\.com/.+/_wiki`), "wiki"},
}

func DetectLinkType(url string) string {
    for _, p := range urlPatterns {
        if p.Pattern.MatchString(url) {
            return p.Type
        }
    }
    return "url" // fallback
}
```

### Background Sync Goroutine (SYNC-01)
```go
// internal/app/syncservice.go
type SyncService struct {
    db       *db.DB
    tokenProv auth.TokenProvider
    app      *application.App
    interval time.Duration
    stopCh   chan struct{}
}

func (s *SyncService) StartBackgroundSync() {
    s.stopCh = make(chan struct{})
    go func() {
        ticker := time.NewTicker(s.interval)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                s.pullChanges() // silent inbound sync
            case <-s.stopCh:
                return
            }
        }
    }()
}

func (s *SyncService) pullChanges() {
    // 1. Get token
    // 2. For each org/project pair, create Client
    // 3. Fetch current ADO state for all linked items
    // 4. Compare with local sync_state
    // 5. Update local DB with remote changes (non-conflicting)
    // 6. Flag conflicts for user resolution
    // 7. Emit "sync:completed" event to frontend
    s.app.Event.Emit("sync:completed", map[string]any{
        "conflicts": conflictCount,
        "updated":   updatedCount,
    })
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| ADO via az cli subprocess | Direct REST API with token | Phase 2 refactor | Faster, no az cli dependency per query |
| Single org/project | Multi-org/project config | Phase 2 | Unified cross-org view |
| Tasks always local | Personal→public model | Phase 2 | ADO-synced tasks with visual distinction |
| No sync | Bidirectional sync with confirmation | Phase 2 | Safe ADO integration |
| No external links | Structured URL links on tasks | Phase 2 | Investigation hub workflow |
| No comments | Local + selective ADO push | Phase 2 | Private/public comment distinction |

**ADO REST API Version:** Use `api-version=7.0` consistently (stable, matches xl patterns). The newer `7.1-preview.*` versions add features not needed here.

**Deprecated/outdated:**
- `az boards query` subprocess pattern in current `adoservice.go` — replace entirely with REST client
- Single `ado.organization`/`ado.project` config — replace with `ado.orgs` list (backward-compat shim)

## Open Questions

1. **ADO Process Template Variations**
   - What we know: XL's state mapping works for Agile/CMMI. Xbox team uses specific process templates.
   - What's unclear: Whether all target orgs use the same process template and state transitions.
   - Recommendation: Use xl's `statusByType` mapping as default. If transition fails, log the error and let user know. Could add custom mapping in config later.

2. **ADO Work Item Tree Depth**
   - What we know: D-14 wants Scenario → Deliverable → Task/Bug/Story tree. D-15 says sync items assigned to me + parent chain.
   - What's unclear: How deep the parent chain can go in practice. Fetching each parent is an additional API call.
   - Recommendation: Fetch assigned items first, then batch-fetch their parents (up to 3 levels). Cache the tree in SQLite. Don't recurse indefinitely.

3. **Multi-Org Token Scope**
   - What we know: D-04a says single token for all orgs (same Azure tenant).
   - What's unclear: Whether the az cli token scope (`499b84ac-...` resource) covers all ADO orgs in the tenant.
   - Recommendation: It should — the resource ID is the ADO API resource, not org-specific. Test with first configured org; flag error clearly if token doesn't work for other orgs.

4. **Saved ADO Query Picker (UX-05)**
   - What we know: ADO has a saved queries API (`/_apis/wit/queries`).
   - What's unclear: Whether users actually have saved queries and how the picker UI should work.
   - Recommendation: Build as a dropdown that lists saved queries. Executing a query returns work item IDs which feed into the batch-fetch pattern. Lower priority — implement after core browser works.

## Sources

### Primary (HIGH confidence)
- XL codebase (`XL-CODE-REFERENCE.md`) — complete ADO client implementation, state mapping, WIQL patterns, auth patterns
- Existing xb-tasks codebase — current schema, services, stores, UI components
- CONTEXT.md decisions D-01 through D-32 — locked user decisions

### Secondary (MEDIUM confidence)
- Azure DevOps REST API (api-version 7.0) — endpoint patterns, JSON Patch format, WIQL query syntax
- ADO resource ID `499b84ac-1321-427f-aa17-267ca6975798` — from xl's proven auth flow

### Tertiary (LOW confidence)
- ADO rate limit specifics (~200 req/5min) — from general ADO documentation, exact limits may vary by tier
- ADO state transition rules — process-template dependent; may need runtime discovery

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all libraries already in project; no new dependencies needed
- Architecture: HIGH — patterns proven in xl codebase, clear migration path from current code
- ADO API patterns: HIGH — xl reference has complete working client code
- Sync/conflict resolution: MEDIUM — logic is clear but edge cases need runtime validation
- Frontend architecture: MEDIUM — UI components exist but significant new views needed
- Multi-org: MEDIUM — pattern is clear but untested across multiple real ADO orgs
- Pitfalls: HIGH — documented from xl's production experience

**Research date:** 2026-04-04
**Valid until:** 2026-05-04 (stable domain — ADO API 7.0 is mature, Go stdlib is stable)
