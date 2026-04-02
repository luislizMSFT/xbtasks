<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useTaskStore } from '../stores/tasks'
import TaskRow from '../components/TaskRow.vue'
import TaskDetail from '../components/TaskDetail.vue'
import { Plus, Filter, ArrowUpDown, ChevronDown, ChevronRight } from 'lucide-vue-next'

const taskStore = useTaskStore()

const showNewTask = ref(false)
const newTaskTitle = ref('')
const newTaskInput = ref<HTMLInputElement | null>(null)
const collapsedSections = ref<Set<string>>(new Set())

const tabs = [
  { id: 'all', label: 'All' },
  { id: 'active', label: 'Active' },
  { id: 'done', label: 'Done' },
  { id: 'blocked', label: 'Blocked' },
]

const sectionOrder = ['in_progress', 'in_review', 'todo', 'blocked', 'done', 'cancelled'] as const

const sectionLabels: Record<string, string> = {
  in_progress: 'In Progress',
  in_review: 'In Review',
  todo: 'To Do',
  blocked: 'Blocked',
  done: 'Done',
  cancelled: 'Cancelled',
}

const sectionColors: Record<string, string> = {
  in_progress: 'text-blue-500',
  in_review: 'text-violet-500',
  todo: 'text-zinc-500',
  blocked: 'text-red-500',
  done: 'text-emerald-500',
  cancelled: 'text-zinc-400',
}

const visibleSections = computed(() => {
  return sectionOrder.filter(s => {
    const tasks = taskStore.grouped[s]
    return tasks && tasks.length > 0
  })
})

function toggleSection(section: string) {
  if (collapsedSections.value.has(section)) {
    collapsedSections.value.delete(section)
  } else {
    collapsedSections.value.add(section)
  }
}

function isCollapsed(section: string) {
  return collapsedSections.value.has(section)
}

async function createTask() {
  const title = newTaskTitle.value.trim()
  if (!title) return
  await taskStore.createTask(title)
  newTaskTitle.value = ''
  showNewTask.value = false
}

function startInlineCreate() {
  showNewTask.value = true
  setTimeout(() => newTaskInput.value?.focus(), 50)
}

function onToggleStatus(id: number) {
  const task = taskStore.tasks.find(t => t.id === id)
  if (!task) return
  const nextStatus = task.status === 'done' ? 'todo' : 'done'
  taskStore.setStatus(id, nextStatus)
}

function selectTask(id: number) {
  taskStore.selectTask(taskStore.selectedTaskId === id ? null : id)
}

function closeDetail() {
  taskStore.selectTask(null)
}

const hasAnyTasks = computed(() => taskStore.tasks.length > 0)

onMounted(() => {
  taskStore.fetchTasks()
})
</script>

<template>
  <div class="flex-1 flex overflow-hidden">
    <!-- Task list -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Header -->
      <div class="flex items-center justify-between px-5 py-3 border-b border-border-default">
        <div class="flex items-center gap-1">
          <h1 class="text-lg font-semibold text-text-primary mr-4">Tasks</h1>
          <!-- Status tabs (pill style) -->
          <div class="flex items-center gap-0.5 bg-surface-tertiary rounded-lg p-0.5">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              @click="taskStore.filterStatus = tab.id"
              class="px-3 py-1 rounded-md text-xs font-medium transition-colors duration-150"
              :class="[
                taskStore.filterStatus === tab.id
                  ? 'bg-surface-primary text-text-primary shadow-sm'
                  : 'text-text-secondary hover:text-text-primary'
              ]"
            >
              {{ tab.label }}
            </button>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs text-text-secondary hover:text-text-primary hover:bg-surface-tertiary transition-colors">
            <Filter :size="14" />
            Filter
          </button>
          <button class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs text-text-secondary hover:text-text-primary hover:bg-surface-tertiary transition-colors">
            <ArrowUpDown :size="14" />
            Sort
          </button>
          <button
            @click="startInlineCreate"
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
          >
            <Plus :size="14" />
            New Task
          </button>
        </div>
      </div>

      <!-- Task list body -->
      <div class="flex-1 overflow-y-auto">
        <!-- Loading -->
        <div v-if="taskStore.loading" class="flex items-center justify-center py-20">
          <div class="w-5 h-5 border-2 border-accent/30 border-t-accent rounded-full animate-spin" />
        </div>

        <!-- Empty state -->
        <div v-else-if="!hasAnyTasks" class="flex flex-col items-center justify-center py-20 gap-3">
          <div class="w-12 h-12 rounded-full bg-surface-tertiary flex items-center justify-center">
            <Plus :size="24" class="text-text-secondary" />
          </div>
          <p class="text-sm font-medium text-text-primary">No tasks yet</p>
          <p class="text-xs text-text-secondary">Create your first task to get started</p>
          <button
            @click="startInlineCreate"
            class="mt-2 px-4 py-2 rounded-md text-sm font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
          >
            Create Task
          </button>
        </div>

        <!-- Inline create -->
        <div v-if="showNewTask" class="px-4 py-2 border-b border-border-default bg-surface-secondary/50">
          <div class="flex items-center gap-3">
            <div class="w-5 h-5 rounded-full border-2 border-zinc-300 dark:border-zinc-600" />
            <input
              ref="newTaskInput"
              v-model="newTaskTitle"
              @keydown.enter="createTask"
              @keydown.esc="showNewTask = false"
              class="flex-1 text-sm bg-transparent border-none outline-none text-text-primary placeholder-text-secondary"
              placeholder="What needs to be done?"
            />
            <button
              @click="createTask"
              class="px-2.5 py-1 rounded-md text-xs font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
            >
              Add
            </button>
            <button
              @click="showNewTask = false"
              class="px-2.5 py-1 rounded-md text-xs text-text-secondary hover:text-text-primary hover:bg-surface-tertiary transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>

        <!-- Grouped sections (Things 3 style) -->
        <div v-if="hasAnyTasks && !taskStore.loading">
          <div v-for="section in visibleSections" :key="section" class="border-b border-border-default/50 last:border-b-0">
            <!-- Section header -->
            <button
              @click="toggleSection(section)"
              class="flex items-center gap-2 w-full px-5 py-2 text-left hover:bg-surface-tertiary/30 transition-colors"
            >
              <component
                :is="isCollapsed(section) ? ChevronRight : ChevronDown"
                :size="14"
                class="text-text-secondary"
              />
              <span
                class="text-[12px] font-semibold uppercase tracking-wider"
                :class="sectionColors[section]"
              >
                {{ sectionLabels[section] }}
              </span>
              <span class="text-[11px] text-text-secondary">
                {{ taskStore.grouped[section]?.length ?? 0 }}
              </span>
            </button>

            <!-- Section tasks -->
            <TransitionGroup
              v-if="!isCollapsed(section)"
              tag="div"
              enter-active-class="transition-all duration-200 ease-out"
              enter-from-class="opacity-0 -translate-y-1"
              enter-to-class="opacity-100 translate-y-0"
              leave-active-class="transition-all duration-150 ease-in"
              leave-from-class="opacity-100"
              leave-to-class="opacity-0 -translate-y-1"
            >
              <TaskRow
                v-for="task in taskStore.grouped[section]"
                :key="task.id"
                :task="task"
                :selected="taskStore.selectedTaskId === task.id"
                @select="selectTask"
                @toggle-status="onToggleStatus"
              />
            </TransitionGroup>
          </div>
        </div>
      </div>
    </div>

    <!-- Detail panel -->
    <TaskDetail
      v-if="taskStore.selectedTask"
      @close="closeDetail"
    />
  </div>
</template>
