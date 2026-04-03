<script setup lang="ts">
import { computed } from 'vue'
import type { Task } from '@/stores/tasks'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import TagChip from '@/components/ui/TagChip.vue'
import AdoBadge from '@/components/ui/AdoBadge.vue'
import {
  Circle,
  CircleDot,
  Eye,
  CheckCircle2,
  Octagon,
  XCircle,
} from 'lucide-vue-next'

const props = defineProps<{
  task: Task
  selected?: boolean
}>()

const emit = defineEmits<{
  select: [id: number]
  toggleStatus: [id: number]
}>()

const statusIcon = computed(() => {
  switch (props.task.status) {
    case 'todo': return Circle
    case 'in_progress': return CircleDot
    case 'in_review': return Eye
    case 'done': return CheckCircle2
    case 'blocked': return Octagon
    case 'cancelled': return XCircle
    default: return Circle
  }
})

const statusColor = computed(() => {
  switch (props.task.status) {
    case 'todo': return 'text-zinc-400 hover:text-zinc-500'
    case 'in_progress': return 'text-blue-500 hover:text-blue-600'
    case 'in_review': return 'text-violet-500 hover:text-violet-600'
    case 'done': return 'text-emerald-500 hover:text-emerald-600'
    case 'blocked': return 'text-red-500 hover:text-red-600'
    case 'cancelled': return 'text-zinc-400 hover:text-zinc-500'
    default: return 'text-zinc-400'
  }
})

const tags = computed(() => {
  if (!props.task.tags) return []
  return props.task.tags.split(',').map(t => t.trim()).filter(Boolean)
})

const visibleTags = computed(() => tags.value.slice(0, 2))
const overflowCount = computed(() => Math.max(0, tags.value.length - 2))

const timeAgo = computed(() => {
  const diff = Date.now() - new Date(props.task.updatedAt).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
})

const isDone = computed(() => props.task.status === 'done' || props.task.status === 'cancelled')

function onCheckClick(e: Event) {
  e.stopPropagation()
  emit('toggleStatus', props.task.id)
}
</script>

<template>
  <div
    @click="emit('select', task.id)"
    class="group flex items-start gap-3 px-4 py-2.5 cursor-pointer transition-all duration-100 relative"
    :class="[
      selected
        ? 'bg-primary/5 border-l-2 border-l-primary'
        : 'border-l-2 border-l-transparent hover:bg-muted/50',
    ]"
  >
    <!-- Status circle (Things 3 style) -->
    <button
      @click="onCheckClick"
      class="flex-shrink-0 mt-0.5 w-5 h-5 flex items-center justify-center rounded-full transition-colors duration-150"
      :class="statusColor"
    >
      <component :is="statusIcon" :size="18" :stroke-width="1.75" />
    </button>

    <!-- Content -->
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2 min-w-0">
        <!-- Title -->
        <span
          class="text-sm truncate flex-1 min-w-0"
          :class="[
            isDone ? 'line-through text-muted-foreground' : 'text-foreground',
          ]"
        >
          {{ task.title }}
        </span>

        <!-- Right side badges -->
        <div class="flex items-center gap-1.5 flex-shrink-0">
          <PriorityBadge :priority="task.priority" />

          <TagChip v-for="tag in visibleTags" :key="tag" :tag="tag" />
          <span
            v-if="overflowCount > 0"
            class="text-[11px] text-muted-foreground font-medium"
          >
            +{{ overflowCount }}
          </span>

          <AdoBadge :ado-id="task.adoId" />

          <span class="text-[11px] text-muted-foreground ml-1 tabular-nums whitespace-nowrap">
            {{ timeAgo }}
          </span>
        </div>
      </div>

      <!-- Blocked subtitle -->
      <div
        v-if="task.status === 'blocked' && task.blockedReason"
        class="mt-0.5 text-xs text-red-500/80 truncate"
      >
        ⊘ {{ task.blockedReason }}
      </div>
    </div>
  </div>
</template>
