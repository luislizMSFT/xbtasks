<script setup lang="ts">
/**
 * Playground: Top bar header for every page.
 *
 * Shows CURRENT production header vs PROPOSED improvements side-by-side,
 * for all 6 pages: Dashboard, Tasks, ADO, Projects, Dependencies, Settings.
 */
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import {
  RefreshCw, Wifi, Clock, CheckCircle2, ArrowLeft,
  Loader2, Search, Plus, CheckSquare, GitPullRequest, GitMerge,
  Network, Settings, FolderKanban, LayoutDashboard,
} from 'lucide-vue-next'

const router = useRouter()

const mockSync = ref({
  connected: true,
  lastSyncedAt: new Date(Date.now() - 3 * 60 * 1000).toISOString(),
  syncing: false,
  pendingChanges: 3,
})

function relativeTime(iso: string): string {
  const ms = Date.now() - new Date(iso).getTime()
  const min = Math.floor(ms / 60000)
  if (min < 1) return 'just now'
  if (min < 60) return `${min}m ago`
  return `${Math.floor(min / 60)}h ago`
}

function toggleConnection() {
  mockSync.value.connected = !mockSync.value.connected
}

function simulateSync() {
  mockSync.value.syncing = true
  setTimeout(() => {
    mockSync.value.syncing = false
    mockSync.value.lastSyncedAt = new Date().toISOString()
    mockSync.value.pendingChanges = 0
  }, 2000)
}
</script>

<template>
  <div class="h-screen w-screen flex flex-col bg-background text-foreground">
    <!-- Playground header -->
    <div class="shrink-0 border-b border-border bg-card px-4 py-2 flex items-center gap-3 titlebar-drag">
      <Button variant="ghost" size="sm" class="h-7 gap-1 titlebar-no-drag" @click="router.back()">
        <ArrowLeft :size="14" /> Back
      </Button>
      <span class="text-sm font-medium">Top Bar — All Pages</span>
      <span class="text-[10px] text-muted-foreground bg-muted px-2 py-0.5 rounded-full">Current vs Proposed</span>
      <div class="flex-1" />
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" class="h-7 text-[10px]" @click="toggleConnection">
          Toggle Connection
        </Button>
        <Button variant="outline" size="sm" class="h-7 text-[10px]" @click="simulateSync">
          Trigger Sync
        </Button>
        <Button variant="outline" size="sm" class="h-7 text-[10px]" @click="mockSync.pendingChanges = mockSync.pendingChanges ? 0 : 5">
          Toggle Pending
        </Button>
      </div>
    </div>

    <ScrollArea class="flex-1 min-h-0">
      <div class="max-w-[1100px] mx-auto px-6 py-6 space-y-10">

        <!-- ═══════════════════════════════════════ -->
        <!-- DASHBOARD PAGE                          -->
        <!-- ═══════════════════════════════════════ -->
        <section>
          <h2 class="text-xs font-bold text-muted-foreground uppercase tracking-wider mb-3 flex items-center gap-2">
            <LayoutDashboard :size="13" /> Dashboard
          </h2>

          <p class="text-[10px] text-muted-foreground mb-1.5">Current — dots only, no labels</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 mb-4">
            <span class="text-sm font-medium shrink-0">Dashboard</span>
            <div class="flex items-center gap-2 ml-3 text-[11px]">
              <span class="tabular-nums text-muted-foreground">5/47</span>
              <span class="w-2 h-2 rounded-full bg-red-500 animate-pulse" title="2 blocked" />
              <span class="w-2 h-2 rounded-full bg-green-500" title="ADO connected" />
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0">
                <RefreshCw :size="12" class="text-muted-foreground" />
              </Button>
            </div>
            <div class="flex-1" />
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>

          <p class="text-[10px] text-emerald-600 mb-1.5">Proposed — sync status with labels, stats in center</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 border-l-2 border-l-emerald-500">
            <span class="text-sm font-medium shrink-0">Dashboard</span>
            <span class="w-px h-4 bg-border mx-3" />
            <div class="flex items-center gap-2 text-[11px] shrink-0">
              <span class="inline-flex items-center gap-1.5" :class="mockSync.connected ? 'text-green-600' : 'text-red-500'">
                <span class="w-1.5 h-1.5 rounded-full" :class="mockSync.connected ? 'bg-green-500' : 'bg-red-500'" />
                {{ mockSync.connected ? 'Synced' : 'Offline' }}
              </span>
              <span v-if="mockSync.connected" class="text-muted-foreground/50 tabular-nums">{{ relativeTime(mockSync.lastSyncedAt) }}</span>
              <Badge v-if="mockSync.pendingChanges > 0" variant="outline" class="text-[9px] h-4 px-1.5 gap-1 text-amber-600 border-amber-500/30 bg-amber-500/10">
                <span class="w-1 h-1 rounded-full bg-amber-500" /> {{ mockSync.pendingChanges }} pending
              </Badge>
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="simulateSync" :disabled="mockSync.syncing">
                <component :is="mockSync.syncing ? Loader2 : RefreshCw" :size="12" :class="mockSync.syncing && 'animate-spin'" class="text-muted-foreground" />
              </Button>
            </div>
            <div class="flex-1 flex items-center justify-center gap-2 text-[10px]">
              <Badge variant="secondary" class="text-[9px] h-4 px-1.5">5 active</Badge>
              <Badge variant="destructive" class="text-[9px] h-4 px-1.5">2 blocked</Badge>
              <span class="text-muted-foreground/40 tabular-nums">31/47 done</span>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>
        </section>

        <!-- ═══════════════════════════════════════ -->
        <!-- TASKS PAGE                              -->
        <!-- ═══════════════════════════════════════ -->
        <section>
          <h2 class="text-xs font-bold text-muted-foreground uppercase tracking-wider mb-3 flex items-center gap-2">
            <CheckSquare :size="13" /> Tasks
          </h2>

          <p class="text-[10px] text-muted-foreground mb-1.5">Current — status chips + sync icon, no sync info</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 mb-4">
            <span class="text-sm font-medium shrink-0">Tasks</span>
            <div class="flex items-center gap-1 ml-3">
              <Badge variant="default" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">All</Badge>
              <Badge variant="outline" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">Active</Badge>
              <Badge variant="outline" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">Blocked</Badge>
              <Badge variant="outline" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">Done</Badge>
              <Button variant="ghost" size="icon" class="h-6 w-6 ml-1"><RefreshCw :size="12" /></Button>
            </div>
            <div class="flex-1" />
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>

          <p class="text-[10px] text-emerald-600 mb-1.5">Proposed — sync cluster + status chips centered</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 border-l-2 border-l-emerald-500">
            <span class="text-sm font-medium shrink-0">Tasks</span>
            <span class="w-px h-4 bg-border mx-3" />
            <div class="flex items-center gap-2 text-[11px] shrink-0">
              <span class="inline-flex items-center gap-1.5" :class="mockSync.connected ? 'text-green-600' : 'text-red-500'">
                <span class="w-1.5 h-1.5 rounded-full" :class="mockSync.connected ? 'bg-green-500' : 'bg-red-500'" />
                {{ mockSync.connected ? 'Synced' : 'Offline' }}
              </span>
              <span v-if="mockSync.connected" class="text-muted-foreground/50 tabular-nums">{{ relativeTime(mockSync.lastSyncedAt) }}</span>
              <Badge v-if="mockSync.pendingChanges > 0" variant="outline" class="text-[9px] h-4 px-1.5 gap-1 text-amber-600 border-amber-500/30 bg-amber-500/10">
                <span class="w-1 h-1 rounded-full bg-amber-500" /> {{ mockSync.pendingChanges }} pending
              </Badge>
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="simulateSync" :disabled="mockSync.syncing">
                <component :is="mockSync.syncing ? Loader2 : RefreshCw" :size="12" :class="mockSync.syncing && 'animate-spin'" class="text-muted-foreground" />
              </Button>
            </div>
            <div class="flex-1 flex items-center justify-center gap-1">
              <Badge variant="default" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">All</Badge>
              <Badge variant="outline" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">Active</Badge>
              <Badge variant="outline" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">Blocked</Badge>
              <Badge variant="outline" class="cursor-pointer text-[10px] px-1.5 py-0 h-5">Done</Badge>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>
        </section>

        <!-- ═══════════════════════════════════════ -->
        <!-- ADO PAGE                                -->
        <!-- ═══════════════════════════════════════ -->
        <section>
          <h2 class="text-xs font-bold text-muted-foreground uppercase tracking-wider mb-3 flex items-center gap-2">
            <AzureDevOpsIcon :size="13" /> Azure DevOps
          </h2>

          <p class="text-[10px] text-muted-foreground mb-1.5">Current — tabs + Connected pill + Sync button + last sync time</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 mb-4">
            <span class="text-sm font-medium shrink-0">Azure DevOps</span>
            <div class="flex items-center gap-2 ml-3">
              <Tabs model-value="browser">
                <TabsList class="h-7 bg-transparent p-0 gap-1">
                  <TabsTrigger value="browser" class="text-[11px] gap-1 px-2.5 h-6 data-[state=active]:bg-muted">
                    <CheckSquare :size="11" /> Browser <Badge variant="secondary" class="text-[9px] h-3.5 px-1">12</Badge>
                  </TabsTrigger>
                  <TabsTrigger value="prs" class="text-[11px] gap-1 px-2.5 h-6">
                    <GitPullRequest :size="11" /> PRs <Badge variant="secondary" class="text-[9px] h-3.5 px-1">5</Badge>
                  </TabsTrigger>
                  <TabsTrigger value="pipelines" class="text-[11px] gap-1 px-2.5 h-6">
                    <GitMerge :size="11" /> Pipelines <Badge variant="secondary" class="text-[9px] h-3.5 px-1">4</Badge>
                  </TabsTrigger>
                </TabsList>
              </Tabs>
              <span class="text-[10px] text-muted-foreground">3m ago</span>
              <span class="flex items-center gap-1 text-[10px] px-1.5 py-0.5 rounded-full text-emerald-600 bg-emerald-500/10">
                <Wifi :size="10" /> Connected
              </span>
              <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1"><RefreshCw :size="11" /> Sync</Button>
            </div>
            <div class="flex-1" />
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>

          <p class="text-[10px] text-emerald-600 mb-1.5">Proposed — unified sync cluster + tabs centered</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 border-l-2 border-l-emerald-500">
            <span class="text-sm font-medium shrink-0">Azure DevOps</span>
            <span class="w-px h-4 bg-border mx-3" />
            <div class="flex items-center gap-2 text-[11px] shrink-0">
              <span class="inline-flex items-center gap-1.5" :class="mockSync.connected ? 'text-green-600' : 'text-red-500'">
                <span class="w-1.5 h-1.5 rounded-full" :class="mockSync.connected ? 'bg-green-500' : 'bg-red-500'" />
                {{ mockSync.connected ? 'Synced' : 'Offline' }}
              </span>
              <span v-if="mockSync.connected" class="text-muted-foreground/50 tabular-nums">{{ relativeTime(mockSync.lastSyncedAt) }}</span>
              <Badge v-if="mockSync.pendingChanges > 0" variant="outline" class="text-[9px] h-4 px-1.5 gap-1 text-amber-600 border-amber-500/30 bg-amber-500/10">
                <span class="w-1 h-1 rounded-full bg-amber-500" /> {{ mockSync.pendingChanges }} pending
              </Badge>
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="simulateSync" :disabled="mockSync.syncing">
                <component :is="mockSync.syncing ? Loader2 : RefreshCw" :size="12" :class="mockSync.syncing && 'animate-spin'" class="text-muted-foreground" />
              </Button>
            </div>
            <div class="flex-1 flex items-center justify-center">
              <Tabs model-value="browser">
                <TabsList class="h-7 bg-transparent p-0 gap-1">
                  <TabsTrigger value="browser" class="text-[11px] gap-1 px-2.5 h-6 data-[state=active]:bg-muted">
                    <CheckSquare :size="11" /> Browser <Badge variant="secondary" class="text-[9px] h-3.5 px-1">12</Badge>
                  </TabsTrigger>
                  <TabsTrigger value="prs" class="text-[11px] gap-1 px-2.5 h-6">
                    <GitPullRequest :size="11" /> PRs <Badge variant="secondary" class="text-[9px] h-3.5 px-1">5</Badge>
                  </TabsTrigger>
                  <TabsTrigger value="pipelines" class="text-[11px] gap-1 px-2.5 h-6">
                    <GitMerge :size="11" /> Pipelines <Badge variant="secondary" class="text-[9px] h-3.5 px-1">4</Badge>
                  </TabsTrigger>
                </TabsList>
              </Tabs>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>
        </section>

        <!-- ═══════════════════════════════════════ -->
        <!-- PROJECTS PAGE                           -->
        <!-- ═══════════════════════════════════════ -->
        <section>
          <h2 class="text-xs font-bold text-muted-foreground uppercase tracking-wider mb-3 flex items-center gap-2">
            <FolderKanban :size="13" /> Projects
          </h2>

          <p class="text-[10px] text-muted-foreground mb-1.5">Current — empty, no page-specific content</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 mb-4">
            <span class="text-sm font-medium shrink-0">Projects</span>
            <div class="flex-1" />
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>

          <p class="text-[10px] text-emerald-600 mb-1.5">Proposed — add sync cluster + project count</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 border-l-2 border-l-emerald-500">
            <span class="text-sm font-medium shrink-0">Projects</span>
            <span class="w-px h-4 bg-border mx-3" />
            <div class="flex items-center gap-2 text-[11px] shrink-0">
              <span class="inline-flex items-center gap-1.5" :class="mockSync.connected ? 'text-green-600' : 'text-red-500'">
                <span class="w-1.5 h-1.5 rounded-full" :class="mockSync.connected ? 'bg-green-500' : 'bg-red-500'" />
                {{ mockSync.connected ? 'Synced' : 'Offline' }}
              </span>
              <span v-if="mockSync.connected" class="text-muted-foreground/50 tabular-nums">{{ relativeTime(mockSync.lastSyncedAt) }}</span>
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="simulateSync" :disabled="mockSync.syncing">
                <component :is="mockSync.syncing ? Loader2 : RefreshCw" :size="12" :class="mockSync.syncing && 'animate-spin'" class="text-muted-foreground" />
              </Button>
            </div>
            <div class="flex-1 flex items-center justify-center gap-2 text-[10px]">
              <span class="text-muted-foreground tabular-nums">4 projects</span>
              <Badge variant="secondary" class="text-[9px] h-4 px-1.5">2 active</Badge>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>
        </section>

        <!-- ═══════════════════════════════════════ -->
        <!-- DEPENDENCIES PAGE                       -->
        <!-- ═══════════════════════════════════════ -->
        <section>
          <h2 class="text-xs font-bold text-muted-foreground uppercase tracking-wider mb-3 flex items-center gap-2">
            <Network :size="13" /> Dependencies
          </h2>

          <p class="text-[10px] text-muted-foreground mb-1.5">Current — empty</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 mb-4">
            <span class="text-sm font-medium shrink-0">Dependencies</span>
            <div class="flex-1" />
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>

          <p class="text-[10px] text-emerald-600 mb-1.5">Proposed — sync cluster + graph stats</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 border-l-2 border-l-emerald-500">
            <span class="text-sm font-medium shrink-0">Dependencies</span>
            <span class="w-px h-4 bg-border mx-3" />
            <div class="flex items-center gap-2 text-[11px] shrink-0">
              <span class="inline-flex items-center gap-1.5" :class="mockSync.connected ? 'text-green-600' : 'text-red-500'">
                <span class="w-1.5 h-1.5 rounded-full" :class="mockSync.connected ? 'bg-green-500' : 'bg-red-500'" />
                {{ mockSync.connected ? 'Synced' : 'Offline' }}
              </span>
              <span v-if="mockSync.connected" class="text-muted-foreground/50 tabular-nums">{{ relativeTime(mockSync.lastSyncedAt) }}</span>
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="simulateSync" :disabled="mockSync.syncing">
                <component :is="mockSync.syncing ? Loader2 : RefreshCw" :size="12" :class="mockSync.syncing && 'animate-spin'" class="text-muted-foreground" />
              </Button>
            </div>
            <div class="flex-1 flex items-center justify-center gap-2 text-[10px]">
              <span class="text-muted-foreground tabular-nums">8 nodes · 5 edges</span>
              <Badge variant="outline" class="text-[9px] h-4 px-1.5 border-red-500/30 text-red-500">1 cycle</Badge>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>
        </section>

        <!-- ═══════════════════════════════════════ -->
        <!-- SETTINGS PAGE                           -->
        <!-- ═══════════════════════════════════════ -->
        <section>
          <h2 class="text-xs font-bold text-muted-foreground uppercase tracking-wider mb-3 flex items-center gap-2">
            <Settings :size="13" /> Settings
          </h2>

          <p class="text-[10px] text-muted-foreground mb-1.5">Current — empty</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 mb-4">
            <span class="text-sm font-medium shrink-0">Settings</span>
            <div class="flex-1" />
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>

          <p class="text-[10px] text-emerald-600 mb-1.5">Proposed — sync cluster only (no page content needed)</p>
          <div class="h-[46px] bg-card border border-border rounded-lg flex items-center px-4 border-l-2 border-l-emerald-500">
            <span class="text-sm font-medium shrink-0">Settings</span>
            <span class="w-px h-4 bg-border mx-3" />
            <div class="flex items-center gap-2 text-[11px] shrink-0">
              <span class="inline-flex items-center gap-1.5" :class="mockSync.connected ? 'text-green-600' : 'text-red-500'">
                <span class="w-1.5 h-1.5 rounded-full" :class="mockSync.connected ? 'bg-green-500' : 'bg-red-500'" />
                {{ mockSync.connected ? 'Synced' : 'Offline' }}
              </span>
              <span v-if="mockSync.connected" class="text-muted-foreground/50 tabular-nums">{{ relativeTime(mockSync.lastSyncedAt) }}</span>
              <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="simulateSync" :disabled="mockSync.syncing">
                <component :is="mockSync.syncing ? Loader2 : RefreshCw" :size="12" :class="mockSync.syncing && 'animate-spin'" class="text-muted-foreground" />
              </Button>
            </div>
            <div class="flex-1" />
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[140px] justify-start">
                <Search :size="14" /><span class="flex-1 text-left">Search...</span><kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
              </Button>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1"><Plus :size="13" /> New</Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-muted-foreground"><Clock :size="15" /></Button>
            </div>
          </div>
        </section>

        <!-- ═══════════════════════════════════════ -->
        <!-- DESIGN NOTES                            -->
        <!-- ═══════════════════════════════════════ -->
        <section class="text-xs text-muted-foreground space-y-1.5 max-w-lg pb-6">
          <p class="font-bold text-foreground text-sm mb-2">Design Pattern</p>
          <p>• <strong>Left zone</strong> — Page name + sync cluster (consistent everywhere)</p>
          <p>• <strong>Center zone</strong> — Per-page content (filters, tabs, stats)</p>
          <p>• <strong>Right zone</strong> — Search ⌘K + New + Activity (always the same)</p>
          <p class="mt-3 font-semibold text-foreground">Key improvements:</p>
          <p>• Sync status visible on <em>every</em> page (currently only Dashboard/Tasks/ADO)</p>
          <p>• Consistent layout: sync left, page content center, actions right</p>
          <p>• ADO page: tabs move to center (currently left-packed, crowded)</p>
          <p>• Dashboard: labeled stats replace cryptic ratio + dots</p>
          <p>• Projects/Deps: get useful page-specific stats in center zone</p>
        </section>

      </div>
    </ScrollArea>
  </div>
</template>
