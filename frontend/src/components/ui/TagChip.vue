<script setup lang="ts">
defineProps<{ tag: string; removable?: boolean }>()
defineEmits<{ remove: [] }>()
</script>

<template>
  <span
    class="group inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[11px] font-medium leading-none bg-zinc-200/80 text-zinc-700 dark:bg-zinc-700/60 dark:text-zinc-300 max-w-[100px] truncate"
  >
    <span class="opacity-50">#</span>{{ tag }}
    <button
      v-if="removable"
      @click.stop="$emit('remove')"
      class="ml-0.5 opacity-0 group-hover:opacity-100 transition-opacity hover:text-red-500"
    >
      ×
    </button>
  </span>
</template>
