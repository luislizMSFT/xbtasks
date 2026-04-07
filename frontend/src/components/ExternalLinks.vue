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

const props = defineProps<{ taskId: number }>()

interface TaskLink {
  id: number
  taskId: number
  url: string
  label: string
  type: string // 'url', 'icm', 'grafana', 'ado', 'wiki'
  createdAt: string
}

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

const typeColors: Record<string, string> = {
  icm: 'text-red-500',
  grafana: 'text-green-500',
  ado: 'text-blue-500',
  wiki: 'text-purple-500',
  url: 'text-muted-foreground',
}

async function fetchLinks() {
  loading.value = true
  try {
    const { ListLinks } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/externallinksservice')
    links.value = (await ListLinks(props.taskId)) as TaskLink[]
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
    const { AddLink } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/externallinksservice')
    const link = (await AddLink(props.taskId, url, newLabel.value.trim())) as TaskLink
    links.value.push(link)
    newUrl.value = ''
    newLabel.value = ''
  } catch (e) {
    console.warn('[ExternalLinks] Failed to add link:', e)
  }
}

async function deleteLink(id: number) {
  try {
    const { DeleteLink } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/externallinksservice')
    await DeleteLink(id)
    links.value = links.value.filter(l => l.id !== id)
  } catch (e) {
    console.warn('[ExternalLinks] Failed to delete link:', e)
  }
}

async function openExternal(url: string) {
  try {
    const { OpenURL } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/browserservice')
    await OpenURL(url)
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
        v-model="newUrl"
        placeholder="Paste URL"
        class="h-7 text-xs flex-1"
        @keydown.enter="addLink"
      />
      <Input
        v-model="newLabel"
        placeholder="Label (optional)"
        class="h-7 text-xs w-24"
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
