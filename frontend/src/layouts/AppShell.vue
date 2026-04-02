<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import CommandPalette from '@/components/CommandPalette.vue'

const route = useRoute()

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
</script>

<template>
  <div class="h-screen w-screen flex overflow-hidden bg-background text-foreground">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Main column -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Titlebar / breadcrumb bar -->
      <div class="h-[28px] flex items-center px-3 border-b border-border bg-card/50 titlebar-drag shrink-0">
        <span class="text-[11px] text-muted-foreground font-medium titlebar-no-drag select-none ml-1">
          {{ breadcrumb }}
        </span>
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
