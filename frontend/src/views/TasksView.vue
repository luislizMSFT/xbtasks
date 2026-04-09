<script setup lang="ts">
import { ref, onMounted, onActivated, onDeactivated, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMagicKeys, whenever } from '@vueuse/core'
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
import { adoTypeColor, adoTypeIcon, statusBgColor, priorityDotBgColor } from '@/lib/styles'
import { ClipboardList, ChevronRight, ChevronDown, Folder } from 'lucide-vue-next'
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

// Auto-expand all groups when groupBy changes or groups load for the first time
watch([() => taskStore.groupBy, groupKeys], () => {
  if (taskStore.groupBy && groupKeys.value.length > 0) {
    expandedGroups.value = new Set(groupKeys.value)
  }
}, { immediate: true })

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

// Task list header count
const taskListCount = computed(() => {
  if (treeView.value) {
    const roots = rootTasks.value.length
    const total = taskStore.enhancedFilteredTasks.length
    return `${roots} roots · ${total} tasks`
  }
  return `${taskStore.enhancedFilteredTasks.length} tasks`
})

function expandAll() {
  if (treeView.value) {
    expandedTreeNodes.value = new Set(taskStore.enhancedFilteredTasks.map(t => t.id))
  } else if (taskStore.groupBy) {
    expandedGroups.value = new Set(groupKeys.value)
  }
}

function collapseAll() {
  if (treeView.value) {
    expandedTreeNodes.value = new Set()
  } else if (taskStore.groupBy) {
    expandedGroups.value = new Set()
  }
}

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
    <!-- Teleport status chips + FilterCycleButton to top bar (only when this page is active) -->
    <Teleport v-if="isActive" to="#topbar-center">
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
        :sort-by="taskStore.sortBy"
        :group-by="taskStore.groupBy"
        :tree-view="treeView"
        :projects="projectStore.projects"
        :syncing="syncStore.syncing"
        @update:filter-status="taskStore.filterStatus = $event"
        @update:filter-priority="taskStore.filterPriority = $event"
        @update:filter-project="taskStore.filterProject = $event"
        @update:filter-due-date="taskStore.filterDueDate = $event"
        @update:sort-by="taskStore.sortBy = $event"
        @update:group-by="taskStore.groupBy = $event"
        @update:tree-view="treeView = $event"
        @sync="syncStore.manualSync()"
      />

      <!-- Task list header bar -->
      <div class="px-3 py-2 border-b border-border/50 flex items-center gap-2">
        <span class="text-xs font-semibold text-muted-foreground flex-1">{{ taskListCount }}</span>
        <template v-if="taskStore.groupBy || treeView">
          <button class="text-[10px] text-muted-foreground hover:text-foreground px-1.5 py-0.5 rounded hover:bg-muted"
            @click="expandAll">Expand All</button>
          <button class="text-[10px] text-muted-foreground hover:text-foreground px-1.5 py-0.5 rounded hover:bg-muted"
            @click="collapseAll">Collapse All</button>
        </template>
        <FilterCycleButton v-model="taskStore.filterAdoLink" />
      </div>

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
          <!-- Group sections -->
          <div v-for="key in groupKeys" :key="key">
            <!-- Project group header — matches TreeTaskRow styling -->
            <div
              v-if="taskStore.groupBy === 'project'"
              class="group cursor-pointer hover:bg-muted/50 transition-colors border-b border-border/30"
              :class="selectedProjectId === Number(key) && key !== 'No Project' ? 'bg-primary/5 border-l-2 border-l-primary' : 'border-l-2 border-l-transparent'"
              @click="selectProject(key)"
            >
              <div class="flex items-center gap-2 px-3 py-2.5">
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

                <!-- ADO type icon or folder icon -->
                <component
                  v-if="key !== 'No Project' && projectAdoLookup.get(Number(key))?.type"
                  :is="typeIcon(projectAdoLookup.get(Number(key))!.type)"
                  :size="16"
                  :class="adoTypeColor(projectAdoLookup.get(Number(key))!.type)"
                  class="shrink-0"
                />
                <Folder v-else :size="14" class="text-muted-foreground/50 shrink-0" />

                <!-- Project name -->
                <span class="text-sm font-medium text-foreground truncate">
                  {{ groupLabel(key) }}
                </span>

                <!-- ADO ID badge -->
                <Badge
                  v-if="key !== 'No Project' && projectAdoLookup.get(Number(key))?.number"
                  variant="outline"
                  class="text-[10px] px-1.5 py-0 h-4 font-mono text-blue-500 border-blue-500/30 shrink-0"
                >
                  {{ projectAdoLookup.get(Number(key))!.number }}
                </Badge>

                <span class="flex-1" />

                <!-- Task count -->
                <Badge variant="secondary" class="text-[10px] px-1.5 py-0 h-4 shrink-0">
                  {{ groupedTasks[key]?.length ?? 0 }} {{ (groupedTasks[key]?.length ?? 0) === 1 ? 'task' : 'tasks' }}
                </Badge>
              </div>
            </div>

            <!-- Status/priority group header -->
            <div
              v-else
              class="flex items-center gap-2.5 px-4 py-2.5 cursor-pointer select-none border-b border-border/30 hover:bg-muted/30 transition-colors"
              @click="expandedGroups.has(key) ? expandedGroups.delete(key) : expandedGroups.add(key)"
            >
              <button class="shrink-0 p-0.5 hover:bg-muted rounded">
                <component
                  :is="expandedGroups.has(key) ? ChevronDown : ChevronRight"
                  :size="14"
                  class="text-muted-foreground"
                />
              </button>

              <!-- Status dot -->
              <span
                v-if="taskStore.groupBy === 'status'"
                class="shrink-0 inline-block w-2.5 h-2.5 rounded-full"
                :class="statusBgColor(key)"
              />

              <!-- Priority dot -->
              <span
                v-if="taskStore.groupBy === 'priority'"
                class="shrink-0 inline-block w-2.5 h-2.5 rounded-full"
                :class="priorityDotBgColor(key)"
              />

              <span class="text-sm font-medium text-foreground truncate flex-1">
                {{ groupLabel(key) }}
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
        <div v-else-if="hasAnyTasks && !taskStore.loading && !treeView">
          <TreeTaskRow
            v-for="task in taskStore.enhancedFilteredTasks"
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
        </div>
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
