---
status: diagnosed
phase: 10-implement-dashboard-redesign-and-unified-header-bars
source: [10-01-SUMMARY.md, 10-02-SUMMARY.md]
started: 2026-04-09T00:30:54Z
updated: 2026-04-09T00:52:00Z
---

## Current Test

[testing complete]

## Tests

### 1. SyncCluster in Top Bar
expected: Every page shows the SyncCluster in the top bar left zone — green/red dot, Synced/Offline/Syncing label, relative time since last sync, and a refresh button with aria-label.
result: issue
reported: "synched is done but then prs appear if it's synched I expect everything to load. When I navigate to azure tab it loads again"
severity: major

### 2. SyncCluster Refresh Button
expected: Clicking the refresh button in the SyncCluster triggers a sync. The button shows a spinner while syncing, and the status updates to "Syncing…" then back to "Synced" with updated time.
result: pass

### 3. Dashboard Center Zone Stats
expected: On the Dashboard page, the top bar center zone shows stat badges — "N active", "N blocked", and "done/total done" counts.
result: pass

### 4. Tasks Page Center Zone
expected: On the Tasks page, the top bar center zone shows status filter chips (replacing any old sync button content).
result: pass

### 5. ADO Page Center Zone
expected: On the ADO page, the top bar center zone shows only the tabs (Board/Backlog/etc). Old connection pill, sync button, and loading spinner are gone from the top bar.
result: pass

### 6. Attention Bar Urgency Nudges
expected: On the Dashboard, when relevant data exists, an Attention Bar appears with urgency nudges — due-soon tasks highlighted in amber, failed pipelines in red, merge-ready PRs in green. When none of these conditions apply, the Attention Bar is hidden.
result: issue
reported: "I put a date on a task and didn't show up"
severity: major

### 7. Today's Focus Richer Task Rows
expected: The Dashboard "Today's Focus" section uses richer task rows showing: a priority dot (colored by priority), ADO type icon (or checkmark for personal tasks), task title, personal badge for non-linked tasks, status badge, and optional due date.
result: pass

### 8. Upcoming Section Replaces Recent Activity
expected: The Dashboard shows an "Upcoming" section (not "Recent Activity") listing tasks due within 3 days. Overdue tasks are highlighted in red.
result: pass

### 9. Blocked Section with Reason Text
expected: The Dashboard "Blocked" section shows blocked tasks. When a task has a blockedReason, it displays the reason text in italic, red-tinted styling below the task.
result: skipped
reason: No UI to add blocked reason to a task — can't verify display

### 10. Dashboard Greeting Size
expected: The Dashboard greeting text ("Good morning, Name") appears in a smaller size (text-lg font-semibold) — noticeably smaller than the previous text-xl styling.
result: pass

## Summary

total: 10
passed: 7
issues: 2
pending: 0
skipped: 1

## Gaps

- truth: "Every page shows the SyncCluster in the top bar left zone with green/red dot, Synced/Offline/Syncing label, relative time, and refresh button"
  status: failed
  reason: "User reported: synched is done but then prs appear if it's synched I expect everything to load. When I navigate to azure tab it loads again"
  severity: major
  test: 1
  root_cause: "SyncCluster label only checks syncStore.syncing (backend sync) — has zero awareness of frontend data fetches (PRs, pipelines, work items). adoStore.connected flips true on first API success while other stores still loading. AdoView.vue unconditionally re-fetches all data on mount without emptiness guards."
  artifacts:
    - path: "frontend/src/components/SyncCluster.vue"
      issue: "Label logic only checks syncStore.syncing, not whether all stores finished loading"
    - path: "frontend/src/stores/sync.ts"
      issue: "No concept of 'all data loaded' — only tracks backend sync lifecycle"
    - path: "frontend/src/views/AdoView.vue"
      issue: "onMounted unconditionally re-fetches data already prefetched by App.vue"
  missing:
    - "Add composite isFullyLoaded computed that aggregates loading states from all stores"
    - "SyncCluster should show Syncing until composite signal indicates all data ready"
    - "Guard AdoView onMounted fetches with emptiness checks like DashboardView does"
  debug_session: ".planning/debug/sync-shows-done-but-data-loading.md"

- truth: "Attention Bar shows urgency nudges when relevant data exists — due-soon tasks in amber, failed pipelines in red, merge-ready PRs in green"
  status: failed
  reason: "User reported: I put a date on a task and didn't show up"
  severity: major
  test: 6
  root_cause: "dueSoonTasks computed uses diffDays >= 0 with timestamp-level comparison. new Date('YYYY-MM-DD') creates midnight UTC, so once local time passes midnight UTC, today's tasks get negative diffDays and are excluded. upcomingTasks uses diffDays >= -7 which is why Test 8 passed."
  artifacts:
    - path: "frontend/src/views/DashboardView.vue"
      issue: "dueSoonTasks filter diffDays >= 0 compares at timestamp level instead of date level"
  missing:
    - "Normalize both due date and current date to start-of-day before computing diffDays"
    - "Use existing startOfDay() helper from stores/tasks.ts"
  debug_session: ".planning/debug/attention-bar-due-soon-missing.md"
