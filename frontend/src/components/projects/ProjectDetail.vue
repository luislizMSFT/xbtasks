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
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import { statusBgColor, adoTypeColor, adoTypeIcon, priorityDotBgColor, statusIcon, adoStateClasses } from '@/lib/styles'
import { cn } from '@/lib/utils'
import { useAdoMeta } from '@/composables/useAdoMeta'
import DOMPurify from 'dompurify'
import {
  X, Link2Off, Link2, ExternalLink, Folder, Pencil, Save, CheckCircle2, Archive, Trash2,
} from 'lucide-vue-next'

const props = defineProps<{ projectId: number }>()
const emit = defineEmits<{ close: []; 'select-task': [id: number] }>()

const projectStore = useProjectStore()
const taskStore = useTaskStore()
const adoStore = useADOStore()
const adoMeta = useAdoMeta()

const editName = ref('')
const editDescription = ref('')
const descriptionOpen = ref(false)

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

// Status breakdown for child tasks
const statusCounts = computed(() => {
  const counts: Record<string, number> = {}
  for (const t of childTasks.value) {
    counts[t.status] = (counts[t.status] || 0) + 1
  }
  return counts
})

// ADO type for linked project
const adoType = computed(() => {
  if (!link.value) return ''
  const wi = adoStore.workItemTree.find(w => w.adoId === link.value?.adoId)
  return wi?.type || ''
})

// ADO number for linked project
const adoNumber = computed(() => {
  if (!link.value) return ''
  return link.value.adoId
})

function typeIcon(type: string) {
  return adoTypeIcon(type)
}

// Get ADO meta for a child task
function taskMeta(task: Task) {
  return adoMeta.getAdoMeta(task.id)
}

// Sanitize description HTML
function sanitizeHtml(html: string): string {
  if (!html) return ''
  return DOMPurify.sanitize(html, {
    FORBID_TAGS: ['style', 'script', 'iframe', 'object', 'embed', 'link'],
    FORBID_ATTR: ['style'],
  })
}

// ADO URL
const adoUrl = ref('')
watch(() => link.value, async (lnk) => {
  if (!lnk) { adoUrl.value = ''; return }
  try {
    const { getCachedWorkItem } = await import('@/api/workitems')
    const wi = await getCachedWorkItem(lnk.adoId) as any
    if (wi?.org && wi?.project) {
      adoUrl.value = `https://${wi.org}.visualstudio.com/${wi.project}/_workitems/edit/${lnk.adoId}`
    } else {
      adoUrl.value = ''
    }
  } catch {
    adoUrl.value = ''
  }
}, { immediate: true })

async function openAdoLink() {
  if (!adoUrl.value) return
  try {
    const { openURL } = await import('@/api/browser')
    await openURL(adoUrl.value)
  } catch { window.open(adoUrl.value, '_blank') }
}

const statusLabels: Record<string, string> = {
  in_progress: 'In Progress', in_review: 'In Review',
  todo: 'To Do', blocked: 'Blocked', done: 'Done', cancelled: 'Cancelled',
}

// Sync fields when project changes
watch(() => props.projectId, () => {
  const p = project.value
  if (p) {
    editName.value = p.name
    editDescription.value = p.description
  }
  descriptionOpen.value = false
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
    await m.Update(project.value.id, editName.value, editDescription.value, project.value.status)
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

async function updateStatus(status: string) {
  if (!project.value) return
  await projectStore.updateProjectStatus(project.value.id, status)
}

// ADO link/unlink
const pickerOpen = ref(false)
const confirmDeleteOpen = ref(false)
const confirmUnlinkOpen = ref(false)

async function onAdoSelected(adoId: string) {
  pickerOpen.value = false
  await projectStore.linkProjectToADO(props.projectId, adoId)
  projectStore.fetchProjectLink(props.projectId)
}

async function unlinkAdo() {
  if (!link.value) return
  await projectStore.unlinkProject(props.projectId, link.value.adoId)
  confirmUnlinkOpen.value = false
}

async function onDeleteProject() {
  await projectStore.deleteProject(props.projectId)
  confirmDeleteOpen.value = false
  emit('close')
}

// Timestamps
const timeCreated = computed(() => {
  if (!project.value?.createdAt) return ''
  return new Date(project.value.createdAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
})

const timeUpdated = computed(() => {
  if (!project.value?.updatedAt) return ''
  return new Date(project.value.updatedAt).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
})

onMounted(() => {
  projectStore.fetchProjectLink(props.projectId)
  projectStore.fetchProjectProgress(props.projectId)
})
</script>

<template>
    <aside
      v-if="project"
      class="flex flex-col h-full min-h-0 bg-background"
    >
      <!-- ─── Header ──── -->
      <div class="shrink-0 border-b border-border px-4 py-3 space-y-2">
        <!-- Row 1: Type icon + Name + close -->
        <div class="flex items-center gap-2">
          <component
            v-if="adoType"
            :is="typeIcon(adoType)"
            :size="16"
            :class="adoTypeColor(adoType)"
            class="shrink-0"
          />
          <Folder v-else :size="16" class="text-muted-foreground shrink-0" />
          <Input
            v-model="editName"
            @blur="save"
            @keydown.enter="($event.target as HTMLInputElement)?.blur()"
            class="text-sm font-semibold border-none shadow-none focus-visible:ring-0 bg-transparent px-0 h-auto flex-1"
            placeholder="Project name..."
          />
          <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="emit('close')">
            <X :size="16" />
          </Button>
        </div>

        <!-- Row 2: Status summary badges -->
        <div v-if="childTasks.length > 0" class="flex items-center gap-1.5 flex-wrap">
          <Badge v-for="(count, status) in statusCounts" :key="status" variant="outline" class="text-[10px] h-5 px-1.5 gap-1">
            <span :class="cn('size-1.5 rounded-full', statusBgColor(String(status)))" />
            {{ statusLabels[String(status)] ?? status }} · {{ count }}
          </Badge>
        </div>

        <!-- Row 3: ADO integration bar -->
        <div v-if="link" class="flex items-center gap-2 px-2 py-1 rounded-md bg-blue-500/5 border border-blue-500/15">
          <AzureDevOpsIcon :size="12" class="text-blue-500 shrink-0" />
          <span class="text-[10px] text-blue-500 tabular-nums font-medium">#{{ adoNumber }}</span>
          <div class="flex-1" />
          <Button v-if="adoUrl" variant="ghost" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5 text-blue-500 hover:text-blue-600" @click="openAdoLink">
            <ExternalLink :size="10" /> Open
          </Button>
          <Button variant="ghost" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5 text-red-400 hover:text-red-500 hover:bg-red-500/10" @click="confirmUnlinkOpen = true">
            <Link2Off :size="10" />
          </Button>
        </div>
        <div v-else class="flex items-center gap-2 px-2 py-1 rounded-md border border-dashed border-border/50">
          <Folder :size="12" class="text-muted-foreground/40 shrink-0" />
          <span class="text-[10px] text-muted-foreground/60 flex-1">Not linked to ADO</span>
          <Button variant="outline" size="sm" class="h-5 text-[9px] gap-0.5 px-1.5" @click="pickerOpen = true">
            <Link2 :size="10" /> Link
          </Button>
        </div>

        <!-- Progress bar -->
        <div v-if="childTasks.length > 0">
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

        <!-- Status actions -->
        <div class="flex items-center gap-1.5">
          <Badge variant="outline" class="text-[10px] h-5 px-1.5 capitalize">{{ project.status }}</Badge>
          <div class="flex-1" />
          <template v-if="project.status === 'active'">
            <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1" @click="updateStatus('completed')">
              <CheckCircle2 :size="11" /> Complete
            </Button>
            <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1" @click="updateStatus('archived')">
              <Archive :size="11" /> Archive
            </Button>
          </template>
          <template v-else-if="project.status === 'completed'">
            <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1" @click="updateStatus('active')">Reopen</Button>
            <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1" @click="updateStatus('archived')">
              <Archive :size="11" /> Archive
            </Button>
          </template>
          <template v-else-if="project.status === 'archived'">
            <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1" @click="updateStatus('active')">Reopen</Button>
          </template>
          <template v-else>
            <Button variant="outline" size="sm" class="h-6 text-[10px] gap-1" @click="updateStatus('active')">Activate</Button>
          </template>
          <Button variant="ghost" size="sm" class="h-6 text-[10px] gap-1 text-red-400 hover:text-red-500 hover:bg-red-500/10" @click="confirmDeleteOpen = true">
            <Trash2 :size="11" />
          </Button>
        </div>
      </div>

      <!-- ─── Scrollable content ─── -->
      <ScrollArea class="flex-1 h-0 min-h-0">
        <div class="flex flex-col">
          <!-- Child Tasks section (above description per user request) -->
          <div class="border-b border-border px-4 py-2">
            <div class="flex items-center justify-between mb-1.5">
              <div class="flex items-center gap-2">
                <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Tasks</h3>
                <Badge v-if="childTasks.length > 0" variant="secondary" class="h-4 text-[10px] px-1.5">
                  {{ progressDone }}/{{ childTasks.length }}
                </Badge>
              </div>
            </div>
            <div v-if="childTasks.length === 0" class="text-[11px] text-muted-foreground/40 italic py-2 text-center">
              No tasks in this project yet.
            </div>
            <div v-else class="flex flex-col gap-px">
              <button
                v-for="task in childTasks"
                :key="task.id"
                class="flex items-center gap-2 py-1.5 px-2 rounded hover:bg-muted/50 transition-colors group text-left"
                :class="taskStore.selectedTaskId === task.id && 'bg-primary/5'"
                @click="emit('select-task', task.id)"
              >
                <!-- Icon: ADO type icon for ADO tasks, status icon for personal -->
                <template v-if="task.adoId && taskMeta(task)?.type">
                  <component
                    :is="adoTypeIcon(taskMeta(task)!.type)"
                    :size="14"
                    :class="adoTypeColor(taskMeta(task)!.type)"
                    class="shrink-0"
                  />
                </template>
                <template v-else>
                  <component
                    :is="statusIcon(task.status)"
                    :size="14"
                    :class="task.status === 'done' ? 'text-emerald-500' : 'text-muted-foreground'"
                    class="shrink-0"
                  />
                </template>
                <!-- Title -->
                <span :class="cn('text-xs truncate flex-1', task.status === 'done' && 'line-through text-muted-foreground')">
                  {{ task.title }}
                </span>
                <!-- ADO state badge -->
                <Badge v-if="task.adoId && taskMeta(task)?.state" variant="outline" :class="adoStateClasses(taskMeta(task)!.state)" class="text-[8px] h-4 px-1">
                  {{ taskMeta(task)!.state }}
                </Badge>
                <!-- Priority dot -->
                <span :class="cn('size-1.5 rounded-full shrink-0', priorityDotBgColor(task.priority))" />
                <!-- ADO ID -->
                <span v-if="task.adoId" class="text-[9px] text-blue-500/60 tabular-nums shrink-0">#{{ task.adoId }}</span>
              </button>
            </div>
          </div>

          <!-- Description section -->
          <div class="border-b border-border px-4 py-3">
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
              <button v-if="!descriptionOpen" class="text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1" @click="descriptionOpen = true">
                <Pencil :size="10" /> Edit
              </button>
              <div v-else class="flex items-center gap-1">
                <button class="text-[10px] text-emerald-600 hover:text-emerald-700 flex items-center gap-0.5" @click="descriptionOpen = false; save()">
                  <Save :size="10" /> Save
                </button>
                <button class="text-[10px] text-muted-foreground hover:text-foreground" @click="descriptionOpen = false">Cancel</button>
              </div>
            </div>
            <Textarea
              v-if="descriptionOpen"
              v-model="editDescription"
              @blur="save"
              :rows="4"
              class="resize-y text-xs"
              placeholder="Add a description..."
            />
            <div v-else-if="editDescription"
              class="text-xs text-foreground prose prose-sm max-w-none [&_*]:text-xs [&_*]:text-foreground cursor-pointer hover:bg-muted/30 rounded p-1 -m-1 transition-colors"
              v-html="sanitizeHtml(editDescription)"
              @click="descriptionOpen = true"
            />
            <p v-else class="text-[11px] text-muted-foreground/40 italic cursor-pointer hover:text-muted-foreground" @click="descriptionOpen = true">
              Click to add description...
            </p>
          </div>
        </div>
      </ScrollArea>

      <!-- ─── Footer: timestamps ──── -->
      <div v-if="timeCreated || timeUpdated" class="shrink-0 border-t border-border flex items-center px-4 py-1.5">
        <span class="text-[11px] text-muted-foreground">
          <template v-if="timeCreated">Created {{ timeCreated }}</template>
          <template v-if="timeCreated && timeUpdated"> · </template>
          <template v-if="timeUpdated">Updated {{ timeUpdated }}</template>
        </span>
      </div>

      <!-- ADO Item Picker dialog -->
      <AdoItemPicker
        :open="pickerOpen"
        title="Link Project to ADO Work Item"
        @update:open="pickerOpen = $event"
        @selected="(adoId) => onAdoSelected(adoId)"
      />

      <!-- Delete confirmation -->
      <Dialog v-model:open="confirmDeleteOpen">
        <DialogContent class="max-w-sm">
          <DialogHeader>
            <DialogTitle>Delete project?</DialogTitle>
            <DialogDescription>
              "{{ project?.name }}" and all its data will be permanently deleted. This cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" size="sm" @click="confirmDeleteOpen = false">Cancel</Button>
            <Button variant="destructive" size="sm" @click="onDeleteProject">Delete</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <!-- Unlink confirmation -->
      <Dialog v-model:open="confirmUnlinkOpen">
        <DialogContent class="max-w-sm">
          <DialogHeader>
            <DialogTitle>Unlink from ADO?</DialogTitle>
            <DialogDescription>
              This will disconnect the project from ADO work item #{{ link?.adoId }}. The local project will be kept.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" size="sm" @click="confirmUnlinkOpen = false">Cancel</Button>
            <Button variant="destructive" size="sm" @click="unlinkAdo">Unlink</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </aside>
</template>
