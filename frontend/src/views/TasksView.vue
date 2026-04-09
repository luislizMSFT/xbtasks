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
import { useAdoMeta } from '@/composables/useAdoMeta'
import TaskDetail from '@/components/tasks/TaskDetail.vue'
import ProjectDetail from '@/components/projects/ProjectDetail.vue'
import TreeTaskRow from '@/components/tasks/TreeTaskRow.vue'
import FilterBar from '@/components/FilterBar.vue'
import FilterCycleButton from '@/components/tasks/FilterCycleButton.vue'
import QuickAddInput from '@/components/tasks/QuickAddInput.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import { adoTypeColor, adoTypeIcon } from '@/lib/styles'
import { ClipboardList, RefreshCw, ChevronRight, ChevronDown } from 'lucide-vue-next'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/EmptyState.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'

const route = useRoute()
const taskStore = useTaskStore()
const projectStore = useProjectStore()
const adoStore = useADOStore()
const syncStore = useSyncStore()
const prStore = usePRStore()
const adoMeta = useAdoMeta()

const isActive = ref(false)
onActivated(() => { isActive.value = true })
onDeactivated(() => { isActive.value = false })
onMounted(() => { isActive.value = true })

const showQuickAdd = ref(false)
const expandedSubtasks = ref<Set<number>>(new Set())
const expandedGroups = ref<Set<string>>(new Set())
const selectedProjectId = ref<number | null>(null)
const treeView = ref(false)
const expandedTreeNodes = ref<Set<number>>(new Set())

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

// Project name lookup for TreeTaskRow
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

function handleQuickAdd(title: string) {
  taskStore.quickAdd(title)
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

// Reactive lookup: project → ADO info (type, adoId, #number)
const projectAdoLookup = computed(() => {
  const result = new Map<number, { type: string; adoId: string; number: string }>()
  const tree = adoStore.workItemTree
  const links = projectStore.projectLinks
  for (const [pid, lnk] of links) {
    const wi = tree.find(w => w.adoId === lnk.adoId)
    const num = lnk.adoId?.match(/\d+/)
    result.set(pid, {
      type: wi?.type || '',
      adoId: lnk.adoId,
      number: num ? `#${num[0]}` : lnk.adoId,
    })
  }
  return result
})

function typeIcon(type: string) {
  return adoTypeIcon(type)
}

const hasAnyTasks = computed(() => taskStore.tasks.length > 0)

// ── Tree view helpers ──
const rootTasks = computed(() =>
  taskStore.enhancedFilteredTasks.filter(t => !t.parentId)
)

function treeChildren(taskId: number): Task[] {
  return taskStore.enhancedFilteredTasks.filter(t => t.parentId === taskId)
}

function hasTreeChildren(taskId: number): boolean {
  return taskStore.enhancedFilteredTasks.some(t => t.parentId === taskId)
}

function toggleTreeNode(id: number) {
  if (expandedTreeNodes.value.has(id)) {
    expandedTreeNodes.value.delete(id)
  } else {
    expandedTreeNodes.value.add(id)
  }
}

function treeSubtaskProgress(taskId: number) {
  const children = treeChildren(taskId)
  if (children.length === 0) return null
  const done = children.filter(c => c.status === 'done').length
  return { done, total: children.length, pct: Math.round((done / children.length) * 100) }
}

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
    <!-- Teleport status chips + FilterCycleButton + sync to top bar (only when this page is active) -->
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
        <FilterCycleButton v-model="taskStore.filterAdoLink" />
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

    <!-- LEFT: Task list panel -->
    <div class="flex flex-col min-w-0 w-[55%]">
      <!-- FilterBar (dropdowns only, no status chips or sync) -->
      <FilterBar
        :filter-status="taskStore.filterStatus"
        :filter-priority="taskStore.filterPriority"
        :filter-project="taskStore.filterProject"
        :filter-due-date="taskStore.filterDueDate"
        :filter-ado-link="taskStore.filterAdoLink"
        :sort-by="taskStore.sortBy"
        :group-by="taskStore.groupBy"
        :tree-view="treeView"
        :projects="projectStore.projects"
        :syncing="syncStore.syncing"
        @update:filter-status="taskStore.filterStatus = $event"
        @update:filter-priority="taskStore.filterPriority = $event"
        @update:filter-project="taskStore.filterProject = $event"
        @update:filter-due-date="taskStore.filterDueDate = $event"
        @update:filter-ado-link="taskStore.filterAdoLink = $event"
        @update:sort-by="taskStore.sortBy = $event"
        @update:group-by="taskStore.groupBy = $event"
        @update:tree-view="treeView = $event"
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
          <QuickAddInput data-quick-add @add="handleQuickAdd" @cancel="() => {}" />
        </div>

        <!-- Grouped rendering -->
        <div v-if="hasAnyTasks && !taskStore.loading && taskStore.groupBy && !treeView">
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
                v-if="taskStore.groupBy === 'project' && projectAdoLookup.get(Number(key))?.type"
                :is="typeIcon(projectAdoLookup.get(Number(key))!.type)"
                :size="14"
                :class="adoTypeColor(projectAdoLookup.get(Number(key))!.type)"
                class="shrink-0"
              />
              <span class="text-sm font-medium text-foreground truncate flex-1" @click="selectProject(key)">
                {{ groupLabel(key) }}
              </span>
              <span
                v-if="taskStore.groupBy === 'project' && projectStore.isLinked(Number(key))"
                class="inline-flex items-center gap-1 text-[10px] text-blue-500 shrink-0"
              >
                <component
                  v-if="projectAdoLookup.get(Number(key))?.type"
                  :is="typeIcon(projectAdoLookup.get(Number(key))!.type)"
                  :size="12"
                  :class="adoTypeColor(projectAdoLookup.get(Number(key))!.type)"
                />
                <AzureDevOpsIcon v-else :size="12" />
                <span class="tabular-nums">{{ projectAdoLookup.get(Number(key))?.number ?? '' }}</span>
              </span>
              <span class="text-[11px] text-muted-foreground tabular-nums shrink-0">
                {{ groupedTasks[key]?.length ?? 0 }}
              </span>
            </div>

            <!-- Task rows in group -->
            <template v-if="expandedGroups.has(key)">
              <TreeTaskRow
                v-for="task in groupedTasks[key]"
                :key="task.id"
                :task="task"
                :is-public="taskStore.isPublic(task.id)"
                :selected="taskStore.selectedTaskId === task.id"
                :ado-meta="adoMeta.getAdoMeta(task.id)"
                :project-name="task.projectId ? projectNameMap[task.projectId] : undefined"
                :subtask-progress="subtaskProgress(task.id)"
                @click="selectTask(task.id)"
                @toggle-status="(id: number) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
              />
            </template>
          </div>

          <!-- Ungrouped tasks (no project) — flat list, no header -->
          <template v-if="groupedTasks['No Project']?.length">
            <TreeTaskRow
              v-for="task in groupedTasks['No Project']"
              :key="task.id"
              :task="task"
              :is-public="taskStore.isPublic(task.id)"
              :selected="taskStore.selectedTaskId === task.id"
              :ado-meta="adoMeta.getAdoMeta(task.id)"
              :project-name="undefined"
              :subtask-progress="subtaskProgress(task.id)"
              @click="selectTask(task.id)"
              @toggle-status="(id: number) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
            />
          </template>
        </div>

        <!-- Tree rendering (hierarchical parent → child nesting via TreeTaskRow with indentLevel) -->
        <div v-else-if="hasAnyTasks && !taskStore.loading && treeView">
          <template v-for="task in rootTasks" :key="task.id">
            <!-- Root node (level 0) -->
            <TreeTaskRow
              :task="task"
              :indent-level="0"
              :selected="taskStore.selectedTaskId === task.id"
              :expanded="expandedTreeNodes.has(task.id)"
              :has-children="hasTreeChildren(task.id)"
              :ado-meta="adoMeta.getAdoMeta(task.id)"
              :project-name="task.projectId ? projectNameMap[task.projectId] : undefined"
              :subtask-progress="treeSubtaskProgress(task.id)"
              :is-public="taskStore.isPublic(task.id)"
              @click="selectTask(task.id)"
              @toggle-expand="toggleTreeNode(task.id)"
              @toggle-status="(id: number) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
            />

            <!-- Children (level 1) -->
            <template v-if="expandedTreeNodes.has(task.id)">
              <template v-for="child in treeChildren(task.id)" :key="child.id">
                <TreeTaskRow
                  :task="child"
                  :indent-level="1"
                  :selected="taskStore.selectedTaskId === child.id"
                  :expanded="expandedTreeNodes.has(child.id)"
                  :has-children="hasTreeChildren(child.id)"
                  :ado-meta="adoMeta.getAdoMeta(child.id)"
                  :project-name="child.projectId ? projectNameMap[child.projectId] : undefined"
                  :subtask-progress="treeSubtaskProgress(child.id)"
                  :is-public="taskStore.isPublic(child.id)"
                  @click="selectTask(child.id)"
                  @toggle-expand="toggleTreeNode(child.id)"
                  @toggle-status="(id: number) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
                />

                <!-- Grandchildren (level 2) -->
                <template v-if="expandedTreeNodes.has(child.id)">
                  <TreeTaskRow
                    v-for="gc in treeChildren(child.id)"
                    :key="gc.id"
                    :task="gc"
                    :indent-level="2"
                    :selected="taskStore.selectedTaskId === gc.id"
                    :ado-meta="adoMeta.getAdoMeta(gc.id)"
                    :project-name="gc.projectId ? projectNameMap[gc.projectId] : undefined"
                    :subtask-progress="subtaskProgress(gc.id)"
                    :is-public="taskStore.isPublic(gc.id)"
                    @click="selectTask(gc.id)"
                    @toggle-status="(id: number) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
                  />
                </template>
              </template>
            </template>
          </template>
        </div>

        <!-- Non-grouped flat rendering (enhanced filtered + sorted, drag & drop) -->
        <draggable
          v-else-if="hasAnyTasks && !taskStore.loading && !treeView"
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
                <TreeTaskRow
                  :task="task"
                  :is-public="taskStore.isPublic(task.id)"
                  :selected="taskStore.selectedTaskId === task.id"
                  :ado-meta="adoMeta.getAdoMeta(task.id)"
                  :project-name="task.projectId ? projectNameMap[task.projectId] : undefined"
                  :subtask-progress="subtaskProgress(task.id)"
                  @click="selectTask(task.id)"
                  @toggle-status="(id: number) => { const t = taskStore.tasks.find(x => x.id === id); if (t) toggleDone(t) }"
                />
              </div>
            </div>
          </template>
        </draggable>
      </ScrollArea>
    </div>

    <!-- RIGHT: Detail panel (permanent, not slide-out) -->
    <div class="w-[45%] shrink-0 border-l border-border flex flex-col min-h-0">
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
      <div v-else class="flex-1 flex items-center justify-center">
        <p class="text-sm text-muted-foreground">Select a task or project</p>
      </div>
    </div>
  </div>
</template>
