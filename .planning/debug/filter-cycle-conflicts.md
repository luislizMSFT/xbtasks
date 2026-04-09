---
status: diagnosed
trigger: "Filter cycle button has conflicts"
created: 2025-01-09T20:00:00Z
updated: 2025-01-09T20:00:00Z
---

## Current Focus

hypothesis: CONFIRMED — Duplicate controls for filterAdoLink with inconsistent labels
test: n/a — root cause confirmed
expecting: n/a
next_action: Return diagnosis

## Symptoms

expected: A small filter button in the toolbar cycles through All → ADO → Personal on each click. The task list filters accordingly — ADO shows only ADO-linked tasks, Personal shows only personal tasks. Active filter is visually highlighted.
actual: User reported "There are conflicts"
errors: None reported
reproduction: Test 4 in UAT - click the filter cycle button in the toolbar
started: Discovered during UAT for phase 11

## Eliminated

- hypothesis: Git merge conflict markers in source files
  evidence: Searched all .vue/.ts/.js files in frontend/src for <<<<<<, >>>>>>, ======= — zero matches
  timestamp: 2025-01-09T20:05:00Z

- hypothesis: TypeScript compilation errors / build conflicts
  evidence: npx vue-tsc --noEmit exits cleanly with code 0 — no type errors
  timestamp: 2025-01-09T20:08:00Z

## Evidence

- timestamp: 2025-01-09T20:03:00Z
  checked: FilterCycleButton.vue source
  found: Clean implementation. Cycles 'all' → 'linked' → 'personal' via v-model on taskStore.filterAdoLink. Labels: All / ADO / Personal.
  implication: Component itself is correct

- timestamp: 2025-01-09T20:04:00Z
  checked: TasksView.vue lines 288-300 (topbar Teleport) and lines 306-326 (FilterBar)
  found: TWO controls for filterAdoLink — FilterCycleButton at line 299 (v-model="taskStore.filterAdoLink") AND FilterBar at lines 311/321 (:filter-ado-link / @update:filter-ado-link)
  implication: Both controls modify the same reactive state simultaneously

- timestamp: 2025-01-09T20:05:00Z
  checked: FilterBar.vue lines 139-152 (Scope dropdown)
  found: FilterBar has "Scope" dropdown with values all/linked/personal but labels All/Public/Private. These differ from FilterCycleButton labels All/ADO/Personal.
  implication: Inconsistent labels for same filter values — "ADO" vs "Public", "Personal" vs "Private"

- timestamp: 2025-01-09T20:06:00Z
  checked: 11-04-PLAN.md lines 150-166
  found: Plan explicitly preserves FilterBar with adoLink dropdown (line 152) AND adds FilterCycleButton to topbar (line 166). Design decision to keep both.
  implication: Duplicate controls were by design, but create the user-reported "conflict"

- timestamp: 2025-01-09T20:07:00Z
  checked: taskStore.filterAdoLink reactivity
  found: Both FilterCycleButton (v-model) and FilterBar Scope dropdown bind to same ref. Changing one updates the other reactively. No state fighting — but user sees two controls change simultaneously with different labels.
  implication: Not a code bug — a UX design conflict: redundant controls with inconsistent terminology

## Resolution

root_cause: The FilterCycleButton (topbar) and FilterBar "Scope" dropdown (filter bar) both control `taskStore.filterAdoLink` with inconsistent labels — FilterCycleButton uses "ADO"/"Personal" while FilterBar uses "Public"/"Private" for the same values. Two visible, simultaneous controls for the same filter create user-reported "conflicts". The 11-04-PLAN deliberately preserved both, but this creates UX confusion.
fix:
verification:
files_changed: []
