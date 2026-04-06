import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

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
  org: string
  project: string
  parentId: number
  changedDate: string
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

export interface SavedQuery {
  id: string
  name: string
  path: string
  isFolder: boolean
}

export const useADOStore = defineStore('ado', () => {
  const workItems = ref<ADOWorkItem[]>([])
  const pipelines = ref<ADOPipeline[]>([])
  const loading = ref(false)
  const connected = ref(false)
  const syncing = ref(false)
  const error = ref('')
  const lastSyncedAt = ref<string | null>(null)

  // Tree browser state
  const workItemTree = ref<ADOWorkItem[]>([])
  const linkedAdoIds = ref<Set<string>>(new Set())
  const hideLinked = ref(false)
  const searchQuery = ref('')
  const filterType = ref<string>('all')
  const filterState = ref<string>('all')
  const filterArea = ref<string>('all')
  const savedQueries = ref<SavedQuery[]>([])
  const selectedItem = ref<ADOWorkItem | null>(null)

  // Computed: available area paths from tree data
  const availableAreaPaths = computed(() => {
    const paths = new Set(workItemTree.value.map(i => i.areaPath).filter(Boolean))
    return Array.from(paths).sort()
  })

  // Computed: tree roots after filtering
  const treeRoots = computed(() => {
    let items = workItemTree.value
    if (hideLinked.value) items = items.filter(i => !isLinked(i.adoId))
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      items = items.filter(i => i.title.toLowerCase().includes(q) || i.adoId.includes(q))
    }
    if (filterType.value !== 'all') items = items.filter(i => i.type === filterType.value)
    if (filterState.value !== 'all') items = items.filter(i => i.state === filterState.value)
    if (filterArea.value !== 'all') items = items.filter(i => i.areaPath.startsWith(filterArea.value))
    // Items with no parentId or parentId not in the filtered set are roots
    const idSet = new Set(items.map(i => i.id))
    return items.filter(i => !i.parentId || !idSet.has(i.parentId))
  })

  function getChildren(parentId: number): ADOWorkItem[] {
    let items = workItemTree.value
    if (hideLinked.value) items = items.filter(i => !isLinked(i.adoId))
    if (filterType.value !== 'all') items = items.filter(i => i.type === filterType.value)
    if (filterState.value !== 'all') items = items.filter(i => i.state === filterState.value)
    if (filterArea.value !== 'all') items = items.filter(i => i.areaPath.startsWith(filterArea.value))
    return items.filter(i => i.parentId === parentId)
  }

  function isLinked(adoId: string): boolean {
    return linkedAdoIds.value.has(adoId)
  }

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

  async function fetchWorkItemTree() {
    loading.value = true
    error.value = ''
    try {
      const { GetWorkItemTree } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
      workItemTree.value = (await GetWorkItemTree()) as ADOWorkItem[]
      connected.value = true
    } catch (e: any) {
      console.warn('[ADOStore] GetWorkItemTree failed:', e)
      error.value = e?.message || 'Failed to fetch work item tree'
      workItemTree.value = []
    } finally {
      loading.value = false
    }
  }

  async function fetchLinkedAdoIds() {
    try {
      const { ListLinkedAdoIDs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
      const ids = (await ListLinkedAdoIDs()) as string[]
      linkedAdoIds.value = new Set(ids || [])
    } catch {
      linkedAdoIds.value = new Set()
    }
  }

  async function fetchSavedQueries() {
    try {
      const { GetSavedQueries } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
      savedQueries.value = (await GetSavedQueries()) as SavedQuery[]
    } catch {
      savedQueries.value = []
    }
  }

  async function runSavedQuery(queryId: string) {
    loading.value = true
    error.value = ''
    try {
      const { RunSavedQuery } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
      workItemTree.value = (await RunSavedQuery(queryId)) as ADOWorkItem[]
      connected.value = true
    } catch (e: any) {
      error.value = e?.message || 'Failed to run saved query'
    } finally {
      loading.value = false
    }
  }

  return {
    // Existing
    workItems, pipelines, loading, connected, syncing, error, lastSyncedAt,
    fetchWorkItems, syncWorkItems, fetchCached, fetchPipelines, fetchAll,
    // Tree browser
    workItemTree, linkedAdoIds, hideLinked, searchQuery, filterType, filterState, filterArea,
    savedQueries, selectedItem,
    treeRoots, availableAreaPaths, getChildren, isLinked,
    fetchWorkItemTree, fetchLinkedAdoIds, fetchSavedQueries, runSavedQuery,
  }
})
