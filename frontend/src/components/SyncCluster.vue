<script setup lang="ts">
import { computed } from 'vue'
import { useADOStore } from '@/stores/ado'
import { useSyncStore } from '@/stores/sync'
import { relativeTime } from '@/lib/date'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { RefreshCw, Loader2 } from 'lucide-vue-next'

const adoStore = useADOStore()
const syncStore = useSyncStore()
const pendingCount = computed(() => syncStore.conflicts.length)
</script>

<template>
  <div class="flex items-center gap-2 text-[11px]">
    <!-- Connection status: dot + label -->
    <span class="inline-flex items-center gap-1.5">
      <span class="w-1.5 h-1.5 rounded-full"
        :class="adoStore.connected ? 'bg-green-500' : 'bg-red-500'" />
      <span :class="[
        !adoStore.connected ? 'text-red-500' :
        syncStore.syncing ? 'text-muted-foreground' : 'text-green-600'
      ]">
        {{ !adoStore.connected ? 'Offline' : syncStore.syncing ? 'Syncing…' : 'Synced' }}
      </span>
    </span>
    <!-- Relative time (only when connected and not syncing) -->
    <span v-if="syncStore.lastSyncedAt && adoStore.connected && !syncStore.syncing"
      class="text-muted-foreground/50 tabular-nums">
      {{ relativeTime(syncStore.lastSyncedAt) }}
    </span>
    <!-- Pending changes badge (only when conflicts > 0) -->
    <Badge v-if="pendingCount > 0" variant="outline"
      class="text-[9px] h-4 px-1.5 gap-1 text-amber-600 border-amber-500/30 bg-amber-500/10">
      <span class="w-1 h-1 rounded-full bg-amber-500" />
      {{ pendingCount }} pending
    </Badge>
    <!-- Refresh button (icon-only, must have aria-label) -->
    <Button variant="ghost" size="sm" class="h-6 w-6 p-0"
      aria-label="Refresh sync"
      @click="syncStore.manualSync()" :disabled="syncStore.syncing">
      <component :is="syncStore.syncing ? Loader2 : RefreshCw"
        :size="12" :class="syncStore.syncing && 'animate-spin'"
        class="text-muted-foreground" />
    </Button>
  </div>
</template>
