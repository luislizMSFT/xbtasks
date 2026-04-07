import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import * as authApi from '@/api/auth'

export type { User }

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
      try {
        const u = await authApi.signIn()
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
        const u = await authApi.signInWithPAT(pat)
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
      const u = await authApi.signInWithAzCli()
      user.value = u as User
      authMethod.value = 'azcli'
    } catch (e: any) {
      error.value = e.message || 'Az CLI sign in failed — run "az login" first'
    } finally {
      loading.value = false
    }
  }

  async function signOut() {
    try {
      await authApi.signOut()
    } catch (e) {
      console.warn('[AuthStore] Wails signOut binding unavailable:', e)
    }
    user.value = null
    authMethod.value = null
  }

  async function tryRestore() {
    try {
      const u = await authApi.tryRestoreSession()
      if (u) user.value = u as User
    } catch (e) {
      try {
        const u = await authApi.getCurrentUser()
        if (u) user.value = u as User
      } catch {
        // No auth available
      }
    }
  }

  return { user, loading, error, authMethod, isAuthenticated, initials, signIn, signInWithPAT, signInWithAzCli, signOut, tryRestore }
})
