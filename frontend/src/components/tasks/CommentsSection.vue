<script setup lang="ts">
import { ref, watch } from 'vue'
import { relativeTime } from '@/lib/date'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
import { Globe, Lock, Upload } from 'lucide-vue-next'
import type { TaskComment } from '@/types'

const props = defineProps<{
  taskId: number
  isPublicTask: boolean // whether the task is linked to ADO
}>()

const comments = ref<TaskComment[]>([])
const newComment = ref('')
const pushing = ref<number | null>(null) // comment ID being pushed
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

async function addComment() {
  const content = newComment.value.trim()
  if (!content) return
  try {
    const { addComment: addCommentApi } = await import('@/api/comments')
    const c = (await addCommentApi(props.taskId, content)) as TaskComment
    comments.value.push(c)
    newComment.value = ''
  } catch (e) {
    console.warn('[CommentsSection] Failed to add comment:', e)
  }
}

async function pushToADO(commentId: number) {
  pushing.value = commentId
  try {
    const { pushCommentToADO } = await import('@/api/comments')
    await pushCommentToADO(commentId)
    const idx = comments.value.findIndex(c => c.id === commentId)
    if (idx !== -1) comments.value[idx].isPublic = true
  } catch (e) {
    console.warn('[CommentsSection] Failed to push comment to ADO:', e)
  } finally {
    pushing.value = null
  }
}

watch(() => props.taskId, () => fetchComments(), { immediate: true })
</script>

<template>
  <div class="space-y-3">
    <div class="text-xs font-semibold text-muted-foreground uppercase">Comments</div>

    <!-- Comments list (bubble style) -->
    <div v-if="comments.length > 0" class="flex flex-col gap-2 max-h-60 overflow-y-auto">
      <div
        v-for="c in comments"
        :key="c.id"
        class="rounded-lg px-3 py-2 max-w-[92%] self-end"
        :class="c.isPublic
          ? 'bg-blue-500/10 border border-blue-500/20'
          : 'bg-muted/60 border border-border/50'"
      >
        <!-- Bubble header: visibility + timestamp -->
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
        <!-- Content -->
        <div class="text-[13px] text-foreground whitespace-pre-wrap leading-relaxed">{{ c.content }}</div>
        <!-- Push to ADO (only private on public tasks) -->
        <div v-if="!c.isPublic && isPublicTask" class="mt-1.5 flex justify-end">
          <Button
            variant="ghost"
            size="sm"
            class="h-5 text-[9px] gap-0.5 px-1.5 text-blue-500 hover:text-blue-600"
            @click="pushToADO(c.id)"
            :disabled="pushing === c.id"
          >
            <Upload :size="10" />
            {{ pushing === c.id ? 'Pushing...' : 'Push to ADO' }}
          </Button>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-if="!comments.length && !loading" class="text-[11px] text-muted-foreground/40 italic text-center py-2">
      No comments yet
    </div>

    <!-- Add comment -->
    <div class="space-y-1.5">
      <Textarea
        v-model="newComment"
        placeholder="Add a comment..."
        class="text-xs min-h-[50px] resize-none"
        :rows="2"
        @keydown.meta.enter="addComment"
        @keydown.ctrl.enter="addComment"
      />
      <div class="flex justify-between items-center">
        <span class="text-[9px] text-muted-foreground/50">Comments are private by default</span>
        <Button
          size="sm"
          class="h-6 text-[10px]"
          @click="addComment"
          :disabled="!newComment.trim()"
        >
          Add Comment
        </Button>
      </div>
    </div>
  </div>
</template>
