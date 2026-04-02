<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import Sidebar from '@/components/Sidebar.vue'
import CommandPalette from '@/components/CommandPalette.vue'
import { Plus, Search, Loader2, Octagon } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const taskStore = useTaskStore()

const breadcrumb = computed(() => {
  const map: Record<string, string> = {
    '/dashboard': 'Dashboard',
    '/tasks': 'Tasks',
    '/projects': 'Projects',
    '/settings': 'Settings',
    '/playground/tasks': '🧪 Task Layouts',
    '/playground/dashboard': '🧪 Dashboard Layouts',
    '/playground/detail': '🧪 Detail Layouts',
  }
  return map[route.path] || 'Team ADO Tool'
})

const showActions = computed(() => !route.path.startsWith('/login'))

function openSearch() {
  // Trigger command palette via keyboard event
  window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }))
}
</script>

<template>
  <div class="h-screen w-screen flex overflow-hidden bg-background text-foreground">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Main column -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Top bar: breadcrumb + summary + actions -->
      <div class="h-[32px] flex items-center gap-2 px-3 border-b border-border bg-card/50 titlebar-drag shrink-0">
        <!-- Breadcrumb -->
        <span class="text-[11px] text-muted-foreground font-medium titlebar-no-drag select-none">
          {{ breadcrumb }}
        </span>

        <!-- Summary badges -->
        <div v-if="showActions" class="flex items-center gap-1.5 ml-2 titlebar-no-drag">
          <Badge v-if="taskStore.stats.inProgress > 0" variant="secondary" class="h-[18px] px-1.5 text-[10px] gap-1 font-normal">
            <Loader2 :size="10" class="text-blue-500" />
            {{ taskStore.stats.inProgress }}
          </Badge>
          <Badge v-if="taskStore.stats.blocked > 0" variant="secondary" class="h-[18px] px-1.5 text-[10px] gap-1 font-normal">
            <Octagon :size="10" class="text-red-500" />
            {{ taskStore.stats.blocked }}
          </Badge>
        </div>

        <!-- Spacer (draggable) -->
        <div class="flex-1" />

        <!-- Actions (not draggable) -->
        <div v-if="showActions" class="flex items-center gap-1 titlebar-no-drag">
          <Button variant="ghost" size="sm" class="h-6 px-2 text-[11px] text-muted-foreground" @click="openSearch">
            <Search :size="12" class="mr-1" />
            Search
            <kbd class="ml-1.5 text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
          </Button>
          <Button size="sm" class="h-6 px-2 text-[11px]" @click="router.push('/tasks')">
            <Plus :size="12" class="mr-1" />
            New Task
          </Button>
        </div>
      </div>

      <!-- Main content -->
      <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
        <slot />
      </main>
    </div>

    <!-- Command Palette -->
    <CommandPalette />
  </div>
</template>
