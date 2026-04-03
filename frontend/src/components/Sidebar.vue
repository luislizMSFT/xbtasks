<script setup lang="ts">
import { computed } from 'vue'
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
import {
  LayoutDashboard,
  CheckSquare,
  FolderKanban,
  Settings,
  Sun,
  Moon,
  GitBranch,
} from 'lucide-vue-next'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { mode, toggle: toggleTheme } = useTheme()

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
  { name: 'ado', icon: AzureDevOpsIcon, path: '/ado', label: 'Azure DevOps' },
  { name: 'dependencies', icon: GitBranch, path: '/dependencies', label: 'Dependencies' },
]

const isActive = (path: string) => route.path === path || route.path.startsWith(path + '/')

function navigate(path: string) {
  router.push(path)
}

const isDark = computed(() => mode.value === 'dark')
</script>

<template>
  <nav class="w-14 flex-shrink-0 flex flex-col items-center py-1 gap-0.5 pt-[44px] border-r border-border bg-card titlebar-no-drag relative z-40">
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
