<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useTheme } from './composables/useTheme'
import { useAuthStore } from './stores/auth'
import { useTaskStore } from './stores/tasks'
import { useADOStore } from './stores/ado'
import { usePRStore } from './stores/prs'
import { useProjectStore } from './stores/projects'
import AppShell from './layouts/AppShell.vue'
import LoginView from './views/LoginView.vue'

useTheme()
const authStore = useAuthStore()
const taskStore = useTaskStore()
const adoStore = useADOStore()
const prStore = usePRStore()
const projectStore = useProjectStore()

onMounted(async () => {
  await authStore.tryRestore()
})

// When auth state changes to authenticated, prefetch all core data
watch(() => authStore.isAuthenticated, (authed) => {
  if (authed) {
    taskStore.fetchTasks()
    projectStore.fetchProjects()
    prStore.fetchAll()
    adoStore.fetchWorkItemTree()
    adoStore.fetchLinkedAdoIds()
    adoStore.fetchPipelines()
    adoStore.fetchSavedQueries()
  }
}, { immediate: true })
</script>

<template>
  <LoginView v-if="!authStore.isAuthenticated" />
  <AppShell v-else>
    <router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component" />
      </keep-alive>
    </router-view>
  </AppShell>
</template>
