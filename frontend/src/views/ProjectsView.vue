<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore, type ProjectADOLink } from '@/stores/projects'
import { useTaskStore, type Task } from '@/stores/tasks'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import PageHeader from '@/components/PageHeader.vue'
import ProjectCard from '@/components/projects/ProjectCard.vue'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/EmptyState.vue'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import AdoItemPicker from '@/components/ado/AdoItemPicker.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import { cn } from '@/lib/utils'
import { Plus, FolderOpen, ListTodo, CheckCircle2, Link2, Unlink } from 'lucide-vue-next'
import { useNotify } from '@/composables/useNotify'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const taskStore = useTaskStore()

const notify = useNotify()
const showCreate = ref(false)
const newName = ref('')
const newDescription = ref('')
const selectedProjectId = ref<number | null>(null)
const showLinkDialog = ref(false)

// Handle ?create=1 query param from top bar "New" dropdown
watch(() => route.query.create, (val) => {
  if (val === '1') showCreate.value = true
}, { immediate: true })

const selectedProject = computed(() =>
  projectStore.projects.find(p => p.id === selectedProjectId.value) ?? null
)

const projectTasks = computed(() =>
  selectedProjectId.value != null
    ? taskStore.tasks.filter(t => t.projectId === selectedProjectId.value)
    : []
)

const tasksByStatus = computed(() => {
  const groups: Record<string, Task[]> = {
    blocked: [], in_progress: [], in_review: [], todo: [], done: [],
  }
  for (const t of projectTasks.value) {
    if (groups[t.status]) groups[t.status].push(t)
  }
  return groups
})

const progress = computed(() =>
  selectedProjectId.value != null
    ? projectStore.projectProgress.get(selectedProjectId.value)
    : undefined
)

const adoLink = computed<ProjectADOLink | undefined>(() =>
  selectedProjectId.value != null
    ? projectStore.projectLinks.get(selectedProjectId.value)
    : undefined
)

const localPercent = computed(() => {
  if (!progress.value?.localTotal) return 0
  return Math.round((progress.value.localDone / progress.value.localTotal) * 100)
})

const adoPercent = computed(() => {
  if (!progress.value?.adoTotal) return 0
  return Math.round((progress.value.adoDone / progress.value.adoTotal) * 100)
})

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

function selectProject(id: number) {
  selectedProjectId.value = selectedProjectId.value === id ? null : id
}

watch(selectedProjectId, async (id) => {
  if (id != null) {
    await Promise.all([
      projectStore.fetchProjectLink(id),
      projectStore.fetchProjectProgress(id),
    ])
  }
})

onMounted(async () => {
  await Promise.all([
    projectStore.fetchProjects(),
    taskStore.fetchTasks(),
    taskStore.fetchPublicTaskIds(),
  ])
  await Promise.all(
    projectStore.projects.map(p =>
      Promise.all([
        projectStore.fetchProjectLink(p.id),
        projectStore.fetchProjectProgress(p.id),
      ])
    )
  )
})

async function handleLink(adoId: string) {
  if (selectedProjectId.value == null) return
  await projectStore.linkProjectToADO(selectedProjectId.value, adoId)
  showLinkDialog.value = false
  await projectStore.fetchProjectProgress(selectedProjectId.value)
  notify.success('Project linked to ADO')
}

async function handleUnlink() {
  const link = adoLink.value
  if (!link || selectedProjectId.value == null) return
  await projectStore.unlinkProject(selectedProjectId.value, link.adoId)
  await projectStore.fetchProjectProgress(selectedProjectId.value)
  notify.info('Project unlinked from ADO')
}

async function createProject() {
  const name = newName.value.trim()
  if (!name) return
  const p = await projectStore.createProject(name, newDescription.value.trim())
  newName.value = ''
  newDescription.value = ''
  showCreate.value = false
  if (p) {
    selectedProjectId.value = p.id
    notify.success('Project created')
  }
}
</script>

<template>
  <div class="p-6 space-y-4">
    <PageHeader>
      <template #left>
        <span class="text-xs text-muted-foreground">{{ projectStore.projects.length }} total</span>
      </template>
      <template #right>
        <Button size="sm" class="h-7 text-xs gap-1" @click="showCreate = !showCreate">
          <Plus :size="13" />
          New Project
        </Button>
      </template>
    </PageHeader>

    <!-- Inline create form -->
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="showCreate" class="max-w-md space-y-2">
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

    <!-- Split layout: projects left, tasks right -->
    <div class="flex gap-6 min-h-0">
      <!-- Left panel: project cards -->
      <div class="w-[45%] min-w-0 space-y-4 overflow-y-auto">
        <!-- Loading skeletons -->
        <div v-if="projectStore.loading" class="grid grid-cols-1 gap-3">
          <Skeleton v-for="n in 4" :key="n" class="h-36 rounded-lg" />
        </div>

        <!-- Empty state -->
        <EmptyState
          v-else-if="projectStore.projects.length === 0"
          :icon="FolderOpen"
          title="No projects yet"
          description="Projects help you group related tasks. Create one to get organized."
        >
          <template #action>
            <Button size="sm" @click="showCreate = true">Create Project</Button>
          </template>
        </EmptyState>

        <template v-else>
          <!-- Pinned projects -->
          <div v-if="projectStore.pinnedProjects.length">
            <div class="text-xs font-semibold text-muted-foreground uppercase mb-2">Pinned</div>
            <div class="grid grid-cols-1 gap-3">
              <ProjectCard
                v-for="p in projectStore.pinnedProjects"
                :key="p.id"
                :project="p"
                :is-linked="projectStore.isLinked(p.id)"
                :progress="projectStore.projectProgress.get(p.id) ?? null"
                :class="selectedProjectId === p.id ? 'ring-2 ring-primary' : ''"
                @click="selectProject(p.id)"
                @pin="projectStore.pinProject(p.id, $event)"
              />
            </div>
          </div>

          <!-- All projects (unpinned) -->
          <div>
            <div class="text-xs font-semibold text-muted-foreground uppercase mb-2">All Projects</div>
            <div class="grid grid-cols-1 gap-3">
              <ProjectCard
                v-for="p in projectStore.unpinnedProjects"
                :key="p.id"
                :project="p"
                :is-linked="projectStore.isLinked(p.id)"
                :progress="projectStore.projectProgress.get(p.id) ?? null"
                :class="selectedProjectId === p.id ? 'ring-2 ring-primary' : ''"
                @click="selectProject(p.id)"
                @pin="projectStore.pinProject(p.id, $event)"
              />
            </div>
          </div>
        </template>
      </div>

      <!-- Right panel: project detail -->
      <div class="w-[55%] min-w-0 border rounded-lg bg-card flex flex-col">
        <!-- No project selected -->
        <div v-if="!selectedProject" class="flex-1 flex flex-col items-center justify-center text-muted-foreground gap-2 p-8">
          <ListTodo :size="32" :stroke-width="1.25" />
          <span class="text-sm">Select a project to view details</span>
        </div>

        <template v-else>
          <!-- Header: name + description + ADO controls -->
          <div class="px-4 py-3 border-b space-y-1">
            <div class="flex items-center justify-between gap-2">
              <div class="min-w-0 flex-1">
                <h2 class="text-sm font-semibold truncate">{{ selectedProject.name }}</h2>
                <p v-if="selectedProject.description" class="text-xs text-muted-foreground truncate">{{ selectedProject.description }}</p>
              </div>
              <div class="flex gap-2 shrink-0">
                <Button v-if="!adoLink" variant="outline" size="sm" class="h-7 text-xs" @click="showLinkDialog = true">
                  <Link2 :size="12" class="mr-1" />
                  Link to ADO
                </Button>
                <template v-else>
                  <Badge variant="outline" class="gap-1 text-blue-500 border-blue-500/30 text-[10px]">
                    <AzureDevOpsIcon :size="10" />
                    ADO Linked
                  </Badge>
                  <Button variant="ghost" size="sm" class="h-7 text-xs" @click="handleUnlink">
                    <Unlink :size="12" class="mr-1" />
                    Unlink
                  </Button>
                </template>
              </div>
            </div>
          </div>

          <!-- ADO link picker dialog -->
          <AdoItemPicker
            :open="showLinkDialog"
            title="Link project to ADO work item"
            @update:open="showLinkDialog = $event"
            @selected="(adoId: string) => handleLink(adoId)"
          />

          <!-- Progress + tasks body -->
          <div class="flex-1 overflow-y-auto p-4 space-y-4">
            <!-- Progress bar -->
            <div class="space-y-2">
              <div class="flex items-center gap-3">
                <span class="text-xs font-medium w-16">Local</span>
                <div class="flex-1 h-1.5 rounded-full bg-muted">
                  <div class="h-full rounded-full bg-primary transition-all" :style="{ width: localPercent + '%' }" />
                </div>
                <span class="text-xs text-muted-foreground tabular-nums">{{ localPercent }}%</span>
              </div>
              <div v-if="adoLink" class="flex items-center gap-3">
                <span class="text-xs font-medium w-16">ADO</span>
                <div class="flex-1 h-1.5 rounded-full bg-muted">
                  <div class="h-full rounded-full bg-blue-500 transition-all" :style="{ width: adoPercent + '%' }" />
                </div>
                <span class="text-xs text-muted-foreground tabular-nums">{{ adoPercent }}%</span>
              </div>
            </div>

            <!-- Task list grouped by status -->
            <div>
              <div class="text-xs font-semibold mb-2">Tasks</div>

              <div v-if="projectTasks.length > 0" class="space-y-1">
                <div v-for="section in visibleSections" :key="section">
                  <!-- Section header -->
                  <div class="flex items-center gap-2 px-1 py-1.5">
                    <span :class="cn('w-1.5 h-1.5 rounded-full', sectionMeta[section].dot)" />
                    <span class="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground/60">
                      {{ sectionMeta[section].label }}
                    </span>
                    <span class="text-[11px] text-muted-foreground/40 tabular-nums">
                      {{ tasksByStatus[section]?.length ?? 0 }}
                    </span>
                    <div class="flex-1 h-px bg-border/50" />
                  </div>

                  <!-- Task rows -->
                  <div
                    v-for="task in tasksByStatus[section]"
                    :key="task.id"
                    class="flex items-center gap-2.5 px-3 py-1.5 rounded hover:bg-muted/40 transition-colors cursor-pointer"
                    @click="taskStore.selectTask(task.id)"
                  >
                    <button
                      @click.stop="taskStore.setStatus(task.id, task.status === 'done' ? 'todo' : 'done')"
                      :class="cn(
                        'size-[14px] rounded-full border-[1.5px] shrink-0 flex items-center justify-center transition-all hover:scale-110',
                        task.status === 'done'
                          ? 'bg-emerald-500 border-emerald-500'
                          : task.status === 'blocked'
                            ? 'border-red-400 hover:border-red-500'
                            : 'border-muted-foreground/60 hover:border-muted-foreground'
                      )"
                    >
                      <CheckCircle2 v-if="task.status === 'done'" :size="8" class="text-white" :stroke-width="3" />
                    </button>
                    <span
                      :class="cn(
                        'text-[13px] font-medium truncate flex-1',
                        task.status === 'done' ? 'text-muted-foreground line-through decoration-muted-foreground/30' : 'text-foreground'
                      )"
                    >
                      {{ task.title }}
                    </span>
                    <PriorityBadge v-if="task.priority" :priority="task.priority" />
                  </div>
                </div>
              </div>

              <div v-else class="text-xs text-muted-foreground text-center py-6">
                No tasks in this project
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
