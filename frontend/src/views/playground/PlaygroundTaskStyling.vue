<script setup lang="ts">
/**
 * Playground: ADO-styled Task Tree List
 *
 * Combines the tree hierarchy (expand/collapse, 3-level nesting) with
 * ADO-style row decorations (type icons, colored state badges, 2-line layout,
 * progress bars, personal-item indicator, star/pin toggle).
 *
 * Side-by-side: current flat TaskRow vs proposed ADO-styled tree.
 */
import { ref, computed } from 'vue'
import type { Task } from '@/types'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import TaskRow from '@/components/tasks/TaskRow.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import {
  Circle, CircleDot, CheckCircle2, Octagon, Eye, XCircle,
  ChevronRight, ChevronDown, Paintbrush, User,
  CalendarDays,
} from 'lucide-vue-next'
import {
  statusColor, statusClasses,
  adoTypeColor, adoTypeIcon, adoStateClasses,
} from '@/lib/styles'

// ADO type metadata per task (simulates what projectAdoLookup would provide)
interface AdoMeta { type: string; state: string; adoId: string }
const adoMeta: Record<number, AdoMeta> = {
  100: { type: 'Deliverable', state: 'Active', adoId: '50001' },
  101: { type: 'User Story', state: 'Active', adoId: '50010' },
  102: { type: 'Task', state: 'New', adoId: '' },
  103: { type: 'Task', state: 'Closed', adoId: '50011' },
  104: { type: 'Task', state: 'Active', adoId: '' },
  200: { type: 'Feature', state: 'Active', adoId: '60001' },
  201: { type: 'Task', state: 'Closed', adoId: '' },
  202: { type: 'Bug', state: 'Active', adoId: '60010' },
  300: { type: '', state: '', adoId: '' },
  301: { type: '', state: '', adoId: '' },
}

// ── Mock task tree ──
const mockTasks: Task[] = [
  {
    id: 100, title: 'Xbox Platform Services', description: 'Top-level deliverable for platform work',
    status: 'in_progress', priority: 'P0', category: '', projectId: 1, area: 'Platform',
    dueDate: '2026-04-20', adoId: '50001', tags: '', blockedReason: '', blockedBy: '',
    parentId: null, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-01T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 101, title: 'Implement auth token refresh', description: 'Auto-refresh tokens before expiry',
    status: 'in_progress', priority: 'P1', category: '', projectId: 1, area: 'Auth',
    dueDate: '2026-04-10', adoId: '50010', tags: 'auth', blockedReason: '', blockedBy: '',
    parentId: 100, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-02T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 102, title: 'Add rate limiting middleware', description: 'Prevent API abuse with token bucket',
    status: 'todo', priority: 'P2', category: '', projectId: 1, area: 'Backend',
    dueDate: '', adoId: '', tags: 'backend', blockedReason: '', blockedBy: '',
    parentId: 100, personalPriority: '', sortOrder: 1,
    createdAt: '2026-04-03T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 103, title: 'Write unit tests for auth module', description: '',
    status: 'done', priority: 'P1', category: '', projectId: 1, area: '',
    dueDate: '', adoId: '50011', tags: 'test', blockedReason: '', blockedBy: '',
    parentId: 101, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-02T10:00:00Z', updatedAt: '2026-04-05T10:00:00Z', completedAt: '2026-04-05T10:00:00Z',
  },
  {
    id: 104, title: 'Handle token expiry edge cases', description: 'Graceful fallback when refresh fails',
    status: 'in_progress', priority: 'P1', category: '', projectId: 1, area: '',
    dueDate: '', adoId: '', tags: '', blockedReason: '', blockedBy: '',
    parentId: 101, personalPriority: '', sortOrder: 1,
    createdAt: '2026-04-04T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 200, title: 'Dashboard Redesign', description: 'Revamp the dashboard layout and widgets',
    status: 'in_progress', priority: 'P1', category: '', projectId: 2, area: 'UI',
    dueDate: '2026-04-15', adoId: '60001', tags: 'frontend', blockedReason: '', blockedBy: '',
    parentId: null, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-01T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 201, title: 'Add stat cards row', description: '',
    status: 'done', priority: 'P2', category: '', projectId: 2, area: '',
    dueDate: '', adoId: '', tags: '', blockedReason: '', blockedBy: '',
    parentId: 200, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-02T10:00:00Z', updatedAt: '2026-04-04T10:00:00Z', completedAt: '2026-04-04T10:00:00Z',
  },
  {
    id: 202, title: 'Implement PR summary widget', description: 'Show PR counts and status breakdown',
    status: 'blocked', priority: 'P1', category: '', projectId: 2, area: 'PRs',
    dueDate: '', adoId: '60010', tags: 'pr', blockedReason: 'Waiting on PR API endpoint', blockedBy: '',
    parentId: 200, personalPriority: '', sortOrder: 1,
    createdAt: '2026-04-03T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  // Personal tasks (no adoId, no project)
  {
    id: 300, title: 'Fix CI pipeline timeout', description: 'Pipeline times out on large repos',
    status: 'todo', priority: 'P0', category: '', projectId: null, area: '',
    dueDate: '2026-04-08', adoId: '', tags: 'ci,urgent', blockedReason: '', blockedBy: '',
    parentId: null, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-06T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 301, title: 'Organize notes and docs', description: 'Clean up personal task backlog',
    status: 'todo', priority: 'P3', category: '', projectId: null, area: '',
    dueDate: '', adoId: '', tags: '', blockedReason: '', blockedBy: '',
    parentId: null, personalPriority: '', sortOrder: 1,
    createdAt: '2026-04-06T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
]

const projectNames: Record<number, string> = { 1: 'Xbox Platform', 2: 'Dashboard' }

// ── Hierarchy helpers ──
const childrenOf = computed(() => {
  const map: Record<number, Task[]> = {}
  for (const t of mockTasks) {
    if (t.parentId) {
      if (!map[t.parentId]) map[t.parentId] = []
      map[t.parentId].push(t)
    }
  }
  return map
})

const rootTasks = computed(() => mockTasks.filter(t => !t.parentId))

function getChildren(id: number): Task[] { return childrenOf.value[id] || [] }
function hasChildren(id: number): boolean { return (childrenOf.value[id]?.length ?? 0) > 0 }

function subtaskProgress(id: number) {
  const ch = getChildren(id)
  if (!ch.length) return null
  const done = ch.filter(c => c.status === 'done').length
  return { done, total: ch.length, pct: Math.round((done / ch.length) * 100) }
}

const expandedNodes = ref<Set<number>>(new Set([100, 101, 200]))
function toggleExpand(id: number) {
  if (expandedNodes.value.has(id)) expandedNodes.value.delete(id)
  else expandedNodes.value.add(id)
}

const selectedId= ref<number | null>(null)
const activeVariant = ref<'current' | 'proposed' | 'split'>('split')

// ── Helpers ──
function meta(id: number): AdoMeta { return adoMeta[id] || { type: '', state: '', adoId: '' } }
function isPersonal(task: Task): boolean { return !task.adoId && !task.projectId }
function adoNumber(adoId: string): string {
  if (!adoId) return ''
  const m = adoId.match(/\d+/)
  return m ? `#${m[0]}` : adoId
}

// Map task status → display label
function statusLabel(s: string): string {
  switch (s) {
    case 'in_progress': return 'In Progress'
    case 'in_review': return 'In Review'
    case 'todo': return 'To Do'
    case 'done': return 'Done'
    case 'blocked': return 'Blocked'
    case 'cancelled': return 'Cancelled'
    default: return s
  }
}

const statusIcon = (s: string) => {
  switch (s) {
    case 'in_progress': return CircleDot
    case 'in_review': return Eye
    case 'done': return CheckCircle2
    case 'blocked': return Octagon
    case 'cancelled': return XCircle
    default: return Circle
  }
}
</script>

<template>
  <div class="h-screen flex flex-col bg-background text-foreground">
    <!-- Header -->
    <div class="h-10 shrink-0 border-b border-border flex items-center px-4 gap-3">
      <Paintbrush :size="14" class="text-orange-500" />
      <span class="text-sm font-semibold">Playground: ADO-Styled Task Tree</span>
      <Badge variant="secondary" class="text-[9px] h-4">Phase 8</Badge>
      <div class="flex-1" />
      <div class="flex gap-1">
        <Button :variant="activeVariant === 'current' ? 'default' : 'outline'" size="sm" class="h-7 text-[10px]" @click="activeVariant = 'current'">Current (Flat)</Button>
        <Button :variant="activeVariant === 'proposed' ? 'default' : 'outline'" size="sm" class="h-7 text-[10px]" @click="activeVariant = 'proposed'">Proposed (ADO Tree)</Button>
        <Button :variant="activeVariant === 'split' ? 'default' : 'outline'" size="sm" class="h-7 text-[10px]" @click="activeVariant = 'split'">Side by Side</Button>
      </div>
      <Button variant="outline" size="sm" class="h-7 text-xs" @click="expandedNodes = new Set(mockTasks.map(t => t.id))">Expand All</Button>
      <Button variant="outline" size="sm" class="h-7 text-xs" @click="expandedNodes = new Set()">Collapse All</Button>
    </div>

    <div class="flex-1 flex min-h-0">
      <!-- ════ LEFT: Current flat TaskRow rendering ════ -->
      <div v-if="activeVariant === 'current' || activeVariant === 'split'" class="flex-1 border-r border-border flex flex-col">
        <div class="px-4 py-2 border-b border-border/50 text-xs font-semibold text-muted-foreground">
          Current: Flat TaskRow ({{ mockTasks.length }} items, no hierarchy)
        </div>
        <ScrollArea class="flex-1 min-h-0">
          <div class="border-b border-border/20" v-for="task in mockTasks" :key="task.id">
            <TaskRow
              :task="task"
              :is-public="!!task.adoId"
              :selected="selectedId === task.id"
              :project-name="task.projectId ? projectNames[task.projectId] : undefined"
              @select="selectedId = $event"
            />
          </div>
        </ScrollArea>
      </div>

      <!-- ════ RIGHT: Proposed ADO-styled tree ════ -->
      <div v-if="activeVariant === 'proposed' || activeVariant === 'split'" class="flex-1 flex flex-col">
        <div class="px-4 py-2 border-b border-border/50 text-xs font-semibold text-muted-foreground">
          Proposed: ADO-styled Tree ({{ rootTasks.length }} roots, nested, rich rows)
        </div>
        <ScrollArea class="flex-1 min-h-0">
          <template v-for="task in rootTasks" :key="task.id">
            <!-- ── Root node ── -->
            <div
              class="group cursor-pointer hover:bg-muted/50 transition-colors"
              :class="selectedId === task.id ? 'bg-primary/5 border-l-2 border-l-primary' : 'border-l-2 border-l-transparent'"
              @click="selectedId = task.id"
            >
              <!-- Row 1: chevron + type icon + title + badges -->
              <div class="flex items-center gap-2 px-3 pt-2.5 pb-0.5">
                <!-- Expand/collapse chevron -->
                <button
                  v-if="hasChildren(task.id)"
                  class="shrink-0 p-0.5 hover:bg-muted rounded"
                  @click.stop="toggleExpand(task.id)"
                >
                  <component :is="expandedNodes.has(task.id) ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
                </button>
                <span v-else class="w-5 shrink-0" />

                <!-- Type icon or personal indicator -->
                <component
                  v-if="meta(task.id).type"
                  :is="adoTypeIcon(meta(task.id).type)"
                  :size="15"
                  :class="adoTypeColor(meta(task.id).type)"
                  class="shrink-0"
                />
                <User v-else-if="isPersonal(task)" :size="14" class="text-muted-foreground/60 shrink-0" />
                <component v-else :is="statusIcon(task.status)" :size="15" :class="statusColor(task.status)" class="shrink-0" />

                <!-- Title -->
                <span
                  class="text-sm truncate flex-1 min-w-0"
                  :class="[
                    task.status === 'done' ? 'line-through text-muted-foreground' : 'text-foreground',
                    hasChildren(task.id) ? 'font-medium' : '',
                  ]"
                >
                  {{ task.title }}
                </span>

                <!-- State badge (ADO state or task status) -->
                <Badge
                  variant="outline"
                  class="text-[10px] h-4 px-1.5 shrink-0"
                  :class="meta(task.id).state ? adoStateClasses(meta(task.id).state) : statusClasses(task.status)"
                >
                  {{ meta(task.id).state || statusLabel(task.status) }}
                </Badge>

                <!-- Subtask progress -->
                <template v-if="subtaskProgress(task.id)">
                  <span class="text-[10px] text-muted-foreground tabular-nums shrink-0">
                    {{ subtaskProgress(task.id)!.done }}/{{ subtaskProgress(task.id)!.total }}
                  </span>
                  <div class="w-12 h-1 bg-muted rounded-full overflow-hidden shrink-0">
                    <div class="h-full bg-green-500 rounded-full transition-all" :style="{ width: subtaskProgress(task.id)!.pct + '%' }" />
                  </div>
                </template>
              </div>

              <!-- Row 2: metadata (ADO ID, project, area, priority, due date) -->
              <div class="flex items-center gap-2 pb-2 text-[11px]" :style="{ paddingLeft: '68px', paddingRight: '12px' }">
                <span v-if="meta(task.id).adoId" class="text-muted-foreground/50 tabular-nums">{{ adoNumber(meta(task.id).adoId) }}</span>
                <Badge v-if="isPersonal(task)" variant="outline" class="text-[9px] h-3.5 px-1 border-dashed text-muted-foreground/60">
                  <User :size="8" class="mr-0.5" /> Personal
                </Badge>
                <span v-if="task.projectId && projectNames[task.projectId]" class="text-muted-foreground/40 truncate">
                  {{ projectNames[task.projectId] }}
                </span>
                <span v-if="task.area" class="text-muted-foreground/30">{{ task.area }}</span>
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

            <!-- ── Children (level 1) ── -->
            <template v-if="expandedNodes.has(task.id)">
              <template v-for="child in getChildren(task.id)" :key="child.id">
                <div
                  class="group cursor-pointer hover:bg-muted/50 transition-colors border-l-2"
                  :class="[
                    selectedId === child.id ? 'bg-primary/5 border-l-primary' : 'border-l-muted-foreground/10',
                  ]"
                  @click="selectedId = child.id"
                >
                  <div class="flex items-center gap-2 pt-2 pb-0.5" :style="{ paddingLeft: '36px', paddingRight: '12px' }">
                    <button
                      v-if="hasChildren(child.id)"
                      class="shrink-0 p-0.5 hover:bg-muted rounded"
                      @click.stop="toggleExpand(child.id)"
                    >
                      <component :is="expandedNodes.has(child.id) ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
                    </button>
                    <span v-else class="w-5 shrink-0" />

                    <component
                      v-if="meta(child.id).type"
                      :is="adoTypeIcon(meta(child.id).type)"
                      :size="14"
                      :class="adoTypeColor(meta(child.id).type)"
                      class="shrink-0"
                    />
                    <User v-else-if="isPersonal(child)" :size="13" class="text-muted-foreground/60 shrink-0" />
                    <component v-else :is="statusIcon(child.status)" :size="14" :class="statusColor(child.status)" class="shrink-0" />

                    <span class="text-sm truncate flex-1 min-w-0" :class="child.status === 'done' ? 'line-through text-muted-foreground' : 'text-foreground'">
                      {{ child.title }}
                    </span>

                    <Badge
                      variant="outline"
                      class="text-[10px] h-4 px-1.5 shrink-0"
                      :class="meta(child.id).state ? adoStateClasses(meta(child.id).state) : statusClasses(child.status)"
                    >
                      {{ meta(child.id).state || statusLabel(child.status) }}
                    </Badge>

                    <template v-if="subtaskProgress(child.id)">
                      <span class="text-[10px] text-muted-foreground tabular-nums shrink-0">
                        {{ subtaskProgress(child.id)!.done }}/{{ subtaskProgress(child.id)!.total }}
                      </span>
                    </template>
                  </div>

                  <!-- Child metadata row -->
                  <div class="flex items-center gap-2 pb-1.5 text-[10px]" :style="{ paddingLeft: '80px', paddingRight: '12px' }">
                    <span v-if="meta(child.id).adoId" class="text-muted-foreground/50 tabular-nums">{{ adoNumber(meta(child.id).adoId) }}</span>
                    <Badge v-if="isPersonal(child)" variant="outline" class="text-[9px] h-3.5 px-1 border-dashed text-muted-foreground/60">
                      <User :size="8" class="mr-0.5" /> Personal
                    </Badge>
                    <span v-if="child.area" class="text-muted-foreground/30">{{ child.area }}</span>
                    <div class="flex-1" />
                    <span
                      v-if="child.dueDate"
                      class="inline-flex items-center gap-0.5 text-[10px] shrink-0"
                      :class="new Date(child.dueDate) < new Date() ? 'text-red-500' : 'text-muted-foreground/50'"
                    >
                      <CalendarDays :size="9" />
                      {{ child.dueDate }}
                    </span>
                    <PriorityBadge :priority="child.priority" />
                  </div>

                  <div
                    v-if="child.status === 'blocked' && child.blockedReason"
                    class="pb-1.5 text-[10px] text-red-500/80 truncate"
                    :style="{ paddingLeft: '80px', paddingRight: '12px' }"
                  >
                    ⚠ {{ child.blockedReason }}
                  </div>
                </div>

                <!-- ── Grandchildren (level 2) ── -->
                <template v-if="expandedNodes.has(child.id)">
                  <div
                    v-for="gc in getChildren(child.id)"
                    :key="gc.id"
                    class="group flex items-center gap-2 py-1.5 cursor-pointer hover:bg-muted/50 transition-colors border-l-2"
                    :class="selectedId === gc.id ? 'bg-primary/5 border-l-primary' : 'border-l-muted-foreground/5'"
                    :style="{ paddingLeft: '60px', paddingRight: '12px' }"
                    @click="selectedId = gc.id"
                  >
                    <span class="w-5 shrink-0" />

                    <component
                      v-if="meta(gc.id).type"
                      :is="adoTypeIcon(meta(gc.id).type)"
                      :size="12"
                      :class="adoTypeColor(meta(gc.id).type)"
                      class="shrink-0"
                    />
                    <User v-else-if="isPersonal(gc)" :size="12" class="text-muted-foreground/60 shrink-0" />
                    <component v-else :is="statusIcon(gc.status)" :size="12" :class="statusColor(gc.status)" class="shrink-0" />

                    <span class="text-xs truncate flex-1 min-w-0" :class="gc.status === 'done' ? 'line-through text-muted-foreground' : 'text-foreground'">
                      {{ gc.title }}
                    </span>

                    <Badge
                      variant="outline"
                      class="text-[9px] h-3.5 px-1 shrink-0"
                      :class="meta(gc.id).state ? adoStateClasses(meta(gc.id).state) : statusClasses(gc.status)"
                    >
                      {{ meta(gc.id).state || statusLabel(gc.status) }}
                    </Badge>

                    <span v-if="meta(gc.id).adoId" class="text-[9px] text-muted-foreground/40 tabular-nums shrink-0">
                      {{ adoNumber(meta(gc.id).adoId) }}
                    </span>
                  </div>
                </template>
              </template>
            </template>
          </template>
        </ScrollArea>
      </div>
    </div>

    <!-- Legend -->
    <div class="shrink-0 border-t border-border px-4 py-2 text-[10px] text-muted-foreground flex gap-4">
      <span><strong>Current:</strong> Flat TaskRow, 1-line, status circles, inline badges</span>
      <span>|</span>
      <span><strong>Proposed:</strong> Hierarchical tree + ADO type icons, state badges, 2-line metadata, personal indicator, progress bars</span>
      <span>|</span>
      <span>{{ rootTasks.length }} roots · {{ mockTasks.length }} total</span>
    </div>
  </div>
</template>
