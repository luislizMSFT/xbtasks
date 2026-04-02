<script setup lang="ts">
import { ref, computed } from 'vue'
import { cn } from '@/lib/utils'
import { useTheme } from '@/composables/useTheme'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  LayoutDashboard,
  ListTodo,
  FolderKanban,
  Settings,
  Sun,
  Moon,
  Search,
  Plus,
  Clock,
  Bug,
  CheckSquare,
  BookOpen,
  ChevronDown,
  ChevronRight,
  Pencil,
  CheckCircle2,
  Circle,
  GitPullRequest,
  Send,
  Trash2,
  Folder,
  CalendarDays,
  Landmark,
  Tag,
  MessageSquare,
  User,
  ExternalLink,
  Activity,
  GitBranch,
  ArrowRight,
} from 'lucide-vue-next'

const { mode, toggle } = useTheme()

// ── Types ──────────────────────────────────────────────────────────────
interface Subtask {
  title: string
  done: boolean
}

interface PullRequest {
  number: number
  title: string
  repo: string
  source: string
  target: string
  status: 'Active' | 'Draft' | 'Merged'
  reviewers?: { name: string; initials: string; approved: boolean }[]
}

interface Comment {
  author: string
  initials: string
  text: string
  time: string
}

interface MockTask {
  id: string
  title: string
  status: 'in_progress' | 'in_review' | 'todo' | 'blocked' | 'done'
  priority: 'P0' | 'P1' | 'P2' | 'P3'
  adoType?: 'Bug' | 'Task' | 'UserStory'
  adoNumber?: number
  timeAgo: string
  group: 'in_progress' | 'todo' | 'blocked' | 'done'
  description?: string
  subtasks?: Subtask[]
  prs?: PullRequest[]
  comments?: Comment[]
  adoComments?: Comment[]
  project?: string
  dueDate?: string
  tags?: string[]
  createdAt?: string
  updatedAt?: string
}

type StyleVariant = 'compact' | 'spacious' | 'dense'

// ── Mock data ──────────────────────────────────────────────────────────
const tasks: MockTask[] = [
  {
    id: 't1',
    title: 'Fix auth redirect loop',
    status: 'in_progress',
    priority: 'P0',
    adoType: 'Bug',
    adoNumber: 48291,
    timeAgo: '30m ago',
    group: 'in_progress',
    description:
      'The OAuth callback is redirecting in a loop when the token expires mid-session. Need to check the refresh token logic and handle edge cases around concurrent requests and stale sessions.',
    subtasks: [
      { title: 'Identify refresh token edge case', done: true },
      { title: 'Add token expiry middleware', done: true },
      { title: 'Update error handling in auth callback', done: false },
      { title: 'Add E2E test for token refresh', done: false },
    ],
    prs: [
      {
        number: 234,
        title: 'Fix auth redirect loop',
        repo: 'xb-services',
        source: 'fix/auth-redirect',
        target: 'main',
        status: 'Active',
        reviewers: [
          { name: 'Alex', initials: 'AK', approved: true },
          { name: 'Sam', initials: 'SL', approved: false },
        ],
      },
      {
        number: 240,
        title: 'Add token expiry middleware',
        repo: 'xb-services',
        source: 'feat/token-middleware',
        target: 'main',
        status: 'Draft',
      },
    ],
    comments: [
      { author: 'You', initials: 'LL', text: 'Root cause is in the token refresh — concurrent requests cause a race condition.', time: '25m ago' },
      { author: 'You', initials: 'LL', text: 'Middleware approach looks cleaner than patching the callback directly.', time: '1h ago' },
    ],
    adoComments: [
      { author: 'Alex K.', initials: 'AK', text: 'Confirmed the loop on staging. Happens when refresh token is within 30s of expiry.', time: '2h ago' },
      { author: 'Jordan M.', initials: 'JM', text: 'Priority bump — this is blocking the demo prep.', time: '4h ago' },
    ],
    project: 'Team ADO Tool',
    dueDate: 'Apr 10, 2026',
    tags: ['bug', 'urgent', 'auth'],
    createdAt: 'Mar 31',
    updatedAt: '2h ago',
  },
  {
    id: 't2',
    title: 'Build sidebar navigation',
    status: 'in_progress',
    priority: 'P1',
    timeAgo: '1h ago',
    group: 'in_progress',
    description: 'Implement the collapsible sidebar with icon rail and full-width modes. Must support keyboard navigation and persist collapsed state.',
    subtasks: [
      { title: 'Icon rail layout', done: true },
      { title: 'Expanded panel with labels', done: false },
      { title: 'Collapse animation', done: false },
    ],
    comments: [
      { author: 'You', initials: 'LL', text: 'Starting with the icon rail. Using 48px width to match macOS sidebar patterns.', time: '45m ago' },
    ],
    adoComments: [],
    project: 'Team ADO Tool',
    tags: ['frontend', 'ui'],
    createdAt: 'Apr 1',
    updatedAt: '1h ago',
  },
  {
    id: 't3',
    title: 'Review PR #234 — schema changes',
    status: 'todo',
    priority: 'P1',
    adoType: 'Task',
    adoNumber: 48350,
    timeAgo: '2h ago',
    group: 'todo',
    description: 'Review the database schema changes in PR #234. Focus on migration safety and backward compatibility.',
    subtasks: [
      { title: 'Review migration scripts', done: false },
      { title: 'Check rollback path', done: false },
    ],
    comments: [],
    adoComments: [
      { author: 'DB Team', initials: 'DB', text: 'Schema changes are backward compatible. No downtime migration.', time: '1h ago' },
    ],
    project: 'Team ADO Tool',
    tags: ['review'],
    createdAt: 'Apr 2',
    updatedAt: '2h ago',
  },
  {
    id: 't4',
    title: 'Add rate limiting to gateway',
    status: 'todo',
    priority: 'P1',
    adoType: 'Task',
    adoNumber: 48400,
    timeAgo: '3h ago',
    group: 'todo',
    description: 'Implement rate limiting middleware on the API gateway. Use sliding window algorithm with Redis backing store.',
    subtasks: [],
    comments: [],
    adoComments: [],
    project: 'Platform',
    tags: ['backend', 'security'],
    createdAt: 'Apr 1',
    updatedAt: '3h ago',
  },
  {
    id: 't5',
    title: 'Migrate auth to MSAL v3',
    status: 'blocked',
    priority: 'P0',
    adoType: 'Task',
    adoNumber: 48100,
    timeAgo: '1d ago',
    group: 'blocked',
    description: 'Upgrade from ADAL to MSAL v3 for authentication. Blocked on IT team providing the new app registration.',
    subtasks: [
      { title: 'Create MSAL config', done: true },
      { title: 'Swap auth provider', done: false },
      { title: 'Update token cache', done: false },
    ],
    comments: [
      { author: 'You', initials: 'LL', text: 'Waiting on IT for the new app registration. Pinged them again.', time: '6h ago' },
    ],
    adoComments: [
      { author: 'IT Ops', initials: 'IT', text: 'App registration is in the approval queue. ETA 2 business days.', time: '1d ago' },
    ],
    project: 'Platform',
    tags: ['auth', 'migration'],
    createdAt: 'Mar 28',
    updatedAt: '1d ago',
  },
  {
    id: 't6',
    title: 'Write E2E tests for task CRUD',
    status: 'done',
    priority: 'P2',
    timeAgo: '2d ago',
    group: 'done',
    description: 'Full E2E test suite covering create, read, update, delete flows for tasks. Using Playwright.',
    subtasks: [
      { title: 'Create task test', done: true },
      { title: 'Edit task test', done: true },
      { title: 'Delete task test', done: true },
    ],
    comments: [
      { author: 'You', initials: 'LL', text: 'All tests passing. Coverage at 94% for CRUD operations.', time: '2d ago' },
    ],
    adoComments: [],
    project: 'Team ADO Tool',
    tags: ['testing'],
    createdAt: 'Mar 25',
    updatedAt: '2d ago',
  },
]

const groups = [
  { key: 'in_progress' as const, label: 'In Progress' },
  { key: 'todo' as const, label: 'To Do' },
  { key: 'blocked' as const, label: 'Blocked' },
  { key: 'done' as const, label: 'Done' },
]

const navItems = [
  { key: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { key: 'tasks', label: 'Tasks', icon: ListTodo },
  { key: 'projects', label: 'Projects', icon: FolderKanban },
  { key: 'settings', label: 'Settings', icon: Settings },
]

// ── Reactive state ─────────────────────────────────────────────────────
const selectedTaskId = ref('t1')
const variant = ref<StyleVariant>('compact')
const expandedGroups = ref<Record<string, boolean>>({
  in_progress: true,
  todo: true,
  blocked: true,
  done: true,
})
const descriptionOpen = ref(true)
const prsOpen = ref(true)
const noteText = ref('')
const commentTab = ref('notes')

// ── Computed ───────────────────────────────────────────────────────────
const selectedTask = computed(() => tasks.find((t) => t.id === selectedTaskId.value) ?? tasks[0])

const subtasksDone = computed(() => selectedTask.value.subtasks?.filter((s) => s.done).length ?? 0)
const subtasksTotal = computed(() => selectedTask.value.subtasks?.length ?? 0)
const subtaskPct = computed(() => (subtasksTotal.value > 0 ? (subtasksDone.value / subtasksTotal.value) * 100 : 0))

const statCounts = computed(() => ({
  in_progress: tasks.filter((t) => t.status === 'in_progress').length,
  in_review: tasks.filter((t) => t.status === 'in_review').length,
  blocked: tasks.filter((t) => t.status === 'blocked').length,
  done: tasks.filter((t) => t.status === 'done').length,
}))

function groupTasks(key: string) {
  return tasks.filter((t) => t.group === key)
}

function toggleGroup(key: string) {
  expandedGroups.value[key] = !expandedGroups.value[key]
}

// ── Style helpers ──────────────────────────────────────────────────────
const statusColor: Record<string, string> = {
  in_progress: 'bg-blue-500',
  in_review: 'bg-violet-500',
  todo: 'bg-zinc-400',
  blocked: 'bg-red-500',
  done: 'bg-emerald-500',
}

const statusLabel: Record<string, string> = {
  in_progress: 'In Progress',
  in_review: 'In Review',
  todo: 'To Do',
  blocked: 'Blocked',
  done: 'Done',
}

const priorityBorder: Record<string, string> = {
  P0: 'border-l-red-500',
  P1: 'border-l-orange-500',
  P2: 'border-l-amber-500',
  P3: 'border-l-zinc-400',
}

const priorityDot: Record<string, string> = {
  P0: 'bg-red-500',
  P1: 'bg-orange-500',
  P2: 'bg-amber-500',
  P3: 'bg-zinc-400',
}

const priorityLabel: Record<string, string> = {
  P0: 'P0 — Critical',
  P1: 'P1 — High',
  P2: 'P2 — Medium',
  P3: 'P3 — Low',
}

function adoIcon(type?: string) {
  if (type === 'Bug') return Bug
  if (type === 'Task') return CheckSquare
  if (type === 'UserStory') return BookOpen
  return null
}

function adoColor(type?: string) {
  if (type === 'Bug') return 'text-red-500'
  if (type === 'Task') return 'text-blue-500'
  if (type === 'UserStory') return 'text-green-500'
  return 'text-muted-foreground'
}

// variant-aware class helpers
const v = computed(() => {
  const s = variant.value
  return {
    sectionGap: s === 'compact' ? 'space-y-2' : s === 'spacious' ? 'space-y-4' : 'space-y-1.5',
    sectionPx: s === 'compact' ? 'px-4' : s === 'spacious' ? 'px-5' : 'px-3',
    sectionPy: s === 'compact' ? 'py-2' : s === 'spacious' ? 'py-3' : 'py-1',
    titleSize: s === 'compact' ? 'text-base' : s === 'spacious' ? 'text-lg' : 'text-sm',
    bodySize: s === 'compact' ? 'text-[13px]' : s === 'spacious' ? 'text-sm' : 'text-xs',
    labelSize: s === 'compact' ? 'text-xs' : s === 'spacious' ? 'text-sm' : 'text-[11px]',
    sectionWrap: s === 'spacious' ? 'border border-border rounded-lg p-4' : '',
    monoClass: s === 'dense' ? 'font-mono' : '',
    rowPy: s === 'compact' ? 'py-1.5' : s === 'spacious' ? 'py-2' : 'py-1',
  }
})
</script>

<template>
  <div class="flex flex-col gap-3 p-4" style="height: calc(100vh - 2rem)">
    <!-- Variant selector tabs -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h2 class="text-sm font-semibold text-foreground">Detail Panel Playground</h2>
        <Separator orientation="vertical" class="h-4" />
        <span class="text-xs text-muted-foreground">Style variant:</span>
      </div>
      <Tabs v-model="variant" class="w-auto">
        <TabsList class="h-8">
          <TabsTrigger value="compact" class="text-xs px-3 h-6">Compact</TabsTrigger>
          <TabsTrigger value="spacious" class="text-xs px-3 h-6">Spacious</TabsTrigger>
          <TabsTrigger value="dense" class="text-xs px-3 h-6">Dense</TabsTrigger>
        </TabsList>
      </Tabs>
    </div>

    <!-- App Shell Mock -->
    <div class="flex flex-1 min-h-0 rounded-xl border border-border bg-background overflow-hidden">
      <!-- Sidebar icon rail -->
      <div class="w-12 shrink-0 flex flex-col items-center bg-muted/40 border-r border-border pt-[48px]">
        <!-- Traffic lights -->
        <div class="absolute top-3 left-3 flex items-center gap-1.5 z-10">
          <span class="size-3 rounded-full bg-[#FF5F57]" />
          <span class="size-3 rounded-full bg-[#FEBC2E]" />
          <span class="size-3 rounded-full bg-[#28C840]" />
        </div>

        <TooltipProvider :delay-duration="300">
          <div class="flex flex-col items-center gap-1 mt-2">
            <Tooltip v-for="item in navItems" :key="item.key">
              <TooltipTrigger as-child>
                <Button
                  variant="ghost"
                  size="icon"
                  :class="cn('size-9 rounded-lg', item.key === 'tasks' && 'bg-accent text-accent-foreground')"
                >
                  <component :is="item.icon" :size="18" :stroke-width="1.75" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="right" class="text-xs">{{ item.label }}</TooltipContent>
            </Tooltip>
          </div>
        </TooltipProvider>

        <div class="mt-auto mb-3">
          <Button variant="ghost" size="icon" class="size-9 rounded-lg" @click="toggle">
            <Sun v-if="mode === 'dark'" :size="16" :stroke-width="1.75" />
            <Moon v-else :size="16" :stroke-width="1.75" />
          </Button>
        </div>
      </div>

      <!-- Main content area -->
      <div class="flex flex-col flex-1 min-w-0">
        <!-- Top bar -->
        <div class="h-10 shrink-0 flex items-center justify-between px-4 border-b border-border bg-background">
          <span class="text-sm font-medium text-foreground">Tasks</span>

          <!-- Stats pills -->
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
              <span class="size-2 rounded-full bg-blue-500" />
              {{ statCounts.in_progress }}
            </div>
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
              <span class="size-2 rounded-full bg-violet-500" />
              {{ statCounts.in_review }}
            </div>
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
              <span class="size-2 rounded-full bg-red-500" />
              {{ statCounts.blocked }}
            </div>
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
              <span class="size-2 rounded-full bg-emerald-500" />
              {{ statCounts.done }}
            </div>
          </div>

          <!-- Right actions -->
          <div class="flex items-center gap-1">
            <Button variant="ghost" size="icon" class="size-8">
              <Search :size="15" :stroke-width="1.75" />
            </Button>
            <Button variant="ghost" size="sm" class="h-7 gap-1 text-xs">
              <Plus :size="14" :stroke-width="2" />
              New
            </Button>
            <Button variant="ghost" size="icon" class="size-8">
              <Activity :size="15" :stroke-width="1.75" />
            </Button>
          </div>
        </div>

        <!-- Content split -->
        <div class="flex flex-1 min-h-0">
          <!-- Left: Task list (~55%) -->
          <div class="w-[55%] border-r border-border flex flex-col min-h-0">
            <ScrollArea class="flex-1">
              <div class="py-1">
                <template v-for="group in groups" :key="group.key">
                  <div v-if="groupTasks(group.key).length > 0">
                    <!-- Group header -->
                    <button
                      class="flex items-center gap-2 w-full px-3 py-1.5 text-xs font-semibold text-muted-foreground hover:bg-accent/50 transition-colors"
                      @click="toggleGroup(group.key)"
                    >
                      <ChevronDown v-if="expandedGroups[group.key]" :size="14" />
                      <ChevronRight v-else :size="14" />
                      <span :class="cn('size-2 rounded-full', statusColor[group.key])" />
                      {{ group.label }}
                      <Badge variant="secondary" class="ml-1 h-4 text-[10px] px-1.5">
                        {{ groupTasks(group.key).length }}
                      </Badge>
                    </button>

                    <!-- Task rows -->
                    <div v-if="expandedGroups[group.key]">
                      <button
                        v-for="task in groupTasks(group.key)"
                        :key="task.id"
                        :class="cn(
                          'flex items-center gap-2.5 w-full px-3 py-2 text-left border-l-[3px] transition-colors',
                          priorityBorder[task.priority],
                          selectedTaskId === task.id ? 'bg-accent' : 'hover:bg-accent/50'
                        )"
                        @click="selectedTaskId = task.id"
                      >
                        <!-- Status dot -->
                        <span :class="cn('size-1.5 rounded-full shrink-0', statusColor[task.status])" />

                        <!-- Title -->
                        <span
                          :class="cn(
                            'text-sm truncate flex-1',
                            task.status === 'done' && 'line-through text-muted-foreground'
                          )"
                        >
                          {{ task.title }}
                        </span>

                        <!-- ADO badge -->
                        <span
                          v-if="task.adoType && task.adoNumber"
                          class="flex items-center gap-1 text-[11px] text-blue-500 shrink-0"
                        >
                          <component :is="adoIcon(task.adoType)" :size="12" :class="adoColor(task.adoType)" />
                          #{{ task.adoNumber }}
                        </span>

                        <!-- Time -->
                        <span class="text-[11px] text-muted-foreground shrink-0">{{ task.timeAgo }}</span>
                      </button>
                    </div>
                  </div>
                </template>
              </div>
            </ScrollArea>
          </div>

          <!-- Right: Detail panel (~45%) -->
          <div class="w-[45%] flex flex-col min-h-0 min-w-[340px]">
            <ScrollArea class="flex-1">
              <div :class="cn('flex flex-col', v.sectionGap)">

                <!-- ─── Header ─────────────────────────────────── -->
                <div :class="cn('border-b border-border', v.sectionPx, v.sectionPy, 'pt-3 pb-3')">
                  <div class="group flex items-start gap-2">
                    <h3 :class="cn('font-semibold leading-snug flex-1', v.titleSize, v.monoClass, selectedTask.status === 'done' && 'line-through text-muted-foreground')">
                      {{ selectedTask.title }}
                    </h3>
                    <Pencil :size="14" class="shrink-0 mt-1 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
                  </div>
                  <div class="flex items-center gap-2 mt-2 flex-wrap">
                    <!-- Status badge -->
                    <span class="inline-flex items-center gap-1.5 text-xs font-medium border border-border rounded-full px-2 py-0.5">
                      <span :class="cn('size-2 rounded-full', statusColor[selectedTask.status])" />
                      {{ statusLabel[selectedTask.status] }}
                    </span>
                    <!-- Priority badge -->
                    <span class="inline-flex items-center gap-1.5 text-xs font-medium border border-border rounded-full px-2 py-0.5">
                      <span :class="cn('size-2 rounded-full', priorityDot[selectedTask.priority])" />
                      {{ selectedTask.priority }}
                    </span>
                    <!-- ADO link -->
                    <span
                      v-if="selectedTask.adoType && selectedTask.adoNumber"
                      class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-500 border border-blue-500/30 rounded-full px-2 py-0.5 cursor-pointer hover:bg-blue-500/10 transition-colors"
                    >
                      <component :is="adoIcon(selectedTask.adoType)" :size="12" />
                      #{{ selectedTask.adoNumber }}
                      <ExternalLink :size="10" />
                    </span>
                  </div>
                </div>

                <!-- ─── Subtasks ───────────────────────────────── -->
                <div v-if="selectedTask.subtasks && selectedTask.subtasks.length > 0" :class="cn(v.sectionPx, v.sectionWrap && 'mx-4')">
                  <div :class="v.sectionWrap">
                    <div class="flex items-center justify-between mb-2">
                      <div class="flex items-center gap-2">
                        <span :class="cn('font-semibold', v.labelSize)">Subtasks</span>
                        <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ subtasksDone }}/{{ subtasksTotal }}</Badge>
                      </div>
                    </div>
                    <!-- Progress bar -->
                    <div class="h-1 w-full rounded-full bg-muted mb-2">
                      <div class="h-1 rounded-full bg-blue-500 transition-all" :style="{ width: subtaskPct + '%' }" />
                    </div>
                    <!-- Subtask rows -->
                    <div :class="cn('flex flex-col', variant === 'dense' ? 'gap-0.5' : 'gap-1')">
                      <div
                        v-for="(sub, i) in selectedTask.subtasks"
                        :key="i"
                        :class="cn('flex items-center gap-2', v.rowPy)"
                      >
                        <CheckCircle2 v-if="sub.done" :size="15" class="text-emerald-500 shrink-0" />
                        <Circle v-else :size="15" class="text-muted-foreground shrink-0" />
                        <span :class="cn(v.bodySize, v.monoClass, sub.done && 'line-through text-muted-foreground')">
                          {{ sub.title }}
                        </span>
                      </div>
                    </div>
                    <button class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground mt-1.5 transition-colors">
                      <Plus :size="12" />
                      Add subtask
                    </button>
                  </div>
                </div>

                <Separator v-if="variant !== 'spacious'" />

                <!-- ─── Description ────────────────────────────── -->
                <div v-if="selectedTask.description" :class="cn(v.sectionPx, v.sectionWrap && 'mx-4')">
                  <div :class="v.sectionWrap">
                    <button class="flex items-center gap-1.5 w-full" @click="descriptionOpen = !descriptionOpen">
                      <ChevronDown v-if="descriptionOpen" :size="14" class="text-muted-foreground" />
                      <ChevronRight v-else :size="14" class="text-muted-foreground" />
                      <span :class="cn('font-semibold', v.labelSize)">Description</span>
                    </button>
                    <p
                      v-if="descriptionOpen"
                      :class="cn('mt-1.5 text-muted-foreground leading-relaxed', v.bodySize, v.monoClass)"
                    >
                      {{ selectedTask.description }}
                    </p>
                  </div>
                </div>

                <Separator v-if="variant !== 'spacious'" />

                <!-- ─── Pull Requests ──────────────────────────── -->
                <div v-if="selectedTask.prs && selectedTask.prs.length > 0" :class="cn(v.sectionPx, v.sectionWrap && 'mx-4')">
                  <div :class="v.sectionWrap">
                    <button class="flex items-center gap-1.5 w-full" @click="prsOpen = !prsOpen">
                      <ChevronDown v-if="prsOpen" :size="14" class="text-muted-foreground" />
                      <ChevronRight v-else :size="14" class="text-muted-foreground" />
                      <span :class="cn('font-semibold', v.labelSize)">Pull Requests</span>
                      <Badge variant="secondary" class="h-4 text-[10px] px-1.5 ml-1">{{ selectedTask.prs.length }}</Badge>
                    </button>
                    <div v-if="prsOpen" :class="cn('flex flex-col mt-2', variant === 'dense' ? 'gap-1' : 'gap-2')">
                      <div
                        v-for="pr in selectedTask.prs"
                        :key="pr.number"
                        :class="cn('flex flex-col gap-1 rounded-md border border-border p-2', v.bodySize)"
                      >
                        <div class="flex items-center gap-2">
                          <GitPullRequest :size="14" class="text-muted-foreground shrink-0" />
                          <span class="font-medium truncate">{{ pr.title }}</span>
                          <Badge
                            :variant="pr.status === 'Draft' ? 'outline' : 'secondary'"
                            class="h-4 text-[10px] px-1.5 ml-auto shrink-0"
                          >
                            {{ pr.status }}
                          </Badge>
                        </div>
                        <div class="flex items-center gap-1.5 text-muted-foreground text-[11px] pl-5">
                          <Badge variant="outline" class="h-4 text-[10px] px-1.5">{{ pr.repo }}</Badge>
                          <GitBranch :size="11" />
                          <span :class="v.monoClass">{{ pr.source }}</span>
                          <ArrowRight :size="10" />
                          <span :class="v.monoClass">{{ pr.target }}</span>
                        </div>
                        <!-- Reviewers -->
                        <div v-if="pr.reviewers && pr.reviewers.length > 0" class="flex items-center gap-1.5 pl-5 mt-0.5">
                          <div
                            v-for="rev in pr.reviewers"
                            :key="rev.name"
                            :class="cn(
                              'size-5 rounded-full flex items-center justify-center text-[9px] font-semibold border',
                              rev.approved
                                ? 'bg-emerald-500/15 text-emerald-600 border-emerald-500/30'
                                : 'bg-muted text-muted-foreground border-border'
                            )"
                          >
                            {{ rev.initials }}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <Separator v-if="variant !== 'spacious'" />

                <!-- ─── Comments ───────────────────────────────── -->
                <div :class="cn(v.sectionPx, v.sectionWrap && 'mx-4')">
                  <div :class="v.sectionWrap">
                    <Tabs v-model="commentTab" class="w-full">
                      <TabsList class="h-7 w-full">
                        <TabsTrigger value="notes" class="text-[11px] h-5 flex-1 gap-1">
                          <MessageSquare :size="11" />
                          Notes
                          <Badge v-if="selectedTask.comments && selectedTask.comments.length > 0" variant="secondary" class="h-3.5 text-[9px] px-1">
                            {{ selectedTask.comments.length }}
                          </Badge>
                        </TabsTrigger>
                        <TabsTrigger value="ado" class="text-[11px] h-5 flex-1 gap-1">
                          <Landmark :size="11" />
                          ADO Comments
                          <Badge v-if="selectedTask.adoComments && selectedTask.adoComments.length > 0" variant="secondary" class="h-3.5 text-[9px] px-1">
                            {{ selectedTask.adoComments.length }}
                          </Badge>
                        </TabsTrigger>
                      </TabsList>

                      <TabsContent value="notes" class="mt-2">
                        <div :class="cn('flex flex-col', variant === 'dense' ? 'gap-1.5' : 'gap-2.5')">
                          <div
                            v-for="(c, i) in selectedTask.comments ?? []"
                            :key="'n' + i"
                            class="flex gap-2"
                          >
                            <div class="size-6 rounded-full bg-muted flex items-center justify-center shrink-0">
                              <span class="text-[9px] font-semibold text-muted-foreground">{{ c.initials }}</span>
                            </div>
                            <div class="flex-1 min-w-0">
                              <div class="flex items-baseline gap-2">
                                <span :class="cn('font-medium', v.bodySize)">{{ c.author }}</span>
                                <span :class="cn('text-muted-foreground', v.monoClass, 'text-[10px]')">{{ c.time }}</span>
                              </div>
                              <p :class="cn('text-muted-foreground mt-0.5', v.bodySize)">{{ c.text }}</p>
                            </div>
                          </div>
                          <p
                            v-if="!selectedTask.comments || selectedTask.comments.length === 0"
                            class="text-xs text-muted-foreground italic"
                          >
                            No notes yet.
                          </p>
                        </div>
                        <!-- Input -->
                        <div class="flex items-center gap-2 mt-3">
                          <Input
                            v-model="noteText"
                            placeholder="Add a note..."
                            class="h-7 text-xs flex-1"
                          />
                          <Button variant="ghost" size="icon" class="size-7 shrink-0">
                            <Send :size="13" />
                          </Button>
                        </div>
                      </TabsContent>

                      <TabsContent value="ado" class="mt-2">
                        <div :class="cn('flex flex-col', variant === 'dense' ? 'gap-1.5' : 'gap-2.5')">
                          <div
                            v-for="(c, i) in selectedTask.adoComments ?? []"
                            :key="'a' + i"
                            class="flex gap-2"
                          >
                            <div class="size-6 rounded-full bg-muted flex items-center justify-center shrink-0">
                              <span class="text-[9px] font-semibold text-muted-foreground">{{ c.initials }}</span>
                            </div>
                            <div class="flex-1 min-w-0">
                              <div class="flex items-baseline gap-2">
                                <span :class="cn('font-medium', v.bodySize)">{{ c.author }}</span>
                                <span :class="cn('text-muted-foreground', v.monoClass, 'text-[10px]')">{{ c.time }}</span>
                              </div>
                              <p :class="cn('text-muted-foreground mt-0.5', v.bodySize)">{{ c.text }}</p>
                            </div>
                          </div>
                          <p
                            v-if="!selectedTask.adoComments || selectedTask.adoComments.length === 0"
                            class="text-xs text-muted-foreground italic"
                          >
                            No ADO comments synced.
                          </p>
                        </div>
                      </TabsContent>
                    </Tabs>
                  </div>
                </div>

                <Separator v-if="variant !== 'spacious'" />

                <!-- ─── Config / Metadata ──────────────────────── -->
                <div :class="cn(v.sectionPx, v.sectionWrap && 'mx-4')">
                  <div :class="v.sectionWrap">
                    <span :class="cn('font-semibold mb-2 block', v.labelSize)">Details</span>
                    <div :class="cn('grid grid-cols-[auto_1fr] items-center', variant === 'dense' ? 'gap-x-3 gap-y-1' : 'gap-x-4 gap-y-2')">
                      <!-- Status -->
                      <span :class="cn('text-muted-foreground', v.bodySize)">Status</span>
                      <span :class="cn('flex items-center gap-1.5', v.bodySize, v.monoClass)">
                        <span :class="cn('size-2 rounded-full', statusColor[selectedTask.status])" />
                        {{ statusLabel[selectedTask.status] }}
                      </span>

                      <!-- Priority -->
                      <span :class="cn('text-muted-foreground', v.bodySize)">Priority</span>
                      <span :class="cn('flex items-center gap-1.5', v.bodySize, v.monoClass)">
                        <span :class="cn('size-2 rounded-full', priorityDot[selectedTask.priority])" />
                        {{ priorityLabel[selectedTask.priority] }}
                      </span>

                      <!-- Project -->
                      <span :class="cn('text-muted-foreground', v.bodySize)">Project</span>
                      <span :class="cn('flex items-center gap-1.5', v.bodySize)">
                        <Folder :size="12" class="text-muted-foreground" />
                        {{ selectedTask.project ?? '—' }}
                      </span>

                      <!-- Due Date -->
                      <span :class="cn('text-muted-foreground', v.bodySize)">Due Date</span>
                      <span :class="cn('flex items-center gap-1.5', v.bodySize, v.monoClass)">
                        <CalendarDays :size="12" class="text-muted-foreground" />
                        {{ selectedTask.dueDate ?? '—' }}
                      </span>

                      <!-- Tags -->
                      <span :class="cn('text-muted-foreground', v.bodySize)">Tags</span>
                      <div class="flex items-center gap-1 flex-wrap">
                        <span
                          v-for="tag in selectedTask.tags ?? []"
                          :key="tag"
                          :class="cn(
                            'inline-flex items-center gap-1 rounded-full border border-border px-2 py-0.5',
                            variant === 'dense' ? 'text-[10px]' : 'text-[11px]'
                          )"
                        >
                          <Tag :size="10" class="text-muted-foreground" />
                          {{ tag }}
                        </span>
                        <span v-if="!selectedTask.tags || selectedTask.tags.length === 0" :class="cn('text-muted-foreground', v.bodySize)">—</span>
                      </div>

                      <!-- ADO Link -->
                      <span :class="cn('text-muted-foreground', v.bodySize)">ADO Link</span>
                      <span v-if="selectedTask.adoType && selectedTask.adoNumber" :class="cn('flex items-center gap-1.5 text-blue-500 cursor-pointer hover:underline', v.bodySize, v.monoClass)">
                        <Landmark :size="12" />
                        {{ selectedTask.adoType }} #{{ selectedTask.adoNumber }}
                        <ExternalLink :size="10" />
                      </span>
                      <span v-else :class="cn('text-muted-foreground', v.bodySize)">—</span>
                    </div>
                  </div>
                </div>

                <!-- ─── Footer ─────────────────────────────────── -->
                <div :class="cn('flex items-center justify-between border-t border-border', v.sectionPx, 'py-2 mt-1')">
                  <span class="text-[11px] text-muted-foreground">
                    Created {{ selectedTask.createdAt ?? '—' }} · Updated {{ selectedTask.updatedAt ?? '—' }}
                  </span>
                  <Button variant="ghost" size="sm" class="h-6 text-[11px] text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1">
                    <Trash2 :size="12" />
                    Delete
                  </Button>
                </div>

                <!-- bottom spacer for scroll -->
                <div class="h-4" />
              </div>
            </ScrollArea>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
