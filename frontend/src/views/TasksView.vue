<script setup lang="ts">
import { ref, onMounted, onActivated, onDeactivated, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMagicKeys, whenever } from '@vueuse/core'
import draggable from 'vuedraggable'
import { useTaskStore, type Task } from '@/stores/tasks'
import { useProjectStore } from '@/stores/projects'
import { useADOStore } from '@/stores/ado'
import { useSyncStore } from '@/stores/sync'
import { usePRStore } from '@/stores/prs'
import TaskDetail from '@/components/tasks/TaskDetail.vue'
import ProjectDetail from '@/components/projects/ProjectDetail.vue'
import TaskRow from '@/components/tasks/TaskRow.vue'
import FilterBar from '@/components/FilterBar.vue'
import { adoTypeColor, adoTypeIcon, adoStateClasses, adoPriorityClasses } from '@/lib/styles'
import {
  ClipboardList, RefreshCw, ChevronRight, ChevronDown,
} from 'lucide-vue-next'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/EmptyState.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'

const route = useRoute()
const taskStore = useTaskStore()
const projectStore = useProjectStore()
const adoStore = useADOStore()
const syncStore = useSyncStore()
const prStore = usePRStore()

const isActive = ref(false)
onActivated(() => { isActive.value = true })
onDeactivated(() => { isActive.value = false })
onMounted(() => { isActive.value = true })

const quickAddTitle = ref('')
const expandedSubtasks = ref<Set<number>>(new Set())
const expandedGroups = ref<Set<string>>(new Set())
const selectedProjectId = ref<number | null>(null)

const statusChips = [
  { id: 'all', label: 'All' },
  { id: 'todo', label: 'Todo' },
  { id: 'in_progress', label: 'In Progress' },
  { id: 'in_review', label: 'In Review' },
  { id: 'done', label: 'Done' },
  { id: 'blocked', label: 'Blocked' },
]

// Keyboard shortcuts: ⌘N / Ctrl+N to toggle create form, Escape to close panels
const { Meta_n, Ctrl_n } = useMagicKeys()

whenever(Meta_n, () => {
  if (taskStore.selectedTaskId) {
    closeDetail()
  }
})

whenever(Ctrl_n, () => {
  if (taskStore.selectedTaskId) {
    closeDetail()
  }
})

// Handle Escape: close detail panel
function onEscape() {
  if (taskStore.selectedTaskId) {
    closeDetail()
  }
}

// Handle ?create=1 query param from Dashboard "Create Task" button
watch(() => route.query.create, (val) => {
  if (val === '1') {
    // Focus quick-add input
    const input = document.querySelector<HTMLInputElement>('[data-quick-add]')
    input?.focus()
  }
}, { immediate: true })

// Project name lookup for TaskRow
const projectNameMap = computed(() => {
  const map: Record<number, string> = {}
  for (const p of projectStore.projects) {
    map[p.id] = p.name
  }
  return map
})

// Compute subtask counts from task store (tasks with parentId matching)
const subtaskCounts = computed(() => {
  const counts: Record<number, { done: number; total: number }> = {}
  for (const t of taskStore.tasks) {
    if (t.parentId) {
      if (!counts[t.parentId]) counts[t.parentId] = { done: 0, total: 0 }
      counts[t.parentId].total++
      if (t.status === 'done') counts[t.parentId].done++
    }
  }
  return counts
})

// Real subtasks from the task store (children of a given parent)
const subtasksByParent = computed(() => {
  const map: Record<number, Task[]> = {}
  for (const t of taskStore.tasks) {
    if (t.parentId) {
      if (!map[t.parentId]) map[t.parentId] = []
      map[t.parentId].push(t)
    }
  }
  return map
})

function subtasksFor(taskId: number): Task[] {
  return subtasksByParent.value[taskId] || []
}

function subtaskProgress(taskId: number) {
  const counts = subtaskCounts.value[taskId]
  if (!counts || counts.total === 0) return null
  return { done: counts.done, total: counts.total, pct: Math.round((counts.done / counts.total) * 100) }
}

function toggleSubtasks(taskId: number) {
  if (expandedSubtasks.value.has(taskId)) {
    expandedSubtasks.value.delete(taskId)
  } else {
    expandedSubtasks.value.add(taskId)
  }
}

async function toggleSubtaskDone(taskId: number, subId: number) {
  const sub = subtasksFor(taskId).find(s => s.id === subId)
  if (!sub) return
  const newStatus = sub.status === 'done' ? 'todo' : 'done'
  await taskStore.setStatus(subId, newStatus)
}

// Grouped tasks for group-by mode
const groupedTasks = computed(() => {
  if (!taskStore.groupBy) return {} as Record<string, Task[]>
  return taskStore.groupedEnhanced
})

const groupKeys = computed(() => {
  if (!taskStore.groupBy) return []
  // Order groups sensibly
  if (taskStore.groupBy === 'status') {
    const order = ['in_progress', 'in_review', 'todo', 'blocked', 'done', 'cancelled']
    return order.filter(k => groupedTasks.value[k]?.length > 0)
  }
  if (taskStore.groupBy === 'priority') {
    const order = ['P0', 'P1', 'P2', 'P3', 'None']
    return order.filter(k => groupedTasks.value[k]?.length > 0)
  }
  // project grouping — sort alphabetically, 'No Project' last
  const keys = Object.keys(groupedTasks.value).filter(k => (groupedTasks.value[k]?.length ?? 0) > 0)
  return keys.sort((a, b) => {
    if (a === 'No Project') return 1
    if (b === 'No Project') return -1
    return a.localeCompare(b)
  })
})

// Resolve group label (project IDs to names)
function groupLabel(key: string): string {
  if (taskStore.groupBy === 'project' && key !== 'No Project') {
    const id = Number(key)
    return projectNameMap.value[id] || `Project ${key}`
  }
  // Status labels
  const statusLabels: Record<string, string> = {
    in_progress: 'In Progress', in_review: 'In Review',
    todo: 'To Do', blocked: 'Blocked', done: 'Done', cancelled: 'Cancelled',
  }
  return statusLabels[key] || key
}

async function toggleDone(task: Task) {
  const newStatus = task.status === 'done' ? 'todo' : 'done'
  await taskStore.setStatus(task.id, newStatus)
}

async function handleQuickAdd() {
  const title = quickAddTitle.value.trim()
  if (!title) return
  await taskStore.quickAdd(title)
  quickAddTitle.value = ''
}

function selectTask(id: number) {
  selectedProjectId.value = null
  taskStore.selectTask(taskStore.selectedTaskId === id ? null : id)
}

function selectProject(key: string) {
  if (key === 'No Project') return
  const id = Number(key)
  taskStore.selectTask(null)
  selectedProjectId.value = selectedProjectId.value === id ? null : id
}

function closeDetail() {
  taskStore.selectTask(null)
  selectedProjectId.value = null
}

function projectAdoType(projectId: number): string {
  const lnk = projectStore.projectLinks.get(projectId)
  if (!lnk) return ''
  const wi = adoStore.workItemTree.find(w => w.adoId === lnk.adoId)
  return wi?.type || ''
}

function projectAdoItem(projectId: number) {
  const lnk = projectStore.projectLinks.get(projectId)
  if (!lnk) return null
  return adoStore.workItemTree.find(w => w.adoId === lnk.adoId) ?? null
}

function typeIcon(type: string) {
  return adoTypeIcon(type)
}

const hasAnyTasks = computed(() => taskStore.tasks.length > 0)

// Writable list for drag & drop — syncs back to store on reorder
const draggableList = computed({
  get: () => taskStore.enhancedFilteredTasks,
  set: (val: Task[]) => {
    // Update local order immediately
    const ids = val.map(t => t.id)
    // Apply new sortOrder to the underlying store tasks
    for (let i = 0; i < val.length; i++) {
      const idx = taskStore.tasks.findIndex(t => t.id === val[i].id)
      if (idx !== -1) taskStore.tasks[idx].sortOrder = i
    }
    // Persist to backend
    taskStore.reorderTasks(ids)
  },
})

const isDragging = ref(false)

onMounted(async () => {
  // Data is fetched by App.vue on auth — only fetch if stores are empty
  if (!taskStore.tasks.length) await taskStore.fetchTasks()
  if (!projectStore.projects.length) projectStore.fetchProjects()
  // Auto-select first task so detail panel is always visible
  if (taskStore.tasks.length > 0 && !taskStore.selectedTaskId) {
    taskStore.selectTask(taskStore.tasks[0].id)
  }
})
</script>

<template>
  <div class="flex-1 flex overflow-hidden" @keydown.esc="onEscape" tabindex="-1">
    <!-- Teleport status chips + sync to top bar (only when this page is active) -->
    <Teleport v-if="isActive" to="#topbar-actions">
      <div class="flex items-center gap-1">
        <Badge
          v-for="chip in statusChips"
          :key="chip.id"
          :variant="taskStore.filterStatus === chip.id ? 'default' : 'outline'"
          class="cursor-pointer text-[10px] px-1.5 py-0 h-5 select-none"
          @click="taskStore.filterStatus = chip.id"
        >
          {{ chip.label }}
        </Badge>
        <Button
          variant="ghost"
          size="icon"
          class="h-6 w-6 ml-1"
          title="Sync with ADO"
          @click="syncStore.manualSync()"
        >
          <RefreshCw :size="12" :class="{ 'animate-spin': syncStore.syncing }" />
        </Button>
      </div>
    </Teleport>

    <!-- Left: Task list -->
    <div class="flex-1 flex flex-col min-w-0 w-[55%]">
      <!-- FilterBar (dropdowns only, no status chips or sync) -->
      <FilterBar
        :filter-status="taskStore.filterStatus"
        :filter-priority="taskStore.filterPriority"
        :filter-project="taskStore.filterProject"
        :filter-due-date="taskStore.filterDueDate"
        :filter-ado-link="taskStore.filterAdoLink"
        :sort-by="taskStore.sortBy"
        :group-by="taskStore.groupBy"
        :projects="projectStore.projects"
        :syncing="syncStore.syncing"
        @update:filter-status="taskStore.filterStatus = $event"
        @update:filter-priority="taskStore.filterPriority = $event"
        @update:filter-project="taskStore.filterProject = $event"
        @update:filter-due-date="taskStore.filterDueDate = $event"
        @update:filter-ado-link="taskStore.filterAdoLink = $event"
        @update:sort-by="taskStore.sortBy = $event"
        @update:group-by="taskStore.groupBy = $event"
        @sync="syncStore.manualSync()"
      />

      <ScrollArea class="flex-1 h-full">
        <!-- Loading skeleton -->
        <div v-if="taskStore.loading" class="px-4 py-3 space-y-2">
          <div v-for="i in 6" :key="i" class="flex items-center gap-2.5 px-4 py-2">
            <Skeleton class="size-[18px] rounded-full shrink-0" />
            <Skeleton class="w-[14px] h-[14px] shrink-0" />
            <Skeleton class="h-4 flex-1" />
            <Skeleton class="h-4 w-10 shrink-0 rounded" />
            <Skeleton class="h-5 w-5 rounded-full shrink-0" />
            <Skeleton class="h-4 w-16 shrink-0 rounded" />
          </div>
        </div>

        <!-- Empty state -->
        <EmptyState
          v-else-if="!hasAnyTasks"
          :icon="ClipboardList"
          title="No tasks yet"
          description="Create your first task to start tracking your work. Press Enter in the quick-add bar above."
        />

        <!-- Quick-add input (always at top) -->
        <div class="px-4 py-2 border-b border-border/50">
          <Input
            v-model="quickAddTitle"
            data-quick-add
            placeholder="Quick add task — press Enter"
            class="h-8 text-sm"
            @keydown.enter="handleQuickAdd"
          />
        </div>

        <!-- Grouped rendering -->
        <div v-if="hasAnyTasks && !taskStore.loading && taskStore.groupBy">
          <!-- Project groups (with headers) -->
          <div v-for="key in groupKeys.filter(k => k !== 'No Project')" :key="key">
            <!-- Project group header -->
            <div
              class="flex items-center gap-2.5 px-4 py-2.5 cursor-pointer select-none border-b border-border/30 hover:bg-muted/30 transition-colors"
              :class="selectedProjectId === Number(key) && 'bg-primary/5 border-l-2 border-l-primary'"
              @click.self="selectProject(key)"
            >
              <button
                class="shrink-0 p-0.5 hover:bg-muted rounded"
                @click.stop="expandedGroups.has(key) ? expandedGroups.delete(key) : expandedGroups.add(key)"
              >
                <component
                  :is="expandedGroups.has(key) ? ChevronDown : ChevronRight"
                  :size="14"
                  class="text-muted-foreground"
                />
              </button>
              <component
                v-if="taskStore.groupBy === 'project' && projectAdoType(Number(key))"
                :is="typeIcon(projectAdoType(Number(key)))"
                :size="14"
                :class="adoTypeColor(projectAdoType(Number(key)))"
                class="shrink-0"
              />
              <span class="text-sm font-medium text-foreground truncate flex-1" @click="selectProject(key)">
                {{ groupLabel(key) }}
              </span>
              <Badge
                v-if="taskStore.groupBy === 'project' && projectStore.isLinked(Number(key))"
                variant="outline"
                class="text-[9px] h-4 px-1 text-blue-500 border-blue-500/30 shrink-0"
              >ADO</Badge>
              <span class="text-[11px] text-muted-foreground tabular-nums shrink-0">
                {{ groupedTasks[key]?.length ?? 0 }}
              </span>
            </div>

            <!-- Task rows in group -->
            <template v-if="expandedGroups.has(key)">
              <div v-for="task in groupedTasks[key]" :key="task.id">
                <TaskRow
                  :task="task"
                  :is-public="taskStore.isPublic(task.id)"
                  :selected="taskStore.selectedTaskId === task.id"
                  @select="selectTask"
                  @toggle-status="(id) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
                />
              </div>
            </template>
          </div>

          <!-- Ungrouped tasks (no project) — flat list, no header -->
          <template v-if="groupedTasks['No Project']?.length">
            <div v-for="task in groupedTasks['No Project']" :key="task.id">
              <TaskRow
                :task="task"
                :is-public="taskStore.isPublic(task.id)"
                :selected="taskStore.selectedTaskId === task.id"
                @select="selectTask"
                @toggle-status="(id) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
              />
            </div>
          </template>
        </div>

        <!-- Non-grouped rendering (enhanced filtered + sorted, drag & drop) -->
        <draggable
          v-else-if="hasAnyTasks && !taskStore.loading"
          v-model="draggableList"
          item-key="id"
          handle=".drag-handle"
          ghost-class="opacity-30"
          :animation="150"
          @start="isDragging = true"
          @end="isDragging = false"
        >
          <template #item="{ element: task }">
            <div class="flex items-center">
              <!-- Drag handle -->
              <div class="drag-handle cursor-grab active:cursor-grabbing shrink-0 pl-2 text-muted-foreground/0 hover:text-muted-foreground/30 transition-colors">
                <svg width="10" height="14" viewBox="0 0 10 14" fill="currentColor">
                  <circle cx="2" cy="2" r="1.5" /><circle cx="8" cy="2" r="1.5" />
                  <circle cx="2" cy="7" r="1.5" /><circle cx="8" cy="7" r="1.5" />
                  <circle cx="2" cy="12" r="1.5" /><circle cx="8" cy="12" r="1.5" />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <TaskRow
                  :task="task"
                  :is-public="taskStore.isPublic(task.id)"
                  :selected="taskStore.selectedTaskId === task.id"
                  :project-name="task.projectId ? projectNameMap[task.projectId] : undefined"
                  @select="selectTask"
                  @toggle-status="(id) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
                />
              </div>
            </div>
          </template>
        </draggable>
      </ScrollArea>
    </div>

    <!-- Right: Detail panel (always visible) -->
    <ProjectDetail
      v-if="selectedProjectId"
      :project-id="selectedProjectId"
      @close="selectedProjectId = null"
      @select-task="selectTask"
    />
    <TaskDetail
      v-else-if="taskStore.selectedTask"
      @close="closeDetail"
    />
    <div v-else class="w-[45%] shrink-0 border-l border-border flex items-center justify-center">
      <p class="text-sm text-muted-foreground">Select a task or project</p>
    </div>
  </div>
</template>
