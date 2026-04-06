<script setup lang="ts">
import { computed } from 'vue'
import type { AcceptableValue } from 'reka-ui'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  RefreshCw,
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

const statusChips = [
  { id: 'all', label: 'All' },
  { id: 'todo', label: 'Todo' },
  { id: 'in_progress', label: 'In Progress' },
  { id: 'in_review', label: 'In Review' },
  { id: 'done', label: 'Done' },
  { id: 'blocked', label: 'Blocked' },
]

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
  <div class="flex items-center gap-2 px-4 py-2 border-b border-border overflow-x-auto">
    <!-- Status chips -->
    <div class="flex items-center gap-1 shrink-0">
      <Badge
        v-for="chip in statusChips"
        :key="chip.id"
        :variant="filterStatus === chip.id ? 'default' : 'outline'"
        class="cursor-pointer text-[11px] px-2 py-0.5 select-none"
        @click="emit('update:filterStatus', chip.id)"
      >
        {{ chip.label }}
      </Badge>
    </div>

    <Separator orientation="vertical" class="h-5 mx-1" />

    <!-- Priority dropdown -->
    <Select
      :model-value="filterPriority"
      @update:model-value="handlePriorityChange"
    >
      <SelectTrigger size="sm" class="h-7 text-[11px] gap-1 min-w-[70px]">
        <SelectValue placeholder="Priority" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">All Priority</SelectItem>
        <SelectItem value="P0">P0</SelectItem>
        <SelectItem value="P1">P1</SelectItem>
        <SelectItem value="P2">P2</SelectItem>
        <SelectItem value="P3">P3</SelectItem>
      </SelectContent>
    </Select>

    <!-- Project dropdown -->
    <Select
      :model-value="filterProject !== null ? String(filterProject) : 'all'"
      @update:model-value="handleProjectChange"
    >
      <SelectTrigger size="sm" class="h-7 text-[11px] gap-1 min-w-[80px]">
        <SelectValue placeholder="Project" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">All Projects</SelectItem>
        <SelectItem
          v-for="proj in projects"
          :key="proj.id"
          :value="String(proj.id)"
        >
          {{ proj.name }}
        </SelectItem>
      </SelectContent>
    </Select>

    <!-- Due date dropdown -->
    <Select
      :model-value="filterDueDate"
      @update:model-value="handleDueDateChange"
    >
      <SelectTrigger size="sm" class="h-7 text-[11px] gap-1 min-w-[80px]">
        <SelectValue placeholder="Due Date" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">All Dates</SelectItem>
        <SelectItem value="overdue">Overdue</SelectItem>
        <SelectItem value="today">Today</SelectItem>
        <SelectItem value="week">This Week</SelectItem>
        <SelectItem value="none">No Date</SelectItem>
      </SelectContent>
    </Select>

    <!-- ADO link filter -->
    <Select
      :model-value="filterAdoLink"
      @update:model-value="handleAdoLinkChange"
    >
      <SelectTrigger size="sm" class="h-7 text-[11px] gap-1 min-w-[80px]">
        <SelectValue placeholder="ADO" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">All Tasks</SelectItem>
        <SelectItem value="linked">Linked (Public)</SelectItem>
        <SelectItem value="personal">Personal</SelectItem>
      </SelectContent>
    </Select>

    <!-- Active filter count -->
    <Badge
      v-if="activeFilterCount > 0"
      variant="secondary"
      class="text-[10px] px-1.5 py-0 shrink-0"
    >
      {{ activeFilterCount }} filter{{ activeFilterCount > 1 ? 's' : '' }}
    </Badge>

    <div class="flex-1" />

    <!-- Right side controls -->
    <div class="flex items-center gap-1 shrink-0">
      <!-- Sync button -->
      <Button
        variant="ghost"
        size="icon"
        class="h-7 w-7"
        title="Sync with ADO"
        @click="emit('sync')"
      >
        <RefreshCw :size="14" :class="{ 'animate-spin': syncing }" />
      </Button>

      <Separator orientation="vertical" class="h-5 mx-0.5" />

      <!-- Sort dropdown -->
      <Select
        :model-value="sortBy"
        @update:model-value="handleSortByChange"
      >
        <SelectTrigger size="sm" class="h-7 text-[11px] gap-1 border-none shadow-none">
          <ArrowUpDown :size="12" class="shrink-0" />
          <SelectValue placeholder="Sort" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="priority">Priority</SelectItem>
          <SelectItem value="dueDate">Due Date</SelectItem>
          <SelectItem value="title">Title</SelectItem>
          <SelectItem value="status">Status</SelectItem>
        </SelectContent>
      </Select>

      <!-- Group-by dropdown -->
      <Select
        :model-value="groupBy ?? 'none'"
        @update:model-value="handleGroupByChange"
      >
        <SelectTrigger size="sm" class="h-7 text-[11px] gap-1 border-none shadow-none">
          <LayoutGrid :size="12" class="shrink-0" />
          <SelectValue placeholder="Group" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="none">No Grouping</SelectItem>
          <SelectItem value="status">By Status</SelectItem>
          <SelectItem value="priority">By Priority</SelectItem>
          <SelectItem value="project">By Project</SelectItem>
        </SelectContent>
      </Select>
    </div>
  </div>
</template>
