<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useMagicKeys } from '@vueuse/core'
import { useTaskStore } from '@/stores/tasks'
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '@/components/ui/command'
import {
  LayoutDashboard,
  CheckSquare,
  FolderKanban,
  Plus,
  Settings,
  LogOut,
  ArrowRight,
  Bug,
  CircleCheck,
  BookOpen,
  GitPullRequest,
  FlaskConical,
} from 'lucide-vue-next'

const router = useRouter()
const taskStore = useTaskStore()
const isOpen = ref(false)

const isDev = import.meta.env.DEV

interface Action {
  id: string
  label: string
  icon: any
  hint?: string
  shortcut?: string
  action: () => void
}

const navigationActions: Action[] = [
  { id: 'nav-dashboard', label: 'Dashboard', icon: LayoutDashboard, hint: 'Overview', shortcut: '⌘1', action: () => router.push('/dashboard') },
  { id: 'nav-tasks', label: 'Tasks', icon: CheckSquare, hint: 'My tasks', shortcut: '⌘2', action: () => router.push('/tasks') },
  { id: 'nav-projects', label: 'Projects', icon: FolderKanban, hint: 'All projects', shortcut: '⌘3', action: () => router.push('/projects') },
  { id: 'nav-settings', label: 'Settings', icon: Settings, hint: 'Configuration', shortcut: '⌘,', action: () => router.push('/settings') },
]

const devActions: Action[] = [
  { id: 'pg-tasks', label: '🧪 Task Layouts', icon: FlaskConical, hint: 'Compare layouts', action: () => router.push('/playground/tasks') },
  { id: 'pg-dashboard', label: '🧪 Dashboard Layouts', icon: FlaskConical, hint: 'Compare dashboards', action: () => router.push('/playground/dashboard') },
  { id: 'pg-detail', label: '🧪 Detail Layouts', icon: FlaskConical, hint: 'Compare panels', action: () => router.push('/playground/detail') },
]

const createActions: Action[] = [
  { id: 'create-task', label: 'New Task', icon: Plus, hint: 'Create personal task', shortcut: '⌘N', action: () => router.push('/tasks') },
  { id: 'create-project', label: 'New Project', icon: FolderKanban, hint: 'Create project', action: () => router.push('/projects') },
]

const adoTypeIcon = (type: string) => {
  const map: Record<string, any> = { bug: Bug, task: CircleCheck, user_story: BookOpen }
  return map[type] || CircleCheck
}

// Recent tasks as searchable items
const taskActions = computed<Action[]>(() =>
  taskStore.tasks.slice(0, 8).map(t => ({
    id: `task-${t.id}`,
    label: t.title,
    icon: t.adoId ? adoTypeIcon('task') : CheckSquare,
    hint: [t.status.replace('_', ' '), t.priority, t.adoId].filter(Boolean).join(' · '),
    action: () => { router.push('/tasks'); taskStore.selectTask(t.id) },
  }))
)

const accountActions: Action[] = [
  { id: 'signout', label: 'Sign Out', icon: LogOut, hint: 'Log out of your account', action: () => router.push('/login') },
]

function execute(action: Action) {
  isOpen.value = false
  action.action()
}

// ⌘K / Ctrl+K shortcut
const keys = useMagicKeys()
watch(keys['Meta+k']!, (pressed) => { if (pressed) isOpen.value = !isOpen.value })
watch(keys['Ctrl+k']!, (pressed) => { if (pressed) isOpen.value = !isOpen.value })

// Also listen for dispatched events from top bar search button
window.addEventListener('keydown', () => {})
</script>

<template>
  <CommandDialog v-model:open="isOpen">
    <CommandInput placeholder="Jump to anything…" />
    <CommandList>
      <CommandEmpty>No results found.</CommandEmpty>

      <!-- Navigation -->
      <CommandGroup heading="Navigate">
        <CommandItem
          v-for="action in navigationActions"
          :key="action.id"
          :value="action.label"
          @select="execute(action)"
          class="flex items-center gap-2.5"
        >
          <component :is="action.icon" :size="14" :stroke-width="1.75" class="text-muted-foreground" />
          <span class="flex-1 text-sm">{{ action.label }}</span>
          <span v-if="action.hint" class="text-[11px] text-muted-foreground">{{ action.hint }}</span>
          <kbd v-if="action.shortcut" class="text-[10px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded font-mono">{{ action.shortcut }}</kbd>
        </CommandItem>
      </CommandGroup>

      <CommandSeparator />

      <!-- Create -->
      <CommandGroup heading="Create">
        <CommandItem
          v-for="action in createActions"
          :key="action.id"
          :value="action.label"
          @select="execute(action)"
          class="flex items-center gap-2.5"
        >
          <component :is="action.icon" :size="14" :stroke-width="1.75" class="text-muted-foreground" />
          <span class="flex-1 text-sm">{{ action.label }}</span>
          <span v-if="action.hint" class="text-[11px] text-muted-foreground">{{ action.hint }}</span>
          <kbd v-if="action.shortcut" class="text-[10px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded font-mono">{{ action.shortcut }}</kbd>
        </CommandItem>
      </CommandGroup>

      <CommandSeparator />

      <!-- Tasks (searchable) -->
      <CommandGroup v-if="taskActions.length > 0" heading="Tasks">
        <CommandItem
          v-for="action in taskActions"
          :key="action.id"
          :value="action.label"
          @select="execute(action)"
          class="flex items-center gap-2.5"
        >
          <component :is="action.icon" :size="14" :stroke-width="1.75" class="text-muted-foreground" />
          <span class="flex-1 text-sm truncate">{{ action.label }}</span>
          <span class="text-[11px] text-muted-foreground truncate max-w-[180px]">{{ action.hint }}</span>
          <ArrowRight :size="12" class="text-muted-foreground/50" />
        </CommandItem>
      </CommandGroup>

      <!-- Playground (dev only) -->
      <template v-if="isDev">
        <CommandSeparator />
        <CommandGroup heading="🧪 Playground">
          <CommandItem
            v-for="action in devActions"
            :key="action.id"
            :value="action.label"
            @select="execute(action)"
            class="flex items-center gap-2.5"
          >
            <component :is="action.icon" :size="14" :stroke-width="1.75" class="text-amber-500/70" />
            <span class="flex-1 text-sm">{{ action.label }}</span>
            <span class="text-[11px] text-muted-foreground">{{ action.hint }}</span>
          </CommandItem>
        </CommandGroup>
      </template>

      <CommandSeparator />

      <!-- Account -->
      <CommandGroup heading="Account">
        <CommandItem
          v-for="action in accountActions"
          :key="action.id"
          :value="action.label"
          @select="execute(action)"
          class="flex items-center gap-2.5"
        >
          <component :is="action.icon" :size="14" :stroke-width="1.75" class="text-muted-foreground" />
          <span class="flex-1 text-sm">{{ action.label }}</span>
          <span class="text-[11px] text-muted-foreground">{{ action.hint }}</span>
        </CommandItem>
      </CommandGroup>
    </CommandList>
  </CommandDialog>
</template>
