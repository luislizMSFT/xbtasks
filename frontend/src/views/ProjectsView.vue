<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useProjectStore } from '../stores/projects'
import { Plus, Folder, X } from 'lucide-vue-next'

const projectStore = useProjectStore()

const showCreate = ref(false)
const newName = ref('')
const newDescription = ref('')

onMounted(() => {
  projectStore.fetchProjects()
})

async function createProject() {
  const name = newName.value.trim()
  if (!name) return
  await projectStore.createProject(name, newDescription.value.trim())
  newName.value = ''
  newDescription.value = ''
  showCreate.value = false
}

const statusColors: Record<string, string> = {
  active: 'bg-emerald-500',
  paused: 'bg-amber-500',
  completed: 'bg-blue-500',
  archived: 'bg-zinc-400',
  blocked: 'bg-red-500',
}
</script>

<template>
  <div class="flex-1 overflow-y-auto">
    <div class="max-w-4xl mx-auto px-6 py-6 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold text-text-primary">Projects</h1>
          <p class="text-sm text-text-secondary mt-0.5">Organize tasks by project</p>
        </div>
        <button
          @click="showCreate = true"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
        >
          <Plus :size="14" />
          New Project
        </button>
      </div>

      <!-- Create form -->
      <Transition
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-100 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div v-if="showCreate" class="rounded-xl border border-border-default bg-surface-secondary p-4 space-y-3">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-text-primary">New Project</h3>
            <button @click="showCreate = false" class="text-text-secondary hover:text-text-primary">
              <X :size="16" />
            </button>
          </div>
          <input
            v-model="newName"
            @keydown.enter="createProject"
            class="w-full px-3 py-2 rounded-md text-sm bg-surface-tertiary border border-border-default text-text-primary placeholder-text-secondary outline-none focus:ring-1 focus:ring-accent"
            placeholder="Project name"
            autofocus
          />
          <textarea
            v-model="newDescription"
            rows="2"
            class="w-full px-3 py-2 rounded-md text-sm bg-surface-tertiary border border-border-default text-text-primary placeholder-text-secondary outline-none focus:ring-1 focus:ring-accent resize-none"
            placeholder="Description (optional)"
          />
          <div class="flex justify-end gap-2">
            <button
              @click="showCreate = false"
              class="px-3 py-1.5 rounded-md text-xs text-text-secondary hover:text-text-primary hover:bg-surface-tertiary transition-colors"
            >
              Cancel
            </button>
            <button
              @click="createProject"
              class="px-3 py-1.5 rounded-md text-xs font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
            >
              Create
            </button>
          </div>
        </div>
      </Transition>

      <!-- Project grid -->
      <div v-if="projectStore.projects.length > 0" class="grid grid-cols-2 gap-4">
        <div
          v-for="project in projectStore.projects"
          :key="project.id"
          class="rounded-xl border border-border-default bg-surface-secondary p-4 hover:border-accent/30 cursor-pointer transition-colors group"
        >
          <div class="flex items-start gap-3">
            <div class="w-9 h-9 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0">
              <Folder :size="18" class="text-accent" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <h3 class="text-sm font-semibold text-text-primary truncate">{{ project.name }}</h3>
                <span
                  class="w-2 h-2 rounded-full flex-shrink-0"
                  :class="statusColors[project.status] || 'bg-zinc-400'"
                />
              </div>
              <p class="text-xs text-text-secondary mt-1 line-clamp-2">{{ project.description || 'No description' }}</p>
              <div class="flex items-center gap-3 mt-2 text-[11px] text-text-secondary">
                <span>{{ project.taskCount ?? 0 }} tasks</span>
                <span class="capitalize">{{ project.status }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else-if="!projectStore.loading" class="text-center py-16">
        <div class="w-12 h-12 rounded-full bg-surface-tertiary flex items-center justify-center mx-auto mb-3">
          <Folder :size="24" class="text-text-secondary" />
        </div>
        <p class="text-sm font-medium text-text-primary">No projects yet</p>
        <p class="text-xs text-text-secondary mt-1">Create a project to organize your tasks</p>
        <button
          @click="showCreate = true"
          class="mt-3 px-4 py-2 rounded-md text-sm font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
        >
          Create Project
        </button>
      </div>
    </div>
  </div>
</template>
