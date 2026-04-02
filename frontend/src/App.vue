<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed } from 'vue'
import { useTheme } from './composables/useTheme'
import AppShell from './layouts/AppShell.vue'

useTheme()
const route = useRoute()
const hideShell = computed(() => route.meta?.hideShell === true)
</script>

<template>
  <AppShell v-if="!hideShell">
    <router-view />
  </AppShell>
  <router-view v-else />
</template>
