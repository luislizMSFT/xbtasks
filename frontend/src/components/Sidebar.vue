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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  LayoutDashboard,
  CheckSquare,
  FolderKanban,
  Settings,
  Sun,
  Moon,
  Network,
  LogOut,
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
  { name: 'ado', icon: AzureDevOpsIcon, path: '/ado', label: 'Azure DevOps' },
  { name: 'dependencies', icon: Network, path: '/dependencies', label: 'Dependencies' },
]

const isActive = (path: string) => route.path === path || route.path.startsWith(path + '/')

function navigate(path: string) {
  router.push(path)
}

async function handleSignOut() {
  await authStore.signOut()
  router.push('/login')
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
              :aria-label="item.label"
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
              :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
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
              aria-label="Settings"
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

        <!-- User avatar with sign-out dropdown -->
        <DropdownMenu v-if="authStore.isAuthenticated">
          <DropdownMenuTrigger as-child>
            <button
              aria-label="User menu"
              class="w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center text-xs font-semibold mt-1 cursor-pointer hover:bg-primary/30 transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-1"
            >
              {{ authStore.initials }}
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent side="right" align="end" class="w-48">
            <DropdownMenuLabel class="font-normal">
              <div class="flex flex-col space-y-1">
                <p class="text-sm font-medium leading-none">{{ authStore.user?.displayName }}</p>
                <p class="text-xs leading-none text-muted-foreground">{{ authStore.user?.email }}</p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="handleSignOut" class="text-destructive focus:text-destructive">
              <LogOut class="mr-2 h-4 w-4" />
              Sign Out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </TooltipProvider>
  </nav>
</template>
