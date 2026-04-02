<script setup lang="ts">
import { onMounted, computed, ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { cn } from '@/lib/utils'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import AdoBadge from '@/components/ui/AdoBadge.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import {
  Plus,
  AlertTriangle,
  Clock,
  CheckCircle2,
  Loader2,
  Octagon,
  GitPullRequest,
  ThumbsUp,
  ThumbsDown,
  MessageSquare,
  ArrowRight,
  BarChart3,
  Zap,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()

onMounted(() => {
  taskStore.fetchTasks()
})

// --- PR mock data ---

interface PR {
  id: number
  title: string
  status: 'active' | 'draft' | 'completed' | 'abandoned'
  repo: string
  author: string
  isReviewRequest: boolean
  unread: boolean
  votes: { up: number; down: number }
  reviewers: string[]
  updatedAt: string
}

const mockPRs = reactive<PR[]>([
  { id: 1, title: 'Fix auth redirect loop', status: 'active', repo: 'xb-services', author: 'Luis', isReviewRequest: false, unread: true, votes: { up: 1, down: 0 }, reviewers: ['Alex', 'Sam'], updatedAt: '30m ago' },
  { id: 2, title: 'Add rate limiting to gateway', status: 'active', repo: 'xb-gateway', author: 'Alex', isReviewRequest: true, unread: true, votes: { up: 0, down: 0 }, reviewers: ['Luis', 'Jordan'], updatedAt: '1h ago' },
  { id: 3, title: 'Schema migration v2', status: 'active', repo: 'xb-services', author: 'Jordan', isReviewRequest: true, unread: true, votes: { up: 0, down: 0 }, reviewers: ['Luis'], updatedAt: '2h ago' },
  { id: 4, title: 'Update CI pipeline', status: 'completed', repo: 'xb-infra', author: 'Luis', isReviewRequest: false, unread: false, votes: { up: 2, down: 0 }, reviewers: ['Sam', 'Alex'], updatedAt: '1d ago' },
])

const reviewRequests = computed(() => mockPRs.filter(pr => pr.isReviewRequest))
const myPRs = computed(() => mockPRs.filter(pr => !pr.isReviewRequest))
const unreadCount = computed(() => mockPRs.filter(pr => pr.unread).length)

function markViewed(pr: PR) {
  pr.unread = false
}

const prStatusMap: Record<PR['status'], { label: string; variant: 'default' | 'secondary' | 'outline' | 'destructive' }> = {
  active: { label: 'Active', variant: 'default' },
  draft: { label: 'Draft', variant: 'secondary' },
  completed: { label: 'Completed', variant: 'outline' },
  abandoned: { label: 'Abandoned', variant: 'destructive' },
}

function getInitials(name: string) {
  return name.slice(0, 2).toUpperCase()
}

// --- Task data ---

const focusTasks = computed(() =>
  taskStore.tasks.filter(t => ['in_progress', 'in_review'].includes(t.status)).slice(0, 5)
)

const blockedTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'blocked')
)

// --- Mock activity timeline ---

const activityEvents = ref([
  { id: 1, icon: 'check', description: 'Completed "Add unit tests for auth module"', timestamp: '25m ago' },
  { id: 2, icon: 'pr', description: 'PR merged: Fix auth redirect loop', timestamp: '1h ago' },
  { id: 3, icon: 'comment', description: 'Commented on "Schema migration v2"', timestamp: '2h ago' },
  { id: 4, icon: 'move', description: 'Moved "API rate limiting" to In Review', timestamp: '3h ago' },
  { id: 5, icon: 'create', description: 'Created task "Update deployment docs"', timestamp: '5h ago' },
])

function goToTasks() {
  router.push('/tasks')
}

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <div class="max-w-6xl mx-auto px-6 py-6 space-y-4">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold text-foreground">Dashboard</h1>
          <p class="text-sm text-muted-foreground mt-0.5">Your task overview at a glance</p>
        </div>
        <Button size="sm" @click="goToTasks">
          <Plus :size="14" />
          Create Task
        </Button>
      </div>

      <!-- Row 1: Pull Requests (full width) -->
      <Card class="transition-shadow hover:shadow-md">
        <CardHeader class="px-4 py-3">
          <div class="flex items-center justify-between">
            <CardTitle class="text-sm flex items-center gap-2">
              <GitPullRequest :size="16" class="text-muted-foreground" />
              Pull Requests
            </CardTitle>
            <Badge v-if="unreadCount > 0" variant="default" class="text-[11px]">
              {{ unreadCount }} unread
            </Badge>
          </div>
        </CardHeader>
        <Separator />
        <CardContent class="p-0">
          <!-- Needs Your Review -->
          <template v-if="reviewRequests.length > 0">
            <div class="px-4 py-2 bg-muted/30">
              <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Needs Your Review</span>
            </div>
            <div
              v-for="pr in reviewRequests"
              :key="pr.id"
              :class="cn(
                'flex items-center gap-3 px-4 py-2.5 cursor-pointer transition-colors hover:bg-muted/50',
                pr.unread && 'bg-blue-500/5'
              )"
              @click="markViewed(pr)"
            >
              <!-- Unread dot -->
              <div class="w-2 flex-shrink-0 flex justify-center">
                <span v-if="pr.unread" class="block w-1.5 h-1.5 rounded-full bg-blue-500" />
              </div>
              <!-- Title -->
              <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
              <!-- Status badge -->
              <Badge :variant="prStatusMap[pr.status].variant" class="text-[10px] px-1.5 py-0">
                {{ prStatusMap[pr.status].label }}
              </Badge>
              <!-- Repo -->
              <span class="text-[11px] text-muted-foreground font-mono hidden sm:inline">{{ pr.repo }}</span>
              <!-- Reviewer avatars -->
              <div class="flex -space-x-1.5">
                <div
                  v-for="reviewer in pr.reviewers"
                  :key="reviewer"
                  class="w-5 h-5 rounded-full bg-muted border border-background flex items-center justify-center"
                  :title="reviewer"
                >
                  <span class="text-[9px] font-medium text-muted-foreground">{{ getInitials(reviewer) }}</span>
                </div>
              </div>
              <!-- Votes -->
              <div class="flex items-center gap-1 text-[11px] text-muted-foreground tabular-nums min-w-[3rem] justify-end">
                <ThumbsUp v-if="pr.votes.up > 0" :size="11" class="text-emerald-500" />
                <span v-if="pr.votes.up > 0" class="text-emerald-600 dark:text-emerald-400">{{ pr.votes.up }}</span>
                <ThumbsDown v-if="pr.votes.down > 0" :size="11" class="text-red-500" />
                <span v-if="pr.votes.down > 0" class="text-red-600 dark:text-red-400">{{ pr.votes.down }}</span>
              </div>
              <!-- Timestamp -->
              <span class="text-[11px] text-muted-foreground tabular-nums w-12 text-right">{{ pr.updatedAt }}</span>
            </div>
          </template>

          <!-- Your PRs -->
          <template v-if="myPRs.length > 0">
            <div class="px-4 py-2 bg-muted/30">
              <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Your PRs</span>
            </div>
            <div
              v-for="pr in myPRs"
              :key="pr.id"
              :class="cn(
                'flex items-center gap-3 px-4 py-2.5 cursor-pointer transition-colors hover:bg-muted/50',
                pr.unread && 'bg-blue-500/5'
              )"
              @click="markViewed(pr)"
            >
              <div class="w-2 flex-shrink-0 flex justify-center">
                <span v-if="pr.unread" class="block w-1.5 h-1.5 rounded-full bg-blue-500" />
              </div>
              <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
              <Badge :variant="prStatusMap[pr.status].variant" class="text-[10px] px-1.5 py-0">
                {{ prStatusMap[pr.status].label }}
              </Badge>
              <span class="text-[11px] text-muted-foreground font-mono hidden sm:inline">{{ pr.repo }}</span>
              <div class="flex -space-x-1.5">
                <div
                  v-for="reviewer in pr.reviewers"
                  :key="reviewer"
                  class="w-5 h-5 rounded-full bg-muted border border-background flex items-center justify-center"
                  :title="reviewer"
                >
                  <span class="text-[9px] font-medium text-muted-foreground">{{ getInitials(reviewer) }}</span>
                </div>
              </div>
              <div class="flex items-center gap-1 text-[11px] text-muted-foreground tabular-nums min-w-[3rem] justify-end">
                <ThumbsUp v-if="pr.votes.up > 0" :size="11" class="text-emerald-500" />
                <span v-if="pr.votes.up > 0" class="text-emerald-600 dark:text-emerald-400">{{ pr.votes.up }}</span>
                <ThumbsDown v-if="pr.votes.down > 0" :size="11" class="text-red-500" />
                <span v-if="pr.votes.down > 0" class="text-red-600 dark:text-red-400">{{ pr.votes.down }}</span>
              </div>
              <span class="text-[11px] text-muted-foreground tabular-nums w-12 text-right">{{ pr.updatedAt }}</span>
            </div>
          </template>
        </CardContent>
      </Card>

      <!-- Row 2: Focus Tasks + Stats -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Focus Tasks -->
        <Card class="transition-shadow hover:shadow-md">
          <CardHeader class="px-4 py-3">
            <CardTitle class="text-sm flex items-center gap-2">
              <Zap :size="16" class="text-amber-500" />
              Focus Tasks
            </CardTitle>
          </CardHeader>
          <Separator />
          <CardContent class="p-0">
            <template v-if="focusTasks.length > 0">
              <div
                v-for="task in focusTasks"
                :key="task.id"
                class="flex items-center gap-3 px-4 py-2.5 hover:bg-muted/50 cursor-pointer transition-colors"
                @click="router.push('/tasks')"
              >
                <StatusBadge :status="task.status" />
                <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
                <PriorityBadge :priority="task.priority" />
                <AdoBadge :ado-id="task.adoId" />
              </div>
            </template>
            <div v-else class="px-4 py-8 text-center">
              <p class="text-sm text-muted-foreground">No tasks in progress</p>
            </div>
          </CardContent>
        </Card>

        <!-- Stats -->
        <Card class="transition-shadow hover:shadow-md">
          <CardHeader class="px-4 py-3">
            <CardTitle class="text-sm flex items-center gap-2">
              <BarChart3 :size="16" class="text-muted-foreground" />
              Stats
            </CardTitle>
          </CardHeader>
          <Separator />
          <CardContent class="p-4">
            <div class="grid grid-cols-2 gap-3">
              <div class="rounded-lg bg-blue-500/10 p-3">
                <div class="flex items-center gap-2 mb-1">
                  <Loader2 :size="14" class="text-blue-500" />
                  <span class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">In Progress</span>
                </div>
                <p class="text-2xl font-bold text-foreground tabular-nums">{{ taskStore.stats.inProgress }}</p>
              </div>
              <div class="rounded-lg bg-red-500/10 p-3">
                <div class="flex items-center gap-2 mb-1">
                  <Octagon :size="14" class="text-red-500" />
                  <span class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Blocked</span>
                </div>
                <p class="text-2xl font-bold text-foreground tabular-nums">{{ taskStore.stats.blocked }}</p>
              </div>
              <div class="rounded-lg bg-emerald-500/10 p-3">
                <div class="flex items-center gap-2 mb-1">
                  <CheckCircle2 :size="14" class="text-emerald-500" />
                  <span class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Done</span>
                </div>
                <p class="text-2xl font-bold text-foreground tabular-nums">{{ taskStore.stats.done }}</p>
              </div>
              <div class="rounded-lg bg-zinc-500/10 p-3">
                <div class="flex items-center gap-2 mb-1">
                  <Clock :size="14" class="text-zinc-500" />
                  <span class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Total</span>
                </div>
                <p class="text-2xl font-bold text-foreground tabular-nums">{{ taskStore.stats.total }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Row 3: Blocked + Recent Activity -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Blocked -->
        <Card
          v-if="blockedTasks.length > 0"
          class="border-red-500/20 bg-red-500/5 transition-shadow hover:shadow-md"
        >
          <CardHeader class="px-4 py-3">
            <CardTitle class="text-sm flex items-center gap-2">
              <AlertTriangle :size="16" class="text-red-500" />
              <span class="text-red-600 dark:text-red-400">Blocked</span>
            </CardTitle>
          </CardHeader>
          <Separator class="bg-red-500/10" />
          <CardContent class="p-0">
            <div
              v-for="task in blockedTasks"
              :key="task.id"
              class="flex items-center gap-3 px-4 py-2.5 cursor-pointer transition-colors hover:bg-red-500/5"
              @click="router.push('/tasks')"
            >
              <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
              <span class="text-xs text-red-500/80 truncate max-w-[10rem]">{{ task.blockedReason || 'Blocked' }}</span>
              <AdoBadge :ado-id="task.adoId" />
            </div>
          </CardContent>
        </Card>

        <!-- Recent Activity -->
        <Card :class="cn('transition-shadow hover:shadow-md', blockedTasks.length === 0 && 'md:col-span-2')">
          <CardHeader class="px-4 py-3">
            <CardTitle class="text-sm flex items-center gap-2">
              <Clock :size="16" class="text-muted-foreground" />
              Recent Activity
            </CardTitle>
          </CardHeader>
          <Separator />
          <CardContent class="p-0">
            <div
              v-for="event in activityEvents"
              :key="event.id"
              class="flex items-center gap-3 px-4 py-2.5 transition-colors hover:bg-muted/50"
            >
              <div class="w-6 h-6 rounded-full bg-muted flex items-center justify-center flex-shrink-0">
                <CheckCircle2 v-if="event.icon === 'check'" :size="12" class="text-emerald-500" />
                <GitPullRequest v-else-if="event.icon === 'pr'" :size="12" class="text-purple-500" />
                <MessageSquare v-else-if="event.icon === 'comment'" :size="12" class="text-blue-500" />
                <ArrowRight v-else-if="event.icon === 'move'" :size="12" class="text-amber-500" />
                <Plus v-else :size="12" class="text-muted-foreground" />
              </div>
              <span class="text-sm text-foreground flex-1 truncate">{{ event.description }}</span>
              <span class="text-[11px] text-muted-foreground tabular-nums">{{ event.timestamp }}</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Empty state -->
      <div v-if="taskStore.tasks.length === 0 && !taskStore.loading" class="text-center py-16">
        <p class="text-lg font-medium text-foreground">Welcome to Team ADO Tool</p>
        <p class="text-sm text-muted-foreground mt-1">Create your first task to get started</p>
        <Button class="mt-4" @click="goToTasks">
          Get Started
        </Button>
      </div>
    </div>
  </ScrollArea>
</template>
