<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Link } from 'lucide-vue-next'

const props = defineProps<{ adoId: string }>()

const display = computed(() => {
  if (!props.adoId) return ''
  return props.adoId.startsWith('ADO-') ? props.adoId : `ADO-${props.adoId}`
})
</script>

<template>
  <Badge
    v-if="adoId"
    variant="outline"
    class="gap-1 text-[11px] font-medium bg-blue-500/15 text-blue-600 dark:text-blue-400 border-blue-500/20 cursor-pointer hover:bg-blue-500/25 transition-colors"
  >
    <Link :size="10" :stroke-width="2.5" />
    {{ display }}
  </Badge>
</template>
