<script setup lang="ts">
/**
 * Playground: Task hierarchy tree rendering.
 *
 * Compares the current flat task list with a proposed hierarchical tree
 * (similar to how AdoTreeBrowser works). Tasks with parentId are shown
 * nested under their parent instead of flat.
 *
 * Tests:
 *  1. Hierarchical indentation with expand/collapse
 *  2. Deliverable → Task nesting (like ADO's Epic → Feature → Story)
 *  3. Side-by-side comparison: flat vs tree
 */
import { ref, computed, defineComponent, h } from 'vue'
import type { Task } from '@/types'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import {
  Circle, CircleDot, CheckCircle2, Octagon, Eye, XCircle,
  ChevronRight, ChevronDown, Folder, FileText, Layers,
} from 'lucide-vue-next'
import { statusColor } from '@/lib/styles'

// ── Mock task tree with hierarchy ──
const mockTasks: Task[] = [
  // Project: Xbox Platform
  {
    id: 100, title: 'Xbox Platform Services', description: 'Top-level deliverable for platform work',
    status: 'in_progress', priority: 'P0', category: '', projectId: 1, area: '',
    dueDate: '', adoId: '50001', tags: '', blockedReason: '', blockedBy: '',
    parentId: null, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-01T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 101, title: 'Implement auth token refresh', description: 'Auto-refresh tokens before expiry',
    status: 'in_progress', priority: 'P1', category: '', projectId: 1, area: '',
    dueDate: '2026-04-10', adoId: '50010', tags: 'auth', blockedReason: '', blockedBy: '',
    parentId: 100, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-02T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 102, title: 'Add rate limiting middleware', description: 'Prevent API abuse',
    status: 'todo', priority: 'P2', category: '', projectId: 1, area: '',
    dueDate: '', adoId: '', tags: 'backend', blockedReason: '', blockedBy: '',
    parentId: 100, personalPriority: '', sortOrder: 1,
    createdAt: '2026-04-03T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  {
    id: 103, title: 'Write unit tests for auth', description: '',
    status: 'done', priority: 'P1', category: '', projectId: 1, area: '',
    dueDate: '', adoId: '50011', tags: 'test', blockedReason: '', blockedBy: '',
    parentId: 101, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-02T10:00:00Z', updatedAt: '2026-04-05T10:00:00Z', completedAt: '2026-04-05T10:00:00Z',
  },
  {
    id: 104, title: 'Handle token expiry edge cases', description: '',
    status: 'in_progress', priority: 'P1', category: '', projectId: 1, area: '',
    dueDate: '', adoId: '', tags: '', blockedReason: '', blockedBy: '',
    parentId: 101, personalPriority: '', sortOrder: 1,
    createdAt: '2026-04-04T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  // Project: Dashboard
  {
    id: 200, title: 'Dashboard Redesign', description: 'Revamp the dashboard layout',
    status: 'in_progress', priority: 'P1', category: '', projectId: 2, area: '',
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
    id: 202, title: 'Implement PR summary widget', description: 'Show PR counts and status',
    status: 'blocked', priority: 'P1', category: '', projectId: 2, area: '',
    dueDate: '', adoId: '60010', tags: 'pr', blockedReason: 'Waiting on PR API endpoint', blockedBy: '',
    parentId: 200, personalPriority: '', sortOrder: 1,
    createdAt: '2026-04-03T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
  // Orphan task (no parent, no project)
  {
    id: 300, title: 'Fix CI pipeline timeout', description: 'Pipeline times out on large repos',
    status: 'todo', priority: 'P0', category: '', projectId: null, area: '',
    dueDate: '2026-04-08', adoId: '', tags: 'ci,urgent', blockedReason: '', blockedBy: '',
    parentId: null, personalPriority: '', sortOrder: 0,
    createdAt: '2026-04-06T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
  },
]

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

function getChildren(taskId: number): Task[] {
  return childrenOf.value[taskId] || []
}

function hasChildren(taskId: number): boolean {
  return (childrenOf.value[taskId]?.length ?? 0) > 0
}

function subtaskProgress(taskId: number) {
  const children = getChildren(taskId)
  if (children.length === 0) return null
  const done = children.filter(c => c.status === 'done').length
  return { done, total: children.length, pct: Math.round((done / children.length) * 100) }
}

const expandedNodes = ref<Set<number>>(new Set([100, 101, 200]))

function toggleExpand(id: number) {
  if (expandedNodes.value.has(id)) {
    expandedNodes.value.delete(id)
  } else {
    expandedNodes.value.add(id)
  }
}

const selectedId = ref<number | null>(null)

const statusIcon = (status: string) => {
  switch (status) {
    case 'in_progress': return CircleDot
    case 'in_review': return Eye
    case 'done': return CheckCircle2
    case 'blocked': return Octagon
    case 'cancelled': return XCircle
    default: return Circle
  }
}

function adoNumber(adoId: string): string {
  if (!adoId) return ''
  const match = adoId.match(/\d+/)
  return match ? `#${match[0]}` : adoId
}

const viewMode = ref<'flat' | 'tree' | 'split'>('split')
</script>

<template>
  <div class="h-screen flex flex-col bg-background text-foreground">
    <!-- Header -->
    <div class="h-10 shrink-0 border-b border-border flex items-center px-4 gap-3">
      <Layers :size="14" class="text-purple-500" />
      <span class="text-sm font-semibold">Playground: Task Hierarchy Tree</span>
      <Badge variant="secondary" class="text-[9px] h-4">Phase 7</Badge>
      <div class="flex-1" />
      <div class="flex gap-1">
        <Button :variant="viewMode === 'flat' ? 'default' : 'outline'" size="sm" class="h-7 text-[10px]" @click="viewMode = 'flat'">Flat List</Button>
        <Button :variant="viewMode === 'tree' ? 'default' : 'outline'" size="sm" class="h-7 text-[10px]" @click="viewMode = 'tree'">Tree View</Button>
        <Button :variant="viewMode === 'split' ? 'default' : 'outline'" size="sm" class="h-7 text-[10px]" @click="viewMode = 'split'">Side by Side</Button>
      </div>
      <Button variant="outline" size="sm" class="h-7 text-xs" @click="expandedNodes = new Set(mockTasks.map(t => t.id))">Expand All</Button>
      <Button variant="outline" size="sm" class="h-7 text-xs" @click="expandedNodes = new Set()">Collapse All</Button>
    </div>

    <div class="flex-1 flex min-h-0">
      <!-- Flat list view -->
      <div v-if="viewMode === 'flat' || viewMode === 'split'" class="flex-1 border-r border-border flex flex-col">
        <div class="px-4 py-2 border-b border-border/50 text-xs font-semibold text-muted-foreground">
          Current: Flat List ({{ mockTasks.length }} tasks, no nesting)
        </div>
        <ScrollArea class="flex-1 min-h-0">
          <div
            v-for="task in mockTasks"
            :key="task.id"
            class="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors border-b border-border/20"
            :class="selectedId === task.id && 'bg-primary/5'"
            @click="selectedId = task.id"
          >
            <button class="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full" :class="statusColor(task.status)">
              <component :is="statusIcon(task.status)" :size="16" :stroke-width="1.75" />
            </button>
            <span class="text-sm truncate flex-1" :class="task.status === 'done' ? 'line-through text-muted-foreground' : ''">
              {{ task.title }}
            </span>
            <Badge v-if="task.adoId" variant="outline" class="text-[9px] h-4 px-1 text-blue-500 border-blue-500/30 shrink-0">
              {{ adoNumber(task.adoId) }}
            </Badge>
            <Badge variant="secondary" class="text-[9px] h-4 px-1 shrink-0">{{ task.priority }}</Badge>
          </div>
        </ScrollArea>
      </div>

      <!-- Tree view -->
      <div v-if="viewMode === 'tree' || viewMode === 'split'" class="flex-1 flex flex-col">
        <div class="px-4 py-2 border-b border-border/50 text-xs font-semibold text-muted-foreground">
          Proposed: Hierarchical Tree ({{ rootTasks.length }} roots, nested children)
        </div>
        <ScrollArea class="flex-1 min-h-0">
          <!-- Recursive tree rendering -->
          <template v-for="task in rootTasks" :key="task.id">
            <!-- Root node -->
            <div
              class="flex items-center gap-2 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
              :class="selectedId === task.id && 'bg-primary/5'"
              :style="{ paddingLeft: '16px' }"
              @click="selectedId = task.id"
            >
              <!-- Expand/collapse -->
              <button
                v-if="hasChildren(task.id)"
                class="shrink-0 p-0.5 hover:bg-muted rounded"
                @click.stop="toggleExpand(task.id)"
              >
                <component :is="expandedNodes.has(task.id) ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
              </button>
              <span v-else class="w-5 shrink-0" />

              <!-- Status -->
              <button class="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full" :class="statusColor(task.status)">
                <component :is="statusIcon(task.status)" :size="16" :stroke-width="1.75" />
              </button>

              <!-- Content -->
              <span class="text-sm font-medium truncate flex-1" :class="task.status === 'done' ? 'line-through text-muted-foreground' : ''">
                {{ task.title }}
              </span>

              <!-- Subtask progress -->
              <template v-if="subtaskProgress(task.id)">
                <span class="text-[10px] text-muted-foreground tabular-nums shrink-0">
                  {{ subtaskProgress(task.id)!.done }}/{{ subtaskProgress(task.id)!.total }}
                </span>
                <div class="w-12 h-1 bg-muted rounded-full overflow-hidden shrink-0">
                  <div class="h-full bg-green-500 rounded-full" :style="{ width: subtaskProgress(task.id)!.pct + '%' }" />
                </div>
              </template>

              <!-- Badges -->
              <Badge v-if="task.adoId" variant="outline" class="text-[9px] h-4 px-1 text-blue-500 border-blue-500/30 shrink-0">
                <AzureDevOpsIcon :size="9" class="mr-0.5" />
                {{ adoNumber(task.adoId) }}
              </Badge>
              <Badge variant="secondary" class="text-[9px] h-4 px-1 shrink-0">{{ task.priority }}</Badge>
            </div>

            <!-- Children (level 1) -->
            <template v-if="expandedNodes.has(task.id)">
              <template v-for="child in getChildren(task.id)" :key="child.id">
                <div
                  class="flex items-center gap-2 py-2 cursor-pointer hover:bg-muted/50 transition-colors border-l-2 border-l-muted-foreground/10"
                  :class="selectedId === child.id && 'bg-primary/5'"
                  :style="{ paddingLeft: '40px', paddingRight: '16px' }"
                  @click="selectedId = child.id"
                >
                  <button
                    v-if="hasChildren(child.id)"
                    class="shrink-0 p-0.5 hover:bg-muted rounded"
                    @click.stop="toggleExpand(child.id)"
                  >
                    <component :is="expandedNodes.has(child.id) ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
                  </button>
                  <span v-else class="w-5 shrink-0" />

                  <button class="flex-shrink-0 w-4 h-4 flex items-center justify-center rounded-full" :class="statusColor(child.status)">
                    <component :is="statusIcon(child.status)" :size="14" :stroke-width="1.75" />
                  </button>

                  <span class="text-sm truncate flex-1" :class="child.status === 'done' ? 'line-through text-muted-foreground' : ''">
                    {{ child.title }}
                  </span>

                  <template v-if="subtaskProgress(child.id)">
                    <span class="text-[10px] text-muted-foreground tabular-nums shrink-0">
                      {{ subtaskProgress(child.id)!.done }}/{{ subtaskProgress(child.id)!.total }}
                    </span>
                  </template>

                  <Badge v-if="child.adoId" variant="outline" class="text-[9px] h-4 px-1 text-blue-500 border-blue-500/30 shrink-0">
                    {{ adoNumber(child.adoId) }}
                  </Badge>
                  <Badge v-if="child.status === 'blocked'" variant="destructive" class="text-[9px] h-4 px-1 shrink-0">blocked</Badge>
                </div>

                <!-- Grandchildren (level 2) -->
                <template v-if="expandedNodes.has(child.id)">
                  <div
                    v-for="grandchild in getChildren(child.id)"
                    :key="grandchild.id"
                    class="flex items-center gap-2 py-1.5 cursor-pointer hover:bg-muted/50 transition-colors border-l-2 border-l-muted-foreground/5"
                    :class="selectedId === grandchild.id && 'bg-primary/5'"
                    :style="{ paddingLeft: '64px', paddingRight: '16px' }"
                    @click="selectedId = grandchild.id"
                  >
                    <span class="w-5 shrink-0" />
                    <button class="flex-shrink-0 w-4 h-4 flex items-center justify-center rounded-full" :class="statusColor(grandchild.status)">
                      <component :is="statusIcon(grandchild.status)" :size="12" :stroke-width="1.75" />
                    </button>
                    <span class="text-xs truncate flex-1" :class="grandchild.status === 'done' ? 'line-through text-muted-foreground' : ''">
                      {{ grandchild.title }}
                    </span>
                    <Badge v-if="grandchild.adoId" variant="outline" class="text-[9px] h-4 px-1 text-blue-500 border-blue-500/30 shrink-0">
                      {{ adoNumber(grandchild.adoId) }}
                    </Badge>
                  </div>
                </template>
              </template>
            </template>
          </template>
        </ScrollArea>
      </div>
    </div>

    <!-- Legend -->
    <div class="shrink-0 border-t border-border px-4 py-2 flex items-center gap-4 text-[10px] text-muted-foreground">
      <span>Hierarchy: Deliverable → Task → Subtask (3 levels)</span>
      <span>|</span>
      <span>Root tasks: {{ rootTasks.length }}</span>
      <span>|</span>
      <span>Total tasks: {{ mockTasks.length }}</span>
      <span>|</span>
      <span>Children expand/collapse with chevron</span>
    </div>
  </div>
</template>
