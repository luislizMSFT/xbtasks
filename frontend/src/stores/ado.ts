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

const MOCK_WORK_ITEMS: ADOWorkItem[] = [
  {
    id: 1, adoId: 'ADO-48291', title: 'Fix auth redirect loop', state: 'Active', type: 'Bug',
    assignedTo: 'Jane Smith', priority: 1, areaPath: 'Xbox\\Live Services\\Auth',
    description: 'Auth redirect creates an infinite loop when the session token expires.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48291',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
  {
    id: 2, adoId: 'ADO-48292', title: 'Token expiry timing off by 5 minutes', state: 'Active', type: 'Bug',
    assignedTo: 'John Doe', priority: 2, areaPath: 'Xbox\\Live Services\\Auth',
    description: 'Tokens expire 5 minutes earlier than expected due to clock skew.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48292',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
  {
    id: 3, adoId: 'ADO-48350', title: 'Build new dashboard layout', state: 'Active', type: 'Task',
    assignedTo: 'Alex Johnson', priority: 2, areaPath: 'Xbox\\Live Services\\Dashboard',
    description: 'Create the new dashboard layout with grid-based widget system.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48350',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
  {
    id: 4, adoId: 'ADO-48300', title: 'SSO integration for all Xbox services', state: 'New', type: 'User Story',
    assignedTo: 'Unassigned', priority: 3, areaPath: 'Xbox\\Live Services\\Auth',
    description: 'As a user, I want to log in once and access all Xbox services.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48300',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
  {
    id: 5, adoId: 'ADO-48370', title: 'Activity feed real-time updates', state: 'New', type: 'User Story',
    assignedTo: 'Unassigned', priority: 3, areaPath: 'Xbox\\Live Services\\Dashboard',
    description: 'As a team lead, I want a real-time activity feed showing recent changes.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48370',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
  {
    id: 6, adoId: 'ADO-48360', title: 'Chart data mismatch in reports', state: 'Resolved', type: 'Bug',
    assignedTo: 'Maria Garcia', priority: 2, areaPath: 'Xbox\\Live Services\\Dashboard',
    description: 'Charts show stale data due to caching layer returning previous day totals.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48360',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
  {
    id: 7, adoId: 'ADO-48400', title: 'Add rate limiting to CLI', state: 'Closed', type: 'Task',
    assignedTo: 'Chris Lee', priority: 1, areaPath: 'Xbox\\Developer Tooling\\CLI',
    description: 'Implement client-side rate limiting with exponential backoff for API calls.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48400',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
  {
    id: 8, adoId: 'ADO-48100', title: 'Migrate to MSAL v3', state: 'Closed', type: 'Task',
    assignedTo: 'Sam Wilson', priority: 1, areaPath: 'Xbox\\Developer Tooling\\CLI',
    description: 'Migrate auth from ADAL to MSAL v3.',
    url: 'https://dev.azure.com/xbox/xb-tasks/_workitems/edit/48100',
    syncedAt: new Date(Date.now() - 600000).toISOString(),
  },
]

const MOCK_PIPELINES: ADOPipeline[] = [
  {
    id: 1, name: 'Deploy xb-services', status: 'completed', result: 'succeeded',
    url: 'https://dev.azure.com/xbox/xb-tasks/_build/results?buildId=1245',
    sourceBranch: 'refs/heads/main', queueTime: new Date(Date.now() - 3600000).toISOString(),
    finishTime: new Date(Date.now() - 3300000).toISOString(),
  },
  {
    id: 2, name: 'Deploy xb-data', status: 'completed', result: 'failed',
    url: 'https://dev.azure.com/xbox/xb-tasks/_build/results?buildId=1244',
    sourceBranch: 'refs/heads/feature/schema-v2', queueTime: new Date(Date.now() - 7200000).toISOString(),
    finishTime: new Date(Date.now() - 6900000).toISOString(),
  },
  {
    id: 3, name: 'CI xb-tasks', status: 'inProgress', result: '',
    url: 'https://dev.azure.com/xbox/xb-tasks/_build/results?buildId=1246',
    sourceBranch: 'refs/heads/feature/ado-view', queueTime: new Date(Date.now() - 120000).toISOString(),
    finishTime: null,
  },
  {
    id: 4, name: 'Deploy xb-telemetry', status: 'completed', result: 'succeeded',
    url: 'https://dev.azure.com/xbox/xb-tasks/_build/results?buildId=1240',
    sourceBranch: 'refs/heads/main', queueTime: new Date(Date.now() - 86400000).toISOString(),
    finishTime: new Date(Date.now() - 86100000).toISOString(),
  },
  {
    id: 5, name: 'CI xb-gateway', status: 'notStarted', result: '',
    url: 'https://dev.azure.com/xbox/xb-tasks/_build/results?buildId=1247',
    sourceBranch: 'refs/heads/feature/rate-limiting', queueTime: new Date(Date.now() - 60000).toISOString(),
    finishTime: null,
  },
]

let useMock = false

export const useADOStore = defineStore('ado', () => {
  const workItems = ref<ADOWorkItem[]>([])
  const pipelines = ref<ADOPipeline[]>([])
  const loading = ref(false)
  const connected = ref(false)
  const syncing = ref(false)
  const lastSyncedAt = ref<string | null>(null)

  async function fetchWorkItems() {
    try {
      if (!useMock) {
        const { ListMyWorkItems } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
        workItems.value = (await ListMyWorkItems()) as ADOWorkItem[]
        connected.value = true
      } else {
        await new Promise(r => setTimeout(r, 150))
        workItems.value = [...MOCK_WORK_ITEMS]
      }
    } catch (e) {
      console.warn('[ADOStore] Wails binding unavailable for ListMyWorkItems, using mock data:', e)
      useMock = true
      workItems.value = [...MOCK_WORK_ITEMS]
    }
  }

  async function syncWorkItems() {
    syncing.value = true
    try {
      if (!useMock) {
        const { SyncWorkItems } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
        await SyncWorkItems()
        connected.value = true
      }
      await fetchWorkItems()
      lastSyncedAt.value = new Date().toISOString()
    } catch (e) {
      console.warn('[ADOStore] SyncWorkItems failed, fetching from cache:', e)
      await fetchCached()
      lastSyncedAt.value = new Date().toISOString()
    } finally {
      syncing.value = false
    }
  }

  async function fetchCached() {
    try {
      if (!useMock) {
        const { GetCachedWorkItems } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
        workItems.value = (await GetCachedWorkItems()) as ADOWorkItem[]
        connected.value = true
      } else {
        await new Promise(r => setTimeout(r, 100))
        workItems.value = [...MOCK_WORK_ITEMS]
      }
    } catch (e) {
      console.warn('[ADOStore] GetCachedWorkItems unavailable, using mock data:', e)
      useMock = true
      workItems.value = [...MOCK_WORK_ITEMS]
    }
  }

  async function fetchPipelines() {
    try {
      if (!useMock) {
        const { ListRecentRuns } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/pipelineservice')
        pipelines.value = (await ListRecentRuns()) as ADOPipeline[]
        connected.value = true
      } else {
        await new Promise(r => setTimeout(r, 150))
        pipelines.value = [...MOCK_PIPELINES]
      }
    } catch (e) {
      console.warn('[ADOStore] Wails binding unavailable for ListRecentRuns, using mock data:', e)
      useMock = true
      pipelines.value = [...MOCK_PIPELINES]
    }
  }

  async function fetchAll() {
    loading.value = true
    try {
      await Promise.all([fetchWorkItems(), fetchPipelines()])
      if (!lastSyncedAt.value) {
        lastSyncedAt.value = new Date().toISOString()
      }
    } finally {
      loading.value = false
    }
  }

  return {
    workItems, pipelines, loading, connected, syncing, lastSyncedAt,
    fetchWorkItems, syncWorkItems, fetchCached, fetchPipelines, fetchAll,
  }
})
