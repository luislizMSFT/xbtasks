# Stack Research: Team ADO Tool

## Recommended Stack (FINAL)

### Native Shell: Wails v3
**Confidence: High**

- **Wails v3** — Go-native desktop framework with native webview (not Chromium)
- System tray, notifications, multiple windows — all native Go APIs
- Cross-platform (Windows primary, Mac/Linux possible)
- Small binary (~5-10MB vs 100MB+ Electron)
- Go functions directly callable from frontend (auto-generated bindings)
- v3 is alpha but usable in production; v2 is stable fallback

**Why Wails over Tauri:** Team knows Go, doesn't know Rust. Tauri requires Rust for all backend logic (ADO API, SQLite, auth). Wails lets you write everything in Go — directly port xl's ADO patterns.

**Why Wails over Electron:** 10-20x smaller binary, native webview, Go backend is in-process.

### Frontend: Vue 3 + TypeScript (thin shell)
**Confidence: High**

- **Vue 3** with Composition API — thin display layer, not a heavy SPA
- **TypeScript** — type safety for Wails bindings
- **Vite** — fast dev server and build (Wails default)
- **Tailwind CSS** — utility-first rapid UI
- **PrimeVue or Naive UI** — component library for data tables, tree views

Vue is intentionally thin — just rendering and user interaction. All logic lives in Go.

### Backend: Go (Wails bound functions)
**Confidence: High**

- **Go** — team's primary language, xl codebase to port from
- Wails bindings expose Go functions to Vue with auto-generated TypeScript types
- **go-sqlite3 or modernc.org/sqlite** — SQLite access (same as xl)
- **net/http** — ADO REST API client (port xl's pkg/ado/)
- **golang.org/x/oauth2** — Entra ID OAuth2 PKCE flow
- **keyring** — OS-level credential storage

**Porting xl:** xl's entire `pkg/ado/` (client, auth, models, push, pull, sync, query) and `pkg/db/` (schema, queries, models) port directly. The data model is proven with 29 tasks, 10 PRs, 158 memory entries.

### Database: SQLite + FTS5
**Confidence: High**

- **SQLite** — local-first, proven in xl
- Each user has their own local database
- **FTS5** for full-text search (same as xl)
- Schema derived from xl's 17-entity model
- No server, no deployment, no connection strings

### Authentication: Microsoft Entra ID
**Confidence: High**

- **OAuth 2.0 Authorization Code Flow with PKCE** for desktop app
- **Scope:** `https://app.vssps.visualstudio.com/.default` for ADO API
- Token stored via OS keychain (keyring library)
- Auto-refresh on expiry
- **PAT fallback** for users who prefer it (like xl currently does)

### ADO Integration: Azure DevOps REST API
**Confidence: High**

- **Polling** for data sync (desktop app can't receive webhooks)
- **WIQL queries** for complex work item queries
- **Rate limit handling:** exponential backoff, `Retry-After` headers
- **API version pinning**
- **Background sync:** configurable interval with change detection
- Port xl's state mapping: todo→Proposed, in_progress→Active, done→Completed, blocked→Blocked

### Distribution
**Confidence: Medium**

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

---
*Researched: 2026-03-31, Updated: 2026-04-01 after stack finalization*
