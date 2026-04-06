<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import { useTaskStore, type Task } from '@/stores/tasks'
import { useProjectStore } from '@/stores/projects'
import { usePRStore } from '@/stores/prs'
import TagChip from '@/components/ui/TagChip.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ScrollArea } from '@/components/ui/scroll-area'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import ExternalLinks from '@/components/ExternalLinks.vue'
import CommentsSection from '@/components/CommentsSection.vue'
import SyncConfirmDialog from '@/components/SyncConfirmDialog.vue'
import ConflictResolver from '@/components/ConflictResolver.vue'
import { useSyncStore } from '@/stores/sync'
import {
  X,
  Trash2,
  Plus,
  ChevronDown,
  ChevronRight,
  ExternalLink,
  GitPullRequest,
  MessageSquare,
  Activity,
  Send,
  Upload,
  Folder,
  CalendarDays,
} from 'lucide-vue-next'
import { cn } from '@/lib/utils'

// ── Store & emits ──
const taskStore = useTaskStore()
const projectStore = useProjectStore()
const prStore = usePRStore()
const syncStore = useSyncStore()
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

// ── Helpers ──
function openUrl(url: string) { window.open(url, '_blank') }

// ── Discussion panel ──
const showDiscussion = ref(false)

// ── Subtasks (real) ──
const subtasks = ref<Task[]>([])
const newSubtaskTitle = ref('')

async function loadSubtasks() {
  if (!task.value) return
  try {
    const { GetSubtasks } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    subtasks.value = await GetSubtasks(task.value.id) as Task[]
  } catch {
    subtasks.value = []
  }
}

const subtaskProgress = computed(() => {
  const total = subtasks.value.length
  const done = subtasks.value.filter(s => s.status === 'done').length
  return { done, total, percent: total ? (done / total) * 100 : 0 }
})

async function toggleSubtask(id: number) {
  const st = subtasks.value.find(s => s.id === id)
  if (!st) return
  const newStatus = st.status === 'done' ? 'todo' : 'done'
  try {
    const { SetStatus } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    await SetStatus(id, newStatus)
    st.status = newStatus
  } catch {
    st.status = newStatus // optimistic update
  }
}

async function addSubtask() {
  if (!task.value || !newSubtaskTitle.value.trim()) return
  try {
    const { CreateSubtask } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    const sub = await CreateSubtask(task.value.id, newSubtaskTitle.value.trim(), '', 'P2') as Task
    subtasks.value.push(sub)
    newSubtaskTitle.value = ''
  } catch {
    subtasks.value.push({
      id: Date.now(),
      title: newSubtaskTitle.value.trim(),
      status: 'todo',
    } as Task)
    newSubtaskTitle.value = ''
  }
}

watch(task, () => { loadSubtasks() }, { immediate: true })

// ── Pull Requests (from store) ──
const taskPRs = computed(() => {
  if (!task.value) return []
  return [...prStore.myPRs, ...prStore.reviewPRs].filter(pr =>
    pr.taskId === task.value!.id ||
    (task.value!.adoId && pr.adoId === task.value!.adoId)
  )
})
const prsOpen = ref(false)

function prIconColor(status: string) {
  return status === 'completed' ? 'text-violet-500' : status === 'draft' ? 'text-zinc-400' : 'text-emerald-500'
}

function prStatusClasses(status: string) {
  switch (status) {
    case 'active': return 'bg-emerald-500/15 text-emerald-600 border-emerald-500/20'
    case 'completed': return 'bg-violet-500/15 text-violet-600 border-violet-500/20'
    case 'draft': return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
    case 'abandoned': return 'bg-red-500/15 text-red-600 border-red-500/20'
    default: return 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20'
  }
}

function prStatusLabel(status: string) {
  if (status === 'completed') return 'merged'
  return status
}

onMounted(() => {
  projectStore.fetchProjects()
  prStore.fetchAll()
})

// ── Description ──
const descriptionOpen = ref(true)

// ── Comments (mock) ──
// TODO: Wire to real comment backend when available
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
// TODO: Wire to real activity backend when available
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
      class="border-l border-border bg-background flex flex-col h-full min-h-0 relative"
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
          <Select :model-value="editStatus" @update:model-value="(v) => onStatusChange(String(v))">
            <SelectTrigger class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 h-auto border-border shadow-none [&_svg.lucide-chevron-down]:size-3">
              <span :class="cn('size-2 rounded-full', statusColor[editStatus] ?? 'bg-zinc-400')" />
              <span>{{ statusLabel[editStatus] ?? editStatus }}</span>
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="s in statuses" :key="s" :value="s">
                {{ statusLabel[s] ?? s }}
              </SelectItem>
            </SelectContent>
          </Select>
          <Select :model-value="editPriority" @update:model-value="(v) => onPriorityChange(String(v))">
            <SelectTrigger class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 h-auto border-border shadow-none [&_svg.lucide-chevron-down]:size-3">
              <span :class="cn('size-2 rounded-full', priorityDot[editPriority] ?? 'bg-zinc-400')" />
              <span>{{ editPriority }}</span>
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="p in priorities" :key="p" :value="p">{{ p }}</SelectItem>
            </SelectContent>
          </Select>
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
        <div v-if="subtasks.length > 0" class="mt-2.5">
          <div class="h-1 w-full rounded-full bg-muted">
            <div
              class="h-1 rounded-full bg-blue-500 transition-all duration-300"
              :style="{ width: subtaskProgress.percent + '%' }"
            />
          </div>
        </div>
      </div>

      <!-- ─── Main Content ─────────────────────────────── -->
      <ScrollArea class="flex-1 min-h-0">
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
                    v-for="st in subtasks"
                    :key="st.id"
                    class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50 transition-colors"
                    @click="toggleSubtask(st.id)"
                  >
                    <span
                      :class="cn(
                        'size-3.5 rounded-[3px] border-[1.5px] shrink-0 flex items-center justify-center',
                        st.status === 'done'
                          ? 'bg-emerald-500 border-emerald-500'
                          : 'border-muted-foreground/40'
                      )"
                    >
                      <svg v-if="st.status === 'done'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 12 12" class="size-2.5 text-white" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="2 6 5 9 10 3" />
                      </svg>
                    </span>
                    <span :class="cn('text-[13px]', st.status === 'done' && 'line-through text-muted-foreground')">
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
                    {{ taskPRs.length }}
                  </Badge>
                </button>
                <div v-if="prsOpen" class="flex flex-col mt-2 gap-1">
                  <div
                    v-for="pr in taskPRs"
                    :key="pr.id"
                    class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50 transition-colors cursor-pointer"
                    @click="pr.prUrl && openUrl(pr.prUrl)"
                  >
                    <GitPullRequest :size="14" :class="cn('shrink-0', prIconColor(pr.status))" />
                    <span class="text-[13px] truncate flex-1">{{ pr.title }}</span>
                    <span class="text-[11px] text-muted-foreground shrink-0">#{{ pr.prNumber }}</span>
                    <span class="text-[10px] text-muted-foreground shrink-0 truncate max-w-[100px]">{{ pr.repo }}</span>
                    <Badge
                      variant="outline"
                      :class="cn('text-[10px] capitalize shrink-0 h-4 px-1.5', prStatusClasses(pr.status))"
                    >
                      {{ prStatusLabel(pr.status) }}
                    </Badge>
                  </div>
                  <p v-if="taskPRs.length === 0" class="text-xs text-muted-foreground italic py-1 px-1">
                    No linked pull requests
                  </p>
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

              <Separator />

              <!-- External Links -->
              <div class="px-4">
                <ExternalLinks :task-id="task.id" />
              </div>

              <Separator />

              <!-- Comments -->
              <div class="px-4">
                <CommentsSection :task-id="task.id" :is-public-task="taskStore.isPublic(task.id)" />
              </div>

              <!-- Push to ADO (visible only for linked tasks) -->
              <div v-if="taskStore.isPublic(task.id)" class="px-4">
                <Separator class="mb-3" />
                <div class="flex gap-2">
                  <Button variant="outline" size="sm" @click="syncStore.generateOutboundDiff(task.id)">
                    <Upload :size="14" class="mr-1" />
                    Push to ADO
                  </Button>
                </div>
              </div>

              <!-- Discussion trigger -->
              <div class="px-4 pb-3">
                <button
                  class="flex items-center gap-2 w-full text-xs text-muted-foreground hover:text-foreground transition-colors py-2 px-3 rounded-md hover:bg-muted/50 border border-border"
                  @click="showDiscussion = true"
                >
                  <MessageSquare :size="14" />
                  <span class="font-medium">Discussion</span>
                  <Badge variant="secondary" class="h-4 text-[10px] px-1.5 ml-auto">
                    {{ mockNotes.length + mockAdoComments.length }}
                  </Badge>
                </button>
              </div>
        </div>
      </ScrollArea>

      <!-- ─── Activity Bar ─────────────────────────────── -->
      <div v-if="mockActivity.length > 0" class="shrink-0 border-t border-border px-4 py-1.5">
        <div class="flex items-center gap-2 overflow-x-auto scrollbar-none">
          <Activity :size="12" class="text-muted-foreground shrink-0" />
          <div class="flex items-center gap-3 text-[11px] text-muted-foreground whitespace-nowrap">
            <span v-for="(item, i) in mockActivity" :key="'act' + i" class="inline-flex items-center gap-1">
              <span class="size-1 rounded-full bg-muted-foreground/50 shrink-0" />
              {{ item.event }}
              <span class="text-muted-foreground/50">{{ item.time }}</span>
            </span>
          </div>
        </div>
      </div>

      <!-- ─── Details Footer ─────────────────────────────── -->
      <div class="shrink-0 border-t border-border">
        <div class="px-4 py-2">
          <span class="font-semibold text-xs mb-2 block">Details</span>
          <div class="grid grid-cols-2 gap-x-4 gap-y-3">

            <!-- Project -->
            <div class="flex flex-col gap-0.5">
              <span class="text-muted-foreground text-[11px] font-medium">Project</span>
              <Select v-model="editProject">
                <SelectTrigger class="h-7 text-[13px] border-none shadow-none px-1 -ml-1">
                  <span class="flex items-center gap-1.5">
                    <Folder :size="12" class="text-muted-foreground" />
                    <SelectValue placeholder="Select project..." />
                  </span>
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="proj in projectStore.projects" :key="proj.id" :value="String(proj.id)">
                    {{ proj.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- Due Date -->
            <div class="flex flex-col gap-0.5">
              <span class="text-muted-foreground text-[11px] font-medium">Due Date</span>
              <div class="flex items-center gap-1.5">
                <CalendarDays :size="12" class="text-muted-foreground" />
                <Input
                  v-model="editDueDate"
                  type="date"
                  @change="save"
                  class="h-7 text-[13px] border-none shadow-none px-1 -ml-1 flex-1"
                />
              </div>
            </div>

            <!-- ADO Link (full width) -->
            <div class="flex flex-col gap-0.5 col-span-2">
              <span class="text-muted-foreground text-[11px] font-medium">ADO Link</span>
              <span v-if="task.adoId" class="flex items-center gap-1.5 text-blue-500 cursor-pointer hover:underline text-[13px]">
                <AzureDevOpsIcon :size="12" />
                {{ adoNumber(task.adoId) }}
                <ExternalLink :size="10" />
              </span>
              <span v-else class="text-muted-foreground text-[13px]">Not linked</span>
            </div>
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

      <!-- ─── Discussion Slide-up Panel ───────────────────── -->
      <Transition
        enter-active-class="transition-transform duration-300 ease-out"
        enter-from-class="translate-y-full"
        enter-to-class="translate-y-0"
        leave-active-class="transition-transform duration-200 ease-in"
        leave-from-class="translate-y-0"
        leave-to-class="translate-y-full"
      >
        <div v-if="showDiscussion" class="absolute inset-0 bg-background flex flex-col z-50">
          <!-- Panel header -->
          <div class="shrink-0 flex items-center justify-between border-b border-border px-4 py-2.5">
            <div class="flex items-center gap-2">
              <MessageSquare :size="14" />
              <span class="font-semibold text-sm">Discussion</span>
            </div>
            <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="showDiscussion = false">
              <ChevronDown :size="16" />
            </Button>
          </div>

          <!-- Scrollable discussion content -->
          <ScrollArea class="flex-1 min-h-0">
            <!-- Notes section -->
            <div class="px-4 pt-3">
              <div class="flex items-center gap-2 mb-2">
                <span class="font-semibold text-xs">Notes</span>
                <Badge v-if="mockNotes.length > 0" variant="secondary" class="h-4 text-[10px] px-1.5">
                  {{ mockNotes.length }}
                </Badge>
              </div>
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
            </div>

            <Separator class="my-3" />

            <!-- ADO Comments section -->
            <div class="px-4">
              <div class="flex items-center gap-2 mb-2">
                <AzureDevOpsIcon :size="12" />
                <span class="font-semibold text-xs">ADO Comments</span>
                <Badge v-if="mockAdoComments.length > 0" variant="secondary" class="h-4 text-[10px] px-1.5">
                  {{ mockAdoComments.length }}
                </Badge>
              </div>
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
            </div>

            <Separator class="my-3" />

            <!-- Activity timeline -->
            <div class="px-4 pb-4">
              <div class="flex items-center gap-2 mb-2">
                <Activity :size="12" class="text-muted-foreground" />
                <span class="font-semibold text-xs">Activity</span>
              </div>
              <div class="flex flex-col gap-2">
                <div v-for="(item, i) in mockActivity" :key="'disc-act' + i" class="flex items-start gap-2">
                  <span class="size-1.5 rounded-full bg-muted-foreground/50 shrink-0 mt-1.5" />
                  <div class="flex-1 min-w-0">
                    <span class="text-[13px]">{{ item.event }}</span>
                    <span class="text-muted-foreground text-[10px] ml-2">{{ item.time }}</span>
                  </div>
                </div>
              </div>
            </div>
          </ScrollArea>

          <!-- Note input at bottom -->
          <div class="shrink-0 border-t border-border px-4 py-2.5">
            <div class="flex items-center gap-2">
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
          </div>
        </div>
      </Transition>

      <!-- Sync Confirm Dialog (pushed changes preview) -->
      <SyncConfirmDialog />
      <!-- Conflict Resolution Dialog (per-field local vs ADO picker) -->
      <ConflictResolver />
    </aside>
  </Transition>
</template>
