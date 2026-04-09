<script setup lang="ts">
import { computed } from 'vue'
import { ListChecks, Lock, User } from 'lucide-vue-next'

const props = defineProps<{
  modelValue: string  // 'all' | 'linked' | 'personal'
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const label = computed(() => {
  switch (props.modelValue) {
    case 'linked': return 'ADO'
    case 'personal': return 'Personal'
    default: return 'All'
  }
})

const icon = computed(() => {
  switch (props.modelValue) {
    case 'linked': return Lock
    case 'personal': return User
    default: return ListChecks
  }
})

const isActive = computed(() => props.modelValue !== 'all')

function cycle() {
  const order = ['all', 'linked', 'personal']
  const idx = order.indexOf(props.modelValue)
  const next = order[(idx + 1) % order.length]
  emit('update:modelValue', next)
}
</script>

<template>
  <button
    class="text-[10px] flex items-center gap-1 px-1.5 py-0.5 rounded transition-colors"
    :class="isActive ? 'bg-primary/10 text-primary font-medium' : 'text-muted-foreground hover:text-foreground'"
    :title="`Filter: ${label}`"
    @click="cycle"
  >
    <component :is="icon" :size="10" />
    {{ label }}
  </button>
</template>
