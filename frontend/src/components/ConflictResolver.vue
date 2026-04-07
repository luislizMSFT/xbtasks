<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { AlertTriangle } from 'lucide-vue-next'
import { useSyncStore } from '@/stores/sync'

const syncStore = useSyncStore()

// Track user's choice per field: 'local' or 'remote'
const resolutions = ref<Record<string, string>>({})

// Reset resolutions when conflict changes
watch(
  () => syncStore.currentConflict,
  (conflict) => {
    resolutions.value = {}
    if (conflict) {
      conflict.fields.forEach((f) => {
        resolutions.value[f.field] = '' // not yet picked
      })
    }
  },
)

const allResolved = computed(() => {
  if (!syncStore.currentConflict) return false
  return syncStore.currentConflict.fields.every((f) => resolutions.value[f.field])
})

async function confirmResolve() {
  if (!syncStore.currentConflict || !allResolved.value) return
  await syncStore.resolveConflict(syncStore.currentConflict.taskId, resolutions.value)
}
</script>

<template>
  <Dialog v-model:open="syncStore.showConflictResolver">
    <DialogContent class="max-w-xl">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2">
          <AlertTriangle :size="20" class="text-amber-500" />
          Conflict Detected
        </DialogTitle>
        <DialogDescription>
          Both local and ADO changes were made. Pick which value to keep for each field.
        </DialogDescription>
      </DialogHeader>

      <div v-if="syncStore.currentConflict" class="space-y-4">
        <div class="text-sm text-muted-foreground">
          Task #{{ syncStore.currentConflict.taskId }} &#x2194; ADO
          {{ syncStore.currentConflict.adoId }}
        </div>

        <!-- Per-field resolution -->
        <div
          v-for="field in syncStore.currentConflict.fields"
          :key="field.field"
          class="rounded-md border border-border p-4 space-y-3"
        >
          <div class="text-sm font-semibold capitalize">{{ field.field }}</div>

          <div class="grid grid-cols-2 gap-3">
            <!-- Local value -->
            <button
              @click="resolutions[field.field] = 'local'"
              :class="[
                'p-3 rounded-md border text-left text-sm transition-colors',
                resolutions[field.field] === 'local'
                  ? 'border-green-500 bg-green-500/10 ring-2 ring-green-500/30'
                  : 'border-border hover:border-green-500/50',
              ]"
            >
              <div class="text-xs font-medium text-green-600 mb-1">Keep Local</div>
              <div class="text-foreground">{{ field.local || '(empty)' }}</div>
            </button>

            <!-- Remote value -->
            <button
              @click="resolutions[field.field] = 'remote'"
              :class="[
                'p-3 rounded-md border text-left text-sm transition-colors',
                resolutions[field.field] === 'remote'
                  ? 'border-blue-500 bg-blue-500/10 ring-2 ring-blue-500/30'
                  : 'border-border hover:border-blue-500/50',
              ]"
            >
              <div class="text-xs font-medium text-blue-600 mb-1">Use ADO</div>
              <div class="text-foreground">{{ field.remote || '(empty)' }}</div>
            </button>
          </div>
        </div>

        <!-- Remaining conflicts badge -->
        <div v-if="syncStore.conflictCount > 1" class="text-xs text-muted-foreground">
          {{ syncStore.conflictCount - 1 }} more conflict(s) remaining after this one
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="syncStore.showConflictResolver = false">
          Skip for Now
        </Button>
        <Button @click="confirmResolve" :disabled="!allResolved"> Resolve Conflict </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
