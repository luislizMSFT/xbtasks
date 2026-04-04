<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
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
  // TODO: Call AuthService.SignInWithPAT(token) when binding is ready
  await authStore.signIn()
  if (authStore.isAuthenticated) {
    router.push('/tasks')
  }
}
</script>

<template>
  <div class="h-screen w-screen flex items-center justify-center bg-background">
    <Card class="w-full max-w-sm">
      <CardHeader class="text-center">
        <div class="flex justify-center mb-2">
          <div class="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center">
            <svg class="w-8 h-8 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 11l3 3L22 4" />
              <path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11" />
            </svg>
          </div>
        </div>
        <CardTitle class="text-[20px] font-semibold">Team ADO Tool</CardTitle>
        <CardDescription class="text-[14px] font-normal text-muted-foreground">All your work in one place</CardDescription>
      </CardHeader>

      <CardContent class="space-y-4">
        <!-- Error -->
        <div v-if="authStore.error" class="px-3 py-2 rounded-md text-[12px] font-semibold bg-destructive/10 text-destructive border border-destructive/20">
          Authentication failed — check your network connection and try again.
        </div>

        <!-- Sign in buttons -->
        <div class="space-y-3">
          <Button
            @click="signInWithMicrosoft"
            :disabled="authStore.loading"
            class="w-full"
            size="lg"
          >
            <Loader2 v-if="authStore.loading" :size="16" class="animate-spin" />
            <template v-if="authStore.loading">Signing in...</template>
            <template v-else>
              <svg class="w-4 h-4" viewBox="0 0 21 21" fill="none">
                <rect width="10" height="10" fill="#F25022"/>
                <rect x="11" width="10" height="10" fill="#7FBA00"/>
                <rect y="11" width="10" height="10" fill="#00A4EF"/>
                <rect x="11" y="11" width="10" height="10" fill="#FFB900"/>
              </svg>
              Sign in with Microsoft
            </template>
          </Button>

          <Button
            v-if="!showPAT"
            @click="showPAT = true"
            variant="ghost"
            class="w-full"
            size="lg"
            :disabled="authStore.loading"
          >
            <KeyRound :size="16" />
            Use Personal Access Token
          </Button>

          <!-- PAT input -->
          <Transition
            enter-active-class="transition duration-150 ease-out"
            enter-from-class="opacity-0 -translate-y-1"
            enter-to-class="opacity-100 translate-y-0"
          >
            <div v-if="showPAT" class="space-y-2">
              <Input
                v-model="patInput"
                type="password"
                placeholder="Paste your personal access token"
                class="h-10"
                @keydown.enter="signInWithPAT"
                autofocus
              />
              <div class="flex gap-2">
                <Button
                  @click="showPAT = false"
                  variant="outline"
                  size="sm"
                  class="flex-1"
                >
                  Cancel
                </Button>
                <Button
                  @click="signInWithPAT"
                  :disabled="authStore.loading || !patInput.trim()"
                  size="sm"
                  class="flex-1"
                >
                  <Loader2 v-if="authStore.loading" :size="14" class="animate-spin" />
                  <template v-if="authStore.loading">Signing in...</template>
                  <template v-else>Connect</template>
                </Button>
              </div>
            </div>
          </Transition>
        </div>

        <!-- Footer -->
        <p class="text-center text-[11px] text-muted-foreground">
          v0.1.0 · Press <kbd class="px-1 py-0.5 rounded text-[10px] bg-muted border border-border">⌘K</kbd> anytime for commands
        </p>
      </CardContent>
    </Card>
  </div>
</template>
