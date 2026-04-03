<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import {
  GitBranch,
  AlertCircle,
  TreePine,
  List,
  Users,
} from 'lucide-vue-next'

const router = useRouter()

// --- Types ---
interface GraphNode {
  id: number
  title: string
  status: string
  priority: string
  description: string
  dependsOn: number[]
  children: GraphNode[]
}

// --- Mock Data ---
const rawTasks: Omit<GraphNode, 'children'>[] = [
  // Chain of 4: root blocked → cascade
  { id: 1, title: 'Provision staging database cluster', status: 'blocked', priority: 'P0', description: 'Set up PostgreSQL cluster for staging environment. Waiting on infrastructure team approval.', dependsOn: [] },
  { id: 2, title: 'Run database schema migrations', status: 'blocked', priority: 'P0', description: 'Apply v2.3 schema changes including new indexes and partitioning.', dependsOn: [1] },
  { id: 3, title: 'Deploy auth service to staging', status: 'blocked', priority: 'P1', description: 'Deploy authentication microservice with new OAuth2 provider support.', dependsOn: [2] },
  { id: 4, title: 'Run integration test suite', status: 'blocked', priority: 'P1', description: 'Execute full integration test suite against staging environment.', dependsOn: [3] },

  // Fan-out: one task blocks multiple
  { id: 5, title: 'Design API v3 specification', status: 'in_progress', priority: 'P1', description: 'Create OpenAPI spec for the v3 REST API including new endpoints for batch operations.', dependsOn: [] },
  { id: 6, title: 'Implement user endpoints', status: 'todo', priority: 'P2', description: 'Build CRUD endpoints for user management following v3 spec.', dependsOn: [5] },
  { id: 7, title: 'Implement project endpoints', status: 'todo', priority: 'P2', description: 'Build CRUD endpoints for project management following v3 spec.', dependsOn: [5] },
  { id: 8, title: 'Implement webhook endpoints', status: 'todo', priority: 'P3', description: 'Build webhook registration and delivery endpoints following v3 spec.', dependsOn: [5] },

  // Fan-in: one task depends on multiple parents
  { id: 9, title: 'Write API documentation', status: 'todo', priority: 'P2', description: 'Generate and review API documentation for all v3 endpoints.', dependsOn: [6, 7, 8] },

  // Independent tasks
  { id: 10, title: 'Update CI pipeline configuration', status: 'in_progress', priority: 'P2', description: 'Migrate from Jenkins to GitHub Actions for main build pipeline.', dependsOn: [] },
  { id: 11, title: 'Refactor logging middleware', status: 'done', priority: 'P3', description: 'Replace custom logger with structured logging using pino.', dependsOn: [] },

  // Another small chain
  { id: 12, title: 'Set up monitoring dashboards', status: 'done', priority: 'P2', description: 'Create Grafana dashboards for API latency, error rates, and throughput.', dependsOn: [] },
  { id: 13, title: 'Configure alerting rules', status: 'in_progress', priority: 'P1', description: 'Define PagerDuty alert thresholds based on dashboard metrics.', dependsOn: [12] },
  { id: 14, title: 'Run load test baseline', status: 'todo', priority: 'P2', description: 'Execute k6 load test to establish performance baseline before v3 rollout.', dependsOn: [13, 4] },
]

// --- Build tree ---
const allNodes = computed<GraphNode[]>(() => {
  const map = new Map<number, GraphNode>()
  for (const t of rawTasks) {
    map.set(t.id, { ...t, children: [] })
  }
  for (const node of map.values()) {
    for (const depId of node.dependsOn) {
      const parent = map.get(depId)
      if (parent) {
        parent.children.push(node)
      }
    }
  }
  return Array.from(map.values())
})

const rootNodes = computed(() => allNodes.value.filter(n => n.dependsOn.length === 0))

// --- View / Filter state ---
type ViewMode = 'tree' | 'flat'
const viewMode = ref<ViewMode>('tree')
const statusFilter = ref('all')
const expandedNodes = ref(new Set<number>([1, 2, 3, 5, 12, 13]))
const selectedNodeId = ref<number | null>(null)

// --- Chain highlighting ---
const highlightedIds = computed<Set<number>>(() => {
  if (selectedNodeId.value === null) return new Set()
  const ids = new Set<number>()
  const nodeMap = new Map(allNodes.value.map(n => [n.id, n]))

  function walkUp(id: number) {
    if (ids.has(id)) return
    ids.add(id)
    const node = nodeMap.get(id)
    if (node) {
      for (const depId of node.dependsOn) walkUp(depId)
    }
  }

  function walkDown(id: number) {
    if (ids.has(id)) return
    ids.add(id)
    const node = nodeMap.get(id)
    if (node) {
      for (const child of node.children) walkDown(child.id)
    }
  }

  walkUp(selectedNodeId.value)
  walkDown(selectedNodeId.value)
  return ids
})

// --- Dependent count ---
function getDependentCount(node: GraphNode): number {
  let count = node.children.length
  for (const child of node.children) {
    count += getDependentCount(child)
  }
  return count
}

// --- Filtered flat list ---
const flatNodes = computed(() => {
  let nodes = allNodes.value
  if (statusFilter.value !== 'all') {
    nodes = nodes.filter(n => n.status === statusFilter.value)
  }
  return [...nodes].sort((a, b) => {
    const priorityOrder: Record<string, number> = { P0: 0, P1: 1, P2: 2, P3: 3 }
    return (priorityOrder[a.priority] ?? 4) - (priorityOrder[b.priority] ?? 4)
  })
})

// --- Filtered root nodes ---
function nodeMatchesFilter(node: GraphNode): boolean {
  if (statusFilter.value === 'all') return true
  if (node.status === statusFilter.value) return true
  return node.children.some(c => nodeMatchesFilter(c))
}
const filteredRoots = computed(() => rootNodes.value.filter(n => nodeMatchesFilter(n)))

// --- Expand / Collapse ---
function toggleExpand(id: number) {
  const next = new Set(expandedNodes.value)
  if (next.has(id)) {
    next.delete(id)
  } else {
    next.add(id)
  }
  expandedNodes.value = next
}

function selectNode(id: number) {
  selectedNodeId.value = selectedNodeId.value === id ? null : id
}

function navigateToTask(id: number) {
  router.push({ name: 'tasks', query: { taskId: String(id) } })
}

// --- Stats ---
const stats = computed(() => {
  const all = allNodes.value
  return {
    total: all.length,
    blocked: all.filter(n => n.status === 'blocked').length,
    inProgress: all.filter(n => n.status === 'in_progress').length,
  }
})
</script>

<template>
  <div class="flex h-full flex-col overflow-hidden">
    <!-- Top Bar -->
    <div class="flex-none border-b border-border bg-background px-6 py-4">
      <div class="flex items-start justify-between gap-4">
        <div class="min-w-0">
          <div class="flex items-center gap-2">
            <GitBranch :size="20" class="text-muted-foreground" />
            <h1 class="text-lg font-semibold text-foreground">Dependency Graph</h1>
          </div>
          <p class="mt-0.5 text-sm text-muted-foreground">
            Trace blocked tasks and dependency chains
          </p>
        </div>

        <div class="flex items-center gap-3">
          <div class="hidden items-center gap-2 md:flex">
            <Badge variant="outline" class="gap-1 text-xs text-muted-foreground">
              {{ stats.total }} tasks
            </Badge>
            <Badge
              v-if="stats.blocked > 0"
              variant="outline"
              class="gap-1 border-red-500/20 bg-red-500/10 text-xs text-red-600 dark:text-red-400"
            >
              <AlertCircle :size="12" />
              {{ stats.blocked }} blocked
            </Badge>
          </div>

          <Select v-model="statusFilter">
            <SelectTrigger class="h-8 w-[140px] text-xs">
              <SelectValue placeholder="Filter status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All statuses</SelectItem>
              <SelectItem value="blocked">Blocked</SelectItem>
              <SelectItem value="in_progress">In Progress</SelectItem>
              <SelectItem value="todo">To Do</SelectItem>
              <SelectItem value="done">Done</SelectItem>
            </SelectContent>
          </Select>

          <div class="flex rounded-md border border-border">
            <Button
              variant="ghost"
              size="sm"
              :class="[
                'h-8 rounded-r-none px-2.5',
                viewMode === 'tree' && 'bg-accent text-accent-foreground',
              ]"
              @click="viewMode = 'tree'"
            >
              <TreePine :size="14" />
              <span class="ml-1 hidden text-xs sm:inline">Tree</span>
            </Button>
            <Button
              variant="ghost"
              size="sm"
              :class="[
                'h-8 rounded-l-none border-l border-border px-2.5',
                viewMode === 'flat' && 'bg-accent text-accent-foreground',
              ]"
              @click="viewMode = 'flat'"
            >
              <List :size="14" />
              <span class="ml-1 hidden text-xs sm:inline">List</span>
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main content -->
    <ScrollArea class="flex-1">
      <div class="p-6">
        <!-- Tree View -->
        <div v-if="viewMode === 'tree'" class="space-y-0.5">
          <div v-if="filteredRoots.length === 0" class="py-12 text-center text-sm text-muted-foreground">
            No tasks match the current filter.
          </div>
          <TreeNodeItem
            v-for="node in filteredRoots"
            :key="node.id"
            :node="node"
            :depth="0"
            :expanded-nodes="expandedNodes"
            :selected-node-id="selectedNodeId"
            :highlighted-ids="highlightedIds"
            :status-filter="statusFilter"
            @toggle-expand="toggleExpand"
            @select-node="selectNode"
            @navigate="navigateToTask"
          />
        </div>

        <!-- Flat List View -->
        <div v-else class="space-y-2">
          <div v-if="flatNodes.length === 0" class="py-12 text-center text-sm text-muted-foreground">
            No tasks match the current filter.
          </div>
          <TooltipProvider>
            <div
              v-for="node in flatNodes"
              :key="node.id"
              class="group flex items-center gap-3 rounded-lg border border-border bg-card p-3 transition-colors hover:bg-accent/50"
              :class="{
                'ring-2 ring-primary/50': selectedNodeId === node.id,
                'border-red-500/30 bg-red-500/5': node.status === 'blocked',
              }"
              @click="selectNode(node.id)"
            >
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2">
                  <button
                    class="truncate text-sm font-medium text-foreground hover:underline"
                    @click.stop="navigateToTask(node.id)"
                  >
                    {{ node.title }}
                  </button>
                </div>
                <p class="mt-0.5 truncate text-xs text-muted-foreground">
                  {{ node.description }}
                </p>
              </div>
              <div class="flex shrink-0 items-center gap-2">
                <StatusBadge :status="node.status" />
                <PriorityBadge :priority="node.priority" />
                <Tooltip v-if="node.children.length > 0">
                  <TooltipTrigger as-child>
                    <Badge variant="outline" class="gap-1 text-[10px]">
                      <Users :size="10" />
                      {{ getDependentCount(node) }}
                    </Badge>
                  </TooltipTrigger>
                  <TooltipContent side="top">
                    <p class="text-xs">{{ getDependentCount(node) }} dependent task(s)</p>
                  </TooltipContent>
                </Tooltip>
                <Badge
                  v-if="node.dependsOn.length > 0"
                  variant="outline"
                  class="gap-1 text-[10px] text-muted-foreground"
                >
                  {{ node.dependsOn.length }} dep{{ node.dependsOn.length > 1 ? 's' : '' }}
                </Badge>
              </div>
            </div>
          </TooltipProvider>
        </div>
      </div>
    </ScrollArea>
  </div>
</template>

<!-- TreeNodeItem: recursive component for tree nodes -->
<script lang="ts">
import { defineComponent, type PropType } from 'vue'

interface TreeGraphNode {
  id: number
  title: string
  status: string
  priority: string
  description: string
  dependsOn: number[]
  children: TreeGraphNode[]
}

function getDescendantCount(node: TreeGraphNode): number {
  let count = node.children.length
  for (const child of node.children) {
    count += getDescendantCount(child)
  }
  return count
}

function treeNodeMatchesFilter(node: TreeGraphNode, filter: string): boolean {
  if (filter === 'all') return true
  if (node.status === filter) return true
  return node.children.some(c => treeNodeMatchesFilter(c, filter))
}

const TreeNodeItem = defineComponent({
  name: 'TreeNodeItem',
  props: {
    node: { type: Object as PropType<TreeGraphNode>, required: true },
    depth: { type: Number, required: true },
    expandedNodes: { type: Object as PropType<Set<number>>, required: true },
    selectedNodeId: { type: Number as PropType<number | null>, default: null },
    highlightedIds: { type: Object as PropType<Set<number>>, required: true },
    statusFilter: { type: String, required: true },
  },
  emits: ['toggle-expand', 'select-node', 'navigate'],
  computed: {
    isExpanded(): boolean {
      return this.expandedNodes.has(this.node.id)
    },
    hasChildren(): boolean {
      return this.node.children.length > 0
    },
    isSelected(): boolean {
      return this.selectedNodeId === this.node.id
    },
    isHighlighted(): boolean {
      return this.highlightedIds.has(this.node.id)
    },
    isBlocked(): boolean {
      return this.node.status === 'blocked'
    },
    descendantCount(): number {
      return getDescendantCount(this.node)
    },
    filteredChildren(): TreeGraphNode[] {
      return this.node.children.filter(c => treeNodeMatchesFilter(c, this.statusFilter))
    },
    statusColor(): string {
      switch (this.node.status) {
        case 'blocked': return 'bg-red-500'
        case 'in_progress': return 'bg-blue-500'
        case 'done': return 'bg-emerald-500'
        case 'todo': return 'bg-zinc-400'
        default: return 'bg-zinc-400'
      }
    },
    statusLabel(): string {
      switch (this.node.status) {
        case 'blocked': return 'Blocked'
        case 'in_progress': return 'In Progress'
        case 'done': return 'Done'
        case 'todo': return 'To Do'
        default: return this.node.status
      }
    },
    priorityColor(): string {
      switch (this.node.priority) {
        case 'P0': return 'text-red-600 dark:text-red-400'
        case 'P1': return 'text-orange-600 dark:text-orange-400'
        case 'P2': return 'text-amber-600 dark:text-amber-400'
        case 'P3': return 'text-zinc-500'
        default: return 'text-zinc-500'
      }
    },
  },
  template: `
    <div>
      <div
        class="tree-node group relative flex items-center gap-2 rounded-md px-2 py-1.5 transition-colors hover:bg-accent/50"
        :class="{
          'ring-2 ring-primary/50 bg-primary/5': isSelected,
          'bg-yellow-500/5': isHighlighted && !isSelected,
          'border-l-2 border-l-red-500': isBlocked && depth === 0,
        }"
        :style="{ paddingLeft: (depth * 24 + 8) + 'px' }"
        @click="$emit('select-node', node.id)"
      >
        <!-- Connector lines for nested items -->
        <div
          v-if="depth > 0"
          class="connector-line absolute top-0 bottom-0"
          :style="{ left: ((depth - 1) * 24 + 19) + 'px', width: '24px' }"
        >
          <div class="absolute top-1/2 left-0 h-px w-4 bg-border"></div>
          <div class="absolute top-0 left-0 h-1/2 w-px bg-border"></div>
        </div>

        <!-- Expand/collapse chevron -->
        <button
          v-if="hasChildren"
          class="flex h-5 w-5 shrink-0 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground"
          @click.stop="$emit('toggle-expand', node.id)"
        >
          <svg v-if="isExpanded" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
        </button>
        <div v-else class="w-5 shrink-0"></div>

        <!-- Status dot -->
        <div class="flex h-2.5 w-2.5 shrink-0 rounded-full" :class="statusColor"></div>

        <!-- Task title -->
        <button
          class="min-w-0 flex-1 truncate text-left text-sm font-medium text-foreground hover:underline"
          :title="node.description"
          @click.stop="$emit('navigate', node.id)"
        >
          {{ node.title }}
        </button>

        <!-- Badges -->
        <div class="flex shrink-0 items-center gap-1.5 opacity-80 group-hover:opacity-100">
          <span
            class="rounded px-1.5 py-0.5 text-[10px] font-medium"
            :class="{
              'bg-red-500/15 text-red-600 dark:text-red-400': node.status === 'blocked',
              'bg-blue-500/15 text-blue-600 dark:text-blue-400': node.status === 'in_progress',
              'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400': node.status === 'done',
              'bg-zinc-500/15 text-zinc-500': node.status === 'todo',
            }"
          >{{ statusLabel }}</span>
          <span
            class="text-[10px] font-semibold"
            :class="priorityColor"
          >{{ node.priority }}</span>
          <span
            v-if="descendantCount > 0"
            class="rounded bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground"
            :title="descendantCount + ' dependent task(s)'"
          >{{ descendantCount }} dep{{ descendantCount > 1 ? 's' : '' }}</span>
        </div>
      </div>

      <!-- Children (recursive) -->
      <div v-if="hasChildren && isExpanded" class="relative">
        <div
          class="absolute top-0 bottom-2"
          :style="{ left: (depth * 24 + 19) + 'px' }"
        >
          <div class="h-full w-px bg-border"></div>
        </div>
        <TreeNodeItem
          v-for="child in filteredChildren"
          :key="child.id"
          :node="child"
          :depth="depth + 1"
          :expanded-nodes="expandedNodes"
          :selected-node-id="selectedNodeId"
          :highlighted-ids="highlightedIds"
          :status-filter="statusFilter"
          @toggle-expand="$emit('toggle-expand', $event)"
          @select-node="$emit('select-node', $event)"
          @navigate="$emit('navigate', $event)"
        />
      </div>
    </div>
  `,
})

export default {
  components: { TreeNodeItem },
}
</script>

<style scoped>
.tree-node {
  min-height: 36px;
}
</style>
