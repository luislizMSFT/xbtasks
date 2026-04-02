# Technology Stack

**Analysis Date:** 2025-07-15

## Languages

**Primary:**
- Go 1.25+ (go.mod declares `go 1.25.0`; local runtime is Go 1.26.1) — Backend services, database, auth, config, ADO integration
- TypeScript (ESNext target) — Frontend SPA in Vue 3 SFCs and Pinia stores

**Secondary:**
- SQL (SQLite dialect) — Inline schema DDL and queries in Go source (`internal/db/db.go`)
- CSS (Tailwind v4 + custom properties) — Styling via `frontend/src/style.css`

## Runtime

**Environment:**
- **Desktop:** Wails v3 (alpha-74) — embeds a native webview; Go process serves embedded frontend assets via `embed.FS`
- **Go:** 1.25+ required by `go.mod`
- **Node.js:** 25.x (local is v25.6.1) — used only for frontend dev/build tooling
- **OS webview:** macOS WebKit, Windows WebView2, Linux WebKitGTK (provided by Wails runtime)

**Package Manager:**
- Go modules (`go.mod` / `go.sum`)
- npm 11.x (`frontend/package.json` / `frontend/package-lock.json`)
- Lockfiles: both present and committed

## Frameworks

**Core:**
- Wails v3 `v3.0.0-alpha.74` (`github.com/wailsapp/wails/v3`) — Desktop app framework; bridges Go services to frontend via auto-generated TypeScript bindings
- Vue 3 `^3.2.45` — Frontend SPA framework (Composition API, `<script setup>`)
- Pinia `^3.0.4` — State management (stores in `frontend/src/stores/`)
- Vue Router `^4.6.4` — Client-side routing with `createWebHashHistory()` (hash-based for webview compat)

**UI:**
- shadcn-vue (New York style) — Component library configured via `frontend/components.json`; components live in `frontend/src/components/ui/`
- Reka UI `^2.9.3` — Headless primitive components (underlying shadcn-vue)
- Headless UI `^1.7.23` — Additional headless Vue components
- Lucide Vue Next `^1.0.0` — Icon library
- Tailwind CSS v4 `^4.2.2` — Utility-first CSS via Vite plugin (`@tailwindcss/vite`)
- tw-animate-css `^1.4.0` — Tailwind animation utilities
- class-variance-authority `^0.7.1` — Variant-based component styling
- clsx `^2.1.1` + tailwind-merge `^3.5.0` — Class name utilities (`cn()` helper in `frontend/src/lib/utils.ts`)

**Build/Dev:**
- Vite `^5.0.0` — Frontend bundler and dev server
- vue-tsc `^1.0.11` — TypeScript type-checking for Vue SFCs
- Task (go-task.dev) v3 — Build orchestration via `Taskfile.yml` (replaces Make)
- Wails CLI (`wails3`) — Dev mode, binding generation, build, packaging

## Key Dependencies

**Critical (Go):**
- `github.com/mattn/go-sqlite3` v1.14.38 — SQLite driver (CGO required); the sole data store
- `github.com/wailsapp/wails/v3` v3.0.0-alpha.74 — Desktop runtime, service binding, window management, system tray, events
- `golang.org/x/oauth2` v0.36.0 — OAuth2 PKCE flow for Microsoft Entra ID authentication
- `github.com/zalando/go-keyring` v0.2.8 — OS keychain storage (macOS Keychain, Windows Credential Manager, Linux Secret Service) for tokens
- `github.com/spf13/viper` v1.21.0 — Configuration management (YAML file, env vars, defaults)
- `github.com/pkg/browser` v0.0.0-20240102092130 — Opens system browser for OAuth sign-in

**Critical (Frontend):**
- `@wailsio/runtime` latest — Wails frontend runtime (events, window control)
- `vue` ^3.2.45 — Core framework
- `pinia` ^3.0.4 — State management
- `vue-router` ^4.6.4 — Routing
- `@vueuse/core` ^14.2.1 — Vue composition utilities (used for `useColorMode` in theme composable)

**Infrastructure (Go, indirect):**
- `github.com/fsnotify/fsnotify` — File watching (used by Viper for config reload)
- `github.com/go-git/go-git/v5` — Git operations (Wails dependency)
- `github.com/danieljoos/wincred` — Windows credential store (go-keyring backend)
- `github.com/godbus/dbus/v5` — Linux D-Bus (go-keyring backend)

## Configuration

**Application Config:**
- Managed by Viper (`internal/config/config.go`)
- Config file: `config.yaml` in OS-appropriate config directory
  - macOS: `~/Library/Application Support/team-ado-tool/config.yaml`
  - Windows: `%APPDATA%\team-ado-tool\config.yaml`
  - Linux: `~/.config/team-ado-tool/config.yaml` (or `$XDG_CONFIG_HOME`)
- Env prefix: `XBT` (e.g., `XBT_DB_PATH` overrides `db.path`)
- Auto-creates config file with defaults on first run

**Config Keys & Defaults:**
| Key | Default | Purpose |
|-----|---------|---------|
| `db.path` | `<DataDir>/data.db` | SQLite database location |
| `theme` | `system` | UI theme (system/light/dark) |
| `window.width` | `1200` | Main window width |
| `window.height` | `800` | Main window height |
| `ado.organization` | `""` | Azure DevOps organization name |
| `ado.project` | `""` | Azure DevOps project name |
| `ado.pat_keychain_key` | `xbt-ado-pat` | Keychain key for PAT storage |
| `sync.interval_minutes` | `15` | ADO sync interval |
| `log.level` | `info` | Log verbosity |

**Auth Environment Variables:**
- `TEAM_ADO_TENANT_ID` — Microsoft Entra tenant ID (falls back to placeholder)
- `TEAM_ADO_CLIENT_ID` — OAuth2 client/app ID (falls back to placeholder)

**Build Config:**
- `build/config.yml` — Wails build configuration (app metadata, dev mode settings, file associations)
- `Taskfile.yml` — Root task runner config; includes platform-specific taskfiles from `build/`
- `frontend/vite.config.ts` — Vite bundler config with `@` path alias to `frontend/src/`
- `frontend/tsconfig.json` — TypeScript config (ESNext target, strict mode, `@/*` path alias)
- `frontend/components.json` — shadcn-vue component registry config (New York style, zinc base color, CSS variables enabled)

## Platform Requirements

**Development:**
- Go 1.25+ (CGO enabled — required for `go-sqlite3`)
- Node.js 18+ (npm required for frontend)
- Wails CLI v3 (`wails3` command)
- Task runner (`task` command from go-task.dev)
- macOS: Xcode command-line tools (for CGO/WebKit)
- Windows: WebView2 runtime, MSVC build tools
- Linux: GCC, GTK3, WebKitGTK dev headers

**Production:**
- Self-contained native desktop binary per platform
- macOS: `.app` bundle (built via `build/darwin/Taskfile.yml`)
- Windows: `.exe` with NSIS installer support (`build/windows/`)
- Linux: AppImage support (`build/linux/`)
- Docker/server mode also available (`task build:server`)

**Dev Commands:**
```bash
task dev                    # Full dev mode (Go + Vite hot-reload)
task build                  # Production build for current OS
task package                # Package for distribution
wails3 generate bindings    # Regenerate Go→TS bindings
cd frontend && npm run dev  # Frontend-only dev server
```

---

*Stack analysis: 2025-07-15*
