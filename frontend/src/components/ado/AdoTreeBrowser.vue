<script setup lang="ts">
import { ref, computed } from 'vue'
import type { ADOWorkItem } from '@/stores/ado'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { adoTypeColor, adoTypeIcon, adoStateClasses } from '@/lib/styles'
import {
  ChevronRight,
  ChevronDown,
  ExternalLink,
  Link,
  Download,
  CheckCircle2,
} from 'lucide-vue-next'

const props = defineProps<{
  items: ADOWorkItem[]
  getChildren: (id: number) => ADOWorkItem[]
  isLinked: (adoId: string) => boolean
  selectedId?: string
}>()

const emit = defineEmits<{
  'select': [item: ADOWorkItem]
  'import': [item: ADOWorkItem]
  'link': [item: ADOWorkItem]
}>()

const expandedNodes = ref<Set<number>>(new Set())

function toggleExpand(id: number) {
  if (expandedNodes.value.has(id)) {
    expandedNodes.value.delete(id)
  } else {
    expandedNodes.value.add(id)
  }
}

// Use centralized mapping from styles.ts
const typeIcon = adoTypeIcon
</script>

<template>
  <div class="space-y-0.5">
    <template v-for="item in items" :key="item.id">
      <AdoTreeNode
        :item="item"
        :depth="0"
        :expanded-nodes="expandedNodes"
        :get-children="getChildren"
        :is-linked="isLinked"
        :selected-id="selectedId"
        :type-icon="typeIcon"
        :type-color="adoTypeColor"
        :state-classes="adoStateClasses"
        @toggle="toggleExpand"
        @select="(i: ADOWorkItem) => emit('select', i)"
        @import="(i: ADOWorkItem) => emit('import', i)"
        @link="(i: ADOWorkItem) => emit('link', i)"
      />
    </template>
    <p v-if="items.length === 0" class="text-[11px] text-muted-foreground/40 py-4 text-center">
      No work items found
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent, h, type PropType, type Component } from 'vue'

// Recursive tree node component defined inline to avoid separate file
const AdoTreeNode = defineComponent({
  name: 'AdoTreeNode',
  props: {
    item: { type: Object as PropType<ADOWorkItem>, required: true },
    depth: { type: Number, required: true },
    expandedNodes: { type: Object as PropType<Set<number>>, required: true },
    getChildren: { type: Function as PropType<(id: number) => ADOWorkItem[]>, required: true },
    isLinked: { type: Function as PropType<(adoId: string) => boolean>, required: true },
    selectedId: { type: String, default: '' },
    typeIcon: { type: Function as PropType<(type: string) => Component>, required: true },
    typeColor: { type: Function as PropType<(type: string) => string>, required: true },
    stateClasses: { type: Function as PropType<(state: string) => string>, required: true },
  },
  emits: ['toggle', 'select', 'import', 'link'],
  setup(props, { emit }) {
    return () => {
      const item = props.item!
      const children = props.getChildren!(item.id)
      const hasChildren = children.length > 0
      const isExpanded = props.expandedNodes!.has(item.id)
      const linked = props.isLinked!(item.adoId)
      const IconComp = props.typeIcon!(item.type)
      const iconColorClass = props.typeColor!(item.type)
      const stateClass = props.stateClasses!(item.state)
      const indent = props.depth! * 24 + 8

      const rowChildren = [
        // Expand chevron
        hasChildren
          ? h('button', {
              class: 'flex h-5 w-5 shrink-0 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground self-start mt-0.5',
              onClick: (e: Event) => { e.stopPropagation(); emit('toggle', item.id) },
            }, [h(isExpanded ? ChevronDown : ChevronRight, { size: 14 })])
          : h('div', { class: 'w-5 shrink-0' }),

        // Two-line content
        h('div', { class: 'flex-1 min-w-0' }, [
          // Line 1: icon + title + state
          h('div', { class: 'flex items-center gap-1.5' }, [
            h(IconComp, { size: 14, class: `shrink-0 ${iconColorClass}` }),
            h('span', {
              class: 'min-w-0 flex-1 truncate text-sm text-foreground',
              title: item.title,
            }, item.title),
            h(Badge, {
              variant: 'outline',
              class: `text-[10px] h-4 px-1.5 shrink-0 ${stateClass}`,
            }, () => item.state),
            linked
              ? h(CheckCircle2, { size: 12, class: 'shrink-0 text-emerald-500' })
              : null,
          ]),
          // Line 2: metadata + actions (compact)
          h('div', { class: 'flex items-center gap-2 mt-0.5' }, [
            h('span', { class: 'text-[10px] text-muted-foreground/50 tabular-nums' }, `#${item.adoId}`),
            (item.org || item.project) ? h('span', {
              class: 'text-[10px] text-muted-foreground/40 truncate max-w-[8rem]',
            }, `${item.org}/${item.project}`) : null,
            h('div', { class: 'flex-1' }),
            // Action buttons (only on hover)
            !linked ? h('div', {
              class: 'flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity',
            }, [
              h(Button, {
                variant: 'ghost', size: 'sm',
                class: 'h-5 px-1 text-[9px] gap-0.5',
                onClick: (e: Event) => { e.stopPropagation(); emit('import', item) },
              }, () => [h(Download, { size: 9 }), 'Import']),
              h(Button, {
                variant: 'ghost', size: 'sm',
                class: 'h-5 px-1 text-[9px] gap-0.5',
                onClick: (e: Event) => { e.stopPropagation(); emit('link', item) },
              }, () => [h(Link, { size: 9 }), 'Link']),
            ]) : null,
          ]),
        ]),
      ]

      const isSelected = item.adoId === props.selectedId
      const row = h('div', {
        class: `group flex items-start gap-1.5 rounded-md px-2 py-1.5 transition-colors cursor-pointer ${isSelected ? 'bg-primary/10 ring-1 ring-primary/30' : 'hover:bg-accent/50'}`,
        style: { paddingLeft: indent + 'px' },
        onClick: () => emit('select', item),
      }, rowChildren.filter(Boolean))

      // Recursive children
      const childElements = hasChildren && isExpanded
        ? children.map(child =>
            h(AdoTreeNode, {
              key: child.id,
              item: child,
              depth: props.depth! + 1,
              expandedNodes: props.expandedNodes,
              getChildren: props.getChildren,
              isLinked: props.isLinked,
              selectedId: props.selectedId,
              typeIcon: props.typeIcon,
              typeColor: props.typeColor,
              stateClasses: props.stateClasses,
              onToggle: (id: number) => emit('toggle', id),
              onSelect: (i: ADOWorkItem) => emit('select', i),
              onImport: (i: ADOWorkItem) => emit('import', i),
              onLink: (i: ADOWorkItem) => emit('link', i),
            })
          )
        : []

      return h('div', {}, [row, ...childElements])
    }
  },
})
</script>
