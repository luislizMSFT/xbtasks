import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface FieldDiff {
    field: string      // 'title', 'status', 'description'
    local: string      // current local value
    remote: string     // current ADO value
    proposed: string   // what will be written
}

export interface SyncDiff {
    taskId: number
    adoId: string
    changes: FieldDiff[]
    direction: string  // 'push', 'pull', or 'conflict'
}

export interface Conflict {
    taskId: number
    adoId: string
    fields: FieldDiff[]
}

export const useSyncStore = defineStore('sync', () => {
    const syncing = ref(false)
    const lastSyncedAt = ref<string | null>(null)
    const conflicts = ref<Conflict[]>([])
    const pendingDiff = ref<SyncDiff | null>(null)
    const error = ref<string | null>(null)
    const showConfirmDialog = ref(false)
    const showConflictResolver = ref(false)
    const currentConflict = ref<Conflict | null>(null)

    const hasConflicts = computed(() => conflicts.value.length > 0)
    const conflictCount = computed(() => conflicts.value.length)

    async function manualSync() {
        syncing.value = true
        error.value = null
        try {
            const { ManualSync } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
            const diffs = await ManualSync() as SyncDiff[]
            // Separate conflicts (both sides changed) from clean pulls
            conflicts.value = diffs
                .filter(d => d.direction === 'conflict')
                .map(d => ({ taskId: d.taskId, adoId: d.adoId, fields: d.changes }))
            lastSyncedAt.value = new Date().toISOString()
            if (conflicts.value.length > 0) {
                showConflictResolver.value = true
                currentConflict.value = conflicts.value[0]
            }
        } catch (e: any) {
            error.value = e?.message || 'Sync failed'
        } finally {
            syncing.value = false
        }
    }

    async function generateOutboundDiff(taskId: number) {
        try {
            const { GenerateOutboundDiff } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
            pendingDiff.value = await GenerateOutboundDiff(taskId) as SyncDiff
            showConfirmDialog.value = true
        } catch (e: any) {
            error.value = e?.message || 'Failed to generate diff'
        }
    }

    async function confirmPush(taskId: number) {
        try {
            const { PushChanges } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
            await PushChanges(taskId)
            pendingDiff.value = null
            showConfirmDialog.value = false
        } catch (e: any) {
            error.value = e?.message || 'Push failed'
        }
    }

    async function resolveConflict(taskId: number, resolutions: Record<string, string>) {
        try {
            const { ResolveConflict } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
            await ResolveConflict(taskId, resolutions)
            conflicts.value = conflicts.value.filter(c => c.taskId !== taskId)
            if (conflicts.value.length > 0) {
                currentConflict.value = conflicts.value[0]
            } else {
                showConflictResolver.value = false
                currentConflict.value = null
            }
        } catch (e: any) {
            error.value = e?.message || 'Failed to resolve conflict'
        }
    }

    function cancelPush() {
        pendingDiff.value = null
        showConfirmDialog.value = false
    }

    return {
        syncing, lastSyncedAt, conflicts, pendingDiff, error,
        showConfirmDialog, showConflictResolver, currentConflict,
        hasConflicts, conflictCount,
        manualSync, generateOutboundDiff, confirmPush, resolveConflict, cancelPush,
    }
})
