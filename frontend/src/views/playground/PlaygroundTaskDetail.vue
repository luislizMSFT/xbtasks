<script setup lang="ts">
/**
 * Playground: ADO-inspired TaskDetail layout.
 * Pure UI — no stores, no API calls, just mock data for visual comparison.
 *
 * Left: "Current" layout (how TaskDetail looks today)
 * Right: "ADO-style" layout (adapted from AdoView's detail panel)
 */
import { ref, computed } from 'vue'
import { cn } from '@/lib/utils'
import { statusBgColor, adoTypeIcon, adoTypeColor, adoStateClasses, adoPriorityClasses, prStatusClasses } from '@/lib/styles'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import {
  X, Trash2, Plus, ChevronDown, ChevronRight, ExternalLink,
  GitPullRequest, Lock, Upload, Folder, CalendarDays, Link2Off,
  Loader2, CheckCircle2, Download, Link, Search, User,
  Pencil, CirclePlay, CircleCheck, CircleX, Clock, Save,
  ListChecks, SquareCheckBig,
} from 'lucide-vue-next'

// ── Mock data ──
interface MockSubtask {
  id: number; title: string; status: string; priority: string
  adoId: string; adoType: string; adoState: string; assignedTo: string
  children?: MockSubtask[]
  syncStatus?: 'synced' | 'pending' | 'not-pulled'
  source?: 'ado' | 'personal' // personal = local-only, not linked to ADO
}

interface MockTask {
  id: number
  title: string
  description: string
  status: string
  priority: string
  projectId: number | null
  projectName: string
  dueDate: string
  adoId: string
  adoType: string
  adoState: string
  subtasks: MockSubtask[]
  prs: { id: number; title: string; prNumber: number; repo: string; status: string; sourceBranch: string }[]
  notes: { id: number; text: string; timestamp: string }[]
  createdAt: string
  updatedAt: string
  dirtyFields: string[] // fields modified locally but not yet pushed to ADO
}

// Pipeline status linked to PR branches
interface MockPipeline {
  id: number; name: string; status: string; result: string
  sourceBranch: string; finishTime: string | null; url: string
}
const mockPipelines: MockPipeline[] = [
  { id: 1, name: 'CI Build', status: 'completed', result: 'succeeded', sourceBranch: 'refs/heads/feat/auth-refresh', finishTime: '2026-04-07T14:30:00Z', url: '#' },
  { id: 2, name: 'Integration Tests', status: 'completed', result: 'failed', sourceBranch: 'refs/heads/feat/auth-refresh', finishTime: '2026-04-07T14:35:00Z', url: '#' },
  { id: 3, name: 'CI Build', status: 'inprogress', result: '', sourceBranch: 'refs/heads/fix/dashboard-perf', finishTime: null, url: '#' },
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

const mockTasks: MockTask[] = [
  {
    id: 1,
    title: 'Implement user authentication flow with SSO integration',
    description: '<div><p>We need to implement the full authentication flow including:</p><ul><li>Azure AD SSO integration</li><li>Token refresh logic</li><li>Session management</li></ul><p>See <a href="#">design doc</a> for details.</p></div>',
    status: 'in_progress', priority: 'P0',
    projectId: 1, projectName: 'Platform Modernization',
    dueDate: '2026-04-15', adoId: '48291', adoType: 'User Story', adoState: 'Active',
    subtasks: [
      { id: 10, title: 'Setup Azure AD app registration', status: 'done', priority: 'P1', adoId: '48292', adoType: 'Task', adoState: 'Closed', assignedTo: 'Luis L.', children: [], syncStatus: 'synced', source: 'ado' },
      { id: 11, title: 'Implement token refresh middleware', status: 'in_progress', priority: 'P0', adoId: '48293', adoType: 'Task', adoState: 'Active', assignedTo: 'Luis L.', children: [
        { id: 110, title: 'Add MSAL.js v2 browser flow', status: 'done', priority: 'P0', adoId: '48295', adoType: 'Task', adoState: 'Closed', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'ado' },
        { id: 111, title: 'Handle silent refresh fallback', status: 'in_progress', priority: 'P1', adoId: '48296', adoType: 'Task', adoState: 'Active', assignedTo: 'Sarah K.', syncStatus: 'pending', source: 'ado' },
      ], syncStatus: 'pending', source: 'ado' },
      { id: 12, title: 'Add session timeout handling', status: 'todo', priority: 'P2', adoId: '48294', adoType: 'Task', adoState: 'New', assignedTo: 'Mike R.', syncStatus: 'not-pulled', source: 'ado' },
      { id: 13, title: 'Write integration tests', status: 'todo', priority: 'P2', adoId: '', adoType: '', adoState: '', assignedTo: '', syncStatus: 'synced', source: 'ado' },
      { id: 14, title: 'Review security threat model', status: 'in_progress', priority: 'P1', adoId: '48297', adoType: 'Task', adoState: 'Active', assignedTo: 'Dana W.', syncStatus: 'synced', source: 'ado' },
      // Personal subtasks (local only, not in ADO)
      { id: 15, title: 'Draft quick architecture diagram for team meeting', status: 'todo', priority: 'P2', adoId: '', adoType: '', adoState: '', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'personal' },
      { id: 16, title: 'Test SSO flow on staging with VPN', status: 'done', priority: 'P1', adoId: '', adoType: '', adoState: '', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'personal' },
      { id: 17, title: 'Ping Sarah for code review on token PR', status: 'in_progress', priority: 'P1', adoId: '', adoType: '', adoState: '', assignedTo: 'Luis L.', syncStatus: 'synced', source: 'personal' },
    ],
    prs: [
      { id: 100, title: 'feat: add Azure AD provider', prNumber: 342, repo: 'platform-api', status: 'active', sourceBranch: 'refs/heads/feat/auth-refresh' },
      { id: 101, title: 'feat: token refresh middleware', prNumber: 343, repo: 'platform-api', status: 'draft', sourceBranch: 'refs/heads/feat/token-refresh' },
    ],
    notes: [
      { id: 1, text: 'Discussed with the team — we\'ll use MSAL.js v2 for the browser flow. Need to check if the redirect URI config supports Wails.', timestamp: '2026-04-07T15:30:00Z' },
      { id: 2, text: 'Token refresh interval should be 5 min before expiry. Check the Azure AD docs for silent refresh limits.', timestamp: '2026-04-06T10:00:00Z' },
    ],
    createdAt: '2026-03-20T10:00:00Z', updatedAt: '2026-04-07T15:30:00Z',
    dirtyFields: ['description', 'priority'], // these were edited locally
  },
  {
    id: 2,
    title: 'Fix dashboard loading performance',
    description: 'The dashboard takes 4+ seconds to load. Profile and optimize the initial data fetch.',
    status: 'todo', priority: 'P1',
    projectId: 2, projectName: 'Developer Experience',
    dueDate: '', adoId: '48350', adoType: 'Bug', adoState: 'New',
    subtasks: [
      { id: 20, title: 'Profile initial load with DevTools', status: 'todo', priority: 'P1', adoId: '48351', adoType: 'Task', adoState: 'New', assignedTo: '' },
    ],
    prs: [],
    notes: [
      { id: 3, text: 'Suspect it\'s the ADO sync call blocking the render. Try deferring it to after mount.', timestamp: '2026-04-07T10:00:00Z' },
    ],
    createdAt: '2026-04-05T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z',
    dirtyFields: ['status'], // status changed locally
  },
  {
    id: 3,
    title: 'Personal: review team retro notes',
    description: '',
    status: 'todo', priority: 'P3',
    projectId: null, projectName: '',
    dueDate: '2026-04-10', adoId: '', adoType: '', adoState: '',
    subtasks: [
      { id: 30, title: 'Read feedback summary', status: 'done', priority: 'P3', adoId: '', adoType: '', adoState: '', assignedTo: '' },
      { id: 31, title: 'Draft action items', status: 'todo', priority: 'P3', adoId: '', adoType: '', adoState: '', assignedTo: '' },
    ],
    prs: [],
    notes: [],
    createdAt: '2026-04-06T10:00:00Z', updatedAt: '2026-04-07T08:00:00Z',
    dirtyFields: [], // personal task, no ADO sync needed
  },
]

const selectedId = ref<number>(1)
const task = computed(() => mockTasks.find(t => t.id === selectedId.value)!)

const currentUser = 'Luis L.'
const subtaskFilter = ref<'all' | 'ado' | 'personal' | 'mine'>('all')
const showOnlyMine = computed(() => subtaskFilter.value === 'mine')

function isMine(st: MockSubtask): boolean {
  return !st.assignedTo || st.assignedTo === currentUser
}

function isOtherPerson(st: MockSubtask): boolean {
  return !!st.assignedTo && st.assignedTo !== currentUser
}

function isPersonal(st: MockSubtask): boolean {
  return st.source === 'personal'
}

const filteredSubtasks = computed(() => {
  const subs = task.value.subtasks
  switch (subtaskFilter.value) {
    case 'mine': return subs.filter(st => isMine(st))
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

function togglePersonalDone(st: MockSubtask) {
  st.status = st.status === 'done' ? 'todo' : 'done'
}

const statuses = ['todo', 'in_progress', 'in_review', 'done', 'blocked', 'cancelled']
const priorities = ['P0', 'P1', 'P2', 'P3']
const statusLabel: Record<string, string> = {
  in_progress: 'In Progress', in_review: 'In Review', todo: 'To Do',
  blocked: 'Blocked', done: 'Done', cancelled: 'Cancelled',
}
const priorityDot: Record<string, string> = {
  P0: 'bg-red-500', P1: 'bg-orange-500', P2: 'bg-amber-500', P3: 'bg-zinc-400',
}

function adoNumber(adoId: string) {
  return adoId ? `#${adoId}` : ''
}

function typeIcon(type: string) {
  return adoTypeIcon(type)
}

const subtaskProgress = computed(() => {
  const total = task.value.subtasks.length
  const done = task.value.subtasks.filter(s => s.status === 'done').length
  return { done, total, percent: total ? (done / total) * 100 : 0 }
})

function prIconColor(status: string) {
  return status === 'completed' ? 'text-violet-500' : status === 'draft' ? 'text-zinc-400' : 'text-emerald-500'
}

const descriptionOpen = ref(true)
const detailsOpen = ref(true)
const newNote = ref('')
const addingSubtask = ref(false)
const newSubtaskTitle = ref('')

function addSubtask() {
  addingSubtask.value = true
  newSubtaskTitle.value = ''
}

function confirmAddSubtask() {
  if (!newSubtaskTitle.value.trim()) return
  task.value.subtasks.push({
    id: Date.now(),
    title: newSubtaskTitle.value.trim(),
    status: 'todo', priority: 'P2',
    adoId: '', adoType: '', adoState: '',
    assignedTo: currentUser,
    syncStatus: 'synced',
    source: 'personal',
  })
  addingSubtask.value = false
  newSubtaskTitle.value = ''
}

// Editing state
const editingDescription = ref(false)
const editedDescription = ref('')
const editingNoteId = ref<number | null>(null)
const editedNoteText = ref('')

function startEditDescription() {
  editedDescription.value = task.value.description.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim()
  editingDescription.value = true
}

function saveDescription() {
  task.value.description = editedDescription.value
  editingDescription.value = false
}

function addNote() {
  if (!newNote.value.trim()) return
  task.value.notes.unshift({
    id: Date.now(),
    text: newNote.value.trim(),
    timestamp: new Date().toISOString(),
  })
  newNote.value = ''
}

function startEditNote(note: { id: number; text: string }) {
  editingNoteId.value = note.id
  editedNoteText.value = note.text
}

function saveNote(noteId: number) {
  const note = task.value.notes.find(n => n.id === noteId)
  if (note) note.text = editedNoteText.value.trim()
  editingNoteId.value = null
}

function deleteNote(noteId: number) {
  task.value.notes = task.value.notes.filter(n => n.id !== noteId)
}

// Nudge: suggest status update when PR is completed but task still in_progress
const nudge = computed(() => {
  const t = task.value
  const hasMergedPR = t.prs.some(pr => pr.status === 'completed')
  if (hasMergedPR && t.status === 'in_progress') {
    return 'A linked PR has been completed. Update task status?'
  }
  const hasFailedPipeline = t.prs.some(pr => pipelinesForPr(pr.sourceBranch).some(p => p.result === 'failed'))
  if (hasFailedPipeline && t.status !== 'blocked') {
    return 'A pipeline has failed. Mark task as blocked?'
  }
  return null
})

const mockProjects = [
  { id: 1, name: 'Platform Modernization' },
  { id: 2, name: 'Developer Experience' },
  { id: 3, name: 'Infrastructure & Reliability' },
]
</script>

<template>
  <div class="h-screen flex flex-col bg-background text-foreground">
    <!-- Top bar -->
    <div class="shrink-0 border-b border-border px-6 py-3 flex items-center gap-4">
      <h1 class="text-sm font-semibold">Playground: TaskDetail — ADO-style</h1>
      <div class="flex-1" />
      <!-- Task selector -->
      <div class="flex items-center gap-2">
        <span class="text-xs text-muted-foreground">Test task:</span>
        <button
          v-for="t in mockTasks" :key="t.id"
          class="px-2 py-0.5 text-xs rounded border transition-colors"
          :class="selectedId === t.id
            ? 'bg-primary text-primary-foreground border-primary'
            : 'border-border hover:bg-muted'"
          @click="selectedId = t.id"
        >
          {{ t.title.slice(0, 25) }}...
        </button>
      </div>
    </div>

    <!-- Side-by-side comparison -->
    <div class="flex-1 flex min-h-0">

      <!-- ═══════ LEFT: Current layout ═══════ -->
      <div class="w-1/2 border-r border-border flex flex-col min-h-0">
        <div class="shrink-0 px-4 py-1.5 bg-muted/50 border-b border-border">
          <span class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Current Layout</span>
        </div>
        <div class="flex-1 flex flex-col min-h-0">
          <!-- Header -->
          <div class="shrink-0 border-b border-border px-4 pt-3 pb-3">
            <div class="flex items-start gap-2">
              <Input
                :model-value="task.title"
                :class="cn(
                  'text-base font-semibold border-none shadow-none focus-visible:ring-0 bg-transparent px-0 h-auto flex-1',
                  task.status === 'done' && 'line-through text-muted-foreground'
                )"
                placeholder="Task title..."
                readonly
              />
              <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0 -mt-0.5">
                <X :size="16" />
              </Button>
            </div>

            <!-- Current badge row -->
            <div class="flex items-center gap-x-2 gap-y-1.5 mt-2 flex-wrap">
              <!-- Status pill -->
              <div class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 border border-border">
                <span :class="cn('size-2 rounded-full', statusBgColor(task.status))" />
                <span>{{ statusLabel[task.status] ?? task.status }}</span>
                <ChevronDown :size="10" class="text-muted-foreground" />
              </div>
              <!-- Priority pill -->
              <div class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 border border-border">
                <span :class="cn('size-2 rounded-full', priorityDot[task.priority] ?? 'bg-zinc-400')" />
                <span>{{ task.priority }}</span>
                <ChevronDown :size="10" class="text-muted-foreground" />
              </div>

              <span class="hidden sm:block w-px h-4 bg-border" />

              <!-- Project pill (current cramped dropdown) -->
              <div class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 border border-border">
                <Folder :size="11" class="text-muted-foreground" />
                <span class="truncate max-w-[120px]">{{ task.projectName || 'None' }}</span>
                <ChevronDown :size="10" class="text-muted-foreground" />
              </div>
              <!-- Due date pill -->
              <div class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 border border-border">
                <CalendarDays :size="11" class="text-muted-foreground" />
                <span :class="!task.dueDate && 'text-muted-foreground'">{{ task.dueDate || 'No date' }}</span>
              </div>
              <!-- ADO pill -->
              <span
                v-if="task.adoId"
                class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-500 border border-blue-500/30 rounded-full px-2 py-0.5"
              >
                <AzureDevOpsIcon :size="12" />
                {{ adoNumber(task.adoId) }}
                <ChevronDown :size="10" />
              </span>
              <span v-if="task.adoId" class="text-[9px] text-muted-foreground/60">local copy</span>
            </div>

            <!-- Subtask progress -->
            <div v-if="task.subtasks.length > 0" class="mt-2.5">
              <div class="h-1 w-full rounded-full bg-muted">
                <div class="h-1 rounded-full bg-blue-500 transition-all" :style="{ width: subtaskProgress.percent + '%' }" />
              </div>
            </div>
          </div>

          <!-- Content (current) -->
          <ScrollArea class="flex-1 min-h-0">
            <div class="flex flex-col px-4 pt-3 pb-3 space-y-4">
              <!-- Description (collapsible) -->
              <div>
                <button class="flex items-center gap-1.5 w-full" @click="descriptionOpen = !descriptionOpen">
                  <ChevronRight :size="14" class="text-muted-foreground transition-transform" :class="descriptionOpen && 'rotate-90'" />
                  <h3 class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
                </button>
                <div v-if="descriptionOpen && task.description" class="mt-2 text-[13px] text-muted-foreground border rounded p-2 bg-muted/30">
                  {{ task.description.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim() || 'No description' }}
                </div>
              </div>

              <!-- Subtasks -->
              <div>
                <div class="flex items-center gap-2 mb-2">
                  <h3 class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Subtasks</h3>
                  <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ subtaskProgress.done }}/{{ subtaskProgress.total }}</Badge>
                </div>
                <div class="flex flex-col gap-0.5">
                  <div v-for="st in task.subtasks" :key="st.id" class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50">
                    <span :class="cn('size-3.5 rounded-[3px] border-[1.5px] shrink-0 flex items-center justify-center',
                      st.status === 'done' ? 'bg-emerald-500 border-emerald-500' : 'border-muted-foreground/60')">
                      <svg v-if="st.status === 'done'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 12 12" class="size-2.5 text-white" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="2 6 5 9 10 3" />
                      </svg>
                    </span>
                    <span :class="cn('text-[13px]', st.status === 'done' && 'line-through text-muted-foreground')">{{ st.title }}</span>
                  </div>
                </div>
              </div>

              <!-- Details (collapsible) -->
              <div>
                <button class="flex items-center gap-1.5 w-full" @click="detailsOpen = !detailsOpen">
                  <ChevronRight :size="14" class="text-muted-foreground transition-transform" :class="detailsOpen && 'rotate-90'" />
                  <h3 class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Details</h3>
                </button>
                <div v-if="detailsOpen" class="mt-2 space-y-3">
                  <!-- PRs -->
                  <div>
                    <div class="flex items-center gap-2 mb-2">
                      <span class="text-xs font-medium text-muted-foreground">Pull Requests</span>
                      <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ task.prs.length }}</Badge>
                    </div>
                    <div class="flex flex-col gap-1">
                      <div v-for="pr in task.prs" :key="pr.id" class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50">
                        <GitPullRequest :size="14" :class="cn('shrink-0', prIconColor(pr.status))" />
                        <span class="text-[13px] truncate flex-1">{{ pr.title }}</span>
                        <span class="text-[11px] text-muted-foreground shrink-0">#{{ pr.prNumber }}</span>
                      </div>
                      <p v-if="task.prs.length === 0" class="text-xs text-muted-foreground italic py-1 px-1">No linked PRs</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </ScrollArea>

          <!-- Footer -->
          <div class="shrink-0 border-t border-border flex items-center justify-between px-4 py-1.5">
            <span class="text-[11px] text-muted-foreground">Created {{ new Date(task.createdAt).toLocaleDateString() }} · Updated {{ new Date(task.updatedAt).toLocaleDateString() }}</span>
            <Button variant="ghost" size="sm" class="h-6 text-[11px] text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1">
              <Trash2 :size="12" /> Delete
            </Button>
          </div>
        </div>
      </div>

      <!-- ═══════ RIGHT: ADO-style layout ═══════ -->
      <div class="w-1/2 flex flex-col min-h-0">
        <div class="shrink-0 px-4 py-1.5 bg-primary/5 border-b border-border">
          <span class="text-[10px] font-medium text-primary uppercase tracking-wider">New ADO-style Layout</span>
        </div>
        <div class="flex-1 flex flex-col min-h-0">

          <!-- Header: type icon + title + state + close (like ADO panel) -->
          <div class="shrink-0 border-b border-border px-4 py-3 space-y-2">
            <div class="flex items-center gap-2">
              <component
                v-if="task.adoType"
                :is="typeIcon(task.adoType)"
                :size="16"
                :class="adoTypeColor(task.adoType)"
              />
              <h2 class="text-sm font-semibold text-foreground leading-snug flex-1 line-clamp-2">{{ task.title }}</h2>
              <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0">
                <X :size="16" />
              </Button>
            </div>

            <!-- Inline metadata row (like ADO) — dirty fields get amber dot -->
            <div class="flex items-center gap-2 flex-wrap">
              <Badge variant="outline" :class="['text-[10px] h-5 px-1.5', statusBgColor(task.status) ? '' : '']">
                <span :class="cn('size-1.5 rounded-full mr-1', statusBgColor(task.status))" />
                {{ statusLabel[task.status] ?? task.status }}
                <span v-if="task.dirtyFields.includes('status')" class="size-1 rounded-full bg-amber-500 ml-1" title="Modified locally" />
              </Badge>
              <Badge variant="outline" :class="['text-[10px] h-5 px-1.5']">
                <span :class="cn('size-1.5 rounded-full mr-1', priorityDot[task.priority] ?? 'bg-zinc-400')" />
                {{ task.priority }}
                <span v-if="task.dirtyFields.includes('priority')" class="size-1 rounded-full bg-amber-500 ml-1" title="Modified locally" />
              </Badge>
              <span v-if="task.dueDate" class="text-[10px] text-muted-foreground flex items-center gap-1">
                <CalendarDays :size="10" />
                {{ task.dueDate }}
              </span>
            </div>

            <!-- Project context (below metadata, full width, not a cramped pill) -->
            <div v-if="task.projectName" class="flex items-center gap-2 px-2 py-1.5 rounded-md bg-muted/40 border border-border/50">
              <component
                v-if="task.adoType"
                :is="typeIcon(task.adoType)"
                :size="13"
                :class="adoTypeColor(task.adoType)"
                class="shrink-0"
              />
              <Folder v-else :size="13" class="text-muted-foreground shrink-0" />
              <span class="text-xs font-medium text-foreground truncate flex-1">{{ task.projectName }}</span>
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
              <Folder :size="13" class="text-muted-foreground/40 shrink-0" />
              <span class="text-[11px] text-muted-foreground/60 flex-1">No project</span>
              <Select model-value="none">
                <SelectTrigger size="sm" class="h-5 text-[10px] gap-0.5 w-auto px-1.5 border-none shadow-none bg-transparent hover:bg-muted">
                  <span class="text-muted-foreground">Assign</span>
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">None</SelectItem>
                  <SelectItem v-for="p in mockProjects" :key="p.id" :value="String(p.id)">{{ p.name }}</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- ADO integration row -->
            <div v-if="task.adoId" class="flex items-center gap-2 px-2 py-1 rounded-md bg-blue-500/5 border border-blue-500/15">
              <AzureDevOpsIcon :size="12" class="text-blue-500 shrink-0" />
              <span class="text-[10px] text-blue-500 tabular-nums font-medium">#{{ task.adoId }}</span>
              <span class="text-[9px] text-muted-foreground/50">·</span>
              <span v-if="task.dirtyFields.length === 0" class="text-[9px] text-emerald-600 flex items-center gap-0.5">
                <CheckCircle2 :size="9" /> Synced
              </span>
              <span v-else class="text-[9px] text-amber-600 flex items-center gap-0.5">
                <span class="size-1 rounded-full bg-amber-500" />
                {{ task.dirtyFields.length }} pending
              </span>
              <div class="flex-1" />
              <Button
                v-if="task.dirtyFields.length > 0"
                variant="outline"
                size="sm"
                class="h-5 text-[9px] gap-0.5 px-1.5 text-amber-600 border-amber-500/30 hover:bg-amber-500/10"
              >
                <Upload :size="10" /> Push {{ task.dirtyFields.length }}
              </Button>
              <Button v-else variant="ghost" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5 text-muted-foreground hover:text-foreground">
                <Upload :size="10" /> Push
              </Button>
              <Button variant="ghost" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5 text-blue-500 hover:text-blue-600">
                <ExternalLink :size="10" /> Open
              </Button>
            </div>
            <div v-else class="flex items-center gap-2 px-2 py-1 rounded-md border border-dashed border-border/50">
              <User :size="12" class="text-muted-foreground/40 shrink-0" />
              <span class="text-[10px] text-muted-foreground/60">Personal task — not linked to ADO</span>
            </div>

            <!-- Subtask progress bar -->
            <div v-if="task.subtasks.length > 0">
              <div class="h-1 w-full rounded-full bg-muted">
                <div class="h-1 rounded-full bg-blue-500 transition-all" :style="{ width: subtaskProgress.percent + '%' }" />
              </div>
            </div>
          </div>

          <!-- Scrollable content (subtasks + description + PRs + notes) -->
          <ScrollArea class="flex-1 min-h-0">
            <div class="flex flex-col">

          <!-- Subtasks (always shown) -->
          <div class="border-b border-border px-4 py-2">
            <div class="flex items-center justify-between mb-1.5">
              <div class="flex items-center gap-2">
                <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Subtasks</h3>
                <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ subtaskProgress.done }}/{{ subtaskProgress.total }}</Badge>
              </div>
              <div class="flex items-center gap-2">
                <!-- Filter toggle: All → Mine → ADO → Personal -->
                <button
                  class="text-[10px] flex items-center gap-1 px-1.5 py-0.5 rounded transition-colors"
                  :class="subtaskFilter !== 'all'
                    ? 'bg-primary/10 text-primary font-medium'
                    : 'text-muted-foreground hover:text-foreground'"
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
              <div
                v-for="st in filteredSubtasks" :key="st.id"
                class="rounded group"
                :class="st.syncStatus === 'not-pulled' ? 'opacity-50' : ''"
              >
                <!-- Subtask row -->
                <div
                  class="flex items-center gap-2 py-1.5 px-2 cursor-pointer transition-colors"
                  :class="st.syncStatus === 'not-pulled' ? 'hover:bg-muted/30' : 'hover:bg-muted/50'"
                >
                  <!-- Icon: ADO type icon for ADO items, clickable checkbox for personal -->
                  <template v-if="isPersonal(st)">
                    <button
                      class="size-3.5 rounded-[3px] border-[1.5px] shrink-0 flex items-center justify-center transition-colors"
                      :class="st.status === 'done' ? 'bg-emerald-500 border-emerald-500' : 'border-primary/60 hover:border-primary'"
                      @click.stop="togglePersonalDone(st)"
                      title="Toggle done"
                    >
                      <svg v-if="st.status === 'done'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 12 12" class="size-2.5 text-white" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="2 6 5 9 10 3" />
                      </svg>
                    </button>
                  </template>
                  <template v-else>
                    <component
                      v-if="st.adoType"
                      :is="typeIcon(st.adoType)"
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
                  <!-- Title -->
                  <span :class="cn('text-xs truncate flex-1',
                    st.status === 'done' && 'line-through text-muted-foreground',
                    st.syncStatus === 'not-pulled' && 'text-muted-foreground/60')">
                    {{ st.title }}
                  </span>
                  <!-- Personal badge -->
                  <span v-if="isPersonal(st)" class="text-[8px] px-1 py-0.5 rounded bg-primary/8 text-primary/70 shrink-0 border border-primary/10">personal</span>
                  <!-- State badge (ADO only) -->
                  <Badge v-if="st.adoState && !isPersonal(st)" variant="outline" :class="['text-[9px] h-3.5 px-1', adoStateClasses(st.adoState)]">
                    {{ st.adoState }}
                  </Badge>
                  <span v-else-if="!isPersonal(st)" :class="cn('size-1.5 rounded-full shrink-0', statusBgColor(st.status))" />
                  <!-- Sync indicator (ADO only) -->
                  <span v-if="st.syncStatus === 'pending' && !isPersonal(st)" class="size-1 rounded-full bg-amber-500 shrink-0" title="Pending push" />
                  <span v-if="st.syncStatus === 'not-pulled'" class="text-[8px] text-muted-foreground/40 shrink-0">not pulled</span>
                  <!-- Priority -->
                  <span :class="cn('size-1.5 rounded-full shrink-0', priorityDot[st.priority] ?? 'bg-zinc-400')" />
                  <!-- ADO ID (ADO only) -->
                  <span v-if="st.adoId && !isPersonal(st)" class="text-[9px] text-blue-500/60 tabular-nums shrink-0">#{{ st.adoId }}</span>
                  <!-- Assignee -->
                  <span
                    v-if="st.assignedTo"
                    class="text-[9px] shrink-0 rounded-full h-4 px-1.5 flex items-center gap-0.5"
                    :class="isOtherPerson(st)
                      ? 'bg-violet-500/10 text-violet-600 dark:text-violet-400 border border-violet-500/20'
                      : 'bg-muted text-muted-foreground'"
                    :title="st.assignedTo"
                  >
                    {{ st.assignedTo }}
                  </span>
                  <span v-else class="text-[9px] text-muted-foreground/30 shrink-0">Unassigned</span>
                  <!-- Delete personal subtask -->
                  <button
                    v-if="isPersonal(st)"
                    class="opacity-0 group-hover:opacity-100 text-muted-foreground/40 hover:text-red-500 transition-all shrink-0"
                    @click.stop="task.subtasks = task.subtasks.filter(s => s.id !== st.id)"
                    title="Remove"
                  >
                    <Trash2 :size="11" />
                  </button>
                </div>
                <!-- Push action for pending ADO items -->
                <div
                  v-if="st.syncStatus === 'pending' && !isPersonal(st)"
                  class="flex items-center gap-1.5 pl-7 pb-1 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <button class="text-[9px] text-amber-600 hover:text-amber-700 flex items-center gap-0.5">
                    <Upload :size="9" /> Push to ADO
                  </button>
                </div>
                <!-- Pull action for not-pulled items -->
                <div
                  v-if="st.syncStatus === 'not-pulled'"
                  class="flex items-center gap-1.5 pl-7 pb-1 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <button class="text-[9px] text-blue-500 hover:text-blue-600 flex items-center gap-0.5">
                    <Download :size="9" /> Pull from ADO
                  </button>
                </div>
              </div>
            </div>
            <!-- Add subtask inline input -->
            <div v-if="addingSubtask" class="flex items-center gap-2 mt-1.5 px-2">
              <button
                class="size-3.5 rounded-[3px] border-[1.5px] border-primary/40 shrink-0 flex items-center justify-center"
              />
              <Input
                v-model="newSubtaskTitle"
                placeholder="New personal subtask..."
                class="h-7 text-xs flex-1"
                @keydown.enter="confirmAddSubtask"
                @keydown.escape="addingSubtask = false"
              />
              <Button variant="outline" size="sm" class="h-7 px-2 text-[10px] shrink-0" @click="confirmAddSubtask" :disabled="!newSubtaskTitle.trim()">Add</Button>
              <Button variant="ghost" size="sm" class="h-7 px-1.5 text-[10px]" @click="addingSubtask = false">Cancel</Button>
            </div>
            <p v-if="filteredSubtasks.length === 0 && !addingSubtask" class="text-[11px] text-muted-foreground/40 italic py-2 text-center">
              {{ subtaskFilter === 'mine' ? 'No subtasks assigned to you' : subtaskFilter === 'personal' ? 'No personal subtasks' : subtaskFilter === 'ado' ? 'No ADO subtasks' : 'No subtasks yet — click Add to create one' }}
            </p>
          </div>

          <!-- Description (editable) -->
          <div class="border-b border-border px-4 py-3">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-1.5">
                <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
                <span v-if="task.dirtyFields.includes('description')" class="size-1 rounded-full bg-amber-500" title="Modified locally" />
              </div>
              <button
                v-if="!editingDescription"
                class="text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1"
                @click="startEditDescription"
              >
                <Pencil :size="10" /> Edit
              </button>
              <div v-else class="flex items-center gap-1">
                <button class="text-[10px] text-emerald-600 hover:text-emerald-700 flex items-center gap-0.5" @click="saveDescription">
                  <Save :size="10" /> Save
                </button>
                <button class="text-[10px] text-muted-foreground hover:text-foreground" @click="editingDescription = false">
                  Cancel
                </button>
              </div>
            </div>
            <div v-if="editingDescription">
              <textarea
                v-model="editedDescription"
                class="w-full min-h-[80px] text-xs bg-background border border-border rounded-md p-2 resize-y focus:outline-none focus:ring-1 focus:ring-primary"
                placeholder="Add a description..."
              />
            </div>
            <div v-else-if="task.description"
              class="text-xs text-foreground prose prose-sm max-w-none [&_*]:text-xs [&_*]:text-foreground cursor-pointer hover:bg-muted/30 rounded p-1 -m-1 transition-colors"
              v-html="task.description"
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
                  <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ task.prs.length }}</Badge>
                </div>
                <div class="flex flex-col gap-1">
                  <div v-for="pr in task.prs" :key="pr.id" class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50 cursor-pointer">
                    <GitPullRequest :size="14" :class="cn('shrink-0', prIconColor(pr.status))" />
                    <span class="text-[13px] truncate flex-1">{{ pr.title }}</span>
                    <span class="text-[11px] text-muted-foreground shrink-0">#{{ pr.prNumber }}</span>
                    <Badge variant="outline" :class="cn('text-[10px] capitalize shrink-0 h-4 px-1.5', prStatusClasses(pr.status))">
                      {{ pr.status }}
                    </Badge>
                  </div>
                  <p v-if="task.prs.length === 0" class="text-xs text-muted-foreground italic py-1 px-1">No linked PRs</p>
                </div>
              </div>

              <!-- Notes -->
              <div class="px-4 py-3">
                <div class="flex items-center gap-2 mb-2">
                  <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Notes</h3>
                  <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ task.notes.length }}</Badge>
                  <!-- ADO Discussion link inline with title -->
                  <button
                    v-if="task.adoId"
                    class="inline-flex items-center gap-1 text-[9px] text-blue-500 hover:text-blue-600 transition-colors ml-1"
                  >
                    <AzureDevOpsIcon :size="9" />
                    ADO Discussion
                    <ExternalLink :size="8" />
                  </button>
                </div>
                <!-- Note input -->
                <div class="mb-2 flex gap-1.5">
                  <Input
                    placeholder="Add a note..."
                    class="h-7 text-xs flex-1"
                    :model-value="newNote"
                    @update:model-value="newNote = String($event)"
                    @keydown.enter="addNote"
                  />
                  <Button variant="outline" size="sm" class="h-7 px-2 text-[10px] shrink-0" @click="addNote" :disabled="!newNote.trim()">
                    <Plus :size="11" class="mr-0.5" /> Add
                  </Button>
                </div>
                <!-- Note list -->
                <div class="flex flex-col gap-2">
                  <div
                    v-for="note in task.notes" :key="note.id"
                    class="group relative px-3 py-2 rounded-md bg-muted/40 border border-border/50"
                  >
                    <!-- Edit mode -->
                    <div v-if="editingNoteId === note.id">
                      <textarea
                        v-model="editedNoteText"
                        class="w-full min-h-[48px] text-xs bg-background border border-border rounded p-1.5 resize-y focus:outline-none focus:ring-1 focus:ring-primary"
                      />
                      <div class="flex items-center gap-1.5 mt-1">
                        <button class="text-[9px] text-emerald-600 hover:text-emerald-700 flex items-center gap-0.5" @click="saveNote(note.id)">
                          <Save :size="9" /> Save
                        </button>
                        <button class="text-[9px] text-muted-foreground hover:text-foreground" @click="editingNoteId = null">Cancel</button>
                      </div>
                    </div>
                    <!-- View mode -->
                    <div v-else>
                      <p class="text-xs text-foreground leading-relaxed pr-10">{{ note.text }}</p>
                      <span class="text-[9px] text-muted-foreground/60 mt-1 block">
                        {{ new Date(note.timestamp).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) }}
                        · {{ new Date(note.timestamp).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' }) }}
                      </span>
                      <!-- Edit/delete on hover -->
                      <div class="absolute top-1.5 right-1.5 flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button class="p-1 rounded hover:bg-muted text-muted-foreground hover:text-foreground" @click="startEditNote(note)">
                          <Pencil :size="10" />
                        </button>
                        <button class="p-1 rounded hover:bg-red-500/10 text-muted-foreground hover:text-red-500" @click="deleteNote(note.id)">
                          <Trash2 :size="10" />
                        </button>
                      </div>
                    </div>
                  </div>
                  <p v-if="task.notes.length === 0" class="text-xs text-muted-foreground/40 text-center py-3 italic">
                    No notes yet. Add one above.
                  </p>
                </div>
              </div>

            </div>
          </ScrollArea>

          <!-- Footer -->
          <div class="shrink-0 border-t border-border flex items-center justify-between px-4 py-1.5">
            <span class="text-[11px] text-muted-foreground">
              Created {{ new Date(task.createdAt).toLocaleDateString() }} · Updated {{ new Date(task.updatedAt).toLocaleDateString() }}
            </span>
            <Button variant="ghost" size="sm" class="h-6 text-[11px] text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1">
              <Trash2 :size="12" /> Delete
            </Button>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>
