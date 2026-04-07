<script setup lang="ts">
import { ref, watch } from 'vue'
import { relativeTime } from '@/lib/date'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Textarea } from '@/components/ui/textarea'
import { Globe, Lock, Upload } from 'lucide-vue-next'

const props = defineProps<{
  taskId: number
  isPublicTask: boolean // whether the task is linked to ADO
}>()

interface TaskComment {
  id: number
  taskId: number
  content: string
  isPublic: boolean
  adoCommentId: string
  createdAt: string
  updatedAt: string
}

const comments = ref<TaskComment[]>([])
const newComment = ref('')
const pushing = ref<number | null>(null) // comment ID being pushed
const loading = ref(false)

async function fetchComments() {
  loading.value = true
  try {
    const { ListComments } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
    comments.value = (await ListComments(props.taskId)) as TaskComment[]
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
    const { AddComment } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
    const c = (await AddComment(props.taskId, content)) as TaskComment
    comments.value.push(c)
    newComment.value = ''
  } catch (e) {
    console.warn('[CommentsSection] Failed to add comment:', e)
  }
}

async function pushToADO(commentId: number) {
  pushing.value = commentId
  try {
    const { PushCommentToADO } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
    await PushCommentToADO(commentId)
    // Mark as public locally
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

    <ScrollArea class="max-h-60">
      <div
        v-for="c in comments"
        :key="c.id"
        class="px-2 py-2 border-b border-border last:border-0"
      >
        <!-- Public/Private indicator -->
        <div class="flex items-center gap-1 mb-1">
          <Badge
            v-if="c.isPublic"
            variant="outline"
            class="text-[10px] h-4 border-blue-500/30 text-blue-500"
          >
            <Globe :size="10" class="mr-0.5" /> Public
          </Badge>
          <Badge
            v-else
            variant="outline"
            class="text-[10px] h-4 border-muted-foreground/30 text-muted-foreground"
          >
            <Lock :size="10" class="mr-0.5" /> Private
          </Badge>
          <span class="text-[10px] text-muted-foreground ml-auto">
            {{ relativeTime(c.createdAt) }}
          </span>
        </div>
        <!-- Content -->
        <div class="text-sm text-foreground whitespace-pre-wrap">{{ c.content }}</div>
        <!-- Push to ADO button (only for private comments on public tasks) -->
        <div v-if="!c.isPublic && isPublicTask" class="mt-1">
          <Button
            variant="ghost"
            size="sm"
            class="h-6 text-xs"
            @click="pushToADO(c.id)"
            :disabled="pushing === c.id"
          >
            <Upload :size="12" class="mr-1" />
            {{ pushing === c.id ? 'Pushing...' : 'Push to ADO' }}
          </Button>
        </div>
      </div>
    </ScrollArea>

    <!-- Empty state -->
    <div v-if="!comments.length && !loading" class="text-xs text-muted-foreground py-2">
      No comments yet
    </div>

    <!-- Add comment -->
    <div class="flex gap-2">
      <Textarea
        v-model="newComment"
        placeholder="Add a comment..."
        class="text-xs min-h-[60px] resize-none"
        :rows="2"
      />
    </div>
    <div class="flex justify-between items-center">
      <span class="text-[10px] text-muted-foreground">Comments are private by default</span>
      <Button
        size="sm"
        class="h-7 text-xs"
        @click="addComment"
        :disabled="!newComment.trim()"
      >
        Add Comment
      </Button>
    </div>
  </div>
</template>
