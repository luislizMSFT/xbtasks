---
phase: 01-foundation-auth-personal-tasks
plan: 05
subsystem: ui
tags: [pinia, vue, css-custom-properties, wails-bindings, login, badges, task-components]

# Dependency graph
requires:
  - phase: 01-foundation-auth-personal-tasks/plan-03
    provides: "Frontend scaffold with Wails v3 + Vue 3, auth store pattern, style.css with semantic CSS tokens"
provides:
  - "Pinia task store wrapping TaskService + DependencyService Go bindings"
  - "StatusBadge and PriorityBadge using semantic --color-status-*/--color-priority-* CSS tokens"
  - "TagChip with removable capability"
  - "TaskRow and TaskDetail components ready for view consumption"
  - "LoginView with Microsoft SSO + PAT fallback per UI-SPEC copywriting contract"
affects: [01-06-PLAN, 02-ado-integration]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Pinia stores use dynamic import() for Wails bindings with graceful fallback"
    - "Badge components use CSS custom properties via inline style binding for UI-SPEC color contract"
    - "color-mix() for opacity variants on semantic tokens (10% bg, 20% border)"

key-files:
  created: []
  modified:
    - "frontend/src/stores/tasks.ts"
    - "frontend/src/components/ui/StatusBadge.vue"
    - "frontend/src/components/ui/PriorityBadge.vue"
    - "frontend/src/views/LoginView.vue"

key-decisions:
  - "Used CSS custom properties via inline style binding instead of Tailwind utility classes for badge colors — makes UI-SPEC contract explicit and supports dark/light mode via CSS vars"
  - "Used color-mix(in srgb, ...) for background/border opacity variants on semantic tokens — cleaner than Tailwind arbitrary values"
  - "PAT button uses ghost variant per UI-SPEC (was outline)"

patterns-established:
  - "Badge color contract: --color-status-* and --color-priority-* tokens applied via computed style objects"
  - "DependencyService bindings imported from separate Wails binding path"

requirements-completed: [AUTH-01, TASK-01, TASK-02, TASK-03, TASK-04, TASK-05, TASK-06, TASK-07]

# Metrics
duration: 4min
completed: 2026-04-04
---

# Phase 01 Plan 05: Frontend — Pinia Store, UI Primitives & Login View Summary

**Pinia task store with DependencyService bindings, badge components using UI-SPEC semantic color tokens, and LoginView aligned to copywriting contract**

## Performance

- **Duration:** 4 min
- **Started:** 2026-04-04T00:06:28Z
- **Completed:** 2026-04-04T00:10:32Z
- **Tasks:** 3
- **Files modified:** 4

## Accomplishments
- Added 6 missing Pinia store methods (setPersonalPriority, getSubtasks, getDependencies, addDependency, removeDependency, getAllTags) wrapping both TaskService and DependencyService Go bindings
- Migrated StatusBadge and PriorityBadge from hardcoded Tailwind classes to semantic --color-status-*/--color-priority-* CSS custom property tokens via style bindings
- Aligned LoginView to UI-SPEC copywriting contract: tagline, error text, loading states, ghost PAT button

## Task Commits

Each task was committed atomically:

1. **Task 05.1: Create Pinia task store wrapping Go service calls** - `189d00e` (feat)
2. **Task 05.2: UI components use semantic color-status/color-priority tokens** - `8194e1b` (feat)
3. **Task 05.3: LoginView aligned with UI-SPEC** - `27d8e7b` (feat)

## Files Created/Modified
- `frontend/src/stores/tasks.ts` - Added setPersonalPriority, getSubtasks, getDependencies, addDependency, removeDependency, getAllTags methods
- `frontend/src/components/ui/StatusBadge.vue` - Migrated to --color-status-* CSS custom properties via computed style binding
- `frontend/src/components/ui/PriorityBadge.vue` - Migrated to --color-priority-* CSS custom properties via computed style binding
- `frontend/src/views/LoginView.vue` - UI-SPEC copywriting (tagline, error text, loading text), ghost PAT button, disabled state on empty PAT

## Decisions Made
- Used CSS custom properties via inline style binding instead of Tailwind utility classes for badge colors — makes the UI-SPEC color contract explicit and works correctly with dark/light mode via CSS variables
- Used `color-mix(in srgb, var(...) N%, transparent)` for background/border opacity variants — avoids Tailwind arbitrary value complexity
- PAT button changed from `variant="outline"` to `variant="ghost"` per UI-SPEC specification

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] LoginView tagline didn't match UI-SPEC**
- **Found during:** Task 3 (LoginView)
- **Issue:** Tagline was "Tasks · ADO · PRs — one pane of glass" instead of UI-SPEC specified "All your work in one place"
- **Fix:** Updated to match UI-SPEC copywriting contract
- **Files modified:** frontend/src/views/LoginView.vue
- **Committed in:** 27d8e7b

**2. [Rule 1 - Bug] LoginView PAT signIn didn't await result**
- **Found during:** Task 3 (LoginView)
- **Issue:** `signInWithPAT()` called `authStore.signIn()` without `await`, could route before auth completes
- **Fix:** Added `await` and disabled Connect button when PAT input empty
- **Files modified:** frontend/src/views/LoginView.vue
- **Committed in:** 27d8e7b

**3. [Rule 2 - Missing Critical] No loading text feedback**
- **Found during:** Task 3 (LoginView)
- **Issue:** Loading state showed spinner but no text; UI-SPEC specifies "Signing in..."
- **Fix:** Added "Signing in..." text to both sign-in buttons during loading
- **Files modified:** frontend/src/views/LoginView.vue
- **Committed in:** 27d8e7b

---

**Total deviations:** 3 auto-fixed (2 bugs, 1 missing critical)
**Impact on plan:** All fixes align implementation with UI-SPEC contract. No scope creep.

## Issues Encountered
None — existing codebase had all 7 files pre-created from earlier iterations. Plan execution focused on filling gaps and aligning with UI-SPEC contract.

## User Setup Required
None — no external service configuration required.

## Next Phase Readiness
- All UI primitives (StatusBadge, PriorityBadge, TagChip, TaskRow, TaskDetail) ready for consumption by Plan 06 page views
- Pinia task store complete with all CRUD + dependency + tag methods
- LoginView ready for auth integration
- Frontend builds successfully (2548 modules, 6.6s build time)

---
*Phase: 01-foundation-auth-personal-tasks*
*Completed: 2026-04-04*
