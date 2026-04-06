<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg'
  label?: string
  class?: HTMLAttributes['class']
}>(), {
  size: 'md',
})

const sizeMap = {
  sm: 16,
  md: 24,
  lg: 32,
} as const
</script>

<template>
  <div :class="cn('flex items-center justify-center gap-2', props.class)">
    <Loader2 :size="sizeMap[size]" class="animate-spin text-muted-foreground" />
    <span v-if="label" class="text-sm text-muted-foreground">{{ label }}</span>
  </div>
</template>
