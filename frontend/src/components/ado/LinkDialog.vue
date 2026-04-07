<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useTaskStore } from '@/stores/tasks'
import type { Task } from '@/stores/tasks'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Search, Link, Loader2, AlertCircle,
} from 'lucide-vue-next'
import { statusClasses, priorityClasses } from '@/lib/styles'
import { useNotify } from '@/composables/useNotify'

const props = defineProps<{
  open: boolean
  adoId: string
  adoTitle: string
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  'linked': [adoId: string]
}>()

const taskStore = useTaskStore()
const notify = useNotify()

const searchText = ref('')
const selectedTask = ref<Task | null>(null)
const linking = ref(false)
const linkError = ref('')

const personalTasks = computed(() =>
  taskStore.tasks.filter(t => !taskStore.isPublic(t.id))
)

const filteredTasks = computed(() => {
  if (!searchText.value) return personalTasks.value.slice(0, 20)
  const q = searchText.value.toLowerCase()
  return personalTasks.value
    .filter(t => t.title.toLowerCase().includes(q) || String(t.id).includes(q))
    .slice(0, 20)
})

function selectTask(task: Task) {
  selectedTask.value = task
}

async function confirmLink() {
  if (!selectedTask.value) return
  linking.value = true
  linkError.value = ''
  try {
    const { linkTask } = await import('@/api/links')
    await linkTask(selectedTask.value.id, props.adoId)
    await taskStore.fetchTasks()
    emit('linked', props.adoId)
    emit('update:open', false)
    searchText.value = ''
    selectedTask.value = null
    notify.success('Task linked to ADO')
  } catch (e: any) {
    linkError.value = e?.message || 'Failed to link task'
    notify.error('Link failed')
  } finally {
    linking.value = false
  }
}

watch(() => props.open, (val) => {
  if (!val) {
    searchText.value = ''
    selectedTask.value = null
    linkError.value = ''
  }
})
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2 text-sm">
          <Link :size="16" />
          Link to Personal Task
        </DialogTitle>
        <DialogDescription class="text-xs">
          Connect ADO item "{{ adoTitle }}" to an existing local task.
        </DialogDescription>
      </DialogHeader>

      <!-- Search input -->
      <div class="relative">
        <Search :size="14" class="absolute left-2.5 top-2.5 text-muted-foreground" />
        <Input
          v-model="searchText"
          placeholder="Search by task ID or title..."
          class="h-9 pl-8 text-sm"
        />
      </div>

      <!-- Results list -->
      <ScrollArea class="max-h-[240px] -mx-1 px-1">
        <div v-if="filteredTasks.length > 0" class="space-y-0.5">
          <button
            v-for="task in filteredTasks"
            :key="task.id"
            class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left transition-colors hover:bg-accent/50"
            :class="selectedTask?.id === task.id && 'bg-accent ring-1 ring-primary/30'"
            @click="selectTask(task)"
          >
            <span class="text-[10px] text-muted-foreground/50 shrink-0 tabular-nums">#{{ task.id }}</span>
            <span class="text-sm text-foreground flex-1 truncate">{{ task.title }}</span>
            <Badge v-if="task.priority" variant="outline" :class="['text-[10px] h-4 px-1.5 shrink-0', priorityClasses(task.priority)]">
              {{ task.priority }}
            </Badge>
            <Badge variant="outline" :class="['text-[10px] h-4 px-1.5 shrink-0', statusClasses(task.status)]">
              {{ task.status }}
            </Badge>
          </button>
        </div>
        <p v-else-if="searchText" class="text-[11px] text-muted-foreground/40 py-4 text-center">
          No matching personal tasks found
        </p>
        <p v-else class="text-[11px] text-muted-foreground/40 py-4 text-center">
          No personal tasks available
        </p>
      </ScrollArea>

      <!-- Error -->
      <div v-if="linkError" class="flex items-center gap-2 text-xs text-destructive">
        <AlertCircle :size="12" />
        {{ linkError }}
      </div>

      <!-- Selected task preview -->
      <div v-if="selectedTask" class="rounded-md border border-primary/20 bg-primary/5 px-3 py-2">
        <div class="flex items-center gap-2 text-xs">
          <span class="font-medium">#{{ selectedTask.id }}</span>
          <span class="truncate text-foreground">{{ selectedTask.title }}</span>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" size="sm" class="text-xs" @click="emit('update:open', false)">
          Cancel
        </Button>
        <Button
          size="sm"
          class="text-xs gap-1.5"
          :disabled="!selectedTask || linking"
          @click="confirmLink"
        >
          <Loader2 v-if="linking" :size="12" class="animate-spin" />
          <Link v-else :size="12" />
          Link
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
