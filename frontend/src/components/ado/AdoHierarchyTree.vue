<script setup lang="ts">
import { computed } from 'vue'
import type { ADOWorkItem } from '@/stores/ado'
import { Badge } from '@/components/ui/badge'
import { adoTypeColor, adoTypeIcon, adoStateClasses } from '@/lib/styles'

const props = defineProps<{
  item: ADOWorkItem
  getChildren: (parentId: number) => ADOWorkItem[]
  allItems: ADOWorkItem[]
}>()

const emit = defineEmits<{
  select: [item: ADOWorkItem]
}>()

const parent = computed(() =>
  props.allItems.find(w => w.adoId === String(props.item.parentId))
)

const children = computed(() =>
  props.getChildren(Number(props.item.adoId))
)

const hasHierarchy = computed(() => !!parent.value || children.value.length > 0)
</script>

<template>
  <div v-if="!hasHierarchy" class="text-[11px] text-muted-foreground/50 italic">
    No hierarchy
  </div>
  <div v-else class="space-y-0">
    <!-- Parent -->
    <div
      v-if="parent"
      class="flex items-center gap-1.5 h-7 px-1.5 rounded cursor-pointer hover:bg-muted/50 transition-colors"
      @click="emit('select', parent!)"
    >
      <span class="text-muted-foreground text-[10px] select-none shrink-0">┌</span>
      <component :is="adoTypeIcon(parent.type)" :size="11" :class="adoTypeColor(parent.type)" class="shrink-0" />
      <span class="text-xs text-muted-foreground truncate">{{ parent.title }}</span>
    </div>

    <!-- Connector -->
    <div v-if="parent" class="pl-[7px]">
      <div class="border-l-2 border-muted h-1.5" />
    </div>

    <!-- Selected item -->
    <div class="flex items-center gap-1.5 h-7 px-1.5 rounded bg-primary/10 border border-primary/20">
      <span class="text-muted-foreground text-[10px] select-none shrink-0">{{ parent ? '├─' : '●' }}</span>
      <component :is="adoTypeIcon(item.type)" :size="11" :class="adoTypeColor(item.type)" class="shrink-0" />
      <span class="text-xs font-semibold text-foreground truncate">{{ item.title }}</span>
    </div>

    <!-- Children connector -->
    <div v-if="children.length > 0" class="pl-[7px]">
      <div class="border-l-2 border-muted h-1.5" />
    </div>

    <!-- Children -->
    <div v-if="children.length > 0" class="ml-4">
      <div
        v-for="(child, idx) in children"
        :key="child.adoId"
        class="flex items-center gap-1.5 h-7 px-1.5 rounded cursor-pointer hover:bg-muted/50 transition-colors"
        @click="emit('select', child)"
      >
        <span class="text-muted-foreground text-[10px] select-none shrink-0">
          {{ idx === children.length - 1 ? '└─' : '├─' }}
        </span>
        <component :is="adoTypeIcon(child.type)" :size="11" :class="adoTypeColor(child.type)" class="shrink-0" />
        <span class="text-xs text-foreground truncate flex-1 min-w-0">{{ child.title }}</span>
        <Badge variant="outline" :class="['text-[9px] h-3.5 px-1 shrink-0', adoStateClasses(child.state)]">
          {{ child.state }}
        </Badge>
      </div>
    </div>
  </div>
</template>
