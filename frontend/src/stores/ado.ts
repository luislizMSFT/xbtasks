import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import type { ADOWorkItem, ADOPipeline, SavedQuery } from '@/types'
import * as workitemsApi from '@/api/workitems'
import { listRecentRuns } from '@/api/pipelines'
import { listLinkedAdoIDs } from '@/api/links'

export type { ADOWorkItem, ADOPipeline, SavedQuery }

export const useADOStore = defineStore('ado', () => {
  const workItems = ref<ADOWorkItem[]>([])
  const pipelines = ref<ADOPipeline[]>([])
  const loading = ref(false)
  const pipelinesLoading = ref(false)
  const connected = ref(false)
  const syncing = ref(false)
  const error = ref('')
  const lastSyncedAt = ref<string | null>(null)

  let fetchTreeInFlight = false

  // Tree browser state
  const workItemTree = ref<ADOWorkItem[]>([])
  const linkedAdoIds = ref<Set<string>>(new Set())
  const hideLinked = ref(false)
  const searchQuery = ref('')
  const debouncedSearchQuery = ref('')
  const filterType = ref<string>('all')
  const filterState = ref<string>('all')
  const filterArea = ref<string>('all')
  const savedQueries = ref<SavedQuery[]>([])
  const selectedItem = ref<ADOWorkItem | null>(null)

  // Debounce search input so treeRoots doesn't recalc on every keystroke
  const applySearch = useDebounceFn((q: string) => {
    debouncedSearchQuery.value = q
  }, 200)

  watch(searchQuery, (q) => applySearch(q))

  // Computed: available area paths from tree data
  const availableAreaPaths = computed(() => {
    const paths = new Set(workItemTree.value.map(i => i.areaPath).filter(Boolean))
    return Array.from(paths).sort()
  })

  // Computed: tree roots after filtering (uses debounced search for perf)
  const treeRoots = computed(() => {
    let items = workItemTree.value
    if (hideLinked.value) items = items.filter(i => !isLinked(i.adoId))
    if (debouncedSearchQuery.value) {
      const q = debouncedSearchQuery.value.toLowerCase()
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
    // Don't filter children by type/state/area — show all children regardless
    // so the hierarchy is always visible even when filtering roots
    return items.filter(i => i.parentId === parentId)
  }

  function isLinked(adoId: string): boolean {
    return linkedAdoIds.value.has(adoId)
  }

  async function fetchWorkItems() {
    error.value = ''
    try {
      workItems.value = (await workitemsApi.listMyWorkItems()) as ADOWorkItem[]
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
      await workitemsApi.syncWorkItems()
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
      workItems.value = (await workitemsApi.getCachedWorkItems()) as ADOWorkItem[]
    } catch {
      workItems.value = []
    }
  }

  async function fetchPipelines() {
    pipelinesLoading.value = true
    try {
      pipelines.value = (await listRecentRuns()) as ADOPipeline[]
      connected.value = true
    } catch (e: any) {
      console.warn('[ADOStore] ListRecentRuns failed:', e)
      if (!error.value) error.value = e?.message || 'Failed to fetch pipelines'
      pipelines.value = []
    } finally {
      pipelinesLoading.value = false
    }
  }

  async function fetchAll() {
    loading.value = true
    error.value = ''
    try {
      await Promise.allSettled([fetchWorkItems(), fetchPipelines()])
      if (!lastSyncedAt.value && connected.value) {
        lastSyncedAt.value = new Date().toISOString()
      }
    } finally {
      loading.value = false
    }
  }

  async function fetchWorkItemTree() {
    if (fetchTreeInFlight) return
    fetchTreeInFlight = true
    loading.value = true
    error.value = ''
    try {
      const items = (await workitemsApi.getWorkItemTree()) as ADOWorkItem[]
      // Normalize: ensure id = Number(adoId) for proper parent-child matching.
      // GetWorkItemTree returns items from the API where domain.WorkItem.ID is
      // unset (0) but AdoID and ParentID are correct ADO IDs.
      for (const item of items) {
        if (!item.id && item.adoId) {
          item.id = Number(item.adoId)
        }
      }
      workItemTree.value = items
      connected.value = true
    } catch (e: any) {
      console.warn('[ADOStore] GetWorkItemTree failed:', e)
      error.value = e?.message || 'Failed to fetch work item tree'
      workItemTree.value = []
    } finally {
      loading.value = false
      fetchTreeInFlight = false
    }
  }

  async function fetchLinkedAdoIds() {
    try {
      const ids = (await listLinkedAdoIDs()) as string[]
      linkedAdoIds.value = new Set(ids || [])
    } catch {
      linkedAdoIds.value = new Set()
    }
  }

  async function fetchSavedQueries() {
    try {
      savedQueries.value = (await workitemsApi.getSavedQueries()) as SavedQuery[]
    } catch {
      savedQueries.value = []
    }
  }

  async function runSavedQuery(queryId: string) {
    loading.value = true
    error.value = ''
    try {
      const items = (await workitemsApi.runSavedQuery(queryId)) as ADOWorkItem[]
      for (const item of items) {
        if (!item.id && item.adoId) {
          item.id = Number(item.adoId)
        }
      }
      workItemTree.value = items
      connected.value = true
    } catch (e: any) {
      error.value = e?.message || 'Failed to run saved query'
    } finally {
      loading.value = false
    }
  }

  return {
    // Existing
    workItems, pipelines, loading, pipelinesLoading, connected, syncing, error, lastSyncedAt,
    fetchWorkItems, syncWorkItems, fetchCached, fetchPipelines, fetchAll,
    // Tree browser
    workItemTree, linkedAdoIds, hideLinked, searchQuery, filterType, filterState, filterArea,
    savedQueries, selectedItem,
    treeRoots, availableAreaPaths, getChildren, isLinked,
    fetchWorkItemTree, fetchLinkedAdoIds, fetchSavedQueries, runSavedQuery,
  }
})
