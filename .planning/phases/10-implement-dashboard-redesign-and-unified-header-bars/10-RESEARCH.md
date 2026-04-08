# Phase 10: Implement Dashboard Redesign and Unified Header Bars - Research

**Researched:** 2026-04-08
**Domain:** Vue 3 frontend UI — AppShell refactor, Teleport pattern, dashboard layout
**Confidence:** HIGH (all from direct codebase inspection + playground prototypes)

## Summary

Phase 10 is a **frontend-only visual/layout overhaul** — no new backend functionality, no API changes. The approved design is already fully prototyped in two playground files (`PlaygroundDashboard.vue` and `PlaygroundDashboardHeader.vue`). The work is: (1) refactor `AppShell.vue`'s top bar from 2-zone to 3-zone layout, (2) add a second teleport target for the center zone, (3) update every view to teleport its center content there, and (4) rewrite `DashboardView.vue` with the Attention Bar, Upcoming section, and richer task rows.

The core pattern is already established in the codebase: `DashboardView`, `TasksView`, and `AdoView` already use `<Teleport to="#topbar-actions">` for page-specific header content. This pattern extends naturally — AppShell gets a second div `#topbar-center`, and the sync cluster in the left zone becomes a shared component (`SyncCluster.vue`) used by all 6 pages.

**Primary recommendation:** Extract sync cluster markup from the playground into `SyncCluster.vue`, add `#topbar-center` div to AppShell, then update each view (Dashboard → Tasks → ADO → Projects → Deps → Settings) one at a time.

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Dashboard Layout:**
- Keep 2-column layout: tasks left (3/5 width), PRs+Pipelines right (2/5 width)
- Keep greeting + summary line at top ("Good morning, Luis" + "3 tasks in progress · 1 blocked")
- Keep inline compact stats line below greeting ("3 in progress · 1 blocked · 5 done of 47")
- Activity sidebar stays as the global AppShell panel (240px, toggleable via clock icon) — not dashboard-specific

**Dashboard Sections (Left Column):**
- **Today's Focus** — tasks with status `in_progress` or `in_review` (kept from current)
- **Upcoming** — tasks due within 3 days (replaces "Recent Activity" section). Overdue tasks (past due date) highlighted with red indicator
- **Blocked** — tasks with status `blocked`, shown with red left border accent + blocked reason text (from playground)

**Attention Bar:**
- Horizontal scrollable row of urgency nudges between greeting and stats line
- Three nudge types:
  - **Due soon**: amber background, CalendarDays icon, "3 due within 3 days" with first 2 task title previews
  - **Pipeline failure**: red background, XCircle icon, pipeline name + branch name
  - **PR approval ready**: green background, GitMerge icon, "PR #N has N approvals — ready to merge"
- Each nudge is a pill-shaped card with icon + text, horizontally scrollable on overflow
- Only shown when relevant (conditional rendering per nudge type)

**Dashboard Task Rows (Richer):**
- Enhanced rows in Today's Focus and Upcoming sections (not reusing existing TaskRow component):
  - Priority dot (color: P0=red, P1=orange, P2=amber, P3=zinc)
  - ADO type icon (from `adoTypeIcon()`) or SquareCheckBig for personal tasks
  - Task title (truncated)
  - "personal" badge for personal tasks (small text, primary/8 bg)
  - Pending sync amber dot indicator
  - Status badge (blue for active, violet for in_review)
  - Due date (amber highlight if within 2 days)
- Blocked section shows blocked reason text below task title (italic, red-tinted)

**Right Column (PRs + Pipelines):**
- Keep current 3-section structure: Needs Your Review, Your Pull Requests, Pipelines
- No structural changes to right column — matches current production design

**Unified Header Bar Pattern:**
- Apply to all 6 pages: Dashboard, Tasks, ADO, Projects, Dependencies, Settings
- Modify AppShell top bar to support the 3-zone pattern:
  - **Left zone**: Page name + vertical divider + sync cluster
    - Sync cluster: colored dot (green=connected, red=offline) + "Synced"/"Offline" label + relative time ("3m ago") + pending changes badge (amber, "N pending") + refresh button
    - Sync cluster is consistent and identical on every page
  - **Center zone**: Page-specific content (teleported from each page)
    - Dashboard: stat badges ("5 active", "2 blocked", "31/47 done")
    - Tasks: status filter chips (All, Active, Blocked, Done)
    - ADO: tabs (Browser, PRs, Pipelines) with count badges
    - Projects: project count + active badge ("4 projects", "2 active")
    - Dependencies: graph stats ("8 nodes · 5 edges") + cycle badge
    - Settings: empty (no center content)
  - **Right zone**: Search ⌘K + New dropdown + Activity toggle (always the same, unchanged)

### Claude's Discretion
- Exact implementation of sync cluster component (shared component vs inline)
- Whether to extract the Attention Bar nudges into separate sub-components
- Exact responsive breakpoints for Attention Bar horizontal scroll
- Whether dashboard task rows should be a new DashboardTaskRow component or inline in DashboardView
- Animation/transition details for Attention Bar nudges appearing/disappearing

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope
</user_constraints>

---

## Standard Stack

### Core (already in project — no new installs)
| Library | Current Usage | Phase Role |
|---------|--------------|------------|
| Vue 3 + Composition API | All components | Component authoring |
| Tailwind v4 + shadcn-vue | All UI | Styling, Badge, Button, Tabs, ScrollArea |
| lucide-vue-next | All icons | CalendarDays, XCircle, GitMerge, RefreshCw, Loader2 |
| Pinia | All stores | `useSyncStore`, `useTaskStore`, `useADOStore`, `usePRStore` |
| `@/lib/styles.ts` | Task rows | `adoTypeIcon()`, `adoTypeColor()`, `priorityColor()` |
| `@/lib/date.ts` | Timestamps | `relativeTime()` |

### No New Dependencies
This phase requires **zero new package installs**. All libraries are already present.

---

## Architecture Patterns

### Current AppShell Top Bar Structure
```
[breadcrumb] [#topbar-actions ← views teleport here] [flex-1 spacer] [search+new+clock]
```

### Proposed AppShell Top Bar Structure (3-zone)
```
[breadcrumb] [divider] [sync-cluster] | [#topbar-center ← views teleport here] | [search+new+clock]
```

The key change in AppShell.vue template:
```html
<!-- LEFT ZONE -->
<span class="text-sm font-medium shrink-0">{{ breadcrumb }}</span>
<span class="w-px h-4 bg-border mx-3" />
<SyncCluster />  <!-- or inline sync cluster markup -->

<!-- CENTER ZONE — views teleport here -->
<div id="topbar-center" class="flex-1 flex items-center justify-center gap-2 titlebar-no-drag" />

<!-- RIGHT ZONE — unchanged -->
<div v-if="showActions" class="flex items-center gap-2 titlebar-no-drag shrink-0">
  <!-- search + new dropdown + clock — NO CHANGES HERE -->
</div>
```

### Teleport Pattern (established, to extend)

**Current (works today):**
```html
<!-- In DashboardView.vue, TasksView.vue, AdoView.vue -->
<Teleport v-if="isActive" to="#topbar-actions">
  <!-- page-specific content -->
</Teleport>
```

**New pattern (center zone):**
```html
<!-- Each view teleports to #topbar-center -->
<Teleport v-if="isActive" to="#topbar-center">
  <!-- page-specific center content -->
</Teleport>
```

**Migration path per view:**
- `DashboardView` — change target `#topbar-actions` → `#topbar-center`, update content to stat badges
- `TasksView` — change target `#topbar-actions` → `#topbar-center`, keep status chips (remove sync button — now in left zone)
- `AdoView` — change target `#topbar-actions` → `#topbar-center`, keep tabs only (remove Connected/Sync/time — now in sync cluster)
- `ProjectsView` — add new Teleport to `#topbar-center` with project count + active badge
- `DependencyGraphView` — add new Teleport to `#topbar-center` with graph stats + cycle badge
- `SettingsView` — add new Teleport to `#topbar-center` (empty — no center content needed)

### SyncCluster Component (recommended extraction)

Extract the sync cluster markup into `frontend/src/components/SyncCluster.vue`:

```html
<!-- Source: PlaygroundDashboardHeader.vue proposed sync cluster markup -->
<template>
  <div class="flex items-center gap-2 text-[11px] shrink-0 titlebar-no-drag">
    <span class="inline-flex items-center gap-1.5"
      :class="adoStore.connected ? 'text-green-600' : 'text-red-500'">
      <span class="w-1.5 h-1.5 rounded-full"
        :class="adoStore.connected ? 'bg-green-500' : 'bg-red-500'" />
      {{ adoStore.connected ? 'Synced' : 'Offline' }}
    </span>
    <span v-if="syncStore.lastSyncedAt" class="text-muted-foreground/50 tabular-nums">
      {{ relativeTime(syncStore.lastSyncedAt) }}
    </span>
    <Badge v-if="pendingCount > 0" variant="outline"
      class="text-[9px] h-4 px-1.5 gap-1 text-amber-600 border-amber-500/30 bg-amber-500/10">
      <span class="w-1 h-1 rounded-full bg-amber-500" /> {{ pendingCount }} pending
    </Badge>
    <Button variant="ghost" size="sm" class="h-6 w-6 p-0"
      @click="syncStore.manualSync()" :disabled="syncStore.syncing">
      <component :is="syncStore.syncing ? Loader2 : RefreshCw"
        :size="12" :class="syncStore.syncing && 'animate-spin'" class="text-muted-foreground" />
    </Button>
  </div>
</template>
```

**Data sources for SyncCluster:**
- `syncStore.syncing` → spinner state
- `syncStore.lastSyncedAt` → relative time display
- `adoStore.connected` → green/red dot
- `syncStore.manualSync()` → refresh button action
- **`pendingCount`** — NOTE: The sync store does NOT have a raw pending changes count. `pendingDiff` is a single `SyncDiff | null`. For the pending badge, count tasks with unsynced local changes OR use `syncStore.conflicts.length`. This needs a design decision — Claude's discretion.

### Dashboard Richer Task Rows

The playground (`PlaygroundDashboard.vue`) has the exact markup for the richer task rows:

```html
<!-- Priority dot + ADO type icon + title + badges + due date -->
<div class="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-muted/50">
  <span :class="cn('size-2 rounded-full shrink-0', priorityColor(task.priority))" />
  <component v-if="task.adoType" :is="adoTypeIcon(task.adoType)"
    :size="14" :class="adoTypeColor(task.adoType)" class="shrink-0" />
  <SquareCheckBig v-else :size="14" class="text-primary/60 shrink-0" />
  <span class="text-sm flex-1 truncate text-foreground">{{ task.title }}</span>
  <span v-if="task.isPersonal"
    class="text-[8px] px-1 py-0.5 rounded bg-primary/8 text-primary/70 border border-primary/10 shrink-0">
    personal
  </span>
  <span v-if="task.pendingSync" class="size-1.5 rounded-full bg-amber-500 shrink-0" />
  <Badge variant="outline" class="text-[9px] h-4 px-1.5 shrink-0"
    :class="task.status === 'in_review'
      ? 'border-violet-500/30 text-violet-600'
      : 'border-blue-500/30 text-blue-600'">
    {{ task.status === 'in_review' ? 'Review' : 'Active' }}
  </Badge>
  <span v-if="task.dueDate" class="text-[10px] tabular-nums shrink-0"
    :class="isDueSoon(task.dueDate) ? 'text-amber-600 font-medium' : 'text-muted-foreground'">
    {{ formatDueDate(task.dueDate) }}
  </span>
</div>
```

**Note on `task.pendingSync`:** The `Task` type from `stores/tasks.ts` does NOT currently have a `pendingSync` field. Playground used mock data. Either add this field to the Task type + backend, or omit the pending sync dot from richer rows. Claude's discretion.

**Note on `task.isPersonal`:** Tasks are "public" when they have ADO links (`isPublic` is computed from `task_ado_links` presence). `isPersonal = !task.isPublic`. The task store exposes `isPublicTask(taskId)` method. Map this in computed.

### Upcoming Section (replaces Recent Activity)

Filter logic for Upcoming:
```typescript
const upcomingTasks = computed(() =>
  taskStore.tasks.filter(t => {
    if (!t.dueDate || t.status === 'done' || t.status === 'cancelled') return false
    const due = new Date(t.dueDate)
    const now = new Date()
    const diff = (due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
    return diff >= -1 && diff <= 3  // include overdue (negative diff)
  })
    .sort((a, b) => new Date(a.dueDate!).getTime() - new Date(b.dueDate!).getTime())
)
```

Overdue highlight: `new Date(task.dueDate) < new Date()` → red indicator dot or red text.

### Attention Bar — Data Sources

| Nudge | Data Source | Condition |
|-------|-------------|-----------|
| Due soon | `taskStore.tasks` filter by dueDate ≤ 3 days | `dueSoonTasks.length > 0` |
| Pipeline failure | `adoStore.pipelines` filter by `result === 'failed'` | `.some(p => p.result === 'failed')` |
| PR approval ready | `prStore.myPRs` filter by `votes >= 2` | `.some(pr => pr.votes >= 2)` |

All three nudges use **conditional rendering** (`v-if`) — shown only when relevant data exists.

### Recommended Project Structure Changes

```
frontend/src/
├── components/
│   ├── SyncCluster.vue        ← NEW: shared sync status + refresh
│   └── tasks/
│       └── DashboardTaskRow.vue  ← OPTIONAL: richer row component (or inline)
├── layouts/
│   └── AppShell.vue           ← MODIFY: 3-zone top bar, add #topbar-center
└── views/
    ├── DashboardView.vue      ← MAJOR REWRITE: Attention Bar, Upcoming, richer rows
    ├── TasksView.vue          ← MINOR: change teleport target to #topbar-center
    ├── AdoView.vue            ← MINOR: change teleport target, remove redundant sync UI
    ├── ProjectsView.vue       ← ADD: teleport center zone content
    ├── DependencyGraphView.vue ← ADD: teleport center zone content
    └── SettingsView.vue       ← ADD: teleport (empty center)
```

### Anti-Patterns to Avoid
- **Don't use both `#topbar-actions` and `#topbar-center` simultaneously** — rename the existing target to `#topbar-center` so all views use one consistent target. Update existing teleports in Dashboard, Tasks, ADO at the same time AppShell is refactored.
- **Don't duplicate sync cluster markup across 6 views** — extract to `SyncCluster.vue` in AppShell, render once.
- **Don't put the sync cluster inside each view's teleport** — the sync cluster is AppShell-level (lives in AppShell, not teleported from views).
- **Don't reuse `TaskRow.vue` for dashboard richer rows** — dashboard rows have different layout (priority dot + ADO icon + personal badge) that doesn't fit TaskRow's API.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead |
|---------|-------------|-------------|
| Horizontal scroll for Attention Bar | Custom JS scroll logic | `overflow-x-auto` + `flex gap-3` — CSS handles it |
| Sync status time formatting | Custom time formatter | `relativeTime()` from `@/lib/date.ts` |
| Priority dot color | Inline color mapping | `priorityColor()` in playground (P0=red/P1=orange/P2=amber/P3=zinc) — add to `styles.ts` as `priorityBgColor()` |
| ADO type icon lookup | Inline switch | `adoTypeIcon()` + `adoTypeColor()` from `@/lib/styles.ts` |
| Teleport deduplication | Conditional logic | `v-if="isActive"` pattern (already established in Dashboard/Tasks/ADO) |
| Spinner during sync | Custom spinner | `<component :is="syncing ? Loader2 : RefreshCw" :class="syncing && 'animate-spin'" />` |

---

## Common Pitfalls

### Pitfall 1: Teleport Target Rename Timing
**What goes wrong:** Renaming `#topbar-actions` to `#topbar-center` in AppShell while views still teleport to `#topbar-actions` → content disappears silently (no error, just vanishes).
**Why it happens:** Vue teleport to non-existent target fails silently.
**How to avoid:** Update AppShell AND all three existing teleport views (Dashboard, Tasks, ADO) in the same commit/wave. Never leave them out of sync.
**Warning signs:** Top bar shows no page-specific content on any page after AppShell change.

### Pitfall 2: `pendingSync` Field Missing from Task Type
**What goes wrong:** Dashboard task rows reference `task.pendingSync` → TypeScript error, runtime undefined.
**Why it happens:** Playground used mock data with this field; production `Task` type doesn't have it.
**How to avoid:** Either add `pendingSync?: boolean` to the Task domain type (requires backend field) OR omit the amber dot from richer rows entirely. Simplest fix: omit for now.

### Pitfall 3: `pendingChanges` Count for Sync Cluster Badge
**What goes wrong:** Sync cluster shows "N pending" badge but `syncStore` doesn't have a raw pending count.
**Why it happens:** The playground used mock data `mockSync.pendingChanges`. In production, `syncStore.pendingDiff` is a single diff object, not a count.
**How to avoid:** Use `syncStore.conflicts.length` for the pending badge count, OR drive it from tasks with unsynced state. Simplest: omit the pending badge if count cannot be derived, OR always show 0 unless `conflicts.length > 0`.

### Pitfall 4: `isPersonal` Not on Task Type Directly
**What goes wrong:** Dashboard rows reference `task.isPersonal` → undefined.
**Why it happens:** `isPersonal` is not a stored field — it's derived from whether the task has ADO links.
**How to avoid:** Use `taskStore.isPublicTask(task.id)` to derive it, or compute a `personalTaskIds` set in DashboardView.

### Pitfall 5: AppShell `titlebar-drag` Conflict with Interactive Center Zone
**What goes wrong:** Center zone tabs/badges/buttons inside the top bar don't respond to clicks because the top bar div has `titlebar-drag` class (Wails makes the whole bar draggable).
**Why it happens:** AppShell top bar has `titlebar-drag` on the outer div; interactive elements inside need `titlebar-no-drag`.
**How to avoid:** Add `titlebar-no-drag` class to `#topbar-center` div (already done for `#topbar-actions`).

### Pitfall 6: Stats Line Ordering (greeting vs Attention Bar vs stats)
**What goes wrong:** Placing Attention Bar after stats line instead of between greeting and stats.
**Design spec order:** Greeting → Attention Bar → Stats line → Today's Focus/Upcoming/Blocked sections.
**How to avoid:** Follow PlaygroundDashboard.vue's section ordering exactly.

---

## Code Examples

### AppShell 3-Zone Top Bar
```html
<!-- Source: PlaygroundDashboardHeader.vue proposed markup -->
<div class="h-[46px] flex items-center gap-2 px-4 border-b border-border titlebar-drag shrink-0">
  <!-- LEFT ZONE: page name + divider + sync cluster -->
  <span class="text-sm font-medium text-foreground titlebar-no-drag select-none shrink-0">
    {{ breadcrumb }}
  </span>
  <span class="w-px h-4 bg-border mx-1 shrink-0" />
  <SyncCluster class="titlebar-no-drag" />

  <!-- CENTER ZONE: page-specific teleport target -->
  <div id="topbar-center" class="flex-1 flex items-center justify-center gap-2 titlebar-no-drag" />

  <!-- RIGHT ZONE: search + new + activity (unchanged) -->
  <div v-if="showActions" class="flex items-center gap-2 titlebar-no-drag shrink-0">
    <!-- ... existing search/new/clock buttons ... -->
  </div>
</div>
```

### SyncCluster.vue (extracted component)
```html
<!-- Uses: adoStore.connected, syncStore.syncing, syncStore.lastSyncedAt, syncStore.manualSync() -->
<template>
  <div class="flex items-center gap-2 text-[11px]">
    <span class="inline-flex items-center gap-1.5"
      :class="adoStore.connected ? 'text-green-600' : 'text-red-500'">
      <span class="w-1.5 h-1.5 rounded-full"
        :class="adoStore.connected ? 'bg-green-500' : 'bg-red-500'" />
      {{ adoStore.connected ? 'Synced' : 'Offline' }}
    </span>
    <span v-if="syncStore.lastSyncedAt && adoStore.connected"
      class="text-muted-foreground/50 tabular-nums">
      {{ relativeTime(syncStore.lastSyncedAt) }}
    </span>
    <Button variant="ghost" size="sm" class="h-6 w-6 p-0"
      @click="syncStore.manualSync()" :disabled="syncStore.syncing">
      <component :is="syncStore.syncing ? Loader2 : RefreshCw"
        :size="12" :class="syncStore.syncing && 'animate-spin'"
        class="text-muted-foreground" />
    </Button>
  </div>
</template>
```

### Attention Bar (from PlaygroundDashboard.vue)
```html
<!-- Horizontal scrollable urgency nudges — conditional rendering per type -->
<div class="flex gap-3 mb-5 overflow-x-auto pb-1">
  <!-- Due soon -->
  <div v-if="dueSoonTasks.length"
    class="flex items-center gap-2 px-3 py-2 rounded-lg bg-amber-500/8 border border-amber-500/20 shrink-0">
    <CalendarDays :size="14" class="text-amber-600 shrink-0" />
    <span class="text-xs text-amber-700 dark:text-amber-400">
      <strong>{{ dueSoonTasks.length }}</strong> due within 3 days
    </span>
    <div class="flex gap-1 ml-1">
      <span v-for="t in dueSoonTasks.slice(0, 2)" :key="t.id"
        class="text-[10px] bg-amber-500/10 px-1.5 py-0.5 rounded text-amber-700 dark:text-amber-400 truncate max-w-[120px]">
        {{ t.title }}
      </span>
    </div>
  </div>
  <!-- Pipeline failure -->
  <div v-if="failedPipeline"
    class="flex items-center gap-2 px-3 py-2 rounded-lg bg-red-500/8 border border-red-500/20 shrink-0">
    <XCircle :size="14" class="text-red-500 shrink-0" />
    <span class="text-xs text-red-600 dark:text-red-400">
      <strong>{{ failedPipeline.name }}</strong> failed on {{ failedPipeline.sourceBranch }}
    </span>
  </div>
  <!-- PR approval ready -->
  <div v-if="mergeReadyPR"
    class="flex items-center gap-2 px-3 py-2 rounded-lg bg-emerald-500/8 border border-emerald-500/20 shrink-0">
    <GitMerge :size="14" class="text-emerald-600 shrink-0" />
    <span class="text-xs text-emerald-700 dark:text-emerald-400">
      PR #{{ mergeReadyPR.prNumber }} has {{ mergeReadyPR.votes }} approvals — ready to merge
    </span>
  </div>
</div>
```

### Dashboard Center Zone (teleported stat badges)
```html
<!-- Source: PlaygroundDashboardHeader.vue proposed dashboard center zone -->
<Teleport v-if="isActive" to="#topbar-center">
  <div class="flex items-center gap-2 text-[10px]">
    <Badge variant="secondary" class="text-[9px] h-4 px-1.5">
      {{ taskStore.stats.inProgress + taskStore.stats.inReview }} active
    </Badge>
    <Badge variant="destructive" class="text-[9px] h-4 px-1.5"
      v-if="taskStore.stats.blocked">
      {{ taskStore.stats.blocked }} blocked
    </Badge>
    <span class="text-muted-foreground/40 tabular-nums">
      {{ taskStore.stats.done }}/{{ taskStore.stats.total }} done
    </span>
  </div>
</Teleport>
```

### Blocked Section with Reason Text (from PlaygroundDashboard.vue)
```html
<div class="rounded-lg overflow-hidden border border-border border-l-2 border-l-red-500">
  <div v-for="task in blockedTasks" :key="task.id"
    class="flex flex-col gap-1 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors">
    <div class="flex items-center gap-2">
      <span :class="cn('size-2 rounded-full shrink-0', priorityBgColor(task.priority))" />
      <span class="text-sm flex-1 truncate">{{ task.title }}</span>
    </div>
    <p v-if="task.blockedReason"
      class="text-[11px] text-red-500/70 pl-4 italic">
      {{ task.blockedReason }}
    </p>
  </div>
</div>
```

**Note:** `task.blockedReason` — check if this field exists on the production `Task` type. If not, add it or omit the reason text.

---

## Key Implementation Facts

### Files to Modify (exhaustive list)
| File | Change Type | Scope |
|------|-------------|-------|
| `frontend/src/layouts/AppShell.vue` | Structural refactor | Replace `#topbar-actions` + spacer with 3-zone layout; rename div to `#topbar-center`; add `SyncCluster` |
| `frontend/src/views/DashboardView.vue` | Major rewrite | Add Attention Bar, replace Recent Activity → Upcoming, replace TaskRow with richer rows, update teleport target |
| `frontend/src/views/TasksView.vue` | Minor update | Change teleport target, remove sync button (now in SyncCluster) |
| `frontend/src/views/AdoView.vue` | Minor update | Change teleport target, remove redundant Connected/Sync/time UI |
| `frontend/src/views/ProjectsView.vue` | Add teleport | Add `#topbar-center` teleport with project count + active badge |
| `frontend/src/views/DependencyGraphView.vue` | Add teleport | Add `#topbar-center` teleport with graph stats + cycle badge |
| `frontend/src/views/SettingsView.vue` | Add teleport | Add `#topbar-center` teleport (empty or minimal) |
| `frontend/src/components/SyncCluster.vue` | New file | Sync status component for AppShell left zone |

### Files to Check (may need minor updates)
| File | Why |
|------|-----|
| `frontend/src/components/PageHeader.vue` | May become redundant if unified header absorbs page-specific content — keep for now, don't delete |
| `domain/task.go` (or equivalent) | Check if `blockedReason` and `pendingSync` fields exist |

### Store Data Available (confirmed)
| Data | Source |
|------|--------|
| `syncStore.syncing` | `useSyncStore()` |
| `syncStore.lastSyncedAt` | `useSyncStore()` |
| `syncStore.manualSync()` | `useSyncStore()` |
| `adoStore.connected` | `useADOStore()` |
| `adoStore.pipelines` | `useADOStore()` — filter by `result === 'failed'` |
| `taskStore.tasks` | `useTaskStore()` — filter by dueDate, status |
| `taskStore.stats` | `useTaskStore()` — inProgress, blocked, done, total |
| `prStore.myPRs` | `usePRStore()` — filter by `votes >= 2` |
| `prStore.reviewPRs` | `usePRStore()` |

### `priorityBgColor()` helper needed
The playground defines this inline as `priorityColor()` but `styles.ts` only has `priorityColor()` (text color) and `priorityClasses()` (badge). The background dot color for P0-P3 needs a `priorityBgColor()` helper added to `styles.ts`:
```typescript
export function priorityBgColor(priority: string): string {
  switch (priority) {
    case 'P0': return 'bg-red-500'
    case 'P1': return 'bg-orange-500'
    case 'P2': return 'bg-amber-500'
    default: return 'bg-zinc-400'
  }
}
```

---

## State of the Art

| Old Approach | New Approach | Impact |
|--------------|-------------|--------|
| `#topbar-actions` div (left of spacer) | `#topbar-center` div (center of 3-zone bar) | All 6 pages get consistent header structure |
| Sync status per-page (Dashboard/Tasks/ADO only) | `SyncCluster` in AppShell left zone (every page) | Sync status always visible regardless of page |
| `Recent Activity` section (recent by time) | `Upcoming` section (due within 3 days) | Actionable priority focus vs passive history |
| `TaskRow` component in Dashboard sections | Richer inline dashboard rows | Priority dot + ADO type + personal badge + due date |
| ADO Connected/Sync button in ADO view's teleport | Moved to universal SyncCluster | Reduces duplication in AdoView |

---

## Open Questions

1. **`task.blockedReason` field existence**
   - What we know: Playground used it in mock data. Dashboard shows blocked reason text.
   - What's unclear: Does the production `Task` type and backend have this field?
   - Recommendation: Check `domain/task.go` before implementing blocked row. If missing, add the field in this phase or omit the reason text.

2. **`pendingChanges` count for Sync Cluster badge**
   - What we know: Playground showed "N pending" amber badge. `syncStore` has `conflicts.length` and `pendingDiff` (single diff).
   - What's unclear: Should "pending" mean conflict count or something else?
   - Recommendation: Use `syncStore.conflicts.length` for the pending badge, or omit if 0 most of the time and the badge adds noise.

3. **`task.pendingSync` field for richer rows**
   - What we know: Playground had amber dot for tasks with pending sync. Production Task type lacks this field.
   - What's unclear: Is there a way to derive it without a new backend field?
   - Recommendation: Omit the amber sync dot from richer rows for now (no backend support), or derive from `syncStore.conflicts` matching task IDs.

---

## Validation Architecture

> `nyquist_validation` is `false` in `.planning/config.json` — skip formal test infrastructure section.

---

## Sources

### Primary (HIGH confidence — direct codebase inspection)
- `frontend/src/views/playground/PlaygroundDashboardHeader.vue` — approved proposed header markup for all 6 pages
- `frontend/src/views/playground/PlaygroundDashboard.vue` — approved dashboard layout with Attention Bar, richer rows
- `frontend/src/layouts/AppShell.vue` — current production AppShell (lines 82–128)
- `frontend/src/views/DashboardView.vue` — current dashboard (lines 162–434)
- `frontend/src/views/TasksView.vue` — current Tasks view teleport (lines 296–318)
- `frontend/src/views/AdoView.vue` — current ADO view teleport (lines 244–296)
- `frontend/src/stores/sync.ts` — sync store API surface
- `frontend/src/stores/ado.ts` — `connected` state
- `frontend/src/lib/styles.ts` — confirmed `adoTypeIcon()`, `adoTypeColor()`, `priorityColor()`

### Secondary (HIGH confidence)
- `.planning/phases/10-implement-dashboard-redesign-and-unified-header-bars/10-CONTEXT.md` — locked decisions
- `.planning/STATE.md` — project context and decisions

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — direct file inspection, all libs already installed
- Architecture patterns: HIGH — playground prototypes are the exact source of truth
- Pitfalls: HIGH — derived from actual type/store gaps found in code inspection
- Implementation order: HIGH — clearly sequenced (AppShell first, then views)

**Research date:** 2026-04-08
**Valid until:** Stable (no external dependencies — all internal)
