<script setup lang="ts">
import { Badge } from '@/components/ui/badge'

defineProps<{ tag: string; removable?: boolean }>()
defineEmits<{ remove: [] }>()
</script>

<template>
  <Badge
    variant="secondary"
    class="group gap-0.5 text-[11px] font-medium max-w-[100px] truncate"
  >
    <span class="opacity-50">#</span>{{ tag }}
    <button
      v-if="removable"
      @click.stop="$emit('remove')"
      class="ml-0.5 opacity-0 group-hover:opacity-100 transition-opacity hover:text-red-500"
    >
      ×
    </button>
  </Badge>
</template>
