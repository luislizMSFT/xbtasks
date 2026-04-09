<script setup lang="ts">
import { computed } from 'vue'
import type { Task } from '@/types'
import { priorityBgColor, adoTypeIcon, adoTypeColor } from '@/lib/styles'
import { Badge } from '@/components/ui/badge'
import { SquareCheckBig } from 'lucide-vue-next'

interface Props {
  task: Task
  isPersonal?: boolean
  showDueDate?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isPersonal: true,
  showDueDate: false,
})

const emit = defineEmits<{
  click: [taskId: number]
}>()

// Determine ADO type from task — if adoId is set, task is linked to ADO
// The adoType is stored in task.category for imported items, or defaults to null
const hasAdoLink = computed(() => !!props.task.adoId)

// Due date formatting
const formattedDueDate = computed(() => {
  if (!props.task.dueDate) return ''
  const d = new Date(props.task.dueDate)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
})

// Due date is "soon" if within 2 days from now
const isDueSoon = computed(() => {
  if (!props.task.dueDate) return false
  const diff = (new Date(props.task.dueDate).getTime() - Date.now()) / (1000 * 60 * 60 * 24)
  return diff <= 2
})

// Is overdue (past due date)
const isOverdue = computed(() => {
  if (!props.task.dueDate) return false
  return new Date(props.task.dueDate) < new Date()
})

// Status badge label and style
const statusLabel = computed(() => {
  if (props.task.status === 'in_review') return 'Review'
  if (props.task.status === 'in_progress') return 'Active'
  return props.task.status
})

const statusClass = computed(() => {
  if (props.task.status === 'in_review') return 'border-violet-500/30 text-violet-600'
  return 'border-blue-500/30 text-blue-600'
})
</script>

<template>
  <div
    class="flex items-center gap-3 px-4 py-2 cursor-pointer hover:bg-muted/50 transition-colors"
    @click="emit('click', task.id)"
  >
    <!-- Priority dot (8px) -->
    <span
      class="size-2 rounded-full shrink-0"
      :class="priorityBgColor(task.priority)"
    />

    <!-- ADO type icon OR personal SquareCheckBig -->
    <component
      v-if="hasAdoLink && task.category"
      :is="adoTypeIcon(task.category)"
      :size="14"
      :class="adoTypeColor(task.category)"
      class="shrink-0"
    />
    <SquareCheckBig v-else :size="14" class="text-primary/60 shrink-0" />

    <!-- Task title (truncated, flex-1) -->
    <span class="text-sm flex-1 truncate text-foreground">{{ task.title }}</span>

    <!-- Personal badge (only for personal/non-linked tasks) -->
    <span
      v-if="isPersonal"
      class="text-[8px] px-1 py-0.5 rounded bg-primary/8 text-primary/70 border border-primary/10 shrink-0"
    >
      personal
    </span>

    <!-- Status badge -->
    <Badge
      variant="outline"
      class="text-[9px] h-4 px-1.5 shrink-0"
      :class="statusClass"
    >
      {{ statusLabel }}
    </Badge>

    <!-- Due date (optional, amber if due soon, red if overdue) -->
    <span
      v-if="showDueDate && task.dueDate"
      class="text-[11px] tabular-nums shrink-0"
      :class="[
        isOverdue ? 'text-red-600 font-medium' :
        isDueSoon ? 'text-amber-600 font-medium' :
        'text-muted-foreground'
      ]"
    >
      {{ formattedDueDate }}
    </span>
  </div>
</template>
