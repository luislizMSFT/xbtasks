<script setup lang="ts">
import { ref, computed, onMounted, onActivated, onDeactivated } from 'vue'
import type { AcceptableValue } from 'reka-ui'
import { cn } from '@/lib/utils'
import { adoTypeColor, adoTypeIcon, adoStateClasses, adoPriorityClasses, prStatusClasses } from '@/lib/styles'
import { useADOStore } from '@/stores/ado'
import { useNotify } from '@/composables/useNotify'
import type { ADOWorkItem, ADOPipeline } from '@/stores/ado'
import { useTaskStore } from '@/stores/tasks'
import { useProjectStore } from '@/stores/projects'
import { usePRStore, parseReviewers, branchName, voteIcon, relativeTime } from '@/stores/prs'
import type { PullRequest } from '@/stores/prs'
import PageHeader from '@/components/PageHeader.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'
import AdoTreeBrowser from '@/components/ado/AdoTreeBrowser.vue'
import AdoHierarchyTree from '@/components/ado/AdoHierarchyTree.vue'
import LinkDialog from '@/components/ado/LinkDialog.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardContent } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import {
  AlertCircle,
  CheckSquare,
  Circle,
  ExternalLink,
  GitPullRequest,
  CheckCircle2,
  XCircle,
  Loader2,
  Clock,
  GitMerge,
  Search,
  Download,
  Link,
  EyeOff,
  Eye,
  FolderKanban,
  Target,
  Layers,
} from 'lucide-vue-next'

const adoStore = useADOStore()
const taskStore = useTaskStore()
const projectStore = useProjectStore()
const prStore = usePRStore()

const isActive = ref(false)
onActivated(() => { isActive.value = true })
onDeactivated(() => { isActive.value = false })
onMounted(() => { isActive.value = true })

const activeTab = ref('browser')
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

const typeIcon = adoTypeIcon

// --- Tree browser actions ---

function onTreeSelect(item: ADOWorkItem) {
  adoStore.selectedItem = item
}

async function onTreeImport(item: ADOWorkItem) {
  importingAdoItem.value = item
  importChoiceOpen.value = true
}

const notify = useNotify()

async function doImportAsTask() {
  if (!importingAdoItem.value) return
  try {
    const { importWorkItem } = await import('@/api/links')
    await importWorkItem(importingAdoItem.value.adoId)
    await Promise.all([adoStore.fetchLinkedAdoIds(), taskStore.fetchTasks()])
    notify.success('Task imported from ADO')
  } catch (e: any) {
    adoStore.error = e?.message || 'Import failed'
    notify.error(e?.message || 'Import failed')
  } finally {
    importChoiceOpen.value = false
    importingAdoItem.value = null
  }
}

async function doImportAsProject() {
  if (!importingAdoItem.value) return
  try {
    const { importWorkItemAsProject } = await import('@/api/links')
    await importWorkItemAsProject(importingAdoItem.value.adoId)
    await Promise.all([
      adoStore.fetchLinkedAdoIds(),
      taskStore.fetchTasks(),
      projectStore.fetchProjects(),
    ])
    notify.success('Project imported from ADO')
  } catch (e: any) {
    adoStore.error = e?.message || 'Import as project failed'
    notify.error(e?.message || 'Import failed')
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
    const { openURL } = await import('@/api/browser')
    await openURL(url)
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
    <!-- Teleport tabs + sync to top bar (only when active) -->
    <Teleport v-if="isActive" to="#topbar-center">
      <Tabs v-model="activeTab">
        <TabsList class="h-7 bg-transparent p-0 gap-1">
          <TabsTrigger value="browser" class="text-[11px] gap-1 px-2.5 h-6 data-[state=active]:bg-muted">
            <CheckSquare :size="11" />
            Browser
            <Badge v-if="adoStore.workItemTree.length" variant="secondary" class="text-[9px] h-3.5 px-1">{{ adoStore.workItemTree.length }}</Badge>
          </TabsTrigger>
          <TabsTrigger value="pull-requests" class="text-[11px] gap-1 px-2.5 h-6 data-[state=active]:bg-muted">
            <GitPullRequest :size="11" />
            PRs
            <Badge v-if="totalPRCount" variant="secondary" class="text-[9px] h-3.5 px-1">{{ totalPRCount }}</Badge>
          </TabsTrigger>
          <TabsTrigger value="pipelines" class="text-[11px] gap-1 px-2.5 h-6 data-[state=active]:bg-muted">
            <GitMerge :size="11" />
            Pipelines
            <Badge v-if="adoStore.pipelines.length" variant="secondary" class="text-[9px] h-3.5 px-1">{{ adoStore.pipelines.length }}</Badge>
          </TabsTrigger>
        </TabsList>
      </Tabs>
    </Teleport>

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

    <!-- Tab content -->
    <template v-else>
      <div class="flex-1 flex flex-col min-h-0">

        <!-- ═══ ADO Browser Tab ═══ -->
        <template v-if="activeTab === 'browser'">
          <!-- Filter row (above list, not part of it) -->
          <div class="flex items-center gap-1.5 px-4 py-1.5 border-b border-border shrink-0">
            <span class="text-[10px] text-muted-foreground shrink-0">Type</span>
            <Select v-model="adoStore.filterType">
              <SelectTrigger class="h-6 w-[68px] text-[10px] px-1.5">
                <SelectValue placeholder="All" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all" class="text-xs">All</SelectItem>
                <SelectItem value="Task" class="text-xs">Task</SelectItem>
                <SelectItem value="Bug" class="text-xs">Bug</SelectItem>
                <SelectItem value="User Story" class="text-xs">User Story</SelectItem>
                <SelectItem value="Feature" class="text-xs">Feature</SelectItem>
                <SelectItem value="Deliverable" class="text-xs">Deliverable</SelectItem>
                <SelectItem value="Scenario" class="text-xs">Scenario</SelectItem>
                <SelectItem value="Epic" class="text-xs">Epic</SelectItem>
              </SelectContent>
            </Select>

            <span class="text-[10px] text-muted-foreground shrink-0">State</span>
            <Select v-model="adoStore.filterState">
              <SelectTrigger class="h-6 w-[60px] text-[10px] px-1.5">
                <SelectValue placeholder="All" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all" class="text-xs">All</SelectItem>
                <SelectItem value="Active" class="text-xs">Active</SelectItem>
                <SelectItem value="New" class="text-xs">New</SelectItem>
                <SelectItem value="Resolved" class="text-xs">Resolved</SelectItem>
                <SelectItem value="Closed" class="text-xs">Closed</SelectItem>
              </SelectContent>
            </Select>

            <span class="text-[10px] text-muted-foreground shrink-0">Area</span>
            <Select v-model="adoStore.filterArea">
              <SelectTrigger class="h-6 w-[80px] text-[10px] px-1.5">
                <SelectValue placeholder="All" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all" class="text-xs">All</SelectItem>
                <SelectItem
                  v-for="area in adoStore.availableAreaPaths"
                  :key="area"
                  :value="area"
                  class="text-xs"
                >{{ area }}</SelectItem>
              </SelectContent>
            </Select>

            <span class="text-[10px] text-muted-foreground shrink-0">Query</span>
            <Select :model-value="selectedQueryId" @update:model-value="onQueryChange">
              <SelectTrigger class="h-6 w-[90px] text-[10px] px-1.5">
                <SelectValue placeholder="My Items" />
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

            <div class="flex-1" />

            <Button
              :variant="adoStore.hideLinked ? 'default' : 'outline'"
              size="sm"
              class="h-6 text-[10px] gap-1"
              @click="adoStore.hideLinked = !adoStore.hideLinked"
            >
              <EyeOff v-if="adoStore.hideLinked" :size="11" />
              <Eye v-else :size="11" />
              Hide Linked
            </Button>
          </div>

          <!-- Two-column layout: tree + detail -->
          <div class="flex-1 flex min-h-0">
            <!-- Left: Tree browser (60%) -->
            <div class="w-[60%] border-r border-border flex flex-col">
              <!-- Search inside tree panel -->
              <div class="px-2 py-1.5 border-b border-border shrink-0">
                <div class="relative">
                  <Search :size="12" class="absolute left-2 top-1.5 text-muted-foreground" />
                  <Input
                    v-model="adoStore.searchQuery"
                    placeholder="Search work items..."
                    class="h-6 pl-7 text-[10px]"
                  />
                </div>
              </div>
              <ScrollArea class="flex-1">
                <!-- Loading skeleton for tree -->
                <div v-if="adoStore.loading" class="p-3 space-y-2">
                  <div v-for="i in 8" :key="i" class="flex items-center gap-2 px-2 py-1.5">
                    <Skeleton class="h-3 w-3 rounded shrink-0" />
                    <Skeleton class="h-3 flex-1 rounded" />
                    <Skeleton class="h-3 w-12 rounded" />
                  </div>
                </div>
                <div v-else class="py-2 px-2">
                  <EmptyState
                    v-if="adoStore.treeRoots.length === 0"
                    :icon="CheckSquare"
                    title="No work items"
                    description="No ADO work items match your current filters."
                    class="py-8"
                  />
                  <AdoTreeBrowser
                    v-else
                    :items="adoStore.treeRoots"
                    :get-children="adoStore.getChildren"
                    :is-linked="adoStore.isLinked"
                    :selected-id="adoStore.selectedItem?.adoId"
                    @select="onTreeSelect"
                    @import="onTreeImport"
                    @link="onTreeLink"
                  />
                </div>
              </ScrollArea>
            </div>

            <!-- Right: Detail panel (40%) -->
            <div class="w-[40%] flex flex-col">
              <template v-if="adoStore.selectedItem">
                <!-- Header: type + title + metadata (compact) -->
                <div class="shrink-0 border-b border-border px-4 py-3 space-y-2">
                  <div class="flex items-center gap-2">
                    <component :is="typeIcon(adoStore.selectedItem.type)" :size="16" :class="adoTypeColor(adoStore.selectedItem.type)" />
                    <h2 class="text-sm font-semibold text-foreground leading-snug flex-1 truncate">{{ adoStore.selectedItem.title }}</h2>
                    <Badge variant="outline" :class="['text-[10px] h-4 px-1.5', adoStateClasses(adoStore.selectedItem.state)]">
                      {{ adoStore.selectedItem.state }}
                    </Badge>
                    <span class="text-[10px] text-muted-foreground tabular-nums">#{{ adoStore.selectedItem.adoId }}</span>
                  </div>
                  <!-- Inline metadata row -->
                  <div class="flex items-center gap-3 text-[10px] text-muted-foreground flex-wrap">
                    <span v-if="adoStore.selectedItem.assignedTo">{{ adoStore.selectedItem.assignedTo }}</span>
                    <Badge variant="outline" :class="['text-[9px] h-3.5 px-1', adoPriorityClasses(adoStore.selectedItem.priority)]">
                      P{{ adoStore.selectedItem.priority }}
                    </Badge>
                    <span v-if="adoStore.selectedItem.areaPath" class="truncate max-w-[10rem]">{{ adoStore.selectedItem.areaPath }}</span>
                    <span v-if="adoStore.selectedItem.org">{{ adoStore.selectedItem.org }}/{{ adoStore.selectedItem.project }}</span>
                  </div>
                  <!-- Action buttons (compact row) -->
                  <div class="flex items-center gap-1.5">
                    <Button
                      v-if="!adoStore.isLinked(adoStore.selectedItem.adoId)"
                      size="sm"
                      class="h-6 text-[10px] gap-1"
                      @click="onTreeImport(adoStore.selectedItem)"
                    >
                      <Download :size="11" />
                      Import
                    </Button>
                    <Button
                      v-if="!adoStore.isLinked(adoStore.selectedItem.adoId)"
                      variant="outline"
                      size="sm"
                      class="h-6 text-[10px] gap-1"
                      @click="onTreeLink(adoStore.selectedItem)"
                    >
                      <Link :size="11" />
                      Link
                    </Button>
                    <Badge v-if="adoStore.isLinked(adoStore.selectedItem.adoId)" class="bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/25 text-[10px]">
                      <CheckCircle2 :size="10" class="mr-1" />
                      Linked
                    </Badge>
                    <div class="flex-1" />
                    <Button
                      variant="ghost"
                      size="sm"
                      class="h-6 text-[10px] gap-1"
                      @click="adoStore.selectedItem && openExternal(adoStore.selectedItem.url)"
                    >
                      <ExternalLink :size="11" />
                      Open in ADO
                    </Button>
                  </div>
                </div>

                <!-- Hierarchy (always visible, compact) -->
                <div class="shrink-0 border-b border-border px-4 py-2">
                  <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider mb-1">Hierarchy</h3>
                  <AdoHierarchyTree
                    :item="adoStore.selectedItem"
                    :get-children="adoStore.getChildren"
                    :all-items="adoStore.workItemTree"
                    @select="onTreeSelect"
                  />
                </div>

                <!-- Description (fills remaining height) -->
                <div class="flex-1 min-h-0">
                  <ScrollArea class="h-full">
                    <div v-if="adoStore.selectedItem.description" class="px-4 py-3">
                      <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider mb-2">Description</h3>
                      <div
                        class="text-xs text-foreground prose prose-sm max-w-none [&_*]:text-xs [&_*]:text-foreground"
                        v-html="adoStore.selectedItem.description"
                      />
                    </div>
                    <div v-else class="px-4 py-6 text-center">
                      <p class="text-[11px] text-muted-foreground/40">No description</p>
                    </div>
                  </ScrollArea>
                </div>
              </template>

              <!-- Empty state -->
              <div v-else class="flex-1 flex items-center justify-center">
                <p class="text-[11px] text-muted-foreground/40 text-center">Select a work item to view details</p>
              </div>
            </div>
          </div>
        </template>

        <!-- ═══ Pull Requests Tab ═══ -->
        <template v-if="activeTab === 'pull-requests'">
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
        </template>

        <!-- ═══ Pipelines Tab ═══ -->
        <template v-if="activeTab === 'pipelines'">
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
        </template>
      </div>
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
