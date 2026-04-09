<script setup lang="ts">
import { onMounted, onActivated, onDeactivated, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { useAuthStore } from '@/stores/auth'
import { usePRStore, branchName, voteIcon, parseReviewers } from '@/stores/prs'
import { relativeTime } from '@/lib/date'
import type { PullRequest } from '@/stores/prs'
import { useADOStore } from '@/stores/ado'
import { useSyncStore } from '@/stores/sync'
import type { ADOPipeline } from '@/stores/ado'
import DashboardTaskRow from '@/components/tasks/DashboardTaskRow.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/EmptyState.vue'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Badge } from '@/components/ui/badge'
import { priorityBgColor } from '@/lib/styles'
import {
  Plus, GitPullRequest, Play, CheckCircle2, XCircle, Clock,
  ExternalLink, GitBranch, ClipboardList, Loader2,
  CalendarDays, GitMerge,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()
const authStore = useAuthStore()
const prStore = usePRStore()
const adoStore = useADOStore()
const syncStore = useSyncStore()

const isActive = ref(false)
onActivated(() => { isActive.value = true })
onDeactivated(() => { isActive.value = false })
onMounted(() => { isActive.value = true })

const greeting = computed(() => {
  const hour = new Date().getHours()
  const timeOfDay = hour < 12 ? 'morning' : hour < 17 ? 'afternoon' : 'evening'
  const name = authStore.user?.displayName?.split(' ')[0] || ''
  return name ? `Good ${timeOfDay}, ${name}` : `Good ${timeOfDay}`
})

const summaryLine = computed(() => {
  const inProgress = taskStore.tasks.filter(t => t.status === 'in_progress').length
  const blocked = taskStore.tasks.filter(t => t.status === 'blocked').length
  const toReview = prStore.reviewPRs.length
  const parts: string[] = []
  if (inProgress) parts.push(`${inProgress} task${inProgress > 1 ? 's' : ''} in progress`)
  if (blocked) parts.push(`${blocked} blocked`)
  if (toReview) parts.push(`${toReview} PR${toReview > 1 ? 's' : ''} to review`)
  return parts.length ? parts.join(' · ') : 'All clear — no urgent items'
})

onMounted(() => {
  // Data is fetched by App.vue on auth — only fetch if stores are empty
  if (!taskStore.tasks.length) taskStore.fetchTasks()
  if (!prStore.myPRs.length && !prStore.reviewPRs.length) prStore.fetchAll()
  if (!adoStore.pipelines.length) adoStore.fetchPipelines()
})

const focusTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'in_progress' || t.status === 'in_review')
)

// Upcoming: tasks due within 3 days (replaces former recent-activity section)
const upcomingTasks = computed(() =>
  taskStore.tasks.filter(t => {
    if (!t.dueDate || t.status === 'done' || t.status === 'cancelled') return false
    const due = new Date(t.dueDate)
    const now = new Date()
    const diffDays = (due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
    return diffDays <= 3 && diffDays >= -7 // include up to 7 days overdue
  })
    .sort((a, b) => new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime())
)

const blockedTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'blocked')
)

// --- Attention Bar nudge data ---

// Due soon: tasks with due date within 3 days (not done/cancelled)
const dueSoonTasks = computed(() =>
  taskStore.tasks.filter(t => {
    if (!t.dueDate || t.status === 'done' || t.status === 'cancelled') return false
    const diffDays = (new Date(t.dueDate).getTime() - Date.now()) / (1000 * 60 * 60 * 24)
    return diffDays >= 0 && diffDays <= 3
  })
)

// Failed pipeline: first failed pipeline (show one nudge)
const failedPipeline = computed(() =>
  adoStore.pipelines.find(p => p.result === 'failed') ?? null
)

// Merge-ready PRs: active PRs with 2+ reviewer approvals (vote >= 5)
const mergeReadyPRs = computed(() =>
  prStore.myPRs
    .filter(pr => pr.status === 'active')
    .map(pr => {
      const reviewerList = parseReviewers(pr.reviewers)
      const approvals = reviewerList.filter(r => r.vote >= 5).length
      return { ...pr, approvals }
    })
    .filter(pr => pr.approvals >= 2)
)

// Master toggle for Attention Bar visibility
const showAttentionBar = computed(() =>
  dueSoonTasks.value.length > 0 ||
  failedPipeline.value !== null ||
  mergeReadyPRs.value.length > 0
)

function selectTask(taskId: number) {
  taskStore.selectTask(taskId)
  router.push('/tasks')
}

const statCards = computed(() => [
  { label: 'Total', value: taskStore.stats.total },
  { label: 'In Progress', value: taskStore.stats.inProgress },
  { label: 'Blocked', value: taskStore.stats.blocked },
  { label: 'Done', value: taskStore.stats.done },
])

const activeReviewPRs = computed(() =>
  prStore.reviewPRs.filter(pr => pr.status === 'active')
)

const activeMyPRs = computed(() =>
  prStore.myPRs.filter(pr => pr.status === 'active')
)

interface PipelineGroup {
  pipelineName: string
  runs: ADOPipeline[]
  latestResult: string
  latestTime: string
}

// Collect branches from user's PRs
const myPRBranches = computed(() => {
  const branches = new Set<string>()
  for (const pr of [...prStore.myPRs, ...prStore.reviewPRs]) {
    const branch = pr.sourceBranch.replace('refs/heads/', '')
    branches.add(branch)
  }
  return branches
})

// Filter pipelines to only those matching user's PR branches, then group by pipeline name
const groupedPipelines = computed<PipelineGroup[]>(() => {
  const branches = myPRBranches.value
  if (branches.size === 0) return []
  const filtered = adoStore.pipelines.filter(p => branches.has(p.sourceBranch))
  const groups = new Map<string, ADOPipeline[]>()
  for (const p of filtered) {
    const existing = groups.get(p.name) || []
    existing.push(p)
    groups.set(p.name, existing)
  }
  return Array.from(groups.entries()).map(([name, runs]) => {
    const latest = runs[0]
    return {
      pipelineName: name,
      runs: runs.slice(0, 3),
      latestResult: latest.result || latest.status,
      latestTime: latest.finishTime || latest.queueTime,
    }
  })
})

function pipelineIcon(result: string) {
  if (result === 'succeeded') return CheckCircle2
  if (result === 'failed') return XCircle
  return Play
}

function pipelineColor(result: string) {
  if (result === 'succeeded') return 'text-green-500'
  if (result === 'failed') return 'text-red-500'
  return 'text-yellow-500'
}

function goCreateTask() {
  router.push({ path: '/tasks', query: { create: '1' } })
}

async function openPR(pr: PullRequest) {
  try {
    const { openURL } = await import('@/api/browser')
    await openURL(pr.prUrl)
  } catch { window.open(pr.prUrl, '_blank') }
}

async function openPipeline(p: ADOPipeline) {
  try {
    const { openURL } = await import('@/api/browser')
    await openURL(p.url)
  } catch { window.open(p.url, '_blank') }
}
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <!-- Dashboard top bar: summary stats -->
    <Teleport v-if="isActive" to="#topbar-center">
      <div class="flex items-center gap-2 text-[10px]">
        <Badge variant="secondary" class="text-[9px] h-4 px-1.5">
          {{ taskStore.stats.inProgress + taskStore.stats.inReview }} active
        </Badge>
        <Badge v-if="taskStore.stats.blocked" variant="destructive" class="text-[9px] h-4 px-1.5">
          {{ taskStore.stats.blocked }} blocked
        </Badge>
        <span class="text-muted-foreground/40 tabular-nums">
          {{ taskStore.stats.done }}/{{ taskStore.stats.total }} done
        </span>
      </div>
    </Teleport>
    <div class="px-6 py-5">
      <!-- Greeting + summary -->
      <div class="mb-3">
        <h1 class="text-lg font-semibold text-foreground">{{ greeting }}</h1>
        <p v-if="summaryLine" class="text-sm text-muted-foreground mt-1">{{ summaryLine }}</p>
      </div>

      <!-- Attention Bar (conditional — shown only when urgency nudges exist) -->
      <div v-if="showAttentionBar" class="flex gap-3 mb-5 overflow-x-auto pb-1">
        <!-- Due soon nudge (amber) -->
        <div v-if="dueSoonTasks.length > 0"
          class="flex items-center gap-2 px-3 py-2 rounded-lg border shrink-0"
          style="background: rgb(245 158 11 / 0.08); border-color: rgb(245 158 11 / 0.2);">
          <CalendarDays :size="14" class="text-amber-600 dark:text-amber-400 shrink-0" />
          <span class="text-sm text-amber-700 dark:text-amber-400">
            <strong class="font-semibold">{{ dueSoonTasks.length }}</strong> due within 3 days
          </span>
          <div class="flex gap-1 ml-1">
            <span v-for="t in dueSoonTasks.slice(0, 2)" :key="t.id"
              class="text-[10px] px-1.5 py-0.5 rounded truncate max-w-[120px] text-amber-700 dark:text-amber-400"
              style="background: rgb(245 158 11 / 0.1);">
              {{ t.title }}
            </span>
          </div>
        </div>

        <!-- Pipeline failure nudge (red) -->
        <div v-if="failedPipeline"
          class="flex items-center gap-2 px-3 py-2 rounded-lg border shrink-0"
          style="background: rgb(239 68 68 / 0.08); border-color: rgb(239 68 68 / 0.2);">
          <XCircle :size="14" class="text-red-500 shrink-0" />
          <span class="text-sm text-red-600 dark:text-red-400">
            <strong class="font-semibold">{{ failedPipeline.name }}</strong> failed on {{ failedPipeline.sourceBranch.replace('refs/heads/', '') }}
          </span>
        </div>

        <!-- PR approval ready nudge (emerald) -->
        <div v-for="pr in mergeReadyPRs.slice(0, 2)" :key="`merge-${pr.id}`"
          class="flex items-center gap-2 px-3 py-2 rounded-lg border shrink-0"
          style="background: rgb(16 185 129 / 0.08); border-color: rgb(16 185 129 / 0.2);">
          <GitMerge :size="14" class="text-emerald-600 dark:text-emerald-400 shrink-0" />
          <span class="text-sm text-emerald-700 dark:text-emerald-400">
            PR #{{ pr.prNumber }} has {{ pr.approvals }} approvals — ready to merge
          </span>
        </div>
      </div>

      <!-- Loading state -->
      <template v-if="taskStore.loading && taskStore.tasks.length === 0">
        <div class="flex items-center gap-2 text-sm text-muted-foreground py-8">
          <Loader2 :size="16" class="animate-spin" />
          Loading your dashboard...
        </div>
      </template>

      <!-- Dashboard content (reactive — fills as data loads) -->
      <template v-else-if="taskStore.tasks.length === 0 && activeMyPRs.length === 0 && activeReviewPRs.length === 0 && !taskStore.loading">
        <EmptyState
          :icon="ClipboardList"
          title="Welcome to XB Tasks"
          description="Your dashboard will show today's focus, PRs, pipelines, and blocked items. Create a task to get started."
        >
          <template #action>
            <Button @click="goCreateTask" class="gap-1.5">
              <Plus :size="14" />
              Create Task
            </Button>
          </template>
        </EmptyState>
      </template>

      <template v-else>
        <!-- Compact stats line -->
        <div class="flex items-center gap-4 mb-5 text-sm text-muted-foreground">
          <span><strong class="tabular-nums text-foreground">{{ taskStore.stats.inProgress }}</strong> in progress</span>
          <span><strong class="tabular-nums text-foreground">{{ taskStore.stats.blocked }}</strong> blocked</span>
          <span><strong class="tabular-nums text-foreground">{{ taskStore.stats.done }}</strong> done</span>
          <span class="text-muted-foreground/50">of {{ taskStore.stats.total }}</span>
        </div>

        <!-- Two-column layout: Tasks left, PRs + Pipelines right -->
        <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
          <!-- Left column (3/5): Focus + Upcoming + Blocked -->
          <div class="lg:col-span-3 space-y-6">
            <!-- Today's Focus -->
            <div>
              <h2 class="text-sm font-semibold text-foreground mb-3">
                Today's Focus
              </h2>
              <div
                v-if="focusTasks.length > 0"
                class="rounded-lg overflow-hidden border border-border"
              >
                <DashboardTaskRow
                  v-for="(task, idx) in focusTasks"
                  :key="task.id"
                  :task="task"
                  :is-personal="!taskStore.isPublic(task.id)"
                  :class="idx < focusTasks.length - 1 ? 'border-b border-border' : ''"
                  @click="selectTask"
                />
              </div>
              <p v-else class="text-sm text-muted-foreground">
                All caught up
              </p>
            </div>

            <!-- Upcoming -->
            <div>
              <h2 class="text-sm font-semibold text-foreground mb-3">
                Upcoming
              </h2>
              <div
                v-if="upcomingTasks.length > 0"
                class="rounded-lg overflow-hidden border border-border"
              >
                <DashboardTaskRow
                  v-for="(task, idx) in upcomingTasks"
                  :key="task.id"
                  :task="task"
                  :is-personal="!taskStore.isPublic(task.id)"
                  show-due-date
                  :class="idx < upcomingTasks.length - 1 ? 'border-b border-border' : ''"
                  @click="selectTask"
                />
              </div>
              <p v-else class="text-sm text-muted-foreground">
                No tasks due soon
              </p>
            </div>

            <!-- Blocked -->
            <div v-if="blockedTasks.length > 0">
              <h2 class="text-sm font-semibold text-foreground mb-3">
                Blocked
              </h2>
              <div
                class="rounded-lg overflow-hidden border border-border border-l-2 border-l-red-500"
              >
                <div
                  v-for="(task, idx) in blockedTasks"
                  :key="task.id"
                  class="flex flex-col gap-1 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
                  :class="idx < blockedTasks.length - 1 ? 'border-b border-border' : ''"
                  @click="selectTask(task.id)"
                >
                  <div class="flex items-center gap-2">
                    <span :class="['size-2 rounded-full shrink-0', priorityBgColor(task.priority)]" />
                    <span class="text-sm flex-1 truncate text-foreground">{{ task.title }}</span>
                  </div>
                  <p v-if="task.blockedReason"
                    class="text-[11px] text-red-500/70 pl-4 italic">
                    {{ task.blockedReason }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Right column (2/5): PRs + Pipelines -->
          <div class="lg:col-span-2 space-y-6">
            <!-- Needs Your Review -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <GitPullRequest :size="14" class="text-muted-foreground" />
                <h2 class="text-sm font-semibold text-foreground">
                  Needs Your Review
                </h2>
                <Badge v-if="activeReviewPRs.length > 0" variant="secondary" class="text-[10px] px-1.5 py-0 h-4">
                  {{ activeReviewPRs.length }}
                </Badge>
              </div>
              <div
                v-if="activeReviewPRs.length > 0"
                class="rounded-lg overflow-hidden space-y-px border border-border"
              >
                <div
                  v-for="pr in activeReviewPRs"
                  :key="`review-${pr.repo}-${pr.id}`"
                  class="flex flex-col gap-1 px-3 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors border-b border-border"
                  @click="openPR(pr)"
                >
                  <div class="flex items-center gap-2">
                    <GitPullRequest :size="13" class="text-green-500 shrink-0" />
                    <span class="text-sm truncate flex-1 text-foreground">{{ pr.title }}</span>
                    <span v-if="pr.votes" class="shrink-0 text-xs">{{ voteIcon(pr.votes) }}</span>
                    <ExternalLink :size="11" class="text-muted-foreground shrink-0 opacity-0 group-hover:opacity-100" />
                  </div>
                  <div class="flex items-center gap-2 pl-5 text-[11px] text-muted-foreground">
                    <span class="truncate">{{ pr.repo }}</span>
                    <span>·</span>
                    <GitBranch :size="11" class="shrink-0" />
                    <span class="truncate">{{ branchName(pr.sourceBranch) }}</span>
                    <span class="ml-auto shrink-0 tabular-nums">{{ relativeTime(pr.updatedAt) }}</span>
                  </div>
                </div>
              </div>
              <p v-else class="text-[13px] text-muted-foreground">
                No PRs to review.
              </p>
            </div>

            <!-- Your Pull Requests -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <GitPullRequest :size="14" class="text-muted-foreground" />
                <h2 class="text-sm font-semibold text-foreground">
                  Your Pull Requests
                </h2>
                <Badge v-if="activeMyPRs.length > 0" variant="secondary" class="text-[10px] px-1.5 py-0 h-4">
                  {{ activeMyPRs.length }}
                </Badge>
              </div>
              <div
                v-if="activeMyPRs.length > 0"
                class="rounded-lg overflow-hidden space-y-px border border-border"
              >
                <div
                  v-for="pr in activeMyPRs"
                  :key="`my-${pr.repo}-${pr.id}`"
                  class="flex flex-col gap-1 px-3 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors border-b border-border"
                  @click="openPR(pr)"
                >
                  <div class="flex items-center gap-2">
                    <GitPullRequest :size="13" class="text-green-500 shrink-0" />
                    <span class="text-sm truncate flex-1 text-foreground">{{ pr.title }}</span>
                    <span v-if="pr.votes" class="shrink-0 text-xs">{{ voteIcon(pr.votes) }}</span>
                    <ExternalLink :size="11" class="text-muted-foreground shrink-0 opacity-0 group-hover:opacity-100" />
                  </div>
                  <div class="flex items-center gap-2 pl-5 text-[11px] text-muted-foreground">
                    <span class="truncate">{{ pr.repo }}</span>
                    <span>·</span>
                    <GitBranch :size="11" class="shrink-0" />
                    <span class="truncate">{{ branchName(pr.sourceBranch) }}</span>
                    <span class="ml-auto shrink-0 tabular-nums">{{ relativeTime(pr.updatedAt) }}</span>
                  </div>
                </div>
              </div>
              <p v-else class="text-[13px] text-muted-foreground">
                No active PRs.
              </p>
            </div>

            <!-- Pipelines (grouped by pipeline name, filtered to user's PR branches) -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <Play :size="14" class="text-muted-foreground" />
                <h2 class="text-sm font-semibold text-foreground">
                  Pipelines
                </h2>
                <Badge v-if="groupedPipelines.length > 0" variant="secondary" class="text-[10px] px-1.5 py-0 h-4">
                  {{ groupedPipelines.length }}
                </Badge>
              </div>
              <div
                v-if="groupedPipelines.length > 0"
                class="space-y-2"
              >
                <div
                  v-for="group in groupedPipelines"
                  :key="group.pipelineName"
                  class="rounded-lg overflow-hidden border border-border"
                >
                  <div class="px-3 py-1.5 bg-muted border-b border-border">
                    <span class="text-xs font-medium text-foreground">{{ group.pipelineName }}</span>
                  </div>
                  <div
                    v-for="(run, idx) in group.runs"
                    :key="run.id"
                    class="flex items-center gap-2.5 px-3 py-1.5 cursor-pointer hover:bg-muted/50 transition-colors"
                    :class="idx < group.runs.length - 1 ? 'border-b border-border' : ''"
                    @click="openPipeline(run)"
                  >
                    <component :is="pipelineIcon(run.result)" :size="13" :class="pipelineColor(run.result)" class="shrink-0" />
                    <span class="text-xs truncate flex-1 text-muted-foreground">{{ run.sourceBranch }}</span>
                    <Badge
                      :variant="run.result === 'succeeded' ? 'secondary' : run.result === 'failed' ? 'destructive' : 'outline'"
                      class="text-[10px] px-1.5 py-0 h-4 shrink-0"
                    >{{ run.result || run.status }}</Badge>
                    <span class="text-[11px] tabular-nums shrink-0 text-muted-foreground">
                      {{ run.finishTime ? relativeTime(run.finishTime) : relativeTime(run.queueTime) }}
                    </span>
                  </div>
                </div>
              </div>
              <p v-else class="text-[13px] text-muted-foreground">
                No pipeline runs for your PRs.
              </p>
            </div>
          </div>
        </div>
      </template>


    </div>
  </ScrollArea>
</template>
