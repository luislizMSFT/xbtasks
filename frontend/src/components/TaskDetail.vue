<script setup lang="ts">
import { ref, watch, computed, reactive } from 'vue'
import { useTaskStore, type Task } from '@/stores/tasks'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import AdoBadge from '@/components/ui/AdoBadge.vue'
import TagChip from '@/components/ui/TagChip.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  X,
  Trash2,
  Plus,
  ChevronDown,
  ChevronRight,
  Check,
  XIcon,
  ExternalLink,
  GitPullRequest,
  MessageSquare,
  RefreshCw,
  CircleDot,
  Link,
  Bug,
  CheckSquare,
  BookOpen,
  Square,
  SquareCheck,
} from 'lucide-vue-next'
import { cn } from '@/lib/utils'

// ── Store & emits ──
const taskStore = useTaskStore()
const emit = defineEmits<{ close: [] }>()

// ── Editable fields ──
const editTitle = ref('')
const editDescription = ref('')
const editStatus = ref('')
const editPriority = ref('')
const editTags = ref<string[]>([])
const editDueDate = ref('')
const editProject = ref('')
const newTag = ref('')

const statuses = ['todo', 'in_progress', 'in_review', 'done', 'blocked', 'cancelled']
const priorities = ['P0', 'P1', 'P2', 'P3']
const mockProjects = ['xb-tasks', 'xb-services', 'xb-infra', 'xb-docs']

const task = computed(() => taskStore.selectedTask)

watch(task, (t) => {
  if (t) {
    editTitle.value = t.title
    editDescription.value = t.description
    editStatus.value = t.status
    editPriority.value = t.priority
    editDueDate.value = t.dueDate
    editProject.value = t.projectId ? String(t.projectId) : ''
    editTags.value = t.tags ? t.tags.split(',').map(s => s.trim()).filter(Boolean) : []
  }
}, { immediate: true })

// ── Save / delete ──
async function save() {
  if (!task.value) return
  const updated: Task = {
    ...task.value,
    title: editTitle.value,
    description: editDescription.value,
    status: editStatus.value,
    priority: editPriority.value,
    dueDate: editDueDate.value,
    tags: editTags.value.join(','),
  }
  await taskStore.updateTask(updated)
}

async function onStatusChange(status: string) {
  editStatus.value = status
  await save()
}

async function onPriorityChange(priority: string) {
  editPriority.value = priority
  await save()
}

async function deleteTask() {
  if (!task.value) return
  await taskStore.deleteTask(task.value.id)
  emit('close')
}

// ── Tags ──
function addTag() {
  const tag = newTag.value.trim()
  if (tag && !editTags.value.includes(tag)) {
    editTags.value.push(tag)
    newTag.value = ''
    save()
  }
}

function removeTag(tag: string) {
  editTags.value = editTags.value.filter(t => t !== tag)
  save()
}

// ── Subtasks (mock) ──
interface Subtask { id: number; title: string; done: boolean }
const subtasks = ref<Subtask[]>([
  { id: 1, title: 'Validate token refresh flow', done: true },
  { id: 2, title: 'Add redirect guard', done: false },
  { id: 3, title: 'Write integration test', done: false },
])
const newSubtaskTitle = ref('')

const subtaskProgress = computed(() => {
  const total = subtasks.value.length
  const done = subtasks.value.filter(s => s.done).length
  return { done, total, percent: total ? (done / total) * 100 : 0 }
})

function toggleSubtask(id: number) {
  const st = subtasks.value.find(s => s.id === id)
  if (st) st.done = !st.done
}

function addSubtask() {
  const title = newSubtaskTitle.value.trim()
  if (!title) return
  subtasks.value.push({
    id: Math.max(0, ...subtasks.value.map(s => s.id)) + 1,
    title,
    done: false,
  })
  newSubtaskTitle.value = ''
}

// ── Pull Requests (mock) ──
interface PRReviewer { name: string; vote: string }
interface PullRequest {
  id: number
  title: string
  status: 'active' | 'draft' | 'completed'
  repo: string
  sourceBranch: string
  targetBranch: string
  reviewers: PRReviewer[]
  votes: { up: number; down: number }
  url: string
}
const mockPRs = ref<PullRequest[]>([
  {
    id: 1, title: 'Fix auth redirect loop', status: 'active',
    repo: 'xb-services', sourceBranch: 'fix/auth-redirect', targetBranch: 'main',
    reviewers: [{ name: 'Alex', vote: 'approved' }, { name: 'Sam', vote: 'pending' }],
    votes: { up: 1, down: 0 }, url: '#',
  },
  {
    id: 2, title: 'Add task dependency API', status: 'completed',
    repo: 'xb-services', sourceBranch: 'feat/deps', targetBranch: 'main',
    reviewers: [{ name: 'Jordan', vote: 'approved' }],
    votes: { up: 1, down: 0 }, url: '#',
  },
])
const prSectionOpen = ref(true)
const expandedPRs = reactive<Record<number, boolean>>({})

function togglePR(id: number) {
  expandedPRs[id] = !expandedPRs[id]
}

function prStatusVariant(status: string) {
  switch (status) {
    case 'active': return 'bg-blue-500/15 text-blue-600 border-blue-500/20'
    case 'draft': return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
    case 'completed': return 'bg-emerald-500/15 text-emerald-600 border-emerald-500/20'
    default: return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
  }
}

function voteIcon(vote: string) {
  return vote === 'approved' ? '✓' : vote === 'rejected' ? '✗' : '·'
}

function voteColor(vote: string) {
  return vote === 'approved' ? 'text-emerald-500' : vote === 'rejected' ? 'text-red-500' : 'text-muted-foreground'
}

// ── Comments (mock) ──
interface Comment { id: number; type: 'personal' | 'ado'; text: string; createdAt: string }
const mockComments = ref<Comment[]>([
  { id: 1, type: 'personal', text: 'Need to check the token refresh logic', createdAt: '2026-04-01T10:00:00Z' },
  { id: 2, type: 'ado', text: 'Updated the work item with latest findings', createdAt: '2026-04-02T14:30:00Z' },
])
const newCommentText = ref('')
const commentTab = ref('personal')

const personalComments = computed(() => mockComments.value.filter(c => c.type === 'personal'))
const adoComments = computed(() => mockComments.value.filter(c => c.type === 'ado'))

function addComment() {
  const text = newCommentText.value.trim()
  if (!text) return
  mockComments.value.push({
    id: Math.max(0, ...mockComments.value.map(c => c.id)) + 1,
    type: commentTab.value as 'personal' | 'ado',
    text,
    createdAt: new Date().toISOString(),
  })
  newCommentText.value = ''
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit',
  })
}

// ── Timeline (mock) ──
interface TimelineEvent { type: 'status' | 'pr' | 'comment' | 'ado'; description: string; time: string }
const mockTimeline: TimelineEvent[] = [
  { type: 'status', description: 'Status changed to In Progress', time: '2h ago' },
  { type: 'pr', description: 'PR #234 linked — Fix auth redirect', time: '1h ago' },
  { type: 'comment', description: 'Added personal note', time: '45m ago' },
  { type: 'ado', description: 'Synced with ADO-48291', time: '30m ago' },
]

function timelineIcon(type: string) {
  switch (type) {
    case 'status': return CircleDot
    case 'pr': return GitPullRequest
    case 'comment': return MessageSquare
    case 'ado': return RefreshCw
    default: return CircleDot
  }
}

function timelineColor(type: string) {
  switch (type) {
    case 'status': return 'text-blue-500'
    case 'pr': return 'text-violet-500'
    case 'comment': return 'text-amber-500'
    case 'ado': return 'text-emerald-500'
    default: return 'text-muted-foreground'
  }
}

// ── ADO icon helper ──
function adoTypeIcon(adoId: string) {
  if (adoId.includes('Bug') || adoId.toLowerCase().includes('bug')) return '🐛'
  if (adoId.includes('Test')) return '✅'
  return '📖'
}

// ── Timestamps ──
const timeCreated = computed(() => {
  if (!task.value) return ''
  return new Date(task.value.createdAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric', hour: 'numeric', minute: '2-digit',
  })
})

const timeUpdated = computed(() => {
  if (!task.value) return ''
  return new Date(task.value.updatedAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric', hour: 'numeric', minute: '2-digit',
  })
})
</script>

<template>
  <Transition
    enter-active-class="transition-transform duration-200 ease-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-150 ease-in"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <aside
      v-if="task"
      class="w-[400px] border-l border-border bg-background flex flex-col"
    >
      <!-- Close button -->
      <div class="flex items-center justify-end px-4 py-2 border-b border-border">
        <Button variant="ghost" size="icon" class="h-7 w-7" @click="emit('close')">
          <X :size="16" />
        </Button>
      </div>

      <ScrollArea class="flex-1">
        <div class="px-5 py-4 space-y-5">

          <!-- ═══════ TOP: Title ═══════ -->
          <Input
            v-model="editTitle"
            @blur="save"
            @keydown.enter="($event.target as HTMLInputElement)?.blur()"
            class="text-lg font-semibold border-none shadow-none focus-visible:ring-0 bg-transparent px-0"
            placeholder="Task title…"
          />

          <!-- ═══════ TOP: Subtasks ═══════ -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">
                Subtasks
              </label>
              <span class="text-[11px] text-muted-foreground">
                {{ subtaskProgress.done }}/{{ subtaskProgress.total }}
              </span>
            </div>

            <!-- Progress bar -->
            <div class="h-[2px] w-full bg-muted rounded-full overflow-hidden">
              <div
                class="h-full bg-emerald-500 transition-all duration-300"
                :style="{ width: `${subtaskProgress.percent}%` }"
              />
            </div>

            <!-- Subtask list -->
            <div class="space-y-1">
              <button
                v-for="st in subtasks"
                :key="st.id"
                class="flex items-center gap-2 w-full text-left px-1 py-1 rounded hover:bg-muted/50 transition-colors group"
                @click="toggleSubtask(st.id)"
              >
                <component
                  :is="st.done ? SquareCheck : Square"
                  :size="15"
                  :class="st.done ? 'text-emerald-500' : 'text-muted-foreground'"
                />
                <span :class="cn('text-sm', st.done && 'line-through text-muted-foreground')">
                  {{ st.title }}
                </span>
              </button>
            </div>

            <!-- Add subtask -->
            <div class="flex items-center gap-1.5">
              <Input
                v-model="newSubtaskTitle"
                @keydown.enter.prevent="addSubtask"
                class="h-7 text-xs flex-1"
                placeholder="Add subtask…"
              />
              <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="addSubtask">
                <Plus :size="14" />
              </Button>
            </div>
          </div>

          <!-- ═══════ PR SECTION (collapsible) ═══════ -->
          <div class="space-y-2">
            <button
              class="flex items-center gap-2 w-full text-left"
              @click="prSectionOpen = !prSectionOpen"
            >
              <component :is="prSectionOpen ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
              <span class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">
                Pull Requests
              </span>
              <Badge variant="secondary" class="text-[10px] h-4 px-1.5">
                {{ mockPRs.length }}
              </Badge>
            </button>

            <div v-if="prSectionOpen" class="space-y-1.5">
              <div
                v-for="pr in mockPRs"
                :key="pr.id"
                class="rounded-md border border-border overflow-hidden"
              >
                <!-- Compact PR row -->
                <button
                  class="flex items-center gap-2 w-full text-left px-3 py-2 hover:bg-muted/50 transition-colors"
                  @click="togglePR(pr.id)"
                >
                  <Badge
                    variant="outline"
                    :class="cn('text-[10px] capitalize shrink-0', prStatusVariant(pr.status))"
                  >
                    {{ pr.status }}
                  </Badge>
                  <span class="text-sm truncate flex-1">{{ pr.title }}</span>
                  <span class="flex items-center gap-1 text-[11px] text-muted-foreground shrink-0">
                    <span class="text-emerald-500">✓{{ pr.votes.up }}</span>
                    <span v-if="pr.votes.down" class="text-red-500">✗{{ pr.votes.down }}</span>
                    <span class="ml-1 opacity-60">{{ pr.reviewers.length }}</span>
                  </span>
                </button>

                <!-- Expanded details -->
                <div v-if="expandedPRs[pr.id]" class="px-3 pb-3 pt-1 space-y-2 border-t border-border bg-muted/30">
                  <div class="text-[11px] text-muted-foreground">
                    <span class="font-medium">{{ pr.repo }}</span>
                    <span class="mx-1">·</span>
                    <span>{{ pr.sourceBranch }}</span>
                    <span class="mx-1">→</span>
                    <span>{{ pr.targetBranch }}</span>
                  </div>
                  <div class="space-y-1">
                    <div
                      v-for="rev in pr.reviewers"
                      :key="rev.name"
                      class="flex items-center gap-2 text-xs"
                    >
                      <span :class="voteColor(rev.vote)" class="font-mono text-sm">
                        {{ voteIcon(rev.vote) }}
                      </span>
                      <span>{{ rev.name }}</span>
                      <span class="text-muted-foreground capitalize text-[10px]">({{ rev.vote }})</span>
                    </div>
                  </div>
                  <Button variant="outline" size="sm" class="h-7 text-xs gap-1" as="a" :href="pr.url" target="_blank">
                    <ExternalLink :size="12" />
                    Open in ADO
                  </Button>
                </div>
              </div>

              <!-- Link PR button -->
              <Button variant="ghost" size="sm" class="gap-1.5 text-xs text-muted-foreground w-full justify-start">
                <Link :size="14" />
                Link PR
              </Button>
            </div>
          </div>

          <!-- ═══════ MIDDLE: Description ═══════ -->
          <Separator />
          <Textarea
            v-model="editDescription"
            @blur="save"
            :rows="4"
            class="resize-none"
            placeholder="Add a description…"
          />

          <!-- ═══════ MIDDLE: Comments ═══════ -->
          <Separator />
          <div class="space-y-3">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">
              Comments
            </label>
            <Tabs v-model="commentTab" class="w-full">
              <TabsList class="w-full">
                <TabsTrigger value="personal" class="flex-1 text-xs">Personal</TabsTrigger>
                <TabsTrigger value="ado" class="flex-1 text-xs">ADO</TabsTrigger>
              </TabsList>

              <TabsContent value="personal" class="mt-3 space-y-2">
                <div
                  v-for="c in personalComments"
                  :key="c.id"
                  class="rounded-md bg-muted/50 px-3 py-2 text-sm space-y-1"
                >
                  <p>{{ c.text }}</p>
                  <span class="text-[10px] text-muted-foreground">{{ formatDate(c.createdAt) }}</span>
                </div>
                <p v-if="!personalComments.length" class="text-xs text-muted-foreground italic">
                  No personal notes yet.
                </p>
              </TabsContent>

              <TabsContent value="ado" class="mt-3 space-y-2">
                <div
                  v-for="c in adoComments"
                  :key="c.id"
                  class="rounded-md bg-muted/50 px-3 py-2 text-sm space-y-1"
                >
                  <p>{{ c.text }}</p>
                  <span class="text-[10px] text-muted-foreground">{{ formatDate(c.createdAt) }}</span>
                </div>
                <p v-if="!adoComments.length" class="text-xs text-muted-foreground italic">
                  No ADO comments synced.
                </p>
              </TabsContent>
            </Tabs>

            <!-- Comment input -->
            <div class="flex flex-col gap-2">
              <Textarea
                v-model="newCommentText"
                :rows="2"
                class="resize-none text-sm"
                placeholder="Write a comment…"
                @keydown.meta.enter.prevent="addComment"
              />
              <Button size="sm" class="self-end text-xs" @click="addComment">
                Add comment
              </Button>
            </div>
          </div>

          <!-- ═══════ MIDDLE: Activity Timeline ═══════ -->
          <Separator />
          <div class="space-y-3">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">
              Activity
            </label>
            <div class="space-y-2">
              <div
                v-for="(evt, i) in mockTimeline"
                :key="i"
                class="flex items-start gap-2.5"
              >
                <component
                  :is="timelineIcon(evt.type)"
                  :size="14"
                  :class="cn('mt-0.5 shrink-0', timelineColor(evt.type))"
                />
                <div class="flex-1 min-w-0">
                  <p class="text-sm leading-snug">{{ evt.description }}</p>
                  <span class="text-[10px] text-muted-foreground">{{ evt.time }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- ═══════ BOTTOM: Config Section ═══════ -->
          <Separator />

          <!-- Status & Priority grid -->
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1.5">
              <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Status</label>
              <Select :model-value="editStatus" @update:model-value="(v) => onStatusChange(String(v))">
                <SelectTrigger class="w-full">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="s in statuses" :key="s" :value="s">
                    {{ s.replace('_', ' ') }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-1.5">
              <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Priority</label>
              <Select :model-value="editPriority" @update:model-value="(v) => onPriorityChange(String(v))">
                <SelectTrigger class="w-full">
                  <SelectValue placeholder="Priority" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="p in priorities" :key="p" :value="p">{{ p }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <!-- Tags -->
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Tags</label>
            <div class="flex flex-wrap gap-1.5">
              <TagChip
                v-for="tag in editTags"
                :key="tag"
                :tag="tag"
                removable
                @remove="removeTag(tag)"
              />
              <Input
                v-model="newTag"
                @keydown.enter.prevent="addTag"
                class="inline-flex h-6 w-20 px-1.5 text-[12px] border-none shadow-none focus-visible:ring-0 bg-transparent"
                placeholder="+ tag"
              />
            </div>
          </div>

          <!-- Project selector -->
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Project</label>
            <Select v-model="editProject">
              <SelectTrigger class="w-full">
                <SelectValue placeholder="Select project…" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="proj in mockProjects" :key="proj" :value="proj">
                  {{ proj }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- ADO Link -->
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">ADO Link</label>
            <div v-if="task.adoId" class="flex items-center gap-2">
              <span class="text-base">{{ adoTypeIcon(task.adoId) }}</span>
              <AdoBadge :ado-id="task.adoId" />
              <span class="text-xs text-muted-foreground truncate flex-1">{{ task.title }}</span>
              <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1 shrink-0">
                <ExternalLink :size="10" />
                Open in ADO
              </Button>
            </div>
            <div v-else class="flex items-center gap-2">
              <span class="text-base opacity-30">📋</span>
              <Button variant="ghost" size="sm" class="gap-1.5 text-xs text-muted-foreground">
                <Link :size="14" />
                Link to ADO
              </Button>
            </div>
          </div>

          <!-- Due date -->
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Due Date</label>
            <Input
              v-model="editDueDate"
              type="date"
              @change="save"
              class="w-full"
            />
          </div>

          <!-- Footer -->
          <Separator />
          <div class="space-y-2">
            <div class="flex justify-between text-[11px] text-muted-foreground">
              <span>Created {{ timeCreated }}</span>
              <span>Updated {{ timeUpdated }}</span>
            </div>
            <Button
              variant="ghost"
              size="sm"
              class="gap-1.5 text-xs text-destructive hover:text-destructive hover:bg-destructive/10"
              @click="deleteTask"
            >
              <Trash2 :size="13" />
              Delete task
            </Button>
          </div>

        </div>
      </ScrollArea>
    </aside>
  </Transition>
</template>
