import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface User {
  id: string
  displayName: string
  email: string
  avatarUrl: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)
  const initials = computed(() => {
    if (!user.value) return '?'
    const parts = user.value.displayName.split(' ')
    return parts.map(p => p[0]).join('').toUpperCase().slice(0, 2)
  })

  async function signIn() {
    loading.value = true
    error.value = null
    try {
      // Try Wails binding first, fall back to mock
      try {
        const { SignIn } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
        const u = await SignIn()
        user.value = u as User
      } catch (e) {
        console.warn('[AuthStore] Wails binding unavailable, using mock auth:', e)
        user.value = {
          id: 'mock-user-1',
          displayName: 'Dev User',
          email: 'dev@example.com',
          avatarUrl: '',
        }
      }
    } catch (e: any) {
      error.value = e.message || 'Sign in failed'
    } finally {
      loading.value = false
    }
  }

  async function signOut() {
    try {
      const { SignOut } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
      await SignOut()
    } catch (e) {
      console.warn('[AuthStore] Wails signOut binding unavailable:', e)
    }
    user.value = null
  }

  async function tryRestore() {
    try {
      const { GetCurrentUser } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
      const u = await GetCurrentUser()
      if (u) user.value = u as User
    } catch (e) {
      console.warn('[AuthStore] Wails restore binding unavailable:', e)
    }
  }

  return { user, loading, error, isAuthenticated, initials, signIn, signOut, tryRestore }
})
