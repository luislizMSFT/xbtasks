<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { cn } from '@/lib/utils'
import { useADOStore } from '@/stores/ado'
import type { ADOWorkItem, ADOPipeline } from '@/stores/ado'
import { usePRStore, parseReviewers, branchName, voteIcon, relativeTime } from '@/stores/prs'
import type { PullRequest } from '@/stores/prs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Card, CardContent } from '@/components/ui/card'
import {
  AlertCircle,
  Bug,
  BookOpen,
  CheckSquare,
  ChevronRight,
  ChevronDown,
  ExternalLink,
  GitPullRequest,
  CheckCircle2,
  XCircle,
  Loader2,
  RefreshCw,
  Wifi,
  WifiOff,
  Clock,
  Circle,
} from 'lucide-vue-next'

const adoStore = useADOStore()
const prStore = usePRStore()

const closedExpanded = ref(false)
const showAllPipelines = ref(false)
const showAllTeamPRs = ref(false)

onMounted(async () => {
  await Promise.all([
    adoStore.fetchAll(),
    prStore.fetchAll(),
  ])
})

// --- Work item helpers ---

const STATE_ORDER: Record<string, number> = { Active: 0, New: 1, Resolved: 2, Closed: 3 }

const groupedWorkItems = computed(() => {
  const groups: Record<string, ADOWorkItem[]> = { Active: [], New: [], Resolved: [], Closed: [] }
  for (const item of adoStore.workItems) {
    const key = groups[item.state] ? item.state : 'Active'
    groups[key].push(item)
  }
  return Object.entries(groups)
    .filter(([, items]) => items.length > 0)
    .sort(([a], [b]) => (STATE_ORDER[a] ?? 99) - (STATE_ORDER[b] ?? 99))
})

function typeIcon(type: string) {
  switch (type.toLowerCase()) {
    case 'bug': return Bug
    case 'task': return CheckSquare
    case 'user story': return BookOpen
    default: return Circle
  }
}

function typeColor(type: string) {
  switch (type.toLowerCase()) {
    case 'bug': return 'text-red-500'
    case 'task': return 'text-blue-500'
    case 'user story': return 'text-green-500'
    default: return 'text-muted-foreground'
  }
}

function stateClasses(state: string) {
  switch (state) {
    case 'Active': return 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25'
    case 'New': return 'bg-muted text-muted-foreground border-border'
    case 'Resolved': return 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25'
    case 'Closed': return 'bg-zinc-500/15 text-zinc-600 dark:text-zinc-400 border-zinc-500/25'
    default: return ''
  }
}

function priorityClasses(p: number) {
  switch (p) {
    case 1: return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    case 2: return 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/25'
    case 3: return 'bg-yellow-500/15 text-yellow-700 dark:text-yellow-400 border-yellow-500/25'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}

function openItem(item: ADOWorkItem) {
  if (item.url) window.open(item.url, '_blank')
}

// --- PR helpers ---

const activeReviewPRs = computed(() =>
  prStore.reviewPRs.filter(pr => pr.status === 'active')
)

function prStatusClasses(status: string) {
  switch (status) {
    case 'active': return 'bg-blue-500/15 text-blue-600 dark:text-blue-400 border-blue-500/20'
    case 'draft': return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
    case 'completed': return 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
    case 'abandoned': return 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20'
    default: return ''
  }
}

function openPR(pr: PullRequest) {
  if (pr.prUrl) window.open(pr.prUrl, '_blank')
}

function getInitials(name: string) { return name.slice(0, 2).toUpperCase() }

// --- Pipeline helpers ---

function pipelineStatusIcon(pipe: ADOPipeline) {
  const r = pipe.result.toLowerCase()
  const s = pipe.status.toLowerCase()
  if (r === 'succeeded') return { icon: CheckCircle2, color: 'text-green-500', label: 'succeeded' }
  if (r === 'failed') return { icon: XCircle, color: 'text-red-500', label: 'failed' }
  if (s === 'inprogress' || s === 'in_progress') return { icon: Loader2, color: 'text-blue-500 animate-spin', label: 'running' }
  if (s === 'notstarted' || s === 'not_started') return { icon: Clock, color: 'text-amber-500', label: 'queued' }
  return { icon: Circle, color: 'text-muted-foreground', label: s || r }
}

function pipelineBranch(ref: string) {
  return ref.replace('refs/heads/', '')
}

function pipelineDuration(pipe: ADOPipeline) {
  if (!pipe.finishTime || !pipe.queueTime) return ''
  const ms = new Date(pipe.finishTime).getTime() - new Date(pipe.queueTime).getTime()
  if (ms < 0) return ''
  const mins = Math.floor(ms / 60000)
  const secs = Math.floor((ms % 60000) / 1000)
  return `${mins}m ${secs}s`
}

function openPipeline(pipe: ADOPipeline) {
  if (pipe.url) window.open(pipe.url, '_blank')
}

// --- Sync ---

async function handleSync() {
  await adoStore.syncWorkItems()
}

const lastSyncLabel = computed(() => {
  if (!adoStore.lastSyncedAt) return ''
  return relativeTime(adoStore.lastSyncedAt)
})

const pipelinesToShow = computed(() =>
  showAllPipelines.value ? adoStore.pipelines : adoStore.pipelines.slice(0, 5)
)

const teamPRsToShow = computed(() =>
  showAllTeamPRs.value ? prStore.teamPRs : prStore.teamPRs.slice(0, 5)
)
</script>

<template>
  <div class="flex-1 flex flex-col overflow-hidden">
    <!-- Top bar -->
    <div class="px-4 py-2.5 flex items-center gap-3 border-b border-border shrink-0">
      <h1 class="text-sm font-semibold text-foreground">Azure DevOps</h1>
      <Button
        variant="outline"
        size="sm"
        class="h-7 text-xs gap-1.5"
        :disabled="adoStore.syncing"
        @click="handleSync"
      >
        <RefreshCw :size="12" :class="adoStore.syncing && 'animate-spin'" />
        Sync
      </Button>
      <span
        :class="cn(
          'flex items-center gap-1 text-[10px] px-1.5 py-0.5 rounded-full',
          adoStore.connected
            ? 'text-emerald-600 dark:text-emerald-400 bg-emerald-500/10'
            : 'text-muted-foreground bg-muted/50'
        )"
      >
        <Wifi v-if="adoStore.connected" :size="10" />
        <WifiOff v-else :size="10" />
        {{ adoStore.connected ? 'Connected' : 'Cached' }}
      </span>
      <span v-if="lastSyncLabel" class="text-[10px] text-muted-foreground">
        Last synced {{ lastSyncLabel }}
      </span>
    </div>

    <!-- Error banner -->
    <div v-if="adoStore.error && !adoStore.loading" class="mx-4 mt-2 px-3 py-2 rounded-md bg-destructive/10 border border-destructive/20">
      <div class="flex items-center gap-2">
        <AlertCircle :size="14" class="text-destructive shrink-0" />
        <span class="text-xs text-destructive truncate">{{ adoStore.error }}</span>
      </div>
      <p class="text-[10px] text-muted-foreground mt-1">
        Make sure <code class="bg-muted px-1 rounded text-[10px]">az cli</code> is installed and authenticated. Run <code class="bg-muted px-1 rounded text-[10px]">az login</code> then <code class="bg-muted px-1 rounded text-[10px]">az devops configure --defaults organization=https://dev.azure.com/YOUR_ORG project=YOUR_PROJECT</code>
      </p>
    </div>

    <!-- Loading skeleton -->
    <template v-if="adoStore.loading">
      <div class="flex-1 flex items-center justify-center">
        <div class="text-center space-y-2">
          <Loader2 :size="24" class="mx-auto animate-spin text-muted-foreground" />
          <p class="text-xs text-muted-foreground">Loading ADO data…</p>
        </div>
      </div>
    </template>

    <!-- Main two-column layout -->
    <template v-else>
      <div class="flex-1 flex min-h-0">

        <!-- Left column: My Work Items (~60%) -->
        <ScrollArea class="flex-[3] h-full">
          <div class="px-4 py-3 space-y-3">
            <div class="flex items-center gap-1.5">
              <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">My Work Items</span>
              <span class="text-[11px] text-muted-foreground/50 tabular-nums">{{ adoStore.workItems.length }}</span>
            </div>

            <template v-for="[state, items] in groupedWorkItems" :key="state">
              <!-- Closed is collapsible -->
              <div v-if="state === 'Closed'">
                <button
                  class="flex items-center gap-1 text-[10px] font-medium text-muted-foreground/60 uppercase tracking-wider hover:text-muted-foreground transition-colors mb-1"
                  @click="closedExpanded = !closedExpanded"
                >
                  <ChevronRight v-if="!closedExpanded" :size="10" />
                  <ChevronDown v-else :size="10" />
                  {{ state }}
                  <span class="normal-case text-muted-foreground/40 tabular-nums">{{ items.length }}</span>
                </button>
                <Card v-if="closedExpanded" class="overflow-hidden opacity-70">
                  <CardContent class="p-0">
                    <div
                      v-for="item in items"
                      :key="item.id"
                      class="flex items-center gap-2.5 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0"
                      @click="openItem(item)"
                    >
                      <component :is="typeIcon(item.type)" :size="14" :class="['shrink-0', typeColor(item.type)]" />
                      <span class="text-sm text-foreground flex-1 truncate">{{ item.title }}</span>
                      <Badge variant="outline" :class="cn('text-[10px] h-4 px-1.5 shrink-0', priorityClasses(item.priority))">
                        P{{ item.priority }}
                      </Badge>
                      <span class="text-[10px] text-muted-foreground truncate max-w-[10rem]">{{ item.areaPath }}</span>
                      <span class="text-[10px] text-muted-foreground shrink-0">{{ item.assignedTo }}</span>
                      <ExternalLink :size="10" class="shrink-0 text-muted-foreground/30" />
                    </div>
                  </CardContent>
                </Card>
              </div>

              <!-- Non-closed states -->
              <div v-else>
                <div class="flex items-center gap-1.5 mb-1">
                  <Badge variant="outline" :class="cn('text-[10px] h-4 px-1.5', stateClasses(state))">
                    {{ state }}
                  </Badge>
                  <span class="text-[10px] text-muted-foreground/40 tabular-nums">{{ items.length }}</span>
                </div>
                <Card class="overflow-hidden">
                  <CardContent class="p-0">
                    <div
                      v-for="item in items"
                      :key="item.id"
                      class="flex items-center gap-2.5 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0"
                      @click="openItem(item)"
                    >
                      <component :is="typeIcon(item.type)" :size="14" :class="['shrink-0', typeColor(item.type)]" />
                      <span class="text-sm text-foreground flex-1 truncate">{{ item.title }}</span>
                      <Badge variant="outline" :class="cn('text-[10px] h-4 px-1.5 shrink-0', priorityClasses(item.priority))">
                        P{{ item.priority }}
                      </Badge>
                      <span class="text-[10px] text-muted-foreground truncate max-w-[10rem]">{{ item.areaPath }}</span>
                      <span class="text-[10px] text-muted-foreground shrink-0">{{ item.assignedTo }}</span>
                      <ExternalLink :size="10" class="shrink-0 text-muted-foreground/30" />
                    </div>
                  </CardContent>
                </Card>
              </div>
            </template>

            <p v-if="adoStore.workItems.length === 0" class="text-[11px] text-muted-foreground/40 py-2 px-1">No work items found</p>
          </div>
        </ScrollArea>

        <!-- Right column: PRs + Pipelines (~40%) -->
        <ScrollArea class="flex-[2] h-full border-l border-border">
          <div class="px-4 py-3 space-y-3">

            <!-- Pull Requests -->
            <div class="flex items-center gap-1.5">
              <GitPullRequest :size="14" class="text-muted-foreground" />
              <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Pull Requests</span>
            </div>

            <!-- Loading skeleton for PRs -->
            <template v-if="prStore.loading">
              <div class="space-y-2">
                <div v-for="i in 3" :key="i" class="h-10 rounded bg-muted/50 animate-pulse" />
              </div>
            </template>

            <template v-else>
              <!-- Review Requests -->
              <div>
                <span class="text-[10px] font-medium text-muted-foreground/60 uppercase tracking-wider">Needs Your Review</span>
                <Card v-if="activeReviewPRs.length > 0" class="mt-1.5 overflow-hidden">
                  <CardContent class="p-0">
                    <div
                      v-for="pr in activeReviewPRs"
                      :key="pr.id"
                      class="flex flex-col gap-1 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0"
                      @click="openPR(pr)"
                    >
                      <div class="flex items-center gap-2">
                        <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
                        <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize shrink-0', prStatusClasses(pr.status))">
                          {{ pr.status }}
                        </Badge>
                      </div>
                      <div class="flex items-center gap-2 text-[10px] text-muted-foreground/60">
                        <span class="font-medium">{{ pr.repo }}</span>
                        <span>{{ branchName(pr.sourceBranch) }} → {{ branchName(pr.targetBranch) }}</span>
                        <span class="ml-auto flex items-center gap-1">
                          <template v-for="reviewer in parseReviewers(pr.reviewers)" :key="reviewer.uniqueName">
                            <span v-if="reviewer.vote !== 0" :title="`${reviewer.displayName}: ${voteIcon(reviewer.vote)}`">
                              {{ voteIcon(reviewer.vote) }}
                            </span>
                          </template>
                        </span>
                        <span class="tabular-nums shrink-0">{{ relativeTime(pr.updatedAt) }}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
                <p v-else class="text-[11px] text-muted-foreground/40 py-2 px-1">No PRs need your review</p>
              </div>

              <!-- Your PRs -->
              <div>
                <span class="text-[10px] font-medium text-muted-foreground/60 uppercase tracking-wider">Your PRs</span>
                <Card v-if="prStore.myPRs.length > 0" class="mt-1.5 overflow-hidden">
                  <CardContent class="p-0">
                    <div
                      v-for="pr in prStore.myPRs"
                      :key="pr.id"
                      class="flex flex-col gap-1 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0"
                      @click="openPR(pr)"
                    >
                      <div class="flex items-center gap-2">
                        <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
                        <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize shrink-0', prStatusClasses(pr.status))">
                          {{ pr.status }}
                        </Badge>
                      </div>
                      <div class="flex items-center gap-2 text-[10px] text-muted-foreground/60">
                        <span class="font-medium">{{ pr.repo }}</span>
                        <span>{{ branchName(pr.sourceBranch) }} → {{ branchName(pr.targetBranch) }}</span>
                        <div class="ml-auto flex items-center gap-1">
                          <template v-for="reviewer in parseReviewers(pr.reviewers)" :key="reviewer.uniqueName">
                            <span v-if="reviewer.vote !== 0" :title="`${reviewer.displayName}: ${voteIcon(reviewer.vote)}`">
                              {{ voteIcon(reviewer.vote) }}
                            </span>
                            <div
                              v-else
                              class="w-4 h-4 rounded-full bg-muted border border-background flex items-center justify-center"
                              :title="reviewer.displayName"
                            >
                              <span class="text-[7px] font-medium text-muted-foreground">{{ getInitials(reviewer.displayName) }}</span>
                            </div>
                          </template>
                        </div>
                        <span class="tabular-nums shrink-0">{{ relativeTime(pr.updatedAt) }}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
                <p v-else class="text-[11px] text-muted-foreground/40 py-2 px-1">You have no active PRs</p>
              </div>
            </template>

            <!-- Pipelines -->
            <div class="pt-2">
              <div class="flex items-center gap-1.5 mb-1.5">
                <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Pipelines</span>
                <span class="text-[11px] text-muted-foreground/50 tabular-nums">{{ adoStore.pipelines.length }}</span>
              </div>
              <Card v-if="adoStore.pipelines.length > 0" class="overflow-hidden">
                <CardContent class="p-0">
                  <div
                    v-for="pipe in pipelinesToShow"
                    :key="pipe.id"
                    class="px-3 py-2 border-b border-border/50 last:border-b-0 hover:bg-muted/50 cursor-pointer transition-colors"
                    @click="openPipeline(pipe)"
                  >
                    <div class="flex items-center gap-2">
                      <component
                        :is="pipelineStatusIcon(pipe).icon"
                        :size="14"
                        :class="pipelineStatusIcon(pipe).color"
                      />
                      <span class="text-xs font-medium text-foreground truncate">{{ pipe.name }}</span>
                      <Badge
                        variant="outline"
                        :class="cn(
                          'text-[10px] h-4 px-1.5 capitalize shrink-0',
                          pipelineStatusIcon(pipe).label === 'succeeded' && 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25',
                          pipelineStatusIcon(pipe).label === 'running' && 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25',
                          pipelineStatusIcon(pipe).label === 'failed' && 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25',
                          pipelineStatusIcon(pipe).label === 'queued' && 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/25',
                        )"
                      >{{ pipelineStatusIcon(pipe).label }}</Badge>
                      <span v-if="pipelineDuration(pipe)" class="text-[10px] text-muted-foreground ml-auto shrink-0">{{ pipelineDuration(pipe) }}</span>
                    </div>
                    <div class="text-[10px] text-muted-foreground mt-0.5 pl-5 flex items-center gap-2">
                      <span>{{ pipelineBranch(pipe.sourceBranch) }}</span>
                      <span>{{ relativeTime(pipe.queueTime) }}</span>
                      <ExternalLink :size="9" class="shrink-0 opacity-30" />
                    </div>
                  </div>
                </CardContent>
              </Card>
              <p v-else class="text-[11px] text-muted-foreground/40 py-2 px-1">No pipeline runs found</p>
              <Button
                v-if="adoStore.pipelines.length > 5 && !showAllPipelines"
                variant="ghost" size="sm"
                class="w-full text-xs h-7 mt-1"
                @click="showAllPipelines = true"
              >
                Show all {{ adoStore.pipelines.length }} pipelines
              </Button>
            </div>

            <!-- Team Activity -->
            <div class="pt-2">
              <div class="flex items-center gap-1.5 mb-1.5">
                <GitPullRequest :size="14" class="text-muted-foreground" />
                <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Team Activity</span>
                <span class="text-[11px] text-muted-foreground/50 tabular-nums">{{ prStore.teamPRs.length }}</span>
              </div>
              <Card v-if="prStore.teamPRs.length > 0" class="overflow-hidden">
                <CardContent class="p-0">
                  <div
                    v-for="pr in teamPRsToShow"
                    :key="pr.id"
                    class="flex flex-col gap-1 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0"
                    @click="openPR(pr)"
                  >
                    <div class="flex items-center gap-2">
                      <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
                      <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize shrink-0', prStatusClasses(pr.status))">
                        {{ pr.status }}
                      </Badge>
                    </div>
                    <div class="flex items-center gap-2 text-[10px] text-muted-foreground/60">
                      <span class="font-medium">{{ pr.repo }}</span>
                      <span>{{ branchName(pr.sourceBranch) }} → {{ branchName(pr.targetBranch) }}</span>
                      <div class="ml-auto flex items-center gap-1">
                        <template v-for="reviewer in parseReviewers(pr.reviewers)" :key="reviewer.uniqueName">
                          <span v-if="reviewer.vote !== 0" :title="`${reviewer.displayName}: ${voteIcon(reviewer.vote)}`">
                            {{ voteIcon(reviewer.vote) }}
                          </span>
                          <div
                            v-else
                            class="w-4 h-4 rounded-full bg-muted border border-background flex items-center justify-center"
                            :title="reviewer.displayName"
                          >
                            <span class="text-[7px] font-medium text-muted-foreground">{{ getInitials(reviewer.displayName) }}</span>
                          </div>
                        </template>
                      </div>
                      <span class="tabular-nums shrink-0">{{ relativeTime(pr.updatedAt) }}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
              <p v-else class="text-[11px] text-muted-foreground/40 py-2 px-1">No team PRs found</p>
              <Button
                v-if="prStore.teamPRs.length > 5 && !showAllTeamPRs"
                variant="ghost" size="sm"
                class="w-full text-xs h-7 mt-1"
                @click="showAllTeamPRs = true"
              >
                Show all {{ prStore.teamPRs.length }} team PRs
              </Button>
            </div>

          </div>
        </ScrollArea>
      </div>
    </template>
  </div>
</template>
