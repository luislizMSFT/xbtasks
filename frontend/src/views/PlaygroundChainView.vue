<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Component } from 'vue'
import { cn } from '@/lib/utils'
import PageHeader from '@/components/PageHeader.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import {
  Check,
  Circle,
  CircleDot,
  GitPullRequest,
  Package,
  Rocket,
  X,
  AlertTriangle,
  Bug,
  ClipboardList,
  Clock,
  Bell,
  Zap,
  ExternalLink,
  ChevronDown,
  ChevronRight,
  Eye,
  MessageSquare,
  FileText,
  Minus,
  ArrowRight,
  TriangleAlert,
  Megaphone,
  Merge,
  ListChecks,
  Filter,
} from 'lucide-vue-next'

// ---------- Types ----------

type StageStatus = 'done' | 'in_progress' | 'failed' | 'pending' | 'skipped'

interface ChainStage {
  label: string
  status: StageStatus
  icon: Component
  detail: string
  link?: string
}

interface Blocker {
  reason: string
  owner?: string
  actionLabel?: string
}

interface FollowUp {
  description: string
  type: 'manual' | 'auto'
  status: 'waiting' | 'triggered' | 'expired'
  sourceTask: string
}

interface MockTask {
  id: string
  title: string
  priority: string
  type: 'bug' | 'task'
  adoId: string
  time: string
  chain: ChainStage[]
  blocker?: Blocker
  followUps?: FollowUp[]
}

// ---------- Filter state ----------

const filterMode = ref<'all' | 'blocked' | 'inflight'>('all')
const expandedTask = ref<string | null>(null)

function toggleExpand(id: string) {
  expandedTask.value = expandedTask.value === id ? null : id
}

// ---------- Mock data ----------

const tasks: MockTask[] = [
  {
    id: 't1',
    title: 'Fix auth redirect loop',
    priority: 'P0',
    type: 'bug',
    adoId: '#48291',
    time: '30m',
    chain: [
      { label: 'Task', status: 'in_progress', icon: ClipboardList, detail: 'In Progress' },
      { label: 'PR #234', status: 'in_progress', icon: GitPullRequest, detail: 'Active (Alex \u2713, Sam pending)', link: 'PR #234' },
      { label: 'Build', status: 'pending', icon: Package, detail: 'Queued' },
      { label: 'Deploy', status: 'pending', icon: Rocket, detail: 'Pending' },
    ],
  },
  {
    id: 't2',
    title: 'Schema migration v2',
    priority: 'P1',
    type: 'task',
    adoId: '#48350',
    time: '2h',
    chain: [
      { label: 'Task', status: 'done', icon: ClipboardList, detail: 'Done' },
      { label: 'PR #298', status: 'done', icon: GitPullRequest, detail: 'Approved', link: 'PR #298' },
      { label: 'Build #1245', status: 'done', icon: Package, detail: 'Succeeded', link: 'Build #1245' },
      { label: 'Deploy', status: 'in_progress', icon: Rocket, detail: 'Running (2m elapsed)' },
    ],
  },
  {
    id: 't3',
    title: 'Rate limit middleware',
    priority: 'P1',
    type: 'task',
    adoId: '#48400',
    time: '3h',
    chain: [
      { label: 'Task', status: 'done', icon: ClipboardList, detail: 'Done' },
      { label: 'PR #280', status: 'failed', icon: GitPullRequest, detail: 'Changes Requested', link: 'PR #280' },
      { label: 'Build', status: 'skipped', icon: Package, detail: 'N/A' },
      { label: 'Deploy', status: 'skipped', icon: Rocket, detail: 'N/A' },
    ],
    blocker: {
      reason: 'Changes requested by Alex K \u2014 "Need rate limit tests"',
      owner: 'Alex K',
      actionLabel: 'View PR',
    },
  },
  {
    id: 't4',
    title: 'Migrate auth to MSAL v3',
    priority: 'P0',
    type: 'bug',
    adoId: '#48100',
    time: '1d',
    chain: [
      { label: 'Task', status: 'failed', icon: ClipboardList, detail: 'Blocked (external)' },
      { label: 'PR', status: 'skipped', icon: GitPullRequest, detail: 'N/A' },
      { label: 'Build', status: 'skipped', icon: Package, detail: 'N/A' },
      { label: 'Deploy', status: 'skipped', icon: Rocket, detail: 'N/A' },
    ],
    blocker: {
      reason: 'Blocked by: Bug #48292 (Alex K) \u2014 Token expiry fix',
      owner: 'Alex K',
      actionLabel: 'Open in ADO',
    },
    followUps: [
      {
        description: 'Notify when #48292 resolved',
        type: 'auto',
        status: 'waiting',
        sourceTask: '#48100',
      },
    ],
  },
  {
    id: 't5',
    title: 'Add telemetry hooks',
    priority: 'P2',
    type: 'task',
    adoId: '#48500',
    time: '4h',
    chain: [
      { label: 'Task', status: 'done', icon: ClipboardList, detail: 'Done' },
      { label: 'PR #301', status: 'in_progress', icon: GitPullRequest, detail: 'Active', link: 'PR #301' },
      { label: 'Build #1244', status: 'in_progress', icon: Package, detail: 'Running', link: 'Build #1244' },
      { label: 'Deploy', status: 'pending', icon: Rocket, detail: 'Pending' },
    ],
  },
  {
    id: 't6',
    title: 'Update deployment docs',
    priority: 'P3',
    type: 'task',
    adoId: '',
    time: '45m',
    chain: [
      { label: 'Task', status: 'done', icon: ClipboardList, detail: 'Done' },
      { label: 'Verified', status: 'done', icon: Check, detail: 'Verified' },
    ],
  },
]

// ---------- Blocker cards ----------

interface BlockerCard {
  title: string
  description: string
  owner?: string
  detail: string
  timeStuck?: string
  actions: string[]
}

const blockerCards: BlockerCard[] = [
  {
    title: 'Rate limit middleware \u2192 PR #280',
    description: 'PR has changes requested',
    owner: 'Alex K',
    detail: 'Requested changes on code review',
    timeStuck: '2h',
    actions: ['View PR', 'Ping reviewer'],
  },
  {
    title: 'MSAL v3 migration \u2192 Bug #48292',
    description: 'Depends on another team member\'s ADO item',
    owner: 'Alex K',
    detail: 'Status in ADO: Active \u2014 Follow-up: Reminder active (daily check)',
    actions: ['Open in ADO', 'Create follow-up task'],
  },
  {
    title: 'Deploy pipeline #1240 failed',
    description: 'Pipeline failure blocking deployment',
    detail: 'Error: Test stage failed \u2014 3 tests failing',
    actions: ['View logs', 'Create fix task'],
  },
]

// ---------- Follow-up items ----------

const followUpItems: FollowUp[] = [
  {
    description: 'Notify when Bug #48292 resolved',
    type: 'manual',
    status: 'waiting',
    sourceTask: '#48100',
  },
  {
    description: 'Verify deployment after schema migration merges',
    type: 'auto',
    status: 'waiting',
    sourceTask: '#48350',
  },
]

// ---------- Suggested actions ----------

interface SuggestedAction {
  icon: Component
  description: string
  actionLabel: string
}

const suggestedActions: SuggestedAction[] = [
  { icon: Merge, description: 'PR #298 is approved \u2014 ready to merge and deploy', actionLabel: 'Merge PR' },
  { icon: ListChecks, description: 'Build #1245 succeeded \u2014 verify deployment', actionLabel: 'Create verify task' },
  { icon: Bug, description: '3 tests failing in Pipeline #1240 \u2014 create fix task', actionLabel: 'Create task from failure' },
]

// ---------- Computed ----------

const filteredTasks = computed(() => {
  if (filterMode.value === 'all') return tasks
  if (filterMode.value === 'blocked') {
    return tasks.filter((t) => t.chain.some((s) => s.status === 'failed'))
  }
  // inflight: has at least one in_progress stage
  return tasks.filter((t) => t.chain.some((s) => s.status === 'in_progress'))
})

// ---------- Helpers ----------

function priorityColor(p: string): string {
  switch (p) {
    case 'P0': return 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20'
    case 'P1': return 'bg-orange-500/15 text-orange-600 dark:text-orange-400 border-orange-500/20'
    case 'P2': return 'bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/20'
    default: return 'bg-zinc-500/10 text-zinc-500 border-zinc-500/20'
  }
}

function stageCircleClass(status: StageStatus): string {
  switch (status) {
    case 'done': return 'bg-emerald-500 border-emerald-500 text-white'
    case 'in_progress': return 'bg-blue-500 border-blue-500 text-white animate-pulse'
    case 'failed': return 'bg-red-500 border-red-500 text-white'
    case 'pending': return 'bg-transparent border-zinc-300 dark:border-zinc-600 text-zinc-400'
    case 'skipped': return 'bg-transparent border-zinc-200 dark:border-zinc-700 border-dashed text-zinc-300 dark:text-zinc-600'
  }
}

function stageLineClass(leftStatus: StageStatus): string {
  switch (leftStatus) {
    case 'done': return 'bg-emerald-500'
    case 'in_progress': return 'bg-blue-500'
    case 'failed': return 'bg-red-500'
    default: return 'bg-zinc-200 dark:bg-zinc-700'
  }
}

function stageLineDashed(leftStatus: StageStatus): boolean {
  return leftStatus === 'pending' || leftStatus === 'skipped'
}

function chainSummaryDots(chain: ChainStage[]): StageStatus[] {
  return chain.map((s) => s.status)
}

function isTaskDone(task: MockTask): boolean {
  return task.chain.every((s) => s.status === 'done')
}

function isTaskBlocked(task: MockTask): boolean {
  return task.chain.some((s) => s.status === 'failed')
}
</script>

<template>
  <div class="flex-1 flex flex-col min-h-0 h-full">
    <!-- Page header -->
    <PageHeader>
      <template #left>
        <h1 class="text-sm font-semibold text-foreground">Work Chains</h1>
        <Separator orientation="vertical" class="h-4 mx-1" />
        <span class="text-xs text-muted-foreground">Dependency flow visualization</span>
      </template>
      <template #right>
        <div class="flex items-center gap-1">
          <Button
            v-for="f in (['all', 'blocked', 'inflight'] as const)"
            :key="f"
            :variant="filterMode === f ? 'secondary' : 'ghost'"
            size="sm"
            class="h-7 text-xs px-2.5"
            @click="filterMode = f"
          >
            <Filter v-if="f === 'all'" :size="12" class="mr-1" />
            <AlertTriangle v-if="f === 'blocked'" :size="12" class="mr-1" />
            <Rocket v-if="f === 'inflight'" :size="12" class="mr-1" />
            {{ f === 'all' ? 'Show all' : f === 'blocked' ? 'Blocked only' : 'In flight' }}
          </Button>
        </div>
      </template>
    </PageHeader>

    <!-- Main scrollable content -->
    <ScrollArea class="flex-1">
      <div class="px-5 py-5 max-w-[1400px] mx-auto space-y-6">

        <!-- ============================================================ -->
        <!-- SECTION 1: Task List with Flow Chains                        -->
        <!-- ============================================================ -->
        <section>
          <div class="flex items-center gap-2 mb-3">
            <ClipboardList :size="14" class="text-muted-foreground" />
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Task Chains
            </h2>
            <Badge variant="outline" class="text-[10px] h-4 px-1.5 ml-1">
              {{ filteredTasks.length }}
            </Badge>
          </div>

          <div class="border border-border rounded-lg bg-card divide-y divide-border overflow-hidden">
            <div
              v-for="task in filteredTasks"
              :key="task.id"
            >
              <!-- Task row -->
              <button
                class="w-full flex items-center gap-3 px-3.5 py-2.5 text-left hover:bg-muted/40 transition-colors"
                @click="toggleExpand(task.id)"
              >
                <!-- Expand icon -->
                <component
                  :is="expandedTask === task.id ? ChevronDown : ChevronRight"
                  :size="13"
                  class="text-muted-foreground shrink-0"
                />

                <!-- Done checkbox circle -->
                <div
                  :class="cn(
                    'w-[18px] h-[18px] rounded-full border-[1.5px] flex items-center justify-center shrink-0',
                    isTaskDone(task)
                      ? 'bg-emerald-500 border-emerald-500'
                      : isTaskBlocked(task)
                        ? 'border-red-400 bg-red-500/10'
                        : 'border-zinc-300 dark:border-zinc-600'
                  )"
                >
                  <Check v-if="isTaskDone(task)" :size="11" class="text-white" />
                  <X v-else-if="isTaskBlocked(task)" :size="10" class="text-red-500" />
                </div>

                <!-- Title -->
                <span class="text-[13px] text-foreground truncate flex-1 min-w-0">
                  {{ task.title }}
                </span>

                <!-- ADO Badge -->
                <Badge
                  v-if="task.adoId"
                  variant="outline"
                  class="text-[10px] h-[18px] px-1.5 gap-1 shrink-0 bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20"
                >
                  <component :is="task.type === 'bug' ? Bug : ClipboardList" :size="9" />
                  {{ task.adoId }}
                </Badge>

                <!-- Priority badge -->
                <Badge
                  variant="outline"
                  :class="cn('text-[10px] h-[18px] px-1.5 shrink-0', priorityColor(task.priority))"
                >
                  {{ task.priority }}
                </Badge>

                <!-- Mini flow dots -->
                <TooltipProvider :delay-duration="200">
                  <div class="flex items-center gap-[3px] shrink-0 ml-1">
                    <template v-for="(dotStatus, dIdx) in chainSummaryDots(task.chain)" :key="dIdx">
                      <!-- connecting mini line -->
                      <div
                        v-if="dIdx > 0"
                        :class="cn(
                          'w-[6px] h-[1.5px] rounded-full',
                          stageLineClass(task.chain[dIdx - 1].status),
                          stageLineDashed(task.chain[dIdx - 1].status) && 'opacity-40'
                        )"
                      />
                      <Tooltip>
                        <TooltipTrigger as-child>
                          <div
                            :class="cn(
                              'w-[7px] h-[7px] rounded-full border',
                              dotStatus === 'done' ? 'bg-emerald-500 border-emerald-500' :
                              dotStatus === 'in_progress' ? 'bg-blue-500 border-blue-500' :
                              dotStatus === 'failed' ? 'bg-red-500 border-red-500' :
                              dotStatus === 'skipped' ? 'border-zinc-300 dark:border-zinc-600 border-dashed bg-transparent' :
                              'border-zinc-300 dark:border-zinc-600 bg-transparent'
                            )"
                          />
                        </TooltipTrigger>
                        <TooltipContent side="top" class="text-[10px] px-2 py-1">
                          {{ task.chain[dIdx].label }}: {{ task.chain[dIdx].detail }}
                        </TooltipContent>
                      </Tooltip>
                    </template>
                  </div>
                </TooltipProvider>

                <!-- Time -->
                <span class="text-[11px] text-muted-foreground/60 tabular-nums shrink-0 ml-1">
                  {{ task.time }}
                </span>
              </button>

              <!-- Expanded chain detail -->
              <div
                v-if="expandedTask === task.id"
                class="px-4 pb-4 pt-1 bg-muted/20"
              >
                <!-- Horizontal chain timeline -->
                <div class="flex items-start gap-0 py-3 pl-8 overflow-x-auto">
                  <template v-for="(stage, sIdx) in task.chain" :key="sIdx">
                    <!-- Connecting line (before each stage except first) -->
                    <div
                      v-if="sIdx > 0"
                      class="flex items-center pt-[10px] shrink-0"
                    >
                      <div
                        :class="cn(
                          'w-8 h-[2px] rounded-full',
                          stageLineClass(task.chain[sIdx - 1].status),
                        )"
                        :style="stageLineDashed(task.chain[sIdx - 1].status) ? 'background: repeating-linear-gradient(90deg, currentColor 0 3px, transparent 3px 6px)' : ''"
                      />
                    </div>

                    <!-- Stage column -->
                    <div class="flex flex-col items-center shrink-0 min-w-[80px]">
                      <!-- Stage circle -->
                      <div
                        :class="cn(
                          'w-5 h-5 rounded-full border-[1.5px] flex items-center justify-center',
                          stageCircleClass(stage.status)
                        )"
                      >
                        <Check v-if="stage.status === 'done'" :size="11" />
                        <X v-else-if="stage.status === 'failed'" :size="11" />
                        <Minus v-else-if="stage.status === 'skipped'" :size="9" />
                        <component
                          v-else
                          :is="stage.icon"
                          :size="10"
                        />
                      </div>

                      <!-- Label -->
                      <span class="text-[11px] font-medium text-foreground mt-1.5">
                        {{ stage.label }}
                      </span>

                      <!-- Status detail -->
                      <span
                        :class="cn(
                          'text-[10px] mt-0.5 text-center max-w-[100px]',
                          stage.status === 'done' ? 'text-emerald-600 dark:text-emerald-400' :
                          stage.status === 'in_progress' ? 'text-blue-600 dark:text-blue-400' :
                          stage.status === 'failed' ? 'text-red-600 dark:text-red-400' :
                          'text-muted-foreground/60'
                        )"
                      >
                        {{ stage.detail }}
                      </span>

                      <!-- Link -->
                      <button
                        v-if="stage.link"
                        class="flex items-center gap-0.5 text-[10px] text-blue-500 hover:text-blue-600 mt-1"
                      >
                        <ExternalLink :size="9" />
                        {{ stage.link }}
                      </button>
                    </div>
                  </template>
                </div>

                <!-- Blocker info if present -->
                <div
                  v-if="task.blocker"
                  class="ml-8 mt-2 flex items-start gap-2 rounded-md bg-red-500/5 border border-red-500/10 px-3 py-2"
                >
                  <TriangleAlert :size="13" class="text-red-500 mt-0.5 shrink-0" />
                  <div class="flex-1 min-w-0">
                    <p class="text-[11px] text-red-600 dark:text-red-400 font-medium">
                      {{ task.blocker.reason }}
                    </p>
                    <div class="flex items-center gap-2 mt-1.5">
                      <Button variant="outline" size="sm" class="h-6 text-[10px] px-2 gap-1">
                        <ExternalLink :size="9" />
                        {{ task.blocker.actionLabel }}
                      </Button>
                      <Button variant="ghost" size="sm" class="h-6 text-[10px] px-2 gap-1">
                        <Megaphone :size="9" />
                        Ping owner
                      </Button>
                    </div>
                  </div>
                </div>

                <!-- Follow-up reminders if present -->
                <div
                  v-if="task.followUps?.length"
                  class="ml-8 mt-2 space-y-1.5"
                >
                  <div
                    v-for="(fu, fuIdx) in task.followUps"
                    :key="fuIdx"
                    class="flex items-center gap-2 rounded-md bg-amber-500/5 border border-amber-500/10 px-3 py-2"
                  >
                    <component :is="fu.type === 'manual' ? Bell : Zap" :size="12" class="text-amber-500 shrink-0" />
                    <span class="text-[11px] text-amber-700 dark:text-amber-400 flex-1">
                      {{ fu.description }}
                    </span>
                    <Badge variant="outline" class="text-[9px] h-4 px-1.5 bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20">
                      {{ fu.status }}
                    </Badge>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <Separator />

        <!-- ============================================================ -->
        <!-- SECTION 2: Blockers & Dependencies                           -->
        <!-- ============================================================ -->
        <section>
          <div class="flex items-center gap-2 mb-3">
            <AlertTriangle :size="14" class="text-red-500" />
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Blockers &amp; Dependencies
            </h2>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Active Blockers -->
            <div>
              <div class="flex items-center gap-1.5 mb-2">
                <span class="text-[11px] font-semibold text-foreground">Active Blockers</span>
                <Badge variant="outline" class="text-[9px] h-4 px-1 bg-red-500/10 text-red-500 border-red-500/20">
                  {{ blockerCards.length }}
                </Badge>
              </div>
              <div class="space-y-2.5">
                <div
                  v-for="(card, cIdx) in blockerCards"
                  :key="cIdx"
                  class="border border-border rounded-md bg-card border-l-[3px] border-l-red-500 p-3"
                >
                  <div class="flex items-start gap-2">
                    <X :size="13" class="text-red-500 mt-0.5 shrink-0" />
                    <div class="flex-1 min-w-0">
                      <p class="text-[12px] font-medium text-foreground">{{ card.title }}</p>
                      <p class="text-[11px] text-muted-foreground mt-0.5">{{ card.description }}</p>
                      <p class="text-[10px] text-muted-foreground/70 mt-1">{{ card.detail }}</p>

                      <div class="flex items-center gap-2 mt-2">
                        <div v-if="card.owner" class="flex items-center gap-1 text-[10px] text-muted-foreground">
                          <Circle :size="8" class="text-zinc-400" />
                          {{ card.owner }}
                        </div>
                        <div v-if="card.timeStuck" class="flex items-center gap-1 text-[10px] text-muted-foreground">
                          <Clock :size="9" />
                          Stuck {{ card.timeStuck }}
                        </div>
                      </div>

                      <div class="flex items-center gap-1.5 mt-2">
                        <Button
                          v-for="action in card.actions"
                          :key="action"
                          variant="outline"
                          size="sm"
                          class="h-6 text-[10px] px-2 gap-1"
                        >
                          <ExternalLink v-if="action.includes('View') || action.includes('Open')" :size="9" />
                          <Megaphone v-else-if="action.includes('Ping')" :size="9" />
                          <ClipboardList v-else :size="9" />
                          {{ action }}
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Follow-up Items -->
            <div>
              <div class="flex items-center gap-1.5 mb-2">
                <span class="text-[11px] font-semibold text-foreground">Follow-up Items</span>
                <Badge variant="outline" class="text-[9px] h-4 px-1 bg-amber-500/10 text-amber-500 border-amber-500/20">
                  {{ followUpItems.length }}
                </Badge>
              </div>
              <div class="space-y-2.5">
                <div
                  v-for="(fu, fuIdx) in followUpItems"
                  :key="fuIdx"
                  class="border border-border rounded-md bg-amber-500/5 p-3"
                >
                  <div class="flex items-start gap-2">
                    <component
                      :is="fu.type === 'manual' ? Bell : Zap"
                      :size="13"
                      class="text-amber-500 mt-0.5 shrink-0"
                    />
                    <div class="flex-1 min-w-0">
                      <p class="text-[12px] text-foreground">{{ fu.description }}</p>
                      <div class="flex items-center gap-2 mt-1.5">
                        <Badge
                          variant="outline"
                          class="text-[9px] h-4 px-1.5"
                          :class="fu.type === 'manual'
                            ? 'bg-zinc-500/10 text-zinc-500 border-zinc-500/20'
                            : 'bg-violet-500/10 text-violet-500 border-violet-500/20'
                          "
                        >
                          <component :is="fu.type === 'manual' ? Bell : Zap" :size="8" class="mr-0.5" />
                          {{ fu.type === 'manual' ? 'Manual' : 'Auto-generated' }}
                        </Badge>
                        <Badge variant="outline" class="text-[9px] h-4 px-1.5 bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20">
                          {{ fu.status }}
                        </Badge>
                        <span class="text-[10px] text-muted-foreground/60">
                          Source: {{ fu.sourceTask }}
                        </span>
                      </div>
                      <div class="flex items-center gap-1.5 mt-2">
                        <Button variant="ghost" size="sm" class="h-6 text-[10px] px-2 gap-1">
                          <Clock :size="9" />
                          Snooze
                        </Button>
                        <Button variant="ghost" size="sm" class="h-6 text-[10px] px-2 gap-1 text-muted-foreground/60">
                          <X :size="9" />
                          Dismiss
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <Separator />

        <!-- ============================================================ -->
        <!-- SECTION 3: Suggested Actions                                 -->
        <!-- ============================================================ -->
        <section class="pb-4">
          <div class="flex items-center gap-2 mb-3">
            <Zap :size="14" class="text-blue-500" />
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Suggested Actions
            </h2>
          </div>

          <div class="flex flex-col gap-2">
            <div
              v-for="(action, aIdx) in suggestedActions"
              :key="aIdx"
              class="flex items-center gap-3 rounded-md bg-blue-500/5 border border-blue-500/10 border-l-[3px] border-l-blue-500 px-3 py-2.5"
            >
              <component :is="action.icon" :size="14" class="text-blue-500 shrink-0" />
              <span class="text-[12px] text-foreground flex-1 min-w-0">{{ action.description }}</span>
              <Button variant="outline" size="sm" class="h-6 text-[10px] px-2.5 gap-1 shrink-0 border-blue-500/20 text-blue-600 dark:text-blue-400 hover:bg-blue-500/10">
                <ArrowRight :size="9" />
                {{ action.actionLabel }}
              </Button>
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0 shrink-0 text-muted-foreground/40 hover:text-muted-foreground">
                <X :size="12" />
              </Button>
            </div>
          </div>
        </section>

      </div>
    </ScrollArea>
  </div>
</template>
