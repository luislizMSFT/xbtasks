<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import { useTaskStore, type Task } from '@/stores/tasks'
import { useNotify } from '@/composables/useNotify'
import { useProjectStore } from '@/stores/projects'
import { usePRStore } from '@/stores/prs'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import ExternalLinks from '@/components/tasks/ExternalLinks.vue'
import CommentsSection from '@/components/tasks/CommentsSection.vue'
import SyncConfirmDialog from '@/components/ado/SyncConfirmDialog.vue'
import ConflictResolver from '@/components/ado/ConflictResolver.vue'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useSyncStore } from '@/stores/sync'
import {
  X,
  Trash2,
  Plus,
  ChevronDown,
  ChevronRight,
  ExternalLink,
  GitPullRequest,
  Upload,
  Folder,
  CalendarDays,
  Link2Off,
  Loader2,
} from 'lucide-vue-next'
import { cn } from '@/lib/utils'
import { statusBgColor, statusIcon, priorityDotBgColor, prStatusClasses } from '@/lib/styles'
import { branchName } from '@/stores/prs'

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
const editProject = ref('none')
const newTag = ref('')

const statuses = ['todo', 'in_progress', 'in_review', 'done', 'blocked', 'cancelled']
const priorities = ['P0', 'P1', 'P2', 'P3']

// Keep a stable task reference that persists during deselect.
// Without this, task becomes null mid-render causing "$setup.task.adoId" crash.
const selectedTask = computed(() => taskStore.selectedTask)
const lastTask = ref<Task | null>(null)
const task = computed(() => selectedTask.value ?? lastTask.value)
const isOpen = computed(() => !!selectedTask.value)
const detailLoading = ref(false)

watch(selectedTask, (t, oldT) => {
  if (t) {
    // Show brief loading flash when switching to a different task
    if (oldT && t.id !== oldT.id) {
      saving = false
      detailLoading.value = true
      requestAnimationFrame(() => { detailLoading.value = false })
    }
    lastTask.value = t
    editTitle.value = t.title
    editDescription.value = t.description.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim()
    editStatus.value = t.status
    editPriority.value = t.priority
    editDueDate.value = t.dueDate
    editProject.value = t.projectId ? String(t.projectId) : 'none'
    editTags.value = t.tags ? t.tags.split(',').map(s => s.trim()).filter(Boolean) : []
  }
}, { immediate: true })

// ── Save / delete ──
let saving = false
let saveTaskId: number | null = null
async function save() {
  if (!task.value || saving) return
  const currentId = task.value.id
  // Skip save if fields haven't actually changed
  if (
    editTitle.value === task.value.title &&
    editDescription.value === task.value.description &&
    editStatus.value === task.value.status &&
    editPriority.value === task.value.priority &&
    editDueDate.value === task.value.dueDate &&
    editTags.value.join(',') === task.value.tags &&
    (editProject.value !== 'none' ? Number(editProject.value) : null) === (task.value.projectId ?? null)
  ) return
  saving = true
  saveTaskId = currentId
  try {
    // Abort if task changed during save (user clicked another task)
    if (task.value?.id !== currentId) return
    const updated: Task = {
      ...task.value,
      title: editTitle.value,
      description: editDescription.value,
      status: editStatus.value,
      priority: editPriority.value,
      dueDate: editDueDate.value,
      tags: editTags.value.join(','),
      projectId: editProject.value !== 'none' ? Number(editProject.value) : null,
    }
    await taskStore.updateTask(updated)
  } catch (e) {
    console.warn('[TaskDetail] save failed:', e)
  } finally {
    saving = false
  }
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
async function openUrl(url: string) {
  try {
    const { openURL } = await import('@/api/browser')
    await openURL(url)
  } catch { window.open(url, '_blank') }
}

// ── Subtasks (real) ──
const subtasks = ref<Task[]>([])
const newSubtaskTitle = ref('')

async function loadSubtasks() {
  if (!task.value) return
  try {
    const { getSubtasks } = await import('@/api/tasks')
    subtasks.value = await getSubtasks(task.value.id) as Task[]
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
    const { setStatus } = await import('@/api/tasks')
    await setStatus(id, newStatus)
    st.status = newStatus
  } catch {
    st.status = newStatus // optimistic update
  }
}

async function addSubtask() {
  if (!task.value || !newSubtaskTitle.value.trim()) return
  try {
    const { createSubtask } = await import('@/api/tasks')
    const sub = await createSubtask(task.value.id, newSubtaskTitle.value.trim(), '', 'P2') as Task
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

watch(selectedTask, () => { loadSubtasks() }, { immediate: true })

// ── Pull Requests (from store) ──
const taskPRs = computed(() => {
  if (!task.value) return []
  const t = task.value
  return [...prStore.myPRs, ...prStore.reviewPRs].filter(pr =>
    pr.taskId === t.id ||
    (t.adoId && pr.adoId === t.adoId)
  )
})
function prIconColor(status: string) {
  return status === 'completed' ? 'text-violet-500' : status === 'draft' ? 'text-zinc-400' : 'text-emerald-500'
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
const descriptionOpen = ref(false)

// ── Style helpers ──

const statusLabelMap: Record<string, string> = {
  in_progress: 'In Progress',
  in_review: 'In Review',
  todo: 'To Do',
  blocked: 'Blocked',
  done: 'Done',
  cancelled: 'Cancelled',
}

// ── ADO helpers ──
function adoNumber(adoId: string): string {
  const match = adoId.match(/\d+/)
  return match ? `#${match[0]}` : adoId
}

// Build ADO URL using cached org/project from the work item
const adoUrl = ref('')
const lastAdoId = ref('')
watch(selectedTask, async (t) => {
  const adoId = t?.adoId || ''
  if (adoId === lastAdoId.value) return
  lastAdoId.value = adoId
  if (!adoId) { adoUrl.value = ''; return }
  try {
    const { getCachedWorkItem } = await import('@/api/workitems')
    const wi = await getCachedWorkItem(adoId) as any
    if (wi?.org && wi?.project) {
      adoUrl.value = `https://${wi.org}.visualstudio.com/${wi.project}/_workitems/edit/${adoId}`
    } else {
      adoUrl.value = ''
    }
  } catch {
    adoUrl.value = ''
  }
}, { immediate: true })

async function openAdoLink() {
  if (!adoUrl.value) return
  openUrl(adoUrl.value)
}

const notify = useNotify()

async function unlinkFromADO() {
  if (!task.value?.adoId) return
  try {
    const { unlinkTask } = await import('@/api/links')
    await unlinkTask(task.value.id, task.value.adoId, false)
    await taskStore.fetchTasks()
    notify.info('Task unlinked from ADO')
  } catch (e) {
    console.warn('[TaskDetail] Failed to unlink:', e)
    notify.error('Failed to unlink')
  }
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

// ── Project name helper ──
const projectName = computed(() => {
  if (!task.value?.projectId) return ''
  const proj = projectStore.projects.find(p => p.id === task.value!.projectId)
  return proj?.name || ''
})
</script>

<template>
  <aside
    v-if="isOpen || lastTask"
    class="flex flex-col h-full min-h-0 bg-background"
  >
    <!-- Loading overlay when switching tasks -->
    <Transition
      enter-active-class="transition-opacity duration-100"
      enter-from-class="opacity-100"
      enter-to-class="opacity-0"
      leave-active-class="transition-opacity duration-75"
      leave-from-class="opacity-0"
      leave-to-class="opacity-100"
    >
      <div v-if="detailLoading" class="absolute inset-0 z-10 bg-background/80 flex items-center justify-center">
        <Loader2 :size="20" class="animate-spin text-muted-foreground" />
      </div>
    </Transition>

    <!-- ─── Header (shrink-0, border-b) ─────────────────── -->
    <div class="shrink-0 border-b border-border px-4 pt-3 pb-3">
      <!-- Row 1: Title + close button -->
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

      <!-- Row 2: Status + Priority + Due Date + Delete -->
      <div class="flex items-center gap-x-2 gap-y-1.5 mt-2 flex-wrap">
        <Select :model-value="editStatus" @update:model-value="(v) => onStatusChange(String(v))">
          <SelectTrigger class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 h-auto border-border shadow-none [&_svg.lucide-chevron-down]:size-3">
            <span :class="cn('size-2 rounded-full', statusBgColor(editStatus))" />
            <span>{{ statusLabelMap[editStatus] ?? editStatus }}</span>
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in statuses" :key="s" :value="s">
              {{ statusLabelMap[s] ?? s }}
            </SelectItem>
          </SelectContent>
        </Select>
        <Select :model-value="editPriority" @update:model-value="(v) => onPriorityChange(String(v))">
          <SelectTrigger class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 h-auto border-border shadow-none [&_svg.lucide-chevron-down]:size-3">
            <span :class="cn('size-2 rounded-full', priorityDotBgColor(editPriority))" />
            <span>{{ editPriority }}</span>
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="p in priorities" :key="p" :value="p">{{ p }}</SelectItem>
          </SelectContent>
        </Select>

        <div class="inline-flex items-center gap-1.5 text-xs font-medium rounded-full pl-2 pr-1.5 py-0.5 border border-border">
          <CalendarDays :size="11" class="text-muted-foreground" />
          <input
            v-model="editDueDate"
            type="date"
            @change="save"
            class="bg-transparent text-xs border-none outline-none w-[105px] py-0"
            :class="!editDueDate && 'text-muted-foreground'"
          />
        </div>

        <Button
          variant="ghost"
          size="sm"
          class="h-6 text-[11px] text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-1 ml-auto"
          @click="onDeleteTask"
        >
          <Trash2 :size="12" />
          Delete
        </Button>
      </div>

      <!-- Row 3: ADO badge + actions (if linked) -->
      <div v-if="task?.adoId" class="flex items-center gap-2 mt-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <span
              class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-500 border border-blue-500/30 rounded-full px-2 py-0.5 cursor-pointer hover:bg-blue-500/10 transition-colors"
            >
              <AzureDevOpsIcon :size="12" />
              {{ adoNumber(task.adoId) }}
              <ChevronDown :size="10" />
            </span>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" class="w-40">
            <DropdownMenuItem @click="openAdoLink" class="text-xs gap-2">
              <ExternalLink :size="12" />
              Open in ADO
            </DropdownMenuItem>
            <DropdownMenuItem
              @click="task && syncStore.generateOutboundDiff(task.id)"
              class="text-xs gap-2"
            >
              <Upload :size="12" />
              Push to ADO
            </DropdownMenuItem>
            <DropdownMenuItem @click="unlinkFromADO" class="text-xs gap-2 text-red-500">
              <Link2Off :size="12" />
              Unlink
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <span class="text-[9px] text-muted-foreground/60">local copy</span>
      </div>

      <!-- Row 4: Project context bar (if project set) -->
      <div v-if="projectName" class="flex items-center gap-1.5 mt-2 text-xs text-muted-foreground">
        <Folder :size="11" />
        <span>{{ projectName }}</span>
      </div>
    </div>

    <!-- ─── Description (collapsible) ─────────────────── -->
    <div class="shrink-0 px-4 py-2 border-b border-border">
      <button class="flex items-center gap-1.5 w-full" @click="descriptionOpen = !descriptionOpen">
        <ChevronRight :size="14" class="text-muted-foreground transition-transform" :class="descriptionOpen && 'rotate-90'" />
        <h3 class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
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

    <!-- ─── External Links ──────────────────────────── -->
    <div v-if="task" class="shrink-0 px-4 py-2 border-b border-border">
      <ExternalLinks :task-id="task.id" />
    </div>

    <!-- ─── 3-Tab Section ───────────────────────────── -->
    <Tabs default-value="subtasks" class="flex-1 flex flex-col min-h-0">
      <TabsList class="shrink-0 w-full justify-start border-b bg-transparent px-4 rounded-none h-9">
        <TabsTrigger value="subtasks" class="text-xs">Subtasks</TabsTrigger>
        <TabsTrigger value="prs" class="text-xs">PRs</TabsTrigger>
        <TabsTrigger value="notes" class="text-xs">Notes</TabsTrigger>
      </TabsList>

      <!-- Subtasks tab -->
      <TabsContent value="subtasks" class="flex-1 min-h-0 mt-0">
        <ScrollArea class="h-full">
          <div class="px-4 py-3">
            <!-- Subtask progress -->
            <div v-if="subtasks.length > 0" class="mb-3">
              <div class="flex items-center gap-2 mb-1.5">
                <Badge variant="secondary" class="h-4 text-[10px] px-1.5">
                  {{ subtaskProgress.done }}/{{ subtaskProgress.total }}
                </Badge>
              </div>
              <div class="h-1 w-full rounded-full bg-muted">
                <div
                  class="h-1 rounded-full bg-blue-500 transition-all duration-300"
                  :style="{ width: subtaskProgress.percent + '%' }"
                />
              </div>
            </div>

            <!-- Subtask list -->
            <div class="flex flex-col gap-0.5">
              <button
                v-for="st in subtasks"
                :key="st.id"
                class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50 transition-colors"
                @click="toggleSubtask(st.id)"
              >
                <component
                  :is="statusIcon(st.status)"
                  :size="14"
                  :class="st.status === 'done' ? 'text-emerald-500' : 'text-muted-foreground'"
                  class="shrink-0"
                />
                <span :class="cn('text-[13px] flex-1', st.status === 'done' && 'line-through text-muted-foreground')">
                  {{ st.title }}
                </span>
                <span :class="cn('size-1.5 rounded-full shrink-0', priorityDotBgColor(st.priority))" />
              </button>
            </div>

            <!-- Add subtask input -->
            <button
              v-if="!newSubtaskTitle"
              class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground mt-2 transition-colors"
              @click="newSubtaskTitle = ' '"
            >
              <Plus :size="12" />
              Add subtask
            </button>
            <div v-else class="flex items-center gap-1.5 mt-2">
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
        </ScrollArea>
      </TabsContent>

      <!-- PRs tab -->
      <TabsContent value="prs" class="flex-1 min-h-0 mt-0">
        <ScrollArea class="h-full">
          <div class="px-4 py-3">
            <div class="flex flex-col gap-1">
              <div
                v-for="pr in taskPRs"
                :key="pr.id"
                class="flex items-center gap-2 py-1.5 px-1 rounded hover:bg-muted/50 transition-colors cursor-pointer"
                @click="pr.prUrl && openUrl(pr.prUrl)"
              >
                <GitPullRequest :size="14" :class="cn('shrink-0', prIconColor(pr.status))" />
                <span class="text-[13px] truncate flex-1">{{ pr.title }}</span>
                <Badge
                  variant="outline"
                  :class="cn('text-[10px] capitalize shrink-0 h-4 px-1.5', prStatusClasses(pr.status))"
                >
                  {{ prStatusLabel(pr.status) }}
                </Badge>
              </div>
              <div
                v-for="pr in taskPRs"
                :key="`meta-${pr.id}`"
                class="flex items-center gap-2 pl-7 -mt-1 mb-1 text-[11px] text-muted-foreground"
              >
                <span class="tabular-nums">#{{ pr.prNumber }}</span>
                <span class="truncate max-w-[120px]">{{ pr.repo }}</span>
                <span v-if="pr.sourceBranch" class="truncate max-w-[120px]">{{ branchName(pr.sourceBranch) }}</span>
              </div>
            </div>
            <p v-if="taskPRs.length === 0" class="text-xs text-muted-foreground italic py-4 text-center">
              No pull requests linked to this task
            </p>
          </div>
        </ScrollArea>
      </TabsContent>

      <!-- Notes tab -->
      <TabsContent value="notes" class="flex-1 min-h-0 mt-0">
        <ScrollArea class="h-full">
          <div class="px-4 py-3">
            <CommentsSection v-if="task" :task-id="task.id" :is-public-task="taskStore.isPublic(task.id)" />
          </div>
        </ScrollArea>
      </TabsContent>
    </Tabs>

    <!-- ─── Footer: timestamps ──────────────────────── -->
    <div class="shrink-0 border-t border-border flex items-center px-4 py-1.5">
      <span class="text-[11px] text-muted-foreground">
        Created {{ timeCreated }} · Updated {{ timeUpdated }}
      </span>
    </div>

    <!-- Sync Confirm Dialog (pushed changes preview) -->
    <SyncConfirmDialog />
    <!-- Conflict Resolution Dialog (per-field local vs ADO picker) -->
    <ConflictResolver />
  </aside>
</template>
