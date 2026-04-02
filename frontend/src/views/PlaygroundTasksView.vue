<script setup lang="ts">
import { ref, computed } from 'vue'
import { cn } from '@/lib/utils'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import {
  Circle,
  CircleDot,
  Eye,
  CheckCircle2,
  Octagon,
  GitPullRequest,
  ExternalLink,
  Activity,
  Bug,
  CheckSquare,
  BookOpen,
  Landmark,
} from 'lucide-vue-next'

interface Task {
  id: number
  title: string
  status: string
  priority: string
  tags: string
  adoId: string
  adoType: string
  subtasksDone: number
  subtasksTotal: number
  updatedAt: string
}

const mockTasks: Task[] = [
  { id: 1, title: 'Fix auth redirect loop', status: 'in_progress', priority: 'P0', tags: 'bug,urgent', adoId: 'ADO-48291', adoType: 'bug', subtasksDone: 2, subtasksTotal: 4, updatedAt: '30m ago' },
  { id: 2, title: 'Build sidebar navigation', status: 'in_progress', priority: 'P1', tags: 'ui', adoId: '', adoType: '', subtasksDone: 0, subtasksTotal: 0, updatedAt: '1h ago' },
  { id: 3, title: 'Review PR #234', status: 'in_review', priority: 'P1', tags: 'review', adoId: 'ADO-48350', adoType: 'task', subtasksDone: 0, subtasksTotal: 0, updatedAt: '2h ago' },
  { id: 4, title: 'Add rate limiting to gateway', status: 'todo', priority: 'P1', tags: 'backend', adoId: 'ADO-48400', adoType: 'task', subtasksDone: 0, subtasksTotal: 3, updatedAt: '3h ago' },
  { id: 5, title: 'Design system tokens for dark mode', status: 'todo', priority: 'P2', tags: 'design', adoId: '', adoType: '', subtasksDone: 0, subtasksTotal: 0, updatedAt: '5h ago' },
  { id: 6, title: 'Migrate auth to MSAL v3', status: 'blocked', priority: 'P0', tags: 'auth,blocked', adoId: 'ADO-48100', adoType: 'bug', subtasksDone: 1, subtasksTotal: 2, updatedAt: '1d ago' },
  { id: 7, title: 'Write E2E tests for task CRUD', status: 'done', priority: 'P2', tags: 'testing', adoId: '', adoType: '', subtasksDone: 5, subtasksTotal: 5, updatedAt: '2d ago' },
  { id: 8, title: 'Update team onboarding docs', status: 'done', priority: 'P3', tags: 'docs', adoId: 'ADO-47900', adoType: 'user_story', subtasksDone: 3, subtasksTotal: 3, updatedAt: '3d ago' },
]

const statusConfig: Record<string, { label: string; color: string; dotClass: string; borderClass: string; icon: typeof Circle }> = {
  in_progress: { label: 'In Progress', color: 'text-blue-500', dotClass: 'bg-blue-500', borderClass: 'border-blue-500', icon: CircleDot },
  in_review:   { label: 'In Review',   color: 'text-violet-500', dotClass: 'bg-violet-500', borderClass: 'border-violet-500', icon: Eye },
  todo:        { label: 'To Do',       color: 'text-zinc-400', dotClass: 'bg-zinc-400', borderClass: 'border-zinc-400', icon: Circle },
  blocked:     { label: 'Blocked',     color: 'text-red-500', dotClass: 'bg-red-500', borderClass: 'border-red-500', icon: Octagon },
  done:        { label: 'Done',        color: 'text-emerald-500', dotClass: 'bg-emerald-500', borderClass: 'border-emerald-500', icon: CheckCircle2 },
}

const priorityConfig: Record<string, { classes: string }> = {
  P0: { classes: 'border-red-500/30 text-red-600 dark:text-red-400 bg-red-500/5' },
  P1: { classes: 'border-orange-500/30 text-orange-600 dark:text-orange-400 bg-orange-500/5' },
  P2: { classes: 'border-amber-500/30 text-amber-600 dark:text-amber-400 bg-amber-500/5' },
  P3: { classes: 'border-zinc-400/30 text-zinc-500 bg-zinc-400/5' },
}

const adoTypeConfig: Record<string, { icon: typeof Bug; color: string }> = {
  bug:        { icon: Bug, color: 'text-red-500' },
  task:       { icon: CheckSquare, color: 'text-blue-500' },
  user_story: { icon: BookOpen, color: 'text-green-500' },
}

const mockSubtasks = [
  { id: 1, title: 'Identify refresh token edge case', done: true },
  { id: 2, title: 'Add token expiry middleware', done: true },
  { id: 3, title: 'Update error handling in auth callback', done: false },
  { id: 4, title: 'Add E2E test for token refresh', done: false },
]

const mockPRs = [
  { id: 234, title: 'Fix auth redirect loop', status: 'active', additions: 400, deletions: 120 },
  { id: 240, title: 'Add token expiry middleware', status: 'merged', additions: 80, deletions: 10 },
]

const mockComments = [
  { avatar: 'LL', author: 'You', text: 'Checked logs \u2014 refresh token expiring 5 min early. Clock skew?', time: '2h ago' },
  { avatar: 'AK', author: 'Alex K', text: 'Root cause confirmed as clock skew in token validation.', time: '1h ago' },
  { avatar: 'LL', author: 'You', text: 'Fix deployed to staging, monitoring for recurrence.', time: '30m ago' },
]

const mockActivity = [
  { event: 'Status changed to In Progress', time: '2d ago' },
  { event: 'Linked to ADO Bug #48291', time: '2d ago' },
  { event: 'Subtask completed: Identify refresh token edge case', time: '1d ago' },
  { event: 'PR #234 opened', time: '6h ago' },
  { event: 'Subtask completed: Add token expiry middleware', time: '3h ago' },
]

const previewTaskId = ref<number>(1)
const previewTask = computed(() => mockTasks.find(t => t.id === previewTaskId.value) ?? mockTasks[0])
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <div class="px-6 py-6 max-w-[1400px] mx-auto">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-xl font-semibold text-foreground">Task Layout Playground</h1>
          <p class="text-sm text-muted-foreground mt-0.5">Todoist Hybrid \u2014 compact list + rich preview</p>
        </div>
        <Badge variant="outline" class="text-xs text-muted-foreground">
          {{ mockTasks.length }} tasks
        </Badge>
      </div>

      <div class="flex gap-0 rounded-xl border border-border bg-card overflow-hidden" style="height: 640px">
        <!-- Left: Task List -->
        <div class="w-[55%] border-r border-border flex flex-col">
          <div class="px-4 py-3 border-b border-border bg-muted/30">
            <span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              All Tasks \u00b7 {{ mockTasks.length }}
            </span>
          </div>
          <ScrollArea class="flex-1">
            <div>
              <div
                v-for="task in mockTasks"
                :key="task.id"
                @click="previewTaskId = task.id"
                :class="cn(
                  'group flex items-center gap-2.5 px-3 py-2.5 cursor-pointer transition-all border-b border-border/50',
                  'hover:bg-muted/40',
                  previewTaskId === task.id ? 'bg-blue-500/[0.06]' : ''
                )"
              >
                <div
                  :class="cn(
                    'w-[3px] self-stretch rounded-full shrink-0',
                    task.priority === 'P0' ? 'bg-red-500' :
                    task.priority === 'P1' ? 'bg-orange-500' :
                    task.priority === 'P2' ? 'bg-amber-500' :
                    'bg-zinc-300 dark:bg-zinc-600'
                  )"
                />
                <div
                  :class="cn(
                    'size-4 rounded-full border-[1.5px] shrink-0 flex items-center justify-center',
                    statusConfig[task.status].borderClass,
                    task.status === 'done' ? statusConfig[task.status].dotClass : ''
                  )"
                >
                  <CheckCircle2 v-if="task.status === 'done'" :size="8" class="text-white" :stroke-width="3" />
                </div>
                <span
                  :class="cn(
                    'text-[13px] font-medium truncate flex-1',
                    task.status === 'done' ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground'
                  )"
                >
                  {{ task.title }}
                </span>
                <component
                  v-if="task.adoId && adoTypeConfig[task.adoType]"
                  :is="adoTypeConfig[task.adoType].icon"
                  :size="12"
                  :class="adoTypeConfig[task.adoType].color"
                  class="shrink-0 opacity-60"
                />
                <Landmark v-else-if="task.adoId" :size="12" class="text-blue-500 shrink-0 opacity-60" />
                <Circle v-else :size="12" class="text-muted-foreground/20 shrink-0" />
                <span class="text-[10px] text-muted-foreground/40 tabular-nums shrink-0">{{ task.updatedAt }}</span>
              </div>
            </div>
          </ScrollArea>
        </div>

        <!-- Right: Preview -->
        <div class="w-[45%] flex flex-col bg-card">
          <div class="px-5 py-4 border-b border-border">
            <div class="flex items-start justify-between gap-3">
              <div class="flex-1 min-w-0">
                <h2 class="text-base font-semibold text-foreground leading-snug">{{ previewTask.title }}</h2>
                <div class="flex items-center gap-2 mt-2 flex-wrap">
                  <Badge variant="outline" :class="cn('text-[10px] font-bold uppercase tracking-wider px-1.5 py-0 h-5', priorityConfig[previewTask.priority]?.classes)">
                    {{ previewTask.priority }}
                  </Badge>
                  <Badge variant="outline" :class="cn('text-[10px] font-medium gap-1 h-5', statusConfig[previewTask.status]?.color)">
                    <component :is="statusConfig[previewTask.status].icon" :size="11" :stroke-width="2" />
                    {{ statusConfig[previewTask.status].label }}
                  </Badge>
                  <Badge v-if="previewTask.adoId" variant="outline" class="text-[10px] font-medium gap-1 h-5 bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20">
                    <Landmark :size="10" />
                    {{ previewTask.adoId }}
                  </Badge>
                </div>
              </div>
              <button class="p-1.5 rounded-md hover:bg-muted transition-colors text-muted-foreground">
                <ExternalLink :size="14" />
              </button>
            </div>
            <div v-if="previewTask.subtasksTotal > 0" class="mt-3 flex items-center gap-2.5">
              <div class="h-1.5 flex-1 rounded-full bg-muted overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :class="previewTask.subtasksDone === previewTask.subtasksTotal ? 'bg-emerald-500' : 'bg-blue-500'"
                  :style="{ width: `${(previewTask.subtasksDone / previewTask.subtasksTotal) * 100}%` }"
                />
              </div>
              <span class="text-[11px] text-muted-foreground tabular-nums">{{ previewTask.subtasksDone }}/{{ previewTask.subtasksTotal }}</span>
            </div>
          </div>

          <ScrollArea class="flex-1">
            <div class="px-5 py-4 space-y-5">
              <!-- Subtasks -->
              <div>
                <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">Subtasks</h3>
                <div class="space-y-1">
                  <div v-for="sub in mockSubtasks" :key="sub.id" class="flex items-center gap-2.5 py-1.5 px-2 rounded-md hover:bg-muted/40 transition-colors">
                    <div :class="cn('size-4 rounded-[4px] border-[1.5px] shrink-0 flex items-center justify-center', sub.done ? 'bg-emerald-500 border-emerald-500' : 'border-muted-foreground/30')">
                      <CheckCircle2 v-if="sub.done" :size="10" class="text-white" :stroke-width="3" />
                    </div>
                    <span :class="cn('text-[13px]', sub.done ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground')">{{ sub.title }}</span>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- PRs -->
              <div>
                <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">Pull Requests</h3>
                <div class="space-y-1.5">
                  <div v-for="pr in mockPRs" :key="pr.id" class="flex items-center gap-2.5 py-2 px-3 rounded-lg border border-border hover:bg-muted/30 transition-colors cursor-pointer">
                    <GitPullRequest :size="14" :class="pr.status === 'merged' ? 'text-violet-500' : 'text-emerald-500'" />
                    <div class="flex-1 min-w-0">
                      <span class="text-[13px] font-medium text-foreground truncate block">{{ pr.title }}</span>
                      <span class="text-[10px] text-muted-foreground">#{{ pr.id }} \u00b7 <span class="text-emerald-500">+{{ pr.additions }}</span> <span class="text-red-500">-{{ pr.deletions }}</span></span>
                    </div>
                    <Badge variant="outline" :class="cn('text-[10px] h-5 capitalize', pr.status === 'merged' ? 'bg-violet-500/10 text-violet-500 border-violet-500/20' : 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20')">{{ pr.status }}</Badge>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Comments -->
              <div>
                <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">Comments</h3>
                <div class="space-y-3">
                  <div v-for="(comment, idx) in mockComments" :key="idx" class="flex gap-2.5">
                    <div :class="cn('size-7 rounded-full shrink-0 flex items-center justify-center text-[10px] font-bold', comment.author === 'You' ? 'bg-blue-500/15 text-blue-600 dark:text-blue-400' : 'bg-muted text-muted-foreground')">{{ comment.avatar }}</div>
                    <div class="flex-1 min-w-0">
                      <div class="flex items-center gap-2">
                        <span class="text-[12px] font-semibold text-foreground">{{ comment.author }}</span>
                        <span class="text-[10px] text-muted-foreground/50">{{ comment.time }}</span>
                      </div>
                      <p class="text-[13px] text-muted-foreground leading-relaxed mt-0.5">{{ comment.text }}</p>
                    </div>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Activity -->
              <div>
                <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">Activity</h3>
                <div class="space-y-0">
                  <div v-for="(item, idx) in mockActivity" :key="idx" class="flex items-start gap-2.5 py-1.5 relative">
                    <div v-if="idx < mockActivity.length - 1" class="absolute left-[7px] top-[18px] w-px h-[calc(100%-6px)] bg-border" />
                    <div class="size-[15px] rounded-full bg-muted border border-border shrink-0 flex items-center justify-center z-10">
                      <Activity :size="8" class="text-muted-foreground" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <span class="text-[12px] text-muted-foreground">{{ item.event }}</span>
                    </div>
                    <span class="text-[10px] text-muted-foreground/40 tabular-nums shrink-0">{{ item.time }}</span>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Details -->
              <div>
                <h3 class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-2.5">Details</h3>
                <div class="grid grid-cols-2 gap-y-2 gap-x-4 text-[12px]">
                  <div class="text-muted-foreground/60">Priority</div>
                  <div class="text-foreground font-medium">{{ previewTask.priority }}</div>
                  <div class="text-muted-foreground/60">Status</div>
                  <div class="text-foreground font-medium">{{ statusConfig[previewTask.status].label }}</div>
                  <div class="text-muted-foreground/60">Tags</div>
                  <div class="flex gap-1 flex-wrap">
                    <Badge v-for="tag in previewTask.tags.split(',')" :key="tag" variant="secondary" class="text-[10px] px-1.5 py-0 h-[16px] rounded">{{ tag }}</Badge>
                  </div>
                  <div class="text-muted-foreground/60">Updated</div>
                  <div class="text-foreground font-medium">{{ previewTask.updatedAt }}</div>
                  <template v-if="previewTask.adoId">
                    <div class="text-muted-foreground/60">ADO Item</div>
                    <div class="text-blue-600 dark:text-blue-400 font-medium">{{ previewTask.adoId }}</div>
                  </template>
                </div>
              </div>
            </div>
          </ScrollArea>
        </div>
      </div>
    </div>
  </ScrollArea>
</template>
