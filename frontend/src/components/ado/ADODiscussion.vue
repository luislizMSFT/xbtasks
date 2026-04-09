<script setup lang="ts">
import { ref, watch } from 'vue'
import DOMPurify from 'dompurify'
import { relativeTime } from '@/lib/date'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import EmptyState from '@/components/EmptyState.vue'
import { Globe, RefreshCw } from 'lucide-vue-next'

function sanitize(html: string): string {
  return DOMPurify.sanitize(html, {
    USE_PROFILES: { html: true },
    FORBID_TAGS: ['style', 'script', 'iframe', 'object', 'embed', 'link'],
    FORBID_ATTR: ['style'],
  })
}

const props = defineProps<{
  taskId: number
  adoId: string
}>()

import type { ADOComment } from '@/types'

const comments = ref<ADOComment[]>([])
const loading = ref(false)
const replyText = ref('')
const replying = ref(false)

async function fetchComments() {
  loading.value = true
  try {
    const { fetchADOComments } = await import('@/api/comments')
    comments.value = (await fetchADOComments(props.taskId)) as ADOComment[]
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
    const { replyToADOComment } = await import('@/api/comments')
    const c = await replyToADOComment(props.taskId, content)
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
    <div class="flex items-center gap-2 px-3 py-1.5 border-b border-blue-500/20">
      <Globe :size="13" class="text-blue-500" />
      <span class="font-semibold text-[11px] text-blue-500">ADO Discussion</span>
      <Badge v-if="comments.length > 0" variant="secondary" class="h-4 text-[10px] px-1.5 ml-auto">
        {{ comments.length }}
      </Badge>
      <Button variant="ghost" size="icon" class="h-5 w-5 shrink-0" :class="loading && 'animate-spin'" @click="fetchComments" title="Refresh">
        <RefreshCw :size="11" class="text-blue-500" />
      </Button>
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
      description="Comments from the linked work item appear here."
      class="py-4"
    />

    <!-- Comments list (bubble style, left-aligned for ADO authors) -->
    <div v-else class="px-3 py-2 flex flex-col gap-2 max-h-72 overflow-y-auto">
      <div
        v-for="c in comments"
        :key="c.id"
        class="rounded-lg px-3 py-2 max-w-[92%] self-start bg-blue-500/5 border border-blue-500/10"
      >
        <div class="flex items-center gap-1.5 mb-1">
          <span class="font-medium text-[11px] text-blue-400">{{ c.createdBy }}</span>
          <span class="text-[9px] text-muted-foreground/50 ml-auto tabular-nums">{{ relativeTime(c.createdDate) }}</span>
        </div>
        <div class="text-[13px] text-muted-foreground leading-relaxed prose prose-sm max-w-none [&_*]:text-[13px] [&_*]:text-muted-foreground" v-html="sanitize(c.text)" />
      </div>
    </div>

    <!-- Reply input -->
    <div class="px-3 py-2 border-t border-blue-500/20">
      <Textarea
        v-model="replyText"
        placeholder="Reply on ADO..."
        class="text-xs min-h-[40px] resize-none mb-1.5"
        :rows="2"
      />
      <div class="flex justify-between items-center">
        <span class="text-[9px] text-muted-foreground/50">Posted to ADO</span>
        <Button
          size="sm"
          class="h-6 text-[10px] bg-blue-600 hover:bg-blue-700 text-white gap-1"
          @click="reply"
          :disabled="!replyText.trim() || replying"
        >
          <Globe :size="10" />
          {{ replying ? 'Posting...' : 'Reply on ADO' }}
        </Button>
      </div>
    </div>
  </div>
</template>
