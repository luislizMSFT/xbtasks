<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/tasks'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import Sidebar from '@/components/Sidebar.vue'
import CommandPalette from '@/components/CommandPalette.vue'
import {
  Plus, Search, PanelRightClose, PanelRightOpen,
  CircleDot, Eye, Octagon, CheckCircle2, Clock,
  GitPullRequest, MessageSquare, ArrowRight,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const taskStore = useTaskStore()
const activityOpen = ref(false)

const breadcrumb = computed(() => {
  const map: Record<string, string> = {
    '/dashboard': 'Dashboard',
    '/tasks': 'Tasks',
    '/projects': 'Projects',
    '/settings': 'Settings',
    '/playground/tasks': 'Task Layouts',
    '/playground/dashboard': 'Dashboard Layouts',
    '/playground/detail': 'Detail Layouts',
    '/ado': 'Azure DevOps',
    '/playground/ado': 'ADO Playground',
  }
  return map[route.path] || 'Team ADO Tool'
})

const showActions = computed(() => !route.path.startsWith('/login'))

const statItems = computed(() => [
  { key: 'inProgress', label: 'In Progress', count: taskStore.stats.inProgress, color: 'bg-blue-500', textColor: 'text-blue-600 dark:text-blue-400' },
  { key: 'inReview', label: 'Review', count: taskStore.tasks.filter(t => t.status === 'in_review').length, color: 'bg-violet-500', textColor: 'text-violet-600 dark:text-violet-400' },
  { key: 'blocked', label: 'Blocked', count: taskStore.stats.blocked, color: 'bg-red-500', textColor: 'text-red-600 dark:text-red-400' },
  { key: 'done', label: 'Done', count: taskStore.stats.done, color: 'bg-emerald-500', textColor: 'text-emerald-600 dark:text-emerald-400' },
])

// Global activity feed
const activityEvents = [
  { id: 1, icon: CheckCircle2, color: 'text-emerald-500', description: 'Completed "Add unit tests for auth module"', timestamp: '25m ago' },
  { id: 2, icon: GitPullRequest, color: 'text-violet-500', description: 'PR merged: Fix auth redirect loop', timestamp: '1h ago' },
  { id: 3, icon: MessageSquare, color: 'text-blue-500', description: 'Commented on "Schema migration v2"', timestamp: '2h ago' },
  { id: 4, icon: ArrowRight, color: 'text-amber-500', description: 'Moved "API rate limiting" to In Review', timestamp: '3h ago' },
  { id: 5, icon: Plus, color: 'text-muted-foreground', description: 'Created task "Update deployment docs"', timestamp: '5h ago' },
]

function openSearch() {
  window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }))
}
</script>

<template>
  <!-- App shell: 56px sidebar + content area, themed via CSS custom properties -->
  <div
    class="h-screen w-screen flex overflow-hidden text-foreground"
    :style="{ backgroundColor: 'var(--color-bg-primary, var(--surface-primary))' }"
  >
    <!-- Sidebar (w-14 = 56px) -->
    <Sidebar />

    <!-- Main column -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Top bar: breadcrumb + stats + actions -->
      <div class="h-[40px] flex items-center gap-2 px-4 border-b border-border titlebar-drag shrink-0" :style="{ backgroundColor: 'var(--color-bg-secondary, var(--surface-secondary))' }">
        <!-- Breadcrumb -->
        <span class="text-xs text-muted-foreground font-medium titlebar-no-drag select-none shrink-0">
          {{ breadcrumb }}
        </span>

        <!-- Stats pills (always visible) -->
        <div v-if="showActions" class="flex items-center gap-0.5 ml-2 titlebar-no-drag">
          <div
            v-for="stat in statItems"
            :key="stat.key"
            class="flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] text-muted-foreground"
          >
            <span :class="cn('w-1.5 h-1.5 rounded-full shrink-0', stat.color)" />
            <span class="tabular-nums font-medium" :class="stat.textColor">{{ stat.count }}</span>
            <span class="hidden lg:inline">{{ stat.label }}</span>
          </div>
        </div>

        <!-- Spacer (draggable) -->
        <div class="flex-1" />

        <!-- Actions (not draggable) -->
        <div v-if="showActions" class="flex items-center gap-1 titlebar-no-drag">
          <Button variant="ghost" size="sm" class="h-7 px-2 text-xs text-muted-foreground gap-1" @click="openSearch">
            <Search :size="13" />
            <kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
          </Button>
          <Button size="sm" class="h-7 px-2.5 text-xs gap-1" @click="router.push('/tasks')">
            <Plus :size="13" />
            New
          </Button>
          <Button
            variant="ghost" size="icon"
            class="h-7 w-7 text-muted-foreground"
            @click="activityOpen = !activityOpen"
            :title="activityOpen ? 'Hide activity' : 'Show activity'"
          >
            <Clock :size="14" />
          </Button>
        </div>
      </div>

      <!-- Content + Activity sidebar -->
      <div class="flex-1 flex overflow-hidden">
        <!-- Main content -->
        <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
          <slot />
        </main>

        <!-- Activity sidebar (global, hideable) -->
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="w-0 opacity-0"
          enter-to-class="w-[240px] opacity-100"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="w-[240px] opacity-100"
          leave-to-class="w-0 opacity-0"
        >
          <aside
            v-if="activityOpen"
            class="w-[240px] shrink-0 border-l border-border bg-card/30 flex flex-col overflow-hidden"
          >
            <div class="flex items-center justify-between px-3 py-2 border-b border-border shrink-0">
              <span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Activity</span>
              <Button variant="ghost" size="icon" class="h-6 w-6" @click="activityOpen = false">
                <PanelRightClose :size="12" />
              </Button>
            </div>
            <ScrollArea class="flex-1 h-full">
              <div class="px-3 py-2 space-y-0.5">
                <div
                  v-for="event in activityEvents"
                  :key="event.id"
                  class="flex items-start gap-2 py-2"
                >
                  <div class="w-5 h-5 rounded-full bg-muted flex items-center justify-center flex-shrink-0 mt-0.5">
                    <component :is="event.icon" :size="10" :class="event.color" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-[12px] text-foreground leading-snug">{{ event.description }}</p>
                    <span class="text-[10px] text-muted-foreground/50">{{ event.timestamp }}</span>
                  </div>
                </div>
              </div>
            </ScrollArea>
          </aside>
        </Transition>
      </div>
    </div>

    <!-- Command Palette -->
    <CommandPalette />
  </div>
</template>
