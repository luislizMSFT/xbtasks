<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import DOMPurify from 'dompurify'
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
import ADODiscussion from '@/components/ado/ADODiscussion.vue'
import AdoItemPicker from '@/components/ado/AdoItemPicker.vue'
import ConflictResolver from '@/components/ado/ConflictResolver.vue'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import {
  X,
  Trash2,
  Plus,
  ChevronRight,
  ExternalLink,
  GitPullRequest,
  Folder,
  CalendarDays,
  Link2Off,
  Link2,
  Loader2,
  User,
  Pencil,
  Save,
  ListChecks,
  Lock,
  CheckCircle2,
  CirclePlay,
  CircleCheck,
  CircleX,
  Clock,
} from 'lucide-vue-next'
import { cn } from '@/lib/utils'
import { statusBgColor, statusIcon, priorityDotBgColor, prStatusClasses, adoTypeIcon, adoTypeColor, adoStateClasses } from '@/lib/styles'
import { branchName } from '@/stores/prs'
import { useAdoMeta } from '@/composables/useAdoMeta'
import { useADOStore } from '@/stores/ado'

// ── Store & emits ──
const taskStore = useTaskStore()
const projectStore = useProjectStore()
const prStore = usePRStore()
const adoMeta = useAdoMeta()
const adoStore = useADOStore()
const emit = defineEmits<{ close: [] }>()

// ── Editable fields ──
const editTitle = ref('')
const editDescription = ref('')
const editStatus = ref('todo')
const editPriority = ref('P2')
const editTags = ref<string[]>([])
const editDueDate = ref('')
const editProject = ref('none')
const newTag = ref('')
const notesTab = ref<'notes' | 'ado'>('notes')
const noteInput = ref('')
const commentsRef = ref<InstanceType<typeof CommentsSection> | null>(null)
const adoDiscussionRef = ref<InstanceType<typeof ADODiscussion> | null>(null)

async function submitNote() {
  const content = noteInput.value.trim()
  if (!content) return
  if (notesTab.value === 'ado' && adoDiscussionRef.value) {
    await adoDiscussionRef.value.addReply(content)
  } else if (commentsRef.value) {
    await commentsRef.value.addComment(content)
  }
  noteInput.value = ''
}

const statuses = ['todo', 'in_progress', 'in_review', 'done', 'blocked', 'cancelled']
const priorities = ['P0', 'P1', 'P2', 'P3']

// Keep a stable task reference that persists during deselect.
// Without this, task becomes null mid-render causing "$setup.task.adoId" crash.
const selectedTask = computed(() => taskStore.selectedTask)
const lastTask = ref<Task | null>(null)
const task = computed(() => selectedTask.value ?? lastTask.value)
const isOpen = computed(() => !!selectedTask.value)
const detailLoading = ref(false)
const isPersonal = computed(() => task.value ? !taskStore.isPublic(task.value.id) : true)
const taskAdoMeta = computed(() => task.value ? adoMeta.getAdoMeta(task.value.id) : null)

watch(selectedTask, (t, oldT) => {
  if (t) {
    // Skip re-syncing edit fields when the same task ref changed due to our own save
    // (splice in updateTask replaces the object, triggering this watcher)
    if (saving && saveTaskId === t.id && oldT && t.id === oldT.id) {
      lastTask.value = t
      return
    }
    // Reset saving flag when switching to a different task
    if (oldT && t.id !== oldT.id) {
      saving = false
      saveTaskId = null
    }
    lastTask.value = t
    editTitle.value = t.title
    editDescription.value = t.description
    editStatus.value = t.status || 'todo'
    editPriority.value = t.priority || 'P2'
    editDueDate.value = t.dueDate
    editProject.value = t.projectId ? String(t.projectId) : 'none'
    editTags.value = t.tags ? t.tags.split(',').map(s => s.trim()).filter(Boolean) : []
  }
}, { immediate: true, flush: 'sync' })

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

const confirmDeleteOpen = ref(false)
const confirmUnlinkOpen = ref(false)
const linkPickerOpen = ref(false)

async function onDeleteTask() {
  if (!task.value) return
  try {
    await taskStore.deleteTask(task.value.id)
    confirmDeleteOpen.value = false
    emit('close')
  } catch {
    notify.error('Failed to delete task')
  }
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

function sanitizeHtml(html: string): string {
  if (!html) return ''
  return DOMPurify.sanitize(html, {
    FORBID_TAGS: ['style', 'script', 'iframe', 'object', 'embed', 'link'],
    FORBID_ATTR: ['style'],
  })
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

// ── Subtask filter (playground parity) ──
const subtaskFilter = ref<'all' | 'ado' | 'personal' | 'mine'>('all')

function isPersonalSubtask(st: Task): boolean {
  return !st.adoId
}

function isOtherPerson(_st: Task): boolean {
  // assignedTo not yet on Task type — future-proof
  return false
}

const filteredSubtasks = computed(() => {
  const subs = subtasks.value
  switch (subtaskFilter.value) {
    case 'ado': return subs.filter(st => !!st.adoId)
    case 'personal': return subs.filter(st => !st.adoId)
    case 'mine': return subs  // all are mine for now
    default: return subs
  }
})

function cycleSubtaskFilter() {
  const order: Array<'all' | 'mine' | 'ado' | 'personal'> = ['all', 'mine', 'ado', 'personal']
  const idx = order.indexOf(subtaskFilter.value)
  subtaskFilter.value = order[(idx + 1) % order.length]
}

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

// ── Pipeline status helpers (playground parity) ──
function pipelineIcon(p: { result: string; status: string }) {
  if (p.result === 'succeeded') return CircleCheck
  if (p.result === 'failed') return CircleX
  if (p.status === 'inprogress') return CirclePlay
  return Clock
}

function pipelineColor(p: { result: string; status: string }) {
  if (p.result === 'succeeded') return 'text-emerald-500'
  if (p.result === 'failed') return 'text-red-500'
  if (p.status === 'inprogress') return 'text-blue-500'
  return 'text-amber-500'
}

function pipelinesForPr(sourceBranch: string) {
  const branch = sourceBranch.replace('refs/heads/', '')
  return adoStore.pipelines.filter(p => p.sourceBranch === branch || p.sourceBranch.replace('refs/heads/', '') === branch)
}

onMounted(() => {
  projectStore.fetchProjects()
  prStore.fetchAll()
  if (!adoStore.pipelines.length) adoStore.fetchPipelines()
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
    const { refreshADOMeta } = await import('@/api/adometa')
    await unlinkTask(task.value.id, task.value.adoId, false)
    await Promise.all([
      taskStore.fetchTasks(),
      taskStore.fetchPublicTaskIds(),
      refreshADOMeta(),
    ])
    adoMeta.refresh()
    adoUrl.value = ''
    lastAdoId.value = ''
    confirmUnlinkOpen.value = false
    notify.info('Task unlinked from ADO')
  } catch (e) {
    console.warn('[TaskDetail] Failed to unlink:', e)
    notify.error('Failed to unlink')
  }
}

async function onLinkSelected(adoId: string) {
  if (!task.value) return
  linkPickerOpen.value = false
  try {
    const { linkTask } = await import('@/api/links')
    const { refreshADOMeta } = await import('@/api/adometa')
    await linkTask(task.value.id, adoId)
    await Promise.all([
      taskStore.fetchTasks(),
      taskStore.fetchPublicTaskIds(),
      refreshADOMeta(),
    ])
    adoMeta.refresh()
    lastAdoId.value = '' // force ADO URL refresh
    notify.success('Linked to ADO')
  } catch (e: any) {
    notify.error(e?.message || 'Failed to link')
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

    <!-- ─── Header (shrink-0, border-b) — playground-style ──── -->
    <div class="shrink-0 border-b border-border px-4 py-3 space-y-2">
      <!-- Row 1: ADO type icon + Title + close -->
      <div class="flex items-center gap-2">
        <component
          v-if="taskAdoMeta?.type"
          :is="adoTypeIcon(taskAdoMeta.type)"
          :size="16"
          :class="adoTypeColor(taskAdoMeta.type)"
          class="shrink-0"
        />
        <Input
          v-model="editTitle"
          @blur="save"
          @keydown.enter="($event.target as HTMLInputElement)?.blur()"
          :class="cn(
            'text-sm font-semibold border-none shadow-none focus-visible:ring-0 bg-transparent px-0 h-auto flex-1',
            editStatus === 'done' && 'line-through text-muted-foreground'
          )"
          placeholder="Task title..."
        />
        <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="emit('close')">
          <X :size="16" />
        </Button>
      </div>

      <!-- Row 2: Metadata badges (status, priority, due date, delete) -->
      <div class="flex items-center gap-2 flex-wrap">
        <Select :model-value="editStatus" @update:model-value="(v) => onStatusChange(String(v))">
          <SelectTrigger class="inline-flex items-center gap-1 text-[10px] font-medium h-5 px-1.5 rounded-full border-border shadow-none [&_svg.lucide-chevron-down]:size-2.5">
            <span :class="cn('size-1.5 rounded-full mr-0.5', statusBgColor(editStatus))" />
            <span>{{ statusLabelMap[editStatus] ?? editStatus }}</span>
            <span v-if="task?.adoId && editStatus !== task.status" class="size-1 rounded-full bg-amber-500 ml-0.5" title="Modified locally" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in statuses" :key="s" :value="s">
              {{ statusLabelMap[s] ?? s }}
            </SelectItem>
          </SelectContent>
        </Select>
        <Select :model-value="editPriority" @update:model-value="(v) => onPriorityChange(String(v))">
          <SelectTrigger class="inline-flex items-center gap-1 text-[10px] font-medium h-5 px-1.5 rounded-full border-border shadow-none [&_svg.lucide-chevron-down]:size-2.5">
            <span :class="cn('size-1.5 rounded-full mr-0.5', priorityDotBgColor(editPriority))" />
            <span>{{ editPriority }}</span>
            <span v-if="task?.adoId && editPriority !== task.priority" class="size-1 rounded-full bg-amber-500 ml-0.5" title="Modified locally" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="p in priorities" :key="p" :value="p">{{ p }}</SelectItem>
          </SelectContent>
        </Select>
        <span class="text-[10px] text-muted-foreground flex items-center gap-1">
          <CalendarDays :size="10" />
          <input
            v-model="editDueDate"
            type="date"
            @change="save"
            class="bg-transparent text-[10px] border-none outline-none w-[95px] py-0"
            :class="!editDueDate && 'text-muted-foreground'"
          />
        </span>
        <Button
          variant="ghost"
          size="sm"
          class="h-5 text-[9px] text-red-500 hover:text-red-600 hover:bg-red-500/10 gap-0.5 ml-auto px-1.5"
          @click="confirmDeleteOpen = true"
        >
          <Trash2 :size="10" /> Delete
        </Button>
      </div>

      <!-- Row 3: Project context bar (only when project assigned) -->
      <div v-if="projectName" class="flex items-center gap-2 px-2 py-1.5 rounded-md bg-muted/40 border border-border/50">
        <component
          v-if="taskAdoMeta?.type"
          :is="adoTypeIcon(taskAdoMeta.type)"
          :size="13"
          :class="adoTypeColor(taskAdoMeta.type)"
          class="shrink-0"
        />
        <Folder v-else :size="13" class="text-muted-foreground shrink-0" />
        <span class="text-xs font-medium text-foreground truncate flex-1">{{ projectName }}</span>
      </div>

      <!-- Row 4: ADO integration bar OR personal task indicator -->
      <div v-if="task?.adoId" class="flex items-center gap-2 px-2 py-1 rounded-md bg-blue-500/5 border border-blue-500/15">
        <AzureDevOpsIcon :size="12" class="text-blue-500 shrink-0" />
        <span class="text-[10px] text-blue-500 tabular-nums font-medium">{{ adoNumber(task.adoId) }}</span>
        <div class="flex-1" />
        <Button variant="ghost" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5 text-blue-500 hover:text-blue-600" @click="openAdoLink">
          <ExternalLink :size="10" /> Open
        </Button>
        <Button variant="ghost" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5 text-red-400 hover:text-red-500 hover:bg-red-500/10" @click="confirmUnlinkOpen = true">
          <Link2Off :size="10" />
        </Button>
      </div>
      <div v-else class="flex items-center gap-2 px-2 py-1 rounded-md border border-dashed border-border/50">
        <User :size="12" class="text-muted-foreground/40 shrink-0" />
        <span class="text-[10px] text-muted-foreground/60 flex-1">Personal task — not linked to ADO</span>
        <Button variant="outline" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5" @click="linkPickerOpen = true">
          <Link2 :size="10" /> Link
        </Button>
      </div>

      <!-- Subtask progress bar (in header) -->
      <div v-if="subtasks.length > 0">
        <div class="h-1 w-full rounded-full bg-muted">
          <div class="h-1 rounded-full bg-blue-500 transition-all" :style="{ width: subtaskProgress.percent + '%' }" />
        </div>
      </div>
    </div>

    <!-- ─── Scrollable content (continuous sections like playground) ─── -->
    <ScrollArea class="flex-1 min-h-0">
      <div class="flex flex-col">

        <!-- Subtasks section (playground-parity) -->
        <div class="border-b border-border px-4 py-2">
          <div class="flex items-center justify-between mb-1.5">
            <div class="flex items-center gap-2">
              <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Subtasks</h3>
              <Badge v-if="subtasks.length > 0" variant="secondary" class="h-4 text-[10px] px-1.5">
                {{ subtaskProgress.done }}/{{ subtaskProgress.total }}
              </Badge>
            </div>
            <div class="flex items-center gap-1.5">
              <button
                v-if="subtasks.length > 0"
                class="text-[10px] flex items-center gap-1 px-1.5 py-0.5 rounded transition-colors"
                :class="subtaskFilter !== 'all' ? 'bg-primary/10 text-primary font-medium' : 'text-muted-foreground hover:text-foreground'"
                @click="cycleSubtaskFilter"
                :title="`Filter: ${subtaskFilter}`"
              >
                <component :is="subtaskFilter === 'personal' ? ListChecks : subtaskFilter === 'ado' ? Lock : User" :size="10" />
                {{ subtaskFilter === 'all' ? 'All' : subtaskFilter === 'mine' ? 'Mine' : subtaskFilter === 'ado' ? 'ADO' : 'Personal' }}
              </button>
              <button class="text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1" @click="newSubtaskTitle = ' '">
                <Plus :size="11" /> Add
              </button>
            </div>
          </div>

          <!-- Subtask list (playground-parity) -->
          <div v-if="filteredSubtasks.length > 0" class="flex flex-col gap-px">
            <button
              v-for="st in filteredSubtasks"
              :key="st.id"
              class="flex items-center gap-2 py-1.5 px-2 rounded hover:bg-muted/50 transition-colors group"
              @click="toggleSubtask(st.id)"
            >
              <!-- Icon: ADO type icon for ADO subtasks, checkbox for personal -->
              <template v-if="st.adoId">
                <component
                  :is="adoTypeIcon((st as any).category || 'Task')"
                  :size="14"
                  :class="adoTypeColor((st as any).category || 'Task')"
                  class="shrink-0"
                />
              </template>
              <template v-else>
                <component
                  :is="statusIcon(st.status)"
                  :size="14"
                  :class="st.status === 'done' ? 'text-emerald-500' : 'text-muted-foreground'"
                  class="shrink-0"
                />
              </template>
              <!-- Title -->
              <span :class="cn('text-xs truncate flex-1', st.status === 'done' && 'line-through text-muted-foreground')">
                {{ st.title }}
              </span>
              <!-- Personal badge -->
              <span v-if="isPersonalSubtask(st)" class="text-[8px] px-1 py-0.5 rounded bg-primary/8 text-primary/70 shrink-0 border border-primary/10">personal</span>
              <!-- ADO state badge (when available) -->
              <Badge v-if="st.adoId && (st as any).adoState" variant="outline" :class="adoStateClasses((st as any).adoState)" class="text-[8px] h-4 px-1">
                {{ (st as any).adoState }}
              </Badge>
              <!-- Sync status indicator (when available) -->
              <span v-if="(st as any).syncStatus === 'pending'" class="size-1.5 rounded-full bg-amber-500 shrink-0" title="Pending sync" />
              <span v-if="(st as any).syncStatus === 'not-pulled'" class="text-[8px] text-amber-600 shrink-0">not-pulled</span>
              <!-- Assigned-to badge (when available and assigned to someone else) -->
              <span v-if="(st as any).assignedTo && isOtherPerson(st)" class="text-[8px] px-1 py-0.5 rounded bg-orange-500/10 text-orange-600 shrink-0 border border-orange-500/15">
                {{ (st as any).assignedTo }}
              </span>
              <!-- Priority dot -->
              <span :class="cn('size-1.5 rounded-full shrink-0', priorityDotBgColor(st.priority))" />
              <!-- ADO ID -->
              <span v-if="st.adoId" class="text-[9px] text-blue-500/60 tabular-nums shrink-0">#{{ st.adoId }}</span>
            </button>
          </div>

          <!-- Add subtask inline -->
          <div v-if="newSubtaskTitle" class="flex items-center gap-2 mt-1.5 px-2">
            <Input
              v-model.trim="newSubtaskTitle"
              @keydown.enter.prevent="addSubtask"
              @keydown.escape="newSubtaskTitle = ''"
              class="h-7 text-xs flex-1"
              placeholder="New subtask..."
              autofocus
            />
            <Button variant="outline" size="sm" class="h-7 px-2 text-[10px] shrink-0" @click="addSubtask" :disabled="!newSubtaskTitle.trim()">Add</Button>
            <Button variant="ghost" size="sm" class="h-7 px-1.5 text-[10px]" @click="newSubtaskTitle = ''">Cancel</Button>
          </div>

          <p v-if="filteredSubtasks.length === 0 && !newSubtaskTitle" class="text-[11px] text-muted-foreground/40 italic py-2 text-center">
            {{ subtaskFilter === 'mine' ? 'No subtasks assigned to you' : subtaskFilter === 'personal' ? 'No personal subtasks' : subtaskFilter === 'ado' ? 'No ADO subtasks' : 'No subtasks — click Add to create one' }}
          </p>
        </div>

        <!-- Description section -->
        <div class="border-b border-border px-4 py-3">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-1.5">
              <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
            </div>
            <button v-if="!descriptionOpen" class="text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1" @click="descriptionOpen = true">
              <Pencil :size="10" /> Edit
            </button>
            <div v-else class="flex items-center gap-1">
              <button class="text-[10px] text-emerald-600 hover:text-emerald-700 flex items-center gap-0.5" @click="descriptionOpen = false; save()">
                <Save :size="10" /> Save
              </button>
              <button class="text-[10px] text-muted-foreground hover:text-foreground" @click="descriptionOpen = false">Cancel</button>
            </div>
          </div>
          <Textarea
            v-if="descriptionOpen"
            v-model="editDescription"
            @blur="save"
            :rows="4"
            class="resize-y text-xs"
            placeholder="Add a description..."
          />
          <div v-else-if="editDescription"
            class="text-xs text-foreground prose prose-sm max-w-none [&_*]:text-xs [&_*]:text-foreground cursor-pointer hover:bg-muted/30 rounded p-1 -m-1 transition-colors"
            v-html="sanitizeHtml(editDescription)"
            @click="descriptionOpen = true"
          />
          <p v-else class="text-[11px] text-muted-foreground/40 italic cursor-pointer hover:text-muted-foreground" @click="descriptionOpen = true">
            Click to add description...
          </p>
        </div>

        <!-- External Links section -->
        <div v-if="task" class="border-b border-border px-4 py-2">
          <ExternalLinks :task-id="task.id" />
        </div>

        <!-- Pull Requests section -->
        <div class="border-b border-border px-4 py-3">
          <div class="flex items-center gap-2 mb-2">
            <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Pull Requests</h3>
            <Badge variant="secondary" class="h-4 text-[10px] px-1.5">{{ taskPRs.length }}</Badge>
          </div>
          <div class="flex flex-col gap-1">
            <div
              v-for="pr in taskPRs"
              :key="pr.id"
              class="rounded hover:bg-muted/50 cursor-pointer"
              @click="pr.prUrl && openUrl(pr.prUrl)"
            >
              <div class="flex items-center gap-2 py-1.5 px-1">
                <GitPullRequest :size="14" :class="cn('shrink-0', prIconColor(pr.status))" />
                <span class="text-[13px] truncate flex-1">{{ pr.title }}</span>
                <span class="text-[11px] text-muted-foreground shrink-0">#{{ pr.prNumber }}</span>
                <Badge variant="outline" :class="cn('text-[10px] capitalize shrink-0 h-4 px-1.5', prStatusClasses(pr.status))">
                  {{ prStatusLabel(pr.status) }}
                </Badge>
              </div>
              <!-- Pipeline status for this PR -->
              <div v-if="pipelinesForPr(pr.sourceBranch).length > 0" class="flex items-center gap-3 pl-6 pb-1.5">
                <div v-for="pipe in pipelinesForPr(pr.sourceBranch)" :key="pipe.id" class="flex items-center gap-1 text-[9px]">
                  <component :is="pipelineIcon(pipe)" :size="10" :class="pipelineColor(pipe)" />
                  <span class="text-muted-foreground">{{ pipe.name }}</span>
                </div>
              </div>
            </div>
            <p v-if="taskPRs.length === 0" class="text-xs text-muted-foreground italic py-1 px-1">No linked PRs</p>
          </div>
        </div>

        <!-- Notes section -->
        <div class="px-4 py-3">
          <div class="flex items-center gap-2 mb-2">
            <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Notes</h3>
            <!-- Chip toggle (only when task is linked to ADO) -->
            <div v-if="task?.adoId" class="flex items-center gap-1 ml-auto">
              <button
                class="text-[9px] font-medium px-2 py-0.5 rounded-full transition-colors"
                :class="notesTab === 'notes'
                  ? 'bg-foreground/10 text-foreground'
                  : 'text-muted-foreground/60 hover:text-muted-foreground'"
                @click="notesTab = 'notes'"
              >
                Notes
              </button>
              <button
                class="text-[9px] font-medium px-2 py-0.5 rounded-full transition-colors"
                :class="notesTab === 'ado'
                  ? 'bg-blue-500/15 text-blue-500'
                  : 'text-muted-foreground/60 hover:text-muted-foreground'"
                @click="notesTab = 'ado'"
              >
                ADO
              </button>
            </div>
          </div>

          <!-- Shared add input -->
          <div class="flex items-start gap-1.5 mb-2">
            <Textarea
              v-model="noteInput"
              :placeholder="notesTab === 'ado' ? 'Reply on ADO...' : 'Add a note...'"
              class="text-xs min-h-[36px] resize-none flex-1"
              :rows="1"
              @keydown.meta.enter="submitNote"
              @keydown.ctrl.enter="submitNote"
            />
            <Button
              size="sm"
              class="h-[36px] text-[10px] shrink-0"
              :class="notesTab === 'ado' && 'bg-blue-600 hover:bg-blue-700 text-white'"
              @click="submitNote"
              :disabled="!noteInput.trim()"
            >
              Add
            </Button>
          </div>

          <!-- Tab content -->
          <ADODiscussion v-if="task?.adoId && notesTab === 'ado'" ref="adoDiscussionRef" :task-id="task.id" :ado-id="task.adoId" />
          <CommentsSection v-else-if="task" ref="commentsRef" :task-id="task.id" />
        </div>

      </div>
    </ScrollArea>

    <!-- ─── Footer: timestamps ──────────────────────── -->
    <div class="shrink-0 border-t border-border flex items-center px-4 py-1.5">
      <span class="text-[11px] text-muted-foreground">
        Created {{ timeCreated }} · Updated {{ timeUpdated }}
      </span>
    </div>

    <!-- Conflict Resolution Dialog (per-field local vs ADO picker) -->
    <ConflictResolver />

    <!-- ADO Link Picker -->
    <AdoItemPicker
      :open="linkPickerOpen"
      title="Link Task to ADO Work Item"
      @update:open="linkPickerOpen = $event"
      @selected="onLinkSelected"
    />

    <!-- Delete confirmation -->
    <Dialog v-model:open="confirmDeleteOpen">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>Delete task?</DialogTitle>
          <DialogDescription>
            "{{ task?.title }}" will be permanently deleted. This cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" size="sm" @click="confirmDeleteOpen = false">Cancel</Button>
          <Button variant="destructive" size="sm" @click="onDeleteTask">Delete</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Unlink confirmation -->
    <Dialog v-model:open="confirmUnlinkOpen">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>Unlink from ADO?</DialogTitle>
          <DialogDescription>
            This will disconnect the task from ADO work item #{{ task?.adoId }}. The local task will be kept.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" size="sm" @click="confirmUnlinkOpen = false">Cancel</Button>
          <Button variant="destructive" size="sm" @click="unlinkFromADO">Unlink</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </aside>
</template>
