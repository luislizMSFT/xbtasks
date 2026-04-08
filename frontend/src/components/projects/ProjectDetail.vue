<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import { useProjectStore } from '@/stores/projects'
import { useTaskStore, type Task } from '@/stores/tasks'
import { useADOStore } from '@/stores/ado'
import AdoItemPicker from '@/components/ado/AdoItemPicker.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { statusBgColor, adoTypeColor, adoTypeIcon } from '@/lib/styles'
import { cn } from '@/lib/utils'
import {
  X, Link2Off, Link2,
} from 'lucide-vue-next'

const props = defineProps<{ projectId: number }>()
const emit = defineEmits<{ close: []; 'select-task': [id: number] }>()

const projectStore = useProjectStore()
const taskStore = useTaskStore()
const adoStore = useADOStore()

// Editable fields
const editName = ref('')
const editDescription = ref('')

const project = computed(() =>
  projectStore.projects.find(p => p.id === props.projectId) ?? null
)

const link = computed(() => projectStore.projectLinks.get(props.projectId) ?? null)
const progress = computed(() => projectStore.projectProgress.get(props.projectId) ?? null)

const childTasks = computed(() =>
  taskStore.tasks.filter(t => t.projectId === props.projectId)
)

const progressPct = computed(() => {
  const done = childTasks.value.filter(t => t.status === 'done').length
  const total = childTasks.value.length
  return total > 0 ? Math.round((done / total) * 100) : 0
})

const progressDone = computed(() => childTasks.value.filter(t => t.status === 'done').length)

// ADO type for linked project
const adoType = computed(() => {
  if (!link.value) return ''
  const wi = adoStore.workItemTree.find(w => w.adoId === link.value?.adoId)
  return wi?.type || ''
})

function typeIcon(type: string) {
  return adoTypeIcon(type)
}

// Sync fields when project changes
watch(() => props.projectId, () => {
  const p = project.value
  if (p) {
    editName.value = p.name
    editDescription.value = p.description
  }
  projectStore.fetchProjectLink(props.projectId)
  projectStore.fetchProjectProgress(props.projectId)
}, { immediate: true })

// Save
let saving = false
async function save() {
  if (!project.value || saving) return
  if (editName.value === project.value.name && editDescription.value === project.value.description) return
  saving = true
  try {
    const m = await import('../../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
    const updated = await m.Update(project.value.id, editName.value, editDescription.value, project.value.status)
    const idx = projectStore.projects.findIndex(p => p.id === props.projectId)
    if (idx !== -1) {
      projectStore.projects[idx].name = editName.value
      projectStore.projects[idx].description = editDescription.value
    }
  } catch (e) {
    console.warn('[ProjectDetail] save failed:', e)
  } finally {
    saving = false
  }
}

// ADO link/unlink
const pickerOpen = ref(false)

async function onAdoSelected(adoId: string) {
  pickerOpen.value = false
  await projectStore.linkProjectToADO(props.projectId, adoId)
  projectStore.fetchProjectLink(props.projectId)
}

async function unlinkAdo() {
  if (!link.value) return
  await projectStore.unlinkProject(props.projectId, link.value.adoId)
}

function priorityLabel(p: string) {
  return p || 'None'
}

onMounted(() => {
  projectStore.fetchProjectLink(props.projectId)
  projectStore.fetchProjectProgress(props.projectId)
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
      v-if="project"
      class="w-[45%] shrink-0 border-l border-border bg-background flex flex-col h-full min-h-0 relative"
    >
      <!-- Header -->
      <div class="shrink-0 border-b border-border px-4 pt-3 pb-3">
        <div class="flex items-start gap-2">
          <Input
            v-model="editName"
            @blur="save"
            @keydown.enter="($event.target as HTMLInputElement)?.blur()"
            class="text-base font-semibold border-none shadow-none focus-visible:ring-0 bg-transparent px-0 h-auto flex-1"
            placeholder="Project name..."
          />
          <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0 -mt-0.5" @click="emit('close')">
            <X :size="16" />
          </Button>
        </div>

        <!-- ADO Link -->
        <div class="flex items-center gap-2 mt-2">
          <template v-if="link">
            <span class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-500 border border-blue-500/30 rounded-full px-2 py-0.5">
              <component
                v-if="adoType"
                :is="typeIcon(adoType)"
                :size="12"
                :class="adoTypeColor(adoType)"
              />
              <AzureDevOpsIcon v-else :size="12" />
              Linked to #{{ link.adoId }}
            </span>
            <Button variant="ghost" size="icon" class="h-6 w-6 text-muted-foreground hover:text-red-500" @click="unlinkAdo">
              <Link2Off :size="14" />
            </Button>
          </template>
          <Button v-else variant="outline" size="sm" class="h-7 text-xs gap-1.5" @click="pickerOpen = true">
            <Link2 :size="12" />
            Link to ADO
          </Button>
        </div>

        <!-- Progress bar -->
        <div v-if="childTasks.length > 0" class="mt-2.5">
          <div class="flex items-center justify-between text-[10px] text-muted-foreground mb-1">
            <span>{{ progressDone }}/{{ childTasks.length }} tasks done</span>
            <span>{{ progressPct }}%</span>
          </div>
          <div class="h-1 w-full rounded-full bg-muted">
            <div
              class="h-1 rounded-full bg-blue-500 transition-all duration-300"
              :style="{ width: progressPct + '%' }"
            />
          </div>
        </div>
      </div>

      <!-- Scrollable content -->
      <ScrollArea class="flex-1 h-0 min-h-0">
        <div class="px-4 py-3 space-y-4">
          <!-- Description -->
          <div>
            <h3 class="text-xs font-medium text-muted-foreground mb-1.5">Description</h3>
            <Textarea
              v-model="editDescription"
              @blur="save"
              placeholder="Add a description..."
              class="min-h-[80px] text-sm resize-none border-border/50"
              rows="3"
            />
          </div>

          <!-- Child Tasks -->
          <div>
            <h3 class="text-xs font-medium text-muted-foreground mb-1.5">
              Tasks ({{ childTasks.length }})
            </h3>
            <div v-if="childTasks.length === 0" class="text-xs text-muted-foreground/60 py-2">
              No tasks in this project yet.
            </div>
            <div v-else class="space-y-0.5">
              <button
                v-for="task in childTasks"
                :key="task.id"
                class="w-full flex items-center gap-2 px-2 py-1.5 rounded-md text-left hover:bg-muted/50 transition-colors group"
                :class="taskStore.selectedTaskId === task.id && 'bg-primary/5'"
                @click="emit('select-task', task.id)"
              >
                <span :class="cn('size-2 rounded-full shrink-0', statusBgColor(task.status))" />
                <span class="text-sm truncate flex-1">{{ task.title }}</span>
                <span v-if="task.priority" class="text-[10px] text-muted-foreground tabular-nums shrink-0">
                  {{ task.priority }}
                </span>
              </button>
            </div>
          </div>
        </div>
      </ScrollArea>

      <!-- ADO Item Picker dialog -->
      <AdoItemPicker
        :open="pickerOpen"
        title="Link Project to ADO Work Item"
        @update:open="pickerOpen = $event"
        @selected="(adoId) => onAdoSelected(adoId)"
      />
    </aside>
  </Transition>
</template>
