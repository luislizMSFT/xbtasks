<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import Sidebar from '@/components/Sidebar.vue'
import CommandPalette from '@/components/CommandPalette.vue'
import { useTaskStore } from '@/stores/tasks'
import { relativeTime } from '@/lib/date'
import {
  Plus, Search, CheckCircle2, Circle, Play, AlertCircle, Clock, PanelRightClose,
  CheckSquare, Link,
} from 'lucide-vue-next'
import {
  DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem,
} from '@/components/ui/dropdown-menu'

const route = useRoute()
const router = useRouter()
const activityOpen = ref(false)
const taskStore = useTaskStore()

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

const statusMeta: Record<string, { icon: typeof CheckCircle2; color: string; label: string }> = {
  done:        { icon: CheckCircle2, color: 'text-emerald-500', label: 'Completed' },
  in_progress: { icon: Play,        color: 'text-blue-500',    label: 'In progress' },
  blocked:     { icon: AlertCircle,  color: 'text-red-500',     label: 'Blocked' },
  in_review:   { icon: Circle,      color: 'text-amber-500',   label: 'In review' },
  todo:        { icon: Circle,      color: 'text-muted-foreground', label: 'To do' },
}
const defaultMeta = { icon: Circle, color: 'text-muted-foreground', label: 'Updated' }

const activityEvents = computed(() =>
  [...taskStore.tasks]
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    .slice(0, 10)
    .map(t => {
      const meta = statusMeta[t.status] ?? defaultMeta
      return {
        id: t.id,
        icon: meta.icon,
        color: meta.color,
        description: `${meta.label}: "${t.title}"`,
        timestamp: relativeTime(t.updatedAt),
      }
    })
)

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
      <div class="h-[46px] flex items-center gap-2 px-4 border-b border-border titlebar-drag shrink-0" :style="{ backgroundColor: 'var(--color-bg-secondary, var(--surface-secondary))' }">
        <!-- Breadcrumb -->
        <span class="text-sm font-medium text-foreground titlebar-no-drag select-none shrink-0">
          {{ breadcrumb }}
        </span>

        <!-- Teleport target for page-specific actions -->
        <div id="topbar-actions" class="flex items-center gap-2 ml-3 titlebar-no-drag" />

        <!-- Spacer (draggable) -->
        <div class="flex-1" />

        <!-- Actions (not draggable) -->
        <div v-if="showActions" class="flex items-center gap-2 titlebar-no-drag">
          <Button variant="outline" size="sm" class="h-8 px-3 text-xs text-muted-foreground gap-2 min-w-[160px] justify-start" @click="openSearch">
            <Search :size="14" />
            <span class="flex-1 text-left">Search...</span>
            <kbd class="text-[9px] bg-muted px-1 rounded font-mono">⌘K</kbd>
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button size="sm" class="h-7 px-2.5 text-xs gap-1">
                <Plus :size="13" />
                New
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" class="w-44">
              <DropdownMenuItem class="text-xs gap-2" @click="openSearch()">
                <CheckSquare :size="13" />
                Task
              </DropdownMenuItem>
              <DropdownMenuItem class="text-xs gap-2" @click="router.push('/ado')">
                <Link :size="13" />
                Import from ADO
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <Button
            variant="ghost" size="icon"
            class="h-7 w-7 text-muted-foreground"
            @click="activityOpen = !activityOpen"
            :title="activityOpen ? 'Hide activity' : 'Show activity'"
          >
            <Clock :size="15" />
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
              <div v-if="activityEvents.length === 0" class="px-3 py-6 text-center">
                <p class="text-[11px] text-muted-foreground">No recent activity</p>
              </div>
              <div v-else class="px-3 py-2 space-y-0.5">
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
