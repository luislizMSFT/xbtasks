<script setup lang="ts">
import { ref, watch } from 'vue'
import DOMPurify from 'dompurify'
import { relativeTime } from '@/lib/date'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import { Globe } from 'lucide-vue-next'

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
  <div class="space-y-2">
    <!-- Reply input (on top) -->
    <div class="space-y-1.5">
      <Textarea
        v-model="replyText"
        placeholder="Reply on ADO..."
        class="text-xs min-h-[40px] resize-none"
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
          {{ replying ? 'Posting...' : 'Reply' }}
        </Button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="py-4">
      <LoadingSpinner size="sm" label="Loading ADO comments..." />
    </div>

    <!-- Comments list (bubble style, left-aligned for ADO authors) -->
    <div v-else-if="comments.length > 0" class="flex flex-col gap-2 max-h-60 overflow-y-auto">
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

    <!-- Empty state -->
    <div v-else class="text-[11px] text-muted-foreground/40 italic text-center py-2">
      No ADO comments yet
    </div>
  </div>
</template>
