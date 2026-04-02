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

const config = computed(() => {
  switch (props.status) {
    case 'todo': return { icon: Circle, variant: 'outline' as const, classes: 'bg-zinc-500/15 text-zinc-500 border-zinc-500/20', label: 'To Do' }
    case 'in_progress': return { icon: CircleDot, variant: 'default' as const, classes: 'bg-blue-500/15 text-blue-600 border-blue-500/20', label: 'In Progress' }
    case 'in_review': return { icon: Eye, variant: 'default' as const, classes: 'bg-violet-500/15 text-violet-600 border-violet-500/20', label: 'In Review' }
    case 'done': return { icon: CheckCircle2, variant: 'default' as const, classes: 'bg-emerald-500/15 text-emerald-600 border-emerald-500/20', label: 'Done' }
    case 'blocked': return { icon: Octagon, variant: 'default' as const, classes: 'bg-red-500/15 text-red-600 border-red-500/20', label: 'Blocked' }
    case 'cancelled': return { icon: XCircle, variant: 'outline' as const, classes: 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20', label: 'Cancelled' }
    default: return { icon: Circle, variant: 'outline' as const, classes: 'bg-zinc-500/15 text-zinc-500 border-zinc-500/20', label: props.status }
  }
})
</script>

<template>
  <Badge
    :variant="config.variant"
    :class="cn('inline-flex items-center gap-1 text-[11px] font-medium', config.classes)"
  >
    <component :is="config.icon" :size="13" :stroke-width="2" />
    <span class="leading-none">{{ config.label }}</span>
  </Badge>
</template>
