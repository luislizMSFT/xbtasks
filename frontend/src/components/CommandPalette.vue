<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useMagicKeys } from '@vueuse/core'
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import {
  LayoutDashboard,
  CheckSquare,
  FolderKanban,
  Plus,
  Search,
  LogOut,
} from 'lucide-vue-next'

const router = useRouter()
const isOpen = ref(false)

interface Action {
  id: string
  label: string
  icon: any
  shortcut?: string
  action: () => void
}

const actions: Action[] = [
  { id: 'dashboard', label: 'Go to Dashboard', icon: LayoutDashboard, shortcut: '⌘1', action: () => router.push('/dashboard') },
  { id: 'tasks', label: 'Go to Tasks', icon: CheckSquare, shortcut: '⌘2', action: () => router.push('/tasks') },
  { id: 'projects', label: 'Go to Projects', icon: FolderKanban, shortcut: '⌘3', action: () => router.push('/projects') },
  { id: 'create', label: 'Create Task', icon: Plus, shortcut: '⌘N', action: () => { router.push('/tasks'); /* TODO: trigger inline create */ } },
  { id: 'search', label: 'Search tasks…', icon: Search, action: () => router.push('/tasks') },
  { id: 'signout', label: 'Sign Out', icon: LogOut, action: () => router.push('/login') },
]

function execute(action: Action) {
  isOpen.value = false
  action.action()
}

// ⌘K / Ctrl+K shortcut
const keys = useMagicKeys()
watch(keys['Meta+k']!, (pressed) => { if (pressed) { isOpen.value = !isOpen.value } })
watch(keys['Ctrl+k']!, (pressed) => { if (pressed) { isOpen.value = !isOpen.value } })
</script>

<template>
  <CommandDialog v-model:open="isOpen">
    <CommandInput placeholder="Type a command or search…" />
    <CommandList>
      <CommandEmpty>No results found.</CommandEmpty>
      <CommandGroup heading="Actions">
        <CommandItem
          v-for="action in actions"
          :key="action.id"
          :value="action.label"
          @select="execute(action)"
          class="flex items-center gap-3"
        >
          <component :is="action.icon" :size="16" :stroke-width="1.75" class="opacity-60" />
          <span class="flex-1">{{ action.label }}</span>
          <span v-if="action.shortcut" class="text-xs text-muted-foreground">{{ action.shortcut }}</span>
        </CommandItem>
      </CommandGroup>
    </CommandList>
  </CommandDialog>
</template>
