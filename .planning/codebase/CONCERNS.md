# Codebase Concerns

**Analysis Date:** 2025-07-15

## Tech Debt

### Stores Hardcoded to Mock Data (Critical)
- Issue: Both `useTaskStore` and `useProjectStore` have `useMock = true` as a module-level variable, causing all CRUD operations to use in-memory mock data instead of Wails bindings to the Go backend. The app never hits SQLite in development.
- Files: `frontend/src/stores/tasks.ts` (line 110), `frontend/src/stores/projects.ts` (line 27)
- Impact: All task/project data is ephemeral — lost on page reload. The backend services (`internal/app/tasks.go`, `internal/app/projects.go`) are never exercised from the UI. Bugs in the Go<->frontend data contract will only surface when switching to real bindings.
- Fix approach: Replace the `useMock` boolean with a runtime check (e.g., detect if Wails runtime is available via `window.__wails__` or similar). Alternatively, make this a build-time env flag (`import.meta.env.VITE_USE_MOCK`). Remove the `MOCK_TASKS` and `MOCK_PROJECTS` arrays after the switch.

### Auth Store Silent Fallback to Mock User
- Issue: `useAuthStore.signIn()` catches any error from the Wails `SignIn` binding and silently falls back to a hardcoded mock user (`id: 'mock-user-1'`). There is no way for the UI to distinguish between "running outside Wails" and "real auth failure."
- Files: `frontend/src/stores/auth.ts` (lines 27–40)
- Impact: Auth failures are invisible. Users may believe they are signed in when they are not actually authenticated against Azure AD. The mock user has no real token, so any ADO API calls would fail.
- Fix approach: Separate the "Wails not available" case from actual auth errors. Use a dedicated `isWailsAvailable()` check. Only fall back to mock in dev mode.

### Duplicated Type Definitions (Task, Project)
- Issue: `Task` and `Project` types are manually defined as TypeScript interfaces in the stores (`frontend/src/stores/tasks.ts` lines 4–23, `frontend/src/stores/projects.ts` lines 4–12) AND auto-generated as classes in Wails bindings (`frontend/bindings/dev.azure.com/xbox/xb-tasks/domain/models.ts` lines 12, 59). The store interfaces use `string` for date fields while bindings use `time$0.Time`.
- Files: `frontend/src/stores/tasks.ts`, `frontend/src/stores/projects.ts`, `frontend/bindings/dev.azure.com/xbox/xb-tasks/domain/models.ts`
- Impact: If Go domain types change, the manually-defined store interfaces will silently drift from the binding-generated types. Date handling will differ (plain strings vs. Wails `Time` wrapper). The `as Task` casts in stores (e.g., line 162 of `tasks.ts`) mask type mismatches.
- Fix approach: Import and re-export the generated binding types in stores instead of redefining them. Map date fields explicitly if the Wails `Time` type needs conversion.

### Hardcoded Mock Data in View Components
- Issue: Multiple view components contain their own inline mock data that is completely independent of stores or the backend.
- Files:
  - `frontend/src/views/DashboardView.vue` (lines 26–49): Hardcoded `mockPRs` array for PR section
  - `frontend/src/components/TaskDetail.vue` (lines 60, 122–128, 153–179, 204–211, 234–241): Hardcoded mock projects, subtasks, PRs, comments, and timeline events
  - `frontend/src/layouts/AppShell.vue` (lines 45–51): Hardcoded `activityEvents` array
- Impact: These sections display static fake data to the user. Subtasks, comments, PRs, and activity feeds appear functional but are completely non-persistent. User might add a comment and lose it on navigation.
- Fix approach: Build stores/services for PRs, subtasks, comments. The Go backend already has DB tables for `pull_requests` and `task_deps` but no corresponding frontend stores. Wire up the existing schema.

### Playground Views Are Dead Weight (1,935 lines)
- Issue: Four playground views exist for UI prototyping. Two are stub placeholders, two are massive (1,226 and 725 lines) containing full duplicate UI implementations with their own hardcoded mock data.
- Files:
  - `frontend/src/views/PlaygroundDetailView.vue` (1,226 lines)
  - `frontend/src/views/PlaygroundTasksView.vue` (725 lines)
  - `frontend/src/views/PlaygroundDashboardView.vue` (11 lines — stub)
  - `frontend/src/views/PlaygroundAdoView.vue` (13 lines — stub)
- Impact: ~1,935 lines of code that ship in the router. Adds bundle size and maintenance confusion. The playground views have their own component logic that can diverge from the real views.
- Fix approach: Extract any useful design patterns into the real views, then delete playground files. Remove their routes from `frontend/src/router/index.ts` (lines 31–49). The sidebar already gates playground nav items behind `import.meta.env.DEV` (good), but the routes are always registered.

### ADO Package Is Types-Only Skeleton
- Issue: `pkg/ado/` contains only a `domain/types.go` file with type definitions. No ADO API client, no sync logic, no work item fetching.
- Files: `pkg/ado/domain/types.go`
- Impact: The ADO integration promised by the app name and UI doesn't exist yet. DB tables for `ado_work_items`, `pull_requests`, and `task_ado_links` are created but never populated.
- Fix approach: Implement an ADO REST API client in `pkg/ado/` using the PAT or OAuth token from the keyring. Build a sync service that runs on the configured interval (`sync.interval_minutes` from config).

## Known Bugs

### PAT Sign-In Does Not Use PAT Token
- Symptoms: `signInWithPAT()` in `LoginView.vue` collects the PAT input but calls `authStore.signIn()` (the OAuth flow) instead of passing the token to `SignInWithPAT(pat)`. The PAT is collected but discarded.
- Files: `frontend/src/views/LoginView.vue` (lines 22–29)
- Trigger: Click "Sign In" after entering a PAT on the login page.
- Workaround: None — PAT auth is broken from the UI. The Go backend `SignInWithPAT` method works correctly but is never called.
- Fix: Change line 25 from `authStore.signIn()` to call a new `authStore.signInWithPAT(token)` method that invokes the `SignInWithPAT` Wails binding.

### signInWithPAT Missing `await`
- Symptoms: Even if the PAT call were correct, `signInWithPAT()` is `async` but does not `await` the `authStore.signIn()` call on line 25. The `isAuthenticated` check on line 26 will always evaluate before sign-in completes.
- Files: `frontend/src/views/LoginView.vue` (line 25)
- Trigger: Any use of PAT sign-in flow.
- Fix: Add `await` before the auth store call.

### Settings Save Sends Multiple Sequential Requests
- Symptoms: `saveConfig()` in SettingsView calls `Set()` 8 times sequentially (one per config key). Each call writes the entire Viper config to disk. This is slow and risks partial writes if the app crashes mid-save.
- Files: `frontend/src/views/SettingsView.vue` (lines 61–77)
- Trigger: Click "Save" in settings.
- Fix: Add a `SetAll()` or `SetBatch()` method to `ConfigService` that accepts all values and writes once.

## Security Considerations

### OAuth State Parameter Not Validated
- Risk: The OAuth PKCE flow uses a hardcoded `"state"` string (line 79 of `auth.go`) instead of a cryptographically random value. The callback handler does not validate the state parameter at all.
- Files: `internal/auth/auth.go` (line 79, lines 88–97)
- Current mitigation: The PKCE flow itself (code_verifier/code_challenge) provides some CSRF protection. The callback server runs on localhost with a random port, limiting the attack surface.
- Recommendations: Generate a random state value per auth attempt. Validate it in the callback handler before exchanging the code.

### Placeholder OAuth Client IDs in Source
- Risk: Default constants `YOUR_TENANT_ID` and `YOUR_CLIENT_ID` are hardcoded in `auth.go`. If env vars are not set, the app will attempt OAuth with these invalid values, which will fail but may leak the placeholder to Microsoft's auth endpoint in the URL.
- Files: `internal/auth/auth.go` (lines 27–28)
- Current mitigation: Auth will fail with invalid credentials. Env vars override the defaults.
- Recommendations: Detect placeholder values at startup and show a clear configuration error instead of attempting OAuth. Add a config validation step.

### http.DefaultClient Has No Timeout
- Risk: `fetchUserProfile()` uses `http.DefaultClient` which has no timeout. A hung Microsoft Graph API call would block the goroutine indefinitely.
- Files: `internal/auth/auth.go` (line 291)
- Current mitigation: None.
- Recommendations: Create a custom `http.Client` with a 10–15 second timeout for all external API calls.

### No Input Validation on Frontend-to-Backend Boundary
- Risk: Wails-bound service methods accept raw strings from the frontend without validation beyond empty checks. The `status` and `priority` fields have CHECK constraints in SQLite but no Go-level validation — errors are returned as generic SQL errors.
- Files: `internal/app/tasks.go` (all methods), `internal/app/projects.go` (all methods)
- Current mitigation: SQLite CHECK constraints catch invalid enum values. Parameterized queries prevent SQL injection.
- Recommendations: Add Go-level validation for status/priority enums with clear error messages. Validate string lengths for title/description.

## Performance Bottlenecks

### DependencyService Circular Check Uses N+1 Queries
- Problem: `hasCircularDep()` performs DFS by issuing a separate SQL query for each node visited in the dependency graph.
- Files: `internal/app/deps.go` (lines 55–93)
- Cause: Each iteration of the DFS loop executes `SELECT depends_on FROM task_deps WHERE task_id = ?` for the current node.
- Improvement path: Load the entire `task_deps` table into memory (it will be small for a personal tool) and perform the DFS in Go. Alternatively, use a recursive CTE: `WITH RECURSIVE deps AS (...)`.

### Repeated Scan Boilerplate Across All Go Services
- Problem: The same 18-field `Scan()` call for `domain.Task` is repeated 6 times across `tasks.go` and `deps.go`. Any schema change requires updating all 6 locations.
- Files: `internal/app/tasks.go` (lines 50–56, 84–90, 147–153, 225–231), `internal/app/deps.go` (lines 153–159)
- Cause: No shared scan helper was used until `deps.go` introduced `scanTasks()` (line 146), but `tasks.go` still has inline scans.
- Improvement path: Use the `scanTasks()` helper from `deps.go` everywhere, or adopt a lightweight scan library. Move `scanTasks` to a shared location.

## Fragile Areas

### TaskDetail.vue (658 lines, Heavily Mocked)
- Files: `frontend/src/components/TaskDetail.vue`
- Why fragile: This component owns editing, subtasks, PRs, comments, and timeline — all with independent mock data. It mixes real store operations (`taskStore.updateTask`) with purely local mock state (subtasks, comments). Any refactor to wire up real data will require rewriting most of the component.
- Safe modification: Changes to the task save/status logic are safe (they use the store). Changes to subtasks, comments, PRs, or timeline will require building new stores first.
- Test coverage: Zero tests.

### Store Mock/Real Code Paths
- Files: `frontend/src/stores/tasks.ts`, `frontend/src/stores/projects.ts`
- Why fragile: Every store method has an `if (useMock)` branch creating parallel code paths. The real (non-mock) paths are untested since `useMock` is always `true`. When switching to real bindings, the mock code paths become dead code that clutters the stores.
- Safe modification: Changes to mock behavior are safe but meaningless. Changes to the real code paths can't be verified without running in Wails context.
- Test coverage: Zero tests.

### Database Schema Migration Strategy
- Files: `internal/db/db.go` (lines 39–136)
- Why fragile: The migration is a single `CREATE TABLE IF NOT EXISTS` block. There is no versioning, no ALTER TABLE support, no migration tracking. Adding a column to an existing table requires the user to delete their database.
- Safe modification: Adding new tables is safe. Modifying existing tables is destructive.
- Test coverage: Zero tests.

## Scaling Limits

### SQLite Single-Writer
- Current capacity: Single user, single process — fine for a desktop app.
- Limit: WAL mode allows concurrent reads but only one writer. Not a real concern for this use case.
- Scaling path: Not applicable — this is a personal desktop tool. SQLite is the correct choice.

### No Pagination on Task Lists
- Current capacity: Works fine with <100 tasks.
- Limit: `TaskService.List()` and `ListFiltered()` return ALL matching tasks. With thousands of tasks, this will cause UI lag.
- Files: `internal/app/tasks.go` (lines 63–97, 192–238), `frontend/src/stores/tasks.ts` (line 165)
- Scaling path: Add `LIMIT`/`OFFSET` or cursor-based pagination to Go services. Add infinite scroll or virtual list to the frontend.

## Dependencies at Risk

### Wails v3 Alpha
- Risk: `github.com/wailsapp/wails/v3 v3.0.0-alpha.74` — this is pre-release software. APIs may change between alpha versions. Binding generation, event system, and window management may have breaking changes.
- Impact: Any Wails upgrade could require code changes across `main.go`, all service bindings, and the frontend binding imports.
- Migration plan: Pin the exact version. Monitor Wails v3 changelog before upgrading. Consider wrapping Wails-specific APIs behind thin abstractions.

### Go 1.25 (Unreleased)
- Risk: `go.mod` specifies `go 1.25.0` which is not yet a stable release. This may cause build issues on standard toolchains.
- Files: `go.mod` (line 3)
- Impact: Contributors with standard Go installations cannot build the project.
- Migration plan: Verify if this is intentional (using a Go tip build) or should be `go 1.24`.

## Missing Critical Features

### No Route Guards / Auth Protection
- Problem: The Vue router has no navigation guards. All routes (tasks, projects, settings, dashboard) are accessible without authentication. The login page exists but is never enforced.
- Files: `frontend/src/router/index.ts`
- Blocks: Cannot enforce authentication. Users can access the full app without signing in.

### No Error Feedback to Users
- Problem: All `catch` blocks in stores and views are empty or silently fall back to defaults. The user is never informed when a save fails, a fetch errors, or settings don't persist.
- Files: `frontend/src/stores/tasks.ts` (line 167), `frontend/src/stores/projects.ts` (line 43), `frontend/src/views/SettingsView.vue` (lines 54, 74)
- Blocks: Users cannot distinguish between "data saved" and "save failed silently."

### No ADO Sync Implementation
- Problem: The entire ADO integration layer is missing. DB tables exist (`ado_work_items`, `pull_requests`, `task_ado_links`) but there is no sync service, no API client, no import/export logic.
- Files: `pkg/ado/domain/types.go` (types only), `internal/db/db.go` (schema lines 99–119)
- Blocks: The core value proposition — unifying ADO work items with local tasks — is not functional.

## Test Coverage Gaps

### No Tests Exist Anywhere
- What's not tested: The entire codebase — zero test files found. No `*_test.go` files, no `*.test.ts` files, no `*.spec.ts` files.
- Files: All of `internal/`, `frontend/src/`, `pkg/`
- Risk: Any refactor (especially wiring stores to real bindings, adding ADO sync, or changing the DB schema) has no safety net. Regressions will only be caught by manual testing.
- Priority: **High** — Start with Go service tests (`TaskService`, `ProjectService`, `DependencyService`) since they have clean interfaces and a testable SQLite backend. Then add Vue component tests for stores.

---

*Concerns audit: 2025-07-15*
