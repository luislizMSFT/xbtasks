<script setup lang="ts">
import { ref, computed } from 'vue'
import { cn } from '@/lib/utils'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipProvider, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'
import {
  ChevronDown,
  ChevronRight,
  Circle,
  CircleDot,
  Eye,
  CheckCircle2,
  Octagon,
  GitPullRequest,
  ExternalLink,
  Activity,
} from 'lucide-vue-next'

// ─── Mock Data ──────────────────────────────────────────────────────────────────

interface Task {
  id: number
  title: string
  status: string
  priority: string
  tags: string
  adoId: string
  adoType: string
  subtasksDone: number
  subtasksTotal: number
  updatedAt: string
  blockedReason?: string
}

const mockTasks: Task[] = [
  { id: 1, title: 'Fix auth redirect loop', status: 'in_progress', priority: 'P0', tags: 'bug,urgent', adoId: 'ADO-48291', adoType: 'bug', subtasksDone: 2, subtasksTotal: 4, updatedAt: '30m ago' },
  { id: 2, title: 'Build sidebar navigation', status: 'in_progress', priority: 'P1', tags: 'ui', adoId: '', adoType: '', subtasksDone: 0, subtasksTotal: 0, updatedAt: '1h ago' },
  { id: 3, title: 'Review PR #234 — schema changes', status: 'in_review', priority: 'P1', tags: 'review', adoId: 'ADO-48350', adoType: 'task', subtasksDone: 0, subtasksTotal: 0, updatedAt: '2h ago' },
  { id: 4, title: 'Add rate limiting to gateway', status: 'todo', priority: 'P1', tags: 'backend', adoId: 'ADO-48400', adoType: 'task', subtasksDone: 0, subtasksTotal: 3, updatedAt: '3h ago' },
  { id: 5, title: 'Design system tokens for dark mode', status: 'todo', priority: 'P2', tags: 'design', adoId: '', adoType: '', subtasksDone: 0, subtasksTotal: 0, updatedAt: '5h ago' },
  { id: 6, title: 'Migrate auth to MSAL v3', status: 'blocked', priority: 'P0', tags: 'auth,blocked', adoId: 'ADO-48100', adoType: 'bug', subtasksDone: 1, subtasksTotal: 2, blockedReason: 'Waiting on SDK release', updatedAt: '1d ago' },
  { id: 7, title: 'Write E2E tests for task CRUD', status: 'done', priority: 'P2', tags: 'testing', adoId: '', adoType: '', subtasksDone: 5, subtasksTotal: 5, updatedAt: '2d ago' },
  { id: 8, title: 'Update team onboarding docs', status: 'done', priority: 'P3', tags: 'docs', adoId: 'ADO-47900', adoType: 'user_story', subtasksDone: 3, subtasksTotal: 3, updatedAt: '3d ago' },
]

// ─── Shared Helpers ─────────────────────────────────────────────────────────────

const statusOrder = ['in_progress', 'in_review', 'todo', 'blocked', 'done'] as const

const statusConfig: Record<string, { label: string; color: string; dotClass: string; borderClass: string; icon: typeof Circle }> = {
  in_progress: { label: 'In Progress', color: 'text-blue-500', dotClass: 'bg-blue-500', borderClass: 'border-blue-500', icon: CircleDot },
  in_review:   { label: 'In Review',   color: 'text-violet-500', dotClass: 'bg-violet-500', borderClass: 'border-violet-500', icon: Eye },
  todo:        { label: 'To Do',       color: 'text-zinc-400', dotClass: 'bg-zinc-400', borderClass: 'border-zinc-400', icon: Circle },
  blocked:     { label: 'Blocked',     color: 'text-red-500', dotClass: 'bg-red-500', borderClass: 'border-red-500', icon: Octagon },
  done:        { label: 'Done',        color: 'text-emerald-500', dotClass: 'bg-emerald-500', borderClass: 'border-emerald-500', icon: CheckCircle2 },
}

const priorityConfig: Record<string, { classes: string }> = {
  P0: { classes: 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20' },
  P1: { classes: 'bg-orange-500/15 text-orange-600 dark:text-orange-400 border-orange-500/20' },
  P2: { classes: 'bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/20' },
  P3: { classes: 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20' },
}

const adoTypeIcons: Record<string, string> = {
  bug: '🐛',
  task: '✅',
  user_story: '📖',
  deliverable: '📦',
  scenario: '🎯',
}

function groupedByStatus(tasks: Task[]) {
  const groups: Record<string, Task[]> = {}
  for (const s of statusOrder) groups[s] = []
  for (const t of tasks) {
    if (groups[t.status]) groups[t.status].push(t)
  }
  return groups
}

const grouped = computed(() => groupedByStatus(mockTasks))

// ─── Variant A State ────────────────────────────────────────────────────────────

const collapsedSections = ref<Record<string, boolean>>({})
const selectedTaskId = ref<number | null>(1)

function toggleSection(status: string) {
  collapsedSections.value[status] = !collapsedSections.value[status]
}

// ─── Variant C State ────────────────────────────────────────────────────────────

const previewTaskId = ref<number>(1)
const previewTask = computed(() => mockTasks.find(t => t.id === previewTaskId.value) ?? mockTasks[0])

const mockSubtasks = [
  { id: 1, title: 'Identify redirect loop entry point', done: true },
  { id: 2, title: 'Add token refresh guard', done: true },
  { id: 3, title: 'Write regression test', done: false },
  { id: 4, title: 'Verify with QA on staging', done: false },
]

const mockPRs = [
  { id: 234, title: 'fix: auth redirect loop guard', status: 'open', additions: 47, deletions: 12 },
  { id: 228, title: 'refactor: extract token service', status: 'merged', additions: 120, deletions: 88 },
]

const mockComments = [
  { author: 'Sarah K.', avatar: 'SK', time: '2h ago', text: 'Looks like the root cause is in the MSAL callback handler. The redirect URI is being double-encoded.' },
  { author: 'You', avatar: 'YO', time: '1h ago', text: 'Good catch — I\'ve pushed a fix to the token guard. Can you re-test on staging?' },
  { author: 'Sarah K.', avatar: 'SK', time: '30m ago', text: 'Confirmed fixed on staging. LGTM 👍' },
]

const mockActivity = [
  { time: '30m ago', event: 'Status changed to In Progress' },
  { time: '1h ago', event: 'PR #234 opened' },
  { time: '2h ago', event: 'Comment added by Sarah K.' },
  { time: '1d ago', event: 'Created from ADO-48291' },
]
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <div class="px-6 py-6 max-w-[1400px] mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-xl font-semibold text-foreground">Task Layout Playground</h1>
          <p class="text-sm text-muted-foreground mt-0.5">Compare layouts and pick your favorite</p>
        </div>
        <Badge variant="outline" class="text-xs text-muted-foreground">
          {{ mockTasks.length }} tasks · 3 variants
        </Badge>
      </div>

      <!-- Tabs -->
      <Tabs default-value="variant-a" class="w-full">
        <TabsList class="mb-6">
          <TabsTrigger value="variant-a">Things 3 Classic</TabsTrigger>
          <TabsTrigger value="variant-b">Linear / Issues</TabsTrigger>
          <TabsTrigger value="variant-c">Todoist Hybrid</TabsTrigger>
        </TabsList>

        <!-- ════════════════════════════════════════════════════════════════════ -->
        <!-- VARIANT A: Things 3 Classic                                        -->
        <!-- ════════════════════════════════════════════════════════════════════ -->
        <TabsContent value="variant-a">
          <div class="rounded-xl border border-border bg-card">
            <template v-for="(status, sIdx) in statusOrder" :key="status">
              <div v-if="grouped[status].length > 0">
                <!-- Section Header -->
                <button
                  @click="toggleSection(status)"
                  :class="cn(
                    'w-full flex items-center gap-2.5 px-5 py-3 hover:bg-muted/50 transition-colors',
                    sIdx > 0 ? 'border-t border-border' : ''
                  )"
                >
                  <component
                    :is="collapsedSections[status] ? ChevronRight : ChevronDown"
                    :size="14"
                    class="text-muted-foreground/60 shrink-0"
                  />
                  <span
                    :class="cn('size-2 rounded-full shrink-0', statusConfig[status].dotClass)"
                  />
                  <span class="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground">
                    {{ statusConfig[status].label }}
                  </span>
                  <span class="text-[11px] text-muted-foreground/50 font-medium">
                    {{ grouped[status].length }}
                  </span>
                </button>

                <!-- Task Rows -->
                <div v-if="!collapsedSections[status]">
                  <div
                    v-for="task in grouped[status]"
                    :key="task.id"
                    @click="selectedTaskId = task.id"
                    :class="cn(
                      'group flex items-center gap-3 px-5 py-3 cursor-pointer transition-all',
                      'hover:bg-muted/40',
                      selectedTaskId === task.id
                        ? 'border-l-[3px] border-l-blue-500 bg-blue-500/[0.04] pl-[17px]'
                        : 'border-l-[3px] border-l-transparent'
                    )"
                  >
                    <!-- Status Checkbox Circle -->
                    <TooltipProvider :delay-duration="300">
                      <Tooltip>
                        <TooltipTrigger as-child>
                          <div
                            :class="cn(
                              'size-[18px] rounded-full border-2 shrink-0 flex items-center justify-center transition-colors',
                              statusConfig[task.status].borderClass,
                              task.status === 'done' ? statusConfig[task.status].dotClass : ''
                            )"
                          >
                            <CheckCircle2
                              v-if="task.status === 'done'"
                              :size="10"
                              class="text-white"
                              :stroke-width="3"
                            />
                            <div
                              v-else-if="task.status === 'in_progress'"
                              class="size-2 rounded-full bg-blue-500"
                            />
                          </div>
                        </TooltipTrigger>
                        <TooltipContent side="left" :side-offset="8">
                          <p class="text-xs">{{ statusConfig[task.status].label }}</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>

                    <!-- Title + Progress -->
                    <div class="flex-1 min-w-0">
                      <span
                        :class="cn(
                          'text-sm font-medium block truncate',
                          task.status === 'done' ? 'text-muted-foreground line-through decoration-muted-foreground/40' : 'text-foreground'
                        )"
                      >
                        {{ task.title }}
                      </span>
                      <div
                        v-if="task.subtasksTotal > 0"
                        class="mt-1.5 flex items-center gap-2"
                      >
                        <div class="h-[3px] w-24 rounded-full bg-muted overflow-hidden">
                          <div
                            class="h-full rounded-full transition-all duration-500"
                            :class="task.subtasksDone === task.subtasksTotal ? 'bg-emerald-500' : 'bg-blue-500'"
                            :style="{ width: `${(task.subtasksDone / task.subtasksTotal) * 100}%` }"
                          />
                        </div>
                        <span class="text-[10px] text-muted-foreground/60 tabular-nums">
                          {{ task.subtasksDone }}/{{ task.subtasksTotal }}
                        </span>
                      </div>
                    </div>

                    <!-- Priority Badge -->
                    <Badge
                      variant="outline"
                      :class="cn(
                        'text-[10px] font-bold uppercase tracking-wider px-1.5 py-0 h-5 shrink-0',
                        priorityConfig[task.priority]?.classes
                      )"
                    >
                      {{ task.priority }}
                    </Badge>

                    <!-- ADO Link Icon -->
                    <TooltipProvider :delay-duration="300">
                      <Tooltip>
                        <TooltipTrigger as-child>
                          <span class="text-sm shrink-0 w-5 text-center">
                            <template v-if="task.adoId">
                              {{ adoTypeIcons[task.adoType] || '🔗' }}
                            </template>
                            <template v-else>
                              <Circle :size="14" class="text-muted-foreground/25 inline-block" />
                            </template>
                          </span>
                        </TooltipTrigger>
                        <TooltipContent side="top" :side-offset="4">
                          <p class="text-xs">{{ task.adoId || 'Not linked to ADO' }}</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>

                    <!-- Relative Time -->
                    <span class="text-[11px] text-muted-foreground/50 tabular-nums w-14 text-right shrink-0">
                      {{ task.updatedAt }}
                    </span>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </TabsContent>

        <!-- ════════════════════════════════════════════════════════════════════ -->
        <!-- VARIANT B: Linear / GitHub Issues                                  -->
        <!-- ════════════════════════════════════════════════════════════════════ -->
        <TabsContent value="variant-b">
          <div class="space-y-6">
            <template v-for="status in statusOrder" :key="status">
              <div v-if="grouped[status].length > 0">
                <!-- Section Label -->
                <div class="flex items-center gap-2 mb-2">
                  <Separator class="flex-1" />
                  <span class="text-[11px] font-medium uppercase tracking-wider text-muted-foreground/60 shrink-0 flex items-center gap-1.5">
                    <span :class="cn('size-1.5 rounded-[2px]', statusConfig[status].dotClass)" />
                    {{ statusConfig[status].label }}
                    <span class="text-muted-foreground/40">{{ grouped[status].length }}</span>
                  </span>
                  <Separator class="flex-1" />
                </div>

                <!-- Task Cards -->
                <div class="space-y-1">
                  <div
                    v-for="task in grouped[status]"
                    :key="task.id"
                    :class="cn(
                      'group flex items-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-card',
                      'hover:bg-muted/30 hover:border-border/80 transition-all cursor-pointer',
                      task.status === 'blocked' ? 'border-l-[3px] border-l-red-500 pl-[13px]' : ''
                    )"
                  >
                    <!-- Status Icon -->
                    <component
                      :is="statusConfig[task.status].icon"
                      :size="16"
                      :class="cn('shrink-0', statusConfig[task.status].color)"
                      :stroke-width="2"
                    />

                    <!-- Title -->
                    <span
                      :class="cn(
                        'text-sm font-semibold truncate',
                        task.status === 'done' ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground'
                      )"
                    >
                      {{ task.title }}
                    </span>

                    <!-- Blocked Reason -->
                    <Badge
                      v-if="task.blockedReason"
                      variant="outline"
                      class="text-[10px] bg-red-500/10 text-red-500 border-red-500/20 shrink-0"
                    >
                      {{ task.blockedReason }}
                    </Badge>

                    <!-- Inline Tags -->
                    <div class="flex items-center gap-1 shrink-0">
                      <Badge
                        v-for="tag in task.tags.split(',')"
                        :key="tag"
                        variant="secondary"
                        class="text-[10px] font-normal px-1.5 py-0 h-[18px] rounded-md"
                      >
                        <span class="opacity-40">#</span>{{ tag }}
                      </Badge>
                    </div>

                    <div class="flex-1" />

                    <!-- Priority -->
                    <Badge
                      variant="outline"
                      :class="cn(
                        'text-[10px] font-bold uppercase tracking-wider px-1.5 py-0 h-5 shrink-0',
                        priorityConfig[task.priority]?.classes
                      )"
                    >
                      {{ task.priority }}
                    </Badge>

                    <!-- ADO Type Badge -->
                    <Badge
                      v-if="task.adoId"
                      variant="outline"
                      class="text-[10px] font-medium gap-1 bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20 shrink-0"
                    >
                      {{ adoTypeIcons[task.adoType] || '🔗' }}
                      {{ task.adoId }}
                    </Badge>

                    <!-- Progress Text -->
                    <span
                      v-if="task.subtasksTotal > 0"
                      :class="cn(
                        'text-[11px] tabular-nums shrink-0 font-medium',
                        task.subtasksDone === task.subtasksTotal ? 'text-emerald-500' : 'text-muted-foreground'
                      )"
                    >
                      {{ task.subtasksDone }}/{{ task.subtasksTotal }}
                    </span>

                    <!-- Time -->
                    <span class="text-[11px] text-muted-foreground/50 tabular-nums shrink-0 w-12 text-right">
                      {{ task.updatedAt }}
                    </span>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </TabsContent>

        <!-- ════════════════════════════════════════════════════════════════════ -->
        <!-- VARIANT C: Todoist / TickTick Hybrid                               -->
        <!-- ════════════════════════════════════════════════════════════════════ -->
        <TabsContent value="variant-c">
          <div class="flex gap-0 rounded-xl border border-border bg-card overflow-hidden" style="height: 640px">
            <!-- Left Panel: Task List -->
            <div class="w-[55%] border-r border-border flex flex-col">
              <div class="px-4 py-3 border-b border-border bg-muted/30">
                <span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  All Tasks · {{ mockTasks.length }}
                </span>
              </div>
              <ScrollArea class="flex-1">
                <div>
                  <div
                    v-for="task in mockTasks"
                    :key="task.id"
                    @click="previewTaskId = task.id"
                    :class="cn(
                      'group flex items-center gap-2.5 px-3 py-2.5 cursor-pointer transition-all border-b border-border/50',
                      'hover:bg-muted/40',
                      previewTaskId === task.id ? 'bg-blue-500/[0.06]' : ''
                    )"
                  >
                    <!-- Priority Stripe -->
                    <div
                      :class="cn(
                        'w-[3px] self-stretch rounded-full shrink-0',
                        task.priority === 'P0' ? 'bg-red-500' :
                        task.priority === 'P1' ? 'bg-orange-500' :
                        task.priority === 'P2' ? 'bg-amber-500' :
                        'bg-zinc-300 dark:bg-zinc-600'
                      )"
                    />

                    <!-- Status Checkbox -->
                    <div
                      :class="cn(
                        'size-4 rounded-full border-[1.5px] shrink-0 flex items-center justify-center',
                        statusConfig[task.status].borderClass,
                        task.status === 'done' ? statusConfig[task.status].dotClass : ''
                      )"
                    >
                      <CheckCircle2
                        v-if="task.status === 'done'"
                        :size="8"
                        class="text-white"
                        :stroke-width="3"
                      />
                    </div>

                    <!-- Title -->
                    <span
                      :class="cn(
                        'text-[13px] font-medium truncate flex-1',
                        task.status === 'done' ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground'
                      )"
                    >
                      {{ task.title }}
                    </span>

                    <!-- ADO Icon -->
                    <span v-if="task.adoId" class="text-xs shrink-0 opacity-60">
                      {{ adoTypeIcons[task.adoType] || '🔗' }}
                    </span>
                    <Circle v-else :size="12" class="text-muted-foreground/20 shrink-0" />

                    <!-- Time -->
                    <span class="text-[10px] text-muted-foreground/40 tabular-nums shrink-0">
                      {{ task.updatedAt }}
                    </span>
                  </div>
                </div>
              </ScrollArea>
            </div>

            <!-- Right Panel: Preview -->
            <div class="w-[45%] flex flex-col bg-card">
              <!-- Preview Header -->
              <div class="px-5 py-4 border-b border-border">
                <div class="flex items-start justify-between gap-3">
                  <div class="flex-1 min-w-0">
                    <h2 class="text-base font-semibold text-foreground leading-snug">
                      {{ previewTask.title }}
                    </h2>
                    <div class="flex items-center gap-2 mt-2 flex-wrap">
                      <Badge
                        variant="outline"
                        :class="cn(
                          'text-[10px] font-bold uppercase tracking-wider px-1.5 py-0 h-5',
                          priorityConfig[previewTask.priority]?.classes
                        )"
                      >
                        {{ previewTask.priority }}
                      </Badge>
                      <Badge
                        variant="outline"
                        :class="cn(
                          'text-[10px] font-medium gap-1 h-5',
                          statusConfig[previewTask.status]?.color
                        )"
                      >
                        <component :is="statusConfig[previewTask.status].icon" :size="11" :stroke-width="2" />
                        {{ statusConfig[previewTask.status].label }}
                      </Badge>
                      <Badge
                        v-if="previewTask.adoId"
                        variant="outline"
                        class="text-[10px] font-medium gap-1 h-5 bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20"
                      >
                        {{ adoTypeIcons[previewTask.adoType] }}
                        {{ previewTask.adoId }}
                      </Badge>
                    </div>
                  </div>
                  <button class="p-1.5 rounded-md hover:bg-muted transition-colors text-muted-foreground">
                    <ExternalLink :size="14" />
                  </button>
                </div>

                <!-- Progress Bar in Header -->
                <div
                  v-if="previewTask.subtasksTotal > 0"
                  class="mt-3 flex items-center gap-2.5"
                >
                  <div class="h-1.5 flex-1 rounded-full bg-muted overflow-hidden">
                    <div
                      class="h-full rounded-full transition-all duration-500"
                      :class="previewTask.subtasksDone === previewTask.subtasksTotal ? 'bg-emerald-500' : 'bg-blue-500'"
                      :style="{ width: `${(previewTask.subtasksDone / previewTask.subtasksTotal) * 100}%` }"
                    />
                  </div>
                  <span class="text-[11px] text-muted-foreground tabular-nums">
                    {{ previewTask.subtasksDone }}/{{ previewTask.subtasksTotal }}
                  </span>
                </div>
              </div>

              <!-- Preview Content -->
              <ScrollArea class="flex-1">
                <div class="px-5 py-4 space-y-5">
                  <!-- Subtasks -->
                  <div>
                    <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">
                      Subtasks
                    </h3>
                    <div class="space-y-1">
                      <div
                        v-for="sub in mockSubtasks"
                        :key="sub.id"
                        class="flex items-center gap-2.5 py-1.5 px-2 rounded-md hover:bg-muted/40 transition-colors"
                      >
                        <div
                          :class="cn(
                            'size-4 rounded-[4px] border-[1.5px] shrink-0 flex items-center justify-center transition-colors',
                            sub.done
                              ? 'bg-emerald-500 border-emerald-500'
                              : 'border-muted-foreground/30'
                          )"
                        >
                          <CheckCircle2 v-if="sub.done" :size="10" class="text-white" :stroke-width="3" />
                        </div>
                        <span
                          :class="cn(
                            'text-[13px]',
                            sub.done ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground'
                          )"
                        >
                          {{ sub.title }}
                        </span>
                      </div>
                    </div>
                  </div>

                  <Separator />

                  <!-- Pull Requests -->
                  <div>
                    <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">
                      Pull Requests
                    </h3>
                    <div class="space-y-1.5">
                      <div
                        v-for="pr in mockPRs"
                        :key="pr.id"
                        class="flex items-center gap-2.5 py-2 px-3 rounded-lg border border-border hover:bg-muted/30 transition-colors cursor-pointer"
                      >
                        <GitPullRequest
                          :size="14"
                          :class="pr.status === 'merged' ? 'text-violet-500' : 'text-emerald-500'"
                        />
                        <div class="flex-1 min-w-0">
                          <span class="text-[13px] font-medium text-foreground truncate block">
                            {{ pr.title }}
                          </span>
                          <span class="text-[10px] text-muted-foreground">
                            #{{ pr.id }} ·
                            <span class="text-emerald-500">+{{ pr.additions }}</span>
                            <span class="text-red-500 ml-0.5">-{{ pr.deletions }}</span>
                          </span>
                        </div>
                        <Badge
                          variant="outline"
                          :class="cn(
                            'text-[10px] h-5 capitalize',
                            pr.status === 'merged'
                              ? 'bg-violet-500/10 text-violet-500 border-violet-500/20'
                              : 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20'
                          )"
                        >
                          {{ pr.status }}
                        </Badge>
                      </div>
                    </div>
                  </div>

                  <Separator />

                  <!-- Comments -->
                  <div>
                    <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">
                      Comments
                    </h3>
                    <div class="space-y-3">
                      <div
                        v-for="(comment, idx) in mockComments"
                        :key="idx"
                        class="flex gap-2.5"
                      >
                        <div
                          :class="cn(
                            'size-7 rounded-full shrink-0 flex items-center justify-center text-[10px] font-bold',
                            comment.author === 'You'
                              ? 'bg-blue-500/15 text-blue-600 dark:text-blue-400'
                              : 'bg-muted text-muted-foreground'
                          )"
                        >
                          {{ comment.avatar }}
                        </div>
                        <div class="flex-1 min-w-0">
                          <div class="flex items-center gap-2">
                            <span class="text-[12px] font-semibold text-foreground">{{ comment.author }}</span>
                            <span class="text-[10px] text-muted-foreground/50">{{ comment.time }}</span>
                          </div>
                          <p class="text-[13px] text-muted-foreground leading-relaxed mt-0.5">
                            {{ comment.text }}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <Separator />

                  <!-- Activity Timeline -->
                  <div>
                    <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">
                      Activity
                    </h3>
                    <div class="space-y-0">
                      <div
                        v-for="(item, idx) in mockActivity"
                        :key="idx"
                        class="flex items-start gap-2.5 py-1.5 relative"
                      >
                        <div
                          v-if="idx < mockActivity.length - 1"
                          class="absolute left-[7px] top-[18px] w-px h-[calc(100%-6px)] bg-border"
                        />
                        <div class="size-[15px] rounded-full bg-muted border border-border shrink-0 flex items-center justify-center z-10">
                          <Activity :size="8" class="text-muted-foreground" />
                        </div>
                        <div class="flex-1 min-w-0">
                          <span class="text-[12px] text-muted-foreground">{{ item.event }}</span>
                        </div>
                        <span class="text-[10px] text-muted-foreground/40 tabular-nums shrink-0">
                          {{ item.time }}
                        </span>
                      </div>
                    </div>
                  </div>

                  <Separator />

                  <!-- Config / Meta -->
                  <div>
                    <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">
                      Details
                    </h3>
                    <div class="grid grid-cols-2 gap-y-2 gap-x-4 text-[12px]">
                      <div class="text-muted-foreground/60">Priority</div>
                      <div class="text-foreground font-medium">{{ previewTask.priority }}</div>
                      <div class="text-muted-foreground/60">Status</div>
                      <div class="text-foreground font-medium">{{ statusConfig[previewTask.status].label }}</div>
                      <div class="text-muted-foreground/60">Tags</div>
                      <div class="flex gap-1 flex-wrap">
                        <Badge
                          v-for="tag in previewTask.tags.split(',')"
                          :key="tag"
                          variant="secondary"
                          class="text-[10px] px-1.5 py-0 h-[16px] rounded"
                        >
                          {{ tag }}
                        </Badge>
                      </div>
                      <div class="text-muted-foreground/60">Updated</div>
                      <div class="text-foreground font-medium">{{ previewTask.updatedAt }}</div>
                      <template v-if="previewTask.adoId">
                        <div class="text-muted-foreground/60">ADO Item</div>
                        <div class="text-blue-600 dark:text-blue-400 font-medium">{{ previewTask.adoId }}</div>
                      </template>
                    </div>
                  </div>
                </div>
              </ScrollArea>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  </ScrollArea>
</template>
