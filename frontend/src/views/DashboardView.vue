<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '../stores/tasks'
import StatusBadge from '../components/ui/StatusBadge.vue'
import PriorityBadge from '../components/ui/PriorityBadge.vue'
import AdoBadge from '../components/ui/AdoBadge.vue'
import {
  Plus,
  AlertTriangle,
  Clock,
  CheckCircle2,
  Loader2,
  Octagon,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()

onMounted(() => {
  taskStore.fetchTasks()
})

const focusTasks = computed(() =>
  taskStore.tasks.filter(t => ['in_progress', 'in_review'].includes(t.status)).slice(0, 5)
)

const recentTasks = computed(() =>
  [...taskStore.tasks]
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    .slice(0, 5)
)

const blockedTasks = computed(() =>
  taskStore.tasks.filter(t => t.status === 'blocked')
)

function goToTasks() {
  router.push('/tasks')
}

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}
</script>

<template>
  <div class="flex-1 overflow-y-auto">
    <div class="max-w-4xl mx-auto px-6 py-6 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold text-text-primary">Dashboard</h1>
          <p class="text-sm text-text-secondary mt-0.5">Your task overview at a glance</p>
        </div>
        <button
          @click="goToTasks"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
        >
          <Plus :size="14" />
          Create Task
        </button>
      </div>

      <!-- Stat cards -->
      <div class="grid grid-cols-4 gap-3">
        <div class="rounded-xl border border-border-default bg-surface-secondary p-4">
          <div class="flex items-center gap-2 mb-2">
            <Loader2 :size="16" class="text-blue-500" />
            <span class="text-xs font-medium text-text-secondary uppercase tracking-wider">In Progress</span>
          </div>
          <p class="text-2xl font-bold text-text-primary tabular-nums">{{ taskStore.stats.inProgress }}</p>
        </div>

        <div class="rounded-xl border border-border-default bg-surface-secondary p-4">
          <div class="flex items-center gap-2 mb-2">
            <Octagon :size="16" class="text-red-500" />
            <span class="text-xs font-medium text-text-secondary uppercase tracking-wider">Blocked</span>
          </div>
          <p class="text-2xl font-bold text-text-primary tabular-nums">{{ taskStore.stats.blocked }}</p>
        </div>

        <div class="rounded-xl border border-border-default bg-surface-secondary p-4">
          <div class="flex items-center gap-2 mb-2">
            <CheckCircle2 :size="16" class="text-emerald-500" />
            <span class="text-xs font-medium text-text-secondary uppercase tracking-wider">Done</span>
          </div>
          <p class="text-2xl font-bold text-text-primary tabular-nums">{{ taskStore.stats.done }}</p>
        </div>

        <div class="rounded-xl border border-border-default bg-surface-secondary p-4">
          <div class="flex items-center gap-2 mb-2">
            <Clock :size="16" class="text-zinc-500" />
            <span class="text-xs font-medium text-text-secondary uppercase tracking-wider">Total</span>
          </div>
          <p class="text-2xl font-bold text-text-primary tabular-nums">{{ taskStore.stats.total }}</p>
        </div>
      </div>

      <!-- Today's Focus -->
      <div v-if="focusTasks.length > 0" class="rounded-xl border border-border-default bg-surface-secondary">
        <div class="px-4 py-3 border-b border-border-default">
          <h2 class="text-sm font-semibold text-text-primary">Today's Focus</h2>
        </div>
        <div>
          <div
            v-for="task in focusTasks"
            :key="task.id"
            class="flex items-center gap-3 px-4 py-2.5 hover:bg-surface-tertiary/50 cursor-pointer transition-colors border-b border-border-default/50 last:border-b-0"
            @click="router.push('/tasks')"
          >
            <StatusBadge :status="task.status" />
            <span class="text-sm text-text-primary flex-1 truncate">{{ task.title }}</span>
            <PriorityBadge :priority="task.priority" />
            <AdoBadge :ado-id="task.adoId" />
          </div>
        </div>
      </div>

      <!-- Blocked section -->
      <div v-if="blockedTasks.length > 0" class="rounded-xl border border-red-500/20 bg-red-500/5">
        <div class="px-4 py-3 border-b border-red-500/10">
          <div class="flex items-center gap-2">
            <AlertTriangle :size="16" class="text-red-500" />
            <h2 class="text-sm font-semibold text-red-600 dark:text-red-400">Blocked</h2>
          </div>
        </div>
        <div>
          <div
            v-for="task in blockedTasks"
            :key="task.id"
            class="flex items-center gap-3 px-4 py-2.5 cursor-pointer transition-colors border-b border-red-500/10 last:border-b-0 hover:bg-red-500/5"
            @click="router.push('/tasks')"
          >
            <span class="text-sm text-text-primary flex-1 truncate">{{ task.title }}</span>
            <AdoBadge :ado-id="task.adoId" />
            <span class="text-xs text-red-500/80">{{ task.blockedReason || 'Blocked' }}</span>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="rounded-xl border border-border-default bg-surface-secondary">
        <div class="px-4 py-3 border-b border-border-default">
          <h2 class="text-sm font-semibold text-text-primary">Recent Activity</h2>
        </div>
        <div>
          <div
            v-for="task in recentTasks"
            :key="task.id"
            class="flex items-center gap-3 px-4 py-2.5 hover:bg-surface-tertiary/50 cursor-pointer transition-colors border-b border-border-default/50 last:border-b-0"
            @click="router.push('/tasks')"
          >
            <StatusBadge :status="task.status" />
            <span class="text-sm text-text-primary flex-1 truncate">{{ task.title }}</span>
            <span class="text-[11px] text-text-secondary tabular-nums">{{ timeAgo(task.updatedAt) }}</span>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-if="taskStore.tasks.length === 0 && !taskStore.loading" class="text-center py-16">
        <p class="text-lg font-medium text-text-primary">Welcome to Team ADO Tool</p>
        <p class="text-sm text-text-secondary mt-1">Create your first task to get started</p>
        <button
          @click="goToTasks"
          class="mt-4 px-4 py-2 rounded-md text-sm font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
        >
          Get Started
        </button>
      </div>
    </div>
  </div>
</template>
