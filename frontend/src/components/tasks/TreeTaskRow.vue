<script setup lang="ts">
import { computed } from 'vue'
import type { Task } from '@/types'
import { Badge } from '@/components/ui/badge'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import { ChevronRight, ChevronDown, User, CalendarDays } from 'lucide-vue-next'
import { statusIcon, statusColor, adoTypeIcon, adoTypeColor, adoStateClasses, statusClasses } from '@/lib/styles'
import { formatDate } from '@/lib/date'

interface Props {
  task: Task
  indentLevel?: number
  selected?: boolean
  expanded?: boolean
  hasChildren?: boolean
  adoMeta?: { type: string; state: string } | null
  projectName?: string
  subtaskProgress?: { done: number; total: number; pct: number } | null
  isPublic?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  indentLevel: 0,
  selected: false,
  expanded: false,
  hasChildren: false,
  adoMeta: null,
  projectName: undefined,
  subtaskProgress: null,
  isPublic: false,
})

const emit = defineEmits<{
  click: []
  toggleExpand: []
  toggleStatus: [id: number]
}>()

const paddingLeft = computed(() =>
  props.indentLevel > 0 ? `${props.indentLevel * 24}px` : undefined,
)

const isPersonal = computed(() => !props.isPublic)

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

function statusLabel(status: string): string {
  switch (status) {
    case 'todo': return 'Todo'
    case 'in_progress': return 'In Progress'
    case 'in_review': return 'In Review'
    case 'done': return 'Done'
    case 'blocked': return 'Blocked'
    case 'cancelled': return 'Cancelled'
    default: return status
  }
}

function onStatusClick(e: Event) {
  e.stopPropagation()
  emit('toggleStatus', props.task.id)
}

function onExpandClick(e: Event) {
  e.stopPropagation()
  emit('toggleExpand')
}
</script>

<template>
  <div
    class="group cursor-pointer hover:bg-muted/50 transition-colors"
    :class="selected ? 'bg-primary/5 border-l-2 border-l-primary' : 'border-l-2 border-l-transparent'"
    :style="paddingLeft ? { paddingLeft } : undefined"
    @click="emit('click')"
  >
    <!-- Row 1: chevron + status + type icon + title + badges -->
    <div class="flex items-center gap-2 px-3 pt-2.5 pb-0.5">
      <!-- Expand button -->
      <button
        v-if="hasChildren"
        class="shrink-0 p-0.5 hover:bg-muted rounded"
        @click="onExpandClick"
      >
        <component :is="expanded ? ChevronDown : ChevronRight" :size="14" class="text-muted-foreground" />
      </button>
      <span v-else class="w-5 shrink-0" />

      <!-- Status icon -->
      <button
        class="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full"
        :class="statusColor(task.status)"
        @click="onStatusClick"
      >
        <component :is="statusIcon(task.status)" :size="16" :stroke-width="1.75" />
      </button>

      <!-- Type icon (ADO type or personal indicator) -->
      <component
        v-if="adoMeta?.type"
        :is="adoTypeIcon(adoMeta.type)"
        :size="14"
        :class="adoTypeColor(adoMeta.type)"
        class="shrink-0 -ml-1"
      />
      <User v-else-if="isPersonal" :size="12" class="text-muted-foreground/40 shrink-0 -ml-1" />

      <!-- Title -->
      <span
        class="text-sm truncate flex-1 min-w-0"
        :class="[
          isDone ? 'line-through text-muted-foreground' : 'text-foreground',
          hasChildren ? 'font-medium' : '',
        ]"
      >
        {{ task.title }}
      </span>

      <!-- State badge -->
      <Badge
        variant="outline"
        class="text-[10px] h-4 px-1.5 shrink-0"
        :class="adoMeta?.state ? adoStateClasses(adoMeta.state) : statusClasses(task.status)"
      >
        {{ adoMeta?.state || statusLabel(task.status) }}
      </Badge>

      <!-- Subtask progress -->
      <template v-if="subtaskProgress">
        <span class="text-[10px] text-muted-foreground tabular-nums shrink-0">
          {{ subtaskProgress.done }}/{{ subtaskProgress.total }}
        </span>
        <div class="w-12 h-1 bg-muted rounded-full overflow-hidden shrink-0">
          <div class="h-full bg-green-500 rounded-full transition-all" :style="{ width: subtaskProgress.pct + '%' }" />
        </div>
      </template>
    </div>

    <!-- Row 2: metadata -->
    <div class="flex items-center gap-2 pb-2 text-[11px]" :style="{ paddingLeft: '68px', paddingRight: '12px' }">
      <!-- ADO badge -->
      <span v-if="isPublic && task.adoId" class="text-muted-foreground/50 tabular-nums">
        {{ task.adoId }}
      </span>
      <Badge v-if="isPersonal" variant="outline" class="text-[9px] h-3.5 px-1 border-dashed text-muted-foreground/60">
        <User :size="8" class="mr-0.5" /> Personal
      </Badge>

      <!-- Project name -->
      <span v-if="projectName" class="text-muted-foreground/40 truncate">
        {{ projectName }}
      </span>

      <!-- Area -->
      <span v-if="task.area" class="text-muted-foreground/30">{{ task.area }}</span>

      <div class="flex-1" />

      <!-- Due date -->
      <span
        v-if="dueDateDisplay"
        class="inline-flex items-center gap-0.5 text-[10px] shrink-0"
        :class="dueDateDisplay.overdue ? 'text-red-500' : 'text-muted-foreground/50'"
      >
        <CalendarDays :size="10" />
        {{ dueDateDisplay.text }}
      </span>

      <!-- Priority badge -->
      <PriorityBadge :priority="task.priority" />
    </div>

    <!-- Row 3: blocked banner -->
    <div
      v-if="task.status === 'blocked' && task.blockedReason"
      class="pb-2 text-[10px] text-red-500/80 truncate"
      :style="{ paddingLeft: '68px', paddingRight: '12px' }"
    >
      ⚠ {{ task.blockedReason }}
    </div>
  </div>
</template>
