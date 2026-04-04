<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import TaskRow from '@/components/TaskRow.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Plus } from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()

onMounted(() => taskStore.fetchTasks())

// Today's Focus: tasks with status in_progress or in_review (per UI-SPEC)
const focusTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'in_progress' || t.status === 'in_review')
)

// Recent Activity: last 5 updated tasks, sorted by updatedAt DESC
const recentTasks = computed(() =>
  taskStore.tasks
    .slice()
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    .slice(0, 5)
)

// Blocked tasks — section hidden entirely if empty (per UI-SPEC)
const blockedTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'blocked')
)

// Stats for stat cards (uses taskStore.stats computed)
const statCards = computed(() => [
  { label: 'Total', value: taskStore.stats.total },
  { label: 'In Progress', value: taskStore.stats.inProgress },
  { label: 'Blocked', value: taskStore.stats.blocked },
  { label: 'Done', value: taskStore.stats.done },
])

function relativeTime(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

function goCreateTask() {
  router.push({ path: '/tasks', query: { create: '1' } })
}
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <div class="px-6 py-5 max-w-4xl">
      <!-- Page header: "Dashboard" + "+ Create Task" (per UI-SPEC) -->
      <div class="flex items-center justify-between mb-6">
        <h1 style="font-size: 20px; font-weight: 600; color: var(--color-text-primary)">Dashboard</h1>
        <Button @click="goCreateTask" class="gap-1.5">
          <Plus :size="14" />
          Create Task
        </Button>
      </div>

      <!-- Empty dashboard state (per UI-SPEC copywriting contract) -->
      <div v-if="taskStore.tasks.length === 0 && !taskStore.loading" class="text-center py-20">
        <h2 style="font-size: 20px; font-weight: 600; color: var(--color-text-primary)">Welcome to Team ADO Tool</h2>
        <p style="font-size: 14px; font-weight: 400; color: var(--color-text-secondary)" class="mt-2">
          Your dashboard will show today's focus, recent activity, and blocked items. Create a task to get started.
        </p>
        <Button @click="goCreateTask" class="mt-4 gap-1.5">
          <Plus :size="14" />
          Create Task
        </Button>
      </div>

      <!-- Dashboard content (when tasks exist) -->
      <template v-if="taskStore.tasks.length > 0">
        <!-- Stat cards row — focal point per UI-SPEC: 28px display values draw the eye -->
        <div class="grid grid-cols-4 gap-4 mb-8">
          <div
            v-for="card in statCards"
            :key="card.label"
            class="rounded-lg p-4"
            style="background: var(--color-bg-secondary); border: 1px solid var(--color-border-default)"
          >
            <div style="font-size: 28px; font-weight: 600; color: var(--color-text-primary)" class="tabular-nums">
              {{ card.value }}
            </div>
            <div style="font-size: 12px; font-weight: 600; color: var(--color-text-secondary)">
              {{ card.label }}
            </div>
          </div>
        </div>

        <!-- Today's Focus section — in_progress + in_review tasks -->
        <div class="mb-6">
          <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)" class="mb-3">
            Today's Focus
          </h2>
          <div
            v-if="focusTasks.length > 0"
            class="rounded-lg overflow-hidden"
            style="border: 1px solid var(--color-border-default)"
          >
            <TaskRow
              v-for="task in focusTasks"
              :key="task.id"
              :task="task"
              @select="(id) => { taskStore.selectTask(id); router.push('/tasks') }"
              @toggle-status="(id) => taskStore.setStatus(id, 'done')"
            />
          </div>
          <p
            v-else
            style="font-size: 14px; font-weight: 400; color: var(--color-text-tertiary)"
          >
            No tasks in progress — pick something to work on.
          </p>
        </div>

        <!-- Recent Activity section — last 5 modified tasks -->
        <div class="mb-6">
          <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)" class="mb-3">
            Recent Activity
          </h2>
          <div
            v-if="recentTasks.length > 0"
            class="rounded-lg overflow-hidden"
            style="border: 1px solid var(--color-border-default)"
          >
            <div
              v-for="(task, idx) in recentTasks"
              :key="task.id"
              class="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
              :style="idx < recentTasks.length - 1 ? 'border-bottom: 1px solid var(--color-border-default)' : ''"
              @click="() => { taskStore.selectTask(task.id); router.push('/tasks') }"
            >
              <StatusBadge :status="task.status" />
              <span class="text-sm flex-1 truncate" style="color: var(--color-text-primary)">{{ task.title }}</span>
              <span class="text-xs tabular-nums" style="color: var(--color-text-secondary)">{{ relativeTime(task.updatedAt) }}</span>
            </div>
          </div>
        </div>

        <!-- Blocked section — hidden entirely when no blocked tasks (per UI-SPEC) -->
        <div v-if="blockedTasks.length > 0">
          <h2 style="font-size: 14px; font-weight: 600; color: var(--color-text-primary)" class="mb-3">
            Blocked
          </h2>
          <div
            class="rounded-lg overflow-hidden"
            style="border: 1px solid var(--color-border-default); border-left: 2px solid var(--color-status-blocked)"
          >
            <TaskRow
              v-for="task in blockedTasks"
              :key="task.id"
              :task="task"
              @select="(id) => { taskStore.selectTask(id); router.push('/tasks') }"
            />
          </div>
        </div>
      </template>

      <!-- Loading state -->
      <div v-if="taskStore.loading" class="flex items-center justify-center py-20">
        <div class="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
      </div>
    </div>
  </ScrollArea>
</template>
