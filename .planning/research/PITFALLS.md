# Pitfalls Research

**Domain:** Developer productivity dashboard with Azure DevOps integration
**Researched:** 2025-07-14
**Confidence:** HIGH (verified against official Microsoft docs)

## Critical Pitfalls

### Pitfall 1: WIQL Returns Only IDs — Not Full Work Items

**What goes wrong:**
Developers assume the WIQL query endpoint returns complete work item data (title, state, assignedTo, etc.). It does not. The WIQL API (`POST /_apis/wit/wiql`) returns **only work item IDs**, regardless of which fields you include in the `SELECT` statement. You must make a second call to `GET /_apis/wit/workitems?ids=...` to get actual field data.

**Why it happens:**
The WIQL syntax looks like SQL with `SELECT [System.Title], [System.State]...`, which strongly implies the response includes those fields. The official docs explicitly warn: *"The API only returns work item IDs, regardless of which fields you include in the SELECT statement."*

**How to avoid:**
- Always implement a two-step fetch: (1) WIQL query → get IDs, (2) Batch work items API → get fields
- The batch work items API accepts up to 200 IDs per call (`GET /_apis/wit/workitems?ids=1,2,3...&$expand=relations`)
- Use the `fields` parameter to request only what you need, reducing payload size
- Cache work item data locally with a `rev` (revision) field to skip unchanged items on subsequent fetches

**Warning signs:**
- API responses that seem empty or missing expected fields
- Building a UI that requires a separate API call for every single work item
- Page load times growing linearly with work item count

**Phase to address:**
Phase 1 (Core ADO Integration) — This is foundational. Get the two-step fetch pattern right from day one or the entire sync layer will be wrong.

---

### Pitfall 2: Rate Limits Use Opaque "TSTUs" — Not Simple Request Counts

**What goes wrong:**
Developers implement rate limiting based on request count (e.g., "stay under 100 requests/minute"). ADO doesn't work that way. Azure DevOps uses **Azure DevOps Throughput Units (TSTUs)** — an abstract blend of database DTUs, CPU, memory, and storage bandwidth. A complex WIQL query costs vastly more TSTUs than a simple `GET` by ID. You can hit limits with 10 expensive queries while 200 cheap ones are fine.

**Why it happens:**
Most REST APIs have simple rate limits (X requests per minute). ADO's TSTU system is intentionally opaque — *"You can't calculate usage in TSTUs for an action with a formula"* (official docs). The global limit is 200 TSTUs per 5-minute sliding window, but you can't predict what a TSTU costs until you hit the limit.

**How to avoid:**
- **Honor `Retry-After` header** — this is the primary mechanism. When ADO returns it, wait the specified seconds before retrying
- **Monitor `X-RateLimit-Remaining` and `X-RateLimit-Limit` headers** proactively to detect approaching thresholds
- **Batch changes** — single-field updates are the "leading cause of performance issues and rate limiting" per Microsoft
- **Use reporting APIs** (`/_apis/wit/reporting/workitemrevisions`) for bulk data retrieval instead of WIQL + individual fetches
- **Never poll aggressively** — use service hooks (webhooks) for real-time updates instead

**Warning signs:**
- HTTP 429 responses with `TF400733` error codes
- Emails from ADO about usage limits (sent to user or Project Collection Administrators)
- `X-RateLimit-Remaining` approaching 0 in response headers
- Users reporting the dashboard becoming sluggish during peak hours

**Phase to address:**
Phase 1 (Core ADO Integration) — Build the rate-limit-aware HTTP client from the start. Retrofit is painful because every API call needs to be instrumented.

---

### Pitfall 3: Bidirectional Sync Creates Infinite Update Loops

**What goes wrong:**
Local task updated → syncs to ADO → ADO webhook fires "work item updated" → tool pulls change → detects "change" → syncs back to ADO → ADO fires webhook again → infinite loop. Even without webhooks, polling-based sync can create ping-pong updates where each sync cycle "discovers" the change it just made.

**Why it happens:**
The tool and ADO are both sources of truth for the same data. Any write to one side triggers a read-detect-write cycle on the other. Without loop detection, the system oscillates.

**How to avoid:**
- **Track sync provenance**: Every update record must store whether the change originated locally or from ADO. Never re-sync a change you just pushed.
- **Use revision numbers**: ADO work items have a `rev` field that increments on every update. Before pushing, check if the current `rev` matches what you last saw. If `rev` is exactly `your_last_pushed_rev + 1`, you made the only change — skip.
- **Debounce sync windows**: Don't sync immediately after writing. Wait 5-10 seconds to let ADO settle, then read back the authoritative state.
- **Write-lock per item**: When pushing a change to ADO, mark the local item as "sync in progress" and ignore incoming changes for that item until the push completes.
- **Store `pushed_at` timestamps**: The schema already includes this (good). Use it to detect and suppress echo updates.
- **Clear ownership rules:** ADO owns organizational fields (state, assigned_to, area_path). Local owns personal fields (personal priority, notes, subtasks). xl's approach works: status mapping is one-directional per action (push = local→ADO, pull = ADO→local).

**Warning signs:**
- `updated_at` timestamps on local items changing without user action
- ADO work item revision counts climbing rapidly (approaching the 10,000 REST API revision limit)
- Network traffic to ADO API during idle periods
- Users seeing "phantom" state changes that revert

**Phase to address:**
Phase 2 (Bidirectional Sync) — This is the single hardest technical problem in the project. Dedicate serious design time before coding.

---

### Pitfall 4: Scope Creep — "Single Pane of Glass" Becomes "Reinventing ADO"

**What goes wrong:**
The tool starts as a lightweight productivity layer but gradually absorbs ADO features: sprint planning, backlog grooming, capacity planning, custom fields, work item templates, reporting. Each addition seems small, but collectively the tool becomes a second ADO that's always behind the real one and never quite right.

**Why it happens:**
Users request features they already have in ADO ("can we also see velocity charts?"). Developers oblige because the requests seem reasonable. The tool's value was supposed to be personal task management layered on ADO, not a complete replacement.

**How to avoid:**
- **Define the "link out" boundary**: For any feature that ADO already does, deep-link to ADO instead of reimplementing. E.g., clicking a work item opens it in ADO, not in a local editor.
- **Apply the 80/20 rule ruthlessly**: The tool shows ADO data (read), lets you set personal priority/notes (write local), and promotes/pulls tasks (sync). That's it. Everything else is a link to ADO.
- **Maintain an explicit "Anti-Features" list** in PROJECT.md (already started in Out of Scope — keep enforcing)
- **Time-box feature discussions**: "Would this take less than 2 hours AND does ADO not already do it?" If no to either, defer.
- **ADO is source of truth** (already decided): Enforce this by never allowing the tool to modify ADO fields that ADO's own UI manages (sprint assignment, area path, work item type)

**Warning signs:**
- Users asking "why can't I edit the description here?" (because they should be in ADO)
- Building custom form fields that mirror ADO fields
- PR comments adding features that exist in ADO but "would be nice to have here too"
- Feature backlog growing faster than completion rate
- More than 50% of features duplicate ADO functionality

**Phase to address:**
Every phase — but especially Phase 3+ when the tool is functional and feature requests accelerate. Establish the boundary in Phase 1 architecture decisions.

---

### Pitfall 5: Auth Token Lifecycle Mismanagement

**What goes wrong:**
Application uses PATs (Personal Access Tokens) for development convenience and never migrates to proper OAuth/Entra ID. PATs expire, get committed to repos, get shared between team members, and have no revocation audit trail. Alternatively: the app implements OAuth but doesn't handle token refresh, leading to random "401 Unauthorized" errors mid-session.

**Why it happens:**
PATs are easy — one header, no redirect flows, works immediately. Entra ID OAuth requires app registration, redirect URIs, consent flows, and token refresh logic. For an internal tool, PATs feel "good enough." Microsoft is actively pushing away from PATs: *"Use personal access tokens sparingly, and only when Microsoft Entra ID isn't available"* and plans further PAT reduction.

**How to avoid:**
- **Phase 1: Use PATs for dogfooding** (Luis-only). This is acceptable for a single developer building the tool.
- **Phase 2: Implement Entra ID OAuth before team rollout.** Use MSAL (Microsoft Authentication Library) — it handles token caching, refresh, and silent acquisition automatically.
- **Never store PATs in code or config files.** Use Azure CLI auth (`az account get-access-token --resource 499b84ac-1321-427f-aa17-267ca6975798`) as xl already does.
- **Treat tokens as opaque** — Microsoft is encrypting token payloads starting 2025. Never decode or inspect JWT claims from ADO tokens.
- **Request minimum scopes**: `vso.work` for work items (read), `vso.work_write` for updates. Add `vso.code` only if showing PRs.
- **Handle 401 gracefully**: Redirect to re-auth, don't just show an error.

**Warning signs:**
- PATs hard-coded in `.env` files
- Token expiry causing "it worked yesterday" bugs
- Team members sharing a single PAT
- 401 errors in production logs without automatic recovery
- "Works for Luis, broken for others" (different tenant configs)

**Phase to address:**
Phase 1 (PAT for dogfooding), Phase 3 (Entra ID OAuth before team rollout). The auth refactor must happen before anyone besides Luis uses the tool.

---

### Pitfall 6: Adoption Death — Tool Isn't Useful in Week One

**What goes wrong:**
The tool launches with a grand vision but requires users to change their workflow significantly before seeing value. Engineers try it once, find it's slower than their current approach (opening ADO directly), and never return. The "set up your sync, configure your filters, link your tasks" onboarding takes 30+ minutes. Dead on arrival.

**Why it happens:**
Internal tools compete with inertia, not with nothing. Every team member already has a working (if fragmented) system. The bar for switching is: "this must be immediately, obviously better than what I do now." Most internal tools fail this bar because they optimize for the completed vision rather than the day-one experience.

**How to avoid:**
- **Zero-config first view**: On first login, immediately show the user's assigned ADO work items and open PRs. No setup required — use their identity to query ADO.
- **Value before contribution**: The tool should provide value (seeing all your work) before asking anything (creating tasks, configuring sync).
- **10-second test**: A new user should see useful data within 10 seconds of their first login.
- **Don't require local task creation**: The dashboard must be valuable as a read-only ADO viewer with personal priority overlays. Task creation is an advanced feature.
- **Dogfood ruthlessly**: Luis must use this as his primary work surface for 2+ weeks before team rollout. If Luis switches back to `xl`, something is wrong.
- **Friction audit**: Every action should take fewer steps than doing it in ADO directly. If it's harder → don't ship it.
- **Don't force workflow**: Let people use just the PR view, or just the task view. Don't require adoption of the whole system.

**Warning signs:**
- Users logging in once and never returning (track with simple analytics)
- "It's cool but I'll just use ADO" feedback
- Setup guide longer than one page
- Features that only work after extensive configuration
- Only 1-2 people using it after 2 weeks of team access

**Phase to address:**
Phase 1 (read-only dashboard), Phase 3 (team rollout). The read-only view must be compelling standalone before adding write features.

---

### Pitfall 7: ADO Work Item Revision Limit (10,000 via REST API)

**What goes wrong:**
Automated sync creates frequent, small updates to work items. Each update via REST API increments the revision counter. At 10,000 revisions, the REST API stops accepting updates. The work item is effectively frozen for programmatic updates (web portal updates still work).

**Why it happens:**
Microsoft explicitly enforces this: *"A work item revision limit of 10,000 is in effect for updates made through the REST API."* This is a hard limit. With aggressive sync (updating status every few minutes, syncing notes, priority changes), a busy work item can burn through revisions in months.

**How to avoid:**
- **Batch field updates**: Never update one field at a time. Accumulate changes and send a single JSON Patch with all modifications.
- **Minimize sync frequency**: Sync on explicit user action (push/pull), not on a timer. If you must poll, sync state only, not metadata.
- **Track revision consumption**: Monitor `rev` numbers. Alert when items approach 5,000 revisions.
- **Avoid syncing personal-only fields to ADO**: Personal priority, local notes, and local subtasks should NEVER be written to ADO fields. They live only in the local database.
- **Use tags sparingly**: Each tag add/remove is a revision. Don't use ADO tags as sync markers.

**Warning signs:**
- `rev` numbers on work items climbing faster than expected
- REST API 400/403 errors on work item updates for specific items
- Some items updatable via web portal but failing via API

**Phase to address:**
Phase 2 (Bidirectional Sync) — Design the sync protocol to minimize writes from day one.

---

### Pitfall 8: WIQL Query Performance Cliff

**What goes wrong:**
Queries that work fine during development with 50 work items become 30-second timeouts in production with 5,000+ items. The dashboard shows a loading spinner or errors. Users lose trust.

**Why it happens:**
WIQL query execution has a **hard 30-second timeout** (Azure DevOps Services). Results are capped at **20,000 items**. WIQL has no OFFSET/LIMIT pagination — it only has `$top`. Certain operators (`Ever`, `Contains` on long text, `<>`, `Not`, `In Group` with large groups) are expensive. The query planner doesn't expose explain plans, so you can't diagnose slow queries without trial and error.

**How to avoid:**
- **Always include date/range filters** at the top of queries: `[System.ChangedDate] >= @Today - 30`
- **Scope to project**: Always include `[System.TeamProject] = @Project`. Unscoped queries search the entire collection.
- **Avoid `Contains` on long text fields** — use `Contains Words` instead
- **Minimize `Ever` operator** — it's a full-history scan
- **Don't use WIQL for reporting** — use the reporting APIs (`/_apis/wit/reporting/workitemrevisions`) instead
- **Pre-filter on assignment**: `[System.AssignedTo] = @Me` dramatically reduces result sets for personal dashboards
- **Cache query results** with TTL: Most ADO data doesn't change minute-to-minute
- **WIQL character limit: 32,000 characters** — dynamic queries with many IDs can exceed this

**Warning signs:**
- `VS402335: The timeout period (30 seconds) elapsed` errors
- Dashboard taking 5+ seconds to load
- Queries returning max 20,000 results (you've hit the cap and are missing items)

**Phase to address:**
Phase 1 (Core ADO Integration) — Get queries right early. Bad query patterns baked into the codebase are painful to refactor.

---

### Pitfall 9: Service Hooks (Webhooks) Are Unreliable for Critical Data Sync

**What goes wrong:**
Team builds real-time sync around ADO service hooks (webhooks), assuming they're reliable event delivery. They're not. Service hooks have complex failure states: transient failures retry up to 8 times over ~3 minutes, enduring failures put the subscription on "probation" (events are **lost** during probation), and terminal failures (410) permanently disable the subscription. A brief outage of your webhook endpoint can cause you to miss hours of ADO events.

**Why it happens:**
Webhook-based architectures assume at-least-once delivery. ADO service hooks provide at-most-once during failure states. During probation, events are dropped silently — there's no replay mechanism and no dead-letter queue.

**How to avoid:**
- **Webhooks for acceleration, not correctness**: Use service hooks to trigger immediate UI refreshes, but never as the sole data source. Always have a polling fallback.
- **Implement periodic reconciliation**: Every 5-15 minutes, poll ADO for changes since last successful sync, regardless of webhooks. This catches anything missed during probation.
- **Monitor subscription health**: Check the `/_apis/hooks/subscriptions` endpoint to detect `Disabled` or `Restricted` states. Re-create subscriptions that die.
- **Idempotent event processing**: Events may arrive out of order or be duplicated (during retry). Use `rev` numbers to determine if an event is stale.
- **Payload size limit: 2 MB** — large work items with many attachments/comments can exceed this.

**Warning signs:**
- Dashboard data is stale despite webhook subscription being "active"
- Service hook subscriptions showing "Enabled (restricted)" state in ADO
- Gaps in event history visible in ADO's service hooks admin page
- Missing events after network blips or deployment restarts

**Phase to address:**
Phase 2+ (if webhooks are used). For Phase 1, polling-only is simpler and more reliable. Add webhooks as an optimization in Phase 3+.

---

### Pitfall 10: Building for the Team Before the Individual Works

**What goes wrong:**
The tool is designed for multi-user from day one: role-based views, manager dashboards, team analytics, shared task boards. This adds massive complexity (permissions, data isolation, concurrent editing) while the core value proposition — "a better personal work view" — remains untested.

**Why it happens:**
The vision is a team tool, so it feels right to build for the team. But the PROJECT.md wisely states: *"Start with Luis (dogfooding), grow to team — must be useful for one person before it scales."* Ignoring this leads to premature abstraction.

**How to avoid:**
- **Phase 1-2: Single user, local-first.** SQLite on the developer's machine. No server. No shared state.
- **Phase 3: Add multi-user only after dogfooding proves value.** The server/shared database comes when team rollout begins.
- **Don't build permissions until there are multiple users.** "Everyone can see everything" is fine for a 6-15 person team in Phase 3.
- **Manager views are Phase 4+ at earliest.** ICs must love the tool before managers even know it exists.

**Warning signs:**
- Database schema discussions involving `user_id` foreign keys in Phase 1
- Auth middleware before the core dashboard works
- "But what about when the team uses it?" blocking simple local features

**Phase to address:**
Phase 1 (architecture decision) — Choose local-first architecture. Multi-user is a Phase 3 concern.

---

### Pitfall 11: Data Model Drift from xl

**What goes wrong:**
xl's schema is proven but was designed for single-user. Naively porting it to multi-user can create issues (e.g., `tasks` table has no `user_id`, `memory` table is personal, `sessions` are per-CLI-session).

**Why it happens:**
xl manages 17 entities in SQLite for a single user. The vocabulary and field names are battle-tested (29 tasks, 10 PRs, 53 daily logs, 158 memory entries), but the architecture assumes one user per database.

**How to avoid:**
- Audit each xl table for multi-user readiness before porting
- Add `user_id` foreign key to all personal data tables (Phase 3, not Phase 1)
- `memory` table becomes per-user preferences
- `sessions` becomes web session tracking (different from CLI session logs)
- Keep xl's field names and statuses — the vocabulary is proven
- Reuse xl's ADO state mapping (`todo→Proposed`, `in_progress→Active`, `done→Completed`, `blocked→Blocked`)

**Warning signs:**
- Users seeing each other's personal tasks
- Memory/notes leaking across users
- Session model doesn't map to web sessions

**Phase to address:**
Phase 2 (schema design) for single-user, Phase 3 (multi-user migration) for team rollout.

---

## Technical Debt Patterns

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| PATs instead of OAuth | Instant auth, no setup | Security risk, no audit trail, Microsoft deprecating | Phase 1 dogfooding only |
| Polling instead of webhooks | Simpler, more reliable | Higher API cost, slower updates | Acceptable through Phase 2; add webhooks Phase 3+ |
| Single WIQL query for all data | Simple codebase | Performance cliff at scale, hits rate limits | Never — separate queries for work items, PRs, comments from day one |
| Storing ADO field copies locally | Offline access, faster UI | Data staleness, sync complexity | Acceptable IF you store `rev` numbers and treat local copies as cache, not truth |
| No caching layer | Simpler architecture | Every page load hits ADO API, rate limits | Never — even a 60-second in-memory TTL cache is essential |
| Monolith instead of service layers | Faster to build | Hard to add VS Code extension later | Phase 1-2 only; extract API layer before adding second surface |

## Integration Gotchas

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| ADO WIQL API | Assuming SELECT returns field data | WIQL returns IDs only → batch fetch with Work Items API |
| ADO Work Item Updates | Updating one field at a time | Use JSON Patch with all field changes in a single operation |
| ADO Work Item Updates | Not checking `rev` before writing | Always read current `rev`, include in If-Match header to detect conflicts |
| ADO PR API | Fetching all PRs across repos | Filter by `reviewerId` or `creatorId` to get only relevant PRs |
| ADO Service Hooks | Relying on webhooks as sole event source | Webhooks for speed + polling for correctness (belt and suspenders) |
| ADO Identity | Hardcoding identity format | Identities can be email, display name, or GUID — normalize on first use |
| Azure CLI Auth | Not handling token expiry | `az account get-access-token` can fail silently if login expired; catch and prompt re-auth |
| MSAL Token Cache | Building custom token storage | Use MSAL's built-in token cache (persistence plugins for Node.js: `@azure/msal-node-extensions`) |

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Unbounded WIQL queries | 30-second timeouts, dashboard hangs | Always include date range + project scope + assignment filter | >500 work items in scope |
| N+1 work item fetches | Page load takes 5+ seconds | Batch IDs into groups of 200, fetch in parallel | >20 work items displayed |
| Fetching full work item payloads | Large API responses, high bandwidth | Use `fields` parameter to request only displayed fields | >100 work items per sync |
| Real-time polling | Rate limit hits, 429 responses | Poll every 5 minutes max; use in-memory cache for UI | >3 concurrent users polling |
| Fetching PR threads eagerly | Slow initial load, wasted API calls | Lazy-load comment threads on expand/click, not on page load | >10 active PRs per user |
| No response compression | Slow on corporate VPN | Request gzip encoding (`Accept-Encoding: gzip`) on all API calls | Always, especially over VPN |

## Security Mistakes

| Mistake | Risk | Prevention |
|---------|------|------------|
| PATs in source code or `.env` files | Token compromise, unauthorized ADO access | Use environment variables, Azure CLI auth, or OS keychain. Never commit tokens. |
| Overly broad OAuth scopes | App has more access than needed | Request only `vso.work` (read) and `vso.work_write` (write). Add `vso.code` only if showing PRs. |
| Decoding/inspecting ADO JWT tokens | Breaking change: Microsoft encrypting tokens in 2025 | Treat all tokens as opaque strings. Use REST APIs for user info. |
| Shared service account for all users | No per-user audit trail, permission escalation | Each user authenticates individually via Entra ID OAuth |
| Storing ADO data without access control | User A sees User B's private work items | Respect ADO permissions — only display items the authenticated user can see |
| Long-lived token without refresh | Stale credentials, security window | Use MSAL's automatic token refresh; detect 401 and trigger silent re-auth |

## UX Pitfalls

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Requiring setup before showing value | Users bounce immediately | Show assigned work items on first load with zero config |
| Duplicating ADO's edit UI | Confusing — which is the "real" version? | Deep-link to ADO for editing. Tool is for viewing + personal annotations. |
| Showing stale data without indicating it | Users make decisions on outdated info | Show "last synced: 3 minutes ago" timestamp on every data section |
| Overwhelming "single pane of glass" | Information overload, users can't find what matters | Default to "my stuff today" — expand to team view on demand |
| No empty states | New users see blank screens, think it's broken | "No work items assigned. Here's how to get started..." with a link to ADO |
| Slow perceived performance | Users assume the tool is broken | Optimistic UI: show cached data immediately, refresh in background, animate transitions |
| Requiring bidirectional sync to be useful | Creates adoption barrier | Read-only mode must be fully functional. Sync is opt-in advanced feature. |

## "Looks Done But Isn't" Checklist

- [ ] **ADO Sync:** Often missing conflict detection — verify both sides of a sync can't clobber each other by testing simultaneous edits
- [ ] **WIQL Queries:** Often missing project scoping — verify queries include `[System.TeamProject] = @Project` or they search the entire org
- [ ] **Auth Flow:** Often missing token refresh — verify the app handles expired tokens gracefully (not just initial auth)
- [ ] **Error Handling:** Often missing rate limit handling — verify the app respects `Retry-After` and `X-RateLimit-Remaining` headers
- [ ] **Pagination:** Often missing — verify work item lists handle more than 200 results (the batch API limit per call)
- [ ] **PR Display:** Often missing thread state — verify resolved/active/won't fix comment states are displayed, not just text
- [ ] **Identity Resolution:** Often hardcoded to email — verify display names, GUIDs, and email formats all resolve correctly
- [ ] **Offline Behavior:** Often missing — verify the app degrades gracefully when ADO is unreachable (show cached data, not an error page)
- [ ] **Empty States:** Often missing — verify every list has an empty state message, not a blank white space
- [ ] **Link Behavior:** Often opening in-app — verify ADO links open ADO in the browser, not a broken in-app render

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Infinite sync loop | MEDIUM | Add sync provenance tracking, add debounce, deploy immediately. Check for rev inflation on affected work items. |
| Rate limit hit | LOW | Respect Retry-After header, add exponential backoff, reduce polling frequency. Recovers automatically within 5 minutes. |
| PAT leaked in repo | HIGH | Immediately revoke PAT in ADO, audit access logs, rotate all credentials, migrate to Entra ID OAuth. |
| WIQL timeout in production | MEDIUM | Add date range filters, scope to project, cache results. May need to split into multiple focused queries. |
| Service hook subscription disabled | LOW | Re-create subscription via API. Trigger full reconciliation poll to catch missed events. |
| Work item at 10K revision limit | HIGH | Cannot recover via API. Must use web portal for future edits. Avoid by batching updates and minimizing write frequency. |
| Scope creep beyond ADO layer | HIGH | Audit feature list against "Anti-Features." Remove or convert to deep-links. Requires product discipline, not just code changes. |
| User adoption failure | HIGH | Requires re-evaluation of core value prop. Simplify to zero-config read-only view. Interview lapsed users. May need feature pivot. |

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| WIQL returns IDs only | Phase 1 | Unit test confirming two-step fetch pattern in ADO client |
| TSTU rate limits | Phase 1 | Integration test with rate limit header monitoring; backoff logic in HTTP client |
| Infinite sync loop | Phase 2 | Test: update local → push → verify no re-pull of same change; check rev numbers |
| Scope creep | Every phase | Architecture review: any feature that reimplements ADO UI is flagged for removal |
| Auth token lifecycle | Phase 1 (PAT), Phase 3 (OAuth) | Test: token expires mid-session → app recovers without user seeing raw error |
| Adoption death | Phase 1 (read-only), Phase 3 (team) | Metric: time-to-first-useful-data < 10 seconds for new user |
| Revision limit | Phase 2 | Monitor: alert when any work item exceeds 5,000 revisions |
| WIQL performance | Phase 1 | Load test: query with 1,000+ items completes in < 5 seconds |
| Webhook unreliability | Phase 3+ | Test: kill webhook endpoint for 5 minutes → verify reconciliation poll catches up |
| Premature multi-user | Phase 1 (architecture) | Review: Phase 1-2 code has zero multi-tenancy abstractions |
| Data model drift from xl | Phase 2 (schema), Phase 3 (multi-user) | Audit: each xl entity reviewed for multi-user readiness before porting |

## Sources

- [Azure DevOps Rate and Usage Limits](https://learn.microsoft.com/en-us/azure/devops/integrate/concepts/rate-limits) — TSTU system, 200 TSTU/5min limit, response headers (doc date: 2025-09-15)
- [WIQL Syntax Reference](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax) — Returns IDs only, 32K char limit, ASOF clause (doc date: 2026-02-28)
- [Integration Best Practices](https://learn.microsoft.com/en-us/azure/devops/integrate/concepts/integration-bestpractices) — Batch changes, 10K revision limit, query optimization, reporting APIs (doc date: 2026-03-02)
- [Work Tracking Object Limits](https://learn.microsoft.com/en-us/azure/devops/organizations/settings/work/object-limits) — 30-second query timeout, 20K result cap, 10K revision limit (doc date: 2025-11-06)
- [Service Hooks Overview](https://learn.microsoft.com/en-us/azure/devops/service-hooks/overview) — Event model, available services, webhook setup (doc date: 2025-06-04)
- [Service Hooks Troubleshooting](https://learn.microsoft.com/en-us/azure/devops/service-hooks/troubleshoot) — Failure types, probation system, retry behavior, 2MB payload limit (doc date: 2025-12-03)
- [Authentication Methods](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/authentication-guidance) — Entra ID recommended, PAT reduction, token opacity (doc date: 2026-03-18)
- [Entra ID OAuth for ADO](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/entra-oauth) — App registration, resource ID, MSA limitations (doc date: 2025-01-08)
- PROJECT.md — Team context, xl's proven ADO patterns, dogfooding strategy
- REFERENCE.md — xl's existing ADO client, schema design, sync patterns

---
*Pitfalls research for: Team ADO Tool — Developer Productivity Dashboard*
*Researched: 2025-07-14*
