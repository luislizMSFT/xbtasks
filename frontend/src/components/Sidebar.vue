<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useTheme } from '../composables/useTheme'
import {
  LayoutDashboard,
  CheckSquare,
  FolderKanban,
  Settings,
  Sun,
  Moon,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { mode, toggle: toggleTheme } = useTheme()

const hoveredItem = ref<string | null>(null)
const tooltipTimeout = ref<ReturnType<typeof setTimeout> | null>(null)
const showTooltip = ref(false)

interface NavItem {
  name: string
  icon: any
  path: string
  label: string
}

const navItems: NavItem[] = [
  { name: 'dashboard', icon: LayoutDashboard, path: '/dashboard', label: 'Dashboard' },
  { name: 'tasks', icon: CheckSquare, path: '/tasks', label: 'Tasks' },
  { name: 'projects', icon: FolderKanban, path: '/projects', label: 'Projects' },
]

const isActive = (path: string) => route.path === path

function navigate(path: string) {
  router.push(path)
}

function onMouseEnter(name: string) {
  hoveredItem.value = name
  tooltipTimeout.value = setTimeout(() => { showTooltip.value = true }, 500)
}

function onMouseLeave() {
  if (tooltipTimeout.value) clearTimeout(tooltipTimeout.value)
  hoveredItem.value = null
  showTooltip.value = false
}

const isDark = computed(() => mode.value === 'dark')
</script>

<template>
  <nav class="w-14 flex-shrink-0 flex flex-col items-center py-3 gap-1 pt-[56px] border-r border-border-default bg-surface-secondary titlebar-no-drag relative z-40">
    <!-- Nav items -->
    <div class="flex flex-col items-center gap-1">
      <div
        v-for="item in navItems"
        :key="item.name"
        class="relative group"
        @mouseenter="onMouseEnter(item.name)"
        @mouseleave="onMouseLeave"
      >
        <button
          @click="navigate(item.path)"
          class="w-10 h-10 flex items-center justify-center rounded-lg transition-colors duration-150"
          :class="[
            isActive(item.path)
              ? 'bg-accent/15 text-accent'
              : 'text-text-secondary hover:text-text-primary hover:bg-surface-tertiary'
          ]"
        >
          <component :is="item.icon" :size="20" :stroke-width="1.75" />
        </button>

        <!-- Tooltip -->
        <Transition
          enter-active-class="transition duration-150 ease-out"
          enter-from-class="opacity-0 translate-x-[-4px]"
          enter-to-class="opacity-100 translate-x-0"
          leave-active-class="transition duration-100 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div
            v-if="hoveredItem === item.name && showTooltip"
            class="absolute left-full ml-2 top-1/2 -translate-y-1/2 px-2.5 py-1 rounded-md text-xs font-medium whitespace-nowrap z-50 bg-zinc-800 text-zinc-100 dark:bg-zinc-200 dark:text-zinc-900 shadow-lg"
          >
            {{ item.label }}
          </div>
        </Transition>
      </div>
    </div>

    <!-- Spacer -->
    <div class="flex-1" />

    <!-- Bottom actions -->
    <div class="flex flex-col items-center gap-1 mb-2">
      <!-- Theme toggle -->
      <button
        @click="toggleTheme"
        class="w-10 h-10 flex items-center justify-center rounded-lg text-text-secondary hover:text-text-primary hover:bg-surface-tertiary transition-colors duration-150"
      >
        <Moon v-if="!isDark" :size="18" :stroke-width="1.75" />
        <Sun v-else :size="18" :stroke-width="1.75" />
      </button>

      <!-- Settings -->
      <button
        class="w-10 h-10 flex items-center justify-center rounded-lg text-text-secondary hover:text-text-primary hover:bg-surface-tertiary transition-colors duration-150"
      >
        <Settings :size="18" :stroke-width="1.75" />
      </button>

      <!-- User avatar -->
      <div
        v-if="authStore.isAuthenticated"
        class="w-6 h-6 rounded-full bg-accent/20 text-accent flex items-center justify-center text-[10px] font-semibold mt-1 cursor-pointer"
      >
        {{ authStore.initials }}
      </div>
    </div>
  </nav>
</template>
