<script setup lang="ts">
import { computed } from 'vue'
import { Link } from 'lucide-vue-next'

const props = defineProps<{ adoId: string }>()

const display = computed(() => {
  if (!props.adoId) return ''
  return props.adoId.startsWith('ADO-') ? props.adoId : `ADO-${props.adoId}`
})
</script>

<template>
  <span
    v-if="adoId"
    class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] font-medium leading-none bg-blue-500/10 text-blue-600 dark:text-blue-400 cursor-pointer hover:bg-blue-500/20 transition-colors"
  >
    <Link :size="10" :stroke-width="2.5" />
    {{ display }}
  </span>
</template>
