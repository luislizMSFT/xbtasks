<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { Circle } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

const emit = defineEmits<{
  add: [title: string]
  cancel: []
}>()

const title = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

function handleAdd() {
  const trimmed = title.value.trim()
  if (trimmed) {
    emit('add', trimmed)
    title.value = ''
    // Re-focus for rapid entry
    nextTick(() => inputRef.value?.focus())
  }
}

onMounted(() => {
  nextTick(() => inputRef.value?.focus())
})
</script>

<template>
  <div class="flex items-center gap-2 px-3 py-2 border-t border-border/30">
    <Circle :size="16" class="text-zinc-400 shrink-0" />
    <input
      ref="inputRef"
      v-model="title"
      placeholder="New task title… (Enter to add, Esc to cancel)"
      class="flex-1 text-sm bg-transparent border-none outline-none placeholder:text-muted-foreground/40"
      @keydown.enter="handleAdd"
      @keydown.escape="emit('cancel')"
    />
    <Button variant="outline" size="sm" class="h-7 px-2 text-[10px] shrink-0" @click="handleAdd" :disabled="!title.trim()">Add</Button>
    <Button variant="ghost" size="sm" class="h-7 px-1.5 text-[10px]" @click="emit('cancel')">Cancel</Button>
  </div>
</template>
