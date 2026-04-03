<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useTaskStore, type Task } from '@/stores/tasks'
import TagChip from '@/components/ui/TagChip.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
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
  ExternalLink,
  GitPullRequest,
  MessageSquare,
  ListTodo,
  Activity,
  Send,
  Tag,
  Folder,
  CalendarDays,
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

async function onDeleteTask() {
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

// ── Tabs ──
const detailTab = ref('work')
const commentTab = ref('notes')

// ── Subtasks (mock) ──
interface Subtask { id: number; title: string; done: boolean }
const mockSubtasks = ref<Subtask[]>([
  { id: 1, title: 'Identify refresh token edge case', done: true },
  { id: 2, title: 'Add token expiry middleware', done: true },
  { id: 3, title: 'Update error handling in auth callback', done: false },
  { id: 4, title: 'Add E2E test for token refresh', done: false },
])
const newSubtaskTitle = ref('')

const subtaskProgress = computed(() => {
  const total = mockSubtasks.value.length
  const done = mockSubtasks.value.filter(s => s.done).length
  return { done, total, percent: total ? (done / total) * 100 : 0 }
})

function toggleSubtask(id: number) {
  const st = mockSubtasks.value.find(s => s.id === id)
  if (st) st.done = !st.done
}

function addSubtask() {
  const title = newSubtaskTitle.value.trim()
  if (!title) return
  mockSubtasks.value.push({
    id: Math.max(0, ...mockSubtasks.value.map(s => s.id)) + 1,
    title,
    done: false,
  })
  newSubtaskTitle.value = ''
}

// ── Pull Requests (mock) ──
interface MockPR {
  id: number
  title: string
  status: 'active' | 'merged' | 'draft'
  repo: string
  additions: number
  deletions: number
}
const mockPRs = ref<MockPR[]>([
  { id: 234, title: 'Fix auth redirect loop', status: 'active', repo: 'xb-services', additions: 400, deletions: 120 },
  { id: 240, title: 'Add token expiry middleware', status: 'merged', repo: 'xb-services', additions: 80, deletions: 10 },
])
const prsOpen = ref(false)

function prIconColor(status: string) {
  return status === 'merged' ? 'text-violet-500' : 'text-emerald-500'
}

function prStatusClasses(status: string) {
  switch (status) {
    case 'active': return 'bg-emerald-500/15 text-emerald-600 border-emerald-500/20'
    case 'merged': return 'bg-violet-500/15 text-violet-600 border-violet-500/20'
    case 'draft': return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
    default: return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
  }
}

// ── Description ──
const descriptionOpen = ref(true)

// ── Comments (mock) ──
interface MockComment { id: number; author: string; initials: string; text: string; time: string }
const mockNotes = ref<MockComment[]>([
  { id: 1, author: 'You', initials: 'LL', text: 'Root cause is in the token refresh — concurrent requests cause a race condition.', time: '25m ago' },
  { id: 2, author: 'You', initials: 'LL', text: 'Middleware approach looks cleaner than patching the callback directly.', time: '1h ago' },
])
const mockAdoComments = ref<MockComment[]>([
  { id: 1, author: 'Alex K.', initials: 'AK', text: 'Confirmed the loop on staging. Happens when refresh token is within 30s of expiry.', time: '2h ago' },
  { id: 2, author: 'Jordan M.', initials: 'JM', text: 'Priority bump — this is blocking the demo prep.', time: '4h ago' },
])
const noteText = ref('')

function addNote() {
  const text = noteText.value.trim()
  if (!text) return
  mockNotes.value.unshift({
    id: Math.max(0, ...mockNotes.value.map(c => c.id)) + 1,
    author: 'You',
    initials: 'LL',
    text,
    time: 'just now',
  })
  noteText.value = ''
}

// ── Activity (mock) ──
const mockActivity = [
  { event: 'Status changed to In Progress', time: '2d ago' },
  { event: 'Linked to ADO Bug #48291', time: '2d ago' },
  { event: 'Subtask completed: Identify refresh token edge case', time: '1d ago' },
  { event: 'PR #234 opened', time: '6h ago' },
]

// ── Style helpers ──
const statusColor: Record<string, string> = {
  in_progress: 'bg-blue-500',
  in_review: 'bg-violet-500',
  todo: 'bg-zinc-400',
  blocked: 'bg-red-500',
  done: 'bg-emerald-500',
  cancelled: 'bg-zinc-400',
}

const statusLabel: Record<string, string> = {
  in_progress: 'In Progress',
  in_review: 'In Review',
  todo: 'To Do',
  blocked: 'Blocked',
  done: 'Done',
  cancelled: 'Cancelled',
}

const priorityDot: Record<string, string> = {
  P0: 'bg-red-500',
  P1: 'bg-orange-500',
  P2: 'bg-amber-500',
  P3: 'bg-zinc-400',
}

// ── ADO helpers ──
function adoNumber(adoId: string): string {
  const match = adoId.match(/\d+/)
  return match ? `#${match[0]}` : adoId
}

// ── Timestamps ──
const timeCreated = computed(() => {
  if (!task.value) return ''
  return new Date(task.value.createdAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
})

const timeUpdated = computed(() => {
  if (!task.value) return ''
  return new Date(task.value.updatedAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
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
      class="border-l border-border bg-background flex flex-col h-full min-h-0"
    >
      <!-- ─── Header (shrink-0, border-b) ─────────────────── -->
      <div class="shrink-0 border-b border-border px-4 pt-3 pb-3">
        <div class="flex items-start gap-2">
          <Input
            v-model="editTitle"
            @blur="save"
            @keydown.enter="($event.target as HTMLInputElement)?.blur()"
            :class="cn(
              'text-base font-semibold border-none shadow-none focus-visible:ring-0 bg-transparent px-0 h-auto flex-1',
              editStatus === 'done' && 'line-through text-muted-foreground'
            )"
            placeholder="Task title..."
          />
          <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0 -mt-0.5" @click="emit('close')">
            <X :size="16" />
          </Button>
        </div>

        <!-- Badges row -->
        <div class="flex items-center gap-2 mt-2 flex-wrap">
          <span class="inline-flex items-center gap-1.5 text-xs font-medium border border-border rounded-full px-2 py-0.5">
            <span :class="cn('size-2 rounded-full', statusColor[editStatus] ?? 'bg-zinc-400')" />
            {{ statusLabel[editStatus] ?? editStatus }}
          </span>
          <span class="inline-flex items-center gap-1.5 text-xs font-medium border border-border rounded-full px-2 py-0.5">
            <span :class="cn('size-2 rounded-full', priorityDot[editPriority] ?? 'bg-zinc-400')" />
            {{ editPriority }}
          </span>
          <span
            v-if="task.adoId"
            class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-500 border border-blue-500/30 rounded-full px-2 py-0.5 cursor-pointer hover:bg-blue-500/10 transition-colors"
          >
            <AzureDevOpsIcon :size="12" />
            {{ adoNumber(task.adoId) }}
            <ExternalLink :size="10" />
          </span>
        </div>

        <!-- Subtask progress bar -->
        <div v-if="mockSubtasks.length > 0" class="mt-2.5">
          <div class="h-1 w-full rounded-full bg-muted">
            <div
              class="h-1 rounded-full bg-blue-500 transition-all duration-300"
              :style="{ width: subtaskProgress.percent + '%' }"
            />
          </div>
        </div>
      </div>

      <!-- ─── Tabs: Work / Discussion ────────────────────── -->
      <Tabs v-model="detailTab" class="flex flex-col flex-1 min-h-0">
        <div class="shrink-0 px-4 border-b border-border">
          <TabsList class="h-8 w-full">
            <TabsTrigger value="work" class="text-xs h-6 flex-1 gap-1">
              <ListTodo :size="13" />
              Work
            </TabsTrigger>
            <TabsTrigger value="discussion" class="text-xs h-6 flex-1 gap-1">
              <MessageSquare :size="13" />
              Discussion
            </TabsTrigger>
          </TabsList>
        </div>

        <!-- Scrollable tab content -->
        <ScrollArea class="flex-1 min-h-0">

          <!-- ── Work tab ──────────────────────────────── -->
          <TabsContent value="work" class="mt-0 p-0">
            <div class="flex flex-col space-y-2">

              <!-- Subtasks -->
              <div class="px-4 pt-3">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <span class="font-semibold text-xs">Subtasks</span>
                    <Badge variant="secondary" class="h-4 text-[10px] px-1.5">
                      {{ subtaskProgress.done }}/{{ subtaskProgress.total }}
                    </Badge>
                  </div>
                </div>
                <div class="flex flex-col gap-0.5">
                  <button
                    v-for="st in mockSubtasks"
                    :key="st.id"
                    class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50 transition-colors"
                    @click="toggleSubtask(st.id)"
                  >
                    <span
                      :class="cn(
                        'size-3.5 rounded-[3px] border-[1.5px] shrink-0 flex items-center justify-center',
                        st.done
                          ? 'bg-emerald-500 border-emerald-500'
                          : 'border-muted-foreground/40'
                      )"
                    >
                      <svg v-if="st.done" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 12 12" class="size-2.5 text-white" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="2 6 5 9 10 3" />
                      </svg>
                    </span>
                    <span :class="cn('text-[13px]', st.done && 'line-through text-muted-foreground')">
                      {{ st.title }}
                    </span>
                  </button>
                </div>
                <button
                  v-if="!newSubtaskTitle"
                  class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground mt-1.5 transition-colors"
                  @click="newSubtaskTitle = ' '"
                >
                  <Plus :size="12" />
                  Add subtask
                </button>
                <div v-else class="flex items-center gap-1.5 mt-1.5">
                  <Input
                    v-model.trim="newSubtaskTitle"
                    @keydown.enter.prevent="addSubtask"
                    @blur="newSubtaskTitle = ''"
                    class="h-7 text-xs flex-1"
                    placeholder="Subtask title..."
                    autofocus
                  />
                  <Button variant="ghost" size="icon" class="size-7 shrink-0" @click="addSubtask">
                    <Plus :size="14" />
                  </Button>
                </div>
              </div>

              <Separator />

              <!-- Pull Requests (collapsible, slim) -->
              <div class="px-4">
                <button class="flex items-center gap-1.5 w-full" @click="prsOpen = !prsOpen">
                  <ChevronDown v-if="prsOpen" :size="14" class="text-muted-foreground" />
                  <ChevronRight v-else :size="14" class="text-muted-foreground" />
                  <span class="font-semibold text-xs">Pull Requests</span>
                  <Badge variant="secondary" class="h-4 text-[10px] px-1.5 ml-1">
                    {{ mockPRs.length }}
                  </Badge>
                </button>
                <div v-if="prsOpen" class="flex flex-col mt-2 gap-1">
                  <div
                    v-for="pr in mockPRs"
                    :key="pr.id"
                    class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50 transition-colors"
                  >
                    <GitPullRequest :size="14" :class="cn('shrink-0', prIconColor(pr.status))" />
                    <span class="text-[13px] truncate flex-1">{{ pr.title }}</span>
                    <span class="text-[11px] text-muted-foreground shrink-0">#{{ pr.id }}</span>
                    <span class="flex items-center gap-1 text-[10px] shrink-0">
                      <span class="text-emerald-500">+{{ pr.additions }}</span>
                      <span class="text-red-500">-{{ pr.deletions }}</span>
                    </span>
                    <Badge
                      variant="outline"
                      :class="cn('text-[10px] capitalize shrink-0 h-4 px-1.5', prStatusClasses(pr.status))"
                    >
                      {{ pr.status }}
                    </Badge>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Description (collapsible) -->
              <div class="px-4 pb-4">
                <button class="flex items-center gap-1.5 w-full" @click="descriptionOpen = !descriptionOpen">
                  <ChevronDown v-if="descriptionOpen" :size="14" class="text-muted-foreground" />
                  <ChevronRight v-else :size="14" class="text-muted-foreground" />
                  <span class="font-semibold text-xs">Description</span>
                </button>
                <Textarea
                  v-if="descriptionOpen"
                  v-model="editDescription"
                  @blur="save"
                  :rows="4"
                  class="resize-none mt-2 text-[13px]"
                  placeholder="Add a description..."
                />
              </div>

              <div class="h-4" />
            </div>
          </TabsContent>

          <!-- ── Discussion tab ────────────────────────── -->
          <TabsContent value="discussion" class="mt-0 p-0">
            <div class="flex flex-col space-y-2">

              <!-- Comments (Notes / ADO sub-tabs) -->
              <div class="px-4 pt-3">
                <Tabs v-model="commentTab" class="w-full">
                  <TabsList class="h-7 w-full">
                    <TabsTrigger value="notes" class="text-[11px] h-5 flex-1 gap-1">
                      <MessageSquare :size="11" />
                      Notes
                      <Badge v-if="mockNotes.length > 0" variant="secondary" class="h-3.5 text-[9px] px-1">
                        {{ mockNotes.length }}
                      </Badge>
                    </TabsTrigger>
                    <TabsTrigger value="ado" class="text-[11px] h-5 flex-1 gap-1">
                      <AzureDevOpsIcon :size="11" />
                      ADO Comments
                      <Badge v-if="mockAdoComments.length > 0" variant="secondary" class="h-3.5 text-[9px] px-1">
                        {{ mockAdoComments.length }}
                      </Badge>
                    </TabsTrigger>
                  </TabsList>

                  <TabsContent value="notes" class="mt-2">
                    <div class="flex flex-col gap-2.5">
                      <div
                        v-for="c in mockNotes"
                        :key="'n' + c.id"
                        class="flex gap-2"
                      >
                        <div class="size-6 rounded-full bg-muted flex items-center justify-center shrink-0">
                          <span class="text-[9px] font-semibold text-muted-foreground">{{ c.initials }}</span>
                        </div>
                        <div class="flex-1 min-w-0">
                          <div class="flex items-baseline gap-2">
                            <span class="font-medium text-[13px]">{{ c.author }}</span>
                            <span class="text-muted-foreground text-[10px]">{{ c.time }}</span>
                          </div>
                          <p class="text-muted-foreground mt-0.5 text-[13px]">{{ c.text }}</p>
                        </div>
                      </div>
                      <p v-if="mockNotes.length === 0" class="text-xs text-muted-foreground italic">
                        No notes yet.
                      </p>
                    </div>
                    <div class="flex items-center gap-2 mt-3">
                      <Input
                        v-model="noteText"
                        placeholder="Add a note..."
                        class="h-7 text-xs flex-1"
                        @keydown.enter.prevent="addNote"
                      />
                      <Button variant="ghost" size="icon" class="size-7 shrink-0" @click="addNote">
                        <Send :size="13" />
                      </Button>
                    </div>
                  </TabsContent>

                  <TabsContent value="ado" class="mt-2">
                    <div class="flex flex-col gap-2.5">
                      <div
                        v-for="c in mockAdoComments"
                        :key="'a' + c.id"
                        class="flex gap-2"
                      >
                        <div class="size-6 rounded-full bg-muted flex items-center justify-center shrink-0">
                          <span class="text-[9px] font-semibold text-muted-foreground">{{ c.initials }}</span>
                        </div>
                        <div class="flex-1 min-w-0">
                          <div class="flex items-baseline gap-2">
                            <span class="font-medium text-[13px]">{{ c.author }}</span>
                            <span class="text-muted-foreground text-[10px]">{{ c.time }}</span>
                          </div>
                          <p class="text-muted-foreground mt-0.5 text-[13px]">{{ c.text }}</p>
                        </div>
                      </div>
                      <p v-if="mockAdoComments.length === 0" class="text-xs text-muted-foreground italic">
                        No ADO comments synced.
                      </p>
                    </div>
                  </TabsContent>
                </Tabs>
              </div>

              <Separator />

              <!-- Activity Timeline -->
              <div class="px-4 pb-4">
                <div class="flex items-center gap-1.5 mb-3">
                  <Activity :size="13" class="text-muted-foreground" />
                  <span class="font-semibold text-xs">Activity</span>
                </div>
                <div class="flex flex-col">
                  <div
                    v-for="(item, i) in mockActivity"
                    :key="'act' + i"
                    class="flex gap-2.5 relative"
                  >
                    <div class="flex flex-col items-center w-3 shrink-0">
                      <div class="size-1.5 rounded-full bg-muted-foreground/50 mt-[7px] shrink-0 z-10" />
                      <div
                        v-if="i < mockActivity.length - 1"
                        class="w-px flex-1 bg-border"
                      />
                    </div>
                    <div class="flex-1 min-w-0 pb-3">
                      <span class="text-[13px] text-muted-foreground">{{ item.event }}</span>
                      <span class="text-[10px] text-muted-foreground/60 ml-1.5">&mdash; {{ item.time }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="h-4" />
            </div>
          </TabsContent>
        </ScrollArea>
      </Tabs>

      <!-- ─── Sticky Config Footer ───────────────────────── -->
      <div class="shrink-0 border-t border-border">
        <div class="px-4 py-2">
          <span class="font-semibold text-xs mb-1.5 block">Details</span>
          <div class="grid grid-cols-[auto_1fr] items-center gap-x-4 gap-y-1.5">

            <!-- Status -->
            <span class="text-muted-foreground text-[13px]">Status</span>
            <Select :model-value="editStatus" @update:model-value="(v) => onStatusChange(String(v))">
              <SelectTrigger class="h-7 text-[13px] border-none shadow-none px-1 -ml-1">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="s in statuses" :key="s" :value="s">
                  {{ statusLabel[s] ?? s }}
                </SelectItem>
              </SelectContent>
            </Select>

            <!-- Priority -->
            <span class="text-muted-foreground text-[13px]">Priority</span>
            <Select :model-value="editPriority" @update:model-value="(v) => onPriorityChange(String(v))">
              <SelectTrigger class="h-7 text-[13px] border-none shadow-none px-1 -ml-1">
                <SelectValue placeholder="Priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="p in priorities" :key="p" :value="p">{{ p }}</SelectItem>
              </SelectContent>
            </Select>

            <!-- Project -->
            <span class="text-muted-foreground text-[13px]">Project</span>
            <Select v-model="editProject">
              <SelectTrigger class="h-7 text-[13px] border-none shadow-none px-1 -ml-1">
                <span class="flex items-center gap-1.5">
                  <Folder :size="12" class="text-muted-foreground" />
                  <SelectValue placeholder="Select project..." />
                </span>
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="proj in mockProjects" :key="proj" :value="proj">
                  {{ proj }}
                </SelectItem>
              </SelectContent>
            </Select>

            <!-- Due Date -->
            <span class="text-muted-foreground text-[13px]">Due Date</span>
            <div class="flex items-center gap-1.5">
              <CalendarDays :size="12" class="text-muted-foreground" />
              <Input
                v-model="editDueDate"
                type="date"
                @change="save"
                class="h-7 text-[13px] border-none shadow-none px-1 -ml-1 flex-1"
              />
            </div>

            <!-- Tags -->
            <span class="text-muted-foreground text-[13px]">Tags</span>
            <div class="flex items-center gap-1 flex-wrap">
              <span
                v-for="tag in editTags"
                :key="tag"
                class="inline-flex items-center gap-1 rounded-full border border-border px-2 py-0.5 text-[11px] group"
              >
                <Tag :size="10" class="text-muted-foreground" />
                {{ tag }}
                <button
                  class="opacity-0 group-hover:opacity-100 transition-opacity hover:text-red-500 ml-0.5"
                  @click="removeTag(tag)"
                >
                  <X :size="8" />
                </button>
              </span>
              <Input
                v-model="newTag"
                @keydown.enter.prevent="addTag"
                class="inline-flex h-5 w-16 px-1 text-[11px] border-none shadow-none focus-visible:ring-0 bg-transparent"
                placeholder="+ tag"
              />
            </div>

            <!-- ADO Link -->
            <span class="text-muted-foreground text-[13px]">ADO Link</span>
            <span v-if="task.adoId" class="flex items-center gap-1.5 text-blue-500 cursor-pointer hover:underline text-[13px]">
              <AzureDevOpsIcon :size="12" />
              {{ adoNumber(task.adoId) }}
              <ExternalLink :size="10" />
            </span>
            <span v-else class="text-muted-foreground text-[13px]">Not linked</span>
          </div>
        </div>

        <!-- Footer bar -->
        <div class="flex items-center justify-between border-t border-border px-4 py-2">
          <span class="text-[11px] text-muted-foreground">
            Created {{ timeCreated }} · Updated {{ timeUpdated }}
          </span>
          <Button
            variant="ghost"
            size="sm"
            class="h-6 text-[11px] text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1"
            @click="onDeleteTask"
          >
            <Trash2 :size="12" />
            Delete
          </Button>
        </div>
      </div>
    </aside>
  </Transition>
</template>
