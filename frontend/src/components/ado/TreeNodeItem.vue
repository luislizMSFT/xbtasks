<script setup lang="ts">
import { computed } from 'vue'
import { statusBgColor, priorityColor } from '@/lib/styles'
import type { SubtaskItem, GraphNode } from '@/types'

export type { SubtaskItem, GraphNode }

const props = defineProps<{
  node: GraphNode
  depth: number
  expandedNodes: Set<number>
  selectedNodeId: number | null
  highlightedIds: Set<number>
  statusFilter: string
}>()

const emit = defineEmits<{
  'toggle-expand': [id: number]
  'select-node': [id: number]
  'navigate': [id: number]
}>()

function getDescendantCount(node: GraphNode): number {
  let count = node.children.length
  for (const child of node.children) {
    count += getDescendantCount(child)
  }
  return count
}

function nodeMatchesFilter(node: GraphNode, filter: string): boolean {
  if (filter === 'all') return true
  if (node.status === filter) return true
  return node.children.some(c => nodeMatchesFilter(c, filter))
}

const isExpanded = computed(() => props.expandedNodes.has(props.node.id))
const hasChildren = computed(() => props.node.children.length > 0 || (props.node.subtasks != null && props.node.subtasks.length > 0))
const isSelected = computed(() => props.selectedNodeId === props.node.id)
const isHighlighted = computed(() => props.highlightedIds.has(props.node.id))
const isBlocked = computed(() => props.node.status === 'blocked')
const descendantCount = computed(() => getDescendantCount(props.node))
const filteredChildren = computed(() => props.node.children.filter(c => nodeMatchesFilter(c, props.statusFilter)))

const statusLabel= computed(() => {
  switch (props.node.status) {
    case 'blocked': return 'Blocked'
    case 'in_progress': return 'In Progress'
    case 'done': return 'Done'
    case 'todo': return 'To Do'
    default: return props.node.status
  }
})
</script>

<template>
  <div>
    <div
      class="tree-node group relative flex items-center gap-2 rounded-md px-2 py-1.5 transition-colors hover:bg-accent/50"
      :class="{
        'ring-2 ring-primary/50 bg-primary/5': isSelected,
        'bg-yellow-500/5': isHighlighted && !isSelected,
        'border-l-2 border-l-red-500': isBlocked && depth === 0,
      }"
      :style="{ paddingLeft: (depth * 24 + 8) + 'px' }"
      @click="emit('select-node', node.id)"
    >
      <!-- Connector lines for nested items -->
      <div
        v-if="depth > 0"
        class="connector-line absolute top-0 bottom-0"
        :style="{ left: ((depth - 1) * 24 + 19) + 'px', width: '24px' }"
      >
        <div class="absolute top-1/2 left-0 h-px w-4 bg-border"></div>
        <div class="absolute top-0 left-0 h-1/2 w-px bg-border"></div>
      </div>

      <!-- Expand/collapse chevron -->
      <button
        v-if="hasChildren"
        class="flex h-5 w-5 shrink-0 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground"
        @click.stop="emit('toggle-expand', node.id)"
      >
        <svg v-if="isExpanded" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
      </button>
      <div v-else class="w-5 shrink-0"></div>

      <!-- Status dot -->
      <div class="flex h-2.5 w-2.5 shrink-0 rounded-full" :class="statusBgColor(node.status)"></div>

      <!-- Task title -->
      <button
        class="min-w-0 flex-1 truncate text-left text-sm font-medium text-foreground hover:underline"
        :title="node.description"
        @click.stop="emit('navigate', node.id)"
      >
        {{ node.title }}
      </button>

      <!-- Badges -->
      <div class="flex shrink-0 items-center gap-1.5 opacity-80 group-hover:opacity-100">
        <span
          class="rounded px-1.5 py-0.5 text-[10px] font-medium"
          :class="{
            'bg-red-500/15 text-red-600 dark:text-red-400': node.status === 'blocked',
            'bg-blue-500/15 text-blue-600 dark:text-blue-400': node.status === 'in_progress',
            'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400': node.status === 'done',
            'bg-zinc-500/15 text-zinc-500': node.status === 'todo',
          }"
        >{{ statusLabel }}</span>
        <span
          class="text-[10px] font-semibold"
          :class="priorityColor(node.priority)"
        >{{ node.priority }}</span>
        <span
          v-if="descendantCount > 0"
          class="rounded bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground"
          :title="descendantCount + ' dependent task(s)'"
        >{{ descendantCount }} dep{{ descendantCount > 1 ? 's' : '' }}</span>
      </div>
    </div>

    <!-- Subtasks (leaf items) -->
    <div v-if="isExpanded && node.subtasks && node.subtasks.length > 0" class="relative">
      <div
        class="absolute top-0 bottom-2"
        :style="{ left: (depth * 24 + 19) + 'px' }"
      >
        <div class="h-full w-px bg-border"></div>
      </div>
      <div
        v-for="st in node.subtasks"
        :key="'st-' + st.id"
        class="relative flex items-center gap-2 rounded-md px-2 py-1 text-muted-foreground"
        :style="{ paddingLeft: ((depth + 1) * 24 + 8) + 'px' }"
      >
        <div
          class="connector-line absolute top-0 bottom-0"
          :style="{ left: (depth * 24 + 19) + 'px', width: '24px' }"
        >
          <div class="absolute top-1/2 left-0 h-px w-4 bg-border"></div>
          <div class="absolute top-0 left-0 h-1/2 w-px bg-border"></div>
        </div>
        <div class="w-5 shrink-0"></div>
        <!-- Checkbox icon: filled green if done, empty circle if not -->
        <svg v-if="st.done" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="shrink-0 text-emerald-500"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="m9 11 3 3L22 4"/></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="shrink-0 text-zinc-400"><circle cx="12" cy="12" r="10"/></svg>
        <span class="truncate text-xs" :class="st.done ? 'line-through opacity-60' : 'opacity-80'">{{ st.title }}</span>
      </div>
    </div>

    <!-- Children (recursive) -->
    <div v-if="hasChildren && isExpanded && filteredChildren.length > 0" class="relative">
      <div
        class="absolute top-0 bottom-2"
        :style="{ left: (depth * 24 + 19) + 'px' }"
      >
        <div class="h-full w-px bg-border"></div>
      </div>
      <TreeNodeItem
        v-for="child in filteredChildren"
        :key="child.id"
        :node="child"
        :depth="depth + 1"
        :expanded-nodes="expandedNodes"
        :selected-node-id="selectedNodeId"
        :highlighted-ids="highlightedIds"
        :status-filter="statusFilter"
        @toggle-expand="emit('toggle-expand', $event)"
        @select-node="emit('select-node', $event)"
        @navigate="emit('navigate', $event)"
      />
    </div>
  </div>
</template>

<style scoped>
.tree-node {
  min-height: 36px;
}
</style>
