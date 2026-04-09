<script setup lang="ts">
import { ref, watch } from 'vue'
import DOMPurify from 'dompurify'
import { relativeTime } from '@/lib/date'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { ADOComment } from '@/types'

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

const comments = ref<ADOComment[]>([])
const loading = ref(false)

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

defineExpose({ fetchComments, addReply })

async function addReply(content: string) {
  if (!content.trim()) return
  try {
    const { replyToADOComment } = await import('@/api/comments')
    const c = await replyToADOComment(props.taskId, content)
    if (c) comments.value.push(c as ADOComment)
  } catch (e) {
    console.warn('[ADODiscussion] Failed to reply:', e)
  }
}

watch(() => props.taskId, () => fetchComments(), { immediate: true })
</script>

<template>
  <div class="space-y-2">
    <!-- Loading -->
    <div v-if="loading" class="py-4">
      <LoadingSpinner size="sm" label="Loading ADO comments..." />
    </div>

    <!-- Comments list (full width bubbles) -->
    <div v-else-if="comments.length > 0" class="flex flex-col gap-2 max-h-60 overflow-y-auto">
      <div
        v-for="c in comments"
        :key="c.id"
        class="rounded-lg px-3 py-2 bg-blue-500/5 border border-blue-500/10"
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
