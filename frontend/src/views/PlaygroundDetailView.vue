<script setup lang="ts">
import { ref, computed } from 'vue'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import TagChip from '@/components/ui/TagChip.vue'
import {
  X,
  Plus,
  Check,
  Circle,
  ChevronDown,
  ChevronRight,
  GitPullRequest,
  ExternalLink,
  Bug,
  MessageSquare,
  Clock,
  Folder,
  CalendarDays,
  Trash2,
  Link as LinkIcon,
  CircleDot,
  FileEdit,
  Send,
  Activity,
  User,
} from 'lucide-vue-next'

// ─── Mock Data ──────────────────────────────────────────────

const task = ref({
  title: 'Fix auth redirect loop',
  status: 'in_progress',
  priority: 'P0',
  description:
    'The OAuth callback is redirecting in a loop when the token expires mid-session. Need to check the refresh token logic and handle edge cases.',
  tags: ['bug', 'urgent', 'auth'],
  adoId: 'ADO-48291',
  adoType: 'bug',
  adoTitle: 'Auth redirect loop on token expiry',
  adoState: 'Active',
  adoUrl: '#',
  project: 'Team ADO Tool',
  dueDate: '2026-04-10',
  createdAt: 'Mar 31, 2026 10:00 AM',
  updatedAt: 'Apr 2, 2026 2:30 PM',
})

const subtasks = ref([
  { id: 1, title: 'Identify refresh token edge case', done: true },
  { id: 2, title: 'Add token expiry middleware', done: true },
  { id: 3, title: 'Update error handling in auth callback', done: false },
  { id: 4, title: 'Add E2E test for token refresh', done: false },
])

const prs = ref([
  {
    id: 234,
    title: 'Fix auth redirect loop',
    status: 'active',
    repo: 'xb-services',
    branch: 'fix/auth-redirect → main',
    reviewers: [
      { name: 'Alex', vote: 'approved' },
      { name: 'Sam', vote: 'pending' },
    ],
  },
  {
    id: 240,
    title: 'Add token expiry middleware',
    status: 'draft',
    repo: 'xb-services',
    branch: 'feat/token-middleware → main',
    reviewers: [],
  },
])

const comments = ref([
  {
    id: 1,
    type: 'personal',
    text: 'Checked logs — the refresh token is expiring 5 min before expected. Might be clock skew.',
    time: '2h ago',
  },
  {
    id: 2,
    type: 'ado',
    text: 'Updated work item: root cause identified as clock skew in token validation.',
    time: '1h ago',
  },
  {
    id: 3,
    type: 'personal',
    text: 'Fix deployed to staging, monitoring for recurrence.',
    time: '30m ago',
  },
])

const timeline = ref([
  { type: 'created', text: 'Task created', time: 'Mar 31' },
  { type: 'ado', text: 'Linked to ADO-48291', time: 'Mar 31' },
  { type: 'status', text: 'Status → In Progress', time: 'Apr 1' },
  { type: 'pr', text: 'PR #234 opened', time: 'Apr 1' },
  { type: 'comment', text: 'Personal note added', time: 'Apr 2' },
  { type: 'pr', text: 'PR #240 opened (draft)', time: 'Apr 2' },
  { type: 'ado', text: 'ADO synced — state: Active', time: 'Apr 2' },
])

// ─── Computed ───────────────────────────────────────────────

const subtasksDone = computed(() => subtasks.value.filter((s) => s.done).length)
const subtasksTotal = computed(() => subtasks.value.length)
const subtaskProgress = computed(() =>
  subtasksTotal.value > 0
    ? (subtasksDone.value / subtasksTotal.value) * 100
    : 0,
)

const personalComments = computed(() =>
  comments.value.filter((c) => c.type === 'personal'),
)
const adoComments = computed(() =>
  comments.value.filter((c) => c.type === 'ado'),
)

// ─── UI State ───────────────────────────────────────────────

const showPrsA = ref(true)
const showPrsB = ref(true)
const showPrsC = ref(true)
const newComment = ref('')

// ─── Helpers ────────────────────────────────────────────────

function prStatusDot(status: string) {
  return status === 'active'
    ? 'bg-emerald-500'
    : status === 'draft'
      ? 'bg-zinc-400'
      : 'bg-blue-500'
}

function prStatusLabel(status: string) {
  return status === 'active'
    ? 'Active'
    : status === 'draft'
      ? 'Draft'
      : status.charAt(0).toUpperCase() + status.slice(1)
}

function reviewerVoteColor(vote: string) {
  return vote === 'approved'
    ? 'bg-emerald-500/15 text-emerald-600 border-emerald-500/20'
    : 'bg-zinc-500/10 text-muted-foreground border-zinc-500/20'
}

function timelineIcon(type: string) {
  const map: Record<string, any> = {
    created: Plus,
    ado: LinkIcon,
    status: CircleDot,
    pr: GitPullRequest,
    comment: MessageSquare,
  }
  return map[type] || Circle
}

function timelineIconColor(type: string) {
  const map: Record<string, string> = {
    created: 'text-emerald-500',
    ado: 'text-blue-500',
    status: 'text-blue-500',
    pr: 'text-violet-500',
    comment: 'text-amber-500',
  }
  return map[type] || 'text-muted-foreground'
}

function commentAvatar(type: string) {
  return type === 'personal' ? 'You' : 'ADO'
}

function commentAvatarBg(type: string) {
  return type === 'personal'
    ? 'bg-blue-500/15 text-blue-600'
    : 'bg-orange-500/15 text-orange-600'
}
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <div class="px-6 py-6">
      <!-- Page Header -->
      <div class="mb-6">
        <h1 class="text-xl font-semibold text-foreground">
          Detail Panel Playground
        </h1>
        <p class="text-sm text-muted-foreground mt-0.5">
          Compare task detail layouts — 3 variants for design review
        </p>
      </div>

      <Tabs default-value="variant-a">
        <TabsList class="mb-6">
          <TabsTrigger value="variant-a">Linear Style</TabsTrigger>
          <TabsTrigger value="variant-b">Notion Style</TabsTrigger>
          <TabsTrigger value="variant-c">GitHub Issue</TabsTrigger>
        </TabsList>

        <!-- ═══════════════════════════════════════════════════════════
             VARIANT A — "Linear Style" (single column, top to bottom)
             ═══════════════════════════════════════════════════════════ -->
        <TabsContent value="variant-a">
          <Card class="max-w-[480px] py-0 gap-0">
            <div class="px-5 py-4 space-y-5">
              <!-- Header: Close + Status -->
              <div class="flex items-center justify-between">
                <StatusBadge :status="task.status" />
                <Button variant="ghost" size="icon-sm" class="h-7 w-7">
                  <X :size="16" />
                </Button>
              </div>

              <!-- Title -->
              <Input
                :model-value="task.title"
                class="text-lg font-semibold border-none shadow-none px-0 h-auto focus-visible:ring-0 focus-visible:border-none"
              />

              <!-- Subtasks -->
              <div class="space-y-2.5">
                <div class="flex items-center justify-between">
                  <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                    Subtasks
                  </span>
                  <span class="text-xs text-muted-foreground">
                    {{ subtasksDone }}/{{ subtasksTotal }}
                  </span>
                </div>
                <div class="h-[2px] w-full bg-muted rounded-full overflow-hidden">
                  <div
                    class="h-full bg-blue-500 rounded-full transition-all duration-300"
                    :style="{ width: `${subtaskProgress}%` }"
                  />
                </div>
                <div class="space-y-1">
                  <div
                    v-for="sub in subtasks"
                    :key="sub.id"
                    class="flex items-center gap-2.5 py-1 group"
                  >
                    <button
                      :class="cn(
                        'flex-shrink-0 w-[18px] h-[18px] rounded-full border-[1.5px] flex items-center justify-center transition-all',
                        sub.done
                          ? 'bg-blue-500 border-blue-500 text-white'
                          : 'border-border hover:border-blue-400',
                      )"
                    >
                      <Check v-if="sub.done" :size="11" :stroke-width="3" />
                    </button>
                    <span
                      :class="cn(
                        'text-sm leading-tight',
                        sub.done && 'line-through text-muted-foreground',
                      )"
                    >
                      {{ sub.title }}
                    </span>
                  </div>
                </div>
                <Button variant="ghost" size="sm" class="h-7 text-xs text-muted-foreground gap-1 px-1.5">
                  <Plus :size="14" />
                  Add subtask
                </Button>
              </div>

              <Separator />

              <!-- PRs — Collapsible -->
              <div class="space-y-2">
                <button
                  class="flex items-center gap-1.5 text-xs font-medium text-muted-foreground uppercase tracking-wide hover:text-foreground transition-colors w-full"
                  @click="showPrsA = !showPrsA"
                >
                  <component
                    :is="showPrsA ? ChevronDown : ChevronRight"
                    :size="14"
                  />
                  Pull Requests
                  <Badge variant="secondary" class="ml-auto text-[10px] h-4 px-1.5">
                    {{ prs.length }}
                  </Badge>
                </button>
                <div v-if="showPrsA" class="space-y-1.5">
                  <div
                    v-for="pr in prs"
                    :key="pr.id"
                    class="flex items-start gap-2.5 p-2.5 rounded-lg bg-muted/40 hover:bg-muted/70 transition-colors"
                  >
                    <div
                      :class="cn(
                        'w-2 h-2 rounded-full mt-1.5 flex-shrink-0',
                        prStatusDot(pr.status),
                      )"
                    />
                    <div class="flex-1 min-w-0 space-y-1">
                      <div class="flex items-center gap-1.5">
                        <span class="text-sm font-medium truncate">
                          {{ pr.title }}
                        </span>
                        <span class="text-[10px] text-muted-foreground">
                          #{{ pr.id }}
                        </span>
                      </div>
                      <div class="text-[11px] text-muted-foreground font-mono truncate">
                        {{ pr.repo }} · {{ pr.branch }}
                      </div>
                      <div v-if="pr.reviewers.length" class="flex items-center gap-1 mt-0.5">
                        <Badge
                          v-for="rev in pr.reviewers"
                          :key="rev.name"
                          variant="outline"
                          :class="cn('text-[10px] h-4 px-1.5', reviewerVoteColor(rev.vote))"
                        >
                          {{ rev.name }}
                          <Check v-if="rev.vote === 'approved'" :size="10" class="ml-0.5" />
                        </Badge>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Description -->
              <div class="space-y-2">
                <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                  Description
                </span>
                <Textarea
                  :model-value="task.description"
                  class="min-h-[80px] text-sm resize-none"
                />
              </div>

              <Separator />

              <!-- Comments -->
              <div class="space-y-3">
                <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                  Comments
                </span>
                <Tabs default-value="personal">
                  <TabsList class="h-8">
                    <TabsTrigger value="personal" class="text-xs h-6 px-2.5">
                      Personal
                      <Badge variant="secondary" class="ml-1 text-[10px] h-4 px-1">
                        {{ personalComments.length }}
                      </Badge>
                    </TabsTrigger>
                    <TabsTrigger value="ado" class="text-xs h-6 px-2.5">
                      ADO
                      <Badge variant="secondary" class="ml-1 text-[10px] h-4 px-1">
                        {{ adoComments.length }}
                      </Badge>
                    </TabsTrigger>
                  </TabsList>
                  <TabsContent value="personal" class="mt-3 space-y-3">
                    <div
                      v-for="c in personalComments"
                      :key="c.id"
                      class="flex gap-2.5"
                    >
                      <div
                        :class="cn(
                          'flex-shrink-0 w-7 h-7 rounded-full flex items-center justify-center text-[10px] font-semibold',
                          commentAvatarBg(c.type),
                        )"
                      >
                        {{ commentAvatar(c.type).charAt(0) }}
                      </div>
                      <div class="flex-1 min-w-0 space-y-0.5">
                        <p class="text-sm leading-relaxed">{{ c.text }}</p>
                        <span class="text-[11px] text-muted-foreground">{{ c.time }}</span>
                      </div>
                    </div>
                  </TabsContent>
                  <TabsContent value="ado" class="mt-3 space-y-3">
                    <div
                      v-for="c in adoComments"
                      :key="c.id"
                      class="flex gap-2.5"
                    >
                      <div
                        :class="cn(
                          'flex-shrink-0 w-7 h-7 rounded-full flex items-center justify-center text-[10px] font-semibold',
                          commentAvatarBg(c.type),
                        )"
                      >
                        {{ commentAvatar(c.type).charAt(0) }}
                      </div>
                      <div class="flex-1 min-w-0 space-y-0.5">
                        <p class="text-sm leading-relaxed">{{ c.text }}</p>
                        <span class="text-[11px] text-muted-foreground">{{ c.time }}</span>
                      </div>
                    </div>
                  </TabsContent>
                </Tabs>
                <div class="flex gap-2">
                  <Input
                    v-model="newComment"
                    placeholder="Add a comment..."
                    class="text-sm h-8"
                  />
                  <Button variant="ghost" size="icon-sm" class="flex-shrink-0">
                    <Send :size="14" />
                  </Button>
                </div>
              </div>

              <Separator />

              <!-- Timeline -->
              <div class="space-y-2.5">
                <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                  Timeline
                </span>
                <div class="relative space-y-0">
                  <div
                    v-for="(event, i) in timeline"
                    :key="i"
                    class="flex items-start gap-3 relative"
                  >
                    <div class="flex flex-col items-center">
                      <div
                        :class="cn(
                          'flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center bg-muted/60',
                          timelineIconColor(event.type),
                        )"
                      >
                        <component :is="timelineIcon(event.type)" :size="12" />
                      </div>
                      <div
                        v-if="i < timeline.length - 1"
                        class="w-px h-4 bg-border"
                      />
                    </div>
                    <div class="flex items-center gap-2 pb-1 pt-0.5">
                      <span class="text-sm">{{ event.text }}</span>
                      <span class="text-[11px] text-muted-foreground">{{ event.time }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Config Grid -->
              <div class="space-y-3">
                <div class="grid grid-cols-2 gap-3">
                  <div class="space-y-1.5">
                    <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                      Status
                    </label>
                    <Select :default-value="task.status">
                      <SelectTrigger class="h-8 text-xs">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="todo">To Do</SelectItem>
                        <SelectItem value="in_progress">In Progress</SelectItem>
                        <SelectItem value="in_review">In Review</SelectItem>
                        <SelectItem value="blocked">Blocked</SelectItem>
                        <SelectItem value="done">Done</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div class="space-y-1.5">
                    <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                      Priority
                    </label>
                    <Select :default-value="task.priority">
                      <SelectTrigger class="h-8 text-xs">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="P0">P0 — Critical</SelectItem>
                        <SelectItem value="P1">P1 — High</SelectItem>
                        <SelectItem value="P2">P2 — Medium</SelectItem>
                        <SelectItem value="P3">P3 — Low</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Tags
                  </label>
                  <div class="flex flex-wrap gap-1.5">
                    <TagChip v-for="tag in task.tags" :key="tag" :tag="tag" removable />
                    <Button variant="ghost" size="sm" class="h-5 text-[11px] px-1.5 text-muted-foreground">
                      <Plus :size="12" />
                    </Button>
                  </div>
                </div>
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Project
                  </label>
                  <div class="flex items-center gap-1.5 text-sm">
                    <Folder :size="14" class="text-muted-foreground" />
                    {{ task.project }}
                  </div>
                </div>
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Azure DevOps
                  </label>
                  <div class="flex items-center gap-2 p-2 rounded-md bg-muted/40">
                    <Bug :size="14" class="text-red-500 flex-shrink-0" />
                    <div class="flex-1 min-w-0">
                      <div class="flex items-center gap-1.5">
                        <Badge variant="outline" class="text-[10px] h-4 px-1.5 bg-blue-500/10 text-blue-600 border-blue-500/20">
                          {{ task.adoId }}
                        </Badge>
                        <Badge variant="outline" class="text-[10px] h-4 px-1.5">
                          {{ task.adoState }}
                        </Badge>
                      </div>
                      <p class="text-xs text-muted-foreground mt-0.5 truncate">
                        {{ task.adoTitle }}
                      </p>
                    </div>
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger as-child>
                          <Button variant="ghost" size="icon-sm" class="h-6 w-6 flex-shrink-0" as="a" :href="task.adoUrl">
                            <ExternalLink :size="12" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>Open in ADO</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  </div>
                </div>
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Due Date
                  </label>
                  <div class="flex items-center gap-1.5 text-sm">
                    <CalendarDays :size="14" class="text-muted-foreground" />
                    {{ task.dueDate }}
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Footer -->
              <div class="flex items-center justify-between pb-1">
                <div class="text-[11px] text-muted-foreground space-y-0.5">
                  <p>Created {{ task.createdAt }}</p>
                  <p>Updated {{ task.updatedAt }}</p>
                </div>
                <Button variant="ghost" size="sm" class="text-xs text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1">
                  <Trash2 :size="13" />
                  Delete
                </Button>
              </div>
            </div>
          </Card>
        </TabsContent>

        <!-- ═══════════════════════════════════════════════════════════
             VARIANT B — "Notion Style" (property bar + content)
             ═══════════════════════════════════════════════════════════ -->
        <TabsContent value="variant-b">
          <Card class="max-w-[540px] py-0 gap-0">
            <!-- Top property bar -->
            <div class="px-5 py-3 border-b bg-muted/30 rounded-t-xl">
              <div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-xs">
                <div class="flex items-center gap-1.5">
                  <span class="text-muted-foreground">Status</span>
                  <StatusBadge :status="task.status" />
                </div>
                <Separator orientation="vertical" class="h-4" />
                <div class="flex items-center gap-1.5">
                  <span class="text-muted-foreground">Priority</span>
                  <PriorityBadge :priority="task.priority" />
                </div>
                <Separator orientation="vertical" class="h-4" />
                <div class="flex items-center gap-1.5">
                  <span class="text-muted-foreground">Project</span>
                  <span class="font-medium text-foreground">{{ task.project }}</span>
                </div>
                <Separator orientation="vertical" class="h-4" />
                <div class="flex items-center gap-1.5">
                  <CalendarDays :size="12" class="text-muted-foreground" />
                  <span class="text-foreground">{{ task.dueDate }}</span>
                </div>
                <Separator orientation="vertical" class="h-4" />
                <a :href="task.adoUrl" class="flex items-center gap-1 text-blue-600 hover:text-blue-700 transition-colors">
                  <Bug :size="12" />
                  <span class="font-medium">{{ task.adoId }}</span>
                  <Badge variant="outline" class="text-[10px] h-4 px-1 ml-0.5">
                    {{ task.adoState }}
                  </Badge>
                </a>
              </div>
            </div>

            <!-- Content area -->
            <div class="px-5 py-5 space-y-5">
              <!-- Title + Tags -->
              <div class="space-y-2">
                <Input
                  :model-value="task.title"
                  class="text-xl font-semibold border-none shadow-none px-0 h-auto focus-visible:ring-0 focus-visible:border-none"
                />
                <div class="flex flex-wrap items-center gap-1.5">
                  <TagChip v-for="tag in task.tags" :key="tag" :tag="tag" removable />
                  <Button variant="ghost" size="sm" class="h-5 text-[10px] px-1.5 text-muted-foreground gap-0.5">
                    <Plus :size="11" />
                    tag
                  </Button>
                </div>
              </div>

              <!-- Subtasks -->
              <div class="space-y-2.5">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-medium">Subtasks</span>
                    <span class="text-xs text-muted-foreground bg-muted rounded px-1.5 py-0.5">
                      {{ subtasksDone }}/{{ subtasksTotal }}
                    </span>
                  </div>
                </div>
                <div class="h-[2px] w-full bg-muted rounded-full overflow-hidden">
                  <div
                    class="h-full bg-blue-500 rounded-full transition-all duration-300"
                    :style="{ width: `${subtaskProgress}%` }"
                  />
                </div>
                <div class="space-y-0.5">
                  <div
                    v-for="sub in subtasks"
                    :key="sub.id"
                    class="flex items-center gap-2.5 py-1.5 px-2 -mx-2 rounded-md hover:bg-muted/50 group transition-colors"
                  >
                    <button
                      :class="cn(
                        'flex-shrink-0 w-[18px] h-[18px] rounded-full border-[1.5px] flex items-center justify-center transition-all',
                        sub.done
                          ? 'bg-blue-500 border-blue-500 text-white'
                          : 'border-border hover:border-blue-400',
                      )"
                    >
                      <Check v-if="sub.done" :size="11" :stroke-width="3" />
                    </button>
                    <span
                      :class="cn(
                        'text-sm',
                        sub.done && 'line-through text-muted-foreground',
                      )"
                    >
                      {{ sub.title }}
                    </span>
                  </div>
                </div>
                <Button variant="ghost" size="sm" class="h-7 text-xs text-muted-foreground gap-1 px-1.5">
                  <Plus :size="14" />
                  Add subtask
                </Button>
              </div>

              <!-- Description -->
              <div class="space-y-2">
                <span class="text-sm font-medium">Description</span>
                <Textarea
                  :model-value="task.description"
                  class="min-h-[80px] text-sm resize-none border-none bg-muted/30 focus-visible:ring-1"
                />
              </div>

              <!-- PRs — Collapsible -->
              <div class="space-y-2">
                <button
                  class="flex items-center gap-1.5 text-sm font-medium hover:text-foreground transition-colors w-full"
                  @click="showPrsB = !showPrsB"
                >
                  <component
                    :is="showPrsB ? ChevronDown : ChevronRight"
                    :size="14"
                    class="text-muted-foreground"
                  />
                  <GitPullRequest :size="14" class="text-muted-foreground" />
                  Pull Requests
                  <Badge variant="secondary" class="ml-auto text-[10px] h-4 px-1.5">
                    {{ prs.length }}
                  </Badge>
                </button>
                <div v-if="showPrsB" class="space-y-1.5 pl-1">
                  <div
                    v-for="pr in prs"
                    :key="pr.id"
                    class="flex items-start gap-2.5 p-2.5 rounded-lg border border-border/50 hover:border-border transition-colors"
                  >
                    <GitPullRequest :size="14" :class="cn('mt-0.5 flex-shrink-0', pr.status === 'active' ? 'text-emerald-500' : 'text-muted-foreground')" />
                    <div class="flex-1 min-w-0 space-y-1">
                      <div class="flex items-center gap-1.5">
                        <span class="text-sm font-medium truncate">{{ pr.title }}</span>
                        <Badge variant="outline" class="text-[10px] h-4 px-1">
                          {{ prStatusLabel(pr.status) }}
                        </Badge>
                      </div>
                      <div class="text-[11px] text-muted-foreground font-mono">
                        {{ pr.repo }} · {{ pr.branch }}
                      </div>
                      <div v-if="pr.reviewers.length" class="flex items-center gap-1 mt-0.5">
                        <Badge
                          v-for="rev in pr.reviewers"
                          :key="rev.name"
                          variant="outline"
                          :class="cn('text-[10px] h-4 px-1.5', reviewerVoteColor(rev.vote))"
                        >
                          {{ rev.name }}
                          <Check v-if="rev.vote === 'approved'" :size="10" class="ml-0.5" />
                        </Badge>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Comments -->
              <div class="space-y-3">
                <span class="text-sm font-medium">Comments</span>
                <Tabs default-value="personal">
                  <TabsList class="h-8">
                    <TabsTrigger value="personal" class="text-xs h-6 px-2.5">
                      Personal
                      <Badge variant="secondary" class="ml-1 text-[10px] h-4 px-1">
                        {{ personalComments.length }}
                      </Badge>
                    </TabsTrigger>
                    <TabsTrigger value="ado" class="text-xs h-6 px-2.5">
                      ADO
                      <Badge variant="secondary" class="ml-1 text-[10px] h-4 px-1">
                        {{ adoComments.length }}
                      </Badge>
                    </TabsTrigger>
                  </TabsList>
                  <TabsContent value="personal" class="mt-3 space-y-3">
                    <div
                      v-for="c in personalComments"
                      :key="c.id"
                      class="flex gap-2.5"
                    >
                      <div
                        :class="cn(
                          'flex-shrink-0 w-7 h-7 rounded-full flex items-center justify-center text-[10px] font-semibold',
                          commentAvatarBg(c.type),
                        )"
                      >
                        {{ commentAvatar(c.type).charAt(0) }}
                      </div>
                      <div class="flex-1 min-w-0 space-y-0.5">
                        <p class="text-sm leading-relaxed">{{ c.text }}</p>
                        <span class="text-[11px] text-muted-foreground">{{ c.time }}</span>
                      </div>
                    </div>
                  </TabsContent>
                  <TabsContent value="ado" class="mt-3 space-y-3">
                    <div
                      v-for="c in adoComments"
                      :key="c.id"
                      class="flex gap-2.5"
                    >
                      <div
                        :class="cn(
                          'flex-shrink-0 w-7 h-7 rounded-full flex items-center justify-center text-[10px] font-semibold',
                          commentAvatarBg(c.type),
                        )"
                      >
                        {{ commentAvatar(c.type).charAt(0) }}
                      </div>
                      <div class="flex-1 min-w-0 space-y-0.5">
                        <p class="text-sm leading-relaxed">{{ c.text }}</p>
                        <span class="text-[11px] text-muted-foreground">{{ c.time }}</span>
                      </div>
                    </div>
                  </TabsContent>
                </Tabs>
                <div class="flex gap-2">
                  <Input
                    v-model="newComment"
                    placeholder="Write a comment..."
                    class="text-sm h-8"
                  />
                  <Button variant="ghost" size="icon-sm" class="flex-shrink-0">
                    <Send :size="14" />
                  </Button>
                </div>
              </div>

              <Separator />

              <!-- Timeline -->
              <div class="space-y-2.5">
                <span class="text-sm font-medium">Activity</span>
                <div class="relative space-y-0">
                  <div
                    v-for="(event, i) in timeline"
                    :key="i"
                    class="flex items-start gap-3 relative"
                  >
                    <div class="flex flex-col items-center">
                      <div
                        :class="cn(
                          'flex-shrink-0 w-5 h-5 rounded-full flex items-center justify-center',
                          timelineIconColor(event.type),
                        )"
                      >
                        <component :is="timelineIcon(event.type)" :size="11" />
                      </div>
                      <div
                        v-if="i < timeline.length - 1"
                        class="w-px h-3 bg-border"
                      />
                    </div>
                    <div class="flex items-center gap-2 pt-0.5 pb-0.5">
                      <span class="text-xs text-muted-foreground">{{ event.text }}</span>
                      <span class="text-[10px] text-muted-foreground/60">{{ event.time }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </Card>
        </TabsContent>

        <!-- ═══════════════════════════════════════════════════════════
             VARIANT C — "GitHub Issue Style" (main + sidebar)
             ═══════════════════════════════════════════════════════════ -->
        <TabsContent value="variant-c">
          <Card class="max-w-[760px] py-0 gap-0">
            <div class="flex">
              <!-- Main content (left ~65%) -->
              <div class="flex-1 min-w-0 px-5 py-5 space-y-5 border-r">
                <!-- Title -->
                <Input
                  :model-value="task.title"
                  class="text-xl font-semibold border-none shadow-none px-0 h-auto focus-visible:ring-0 focus-visible:border-none"
                />

                <Separator />

                <!-- Description -->
                <div class="space-y-2">
                  <div class="flex items-center gap-1.5">
                    <FileEdit :size="14" class="text-muted-foreground" />
                    <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                      Description
                    </span>
                  </div>
                  <Textarea
                    :model-value="task.description"
                    class="min-h-[80px] text-sm resize-none"
                  />
                </div>

                <!-- Subtasks -->
                <div class="space-y-2.5">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-1.5">
                      <Check :size="14" class="text-muted-foreground" />
                      <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                        Subtasks
                      </span>
                    </div>
                    <span class="text-xs text-muted-foreground">
                      {{ subtasksDone }}/{{ subtasksTotal }} complete
                    </span>
                  </div>
                  <div class="h-[2px] w-full bg-muted rounded-full overflow-hidden">
                    <div
                      class="h-full bg-blue-500 rounded-full transition-all duration-300"
                      :style="{ width: `${subtaskProgress}%` }"
                    />
                  </div>
                  <div class="space-y-0.5">
                    <div
                      v-for="sub in subtasks"
                      :key="sub.id"
                      class="flex items-center gap-2.5 py-1.5 px-2 -mx-2 rounded-md hover:bg-muted/50 group transition-colors"
                    >
                      <button
                        :class="cn(
                          'flex-shrink-0 w-[18px] h-[18px] rounded-full border-[1.5px] flex items-center justify-center transition-all',
                          sub.done
                            ? 'bg-blue-500 border-blue-500 text-white'
                            : 'border-border hover:border-blue-400',
                        )"
                      >
                        <Check v-if="sub.done" :size="11" :stroke-width="3" />
                      </button>
                      <span
                        :class="cn(
                          'text-sm',
                          sub.done && 'line-through text-muted-foreground',
                        )"
                      >
                        {{ sub.title }}
                      </span>
                    </div>
                  </div>
                  <Button variant="ghost" size="sm" class="h-7 text-xs text-muted-foreground gap-1 px-1.5">
                    <Plus :size="14" />
                    Add subtask
                  </Button>
                </div>

                <Separator />

                <!-- Comments (combined with type badges) -->
                <div class="space-y-3">
                  <div class="flex items-center gap-1.5">
                    <MessageSquare :size="14" class="text-muted-foreground" />
                    <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                      Comments
                    </span>
                  </div>
                  <div class="space-y-3">
                    <div
                      v-for="c in comments"
                      :key="c.id"
                      class="flex gap-2.5"
                    >
                      <div
                        :class="cn(
                          'flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-[11px] font-semibold',
                          commentAvatarBg(c.type),
                        )"
                      >
                        {{ commentAvatar(c.type).charAt(0) }}
                      </div>
                      <div class="flex-1 min-w-0 p-3 rounded-lg border border-border/60 bg-muted/20">
                        <div class="flex items-center gap-2 mb-1">
                          <span class="text-xs font-medium">
                            {{ c.type === 'personal' ? 'You' : 'Azure DevOps' }}
                          </span>
                          <Badge
                            variant="outline"
                            :class="cn(
                              'text-[10px] h-4 px-1',
                              c.type === 'personal'
                                ? 'bg-blue-500/10 text-blue-600 border-blue-500/20'
                                : 'bg-orange-500/10 text-orange-600 border-orange-500/20',
                            )"
                          >
                            {{ c.type === 'personal' ? 'Personal' : 'ADO' }}
                          </Badge>
                          <span class="text-[11px] text-muted-foreground ml-auto">{{ c.time }}</span>
                        </div>
                        <p class="text-sm leading-relaxed">{{ c.text }}</p>
                      </div>
                    </div>
                  </div>
                  <div class="flex gap-2.5">
                    <div class="flex-shrink-0 w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                      <User :size="14" class="text-muted-foreground" />
                    </div>
                    <div class="flex-1 flex gap-2">
                      <Input
                        v-model="newComment"
                        placeholder="Leave a comment..."
                        class="text-sm h-9"
                      />
                      <Button size="sm" class="gap-1">
                        <Send :size="13" />
                        Comment
                      </Button>
                    </div>
                  </div>
                </div>

                <Separator />

                <!-- Timeline -->
                <div class="space-y-2.5">
                  <div class="flex items-center gap-1.5">
                    <Activity :size="14" class="text-muted-foreground" />
                    <span class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                      Timeline
                    </span>
                  </div>
                  <div class="relative pl-3">
                    <div class="absolute left-[11px] top-1 bottom-1 w-px bg-border" />
                    <div class="space-y-2">
                      <div
                        v-for="(event, i) in timeline"
                        :key="i"
                        class="flex items-center gap-3 relative"
                      >
                        <div
                          :class="cn(
                            'flex-shrink-0 w-5 h-5 rounded-full flex items-center justify-center bg-background border border-border z-10',
                            timelineIconColor(event.type),
                          )"
                        >
                          <component :is="timelineIcon(event.type)" :size="10" />
                        </div>
                        <span class="text-xs">{{ event.text }}</span>
                        <span class="text-[10px] text-muted-foreground ml-auto">{{ event.time }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Sidebar (right ~35%) -->
              <div class="w-[260px] flex-shrink-0 px-4 py-5 space-y-4">
                <!-- Status -->
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Status
                  </label>
                  <Select :default-value="task.status">
                    <SelectTrigger class="h-8 text-xs">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="todo">To Do</SelectItem>
                      <SelectItem value="in_progress">In Progress</SelectItem>
                      <SelectItem value="in_review">In Review</SelectItem>
                      <SelectItem value="blocked">Blocked</SelectItem>
                      <SelectItem value="done">Done</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <!-- Priority -->
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Priority
                  </label>
                  <Select :default-value="task.priority">
                    <SelectTrigger class="h-8 text-xs">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="P0">P0 — Critical</SelectItem>
                      <SelectItem value="P1">P1 — High</SelectItem>
                      <SelectItem value="P2">P2 — Medium</SelectItem>
                      <SelectItem value="P3">P3 — Low</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <!-- Tags -->
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Tags
                  </label>
                  <div class="flex flex-wrap gap-1">
                    <TagChip v-for="tag in task.tags" :key="tag" :tag="tag" removable />
                    <Button variant="ghost" size="sm" class="h-5 text-[10px] px-1 text-muted-foreground">
                      <Plus :size="11" />
                    </Button>
                  </div>
                </div>

                <!-- Project -->
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Project
                  </label>
                  <div class="flex items-center gap-1.5 text-sm">
                    <Folder :size="13" class="text-muted-foreground" />
                    {{ task.project }}
                  </div>
                </div>

                <Separator />

                <!-- ADO Link -->
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Azure DevOps
                  </label>
                  <div class="space-y-2 p-2.5 rounded-lg border border-border/60 bg-muted/20">
                    <div class="flex items-center gap-1.5">
                      <Bug :size="13" class="text-red-500" />
                      <Badge variant="outline" class="text-[10px] h-4 px-1.5 bg-blue-500/10 text-blue-600 border-blue-500/20">
                        {{ task.adoId }}
                      </Badge>
                    </div>
                    <p class="text-xs truncate">{{ task.adoTitle }}</p>
                    <div class="flex items-center justify-between">
                      <Badge variant="outline" class="text-[10px] h-4 px-1.5">
                        {{ task.adoState }}
                      </Badge>
                      <Button variant="ghost" size="sm" class="h-5 text-[10px] px-1.5 gap-1 text-blue-600" as="a" :href="task.adoUrl">
                        <ExternalLink :size="10" />
                        Open in ADO
                      </Button>
                    </div>
                  </div>
                </div>

                <Separator />

                <!-- PRs -->
                <div class="space-y-2">
                  <button
                    class="flex items-center gap-1.5 text-[11px] font-medium text-muted-foreground uppercase tracking-wide hover:text-foreground transition-colors w-full"
                    @click="showPrsC = !showPrsC"
                  >
                    <component
                      :is="showPrsC ? ChevronDown : ChevronRight"
                      :size="12"
                    />
                    Pull Requests
                    <Badge variant="secondary" class="ml-auto text-[10px] h-4 px-1">
                      {{ prs.length }}
                    </Badge>
                  </button>
                  <div v-if="showPrsC" class="space-y-1.5">
                    <div
                      v-for="pr in prs"
                      :key="pr.id"
                      class="p-2 rounded-md border border-border/50 hover:border-border transition-colors space-y-1"
                    >
                      <div class="flex items-center gap-1.5">
                        <div
                          :class="cn('w-1.5 h-1.5 rounded-full flex-shrink-0', prStatusDot(pr.status))"
                        />
                        <span class="text-xs font-medium truncate">{{ pr.title }}</span>
                      </div>
                      <div class="text-[10px] text-muted-foreground font-mono pl-3">
                        #{{ pr.id }} · {{ pr.repo }}
                      </div>
                      <div v-if="pr.reviewers.length" class="flex items-center gap-1 pl-3">
                        <Badge
                          v-for="rev in pr.reviewers"
                          :key="rev.name"
                          variant="outline"
                          :class="cn('text-[9px] h-3.5 px-1', reviewerVoteColor(rev.vote))"
                        >
                          {{ rev.name }}
                          <Check v-if="rev.vote === 'approved'" :size="8" class="ml-0.5" />
                        </Badge>
                      </div>
                    </div>
                  </div>
                </div>

                <Separator />

                <!-- Due Date -->
                <div class="space-y-1.5">
                  <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wide">
                    Due Date
                  </label>
                  <div class="flex items-center gap-1.5 text-sm">
                    <CalendarDays :size="13" class="text-muted-foreground" />
                    {{ task.dueDate }}
                  </div>
                </div>

                <!-- Timestamps -->
                <div class="space-y-1">
                  <div class="text-[11px] text-muted-foreground flex items-center gap-1.5">
                    <Clock :size="11" />
                    Created {{ task.createdAt }}
                  </div>
                  <div class="text-[11px] text-muted-foreground flex items-center gap-1.5">
                    <Clock :size="11" />
                    Updated {{ task.updatedAt }}
                  </div>
                </div>

                <Separator />

                <!-- Delete -->
                <Button variant="ghost" size="sm" class="w-full text-xs text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1.5 justify-start">
                  <Trash2 :size="13" />
                  Delete task
                </Button>
              </div>
            </div>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  </ScrollArea>
</template>
