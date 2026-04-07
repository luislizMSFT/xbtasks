<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore, type ProjectADOLink } from '@/stores/projects'
import { useTaskStore, type Task } from '@/stores/tasks'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/EmptyState.vue'
import { cn } from '@/lib/utils'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import AdoItemPicker from '@/components/ado/AdoItemPicker.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import {
  ArrowLeft, CheckCircle2, ClipboardList, Link2, Unlink,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const taskStore = useTaskStore()

const projectId = computed(() => Number(route.params.id))
const project = computed(() => projectStore.projects.find(p => p.id === projectId.value))
const progress = computed(() => projectStore.projectProgress.get(projectId.value))
const adoLink = computed<ProjectADOLink | undefined>(() => projectStore.projectLinks.get(projectId.value))
const loading = ref(true)
const showLinkDialog = ref(false)

const projectTasks = computed(() =>
  taskStore.tasks.filter(t => t.projectId === projectId.value)
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

async function loadProjectData() {
  const id = projectId.value
  if (!id) return
  loading.value = true
  try {
    if (!projectStore.projects.length) await projectStore.fetchProjects()
    if (!taskStore.tasks.length) await taskStore.fetchTasks()
    await Promise.all([
      projectStore.fetchProjectLink(id),
      projectStore.fetchProjectProgress(id),
    ])
  } finally {
    loading.value = false
  }
}

async function handleLink(adoId: string) {
  await projectStore.linkProjectToADO(projectId.value, adoId)
  showLinkDialog.value = false
  await projectStore.fetchProjectProgress(projectId.value)
}

async function handleUnlink() {
  const link = adoLink.value
  if (!link) return
  await projectStore.unlinkProject(projectId.value, link.adoId)
  await projectStore.fetchProjectProgress(projectId.value)
}

onMounted(loadProjectData)
watch(projectId, loadProjectData)
</script>

<template>
  <div class="flex-1 overflow-auto">
    <div class="p-6 space-y-6 max-w-4xl">
      <!-- ===== Loading skeleton ===== -->
      <template v-if="loading">
        <!-- Header skeleton -->
        <div class="flex items-center gap-3">
          <Skeleton class="h-8 w-8 rounded-md" />
          <div class="flex-1 space-y-2">
            <Skeleton class="h-5 w-48" />
            <Skeleton class="h-4 w-72" />
          </div>
          <Skeleton class="h-8 w-28 rounded-md" />
        </div>

        <!-- Stats skeleton -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Card v-for="i in 4" :key="i">
            <CardContent class="pt-4 space-y-2">
              <Skeleton class="h-7 w-12" />
              <Skeleton class="h-3 w-20" />
            </CardContent>
          </Card>
        </div>

        <!-- Progress bar skeleton -->
        <div class="space-y-2">
          <div class="flex items-center gap-3">
            <Skeleton class="h-4 w-20" />
            <Skeleton class="h-2 flex-1 rounded-full" />
            <Skeleton class="h-4 w-8" />
          </div>
        </div>

        <!-- Task list skeleton -->
        <div class="space-y-1">
          <Skeleton class="h-4 w-16 mb-3" />
          <div v-for="i in 5" :key="i" class="flex items-center gap-2.5 px-3 py-2">
            <Skeleton class="h-4 w-4 rounded-full shrink-0" />
            <Skeleton class="h-4 flex-1" />
            <Skeleton class="h-5 w-14 rounded-full" />
          </div>
        </div>
      </template>

      <!-- ===== Loaded content ===== -->
      <template v-else>
        <!-- Back button + title -->
        <div class="flex items-center gap-3">
          <Button variant="ghost" size="icon" class="h-8 w-8" @click="router.push('/projects')">
            <ArrowLeft :size="16" />
          </Button>
          <div class="flex-1 min-w-0">
            <h1 class="text-lg font-semibold truncate">{{ project?.name }}</h1>
            <p v-if="project?.description" class="text-sm text-muted-foreground truncate">{{ project.description }}</p>
          </div>
          <div class="flex gap-2 shrink-0">
            <!-- ADO link/unlink -->
            <Button v-if="!adoLink" variant="outline" size="sm" @click="showLinkDialog = true">
              <Link2 :size="14" class="mr-1" />
              Link to ADO
            </Button>
            <template v-else>
              <Badge variant="outline" class="gap-1 text-blue-500 border-blue-500/30">
                <AzureDevOpsIcon :size="12" />
                ADO Linked
              </Badge>
              <Button variant="ghost" size="sm" @click="handleUnlink">
                <Unlink :size="14" class="mr-1" />
                Unlink
              </Button>
            </template>
          </div>
        </div>

        <!-- Link dialog -->
        <AdoItemPicker
          :open="showLinkDialog"
          title="Link project to ADO work item"
          @update:open="showLinkDialog = $event"
          @selected="(adoId) => handleLink(adoId)"
        />

        <!-- Stats cards row -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Card>
            <CardContent class="pt-4">
              <div class="text-2xl font-bold">{{ progress?.localTotal ?? 0 }}</div>
              <div class="text-xs text-muted-foreground">Total Tasks</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="pt-4">
              <div class="text-2xl font-bold text-green-500">{{ progress?.localDone ?? 0 }}</div>
              <div class="text-xs text-muted-foreground">Completed</div>
            </CardContent>
          </Card>
          <Card v-if="adoLink">
            <CardContent class="pt-4">
              <div class="text-2xl font-bold text-blue-500">{{ progress?.adoTotal ?? 0 }}</div>
              <div class="text-xs text-muted-foreground">ADO Children</div>
            </CardContent>
          </Card>
          <Card v-if="adoLink">
            <CardContent class="pt-4">
              <div class="text-2xl font-bold text-blue-500">{{ adoPercent }}%</div>
              <div class="text-xs text-muted-foreground">ADO Progress</div>
            </CardContent>
          </Card>
        </div>

        <!-- Dual progress bars (PROJ-06) -->
        <div class="space-y-2">
          <div class="flex items-center gap-3">
            <span class="text-sm font-medium w-20">Local Tasks</span>
            <div class="flex-1 h-2 rounded-full bg-muted">
              <div class="h-full rounded-full bg-primary transition-all" :style="{ width: localPercent + '%' }" />
            </div>
            <span class="text-sm text-muted-foreground tabular-nums">{{ localPercent }}%</span>
          </div>
          <div v-if="adoLink" class="flex items-center gap-3">
            <span class="text-sm font-medium w-20">ADO Items</span>
            <div class="flex-1 h-2 rounded-full bg-muted">
              <div class="h-full rounded-full bg-blue-500 transition-all" :style="{ width: adoPercent + '%' }" />
            </div>
            <span class="text-sm text-muted-foreground tabular-nums">{{ adoPercent }}%</span>
          </div>
        </div>

        <!-- ADO context (if linked) -->
        <Card v-if="adoLink">
          <CardHeader class="pb-2">
            <CardTitle class="text-sm">ADO Context</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="flex items-center gap-2 text-sm">
              <AzureDevOpsIcon :size="14" />
              <span>Linked to ADO item {{ adoLink.adoId }} ({{ adoLink.direction }})</span>
            </div>
          </CardContent>
        </Card>

        <!-- Filtered tasks for this project -->
        <div>
          <div class="text-sm font-semibold mb-3">Tasks</div>

          <div v-if="projectTasks.length > 0" class="space-y-1">
            <div v-for="section in visibleSections" :key="section">
              <!-- Section header -->
              <div class="flex items-center gap-2 px-1 py-2">
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
                class="flex items-center gap-2.5 px-3 py-2 rounded hover:bg-muted/40 transition-colors cursor-pointer"
                @click="taskStore.selectTask(task.id)"
              >
                <!-- Done checkbox -->
                <button
                  @click.stop="taskStore.setStatus(task.id, task.status === 'done' ? 'todo' : 'done')"
                  :class="cn(
                    'size-[16px] rounded-full border-[1.5px] shrink-0 flex items-center justify-center transition-all hover:scale-110',
                    task.status === 'done'
                      ? 'bg-emerald-500 border-emerald-500'
                      : task.status === 'blocked'
                        ? 'border-red-400 hover:border-red-500'
                        : 'border-muted-foreground/60 hover:border-muted-foreground'
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
              </div>
            </div>
          </div>

          <EmptyState
            v-else
            :icon="ClipboardList"
            title="No tasks yet"
            description="This project doesn't have any tasks. Create one to get started."
            class="py-8"
          />
        </div>
      </template>
    </div>
  </div>
</template>
