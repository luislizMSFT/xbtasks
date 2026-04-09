---
phase: 10-implement-dashboard-redesign-and-unified-header-bars
plan: 02
subsystem: frontend-dashboard
tags: [dashboard, attention-bar, task-rows, upcoming, blocked-reasons]
dependency_graph:
  requires: [10-01]
  provides: [DashboardTaskRow, redesigned-dashboard]
  affects: [DashboardView, frontend-components]
tech_stack:
  added: []
  patterns: [attention-bar-nudges, richer-task-rows, urgency-computed-properties]
key_files:
  created:
    - frontend/src/components/tasks/DashboardTaskRow.vue
  modified:
    - frontend/src/views/DashboardView.vue
decisions:
  - DashboardTaskRow is a NEW component (not TaskRow modification) with different layout optimized for dashboard context
  - Attention Bar uses inline rgb() styles for alpha-blended backgrounds (reliable across Tailwind v4 configurations)
  - upcomingTasks includes tasks up to 7 days overdue for visibility
  - mergeReadyPRs threshold is 2+ approvals (vote >= 5)
  - Blocked section uses priority dot + blockedReason (not full DashboardTaskRow) for compact blocked display
metrics:
  duration: 4m
  completed: 2026-04-09
requirements:
  - P10-DASH-01
---

# Phase 10 Plan 02: Dashboard Redesign — Attention Bar, Upcoming, Richer Rows Summary

Redesigned DashboardView with Attention Bar (urgency nudges for due-soon/pipeline-failure/merge-ready), Upcoming section replacing Recent Activity, Blocked section with reason text, and DashboardTaskRow component for richer task context.

## What Was Done

### Task 1: Create DashboardTaskRow component
**Commit:** `6621562`

Created `frontend/src/components/tasks/DashboardTaskRow.vue` — a dedicated dashboard task row with:
- Priority dot (8px) using `priorityBgColor()` from styles.ts
- ADO type icon via `adoTypeIcon()`/`adoTypeColor()` or `SquareCheckBig` fallback for personal tasks
- Task title (truncated, flex-1)
- Personal badge (`text-[8px]`) for non-linked tasks
- Status badge (`text-[9px] h-4`) with outline variant and status-specific colors
- Optional due date display with amber (due-soon) / red (overdue) highlighting
- Click event emission with taskId for parent navigation

### Task 2: Rewrite DashboardView with Attention Bar, Upcoming, Blocked, richer rows
**Commit:** `7b0d844`

Major rewrite of `frontend/src/views/DashboardView.vue`:

**Imports updated:**
- Removed: `TaskRow`, `RefreshCw` (sync moved to SyncCluster in plan 10-01)
- Added: `DashboardTaskRow`, `CalendarDays`, `GitMerge`, `parseReviewers`, `priorityBgColor`

**New computed properties:**
- `upcomingTasks` — tasks due within 3 days (or up to 7 days overdue), sorted by due date
- `dueSoonTasks` — tasks due within 3 days (Attention Bar amber nudge)
- `failedPipeline` — first failed pipeline (Attention Bar red nudge)
- `mergeReadyPRs` — active PRs with 2+ reviewer approvals (Attention Bar green nudge)
- `showAttentionBar` — master toggle for Attention Bar visibility
- `selectTask()` — navigation handler (replaces inline lambdas)

**Template changes:**
- Teleport: stat badges (`N active`, `N blocked`, `done/total done`) replace old ratio/dots/sync-button
- Greeting: downsized to `text-lg font-semibold` (was `text-xl`)
- Attention Bar: conditional urgency nudges (due-soon amber, pipeline-failure red, merge-ready emerald)
- Today's Focus: `DashboardTaskRow` replaces `TaskRow` with `isPersonal` prop
- Upcoming: replaces "Recent Activity" with due-date-aware task list using `DashboardTaskRow`
- Blocked: enhanced with priority dot + `blockedReason` text (italic, red-tinted)
- Right column (PRs + Pipelines): preserved unchanged

**Removed:**
- `recentTasks` computed property
- `TaskRow` usage in all sections

## Deviations from Plan

None — plan executed exactly as written.

## Verification Results

- TypeScript: `vue-tsc --noEmit` passes with zero errors
- All acceptance criteria validated via grep checks (20/20 pass)
- No `TaskRow` or `RefreshCw` references remain in DashboardView.vue
- No "Recent Activity" section heading in DashboardView.vue

## Self-Check: PASSED
