<script setup lang="ts">
import { ref, watch, computed, nextTick } from 'vue'
import { useTaskStore, type Task } from '../stores/tasks'
import StatusBadge from './ui/StatusBadge.vue'
import PriorityBadge from './ui/PriorityBadge.vue'
import TagChip from './ui/TagChip.vue'
import AdoBadge from './ui/AdoBadge.vue'
import { X, Trash2, Plus, Link } from 'lucide-vue-next'

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
      class="w-[40%] min-w-[340px] max-w-[520px] border-l border-border-default bg-surface-primary flex flex-col overflow-y-auto"
    >
      <!-- Header -->
      <div class="flex items-center justify-between px-5 py-3 border-b border-border-default">
        <StatusBadge :status="editStatus" />
        <button
          @click="emit('close')"
          class="w-7 h-7 flex items-center justify-center rounded-md text-text-secondary hover:text-text-primary hover:bg-surface-tertiary transition-colors"
        >
          <X :size="16" />
        </button>
      </div>

      <!-- Body -->
      <div class="flex-1 px-5 py-4 space-y-5">
        <!-- Title -->
        <input
          ref="titleInput"
          v-model="editTitle"
          @blur="save"
          @keydown.enter="($event.target as HTMLInputElement)?.blur()"
          class="w-full text-lg font-semibold bg-transparent border-none outline-none text-text-primary placeholder-text-secondary"
          placeholder="Task title…"
        />

        <!-- Status & Priority row -->
        <div class="flex items-center gap-3">
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-text-secondary uppercase tracking-wider">Status</label>
            <select
              v-model="editStatus"
              @change="onStatusChange(editStatus)"
              class="block w-full px-2.5 py-1.5 rounded-md text-sm bg-surface-tertiary border border-border-default text-text-primary outline-none focus:ring-1 focus:ring-accent"
            >
              <option v-for="s in statuses" :key="s" :value="s">
                {{ s.replace('_', ' ') }}
              </option>
            </select>
          </div>
          <div class="space-y-1.5">
            <label class="text-[11px] font-medium text-text-secondary uppercase tracking-wider">Priority</label>
            <select
              v-model="editPriority"
              @change="onPriorityChange(editPriority)"
              class="block w-full px-2.5 py-1.5 rounded-md text-sm bg-surface-tertiary border border-border-default text-text-primary outline-none focus:ring-1 focus:ring-accent"
            >
              <option v-for="p in priorities" :key="p" :value="p">{{ p }}</option>
            </select>
          </div>
        </div>

        <!-- Description -->
        <div class="space-y-1.5">
          <label class="text-[11px] font-medium text-text-secondary uppercase tracking-wider">Description</label>
          <textarea
            v-model="editDescription"
            @blur="save"
            rows="4"
            class="w-full px-3 py-2 rounded-md text-sm bg-surface-tertiary border border-border-default text-text-primary placeholder-text-secondary outline-none focus:ring-1 focus:ring-accent resize-none"
            placeholder="Add a description…"
          />
        </div>

        <!-- Tags -->
        <div class="space-y-1.5">
          <label class="text-[11px] font-medium text-text-secondary uppercase tracking-wider">Tags</label>
          <div class="flex flex-wrap gap-1.5">
            <TagChip
              v-for="tag in editTags"
              :key="tag"
              :tag="tag"
              :removable="true"
              @remove="removeTag(tag)"
            />
            <input
              v-model="newTag"
              @keydown.enter.prevent="addTag"
              class="inline-flex px-1.5 py-0.5 text-[12px] bg-transparent border-none outline-none text-text-primary placeholder-text-secondary w-20"
              placeholder="+ tag"
            />
          </div>
        </div>

        <!-- ADO Link -->
        <div v-if="task.adoId" class="space-y-1.5">
          <label class="text-[11px] font-medium text-text-secondary uppercase tracking-wider">ADO Link</label>
          <div class="flex items-center gap-2">
            <AdoBadge :ado-id="task.adoId" />
            <span class="text-xs text-text-secondary">Linked work item</span>
          </div>
        </div>

        <!-- Blocked by -->
        <div v-if="task.blockedBy || task.status === 'blocked'" class="space-y-1.5">
          <label class="text-[11px] font-medium text-text-secondary uppercase tracking-wider">Blocked By</label>
          <div class="px-3 py-2 rounded-md text-sm bg-red-500/5 border border-red-500/20 text-red-600 dark:text-red-400">
            {{ task.blockedReason || `Task #${task.blockedBy}` }}
          </div>
        </div>

        <!-- Subtasks placeholder -->
        <div class="space-y-1.5">
          <label class="text-[11px] font-medium text-text-secondary uppercase tracking-wider">Subtasks</label>
          <button class="flex items-center gap-1.5 text-xs text-text-secondary hover:text-accent transition-colors">
            <Plus :size="14" />
            Add subtask
          </button>
        </div>
      </div>

      <!-- Footer -->
      <div class="px-5 py-3 border-t border-border-default space-y-2">
        <div class="flex justify-between text-[11px] text-text-secondary">
          <span>Created {{ timeCreated }}</span>
          <span>Updated {{ timeUpdated }}</span>
        </div>
        <button
          @click="deleteTask"
          class="flex items-center gap-1.5 text-xs text-red-500 hover:text-red-600 transition-colors"
        >
          <Trash2 :size="13" />
          Delete task
        </button>
      </div>
    </aside>
  </Transition>
</template>
