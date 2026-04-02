# Codebase Structure

**Analysis Date:** 2025-07-15

## Directory Layout

```
xb-tasks/
├── main.go                 # Go entry point — app setup, service registration, window creation
├── go.mod                  # Go module: dev.azure.com/xbox/xb-tasks
├── go.sum                  # Go dependency checksums
├── Taskfile.yml            # Task runner config (build, dev, package)
├── REFERENCE.md            # API reference documentation
├── copilot-instructions.md # AI assistant context
├── domain/                 # Shared Go domain types (leaf package)
│   └── types.go
├── internal/               # Go backend — private application code
│   ├── app/                # Business logic services (Wails-bound)
│   │   ├── tasks.go        # TaskService — CRUD, filtering, subtasks, tags
│   │   ├── projects.go     # ProjectService — CRUD
│   │   └── deps.go         # DependencyService — task deps + cycle detection
│   ├── auth/               # Authentication service (Wails-bound)
│   │   └── auth.go         # AuthService — OAuth2 PKCE, PAT, session management
│   ├── config/             # Configuration management
│   │   ├── config.go       # Viper init, defaults, OS paths, accessors
│   │   └── service.go      # ConfigService — Wails-bound config read/write
│   └── db/                 # Database layer
│       └── db.go           # SQLite open, migrate, schema DDL
├── pkg/                    # Go packages intended for broader use
│   └── ado/                # Azure DevOps integration (future)
│       └── domain/
│           └── types.go    # ADOPullRequest, ADOWorkItem types
├── frontend/               # Vue 3 + TypeScript SPA
│   ├── package.json        # NPM dependencies
│   ├── vite.config.ts      # Vite config (Vue plugin, Wails plugin, Tailwind, @ alias)
│   ├── tsconfig.json       # TypeScript config (ESNext, @/* path alias)
│   ├── bindings/           # ⚡ AUTO-GENERATED — Wails TypeScript bindings
│   │   └── dev.azure.com/xbox/xb-tasks/
│   │       ├── domain/models.ts          # TS classes: Task, Project, User
│   │       └── internal/
│   │           ├── app/taskservice.ts     # TaskService binding functions
│   │           ├── app/projectservice.ts  # ProjectService binding functions
│   │           ├── app/dependencyservice.ts # DependencyService binding functions
│   │           ├── auth/authservice.ts    # AuthService binding functions
│   │           └── config/configservice.ts # ConfigService binding functions
│   ├── src/
│   │   ├── main.ts         # Vue app bootstrap (Pinia, Router, mount)
│   │   ├── App.vue         # Root component — theme init, conditional AppShell
│   │   ├── style.css       # Global styles, Tailwind imports, CSS variables, themes
│   │   ├── vite-env.d.ts   # Vite type declarations
│   │   ├── router/
│   │   │   └── index.ts    # Vue Router — hash history, all route definitions
│   │   ├── stores/         # Pinia stores (Composition API)
│   │   │   ├── tasks.ts    # Task state, CRUD, filtering, grouping, mock data
│   │   │   ├── projects.ts # Project state, CRUD, mock data
│   │   │   └── auth.ts     # Auth state, sign in/out, session restore
│   │   ├── views/          # Route-level page components
│   │   │   ├── TasksView.vue              # Main task list + detail panel
│   │   │   ├── DashboardView.vue          # Overview: tasks + PR summary
│   │   │   ├── ProjectsView.vue           # Project grid + create form
│   │   │   ├── SettingsView.vue           # App configuration UI
│   │   │   ├── LoginView.vue              # Auth page (Microsoft + PAT)
│   │   │   ├── PlaygroundTasksView.vue    # Dev: task layout experiments
│   │   │   ├── PlaygroundDashboardView.vue # Dev: dashboard layout experiments
│   │   │   ├── PlaygroundDetailView.vue   # Dev: detail panel experiments
│   │   │   └── PlaygroundAdoView.vue      # Dev: ADO management experiments
│   │   ├── layouts/
│   │   │   └── AppShell.vue # Main layout: sidebar + top bar + activity panel + content slot
│   │   ├── components/     # Reusable app components
│   │   │   ├── Sidebar.vue           # Icon nav rail (dashboard, tasks, projects, settings, playground)
│   │   │   ├── CommandPalette.vue    # ⌘K command dialog (navigate, create, search tasks)
│   │   │   ├── TaskDetail.vue        # Task detail/edit panel (658 lines — largest component)
│   │   │   ├── TaskRow.vue           # Task list row with status icon, badges, time
│   │   │   ├── PageHeader.vue        # Consistent sub-bar (left slot + right slot)
│   │   │   ├── HelloWorld.vue        # Unused scaffold placeholder (9 lines)
│   │   │   └── ui/                   # UI primitive library (shadcn-vue style)
│   │   │       ├── badge/Badge.vue + index.ts
│   │   │       ├── button/Button.vue + index.ts
│   │   │       ├── card/Card.vue, CardHeader.vue, CardContent.vue, ... + index.ts
│   │   │       ├── command/          # Command palette primitives (reka-ui based)
│   │   │       ├── dialog/           # Dialog primitives
│   │   │       ├── dropdown-menu/    # Dropdown menu primitives
│   │   │       ├── input/            # Input component
│   │   │       ├── scroll-area/      # Scroll area wrapper
│   │   │       ├── select/           # Select dropdown
│   │   │       ├── separator/        # Visual separator
│   │   │       ├── tabs/             # Tab components
│   │   │       ├── textarea/         # Textarea component
│   │   │       ├── tooltip/          # Tooltip components
│   │   │       ├── AdoBadge.vue      # ADO work item ID badge
│   │   │       ├── PriorityBadge.vue # P0-P3 priority indicator
│   │   │       ├── StatusBadge.vue   # Task status dot + label
│   │   │       └── TagChip.vue       # Tag pill component
│   │   ├── composables/
│   │   │   └── useTheme.ts # Dark/light mode composable (wraps @vueuse/core)
│   │   └── lib/
│   │       └── utils.ts    # cn() utility (clsx + tailwind-merge)
│   ├── public/             # Static assets served as-is
│   └── dist/               # Build output (embedded into Go binary)
├── build/                  # Wails build configuration per platform
│   ├── config.yml          # Wails build config (referenced by Taskfile)
│   ├── darwin/             # macOS build assets + Taskfile
│   ├── windows/            # Windows build assets + Taskfile
│   ├── linux/              # Linux build assets (AppImage, nfpm, etc.)
│   ├── ios/                # iOS build scripts
│   ├── android/            # Android build config (gradle, manifests)
│   ├── docker/             # Docker build for server mode
│   └── appicon.icon/       # App icon assets
├── bin/                    # Build output directory
│   └── team-ado-tool.dev.app/ # macOS dev app bundle
└── .planning/              # Project planning documents
    ├── PROJECT.md
    ├── REQUIREMENTS.md
    ├── ROADMAP.md
    ├── STATE.md
    ├── config.json
    ├── codebase/           # Generated codebase analysis (this directory)
    ├── phases/             # Implementation phase plans
    │   ├── 01-foundation-auth-personal-tasks/
    │   └── 02-ado-integration-prs-unified-dashboard/
    └── research/           # Research notes
```

## Directory Purposes

**`domain/`:**
- Purpose: Shared domain model types used across Go packages
- Contains: Single `types.go` with `User`, `Project`, `Task`, `TaskADOLink`, `TaskDep` structs
- Key files: `domain/types.go`

**`internal/app/`:**
- Purpose: Core business logic services bound to the Wails frontend
- Contains: One file per service — `tasks.go`, `projects.go`, `deps.go`
- Key files: `internal/app/tasks.go` (268 lines, most complex service)

**`internal/auth/`:**
- Purpose: Authentication and identity management
- Contains: `auth.go` — OAuth2 PKCE flow, PAT auth, keyring storage, MS Graph profile fetch
- Key files: `internal/auth/auth.go` (327 lines)

**`internal/config/`:**
- Purpose: Application configuration via Viper (YAML + env vars)
- Contains: `config.go` (init, defaults, path resolution), `service.go` (Wails-bound wrapper)
- Key files: `internal/config/config.go`, `internal/config/service.go`

**`internal/db/`:**
- Purpose: SQLite database management — connection, migration, schema
- Contains: `db.go` with `DB` struct, `Open()`, embedded schema DDL
- Key files: `internal/db/db.go` (142 lines — schema defines 6 tables + 8 indexes)

**`pkg/ado/domain/`:**
- Purpose: ADO-specific domain types for future integration
- Contains: `types.go` with `ADOPullRequest`, `ADOWorkItem` and their status/type constants
- Key files: `pkg/ado/domain/types.go`

**`frontend/bindings/`:**
- Purpose: Auto-generated TypeScript clients for Go services
- Contains: TypeScript modules mirroring Go package paths
- **Generated: Yes — NEVER edit manually**

**`frontend/src/stores/`:**
- Purpose: Centralized frontend state management via Pinia
- Contains: One store per domain entity (tasks, projects, auth)
- Key files: `frontend/src/stores/tasks.ts` (256 lines — includes mock data + grouping logic)

**`frontend/src/views/`:**
- Purpose: Full-page components mounted by Vue Router
- Contains: 5 production views + 4 playground/experiment views
- Key files: `frontend/src/views/TasksView.vue` (main working view)

**`frontend/src/components/ui/`:**
- Purpose: Reusable UI primitives following shadcn-vue pattern
- Contains: Directory-based primitives (badge/, button/, card/, etc.) + domain-specific loose files (AdoBadge.vue, PriorityBadge.vue, StatusBadge.vue, TagChip.vue)
- Key files: Each directory has `index.ts` barrel export

**`frontend/src/components/`:**
- Purpose: App-level reusable components
- Contains: Sidebar, CommandPalette, TaskDetail, TaskRow, PageHeader
- Key files: `frontend/src/components/TaskDetail.vue` (658 lines — largest Vue component)

**`build/`:**
- Purpose: Platform-specific build configs for Wails
- Contains: Per-OS directories with Taskfile.yml + platform assets
- Generated: Partially — scaffolded by Wails, customized per platform

**`.planning/`:**
- Purpose: Project planning, requirements, roadmap, phase tracking
- Contains: Markdown planning docs + phase directories
- Committed: Yes

## Key File Locations

**Entry Points:**
- `main.go`: Go application entry — initializes everything, creates window, runs event loop
- `frontend/src/main.ts`: Vue app bootstrap — Pinia, Router, mount

**Configuration:**
- `frontend/vite.config.ts`: Vite build config (Vue, Wails, Tailwind plugins + `@` alias)
- `frontend/tsconfig.json`: TypeScript compiler options + `@/*` path alias
- `internal/config/config.go`: Runtime config (Viper YAML + env)
- `Taskfile.yml`: Build/dev task definitions
- `go.mod`: Go module definition and dependencies

**Core Logic:**
- `internal/app/tasks.go`: Task CRUD, filtering, subtasks, personal priority, tags
- `internal/app/projects.go`: Project CRUD
- `internal/app/deps.go`: Dependency management with cycle detection (DFS)
- `internal/auth/auth.go`: Full auth lifecycle
- `internal/db/db.go`: Database schema and connection

**Frontend State:**
- `frontend/src/stores/tasks.ts`: All task state and operations
- `frontend/src/stores/projects.ts`: Project state and operations
- `frontend/src/stores/auth.ts`: Authentication state

**Frontend UI:**
- `frontend/src/layouts/AppShell.vue`: Main layout wrapper
- `frontend/src/components/Sidebar.vue`: Navigation rail
- `frontend/src/components/CommandPalette.vue`: ⌘K command dialog
- `frontend/src/components/TaskDetail.vue`: Task editing panel

**Styling:**
- `frontend/src/style.css`: Global CSS, Tailwind imports, CSS custom properties, light/dark themes

## Naming Conventions

**Go Files:**
- Pattern: lowercase single-word names matching content: `tasks.go`, `projects.go`, `deps.go`, `auth.go`, `config.go`, `service.go`, `db.go`
- One service struct per file; filename matches service domain

**Go Packages:**
- Pattern: short lowercase names: `app`, `auth`, `config`, `db`, `domain`
- `internal/` for private packages, `pkg/` for potentially-shared packages

**Vue Component Files:**
- Pattern: PascalCase `.vue` files: `TaskDetail.vue`, `CommandPalette.vue`, `AppShell.vue`
- Views: PascalCase + `View` suffix: `TasksView.vue`, `DashboardView.vue`, `LoginView.vue`
- Playground views: `Playground` prefix + domain + `View` suffix: `PlaygroundTasksView.vue`
- UI primitives: PascalCase in own directory: `button/Button.vue`, `card/Card.vue`
- Domain badges: PascalCase in `ui/` root: `StatusBadge.vue`, `PriorityBadge.vue`

**TypeScript Files:**
- Stores: camelCase `.ts`: `tasks.ts`, `projects.ts`, `auth.ts`
- Composables: `use` prefix camelCase: `useTheme.ts`
- Utilities: camelCase: `utils.ts`
- Router: `index.ts` in `router/` directory

**Directories:**
- Go: lowercase single-word: `app/`, `auth/`, `config/`, `db/`, `domain/`
- Frontend components: kebab-case for UI primitives: `dropdown-menu/`, `scroll-area/`
- Frontend top-level: lowercase: `stores/`, `views/`, `components/`, `layouts/`, `composables/`, `lib/`

## Where to Add New Code

**New Go Service (Wails-bound):**
1. Create file in `internal/app/` (or new `internal/` subpackage for distinct concerns)
2. Define struct with `*db.DB` field, constructor `NewXxxService(database *db.DB) *XxxService`
3. Add exported methods — these become frontend-callable RPC
4. Register in `main.go`: `wailsApp.RegisterService(application.NewService(xxxService))`
5. Run `wails3 generate bindings` (or restart `wails3 dev`) to regenerate TypeScript clients
6. Import generated binding in store: `import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/xxxservice')`

**New Domain Type:**
- Go struct: Add to `domain/types.go` (core types) or `pkg/ado/domain/types.go` (ADO types)
- DB schema: Add `CREATE TABLE` to `internal/db/db.go` `schema` const
- Bindings will auto-generate TS class in `frontend/bindings/.../domain/models.ts`

**New Frontend Store:**
- Create `frontend/src/stores/{name}.ts`
- Use Composition API: `defineStore('{name}', () => { ... })`
- Include mock data and `useMock` flag following existing pattern
- Import in views/components as needed

**New Frontend View (Route):**
- Create `frontend/src/views/{Name}View.vue`
- Add route to `frontend/src/router/index.ts`
- Use `PageHeader` component for consistent sub-bar
- Add navigation entry to `frontend/src/components/Sidebar.vue` if it's a primary route
- Add to `CommandPalette.vue` navigation actions

**New UI Primitive Component:**
- Create `frontend/src/components/ui/{name}/` directory
- Add `{Name}.vue` + `index.ts` barrel export
- Use `class-variance-authority` for variant patterns, `reka-ui` for headless behavior
- Follow existing shadcn-vue patterns

**New Domain UI Component (badge, chip, etc.):**
- Create `frontend/src/components/ui/{Name}.vue` (loose file, not in subdirectory)
- Accept domain-specific props (status string, priority string, etc.)

**New Composable:**
- Create `frontend/src/composables/use{Name}.ts`
- Return reactive refs and functions
- Follow `useTheme.ts` pattern

**New Utility Function:**
- Add to `frontend/src/lib/utils.ts`

**New DB Table:**
- Add `CREATE TABLE IF NOT EXISTS` + indexes to `schema` const in `internal/db/db.go`
- Migration is automatic on app startup (idempotent `CREATE IF NOT EXISTS`)

## Special Directories

**`frontend/bindings/`:**
- Purpose: Auto-generated Wails TypeScript bindings
- Generated: Yes — by Wails CLI
- Committed: Yes (checked into git)
- Note: NEVER edit manually. Path structure mirrors Go module path (`dev.azure.com/xbox/xb-tasks/...`)

**`frontend/dist/`:**
- Purpose: Vite build output, embedded into Go binary via `//go:embed`
- Generated: Yes — by `vite build`
- Committed: Yes (needed for Go embedding)

**`bin/`:**
- Purpose: Built application bundles
- Generated: Yes — by Wails build
- Committed: Partially (`.app` bundle present)

**`build/`:**
- Purpose: Platform-specific build configuration
- Generated: Scaffolded by Wails, then customized
- Committed: Yes

**`.planning/`:**
- Purpose: Project management — planning docs, phase plans, research
- Generated: No — manually authored
- Committed: Yes

**`.task/`:**
- Purpose: Task runner (Taskfile) cache/state
- Generated: Yes — by `task` CLI
- Committed: Varies

---

*Structure analysis: 2025-07-15*
