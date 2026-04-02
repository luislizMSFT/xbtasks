# External Integrations

**Analysis Date:** 2025-07-15

## APIs & External Services

**Microsoft Graph API:**
- Purpose: Fetch authenticated user profile (display name, email) and avatar photo after OAuth sign-in
- Endpoints used:
  - `GET https://graph.microsoft.com/v1.0/me` — User profile
  - `GET https://graph.microsoft.com/v1.0/me/photo/$value` — User avatar (URL stored, not fetched at auth time)
- Client: Standard `net/http` in `internal/auth/auth.go` (lines 284–318)
- Auth: Bearer token from OAuth2 flow
- Scopes requested: `User.Read`, `offline_access`

**Microsoft Entra ID (Azure AD) — OAuth2:**
- Purpose: User authentication via OAuth2 Authorization Code + PKCE flow
- Endpoints:
  - Auth: `https://login.microsoftonline.com/{tenantId}/oauth2/v2.0/authorize`
  - Token: `https://login.microsoftonline.com/{tenantId}/oauth2/v2.0/token`
- Implementation: `internal/auth/auth.go` — `SignIn()` method (lines 66–143)
- Flow: PKCE (S256 code challenge) → browser redirect → localhost callback → token exchange
- Callback: Ephemeral `localhost:{random_port}/callback` HTTP server with 5-minute timeout
- Token storage: OS keychain via `go-keyring` (access_token, refresh_token, token_expiry)
- Session restore: `TryRestoreSession()` uses stored refresh token to silently re-authenticate
- PAT fallback: `SignInWithPAT()` for development without OAuth (stores PAT in keychain)

**Azure DevOps (ADO) — Planned Integration:**
- Purpose: Sync work items, pull requests, and project data from Azure DevOps
- Current state: Domain types defined (`pkg/ado/domain/types.go`), DB schema includes `ado_work_items`, `pull_requests`, and `task_ado_links` tables, but no ADO REST client is implemented yet
- Config keys ready: `ado.organization`, `ado.project`, `ado.pat_keychain_key` (in `internal/config/config.go`)
- Domain types defined:
  - `ADOWorkItem` — ID, title, state, type (bug/task/user-story/deliverable/scenario), assigned_to, priority, area_path, URL
  - `ADOPullRequest` — ID, title, URL, repo, status (draft/active/completed/abandoned), reviewers, branches, votes, merged_at
  - `TaskADOLink` — Links local tasks to ADO items with direction (promoted/imported/linked)
- DB tables ready: `ado_work_items`, `pull_requests`, `task_ado_links` (in `internal/db/db.go`)
- Planned sync: `sync.interval_minutes` config key (default 15 min) exists but sync loop not implemented

## Data Storage

**Database:**
- Type: SQLite 3 (via `github.com/mattn/go-sqlite3`, CGO required)
- Connection string: `{dbPath}?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON`
- Default location: `<DataDir>/data.db`
  - macOS: `~/Library/Application Support/team-ado-tool/data.db`
  - Windows: `%LOCALAPPDATA%\team-ado-tool\data.db`
  - Linux: `~/.local/share/team-ado-tool/data.db`
- Client: Thin `db.DB` wrapper around `database/sql.DB` in `internal/db/db.go`
- Migration: Single inline `schema` constant executed on every `Open()` — all `CREATE TABLE IF NOT EXISTS` (idempotent)
- Tables:
  - `projects` — Local project grouping
  - `tasks` — Core task management (status, priority P0–P3, category, tags, parent/subtask hierarchy)
  - `task_deps` — Task dependency graph (with cycle detection in Go)
  - `pull_requests` — ADO pull request cache
  - `ado_work_items` — ADO work item cache
  - `task_ado_links` — Bidirectional links between tasks and ADO items
  - `users` — Authenticated user profiles
- Indexes: On `tasks(status)`, `tasks(priority)`, `tasks(project_id)`, `tasks(ado_id)`, `tasks(parent_id)`, `pull_requests(status)`, `pull_requests(task_id)`, `ado_work_items(ado_id)`

**File Storage:**
- No file/blob storage. All data is in SQLite.
- Frontend assets are embedded in the Go binary via `//go:embed all:frontend/dist`

**Caching:**
- No external cache (Redis, etc.)
- SQLite serves as the local cache for ADO data (planned)
- Frontend Pinia stores hold in-memory state

## Authentication & Identity

**Auth Provider:** Microsoft Entra ID (Azure AD)
- Implementation: `internal/auth/auth.go` — `AuthService` struct
- Primary flow: OAuth2 Authorization Code + PKCE
  1. Generate code verifier + S256 challenge
  2. Start ephemeral HTTP server on random localhost port
  3. Open system browser to Entra ID authorize endpoint
  4. Receive auth code via `/callback`
  5. Exchange code + verifier for tokens
  6. Fetch user profile from Microsoft Graph
  7. Store tokens in OS keychain, user in SQLite
- Fallback: Personal Access Token (PAT) sign-in for development (`SignInWithPAT()`)
- Session persistence: Tokens stored in OS keychain (`go-keyring`)
  - Keys: `team-ado-tool/access_token`, `team-ado-tool/refresh_token`, `team-ado-tool/token_expiry`, `team-ado-tool/pat`
- Session restore: On app startup, `TryRestoreSession()` runs in a goroutine to silently refresh tokens
- Frontend events: `auth:state-changed` emitted via Wails event system on sign-in/sign-out
- Frontend store: `frontend/src/stores/auth.ts` — `useAuthStore` with mock fallback for dev

**Required Environment Variables for Auth:**
- `TEAM_ADO_TENANT_ID` — Entra tenant ID (placeholder `YOUR_TENANT_ID` if unset)
- `TEAM_ADO_CLIENT_ID` — App registration client ID (placeholder `YOUR_CLIENT_ID` if unset)

## Monitoring & Observability

**Error Tracking:**
- None. No Sentry, Datadog, or similar service integrated.

**Logs:**
- Go: Standard `log` package (`log.Printf`, `log.Fatalf`); config key `log.level` exists but not wired to a leveled logger
- Frontend: `console` (no structured logging)

## CI/CD & Deployment

**Hosting:**
- Desktop app — distributed as native binaries per platform
- macOS: `.app` bundle (build config in `build/darwin/`)
- Windows: `.exe` + NSIS installer (build config in `build/windows/`)
- Linux: AppImage (build config in `build/linux/`)
- Server mode: Docker-based HTTP server (Dockerfile at `build/docker/Dockerfile.server`)

**CI Pipeline:**
- Not configured. No `.github/workflows/`, `azure-pipelines.yml`, or similar CI config detected.

**Build Orchestration:**
- Task (go-task.dev) via `Taskfile.yml` with platform-specific includes
- Key tasks: `dev`, `build`, `package`, `build:server`, `build:docker`

## Environment Configuration

**Required env vars:**
- `TEAM_ADO_TENANT_ID` — Microsoft Entra tenant ID for OAuth (required for real auth)
- `TEAM_ADO_CLIENT_ID` — OAuth app registration client ID (required for real auth)

**Optional env vars:**
- Any config key can be overridden via `XBT_` prefix (e.g., `XBT_DB_PATH`, `XBT_LOG_LEVEL`)

**Secrets storage:**
- OS keychain via `go-keyring`:
  - `team-ado-tool/access_token` — OAuth access token
  - `team-ado-tool/refresh_token` — OAuth refresh token
  - `team-ado-tool/token_expiry` — Token expiry timestamp
  - `team-ado-tool/pat` — Personal access token (dev fallback)
- No `.env` files detected in the repository

## Wails Bindings (Go → Frontend Bridge)

**Mechanism:**
- Go services registered via `wailsApp.RegisterService()` in `main.go`
- Wails CLI generates TypeScript bindings into `frontend/bindings/` (gitignored, regenerated on build)
- Frontend imports bindings dynamically: `import { SignIn } from '../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice'`
- Binding paths mirror Go module/package structure

**Registered Services:**
| Service | Go Location | Binding Path |
|---------|------------|--------------|
| `ConfigService` | `internal/config/service.go` | `bindings/.../config/configservice` |
| `TaskService` | `internal/app/tasks.go` | `bindings/.../app/taskservice` |
| `ProjectService` | `internal/app/projects.go` | `bindings/.../app/projectservice` |
| `DependencyService` | `internal/app/deps.go` | `bindings/.../app/dependencyservice` |
| `AuthService` | `internal/auth/auth.go` | `bindings/.../auth/authservice` |

**Events (Go → Frontend):**
- `auth:state-changed` — Emitted on sign-in/sign-out with `{authenticated: bool}` payload

## Webhooks & Callbacks

**Incoming:**
- Ephemeral OAuth callback: `http://localhost:{random}/callback` — only active during sign-in flow, self-closes after receiving auth code

**Outgoing:**
- None

---

*Integration audit: 2025-07-15*
