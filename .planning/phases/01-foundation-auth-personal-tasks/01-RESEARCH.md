# Phase 1: Foundation, Auth & Personal Tasks - Research

**Researched:** 2026-04-01
**Domain:** Desktop app shell (Wails v3 + Vue 3), OAuth2 PKCE authentication, personal task management
**Confidence:** HIGH

## Summary

Phase 1 builds the foundational native desktop app where a user signs in with their Microsoft work account and manages personal tasks with full CRUD, hierarchy (subtasks/dependencies), organization (tags, categories, projects), and workflow states. The existing codebase provides a solid foundation: SQLite schema with tasks, projects, task_deps tables; Go services with CRUD operations already bound to the Wails frontend; and auto-generated TypeScript bindings. The primary implementation work falls into three areas: (1) building the OAuth2 PKCE auth flow for Entra ID as a new Go service, (2) extending the existing TaskService with dependency management, tag filtering, and personal priority overlay, and (3) constructing the Vue 3 desktop app shell (sidebar, command palette, dashboard, task views, theming).

The biggest technical risk is the OAuth2 PKCE flow in a desktop context — it requires launching a localhost HTTP server to receive the redirect callback, and the token must be securely stored in the OS keychain for session persistence. The existing codebase already has an empty `internal/auth/` package ready for this. The Vue frontend is a blank slate (only a placeholder `App.vue` and `HelloWorld.vue` exist), so the entire UI must be built from scratch using Tailwind CSS with a component library (Headless UI) for accessibility primitives.

**Primary recommendation:** Build auth first (gate everything behind login), then extend task services (deps, tags, personal priority), then build the UI shell and views in parallel. Use `golang.org/x/oauth2` + custom PKCE for the auth flow (MSAL Go SDK is enterprise-heavy and designed for confidential clients — overkill for a native desktop app). Store tokens via `github.com/zalando/go-keyring` for OS keychain integration.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Slim icon sidebar (Linear/Notion style) — icons for each section, content area takes most of the screen. Modern, desktop-native feel.
- **D-02:** Command palette (⌘K / Ctrl+K) for power-user search and management — navigate anywhere, create tasks, change status, all from keyboard.
- **D-03:** Minimize to system tray on window close — app stays running in background, always available (Slack/1Password pattern). Click tray icon to restore.
- **D-04:** Color scheme follows system preference (auto dark/light). Respect macOS/Windows appearance setting.
- **D-05:** Landing page is a dashboard overview — today's work, recent tasks, quick stats. Not a raw task list.

### Agent's Discretion
- Auth flow details (OAuth2 PKCE implementation, PAT fallback, token storage approach)
- Task list layout (table vs cards, grouping, sort/filter controls)
- Task detail experience (inline editing, side panel, full page)
- Tags UX (free-form vs predefined, autocomplete)
- Dependency visualization (text, tree, graph)

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| AUTH-01 | User can sign in with Microsoft Entra ID (SSO with existing work account) | OAuth2 PKCE flow via `golang.org/x/oauth2` + localhost redirect server; Entra ID app registration requirements documented |
| AUTH-02 | User session persists across app restart (token auto-refresh) | OS keychain storage via `go-keyring`; refresh token flow with `golang.org/x/oauth2`; silent token acquisition on app start |
| AUTH-03 | User profile populated from Entra ID (display name, email, avatar) | Microsoft Graph API `/me` endpoint for profile data + `/me/photo/$value` for avatar |
| TASK-01 | User can create a personal task with title, description, priority (P0-P3), and category | Existing `TaskService.Create()` already handles this — extend to accept tags |
| TASK-02 | User can edit and delete their own tasks | Existing `TaskService.Update()` and `TaskService.Delete()` work — no changes needed |
| TASK-03 | User can set task status (todo, in_progress, in_review, done, blocked, cancelled) | Existing `TaskService.SetStatus()` works — schema already has CHECK constraint for these values |
| TASK-04 | User can create subtasks under a parent task or ADO-linked item | Existing schema has `parent_id` column; `TaskService.Create()` accepts `parentID`; `GetSubtasks()` exists |
| TASK-05 | User can set a personal priority overlay independent of ADO priority | Existing schema has `personal_priority` column; need `TaskService.SetPersonalPriority()` method |
| TASK-06 | User can define dependencies between tasks (task A blocks task B) | Existing `task_deps` table; need new `DependencyService` with Add/Remove/GetDependencies/GetBlockedBy + circular detection |
| TASK-07 | User can add tags to tasks for organization | Existing schema has `tags` TEXT column (CSV); `TaskService.Update()` already accepts tags param; need tag filtering in List |
</phase_requirements>

## 1. Entra ID OAuth2 PKCE for Desktop Apps

### Recommended Approach: `golang.org/x/oauth2` + Custom PKCE
**Confidence: HIGH**

For a native desktop app (not a web app), the OAuth2 Authorization Code Flow with PKCE is the correct pattern. The flow works as follows:

1. **App starts a temporary localhost HTTP server** on a random available port (e.g., `http://localhost:54321/callback`)
2. **App opens the system browser** to the Entra ID authorization endpoint with PKCE challenge
3. **User authenticates in their browser** (existing SSO session usually means zero interaction)
4. **Entra ID redirects to localhost** with the authorization code
5. **App exchanges code for tokens** (access token + refresh token)
6. **App stores tokens in OS keychain** via `go-keyring`
7. **On subsequent launches**, app loads refresh token from keychain and silently acquires a new access token

### Why NOT MSAL Go
The Microsoft Authentication Library for Go (`github.com/AzureAD/microsoft-authentication-library-for-go`) is available but designed primarily for server-side confidential client flows and has complex abstractions. For a native desktop app with a simple PKCE flow, `golang.org/x/oauth2` with manual PKCE is simpler, more transparent, and gives full control. The MSAL Go SDK's public client support exists but is less documented and tested for desktop scenarios.

**Decision: Use `golang.org/x/oauth2` with manual PKCE code verifier/challenge generation.**

### Azure AD App Registration Requirements

| Setting | Value |
|---------|-------|
| Application type | Mobile/Desktop |
| Redirect URI | `http://localhost` (with dynamic port) |
| Supported account types | Single tenant (org only) |
| Platform | Mobile and desktop applications |
| API permissions | `User.Read` (Microsoft Graph — for profile), `https://app.vssps.visualstudio.com/user_impersonation` (for ADO API in Phase 2) |
| Token configuration | Enable refresh tokens |
| Client secret | **None** — PKCE flow is a public client, no secret needed |

### Token Storage Strategy

| Approach | Platform | Implementation |
|----------|----------|----------------|
| OS Keychain | macOS | `go-keyring` → macOS Keychain |
| Credential Manager | Windows | `go-keyring` → Windows Credential Manager |
| Secret Service | Linux | `go-keyring` → D-Bus Secret Service |

Store: `access_token`, `refresh_token`, `token_expiry` as separate keychain entries under service name `team-ado-tool`.

### Token Refresh Flow

```go
// Pseudocode for silent token acquisition on app start
func (a *AuthService) TryRestoreSession() (*User, error) {
    refreshToken, err := keyring.Get("team-ado-tool", "refresh_token")
    if err != nil {
        return nil, nil // No saved session, user needs to sign in
    }
    
    token := &oauth2.Token{RefreshToken: refreshToken}
    tokenSource := a.oauthConfig.TokenSource(context.Background(), token)
    newToken, err := tokenSource.Token() // Auto-refreshes
    if err != nil {
        keyring.Delete("team-ado-tool", "refresh_token") // Expired, clear
        return nil, nil
    }
    
    a.saveTokens(newToken)
    return a.fetchUserProfile(newToken.AccessToken)
}
```

### User Profile (AUTH-03)

Fetch from Microsoft Graph API after obtaining access token with `User.Read` scope:
- `GET https://graph.microsoft.com/v1.0/me` → display name, email
- `GET https://graph.microsoft.com/v1.0/me/photo/$value` → avatar (binary, base64-encode for frontend)

Store in local `users` table (already exists in schema).

### PAT Fallback

For development/dogfooding, support a PAT-based auth bypass:
- User enters PAT in settings
- Store in keychain, use for ADO API calls
- Skip OAuth flow entirely
- Profile data unavailable (PAT doesn't map to Graph API) — use placeholder

## 2. Wails v3 Frontend-Backend Integration

### Service Binding Pattern (Verified from codebase)
**Confidence: HIGH**

Wails v3 alpha.74 (currently in `go.mod`) uses this binding pattern:

```go
// Go: Register services in main.go
wailsApp := application.New(application.Options{
    Services: []application.Service{
        application.NewService(taskService),
        application.NewService(projectService),
        application.NewService(authService),      // NEW
        application.NewService(dependencyService), // NEW
    },
})
```

The Wails CLI generates TypeScript bindings automatically. Each exported method on a bound struct becomes a callable function in the frontend:

```typescript
// Generated: frontend/bindings/dev.azure.com/microsoft/Xbox/xb-tasks/internal/app/taskservice.ts
import * as TaskService from '../bindings/dev.azure.com/microsoft/Xbox/xb-tasks/internal/app/taskservice'
const tasks = await TaskService.List("")
```

**Key constraint:** All Go methods must return `(T, error)` or `error` to generate proper TypeScript signatures. The generated code handles JSON serialization automatically.

### Event System (Backend → Frontend)
**Confidence: HIGH** (verified from Wails v3 source)

Wails v3 has a built-in event system for backend-to-frontend communication:

```go
// Go: Emit event from backend
app.Event.Emit("auth:state-changed", map[string]any{"authenticated": true})
app.Event.Emit("task:updated", taskID)
```

```typescript
// Frontend: Listen for events
import { Events } from '@wailsio/runtime'
Events.On("auth:state-changed", (data) => { ... })
Events.On("task:updated", (taskID) => { ... })
```

Use events for: auth state changes, background operations completing, data invalidation signals.

### System Tray (D-03)
**Confidence: HIGH** (verified from Wails v3 source — `systemtray.go`, `system_tray_manager.go`)

Wails v3 has full system tray support:

```go
// Create system tray
tray := wailsApp.SystemTray.New()
tray.SetIcon(iconBytes)
tray.SetDarkModeIcon(darkIconBytes)
tray.SetTooltip("Team ADO Tool")
tray.SetMenu(trayMenu)

// Click to show/hide window
tray.OnClick(func() {
    mainWindow.Show()
    mainWindow.Focus()
})
```

For "minimize to tray on close" (D-03), use a window hook to intercept the close event:

```go
// Hook into WindowClosing to hide instead of close
mainWindow.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    event.Cancel()      // Prevent actual close
    mainWindow.Hide()   // Hide window instead
})
```

**Critical:** Set `ApplicationShouldTerminateAfterLastWindowClosed: false` in Mac options (currently `true` in `main.go` — must change):

```go
Mac: application.MacOptions{
    ApplicationShouldTerminateAfterLastWindowClosed: false, // Keep running in tray
},
```

Add a "Quit" option in the tray menu for actual termination.

### Window Configuration

The existing `main.go` already has good Mac-specific settings:
- `InvisibleTitleBarHeight: 50` — space for custom titlebar
- `Backdrop: MacBackdropTranslucent` — native macOS feel
- `TitleBar: MacTitleBarHiddenInset` — modern chromeless look

These work well with the slim sidebar layout.

## 3. Vue 3 Desktop App Shell

### App Shell Architecture
**Confidence: HIGH**

```
frontend/src/
├── App.vue                    # Root: sidebar + router-view + command palette
├── main.ts                    # App setup: vue-router, pinia, global styles
├── router/
│   └── index.ts               # Route definitions
├── stores/
│   ├── auth.ts                # Auth state (user, token status)
│   ├── tasks.ts               # Task list state, filters
│   └── ui.ts                  # Theme, sidebar state, command palette visibility
├── composables/
│   ├── useAuth.ts             # Auth actions (sign in, sign out, restore session)
│   ├── useTasks.ts            # Task CRUD wrapper around Wails bindings
│   └── useTheme.ts            # Dark/light theme detection + switching
├── components/
│   ├── layout/
│   │   ├── AppSidebar.vue     # Slim icon sidebar (D-01)
│   │   ├── AppHeader.vue      # Page header with breadcrumbs
│   │   └── UserAvatar.vue     # Profile display
│   ├── command/
│   │   └── CommandPalette.vue # ⌘K command palette (D-02)
│   ├── tasks/
│   │   ├── TaskList.vue       # Task list with filters
│   │   ├── TaskDetail.vue     # Side panel task detail
│   │   ├── TaskForm.vue       # Create/edit form
│   │   ├── TaskStatusBadge.vue
│   │   ├── TaskPriorityBadge.vue
│   │   ├── TaskDependencies.vue
│   │   └── TaskTags.vue       # Tag display + editing
│   ├── projects/
│   │   ├── ProjectList.vue
│   │   └── ProjectDetail.vue
│   └── dashboard/
│       ├── DashboardView.vue  # Landing page (D-05)
│       ├── TodaysSummary.vue
│       └── QuickStats.vue
├── views/
│   ├── LoginView.vue          # Auth screen
│   ├── DashboardView.vue      # Dashboard page
│   ├── TasksView.vue          # Task management page
│   └── ProjectsView.vue       # Projects page
└── assets/
    └── icons/                 # App icons, tray icons
```

### Required Frontend Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `vue-router` | ^4.5 | Client-side routing for views |
| `pinia` | ^3.0 | State management (auth, tasks, ui) |
| `@vueuse/core` | ^14.2 | Composables: `useColorMode`, `useMagicKeys`, `onKeyStroke` |
| `@headlessui/vue` | ^1.7 | Accessible UI primitives: Dialog, Combobox (for command palette), Transition |
| `lucide-vue-next` | ^1.0 | Icon library (lightweight, tree-shakeable — Linear/Notion aesthetic) |

**Install command:**
```bash
cd frontend && npm install vue-router pinia @vueuse/core @headlessui/vue lucide-vue-next
```

### Sidebar (D-01) — Slim Icon Sidebar

Pattern: Fixed-width sidebar (~56px icon-only, expandable to ~200px on hover or toggle). Each icon represents a navigation section.

```vue
<!-- AppSidebar.vue pattern -->
<nav class="flex flex-col items-center w-14 bg-slate-900 border-r border-slate-800 py-4 gap-2">
  <SidebarIcon icon="LayoutDashboard" label="Dashboard" to="/" />
  <SidebarIcon icon="CheckSquare" label="Tasks" to="/tasks" />
  <SidebarIcon icon="FolderKanban" label="Projects" to="/projects" />
  <div class="mt-auto">
    <SidebarIcon icon="Settings" label="Settings" to="/settings" />
    <UserAvatar />
  </div>
</nav>
```

Icons from `lucide-vue-next`: `LayoutDashboard`, `CheckSquare`, `FolderKanban`, `Settings`, `Search`, `Plus`.

### Command Palette (D-02) — ⌘K

Use `@headlessui/vue` `Combobox` + `Dialog` for the command palette:

```typescript
// Keyboard shortcut via @vueuse/core
import { useMagicKeys } from '@vueuse/core'
const { meta_k, ctrl_k } = useMagicKeys()
watch([meta_k, ctrl_k], ([mk, ck]) => {
  if (mk || ck) showCommandPalette.value = true
})
```

Command palette actions for Phase 1:
- Navigate to Dashboard / Tasks / Projects
- Create new task
- Search tasks by title
- Change task status (quick action)
- Sign out

### Theme System (D-04) — Follow System Preference

```typescript
// useTheme.ts using @vueuse/core
import { useColorMode } from '@vueuse/core'
const mode = useColorMode({ attribute: 'class' }) // Adds 'dark' class to <html>
```

Tailwind v4 (installed) supports `dark:` variant natively. Use CSS custom properties for theme tokens:

```css
/* style.css */
@import "tailwindcss";

:root {
  --color-bg-primary: theme(colors.white);
  --color-bg-sidebar: theme(colors.slate.50);
  --color-text-primary: theme(colors.slate.900);
}

.dark {
  --color-bg-primary: theme(colors.slate.900);
  --color-bg-sidebar: theme(colors.slate.950);
  --color-text-primary: theme(colors.slate.100);
}
```

### Dashboard (D-05) — Overview Landing

Phase 1 dashboard shows personal task data only (no ADO integration yet):

- **Today's Focus**: Tasks with status `in_progress` or `in_review`
- **Recent Activity**: Last 5 tasks modified
- **Quick Stats**: Total tasks by status (pie/bar), tasks by priority
- **Blocked Items**: Tasks with status `blocked` (needs attention)
- **Quick Create**: Inline task creation

### Task Detail — Side Panel (Agent's Discretion Recommendation)

**Recommendation: Slide-out side panel** (right side, ~40% width). This works best with sidebar + content area layout because:
- User can see the task list while editing a task
- No full-page context switch
- Matches Linear/Notion patterns
- Works well on desktop screen sizes (1200px+ window width)

### Task List Layout (Agent's Discretion Recommendation)

**Recommendation: Compact table with row actions** (not cards). Reasons:
- Higher information density — critical for a productivity tool
- Priority-sorted table is the proven pattern from xl
- Cards waste space for text-heavy task data
- Table rows with inline status/priority badges + expand for subtasks

## 4. Schema & Data Layer Assessment

### Current Schema — What Already Works

| Table | Status | Notes |
|-------|--------|-------|
| `tasks` | ✅ Complete | Has all needed columns: title, description, status (with CHECK), priority (with CHECK), category, project_id, parent_id, personal_priority, tags, ado_id, blocked_reason, blocked_by |
| `projects` | ✅ Complete | CRUD ready with status enum |
| `task_deps` | ✅ Complete | Composite PK (task_id, depends_on) with cascading deletes |
| `pull_requests` | ✅ Schema ready | Not needed for Phase 1, but won't interfere |
| `ado_work_items` | ✅ Schema ready | Not needed for Phase 1 |
| `task_ado_links` | ✅ Schema ready | Not needed for Phase 1 |
| `users` | ✅ Schema ready | id (TEXT PK), display_name, email, avatar_url — perfect for storing Entra ID profile |

### What's Missing — Schema Additions Needed

**None.** The existing schema already covers all Phase 1 requirements:
- Tags: `tags TEXT` column (CSV) — already in schema
- Subtasks: `parent_id INTEGER REFERENCES tasks(id)` — already in schema
- Dependencies: `task_deps` table — already in schema
- Personal priority: `personal_priority TEXT` — already in schema
- User profile: `users` table — already in schema

The schema is well-designed and needs no modifications for Phase 1.

### Tags Strategy (Agent's Discretion)

**Recommendation: Keep CSV column, add helper functions.** The existing `tags TEXT` column stores comma-separated values. This is simpler than a separate tags table for v1:

- **Pros of CSV**: Already implemented, simple queries with `LIKE`, no JOINs
- **Cons**: No tag-level indexing, harder to query "all tasks with tag X" efficiently
- **For Phase 1**: CSV is fine. At scale (1000+ tasks), migrate to a junction table
- **Tag autocomplete**: Query `SELECT DISTINCT` by splitting tags in Go service, cache suggestions

```go
// Helper to get all unique tags
func (s *TaskService) GetAllTags() ([]string, error) {
    rows, _ := s.db.Query("SELECT DISTINCT tags FROM tasks WHERE tags != ''")
    // Split each CSV row, deduplicate, return sorted
}
```

## 5. Task Dependencies Implementation

### Data Model (Already Exists)

The `task_deps` table with `(task_id, depends_on)` composite PK means: "task_id depends on depends_on" — i.e., `depends_on` **blocks** `task_id`.

### DependencyService Design

```go
type DependencyService struct {
    db *db.DB
}

// AddDependency: task depends on dependency (dependency blocks task)
func (s *DependencyService) AddDependency(taskID, dependsOnID int) error
// RemoveDependency: remove a dependency
func (s *DependencyService) RemoveDependency(taskID, dependsOnID int) error
// GetDependencies: what does this task depend on? (what blocks it)
func (s *DependencyService) GetDependencies(taskID int) ([]models.Task, error)
// GetDependents: what tasks depend on this one? (what does it block)
func (s *DependencyService) GetDependents(taskID int) ([]models.Task, error)
// HasCircularDependency: would adding this dep create a cycle?
func (s *DependencyService) HasCircularDependency(taskID, dependsOnID int) (bool, error)
```

### Circular Dependency Detection

Use a depth-first traversal from `dependsOnID` through the dependency graph. If we reach `taskID`, adding the dependency would create a cycle:

```go
func (s *DependencyService) HasCircularDependency(taskID, dependsOnID int) (bool, error) {
    // Would adding "taskID depends on dependsOnID" create a cycle?
    // Check: can we reach taskID by following dependencies from dependsOnID?
    visited := make(map[int]bool)
    return s.canReach(dependsOnID, taskID, visited)
}

func (s *DependencyService) canReach(from, target int, visited map[int]bool) (bool, error) {
    if from == target { return true, nil }
    if visited[from] { return false, nil }
    visited[from] = true
    
    rows, _ := s.db.Query("SELECT depends_on FROM task_deps WHERE task_id = ?", from)
    defer rows.Close()
    for rows.Next() {
        var depID int
        rows.Scan(&depID)
        if reached, _ := s.canReach(depID, target, visited); reached {
            return true, nil
        }
    }
    return false, nil
}
```

**Performance note:** For Phase 1 with <100 tasks, recursive traversal is fine. No need for a graph database.

### Blocked Status Surfacing

When displaying a task, check if any of its dependencies are not `done`:

```go
func (s *DependencyService) GetBlockingTasks(taskID int) ([]models.Task, error) {
    // Return dependencies that are NOT in 'done' status
    rows, _ := s.db.Query(`
        SELECT t.id, t.title, t.status, t.priority 
        FROM tasks t 
        JOIN task_deps td ON td.depends_on = t.id 
        WHERE td.task_id = ? AND t.status != 'done'`, taskID)
    // ...
}
```

### Dependency Visualization (Agent's Discretion)

**Recommendation: Simple text list for v1.** Show dependencies as:
- "Blocked by: [Task #12: Fix auth], [Task #15: Update schema]" — clickable links
- "Blocks: [Task #20: Deploy], [Task #22: Write tests]"
- Visual indicator (🔴 icon or red badge) on tasks that have unresolved blockers

A graph visualization is overkill for v1 with <100 tasks.

## 6. Existing Code Assessment

### ✅ Reusable Directly (No Changes)

| Component | File | What It Does |
|-----------|------|-------------|
| DB layer | `internal/db/db.go` | SQLite with WAL, foreign keys, auto-migration — production-ready |
| Task CRUD | `internal/app/tasks.go` | Create, GetByID, List, Update, Delete, SetStatus, GetSubtasks — all working |
| Project CRUD | `internal/app/projects.go` | Create, GetByID, List, Update, Delete — all working |
| Models | `pkg/models/models.go` | Task, Project, PullRequest, ADOWorkItem, TaskADOLink, User, TaskDep — comprehensive |
| Schema | `internal/db/db.go` (const schema) | All tables, indexes, constraints — no changes needed |
| Wails bindings | `frontend/bindings/` | Auto-generated TypeScript for TaskService, ProjectService |
| Tailwind setup | `frontend/src/style.css`, `vite.config.ts` | Tailwind v4 with Vite plugin |
| Vite config | `frontend/vite.config.ts` | Vue + Wails + Tailwind plugins configured |

### 🔧 Needs Extending

| Component | What Exists | What's Needed |
|-----------|-------------|---------------|
| `TaskService` | CRUD, status, subtasks | `SetPersonalPriority(id int, priority string)` method |
| `TaskService` | `List(status)` | `ListFiltered(filter TaskFilter)` with tags, priority, category filters |
| `TaskService` | Tags passed in Update | `GetAllTags()` for autocomplete |
| `main.go` | TaskService + ProjectService registered | Add AuthService, DependencyService; add SystemTray; change Mac terminate option |

### 🆕 Must Build From Scratch

| Component | Location | Complexity | Notes |
|-----------|----------|------------|-------|
| **AuthService** | `internal/auth/auth.go` | HIGH | OAuth2 PKCE flow, localhost callback server, token refresh, keychain storage, profile fetch |
| **DependencyService** | `internal/app/dependencies.go` | MEDIUM | Add/Remove deps, circular detection, blocked status queries |
| **Vue Router setup** | `frontend/src/router/` | LOW | 4-5 routes (login, dashboard, tasks, projects, settings) |
| **Pinia stores** | `frontend/src/stores/` | MEDIUM | Auth store, task store, UI store |
| **App Shell** | `frontend/src/App.vue` | MEDIUM | Sidebar + router-view + command palette overlay |
| **AppSidebar** | `frontend/src/components/layout/` | LOW | Icon sidebar with navigation |
| **CommandPalette** | `frontend/src/components/command/` | MEDIUM | Headless UI Combobox + Dialog, action registry |
| **DashboardView** | `frontend/src/views/` | MEDIUM | Today's focus, stats, recent activity |
| **TasksView + components** | `frontend/src/views/`, `components/tasks/` | HIGH | List, detail panel, form, status/priority badges, deps, tags |
| **ProjectsView** | `frontend/src/views/` | LOW | List + basic detail |
| **Theme system** | `frontend/src/composables/useTheme.ts` | LOW | @vueuse/core useColorMode + CSS variables |
| **Login page** | `frontend/src/views/LoginView.vue` | LOW | Sign-in button, loading state, error display |

## Standard Stack

### Core (Already in project)

| Library | Version | Purpose | Status |
|---------|---------|---------|--------|
| Wails v3 | alpha.74 | Native desktop shell | In go.mod |
| Vue 3 | 3.5.31 | Frontend framework | In package.json |
| Tailwind CSS | 4.2.2 | Utility-first CSS | In package.json |
| @wailsio/runtime | 3.0.0-alpha.79 | Wails frontend runtime | In package.json |
| go-sqlite3 | 1.14.38 | SQLite driver | In go.mod |
| TypeScript | 4.9.5 | Type safety | In package.json |
| Vite | 5.4.21 | Build tool | In package.json |

### To Add — Go

| Library | Purpose | Why |
|---------|---------|-----|
| `golang.org/x/oauth2` | OAuth2 PKCE flow | Standard Go OAuth2 library, supports custom PKCE |
| `github.com/zalando/go-keyring` | OS keychain (macOS Keychain, Windows Credential Manager) | Cross-platform, well-maintained, simple API |
| `github.com/pkg/browser` | Open system browser for OAuth | Already an indirect dependency in go.mod |

```bash
go get golang.org/x/oauth2 github.com/zalando/go-keyring
```

### To Add — Frontend

| Library | Version | Purpose | Why |
|---------|---------|---------|-----|
| `vue-router` | ^4.5 | Client-side routing | Required for multi-view app |
| `pinia` | ^3.0 | State management | Vue 3 standard, replaces Vuex |
| `@vueuse/core` | ^14.2 | Composables (useColorMode, useMagicKeys) | Avoid hand-rolling browser APIs |
| `@headlessui/vue` | ^1.7 | Accessible Dialog, Combobox for command palette | Unstyled + accessible primitives |
| `lucide-vue-next` | ^1.0 | Icon library | Lightweight, Linear-style aesthetic, tree-shakeable |

```bash
cd frontend && npm install vue-router@^4.5 pinia@^3.0 @vueuse/core@^14.2 @headlessui/vue@^1.7 lucide-vue-next@^1.0
```

## Architecture Patterns

### Recommended Project Structure (Go Backend)

```
internal/
├── auth/
│   └── auth.go           # AuthService: OAuth2 PKCE, token mgmt, profile
├── app/
│   ├── tasks.go          # TaskService (existing, extend)
│   ├── projects.go       # ProjectService (existing)
│   └── dependencies.go   # DependencyService (new)
├── db/
│   └── db.go             # Database (existing)
pkg/
└── models/
    └── models.go         # All models (existing)
```

### Pattern: Wails Service Registration

Every Go service that the frontend needs must be registered in `main.go`. The Wails build tool auto-generates TypeScript bindings for all exported methods.

```go
// main.go — service registration pattern
authService := auth.NewAuthService(database, wailsApp)
taskService := app.NewTaskService(database)
projectService := app.NewProjectService(database)
depService := app.NewDependencyService(database)

wailsApp := application.New(application.Options{
    Services: []application.Service{
        application.NewService(authService),
        application.NewService(taskService),
        application.NewService(projectService),
        application.NewService(depService),
    },
})
```

### Pattern: Frontend State Management

```
Vue Component → calls Wails binding → Go service method → SQLite
              ← receives return value ← (or error)
              
Backend events → Wails Event.Emit → Frontend Events.On → Pinia store update → Vue reactivity
```

The Vue layer is a **thin shell**. All business logic lives in Go. Vue components:
1. Call Wails bindings (auto-generated TypeScript)
2. Store results in Pinia stores
3. React to events from backend
4. Render UI

### Anti-Patterns to Avoid

- **Don't duplicate business logic in Vue**: All validation, state transitions, and data manipulation happen in Go services
- **Don't use REST/fetch in the frontend**: Use only Wails bindings (auto-generated from Go)
- **Don't store sensitive data in Pinia**: Tokens stay in Go (keychain). Frontend only knows "authenticated: yes/no"
- **Don't hand-roll accessibility**: Use Headless UI for dialogs, comboboxes, menus

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| OAuth2 PKCE code verifier/challenge | Custom crypto | `golang.org/x/oauth2` PKCE extensions | Crypto is easy to get wrong |
| OS keychain access | File-based token storage | `go-keyring` | Cross-platform, secure, battle-tested |
| Dark/light theme detection | `matchMedia` listener | `@vueuse/core useColorMode` | Handles SSR, persistence, system changes |
| Keyboard shortcuts | Manual `addEventListener` | `@vueuse/core useMagicKeys` | Cross-platform key combos, reactive |
| Accessible command palette | Custom modal + input | `@headlessui/vue Combobox + Dialog` | Focus management, ARIA, keyboard navigation |
| System tray | Native bindings | Wails v3 `SystemTray` API | Already in framework |
| TypeScript bindings | Manual interface definitions | Wails auto-generated bindings | Stays in sync with Go types |

## Common Pitfalls

### Pitfall 1: OAuth2 Redirect Port Conflicts
**What goes wrong:** The localhost callback server picks a fixed port that's already in use, causing auth flow to fail silently.
**Why it happens:** Other dev tools or instances of the app bind to the same port.
**How to avoid:** Use port 0 (OS picks available port), include the dynamic port in the redirect URI. Entra ID supports `http://localhost` redirect URIs with any port for native apps.
**Warning signs:** "Connection refused" or "port already in use" errors during sign-in.

### Pitfall 2: Wails Binding Regeneration
**What goes wrong:** Frontend code calls methods that don't exist or have wrong signatures because bindings weren't regenerated after Go changes.
**Why it happens:** Wails generates TypeScript bindings at build time. Changing Go service signatures requires a rebuild.
**How to avoid:** Run `wails3 generate bindings` (or `wails3 dev` which auto-regenerates) after any Go service method change.
**Warning signs:** TypeScript errors about missing methods, or runtime errors about method ID mismatches.

### Pitfall 3: `ApplicationShouldTerminateAfterLastWindowClosed: true`
**What goes wrong:** App terminates when user closes window, defeating the system tray behavior.
**Why it happens:** Current `main.go` has this set to `true`. Must be `false` for tray-persistent behavior.
**How to avoid:** Set to `false` in MacOptions. Add explicit "Quit" menu item in system tray.
**Warning signs:** App disappears from system tray when window is closed.

### Pitfall 4: Token Not Refreshing on macOS Keychain Access Prompt
**What goes wrong:** macOS prompts for keychain access permission, user denies or delays, app treats it as "no token" and forces re-login.
**Why it happens:** First run after install triggers macOS security prompt.
**How to avoid:** Handle keychain access errors gracefully — show explanatory message, retry, don't clear tokens on transient access failures.
**Warning signs:** Users reporting needing to sign in every time on first use.

### Pitfall 5: Circular Dependencies Not Detected Before Insert
**What goes wrong:** User adds A→B→C→A dependency chain, creating an unresolvable deadlock.
**Why it happens:** No pre-insert validation on the dependency graph.
**How to avoid:** Always run `HasCircularDependency()` before `INSERT INTO task_deps`. Return clear error message to UI: "Can't add this dependency — it would create a circular chain: A → B → C → A".
**Warning signs:** Tasks stuck in "blocked" status with no resolution path.

## Code Examples

### AuthService Skeleton

```go
// internal/auth/auth.go
package auth

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "net"
    "net/http"
    
    "github.com/zalando/go-keyring"
    "golang.org/x/oauth2"
    
    "dev.azure.com/microsoft/Xbox/xb-tasks/internal/db"
    "dev.azure.com/microsoft/Xbox/xb-tasks/pkg/models"
)

const (
    serviceName = "team-ado-tool"
    tenantID    = "YOUR_TENANT_ID"
    clientID    = "YOUR_CLIENT_ID"
)

type AuthService struct {
    db          *db.DB
    oauthConfig *oauth2.Config
    currentUser *models.User
}

func NewAuthService(database *db.DB) *AuthService {
    return &AuthService{
        db: database,
        oauthConfig: &oauth2.Config{
            ClientID: clientID,
            Endpoint: oauth2.Endpoint{
                AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
                TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
            },
            Scopes: []string{"User.Read", "offline_access"},
        },
    }
}

func (s *AuthService) SignIn() (*models.User, error) {
    // 1. Generate PKCE verifier + challenge
    verifier := generateCodeVerifier()
    challenge := generateCodeChallenge(verifier)
    
    // 2. Start localhost callback server
    listener, _ := net.Listen("tcp", "localhost:0")
    port := listener.Addr().(*net.TCPAddr).Port
    s.oauthConfig.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)
    
    // 3. Open browser to auth URL
    authURL := s.oauthConfig.AuthCodeURL("state",
        oauth2.SetAuthURLParam("code_challenge", challenge),
        oauth2.SetAuthURLParam("code_challenge_method", "S256"),
    )
    // browser.OpenURL(authURL)
    
    // 4. Wait for callback, exchange code for token
    // 5. Store tokens in keychain
    // 6. Fetch user profile from Microsoft Graph
    // 7. Return user
    return nil, nil // placeholder
}

func (s *AuthService) GetCurrentUser() *models.User {
    return s.currentUser
}

func (s *AuthService) IsAuthenticated() bool {
    return s.currentUser != nil
}

func (s *AuthService) SignOut() error {
    keyring.Delete(serviceName, "access_token")
    keyring.Delete(serviceName, "refresh_token")
    s.currentUser = nil
    return nil
}

func generateCodeVerifier() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
    h := sha256.Sum256([]byte(verifier))
    return base64.RawURLEncoding.EncodeToString(h[:])
}
```

### System Tray + Window Hide on Close

```go
// In main.go setup
tray := wailsApp.SystemTray.New()
tray.SetTemplateIcon(trayIconBytes) // macOS template icon (monochrome)
tray.SetTooltip("Team ADO Tool")

trayMenu := wailsApp.Menu.New()
trayMenu.Add("Show Window").OnClick(func(_ *application.Context) {
    mainWindow.Show()
    mainWindow.Focus()
})
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(_ *application.Context) {
    wailsApp.Quit()
})
tray.SetMenu(trayMenu)

// Intercept window close → hide instead
mainWindow.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    event.Cancel()
    mainWindow.Hide()
})

tray.OnClick(func() {
    mainWindow.Show()
    mainWindow.Focus()
})
```

### Task Filtering Example (Extended List)

```go
// Extended TaskService method
type TaskFilter struct {
    Status   string
    Priority string
    Category string
    Tags     string // Filter tasks containing this tag
    ParentID *int   // Only subtasks of this parent
}

func (s *TaskService) ListFiltered(f TaskFilter) ([]models.Task, error) {
    query := `SELECT id, title, description, status, priority, category, 
              project_id, area, due_date, ado_id, tags, blocked_reason, 
              blocked_by, parent_id, personal_priority, created_at, 
              updated_at, completed_at FROM tasks`
    
    var conditions []string
    var args []any
    
    if f.Status != "" {
        conditions = append(conditions, "status = ?")
        args = append(args, f.Status)
    } else {
        conditions = append(conditions, "status NOT IN ('cancelled')")
    }
    if f.Priority != "" {
        conditions = append(conditions, "priority = ?")
        args = append(args, f.Priority)
    }
    if f.Tags != "" {
        conditions = append(conditions, "(',' || tags || ',') LIKE '%,' || ? || ',%'")
        args = append(args, f.Tags)
    }
    // ... build and execute query
}
```

## Risks & Mitigations

| Risk | Severity | Likelihood | Mitigation |
|------|----------|------------|------------|
| Wails v3 alpha breaking changes | HIGH | MEDIUM | Pin to alpha.74, avoid bleeding-edge APIs, have v2 fallback plan |
| Entra ID app registration requires admin consent | MEDIUM | HIGH | Get IT admin to pre-approve app registration before auth development begins |
| OAuth2 PKCE localhost redirect blocked by corporate proxy/firewall | MEDIUM | LOW | Localhost traffic is typically not proxied; add PAT fallback as escape hatch |
| macOS Keychain prompts frustrating users | LOW | MEDIUM | Request keychain access at app install, handle errors gracefully |
| Vue binding regeneration causing stale types | LOW | HIGH | Add `wails3 generate bindings` to development workflow, pre-build step |

## Open Questions

1. **Entra ID Tenant ID and Client ID**
   - What we know: Need single-tenant app registration for the Xbox Services org
   - What's unclear: Whether Luis has permissions to create the app registration, or needs IT admin
   - Recommendation: Create app registration early (wave 1). If blocked, use PAT fallback for initial development

2. **Wails v3 `wails3` CLI availability**
   - What we know: `wails` CLI not found in PATH on this machine
   - What's unclear: Whether `wails3` is installed separately or needs installation
   - Recommendation: Install via `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` as first task

3. **ADO Scopes for Phase 1**
   - What we know: Phase 1 doesn't need ADO API access (no ADO integration yet)
   - What's unclear: Whether to request ADO scopes now or add them in Phase 2
   - Recommendation: Request `User.Read` + `offline_access` only for Phase 1. Add ADO scopes in Phase 2 (users may need to re-consent)

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Backend compilation | ✓ | 1.26.0 | — |
| Node.js | Frontend build | ✓ | 25.6.1 | — |
| npm | Frontend deps | ✓ | 11.9.0 | — |
| Wails v3 CLI | Build/dev server | ✗ | — | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` |
| SQLite | Database | ✓ (via go-sqlite3) | 1.14.38 | — |
| System browser | OAuth2 flow | ✓ | — | — |
| macOS Keychain | Token storage | ✓ (macOS native) | — | Encrypted file fallback |

**Missing dependencies with no fallback:**
- None (Wails CLI has clear install path)

**Missing dependencies with fallback:**
- Wails v3 CLI: Install with `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

## Sources

### Primary (HIGH confidence)
- Wails v3 source code (alpha.74 in go.mod cache) — system tray API, event system, window management, service binding verified directly from source
- Existing codebase (`internal/db/db.go`, `internal/app/tasks.go`, `pkg/models/models.go`) — schema, models, services verified
- XL-CODE-REFERENCE.md — proven patterns for task deps, filtering, ADO integration
- STACK.md — final stack decisions (Wails v3 + Vue 3 + Go + SQLite)

### Secondary (MEDIUM confidence)
- `golang.org/x/oauth2` — standard Go OAuth2 library, PKCE support is well-documented
- `go-keyring` — widely used for cross-platform keychain access
- Microsoft identity platform documentation — OAuth2 PKCE for native apps

### Tertiary (LOW confidence)
- MSAL Go SDK assessment — based on training data, should verify current public client support if considering

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — verified from go.mod, package.json, and live npm registry
- Architecture: HIGH — verified from Wails v3 source code and existing codebase patterns
- Auth flow: HIGH — OAuth2 PKCE is well-established pattern; Go libraries verified
- Pitfalls: HIGH — verified from Wails source (terminate behavior, hooks) and PITFALLS.md research

**Research date:** 2026-04-01
**Valid until:** 2026-05-01 (30 days — Wails v3 alpha may release new versions)
