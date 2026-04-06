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
  const authMethod = ref<'oauth' | 'pat' | 'azcli' | null>(null)

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
        authMethod.value = 'oauth'
      } catch (e) {
        console.warn('[AuthStore] Wails binding unavailable, using mock auth:', e)
        user.value = {
          id: 'mock-user-1',
          displayName: 'Dev User',
          email: 'dev@example.com',
          avatarUrl: '',
        }
        authMethod.value = 'oauth'
      }
    } catch (e: any) {
      error.value = e.message || 'Sign in failed'
    } finally {
      loading.value = false
    }
  }

  async function signInWithPAT(pat: string) {
    loading.value = true
    error.value = null
    try {
      try {
        const { SignInWithPAT } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
        const u = await SignInWithPAT(pat)
        user.value = u as User
        authMethod.value = 'pat'
      } catch (e) {
        console.warn('[AuthStore] Wails PAT binding unavailable, using mock:', e)
        user.value = {
          id: 'pat-user-1',
          displayName: 'PAT User',
          email: '',
          avatarUrl: '',
        }
        authMethod.value = 'pat'
      }
    } catch (e: any) {
      error.value = e.message || 'PAT sign in failed'
    } finally {
      loading.value = false
    }
  }

  async function signInWithAzCli() {
    loading.value = true
    error.value = null
    try {
      try {
        const { SignInWithAzCli } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
        const u = await SignInWithAzCli()
        user.value = u as User
        authMethod.value = 'azcli'
      } catch (e) {
        console.warn('[AuthStore] Wails AzCli binding unavailable, using mock:', e)
        user.value = {
          id: 'azcli-user-1',
          displayName: 'Az CLI User',
          email: 'azcli@local',
          avatarUrl: '',
        }
        authMethod.value = 'azcli'
      }
    } catch (e: any) {
      error.value = e.message || 'Az CLI sign in failed — run "az login" first'
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
    authMethod.value = null
  }

  async function tryRestore() {
    try {
      const { TryRestoreSession } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
      const u = await TryRestoreSession()
      if (u) user.value = u as User
    } catch (e) {
      // TryRestoreSession not available or failed, try GetCurrentUser as fallback
      try {
        const { GetCurrentUser } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
        const u = await GetCurrentUser()
        if (u) user.value = u as User
      } catch {
        // No auth available
      }
    }
  }

  return { user, loading, error, authMethod, isAuthenticated, initials, signIn, signInWithPAT, signInWithAzCli, signOut, tryRestore }
})
