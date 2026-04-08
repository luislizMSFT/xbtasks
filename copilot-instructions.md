<!-- GSD:project-start source:PROJECT.md -->
## Project

**Team ADO Tool**

A unified productivity dashboard for the Xbox Services team — a native desktop app (Wails v3 + Vue + Go + SQLite) that brings all work — ADO work items, PRs, comments, and personal tasks — into a single view. It replaces the fragmented landscape where team members use different tools (ADO, Wiki, Loop, CLIs, M365 Copilot) with one app everyone gravitates to. Local-first with SQLite, system tray presence, native notifications.

**Core Value:** One pane of glass for all your work — personal tasks linked to ADO items, PRs, and comments — so nobody has to context-switch between tools to know what they need to do.

### Constraints

- **Auth**: Must use Microsoft identity / Entra ID — team is on Microsoft ecosystem. PAT fallback for simplicity.
- **ADO API**: Azure DevOps REST API is the primary integration surface; must handle rate limits and pagination
- **Audience**: Start with Luis (dogfooding), grow to team — must be useful for one person before it scales
- **Stack**: Wails v3 (Go) + Vue 3 (thin shell) + SQLite — native desktop app, not a web app
- **No XAML, no React, no C#**: Team prefers Go backend with lightweight web frontend in native shell
- **Vue is thin shell**: All logic lives in Go. Vue is display + interaction only.
- **Local-first**: SQLite per user, no server dependency for core functionality
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Recommended Stack (FINAL)
### Native Shell: Wails v3
- **Wails v3** — Go-native desktop framework with native webview (not Chromium)
- System tray, notifications, multiple windows — all native Go APIs
- Cross-platform (Windows primary, Mac/Linux possible)
- Small binary (~5-10MB vs 100MB+ Electron)
- Go functions directly callable from frontend (auto-generated bindings)
- v3 is alpha but usable in production; v2 is stable fallback
### Frontend: Vue 3 + TypeScript (thin shell)
- **Vue 3** with Composition API — thin display layer, not a heavy SPA
- **TypeScript** — type safety for Wails bindings
- **Vite** — fast dev server and build (Wails default)
- **Tailwind CSS** — utility-first rapid UI
- **PrimeVue or Naive UI** — component library for data tables, tree views
### Backend: Go (Wails bound functions)
- **Go** — team's primary language, xl codebase to port from
- Wails bindings expose Go functions to Vue with auto-generated TypeScript types
- **go-sqlite3 or modernc.org/sqlite** — SQLite access (same as xl)
- **net/http** — ADO REST API client (port xl's pkg/ado/)
- **golang.org/x/oauth2** — Entra ID OAuth2 PKCE flow
- **keyring** — OS-level credential storage
### Database: SQLite + FTS5
- **SQLite** — local-first, proven in xl
- Each user has their own local database
- **FTS5** for full-text search (same as xl)
- Schema derived from xl's 17-entity model
- No server, no deployment, no connection strings
### Authentication: Microsoft Entra ID
- **OAuth 2.0 Authorization Code Flow with PKCE** for desktop app
- **Scope:** `https://app.vssps.visualstudio.com/.default` for ADO API
- Token stored via OS keychain (keyring library)
- Auto-refresh on expiry
- **PAT fallback** for users who prefer it (like xl currently does)
### ADO Integration: Azure DevOps REST API
- **Polling** for data sync (desktop app can't receive webhooks)
- **WIQL queries** for complex work item queries
- **Rate limit handling:** exponential backoff, `Retry-After` headers
- **API version pinning**
- **Background sync:** configurable interval with change detection
- Port xl's state mapping: todo→Proposed, in_progress→Active, done→Completed, blocked→Blocked
### Distribution
- **Wails bundler** for Windows installer / MSIX
- Internal distribution via SharePoint, Azure Blob, or private GitHub Releases
- Auto-updater (Wails v3 or custom update check)
## What NOT to Use
| Technology | Why Not |
|-----------|---------|
| React / Next.js | User preference — Vue thin shell chosen |
| Tauri | Requires Rust — team knows Go |
| Electron | Bloated, high memory, bundles Chromium |
| Blazor / C# / XAML | Maintenance burden, XAML complexity |
| PostgreSQL | Desktop app is local-first — SQLite per user |
| htmx / templ | Not suitable for rich dashboard |
## Stack Summary
| Layer | Choice | Confidence |
|-------|--------|------------|
| Native Shell | Wails v3 (Go) | High |
| Frontend | Vue 3 + TypeScript + Vite + Tailwind (thin shell) | High |
| Backend | Go (Wails bound functions) | High |
| Database | SQLite + FTS5 (modernc.org/sqlite) | High |
| Auth | Entra ID (OAuth2 PKCE) + PAT fallback | High |
| ADO Integration | REST API + polling | High |
| Distribution | Wails bundler + internal hosting | Medium |
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

### Async & Concurrency Patterns (Go)
- **Wails bindings are synchronous** — the frontend blocks until Go returns. All network-bound methods must have timeouts.
- **az CLI calls**: Always use `RunAzCliCtx` with `exec.CommandContext` and `AzCliTimeout` (15s). Never shell out without a timeout.
- **Token caching**: `CachedTokenProvider` uses RWMutex for reads + separate `refreshMu` for serializing refresh. Double-check pattern: re-verify cache after acquiring refreshMu. Never hold locks during I/O.
- **Sync mutual exclusion**: `SyncService.syncMu` serializes background and manual sync. Use `TryLock()` — background sync skips if locked; manual sync returns "already in progress" error.
- **Lifecycle events**: Backend emits Wails events (`sync:started`, `sync:completed`, `sync:failed`) for frontend reactivity. Use `s.emitEvent()` helper.
- **Multi-org fan-out**: `ListMyWorkItems` uses goroutine-per-org with `time.After(20s)` select. Returns partial results on timeout; only errors if ALL orgs fail.
- **Context propagation**: Timeouts are applied at entry points (service methods). Full context.Context propagation deferred to Phase 9.

### Frontend Patterns (Vue)
- **Wails event listeners**: Use `Events.On('event-name', callback)` from `@wailsio/runtime`. Initialize once on auth via `initEvents()` in stores.
- **Debouncing**: Use `useDebounceFn` from `@vueuse/core` for search/filter inputs (200ms default).
- **Flight guards**: Stores use `isLoading` / `isSyncing` flags to prevent duplicate concurrent calls.
- **Binding imports**: Import from `@/api/` wrapper modules, never from deep `bindings/` paths.
- **Toast notifications**: Stores use `import { toast } from 'vue-sonner'` directly; components use `useNotify` composable.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

### Request Flow
```
Vue Component → await api/xxx.ts → Wails IPC → Go Service Method → (ADO REST / SQLite / az CLI) → return → UI updates
```
Wails bindings are **synchronous from the frontend's perspective** — the UI thread blocks until Go returns. All Go service methods with network I/O must include timeouts.

### Backend Goroutines
1. **Background sync** (`SyncService.StartBackgroundSync`) — 15-min ticker, pulls ADO changes
2. **OAuth callback** (`AuthService`) — HTTP listener for Entra ID redirect
3. **Multi-org fan-out** (`ADOService.ListMyWorkItems`) — goroutine-per-org, 20s timeout

### Concurrency Control
- `CachedTokenProvider`: RWMutex (cache reads) + refreshMu (serialize token refresh)
- `SyncService.syncMu`: Serializes background and manual sync operations
- az CLI calls: 15s `exec.CommandContext` timeout via `RunAzCliCtx`

### Event Communication (Backend → Frontend)
Backend emits Wails events for lifecycle state changes:
- `sync:started` / `sync:completed` / `sync:failed` — sync lifecycle
- Frontend subscribes via `Events.On()` in store `initEvents()` functions
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
