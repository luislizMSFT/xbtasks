<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useTaskStore, type Task } from '@/stores/tasks'
import TaskDetail from '@/components/TaskDetail.vue'
import { cn } from '@/lib/utils'
import {
  Filter, ArrowUpDown, ChevronDown, ChevronRight,
  Bug, CheckSquare, BookOpen, Landmark,
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

// Priority → left stripe color
function priorityStripe(p: string) {
  switch (p) {
    case 'P0': return 'bg-red-500'
    case 'P1': return 'bg-orange-500'
    case 'P2': return 'bg-amber-400'
    default: return 'bg-zinc-300 dark:bg-zinc-600'
  }
}

// Status → dot color for the inline status indicator
function statusDot(s: string) {
  switch (s) {
    case 'in_progress': return 'bg-blue-500'
    case 'in_review': return 'bg-violet-500'
    case 'todo': return 'bg-zinc-400'
    case 'blocked': return 'bg-red-500'
    case 'done': return 'bg-emerald-500'
    case 'cancelled': return 'bg-zinc-300'
    default: return 'bg-zinc-400'
  }
}

// ADO type → icon component
function adoIcon(adoId: string) {
  const lower = adoId.toLowerCase()
  if (lower.includes('bug')) return Bug
  if (lower.includes('story')) return BookOpen
  return CheckSquare
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

onMounted(() => {
  taskStore.fetchTasks()
})
</script>

<template>
  <div class="flex-1 flex overflow-hidden">
    <!-- Left: Task list -->
    <div class="flex-1 flex flex-col min-w-0" :class="taskStore.selectedTask ? 'w-[55%]' : ''">
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

        <!-- Inline create -->
        <div v-if="showNewTask" class="px-4 py-2 border-b border-border bg-card/50">
          <div class="flex items-center gap-3">
            <div class="w-[3px] self-stretch rounded-full bg-zinc-300 dark:bg-zinc-600" />
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
        </div>

        <!-- Grouped sections -->
        <div v-if="hasAnyTasks && !taskStore.loading">
          <div v-for="section in visibleSections" :key="section" class="border-b border-border/50 last:border-b-0">
            <!-- Section header -->
            <button
              @click="toggleSection(section)"
              class="flex items-center gap-2 w-full px-4 py-1.5 text-left hover:bg-muted/30 transition-colors"
            >
              <component
                :is="isCollapsed(section) ? ChevronRight : ChevronDown"
                :size="13"
                class="text-muted-foreground"
              />
              <span class="w-2 h-2 rounded-full" :class="statusDot(section)" />
              <span class="text-[11px] font-semibold uppercase tracking-wider" :class="sectionColors[section]">
                {{ sectionLabels[section] }}
              </span>
              <span class="text-[11px] text-muted-foreground tabular-nums">
                {{ taskStore.grouped[section]?.length ?? 0 }}
              </span>
            </button>

            <!-- Task rows -->
            <div v-if="!isCollapsed(section)">
              <div
                v-for="task in taskStore.grouped[section]"
                :key="task.id"
                @click="selectTask(task.id)"
                :class="cn(
                  'group flex items-center gap-2 px-3 py-2 cursor-pointer transition-all border-b border-border/30 last:border-b-0',
                  'hover:bg-muted/40',
                  taskStore.selectedTaskId === task.id
                    ? 'bg-primary/[0.06]'
                    : ''
                )"
              >
                <!-- Priority stripe -->
                <div :class="cn('w-[3px] self-stretch rounded-full shrink-0', priorityStripe(task.priority))" />

                <!-- Status dot -->
                <div :class="cn('w-2 h-2 rounded-full shrink-0', statusDot(task.status))" />

                <!-- Title + ADO inline -->
                <div class="flex-1 min-w-0 flex items-center gap-1.5">
                  <span
                    :class="cn(
                      'text-[13px] truncate',
                      task.status === 'done' ? 'text-muted-foreground line-through' : 'text-foreground'
                    )"
                  >
                    {{ task.title }}
                  </span>
                  <!-- ADO type icon + number -->
                  <span
                    v-if="task.adoId"
                    class="inline-flex items-center gap-0.5 text-[10px] text-blue-600 dark:text-blue-400 shrink-0"
                  >
                    <component :is="adoIcon(task.adoId)" :size="10" :stroke-width="2" />
                    {{ task.adoId.replace('ADO-', '#') }}
                  </span>
                </div>

                <!-- Blocked reason -->
                <span
                  v-if="task.status === 'blocked' && task.blockedReason"
                  class="text-[10px] text-red-500/70 truncate max-w-[8rem] shrink-0"
                >
                  {{ task.blockedReason }}
                </span>

                <!-- Time -->
                <span class="text-[10px] text-muted-foreground/40 tabular-nums shrink-0">
                  {{ timeAgo(task.updatedAt) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </ScrollArea>
    </div>

    <!-- Right: Detail preview panel (Todoist hybrid) -->
    <TaskDetail
      v-if="taskStore.selectedTask"
      @close="closeDetail"
    />
  </div>
</template>
