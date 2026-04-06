<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/projects'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import PageHeader from '@/components/PageHeader.vue'
import ProjectCard from '@/components/ProjectCard.vue'
import { Plus, Folder } from 'lucide-vue-next'

const router = useRouter()
const projectStore = useProjectStore()

const showCreate = ref(false)
const newName = ref('')
const newDescription = ref('')

onMounted(async () => {
  await projectStore.fetchProjects()
  // Fetch link + progress data for each project
  await Promise.all(
    projectStore.projects.map(p =>
      Promise.all([
        projectStore.fetchProjectLink(p.id),
        projectStore.fetchProjectProgress(p.id),
      ])
    )
  )
})

async function createProject() {
  const name = newName.value.trim()
  if (!name) return
  const p = await projectStore.createProject(name, newDescription.value.trim())
  newName.value = ''
  newDescription.value = ''
  showCreate.value = false
  if (p) router.push('/projects/' + p.id)
}
</script>

<template>
  <div class="p-6 space-y-6">
    <PageHeader>
      <template #left>
        <span class="text-sm font-semibold">Projects</span>
        <span class="text-xs text-muted-foreground ml-2">{{ projectStore.projects.length }} total</span>
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

    <!-- Loading -->
    <div v-if="projectStore.loading" class="flex items-center justify-center py-20">
      <div class="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
    </div>

    <!-- Empty state -->
    <div v-else-if="projectStore.projects.length === 0" class="flex flex-col items-center justify-center py-16 gap-2">
      <Folder :size="24" class="text-muted-foreground/40" />
      <p class="text-sm font-medium" style="color: var(--color-text-primary)">No projects yet</p>
      <p class="text-xs" style="color: var(--color-text-secondary)">Projects help you group related tasks. Create one to get organized.</p>
      <Button size="sm" class="mt-2" @click="showCreate = true">Create Project</Button>
    </div>

    <template v-else>
      <!-- Pinned projects -->
      <div v-if="projectStore.pinnedProjects.length">
        <div class="text-xs font-semibold text-muted-foreground uppercase mb-2">Pinned</div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <ProjectCard
            v-for="p in projectStore.pinnedProjects"
            :key="p.id"
            :project="p"
            :is-linked="projectStore.isLinked(p.id)"
            :progress="projectStore.projectProgress.get(p.id) ?? null"
            @click="router.push('/projects/' + p.id)"
            @pin="projectStore.pinProject(p.id, $event)"
          />
        </div>
      </div>

      <!-- All projects (unpinned) -->
      <div>
        <div class="text-xs font-semibold text-muted-foreground uppercase mb-2">All Projects</div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <ProjectCard
            v-for="p in projectStore.unpinnedProjects"
            :key="p.id"
            :project="p"
            :is-linked="projectStore.isLinked(p.id)"
            :progress="projectStore.projectProgress.get(p.id) ?? null"
            @click="router.push('/projects/' + p.id)"
            @pin="projectStore.pinProject(p.id, $event)"
          />
        </div>
      </div>
    </template>
  </div>
</template>
