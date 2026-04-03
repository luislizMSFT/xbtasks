import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface ADOWorkItem {
  id: number
  adoId: string
  title: string
  state: string
  type: string
  assignedTo: string
  priority: number
  areaPath: string
  description: string
  url: string
  syncedAt: string
}

export interface ADOPipeline {
  id: number
  name: string
  status: string
  result: string
  url: string
  sourceBranch: string
  queueTime: string
  finishTime: string | null
}

export const useADOStore = defineStore('ado', () => {
  const workItems = ref<ADOWorkItem[]>([])
  const pipelines = ref<ADOPipeline[]>([])
  const loading = ref(false)
  const connected = ref(false)
  const syncing = ref(false)
  const error = ref('')
  const lastSyncedAt = ref<string | null>(null)

  async function fetchWorkItems() {
    error.value = ''
    try {
      const { ListMyWorkItems } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
      workItems.value = (await ListMyWorkItems()) as ADOWorkItem[]
      connected.value = true
    } catch (e: any) {
      console.warn('[ADOStore] ListMyWorkItems failed:', e)
      error.value = e?.message || 'Failed to fetch work items. Check az cli auth.'
      connected.value = false
      workItems.value = []
    }
  }

  async function syncWorkItems() {
    syncing.value = true
    error.value = ''
    try {
      const { SyncWorkItems } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
      await SyncWorkItems()
      connected.value = true
      await fetchWorkItems()
      lastSyncedAt.value = new Date().toISOString()
    } catch (e: any) {
      console.warn('[ADOStore] SyncWorkItems failed:', e)
      error.value = e?.message || 'Sync failed. Check az cli auth.'
      await fetchCached()
      lastSyncedAt.value = new Date().toISOString()
    } finally {
      syncing.value = false
    }
  }

  async function fetchCached() {
    try {
      const { GetCachedWorkItems } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
      workItems.value = (await GetCachedWorkItems()) as ADOWorkItem[]
    } catch {
      workItems.value = []
    }
  }

  async function fetchPipelines() {
    try {
      const { ListRecentRuns } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/pipelineservice')
      pipelines.value = (await ListRecentRuns()) as ADOPipeline[]
      connected.value = true
    } catch (e: any) {
      console.warn('[ADOStore] ListRecentRuns failed:', e)
      if (!error.value) error.value = e?.message || 'Failed to fetch pipelines'
      pipelines.value = []
    }
  }

  async function fetchAll() {
    loading.value = true
    error.value = ''
    try {
      await Promise.all([fetchWorkItems(), fetchPipelines()])
      if (!lastSyncedAt.value && connected.value) {
        lastSyncedAt.value = new Date().toISOString()
      }
    } finally {
      loading.value = false
    }
  }

  return {
    workItems, pipelines, loading, connected, syncing, error, lastSyncedAt,
    fetchWorkItems, syncWorkItems, fetchCached, fetchPipelines, fetchAll,
  }
})
