<script setup lang="ts">
import { computed } from 'vue'
import type { Task } from '@/stores/tasks'
import { relativeTime, formatDate } from '@/lib/date'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import TagChip from '@/components/ui/TagChip.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import {
  Circle,
  CircleDot,
  Eye,
  CheckCircle2,
  Octagon,
  XCircle,
  CalendarDays,
} from 'lucide-vue-next'
import { statusColor } from '@/lib/styles'

const props = defineProps<{
  task: Task
  selected?: boolean
  isPublic?: boolean
  projectName?: string
}>()

const emit = defineEmits<{
  select: [id: number]
  toggleStatus: [id: number]
  'link-task': [id: number]
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

const tags= computed(() => {
  if (!props.task.tags) return []
  return props.task.tags.split(',').map(t => t.trim()).filter(Boolean)
})

const visibleTags = computed(() => tags.value.slice(0, 2))
const overflowCount = computed(() => Math.max(0, tags.value.length - 2))

const timeAgo = computed(() => relativeTime(props.task.updatedAt))

const isDone = computed(() => props.task.status === 'done' || props.task.status === 'cancelled')

const dueDateDisplay = computed(() => {
  if (!props.task.dueDate) return null
  const d = new Date(props.task.dueDate)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  if (d < today) return { text: formatDate(props.task.dueDate), overdue: true }
  if (d >= today && d < tomorrow) return { text: 'Today', overdue: false }
  return { text: formatDate(props.task.dueDate), overdue: false }
})

function onCheckClick(e: Event) {
  e.stopPropagation()
  emit('toggleStatus', props.task.id)
}

function onAdoBadgeClick(e: Event) {
  e.stopPropagation()
  if (!props.isPublic) {
    emit('link-task', props.task.id)
  }
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
      :class="statusColor(task.status)"
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
          <!-- ADO Badge: filled if public, hollow if personal -->
          <button
            @click="onAdoBadgeClick"
            class="inline-flex items-center justify-center rounded-full shrink-0 transition-colors"
            :class="isPublic
              ? 'w-5 h-5 bg-blue-500/10 text-blue-500 border border-blue-500/30 hover:bg-blue-500/20'
              : 'w-5 h-5 border border-dashed border-muted-foreground/30 text-muted-foreground/30 hover:border-muted-foreground/60 hover:text-muted-foreground/60'"
            :title="isPublic ? 'Linked to ADO' : 'Personal — click to link'"
          >
            <AzureDevOpsIcon v-if="isPublic" :size="12" />
            <Circle v-else :size="10" :stroke-width="1.5" />
          </button>

          <PriorityBadge :priority="task.priority" />

          <!-- Project tag -->
          <span
            v-if="projectName"
            class="text-[10px] px-1.5 py-0.5 rounded bg-muted text-muted-foreground truncate max-w-[6rem] shrink-0"
          >
            {{ projectName }}
          </span>

          <TagChip v-for="tag in visibleTags" :key="tag" :tag="tag" />
          <span
            v-if="overflowCount > 0"
            class="text-[11px] text-muted-foreground font-medium"
          >
            +{{ overflowCount }}
          </span>

          <!-- Due date -->
          <span
            v-if="dueDateDisplay"
            class="inline-flex items-center gap-0.5 text-[10px] shrink-0"
            :class="dueDateDisplay.overdue ? 'text-red-500' : 'text-muted-foreground'"
          >
            <CalendarDays :size="10" />
            {{ dueDateDisplay.text }}
          </span>

          <span class="text-[11px] text-muted-foreground ml-1 tabular-nums whitespace-nowrap">
            {{ timeAgo }}
          </span>
        </div>
      </div>

      <!-- Description preview (1 line) -->
      <div
        v-if="task.description && !isDone"
        class="mt-0.5 text-xs text-muted-foreground/60 truncate"
      >
        {{ task.description }}
      </div>

      <!-- Blocked subtitle -->
      <div
        v-if="task.status === 'blocked' && task.blockedReason"
        class="mt-0.5 text-xs text-red-500/80 truncate"
      >
        {{ task.blockedReason }}
      </div>
    </div>
  </div>
</template>
