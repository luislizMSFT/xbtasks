<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useTaskStore, type Task } from '@/stores/tasks'
import TaskDetail from '@/components/TaskDetail.vue'
import { cn } from '@/lib/utils'
import {
  Filter, ArrowUpDown, ChevronDown, ChevronRight,
  Bug, CheckSquare, BookOpen, Landmark,
  Circle, CheckCircle2, Plus, GitPullRequest,
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'
import PageHeader from '@/components/PageHeader.vue'

const taskStore = useTaskStore()

const showNewTask = ref(false)
const newTaskTitle = ref('')
const newTaskInput = ref<HTMLInputElement | null>(null)
const expandedSubtasks = ref<Set<number>>(new Set())

const tabs = [
  { id: 'all', label: 'All' },
  { id: 'active', label: 'Active' },
  { id: 'done', label: 'Done' },
  { id: 'blocked', label: 'Blocked' },
]

const sectionOrder = ['in_progress', 'in_review', 'todo', 'blocked', 'done'] as const

const sectionMeta: Record<string, { label: string; dot: string }> = {
  in_progress: { label: 'In Progress', dot: 'bg-blue-500' },
  in_review:   { label: 'In Review',   dot: 'bg-violet-500' },
  todo:        { label: 'To Do',       dot: 'bg-zinc-400' },
  blocked:     { label: 'Blocked',     dot: 'bg-red-500' },
  done:        { label: 'Done',        dot: 'bg-emerald-500' },
}

// Mock subtasks per task (until real subtask backend)
const mockSubtasks: Record<number, { id: number; title: string; done: boolean }[]> = {
  1: [
    { id: 101, title: 'Identify refresh token edge case', done: true },
    { id: 102, title: 'Add token expiry middleware', done: true },
    { id: 103, title: 'Update error handling in auth callback', done: false },
    { id: 104, title: 'Add E2E test for token refresh', done: false },
  ],
  4: [
    { id: 401, title: 'Test Create method', done: false },
    { id: 402, title: 'Test Delete method', done: false },
    { id: 403, title: 'Test GetBlockers method', done: false },
  ],
}

// Mock PR counts per task
const mockPrCounts: Record<number, number> = { 1: 2, 3: 1 }

function adoIcon(adoId: string) {
  const lower = adoId.toLowerCase()
  if (lower.includes('bug')) return Bug
  if (lower.includes('story')) return BookOpen
  return CheckSquare
}

function subtasksFor(taskId: number) {
  return mockSubtasks[taskId] || []
}

function subtaskProgress(taskId: number) {
  const subs = subtasksFor(taskId)
  if (!subs.length) return null
  const done = subs.filter(s => s.done).length
  return { done, total: subs.length, pct: Math.round((done / subs.length) * 100) }
}

function toggleSubtasks(taskId: number) {
  if (expandedSubtasks.value.has(taskId)) {
    expandedSubtasks.value.delete(taskId)
  } else {
    expandedSubtasks.value.add(taskId)
  }
}

function toggleSubtaskDone(taskId: number, subId: number) {
  const subs = mockSubtasks[taskId]
  if (!subs) return
  const sub = subs.find(s => s.id === subId)
  if (sub) sub.done = !sub.done
}

const visibleSections = computed(() => {
  return sectionOrder.filter(s => {
    const tasks = taskStore.grouped[s]
    return tasks && tasks.length > 0
  })
})

async function toggleDone(task: Task) {
  const newStatus = task.status === 'done' ? 'todo' : 'done'
  await taskStore.setStatus(task.id, newStatus)
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

function selectTask(id: number) {
  taskStore.selectTask(taskStore.selectedTaskId === id ? null : id)
}

function closeDetail() {
  taskStore.selectTask(null)
}

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'now'
  if (mins < 60) return `${mins}m`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h`
  const days = Math.floor(hours / 24)
  return `${days}d`
}

const hasAnyTasks = computed(() => taskStore.tasks.length > 0)

onMounted(async () => {
  await taskStore.fetchTasks()
  // Auto-select first task so detail panel is always visible
  if (taskStore.tasks.length > 0 && !taskStore.selectedTaskId) {
    taskStore.selectTask(taskStore.tasks[0].id)
  }
})
</script>

<template>
  <div class="flex-1 flex overflow-hidden">
    <!-- Left: Task list -->
    <div class="flex-1 flex flex-col min-w-0 w-[55%]">
      <PageHeader>
        <template #left>
          <Tabs :default-value="taskStore.filterStatus" @update:model-value="(v) => { if (v) taskStore.filterStatus = String(v) }">
            <TabsList>
              <TabsTrigger v-for="tab in tabs" :key="tab.id" :value="tab.id">
                {{ tab.label }}
              </TabsTrigger>
            </TabsList>
          </Tabs>
        </template>
        <template #right>
          <Button variant="ghost" size="sm" class="h-7 gap-1.5 text-xs">
            <Filter :size="13" />
            Filter
          </Button>
          <Button variant="ghost" size="sm" class="h-7 gap-1.5 text-xs">
            <ArrowUpDown :size="13" />
            Sort
          </Button>
        </template>
      </PageHeader>

      <ScrollArea class="flex-1 h-full">
        <!-- Loading -->
        <div v-if="taskStore.loading" class="flex items-center justify-center py-20">
          <div class="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
        </div>

        <!-- Empty state -->
        <div v-else-if="!hasAnyTasks" class="flex flex-col items-center justify-center py-20 gap-3">
          <p class="text-sm font-medium text-foreground">No tasks yet</p>
          <p class="text-xs text-muted-foreground">Create your first task to get started</p>
          <Button class="mt-2" @click="startInlineCreate">Create Task</Button>
        </div>

        <!-- Add task (always at top) -->
        <div class="px-4 py-2 border-b border-border/50">
          <div v-if="showNewTask" class="flex items-center gap-3">
            <Circle :size="16" class="text-muted-foreground/30 shrink-0" />
            <Input
              ref="newTaskInput"
              v-model="newTaskTitle"
              @keydown.enter="createTask"
              @keydown.esc="showNewTask = false"
              class="flex-1 border-none shadow-none focus-visible:ring-0 text-sm bg-transparent"
              placeholder="What needs to be done?"
            />
            <Button size="sm" class="h-7 text-xs" @click="createTask">Add</Button>
            <Button variant="ghost" size="sm" class="h-7 text-xs" @click="showNewTask = false">Cancel</Button>
          </div>
          <button
            v-else
            @click="startInlineCreate"
            class="flex items-center gap-2 text-[12px] text-muted-foreground/40 hover:text-muted-foreground transition-colors w-full py-0.5"
          >
            <Plus :size="14" />
            Add task
          </button>
        </div>

        <!-- Grouped sections -->
        <div v-if="hasAnyTasks && !taskStore.loading">
          <div v-for="section in visibleSections" :key="section">
            <!-- Subtle section divider -->
            <div class="flex items-center gap-2 px-4 pt-4 pb-1">
              <span :class="cn('w-1.5 h-1.5 rounded-full', sectionMeta[section].dot)" />
              <span class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground/60">
                {{ sectionMeta[section].label }}
              </span>
              <span class="text-[10px] text-muted-foreground/40 tabular-nums">
                {{ taskStore.grouped[section]?.length ?? 0 }}
              </span>
              <div class="flex-1 h-px bg-border/50" />
            </div>

            <!-- Task rows -->
            <div>
              <div v-for="task in taskStore.grouped[section]" :key="task.id">
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

                  <!-- ADO badge -->
                  <span
                    v-if="task.adoId"
                    class="inline-flex items-center gap-0.5 text-[10px] text-blue-600 dark:text-blue-400 bg-blue-500/8 px-1.5 py-0.5 rounded shrink-0"
                  >
                    <component :is="adoIcon(task.adoId)" :size="10" :stroke-width="2" />
                    {{ task.adoId.replace('ADO-', '#') }}
                  </span>

                  <!-- PR count -->
                  <span
                    v-if="mockPrCounts[task.id]"
                    class="inline-flex items-center gap-0.5 text-[10px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded shrink-0"
                  >
                    <GitPullRequest :size="10" />
                    {{ mockPrCounts[task.id] }}
                  </span>

                  <!-- Blocked reason -->
                  <span
                    v-if="task.status === 'blocked' && task.blockedReason"
                    class="text-[10px] text-red-500/70 truncate max-w-[6rem] shrink-0"
                  >
                    {{ task.blockedReason }}
                  </span>

                  <!-- Time -->
                  <span class="text-[10px] text-muted-foreground/40 tabular-nums shrink-0 w-6 text-right">
                    {{ timeAgo(task.updatedAt) }}
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
                        sub.done
                          ? 'bg-emerald-500 border-emerald-500'
                          : 'border-muted-foreground/25 hover:border-muted-foreground/50'
                      )"
                    >
                      <CheckCircle2 v-if="sub.done" :size="8" class="text-white" :stroke-width="3" />
                    </button>
                    <span
                      :class="cn(
                        'text-[12px]',
                        sub.done ? 'text-muted-foreground/50 line-through decoration-muted-foreground/20' : 'text-muted-foreground'
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
