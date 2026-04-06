<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMagicKeys, whenever } from '@vueuse/core'
import { useTaskStore, type Task } from '@/stores/tasks'
import { useProjectStore } from '@/stores/projects'
import { useSyncStore } from '@/stores/sync'
import { usePRStore } from '@/stores/prs'
import TaskDetail from '@/components/TaskDetail.vue'
import TaskRow from '@/components/TaskRow.vue'
import FilterBar from '@/components/FilterBar.vue'
import { cn } from '@/lib/utils'
import {
  ChevronDown, ChevronRight,
  CheckCircle2, Plus,
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'
import PageHeader from '@/components/PageHeader.vue'

const route = useRoute()
const taskStore = useTaskStore()
const projectStore = useProjectStore()
const syncStore = useSyncStore()
const prStore = usePRStore()

const quickAddTitle = ref('')
const expandedSubtasks = ref<Set<number>>(new Set())

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
  taskStore.selectTask(taskStore.selectedTaskId === id ? null : id)
}

function closeDetail() {
  taskStore.selectTask(null)
}

const hasAnyTasks = computed(() => taskStore.tasks.length > 0)

onMounted(async () => {
  await Promise.all([
    taskStore.fetchTasks(),
    projectStore.fetchProjects(),
    prStore.fetchAll(),
  ])
  // Auto-select first task so detail panel is always visible
  if (taskStore.tasks.length > 0 && !taskStore.selectedTaskId) {
    taskStore.selectTask(taskStore.tasks[0].id)
  }
})
</script>

<template>
  <div class="flex-1 flex overflow-hidden" @keydown.esc="onEscape" tabindex="-1">
    <!-- Left: Task list -->
    <div class="flex-1 flex flex-col min-w-0 w-[55%]">
      <PageHeader>
        <template #left>
          <span class="text-sm font-semibold text-foreground">Tasks</span>
        </template>
        <template #right>
          <span v-if="syncStore.lastSyncedAt" class="text-[10px] text-muted-foreground/50 mr-2">
            Last synced {{ new Date(syncStore.lastSyncedAt).toLocaleTimeString() }}
          </span>
        </template>
      </PageHeader>

      <!-- FilterBar -->
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
        <!-- Loading -->
        <div v-if="taskStore.loading" class="flex items-center justify-center py-20">
          <div class="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
        </div>

        <!-- Empty state -->
        <div v-else-if="!hasAnyTasks" class="flex flex-col items-center justify-center py-20 gap-3">
          <p class="text-sm font-medium text-foreground">No tasks yet</p>
          <p class="text-xs text-muted-foreground">Create your first task to start tracking your work. Press Enter in the quick-add bar above.</p>
        </div>

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
          <div v-for="key in groupKeys" :key="key" class="mb-2">
            <!-- Group header -->
            <div class="flex items-center gap-2 px-4 pt-4 pb-1">
              <span class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground/60">
                {{ groupLabel(key) }}
              </span>
              <span class="text-[10px] text-muted-foreground/40 tabular-nums">
                {{ groupedTasks[key]?.length ?? 0 }}
              </span>
              <div class="flex-1 h-px bg-border/50" />
            </div>

            <!-- Task rows in group -->
            <div v-for="task in groupedTasks[key]" :key="task.id">
              <TaskRow
                :task="task"
                :is-public="taskStore.isPublic(task.id)"
                :selected="taskStore.selectedTaskId === task.id"
                :project-name="task.projectId ? projectNameMap[task.projectId] : undefined"
                @select="selectTask"
                @toggle-status="(id) => toggleDone(taskStore.tasks.find(t => t.id === id)!)"
              />
            </div>
          </div>
        </div>

        <!-- Non-grouped rendering (enhanced filtered + sorted) -->
        <div v-else-if="hasAnyTasks && !taskStore.loading">
          <div v-for="task in taskStore.enhancedFilteredTasks" :key="task.id">
            <div>
              <!-- Main row -->
              <div
                @click="selectTask(task.id)"
                :class="cn(
                  'group flex items-center gap-2.5 px-4 py-2 cursor-pointer transition-all',
                  'hover:bg-muted/40',
                  taskStore.selectedTaskId === task.id ? 'bg-primary/[0.06]' : ''
                )"
              >
                <!-- Checkbox -->
                <button
                  @click.stop="toggleDone(task)"
                  :class="cn(
                    'size-[18px] rounded-full border-[1.5px] shrink-0 flex items-center justify-center transition-all hover:scale-110',
                    task.status === 'done'
                      ? 'bg-emerald-500 border-emerald-500'
                      : task.status === 'blocked'
                        ? 'border-red-400 hover:border-red-500'
                        : 'border-muted-foreground/30 hover:border-muted-foreground/60'
                  )"
                >
                  <CheckCircle2 v-if="task.status === 'done'" :size="10" class="text-white" :stroke-width="3" />
                </button>

                <!-- Subtask expand toggle -->
                <button
                  v-if="subtasksFor(task.id).length > 0"
                  @click.stop="toggleSubtasks(task.id)"
                  class="shrink-0 text-muted-foreground/40 hover:text-muted-foreground transition-colors"
                >
                  <component :is="expandedSubtasks.has(task.id) ? ChevronDown : ChevronRight" :size="14" />
                </button>
                <div v-else class="w-[14px] shrink-0" />

                <!-- Title -->
                <span
                  :class="cn(
                    'text-[13px] font-medium truncate flex-1',
                    task.status === 'done' ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground'
                  )"
                >
                  {{ task.title }}
                </span>

                <!-- Subtask progress bar -->
                <div v-if="subtaskProgress(task.id)" class="flex items-center gap-1.5 shrink-0">
                  <div class="w-14 h-[3px] rounded-full bg-muted overflow-hidden">
                    <div
                      class="h-full rounded-full transition-all duration-300"
                      :class="subtaskProgress(task.id)!.pct === 100 ? 'bg-emerald-500' : 'bg-blue-500'"
                      :style="{ width: subtaskProgress(task.id)!.pct + '%' }"
                    />
                  </div>
                  <span class="text-[10px] text-muted-foreground/50 tabular-nums w-7">
                    {{ subtaskProgress(task.id)!.done }}/{{ subtaskProgress(task.id)!.total }}
                  </span>
                </div>

                <!-- ADO badge (personal/public) -->
                <span
                  v-if="taskStore.isPublic(task.id)"
                  class="inline-flex items-center gap-0.5 text-[10px] text-blue-600 dark:text-blue-400 bg-blue-500/8 px-1.5 py-0.5 rounded shrink-0"
                  title="Linked to ADO"
                >
                  ADO
                </span>
                <span
                  v-else
                  class="inline-flex items-center justify-center w-5 h-5 rounded-full border border-dashed border-muted-foreground/20 shrink-0"
                  title="Personal"
                />

                <!-- Project tag -->
                <span
                  v-if="task.projectId && projectNameMap[task.projectId]"
                  class="text-[10px] px-1.5 py-0.5 rounded bg-muted text-muted-foreground truncate max-w-[6rem] shrink-0"
                >
                  {{ projectNameMap[task.projectId] }}
                </span>
              </div>

              <!-- Expanded subtasks -->
              <div v-if="expandedSubtasks.has(task.id) && subtasksFor(task.id).length > 0" class="pl-[52px] pr-4 pb-1">
                <div
                  v-for="sub in subtasksFor(task.id)"
                  :key="sub.id"
                  class="flex items-center gap-2 py-1 group/sub"
                >
                  <button
                    @click="toggleSubtaskDone(task.id, sub.id)"
                    :class="cn(
                      'size-[14px] rounded-[3px] border-[1.5px] shrink-0 flex items-center justify-center transition-all hover:scale-110',
                      sub.status === 'done'
                        ? 'bg-emerald-500 border-emerald-500'
                        : 'border-muted-foreground/25 hover:border-muted-foreground/50'
                    )"
                  >
                    <CheckCircle2 v-if="sub.status === 'done'" :size="8" class="text-white" :stroke-width="3" />
                  </button>
                  <span
                    :class="cn(
                      'text-[12px]',
                      sub.status === 'done' ? 'text-muted-foreground/50 line-through decoration-muted-foreground/20' : 'text-muted-foreground'
                    )"
                  >
                    {{ sub.title }}
                  </span>
                </div>
                <button class="flex items-center gap-1.5 text-[11px] text-muted-foreground/40 hover:text-muted-foreground mt-0.5 py-1 transition-colors">
                  <Plus :size="11" />
                  Add subtask
                </button>
              </div>
            </div>
          </div>
        </div>
      </ScrollArea>
    </div>

    <!-- Right: Detail panel (always visible) -->
    <TaskDetail
      v-if="taskStore.selectedTask"
      @close="closeDetail"
    />
    <div v-else class="w-[45%] shrink-0 border-l border-border flex items-center justify-center">
      <p class="text-sm text-muted-foreground">Select a task to view details</p>
    </div>
  </div>
</template>
