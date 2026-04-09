import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Project, ProjectADOLink, ProjectProgress } from '@/types'
import * as projectsApi from '@/api/projects'

export type { Project, ProjectADOLink, ProjectProgress }

export const useProjectStore = defineStore('projects', () => {
  const projects = ref<Project[]>([])
  const loading = ref(false)
  const projectLinks = ref<Map<number, ProjectADOLink>>(new Map())
  const projectProgress = ref<Map<number, ProjectProgress>>(new Map())

  let fetchInFlight = false

  const pinnedProjects = computed(() => projects.value.filter(p => p.isPinned))
  const unpinnedProjects = computed(() => projects.value.filter(p => !p.isPinned))

  function isLinked(projectId: number): boolean {
    return projectLinks.value.has(projectId)
  }

  async function fetchProjects() {
    if (fetchInFlight) return
    fetchInFlight = true
    loading.value = true
    try {
      projects.value = (await projectsApi.listProjects()) as Project[]
    } catch (e) {
      console.warn('[ProjectStore] Wails binding unavailable:', e)
      projects.value = []
    } finally {
      loading.value = false
      fetchInFlight = false
    }
  }

  async function createProject(name: string, description: string) {
    const p = await projectsApi.createProject(name, description) as Project
    projects.value.push(p)
    return p
  }

  async function deleteProject(id: number) {
    await projectsApi.deleteProject(id)
    projects.value = projects.value.filter(p => p.id !== id)
  }

  async function updateProjectStatus(id: number, status: string) {
    const p = projects.value.find(p => p.id === id)
    if (!p) return
    const updated = await projectsApi.updateProject(id, p.name, p.description, status) as Project
    const idx = projects.value.findIndex(p => p.id === id)
    if (idx !== -1) projects.value[idx] = updated
  }

  async function pinProject(id: number, pinned: boolean) {
    try {
      await projectsApi.pinProject(id, pinned)
    } catch (e) {
      console.warn('[ProjectStore] PinProject binding unavailable:', e)
    }
    const idx = projects.value.findIndex(p => p.id === id)
    if (idx !== -1) projects.value[idx].isPinned = pinned
  }

  async function linkProjectToADO(projectId: number, adoId: string) {
    await projectsApi.linkProjectToADO(projectId, adoId, 'linked')
    projectLinks.value.set(projectId, {
      projectId,
      adoId,
      direction: 'linked',
      createdAt: new Date().toISOString(),
    })
  }

  async function unlinkProject(projectId: number, adoId: string) {
    await projectsApi.unlinkProject(projectId, adoId, false)
    projectLinks.value.delete(projectId)
  }

  async function fetchProjectLink(projectId: number) {
    try {
      const link = await projectsApi.getProjectADOLink(projectId) as ProjectADOLink
      if (link && link.adoId) projectLinks.value.set(projectId, link)
    } catch {
      /* no link */
    }
  }

  async function fetchAllProjectLinks() {
    for (const p of projects.value) {
      await fetchProjectLink(p.id)
    }
  }

  async function fetchProjectProgress(projectId: number) {
    try {
      const progress = await projectsApi.getProjectProgress(projectId) as ProjectProgress
      projectProgress.value.set(projectId, progress)
    } catch {
      /* ignore */
    }
  }

  return {
    projects, loading, projectLinks, projectProgress,
    pinnedProjects, unpinnedProjects,
    isLinked,
    fetchProjects, createProject, deleteProject, updateProjectStatus,
    pinProject, linkProjectToADO, unlinkProject,
    fetchProjectLink, fetchAllProjectLinks, fetchProjectProgress,
  }
})
