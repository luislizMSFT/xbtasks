---
phase: 02-ado-integration-prs-unified-dashboard
plan: 09
subsystem: ui
tags: [vue, wails, lucide, shadcn-vue, external-links, comments, ado-sync]

# Dependency graph
requires:
  - phase: 02-03
    provides: "ADO linking service (LinkService, personal-to-public model)"
  - phase: 02-04
    provides: "ExternalLinksService and CommentService Go backends"
provides:
  - "ExternalLinks.vue — type-detected link display with add/delete/open-in-browser"
  - "CommentsSection.vue — private/public comments with Push to ADO"
  - "TaskDetail integration — links, comments, and sync push button in sidebar"
affects: [02-10, frontend-views, task-detail]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Type-detected icons for external links (ICM→AlertCircle, Grafana→BarChart3, etc.)"
    - "Private-by-default comments with selective ADO push"
    - "Dynamic Wails binding imports in Vue components (same pattern as stores)"

key-files:
  created:
    - "frontend/src/components/ExternalLinks.vue"
    - "frontend/src/components/CommentsSection.vue"
  modified:
    - "frontend/src/components/TaskDetail.vue"

key-decisions:
  - "ExternalLinks uses window.open for real browser opening (UX-02 requirement)"
  - "CommentsSection shows Push to ADO only for private comments on public tasks"
  - "SyncConfirmDialog mounted inside TaskDetail for outbound push preview"

patterns-established:
  - "Component-level dynamic imports for Wails bindings (externallinksservice, commentservice)"
  - "Type-color mapping pattern for link icons with Lucide components"

requirements-completed: [LINK-03, CMT-01, CMT-02, CMT-03, UX-02]

# Metrics
duration: 3min
completed: 2026-04-06
---

# Phase 02 Plan 09: External Links & Comments Sections Summary

**ExternalLinks with auto-detected type icons (ICM/Grafana/ADO/Wiki) and CommentsSection with private/public badges plus selective ADO push, integrated into TaskDetail sidebar**

## Performance

- **Duration:** 3 min
- **Started:** 2026-04-06T19:16:42Z
- **Completed:** 2026-04-06T19:19:54Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments
- ExternalLinks component with type-detected icons (ICM red, Grafana green, ADO blue, Wiki purple), add/delete, and real browser opening
- CommentsSection component with private/public badges (Globe/Lock icons), timestamp formatting, and Push to ADO button for linked tasks
- TaskDetail sidebar extended with ExternalLinks, CommentsSection, sync Push to ADO button, and SyncConfirmDialog

## Task Commits

Each task was committed atomically:

1. **Task 1: ExternalLinks + CommentsSection Components** - `3d0ed24` (feat)
2. **Task 2: Integrate Links + Comments into TaskDetail** - `0a3f1eb` (feat)

## Files Created/Modified
- `frontend/src/components/ExternalLinks.vue` - External links section with type icons, add/delete, browser opening
- `frontend/src/components/CommentsSection.vue` - Comments with private/public badges, Push to ADO, relative time formatting
- `frontend/src/components/TaskDetail.vue` - Integrated ExternalLinks, CommentsSection, Push to ADO button, SyncConfirmDialog

## Decisions Made
- Used `window.open(url, '_blank')` for external link opening (UX-02: real browser, not in-app webview)
- Push to ADO button only shown when `taskStore.isPublic(task.id)` is true (linked tasks only)
- SyncConfirmDialog placed inside TaskDetail aside element (Dialog, so DOM position flexible)
- Comments private by default — Push to ADO is explicit opt-in per the project's outbound-never-auto-push decision

## Deviations from Plan

None - plan executed exactly as written. `isPublic` and `publicTaskIds` already existed in the task store (from earlier plan execution), so no deviation was needed.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- ExternalLinks and CommentsSection are ready for backend wiring (services from Plan 04)
- TaskDetail now has full sidebar integration: subtasks, PRs, description, external links, comments, sync push
- Plan 10 (ADO browser/management view) can proceed independently

---
*Phase: 02-ado-integration-prs-unified-dashboard*
*Completed: 2026-04-06*

## Self-Check: PASSED

All files verified present. All commits verified in git log.
