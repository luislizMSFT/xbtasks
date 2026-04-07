import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { toast } from 'vue-sonner'
import { useTaskStore } from './tasks'
import type { FieldDiff, SyncDiff, Conflict } from '@/types'
import * as syncApi from '@/api/sync'

export type { FieldDiff, SyncDiff, Conflict }

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
            const diffs = await syncApi.manualSync() as SyncDiff[]
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
            const diff = await syncApi.generateOutboundDiff(taskId) as SyncDiff
            pendingDiff.value = diff
            if (!diff.changes?.length) {
                error.value = 'No changes to push — local and ADO are already in sync.'
                return
            }
            showConfirmDialog.value = true
        } catch (e: any) {
            error.value = e?.message || 'Failed to generate diff'
        }
    }

    async function confirmPush(taskId: number) {
        syncing.value = true
        error.value = null
        try {
            await syncApi.pushChanges(taskId)
            const taskStore = useTaskStore()
            await taskStore.fetchTasks()
            pendingDiff.value = null
            showConfirmDialog.value = false
            toast.success('Pushed to ADO', { duration: 3000 })
        } catch (e: any) {
            error.value = e?.message || 'Push failed'
            toast.error('Push to ADO failed', { duration: 5000 })
        } finally {
            syncing.value = false
        }
    }

    async function resolveConflict(taskId: number, resolutions: Record<string, string>) {
        try {
            await syncApi.resolveConflict(taskId, resolutions)
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
