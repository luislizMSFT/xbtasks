<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { Separator } from '@/components/ui/separator'
import {
  LayoutDashboard,
  CheckSquare,
  FolderKanban,
  Settings,
  Sun,
  Moon,
  FlaskConical,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { mode, toggle: toggleTheme } = useTheme()

const isDev = import.meta.env.DEV

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

const playgroundItems: NavItem[] = [
  { name: 'playground-tasks', icon: FlaskConical, path: '/playground/tasks', label: '🧪 Layouts' },
]

const isActive = (path: string) => route.path === path || route.path.startsWith(path + '/')

function navigate(path: string) {
  router.push(path)
}

const isDark = computed(() => mode.value === 'dark')
</script>

<template>
  <nav class="w-14 flex-shrink-0 flex flex-col items-center py-1 gap-0.5 pt-[32px] border-r border-border bg-card titlebar-no-drag relative z-40">
    <!-- Nav items -->
    <TooltipProvider :delay-duration="400">
      <div class="flex flex-col items-center gap-1">
        <Tooltip v-for="item in navItems" :key="item.name">
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              @click="navigate(item.path)"
              :class="[
                'w-10 h-10',
                isActive(item.path)
                  ? 'bg-primary/15 text-primary hover:bg-primary/20 hover:text-primary'
                  : 'text-muted-foreground hover:text-foreground'
              ]"
            >
              <component :is="item.icon" :size="20" :stroke-width="1.75" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="right">
            {{ item.label }}
          </TooltipContent>
        </Tooltip>
      </div>

      <!-- Playground (dev only) -->
      <template v-if="isDev">
        <Separator class="my-1 w-6" />
        <Tooltip v-for="item in playgroundItems" :key="item.name">
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              @click="navigate(item.path)"
              :class="[
                'w-10 h-10',
                isActive(item.path)
                  ? 'bg-amber-500/15 text-amber-500 hover:bg-amber-500/20 hover:text-amber-500'
                  : 'text-muted-foreground/50 hover:text-muted-foreground'
              ]"
            >
              <component :is="item.icon" :size="18" :stroke-width="1.75" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="right">
            {{ item.label }}
          </TooltipContent>
        </Tooltip>
      </template>
    </TooltipProvider>

    <!-- Spacer -->
    <div class="flex-1" />

    <!-- Bottom actions -->
    <TooltipProvider :delay-duration="400">
      <div class="flex flex-col items-center gap-1 mb-2">
        <!-- Theme toggle -->
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              @click="toggleTheme"
              class="w-10 h-10 text-muted-foreground hover:text-foreground"
            >
              <Moon v-if="!isDark" :size="18" :stroke-width="1.75" />
              <Sun v-else :size="18" :stroke-width="1.75" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="right">
            {{ isDark ? 'Light mode' : 'Dark mode' }}
          </TooltipContent>
        </Tooltip>

        <!-- Settings -->
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              @click="navigate('/settings')"
              :class="[
                'w-10 h-10',
                isActive('/settings')
                  ? 'bg-primary/15 text-primary hover:bg-primary/20 hover:text-primary'
                  : 'text-muted-foreground hover:text-foreground'
              ]"
            >
              <Settings :size="18" :stroke-width="1.75" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="right">
            Settings
          </TooltipContent>
        </Tooltip>

        <!-- User avatar -->
        <div
          v-if="authStore.isAuthenticated"
          class="w-6 h-6 rounded-full bg-primary/20 text-primary flex items-center justify-center text-[10px] font-semibold mt-1 cursor-pointer"
        >
          {{ authStore.initials }}
        </div>
      </div>
    </TooltipProvider>
  </nav>
</template>
