import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Project {
  id: number
  name: string
  description: string
  status: string
  createdAt: string
  updatedAt: string
  taskCount?: number
}

export const useProjectStore = defineStore('projects', () => {
  const projects = ref<Project[]>([])
  const loading = ref(false)

  async function fetchProjects() {
    loading.value = true
    try {
      const { List } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
      projects.value = (await List()) as Project[]
    } catch (e) {
      console.warn('[ProjectStore] Wails binding unavailable:', e)
      projects.value = []
    } finally {
      loading.value = false
    }
  }

  async function createProject(name: string, description: string) {
    const { Create } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
    const p = await Create(name, description) as Project
    projects.value.push(p)
    return p
  }

  async function deleteProject(id: number) {
    const { Delete } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
    await Delete(id)
    projects.value = projects.value.filter(p => p.id !== id)
  }

  return { projects, loading, fetchProjects, createProject, deleteProject }
})
