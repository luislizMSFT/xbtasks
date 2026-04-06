<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { AcceptableValue } from 'reka-ui'
import { cn } from '@/lib/utils'
import { useADOStore } from '@/stores/ado'
import type { ADOWorkItem, ADOPipeline } from '@/stores/ado'
import { useTaskStore } from '@/stores/tasks'
import { usePRStore, parseReviewers, branchName, voteIcon, relativeTime } from '@/stores/prs'
import type { PullRequest } from '@/stores/prs'
import PageHeader from '@/components/PageHeader.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'
import AdoTreeBrowser from '@/components/AdoTreeBrowser.vue'
import LinkDialog from '@/components/LinkDialog.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Card, CardContent } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import {
  AlertCircle,
  Bug,
  BookOpen,
  CheckSquare,
  Star,
  Mountain,
  Circle,
  ExternalLink,
  GitPullRequest,
  CheckCircle2,
  XCircle,
  Loader2,
  RefreshCw,
  Wifi,
  WifiOff,
  Clock,
  GitMerge,
  Search,
  Download,
  Link,
  EyeOff,
  Eye,
  FolderKanban,
} from 'lucide-vue-next'

const adoStore = useADOStore()
const taskStore = useTaskStore()
const prStore = usePRStore()

const showAllPipelines = ref(false)
const showAllTeamPRs = ref(false)

// Link dialog state
const linkDialogOpen = ref(false)

// Track the ADO item being linked (when linking from browser to existing task)
const linkingAdoItem = ref<ADOWorkItem | null>(null)

// Import choice dialog state
const importChoiceOpen = ref(false)
const importingAdoItem = ref<ADOWorkItem | null>(null)

onMounted(async () => {
  await Promise.all([
    adoStore.fetchWorkItemTree(),
    adoStore.fetchLinkedAdoIds(),
    adoStore.fetchSavedQueries(),
    prStore.fetchAll(),
  ])
})

// --- Type/state helpers ---

function typeIcon(type: string) {
  switch (type.toLowerCase()) {
    case 'bug': return Bug
    case 'task': return CheckSquare
    case 'user story': return BookOpen
    case 'feature': return Star
    case 'epic': return Mountain
    default: return Circle
  }
}

function typeColor(type: string) {
  switch (type.toLowerCase()) {
    case 'bug': return 'text-red-500'
    case 'task': return 'text-blue-500'
    case 'user story': return 'text-purple-500'
    case 'feature': return 'text-green-500'
    case 'epic': return 'text-orange-500'
    default: return 'text-muted-foreground'
  }
}

function stateClasses(state: string) {
  switch (state) {
    case 'Active': return 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25'
    case 'New': return 'bg-muted text-muted-foreground border-border'
    case 'Resolved': return 'bg-yellow-500/15 text-yellow-700 dark:text-yellow-400 border-yellow-500/25'
    case 'Closed': return 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25'
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

// --- Tree browser actions ---

function onTreeSelect(item: ADOWorkItem) {
  adoStore.selectedItem = item
}

async function onTreeImport(item: ADOWorkItem) {
  importingAdoItem.value = item
  importChoiceOpen.value = true
}

async function doImportAsTask() {
  if (!importingAdoItem.value) return
  try {
    const { ImportWorkItem } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
    await ImportWorkItem(importingAdoItem.value.adoId)
    await Promise.all([adoStore.fetchLinkedAdoIds(), taskStore.fetchTasks()])
  } catch (e: any) {
    adoStore.error = e?.message || 'Import failed'
  } finally {
    importChoiceOpen.value = false
    importingAdoItem.value = null
  }
}

async function doImportAsProject() {
  if (!importingAdoItem.value) return
  try {
    const { ImportWorkItemAsProject } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
    await ImportWorkItemAsProject(importingAdoItem.value.adoId)
    await Promise.all([adoStore.fetchLinkedAdoIds(), taskStore.fetchTasks()])
  } catch (e: any) {
    adoStore.error = e?.message || 'Import as project failed'
  } finally {
    importChoiceOpen.value = false
    importingAdoItem.value = null
  }
}

function onTreeLink(item: ADOWorkItem) {
  linkingAdoItem.value = item
  linkDialogOpen.value = true
}

function onLinked() {
  linkDialogOpen.value = false
  linkingAdoItem.value = null
  Promise.all([adoStore.fetchLinkedAdoIds(), taskStore.fetchTasks()])
}

// --- Saved query picker ---
const selectedQueryId = ref('__my_assignments__')

async function onQueryChange(queryId: AcceptableValue) {
  if (typeof queryId !== 'string') return
  selectedQueryId.value = queryId
  if (queryId === '__my_assignments__') {
    await adoStore.fetchWorkItemTree()
  } else {
    await adoStore.runSavedQuery(queryId)
  }
}

async function openExternal(url: string) {
  try {
    const { OpenURL } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/browserservice')
    await OpenURL(url)
  } catch {
    // Fallback if binding not available
    window.open(url, '_blank')
  }
}

// --- PR helpers ---

const activeReviewPRs = computed(() =>
  prStore.reviewPRs.filter(pr => pr.status === 'active')
)

const totalPRCount = computed(() =>
  activeReviewPRs.value.length + prStore.myPRs.length + prStore.teamPRs.length
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
  if (pr.prUrl) openExternal(pr.prUrl)
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
  if (pipe.url) openExternal(pipe.url)
}

// --- Sync ---

async function handleSync() {
  await adoStore.syncWorkItems()
  await adoStore.fetchWorkItemTree()
  await adoStore.fetchLinkedAdoIds()
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
        {{ adoStore.connected ? 'Connected' : 'Offline' }}
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
        Run <code class="bg-muted px-1 rounded text-[10px]">az login</code> then <code class="bg-muted px-1 rounded text-[10px]">az devops configure --defaults organization=https://dev.azure.com/YOUR_ORG project=YOUR_PROJECT</code>
      </p>
    </div>

    <!-- Loading -->
    <template v-if="adoStore.loading && !adoStore.workItemTree.length">
      <div class="flex-1 flex items-center justify-center">
        <LoadingSpinner label="Loading ADO items..." />
      </div>
    </template>

    <!-- Tabbed content -->
    <template v-else>
      <Tabs default-value="browser" class="flex-1 flex flex-col min-h-0">
        <div class="px-4 pt-2 shrink-0">
          <TabsList class="h-8">
            <TabsTrigger value="browser" class="text-xs gap-1.5 px-3">
              <CheckSquare :size="12" />
              ADO Browser
              <Badge v-if="adoStore.workItemTree.length" variant="secondary" class="text-[10px] h-4 px-1 ml-0.5">{{ adoStore.workItemTree.length }}</Badge>
            </TabsTrigger>
            <TabsTrigger value="pull-requests" class="text-xs gap-1.5 px-3">
              <GitPullRequest :size="12" />
              Pull Requests
              <Badge v-if="totalPRCount" variant="secondary" class="text-[10px] h-4 px-1 ml-0.5">{{ totalPRCount }}</Badge>
            </TabsTrigger>
            <TabsTrigger value="pipelines" class="text-xs gap-1.5 px-3">
              <GitMerge :size="12" />
              Pipelines
              <Badge v-if="adoStore.pipelines.length" variant="secondary" class="text-[10px] h-4 px-1 ml-0.5">{{ adoStore.pipelines.length }}</Badge>
            </TabsTrigger>
          </TabsList>
        </div>

        <!-- ═══ ADO Browser Tab ═══ -->
        <TabsContent value="browser" class="flex-1 min-h-0 mt-0 flex flex-col">
          <!-- Filter toolbar -->
          <PageHeader>
            <template #left>
              <!-- Search -->
              <div class="relative">
                <Search :size="12" class="absolute left-2 top-2 text-muted-foreground" />
                <Input
                  v-model="adoStore.searchQuery"
                  placeholder="Search..."
                  class="h-7 w-40 pl-7 text-[11px]"
                />
              </div>

              <!-- Type filter -->
              <Select v-model="adoStore.filterType">
                <SelectTrigger class="h-7 w-24 text-[11px]">
                  <SelectValue placeholder="Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all" class="text-xs">All Types</SelectItem>
                  <SelectItem value="Task" class="text-xs">Task</SelectItem>
                  <SelectItem value="Bug" class="text-xs">Bug</SelectItem>
                  <SelectItem value="User Story" class="text-xs">User Story</SelectItem>
                  <SelectItem value="Feature" class="text-xs">Feature</SelectItem>
                  <SelectItem value="Epic" class="text-xs">Epic</SelectItem>
                </SelectContent>
              </Select>

              <!-- State filter -->
              <Select v-model="adoStore.filterState">
                <SelectTrigger class="h-7 w-24 text-[11px]">
                  <SelectValue placeholder="State" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all" class="text-xs">All States</SelectItem>
                  <SelectItem value="Active" class="text-xs">Active</SelectItem>
                  <SelectItem value="New" class="text-xs">New</SelectItem>
                  <SelectItem value="Resolved" class="text-xs">Resolved</SelectItem>
                  <SelectItem value="Closed" class="text-xs">Closed</SelectItem>
                </SelectContent>
              </Select>

              <!-- Area Path filter -->
              <Select v-model="adoStore.filterArea">
                <SelectTrigger class="h-7 w-32 text-[11px]">
                  <SelectValue placeholder="Area Path" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all" class="text-xs">All Areas</SelectItem>
                  <SelectItem
                    v-for="area in adoStore.availableAreaPaths"
                    :key="area"
                    :value="area"
                    class="text-xs"
                  >{{ area }}</SelectItem>
                </SelectContent>
              </Select>

              <!-- Saved Queries -->
              <Select :model-value="selectedQueryId" @update:model-value="onQueryChange">
                <SelectTrigger class="h-7 w-36 text-[11px]">
                  <SelectValue placeholder="Saved Queries" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__my_assignments__" class="text-xs">My Assignments</SelectItem>
                  <SelectItem
                    v-for="q in adoStore.savedQueries"
                    :key="q.id"
                    :value="q.id"
                    class="text-xs"
                  >{{ q.name }}</SelectItem>
                </SelectContent>
              </Select>
            </template>

            <template #right>
              <!-- Hide Linked toggle -->
              <Button
                :variant="adoStore.hideLinked ? 'default' : 'outline'"
                size="sm"
                class="h-7 text-[11px] gap-1"
                @click="adoStore.hideLinked = !adoStore.hideLinked"
              >
                <EyeOff v-if="adoStore.hideLinked" :size="12" />
                <Eye v-else :size="12" />
                Hide Linked
              </Button>
            </template>
          </PageHeader>

          <!-- Two-column layout: tree + detail -->
          <div class="flex-1 flex min-h-0">
            <!-- Left: Tree browser (60%) -->
            <div class="w-[60%] border-r border-border">
              <ScrollArea class="h-full">
                <div class="py-2 px-2">
                  <EmptyState
                    v-if="!adoStore.loading && adoStore.treeRoots.length === 0"
                    :icon="CheckSquare"
                    title="No work items"
                    description="No ADO work items match your current filters. Try adjusting your filters or syncing."
                    class="py-8"
                  />
                  <AdoTreeBrowser
                    v-else
                    :items="adoStore.treeRoots"
                    :get-children="adoStore.getChildren"
                    :is-linked="adoStore.isLinked"
                    @select="onTreeSelect"
                    @import="onTreeImport"
                    @link="onTreeLink"
                  />
                </div>
              </ScrollArea>
            </div>

            <!-- Right: Detail panel (40%) -->
            <div class="w-[40%]">
              <ScrollArea class="h-full">
                <div v-if="adoStore.selectedItem" class="p-4 space-y-4">
                  <!-- Title + type badge -->
                  <div class="space-y-2">
                    <div class="flex items-center gap-2">
                      <component :is="typeIcon(adoStore.selectedItem.type)" :size="16" :class="typeColor(adoStore.selectedItem.type)" />
                      <Badge variant="outline" :class="['text-[10px] h-4 px-1.5', stateClasses(adoStore.selectedItem.state)]">
                        {{ adoStore.selectedItem.state }}
                      </Badge>
                      <span class="text-[10px] text-muted-foreground tabular-nums">#{{ adoStore.selectedItem.adoId }}</span>
                    </div>
                    <h2 class="text-sm font-semibold text-foreground leading-snug">{{ adoStore.selectedItem.title }}</h2>
                  </div>

                  <!-- Metadata -->
                  <div class="space-y-1.5">
                    <div v-if="adoStore.selectedItem.assignedTo" class="flex items-center gap-2 text-xs">
                      <span class="text-muted-foreground w-20 shrink-0">Assigned to</span>
                      <span class="text-foreground">{{ adoStore.selectedItem.assignedTo }}</span>
                    </div>
                    <div class="flex items-center gap-2 text-xs">
                      <span class="text-muted-foreground w-20 shrink-0">Priority</span>
                      <Badge variant="outline" :class="['text-[10px] h-4 px-1.5', priorityClasses(adoStore.selectedItem.priority)]">
                        P{{ adoStore.selectedItem.priority }}
                      </Badge>
                    </div>
                    <div v-if="adoStore.selectedItem.areaPath" class="flex items-center gap-2 text-xs">
                      <span class="text-muted-foreground w-20 shrink-0">Area Path</span>
                      <span class="text-foreground truncate">{{ adoStore.selectedItem.areaPath }}</span>
                    </div>
                    <div v-if="adoStore.selectedItem.org || adoStore.selectedItem.project" class="flex items-center gap-2 text-xs">
                      <span class="text-muted-foreground w-20 shrink-0">Org/Project</span>
                      <span class="text-foreground">{{ adoStore.selectedItem.org }}/{{ adoStore.selectedItem.project }}</span>
                    </div>
                  </div>

                  <!-- Description -->
                  <div v-if="adoStore.selectedItem.description" class="space-y-1">
                    <h3 class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
                    <div
                      class="rounded-md border bg-muted/20 px-3 py-2 text-xs text-foreground prose prose-sm max-h-[200px] overflow-y-auto"
                      v-html="adoStore.selectedItem.description"
                    />
                  </div>

                  <!-- Action buttons -->
                  <div class="flex items-center gap-2 pt-2 border-t border-border">
                    <Button
                      v-if="!adoStore.isLinked(adoStore.selectedItem.adoId)"
                      size="sm"
                      class="h-7 text-xs gap-1.5"
                      @click="onTreeImport(adoStore.selectedItem)"
                    >
                      <Download :size="12" />
                      Import as Task
                    </Button>
                    <Button
                      v-if="!adoStore.isLinked(adoStore.selectedItem.adoId)"
                      variant="outline"
                      size="sm"
                      class="h-7 text-xs gap-1.5"
                      @click="onTreeLink(adoStore.selectedItem)"
                    >
                      <Link :size="12" />
                      Link to Existing Task
                    </Button>
                    <Badge v-if="adoStore.isLinked(adoStore.selectedItem.adoId)" class="bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/25 text-[10px]">
                      <CheckCircle2 :size="10" class="mr-1" />
                      Linked
                    </Badge>
                    <Button
                      variant="outline"
                      size="sm"
                      class="h-7 text-xs gap-1.5 ml-auto"
                      @click="openExternal(adoStore.selectedItem!.url)"
                    >
                      <ExternalLink :size="12" />
                      Open in ADO
                    </Button>
                  </div>
                </div>

                <!-- Empty state -->
                <div v-else class="flex-1 flex items-center justify-center p-8">
                  <p class="text-[11px] text-muted-foreground/40 text-center">Select a work item to view details</p>
                </div>
              </ScrollArea>
            </div>
          </div>
        </TabsContent>

        <!-- ═══ Pull Requests Tab ═══ -->
        <TabsContent value="pull-requests" class="flex-1 min-h-0 mt-0">
          <div v-if="prStore.loading" class="flex-1 flex items-center justify-center py-12">
            <LoadingSpinner label="Loading pull requests..." />
          </div>
          <EmptyState
            v-else-if="activeReviewPRs.length === 0 && prStore.myPRs.length === 0 && prStore.teamPRs.length === 0"
            :icon="GitPullRequest"
            title="No pull requests"
            description="No open pull requests found. Your PRs, reviews, and team activity will appear here."
          />
          <ScrollArea v-else class="h-full">
            <div class="px-4 py-3 space-y-4">

              <!-- Needs Your Review -->
              <div>
                <div class="flex items-center gap-1.5 mb-1.5">
                  <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Needs Your Review</span>
                  <Badge v-if="activeReviewPRs.length" variant="secondary" class="text-[10px] h-4 px-1">{{ activeReviewPRs.length }}</Badge>
                </div>
                <Card v-if="activeReviewPRs.length > 0" class="overflow-hidden">
                  <CardContent class="p-0">
                    <div
                      v-for="pr in activeReviewPRs"
                      :key="pr.id"
                      class="flex flex-col gap-1 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0"
                      @click="openPR(pr)"
                    >
                      <div class="flex items-center gap-2">
                        <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
                        <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize shrink-0', prStatusClasses(pr.status))">{{ pr.status }}</Badge>
                      </div>
                      <div class="flex items-center gap-2 text-[10px] text-muted-foreground/60">
                        <span class="font-medium">{{ pr.repo }}</span>
                        <span>{{ branchName(pr.sourceBranch) }} → {{ branchName(pr.targetBranch) }}</span>
                        <span v-if="pr.createdBy" class="italic">by {{ pr.createdBy }}</span>
                        <span class="ml-auto flex items-center gap-1">
                          <template v-for="reviewer in parseReviewers(pr.reviewers)" :key="reviewer.uniqueName">
                            <span v-if="reviewer.vote !== 0" :title="`${reviewer.displayName}: ${voteIcon(reviewer.vote)}`">{{ voteIcon(reviewer.vote) }}</span>
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
                <div class="flex items-center gap-1.5 mb-1.5">
                  <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Your PRs</span>
                  <Badge v-if="prStore.myPRs.length" variant="secondary" class="text-[10px] h-4 px-1">{{ prStore.myPRs.length }}</Badge>
                </div>
                <Card v-if="prStore.myPRs.length > 0" class="overflow-hidden">
                  <CardContent class="p-0">
                    <div
                      v-for="pr in prStore.myPRs"
                      :key="pr.id"
                      class="flex flex-col gap-1 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0"
                      @click="openPR(pr)"
                    >
                      <div class="flex items-center gap-2">
                        <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
                        <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize shrink-0', prStatusClasses(pr.status))">{{ pr.status }}</Badge>
                      </div>
                      <div class="flex items-center gap-2 text-[10px] text-muted-foreground/60">
                        <span class="font-medium">{{ pr.repo }}</span>
                        <span>{{ branchName(pr.sourceBranch) }} → {{ branchName(pr.targetBranch) }}</span>
                        <div class="ml-auto flex items-center gap-1">
                          <template v-for="reviewer in parseReviewers(pr.reviewers)" :key="reviewer.uniqueName">
                            <span v-if="reviewer.vote !== 0" :title="`${reviewer.displayName}: ${voteIcon(reviewer.vote)}`">{{ voteIcon(reviewer.vote) }}</span>
                            <div v-else class="w-4 h-4 rounded-full bg-muted border border-background flex items-center justify-center" :title="reviewer.displayName">
                              <span class="text-[7px] font-medium text-muted-foreground">{{ getInitials(reviewer.displayName) }}</span>
                            </div>
                          </template>
                        </div>
                        <span class="tabular-nums shrink-0">{{ relativeTime(pr.updatedAt) }}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
                <p v-else class="text-[11px] text-muted-foreground/40 py-2 px-1">You have no PRs</p>
              </div>

              <!-- Team Activity -->
              <div>
                <div class="flex items-center gap-1.5 mb-1.5">
                  <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Team Activity</span>
                  <Badge v-if="prStore.teamPRs.length" variant="secondary" class="text-[10px] h-4 px-1">{{ prStore.teamPRs.length }}</Badge>
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
                        <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize shrink-0', prStatusClasses(pr.status))">{{ pr.status }}</Badge>
                      </div>
                      <div class="flex items-center gap-2 text-[10px] text-muted-foreground/60">
                        <span class="font-medium">{{ pr.repo }}</span>
                        <span v-if="pr.createdBy" class="italic">by {{ pr.createdBy }}</span>
                        <span>{{ branchName(pr.sourceBranch) }} → {{ branchName(pr.targetBranch) }}</span>
                        <div class="ml-auto flex items-center gap-1">
                          <template v-for="reviewer in parseReviewers(pr.reviewers)" :key="reviewer.uniqueName">
                            <span v-if="reviewer.vote !== 0" :title="`${reviewer.displayName}: ${voteIcon(reviewer.vote)}`">{{ voteIcon(reviewer.vote) }}</span>
                            <div v-else class="w-4 h-4 rounded-full bg-muted border border-background flex items-center justify-center" :title="reviewer.displayName">
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
        </TabsContent>

        <!-- ═══ Pipelines Tab ═══ -->
        <TabsContent value="pipelines" class="flex-1 min-h-0 mt-0">
          <div v-if="adoStore.pipelinesLoading" class="flex-1 flex items-center justify-center py-12">
            <LoadingSpinner label="Loading pipelines..." />
          </div>
          <EmptyState
            v-else-if="adoStore.pipelines.length === 0"
            :icon="GitMerge"
            title="No pipelines"
            description="No recent pipeline runs found. Pipeline activity will appear here after syncing."
          />
          <ScrollArea v-else class="h-full">
            <div class="px-4 py-3 space-y-3">
              <Card class="overflow-hidden">
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
              <Button
                v-if="adoStore.pipelines.length > 5 && !showAllPipelines"
                variant="ghost" size="sm"
                class="w-full text-xs h-7 mt-1"
                @click="showAllPipelines = true"
              >
                Show all {{ adoStore.pipelines.length }} pipelines
              </Button>
            </div>
          </ScrollArea>
        </TabsContent>
      </Tabs>
    </template>

    <!-- Link Dialog (searches local personal tasks) -->
    <LinkDialog
      v-model:open="linkDialogOpen"
      :ado-id="linkingAdoItem?.adoId ?? ''"
      :ado-title="linkingAdoItem?.title ?? ''"
      @linked="onLinked"
    />

    <!-- Import choice dialog -->
    <Dialog :open="importChoiceOpen" @update:open="(v) => { if (!v) { importChoiceOpen = false; importingAdoItem = null } }">
      <DialogContent class="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle class="text-sm">Import Work Item</DialogTitle>
          <DialogDescription class="text-xs">
            How would you like to import "{{ importingAdoItem?.title }}"?
          </DialogDescription>
        </DialogHeader>
        <div class="flex flex-col gap-2 py-2">
          <Button variant="outline" class="justify-start gap-2 h-10" @click="doImportAsTask">
            <CheckSquare :size="16" class="text-blue-500" />
            <div class="text-left">
              <div class="text-sm font-medium">Import as Task</div>
              <div class="text-[11px] text-muted-foreground">Create a personal task linked to this ADO item</div>
            </div>
          </Button>
          <Button variant="outline" class="justify-start gap-2 h-10" @click="doImportAsProject">
            <FolderKanban :size="16" class="text-purple-500" />
            <div class="text-left">
              <div class="text-sm font-medium">Import as Project</div>
              <div class="text-[11px] text-muted-foreground">Create a project with a linked task under it</div>
            </div>
          </Button>
        </div>
        <DialogFooter>
          <Button variant="ghost" size="sm" class="text-xs" @click="importChoiceOpen = false; importingAdoItem = null">Cancel</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
