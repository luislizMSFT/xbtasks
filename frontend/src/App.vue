<script setup lang="ts">
import { onMounted } from 'vue'
import { useTheme } from './composables/useTheme'
import { useAuthStore } from './stores/auth'
import AppShell from './layouts/AppShell.vue'
import LoginView from './views/LoginView.vue'

useTheme()
const authStore = useAuthStore()

onMounted(async () => {
  // Try to restore session from existing token/keychain
  await authStore.tryRestore()
})
</script>

<template>
  <LoginView v-if="!authStore.isAuthenticated" />
  <AppShell v-else>
    <router-view />
  </AppShell>
</template>
