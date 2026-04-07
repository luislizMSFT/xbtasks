<script setup lang="ts">
import { computed } from 'vue'
import type { AcceptableValue } from 'reka-ui'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  ArrowUpDown,
  LayoutGrid,
} from 'lucide-vue-next'

const props = defineProps<{
  filterStatus: string
  filterPriority: string
  filterProject: number | null
  filterDueDate: string
  filterAdoLink: string
  sortBy: string
  groupBy: string | null
  projects: Array<{ id: number; name: string }>
  syncing: boolean
}>()

const emit = defineEmits<{
  'update:filterStatus': [value: string]
  'update:filterPriority': [value: string]
  'update:filterProject': [value: number | null]
  'update:filterDueDate': [value: string]
  'update:filterAdoLink': [value: string]
  'update:sortBy': [value: string]
  'update:groupBy': [value: string | null]
  'sync': []
}>()

// Count active (non-default) filters
const activeFilterCount = computed(() => {
  let count = 0
  if (props.filterStatus !== 'all') count++
  if (props.filterPriority !== 'all') count++
  if (props.filterProject !== null) count++
  if (props.filterDueDate !== 'all') count++
  if (props.filterAdoLink !== 'all') count++
  return count
})

function handlePriorityChange(value: AcceptableValue) {
  emit('update:filterPriority', String(value))
}

function handleProjectChange(value: AcceptableValue) {
  const v = String(value)
  emit('update:filterProject', v === 'all' ? null : Number(v))
}

function handleDueDateChange(value: AcceptableValue) {
  emit('update:filterDueDate', String(value))
}

function handleAdoLinkChange(value: AcceptableValue) {
  emit('update:filterAdoLink', String(value))
}

function handleSortByChange(value: AcceptableValue) {
  emit('update:sortBy', String(value))
}

function handleGroupByChange(value: AcceptableValue) {
  const v = String(value)
  emit('update:groupBy', v === 'none' ? null : v)
}
</script>

<template>
  <div class="px-4 py-1.5 border-b border-border">
    <div class="flex items-center gap-1.5">
      <span class="text-[10px] text-muted-foreground shrink-0">Priority</span>
      <Select
        :model-value="filterPriority"
        @update:model-value="handlePriorityChange"
      >
        <SelectTrigger size="sm" class="h-6 text-[10px] gap-0.5 w-[52px] px-1.5">
          <SelectValue placeholder="All" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All</SelectItem>
          <SelectItem value="P0">P0</SelectItem>
          <SelectItem value="P1">P1</SelectItem>
          <SelectItem value="P2">P2</SelectItem>
          <SelectItem value="P3">P3</SelectItem>
        </SelectContent>
      </Select>

      <span class="text-[10px] text-muted-foreground shrink-0">Project</span>
      <Select
        :model-value="filterProject !== null ? String(filterProject) : 'all'"
        @update:model-value="handleProjectChange"
      >
        <SelectTrigger size="sm" class="h-6 text-[10px] gap-0.5 w-[72px] px-1.5">
          <SelectValue placeholder="All" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All</SelectItem>
          <SelectItem
            v-for="proj in projects"
            :key="proj.id"
            :value="String(proj.id)"
          >
            {{ proj.name }}
          </SelectItem>
        </SelectContent>
      </Select>

      <span class="text-[10px] text-muted-foreground shrink-0">Due</span>
      <Select
        :model-value="filterDueDate"
        @update:model-value="handleDueDateChange"
      >
        <SelectTrigger size="sm" class="h-6 text-[10px] gap-0.5 w-[60px] px-1.5">
          <SelectValue placeholder="All" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All</SelectItem>
          <SelectItem value="overdue">Overdue</SelectItem>
          <SelectItem value="today">Today</SelectItem>
          <SelectItem value="week">Week</SelectItem>
          <SelectItem value="none">No Date</SelectItem>
        </SelectContent>
      </Select>

      <span class="text-[10px] text-muted-foreground shrink-0">Scope</span>
      <Select
        :model-value="filterAdoLink"
        @update:model-value="handleAdoLinkChange"
      >
        <SelectTrigger size="sm" class="h-6 text-[10px] gap-0.5 w-[60px] px-1.5">
          <SelectValue placeholder="All" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All</SelectItem>
          <SelectItem value="linked">Public</SelectItem>
          <SelectItem value="personal">Private</SelectItem>
        </SelectContent>
      </Select>

      <Badge
        v-if="activeFilterCount > 0"
        variant="secondary"
        class="text-[9px] px-1 py-0 h-4 shrink-0"
      >
        {{ activeFilterCount }}
      </Badge>

      <div class="flex-1" />

      <Select
        :model-value="sortBy"
        @update:model-value="handleSortByChange"
      >
        <SelectTrigger size="sm" class="h-6 text-[10px] gap-0.5 border-none shadow-none px-1">
          <ArrowUpDown :size="10" class="shrink-0" />
          <SelectValue placeholder="Sort" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="priority">Priority</SelectItem>
          <SelectItem value="dueDate">Due Date</SelectItem>
          <SelectItem value="title">Title</SelectItem>
          <SelectItem value="status">Status</SelectItem>
        </SelectContent>
      </Select>

      <Select
        :model-value="groupBy ?? 'none'"
        @update:model-value="handleGroupByChange"
      >
        <SelectTrigger size="sm" class="h-6 text-[10px] gap-0.5 border-none shadow-none px-1">
          <LayoutGrid :size="10" class="shrink-0" />
          <SelectValue placeholder="Group" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="none">None</SelectItem>
          <SelectItem value="status">Status</SelectItem>
          <SelectItem value="priority">Priority</SelectItem>
          <SelectItem value="project">Project</SelectItem>
        </SelectContent>
      </Select>
    </div>
  </div>
</template>
