<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useSyncStore } from '@/stores/sync'

const syncStore = useSyncStore()
</script>

<template>
  <Dialog v-model:open="syncStore.showConfirmDialog">
    <DialogContent class="max-w-lg">
      <DialogHeader>
        <DialogTitle>Push Changes to ADO</DialogTitle>
        <DialogDescription>Review what will change in ADO before confirming.</DialogDescription>
      </DialogHeader>

      <div v-if="syncStore.pendingDiff" class="space-y-4">
        <!-- Task info -->
        <div class="text-sm text-muted-foreground">
          Task #{{ syncStore.pendingDiff.taskId }} → ADO {{ syncStore.pendingDiff.adoId }}
        </div>

        <!-- Field-by-field diff -->
        <div
          v-for="change in syncStore.pendingDiff.changes"
          :key="change.field"
          class="rounded-md border border-border p-3 space-y-2"
        >
          <div class="text-sm font-medium capitalize">{{ change.field }}</div>
          <div class="grid grid-cols-2 gap-2 text-sm">
            <div>
              <div class="text-xs text-muted-foreground mb-1">Local (current)</div>
              <div class="p-2 rounded bg-green-500/10 border border-green-500/20 text-foreground">
                {{ change.local || '(empty)' }}
              </div>
            </div>
            <div>
              <div class="text-xs text-muted-foreground mb-1">ADO (will become)</div>
              <div class="p-2 rounded bg-blue-500/10 border border-blue-500/20 text-foreground">
                {{ change.proposed || '(empty)' }}
              </div>
            </div>
          </div>
          <div v-if="change.remote !== change.local" class="text-xs text-muted-foreground">
            ADO currently: {{ change.remote || '(empty)' }}
          </div>
        </div>

        <!-- No changes -->
        <div
          v-if="!syncStore.pendingDiff.changes.length"
          class="text-center text-muted-foreground py-4"
        >
          No changes to push — local and ADO are in sync.
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="syncStore.cancelPush()">Cancel</Button>
        <Button
          @click="syncStore.confirmPush(syncStore.pendingDiff!.taskId)"
          :disabled="!syncStore.pendingDiff?.changes.length"
        >
          Confirm Push to ADO
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
