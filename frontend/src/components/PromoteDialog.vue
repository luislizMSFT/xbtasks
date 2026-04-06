<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import {
  ArrowUpCircle, Loader2, AlertCircle, AlertTriangle, CheckSquare, Bug, BookOpen,
} from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  taskId: number
  taskTitle: string
  taskDescription: string
  taskStatus: string
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  'promoted': [adoId: string]
}>()

const wiType = ref('Task')
const editableTitle = ref('')
const promoting = ref(false)
const promoteError = ref('')

const WI_TYPES = [
  { value: 'Task', label: 'Task', icon: CheckSquare, color: 'text-blue-500' },
  { value: 'Bug', label: 'Bug', icon: Bug, color: 'text-red-500' },
  { value: 'User Story', label: 'User Story', icon: BookOpen, color: 'text-purple-500' },
]

async function confirmPromote() {
  promoting.value = true
  promoteError.value = ''
  try {
    const { PromoteTask } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
    const link = await PromoteTask(props.taskId, wiType.value)
    // link contains adoId from the TaskADOLink
    const adoId = (link as any)?.adoId || (link as any)?.ado_id || ''
    emit('promoted', adoId)
    emit('update:open', false)
  } catch (e: any) {
    promoteError.value = e?.message || 'Failed to promote task to ADO'
  } finally {
    promoting.value = false
  }
}

// Sync title on open
watch(() => props.open, (val) => {
  if (val) {
    editableTitle.value = props.taskTitle
    wiType.value = 'Task'
    promoteError.value = ''
  }
})
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2 text-sm">
          <ArrowUpCircle :size="16" />
          Promote to ADO Work Item
        </DialogTitle>
        <DialogDescription class="text-xs">
          Create a new Azure DevOps work item from this local task.
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <!-- Preview -->
        <div class="space-y-3">
          <!-- Title (editable) -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted-foreground">Title</label>
            <Input
              v-model="editableTitle"
              class="h-8 text-sm"
            />
          </div>

          <!-- Description preview -->
          <div v-if="taskDescription" class="space-y-1">
            <label class="text-xs font-medium text-muted-foreground">Description</label>
            <div class="rounded-md border bg-muted/30 px-3 py-2 text-xs text-muted-foreground max-h-[80px] overflow-y-auto">
              {{ taskDescription.slice(0, 200) }}{{ taskDescription.length > 200 ? '...' : '' }}
            </div>
          </div>

          <!-- Work item type -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted-foreground">Work Item Type</label>
            <Select v-model="wiType">
              <SelectTrigger class="h-8 text-xs">
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="t in WI_TYPES"
                  :key="t.value"
                  :value="t.value"
                  class="text-xs"
                >
                  <div class="flex items-center gap-2">
                    <component :is="t.icon" :size="12" :class="t.color" />
                    {{ t.label }}
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- Current status -->
          <div class="flex items-center gap-2">
            <span class="text-xs text-muted-foreground">Current status:</span>
            <Badge variant="outline" class="text-[10px] h-4 px-1.5">{{ taskStatus }}</Badge>
          </div>
        </div>

        <!-- Warning -->
        <div class="flex items-start gap-2 rounded-md border border-amber-500/20 bg-amber-500/5 px-3 py-2">
          <AlertTriangle :size="14" class="text-amber-500 shrink-0 mt-0.5" />
          <p class="text-[11px] text-amber-700 dark:text-amber-400">
            Subtasks, personal priority, and local notes will NOT be synced to ADO.
            Only title, description, and status will be pushed.
          </p>
        </div>

        <!-- Error -->
        <div v-if="promoteError" class="flex items-center gap-2 text-xs text-destructive">
          <AlertCircle :size="12" />
          {{ promoteError }}
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" size="sm" class="text-xs" @click="emit('update:open', false)">
          Cancel
        </Button>
        <Button
          size="sm"
          class="text-xs gap-1.5"
          :disabled="!editableTitle.trim() || promoting"
          @click="confirmPromote"
        >
          <Loader2 v-if="promoting" :size="12" class="animate-spin" />
          <ArrowUpCircle v-else :size="12" />
          Create in ADO
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
