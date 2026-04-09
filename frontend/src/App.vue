<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useTheme } from './composables/useTheme'
import { useAuthStore } from './stores/auth'
import { useTaskStore } from './stores/tasks'
import { useADOStore } from './stores/ado'
import { usePRStore } from './stores/prs'
import { useProjectStore } from './stores/projects'
import { useSyncStore } from './stores/sync'
import { Toaster } from 'vue-sonner'
import AppShell from './layouts/AppShell.vue'
import LoginView from './views/LoginView.vue'
import OnboardingView from './views/OnboardingView.vue'

useTheme()
const authStore = useAuthStore()
const taskStore = useTaskStore()
const adoStore = useADOStore()
const prStore = usePRStore()
const projectStore = useProjectStore()
const syncStore = useSyncStore()

const needsOnboarding = ref(false)
const onboardingChecked = ref(false)

onMounted(async () => {
  try {
    const { getOrgProjects } = await import('./api/config')
    const orgs = await getOrgProjects()
    needsOnboarding.value = !orgs || orgs.length === 0
  } catch {
    needsOnboarding.value = false
  }
  onboardingChecked.value = true

  await authStore.tryRestore()
})

function onOnboardingComplete() {
  needsOnboarding.value = false
}

// When auth state changes to authenticated, prefetch all core data in parallel
watch(() => authStore.isAuthenticated, (authed) => {
  if (authed) {
    // Initialize Wails event listeners for backend→frontend communication
    syncStore.initEvents()

    // Fire all fetches in parallel — don't block on any single one
    Promise.allSettled([
      taskStore.fetchTasks(),
      projectStore.fetchProjects().then(() => projectStore.fetchAllProjectLinks()),
      prStore.fetchAll(),
      adoStore.fetchWorkItemTree(),
      adoStore.fetchLinkedAdoIds(),
      adoStore.fetchPipelines(),
      adoStore.fetchSavedQueries(),
    ])
  }
}, { immediate: true })
</script>

<template>
  <template v-if="onboardingChecked">
    <OnboardingView v-if="needsOnboarding" @complete="onOnboardingComplete" />
    <LoginView v-else-if="!authStore.isAuthenticated" />
    <AppShell v-else>
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </AppShell>
  </template>
  <Toaster position="bottom-right" :theme="'system'" richColors />
</template>
