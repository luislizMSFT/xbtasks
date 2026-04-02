<script setup lang="ts">
import { onMounted, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { cn } from '@/lib/utils'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import AdoBadge from '@/components/ui/AdoBadge.vue'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  AlertTriangle,
  Octagon,
  GitPullRequest,
  ThumbsUp,
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

function markViewed(pr: PR) { pr.unread = false }

function prStatusClasses(status: PR['status']) {
  switch (status) {
    case 'active': return 'bg-blue-500/15 text-blue-600 dark:text-blue-400 border-blue-500/20'
    case 'draft': return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
    case 'completed': return 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
    case 'abandoned': return 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20'
    default: return ''
  }
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
    <!-- Left column: Tasks -->
    <ScrollArea class="flex-1 h-full">
      <div class="px-4 py-3 space-y-3">
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
          <Badge v-if="unreadCount > 0" variant="default" class="text-[10px] h-5 px-1.5">
            {{ unreadCount }}
          </Badge>
        </div>

        <!-- Needs Your Review -->
        <div v-if="reviewRequests.length > 0">
          <span class="text-[10px] font-medium text-muted-foreground/60 uppercase tracking-wider">Needs Your Review</span>
          <Card class="mt-1.5 overflow-hidden">
            <CardContent class="p-0">
              <div
                v-for="pr in reviewRequests"
                :key="pr.id"
                :class="cn(
                  'flex items-center gap-2 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0',
                  pr.unread && 'bg-blue-500/[0.03]'
                )"
                @click="markViewed(pr)"
              >
                <span v-if="pr.unread" class="w-1.5 h-1.5 rounded-full bg-blue-500 shrink-0" />
                <span v-else class="w-1.5 shrink-0" />
                <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
                <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize', prStatusClasses(pr.status))">
                  {{ pr.status }}
                </Badge>
                <div class="flex items-center gap-1 text-[11px] tabular-nums shrink-0">
                  <ThumbsUp v-if="pr.votes.up > 0" :size="10" class="text-emerald-500" />
                  <span v-if="pr.votes.up > 0" class="text-emerald-600 dark:text-emerald-400">{{ pr.votes.up }}</span>
                </div>
                <span class="text-[10px] text-muted-foreground/50 tabular-nums shrink-0">{{ pr.updatedAt }}</span>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Your PRs -->
        <div v-if="myPRs.length > 0">
          <span class="text-[10px] font-medium text-muted-foreground/60 uppercase tracking-wider">Your PRs</span>
          <Card class="mt-1.5 overflow-hidden">
            <CardContent class="p-0">
              <div
                v-for="pr in myPRs"
                :key="pr.id"
                :class="cn(
                  'flex items-center gap-2 px-3 py-2 cursor-pointer transition-colors hover:bg-muted/50 border-b border-border/50 last:border-b-0',
                  pr.unread && 'bg-blue-500/[0.03]'
                )"
                @click="markViewed(pr)"
              >
                <span v-if="pr.unread" class="w-1.5 h-1.5 rounded-full bg-blue-500 shrink-0" />
                <span v-else class="w-1.5 shrink-0" />
                <span class="text-sm text-foreground flex-1 truncate">{{ pr.title }}</span>
                <Badge variant="outline" :class="cn('text-[10px] px-1.5 py-0 capitalize', prStatusClasses(pr.status))">
                  {{ pr.status }}
                </Badge>
                <div class="flex -space-x-1.5 shrink-0">
                  <div
                    v-for="reviewer in pr.reviewers.slice(0, 3)"
                    :key="reviewer"
                    class="w-5 h-5 rounded-full bg-muted border border-background flex items-center justify-center"
                    :title="reviewer"
                  >
                    <span class="text-[8px] font-medium text-muted-foreground">{{ getInitials(reviewer) }}</span>
                  </div>
                </div>
                <span class="text-[10px] text-muted-foreground/50 tabular-nums shrink-0">{{ pr.updatedAt }}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </ScrollArea>
  </div>
</template>
