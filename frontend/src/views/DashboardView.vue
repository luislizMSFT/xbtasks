<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { usePRStore, branchName, voteIcon } from '@/stores/prs'
import { relativeTime } from '@/lib/date'
import type { PullRequest } from '@/stores/prs'
import { useADOStore } from '@/stores/ado'
import type { ADOPipeline } from '@/stores/ado'
import TaskRow from '@/components/TaskRow.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Badge } from '@/components/ui/badge'
import {
  Plus, GitPullRequest, Play, CheckCircle2, XCircle, Clock,
  ExternalLink, GitBranch,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()
const prStore = usePRStore()
const adoStore = useADOStore()

onMounted(() => {
  taskStore.fetchTasks()
  prStore.fetchMyPRs()
  prStore.fetchReviewPRs()
  adoStore.fetchPipelines()
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

const activePRs = computed(() =>
  [...prStore.myPRs, ...prStore.reviewPRs].filter(pr => pr.status === 'active')
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
    const { OpenURL } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/browserservice')
    await OpenURL(pr.prUrl)
  } catch { window.open(pr.prUrl, '_blank') }
}

async function openPipeline(p: ADOPipeline) {
  try {
    const { OpenURL } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/browserservice')
    await OpenURL(p.url)
  } catch { window.open(p.url, '_blank') }
}
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <div class="px-6 py-5">
      <!-- Page header -->
      <div class="flex items-center justify-between mb-6">
        <h1 style="font-size: 20px; font-weight: 600; color: var(--color-text-primary)">Dashboard</h1>
        <Button @click="goCreateTask" class="gap-1.5">
          <Plus :size="14" />
          Create Task
        </Button>
      </div>

      <!-- Empty state -->
      <div v-if="taskStore.tasks.length === 0 && !taskStore.loading && activePRs.length === 0" class="text-center py-20">
        <h2 style="font-size: 20px; font-weight: 600; color: var(--color-text-primary)">Welcome to Team ADO Tool</h2>
        <p style="font-size: 14px; font-weight: 400; color: var(--color-text-secondary)" class="mt-2">
          Your dashboard will show today's focus, PRs, pipelines, and blocked items. Create a task to get started.
        </p>
        <Button @click="goCreateTask" class="mt-4 gap-1.5">
          <Plus :size="14" />
          Create Task
        </Button>
      </div>

      <!-- Dashboard content -->
      <template v-if="taskStore.tasks.length > 0 || activePRs.length > 0">
        <!-- Compact stats line -->
        <div class="flex items-center gap-4 mb-5 text-sm" style="color: var(--color-text-secondary)">
          <span><strong class="tabular-nums" style="color: var(--color-text-primary)">{{ taskStore.stats.inProgress }}</strong> in progress</span>
          <span><strong class="tabular-nums" style="color: var(--color-text-primary)">{{ taskStore.stats.blocked }}</strong> blocked</span>
          <span><strong class="tabular-nums" style="color: var(--color-text-primary)">{{ taskStore.stats.done }}</strong> done</span>
          <span class="text-muted-foreground/50">of {{ taskStore.stats.total }}</span>
        </div>

        <!-- Two-column layout: Tasks left, PRs + Pipelines right -->
        <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
          <!-- Left column (3/5): Focus + Recent + Blocked -->
          <div class="lg:col-span-3 space-y-6">
            <!-- Today's Focus -->
            <div>
              <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)" class="mb-3">
                Today's Focus
              </h2>
              <div
                v-if="focusTasks.length > 0"
                class="rounded-lg overflow-hidden"
                style="border: 1px solid var(--color-border-default)"
              >
                <TaskRow
                  v-for="task in focusTasks"
                  :key="task.id"
                  :task="task"
                  @select="(id) => { taskStore.selectTask(id); router.push('/tasks') }"
                  @toggle-status="(id) => taskStore.setStatus(id, 'done')"
                />
              </div>
              <p v-else style="font-size: 14px; font-weight: 400; color: var(--color-text-tertiary)">
                No tasks in progress — pick something to work on.
              </p>
            </div>

            <!-- Recent Activity -->
            <div>
              <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)" class="mb-3">
                Recent Activity
              </h2>
              <div
                v-if="recentTasks.length > 0"
                class="rounded-lg overflow-hidden"
                style="border: 1px solid var(--color-border-default)"
              >
                <div
                  v-for="(task, idx) in recentTasks"
                  :key="task.id"
                  class="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
                  :style="idx < recentTasks.length - 1 ? 'border-bottom: 1px solid var(--color-border-default)' : ''"
                  @click="() => { taskStore.selectTask(task.id); router.push('/tasks') }"
                >
                  <StatusBadge :status="task.status" />
                  <span class="text-sm flex-1 truncate" style="color: var(--color-text-primary)">{{ task.title }}</span>
                  <span class="text-xs tabular-nums" style="color: var(--color-text-secondary)">{{ relativeTime(task.updatedAt) }}</span>
                </div>
              </div>
            </div>

            <!-- Blocked -->
            <div v-if="blockedTasks.length > 0">
              <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)" class="mb-3">
                Blocked
              </h2>
              <div
                class="rounded-lg overflow-hidden"
                style="border: 1px solid var(--color-border-default); border-left: 2px solid var(--color-status-blocked)"
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
            <!-- Pull Requests -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <GitPullRequest :size="14" class="text-muted-foreground" />
                <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)">
                  Pull Requests
                </h2>
                <Badge v-if="activePRs.length > 0" variant="secondary" class="text-[10px] px-1.5 py-0 h-4">
                  {{ activePRs.length }}
                </Badge>
              </div>
              <div
                v-if="activePRs.length > 0"
                class="rounded-lg overflow-hidden space-y-px"
                style="border: 1px solid var(--color-border-default)"
              >
                <div
                  v-for="pr in activePRs"
                  :key="`${pr.repo}-${pr.id}`"
                  class="flex flex-col gap-1 px-3 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
                  style="border-bottom: 1px solid var(--color-border-default)"
                  @click="openPR(pr)"
                >
                  <div class="flex items-center gap-2">
                    <GitPullRequest :size="13" class="text-green-500 shrink-0" />
                    <span class="text-sm truncate flex-1" style="color: var(--color-text-primary)">{{ pr.title }}</span>
                    <ExternalLink :size="11" class="text-muted-foreground shrink-0 opacity-0 group-hover:opacity-100" />
                  </div>
                  <div class="flex items-center gap-2 pl-5 text-[11px]" style="color: var(--color-text-secondary)">
                    <span class="truncate">{{ pr.repo }}</span>
                    <span>·</span>
                    <GitBranch :size="11" class="shrink-0" />
                    <span class="truncate">{{ branchName(pr.sourceBranch) }}</span>
                    <span class="ml-auto shrink-0 tabular-nums">{{ relativeTime(pr.updatedAt) }}</span>
                  </div>
                </div>
              </div>
              <p v-else style="font-size: 13px; font-weight: 400; color: var(--color-text-tertiary)">
                No active pull requests.
              </p>
            </div>

            <!-- Pipelines (grouped by pipeline name, filtered to user's PR branches) -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <Play :size="14" class="text-muted-foreground" />
                <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)">
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
                  class="rounded-lg overflow-hidden"
                  style="border: 1px solid var(--color-border-default)"
                >
                  <div class="px-3 py-1.5" style="background: var(--color-bg-secondary); border-bottom: 1px solid var(--color-border-default)">
                    <span class="text-xs font-medium" style="color: var(--color-text-primary)">{{ group.pipelineName }}</span>
                  </div>
                  <div
                    v-for="(run, idx) in group.runs"
                    :key="run.id"
                    class="flex items-center gap-2.5 px-3 py-1.5 cursor-pointer hover:bg-muted/50 transition-colors"
                    :style="idx < group.runs.length - 1 ? 'border-bottom: 1px solid var(--color-border-default)' : ''"
                    @click="openPipeline(run)"
                  >
                    <component :is="pipelineIcon(run.result)" :size="13" :class="pipelineColor(run.result)" class="shrink-0" />
                    <span class="text-xs truncate flex-1" style="color: var(--color-text-secondary)">{{ run.sourceBranch }}</span>
                    <Badge
                      :variant="run.result === 'succeeded' ? 'secondary' : run.result === 'failed' ? 'destructive' : 'outline'"
                      class="text-[10px] px-1.5 py-0 h-4 shrink-0"
                    >{{ run.result || run.status }}</Badge>
                    <span class="text-[11px] tabular-nums shrink-0" style="color: var(--color-text-secondary)">
                      {{ run.finishTime ? relativeTime(run.finishTime) : relativeTime(run.queueTime) }}
                    </span>
                  </div>
                </div>
              </div>
              <p v-else style="font-size: 13px; font-weight: 400; color: var(--color-text-tertiary)">
                No pipeline runs for your PRs.
              </p>
            </div>
          </div>
        </div>
      </template>

      <!-- Loading state -->
      <div v-if="taskStore.loading" class="flex items-center justify-center py-20">
        <div class="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
      </div>
    </div>
  </ScrollArea>
</template>
