<script setup lang="ts">
import { ref, watch, type Component } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  AlertCircle,
  BarChart3,
  GitPullRequest,
  BookOpen,
  ExternalLink,
  Plus,
  X,
} from 'lucide-vue-next'
import { linkTypeColor } from '@/lib/styles'
import type { TaskLink } from '@/types'

const props = defineProps<{ taskId: number }>()

const links = ref<TaskLink[]>([])
const newUrl = ref('')
const newLabel = ref('')
const loading = ref(false)

// Type icon mapping (Lucide icons)
const typeIcons: Record<string, Component> = {
  icm: AlertCircle,
  grafana: BarChart3,
  ado: GitPullRequest,
  wiki: BookOpen,
  url: ExternalLink,
}

async function fetchLinks() {
  loading.value = true
  try {
    const { listLinks } = await import('@/api/external-links')
    links.value = (await listLinks(props.taskId)) as TaskLink[]
  } catch {
    links.value = []
  } finally {
    loading.value = false
  }
}

async function addLink() {
  const url = newUrl.value.trim()
  if (!url) return
  try {
    const { addLink: addLinkApi } = await import('@/api/external-links')
    const link = (await addLinkApi(props.taskId, url, newLabel.value.trim())) as TaskLink
    links.value.push(link)
    newUrl.value = ''
    newLabel.value = ''
  } catch (e) {
    console.warn('[ExternalLinks] Failed to add link:', e)
  }
}

async function deleteLink(id: number) {
  try {
    const { deleteLink: deleteLinkApi } = await import('@/api/external-links')
    await deleteLinkApi(id)
    links.value = links.value.filter(l => l.id !== id)
  } catch (e) {
    console.warn('[ExternalLinks] Failed to delete link:', e)
  }
}

async function openExternal(url: string) {
  try {
    const { openURL } = await import('@/api/browser')
    await openURL(url)
  } catch { window.open(url, '_blank') }
}

watch(() => props.taskId, () => fetchLinks(), { immediate: true })
</script>

<template>
  <div class="space-y-3">
    <div class="text-xs font-semibold text-muted-foreground uppercase">External Links</div>

    <!-- Link list -->
    <div
      v-for="link in links"
      :key="link.id"
      class="flex items-center gap-2 group px-2 py-1.5 rounded hover:bg-muted"
    >
      <component
        :is="typeIcons[link.type] || ExternalLink"
        :size="14"
        :class="linkTypeColor(link.type)"
      />
      <a
        :href="link.url"
        @click.prevent="openExternal(link.url)"
        class="text-sm text-foreground hover:underline truncate flex-1"
      >
        {{ link.label || link.url }}
      </a>
      <Button
        variant="ghost"
        size="icon"
        class="h-6 w-6 opacity-0 group-hover:opacity-100"
        @click="deleteLink(link.id)"
      >
        <X :size="12" />
      </Button>
    </div>

    <!-- Empty state -->
    <div v-if="!links.length && !loading" class="text-xs text-muted-foreground py-2">
      No links yet — add ICM, dashboard, or wiki URLs
    </div>

    <!-- Add link form -->
    <div class="flex gap-2">
      <Input
        v-model="newLabel"
        placeholder="Label"
        class="h-7 text-xs w-32"
        @keydown.enter="addLink"
      />
      <Input
        v-model="newUrl"
        placeholder="Paste URL"
        class="h-7 text-xs flex-1"
        @keydown.enter="addLink"
      />
      <Button
        variant="ghost"
        size="icon"
        class="h-7 w-7"
        @click="addLink"
        :disabled="!newUrl.trim()"
      >
        <Plus :size="14" />
      </Button>
    </div>
  </div>
</template>
