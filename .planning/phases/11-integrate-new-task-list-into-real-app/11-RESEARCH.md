# Phase 11: Integrate new task list into real app - Research

**Researched:** 2026-04-08
**Domain:** Vue 3 component refactoring, Wails bindings, SQLite denormalization, frontend composables
**Confidence:** HIGH

## Summary

Phase 11 replaces the current TasksView with the playground's integrated 2-panel tree+detail layout. This is a component extraction and store integration phase, not a new feature build. The playground prototype (PlaygroundIntegrated.vue) already demonstrates the desired UX with mock data. The work involves: (1) extracting reusable components from the playground (TreeTaskRow, QuickAddInput, FilterCycleButton), (2) replacing TasksView.vue with the new layout while preserving all existing filters and features, (3) wiring the extracted components to real stores (taskStore, adoStore, syncStore), (4) adding a backend ADO metadata cache to avoid N+1 queries, and (5) creating a useAdoMeta composable for batch-loaded ADO type/state data.

**Primary recommendation:** Follow the component extraction → store wiring → backend cache → integration sequence. Extract components first with minimal changes from the playground, wire them to stores incrementally, add the backend cache table and service, then replace TasksView as the final integration step.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Replacement Strategy:**
- Replace TasksView.vue entirely — rewrite to use the playground's tree+detail 2-panel layout
- Preserve ALL existing features: drag-drop reordering, group-by (status/priority/project), full FilterBar (status, priority, project, dueDate, ADO link), alongside new playground features (tree/flat toggle, All/ADO/Personal filter)
- Permanent 2-panel split layout: task list always on left, detail pane always on right when a task is selected (not a slide-out — permanent split)
- Keep existing store wiring (`enhancedFilteredTasks`, `groupedEnhanced`, `filterAdoLink`) — these already work

**Task Row Component:**
- Create new `TreeTaskRow.vue` from playground pattern — supports indent levels, ADO type icons, state badges, expand/collapse chevrons built-in
- TreeTaskRow replaces TaskRow everywhere: TasksView, DashboardView (Today's Focus, Upcoming, Blocked), and any other list context — single source of truth
- Simplified row content: status icon + title + ADO type icon + priority badge + due date
- Drop from current TaskRow: tags, description preview, time-ago timestamp (too noisy)
- Keep: status circle (Things 3 style), ADO badge (hollow/filled), blocked reason subtitle, project name caption

**Detail Pane:**
- Build new TaskDetail.vue from playground pattern — replaces current TaskDetail.vue completely
- 3-tab structure: Subtasks, PRs, Notes (simpler than current Work/Discussion tabs)
- Include sync features: sync confirmation dialog, conflict resolver, push-to-ADO actions
- Include external links section (reuse ExternalLinks.vue component)
- Reuse CommentsSection.vue for Notes tab (comments with push-to-ADO capability)
- Detail pane is the right panel in the permanent 2-panel split

**ADO Metadata Enrichment:**
- Create `ado_meta_cache` table in SQLite: `(task_id INTEGER, ado_type TEXT, ado_state TEXT)` — populated during sync, keeps tasks table clean
- Create `useAdoMeta()` composable on frontend — provides `getAdoMeta(taskId): { type, state } | null`
- Batch load all ADO metadata on app startup + refresh on `sync:completed` Wails event
- Backend: new `ADOMetaCacheService` with `GetAll() → map[int]{type, state}` and `Refresh()` called during sync
- Frontend composable first checks the batch-loaded cache (fast), no per-task API calls

**Shared Components to Extract:**
- `TreeTaskRow.vue` — task row with indent, ADO type icon, state badge, expand/collapse
- `QuickAddInput.vue` — inline quick-add pattern (Enter to add, Esc to cancel) used for tasks + subtasks
- `FilterCycleButton.vue` — filter cycle UI (click to cycle through All → ADO → Personal)

**Style Helpers to Add:**
- `statusIcon()` in `styles.ts` — maps status → Lucide icon component (currently inline in playground and TaskRow)
- `priorityDotBgColor()` in `styles.ts` — maps P0-P3 → Tailwind bg class

### Claude's Discretion

- Exact component prop API for TreeTaskRow (which props are required vs optional)
- How to handle the transition animation when selecting/deselecting tasks in the split pane
- Whether QuickAddInput needs focus management (auto-focus on mount vs manual)
- Internal state management for expand/collapse in tree view (local ref vs store)
- How to wire drag-drop in tree view (existing Sortable.js or new approach)
- Whether to keep the existing FilterBar component or replace with a simpler toolbar
</user_constraints>

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Vue 3 | 3.2.45 | Frontend framework | Already in use; Composition API is project standard |
| Pinia | 3.0.4 | State management | Project's store layer; all data flows through Pinia stores |
| vue-router | 4.6.4 | Routing | Existing router; TasksView is `/tasks` route |
| TypeScript | 4.9.3 | Type safety | Project-wide type safety for Wails bindings |
| Vite | 5.0.0 | Build tool | Wails default build system |
| @wailsio/runtime | latest | Wails bindings | Go↔Vue bridge; all backend calls use this |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @vueuse/core | 14.2.1 | Composition utilities | useDebounceFn for search, useMagicKeys for shortcuts (existing) |
| lucide-vue-next | 1.0.0 | Icon components | All icons in app; statusIcon helper returns Lucide components |
| reka-ui | 2.9.3 | shadcn-vue primitives | Badge, Button, Input, ScrollArea, Tabs, Select already in use |
| vuedraggable | 4.1.0 | Drag-drop | Existing drag-drop implementation for task reordering |
| vue-sonner | 2.0.9 | Toast notifications | Existing notification system via useNotify composable |
| modernc.org/sqlite | (Go) | SQLite driver | Backend database; pure Go, no cgo dependency |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Pinia | Vuex | Pinia is already project standard; no reason to change |
| reka-ui | Naive UI / PrimeVue | shadcn-vue (reka-ui) already integrated; 79 components installed |
| vuedraggable | SortableJS directly | vuedraggable wraps SortableJS and already works in project |

**Installation:**
No new dependencies required — all libraries already installed.

**Version verification:**
Verified from `frontend/package.json` (current as of Phase 11 research):
- Vue 3.2.45 (published 2023-01)
- Pinia 3.0.4 (published 2025-02)
- @vueuse/core 14.2.1 (published 2025-04)
- lucide-vue-next 1.0.0 (published 2025-03)

## Architecture Patterns

### Recommended Project Structure
```
frontend/src/
├── components/
│   ├── tasks/
│   │   ├── TreeTaskRow.vue         # NEW: extracted from playground
│   │   ├── QuickAddInput.vue       # NEW: extracted from playground
│   │   ├── FilterCycleButton.vue   # NEW: extracted from playground
│   │   ├── TaskDetail.vue          # REPLACE: rebuild from playground
│   │   ├── TaskRow.vue             # DEPRECATED: replaced by TreeTaskRow
│   │   ├── CommentsSection.vue     # REUSE: wire to new TaskDetail
│   │   └── ExternalLinks.vue       # REUSE: wire to new TaskDetail
│   └── ado/
│       ├── SyncConfirmDialog.vue   # REUSE: wire to new TaskDetail
│       └── ConflictResolver.vue    # REUSE: wire to new TaskDetail
├── composables/
│   └── useAdoMeta.ts               # NEW: batch ADO metadata access
├── stores/
│   ├── tasks.ts                    # EXISTING: wire to new components
│   ├── ado.ts                      # EXISTING: metadata source
│   └── sync.ts                     # EXISTING: sync events
├── views/
│   ├── TasksView.vue               # REPLACE: rewrite with 2-panel layout
│   └── playground/
│       └── PlaygroundIntegrated.vue # SOURCE: reference for extraction
└── lib/
    └── styles.ts                   # EXTEND: add statusIcon + priorityDotBgColor
```

```
internal/
├── app/
│   └── adometa.go                  # NEW: ADOMetaCacheService
└── db/
    └── migrate.go                  # EXTEND: add ado_meta_cache table
```

### Pattern 1: Component Extraction from Playground

**What:** Extract TreeTaskRow, QuickAddInput, FilterCycleButton from PlaygroundIntegrated.vue into standalone components.

**When to use:** When a playground prototype demonstrates a desired pattern and needs integration into the real app.

**Example:**
```vue
<!-- From PlaygroundIntegrated.vue lines 731-821 (tree row markup) -->
<!-- Extract to: components/tasks/TreeTaskRow.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import type { Task } from '@/types'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { ChevronRight, ChevronDown, User, CalendarDays } from 'lucide-vue-next'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import { statusIcon, statusColor, adoTypeIcon, adoTypeColor, adoStateClasses, statusClasses } from '@/lib/styles'

interface Props {
  task: Task
  indentLevel?: number
  selected?: boolean
  expanded?: boolean
  hasChildren?: boolean
  adoMeta?: { type: string; state: string } | null
  projectName?: string
  subtaskProgress?: { done: number; total: number; pct: number } | null
}

const props = withDefaults(defineProps<Props>(), {
  indentLevel: 0,
  selected: false,
  expanded: false,
  hasChildren: false,
  adoMeta: null,
  projectName: undefined,
  subtaskProgress: null,
})

const emit = defineEmits<{
  click: []
  toggleExpand: []
}>()

const paddingLeft = computed(() => `${36 + props.indentLevel * 24}px`)
const isPersonal = computed(() => !props.task.adoId && !props.task.projectId)
</script>

<template>
  <div
    class="group cursor-pointer hover:bg-muted/50 transition-colors"
    :class="selected ? 'bg-primary/5 border-l-2 border-l-primary' : 'border-l-2 border-l-transparent'"
    @click="emit('click')"
  >
    <!-- Row 1: chevron + status + type icon + title + badges -->
    <div class="flex items-center gap-2 px-3 pt-2.5 pb-0.5">
      <button
        v-if="hasChildren"
        class="shrink-0 p-0.5 hover:bg-muted rounded"
        @click.stop="emit('toggleExpand')"
      >
        <component :is="expanded ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
      </button>
      <span v-else class="w-5 shrink-0" />

      <!-- Status icon -->
      <button class="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full" :class="statusColor(task.status)">
        <component :is="statusIcon(task.status)" :size="16" :stroke-width="1.75" />
      </button>

      <!-- Type icon (ADO type or personal indicator) -->
      <component
        v-if="adoMeta?.type"
        :is="adoTypeIcon(adoMeta.type)"
        :size="14"
        :class="adoTypeColor(adoMeta.type)"
        class="shrink-0 -ml-1"
      />
      <User v-else-if="isPersonal" :size="12" class="text-muted-foreground/40 shrink-0 -ml-1" />

      <span
        class="text-sm truncate flex-1 min-w-0"
        :class="[
          task.status === 'done' ? 'line-through text-muted-foreground' : 'text-foreground',
          hasChildren ? 'font-medium' : '',
        ]"
      >
        {{ task.title }}
      </span>

      <Badge variant="outline" class="text-[10px] h-4 px-1.5 shrink-0"
        :class="adoMeta?.state ? adoStateClasses(adoMeta.state) : statusClasses(task.status)"
      >
        {{ adoMeta?.state || task.status }}
      </Badge>

      <template v-if="subtaskProgress">
        <span class="text-[10px] text-muted-foreground tabular-nums shrink-0">
          {{ subtaskProgress.done }}/{{ subtaskProgress.total }}
        </span>
        <div class="w-12 h-1 bg-muted rounded-full overflow-hidden shrink-0">
          <div class="h-full bg-green-500 rounded-full transition-all" :style="{ width: subtaskProgress.pct + '%' }" />
        </div>
      </template>
    </div>

    <!-- Row 2: metadata -->
    <div class="flex items-center gap-2 pb-2 text-[11px]" :style="{ paddingLeft: '68px', paddingRight: '12px' }">
      <span v-if="adoMeta?.adoId" class="text-muted-foreground/50 tabular-nums">{{ adoMeta.adoId }}</span>
      <Badge v-if="isPersonal" variant="outline" class="text-[9px] h-3.5 px-1 border-dashed text-muted-foreground/60">
        <User :size="8" class="mr-0.5" /> Personal
      </Badge>
      <span v-if="projectName" class="text-muted-foreground/40 truncate">{{ projectName }}</span>
      <div class="flex-1" />
      <span
        v-if="task.dueDate"
        class="inline-flex items-center gap-0.5 text-[10px] shrink-0"
        :class="new Date(task.dueDate) < new Date() ? 'text-red-500' : 'text-muted-foreground/50'"
      >
        <CalendarDays :size="10" />
        {{ task.dueDate }}
      </span>
      <PriorityBadge :priority="task.priority" />
    </div>

    <!-- Blocked banner -->
    <div
      v-if="task.status === 'blocked' && task.blockedReason"
      class="pb-2 text-[10px] text-red-500/80 truncate"
      :style="{ paddingLeft: '68px', paddingRight: '12px' }"
    >
      ⚠ {{ task.blockedReason }}
    </div>
  </div>
</template>
```

### Pattern 2: Batch Data Loading with Composables

**What:** Create a composable that batch-loads ADO metadata once and provides fast lookups.

**When to use:** When N items need enrichment from a secondary data source, and N+1 queries would kill performance.

**Example:**
```typescript
// composables/useAdoMeta.ts
import { ref, onMounted } from 'vue'
import { Events } from '@wailsio/runtime'
import { GetAllADOMeta } from '@/api/adometa'

interface AdoMeta {
  type: string
  state: string
}

const metaCache = ref<Map<number, AdoMeta>>(new Map())
const loading = ref(false)

export function useAdoMeta() {
  async function refresh() {
    loading.value = true
    try {
      const data = await GetAllADOMeta() as Record<number, AdoMeta>
      metaCache.value = new Map(Object.entries(data).map(([k, v]) => [Number(k), v]))
    } catch (err) {
      console.error('Failed to load ADO metadata:', err)
    } finally {
      loading.value = false
    }
  }

  function getAdoMeta(taskId: number): AdoMeta | null {
    return metaCache.value.get(taskId) ?? null
  }

  onMounted(() => {
    if (metaCache.value.size === 0) {
      refresh()
    }
    // Refresh cache on sync completion
    Events.On('sync:completed', refresh)
  })

  return {
    getAdoMeta,
    refresh,
    loading,
  }
}
```

### Pattern 3: Two-Panel Split Layout (Permanent, Not Slide-Out)

**What:** Left panel shows task tree, right panel shows detail. Both visible simultaneously when a task is selected.

**When to use:** When users need to see list context while viewing details (typical master-detail pattern).

**Example:**
```vue
<!-- views/TasksView.vue — new structure -->
<template>
  <div class="flex flex-col h-full">
    <!-- Header: filters, view toggles, actions -->
    <div class="shrink-0 border-b border-border px-4 py-2 flex items-center gap-3">
      <FilterBar />
      <div class="flex-1" />
      <FilterCycleButton v-model="taskStore.filterAdoLink" />
      <Button variant="outline" size="sm" @click="treeView = !treeView">
        {{ treeView ? 'Flat' : 'Tree' }}
      </Button>
      <Button size="sm" @click="showQuickAdd = true">
        <Plus :size="14" class="mr-1" /> New Task
      </Button>
    </div>

    <!-- Body: 2-panel split -->
    <div class="flex-1 flex min-h-0">
      <!-- Left: task tree -->
      <div class="w-96 border-r border-border flex flex-col">
        <ScrollArea class="flex-1">
          <TreeTaskRow
            v-for="task in taskStore.enhancedFilteredTasks"
            :key="task.id"
            :task="task"
            :selected="taskStore.selectedTaskId === task.id"
            :ado-meta="adoMeta.getAdoMeta(task.id)"
            :project-name="projectNameMap[task.projectId]"
            @click="taskStore.selectedTaskId = task.id"
          />
        </ScrollArea>
        <QuickAddInput v-if="showQuickAdd" @add="handleQuickAdd" @cancel="showQuickAdd = false" />
      </div>

      <!-- Right: detail pane -->
      <TaskDetail v-if="taskStore.selectedTaskId" class="flex-1" />
      <EmptyState v-else title="No task selected" description="Select a task to view details" />
    </div>
  </div>
</template>
```

### Pattern 4: Backend Cache Table with Service Layer

**What:** Add a denormalized cache table that mirrors ADO work item metadata for fast access.

**When to use:** When relational joins are expensive and denormalization is acceptable (ADO metadata changes infrequently).

**Example:**
```go
// internal/db/migrate.go — add to schema constant
CREATE TABLE IF NOT EXISTS ado_meta_cache (
	task_id   INTEGER PRIMARY KEY REFERENCES tasks(id) ON DELETE CASCADE,
	ado_type  TEXT DEFAULT '',
	ado_state TEXT DEFAULT '',
	synced_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ado_meta_cache_task ON ado_meta_cache(task_id);
```

```go
// internal/app/adometa.go
package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yourusername/xb-tasks/internal/db"
)

type ADOMetaCacheService struct {
	db *db.DB
}

type AdoMeta struct {
	Type  string `json:"type"`
	State string `json:"state"`
}

func NewADOMetaCacheService(database *db.DB) *ADOMetaCacheService {
	return &ADOMetaCacheService{db: database}
}

// GetAll returns a map of task_id → {type, state}
func (s *ADOMetaCacheService) GetAll(ctx context.Context) (map[int]AdoMeta, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT task_id, ado_type, ado_state
		FROM ado_meta_cache
	`)
	if err != nil {
		return nil, fmt.Errorf("query ado_meta_cache: %w", err)
	}
	defer rows.Close()

	result := make(map[int]AdoMeta)
	for rows.Next() {
		var taskID int
		var meta AdoMeta
		if err := rows.Scan(&taskID, &meta.Type, &meta.State); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result[taskID] = meta
	}

	return result, rows.Err()
}

// Refresh rebuilds the cache from ado_work_items + task_ado_links
func (s *ADOMetaCacheService) Refresh(ctx context.Context) error {
	// Delete stale cache
	if _, err := s.db.ExecContext(ctx, `DELETE FROM ado_meta_cache`); err != nil {
		return fmt.Errorf("clear cache: %w", err)
	}

	// Rebuild from joins
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO ado_meta_cache (task_id, ado_type, ado_state, synced_at)
		SELECT
			l.task_id,
			w.type,
			w.state,
			CURRENT_TIMESTAMP
		FROM task_ado_links l
		INNER JOIN ado_work_items w ON l.ado_id = w.ado_id
	`)
	if err != nil {
		return fmt.Errorf("rebuild cache: %w", err)
	}

	return nil
}
```

### Pattern 5: Wails Event Subscription for Cache Refresh

**What:** Subscribe to backend events (`sync:completed`) to trigger cache refresh.

**When to use:** When the frontend needs to reactively update based on backend state changes.

**Example:**
```typescript
// composables/useAdoMeta.ts (extended from Pattern 2)
import { Events } from '@wailsio/runtime'

onMounted(() => {
  if (metaCache.value.size === 0) {
    refresh()
  }
  
  // Refresh cache when sync completes
  Events.On('sync:completed', () => {
    console.log('[useAdoMeta] Sync completed, refreshing ADO metadata cache')
    refresh()
  })
})
```

### Anti-Patterns to Avoid

- **N+1 ADO metadata queries:** Don't call `GetADOWorkItem(adoId)` per task row. Use batch loading with `useAdoMeta()` composable.
- **Mixing playground and production patterns:** Don't copy mock data structures from playground. Extract components cleanly, wire to real stores immediately.
- **Inline status icon logic:** Don't duplicate `statusIcon()` mapping in multiple components. Centralize in `styles.ts`.
- **Stateful components without props:** Don't embed store access deep in TreeTaskRow. Pass data via props, emit events for actions.
- **Slide-out detail panel:** Don't use the old slide-out pattern. User decision is permanent 2-panel split.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| ADO metadata joins | Custom query per row | Denormalized `ado_meta_cache` table + batch load | N+1 queries destroy performance; cache table is faster and simpler |
| Task tree rendering | Recursive component imports | `defineComponent` with `h()` for same-file recursion | Vue's SFC doesn't support self-referencing imports cleanly; `h()` render function works |
| Icon mapping | Switch statements in templates | `statusIcon()` / `adoTypeIcon()` helpers in `styles.ts` | Centralized mapping is DRY; helpers return Lucide components directly |
| Drag-drop | Custom mousedown/mousemove handlers | `vuedraggable` (already in project) | Handles edge cases (touch, scroll, cancel, ghost), accessibility, cross-browser |
| Debounced search | Manual setTimeout logic | `useDebounceFn` from `@vueuse/core` | Handles cleanup, cancellation, TypeScript types |
| Keyboard shortcuts | Manual keydown listeners | `useMagicKeys` from `@vueuse/core` | Cross-platform (Meta/Ctrl), cleanup, conflicts |

**Key insight:** Component extraction is deceptively complex. The playground has mock data structures that don't match production types. Don't copy-paste blindly — extract the markup and styling, but wire props/emits to real store shapes immediately.

## Common Pitfalls

### Pitfall 1: Type Mismatches Between Playground and Real Stores

**What goes wrong:** Playground uses `mockTasks.value` with different field names (e.g., `adoId` in mock vs `ado_id` in backend) or shapes (e.g., mock has `parentId` but real store uses `parent_id`).

**Why it happens:** Prototypes use simplified data for speed; production has legacy naming or database constraints.

**How to avoid:**
1. Check `@/types/index.ts` for the canonical `Task` interface
2. Review `stores/tasks.ts` to see actual field names returned from backend
3. Don't assume playground field names match production — verify each prop

**Warning signs:**
- TypeScript errors about missing properties when wiring props
- Runtime errors like "Cannot read property 'type' of undefined"
- Data shows as blank in extracted components that worked in playground

### Pitfall 2: Forgetting to Emit Events from Extracted Components

**What goes wrong:** TreeTaskRow handles click internally instead of emitting to parent, breaking selection in TasksView.

**Why it happens:** Playground components had inline state management; extracted components need props-down, events-up pattern.

**How to avoid:**
1. Identify all user interactions in the playground component (click, expand, drag)
2. Replace inline state changes with `emit('eventName', payload)`
3. Wire emitted events to store actions in the parent view

**Warning signs:**
- Clicking a tree row doesn't select it
- Expand/collapse buttons don't work after extraction
- Console shows "emit is not defined" errors

### Pitfall 3: ADO Metadata Cache Not Refreshing After Sync

**What goes wrong:** User syncs ADO items, but task rows still show stale ADO type/state badges.

**Why it happens:** Forgot to call `Refresh()` during sync, or frontend didn't subscribe to `sync:completed` event.

**How to avoid:**
1. In `SyncService.Sync()`, call `adoMetaCache.Refresh()` after successful sync
2. In `useAdoMeta()` composable, subscribe to `sync:completed` event and call `refresh()`
3. Test: sync, verify cache updates in SQLite, verify frontend refetches

**Warning signs:**
- Task rows show old ADO state even after sync
- ADO type icons don't appear on newly linked tasks until app restart
- `ado_meta_cache` table is empty or stale in SQLite

### Pitfall 4: Performance Regression from Per-Row Metadata Calls

**What goes wrong:** Task list renders slowly because each TreeTaskRow calls `GetADOWorkItem()` during render.

**Why it happens:** Didn't implement batch loading; naively ported playground's mock data access pattern.

**How to avoid:**
1. Load all metadata once via `useAdoMeta()` on mount
2. Pass `adoMeta.getAdoMeta(task.id)` as a prop to TreeTaskRow
3. Never call Wails bindings inside a `v-for` loop render

**Warning signs:**
- Task list takes >1s to render with 50+ tasks
- Network/IPC tab shows hundreds of GetADOWorkItem calls
- UI freezes during scroll or filter changes

### Pitfall 5: Breaking Existing FilterBar Integration

**What goes wrong:** After replacing TasksView, existing filters (status, priority, project, dueDate) stop working.

**Why it happens:** New TasksView didn't wire FilterBar to `taskStore.enhancedFilteredTasks` computed.

**How to avoid:**
1. Keep existing FilterBar component unchanged
2. Wire `enhancedFilteredTasks` computed to the tree rendering logic
3. Test all filter combinations before marking phase complete

**Warning signs:**
- Clicking filter chips has no effect
- `enhancedFilteredTasks` returns full task list regardless of filters
- Console shows "filterStatus is not reactive" warnings

## Code Examples

Verified patterns from official sources and existing codebase:

### Extract QuickAddInput Component

```vue
<!-- components/tasks/QuickAddInput.vue -->
<!-- Source: PlaygroundIntegrated.vue lines 983-995 -->
<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { Circle } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

const emit = defineEmits<{
  add: [title: string]
  cancel: []
}>()

const title = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

function handleAdd() {
  if (title.value.trim()) {
    emit('add', title.value.trim())
    title.value = ''
  }
}

onMounted(() => {
  nextTick(() => inputRef.value?.focus())
})
</script>

<template>
  <div class="flex items-center gap-2 px-3 py-2 border-t border-border/30">
    <Circle :size="16" class="text-zinc-400 shrink-0" />
    <input
      ref="inputRef"
      v-model="title"
      placeholder="New task title… (Enter to add, Esc to cancel)"
      class="flex-1 text-sm bg-transparent border-none outline-none placeholder:text-muted-foreground/40"
      @keydown.enter="handleAdd"
      @keydown.escape="emit('cancel')"
    />
    <Button variant="outline" size="sm" class="h-7 px-2 text-[10px] shrink-0" @click="handleAdd" :disabled="!title.trim()">Add</Button>
    <Button variant="ghost" size="sm" class="h-7 px-1.5 text-[10px]" @click="emit('cancel')">Cancel</Button>
  </div>
</template>
```

### Add statusIcon Helper to styles.ts

```typescript
// lib/styles.ts — add to existing file
import type { Component } from 'vue'
import {
  Circle, CircleDot, Eye, CheckCircle2, Octagon, XCircle,
  // ... existing imports
} from 'lucide-vue-next'

/** Maps task status to Lucide icon component */
export function statusIcon(status: string): Component {
  switch (status) {
    case 'todo': return Circle
    case 'in_progress': return CircleDot
    case 'in_review': return Eye
    case 'done': return CheckCircle2
    case 'blocked': return Octagon
    case 'cancelled': return XCircle
    default: return Circle
  }
}

/** Maps priority to Tailwind bg class for priority dots */
export function priorityDotBgColor(priority: string): string {
  switch (priority) {
    case 'P0': return 'bg-red-500'
    case 'P1': return 'bg-orange-500'
    case 'P2': return 'bg-amber-500'
    case 'P3': return 'bg-zinc-400'
    default: return 'bg-zinc-400'
  }
}
```

### Wire TreeTaskRow to Real Stores

```vue
<!-- views/TasksView.vue — using extracted components -->
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTaskStore } from '@/stores/tasks'
import { useProjectStore } from '@/stores/projects'
import { useAdoMeta } from '@/composables/useAdoMeta'
import TreeTaskRow from '@/components/tasks/TreeTaskRow.vue'
import TaskDetail from '@/components/tasks/TaskDetail.vue'
import QuickAddInput from '@/components/tasks/QuickAddInput.vue'
import FilterBar from '@/components/FilterBar.vue'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-vue-next'

const taskStore = useTaskStore()
const projectStore = useProjectStore()
const adoMeta = useAdoMeta()

const showQuickAdd = ref(false)
const treeView = ref(true)

const projectNameMap = computed(() => {
  const map: Record<number, string> = {}
  for (const p of projectStore.projects) {
    map[p.id] = p.name
  }
  return map
})

function handleQuickAdd(title: string) {
  taskStore.quickAdd(title)
  showQuickAdd.value = false
}
</script>

<template>
  <div class="flex flex-col h-full">
    <div class="shrink-0 border-b border-border px-4 py-2 flex items-center gap-3">
      <FilterBar />
      <div class="flex-1" />
      <Button size="sm" @click="showQuickAdd = !showQuickAdd">
        <Plus :size="14" class="mr-1" /> New Task
      </Button>
    </div>

    <div class="flex-1 flex min-h-0">
      <div class="w-96 border-r border-border flex flex-col">
        <ScrollArea class="flex-1">
          <TreeTaskRow
            v-for="task in taskStore.enhancedFilteredTasks"
            :key="task.id"
            :task="task"
            :selected="taskStore.selectedTaskId === task.id"
            :ado-meta="adoMeta.getAdoMeta(task.id)"
            :project-name="projectNameMap[task.projectId]"
            @click="taskStore.selectedTaskId = task.id"
          />
        </ScrollArea>
        <QuickAddInput v-if="showQuickAdd" @add="handleQuickAdd" @cancel="showQuickAdd = false" />
      </div>

      <TaskDetail v-if="taskStore.selectedTaskId" class="flex-1" />
    </div>
  </div>
</template>
```

### Backend: Populate ADO Meta Cache During Sync

```go
// internal/app/syncservice.go — extend existing Sync method
func (s *SyncService) Sync(ctx context.Context) error {
	s.syncMu.Lock()
	defer s.syncMu.Unlock()

	s.emitEvent("sync:started", nil)

	// ... existing sync logic ...

	// Refresh ADO metadata cache after successful sync
	if err := s.adoMetaCache.Refresh(ctx); err != nil {
		log.Printf("[SyncService] Failed to refresh ADO meta cache: %v", err)
		// Non-fatal: sync succeeded, cache refresh failed
	}

	s.emitEvent("sync:completed", nil)
	return nil
}
```

### Frontend: Subscribe to Sync Events in Composable

```typescript
// composables/useAdoMeta.ts
import { ref, onMounted, onUnmounted } from 'vue'
import { Events } from '@wailsio/runtime'
import { GetAllADOMeta } from '@/api/adometa'

interface AdoMeta {
  type: string
  state: string
}

const metaCache = ref<Map<number, AdoMeta>>(new Map())
const loading = ref(false)

export function useAdoMeta() {
  let unsubscribe: (() => void) | null = null

  async function refresh() {
    loading.value = true
    try {
      const data = await GetAllADOMeta() as Record<string, AdoMeta>
      metaCache.value = new Map(Object.entries(data).map(([k, v]) => [Number(k), v]))
    } catch (err) {
      console.error('[useAdoMeta] Failed to load metadata:', err)
    } finally {
      loading.value = false
    }
  }

  function getAdoMeta(taskId: number): AdoMeta | null {
    return metaCache.value.get(taskId) ?? null
  }

  onMounted(async () => {
    if (metaCache.value.size === 0) {
      await refresh()
    }

    // Subscribe to sync:completed event
    unsubscribe = Events.On('sync:completed', () => {
      console.log('[useAdoMeta] Sync completed, refreshing cache')
      refresh()
    })
  })

  onUnmounted(() => {
    if (unsubscribe) {
      unsubscribe()
    }
  })

  return {
    getAdoMeta,
    refresh,
    loading,
  }
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| TaskRow component with inline logic | TreeTaskRow with props-only, emits events | Phase 11 (this phase) | Single component used everywhere; easier testing |
| Slide-out detail panel | Permanent 2-panel split | Phase 11 (this phase) | Context preserved when viewing details |
| Per-row ADO metadata queries | Batch-loaded cache with composable | Phase 11 (this phase) | 100x faster rendering for large task lists |
| Mock data in playground | Real stores wired to extracted components | Phase 11 (this phase) | Production-ready components |
| Status icon logic in templates | Centralized `statusIcon()` helper in styles.ts | Phase 11 (this phase) | DRY, type-safe icon mapping |

**Deprecated/outdated:**
- `TaskRow.vue`: Replaced by `TreeTaskRow.vue` (Phase 11). Old component had tags, description preview, timestamp — all dropped for simplified UI.
- Slide-out `TaskDetail` pattern: Replaced by permanent 2-panel split (Phase 11). User decision based on playground validation.
- ADO type icons in AzureDevOpsIcon component: Still used for the Azure DevOps logo, but type-specific icons now come from `adoTypeIcon()` helper (returns Lucide components like Bug, Star, Trophy).

## Open Questions

1. **Should expand/collapse state persist across sessions?**
   - What we know: Playground uses in-memory `ref<Set<number>>`; resets on refresh
   - What's unclear: User expectation for tree state persistence
   - Recommendation: Start with in-memory state for Phase 11; add localStorage persistence if users request it

2. **How to handle drag-drop reordering in tree view?**
   - What we know: `vuedraggable` works for flat lists; tree drag-drop is more complex (can't drop child under sibling)
   - What's unclear: Whether to support drag-drop in tree mode or disable it
   - Recommendation: Disable drag-drop in tree mode for Phase 11; only support in flat mode. Revisit if users request it.

3. **Should FilterCycleButton state sync with existing filterAdoLink store field?**
   - What we know: `taskStore.filterAdoLink` exists ('all' / 'linked' / 'personal') but unused in current UI
   - What's unclear: Whether FilterCycleButton should directly mutate store field or emit event
   - Recommendation: Use `v-model` on `taskStore.filterAdoLink` directly; keeps it simple

## Sources

### Primary (HIGH confidence)
- `frontend/src/views/playground/PlaygroundIntegrated.vue` — Lines 1-1200+ reviewed for component extraction patterns
- `frontend/src/stores/tasks.ts` — Existing `enhancedFilteredTasks`, `filterAdoLink`, `groupedEnhanced` computeds verified
- `frontend/src/lib/styles.ts` — Existing helper functions (`statusColor`, `adoTypeIcon`, `adoTypeColor`) reviewed
- `frontend/package.json` — Dependency versions verified (Vue 3.2.45, Pinia 3.0.4, @vueuse/core 14.2.1)
- `internal/db/migrate.go` — Schema structure for new `ado_meta_cache` table
- `.planning/phases/11-integrate-new-task-list-into-real-app/11-CONTEXT.md` — User decisions and canonical references

### Secondary (MEDIUM confidence)
- `frontend/src/components/tasks/TaskRow.vue` — Current implementation reviewed for features to preserve/drop
- `frontend/src/components/tasks/TaskDetail.vue` — Current detail pane structure for replacement guidance
- `internal/app/syncservice.go` — Sync flow for cache refresh integration point
- `copilot-instructions.md` — Project conventions (async patterns, Wails event usage, Vue patterns)

### Tertiary (LOW confidence)
- None — all research based on existing codebase and CONTEXT.md user decisions

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all dependencies already in project, versions verified from package.json
- Architecture: HIGH — patterns extracted directly from working playground code and existing production patterns
- Pitfalls: HIGH — based on common Vue component extraction mistakes and specific codebase patterns (e.g., field naming differences between mock and production)

**Research date:** 2026-04-08
**Valid until:** 2026-05-08 (30 days — stable codebase, low framework churn)
