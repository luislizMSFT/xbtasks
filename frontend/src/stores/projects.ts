import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Project {
  id: number
  name: string
  description: string
  status: string
  isPinned: boolean
  createdAt: string
  updatedAt: string
  taskCount?: number
}

export interface ProjectADOLink {
  projectId: number
  adoId: string
  direction: string
  createdAt: string
}

export interface ProjectProgress {
  localDone: number
  localTotal: number
  adoDone: number
  adoTotal: number
}

export const useProjectStore = defineStore('projects', () => {
  const projects = ref<Project[]>([])
  const loading = ref(false)
  const projectLinks = ref<Map<number, ProjectADOLink>>(new Map())
  const projectProgress = ref<Map<number, ProjectProgress>>(new Map())

  const pinnedProjects = computed(() => projects.value.filter(p => p.isPinned))
  const unpinnedProjects = computed(() => projects.value.filter(p => !p.isPinned))

  function isLinked(projectId: number): boolean {
    return projectLinks.value.has(projectId)
  }

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

  async function pinProject(id: number, pinned: boolean) {
    try {
      const { PinProject } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
      await PinProject(id, pinned)
    } catch (e) {
      console.warn('[ProjectStore] PinProject binding unavailable:', e)
    }
    const idx = projects.value.findIndex(p => p.id === id)
    if (idx !== -1) projects.value[idx].isPinned = pinned
  }

  async function linkProjectToADO(projectId: number, adoId: string) {
    const { LinkProjectToADO } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
    await LinkProjectToADO(projectId, adoId, 'linked')
    projectLinks.value.set(projectId, {
      projectId,
      adoId,
      direction: 'linked',
      createdAt: new Date().toISOString(),
    })
  }

  async function unlinkProject(projectId: number, adoId: string) {
    const { UnlinkProject } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
    await UnlinkProject(projectId, adoId, false)
    projectLinks.value.delete(projectId)
  }

  async function fetchProjectLink(projectId: number) {
    try {
      const { GetProjectADOLink } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
      const link = await GetProjectADOLink(projectId) as ProjectADOLink
      if (link && link.adoId) projectLinks.value.set(projectId, link)
    } catch {
      /* no link */
    }
  }

  async function fetchProjectProgress(projectId: number) {
    try {
      const { GetProjectProgress } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
      const progress = await GetProjectProgress(projectId) as ProjectProgress
      projectProgress.value.set(projectId, progress)
    } catch {
      /* ignore */
    }
  }

  return {
    projects, loading, projectLinks, projectProgress,
    pinnedProjects, unpinnedProjects,
    isLinked,
    fetchProjects, createProject, deleteProject,
    pinProject, linkProjectToADO, unlinkProject,
    fetchProjectLink, fetchProjectProgress,
  }
})
