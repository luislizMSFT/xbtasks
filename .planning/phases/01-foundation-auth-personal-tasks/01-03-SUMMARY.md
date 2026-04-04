---
phase: 01-foundation-auth-personal-tasks
plan: 03
subsystem: ui
tags: [wails, vue-router, pinia, tailwind, system-tray, css-custom-properties, theme]

# Dependency graph
requires:
  - phase: 01-01
    provides: "SQLite database, TaskService, ProjectService"
  - phase: 01-02
    provides: "AuthService with Entra ID OAuth2 PKCE, DependencyService"
provides:
  - "Wails app with all Go services registered (auth, task, project, dep, ado, pr, pipeline, config)"
  - "System tray with show/hide and quit menu"
  - "Window close minimizes to tray (doesn't quit)"
  - "Silent auth restore on startup via TryRestoreSession"
  - "Vue app shell with router (hash history), Pinia stores, theme system"
  - "Semantic CSS custom properties for status and priority colors (UI-SPEC contract)"
  - "Full dark/light theme with system preference detection"
affects: [01-04, 01-05, 01-06, 02-ado-integration]

# Tech tracking
tech-stack:
  added: [vue-router@4, pinia@3, "@headlessui/vue@1.7", lucide-vue-next@1, "@vueuse/core@14", reka-ui, shadcn-vue, tw-animate-css]
  patterns: [wails-service-registration, hash-history-routing, pinia-composition-stores, css-custom-properties-theme, system-tray-minimize-on-close]

key-files:
  created: []
  modified: [main.go, frontend/src/style.css]

key-decisions:
  - "Accent color uses hex values from UI-SPEC (blue-600 light, blue-500 dark) rather than oklch"
  - "Status/priority colors exposed as CSS custom properties for semantic usage across components"

patterns-established:
  - "CSS custom properties for theme: --color-status-* and --color-priority-* tokens"
  - "System tray pattern: close hides window, tray click restores, explicit Quit to exit"
  - "Service registration: create services before app, register via wailsApp.RegisterService"

requirements-completed: [AUTH-01, AUTH-02]

# Metrics
duration: 4min
completed: 2026-04-03
---

# Phase 01 Plan 03: Wails App Shell — Register Services, System Tray, Theme Config Summary

**Wails app with 8 registered Go services, system tray minimize-on-close, Vue shell (router + pinia + theme), and UI-SPEC semantic color tokens**

## Performance

- **Duration:** 4 min
- **Started:** 2026-04-03T23:59:39Z
- **Completed:** 2026-04-04T00:03:43Z
- **Tasks:** 2
- **Files modified:** 1

## Accomplishments
- Verified all Go services registered in main.go (auth, task, project, dep, ado, pr, pipeline, config)
- Verified system tray with show/quit, window close interception, and silent auth restore
- Verified Vue app shell with hash-history router, Pinia state management, and theme composable
- Added semantic status colors (todo, in-progress, in-review, done, blocked, cancelled) as CSS custom properties
- Added priority colors (p0-p3) as CSS custom properties for both light and dark modes
- Added accent-hover, destructive-hover, success, warning, and text-tertiary theme tokens

## Task Commits

Each task was committed atomically:

1. **Task 03.1: Update main.go — register all services, system tray, window close interception** - _already complete_ (all acceptance criteria verified, no changes needed)
2. **Task 03.2: Install frontend dependencies and set up Vue app shell** - `87a3b58` (feat: add semantic status/priority CSS custom properties)

**Plan metadata:** _see final commit_ (docs: complete plan)

## Files Created/Modified
- `frontend/src/style.css` - Added semantic status colors, priority colors, accent-hover, destructive-hover, success, warning, text-tertiary tokens for light and dark modes

## Decisions Made
- Used hex color values from UI-SPEC for accent (blue-600/#2563EB light, blue-500/#3B82F6 dark) replacing the oklch values that didn't match the UI-SPEC contract
- Exposed status and priority colors as `--color-status-*` and `--color-priority-*` CSS custom properties so components can reference them semantically

## Deviations from Plan

### Pre-existing Work

The vast majority of this plan's work was already implemented in prior development iterations before this plan was created. The project has evolved significantly beyond Phase 1 scope:

- **main.go**: Already had all 8 services registered (plan specified 4), system tray, window close interception, and TryRestoreSession — all verified passing acceptance criteria
- **Frontend**: Already had vue-router@4, pinia@3, @headlessui/vue, lucide-vue-next, @vueuse/core installed, plus additional packages (reka-ui, shadcn-vue, @fluentui/svg-icons)
- **Router**: Already configured with hash history and more routes than plan specified (7 routes vs plan's 4)
- **Auth store**: Already more complete than plan (includes Wails binding integration, mock fallback, tryRestore)
- **App.vue**: Already uses AppShell layout with useTheme composable (more advanced than plan's useColorMode)
- **style.css**: Already had comprehensive shadcn-vue theme system with dark/light modes

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Added semantic status/priority CSS custom properties**
- **Found during:** Task 03.2 verification
- **Issue:** UI-SPEC specifies status colors (todo, in-progress, in-review, done, blocked, cancelled) and priority colors (p0-p3) as CSS custom properties, but these were missing from style.css
- **Fix:** Added all status and priority color variables for both light and dark modes, plus text-tertiary, accent-hover, destructive-hover, success, and warning tokens
- **Files modified:** frontend/src/style.css
- **Verification:** `npx vite build` succeeds, grep confirms presence of `--color-status-todo`, `--color-priority-p0`, etc.
- **Committed in:** 87a3b58

---

**Total deviations:** 1 auto-fixed (Rule 2 — missing critical)
**Impact on plan:** Added semantic color tokens that were specified in UI-SPEC but not yet implemented. No scope creep.

## Issues Encountered
- `go build ./...` fails due to iOS build target having no main function — resolved by using `go build .` for current platform only (pre-existing issue, not caused by this plan)

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- All Go services wired up and building
- Frontend shell complete with routing, state management, and full theme system
- Ready for Phase 01 Plans 04-06 (task CRUD views, auth flow UI, dashboard)
- Semantic color tokens available for status badges and priority indicators in upcoming views

---
*Phase: 01-foundation-auth-personal-tasks*
*Completed: 2026-04-03*
