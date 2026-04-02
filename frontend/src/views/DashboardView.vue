<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import AdoBadge from '@/components/ui/AdoBadge.vue'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
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
  <ScrollArea class="flex-1 h-full">
    <div class="max-w-5xl mx-auto px-4 py-4 space-y-3">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold text-foreground">Dashboard</h1>
          <p class="text-sm text-muted-foreground mt-0.5">Your task overview at a glance</p>
        </div>
        <Button size="sm" @click="goToTasks">
          <Plus :size="14" />
          Create Task
        </Button>
      </div>

      <!-- Stat cards -->
      <div class="grid grid-cols-4 gap-2">
        <Card>
          <CardContent class="p-3">
            <div class="flex items-center gap-2 mb-1">
              <Loader2 :size="14" class="text-blue-500" />
              <span class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">In Progress</span>
            </div>
            <p class="text-xl font-bold text-foreground tabular-nums">{{ taskStore.stats.inProgress }}</p>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="p-3">
            <div class="flex items-center gap-2 mb-1">
              <Octagon :size="14" class="text-red-500" />
              <span class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Blocked</span>
            </div>
            <p class="text-xl font-bold text-foreground tabular-nums">{{ taskStore.stats.blocked }}</p>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="p-3">
            <div class="flex items-center gap-2 mb-1">
              <CheckCircle2 :size="14" class="text-emerald-500" />
              <span class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Done</span>
            </div>
            <p class="text-xl font-bold text-foreground tabular-nums">{{ taskStore.stats.done }}</p>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="p-3">
            <div class="flex items-center gap-2 mb-1">
              <Clock :size="14" class="text-zinc-500" />
              <span class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Total</span>
            </div>
            <p class="text-xl font-bold text-foreground tabular-nums">{{ taskStore.stats.total }}</p>
          </CardContent>
        </Card>
      </div>

      <!-- Today's Focus -->
      <Card v-if="focusTasks.length > 0">
        <CardHeader class="px-3 py-2">
          <CardTitle class="text-sm">Today's Focus</CardTitle>
        </CardHeader>
        <Separator />
        <CardContent class="p-0">
          <div
            v-for="(task, i) in focusTasks"
            :key="task.id"
            class="flex items-center gap-3 px-3 py-2 hover:bg-muted/50 cursor-pointer transition-colors"
            @click="router.push('/tasks')"
          >
            <StatusBadge :status="task.status" />
            <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
            <PriorityBadge :priority="task.priority" />
            <AdoBadge :ado-id="task.adoId" />
          </div>
        </CardContent>
      </Card>

      <!-- Blocked section -->
      <Card v-if="blockedTasks.length > 0" class="border-red-500/20 bg-red-500/5">
        <CardHeader class="px-3 py-2">
          <CardTitle class="text-sm flex items-center gap-2">
            <AlertTriangle :size="14" class="text-red-500" />
            <span class="text-red-600 dark:text-red-400">Blocked</span>
          </CardTitle>
        </CardHeader>
        <Separator class="bg-red-500/10" />
        <CardContent class="p-0">
          <div
            v-for="task in blockedTasks"
            :key="task.id"
            class="flex items-center gap-3 px-3 py-2 cursor-pointer transition-colors hover:bg-red-500/5"
            @click="router.push('/tasks')"
          >
            <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
            <AdoBadge :ado-id="task.adoId" />
            <span class="text-xs text-red-500/80">{{ task.blockedReason || 'Blocked' }}</span>
          </div>
        </CardContent>
      </Card>

      <!-- Recent Activity -->
      <Card>
        <CardHeader class="px-3 py-2">
          <CardTitle class="text-sm">Recent Activity</CardTitle>
        </CardHeader>
        <Separator />
        <CardContent class="p-0">
          <div
            v-for="task in recentTasks"
            :key="task.id"
            class="flex items-center gap-3 px-3 py-2 hover:bg-muted/50 cursor-pointer transition-colors"
            @click="router.push('/tasks')"
          >
            <StatusBadge :status="task.status" />
            <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
            <span class="text-[11px] text-muted-foreground tabular-nums">{{ timeAgo(task.updatedAt) }}</span>
          </div>
        </CardContent>
      </Card>

      <!-- Empty state -->
      <div v-if="taskStore.tasks.length === 0 && !taskStore.loading" class="text-center py-16">
        <p class="text-lg font-medium text-foreground">Welcome to Team ADO Tool</p>
        <p class="text-sm text-muted-foreground mt-1">Create your first task to get started</p>
        <Button class="mt-4" @click="goToTasks">
          Get Started
        </Button>
      </div>
    </div>
  </ScrollArea>
</template>
