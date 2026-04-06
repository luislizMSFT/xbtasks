<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
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
import TreeNodeItem from '@/components/TreeNodeItem.vue'
import {
  GitBranch,
  AlertCircle,
  TreePine,
  List,
  Users,
  RefreshCw,
  Loader2,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()

// --- Types ---
interface SubtaskItem {
  id: number
  title: string
  done: boolean
}

interface GraphNode {
  id: number
  title: string
  status: string
  priority: string
  description: string
  dependsOn: number[]
  children: GraphNode[]
  subtasks?: SubtaskItem[]
}

// --- Mock fallback data ---
function buildMockNodes(): Omit<GraphNode, 'children'>[] {
  return [
    { id: 1, title: 'Provision staging database cluster', status: 'blocked', priority: 'P0', description: 'Set up PostgreSQL cluster for staging environment. Waiting on infrastructure team approval.', dependsOn: [], subtasks: [
      { id: 101, title: 'Request infra approval', done: true },
      { id: 102, title: 'Configure cluster settings', done: false },
      { id: 103, title: 'Set up connection pooling', done: false },
    ] },
    { id: 2, title: 'Run database schema migrations', status: 'blocked', priority: 'P0', description: 'Apply v2.3 schema changes including new indexes and partitioning.', dependsOn: [1] },
    { id: 3, title: 'Deploy auth service to staging', status: 'blocked', priority: 'P1', description: 'Deploy authentication microservice with new OAuth2 provider support.', dependsOn: [2], subtasks: [
      { id: 301, title: 'Update OAuth2 client config', done: true },
      { id: 302, title: 'Rotate service credentials', done: false },
    ] },
    { id: 4, title: 'Run integration test suite', status: 'blocked', priority: 'P1', description: 'Execute full integration test suite against staging environment.', dependsOn: [3] },
    { id: 5, title: 'Design API v3 specification', status: 'in_progress', priority: 'P1', description: 'Create OpenAPI spec for the v3 REST API including new endpoints for batch operations.', dependsOn: [], subtasks: [
      { id: 501, title: 'Draft endpoint inventory', done: true },
      { id: 502, title: 'Define request/response schemas', done: true },
      { id: 503, title: 'Review with backend team', done: false },
      { id: 504, title: 'Publish to API portal', done: false },
    ] },
    { id: 6, title: 'Implement user endpoints', status: 'todo', priority: 'P2', description: 'Build CRUD endpoints for user management following v3 spec.', dependsOn: [5] },
    { id: 7, title: 'Implement project endpoints', status: 'todo', priority: 'P2', description: 'Build CRUD endpoints for project management following v3 spec.', dependsOn: [5] },
    { id: 8, title: 'Implement webhook endpoints', status: 'todo', priority: 'P3', description: 'Build webhook registration and delivery endpoints following v3 spec.', dependsOn: [5] },
    { id: 9, title: 'Write API documentation', status: 'todo', priority: 'P2', description: 'Generate and review API documentation for all v3 endpoints.', dependsOn: [6, 7, 8] },
    { id: 10, title: 'Update CI pipeline configuration', status: 'in_progress', priority: 'P2', description: 'Migrate from Jenkins to GitHub Actions for main build pipeline.', dependsOn: [], subtasks: [
      { id: 1001, title: 'Create GitHub Actions workflow file', done: true },
      { id: 1002, title: 'Migrate build secrets', done: false },
      { id: 1003, title: 'Decommission Jenkins jobs', done: false },
    ] },
    { id: 11, title: 'Refactor logging middleware', status: 'done', priority: 'P3', description: 'Replace custom logger with structured logging using pino.', dependsOn: [] },
    { id: 12, title: 'Set up monitoring dashboards', status: 'done', priority: 'P2', description: 'Create Grafana dashboards for API latency, error rates, and throughput.', dependsOn: [] },
    { id: 13, title: 'Configure alerting rules', status: 'in_progress', priority: 'P1', description: 'Define PagerDuty alert thresholds based on dashboard metrics.', dependsOn: [12] },
    { id: 14, title: 'Run load test baseline', status: 'todo', priority: 'P2', description: 'Execute k6 load test to establish performance baseline before v3 rollout.', dependsOn: [13, 4] },
  ]
}

// --- Real data loading ---
const graphNodes = ref<Omit<GraphNode, 'children'>[]>([])
const loadingGraph = ref(false)

async function loadGraph() {
  loadingGraph.value = true
  try {
    // Ensure tasks are loaded
    if (taskStore.tasks.length === 0) {
      await taskStore.fetchTasks()
    }

    // Try loading real dependencies
    const { GetDependencies } = await import(
      '../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/dependencyservice'
    )
    const { GetSubtasks } = await import(
      '../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice'
    )

    // Build nodes from real tasks
    const nodes: Omit<GraphNode, 'children'>[] = taskStore.tasks.map(t => ({
      id: t.id,
      title: t.title,
      status: t.status,
      priority: t.priority,
      description: t.description,
      dependsOn: [] as number[],
      subtasks: [] as SubtaskItem[],
    }))

    // Fetch dependencies and subtasks for each task
    for (const node of nodes) {
      try {
        const deps = await GetDependencies(node.id)
        node.dependsOn = deps.map((d: { id: number }) => d.id)
      } catch { /* no deps */ }

      try {
        const subs = await GetSubtasks(node.id)
        node.subtasks = subs.map((s: { id: number; title: string; status: string }) => ({
          id: s.id,
          title: s.title,
          done: s.status === 'done',
        }))
      } catch { /* no subtasks */ }
    }

    graphNodes.value = nodes
  } catch (e) {
    console.warn('[DependencyGraph] Wails bindings unavailable, using mock data:', e)
    graphNodes.value = buildMockNodes()
  } finally {
    loadingGraph.value = false
  }
}

onMounted(loadGraph)

// --- Build tree ---
const allNodes = computed<GraphNode[]>(() => {
  const map = new Map<number, GraphNode>()
  for (const t of graphNodes.value) {
    map.set(t.id, { ...t, children: [], subtasks: t.subtasks })
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

          <Button
            variant="ghost"
            size="sm"
            class="h-8 px-2.5"
            :disabled="loadingGraph"
            @click="loadGraph"
          >
            <RefreshCw :size="14" :class="{ 'animate-spin': loadingGraph }" />
            <span class="ml-1 hidden text-xs sm:inline">Refresh</span>
          </Button>

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
    <ScrollArea class="min-h-0 flex-1">
      <div class="p-6 pb-12">
        <!-- Loading state -->
        <div v-if="loadingGraph" class="flex items-center justify-center gap-2 py-12 text-sm text-muted-foreground">
          <Loader2 :size="16" class="animate-spin" />
          Loading dependency graph…
        </div>

        <!-- Tree View -->
        <div v-else-if="viewMode === 'tree'" class="space-y-0.5">
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

<style scoped>
.tree-node {
  min-height: 36px;
}
</style>
