<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import {
  Circle,
  CircleDot,
  Eye,
  CheckCircle2,
  Octagon,
  XCircle,
} from 'lucide-vue-next'

const props = defineProps<{ status: string }>()

// Uses semantic --color-status-* tokens from UI-SPEC (defined in style.css @theme)
const config = computed(() => {
  switch (props.status) {
    case 'todo': return { icon: Circle, variant: 'outline' as const, label: 'To Do', colorVar: '--color-status-todo' }
    case 'in_progress': return { icon: CircleDot, variant: 'default' as const, label: 'In Progress', colorVar: '--color-status-in-progress' }
    case 'in_review': return { icon: Eye, variant: 'default' as const, label: 'In Review', colorVar: '--color-status-in-review' }
    case 'done': return { icon: CheckCircle2, variant: 'default' as const, label: 'Done', colorVar: '--color-status-done' }
    case 'blocked': return { icon: Octagon, variant: 'default' as const, label: 'Blocked', colorVar: '--color-status-blocked' }
    case 'cancelled': return { icon: XCircle, variant: 'outline' as const, label: 'Cancelled', colorVar: '--color-status-cancelled' }
    default: return { icon: Circle, variant: 'outline' as const, label: props.status, colorVar: '--color-status-todo' }
  }
})

const badgeStyle = computed(() => ({
  color: `var(${config.value.colorVar})`,
  backgroundColor: `color-mix(in srgb, var(${config.value.colorVar}) 10%, transparent)`,
  borderColor: `color-mix(in srgb, var(${config.value.colorVar}) 20%, transparent)`,
}))
</script>

<template>
  <Badge
    :variant="config.variant"
    :class="cn('inline-flex items-center gap-1 text-[11px] font-medium border', '')"
    :style="badgeStyle"
  >
    <component :is="config.icon" :size="13" :stroke-width="2" />
    <span class="leading-none">{{ config.label }}</span>
  </Badge>
</template>
