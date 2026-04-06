---
phase: 02-ado-integration-prs-unified-dashboard
plan: 10
subsystem: frontend-projects
tags: [projects, card-grid, ado-linking, dual-progress, pin-star]
dependency_graph:
  requires: [02-03, 02-04]
  provides: [ProjectCard, ProjectDetailView, extended-project-store]
  affects: [frontend/src/views/ProjectsView.vue, frontend/src/router/index.ts]
tech_stack:
  added: []
  patterns: [card-grid-layout, dual-progress-bars, pin-star-favorites, ado-badge-linked-hollow]
key_files:
  created:
    - frontend/src/components/ProjectCard.vue
    - frontend/src/views/ProjectDetailView.vue
  modified:
    - frontend/src/stores/projects.ts
    - frontend/src/views/ProjectsView.vue
    - frontend/src/router/index.ts
decisions:
  - Card grid layout replaces previous list+detail split panel for projects page
  - Dual progress bars kept separate (local primary color, ADO blue) per PROJ-06
  - Link dialog is inline card, not modal, for lightweight linking flow
  - ProjectDetailView groups tasks by status with collapsible sections
metrics:
  duration: 4m
  completed: 2026-04-06
  tasks_completed: 2
  tasks_total: 2
requirements_completed: [PROJ-01, PROJ-02, PROJ-03, PROJ-04, PROJ-05, PROJ-06, PROJ-07, UX-02]
---

# Phase 02 Plan 10: Projects Page & Project Detail Summary

Card grid layout for projects with ADO badge, dual progress bars, pin/star favorites, and project detail dashboard with stats, ADO context, and filtered tasks.

## What Was Done

### Task 1: ProjectCard + ProjectsView Card Grid (c5559ee)

**Extended `projects.ts` store:**
- Added `isPinned` to Project interface
- Added `ProjectADOLink` and `ProjectProgress` interfaces
- Added `projectLinks` and `projectProgress` reactive Maps
- Added methods: `pinProject`, `linkProjectToADO`, `unlinkProject`, `fetchProjectLink`, `fetchProjectProgress`
- Added computed: `pinnedProjects`, `unpinnedProjects`
- Added helper: `isLinked()`

**Created `ProjectCard.vue`:**
- Card component with name, Star pin button (amber when pinned), ADO badge (filled blue circle when linked, dashed hollow when local-only)
- Dual progress bars: Local (primary color) + ADO (blue, only shown when linked)
- Task count display from progress data
- Emits `click` and `pin` events

**Overhauled `ProjectsView.vue`:**
- Replaced previous list+detail split panel with responsive card grid (1/2/3 columns)
- Pinned projects section at top with "Pinned" label
- All Projects (unpinned) section below
- Inline create form preserved
- On mount: fetches projects then loads link + progress for each
- Router navigation to `/projects/:id` on card click

### Task 2: ProjectDetailView + Router Update (e0ca654)

**Created `ProjectDetailView.vue`:**
- Back button navigating to `/projects`
- Header with project name, description, link/unlink buttons
- Inline link dialog (not modal) for entering ADO work item ID
- Stats cards row: Total Tasks, Completed (green), ADO Children (blue, if linked), ADO Progress % (blue, if linked)
- Dual progress bars: Local Tasks (primary) and ADO Items (blue), separate per PROJ-06
- ADO Context card showing linked item ID and direction
- Task list filtered to project, grouped by status (blocked, in progress, in review, todo, done) with checkboxes and priority badges
- Loads project data on mount and watches route param changes

**Updated `router/index.ts`:**
- Added `/projects/:id` route with `project-detail` name pointing to `ProjectDetailView.vue`

## Deviations from Plan

None — plan executed exactly as written.

## Requirements Completed

| Requirement | Description | Status |
|-------------|-------------|--------|
| PROJ-01 | Card grid layout | Done |
| PROJ-02 | Pin/star favorites to top | Done |
| PROJ-03 | Projects flat, no sub-projects | Done |
| PROJ-04 | Local-only or ADO-linked | Done |
| PROJ-05 | Project dashboard with stats and filtered tasks | Done |
| PROJ-06 | Dual progress: local + ADO, kept separate | Done |
| PROJ-07 | Link/unlink project to ADO | Done |
| UX-02 | Consistent UI patterns | Done |

## Commits

| Task | Commit | Description |
|------|--------|-------------|
| 1 | c5559ee | ProjectCard card grid with ADO badge, dual progress, pin/star |
| 2 | e0ca654 | ProjectDetailView with stats, dual progress, ADO context, filtered tasks |
