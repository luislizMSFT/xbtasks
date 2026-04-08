<script setup lang="ts">
/**
 * Playground: Integrated Task List + Detail Pane
 *
 * Combines the ADO-styled hierarchical tree list (Phase 8) with the
 * ADO-style detail panel into a single unified view.
 *
 * Left panel:  ADO-styled tree with expand/collapse, type icons, state badges
 * Right panel: ADO-style detail pane with subtasks, PRs, notes, sync status
 *
 * Selecting a task in the tree drives the detail pane.
 */
import { ref, computed } from 'vue'
import type { Task } from '@/types'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Select, SelectContent, SelectItem, SelectTrigger,
} from '@/components/ui/select'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import {
  Circle, CircleDot, CheckCircle2, Octagon, Eye, XCircle,
  ChevronRight, ChevronDown, User, CalendarDays, Paintbrush,
  X, Trash2, Plus, ExternalLink, GitPullRequest, Folder,
  Upload, Download, Lock, Pencil, Save, Link2Off, Search,
  ListChecks, CirclePlay, CircleCheck, CircleX, Clock,
} from 'lucide-vue-next'
import {
  statusColor, statusClasses, statusBgColor,
  adoTypeColor, adoTypeIcon, adoStateClasses, prStatusClasses,
} from '@/lib/styles'

// ── ADO metadata per task ──
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
const mockProjects = [
  { id: 1, name: 'Xbox Platform' },
  { id: 2, name: 'Dashboard' },
  { id: 3, name: 'Infrastructure' },
]

// ── Detail data: rich info for each task ──
interface DetailSubtask {
  id: number; title: string; status: string; priority: string
  adoId: string; adoType: string; adoState: string; assignedTo: string
  syncStatus?: 'synced' | 'pending' | 'not-pulled'
  source?: 'ado' | 'personal'
}
interface DetailPR {
  id: number; title: string; prNumber: number; repo: string; status: string; sourceBranch: string
}
interface DetailNote {
  id: number; text: string; timestamp: string
}
interface TaskDetail {
  description: string
  subtasks: DetailSubtask[]
  prs: DetailPR[]
  notes: DetailNote[]
  dirtyFields: string[]
}

const taskDetails: Record<number, TaskDetail> = {
  100: {
    description: '<p>Top-level deliverable tracking all Xbox Platform Services work. Includes auth, rate limiting, and backend hardening.</p>',
    subtasks: [
      { id: 101, title: 'Implement auth token refresh', status: 'in_progress', priority: 'P1', adoId: '50010', adoType: 'User Story', adoState: 'Active', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'ado' },
      { id: 102, title: 'Add rate limiting middleware', status: 'todo', priority: 'P2', adoId: '', adoType: 'Task', adoState: 'New', assignedTo: '', syncStatus: 'synced', source: 'ado' },
    ],
    prs: [],
    notes: [
      { id: 1, text: 'Kickoff meeting scheduled. Prioritize auth refresh over rate limiting.', timestamp: '2026-04-02T09:00:00Z' },
    ],
    dirtyFields: [],
  },
  101: {
    description: '<p>Auto-refresh tokens before expiry. Uses MSAL.js v2 browser flow with silent refresh fallback.</p>',
    subtasks: [
      { id: 103, title: 'Write unit tests for auth module', status: 'done', priority: 'P1', adoId: '50011', adoType: 'Task', adoState: 'Closed', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'ado' },
      { id: 104, title: 'Handle token expiry edge cases', status: 'in_progress', priority: 'P1', adoId: '', adoType: 'Task', adoState: 'Active', assignedTo: 'Sarah K.', syncStatus: 'pending', source: 'ado' },
      { id: 1010, title: 'Draft architecture diagram for team', status: 'todo', priority: 'P2', adoId: '', adoType: '', adoState: '', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'personal' },
    ],
    prs: [
      { id: 100, title: 'feat: add Azure AD provider', prNumber: 342, repo: 'platform-api', status: 'active', sourceBranch: 'refs/heads/feat/auth-refresh' },
      { id: 101, title: 'feat: token refresh middleware', prNumber: 343, repo: 'platform-api', status: 'draft', sourceBranch: 'refs/heads/feat/token-refresh' },
    ],
    notes: [
      { id: 2, text: 'Using MSAL.js v2 for browser flow. Need to check redirect URI config for Wails.', timestamp: '2026-04-07T15:30:00Z' },
      { id: 3, text: 'Token refresh interval: 5 min before expiry. Check Azure AD docs for silent refresh limits.', timestamp: '2026-04-06T10:00:00Z' },
    ],
    dirtyFields: ['description', 'priority'],
  },
  200: {
    description: '<p>Revamp the dashboard layout with stat cards, activity feed, and PR summary widget.</p>',
    subtasks: [
      { id: 201, title: 'Add stat cards row', status: 'done', priority: 'P2', adoId: '', adoType: 'Task', adoState: 'Closed', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'ado' },
      { id: 202, title: 'Implement PR summary widget', status: 'blocked', priority: 'P1', adoId: '60010', adoType: 'Bug', adoState: 'Active', assignedTo: 'Mike R.', syncStatus: 'not-pulled', source: 'ado' },
    ],
    prs: [
      { id: 200, title: 'feat: dashboard stat cards', prNumber: 350, repo: 'xb-tasks', status: 'completed', sourceBranch: 'refs/heads/fix/dashboard-perf' },
    ],
    notes: [
      { id: 4, text: 'Suspect ADO sync call is blocking the render. Try deferring it to after mount.', timestamp: '2026-04-07T10:00:00Z' },
    ],
    dirtyFields: ['status'],
  },
  300: {
    description: 'Pipeline times out on large repos. Investigate runner config and caching strategy.',
    subtasks: [],
    prs: [],
    notes: [
      { id: 5, text: 'Could be related to the npm cache not being restored correctly.', timestamp: '2026-04-07T14:00:00Z' },
    ],
    dirtyFields: [],
  },
}

// Fallback for tasks without detail data
const emptyDetail: TaskDetail = { description: '', subtasks: [], prs: [], notes: [], dirtyFields: [] }

// ── Mock pipelines ──
interface MockPipeline { id: number; name: string; status: string; result: string; sourceBranch: string }
const mockPipelines: MockPipeline[] = [
  { id: 1, name: 'CI Build', status: 'completed', result: 'succeeded', sourceBranch: 'refs/heads/feat/auth-refresh' },
  { id: 2, name: 'Integration Tests', status: 'completed', result: 'failed', sourceBranch: 'refs/heads/feat/auth-refresh' },
  { id: 3, name: 'CI Build', status: 'inprogress', result: '', sourceBranch: 'refs/heads/fix/dashboard-perf' },
]

function pipelinesForPr(branch: string): MockPipeline[] {
  return mockPipelines.filter(p => p.sourceBranch === branch)
}
function pipelineIcon(p: MockPipeline) {
  if (p.result === 'succeeded') return CircleCheck
  if (p.result === 'failed') return CircleX
  if (p.status === 'inprogress') return CirclePlay
  return Clock
}
function pipelineColor(p: MockPipeline) {
  if (p.result === 'succeeded') return 'text-emerald-500'
  if (p.result === 'failed') return 'text-red-500'
  if (p.status === 'inprogress') return 'text-blue-500'
  return 'text-amber-500'
}

// ── Tree helpers ──
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

// ── Selection ──
const selectedId = ref<number>(101)
const selectedTask = computed(() => mockTasks.find(t => t.id === selectedId.value)!)
const detail = computed(() => taskDetails[selectedId.value] || emptyDetail)

// ── Detail helpers ──
function meta(id: number): AdoMeta { return adoMeta[id] || { type: '', state: '', adoId: '' } }
function isPersonalTask(task: Task): boolean { return !task.adoId && !task.projectId }

function adoNumber(adoId: string): string {
  if (!adoId) return ''
  const m = adoId.match(/\d+/)
  return m ? `#${m[0]}` : adoId
}

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

const priorityDot: Record<string, string> = {
  P0: 'bg-red-500', P1: 'bg-orange-500', P2: 'bg-amber-500', P3: 'bg-zinc-400',
}

function prIconColor(status: string) {
  return status === 'completed' ? 'text-violet-500' : status === 'draft' ? 'text-zinc-400' : 'text-emerald-500'
}

const currentUser = 'Luis L.'
const subtaskFilter = ref<'all' | 'ado' | 'personal' | 'mine'>('all')

function isPersonalSubtask(st: DetailSubtask): boolean { return st.source === 'personal' }
function isOtherPerson(st: DetailSubtask): boolean { return !!st.assignedTo && st.assignedTo !== currentUser }

const filteredSubtasks = computed(() => {
  const subs = detail.value.subtasks
  switch (subtaskFilter.value) {
    case 'mine': return subs.filter(st => !st.assignedTo || st.assignedTo === currentUser)
    case 'ado': return subs.filter(st => st.source !== 'personal')
    case 'personal': return subs.filter(st => st.source === 'personal')
    default: return subs
  }
})

function cycleSubtaskFilter() {
  const order: Array<'all' | 'ado' | 'personal' | 'mine'> = ['all', 'mine', 'ado', 'personal']
  const idx = order.indexOf(subtaskFilter.value)
  subtaskFilter.value = order[(idx + 1) % order.length]
}

function togglePersonalDone(st: DetailSubtask) {
  st.status = st.status === 'done' ? 'todo' : 'done'
}

const detailSubtaskProgress = computed(() => {
  const total = detail.value.subtasks.length
  const done = detail.value.subtasks.filter(s => s.status === 'done').length
  return { done, total, percent: total ? (done / total) * 100 : 0 }
})

// Inline add subtask
const addingSubtask = ref(false)
const newSubtaskTitle = ref('')
function addSubtask() { addingSubtask.value = true; newSubtaskTitle.value = '' }
function confirmAddSubtask() {
  if (!newSubtaskTitle.value.trim()) return
  detail.value.subtasks.push({
    id: Date.now(), title: newSubtaskTitle.value.trim(),
    status: 'todo', priority: 'P2', adoId: '', adoType: '', adoState: '',
    assignedTo: currentUser, syncStatus: 'synced', source: 'personal',
  })
  addingSubtask.value = false; newSubtaskTitle.value = ''
}

// Notes
const newNote = ref('')
function addNote() {
  if (!newNote.value.trim()) return
  detail.value.notes.unshift({ id: Date.now(), text: newNote.value.trim(), timestamp: new Date().toISOString() })
  newNote.value = ''
}

// Description editing
const editingDescription = ref(false)
const editedDescription = ref('')
function startEditDescription() {
  editedDescription.value = detail.value.description.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim()
  editingDescription.value = true
}
function saveDescription() {
  detail.value.description = editedDescription.value
  editingDescription.value = false
}

// Detail panel visibility
const showDetail = ref(true)
</script>

<template>
  <div class="h-screen flex flex-col bg-background text-foreground">
    <!-- Header bar -->
    <div class="h-10 shrink-0 border-b border-border flex items-center px-4 gap-3">
      <Paintbrush :size="14" class="text-orange-500" />
      <span class="text-sm font-semibold">Playground: Integrated Task View</span>
      <Badge variant="secondary" class="text-[9px] h-4">Phase 8</Badge>
      <div class="flex-1" />
      <Button variant="outline" size="sm" class="h-7 text-xs" @click="showDetail = !showDetail">
        {{ showDetail ? 'Hide Detail' : 'Show Detail' }}
      </Button>
      <Button variant="outline" size="sm" class="h-7 text-xs" @click="expandedNodes = new Set(mockTasks.map(t => t.id))">Expand All</Button>
      <Button variant="outline" size="sm" class="h-7 text-xs" @click="expandedNodes = new Set()">Collapse All</Button>
    </div>

    <div class="flex-1 flex min-h-0">
      <!-- ═══════ LEFT: ADO-styled task tree ═══════ -->
      <div class="flex flex-col min-h-0 border-r border-border" :class="showDetail ? 'w-[420px] shrink-0' : 'flex-1'">
        <div class="px-4 py-2 border-b border-border/50 text-xs font-semibold text-muted-foreground">
          Task List ({{ rootTasks.length }} roots · {{ mockTasks.length }} total)
        </div>
        <ScrollArea class="flex-1 min-h-0">
          <template v-for="task in rootTasks" :key="task.id">
            <!-- Root node -->
            <div
              class="group cursor-pointer hover:bg-muted/50 transition-colors"
              :class="selectedId === task.id ? 'bg-primary/5 border-l-2 border-l-primary' : 'border-l-2 border-l-transparent'"
              @click="selectedId = task.id; showDetail = true"
            >
              <!-- Row 1: chevron + type icon + title + badges -->
              <div class="flex items-center gap-2 px-3 pt-2.5 pb-0.5">
                <button
                  v-if="hasChildren(task.id)"
                  class="shrink-0 p-0.5 hover:bg-muted rounded"
                  @click.stop="toggleExpand(task.id)"
                >
                  <component :is="expandedNodes.has(task.id) ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
                </button>
                <span v-else class="w-5 shrink-0" />

                <component
                  v-if="meta(task.id).type"
                  :is="adoTypeIcon(meta(task.id).type)"
                  :size="15"
                  :class="adoTypeColor(meta(task.id).type)"
                  class="shrink-0"
                />
                <User v-else-if="isPersonalTask(task)" :size="14" class="text-muted-foreground/60 shrink-0" />
                <component v-else :is="statusIcon(task.status)" :size="15" :class="statusColor(task.status)" class="shrink-0" />

                <span
                  class="text-sm truncate flex-1 min-w-0"
                  :class="[
                    task.status === 'done' ? 'line-through text-muted-foreground' : 'text-foreground',
                    hasChildren(task.id) ? 'font-medium' : '',
                  ]"
                >
                  {{ task.title }}
                </span>

                <Badge variant="outline" class="text-[10px] h-4 px-1.5 shrink-0"
                  :class="meta(task.id).state ? adoStateClasses(meta(task.id).state) : statusClasses(task.status)"
                >
                  {{ meta(task.id).state || statusLabel(task.status) }}
                </Badge>

                <template v-if="subtaskProgress(task.id)">
                  <span class="text-[10px] text-muted-foreground tabular-nums shrink-0">
                    {{ subtaskProgress(task.id)!.done }}/{{ subtaskProgress(task.id)!.total }}
                  </span>
                  <div class="w-12 h-1 bg-muted rounded-full overflow-hidden shrink-0">
                    <div class="h-full bg-green-500 rounded-full transition-all" :style="{ width: subtaskProgress(task.id)!.pct + '%' }" />
                  </div>
                </template>
              </div>

              <!-- Row 2: metadata -->
              <div class="flex items-center gap-2 pb-2 text-[11px]" :style="{ paddingLeft: '68px', paddingRight: '12px' }">
                <span v-if="meta(task.id).adoId" class="text-muted-foreground/50 tabular-nums">{{ adoNumber(meta(task.id).adoId) }}</span>
                <Badge v-if="isPersonalTask(task)" variant="outline" class="text-[9px] h-3.5 px-1 border-dashed text-muted-foreground/60">
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

            <!-- Children (level 1) -->
            <template v-if="expandedNodes.has(task.id)">
              <template v-for="child in getChildren(task.id)" :key="child.id">
                <div
                  class="group cursor-pointer hover:bg-muted/50 transition-colors border-l-2"
                  :class="selectedId === child.id ? 'bg-primary/5 border-l-primary' : 'border-l-muted-foreground/10'"
                  @click="selectedId = child.id; showDetail = true"
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
                    <User v-else-if="isPersonalTask(child)" :size="13" class="text-muted-foreground/60 shrink-0" />
                    <component v-else :is="statusIcon(child.status)" :size="14" :class="statusColor(child.status)" class="shrink-0" />

                    <span class="text-sm truncate flex-1 min-w-0" :class="child.status === 'done' ? 'line-through text-muted-foreground' : 'text-foreground'">
                      {{ child.title }}
                    </span>

                    <Badge variant="outline" class="text-[10px] h-4 px-1.5 shrink-0"
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

                  <div class="flex items-center gap-2 pb-1.5 text-[10px]" :style="{ paddingLeft: '80px', paddingRight: '12px' }">
                    <span v-if="meta(child.id).adoId" class="text-muted-foreground/50 tabular-nums">{{ adoNumber(meta(child.id).adoId) }}</span>
                    <Badge v-if="isPersonalTask(child)" variant="outline" class="text-[9px] h-3.5 px-1 border-dashed text-muted-foreground/60">
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

                <!-- Grandchildren (level 2) -->
                <template v-if="expandedNodes.has(child.id)">
                  <div
                    v-for="gc in getChildren(child.id)"
                    :key="gc.id"
                    class="group flex items-center gap-2 py-1.5 cursor-pointer hover:bg-muted/50 transition-colors border-l-2"
                    :class="selectedId === gc.id ? 'bg-primary/5 border-l-primary' : 'border-l-muted-foreground/5'"
                    :style="{ paddingLeft: '60px', paddingRight: '12px' }"
                    @click="selectedId = gc.id; showDetail = true"
                  >
                    <span class="w-5 shrink-0" />
                    <component
                      v-if="meta(gc.id).type"
                      :is="adoTypeIcon(meta(gc.id).type)"
                      :size="12"
                      :class="adoTypeColor(meta(gc.id).type)"
                      class="shrink-0"
                    />
                    <User v-else-if="isPersonalTask(gc)" :size="12" class="text-muted-foreground/60 shrink-0" />
                    <component v-else :is="statusIcon(gc.status)" :size="12" :class="statusColor(gc.status)" class="shrink-0" />

                    <span class="text-xs truncate flex-1 min-w-0" :class="gc.status === 'done' ? 'line-through text-muted-foreground' : 'text-foreground'">
                      {{ gc.title }}
                    </span>

                    <Badge variant="outline" class="text-[9px] h-3.5 px-1 shrink-0"
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

      <!-- ═══════ RIGHT: ADO-style detail pane ═══════ -->
      <div v-if="showDetail && selectedTask" class="flex-1 flex flex-col min-h-0">
        <div class="flex-1 flex flex-col min-h-0">
          <!-- Header: type icon + title + close -->
          <div class="shrink-0 border-b border-border px-4 py-3 space-y-2">
            <div class="flex items-center gap-2">
              <component
                v-if="meta(selectedTask.id).type"
                :is="adoTypeIcon(meta(selectedTask.id).type)"
                :size="16"
                :class="adoTypeColor(meta(selectedTask.id).type)"
              />
              <h2 class="text-sm font-semibold text-foreground leading-snug flex-1 line-clamp-2">{{ selectedTask.title }}</h2>
              <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="showDetail = false">
                <X :size="16" />
              </Button>
            </div>

            <!-- Metadata badges -->
            <div class="flex items-center gap-2 flex-wrap">
              <Badge variant="outline" class="text-[10px] h-5 px-1.5">
                <span :class="cn('size-1.5 rounded-full mr-1', statusBgColor(selectedTask.status))" />
                {{ statusLabel(selectedTask.status) }}
                <span v-if="detail.dirtyFields.includes('status')" class="size-1 rounded-full bg-amber-500 ml-1" title="Modified locally" />
              </Badge>
              <Badge variant="outline" class="text-[10px] h-5 px-1.5">
                <span :class="cn('size-1.5 rounded-full mr-1', priorityDot[selectedTask.priority] ?? 'bg-zinc-400')" />
                {{ selectedTask.priority }}
                <span v-if="detail.dirtyFields.includes('priority')" class="size-1 rounded-full bg-amber-500 ml-1" title="Modified locally" />
              </Badge>
              <span v-if="selectedTask.dueDate" class="text-[10px] text-muted-foreground flex items-center gap-1">
                <CalendarDays :size="10" />
                {{ selectedTask.dueDate }}
              </span>
            </div>

            <!-- Project context -->
            <div v-if="selectedTask.projectId && projectNames[selectedTask.projectId]" class="flex items-center gap-2 px-2 py-1.5 rounded-md bg-muted/40 border border-border/50">
              <component
                v-if="meta(selectedTask.id).type"
                :is="adoTypeIcon(meta(selectedTask.id).type)"
                :size="13"
                :class="adoTypeColor(meta(selectedTask.id).type)"
                class="shrink-0"
              />
              <Folder v-else :size="13" class="text-muted-foreground shrink-0" />
              <span class="text-xs font-medium text-foreground truncate flex-1">{{ projectNames[selectedTask.projectId] }}</span>
              <Select model-value="keep">
                <SelectTrigger size="sm" class="h-5 text-[10px] gap-0.5 w-auto px-1.5 border-none shadow-none bg-transparent hover:bg-muted">
                  <Search :size="10" class="shrink-0 mr-0.5" />
                  <span class="text-muted-foreground">Switch</span>
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="p in mockProjects" :key="p.id" :value="String(p.id)">{{ p.name }}</SelectItem>
                </SelectContent>
              </Select>
              <button class="text-[10px] text-muted-foreground hover:text-destructive transition-colors px-1" title="Unlink project">
                <Link2Off :size="11" />
              </button>
            </div>
            <div v-else class="flex items-center gap-2 px-2 py-1.5 rounded-md border border-dashed border-border/50">
              <User :size="12" class="text-muted-foreground/40 shrink-0" />
              <span class="text-[10px] text-muted-foreground/60">Personal task — not linked to a project</span>
            </div>

            <!-- ADO integration row -->
            <div v-if="selectedTask.adoId" class="flex items-center gap-2 px-2 py-1 rounded-md bg-blue-500/5 border border-blue-500/15">
              <AzureDevOpsIcon :size="12" class="text-blue-500 shrink-0" />
              <span class="text-[10px] text-blue-500 tabular-nums font-medium">#{{ selectedTask.adoId }}</span>
              <span class="text-[9px] text-muted-foreground/50">·</span>
              <span v-if="detail.dirtyFields.length === 0" class="text-[9px] text-emerald-600 flex items-center gap-0.5">
                <CircleCheck :size="9" /> Synced
              </span>
              <span v-else class="text-[9px] text-amber-600 flex items-center gap-0.5">
                <span class="size-1 rounded-full bg-amber-500" />
                {{ detail.dirtyFields.length }} pending
              </span>
              <div class="flex-1" />
              <Button
                v-if="detail.dirtyFields.length > 0"
                variant="outline" size="sm"
                class="h-5 text-[9px] gap-0.5 px-1.5 text-amber-600 border-amber-500/30 hover:bg-amber-500/10"
              >
                <Upload :size="10" /> Push {{ detail.dirtyFields.length }}
              </Button>
              <Button variant="ghost" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5 text-blue-500 hover:text-blue-600">
                <ExternalLink :size="10" /> Open
              </Button>
            </div>
            <div v-else-if="isPersonalTask(selectedTask)" class="flex items-center gap-2 px-2 py-1 rounded-md border border-dashed border-border/50">
              <User :size="12" class="text-muted-foreground/40 shrink-0" />
              <span class="text-[10px] text-muted-foreground/60">Personal task — not linked to ADO</span>
            </div>

            <!-- Subtask progress bar -->
            <div v-if="detail.subtasks.length > 0">
              <div class="h-1 w-full rounded-full bg-muted">
                <div class="h-1 rounded-full bg-blue-500 transition-all" :style="{ width: detailSubtaskProgress.percent + '%' }" />
              </div>
            </div>
          </div>

          <!-- Scrollable content -->
          <ScrollArea class="flex-1 min-h-0">
            <div class="flex flex-col">

              <!-- Subtasks -->
              <div v-if="detail.subtasks.length > 0 || addingSubtask" class="border-b border-border px-4 py-2">
                <div class="flex items-center justify-between mb-1.5">
                  <div class="flex items-center gap-2">
                    <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Subtasks</h3>
                    <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ detailSubtaskProgress.done }}/{{ detailSubtaskProgress.total }}</Badge>
                  </div>
                  <div class="flex items-center gap-2">
                    <button
                      class="text-[10px] flex items-center gap-1 px-1.5 py-0.5 rounded transition-colors"
                      :class="subtaskFilter !== 'all' ? 'bg-primary/10 text-primary font-medium' : 'text-muted-foreground hover:text-foreground'"
                      @click="cycleSubtaskFilter"
                      :title="`Filter: ${subtaskFilter}`"
                    >
                      <component :is="subtaskFilter === 'personal' ? ListChecks : subtaskFilter === 'ado' ? Lock : User" :size="10" />
                      {{ subtaskFilter === 'all' ? 'All' : subtaskFilter === 'mine' ? 'Mine' : subtaskFilter === 'ado' ? 'ADO' : 'Personal' }}
                    </button>
                    <button class="text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1" @click="addSubtask">
                      <Plus :size="11" /> Add
                    </button>
                  </div>
                </div>
                <div v-if="filteredSubtasks.length > 0" class="flex flex-col gap-px">
                  <div v-for="st in filteredSubtasks" :key="st.id" class="rounded group"
                    :class="st.syncStatus === 'not-pulled' ? 'opacity-50' : ''"
                  >
                    <div
                      class="flex items-center gap-2 py-1.5 px-2 cursor-pointer transition-colors"
                      :class="st.syncStatus === 'not-pulled' ? 'hover:bg-muted/30' : 'hover:bg-muted/50'"
                    >
                      <template v-if="isPersonalSubtask(st)">
                        <button
                          class="size-3.5 rounded-[3px] border-[1.5px] shrink-0 flex items-center justify-center transition-colors"
                          :class="st.status === 'done' ? 'bg-emerald-500 border-emerald-500' : 'border-primary/60 hover:border-primary'"
                          @click.stop="togglePersonalDone(st)"
                        >
                          <svg v-if="st.status === 'done'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 12 12" class="size-2.5 text-white" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="2 6 5 9 10 3" />
                          </svg>
                        </button>
                      </template>
                      <template v-else>
                        <component
                          v-if="st.adoType"
                          :is="adoTypeIcon(st.adoType)"
                          :size="13"
                          :class="st.syncStatus === 'not-pulled' ? 'text-muted-foreground/40' : adoTypeColor(st.adoType)"
                          class="shrink-0"
                        />
                        <span v-else :class="cn('size-3.5 rounded-[3px] border-[1.5px] shrink-0 flex items-center justify-center',
                          st.status === 'done' ? 'bg-emerald-500 border-emerald-500' : 'border-muted-foreground/60')">
                          <svg v-if="st.status === 'done'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 12 12" class="size-2.5 text-white" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="2 6 5 9 10 3" />
                          </svg>
                        </span>
                      </template>
                      <span :class="cn('text-xs truncate flex-1',
                        st.status === 'done' && 'line-through text-muted-foreground',
                        st.syncStatus === 'not-pulled' && 'text-muted-foreground/60')">
                        {{ st.title }}
                      </span>
                      <span v-if="isPersonalSubtask(st)" class="text-[8px] px-1 py-0.5 rounded bg-primary/8 text-primary/70 shrink-0 border border-primary/10">personal</span>
                      <Badge v-if="st.adoState && !isPersonalSubtask(st)" variant="outline" :class="['text-[9px] h-3.5 px-1', adoStateClasses(st.adoState)]">
                        {{ st.adoState }}
                      </Badge>
                      <span v-if="st.syncStatus === 'pending' && !isPersonalSubtask(st)" class="size-1 rounded-full bg-amber-500 shrink-0" title="Pending push" />
                      <span v-if="st.syncStatus === 'not-pulled'" class="text-[8px] text-muted-foreground/40 shrink-0">not pulled</span>
                      <span :class="cn('size-1.5 rounded-full shrink-0', priorityDot[st.priority] ?? 'bg-zinc-400')" />
                      <span v-if="st.adoId && !isPersonalSubtask(st)" class="text-[9px] text-blue-500/60 tabular-nums shrink-0">#{{ st.adoId }}</span>
                      <span
                        v-if="st.assignedTo"
                        class="text-[9px] shrink-0 rounded-full h-4 px-1.5 flex items-center gap-0.5"
                        :class="isOtherPerson(st)
                          ? 'bg-violet-500/10 text-violet-600 dark:text-violet-400 border border-violet-500/20'
                          : 'bg-muted text-muted-foreground'"
                      >
                        {{ st.assignedTo }}
                      </span>
                      <button
                        v-if="isPersonalSubtask(st)"
                        class="opacity-0 group-hover:opacity-100 text-muted-foreground/40 hover:text-red-500 transition-all shrink-0"
                        @click.stop="detail.subtasks = detail.subtasks.filter(s => s.id !== st.id)"
                        title="Remove"
                      >
                        <Trash2 :size="11" />
                      </button>
                    </div>
                    <!-- Sync actions on hover -->
                    <div v-if="st.syncStatus === 'pending' && !isPersonalSubtask(st)"
                      class="flex items-center gap-1.5 pl-7 pb-1 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <button class="text-[9px] text-amber-600 hover:text-amber-700 flex items-center gap-0.5">
                        <Upload :size="9" /> Push to ADO
                      </button>
                    </div>
                    <div v-if="st.syncStatus === 'not-pulled'"
                      class="flex items-center gap-1.5 pl-7 pb-1 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <button class="text-[9px] text-blue-500 hover:text-blue-600 flex items-center gap-0.5">
                        <Download :size="9" /> Pull from ADO
                      </button>
                    </div>
                  </div>
                </div>
                <!-- Add subtask inline -->
                <div v-if="addingSubtask" class="flex items-center gap-2 mt-1.5 px-2">
                  <button class="size-3.5 rounded-[3px] border-[1.5px] border-primary/40 shrink-0 flex items-center justify-center" />
                  <Input v-model="newSubtaskTitle" placeholder="New personal subtask..." class="h-7 text-xs flex-1"
                    @keydown.enter="confirmAddSubtask" @keydown.escape="addingSubtask = false"
                  />
                  <Button variant="outline" size="sm" class="h-7 px-2 text-[10px] shrink-0" @click="confirmAddSubtask" :disabled="!newSubtaskTitle.trim()">Add</Button>
                  <Button variant="ghost" size="sm" class="h-7 px-1.5 text-[10px]" @click="addingSubtask = false">Cancel</Button>
                </div>
                <p v-if="filteredSubtasks.length === 0 && !addingSubtask" class="text-[11px] text-muted-foreground/40 italic py-2 text-center">
                  {{ subtaskFilter === 'mine' ? 'No subtasks assigned to you' : subtaskFilter === 'personal' ? 'No personal subtasks' : subtaskFilter === 'ado' ? 'No ADO subtasks' : 'No subtasks' }}
                </p>
              </div>

              <!-- Description -->
              <div class="border-b border-border px-4 py-3">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-1.5">
                    <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
                    <span v-if="detail.dirtyFields.includes('description')" class="size-1 rounded-full bg-amber-500" title="Modified locally" />
                  </div>
                  <button v-if="!editingDescription" class="text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1" @click="startEditDescription">
                    <Pencil :size="10" /> Edit
                  </button>
                  <div v-else class="flex items-center gap-1">
                    <button class="text-[10px] text-emerald-600 hover:text-emerald-700 flex items-center gap-0.5" @click="saveDescription">
                      <Save :size="10" /> Save
                    </button>
                    <button class="text-[10px] text-muted-foreground hover:text-foreground" @click="editingDescription = false">Cancel</button>
                  </div>
                </div>
                <div v-if="editingDescription">
                  <textarea v-model="editedDescription"
                    class="w-full min-h-[80px] text-xs bg-background border border-border rounded-md p-2 resize-y focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="Add a description..."
                  />
                </div>
                <div v-else-if="detail.description"
                  class="text-xs text-foreground prose prose-sm max-w-none [&_*]:text-xs [&_*]:text-foreground cursor-pointer hover:bg-muted/30 rounded p-1 -m-1 transition-colors"
                  v-html="detail.description"
                  @click="startEditDescription"
                />
                <p v-else class="text-[11px] text-muted-foreground/40 italic cursor-pointer hover:text-muted-foreground" @click="startEditDescription">
                  Click to add description...
                </p>
              </div>

              <!-- Pull Requests -->
              <div class="border-b border-border px-4 py-3">
                <div class="flex items-center gap-2 mb-2">
                  <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Pull Requests</h3>
                  <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ detail.prs.length }}</Badge>
                </div>
                <div class="flex flex-col gap-1">
                  <div v-for="pr in detail.prs" :key="pr.id" class="rounded hover:bg-muted/50 cursor-pointer">
                    <div class="flex items-center gap-2 py-1.5 px-1">
                      <GitPullRequest :size="14" :class="cn('shrink-0', prIconColor(pr.status))" />
                      <span class="text-[13px] truncate flex-1">{{ pr.title }}</span>
                      <span class="text-[11px] text-muted-foreground shrink-0">#{{ pr.prNumber }}</span>
                      <Badge variant="outline" :class="cn('text-[10px] capitalize shrink-0 h-4 px-1.5', prStatusClasses(pr.status))">
                        {{ pr.status }}
                      </Badge>
                    </div>
                    <!-- Pipeline status for this PR -->
                    <div v-if="pipelinesForPr(pr.sourceBranch).length > 0" class="flex items-center gap-3 pl-6 pb-1.5">
                      <div v-for="pipe in pipelinesForPr(pr.sourceBranch)" :key="pipe.id" class="flex items-center gap-1 text-[9px]">
                        <component :is="pipelineIcon(pipe)" :size="10" :class="pipelineColor(pipe)" />
                        <span class="text-muted-foreground">{{ pipe.name }}</span>
                      </div>
                    </div>
                  </div>
                  <p v-if="detail.prs.length === 0" class="text-xs text-muted-foreground italic py-1 px-1">No linked PRs</p>
                </div>
              </div>

              <!-- Notes -->
              <div class="px-4 py-3">
                <div class="flex items-center gap-2 mb-2">
                  <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Notes</h3>
                  <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ detail.notes.length }}</Badge>
                  <button v-if="selectedTask.adoId" class="inline-flex items-center gap-1 text-[9px] text-blue-500 hover:text-blue-600 transition-colors ml-1">
                    <AzureDevOpsIcon :size="9" /> ADO Discussion <ExternalLink :size="8" />
                  </button>
                </div>
                <div class="mb-2 flex gap-1.5">
                  <Input placeholder="Add a note..." class="h-7 text-xs flex-1"
                    :model-value="newNote" @update:model-value="newNote = String($event)" @keydown.enter="addNote"
                  />
                  <Button variant="outline" size="sm" class="h-7 px-2 text-[10px] shrink-0" @click="addNote" :disabled="!newNote.trim()">
                    <Plus :size="11" class="mr-0.5" /> Add
                  </Button>
                </div>
                <div class="flex flex-col gap-2">
                  <div v-for="note in detail.notes" :key="note.id" class="group relative px-3 py-2 rounded-md bg-muted/40 border border-border/50">
                    <p class="text-xs text-foreground leading-relaxed pr-10">{{ note.text }}</p>
                    <span class="text-[9px] text-muted-foreground/60 mt-1 block">
                      {{ new Date(note.timestamp).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) }}
                      · {{ new Date(note.timestamp).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' }) }}
                    </span>
                    <div class="absolute top-1.5 right-1.5 flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button class="p-1 rounded hover:bg-red-500/10 text-muted-foreground hover:text-red-500"
                        @click="detail.notes = detail.notes.filter(n => n.id !== note.id)"
                      >
                        <Trash2 :size="10" />
                      </button>
                    </div>
                  </div>
                  <p v-if="detail.notes.length === 0" class="text-xs text-muted-foreground/40 text-center py-3 italic">
                    No notes yet. Add one above.
                  </p>
                </div>
              </div>

            </div>
          </ScrollArea>

          <!-- Footer -->
          <div class="shrink-0 border-t border-border flex items-center justify-between px-4 py-1.5">
            <span class="text-[11px] text-muted-foreground">
              Created {{ new Date(selectedTask.createdAt).toLocaleDateString() }} · Updated {{ new Date(selectedTask.updatedAt).toLocaleDateString() }}
            </span>
            <Button variant="ghost" size="sm" class="h-6 text-[11px] text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1">
              <Trash2 :size="12" /> Delete
            </Button>
          </div>
        </div>
      </div>

      <!-- No selection placeholder -->
      <div v-else-if="showDetail && !selectedTask" class="flex-1 flex items-center justify-center text-muted-foreground text-sm">
        Select a task to view details
      </div>
    </div>

    <!-- Footer legend -->
    <div class="shrink-0 border-t border-border px-4 py-1.5 flex items-center gap-4 text-[10px] text-muted-foreground">
      <span><strong>Left:</strong> ADO-styled hierarchical tree</span>
      <span>|</span>
      <span><strong>Right:</strong> ADO-style detail pane with subtasks, PRs, notes, sync status</span>
      <span>|</span>
      <span>{{ rootTasks.length }} roots · {{ mockTasks.length }} total</span>
    </div>
  </div>
</template>
