<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useMagicKeys } from '@vueuse/core'
import {
  LayoutDashboard,
  CheckSquare,
  FolderKanban,
  Plus,
  Search,
  LogOut,
  Command,
} from 'lucide-vue-next'

const router = useRouter()
const isOpen = ref(false)
const query = ref('')
const selectedIndex = ref(0)
const inputRef = ref<HTMLInputElement | null>(null)

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

const filtered = computed(() => {
  if (!query.value) return actions
  const q = query.value.toLowerCase()
  return actions.filter(a => a.label.toLowerCase().includes(q))
})

watch(filtered, () => { selectedIndex.value = 0 })

function open() {
  isOpen.value = true
  query.value = ''
  selectedIndex.value = 0
  nextTick(() => inputRef.value?.focus())
}

function close() {
  isOpen.value = false
}

function execute(action: Action) {
  close()
  action.action()
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    selectedIndex.value = Math.min(selectedIndex.value + 1, filtered.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
  } else if (e.key === 'Enter') {
    e.preventDefault()
    const action = filtered.value[selectedIndex.value]
    if (action) execute(action)
  }
}

// ⌘K / Ctrl+K shortcut
const keys = useMagicKeys()
watch(keys['Meta+k']!, (pressed) => { if (pressed) { isOpen.value ? close() : open() } })
watch(keys['Ctrl+k']!, (pressed) => { if (pressed) { isOpen.value ? close() : open() } })
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 z-[100] flex justify-center"
        style="padding-top: 20vh"
        @click.self="close"
      >
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/50 dark:bg-black/70" @click="close" />

        <!-- Panel -->
        <Transition
          enter-active-class="transition duration-150 ease-out"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition duration-100 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
          appear
        >
          <div
            v-if="isOpen"
            class="relative w-[560px] max-h-[400px] rounded-xl border border-border-default bg-surface-primary shadow-2xl overflow-hidden"
            @keydown="onKeydown"
          >
            <!-- Search input -->
            <div class="flex items-center gap-3 px-4 py-3 border-b border-border-default">
              <Search :size="18" class="text-text-secondary flex-shrink-0" />
              <input
                ref="inputRef"
                v-model="query"
                type="text"
                placeholder="Type a command or search…"
                class="flex-1 bg-transparent text-sm text-text-primary placeholder-text-secondary outline-none"
                @keydown.esc="close"
              />
              <kbd class="hidden sm:flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[10px] font-medium text-text-secondary bg-surface-tertiary border border-border-default">
                esc
              </kbd>
            </div>

            <!-- Actions list -->
            <div class="py-1.5 overflow-y-auto max-h-[340px]">
              <div
                v-for="(action, index) in filtered"
                :key="action.id"
                @click="execute(action)"
                @mouseenter="selectedIndex = index"
                class="flex items-center gap-3 px-4 py-2.5 cursor-pointer transition-colors duration-75"
                :class="[
                  index === selectedIndex
                    ? 'bg-accent/10 text-accent'
                    : 'text-text-primary hover:bg-surface-tertiary'
                ]"
              >
                <component :is="action.icon" :size="16" :stroke-width="1.75" class="flex-shrink-0 opacity-60" />
                <span class="flex-1 text-sm">{{ action.label }}</span>
                <span v-if="action.shortcut" class="text-xs text-text-secondary opacity-60">{{ action.shortcut }}</span>
              </div>

              <div v-if="filtered.length === 0" class="px-4 py-8 text-center text-sm text-text-secondary">
                No results for "{{ query }}"
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
