<script setup lang="ts">
/**
 * LifecycleModal — Dialog showing full lifecycle progress for a single work item.
 * Opened when clicking a row in the traceability table.
 */
import { computed } from 'vue'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  CheckCircle2, Circle, CircleDot, GitPullRequest, GitMerge,
  GitBranch, Rocket, Package, AlertTriangle, Clock, ArrowRight,
  Check, X, Loader2, XCircle,
} from 'lucide-vue-next'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import { statusClasses, prStatusClasses, priorityClasses } from '@/lib/styles'

// ── Types ──

interface LifecycleItem {
  task: {
    id: number; title: string; description: string; status: string; priority: string
    area: string; adoId: string; dueDate: string
    createdAt: string; updatedAt: string; completedAt: string | null
  }
  prs: Array<{
    id: number; title: string; prNumber: number; repo: string; status: string
    sourceBranch: string; targetBranch: string; reviewers: string
    votes: number; createdAt: string; mergedAt: string | null
  }>
  pipelines: Array<{
    id: number; name: string; status: 'running' | 'succeeded' | 'failed' | 'queued' | 'cancelled'
    branch: string; startedAt: string; duration: string | null
  }>
  deployments: Array<{
    id: number; environment: string; status: 'succeeded' | 'failed' | 'in_progress' | 'pending'
    deployedAt: string
  }>
}

const props = defineProps<{
  open: boolean
  item: LifecycleItem | null
}>()

const emit = defineEmits<{ 'update:open': [value: boolean] }>()

// ── Stage progress ──

const stages = computed(() => {
  if (!props.item) return []
  const { task, prs, pipelines, deployments } = props.item

  const hasTask = true
  const hasPR = prs.length > 0
  const hasBuild = pipelines.length > 0
  const hasDeploy = deployments.length > 0

  const prActive = hasPR && prs.some(p => p.status === 'active')
  const buildRunning = hasBuild && pipelines.some(p => p.status === 'running' || p.status === 'queued')
  const deployActive = hasDeploy && deployments.some(d => d.status === 'in_progress' || d.status === 'pending')

  // Determine which stage is the current "active" one
  const activeStage =
    deployActive ? 3
    : hasDeploy ? -1
    : buildRunning ? 2
    : hasBuild ? (hasDeploy ? -1 : 2)
    : prActive ? 1
    : hasPR ? 1
    : 0

  return [
    { label: 'Task', reached: hasTask, active: activeStage === 0 && task.status !== 'done' },
    { label: 'PR', reached: hasPR, active: activeStage === 1 },
    { label: 'Build', reached: hasBuild, active: activeStage === 2 },
    { label: 'Deploy', reached: hasDeploy, active: activeStage === 3 },
  ]
})

// ── Helpers ──

function pipelineStatusColor(status: string): string {
  switch (status) {
    case 'succeeded': return 'text-emerald-500'
    case 'failed': return 'text-red-500'
    case 'running': return 'text-blue-500'
    case 'queued': return 'text-amber-500'
    case 'cancelled': return 'text-zinc-400'
    default: return 'text-zinc-400'
  }
}

function pipelineStatusIcon(status: string) {
  switch (status) {
    case 'succeeded': return Check
    case 'failed': return X
    case 'running': return Loader2
    case 'queued': return Clock
    case 'cancelled': return XCircle
    default: return Circle
  }
}

function deployStatusColor(status: string): string {
  switch (status) {
    case 'succeeded': return 'text-emerald-500'
    case 'failed': return 'text-red-500'
    case 'in_progress': return 'text-blue-500'
    case 'pending': return 'text-amber-500'
    default: return 'text-zinc-400'
  }
}

function voteColor(vote: number): string {
  if (vote >= 10) return 'text-emerald-500'
  if (vote === 5) return 'text-emerald-400'
  if (vote === -5) return 'text-amber-500'
  if (vote <= -10) return 'text-red-500'
  return 'text-zinc-400'
}

function voteIcon(vote: number) {
  if (vote >= 5) return Check
  if (vote <= -5) return X
  return Clock
}

function parseReviewers(json: string): Array<{ displayName: string; uniqueName: string; vote: number }> {
  try { return JSON.parse(json) } catch { return [] }
}

function relativeTime(iso: string): string {
  const ms = Date.now() - new Date(iso).getTime()
  const min = Math.floor(ms / 60000)
  if (min < 1) return 'just now'
  if (min < 60) return `${min}m ago`
  const hrs = Math.floor(min / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function prIcon(pr: { status: string; mergedAt: string | null }) {
  if (pr.mergedAt) return GitMerge
  if (pr.status === 'draft') return GitBranch
  return GitPullRequest
}

function prIconColor(pr: { status: string; mergedAt: string | null }): string {
  if (pr.mergedAt) return 'text-violet-500'
  if (pr.status === 'active') return 'text-emerald-500'
  if (pr.status === 'abandoned') return 'text-red-500'
  return 'text-zinc-400'
}

function initials(name: string): string {
  return name.split(' ').map(p => p[0]).join('').toUpperCase().slice(0, 2)
}
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent class="sm:max-w-2xl p-0 gap-0" :show-close-button="true">
      <template v-if="item">
        <!-- Header -->
        <DialogHeader class="px-5 pt-5 pb-3 border-b border-border space-y-2">
          <DialogTitle class="flex items-center gap-2 text-sm font-semibold leading-tight">
            <CircleDot :size="14" class="text-blue-500 shrink-0" />
            <span class="truncate">{{ item.task.title }}</span>
          </DialogTitle>
          <DialogDescription class="flex flex-wrap items-center gap-1.5">
            <Badge variant="outline" :class="['text-[10px] h-4.5 px-1.5 border', statusClasses(item.task.status)]">
              {{ item.task.status.replace(/_/g, ' ') }}
            </Badge>
            <Badge variant="outline" :class="['text-[10px] h-4.5 px-1.5 border', priorityClasses(item.task.priority)]">
              {{ item.task.priority }}
            </Badge>
            <Badge v-if="item.task.area" variant="secondary" class="text-[10px] h-4.5 px-1.5">
              {{ item.task.area }}
            </Badge>
            <Badge v-if="item.task.adoId" variant="outline" class="text-[10px] h-4.5 px-1.5 gap-1 text-blue-600 dark:text-blue-400 border-blue-500/25">
              <AzureDevOpsIcon :size="10" />
              {{ item.task.adoId }}
            </Badge>
          </DialogDescription>
        </DialogHeader>

        <ScrollArea class="max-h-[70vh]">
          <div class="p-5 space-y-5">
            <!-- Progress stages -->
            <div class="flex items-center justify-center gap-0 py-2">
              <template v-for="(stage, idx) in stages" :key="stage.label">
                <!-- Connector line -->
                <div
                  v-if="idx > 0"
                  class="h-0.5 w-8 sm:w-12"
                  :class="stage.reached ? 'bg-blue-500/60' : 'border-t-2 border-dashed border-zinc-300 dark:border-zinc-600'"
                />
                <!-- Stage circle -->
                <div class="flex flex-col items-center gap-1">
                  <div
                    class="w-8 h-8 rounded-full flex items-center justify-center border-2 transition-colors"
                    :class="[
                      stage.reached
                        ? 'bg-blue-500/15 border-blue-500 text-blue-600 dark:text-blue-400'
                        : 'bg-muted border-zinc-300 dark:border-zinc-600 text-zinc-400',
                      stage.active ? 'ring-2 ring-blue-500/40 ring-offset-1 ring-offset-background' : '',
                    ]"
                  >
                    <CheckCircle2 v-if="stage.reached && !stage.active" :size="14" />
                    <Loader2 v-else-if="stage.active" :size="14" class="animate-spin" />
                    <Circle v-else :size="14" />
                  </div>
                  <span
                    class="text-[10px] font-medium"
                    :class="stage.reached ? 'text-foreground' : 'text-muted-foreground'"
                  >{{ stage.label }}</span>
                </div>
              </template>
            </div>

            <!-- Task Info -->
            <section>
              <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Task Info</h3>
              <div class="rounded-lg border border-border bg-muted/30 p-3 space-y-2">
                <p v-if="item.task.description" class="text-sm text-foreground leading-relaxed">
                  {{ item.task.description }}
                </p>
                <div class="flex flex-wrap gap-x-4 gap-y-1 text-[11px] text-muted-foreground">
                  <span>Created {{ formatDate(item.task.createdAt) }}</span>
                  <span>Updated {{ relativeTime(item.task.updatedAt) }}</span>
                  <span v-if="item.task.completedAt" class="text-emerald-500">
                    Completed {{ formatDate(item.task.completedAt) }}
                  </span>
                  <span v-if="item.task.dueDate">
                    Due {{ formatDate(item.task.dueDate) }}
                  </span>
                </div>
              </div>
            </section>

            <!-- Pull Requests -->
            <section v-if="item.prs.length > 0">
              <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                Pull Requests
                <span class="text-[10px] font-normal ml-1">({{ item.prs.length }})</span>
              </h3>
              <div class="space-y-2">
                <div
                  v-for="pr in item.prs" :key="pr.id"
                  class="rounded-lg border border-border bg-muted/30 p-3"
                >
                  <!-- PR title row -->
                  <div class="flex items-center gap-2 mb-2">
                    <component :is="prIcon(pr)" :size="14" :class="prIconColor(pr)" />
                    <span class="text-sm font-medium truncate">!{{ pr.prNumber }} {{ pr.title }}</span>
                    <div class="flex-1" />
                    <Badge variant="outline" :class="['text-[10px] h-4.5 px-1.5 border shrink-0', prStatusClasses(pr.status)]">
                      {{ pr.status }}
                    </Badge>
                  </div>

                  <!-- Branch info -->
                  <div class="flex items-center gap-1.5 text-[11px] text-muted-foreground mb-2">
                    <GitBranch :size="11" />
                    <span class="font-mono text-[10px]">{{ pr.sourceBranch }}</span>
                    <ArrowRight :size="10" />
                    <span class="font-mono text-[10px]">{{ pr.targetBranch }}</span>
                    <span class="ml-2 text-[10px]">{{ pr.repo }}</span>
                  </div>

                  <!-- Reviewers -->
                  <div v-if="parseReviewers(pr.reviewers).length" class="flex items-center gap-2 flex-wrap">
                    <div
                      v-for="reviewer in parseReviewers(pr.reviewers)"
                      :key="reviewer.uniqueName"
                      class="flex items-center gap-1"
                    >
                      <div class="w-5 h-5 rounded-full bg-zinc-200 dark:bg-zinc-700 flex items-center justify-center text-[8px] font-medium text-foreground">
                        {{ initials(reviewer.displayName) }}
                      </div>
                      <component :is="voteIcon(reviewer.vote)" :size="10" :class="voteColor(reviewer.vote)" />
                    </div>
                    <span class="text-[10px] text-muted-foreground ml-1">{{ relativeTime(pr.createdAt) }}</span>
                  </div>
                </div>
              </div>
            </section>

            <!-- Pipeline Runs -->
            <section v-if="item.pipelines.length > 0">
              <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                Pipeline Runs
                <span class="text-[10px] font-normal ml-1">({{ item.pipelines.length }})</span>
              </h3>
              <div class="space-y-1.5">
                <div
                  v-for="run in item.pipelines" :key="run.id"
                  class="flex items-center gap-2.5 rounded-lg border border-border bg-muted/30 px-3 py-2"
                >
                  <component
                    :is="pipelineStatusIcon(run.status)"
                    :size="14"
                    :class="[pipelineStatusColor(run.status), run.status === 'running' ? 'animate-spin' : '']"
                  />
                  <span class="text-sm font-medium truncate">{{ run.name }}</span>
                  <div class="flex-1" />
                  <span class="text-[10px] text-muted-foreground font-mono">{{ run.branch }}</span>
                  <span v-if="run.duration" class="text-[10px] text-muted-foreground tabular-nums">{{ run.duration }}</span>
                  <span class="text-[10px] text-muted-foreground">{{ relativeTime(run.startedAt) }}</span>
                </div>
              </div>
            </section>

            <!-- Deployments -->
            <section v-if="item.deployments.length > 0">
              <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                Deployments
                <span class="text-[10px] font-normal ml-1">({{ item.deployments.length }})</span>
              </h3>
              <div class="space-y-1.5">
                <div
                  v-for="dep in item.deployments" :key="dep.id"
                  class="flex items-center gap-2.5 rounded-lg border border-border bg-muted/30 px-3 py-2"
                >
                  <Rocket :size="14" :class="deployStatusColor(dep.status)" />
                  <span class="text-sm font-medium capitalize">{{ dep.environment }}</span>
                  <div class="flex-1" />
                  <Badge
                    variant="outline"
                    :class="['text-[10px] h-4.5 px-1.5 border',
                      dep.status === 'succeeded' ? 'text-emerald-500 border-emerald-500/25' :
                      dep.status === 'failed' ? 'text-red-500 border-red-500/25' :
                      dep.status === 'in_progress' ? 'text-blue-500 border-blue-500/25' :
                      'text-zinc-400 border-zinc-400/25'
                    ]"
                  >
                    {{ dep.status.replace(/_/g, ' ') }}
                  </Badge>
                  <span class="text-[10px] text-muted-foreground">{{ formatDate(dep.deployedAt) }}</span>
                </div>
              </div>
            </section>
          </div>
        </ScrollArea>
      </template>
    </DialogContent>
  </Dialog>
</template>
