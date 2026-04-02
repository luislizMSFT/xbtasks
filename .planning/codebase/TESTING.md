# Testing Patterns

**Analysis Date:** 2025-06-26

## Test Framework

**Status: No tests exist in this codebase.**

No test files were found anywhere in the project:
- No `*_test.go` files in any Go package
- No `*.test.ts`, `*.spec.ts`, `*.test.vue`, or `*.spec.vue` files in the frontend
- No test configuration files (`jest.config.*`, `vitest.config.*`, `*.test.js`)
- No test runner dependencies in `frontend/package.json`
- No test commands in `frontend/package.json` scripts
- No test tasks in `Taskfile.yml`

## Current Test Infrastructure

**Go:**
- Go's built-in `testing` package is available (no additional dependencies needed)
- No test helper utilities, fixtures, or mocks exist
- The `internal/db/db.go` uses SQLite with `database/sql`, which is testable with in-memory databases (`":memory:"`)

**Frontend (Vue/TypeScript):**
- No test runner installed (no vitest, jest, or @vue/test-utils in `package.json`)
- No test configuration files

## Recommended Setup

### Go Tests

**Runner:** Go built-in `testing` package

**Test File Naming:** `*_test.go` co-located with source files

**Recommended Structure:**
```
internal/
├── app/
│   ├── tasks.go
│   ├── tasks_test.go       # Unit tests for TaskService
│   ├── projects.go
│   ├── projects_test.go    # Unit tests for ProjectService
│   ├── deps.go
│   └── deps_test.go        # Unit tests for DependencyService
├── auth/
│   ├── auth.go
│   └── auth_test.go        # Unit tests for AuthService
├── config/
│   ├── config.go
│   ├── config_test.go
│   ├── service.go
│   └── service_test.go
└── db/
    ├── db.go
    └── db_test.go           # Schema migration tests
```

**Run Commands (when tests exist):**
```bash
go test ./...               # Run all tests
go test ./internal/app/...  # Run app service tests
go test -v ./...            # Verbose output
go test -cover ./...        # With coverage
go test -race ./...         # With race detection
```

**Test Database Pattern:**
The `internal/db` package opens SQLite with a file path. For testing, use in-memory:
```go
func testDB(t *testing.T) *db.DB {
    t.Helper()
    database, err := db.Open(":memory:")
    if err != nil {
        t.Fatalf("open test db: %v", err)
    }
    t.Cleanup(func() { database.Close() })
    return database
}
```

**Priority Test Targets:**
1. `internal/app/deps.go` — `DependencyService` has cycle detection logic (`hasCircularDep`) that needs unit tests
2. `internal/app/tasks.go` — `TaskService` CRUD operations, especially `ListFiltered` with multiple filter combinations
3. `internal/db/db.go` — Schema migration and constraint validation
4. `internal/auth/auth.go` — Token management (though OAuth flow requires mocking HTTP)

### Frontend Tests

**Recommended Runner:** Vitest (aligns with Vite build system already in use)

**Recommended Structure:**
```
frontend/src/
├── stores/
│   ├── tasks.ts
│   ├── tasks.test.ts        # Store logic tests
│   ├── projects.ts
│   └── projects.test.ts
├── components/
│   ├── TaskRow.vue
│   ├── TaskRow.test.ts      # Component tests
│   └── ...
└── lib/
    ├── utils.ts
    └── utils.test.ts         # Utility function tests
```

**Install Commands (when ready):**
```bash
cd frontend
npm install -D vitest @vue/test-utils happy-dom
```

**Priority Test Targets:**
1. Pinia stores — `useTaskStore` computed properties (`filteredTasks`, `grouped`, `stats`) have pure logic testable without DOM
2. `cn()` utility in `frontend/src/lib/utils.ts`
3. Component rendering for `StatusBadge.vue`, `PriorityBadge.vue` — computed switch logic

## Mock Data as Test Fixtures

The codebase already contains extensive mock data that can serve as test fixtures:

**`frontend/src/stores/tasks.ts`** — `MOCK_TASKS` array (10 tasks covering all statuses: `todo`, `in_progress`, `in_review`, `done`, `blocked`, `cancelled`)

**`frontend/src/stores/projects.ts`** — `MOCK_PROJECTS` array (2 projects)

**`frontend/src/views/DashboardView.vue`** — `mockPRs` array (4 PRs with various statuses)

**`frontend/src/components/TaskDetail.vue`** — Mock subtasks, mock pull requests, mock related ADO items, mock activity events

These mock datasets encode domain assumptions (valid statuses, priority levels, tag formats) and should be extracted to shared fixtures when building tests.

## Test Coverage

**Current Coverage:** 0% — no tests exist

**Recommended Minimum Targets:**
- Go services: 70%+ (CRUD operations, validation, cycle detection)
- Frontend stores: 60%+ (computed properties, action logic)
- Frontend components: 30%+ (critical rendering paths)

## What to Test First

**High-value, low-effort targets:**

1. **`DependencyService.hasCircularDep()`** in `internal/app/deps.go` — Pure graph algorithm, critical correctness requirement, easy to test with in-memory DB
2. **`TaskService.ListFiltered()`** in `internal/app/tasks.go` — Multiple filter combinations, SQL query building logic
3. **`TaskService.Create()` validation** in `internal/app/tasks.go` — Empty title rejection, default priority
4. **Store computed properties** in `frontend/src/stores/tasks.ts` — `filteredTasks`, `grouped`, `stats` are pure functions of state
5. **DB schema migration** in `internal/db/db.go` — Verify all tables and indexes are created, constraints work

## Test Gaps and Risks

**Untested Areas with Highest Risk:**

| Area | Files | Risk | Priority |
|------|-------|------|----------|
| Cycle detection in dependencies | `internal/app/deps.go` | Infinite loops or missed cycles | **High** |
| SQL query building (ListFiltered) | `internal/app/tasks.go` | SQL injection if filter logic changes | **High** |
| Auth token refresh flow | `internal/auth/auth.go` | Silent auth failures, token leaks | **Medium** |
| Schema migration idempotency | `internal/db/db.go` | Data loss on re-migration | **Medium** |
| Store mock/real switching | `frontend/src/stores/tasks.ts` | Mock data served in production | **Medium** |
| Config persistence | `internal/config/config.go` | Lost settings on upgrade | **Low** |

---

*Testing analysis: 2025-06-26*
