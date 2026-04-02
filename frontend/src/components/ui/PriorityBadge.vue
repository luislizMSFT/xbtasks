<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

const props = defineProps<{ priority: string }>()

const config = computed(() => {
  switch (props.priority) {
    case 'P0': return { bg: 'bg-red-500/15', text: 'text-red-600 dark:text-red-400' }
    case 'P1': return { bg: 'bg-orange-500/15', text: 'text-orange-600 dark:text-orange-400' }
    case 'P2': return { bg: 'bg-amber-500/15', text: 'text-amber-600 dark:text-amber-400' }
    case 'P3': return { bg: 'bg-zinc-500/15', text: 'text-zinc-500 dark:text-zinc-400' }
    default: return { bg: 'bg-zinc-500/15', text: 'text-zinc-500' }
  }
})
</script>

<template>
  <Badge
    variant="outline"
    :class="cn('border-transparent text-[11px] font-semibold uppercase tracking-wide', config.bg, config.text)"
  >
    {{ priority }}
  </Badge>
</template>
