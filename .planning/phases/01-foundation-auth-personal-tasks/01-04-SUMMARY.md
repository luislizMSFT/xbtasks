---
phase: 01-foundation-auth-personal-tasks
plan: 04
subsystem: ui
tags: [vue, appshell, sidebar, command-palette, accessibility, css-custom-properties, wails]

# Dependency graph
requires:
  - phase: 01-foundation-auth-personal-tasks
    plan: 03
    provides: "Frontend scaffold with Wails v3 + Vue 3, router, auth store, theme composable, shadcn-vue components"
provides:
  - "AppShell layout with 56px icon sidebar and content area"
  - "Sidebar with navigation icons (Dashboard, Tasks, Projects, ADO, Dependencies) and active state"
  - "Accessible sidebar with aria-labels on all interactive elements"
  - "User avatar dropdown with sign-out action"
  - "Command palette (⌘K/Ctrl+K) with navigation, task creation, task search, and sign-out"
  - "Auth-gated App.vue showing LoginView when unauthenticated"
  - "CSS custom property aliases (--color-bg-primary/secondary/tertiary) for theme-aware backgrounds"
affects: [01-05-PLAN, 01-06-PLAN, ui-views, auth-flow]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Auth-gated shell: App.vue shows LoginView vs AppShell based on authStore.isAuthenticated"
    - "CSS custom property theming: --color-bg-primary/secondary/tertiary aliases for surface vars"
    - "Accessible nav buttons: every interactive sidebar element has aria-label"
    - "DropdownMenu for user actions (sign-out) at sidebar bottom"

key-files:
  created: []
  modified:
    - "frontend/src/App.vue"
    - "frontend/src/layouts/AppShell.vue"
    - "frontend/src/components/Sidebar.vue"
    - "frontend/src/components/CommandPalette.vue"
    - "frontend/src/style.css"

key-decisions:
  - "Kept AppShell in layouts/ directory (Vue convention) rather than components/ as plan specified"
  - "Used shadcn-vue Command (cmdk-vue) for command palette instead of raw Headless UI Dialog+Combobox — superior UX"
  - "Auth guard in App.vue with tryRestore() on mount for session persistence"
  - "Added --color-bg-* as aliases to existing --surface-* vars for forward compatibility"

patterns-established:
  - "Auth-gated routing: App.vue checks authStore.isAuthenticated to show LoginView vs AppShell"
  - "Sidebar aria-labels: every nav button, toggle, and action has descriptive aria-label"
  - "DropdownMenu for avatar: sign-out and user info accessible from sidebar bottom"

requirements-completed: [TASK-01]

# Metrics
duration: 7min
completed: 2026-04-04
---

# Phase 01 Plan 04: Frontend — App Shell Layout, Sidebar, Command Palette Summary

**Auth-gated AppShell with 56px icon sidebar, accessible nav buttons, user dropdown, and ⌘K command palette using shadcn-vue**

## Performance

- **Duration:** 7 min
- **Started:** 2026-04-04T00:06:55Z
- **Completed:** 2026-04-04T00:13:07Z
- **Tasks:** 3
- **Files modified:** 5

## Accomplishments

- App.vue now auth-gates the entire shell: shows LoginView when not authenticated, AppShell with router-view when authenticated
- Sidebar navigation icons all have aria-labels for accessibility compliance
- User avatar at sidebar bottom opens DropdownMenu with display name, email, and Sign Out action
- Command palette sign-out action properly calls authStore.signOut() before navigating
- CSS custom property aliases added for theme-aware backgrounds (--color-bg-primary/secondary/tertiary)
- Frontend builds cleanly (2548 modules, 6.67s)

## Task Commits

Each task was committed atomically:

1. **Task 04.1: Create AppShell layout component** — `551e908` (feat)
2. **Task 04.2: Create Sidebar navigation component** — `19e3b10` (feat)
3. **Task 04.3: Create Command Palette component** — `303ca36` (feat)

## Files Created/Modified

- `frontend/src/App.vue` — Auth-gated rendering: LoginView vs AppShell based on isAuthenticated
- `frontend/src/layouts/AppShell.vue` — Added CSS custom property backgrounds, sidebar width comment
- `frontend/src/components/Sidebar.vue` — Added aria-labels, DropdownMenu sign-out, LogOut icon
- `frontend/src/components/CommandPalette.vue` — Sign-out uses authStore.signOut(), imported useAuthStore
- `frontend/src/style.css` — Added --color-bg-primary/secondary/tertiary aliases

## Decisions Made

1. **Kept AppShell at layouts/AppShell.vue** — Plan specified components/AppShell.vue but layouts/ is the standard Vue convention for layout components and was already established in the codebase.
2. **Used shadcn-vue Command instead of Headless UI Dialog+Combobox** — The cmdk-vue based Command component was already integrated and provides superior search/filter UX with keyboard navigation, grouping, and empty state handling.
3. **Auth guard via App.vue instead of route meta** — Plan's approach (v-if="!authStore.isAuthenticated") is more secure than the previous hideShell route meta approach, which allowed direct navigation to protected routes.
4. **Added tryRestore() on mount** — Session persistence: on app load, attempts to restore existing auth from Wails bindings/keychain before showing login.

## Deviations from Plan

### Structural Deviations

**1. AppShell file location: layouts/ vs components/**
- **Plan specified:** `frontend/src/components/AppShell.vue`
- **Actual:** `frontend/src/layouts/AppShell.vue` (pre-existing location)
- **Rationale:** layouts/ is the established Vue convention. Moving would break existing imports and the layouts/ convention was already established. No functional difference.

**2. Command palette implementation: shadcn-vue Command vs Headless UI Dialog+Combobox**
- **Plan specified:** Headless UI `Dialog` + `Combobox` components
- **Actual:** shadcn-vue `CommandDialog` (based on cmdk-vue)
- **Rationale:** cmdk-vue was already integrated and provides better UX (fuzzy search, grouping, keyboard hints, empty states). Both @headlessui/vue and reka-ui are in dependencies.

**3. CSS variable naming: --color-bg-* vs existing --surface-***
- **Plan specified:** `var(--color-bg-primary)`, `var(--color-bg-secondary)`
- **Actual:** Added as aliases mapping to existing `--surface-primary` / `--surface-secondary` variables
- **Rationale:** The codebase already uses --surface-* naming from the UI-SPEC. Aliases maintain backward compatibility while satisfying the plan's naming convention.

---

**Total deviations:** 3 structural (file location, component library, variable naming)
**Impact on plan:** All deviations maintain or improve upon plan intent. Existing implementations exceed plan requirements.

## Issues Encountered

None — all components were pre-existing in an advanced state. Enhancements were additive (aria-labels, sign-out dropdown, auth guard, CSS variable aliases).

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- App shell layout complete with auth-gated routing
- All views (Dashboard, Tasks, Projects, Settings, ADO, Dependencies) accessible via sidebar navigation
- Command palette provides keyboard-first navigation (⌘K/Ctrl+K)
- Ready for Plan 05 (task views) and Plan 06 (remaining features)
- Frontend builds successfully with 0 errors

---
*Phase: 01-foundation-auth-personal-tasks*
*Completed: 2026-04-04*
