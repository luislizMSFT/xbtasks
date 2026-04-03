<script setup lang="ts">
/**
 * PageHeader — consistent sub-bar below the top bar.
 * Renders left slot (tabs, filters) and right slot (actions).
 * Used by every page to avoid duplicate titles and inconsistent padding.
 */
defineProps<{
  bordered?: boolean
}>()
</script>

<template>
  <div
    :class="[
      'flex items-center justify-between px-4 py-2.5 shrink-0',
      bordered !== false ? 'border-b border-border' : ''
    ]"
  >
    <!-- Left: tabs, filters, labels -->
    <div class="flex items-center gap-2 min-w-0">
      <slot name="left" />
    </div>

    <!-- Right: actions -->
    <div class="flex items-center gap-1.5 shrink-0">
      <slot name="right" />
    </div>
  </div>
</template>
