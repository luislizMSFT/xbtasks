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
import TreeNodeItem from '@/components/ado/TreeNodeItem.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'
import ForceGraph from '@/components/ForceGraph.vue'
import TaskDetail from '@/components/tasks/TaskDetail.vue'
import {
  GitBranch,
  AlertCircle,
  TreePine,
  List,
  Share2,
  Users,
  RefreshCw,
  ExternalLink,
  X,
} from 'lucide-vue-next'
import { adoTypeColor, adoTypeIcon, adoTypeHex } from '@/lib/styles'

const router = useRouter()
const taskStore = useTaskStore()

// --- Types ---
import type { SubtaskItem, GraphNode, WorkItem } from '@/types'

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
const loadingAdo = ref(false)

// ADO-only nodes displayed in the graph (not linked to any local task)
const adoNodes = ref<Map<number, { id: number; title: string; status: string; priority: string; type: string; adoId: string; adoUrl: string; parentId: number }>>(new Map())
// Maps local task ID → ADO parent node ID (negative) for hierarchy edges
const localToAdoParent = ref<Map<number, number>>(new Map())

async function loadGraph() {
  loadingGraph.value = true
  try {
    // Ensure tasks are loaded
    if (taskStore.tasks.length === 0) {
      await taskStore.fetchTasks()
    }

    // Try loading real dependencies
    const { getDependencies } = await import('@/api/dependencies')
    const { getSubtasks } = await import('@/api/tasks')

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
        const deps = await getDependencies(node.id)
        node.dependsOn = deps.map((d: { id: number }) => d.id)
      } catch { /* no deps */ }

      try {
        const subs = await getSubtasks(node.id)
        node.subtasks = subs.map((s: { id: number; title: string; status: string }) => ({
          id: s.id,
          title: s.title,
          done: s.status === 'done',
        }))
      } catch { /* no subtasks */ }
    }

    graphNodes.value = nodes

    // Phase 2: progressively load ADO parents in background
    loadAdoParents()
  } catch (e) {
    console.warn('[DependencyGraph] Wails bindings unavailable, using mock data:', e)
    graphNodes.value = buildMockNodes()
  } finally {
    loadingGraph.value = false
  }
}

// Map ADO state → local status for visual consistency
function adoStateToStatus(state: string): string {
  const s = state.toLowerCase()
  if (s === 'closed' || s === 'done' || s === 'resolved' || s === 'completed') return 'done'
  if (s === 'active' || s === 'committed') return 'in_progress'
  if (s === 'new' || s === 'proposed') return 'todo'
  if (s === 'removed') return 'cancelled'
  return 'todo'
}

async function loadAdoParents() {
  loadingAdo.value = true
  try {
    const { getCachedWorkItems, getWorkItemTree } = await import('@/api/workitems')

    // Collect local tasks that have ADO links
    const linkedAdoIds = new Set(taskStore.tasks.map(t => t.adoId).filter(Boolean))
    const localTaskByAdoId = new Map<string, number>()
    for (const t of taskStore.tasks) {
      if (t.adoId) localTaskByAdoId.set(t.adoId, t.id)
    }

    // Strategy 1: Use cached ADO work items (fast, offline-capable)
    // This covers items already synced to SQLite
    let allItems: WorkItem[] = []
    try {
      allItems = (await getCachedWorkItems()) as WorkItem[]
    } catch {
      console.warn('[DependencyGraph] Cached work items unavailable')
    }

    // Strategy 2: Also try getWorkItemTree for discovering ancestors not in cache
    try {
      const treeItems = (await getWorkItemTree()) as WorkItem[]
      // Merge tree items into allItems (deduplicate by adoId)
      const seen = new Set(allItems.map(i => String(i.adoId || i.id)))
      for (const wi of treeItems) {
        const adoIdStr = String(wi.adoId || wi.id)
        if (!seen.has(adoIdStr)) {
          allItems.push(wi)
          seen.add(adoIdStr)
        }
      }
    } catch {
      console.warn('[DependencyGraph] Work item tree unavailable, using cached data only')
    }

    // Build ADO-only nodes (parents not linked to any local task)
    const newAdoNodes = new Map<number, { id: number; title: string; status: string; priority: string; type: string; adoId: string; adoUrl: string; parentId: number }>()
    for (const wi of allItems) {
      const adoIdStr = String(wi.adoId || wi.id)
      if (!linkedAdoIds.has(adoIdStr)) {
        const nodeId = -Number(adoIdStr)
        newAdoNodes.set(nodeId, {
          id: nodeId,
          title: wi.title,
          status: adoStateToStatus(wi.state),
          priority: wi.priority ? `P${Math.min(wi.priority, 3)}` : 'P2',
          type: wi.type,
          adoId: adoIdStr,
          adoUrl: wi.url,
          parentId: wi.parentId,
        })
      }
    }

    adoNodes.value = newAdoNodes

    // Build local task → parent mapping from ADO hierarchy
    // For each cached/tree item that IS a local task, connect it to its ADO parent
    const parentMap = new Map<number, number>()
    for (const wi of allItems) {
      const adoIdStr = String(wi.adoId || wi.id)
      const localTaskId = localTaskByAdoId.get(adoIdStr)
      if (localTaskId !== undefined && wi.parentId > 0) {
        const parentAdoIdStr = String(wi.parentId)
        // Check if parent is another local task or an ADO-only node
        const parentLocalId = localTaskByAdoId.get(parentAdoIdStr)
        if (parentLocalId !== undefined) {
          parentMap.set(localTaskId, parentLocalId)
        } else {
          const parentNodeId = -Number(parentAdoIdStr)
          if (newAdoNodes.has(parentNodeId)) {
            parentMap.set(localTaskId, parentNodeId)
          }
        }
      }
    }
    localToAdoParent.value = parentMap
  } catch (e) {
    console.warn('[DependencyGraph] ADO parent loading failed:', e)
  } finally {
    loadingAdo.value = false
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
type ViewMode = 'graph' | 'tree' | 'flat'
const viewMode = ref<ViewMode>('graph')
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

// --- Graph-mode data ---
const forceGraphNodes = computed(() => {
  // Local task nodes
  const localNodes = allNodes.value.map(n => ({
    id: n.id,
    title: n.title,
    status: n.status,
    priority: n.priority,
    projectId: taskStore.tasks.find(t => t.id === n.id)?.projectId ?? null,
    source: 'local' as const,
    adoId: taskStore.tasks.find(t => t.id === n.id)?.adoId,
  }))

  // ADO-only parent nodes
  const adoGraphNodes = Array.from(adoNodes.value.values()).map(n => ({
    id: n.id,
    title: n.title,
    status: n.status,
    priority: n.priority,
    projectId: null,
    source: 'ado' as const,
    adoId: n.adoId,
    adoUrl: n.adoUrl,
    type: n.type,
  }))

  return [...localNodes, ...adoGraphNodes]
})

const forceGraphEdges = computed(() => {
  const edges: { source: number; target: number; kind: 'dependency' | 'parent' }[] = []
  const edgeSet = new Set<string>() // dedup key: "source→target"

  function addEdge(source: number, target: number, kind: 'dependency' | 'parent') {
    const key = `${source}→${target}`
    if (!edgeSet.has(key)) {
      edgeSet.add(key)
      edges.push({ source, target, kind })
    }
  }

  // 1) Local dependency edges (task_deps table)
  for (const node of allNodes.value) {
    for (const depId of node.dependsOn) {
      addEdge(depId, node.id, 'dependency')
    }
  }

  // 2) Local parent-child edges (task.parentId — from CreateSubtask)
  for (const task of taskStore.tasks) {
    if (task.parentId) {
      addEdge(task.parentId, task.id, 'parent')
    }
  }

  // 3) ADO parent→child hierarchy edges for ADO-only nodes
  const adoIdToNodeId = new Map<string, number>()
  for (const t of taskStore.tasks) {
    if (t.adoId) adoIdToNodeId.set(t.adoId, t.id)
  }
  for (const [nodeId, n] of adoNodes.value) {
    adoIdToNodeId.set(n.adoId, nodeId)
  }

  for (const [nodeId, n] of adoNodes.value) {
    if (n.parentId > 0) {
      const parentNodeId = adoIdToNodeId.get(String(n.parentId))
      if (parentNodeId !== undefined) {
        addEdge(parentNodeId, nodeId, 'parent')
      }
    }
  }

  // 4) Local task → ADO parent hierarchy edges (from cached/tree data)
  for (const [localId, parentNodeId] of localToAdoParent.value) {
    addEdge(parentNodeId, localId, 'parent')
  }

  return edges
})

// Selected ADO-only node info (when clicking a non-local node)
const selectedAdoNode = computed(() => {
  if (selectedNodeId.value === null || selectedNodeId.value >= 0) return null
  return adoNodes.value.get(selectedNodeId.value) ?? null
})

// Derive legend ADO types from actual data in the graph
const legendAdoTypes = computed(() => {
  const types = new Set<string>()
  for (const n of adoNodes.value.values()) {
    if (n.type) types.add(n.type)
  }
  // Fallback: show common org types if no ADO data loaded yet
  if (types.size === 0) {
    return ['Bug', 'Task', 'Scenario', 'Deliverable']
  }
  return Array.from(types).sort()
})

function onGraphSelect(id: number) {
  if (id === -1) {
    taskStore.selectTask(null)
    selectedNodeId.value = null
  } else if (id >= 0) {
    // Local task
    selectedNodeId.value = id
    taskStore.selectTask(id)
  } else {
    // ADO-only node (negative ID)
    selectedNodeId.value = id
    taskStore.selectTask(null)
  }
}
</script>

<template>
  <div class="flex h-full flex-col overflow-hidden">
    <!-- Top Bar -->
    <div class="flex-none border-b border-border bg-background px-4 py-2.5">
      <div class="flex items-center justify-between gap-4">
        <div class="flex items-center gap-2">
          <p class="text-sm text-muted-foreground">
            Trace blocked tasks and dependency chains
          </p>
        </div>

        <div class="flex items-center gap-2">
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
                viewMode === 'graph' && 'bg-accent text-accent-foreground',
              ]"
              @click="viewMode = 'graph'"
            >
              <Share2 :size="14" />
              <span class="ml-1 hidden text-xs sm:inline">Graph</span>
            </Button>
            <Button
              variant="ghost"
              size="sm"
              :class="[
                'h-8 rounded-none border-l border-border px-2.5',
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

    <!-- Graph View -->
    <div v-if="viewMode === 'graph'" class="flex-1 min-h-0 relative">
      <LoadingSpinner v-if="loadingGraph" label="Building dependency graph..." size="sm" class="flex-1 py-12" />
      <EmptyState
        v-else-if="allNodes.length === 0"
        :icon="GitBranch"
        title="No dependencies found"
        description="There are no tasks with dependency relationships to display."
        class="flex-1"
      />
      <template v-else>
        <ForceGraph
          :nodes="forceGraphNodes"
          :edges="forceGraphEdges"
          :selected-id="selectedNodeId"
          @select="onGraphSelect"
        />

        <!-- ADO loading indicator -->
        <div
          v-if="loadingAdo"
          class="absolute bottom-3 left-3 flex items-center gap-2 rounded-md bg-background/90 border border-border px-3 py-1.5 text-xs text-muted-foreground shadow-sm backdrop-blur-sm"
        >
          <RefreshCw :size="12" class="animate-spin" />
          Loading ADO parents…
        </div>

        <!-- Floating detail panel: local task -->
        <aside
          v-if="taskStore.selectedTask"
          class="absolute top-3 right-3 bottom-3 w-[380px] max-w-[50%] rounded-xl border border-border bg-background/95 backdrop-blur-md shadow-xl overflow-hidden flex flex-col z-10 floating-detail"
        >
          <TaskDetail @close="onGraphSelect(-1)" />
        </aside>

        <!-- Floating detail panel: ADO-only node -->
        <aside
          v-if="selectedAdoNode && !taskStore.selectedTask"
          class="absolute top-3 right-3 w-80 rounded-xl border border-border bg-background/95 backdrop-blur-md shadow-xl overflow-hidden z-10"
        >
          <div class="p-4">
              <!-- Header with type icon + close -->
              <div class="flex items-center justify-between mb-3">
                <div class="flex items-center gap-2">
                  <component
                    :is="adoTypeIcon(selectedAdoNode.type)"
                    :size="16"
                    :class="adoTypeColor(selectedAdoNode.type)"
                  />
                  <Badge variant="outline" class="text-xs gap-1" :class="adoTypeColor(selectedAdoNode.type)">
                    {{ selectedAdoNode.type || 'Work Item' }}
                  </Badge>
                </div>
                <button
                  class="rounded-md p-1 text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
                  @click="onGraphSelect(-1)"
                >
                  <X :size="14" />
                </button>
              </div>

              <!-- Title -->
              <h3 class="text-sm font-semibold mb-3 leading-snug">{{ selectedAdoNode.title }}</h3>

              <!-- Info rows -->
              <div class="space-y-2 text-xs text-muted-foreground">
                <div class="flex items-center gap-2">
                  <span class="font-medium text-foreground w-16">ADO ID</span>
                  <span class="font-mono">{{ selectedAdoNode.adoId }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="font-medium text-foreground w-16">Status</span>
                  <StatusBadge :status="selectedAdoNode.status" />
                </div>
                <div class="flex items-center gap-2">
                  <span class="font-medium text-foreground w-16">Priority</span>
                  <PriorityBadge :priority="selectedAdoNode.priority" />
                </div>
              </div>

              <!-- ADO link -->
              <a
                v-if="selectedAdoNode.adoUrl"
                :href="selectedAdoNode.adoUrl"
                target="_blank"
                class="inline-flex items-center gap-1.5 mt-4 text-xs text-blue-500 hover:text-blue-400 hover:underline transition-colors"
              >
                <ExternalLink :size="12" />
                Open in Azure DevOps
              </a>

              <p class="mt-3 text-[10px] text-muted-foreground/50 italic">
                ADO-only node · not linked to a local task
              </p>
            </div>
          </aside>
        <!-- Graph Legend -->
        <div class="absolute bottom-3 right-3 rounded-xl border border-border bg-background/90 backdrop-blur-sm px-3 py-2.5 text-[10px] text-muted-foreground shadow-sm select-none z-10">
          <div class="font-medium text-foreground text-[11px] mb-2">Legend</div>

          <!-- Node types -->
          <div class="mb-2 space-y-1">
            <div class="flex items-center gap-2">
              <svg width="12" height="12"><circle cx="6" cy="6" r="5" fill="#3b82f6" /></svg>
              <span>Local task</span>
            </div>
          </div>

          <!-- ADO work item types (derived from actual data) -->
          <div class="mb-2 space-y-1">
            <div v-for="t in legendAdoTypes" :key="t" class="flex items-center gap-2">
              <component :is="adoTypeIcon(t)" :size="11" :class="adoTypeColor(t)" />
              <span>{{ t }}</span>
            </div>
          </div>

          <!-- Edge types -->
          <div class="space-y-1 pt-1 border-t border-border">
            <div class="flex items-center gap-2 mt-1">
              <svg width="24" height="8"><line x1="0" y1="4" x2="24" y2="4" stroke="#3f3f46" stroke-width="1" stroke-dasharray="5 3" class="edge-flowing" /></svg>
              <span>Dependency</span>
            </div>
            <div class="flex items-center gap-2">
              <svg width="24" height="8"><line x1="0" y1="4" x2="24" y2="4" stroke="#ef4444" stroke-width="1.5" stroke-dasharray="4 3" stroke-opacity="0.6" /></svg>
              <span>Blocked</span>
            </div>
            <div class="flex items-center gap-2">
              <svg width="24" height="8"><line x1="0" y1="4" x2="24" y2="4" stroke="#6366f1" stroke-width="1" stroke-dasharray="2 2" stroke-opacity="0.5" /></svg>
              <span>Parent / child</span>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Tree / Flat Views -->
    <ScrollArea v-else class="min-h-0 flex-1">
      <div class="p-6 pb-12">
        <!-- Loading state -->
        <LoadingSpinner v-if="loadingGraph" label="Building dependency graph..." size="sm" class="py-12" />

        <!-- Empty state: no data at all -->
        <EmptyState
          v-else-if="allNodes.length === 0"
          :icon="GitBranch"
          title="No dependencies found"
          description="There are no tasks with dependency relationships to display."
        />

        <!-- Tree View -->
        <div v-else-if="viewMode === 'tree'" class="space-y-0.5">
          <EmptyState
            v-if="filteredRoots.length === 0"
            title="No tasks match the current filter"
            description="Try selecting a different status filter."
          />
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
          <EmptyState
            v-if="flatNodes.length === 0"
            title="No tasks match the current filter"
            description="Try selecting a different status filter."
          />
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

.edge-flowing {
  animation: edge-flow 1.5s linear infinite;
}

@keyframes edge-flow {
  from { stroke-dashoffset: 24; }
  to { stroke-dashoffset: 0; }
}

/* Override TaskDetail's own sizing when floating inside graph */
.floating-detail :deep(aside) {
  width: 100% !important;
  max-width: 100% !important;
  border-left: none !important;
  background: transparent !important;
  transform: none !important;
  transition: none !important;
}
</style>
