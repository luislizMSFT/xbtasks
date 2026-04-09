import { ref, onMounted } from 'vue'
import { Events } from '@wailsio/runtime'
import { getAllADOMeta, type AdoMeta } from '@/api/adometa'

// Module-level shared state — all consumers share the same cache
const metaCache = ref<Map<number, AdoMeta>>(new Map())
const loading = ref(false)
let initialized = false
let eventsBound = false

/**
 * Composable for batch-loaded ADO metadata access.
 * Loads all ADO type/state data in a single call on first mount.
 * Refreshes automatically when `sync:completed` Wails event fires.
 *
 * Usage:
 *   const { getAdoMeta } = useAdoMeta()
 *   const meta = getAdoMeta(task.id) // { type: 'Bug', state: 'Active' } | null
 */
export function useAdoMeta() {
  async function refresh() {
    loading.value = true
    try {
      const data = await getAllADOMeta()
      const map = new Map<number, AdoMeta>()
      for (const [k, v] of Object.entries(data)) {
        map.set(Number(k), v)
      }
      metaCache.value = map
    } catch (err) {
      console.error('[useAdoMeta] Failed to load ADO metadata:', err)
    } finally {
      loading.value = false
    }
  }

  function getAdoMeta(taskId: number): AdoMeta | null {
    return metaCache.value.get(taskId) ?? null
  }

  onMounted(() => {
    // Load cache on first component mount
    if (!initialized) {
      initialized = true
      refresh()
    }

    // Subscribe to sync:completed event (once globally)
    if (!eventsBound) {
      eventsBound = true
      Events.On('sync:completed', () => {
        console.log('[useAdoMeta] Sync completed, refreshing ADO metadata cache')
        refresh()
      })
    }
  })

  return {
    getAdoMeta,
    refresh,
    loading,
    metaCache,
  }
}
