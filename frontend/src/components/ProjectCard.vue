<script setup lang="ts">
import { computed } from 'vue'
import type { Project, ProjectProgress } from '@/stores/projects'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Star } from 'lucide-vue-next'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'

const props = defineProps<{
  project: Project
  isLinked: boolean
  progress: ProjectProgress | null
}>()

const emit = defineEmits<{
  click: []
  pin: [pinned: boolean]
}>()

const localPercent = computed(() => {
  if (!props.progress || !props.progress.localTotal) return 0
  return Math.round((props.progress.localDone / props.progress.localTotal) * 100)
})

const adoPercent = computed(() => {
  if (!props.progress || !props.progress.adoTotal) return 0
  return Math.round((props.progress.adoDone / props.progress.adoTotal) * 100)
})
</script>

<template>
  <Card class="cursor-pointer hover:border-primary/50 transition-colors" @click="emit('click')">
    <CardHeader class="pb-2">
      <div class="flex items-center justify-between">
        <CardTitle class="text-sm font-semibold truncate">{{ project.name }}</CardTitle>
        <div class="flex items-center gap-1">
          <!-- Pin button -->
          <Button
            variant="ghost"
            size="icon"
            class="h-6 w-6"
            @click.stop="emit('pin', !project.isPinned)"
          >
            <Star
              :size="14"
              :class="project.isPinned ? 'fill-amber-400 text-amber-400' : 'text-muted-foreground'"
            />
          </Button>
          <!-- ADO badge -->
          <div
            v-if="isLinked"
            class="w-5 h-5 rounded-full bg-blue-500/10 flex items-center justify-center"
            title="Linked to ADO"
          >
            <AzureDevOpsIcon :size="12" class="text-blue-500" />
          </div>
          <div
            v-else
            class="w-5 h-5 rounded-full border border-dashed border-muted-foreground/30 flex items-center justify-center"
            title="Local only"
          />
        </div>
      </div>
    </CardHeader>
    <CardContent class="pb-3 space-y-2">
      <!-- Task count -->
      <div class="text-xs text-muted-foreground">
        {{ progress?.localTotal ?? 0 }} tasks
      </div>

      <!-- Dual progress bars -->
      <div v-if="progress" class="space-y-1.5">
        <!-- Local task progress -->
        <div class="flex items-center gap-2">
          <span class="text-[10px] text-muted-foreground w-12">Local</span>
          <div class="flex-1 h-1.5 rounded-full bg-muted">
            <div
              class="h-full rounded-full bg-primary transition-all"
              :style="{ width: localPercent + '%' }"
            />
          </div>
          <span class="text-[10px] text-muted-foreground w-8 text-right">{{ localPercent }}%</span>
        </div>
        <!-- ADO progress (only if linked) -->
        <div v-if="isLinked" class="flex items-center gap-2">
          <span class="text-[10px] text-muted-foreground w-12">ADO</span>
          <div class="flex-1 h-1.5 rounded-full bg-muted">
            <div
              class="h-full rounded-full bg-blue-500 transition-all"
              :style="{ width: adoPercent + '%' }"
            />
          </div>
          <span class="text-[10px] text-muted-foreground w-8 text-right">{{ adoPercent }}%</span>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
