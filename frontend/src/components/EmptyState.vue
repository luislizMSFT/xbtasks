<script setup lang="ts">
import type { Component } from 'vue'
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  icon?: Component
  title: string
  description?: string
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div :class="cn('flex flex-col items-center justify-center py-12 text-center', props.class)">
    <component
      v-if="icon"
      :is="icon"
      :size="48"
      class="mb-4 text-muted-foreground/50"
    />
    <h3 class="text-lg font-medium text-muted-foreground">{{ title }}</h3>
    <p v-if="description" class="mt-1 text-sm text-muted-foreground/70 max-w-sm">
      {{ description }}
    </p>
    <div v-if="$slots.action" class="mt-4">
      <slot name="action" />
    </div>
  </div>
</template>
