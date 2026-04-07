<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

const props = defineProps<{ priority: string }>()

// Uses semantic --color-priority-* tokens from UI-SPEC (defined in style.css @theme)
const colorVar = computed(() => {
  switch (props.priority) {
    case 'P0': return '--color-priority-p0'
    case 'P1': return '--color-priority-p1'
    case 'P2': return '--color-priority-p2'
    case 'P3': return '--color-priority-p3'
    default: return '--color-priority-p3'
  }
})

const badgeStyle = computed(() => ({
  color: `var(${colorVar.value})`,
  backgroundColor: `color-mix(in srgb, var(${colorVar.value}) 15%, transparent)`,
  borderColor: `color-mix(in srgb, var(${colorVar.value}) 20%, transparent)`,
}))
</script>

<template>
  <Badge
    variant="outline"
    :class="cn('text-[11px] font-semibold uppercase tracking-wide border', '')"
    :style="badgeStyle"
  >
    {{ priority }}
  </Badge>
</template>
