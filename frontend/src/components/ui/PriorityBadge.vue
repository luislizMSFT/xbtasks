<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

const props = defineProps<{ priority: string }>()

const config = computed(() => {
  switch (props.priority) {
    case 'P0': return { variant: 'outline' as const, classes: 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20' }
    case 'P1': return { variant: 'outline' as const, classes: 'bg-orange-500/15 text-orange-600 dark:text-orange-400 border-orange-500/20' }
    case 'P2': return { variant: 'outline' as const, classes: 'bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/20' }
    case 'P3': return { variant: 'outline' as const, classes: 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20' }
    default: return { variant: 'outline' as const, classes: 'bg-zinc-500/15 text-muted-foreground border-zinc-500/20' }
  }
})
</script>

<template>
  <Badge
    :variant="config.variant"
    :class="cn('text-[11px] font-semibold uppercase tracking-wide', config.classes)"
  >
    {{ priority }}
  </Badge>
</template>
