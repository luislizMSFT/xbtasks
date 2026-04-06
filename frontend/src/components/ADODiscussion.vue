<script setup lang="ts">
import { ref, watch } from 'vue'
import { relativeTime } from '@/lib/date'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Textarea } from '@/components/ui/textarea'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'
import { Globe } from 'lucide-vue-next'

const props = defineProps<{
  taskId: number
  adoId: string
}>()

interface ADOComment {
  id: number
  text: string
  createdBy: string
  createdDate: string
}

const comments = ref<ADOComment[]>([])
const loading = ref(false)
const replyText = ref('')
const replying = ref(false)

async function fetchComments() {
  loading.value = true
  try {
    const { FetchADOComments } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
    comments.value = (await FetchADOComments(props.taskId)) as ADOComment[]
  } catch {
    comments.value = []
  } finally {
    loading.value = false
  }
}

async function reply() {
  const content = replyText.value.trim()
  if (!content) return
  replying.value = true
  try {
    const { ReplyToADOComment } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
    const c = await ReplyToADOComment(props.taskId, content)
    if (c) comments.value.push(c as ADOComment)
    replyText.value = ''
  } catch (e) {
    console.warn('[ADODiscussion] Failed to reply:', e)
  } finally {
    replying.value = false
  }
}

watch(() => props.taskId, () => fetchComments(), { immediate: true })
</script>

<template>
  <div class="border border-blue-500/20 rounded-lg bg-blue-50/5">
    <!-- Header -->
    <div class="flex items-center gap-2 px-3 py-2 border-b border-blue-500/20">
      <Globe :size="14" class="text-blue-500" />
      <span class="font-semibold text-xs text-blue-500">ADO Discussion</span>
      <Badge v-if="comments.length > 0" variant="secondary" class="h-4 text-[10px] px-1.5 ml-auto">
        {{ comments.length }}
      </Badge>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="py-6">
      <LoadingSpinner size="sm" label="Loading ADO comments..." />
    </div>

    <!-- Empty state -->
    <EmptyState
      v-else-if="comments.length === 0"
      :icon="Globe"
      title="No ADO discussion yet"
      description="Comments from the linked ADO work item will appear here."
      class="py-6"
    />

    <!-- Comments list -->
    <ScrollArea v-else class="max-h-72">
      <div class="px-3 py-2 space-y-3">
        <div
          v-for="c in comments"
          :key="c.id"
          class="pb-3 border-b border-border/50 last:border-0 last:pb-0"
        >
          <div class="flex items-baseline gap-2 mb-0.5">
            <span class="font-medium text-[13px]">{{ c.createdBy }}</span>
            <span class="text-muted-foreground text-[10px]">{{ relativeTime(c.createdDate) }}</span>
          </div>
          <!-- Render HTML content safely -->
          <div class="text-[13px] text-muted-foreground prose prose-sm max-w-none [&_*]:text-[13px] [&_*]:text-muted-foreground" v-html="c.text" />
        </div>
      </div>
    </ScrollArea>

    <!-- Reply input -->
    <div class="px-3 py-2 border-t border-blue-500/20">
      <Textarea
        v-model="replyText"
        placeholder="Reply on ADO..."
        class="text-xs min-h-[50px] resize-none mb-2"
        :rows="2"
      />
      <div class="flex justify-between items-center">
        <span class="text-[10px] text-muted-foreground">This reply will be posted to ADO</span>
        <Button
          size="sm"
          class="h-7 text-xs bg-blue-600 hover:bg-blue-700 text-white gap-1"
          @click="reply"
          :disabled="!replyText.trim() || replying"
        >
          <Globe :size="12" />
          {{ replying ? 'Posting...' : 'Reply on ADO' }}
        </Button>
      </div>
    </div>
  </div>
</template>
