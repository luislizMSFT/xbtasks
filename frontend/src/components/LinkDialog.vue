<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useADOStore } from '@/stores/ado'
import type { ADOWorkItem } from '@/stores/ado'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Bug, CheckSquare, BookOpen, Star, Mountain, Circle,
  Search, Link, Loader2, AlertCircle,
} from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  taskId: number
  taskTitle: string
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  'linked': [adoId: string]
}>()

const adoStore = useADOStore()

const searchText = ref('')
const selectedItem = ref<ADOWorkItem | null>(null)
const linking = ref(false)
const linkError = ref('')

const filteredItems = computed(() => {
  if (!searchText.value) return []
  const q = searchText.value.toLowerCase()
  // Search in both cached workItems and tree items
  const allItems = [...adoStore.workItems, ...adoStore.workItemTree]
  const seen = new Set<string>()
  return allItems.filter(i => {
    if (seen.has(i.adoId)) return false
    seen.add(i.adoId)
    return i.adoId.includes(q) || i.title.toLowerCase().includes(q)
  }).slice(0, 20)
})

function typeIcon(type: string) {
  switch (type.toLowerCase()) {
    case 'bug': return Bug
    case 'task': return CheckSquare
    case 'user story': return BookOpen
    case 'feature': return Star
    case 'epic': return Mountain
    default: return Circle
  }
}

function typeColor(type: string) {
  switch (type.toLowerCase()) {
    case 'bug': return 'text-red-500'
    case 'task': return 'text-blue-500'
    case 'user story': return 'text-purple-500'
    case 'feature': return 'text-green-500'
    case 'epic': return 'text-orange-500'
    default: return 'text-muted-foreground'
  }
}

function stateClasses(state: string) {
  switch (state) {
    case 'Active': return 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/25'
    case 'New': return 'bg-muted text-muted-foreground border-border'
    case 'Resolved': return 'bg-yellow-500/15 text-yellow-700 dark:text-yellow-400 border-yellow-500/25'
    case 'Closed': return 'bg-green-500/15 text-green-700 dark:text-green-400 border-green-500/25'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}

function selectItem(item: ADOWorkItem) {
  selectedItem.value = item
}

async function confirmLink() {
  if (!selectedItem.value) return
  linking.value = true
  linkError.value = ''
  try {
    const { LinkTask } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
    await LinkTask(props.taskId, selectedItem.value.adoId)
    await adoStore.fetchLinkedAdoIds()
    emit('linked', selectedItem.value.adoId)
    emit('update:open', false)
    // Reset
    searchText.value = ''
    selectedItem.value = null
  } catch (e: any) {
    linkError.value = e?.message || 'Failed to link task'
  } finally {
    linking.value = false
  }
}

// Reset on close
watch(() => props.open, (val) => {
  if (!val) {
    searchText.value = ''
    selectedItem.value = null
    linkError.value = ''
  }
})
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2 text-sm">
          <Link :size="16" />
          Link to ADO Work Item
        </DialogTitle>
        <DialogDescription class="text-xs">
          Connect "{{ taskTitle }}" to an existing Azure DevOps work item.
        </DialogDescription>
      </DialogHeader>

      <!-- Search input -->
      <div class="relative">
        <Search :size="14" class="absolute left-2.5 top-2.5 text-muted-foreground" />
        <Input
          v-model="searchText"
          placeholder="Search by ADO ID or title..."
          class="h-9 pl-8 text-sm"
        />
      </div>

      <!-- Results list -->
      <ScrollArea class="max-h-[240px] -mx-1 px-1">
        <div v-if="filteredItems.length > 0" class="space-y-0.5">
          <button
            v-for="item in filteredItems"
            :key="item.adoId"
            class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left transition-colors hover:bg-accent/50"
            :class="selectedItem?.adoId === item.adoId && 'bg-accent ring-1 ring-primary/30'"
            @click="selectItem(item)"
          >
            <component :is="typeIcon(item.type)" :size="14" :class="['shrink-0', typeColor(item.type)]" />
            <span class="text-[10px] text-muted-foreground/50 shrink-0 tabular-nums">#{{ item.adoId }}</span>
            <span class="text-sm text-foreground flex-1 truncate">{{ item.title }}</span>
            <Badge variant="outline" :class="['text-[10px] h-4 px-1.5 shrink-0', stateClasses(item.state)]">
              {{ item.state }}
            </Badge>
          </button>
        </div>
        <p v-else-if="searchText" class="text-[11px] text-muted-foreground/40 py-4 text-center">
          No matching work items found
        </p>
        <p v-else class="text-[11px] text-muted-foreground/40 py-4 text-center">
          Type to search ADO work items
        </p>
      </ScrollArea>

      <!-- Error -->
      <div v-if="linkError" class="flex items-center gap-2 text-xs text-destructive">
        <AlertCircle :size="12" />
        {{ linkError }}
      </div>

      <!-- Selected item preview -->
      <div v-if="selectedItem" class="rounded-md border border-primary/20 bg-primary/5 px-3 py-2">
        <div class="flex items-center gap-2 text-xs">
          <component :is="typeIcon(selectedItem.type)" :size="12" :class="typeColor(selectedItem.type)" />
          <span class="font-medium">#{{ selectedItem.adoId }}</span>
          <span class="truncate text-foreground">{{ selectedItem.title }}</span>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" size="sm" class="text-xs" @click="emit('update:open', false)">
          Cancel
        </Button>
        <Button
          size="sm"
          class="text-xs gap-1.5"
          :disabled="!selectedItem || linking"
          @click="confirmLink"
        >
          <Loader2 v-if="linking" :size="12" class="animate-spin" />
          <Link v-else :size="12" />
          Link
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
