<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useTaskStore, type Task } from '@/stores/tasks'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import AdoBadge from '@/components/ui/AdoBadge.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { X, Trash2, Plus } from 'lucide-vue-next'

const taskStore = useTaskStore()
const emit = defineEmits<{ close: [] }>()

const editTitle = ref('')
const editDescription = ref('')
const editStatus = ref('')
const editPriority = ref('')
const editTags = ref<string[]>([])
const newTag = ref('')
const titleInput = ref<HTMLInputElement | null>(null)

const statuses = ['todo', 'in_progress', 'in_review', 'done', 'blocked', 'cancelled']
const priorities = ['P0', 'P1', 'P2', 'P3']

const task = computed(() => taskStore.selectedTask)

watch(task, (t) => {
  if (t) {
    editTitle.value = t.title
    editDescription.value = t.description
    editStatus.value = t.status
    editPriority.value = t.priority
    editTags.value = t.tags ? t.tags.split(',').map(s => s.trim()).filter(Boolean) : []
  }
}, { immediate: true })

async function save() {
  if (!task.value) return
  const updated: Task = {
    ...task.value,
    title: editTitle.value,
    description: editDescription.value,
    status: editStatus.value,
    priority: editPriority.value,
    tags: editTags.value.join(','),
  }
  await taskStore.updateTask(updated)
}

function addTag() {
  const tag = newTag.value.trim()
  if (tag && !editTags.value.includes(tag)) {
    editTags.value.push(tag)
    newTag.value = ''
    save()
  }
}

function removeTag(tag: string) {
  editTags.value = editTags.value.filter(t => t !== tag)
  save()
}

async function onStatusChange(status: string) {
  editStatus.value = status
  await save()
}

async function onPriorityChange(priority: string) {
  editPriority.value = priority
  await save()
}

async function deleteTask() {
  if (!task.value) return
  await taskStore.deleteTask(task.value.id)
  emit('close')
}

const timeCreated = computed(() => {
  if (!task.value) return ''
  return new Date(task.value.createdAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric', hour: 'numeric', minute: '2-digit'
  })
})

const timeUpdated = computed(() => {
  if (!task.value) return ''
  return new Date(task.value.updatedAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric', hour: 'numeric', minute: '2-digit'
  })
})
</script>

<template>
  <Transition
    enter-active-class="transition-transform duration-200 ease-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-150 ease-in"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <aside
      v-if="task"
      class="w-[40%] min-w-[340px] max-w-[520px] border-l border-border bg-background flex flex-col"
    >
      <!-- Header -->
      <div class="flex items-center justify-between px-5 py-3 border-b border-border">
        <StatusBadge :status="editStatus" />
        <Button variant="ghost" size="icon" class="h-7 w-7" @click="emit('close')">
          <X :size="16" />
        </Button>
      </div>

      <!-- Body -->
      <ScrollArea class="flex-1 h-full">
      <div class="px-5 py-4 space-y-5">
        <!-- Title -->
        <Input
          ref="titleInput"
          v-model="editTitle"
          @blur="save"
          @keydown.enter="($event.target as HTMLInputElement)?.blur()"
          class="text-lg font-semibold border-none shadow-none focus-visible:ring-0 bg-transparent px-0"
          placeholder="Task title…"
        />

        <Separator />

        <!-- Status & Priority row -->
        <div class="flex items-center gap-3">
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Status</label>
            <Select :model-value="editStatus" @update:model-value="(v) => onStatusChange(String(v))">
              <SelectTrigger class="w-[140px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="s in statuses" :key="s" :value="s">
                  {{ s.replace('_', ' ') }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Priority</label>
            <Select :model-value="editPriority" @update:model-value="(v) => onPriorityChange(String(v))">
              <SelectTrigger class="w-[100px]">
                <SelectValue placeholder="Priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="p in priorities" :key="p" :value="p">{{ p }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <Separator />

        <!-- Description -->
        <div class="space-y-1.5">
          <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Description</label>
          <Textarea
            v-model="editDescription"
            @blur="save"
            :rows="4"
            class="resize-none"
            placeholder="Add a description…"
          />
        </div>

        <Separator />

        <!-- Tags -->
        <div class="space-y-1.5">
          <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Tags</label>
          <div class="flex flex-wrap gap-1.5">
            <Badge
              v-for="tag in editTags"
              :key="tag"
              variant="secondary"
              class="cursor-pointer gap-1"
              @click="removeTag(tag)"
            >
              {{ tag }}
              <X :size="12" class="opacity-60 hover:opacity-100" />
            </Badge>
            <Input
              v-model="newTag"
              @keydown.enter.prevent="addTag"
              class="inline-flex h-6 w-20 px-1.5 text-[12px] border-none shadow-none focus-visible:ring-0 bg-transparent"
              placeholder="+ tag"
            />
          </div>
        </div>

        <Separator />

        <!-- ADO Link -->
        <div v-if="task.adoId" class="space-y-1.5">
          <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">ADO Link</label>
          <div class="flex items-center gap-2">
            <AdoBadge :ado-id="task.adoId" />
            <span class="text-xs text-muted-foreground">Linked work item</span>
          </div>
        </div>

        <!-- Blocked by -->
        <div v-if="task.blockedBy || task.status === 'blocked'" class="space-y-1.5">
          <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Blocked By</label>
          <div class="px-3 py-2 rounded-md text-sm bg-red-500/5 border border-red-500/20 text-red-600 dark:text-red-400">
            {{ task.blockedReason || `Task #${task.blockedBy}` }}
          </div>
        </div>

        <!-- Subtasks placeholder -->
        <div class="space-y-1.5">
          <label class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">Subtasks</label>
          <Button variant="ghost" size="sm" class="gap-1.5 text-xs text-muted-foreground">
            <Plus :size="14" />
            Add subtask
          </Button>
        </div>
      </div>
      </ScrollArea>

      <!-- Footer -->
      <div class="px-5 py-3 border-t border-border space-y-2">
        <div class="flex justify-between text-[11px] text-muted-foreground">
          <span>Created {{ timeCreated }}</span>
          <span>Updated {{ timeUpdated }}</span>
        </div>
        <Button variant="ghost" size="sm" class="gap-1.5 text-xs text-destructive hover:text-destructive hover:bg-destructive/10" @click="deleteTask">
          <Trash2 :size="13" />
          Delete task
        </Button>
      </div>
    </aside>
  </Transition>
</template>
