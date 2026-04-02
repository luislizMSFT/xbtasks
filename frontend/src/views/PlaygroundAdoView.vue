<script setup lang="ts">
import { ref, computed, type Component } from 'vue'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Separator } from '@/components/ui/separator'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import {
  Landmark,
  Search,
  Plus,
  Layers,
  Package,
  Bug,
  BookOpen,
  CheckSquare,
  ChevronRight,
  ChevronDown,
  Check,
  Circle,
  Link2,
  ExternalLink,
  User,
  GitPullRequest,
  GitBranch,
  ArrowRight,
  CheckCircle2,
  XCircle,
  Loader2,
  Filter,
  FileCode,
  MessageSquare,
  Zap,
  Bell,
  RefreshCw,
  X,
  CircleDot,
  Link2Off,
  Square,
} from 'lucide-vue-next'

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface WorkItem {
  id: number
  title: string
  type: 'scenario' | 'deliverable' | 'task' | 'bug' | 'userStory'
  state: string
  priority: number
  linked: boolean
  areaPath: string
  iteration: string
  assignedTo: string
  description: string
  children: WorkItem[]
}

interface FlatNode {
  item: WorkItem
  depth: number
  hasChildren: boolean
  expanded: boolean
}

interface Reviewer {
  initials: string
  color: string
}

interface PullRequest {
  id: number
  title: string
  repo: string
  sourceBranch: string
  targetBranch: string
  reviewers: Reviewer[]
  additions: number
  deletions: number
  timeAgo: string
  status: 'review-requested' | 'approved' | 'changes-requested' | 'merged'
  description: string
  filesChanged: number
  checksPass: boolean
}

interface Pipeline {
  name: string
  runId: number
  status: 'succeeded' | 'running' | 'failed'
  duration: string
  trigger: string
}

interface QueuedOperation {
  id: number
  type: 'link' | 'unlink' | 'status-change' | 'follow-up' | 'sync'
  description: string
  target: string
  icon: Component
  color: string
  targetId?: number
}

// ---------------------------------------------------------------------------
// Mock data — work items
// ---------------------------------------------------------------------------

const workItems: WorkItem[] = [
  {
    id: 50000,
    title: 'Xbox Live Services Q2',
    type: 'scenario',
    state: 'Active',
    priority: 0,
    linked: false,
    areaPath: 'Xbox\\Live Services',
    iteration: 'Q2 2025',
    assignedTo: 'Team Lead',
    description: 'Quarterly initiative for Xbox Live Services covering auth improvements and dashboard modernization.',
    children: [
      {
        id: 50100,
        title: 'Auth Improvements',
        type: 'deliverable',
        state: 'Active',
        priority: 0,
        linked: false,
        areaPath: 'Xbox\\Live Services\\Auth',
        iteration: 'Q2 2025\\Sprint 1',
        assignedTo: 'Auth Team',
        description: 'Improve authentication flows, fix known issues, and integrate SSO.',
        children: [
          {
            id: 48291,
            title: 'Fix auth redirect loop',
            type: 'task',
            state: 'Active',
            priority: 0,
            linked: true,
            areaPath: 'Xbox\\Live Services\\Auth',
            iteration: 'Q2 2025\\Sprint 1',
            assignedTo: 'Jane Smith',
            description: 'Auth redirect creates an infinite loop when the session token expires and the refresh token flow fails silently.',
            children: [],
          },
          {
            id: 48292,
            title: 'Token expiry timing',
            type: 'bug',
            state: 'Active',
            priority: 1,
            linked: true,
            areaPath: 'Xbox\\Live Services\\Auth',
            iteration: 'Q2 2025\\Sprint 1',
            assignedTo: 'John Doe',
            description: 'Tokens expire 5 minutes earlier than expected due to clock skew between services.',
            children: [],
          },
          {
            id: 48300,
            title: 'SSO integration',
            type: 'userStory',
            state: 'New',
            priority: 2,
            linked: false,
            areaPath: 'Xbox\\Live Services\\Auth',
            iteration: 'Q2 2025\\Sprint 2',
            assignedTo: 'Unassigned',
            description: 'As a user, I want to log in once and access all Xbox services without re-authenticating.',
            children: [],
          },
        ],
      },
      {
        id: 50200,
        title: 'Dashboard Modernization',
        type: 'deliverable',
        state: 'Active',
        priority: 1,
        linked: false,
        areaPath: 'Xbox\\Live Services\\Dashboard',
        iteration: 'Q2 2025\\Sprint 2',
        assignedTo: 'Frontend Team',
        description: 'Rebuild the dashboard with modern components, real-time data, and activity feed.',
        children: [
          {
            id: 48350,
            title: 'Build new dashboard layout',
            type: 'task',
            state: 'Active',
            priority: 1,
            linked: true,
            areaPath: 'Xbox\\Live Services\\Dashboard',
            iteration: 'Q2 2025\\Sprint 2',
            assignedTo: 'Alex Johnson',
            description: 'Create the new dashboard layout with grid-based widget system and responsive breakpoints.',
            children: [],
          },
          {
            id: 48360,
            title: 'Chart data mismatch',
            type: 'bug',
            state: 'Resolved',
            priority: 2,
            linked: false,
            areaPath: 'Xbox\\Live Services\\Dashboard',
            iteration: 'Q2 2025\\Sprint 1',
            assignedTo: 'Maria Garcia',
            description: 'Charts show stale data due to caching layer returning previous day totals.',
            children: [],
          },
          {
            id: 48370,
            title: 'Activity feed',
            type: 'userStory',
            state: 'New',
            priority: 3,
            linked: false,
            areaPath: 'Xbox\\Live Services\\Dashboard',
            iteration: 'Q2 2025\\Sprint 3',
            assignedTo: 'Unassigned',
            description: 'As a team lead, I want a real-time activity feed showing recent changes across all projects.',
            children: [],
          },
        ],
      },
    ],
  },
  {
    id: 51000,
    title: 'Developer Tooling Q2',
    type: 'scenario',
    state: 'Active',
    priority: 0,
    linked: false,
    areaPath: 'Xbox\\Developer Tooling',
    iteration: 'Q2 2025',
    assignedTo: 'Platform Team',
    description: 'Quarterly initiative for developer tooling improvements including CLI tools and SDK updates.',
    children: [
      {
        id: 51100,
        title: 'CLI Tools',
        type: 'deliverable',
        state: 'Active',
        priority: 0,
        linked: false,
        areaPath: 'Xbox\\Developer Tooling\\CLI',
        iteration: 'Q2 2025\\Sprint 1',
        assignedTo: 'CLI Team',
        description: 'Improve CLI tools with rate limiting, better auth, and MSAL v3 migration.',
        children: [
          {
            id: 48400,
            title: 'Add rate limiting',
            type: 'task',
            state: 'Active',
            priority: 1,
            linked: false,
            areaPath: 'Xbox\\Developer Tooling\\CLI',
            iteration: 'Q2 2025\\Sprint 2',
            assignedTo: 'Chris Lee',
            description: 'Implement client-side rate limiting with exponential backoff for API calls.',
            children: [],
          },
          {
            id: 48100,
            title: 'Migrate to MSAL v3',
            type: 'task',
            state: 'Blocked',
            priority: 0,
            linked: true,
            areaPath: 'Xbox\\Developer Tooling\\CLI',
            iteration: 'Q2 2025\\Sprint 1',
            assignedTo: 'Sam Wilson',
            description: 'Migrate auth from ADAL to MSAL v3. Blocked on MSAL v3 GA release.',
            children: [],
          },
        ],
      },
    ],
  },
]

// ---------------------------------------------------------------------------
// Mock data — pull requests
// ---------------------------------------------------------------------------

const reviewRequestedPrs: PullRequest[] = [
  {
    id: 234,
    title: 'Fix auth redirect loop',
    repo: 'xb-services',
    sourceBranch: 'fix/auth-redirect',
    targetBranch: 'main',
    reviewers: [
      { initials: 'JS', color: 'bg-blue-500' },
      { initials: 'AJ', color: 'bg-emerald-500' },
    ],
    additions: 400,
    deletions: 120,
    timeAgo: '3h ago',
    status: 'review-requested',
    description: 'Fixes the infinite redirect loop when auth tokens expire. Adds proper error boundary and fallback to login page.',
    filesChanged: 12,
    checksPass: true,
  },
  {
    id: 301,
    title: 'Add telemetry hooks',
    repo: 'xb-telemetry',
    sourceBranch: 'feat/hooks',
    targetBranch: 'main',
    reviewers: [{ initials: 'MG', color: 'bg-purple-500' }],
    additions: 80,
    deletions: 10,
    timeAgo: '1d ago',
    status: 'review-requested',
    description: 'Adds hooks for telemetry event tracking with automatic context propagation.',
    filesChanged: 4,
    checksPass: true,
  },
]

const yourPrs: PullRequest[] = [
  {
    id: 298,
    title: 'Schema migration v2',
    repo: 'xb-data',
    sourceBranch: 'feat/schema-v2',
    targetBranch: 'main',
    reviewers: [
      { initials: 'CL', color: 'bg-amber-500' },
      { initials: 'SW', color: 'bg-rose-500' },
    ],
    additions: 250,
    deletions: 80,
    timeAgo: '2d ago',
    status: 'approved',
    description: 'Major schema migration adding new entity types and relations for v2 data model.',
    filesChanged: 18,
    checksPass: true,
  },
  {
    id: 280,
    title: 'Rate limit middleware',
    repo: 'xb-services',
    sourceBranch: 'feat/rate-limit',
    targetBranch: 'main',
    reviewers: [{ initials: 'JD', color: 'bg-cyan-500' }],
    additions: 180,
    deletions: 20,
    timeAgo: '3d ago',
    status: 'changes-requested',
    description: 'Adds configurable rate limiting middleware with sliding window algorithm.',
    filesChanged: 8,
    checksPass: false,
  },
  {
    id: 250,
    title: 'Update deployment docs',
    repo: 'xb-infra',
    sourceBranch: 'docs/deploy',
    targetBranch: 'main',
    reviewers: [{ initials: 'AJ', color: 'bg-emerald-500' }],
    additions: 60,
    deletions: 30,
    timeAgo: '5d ago',
    status: 'merged',
    description: 'Updates deployment documentation with new Kubernetes configuration and rollback procedures.',
    filesChanged: 3,
    checksPass: true,
  },
]

// ---------------------------------------------------------------------------
// Mock data — pipelines
// ---------------------------------------------------------------------------

const pipelines: Pipeline[] = [
  { name: 'Deploy xb-services', runId: 1245, status: 'succeeded', duration: '4m 32s', trigger: 'PR #250 merge' },
  { name: 'Deploy xb-data', runId: 1244, status: 'running', duration: '2m 10s', trigger: 'PR #298 merge' },
  { name: 'Deploy xb-telemetry', runId: 1240, status: 'failed', duration: '3m 15s', trigger: 'Test stage failed' },
]

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

const darkMode = ref(false)
const activityClockActive = ref(false)
const searchQuery = ref('')
const typeFilter = ref('All')
const stateFilter = ref('All')
const assignedToMe = ref(false)
const selectedItem = ref<WorkItem | null>(null)
const expandedNodes = ref<Record<number, boolean>>({
  50000: true,
  50100: true,
  51000: true,
  51100: true,
})
const expandedPrs = ref<Record<number, boolean>>({})

// Operations queue (GParted-style queue & confirm)
const operationQueue = ref<QueuedOperation[]>([
  { id: 1, type: 'link', description: 'Link to local task', target: '"Fix auth redirect loop" \u2192 ADO Bug #48291', icon: Link2, color: 'text-blue-500', targetId: 48291 },
  { id: 2, type: 'status-change', description: 'Change state to Resolved', target: 'Bug #48292 "Token expiry timing"', icon: CircleDot, color: 'text-emerald-500', targetId: 48292 },
  { id: 3, type: 'follow-up', description: 'Set reminder', target: 'Task #48400 "Add rate limiting"', icon: Bell, color: 'text-amber-500', targetId: 48400 },
  { id: 4, type: 'sync', description: 'Pull latest comments', target: 'ADO #48350 "Build new dashboard layout"', icon: RefreshCw, color: 'text-violet-500', targetId: 48350 },
])
const showQueueReview = ref(false)
const checkedOperations = ref<Record<number, boolean>>({})
let nextOpId = 5

// ---------------------------------------------------------------------------
// Computed
// ---------------------------------------------------------------------------

const flatTree = computed<FlatNode[]>(() => {
  const result: FlatNode[] = []
  function walk(items: WorkItem[], depth: number) {
    for (const item of items) {
      const hasChildren = item.children.length > 0
      const expanded = !!expandedNodes.value[item.id]
      result.push({ item, depth, hasChildren, expanded })
      if (expanded && hasChildren) {
        walk(item.children, depth + 1)
      }
    }
  }
  walk(workItems, 0)
  return result
})

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

function toggleNode(id: number) {
  expandedNodes.value[id] = !expandedNodes.value[id]
}

function selectItem(item: WorkItem) {
  selectedItem.value = selectedItem.value?.id === item.id ? null : item
}

function togglePr(id: number) {
  expandedPrs.value[id] = !expandedPrs.value[id]
}

function cycleTypeFilter() {
  const types = ['All', 'Task', 'Bug', 'User Story', 'Deliverable', 'Scenario']
  const idx = types.indexOf(typeFilter.value)
  typeFilter.value = types[(idx + 1) % types.length]
}

function cycleStateFilter() {
  const states = ['All', 'Active', 'New', 'Resolved', 'Blocked']
  const idx = states.indexOf(stateFilter.value)
  stateFilter.value = states[(idx + 1) % states.length]
}

function getTypeInfo(type: string) {
  switch (type) {
    case 'scenario': return { icon: Layers, color: 'text-orange-500', label: 'Scenario' }
    case 'deliverable': return { icon: Package, color: 'text-violet-500', label: 'Deliverable' }
    case 'task': return { icon: CheckSquare, color: 'text-blue-500', label: 'Task' }
    case 'bug': return { icon: Bug, color: 'text-red-500', label: 'Bug' }
    case 'userStory': return { icon: BookOpen, color: 'text-green-500', label: 'User Story' }
    default: return { icon: Circle, color: 'text-muted-foreground', label: type }
  }
}

function getStateClass(state: string) {
  switch (state) {
    case 'Active': return 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25'
    case 'New': return 'bg-muted text-muted-foreground border-border'
    case 'Resolved': return 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25'
    case 'Blocked': return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    default: return ''
  }
}

function getPriorityClass(p: number) {
  switch (p) {
    case 0: return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    case 1: return 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/25'
    case 2: return 'bg-yellow-500/15 text-yellow-700 dark:text-yellow-400 border-yellow-500/25'
    case 3: return 'bg-muted text-muted-foreground border-border'
    default: return ''
  }
}

function getPrStatusInfo(status: string) {
  switch (status) {
    case 'approved': return { label: 'Approved', cls: 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25' }
    case 'changes-requested': return { label: 'Changes requested', cls: 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/25' }
    case 'merged': return { label: 'Merged', cls: 'bg-violet-500/15 text-violet-700 dark:text-violet-400 border-violet-500/25' }
    case 'review-requested': return { label: 'Review requested', cls: 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25' }
    default: return { label: status, cls: '' }
  }
}

function getPipelineIcon(status: string) {
  switch (status) {
    case 'succeeded': return { icon: CheckCircle2, color: 'text-green-500' }
    case 'running': return { icon: Loader2, color: 'text-blue-500 animate-spin' }
    case 'failed': return { icon: XCircle, color: 'text-red-500' }
    default: return { icon: Circle, color: 'text-muted-foreground' }
  }
}

// ---------------------------------------------------------------------------
// Operations queue helpers
// ---------------------------------------------------------------------------

const pendingItemIds = computed(() => {
  const ids = new Set<number>()
  for (const op of operationQueue.value) {
    if (op.targetId != null) ids.add(op.targetId)
  }
  return ids
})

const checkedCount = computed(() =>
  operationQueue.value.filter(op => checkedOperations.value[op.id] !== false).length,
)

function initChecked() {
  const c: Record<number, boolean> = {}
  for (const op of operationQueue.value) c[op.id] = true
  checkedOperations.value = c
}

function queueLink(item: WorkItem) {
  if (pendingItemIds.value.has(item.id)) return
  const info = getTypeInfo(item.type)
  operationQueue.value.push({
    id: nextOpId++,
    type: 'link',
    description: 'Link to local task',
    target: `"${item.title}" \u2192 ADO ${info.label} #${item.id}`,
    icon: Link2,
    color: 'text-blue-500',
    targetId: item.id,
  })
}

function queueUnlink(item: WorkItem) {
  if (pendingItemIds.value.has(item.id)) return
  const info = getTypeInfo(item.type)
  operationQueue.value.push({
    id: nextOpId++,
    type: 'unlink',
    description: 'Remove link',
    target: `${info.label} #${item.id} "${item.title}"`,
    icon: Link2Off,
    color: 'text-red-500',
    targetId: item.id,
  })
}

function removeOperation(id: number) {
  operationQueue.value = operationQueue.value.filter(op => op.id !== id)
  delete checkedOperations.value[id]
}

function discardAll() {
  operationQueue.value = []
  checkedOperations.value = {}
  showQueueReview.value = false
}

function applyAll() {
  operationQueue.value = operationQueue.value.filter(op => checkedOperations.value[op.id] === false)
  checkedOperations.value = {}
  showQueueReview.value = false
}

function openReview() {
  initChecked()
  showQueueReview.value = true
}

function toggleOpChecked(id: number) {
  checkedOperations.value[id] = !(checkedOperations.value[id] !== false)
}

function getPendingTooltip(itemId: number): string {
  const ops = operationQueue.value.filter(op => op.targetId === itemId)
  return ops.map(op => `Pending: ${op.description}`).join(', ')
}
</script>

<template>
  <TooltipProvider :delay-duration="200">
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Tabs container -->
      <Tabs default-value="management" class="flex-1 flex flex-col min-h-0">
        <div class="px-4 pt-3 pb-1">
          <TabsList class="h-8">
            <TabsTrigger value="management" class="text-xs px-3 h-7">Management</TabsTrigger>
            <TabsTrigger value="devops" class="text-xs px-3 h-7">DevOps</TabsTrigger>
          </TabsList>
        </div>
            </div>
          </div>

          <!-- Tabs container -->
          <Tabs default-value="management" class="flex-1 flex flex-col min-h-0">
            <div class="px-4 pt-3 pb-1">
              <TabsList class="h-8">
                <TabsTrigger value="management" class="text-xs px-3 h-7">Management</TabsTrigger>
                <TabsTrigger value="devops" class="text-xs px-3 h-7">DevOps</TabsTrigger>
              </TabsList>
            </div>

            <!-- ============================================================ -->
            <!-- TAB 1 — Management (Work Items)                              -->
            <!-- ============================================================ -->
            <TabsContent value="management" class="flex-1 min-h-0 mt-0 flex flex-col">

              <!-- Filter bar -->
              <div class="px-4 py-2 flex items-center gap-2 border-b border-border shrink-0">
                <Button variant="outline" size="sm" class="h-7 text-xs gap-1.5" @click="cycleTypeFilter">
                  <Filter :size="12" />
                  Type: {{ typeFilter }}
                </Button>
                <Button variant="outline" size="sm" class="h-7 text-xs gap-1.5" @click="cycleStateFilter">
                  <Filter :size="12" />
                  State: {{ stateFilter }}
                </Button>
                <div class="relative flex-1 max-w-[240px]">
                  <Search :size="14" class="absolute left-2 top-1/2 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    v-model="searchQuery"
                    placeholder="Search work items..."
                    class="h-7 text-xs pl-7"
                  />
                </div>
                <Button
                  :variant="assignedToMe ? 'default' : 'outline'"
                  size="sm"
                  class="h-7 text-xs gap-1.5 ml-auto"
                  @click="assignedToMe = !assignedToMe"
                >
                  <User :size="12" />
                  Assigned to me
                </Button>
              </div>

              <!-- Two-column: tree + detail -->
              <div class="flex-1 flex min-h-0">
                <!-- Work item tree -->
                <ScrollArea class="flex-1 min-w-0">
                  <div class="py-1">
                    <button
                      v-for="node in flatTree"
                      :key="node.item.id"
                      :class="cn(
                        'w-full flex items-center gap-1.5 px-3 py-1.5 text-left hover:bg-muted/50 transition-colors text-sm',
                        selectedItem?.id === node.item.id && 'bg-muted',
                        pendingItemIds.has(node.item.id) && 'border-l-2 border-dashed border-amber-500',
                      )"
                      :style="{ paddingLeft: `${node.depth * 20 + 12}px` }"
                      @click="selectItem(node.item)"
                    >
                      <!-- Expand / collapse chevron -->
                      <span
                        v-if="node.hasChildren"
                        class="w-4 h-4 flex items-center justify-center shrink-0 text-muted-foreground hover:text-foreground cursor-pointer"
                        @click.stop="toggleNode(node.item.id)"
                      >
                        <ChevronDown v-if="node.expanded" :size="14" />
                        <ChevronRight v-else :size="14" />
                      </span>
                      <span v-else class="w-4 shrink-0" />

                      <!-- Type icon -->
                      <component
                        :is="getTypeInfo(node.item.type).icon"
                        :size="14"
                        :class="['shrink-0', getTypeInfo(node.item.type).color]"
                      />

                      <!-- Title -->
                      <span class="truncate text-foreground">{{ node.item.title }}</span>

                      <!-- ID -->
                      <span class="text-xs text-muted-foreground shrink-0">#{{ node.item.id }}</span>

                      <!-- State badge -->
                      <Badge
                        variant="outline"
                        :class="cn('text-[10px] h-4 px-1.5 shrink-0', getStateClass(node.item.state))"
                      >
                        {{ node.item.state }}
                      </Badge>

                      <!-- Priority (leaf items only) -->
                      <Badge
                        v-if="node.item.children.length === 0"
                        variant="outline"
                        :class="cn('text-[10px] h-4 px-1.5 shrink-0', getPriorityClass(node.item.priority))"
                      >
                        P{{ node.item.priority }}
                      </Badge>

                      <!-- Linked indicator -->
                      <Tooltip v-if="node.item.children.length === 0">
                        <TooltipTrigger as-child>
                          <span
                            v-if="node.item.linked"
                            class="w-4 h-4 flex items-center justify-center shrink-0 text-blue-500 cursor-pointer"
                            @click.stop="queueUnlink(node.item)"
                          >
                            <Check :size="12" />
                          </span>
                          <span
                            v-else
                            class="w-4 h-4 flex items-center justify-center shrink-0 text-muted-foreground/50 hover:text-blue-500 transition-colors cursor-pointer"
                            @click.stop="queueLink(node.item)"
                          >
                            <Circle :size="12" />
                          </span>
                        </TooltipTrigger>
                        <TooltipContent>
                          <template v-if="pendingItemIds.has(node.item.id)">
                            {{ getPendingTooltip(node.item.id) }}
                          </template>
                          <template v-else>
                            {{ node.item.linked ? 'Linked to local task' : 'Click to queue link' }}
                          </template>
                        </TooltipContent>
                      </Tooltip>
                      <!-- Pending operation indicator -->
                      <span
                        v-if="node.item.children.length === 0 && pendingItemIds.has(node.item.id)"
                        class="w-2 h-2 rounded-full bg-amber-500 shrink-0 animate-pulse"
                      />
                    </button>
                  </div>
                </ScrollArea>

                <!-- Detail panel -->
                <div
                  v-if="selectedItem"
                  class="w-[350px] border-l border-border shrink-0 flex flex-col min-h-0"
                >
                  <ScrollArea class="flex-1">
                    <div class="p-4 space-y-4">

                      <!-- Header -->
                      <div class="space-y-1.5">
                        <div class="flex items-center gap-2">
                          <component
                            :is="getTypeInfo(selectedItem.type).icon"
                            :size="16"
                            :class="getTypeInfo(selectedItem.type).color"
                          />
                          <span class="text-xs text-muted-foreground">
                            {{ getTypeInfo(selectedItem.type).label }}
                          </span>
                          <Badge variant="outline" class="text-[10px] h-4 px-1.5 ml-auto">
                            #{{ selectedItem.id }}
                          </Badge>
                        </div>
                        <h3 class="text-sm font-semibold text-foreground leading-snug">
                          {{ selectedItem.title }}
                        </h3>
                      </div>

                      <Separator />

                      <!-- State & Priority -->
                      <div class="flex items-center gap-2">
                        <div class="space-y-1">
                          <span class="text-[10px] text-muted-foreground uppercase tracking-wider">State</span>
                          <Badge
                            variant="outline"
                            :class="cn('text-xs h-5 px-2', getStateClass(selectedItem.state))"
                          >
                            {{ selectedItem.state }}
                          </Badge>
                        </div>
                        <div class="space-y-1 ml-6">
                          <span class="text-[10px] text-muted-foreground uppercase tracking-wider">Priority</span>
                          <Badge
                            variant="outline"
                            :class="cn('text-xs h-5 px-2', getPriorityClass(selectedItem.priority))"
                          >
                            P{{ selectedItem.priority }}
                          </Badge>
                        </div>
                      </div>

                      <Separator />

                      <!-- Metadata -->
                      <div class="space-y-2.5 text-xs">
                        <div class="flex items-start gap-2">
                          <FolderKanban :size="13" class="text-muted-foreground mt-0.5 shrink-0" />
                          <div>
                            <div class="text-muted-foreground text-[10px] uppercase tracking-wider">Area Path</div>
                            <div class="text-foreground">{{ selectedItem.areaPath }}</div>
                          </div>
                        </div>
                        <div class="flex items-start gap-2">
                          <Clock :size="13" class="text-muted-foreground mt-0.5 shrink-0" />
                          <div>
                            <div class="text-muted-foreground text-[10px] uppercase tracking-wider">Iteration</div>
                            <div class="text-foreground">{{ selectedItem.iteration }}</div>
                          </div>
                        </div>
                        <div class="flex items-start gap-2">
                          <User :size="13" class="text-muted-foreground mt-0.5 shrink-0" />
                          <div>
                            <div class="text-muted-foreground text-[10px] uppercase tracking-wider">Assigned To</div>
                            <div class="text-foreground">{{ selectedItem.assignedTo }}</div>
                          </div>
                        </div>
                      </div>

                      <Separator />

                      <!-- Description -->
                      <div class="space-y-1.5">
                        <span class="text-[10px] text-muted-foreground uppercase tracking-wider">Description</span>
                        <p class="text-xs text-foreground/80 leading-relaxed">
                          {{ selectedItem.description }}
                        </p>
                      </div>

                      <Separator />

                      <!-- Linked local task -->
                      <div class="space-y-1.5">
                        <span class="text-[10px] text-muted-foreground uppercase tracking-wider">Local Task</span>
                        <Card v-if="selectedItem.linked" class="bg-muted/30">
                          <CardContent class="p-2.5 flex items-center gap-2">
                            <Link2 :size="13" class="text-blue-500 shrink-0" />
                            <span class="text-xs text-foreground truncate">
                              Linked to local task
                            </span>
                            <Badge
                              v-if="pendingItemIds.has(selectedItem.id)"
                              variant="outline"
                              class="text-[10px] h-4 px-1.5 ml-auto bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/25"
                            >
                              Pending
                            </Badge>
                            <Badge
                              v-else
                              variant="outline"
                              class="text-[10px] h-4 px-1.5 ml-auto bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25"
                            >
                              Synced
                            </Badge>
                            <Button
                              v-if="!pendingItemIds.has(selectedItem.id)"
                              variant="ghost"
                              size="sm"
                              class="h-5 w-5 p-0 text-muted-foreground hover:text-red-500"
                              @click="queueUnlink(selectedItem)"
                            >
                              <Link2Off :size="10" />
                            </Button>
                          </CardContent>
                        </Card>
                        <div v-else class="flex items-center gap-2 text-xs text-muted-foreground">
                          <Circle :size="13" />
                          <span>Not linked</span>
                          <Button
                            variant="ghost"
                            size="sm"
                            class="h-5 text-[10px] ml-auto gap-1 text-blue-500 hover:text-blue-600"
                            @click="queueLink(selectedItem!)"
                          >
                            <Link2 :size="10" />
                            {{ pendingItemIds.has(selectedItem!.id) ? 'Queued' : 'Link' }}
                          </Button>
                        </div>
                      </div>

                      <!-- Open in ADO -->
                      <Button variant="outline" size="sm" class="w-full h-8 text-xs gap-1.5">
                        <ExternalLink :size="13" />
                        Open in Azure DevOps
                      </Button>
                    </div>
                  </ScrollArea>
                </div>

                <!-- No selection placeholder -->
                <div
                  v-else
                  class="w-[350px] border-l border-border shrink-0 flex items-center justify-center"
                >
                  <div class="text-center text-muted-foreground space-y-1.5">
                    <Landmark :size="24" class="mx-auto opacity-30" />
                    <p class="text-xs">Select an item to view details</p>
                  </div>
                </div>
              </div>

              <!-- Operations queue bottom bar -->
              <Transition
                enter-active-class="transition-all duration-300 ease-out"
                leave-active-class="transition-all duration-200 ease-in"
                enter-from-class="translate-y-full opacity-0"
                enter-to-class="translate-y-0 opacity-100"
                leave-from-class="translate-y-0 opacity-100"
                leave-to-class="translate-y-full opacity-0"
              >
                <div v-if="operationQueue.length > 0" class="shrink-0 bg-card border-t border-border">
                  <!-- Review panel (expanded) -->
                  <div v-if="showQueueReview" class="border-b border-border">
                    <div class="px-4 py-2 flex items-center gap-2">
                      <span class="text-xs font-semibold text-foreground">Review Pending Operations</span>
                    </div>
                    <ScrollArea class="max-h-[200px]">
                      <div class="px-4 pb-2 space-y-1">
                        <div
                          v-for="op in operationQueue"
                          :key="op.id"
                          :class="cn(
                            'flex items-center gap-2 px-2 py-1.5 rounded-md text-sm hover:bg-muted/50 transition-colors',
                            checkedOperations[op.id] === false && 'opacity-40',
                          )"
                        >
                          <button
                            class="w-4 h-4 flex items-center justify-center shrink-0 text-foreground"
                            @click="toggleOpChecked(op.id)"
                          >
                            <CheckSquare v-if="checkedOperations[op.id] !== false" :size="14" />
                            <Square v-else :size="14" />
                          </button>
                          <component :is="op.icon" :size="14" :class="['shrink-0', op.color]" />
                          <span class="text-xs text-foreground truncate">
                            {{ op.description }}: {{ op.target }}
                          </span>
                          <button
                            class="w-5 h-5 flex items-center justify-center shrink-0 ml-auto text-muted-foreground hover:text-red-500 transition-colors rounded"
                            @click="removeOperation(op.id)"
                          >
                            <X :size="12" />
                          </button>
                        </div>
                      </div>
                    </ScrollArea>
                    <div class="px-4 py-2 flex items-center justify-end gap-2 border-t border-border">
                      <Button variant="ghost" size="sm" class="h-7 text-xs" @click="showQueueReview = false">
                        Cancel
                      </Button>
                      <Button size="sm" class="h-7 text-xs gap-1.5" @click="applyAll">
                        <Check :size="12" />
                        Apply All ({{ checkedCount }})
                      </Button>
                    </div>
                  </div>

                  <!-- Collapsed bar -->
                  <div v-if="!showQueueReview" class="h-11 px-4 flex items-center gap-3">
                    <div class="flex items-center gap-2 text-sm text-foreground">
                      <Zap :size="14" class="text-amber-500" />
                      <span class="text-xs font-medium">{{ operationQueue.length }} pending operation{{ operationQueue.length !== 1 ? 's' : '' }}</span>
                    </div>
                    <div class="flex items-center gap-2 ml-auto">
                      <Button size="sm" class="h-7 text-xs gap-1.5" @click="openReview">
                        Review &amp; Apply
                      </Button>
                      <Button variant="ghost" size="sm" class="h-7 text-xs text-destructive hover:text-destructive" @click="discardAll">
                        Discard All
                      </Button>
                    </div>
                  </div>
                </div>
              </Transition>
            </TabsContent>
            <TabsContent value="devops" class="flex-1 min-h-0 mt-0">
              <ScrollArea class="h-full">
                <div class="p-4 space-y-6">

                  <!-- Review Requested -->
                  <section class="space-y-2">
                    <div class="flex items-center gap-2">
                      <h3 class="text-xs font-semibold text-foreground uppercase tracking-wider">
                        Review Requested
                      </h3>
                      <Badge variant="outline" class="text-[10px] h-4 px-1.5">
                        {{ reviewRequestedPrs.length }}
                      </Badge>
                    </div>
                    <div class="space-y-1.5">
                      <div v-for="pr in reviewRequestedPrs" :key="pr.id">
                        <Card
                          class="cursor-pointer hover:bg-muted/30 transition-colors"
                          @click="togglePr(pr.id)"
                        >
                          <CardContent class="p-3">
                            <div class="flex items-center gap-2">
                              <GitPullRequest :size="14" class="text-blue-500 shrink-0" />
                              <span class="text-sm font-medium text-foreground truncate">{{ pr.title }}</span>
                              <Badge variant="outline" class="text-[10px] h-4 px-1.5 shrink-0">
                                #{{ pr.id }}
                              </Badge>
                            </div>

                            <div class="flex items-center gap-2 mt-1.5 text-xs text-muted-foreground">
                              <Badge variant="secondary" class="text-[10px] h-4 px-1.5">{{ pr.repo }}</Badge>
                              <span class="flex items-center gap-1 text-[11px]">
                                <GitBranch :size="11" />
                                {{ pr.sourceBranch }}
                              </span>
                              <ArrowRight :size="10" class="shrink-0" />
                              <span class="text-[11px]">{{ pr.targetBranch }}</span>
                              <span class="ml-auto flex items-center gap-1.5">
                                <span
                                  v-for="(r, ri) in pr.reviewers"
                                  :key="ri"
                                  :class="cn('w-5 h-5 rounded-full flex items-center justify-center text-[9px] font-medium text-white', r.color)"
                                >
                                  {{ r.initials }}
                                </span>
                              </span>
                            </div>

                            <div class="flex items-center gap-3 mt-1.5 text-[11px] text-muted-foreground">
                              <span class="text-green-600 dark:text-green-400">+{{ pr.additions }}</span>
                              <span class="text-red-600 dark:text-red-400">-{{ pr.deletions }}</span>
                              <span class="ml-auto">{{ pr.timeAgo }}</span>
                            </div>

                            <!-- Expanded detail -->
                            <div
                              v-if="expandedPrs[pr.id]"
                              class="mt-3 pt-3 border-t border-border space-y-2 text-xs"
                            >
                              <p class="text-foreground/80 leading-relaxed">{{ pr.description }}</p>
                              <div class="flex items-center gap-4 text-muted-foreground">
                                <span class="flex items-center gap-1">
                                  <FileCode :size="12" />
                                  {{ pr.filesChanged }} files
                                </span>
                                <span class="flex items-center gap-1">
                                  <CheckCircle2
                                    :size="12"
                                    :class="pr.checksPass ? 'text-green-500' : 'text-red-500'"
                                  />
                                  Checks {{ pr.checksPass ? 'passing' : 'failing' }}
                                </span>
                              </div>
                            </div>
                          </CardContent>
                        </Card>
                      </div>
                    </div>
                  </section>

                  <Separator />

                  <!-- Your PRs -->
                  <section class="space-y-2">
                    <div class="flex items-center gap-2">
                      <h3 class="text-xs font-semibold text-foreground uppercase tracking-wider">
                        Your PRs
                      </h3>
                      <Badge variant="outline" class="text-[10px] h-4 px-1.5">
                        {{ yourPrs.length }}
                      </Badge>
                    </div>
                    <div class="space-y-1.5">
                      <div v-for="pr in yourPrs" :key="pr.id">
                        <Card
                          class="cursor-pointer hover:bg-muted/30 transition-colors"
                          @click="togglePr(pr.id)"
                        >
                          <CardContent class="p-3">
                            <div class="flex items-center gap-2">
                              <GitPullRequest
                                :size="14"
                                :class="pr.status === 'merged' ? 'text-violet-500' : 'text-blue-500'"
                                class="shrink-0"
                              />
                              <span class="text-sm font-medium text-foreground truncate">{{ pr.title }}</span>
                              <Badge variant="outline" class="text-[10px] h-4 px-1.5 shrink-0">
                                #{{ pr.id }}
                              </Badge>
                              <Badge
                                variant="outline"
                                :class="cn('text-[10px] h-4 px-1.5 shrink-0 ml-auto', getPrStatusInfo(pr.status).cls)"
                              >
                                {{ getPrStatusInfo(pr.status).label }}
                              </Badge>
                            </div>

                            <div class="flex items-center gap-2 mt-1.5 text-xs text-muted-foreground">
                              <Badge variant="secondary" class="text-[10px] h-4 px-1.5">{{ pr.repo }}</Badge>
                              <span class="flex items-center gap-1 text-[11px]">
                                <GitBranch :size="11" />
                                {{ pr.sourceBranch }}
                              </span>
                              <ArrowRight :size="10" class="shrink-0" />
                              <span class="text-[11px]">{{ pr.targetBranch }}</span>
                              <span class="ml-auto flex items-center gap-1.5">
                                <span
                                  v-for="(r, ri) in pr.reviewers"
                                  :key="ri"
                                  :class="cn('w-5 h-5 rounded-full flex items-center justify-center text-[9px] font-medium text-white', r.color)"
                                >
                                  {{ r.initials }}
                                </span>
                              </span>
                            </div>

                            <div class="flex items-center gap-3 mt-1.5 text-[11px] text-muted-foreground">
                              <span class="text-green-600 dark:text-green-400">+{{ pr.additions }}</span>
                              <span class="text-red-600 dark:text-red-400">-{{ pr.deletions }}</span>
                              <span class="ml-auto">{{ pr.timeAgo }}</span>
                            </div>

                            <!-- Merged PR → inline pipeline status -->
                            <div
                              v-if="pr.status === 'merged'"
                              class="mt-2 pt-2 border-t border-border"
                            >
                              <div
                                v-for="pipe in pipelines.filter(p => p.trigger.includes('#' + pr.id))"
                                :key="pipe.runId"
                                class="flex items-center gap-2 text-xs text-muted-foreground"
                              >
                                <component
                                  :is="getPipelineIcon(pipe.status).icon"
                                  :size="13"
                                  :class="getPipelineIcon(pipe.status).color"
                                />
                                <span class="text-foreground">{{ pipe.name }}</span>
                                <span class="text-[10px]">Run #{{ pipe.runId }}</span>
                                <span class="ml-auto text-[10px]">{{ pipe.duration }}</span>
                              </div>
                            </div>

                            <!-- Expanded detail -->
                            <div
                              v-if="expandedPrs[pr.id]"
                              class="mt-3 pt-3 border-t border-border space-y-2 text-xs"
                            >
                              <p class="text-foreground/80 leading-relaxed">{{ pr.description }}</p>
                              <div class="flex items-center gap-4 text-muted-foreground">
                                <span class="flex items-center gap-1">
                                  <FileCode :size="12" />
                                  {{ pr.filesChanged }} files
                                </span>
                                <span class="flex items-center gap-1">
                                  <CheckCircle2
                                    :size="12"
                                    :class="pr.checksPass ? 'text-green-500' : 'text-red-500'"
                                  />
                                  Checks {{ pr.checksPass ? 'passing' : 'failing' }}
                                </span>
                              </div>
                            </div>
                          </CardContent>
                        </Card>
                      </div>
                    </div>
                  </section>

                  <Separator />

                  <!-- Pipelines -->
                  <section class="space-y-2">
                    <div class="flex items-center gap-2">
                      <h3 class="text-xs font-semibold text-foreground uppercase tracking-wider">
                        Pipelines
                      </h3>
                      <Badge variant="outline" class="text-[10px] h-4 px-1.5">
                        {{ pipelines.length }}
                      </Badge>
                    </div>
                    <div class="space-y-1">
                      <div
                        v-for="pipe in pipelines"
                        :key="pipe.runId"
                        class="flex items-center gap-3 px-3 py-2 rounded-md border border-border bg-card text-sm"
                      >
                        <component
                          :is="getPipelineIcon(pipe.status).icon"
                          :size="16"
                          :class="getPipelineIcon(pipe.status).color"
                        />
                        <div class="flex flex-col min-w-0">
                          <span class="text-sm font-medium text-foreground truncate">{{ pipe.name }}</span>
                          <span class="text-[11px] text-muted-foreground">Run #{{ pipe.runId }}</span>
                        </div>
                        <Badge
                          variant="outline"
                          :class="cn(
                            'text-[10px] h-4 px-1.5 capitalize',
                            pipe.status === 'succeeded' && 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25',
                            pipe.status === 'running' && 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25',
                            pipe.status === 'failed' && 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25',
                          )"
                        >
                          {{ pipe.status }}
                        </Badge>
                        <span class="text-xs text-muted-foreground ml-auto shrink-0">{{ pipe.duration }}</span>
                        <Separator orientation="vertical" class="h-4" />
                        <span class="text-[11px] text-muted-foreground shrink-0">{{ pipe.trigger }}</span>
                      </div>
                    </div>
                  </section>

                </div>
              </ScrollArea>
            </TabsContent>
          </Tabs>
        </div>
  </TooltipProvider>
</template>
