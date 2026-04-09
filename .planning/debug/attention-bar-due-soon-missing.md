---
status: diagnosed
trigger: "Attention Bar due-soon nudge doesn't appear when a task has a due date within 3 days"
created: 2026-04-09T01:00:00Z
updated: 2026-04-09T01:00:00Z
---

## Current Focus

hypothesis: dueSoonTasks filter uses `diffDays >= 0` which excludes today's date because `new Date("YYYY-MM-DD")` creates midnight UTC, making diffDays negative once the day has started
test: Compare dueSoonTasks filter (>= 0) with upcomingTasks filter (>= -7) — same data, different results match the user report
expecting: confirmed — Test 8 (Upcoming, uses >= -7) passed, Test 6 (Attention Bar, uses >= 0) failed
next_action: Return root cause diagnosis

## Symptoms

expected: Setting a due date on a task (within 3 days) should make it appear in the Attention Bar as a due-soon nudge in amber
actual: "I put a date on a task and didn't show up" — task appears in Upcoming section (Test 8 passed) but NOT in Attention Bar (Test 6 failed)
errors: None
reproduction: Set a due date on a task to today or within 3 days, navigate to Dashboard, observe Attention Bar is hidden
started: Discovered during UAT Test 6

## Eliminated

(none — first hypothesis confirmed)

## Evidence

- timestamp: 2026-04-09T01:00:00Z
  checked: DashboardView.vue dueSoonTasks computed (lines 86-92)
  found: Filter uses `diffDays >= 0 && diffDays <= 3` — requires strictly non-negative diff
  implication: Any task due today is excluded once the day has started (diffDays goes negative)

- timestamp: 2026-04-09T01:00:00Z
  checked: DashboardView.vue upcomingTasks computed (lines 68-77)
  found: Filter uses `diffDays <= 3 && diffDays >= -7` — 7-day overdue buffer
  implication: Same tasks that fail dueSoonTasks pass upcomingTasks — explains Test 8 pass + Test 6 fail

- timestamp: 2026-04-09T01:00:00Z
  checked: TaskDetail.vue input element (line 428)
  found: `<input type="date">` produces "YYYY-MM-DD" strings (no time component)
  implication: `new Date("YYYY-MM-DD")` creates midnight UTC per ECMAScript spec

- timestamp: 2026-04-09T01:00:00Z
  checked: domain/task.go DueDate field (line 15)
  found: DueDate is `string` — stored as plain text in SQLite (TEXT DEFAULT '')
  implication: Due date is always date-only format, no time component

- timestamp: 2026-04-09T01:00:00Z
  checked: stores/tasks.ts startOfDay helper (lines 12-14)
  found: A `startOfDay()` function exists for date-only comparisons but is NOT used by dueSoonTasks
  implication: The correct pattern exists in the codebase but wasn't applied to the Attention Bar filter

## Resolution

root_cause: `dueSoonTasks` computed property in DashboardView.vue (line 89) uses `diffDays >= 0` which creates a date-precision mismatch. The `<input type="date">` stores dates as "YYYY-MM-DD" (date-only). JavaScript's `new Date("YYYY-MM-DD")` creates midnight UTC. Once the current time passes midnight UTC (immediately for most timezones), `diffDays` for today becomes negative and the task is excluded. The identical data appears correctly in the Upcoming section (Test 8 passed) because that filter uses `diffDays >= -7`.
fix: ""
verification: ""
files_changed: []
