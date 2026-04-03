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

const MOCK_PROJECTS: Project[] = [
  {
    id: 1, name: 'Team ADO Tool', description: 'Unified desktop dashboard for tasks, ADO work items, and PRs.',
    status: 'active', createdAt: new Date(Date.now() - 2592000000).toISOString(),
    updatedAt: new Date(Date.now() - 86400000).toISOString(), taskCount: 6,
  },
  {
    id: 2, name: 'API Gateway', description: 'Microservice gateway handling auth, rate limiting, and routing.',
    status: 'active', createdAt: new Date(Date.now() - 5184000000).toISOString(),
    updatedAt: new Date(Date.now() - 172800000).toISOString(), taskCount: 3,
  },
]

let useMock = false

export const useProjectStore = defineStore('projects', () => {
  const projects = ref<Project[]>([])
  const loading = ref(false)

  async function fetchProjects() {
    loading.value = true
    try {
      if (!useMock) {
        const { List } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
        projects.value = (await List()) as Project[]
      } else {
        await new Promise(r => setTimeout(r, 150))
        projects.value = [...MOCK_PROJECTS]
      }
    } catch (e) {
      console.warn('[ProjectStore] Wails binding unavailable, using mock data:', e)
      useMock = true
      projects.value = [...MOCK_PROJECTS]
    } finally {
      loading.value = false
    }
  }

  async function createProject(name: string, description: string) {
    if (useMock) {
      const p: Project = {
        id: Math.max(0, ...projects.value.map(p => p.id)) + 1,
        name, description, status: 'active',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(), taskCount: 0,
      }
      projects.value.push(p)
      return p
    }
    const { Create } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
    const p = await Create(name, description) as Project
    projects.value.push(p)
    return p
  }

  async function deleteProject(id: number) {
    if (useMock) {
      projects.value = projects.value.filter(p => p.id !== id)
      return
    }
    const { Delete } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
    await Delete(id)
    projects.value = projects.value.filter(p => p.id !== id)
  }

  return { projects, loading, fetchProjects, createProject, deleteProject }
})
