<script setup lang="ts">
import { onMounted, onActivated, onDeactivated, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { useAuthStore } from '@/stores/auth'
import { usePRStore, branchName, voteIcon } from '@/stores/prs'
import { relativeTime } from '@/lib/date'
import type { PullRequest } from '@/stores/prs'
import { useADOStore } from '@/stores/ado'
import type { ADOPipeline } from '@/stores/ado'
import TaskRow from '@/components/tasks/TaskRow.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/EmptyState.vue'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Badge } from '@/components/ui/badge'
import {
  Plus, GitPullRequest, Play, CheckCircle2, XCircle, Clock,
  ExternalLink, GitBranch, ClipboardList, Loader2,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()
const authStore = useAuthStore()
const prStore = usePRStore()
const adoStore = useADOStore()

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

const recentTasks = computed(() =>
  taskStore.tasks
    .slice()
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    .slice(0, 5)
)

const blockedTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'blocked')
)

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
    <Teleport v-if="isActive" to="#topbar-actions">
      <div class="flex items-center gap-2 text-[10px]">
        <span class="text-muted-foreground">{{ taskStore.stats.total }} tasks</span>
        <Badge v-if="taskStore.stats.inProgress" variant="secondary" class="text-[9px] h-4 px-1.5">{{ taskStore.stats.inProgress }} active</Badge>
        <Badge v-if="taskStore.stats.blocked" variant="destructive" class="text-[9px] h-4 px-1.5">{{ taskStore.stats.blocked }} blocked</Badge>
        <Badge v-if="prStore.reviewPRs.length" variant="outline" class="text-[9px] h-4 px-1.5">{{ prStore.reviewPRs.length }} PRs to review</Badge>
      </div>
    </Teleport>
    <div class="px-6 py-5">
      <!-- Greeting + summary -->
      <div class="mb-6">
        <h1 class="text-xl font-semibold text-foreground">{{ greeting }}</h1>
        <p class="text-sm text-muted-foreground mt-1">{{ summaryLine }}</p>
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
          <!-- Left column (3/5): Focus + Recent + Blocked -->
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
                <TaskRow
                  v-for="task in focusTasks"
                  :key="task.id"
                  :task="task"
                  @select="(id) => { taskStore.selectTask(id); router.push('/tasks') }"
                  @toggle-status="(id) => taskStore.setStatus(id, 'done')"
                />
              </div>
              <p v-else class="text-sm text-muted-foreground">
                No tasks in progress — pick something to work on.
              </p>
            </div>

            <!-- Recent Activity -->
            <div>
              <h2 class="text-sm font-semibold text-foreground mb-3">
                Recent Activity
              </h2>
              <div
                v-if="recentTasks.length > 0"
                class="rounded-lg overflow-hidden border border-border"
              >
                <div
                  v-for="(task, idx) in recentTasks"
                  :key="task.id"
                  class="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
                  :class="idx < recentTasks.length - 1 ? 'border-b border-border' : ''"
                  @click="() => { taskStore.selectTask(task.id); router.push('/tasks') }"
                >
                  <StatusBadge :status="task.status" />
                  <span class="text-sm flex-1 truncate text-foreground">{{ task.title }}</span>
                  <span class="text-xs tabular-nums text-muted-foreground">{{ relativeTime(task.updatedAt) }}</span>
                </div>
              </div>
              <p v-else class="text-sm text-muted-foreground">
                No recent activity yet.
              </p>
            </div>

            <!-- Blocked -->
            <div v-if="blockedTasks.length > 0">
              <h2 class="text-sm font-semibold text-foreground mb-3">
                Blocked
              </h2>
              <div
                class="rounded-lg overflow-hidden border border-border border-l-2 border-l-destructive"
              >
                <TaskRow
                  v-for="task in blockedTasks"
                  :key="task.id"
                  :task="task"
                  @select="(id) => { taskStore.selectTask(id); router.push('/tasks') }"
                />
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
