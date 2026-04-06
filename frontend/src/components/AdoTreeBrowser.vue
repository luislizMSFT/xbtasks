<script setup lang="ts">
import { ref, computed } from 'vue'
import type { ADOWorkItem } from '@/stores/ado'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Bug,
  CheckSquare,
  BookOpen,
  Star,
  Mountain,
  Circle,
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
    case 'Removed': return 'bg-red-500/15 text-red-700 dark:text-red-400 border-red-500/25'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}
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
        :type-icon="typeIcon"
        :type-color="typeColor"
        :state-classes="stateClasses"
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
              class: 'flex h-5 w-5 shrink-0 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground',
              onClick: (e: Event) => { e.stopPropagation(); emit('toggle', item.id) },
            }, [h(isExpanded ? ChevronDown : ChevronRight, { size: 14 })])
          : h('div', { class: 'w-5 shrink-0' }),

        // Type icon
        h(IconComp, { size: 14, class: `shrink-0 ${iconColorClass}` }),

        // Title
        h('span', {
          class: 'min-w-0 flex-1 truncate text-sm text-foreground',
          title: item.title,
        }, item.title),

        // ADO ID
        h('span', { class: 'text-[10px] text-muted-foreground/50 shrink-0 tabular-nums' }, `#${item.adoId}`),

        // State badge
        h(Badge, {
          variant: 'outline',
          class: `text-[10px] h-4 px-1.5 shrink-0 ${stateClass}`,
        }, () => item.state),

        // Org/project label
        (item.org || item.project) ? h('span', {
          class: 'text-[10px] text-muted-foreground/40 shrink-0 truncate max-w-[8rem]',
        }, `${item.org}/${item.project}`) : null,

        // Linked indicator or action buttons
        linked
          ? h(CheckCircle2, { size: 14, class: 'shrink-0 text-emerald-500', title: 'Linked to local task' })
          : h('div', {
              class: 'shrink-0 flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity',
            }, [
              h(Button, {
                variant: 'ghost',
                size: 'sm',
                class: 'h-6 px-1.5 text-[10px] gap-1',
                title: 'Import as local task',
                onClick: (e: Event) => { e.stopPropagation(); emit('import', item) },
              }, () => [h(Download, { size: 10 }), 'Import']),
              h(Button, {
                variant: 'ghost',
                size: 'sm',
                class: 'h-6 px-1.5 text-[10px] gap-1',
                title: 'Link to existing local task',
                onClick: (e: Event) => { e.stopPropagation(); emit('link', item) },
              }, () => [h(Link, { size: 10 }), 'Link']),
            ]),

        // Open in ADO
        h(Button, {
          variant: 'ghost',
          size: 'sm',
          class: 'h-6 w-6 p-0 shrink-0 opacity-0 group-hover:opacity-60 hover:!opacity-100',
          title: 'Open in ADO',
          onClick: async (e: Event) => {
            e.stopPropagation()
            if (!item.url) return
            try {
              const { OpenURL } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/browserservice')
              await OpenURL(item.url)
            } catch { window.open(item.url, '_blank') }
          },
        }, () => h(ExternalLink, { size: 10 })),
      ]

      const row = h('div', {
        class: 'group flex items-center gap-2 rounded-md px-2 py-1.5 transition-colors hover:bg-accent/50 cursor-pointer',
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
