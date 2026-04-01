# Architecture Research

**Domain:** Team developer productivity dashboard with Azure DevOps integration
**Researched:** 2025-07-17
**Confidence:** HIGH

## Standard Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTS                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────┐  │
│  │  Web App      │  │  VS Code Ext │  │  MCP / Copilot Chat  │  │
│  │  (React SPA)  │  │  (Future)    │  │  (Future)            │  │
│  └──────┬───────┘  └──────┬───────┘  └───────────┬───────────┘  │
│         │                  │                      │              │
│         └──────────────────┼──────────────────────┘              │
│                            │                                     │
│                     REST API + SSE                               │
├────────────────────────────┼─────────────────────────────────────┤
│                       API SERVER                                 │
│                                                                  │
│  ┌────────────┐  ┌──────────────┐  ┌──────────────────────────┐ │
│  │  Auth       │  │  Task        │  │  ADO Sync Engine         │ │
│  │  Module     │  │  Manager     │  │  ┌────────┐ ┌─────────┐ │ │
│  │  (MSAL/     │  │  (Personal + │  │  │ Poller │ │ Webhook │ │ │
│  │   Entra ID) │  │   Linked)    │  │  │        │ │ Recv.   │ │ │
│  └──────┬─────┘  └──────┬───────┘  │  └────┬───┘ └────┬────┘ │ │
│         │               │          │       │          │       │ │
│         │               │          └───────┼──────────┼───────┘ │
│         │               │                  │          │         │
│  ┌──────┴───────────────┴──────────────────┴──────────┴───────┐ │
│  │                    SSE Event Hub                             │ │
│  │             (broadcasts changes to clients)                 │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│                       DATA LAYER                                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                     PostgreSQL                            │   │
│  │  ┌──────────┐  ┌───────────┐  ┌────────────┐            │   │
│  │  │ Personal │  │ ADO Cache │  │ User/Auth  │            │   │
│  │  │ Tasks    │  │ (WI, PR,  │  │ Profiles   │            │   │
│  │  │          │  │  Comments) │  │ + Tokens   │            │   │
│  │  └──────────┘  └───────────┘  └────────────┘            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│                    EXTERNAL SERVICES                             │
│  ┌──────────────────┐  ┌──────────────────┐                     │
│  │  Azure DevOps     │  │  Microsoft       │                     │
│  │  REST API v7.1    │  │  Entra ID        │                     │
│  │  (Work Items,     │  │  (OAuth 2.0      │                     │
│  │   PRs, Comments,  │  │   Auth Code +    │                     │
│  │   Service Hooks)  │  │   PKCE)          │                     │
│  └──────────────────┘  └──────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| **Web App** | Dashboard UI — unified view of tasks, work items, PRs, comments | React SPA with SSE subscription for real-time updates |
| **Auth Module** | Microsoft Entra ID login, token management, refresh, per-user ADO API credentials | MSAL Node (server) + MSAL Browser (client), OAuth 2.0 auth code + PKCE |
| **Task Manager** | CRUD for personal tasks, linking to ADO items, promotion/import flows | REST endpoints, business logic for task ↔ ADO item bidirectional sync |
| **ADO Sync Engine** | Keeps local ADO cache fresh, detects changes, manages sync state | Background polling workers + optional webhook receiver |
| **ADO Poller** | Periodic WIQL queries + batch work item fetches per user | Per-user polling intervals (60s active, 5min idle), smart change detection |
| **Webhook Receiver** | Receives ADO Service Hook HTTP POST events for instant updates | Express endpoint, validates payloads, updates cache, triggers SSE |
| **SSE Event Hub** | Pushes real-time updates to connected browser/extension clients | In-memory pub/sub with per-user channels, Server-Sent Events protocol |
| **PostgreSQL** | Persistent storage for tasks, cached ADO data, user profiles, encrypted tokens | Multi-tenant single DB with user_id foreign keys throughout |

## Recommended Project Structure

```
team-ado-tool/
├── apps/
│   ├── web/                    # React SPA (Vite)
│   │   ├── src/
│   │   │   ├── components/     # UI components
│   │   │   ├── hooks/          # React hooks (useSSE, useAuth, useTasks)
│   │   │   ├── pages/          # Route pages (Dashboard, PRs, Tasks)
│   │   │   ├── stores/         # Client state (Zustand or similar)
│   │   │   └── api/            # API client (typed, auto-generated)
│   │   └── index.html
│   └── server/                 # Node.js API server
│       ├── src/
│       │   ├── routes/         # Express/Hono route handlers
│       │   ├── services/       # Business logic (TaskService, AdoSyncService)
│       │   ├── ado/            # ADO API client wrapper
│       │   │   ├── client.ts   # REST client (typed)
│       │   │   ├── poller.ts   # Background polling logic
│       │   │   ├── webhooks.ts # Webhook receiver + validation
│       │   │   ├── wiql.ts     # WIQL query builder
│       │   │   └── mappers.ts  # ADO ↔ local data mapping
│       │   ├── auth/           # MSAL config, token management
│       │   ├── sse/            # SSE hub, event broadcasting
│       │   ├── db/             # Database access (Drizzle/Kysely)
│       │   │   ├── schema.ts   # Table definitions
│       │   │   └── migrations/ # Database migrations
│       │   └── middleware/     # Auth, error handling, logging
│       └── index.ts
├── packages/
│   └── shared/                 # Shared TypeScript types
│       ├── types/              # Task, WorkItem, PR, Comment types
│       └── constants/          # Status mappings, ADO field names
├── package.json                # Workspace root (pnpm/npm workspaces)
└── docker-compose.yml          # PostgreSQL + app for local dev
```

### Structure Rationale

- **apps/web/ + apps/server/:** Clear client/server split in a monorepo. Both TypeScript, shared types from `packages/shared`. Monorepo enables shared type safety without publishing packages. Separating web from server (vs. a Next.js BFF monolith) keeps the API consumable by VS Code extension and future MCP server without routing through a web framework.
- **apps/server/src/ado/:** Isolated ADO integration layer. Every ADO interaction goes through this module. If ADO API changes or a different provider is added, changes are contained here.
- **apps/server/src/services/:** Business logic separated from route handlers. Services own the rules (e.g., "promoting a task creates an ADO work item and links them"). Routes are thin — parse request, call service, return response.
- **packages/shared/:** TypeScript types and constants shared between web and server. Prevents drift between what the API returns and what the UI expects. Critical for the VS Code extension later — it imports from the same package.

## Architectural Patterns

### Pattern 1: ADO as Source of Truth + Local Overlay

**What:** ADO owns canonical state for work items and PRs. The local database caches ADO data and layers personal metadata on top (priority overrides, personal notes, custom tags). Personal tasks that exist only locally are a separate entity type that can optionally be "promoted" to ADO.

**When to use:** Always — this is the core data philosophy.

**Trade-offs:**
- ✅ Never fights ADO. Users edit work items in ADO or in-app; changes flow correctly either way.
- ✅ Personal data (notes, priority overrides) is never pushed to ADO unless explicitly promoted.
- ⚠️ Cache staleness — users might see stale ADO data between polling intervals.
- ⚠️ Conflict detection needed when both ADO and local have changed.

**Example:**
```typescript
// Core data model concept
interface UnifiedTask {
  id: string;                     // Local UUID
  source: 'local' | 'ado';       // Where this task originated
  
  // ADO-sourced fields (null for local-only tasks)
  adoId?: number;                 // ADO work item ID
  adoUrl?: string;
  adoState?: string;              // e.g., "Active", "Resolved"
  adoRev?: number;                // ADO revision for conflict detection
  adoAssignedTo?: string;
  adoAreaPath?: string;
  adoIterationPath?: string;
  adoType?: string;               // "Task", "Bug", "User Story"
  
  // Shared fields (local values win for local tasks, ADO wins for ADO tasks)
  title: string;
  description?: string;
  status: TaskStatus;             // Mapped from ADO state or set locally
  
  // Personal overlay (always local, never pushed to ADO)
  personalPriority?: 'P0' | 'P1' | 'P2' | 'P3';
  personalNotes?: string;
  personalTags?: string[];
  
  // Sync metadata
  syncStatus: 'local_only' | 'synced' | 'local_modified' | 'remote_modified' | 'conflict';
  lastSyncedAt?: Date;
  lastLocalModifiedAt: Date;
  lastRemoteModifiedAt?: Date;
  
  // Ownership
  userId: string;                 // Every task belongs to a user
}
```

### Pattern 2: Hybrid Sync (Polling Primary, Webhooks Accelerator)

**What:** Background polling is the **primary** sync mechanism. ADO Service Hooks webhooks are an **optional accelerator** that provides near-instant updates when the server has a reachable endpoint.

**When to use:** This is the recommended sync architecture for this tool.

**Why polling primary (not webhooks):**
- ADO webhooks **cannot target localhost** — confirmed in official docs. During development and for internal deployments without a public URL, webhooks are unavailable.
- Polling is reliable and self-healing. If a webhook is missed, polling catches it.
- At 15-50 users, polling is well within ADO rate limits (200 TSTUs per user per 5-minute window).
- Webhooks require a publicly reachable URL and Service Hook subscription management — added complexity that may not be viable for an internal corporate tool.

**Trade-offs:**
- ✅ Works everywhere — no public endpoint required.
- ✅ Self-healing — missed webhooks are caught by next poll.
- ✅ Rate limits are per-user (each user's OAuth token has its own TSTU budget).
- ⚠️ Polling adds latency (30s-5min depending on interval).
- ⚠️ Wasted API calls when nothing changes (mitigated by smart polling).

**Polling strategy:**
```typescript
// Smart polling with adaptive intervals
interface PollConfig {
  activeInterval: 60_000;     // 60s when user is active on dashboard
  idleInterval: 300_000;      // 5min when user is idle
  backgroundInterval: 900_000; // 15min when user is offline (keep cache warm)
}

// Efficient change detection flow:
// 1. WIQL query: "SELECT [System.Id], [System.Rev] FROM WorkItems 
//    WHERE [System.AssignedTo] = @Me AND [System.ChangedDate] > @lastSync"
// 2. Compare returned revisions against cached revisions
// 3. Batch-fetch only changed items: POST _apis/wit/workitemsbatch
// 4. Update cache, emit SSE events for changed items
```

### Pattern 3: Server-Sent Events for Real-Time Client Updates

**What:** The server pushes updates to connected clients via SSE. When the ADO poller detects changes or a webhook arrives, the server broadcasts events through per-user SSE channels.

**When to use:** For all real-time dashboard updates.

**Trade-offs:**
- ✅ Simpler than WebSockets — HTTP-based, auto-reconnect built into browsers, works through corporate proxies.
- ✅ Unidirectional (server → client) matches our use case. Clients send actions via REST POST.
- ✅ `Last-Event-ID` header enables resuming after disconnect without missing events.
- ⚠️ One-directional only — but we use REST for client → server, so this is fine.
- ⚠️ Limited to ~6 concurrent connections per domain in HTTP/1.1 (not an issue with HTTP/2).

**Example:**
```typescript
// SSE event types
type SSEEvent = 
  | { type: 'work-item-updated'; data: { adoId: number; fields: Partial<WorkItem> } }
  | { type: 'pr-status-changed'; data: { prId: number; status: string; reviewers: Reviewer[] } }
  | { type: 'pr-comment-added'; data: { prId: number; threadId: number; comment: Comment } }
  | { type: 'task-synced'; data: { taskId: string; syncStatus: SyncStatus } }
  | { type: 'sync-error'; data: { message: string; retryAt: Date } };

// Server-side hub (in-memory, no Redis needed at this scale)
class SSEHub {
  private connections: Map<string, Set<Response>>; // userId → active SSE responses
  
  broadcast(userId: string, event: SSEEvent) {
    const conns = this.connections.get(userId);
    conns?.forEach(res => {
      res.write(`event: ${event.type}\n`);
      res.write(`data: ${JSON.stringify(event.data)}\n`);
      res.write(`id: ${this.nextEventId()}\n\n`);
    });
  }
}
```

### Pattern 4: On-Behalf-Of Token Flow for ADO API Calls

**What:** The server makes ADO API calls using each user's delegated OAuth token. Users authenticate with Microsoft Entra ID; the server stores their refresh tokens and uses access tokens to call ADO on their behalf.

**When to use:** All ADO API interactions.

**Why not a service account / PAT:**
- Rate limits are per-identity. A shared service account would hit limits with 50 users.
- Users see only what they have ADO permissions for — no accidental data exposure.
- Matches ADO's intended auth model for delegated apps.
- PATs expire and must be manually rotated per user — OAuth refresh tokens are automatic.

**Scopes needed:**
```
vso.work_write    — Read/create/update work items and queries
vso.code          — Read source code, PRs, branches
vso.threads_full  — Read/write PR comment threads
vso.profile       — Read user profile
vso.graph         — Read team/group membership
```

**Token lifecycle:**
```
User login → Entra ID auth code + PKCE → Server exchanges for access + refresh tokens
  → Access token (1hr): stored in memory, used for ADO API calls
  → Refresh token (long-lived): stored encrypted in PostgreSQL
  → Background: auto-refresh access token before expiry
  → Poller: uses user's access token for their ADO queries
```

**Entra ID resource URI for ADO:** `499b84ac-1321-427f-aa17-267ca6975798`
Use `.default` scope when requesting token to get all configured permissions.

## Data Flow

### Dashboard Load Flow

```
User opens dashboard
    │
    ▼
Browser ──GET /api/auth/me──▶ Server validates session cookie
    │                              │
    │                              ▼
    │                         Return user profile + last sync timestamp
    │
    ▼
Browser ──GET /api/dashboard──▶ Server queries PostgreSQL cache
    │                              │
    │                              ├── Personal tasks (source = 'local')
    │                              ├── Linked ADO items (cached, with overlay)
    │                              ├── Active PRs (cached)
    │                              └── Recent comments (cached)
    │                              │
    │                              ▼
    │                         Return unified dashboard payload
    │
    ▼
Browser ──GET /api/events (SSE)──▶ Server registers SSE connection
    │                                   │
    │                                   ▼
    │                              Poller runs in background
    │                              Detects ADO changes → SSE push
    │
    ▼
Dashboard renders, auto-updates via SSE
```

### ADO Sync Flow (Polling)

```
Background Poller (per-user, on interval)
    │
    ▼
WIQL query: "Changed items since last sync"
    │
    ▼
ADO API ──response──▶ List of (id, rev) pairs
    │
    ▼
Compare against cached revisions in PostgreSQL
    │
    ├── No changes → Sleep until next interval
    │
    ├── Changes detected:
    │       │
    │       ▼
    │   Batch fetch changed items: POST _apis/wit/workitemsbatch
    │       │                       (up to 200 IDs per call)
    │       ▼
    │   Update PostgreSQL cache (upsert work items)
    │       │
    │       ▼
    │   Check for conflicts with local modifications
    │       │
    │       ├── No conflict → Mark synced
    │       └── Conflict → Mark conflict, user resolves in UI
    │       │
    │       ▼
    │   Broadcast SSE events to connected clients
    │
    └── Rate limited (429):
            │
            ▼
        Read Retry-After header → Back off → Retry
```

### Task Promotion Flow (Local → ADO)

```
User clicks "Push to ADO" on a local task
    │
    ▼
Browser ──POST /api/tasks/:id/promote──▶ Server
    │                                         │
    │                                         ▼
    │                                    Map local task → ADO JSON Patch
    │                                    POST _apis/wit/workitems/$Task
    │                                         │
    │                                         ▼
    │                                    ADO creates work item, returns ID + rev
    │                                         │
    │                                         ▼
    │                                    Update local task:
    │                                      source: 'local' → 'ado'
    │                                      adoId: <new ID>
    │                                      adoRev: <rev>
    │                                      syncStatus: 'synced'
    │                                         │
    │                                         ▼
    │                                    Broadcast SSE: task-synced
    │
    ▼
UI updates task card with ADO badge + link
```

### Authentication Flow

```
User visits app (no session)
    │
    ▼
Redirect to Microsoft Entra ID
    │  Authorization Code Flow + PKCE
    │  Resource: 499b84ac-1321-427f-aa17-267ca6975798
    │  Scopes: openid profile vso.work_write vso.code vso.threads_full vso.graph
    │
    ▼
User consents → Redirect back with auth code
    │
    ▼
Server exchanges code for tokens (MSAL Node)
    │
    ├── access_token → In-memory cache (1hr TTL)
    ├── refresh_token → Encrypted in PostgreSQL
    └── id_token → Extract user identity (OID, email, name)
    │
    ▼
Create/update user record in PostgreSQL
    │
    ▼
Set httpOnly session cookie → Redirect to dashboard
    │
    ▼
Subsequent requests: Cookie → Session → User → Cached access token
    │
    ▼
Token expired? → MSAL auto-refresh using stored refresh token
```

### Key Data Flows Summary

1. **ADO → Dashboard (read path):** Poller queries ADO → updates PostgreSQL cache → SSE push to client → React state update → UI re-render. Latency: 0-60s depending on poll interval.
2. **Dashboard → ADO (write path):** User action → REST POST → Server validates → ADO API call with user's token → Update local cache from response → SSE confirmation. Latency: 1-3s (ADO API roundtrip).
3. **Personal task CRUD:** User action → REST POST → PostgreSQL insert/update → SSE push. No ADO involvement. Latency: <100ms.
4. **Webhook accelerated path (when available):** ADO event → HTTP POST to webhook endpoint → Validate + parse → Update cache → SSE push. Latency: 1-5s from ADO event.

## Database Schema (Conceptual)

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   users      │     │  tasks            │     │  ado_work_items  │
│─────────────│     │──────────────────│     │─────────────────│
│ id (PK)      │◄──┐│ id (PK)           │     │ ado_id (PK)      │
│ entra_oid    │   ││ user_id (FK)      │──┐  │ user_id (FK)     │
│ email        │   ││ title             │  │  │ type             │
│ display_name │   ││ description       │  │  │ title            │
│ refresh_token│   ││ status            │  │  │ state            │
│ (encrypted)  │   ││ personal_priority │  │  │ assigned_to      │
│ last_login   │   ││ personal_notes    │  │  │ area_path        │
│ settings     │   ││ personal_tags     │  │  │ iteration_path   │
└─────────────┘   ││ source            │  │  │ priority         │
                   ││ ado_work_item_id  │──┼─▶│ rev              │
                   ││ sync_status       │  │  │ raw_fields (JSON)│
                   │└──────────────────┘  │  │ last_fetched_at  │
                   │                       │  └─────────────────┘
                   │ ┌──────────────────┐  │
                   │ │  ado_pull_reqs    │  │  ┌─────────────────┐
                   │ │──────────────────│  │  │  ado_pr_threads  │
                   └─│ user_id (FK)     │  │  │─────────────────│
                     │ pr_id (PK)       │  │  │ id (PK)          │
                     │ repository       │  │  │ pr_id (FK)       │
                     │ title            │  │  │ thread_id        │
                     │ status           │  │  │ status           │
                     │ source_branch    │  │  │ last_fetched_at  │
                     │ target_branch    │  │  └────────┬────────┘
                     │ created_by       │  │           │
                     │ reviewers (JSON) │  │  ┌────────▼────────┐
                     │ vote_status      │  │  │  ado_pr_comments │
                     │ last_fetched_at  │  │  │─────────────────│
                     └──────────────────┘  │  │ id (PK)          │
                                           │  │ thread_id (FK)   │
                     ┌──────────────────┐  │  │ author           │
                     │  sync_state       │  │  │ content          │
                     │──────────────────│  │  │ published_date   │
                     │ user_id (FK)     │──┘  │ parent_id        │
                     │ entity_type      │     │ last_fetched_at  │
                     │ last_polled_at   │     └─────────────────┘
                     │ last_changed_at  │
                     │ poll_cursor      │
                     └──────────────────┘
```

**Key design decisions:**
- **`ado_work_items` is a cache table**, not a source of truth. It's refreshed by polling and can be fully rebuilt from ADO.
- **`tasks` table stores personal tasks AND personal overlays** for ADO items. A task with `source = 'ado'` and `ado_work_item_id` set is a linked overlay. A task with `source = 'local'` is purely personal.
- **`sync_state` tracks polling cursors** per user per entity type. Enables efficient delta queries ("what changed since X?").
- **`raw_fields` JSON column** on `ado_work_items` stores the full ADO response. Avoids needing to model every ADO field as a column — extract what's needed and keep the rest for future use.
- **PR comments stored in their own table** with thread hierarchy. ADO PR threads have status (active, fixed, won't fix, closed) which is tracked separately from individual comments.

## Scaling Considerations

| Scale | Architecture Adjustments |
|-------|--------------------------|
| 1 user (dogfooding) | Single Node.js process, in-memory SSE hub. PostgreSQL still recommended — avoids migration pain later. |
| 15-50 users (team) | Same single process handles this fine. Each user has own polling worker (lightweight async intervals). Connection pool of 10-20. |
| 50-200 users (org) | Consider Redis for SSE pub/sub if deploying multiple server instances. Move polling workers to a separate process/queue. |
| 200+ users | Not a design target. Would need: Redis, worker queues (BullMQ), horizontal scaling, webhook-primary strategy. |

### Scaling Priorities

1. **First bottleneck: ADO API rate limits.** Each user gets 200 TSTUs per 5-min sliding window. At 60s polling, a user consumes ~5-10 TSTUs per cycle (1 WIQL query + 1 batch fetch). That's 25-50 TSTUs per 5 min — well within limits. **No action needed at target scale.**
2. **Second bottleneck: PostgreSQL connections.** At 50 users with a connection pool of 20, this is fine. Node.js is single-threaded so pool contention is low.
3. **Third bottleneck: SSE connections.** At 50 concurrent SSE connections, trivial for Node.js. Only becomes a concern at thousands.

## Anti-Patterns

### Anti-Pattern 1: Treating Cached ADO Data as Locally Mutable

**What people do:** Store ADO work items in the local DB and start making local-only edits to ADO-sourced fields (status, title, assigned_to) without syncing back.
**Why it's wrong:** The local DB drifts from ADO. Other team members see different data. Sync becomes impossible to reason about.
**Do this instead:** ADO-sourced fields in the cache are **read-only locally**. Edits to ADO fields go through the ADO API, then the cache is updated from the response. Personal overlay fields (`personalPriority`, `personalNotes`) are the only locally-owned fields on ADO items.

### Anti-Pattern 2: Single Service Account for All ADO Calls

**What people do:** Use one PAT or service principal for all ADO API calls.
**Why it's wrong:** Rate limits are shared across all users (200 TSTUs total, not per user). Permission boundaries collapse — every user sees everything the service account can access. ADO audit trail shows one identity for all actions.
**Do this instead:** Use delegated OAuth (on-behalf-of). Each user's token has its own rate limit budget and permission scope.

### Anti-Pattern 3: WebSocket Complexity for a Dashboard

**What people do:** Implement full WebSocket infrastructure with reconnection logic, heartbeats, binary framing, load balancer affinity.
**Why it's wrong:** A dashboard is read-heavy with server-push semantics. WebSockets add complexity (connection management, proxy configuration in corporate environments) for no benefit over SSE.
**Do this instead:** SSE for server → client pushes. REST for client → server actions. Simple, debuggable, works through corporate proxies without special config.

### Anti-Pattern 4: Live ADO API Calls on Every Page Load

**What people do:** When a user loads the dashboard, make live ADO API calls to fetch current data.
**Why it's wrong:** ADO API calls take 200-500ms each. A dashboard needs work items, PRs, and comments — that's 1-2s of waterfall API calls per page load. Users navigating the dashboard repeatedly hammer ADO.
**Do this instead:** Background poller keeps the cache fresh. Dashboard reads from PostgreSQL cache (<10ms). SSE pushes updates. Users get instant page loads with eventually-consistent data.

### Anti-Pattern 5: Bidirectional Sync Without Revision Tracking

**What people do:** Push local changes to ADO and pull ADO changes without tracking revision numbers.
**Why it's wrong:** If a user edits a work item locally while someone else edits it in ADO, the last write wins silently. Data is lost.
**Do this instead:** Track `adoRev` (revision number) for every cached work item. Before pushing a local change, check if `adoRev` matches current ADO state. If not, flag a conflict for the user to resolve. ADO's JSON Patch operations support `test` operations for optimistic concurrency.

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| **Azure DevOps REST API v7.1** | OAuth 2.0 delegated tokens, JSON over HTTPS | Rate limit: 200 TSTUs / 5min / user. Batch endpoint for work items (up to 200 IDs per call). WIQL for queries. JSON Patch for updates. |
| **ADO Service Hooks** | HTTP POST webhooks to server endpoint | Events: work item created/updated/deleted/commented, PR created/updated/commented/merge-attempted. **Cannot target localhost** — only works with routable server URL. Filters available per subscription. |
| **Microsoft Entra ID** | OAuth 2.0 Authorization Code + PKCE via MSAL | Resource URI: `499b84ac-1321-427f-aa17-267ca6975798`. Use `.default` scope for all configured permissions. MSAL handles token caching and refresh. |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| Web App ↔ API Server | REST (JSON) + SSE | Shared TypeScript types ensure contract safety. API versioned from day one (`/api/v1/`). |
| API Routes ↔ Services | Direct function calls | Services are injected into routes. Services own business logic; routes own HTTP concerns only. |
| Services ↔ ADO Client | Async function calls | ADO client is a thin typed wrapper around `fetch`. Services decide *what* to sync; ADO client handles *how*. |
| Poller ↔ SSE Hub | In-process event emitter | Poller detects changes → emits events → SSE hub broadcasts. At scale, replace with Redis pub/sub. |
| Server ↔ PostgreSQL | Connection pool (pg + Drizzle/Kysely) | Pool size 10-20 sufficient for 50 users. All queries through typed query builder — no raw SQL outside migrations. |

## Build Order (Dependencies)

The architecture has clear dependency layers that dictate build order:

### Phase 1: Foundation (must come first)
1. **Auth (Entra ID + MSAL)** — Everything depends on authenticated users. Without auth, no ADO API calls, no per-user data.
2. **Database schema + migrations** — Core tables: `users`, `tasks`, `sync_state`. ADO cache tables can come later.
3. **Basic API server skeleton** — Express/Hono with auth middleware, health check, CORS.

### Phase 2: Personal Tasks (no ADO dependency)
4. **Task CRUD API** — Create, read, update, delete personal tasks. This works standalone without ADO.
5. **Basic Web Dashboard** — Render personal tasks. Users get value immediately.

### Phase 3: ADO Read Integration
6. **ADO REST client** — Typed wrapper for work items, PRs, comments endpoints.
7. **ADO Poller** — Background polling with WIQL + batch fetch. Populate cache tables.
8. **Unified Dashboard** — Show ADO items alongside personal tasks. This is the core value proposition.

### Phase 4: ADO Write + Bidirectional Sync
9. **Task ↔ ADO linking** — Import ADO item as linked task. Promote local task to ADO work item.
10. **Conflict detection** — Track revisions, detect and surface conflicts.
11. **SSE real-time updates** — Push poller changes to connected dashboards.

### Phase 5: PR + Comments Surface
12. **PR listing and detail** — Show PRs assigned/created with reviewer status.
13. **PR comment threads** — Deep integration showing threaded comments in-app.

### Phase 6: Acceleration + Polish
14. **Webhook receiver** — Accept ADO Service Hook events for near-instant updates (requires routable URL).
15. **Search** — PostgreSQL full-text across tasks + cached ADO data.
16. **Team views** — See team member work items, capacity, shared dashboards.

### Phase 7 (future): VS Code Extension
17. **VS Code extension** — Consumes the same REST API. Shared types from `packages/shared`.

**Build order rationale:**
- Auth → Database → API skeleton is the universal foundation. Nothing works without it.
- Personal tasks work without ADO — users get value from day 1 even before ADO integration ships.
- Read before write — seeing ADO data is valuable before bidirectional sync adds complexity.
- SSE is tied to the write/sync phase because that's when real-time matters (changes flowing back).
- Webhooks are an **accelerator**, not core — polling provides the baseline, webhooks make it faster.
- VS Code extension is a second client consuming the same API — the API must be stable first.

## Sources

- Azure DevOps REST API v7.1 — https://learn.microsoft.com/en-us/rest/api/azure/devops/ (HIGH confidence)
- ADO Service Hooks events — https://learn.microsoft.com/en-us/azure/devops/service-hooks/events?view=azure-devops (HIGH confidence)
- ADO Webhooks — https://learn.microsoft.com/en-us/azure/devops/service-hooks/services/webhooks?view=azure-devops — confirmed: "Webhooks can't target localhost (loopback) or special range IPv4/IPv6 addresses" (HIGH confidence)
- ADO Rate Limits — https://learn.microsoft.com/en-us/azure/devops/integrate/concepts/rate-limits?view=azure-devops — 200 TSTUs per 5-min sliding window per user (HIGH confidence)
- ADO OAuth Scopes — https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/oauth?view=azure-devops (HIGH confidence)
- Entra ID OAuth for ADO — https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/entra-oauth?view=azure-devops — "recommended" over legacy ADO OAuth (HIGH confidence)
- WIQL Query API — https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql?view=azure-devops-rest-7.1 (HIGH confidence)
- Work Items Batch API — https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-items-batch?view=azure-devops-rest-7.1 — up to 200 IDs per call (HIGH confidence)
- ADO Service Hooks Subscriptions — https://learn.microsoft.com/en-us/azure/devops/service-hooks/create-subscription?view=azure-devops (HIGH confidence)
- xl TUI proven patterns — PROJECT.md context: ADO sync, state mapping, WIQL queries, bidirectional linking (HIGH confidence — firsthand)

---
*Architecture research for: Team ADO developer dashboard*
*Researched: 2025-07-17*
