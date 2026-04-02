<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { Loader2, KeyRound } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const showPAT = ref(false)
const patInput = ref('')

async function signInWithMicrosoft() {
  await authStore.signIn()
  if (authStore.isAuthenticated) {
    router.push('/tasks')
  }
}

async function signInWithPAT() {
  const token = patInput.value.trim()
  if (!token) return
  // For now, treat PAT as mock auth
  authStore.signIn()
  if (authStore.isAuthenticated) {
    router.push('/tasks')
  }
}
</script>

<template>
  <div class="h-screen w-screen flex items-center justify-center bg-surface-primary">
    <div class="w-full max-w-sm text-center space-y-6 px-6">
      <!-- App icon -->
      <div class="flex flex-col items-center gap-3">
        <div class="w-16 h-16 rounded-2xl bg-accent/10 flex items-center justify-center">
          <svg class="w-8 h-8 text-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 11l3 3L22 4" />
            <path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11" />
          </svg>
        </div>
        <div>
          <h1 class="text-xl font-bold text-text-primary">Team ADO Tool</h1>
          <p class="text-sm text-text-secondary mt-1">Tasks · ADO · PRs — one pane of glass</p>
        </div>
      </div>

      <!-- Error -->
      <div v-if="authStore.error" class="px-3 py-2 rounded-md text-sm bg-red-500/10 text-red-600 dark:text-red-400 border border-red-500/20">
        {{ authStore.error }}
      </div>

      <!-- Sign in buttons -->
      <div class="space-y-3">
        <button
          @click="signInWithMicrosoft"
          :disabled="authStore.loading"
          class="w-full flex items-center justify-center gap-2 h-10 rounded-lg text-sm font-medium bg-accent text-white hover:bg-accent/90 disabled:opacity-50 transition-colors"
        >
          <Loader2 v-if="authStore.loading" :size="16" class="animate-spin" />
          <template v-else>
            <svg class="w-4 h-4" viewBox="0 0 21 21" fill="none">
              <rect width="10" height="10" fill="#F25022"/>
              <rect x="11" width="10" height="10" fill="#7FBA00"/>
              <rect y="11" width="10" height="10" fill="#00A4EF"/>
              <rect x="11" y="11" width="10" height="10" fill="#FFB900"/>
            </svg>
            Sign in with Microsoft
          </template>
        </button>

        <button
          v-if="!showPAT"
          @click="showPAT = true"
          class="w-full flex items-center justify-center gap-2 h-10 rounded-lg text-sm font-medium text-text-secondary hover:text-text-primary border border-border-default hover:bg-surface-tertiary transition-colors"
        >
          <KeyRound :size="16" />
          Use Personal Access Token
        </button>

        <!-- PAT input -->
        <Transition
          enter-active-class="transition duration-150 ease-out"
          enter-from-class="opacity-0 -translate-y-1"
          enter-to-class="opacity-100 translate-y-0"
        >
          <div v-if="showPAT" class="space-y-2">
            <input
              v-model="patInput"
              type="password"
              class="w-full h-10 px-3 rounded-lg text-sm bg-surface-tertiary border border-border-default text-text-primary placeholder-text-secondary outline-none focus:ring-1 focus:ring-accent"
              placeholder="Paste your personal access token"
              @keydown.enter="signInWithPAT"
              autofocus
            />
            <div class="flex gap-2">
              <button
                @click="showPAT = false"
                class="flex-1 h-9 rounded-lg text-xs text-text-secondary hover:text-text-primary hover:bg-surface-tertiary border border-border-default transition-colors"
              >
                Cancel
              </button>
              <button
                @click="signInWithPAT"
                class="flex-1 h-9 rounded-lg text-xs font-medium bg-accent text-white hover:bg-accent/90 transition-colors"
              >
                Sign In
              </button>
            </div>
          </div>
        </Transition>
      </div>

      <!-- Footer -->
      <p class="text-[11px] text-text-secondary">
        v0.1.0 · Press <kbd class="px-1 py-0.5 rounded text-[10px] bg-surface-tertiary border border-border-default">⌘K</kbd> anytime for commands
      </p>
    </div>
  </div>
</template>
