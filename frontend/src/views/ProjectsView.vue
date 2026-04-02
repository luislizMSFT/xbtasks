<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useProjectStore } from '@/stores/projects'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { ScrollArea } from '@/components/ui/scroll-area'
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
  <ScrollArea class="flex-1">
    <div class="max-w-4xl mx-auto px-6 py-6 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold text-foreground">Projects</h1>
          <p class="text-sm text-muted-foreground mt-0.5">Organize tasks by project</p>
        </div>
        <Button size="sm" @click="showCreate = true">
          <Plus :size="14" />
          New Project
        </Button>
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
        <Card v-if="showCreate">
          <CardHeader class="p-4 pb-0">
            <div class="flex items-center justify-between">
              <CardTitle class="text-sm">New Project</CardTitle>
              <Button variant="ghost" size="icon-sm" @click="showCreate = false">
                <X :size="16" />
              </Button>
            </div>
          </CardHeader>
          <CardContent class="p-4 space-y-3">
            <Input
              v-model="newName"
              @keydown.enter="createProject"
              placeholder="Project name"
              autofocus
            />
            <Textarea
              v-model="newDescription"
              :rows="2"
              class="resize-none"
              placeholder="Description (optional)"
            />
          </CardContent>
          <CardFooter class="px-4 pb-4 pt-0 flex justify-end gap-2">
            <Button variant="ghost" size="sm" @click="showCreate = false">
              Cancel
            </Button>
            <Button size="sm" @click="createProject">
              Create
            </Button>
          </CardFooter>
        </Card>
      </Transition>

      <!-- Project grid -->
      <div v-if="projectStore.projects.length > 0" class="grid grid-cols-2 gap-4">
        <Card
          v-for="project in projectStore.projects"
          :key="project.id"
          class="hover:border-primary/30 cursor-pointer transition-colors group"
        >
          <CardContent class="p-4">
            <div class="flex items-start gap-3">
              <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                <Folder :size="18" class="text-primary" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <h3 class="text-sm font-semibold text-foreground truncate">{{ project.name }}</h3>
                  <span
                    class="w-2 h-2 rounded-full flex-shrink-0"
                    :class="statusColors[project.status] || 'bg-zinc-400'"
                  />
                </div>
                <p class="text-xs text-muted-foreground mt-1 line-clamp-2">{{ project.description || 'No description' }}</p>
                <div class="flex items-center gap-3 mt-2 text-[11px] text-muted-foreground">
                  <span>{{ project.taskCount ?? 0 }} tasks</span>
                  <span class="capitalize">{{ project.status }}</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Empty state -->
      <div v-else-if="!projectStore.loading" class="text-center py-16">
        <div class="w-12 h-12 rounded-full bg-muted flex items-center justify-center mx-auto mb-3">
          <Folder :size="24" class="text-muted-foreground" />
        </div>
        <p class="text-sm font-medium text-foreground">No projects yet</p>
        <p class="text-xs text-muted-foreground mt-1">Create a project to organize your tasks</p>
        <Button class="mt-3" @click="showCreate = true">
          Create Project
        </Button>
      </div>
    </div>
  </ScrollArea>
</template>
