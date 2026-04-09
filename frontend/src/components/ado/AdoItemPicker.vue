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
import { Search, Loader2 } from 'lucide-vue-next'
import { adoTypeColor, adoStateClasses } from '@/lib/styles'

const props = withDefaults(defineProps<{
  open: boolean
  title?: string
}>(), {
  title: 'Select ADO Work Item',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  'selected': [adoId: string, title: string]
}>()

const adoStore = useADOStore()

const searchText = ref('')
const selectedItem = ref<ADOWorkItem | null>(null)

const filteredItems = computed(() => {
  const items = adoStore.workItemTree
  if (!items.length) return []
  if (!searchText.value) return items.slice(0, 30)
  const q = searchText.value.toLowerCase()
  return items
    .filter(i => i.title.toLowerCase().includes(q) || i.adoId.toLowerCase().includes(q))
    .slice(0, 30)
})

function selectItem(item: ADOWorkItem) {
  selectedItem.value = item
}

function confirmSelection() {
  if (!selectedItem.value) return
  emit('selected', selectedItem.value.adoId, selectedItem.value.title)
  emit('update:open', false)
  searchText.value = ''
  selectedItem.value = null
}

watch(() => props.open, async (val) => {
  if (val) {
    if (!adoStore.workItemTree.length) {
      await adoStore.fetchWorkItemTree()
    }
  } else {
    searchText.value = ''
    selectedItem.value = null
  }
})
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2 text-sm">
          <Search :size="16" />
          {{ title }}
        </DialogTitle>
        <DialogDescription class="text-xs">
          Search and select an ADO work item from cached data.
        </DialogDescription>
      </DialogHeader>

      <!-- Search input -->
      <div class="relative">
        <Search :size="14" class="absolute left-2.5 top-2.5 text-muted-foreground" />
        <Input
          v-model="searchText"
          placeholder="Search by title or ADO ID..."
          class="h-9 pl-8 text-sm"
        />
      </div>

      <!-- Loading state -->
      <div v-if="adoStore.loading" class="flex items-center justify-center gap-2 py-6">
        <Loader2 :size="16" class="animate-spin text-muted-foreground" />
        <span class="text-xs text-muted-foreground">Loading work items…</span>
      </div>

      <!-- Results list -->
      <ScrollArea v-else class="max-h-[300px] -mx-1 px-1">
        <div v-if="filteredItems.length > 0" class="space-y-0.5">
          <button
            v-for="item in filteredItems"
            :key="item.adoId"
            class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left transition-colors hover:bg-accent/50"
            :class="selectedItem?.adoId === item.adoId && 'bg-accent ring-1 ring-primary/30'"
            @click="selectItem(item)"
          >
            <span class="text-[10px] font-medium shrink-0" :class="adoTypeColor(item.type)">
              {{ item.type }}
            </span>
            <span class="text-sm text-foreground flex-1 truncate">{{ item.title }}</span>
            <Badge variant="outline" :class="['text-[10px] h-4 px-1.5 shrink-0', adoStateClasses(item.state)]">
              {{ item.state }}
            </Badge>
            <span class="text-[10px] text-muted-foreground/50 shrink-0 tabular-nums">{{ item.adoId }}</span>
          </button>
        </div>
        <p v-else-if="searchText" class="text-[11px] text-muted-foreground/40 py-4 text-center">
          No matching work items found
        </p>
        <p v-else class="text-[11px] text-muted-foreground/40 py-4 text-center">
          No work items available
        </p>
      </ScrollArea>

      <!-- Selected item preview -->
      <div v-if="selectedItem" class="rounded-md border border-primary/20 bg-primary/5 px-3 py-2">
        <div class="flex items-center gap-2 text-xs">
          <span class="font-medium" :class="adoTypeColor(selectedItem.type)">{{ selectedItem.type }}</span>
          <span class="truncate text-foreground">{{ selectedItem.title }}</span>
          <span class="text-muted-foreground/50 ml-auto shrink-0">{{ selectedItem.adoId }}</span>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" size="sm" class="text-xs" @click="emit('update:open', false)">
          Cancel
        </Button>
        <Button
          size="sm"
          class="text-xs gap-1.5"
          :disabled="!selectedItem"
          @click="confirmSelection"
        >
          Select
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
