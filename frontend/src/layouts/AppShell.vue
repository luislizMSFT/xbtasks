<script setup lang="ts">
import Sidebar from '@/components/Sidebar.vue'
import CommandPalette from '@/components/CommandPalette.vue'
</script>

<template>
  <div class="h-screen w-screen flex overflow-hidden bg-background text-foreground">
    <!-- macOS titlebar drag region -->
    <div class="fixed top-0 left-0 right-0 h-[50px] titlebar-drag z-50" />

    <!-- Sidebar -->
    <Sidebar />

    <!-- Main content -->
    <main class="flex-1 flex flex-col min-w-0 pt-[50px] overflow-hidden">
      <slot />
    </main>

    <!-- Command Palette -->
    <CommandPalette />
  </div>
</template>
