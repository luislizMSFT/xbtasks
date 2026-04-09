<script setup lang="ts">
import { ref, watch } from 'vue'
import { relativeTime } from '@/lib/date'
import { Badge } from '@/components/ui/badge'
import { Globe, Lock } from 'lucide-vue-next'
import type { TaskComment } from '@/types'

const props = defineProps<{
  taskId: number
}>()

const comments = ref<TaskComment[]>([])
const loading = ref(false)

async function fetchComments() {
  loading.value = true
  try {
    const { listComments } = await import('@/api/comments')
    comments.value = (await listComments(props.taskId)) as TaskComment[]
  } catch {
    comments.value = []
  } finally {
    loading.value = false
  }
}

defineExpose({ fetchComments, addComment })

async function addComment(content: string) {
  if (!content.trim()) return
  try {
    const { addComment: addCommentApi } = await import('@/api/comments')
    const c = (await addCommentApi(props.taskId, content)) as TaskComment
    comments.value.push(c)
  } catch (e) {
    console.warn('[CommentsSection] Failed to add comment:', e)
  }
}

watch(() => props.taskId, () => fetchComments(), { immediate: true })
</script>

<template>
  <div class="space-y-2">
    <!-- Comments list (full width bubbles) -->
    <div v-if="comments.length > 0" class="flex flex-col gap-2 max-h-60 overflow-y-auto">
      <div
        v-for="c in comments"
        :key="c.id"
        class="rounded-lg px-3 py-2"
        :class="c.isPublic
          ? 'bg-blue-500/10 border border-blue-500/20'
          : 'bg-muted/60 border border-border/50'"
      >
        <div class="flex items-center gap-1.5 mb-1">
          <Badge
            v-if="c.isPublic"
            variant="outline"
            class="text-[9px] h-3.5 px-1 border-blue-500/30 text-blue-500"
          >
            <Globe :size="8" class="mr-0.5" /> Public
          </Badge>
          <Badge
            v-else
            variant="outline"
            class="text-[9px] h-3.5 px-1 border-muted-foreground/40 text-muted-foreground"
          >
            <Lock :size="8" class="mr-0.5" /> Private
          </Badge>
          <span class="text-[9px] text-muted-foreground/60 ml-auto tabular-nums">
            {{ relativeTime(c.createdAt) }}
          </span>
        </div>
        <div class="text-[13px] text-foreground whitespace-pre-wrap leading-relaxed">{{ c.content }}</div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-if="!comments.length && !loading" class="text-[11px] text-muted-foreground/40 italic text-center py-2">
      No notes yet
    </div>
  </div>
</template>
