<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { cn } from '@/lib/utils'
import { adoTypeIcon, adoTypeColor, prStatusClasses } from '@/lib/styles'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import {
  ArrowLeft, ExternalLink,
  GitPullRequest, Play, CheckCircle2, XCircle, CalendarDays,
  RefreshCw, GitBranch, AlertTriangle,
  SquareCheckBig, GitMerge,
} from 'lucide-vue-next'

const router = useRouter()

// ── Mock data ──
const currentUser = 'Luis L.'
const now = new Date()

interface MockTask {
  id: number; title: string; status: string; priority: string
  projectName: string; adoId: string; adoType: string; adoState: string
  updatedAt: string; dueDate: string; isPersonal: boolean
  pendingSync: boolean; blockedReason?: string
}

interface MockPR {
  id: number; title: string; prNumber: number; repo: string
  status: string; votes: number; sourceBranch: string
  createdBy: string; updatedAt: string
}

const mockTasks: MockTask[] = [
  { id: 1, title: 'Implement user authentication flow with SSO', status: 'in_progress', priority: 'P0', projectName: 'Platform Modernization', adoId: '48291', adoType: 'User Story', adoState: 'Active', updatedAt: '2026-04-08T00:30:00Z', dueDate: '2026-04-10', isPersonal: false, pendingSync: true },
  { id: 2, title: 'Fix dashboard loading performance', status: 'in_progress', priority: 'P1', projectName: 'Developer Experience', adoId: '48350', adoType: 'Bug', adoState: 'Active', updatedAt: '2026-04-07T22:00:00Z', dueDate: '', isPersonal: false, pendingSync: false },
  { id: 3, title: 'Write API migration guide', status: 'in_review', priority: 'P2', projectName: 'Platform Modernization', adoId: '48400', adoType: 'Task', adoState: 'Active', updatedAt: '2026-04-07T18:00:00Z', dueDate: '2026-04-12', isPersonal: false, pendingSync: false },
  { id: 4, title: 'Draft team retro notes', status: 'todo', priority: 'P2', projectName: '', adoId: '', adoType: '', adoState: '', updatedAt: '2026-04-07T14:00:00Z', dueDate: '2026-04-09', isPersonal: true, pendingSync: false },
  { id: 5, title: 'Update Grafana alerts for new service', status: 'blocked', priority: 'P1', projectName: 'Developer Experience', adoId: '48360', adoType: 'Task', adoState: 'Active', updatedAt: '2026-04-07T11:00:00Z', dueDate: '', isPersonal: false, pendingSync: false, blockedReason: 'Waiting on infra team for access' },
  { id: 6, title: 'Review Sarah\'s PR on token handling', status: 'todo', priority: 'P1', projectName: '', adoId: '', adoType: '', adoState: '', updatedAt: '2026-04-07T09:00:00Z', dueDate: '2026-04-08', isPersonal: true, pendingSync: false },
  { id: 7, title: 'Add telemetry to auth flow', status: 'todo', priority: 'P2', projectName: 'Platform Modernization', adoId: '48410', adoType: 'Task', adoState: 'New', updatedAt: '2026-04-06T16:00:00Z', dueDate: '', isPersonal: false, pendingSync: false },
  { id: 8, title: 'Prepare demo for sprint review', status: 'todo', priority: 'P1', projectName: '', adoId: '', adoType: '', adoState: '', updatedAt: '2026-04-06T12:00:00Z', dueDate: '2026-04-10', isPersonal: true, pendingSync: false },
]

const mockPRs: MockPR[] = [
  { id: 100, title: 'feat: add Azure AD provider', prNumber: 342, repo: 'platform-api', status: 'active', votes: 1, sourceBranch: 'feat/auth-refresh', createdBy: 'Luis L.', updatedAt: '2026-04-07T23:00:00Z' },
  { id: 101, title: 'feat: token refresh middleware', prNumber: 343, repo: 'platform-api', status: 'draft', votes: 0, sourceBranch: 'feat/token-refresh', createdBy: 'Luis L.', updatedAt: '2026-04-07T20:00:00Z' },
  { id: 102, title: 'fix: resolve dashboard perf regression', prNumber: 156, repo: 'xb-frontend', status: 'active', votes: 2, sourceBranch: 'fix/dashboard-perf', createdBy: 'Luis L.', updatedAt: '2026-04-07T19:00:00Z' },
  { id: 200, title: 'refactor: extract auth middleware', prNumber: 341, repo: 'platform-api', status: 'active', votes: 0, sourceBranch: 'refactor/auth-middleware', createdBy: 'Sarah K.', updatedAt: '2026-04-07T21:00:00Z' },
  { id: 201, title: 'chore: upgrade deps to latest', prNumber: 155, repo: 'xb-frontend', status: 'active', votes: 0, sourceBranch: 'chore/dep-upgrade', createdBy: 'Mike R.', updatedAt: '2026-04-07T17:00:00Z' },
]

const mockPipelines = [
  { id: 1, name: 'CI Build', result: 'succeeded', sourceBranch: 'feat/auth-refresh', finishTime: '2026-04-07T23:15:00Z', url: '#' },
  { id: 2, name: 'Integration Tests', result: 'failed', sourceBranch: 'feat/auth-refresh', finishTime: '2026-04-07T23:20:00Z', url: '#' },
  { id: 3, name: 'CI Build', result: 'succeeded', sourceBranch: 'fix/dashboard-perf', finishTime: '2026-04-07T19:10:00Z', url: '#' },
  { id: 4, name: 'Deploy Preview', result: 'succeeded', sourceBranch: 'fix/dashboard-perf', finishTime: '2026-04-07T19:25:00Z', url: '#' },
]

// ── Computed ──
const focusTasks = computed(() => mockTasks.filter(t => t.status === 'in_progress' || t.status === 'in_review'))
const blockedTasks = computed(() => mockTasks.filter(t => t.status === 'blocked'))
const myPRs = computed(() => mockPRs.filter(pr => pr.createdBy === currentUser))
const reviewPRs = computed(() => mockPRs.filter(pr => pr.createdBy !== currentUser))
const pendingSyncCount = computed(() => mockTasks.filter(t => t.pendingSync).length)
const dueSoonTasks = computed(() => mockTasks.filter(t => {
  if (!t.dueDate || t.status === 'done') return false
  const due = new Date(t.dueDate)
  const diff = (due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
  return diff >= 0 && diff <= 3
}))

const stats = computed(() => ({
  total: mockTasks.length,
  inProgress: mockTasks.filter(t => t.status === 'in_progress').length,
  inReview: mockTasks.filter(t => t.status === 'in_review').length,
  blocked: mockTasks.filter(t => t.status === 'blocked').length,
  done: 12, // simulated historical
  todo: mockTasks.filter(t => t.status === 'todo').length,
  personal: mockTasks.filter(t => t.isPersonal).length,
}))

// Recent tasks (sorted by updated)
const recentTasks = computed(() =>
  [...mockTasks]
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    .slice(0, 5)
)

// Helpers
function relativeTime(iso: string) {
  const diff = now.getTime() - new Date(iso).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
}

function pipelineIcon(result: string) {
  if (result === 'succeeded') return CheckCircle2
  if (result === 'failed') return XCircle
  return Play
}

function pipelineColor(result: string) {
  if (result === 'succeeded') return 'text-emerald-500'
  if (result === 'failed') return 'text-red-500'
  return 'text-amber-500'
}

function priorityColor(p: string) {
  switch (p) {
    case 'P0': return 'bg-red-500'
    case 'P1': return 'bg-orange-500'
    case 'P2': return 'bg-amber-500'
    default: return 'bg-zinc-400'
  }
}

function statusBgColor(status: string) {
  switch (status) {
    case 'in_progress': return 'bg-blue-500'
    case 'in_review': return 'bg-violet-500'
    case 'done': return 'bg-emerald-500'
    case 'blocked': return 'bg-red-500'
    case 'todo': return 'bg-zinc-400'
    default: return 'bg-zinc-400'
  }
}

function voteDisplay(votes: number) {
  if (votes >= 2) return '✅'
  if (votes === 1) return '👍'
  return ''
}
</script>

<template>
  <div class="h-screen w-screen bg-background text-foreground flex flex-col">
    <!-- Playground header -->
    <div class="shrink-0 border-b border-border bg-card px-4 py-2 flex items-center gap-3 titlebar-drag">
      <Button variant="ghost" size="sm" class="h-7 gap-1 titlebar-no-drag" @click="router.back()">
        <ArrowLeft :size="14" /> Back
      </Button>
      <span class="text-sm font-medium">Dashboard Redesign Playground</span>
      <span class="text-[10px] text-muted-foreground bg-muted px-2 py-0.5 rounded-full">Mock Data</span>
    </div>

    <ScrollArea class="flex-1 min-h-0">
      <div class="max-w-[1200px] mx-auto px-6 py-5">

        <!-- ╔══════════════════════════════════════════════════════════════╗ -->
        <!-- ║  SECTION 1: Top Bar — Greeting + Sync + Quick Actions      ║ -->
        <!-- ╚══════════════════════════════════════════════════════════════╝ -->
        <div class="flex items-start justify-between mb-6">
          <div>
            <h1 class="text-xl font-semibold text-foreground">Good evening, Luis</h1>
            <p class="text-sm text-muted-foreground mt-0.5">
              {{ stats.inProgress }} in progress · {{ stats.blocked }} blocked · {{ reviewPRs.length }} PR{{ reviewPRs.length !== 1 ? 's' : '' }} to review
            </p>
          </div>
          <!-- Sync cluster -->
          <div class="flex items-center gap-2 text-[11px]">
            <div class="flex items-center gap-1.5 px-2 py-1 rounded-md bg-muted/50 border border-border/50">
              <span class="size-1.5 rounded-full bg-emerald-500" />
              <span class="text-muted-foreground">Synced · 3m ago</span>
              <span v-if="pendingSyncCount" class="text-amber-600 font-medium ml-1">· {{ pendingSyncCount }} pending</span>
            </div>
            <Button variant="ghost" size="sm" class="h-7 w-7 p-0">
              <RefreshCw :size="13" class="text-muted-foreground" />
            </Button>
          </div>
        </div>

        <!-- ╔══════════════════════════════════════════════════════════════╗ -->
        <!-- ║  SECTION 2: Attention Bar — Urgent nudges                  ║ -->
        <!-- ╚══════════════════════════════════════════════════════════════╝ -->
        <div class="flex gap-3 mb-5 overflow-x-auto pb-1">
          <!-- Due soon -->
          <div v-if="dueSoonTasks.length" class="flex items-center gap-2 px-3 py-2 rounded-lg bg-amber-500/8 border border-amber-500/20 shrink-0">
            <CalendarDays :size="14" class="text-amber-600 shrink-0" />
            <span class="text-xs text-amber-700 dark:text-amber-400">
              <strong>{{ dueSoonTasks.length }}</strong> due within 3 days
            </span>
            <div class="flex gap-1 ml-1">
              <span v-for="t in dueSoonTasks.slice(0, 2)" :key="t.id" class="text-[10px] bg-amber-500/10 px-1.5 py-0.5 rounded text-amber-700 dark:text-amber-400 truncate max-w-[120px]">{{ t.title }}</span>
            </div>
          </div>
          <!-- Pipeline failure -->
          <div v-if="mockPipelines.some(p => p.result === 'failed')" class="flex items-center gap-2 px-3 py-2 rounded-lg bg-red-500/8 border border-red-500/20 shrink-0">
            <XCircle :size="14" class="text-red-500 shrink-0" />
            <span class="text-xs text-red-600 dark:text-red-400">
              <strong>Integration Tests</strong> failed on feat/auth-refresh
            </span>
          </div>
          <!-- PR approval ready -->
          <div v-if="myPRs.some(pr => pr.votes >= 2)" class="flex items-center gap-2 px-3 py-2 rounded-lg bg-emerald-500/8 border border-emerald-500/20 shrink-0">
            <GitMerge :size="14" class="text-emerald-600 shrink-0" />
            <span class="text-xs text-emerald-700 dark:text-emerald-400">
              PR #156 has 2 approvals — ready to merge
            </span>
          </div>
        </div>

        <!-- ╔══════════════════════════════════════════════════════════════╗ -->
        <!-- ║  SECTION 3: Main Content — simple 2-column like current     ║ -->
        <!-- ╚══════════════════════════════════════════════════════════════╝ -->

        <!-- Compact stats line -->
        <div class="flex items-center gap-4 mb-5 text-sm text-muted-foreground">
          <span><strong class="tabular-nums text-foreground">{{ stats.inProgress }}</strong> in progress</span>
          <span><strong class="tabular-nums text-foreground">{{ stats.blocked }}</strong> blocked</span>
          <span><strong class="tabular-nums text-foreground">{{ stats.done }}</strong> done</span>
          <span class="text-muted-foreground/50">of {{ stats.total + stats.done }}</span>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
          <!-- Left column (3/5): Tasks -->
          <div class="lg:col-span-3 space-y-6">

            <!-- Today's Focus -->
            <div>
              <h2 class="text-sm font-semibold text-foreground mb-3">Today's Focus</h2>
              <div v-if="focusTasks.length > 0" class="rounded-lg overflow-hidden border border-border">
                <div v-for="(task, idx) in focusTasks" :key="task.id"
                  class="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
                  :class="idx < focusTasks.length - 1 ? 'border-b border-border' : ''"
                >
                  <span :class="cn('size-2 rounded-full shrink-0', priorityColor(task.priority))" />
                  <component v-if="task.adoType" :is="adoTypeIcon(task.adoType)" :size="14" :class="adoTypeColor(task.adoType)" class="shrink-0" />
                  <SquareCheckBig v-else :size="14" class="text-primary/60 shrink-0" />
                  <span class="text-sm flex-1 truncate text-foreground">{{ task.title }}</span>
                  <span v-if="task.isPersonal" class="text-[8px] px-1 py-0.5 rounded bg-primary/8 text-primary/70 border border-primary/10 shrink-0">personal</span>
                  <span v-if="task.pendingSync" class="size-1.5 rounded-full bg-amber-500 shrink-0" title="Pending sync" />
                  <Badge variant="outline" class="text-[9px] h-4 px-1.5 shrink-0" :class="task.status === 'in_review' ? 'border-violet-500/30 text-violet-600' : 'border-blue-500/30 text-blue-600'">
                    {{ task.status === 'in_review' ? 'Review' : 'Active' }}
                  </Badge>
                  <span v-if="task.dueDate" class="text-[10px] tabular-nums shrink-0"
                    :class="new Date(task.dueDate) <= new Date(now.getTime() + 2 * 86400000) ? 'text-amber-600 font-medium' : 'text-muted-foreground'">
                    {{ new Date(task.dueDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) }}
                  </span>
                </div>
              </div>
              <p v-else class="text-sm text-muted-foreground">No tasks in progress — pick something to work on.</p>
            </div>

            <!-- Recent Activity -->
            <div>
              <h2 class="text-sm font-semibold text-foreground mb-3">Recent Activity</h2>
              <div class="rounded-lg overflow-hidden border border-border">
                <div v-for="(task, idx) in recentTasks" :key="task.id"
                  class="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
                  :class="idx < recentTasks.length - 1 ? 'border-b border-border' : ''"
                >
                  <span :class="cn('size-2 rounded-full shrink-0', statusBgColor(task.status))" />
                  <span class="text-sm flex-1 truncate text-foreground">{{ task.title }}</span>
                  <span class="text-xs tabular-nums text-muted-foreground">{{ relativeTime(task.updatedAt) }}</span>
                </div>
              </div>
            </div>

            <!-- Blocked -->
            <div v-if="blockedTasks.length > 0">
              <h2 class="text-sm font-semibold text-foreground mb-3">Blocked</h2>
              <div class="rounded-lg overflow-hidden border border-border border-l-2 border-l-red-500">
                <div v-for="task in blockedTasks" :key="task.id"
                  class="flex flex-col gap-1 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
                >
                  <div class="flex items-center gap-2">
                    <span :class="cn('size-2 rounded-full shrink-0', priorityColor(task.priority))" />
                    <span class="text-sm flex-1 truncate">{{ task.title }}</span>
                    <span v-if="task.adoId" class="text-[9px] text-blue-500/60 tabular-nums">#{{ task.adoId }}</span>
                  </div>
                  <p v-if="task.blockedReason" class="text-[11px] text-red-500/70 pl-4 italic">{{ task.blockedReason }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Right column (2/5): PRs + Pipelines -->
          <div class="lg:col-span-2 space-y-6">

            <!-- Needs Your Review -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <GitPullRequest :size="14" class="text-muted-foreground" />
                <h2 class="text-sm font-semibold text-foreground">Needs Your Review</h2>
                <Badge v-if="reviewPRs.length" variant="secondary" class="h-4 text-[10px] px-1.5 bg-violet-500/10 text-violet-600 border border-violet-500/20">{{ reviewPRs.length }}</Badge>
              </div>
              <div v-if="reviewPRs.length > 0" class="rounded-lg overflow-hidden border border-border">
                <div v-for="pr in reviewPRs" :key="`r-${pr.id}`"
                  class="flex flex-col gap-1 px-3 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors border-b border-border last:border-b-0"
                >
                  <div class="flex items-center gap-2">
                    <GitPullRequest :size="13" class="text-emerald-500 shrink-0" />
                    <span class="text-sm truncate flex-1">{{ pr.title }}</span>
                  </div>
                  <div class="flex items-center gap-2 pl-5 text-[11px] text-muted-foreground">
                    <span class="truncate">{{ pr.repo }}</span>
                    <span>·</span>
                    <GitBranch :size="11" class="shrink-0" />
                    <span class="truncate">{{ pr.sourceBranch }}</span>
                    <span>·</span>
                    <span class="shrink-0">{{ pr.createdBy }}</span>
                    <span class="ml-auto shrink-0 tabular-nums">{{ relativeTime(pr.updatedAt) }}</span>
                  </div>
                </div>
              </div>
              <p v-else class="text-[13px] text-muted-foreground">No PRs to review.</p>
            </div>

            <!-- Your Pull Requests -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <GitPullRequest :size="14" class="text-muted-foreground" />
                <h2 class="text-sm font-semibold text-foreground">Your Pull Requests</h2>
                <Badge v-if="myPRs.length" variant="secondary" class="text-[10px] px-1.5 py-0 h-4">{{ myPRs.length }}</Badge>
              </div>
              <div v-if="myPRs.length > 0" class="rounded-lg overflow-hidden border border-border">
                <div v-for="pr in myPRs" :key="`m-${pr.id}`"
                  class="flex flex-col gap-1 px-3 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors border-b border-border last:border-b-0"
                >
                  <div class="flex items-center gap-2">
                    <GitPullRequest :size="13" :class="pr.status === 'draft' ? 'text-zinc-400' : 'text-emerald-500'" class="shrink-0" />
                    <span class="text-sm truncate flex-1">{{ pr.title }}</span>
                    <span v-if="pr.votes" class="text-xs shrink-0">{{ voteDisplay(pr.votes) }}</span>
                  </div>
                  <div class="flex items-center gap-2 pl-5 text-[11px] text-muted-foreground">
                    <span class="truncate">{{ pr.repo }} #{{ pr.prNumber }}</span>
                    <span>·</span>
                    <GitBranch :size="11" class="shrink-0" />
                    <span class="truncate">{{ pr.sourceBranch }}</span>
                    <span class="ml-auto shrink-0 tabular-nums">{{ relativeTime(pr.updatedAt) }}</span>
                  </div>
                </div>
              </div>
              <p v-else class="text-[13px] text-muted-foreground">No active PRs.</p>
            </div>

            <!-- Pipelines -->
            <div>
              <div class="flex items-center gap-2 mb-3">
                <Play :size="14" class="text-muted-foreground" />
                <h2 class="text-sm font-semibold text-foreground">Pipelines</h2>
                <Badge v-if="mockPipelines.length" variant="secondary" class="text-[10px] px-1.5 py-0 h-4">{{ mockPipelines.length }}</Badge>
              </div>
              <div v-if="mockPipelines.length > 0" class="rounded-lg overflow-hidden border border-border">
                <div v-for="(p, idx) in mockPipelines" :key="p.id"
                  class="flex items-center gap-2.5 px-3 py-2 cursor-pointer hover:bg-muted/50 transition-colors"
                  :class="idx < mockPipelines.length - 1 ? 'border-b border-border' : ''"
                >
                  <component :is="pipelineIcon(p.result)" :size="13" :class="pipelineColor(p.result)" class="shrink-0" />
                  <span class="text-xs font-medium truncate">{{ p.name }}</span>
                  <span class="text-[10px] text-muted-foreground truncate flex-1">{{ p.sourceBranch }}</span>
                  <Badge
                    :variant="p.result === 'succeeded' ? 'secondary' : p.result === 'failed' ? 'destructive' : 'outline'"
                    class="text-[9px] px-1.5 h-4 shrink-0 capitalize"
                  >{{ p.result }}</Badge>
                  <span class="text-[10px] tabular-nums text-muted-foreground shrink-0">{{ relativeTime(p.finishTime) }}</span>
                </div>
              </div>
              <p v-else class="text-[13px] text-muted-foreground">No pipeline runs for your PRs.</p>
            </div>

          </div>
        </div>

      </div>
    </ScrollArea>
  </div>
</template>
