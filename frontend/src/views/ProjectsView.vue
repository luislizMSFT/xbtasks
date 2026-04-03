<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useProjectStore } from '@/stores/projects'
import { useTaskStore, type Task } from '@/stores/tasks'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import PageHeader from '@/components/PageHeader.vue'
import {
  Plus, Folder, ChevronDown, ChevronRight,
  CheckCircle2, Pencil,
} from 'lucide-vue-next'

const projectStore = useProjectStore()
const taskStore = useTaskStore()

const selectedProjectId = ref<number | null>(null)
const showCreate = ref(false)
const newName = ref('')
const newDescription = ref('')
const editingName = ref(false)
const editingDescription = ref(false)
const editName = ref('')
const editDescription = ref('')
const collapsedSections = ref<Set<string>>(new Set(['done']))

onMounted(async () => {
  await Promise.all([
    projectStore.fetchProjects(),
    taskStore.fetchTasks(),
  ])
  if (projectStore.projects.length > 0 && !selectedProjectId.value) {
    selectedProjectId.value = projectStore.projects[0].id
  }
})

// --- Selected project ---
const selectedProject = computed(() =>
  projectStore.projects.find(p => p.id === selectedProjectId.value) ?? null
)

// --- Tasks for selected project ---
const projectTasks = computed(() => {
  if (!selectedProjectId.value) return []
  return taskStore.tasks.filter(t => t.projectId === selectedProjectId.value)
})

const tasksByStatus = computed(() => {
  const groups: Record<string, Task[]> = {
    blocked: [], in_progress: [], in_review: [], todo: [], done: [], cancelled: [],
  }
  for (const t of projectTasks.value) {
    if (groups[t.status]) groups[t.status].push(t)
  }
  return groups
})

const projectProgress = computed(() => {
  const total = projectTasks.value.length
  if (total === 0) return 0
  const done = projectTasks.value.filter(t => t.status === 'done').length
  return Math.round((done / total) * 100)
})

// Per-project helpers for left panel
function projectTaskCount(projectId: number) {
  return taskStore.tasks.filter(t => t.projectId === projectId).length
}

function projectProgressPct(projectId: number) {
  const total = projectTaskCount(projectId)
  if (total === 0) return 0
  const done = taskStore.tasks.filter(t => t.projectId === projectId && t.status === 'done').length
  return Math.round((done / total) * 100)
}

// --- Section display ---
const sectionOrder = ['blocked', 'in_progress', 'in_review', 'todo', 'done'] as const

const sectionMeta: Record<string, { label: string; dot: string }> = {
  blocked:     { label: 'Blocked',     dot: 'bg-red-500' },
  in_progress: { label: 'In Progress', dot: 'bg-blue-500' },
  in_review:   { label: 'In Review',   dot: 'bg-violet-500' },
  todo:        { label: 'To Do',       dot: 'bg-zinc-400' },
  done:        { label: 'Done',        dot: 'bg-emerald-500' },
}

const visibleSections = computed(() =>
  sectionOrder.filter(s => (tasksByStatus.value[s]?.length ?? 0) > 0)
)

function toggleSection(section: string) {
  if (collapsedSections.value.has(section)) {
    collapsedSections.value.delete(section)
  } else {
    collapsedSections.value.add(section)
  }
}

// --- Project status ---
const statusColors: Record<string, string> = {
  active: 'bg-emerald-500',
  paused: 'bg-amber-500',
  completed: 'bg-blue-500',
  archived: 'bg-zinc-400',
  blocked: 'bg-red-500',
}

const projectStatuses = ['active', 'paused', 'completed', 'archived', 'blocked']

// --- Create project ---
async function createProject() {
  const name = newName.value.trim()
  if (!name) return
  const p = await projectStore.createProject(name, newDescription.value.trim())
  newName.value = ''
  newDescription.value = ''
  showCreate.value = false
  if (p) selectedProjectId.value = p.id
}

// --- Edit project (mock-only: mutate in place) ---
function startEditName() {
  if (!selectedProject.value) return
  editName.value = selectedProject.value.name
  editingName.value = true
}

function saveName() {
  if (!selectedProject.value) return
  const name = editName.value.trim()
  if (name && name !== selectedProject.value.name) {
    selectedProject.value.name = name
    selectedProject.value.updatedAt = new Date().toISOString()
  }
  editingName.value = false
}

function startEditDescription() {
  if (!selectedProject.value) return
  editDescription.value = selectedProject.value.description
  editingDescription.value = true
}

function saveDescription() {
  if (!selectedProject.value) return
  selectedProject.value.description = editDescription.value.trim()
  selectedProject.value.updatedAt = new Date().toISOString()
  editingDescription.value = false
}

function setProjectStatus(value: string | number | bigint | Record<string, any> | null) {
  if (!selectedProject.value || typeof value !== 'string') return
  selectedProject.value.status = value
  selectedProject.value.updatedAt = new Date().toISOString()
}

// --- Stats badges ---
const statBadges = computed(() => {
  if (!selectedProjectId.value) return []
  return [
    { label: 'Blocked', count: tasksByStatus.value.blocked?.length ?? 0, classes: 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20' },
    { label: 'In Progress', count: tasksByStatus.value.in_progress?.length ?? 0, classes: 'bg-blue-500/15 text-blue-600 dark:text-blue-400 border-blue-500/20' },
    { label: 'In Review', count: tasksByStatus.value.in_review?.length ?? 0, classes: 'bg-violet-500/15 text-violet-600 dark:text-violet-400 border-violet-500/20' },
    { label: 'To Do', count: tasksByStatus.value.todo?.length ?? 0, classes: 'bg-zinc-500/15 text-zinc-500 border-zinc-500/20' },
    { label: 'Done', count: tasksByStatus.value.done?.length ?? 0, classes: 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' },
  ].filter(b => b.count > 0)
})

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'now'
  if (mins < 60) return `${mins}m`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h`
  const days = Math.floor(hours / 24)
  return `${days}d`
}
</script>

<template>
  <div class="flex-1 flex overflow-hidden">
    <!-- Left panel: Project list -->
    <div class="w-[35%] flex flex-col min-w-0 border-r border-border">
      <PageHeader>
        <template #left>
          <span class="text-xs text-muted-foreground">{{ projectStore.projects.length }} projects</span>
        </template>
        <template #right>
          <Button size="sm" class="h-7 text-xs gap-1" @click="showCreate = !showCreate">
            <Plus :size="13" />
            New Project
          </Button>
        </template>
      </PageHeader>

      <ScrollArea class="flex-1 h-full">
        <!-- Inline create form -->
        <Transition
          enter-active-class="transition duration-150 ease-out"
          enter-from-class="opacity-0 -translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-100 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="showCreate" class="p-3 border-b border-border space-y-2">
            <Input
              v-model="newName"
              @keydown.enter="createProject"
              placeholder="Project name"
              autofocus
              class="text-sm"
            />
            <Textarea
              v-model="newDescription"
              :rows="2"
              class="resize-none text-sm"
              placeholder="Description (optional)"
            />
            <div class="flex justify-end gap-2">
              <Button variant="ghost" size="sm" class="h-7 text-xs" @click="showCreate = false">
                Cancel
              </Button>
              <Button size="sm" class="h-7 text-xs" @click="createProject">
                Create
              </Button>
            </div>
          </div>
        </Transition>

        <!-- Project list -->
        <div v-if="projectStore.projects.length > 0">
          <div
            v-for="project in projectStore.projects"
            :key="project.id"
            @click="selectedProjectId = project.id"
            :class="cn(
              'flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors border-b border-border/50',
              'hover:bg-muted/40',
              selectedProjectId === project.id ? 'bg-primary/[0.06]' : ''
            )"
          >
            <span
              class="w-2 h-2 rounded-full shrink-0"
              :class="statusColors[project.status] || 'bg-zinc-400'"
            />
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-foreground truncate">{{ project.name }}</span>
                <span class="text-[10px] text-muted-foreground/50 tabular-nums shrink-0">
                  {{ projectTaskCount(project.id) }}
                </span>
              </div>
              <div v-if="projectTaskCount(project.id) > 0" class="flex items-center gap-2 mt-1.5">
                <div class="flex-1 h-[3px] rounded-full bg-muted overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all duration-300"
                    :class="projectProgressPct(project.id) === 100 ? 'bg-emerald-500' : 'bg-blue-500'"
                    :style="{ width: projectProgressPct(project.id) + '%' }"
                  />
                </div>
                <span class="text-[10px] text-muted-foreground/40 tabular-nums w-7 text-right">
                  {{ projectProgressPct(project.id) }}%
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else-if="!projectStore.loading" class="flex flex-col items-center justify-center py-16 gap-2">
          <Folder :size="24" class="text-muted-foreground/40" />
          <p class="text-sm font-medium text-foreground">No projects yet</p>
          <p class="text-xs text-muted-foreground">Create a project to get started</p>
          <Button size="sm" class="mt-2" @click="showCreate = true">Create Project</Button>
        </div>

        <!-- Loading -->
        <div v-if="projectStore.loading" class="flex items-center justify-center py-20">
          <div class="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
        </div>
      </ScrollArea>
    </div>

    <!-- Right panel: Project detail -->
    <div v-if="selectedProject" class="flex-1 flex flex-col min-w-0">
      <ScrollArea class="flex-1 h-full">
        <div class="px-6 py-5 space-y-5">
          <!-- Project header -->
          <div class="space-y-3">
            <!-- Editable name -->
            <div class="flex items-center gap-2 group">
              <div v-if="editingName" class="flex-1">
                <Input
                  v-model="editName"
                  @keydown.enter="saveName"
                  @keydown.esc="editingName = false"
                  @blur="saveName"
                  autofocus
                  class="text-lg font-semibold"
                />
              </div>
              <h1
                v-else
                @click="startEditName"
                class="text-lg font-semibold text-foreground cursor-pointer hover:text-primary transition-colors"
              >
                {{ selectedProject.name }}
              </h1>
              <button
                v-if="!editingName"
                @click="startEditName"
                class="opacity-0 group-hover:opacity-100 text-muted-foreground hover:text-foreground transition-opacity"
              >
                <Pencil :size="13" />
              </button>
            </div>

            <!-- Editable description -->
            <div class="group">
              <Textarea
                v-if="editingDescription"
                v-model="editDescription"
                @keydown.esc="editingDescription = false"
                @blur="saveDescription"
                :rows="2"
                autofocus
                class="resize-none text-sm"
              />
              <p
                v-else
                @click="startEditDescription"
                class="text-sm text-muted-foreground cursor-pointer hover:text-foreground transition-colors"
              >
                {{ selectedProject.description || 'Add a description...' }}
              </p>
            </div>

            <!-- Status selector -->
            <div class="flex items-center gap-3">
              <span class="text-xs text-muted-foreground">Status:</span>
              <Select :model-value="selectedProject.status" @update:model-value="setProjectStatus">
                <SelectTrigger class="w-32 h-7 text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="s in projectStatuses" :key="s" :value="s" class="text-xs capitalize">
                    {{ s }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <Separator />

          <!-- Stats bar -->
          <div class="flex items-center gap-4 flex-wrap">
            <span class="text-xs text-muted-foreground">{{ projectTasks.length }} tasks</span>
            <div class="flex items-center gap-1.5 flex-wrap">
              <Badge
                v-for="badge in statBadges"
                :key="badge.label"
                variant="outline"
                :class="cn('text-[11px] font-medium', badge.classes)"
              >
                {{ badge.count }} {{ badge.label }}
              </Badge>
            </div>
          </div>

          <!-- Progress bar -->
          <div v-if="projectTasks.length > 0" class="space-y-1.5">
            <div class="flex items-center justify-between">
              <span class="text-xs text-muted-foreground">Progress</span>
              <span
                class="text-xs font-medium tabular-nums"
                :class="projectProgress === 100 ? 'text-emerald-500' : 'text-foreground'"
              >
                {{ projectProgress }}%
              </span>
            </div>
            <div class="w-full h-1.5 rounded-full bg-muted overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500"
                :class="projectProgress === 100 ? 'bg-emerald-500' : 'bg-blue-500'"
                :style="{ width: projectProgress + '%' }"
              />
            </div>
          </div>

          <Separator />

          <!-- Task list grouped by status -->
          <div v-if="projectTasks.length > 0" class="space-y-1">
            <div v-for="section in visibleSections" :key="section">
              <!-- Section header -->
              <button
                @click="toggleSection(section)"
                class="flex items-center gap-2 w-full px-1 py-2 hover:bg-muted/30 rounded transition-colors"
              >
                <component
                  :is="collapsedSections.has(section) ? ChevronRight : ChevronDown"
                  :size="14"
                  class="text-muted-foreground/50"
                />
                <span :class="cn('w-1.5 h-1.5 rounded-full', sectionMeta[section].dot)" />
                <span class="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground/60">
                  {{ sectionMeta[section].label }}
                </span>
                <span class="text-[11px] text-muted-foreground/40 tabular-nums">
                  {{ tasksByStatus[section]?.length ?? 0 }}
                </span>
                <div class="flex-1 h-px bg-border/50" />
              </button>

              <!-- Task rows -->
              <div v-if="!collapsedSections.has(section)">
                <div
                  v-for="task in tasksByStatus[section]"
                  :key="task.id"
                  class="flex items-center gap-2.5 px-3 py-2 rounded hover:bg-muted/40 transition-colors group"
                >
                  <!-- Done checkbox -->
                  <button
                    @click="taskStore.setStatus(task.id, task.status === 'done' ? 'todo' : 'done')"
                    :class="cn(
                      'size-[16px] rounded-full border-[1.5px] shrink-0 flex items-center justify-center transition-all hover:scale-110',
                      task.status === 'done'
                        ? 'bg-emerald-500 border-emerald-500'
                        : task.status === 'blocked'
                          ? 'border-red-400 hover:border-red-500'
                          : 'border-muted-foreground/30 hover:border-muted-foreground/60'
                    )"
                  >
                    <CheckCircle2 v-if="task.status === 'done'" :size="9" class="text-white" :stroke-width="3" />
                  </button>

                  <!-- Title -->
                  <span
                    :class="cn(
                      'text-[13px] font-medium truncate flex-1',
                      task.status === 'done' ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground'
                    )"
                  >
                    {{ task.title }}
                  </span>

                  <!-- Priority badge -->
                  <PriorityBadge v-if="task.priority" :priority="task.priority" />

                  <!-- Time -->
                  <span class="text-[10px] text-muted-foreground/40 tabular-nums shrink-0 w-6 text-right">
                    {{ timeAgo(task.updatedAt) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- No tasks state -->
          <div v-else class="text-center py-10">
            <p class="text-sm text-muted-foreground">No tasks in this project</p>
            <p class="text-xs text-muted-foreground/60 mt-1">Assign tasks to this project to see them here</p>
          </div>
        </div>
      </ScrollArea>
    </div>

    <!-- Empty state: no project selected -->
    <div v-else class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <Folder :size="32" class="text-muted-foreground/30 mx-auto mb-2" />
        <p class="text-sm text-muted-foreground">Select a project to view details</p>
      </div>
    </div>
  </div>
</template>
