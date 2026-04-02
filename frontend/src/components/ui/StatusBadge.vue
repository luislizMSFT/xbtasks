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
    case 'todo': return { icon: Circle, color: 'text-zinc-400', bg: 'bg-zinc-400/10', label: 'To Do' }
    case 'in_progress': return { icon: CircleDot, color: 'text-blue-500', bg: 'bg-blue-500/10', label: 'In Progress' }
    case 'in_review': return { icon: Eye, color: 'text-violet-500', bg: 'bg-violet-500/10', label: 'In Review' }
    case 'done': return { icon: CheckCircle2, color: 'text-emerald-500', bg: 'bg-emerald-500/10', label: 'Done' }
    case 'blocked': return { icon: Octagon, color: 'text-red-500', bg: 'bg-red-500/10', label: 'Blocked' }
    case 'cancelled': return { icon: XCircle, color: 'text-zinc-400', bg: 'bg-zinc-400/10', label: 'Cancelled' }
    default: return { icon: Circle, color: 'text-zinc-400', bg: 'bg-zinc-400/10', label: props.status }
  }
})
</script>

<template>
  <Badge
    variant="outline"
    :class="cn('inline-flex items-center gap-1 border-transparent text-[11px] font-medium', config.bg, config.color)"
  >
    <component :is="config.icon" :size="13" :stroke-width="2" />
    <span class="leading-none">{{ config.label }}</span>
  </Badge>
</template>
