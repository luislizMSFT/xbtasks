<script setup lang="ts">
import { onMounted, computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { usePRStore, parseReviewers, branchName, voteIcon, relativeTime } from '@/stores/prs'
import type { PullRequest } from '@/stores/prs'
import { cn } from '@/lib/utils'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import AdoBadge from '@/components/ui/AdoBadge.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  AlertTriangle,
  Octagon,
  GitPullRequest,
  Plus,
  Inbox,
  ArrowRight,
  X,
  RefreshCw,
  ChevronDown,
  ChevronRight,
  ExternalLink,
  Loader2,
  Wifi,
  WifiOff,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()
const prStore = usePRStore()

onMounted(async () => {
  await Promise.all([
    taskStore.fetchTasks(),
    prStore.fetchAll(),
  ])
})

// --- Inbox / Quick Capture ---
interface InboxItem {
  id: number
  text: string
  createdAt: string
}

const inboxItems = reactive<InboxItem[]>([
  { id: 1, text: 'Look into flaky E2E test on staging', createdAt: '2h ago' },
  { id: 2, text: 'Ask Alex about token refresh approach', createdAt: '5h ago' },
  { id: 3, text: 'Review new ADO permissions model doc', createdAt: '1d ago' },
])

const newInboxText = ref('')
let nextInboxId = 100

function addInboxItem() {
  const text = newInboxText.value.trim()
  if (!text) return
  inboxItems.unshift({ id: nextInboxId++, text, createdAt: 'just now' })
  newInboxText.value = ''
}

function removeInboxItem(id: number) {
  const idx = inboxItems.findIndex(i => i.id === id)
  if (idx !== -1) inboxItems.splice(idx, 1)
}

function promoteToTask(item: InboxItem) {
  taskStore.createTask(item.text)
  removeInboxItem(item.id)
}

// --- PR helpers ---
const teamPRsExpanded = ref(false)
const refreshing = ref(false)

const activeReviewPRs = computed(() =>
  prStore.reviewPRs.filter(pr => pr.status === 'active')
)

async function refreshPRs() {
  refreshing.value = true
  try { await prStore.fetchAll() }
  finally { refreshing.value = false }
}

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

const focusTasks = computed(() =>
  taskStore.tasks.filter(t => !['blocked', 'cancelled'].includes(t.status)).slice(0, 10)
)

const blockedTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'blocked')
)
</script>

<template>
  <div class="flex-1 flex overflow-hidden">
    <!-- Left column: Tasks + Inbox -->
    <ScrollArea class="flex-1 h-full">
      <div class="px-4 py-3 space-y-4">
        <!-- Inbox / Quick Capture -->
        <div>
          <div class="flex items-center gap-1.5 mb-2">
            <Inbox :size="12" class="text-muted-foreground" />
            <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Inbox</span>
            <span class="text-[11px] text-muted-foreground/50 tabular-nums">{{ inboxItems.length }}</span>
          </div>
          <div class="flex items-center gap-2 mb-1.5">
            <input
              v-model="newInboxText"
              @keydown.enter="addInboxItem"
              class="flex-1 h-7 text-xs bg-muted/50 border border-border/50 rounded px-2.5 placeholder:text-muted-foreground/40 focus:outline-none focus:ring-1 focus:ring-primary/30 text-foreground"
              placeholder="Quick capture..."
            />
            <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="addInboxItem">
              <Plus :size="13" />
            </Button>
          </div>
          <div v-if="inboxItems.length > 0" class="space-y-0.5">
            <div
              v-for="item in inboxItems"
              :key="item.id"
              class="group flex items-center gap-2 px-2.5 py-1.5 rounded hover:bg-muted/40 transition-colors"
            >
              <span class="text-[12px] text-foreground flex-1 truncate">{{ item.text }}</span>
              <span class="text-[10px] text-muted-foreground/40 shrink-0">{{ item.createdAt }}</span>
              <button
                class="h-5 w-5 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity shrink-0 rounded hover:bg-primary/10"
                title="Convert to task"
                @click="promoteToTask(item)"
              >
                <ArrowRight :size="11" class="text-primary" />
              </button>
              <button
                class="h-5 w-5 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity shrink-0 rounded hover:bg-muted"
                title="Dismiss"
                @click="removeInboxItem(item.id)"
              >
                <X :size="11" class="text-muted-foreground" />
              </button>
            </div>
          </div>
          <p v-else class="text-[11px] text-muted-foreground/40 py-1 px-2.5">Inbox empty</p>
        </div>

        <!-- Tasks -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Tasks</span>
            <span class="text-[11px] text-muted-foreground tabular-nums">{{ focusTasks.length }}</span>
          </div>
          <Card class="overflow-hidden">
            <CardContent class="p-0">
              <template v-if="focusTasks.length > 0">
                <div
                  v-for="task in focusTasks"
                  :key="task.id"
                  class="flex items-center gap-2.5 px-3 py-2 hover:bg-muted/50 cursor-pointer transition-colors border-b border-border/50 last:border-b-0"
                  @click="router.push('/tasks')"
                >
                  <StatusBadge :status="task.status" />
                  <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
                  <PriorityBadge :priority="task.priority" />
                  <AdoBadge :ado-id="task.adoId" />
                </div>
              </template>
              <div v-else class="px-3 py-6 text-center">
                <p class="text-xs text-muted-foreground">No tasks yet</p>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Blocked at bottom -->
        <div v-if="blockedTasks.length > 0">
          <div class="flex items-center gap-1.5 mb-2">
            <AlertTriangle :size="12" class="text-red-500" />
            <span class="text-[11px] font-semibold text-red-600 dark:text-red-400 uppercase tracking-wider">Blocked</span>
            <span class="text-[11px] text-red-500/60 tabular-nums">{{ blockedTasks.length }}</span>
          </div>
          <Card class="overflow-hidden border-red-500/20 bg-red-500/[0.02]">
            <CardContent class="p-0">
              <div
                v-for="task in blockedTasks"
                :key="task.id"
                class="flex items-center gap-2.5 px-3 py-2 cursor-pointer transition-colors hover:bg-red-500/5 border-b border-red-500/10 last:border-b-0"
                @click="router.push('/tasks')"
              >
                <Octagon :size="14" class="text-red-500 shrink-0" />
                <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
                <span class="text-[11px] text-red-500/70 truncate max-w-[10rem]">{{ task.blockedReason || 'Blocked' }}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </ScrollArea>

    <!-- Right column: PRs -->
    <ScrollArea class="w-[45%] h-full border-l border-border">
      <div class="px-4 py-3 space-y-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-1.5">
            <GitPullRequest :size="14" class="text-muted-foreground" />
            <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Pull Requests</span>
          </div>
          <div class="flex items-center gap-1.5">
            <span
              :class="cn(
                'flex items-center gap-1 text-[10px] px-1.5 py-0.5 rounded-full',
                prStore.connected
                  ? 'text-emerald-600 dark:text-emerald-400 bg-emerald-500/10'
                  : 'text-muted-foreground bg-muted/50'
              )"
            >
              <Wifi v-if="prStore.connected" :size="10" />
              <WifiOff v-else :size="10" />
              {{ prStore.connected ? 'ADO' : 'Cached' }}
            </span>
            <Button
              variant="ghost"
              size="icon"
              class="h-6 w-6"
              :disabled="refreshing"
              @click="refreshPRs"
            >
              <RefreshCw :size="12" :class="refreshing && 'animate-spin'" />
            </Button>
          </div>
        </div>

        <!-- Loading skeleton -->
        <template v-if="prStore.loading">
          <div class="space-y-2">
            <div v-for="i in 3" :key="i" class="h-10 rounded bg-muted/50 animate-pulse" />
          </div>
        </template>

        <template v-else>
          <!-- Needs Your Review -->
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
                    <ExternalLink :size="10" class="shrink-0 opacity-0 group-hover:opacity-100" />
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

          <!-- Team PRs (collapsible) -->
          <div>
            <button
              class="flex items-center gap-1 text-[10px] font-medium text-muted-foreground/60 uppercase tracking-wider hover:text-muted-foreground transition-colors"
              @click="teamPRsExpanded = !teamPRsExpanded"
            >
              <ChevronRight v-if="!teamPRsExpanded" :size="10" />
              <ChevronDown v-else :size="10" />
              Team PRs
              <span class="normal-case text-muted-foreground/40 tabular-nums">{{ prStore.teamPRs.length }}</span>
            </button>
            <template v-if="teamPRsExpanded">
              <Card v-if="prStore.teamPRs.length > 0" class="mt-1.5 overflow-hidden">
                <CardContent class="p-0">
                  <div
                    v-for="pr in prStore.teamPRs"
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
                          <span v-if="reviewer.vote !== 0">{{ voteIcon(reviewer.vote) }}</span>
                        </template>
                      </span>
                      <span class="tabular-nums shrink-0">{{ relativeTime(pr.updatedAt) }}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
              <p v-else class="text-[11px] text-muted-foreground/40 py-2 px-1">No other team PRs</p>
            </template>
          </div>
        </template>
      </div>
    </ScrollArea>
  </div>
</template>
