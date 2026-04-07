<script setup lang="ts">
import { ref, computed } from 'vue'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Globe, Plus, Trash2, ArrowRight, CheckCircle } from 'lucide-vue-next'

interface OrgProject {
  org: string
  projects: string[]
}

const emit = defineEmits<{ complete: [] }>()

const step = ref<'orgs' | 'done'>('orgs')
const orgs = ref<OrgProject[]>([])
const newOrg = ref('')
const newProjects = ref('')
const saving = ref(false)

const canContinue = computed(() => orgs.value.length > 0)

function addOrg() {
  const org = newOrg.value.trim()
  if (!org) return
  const projects = newProjects.value
    .split(',')
    .map(p => p.trim())
    .filter(p => p.length > 0)
  if (projects.length === 0) return

  const existing = orgs.value.find(o => o.org.toLowerCase() === org.toLowerCase())
  if (existing) {
    const newPs = projects.filter(p => !existing.projects.includes(p))
    existing.projects.push(...newPs)
  } else {
    orgs.value.push({ org, projects })
  }

  newOrg.value = ''
  newProjects.value = ''
}

function removeOrg(index: number) {
  orgs.value.splice(index, 1)
}

async function finish() {
  saving.value = true
  try {
    const { setOrgProjects } = await import('@/api/config')
    await setOrgProjects(orgs.value)
    step.value = 'done'
    setTimeout(() => emit('complete'), 800)
  } catch {
    // Bindings unavailable — continue anyway
    emit('complete')
  } finally {
    saving.value = false
  }
}

function skip() {
  emit('complete')
}
</script>

<template>
  <div class="h-screen w-screen flex items-center justify-center bg-background">
    <!-- Step 1: Configure Orgs -->
    <Card v-if="step === 'orgs'" class="w-full max-w-md">
      <CardHeader class="text-center">
        <div class="flex justify-center mb-2">
          <div class="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center">
            <svg class="w-8 h-8 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 11l3 3L22 4" />
              <path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11" />
            </svg>
          </div>
        </div>
        <CardTitle class="text-[20px] font-semibold">Welcome to XB Tasks</CardTitle>
        <CardDescription class="text-[14px]">
          Let's connect your Azure DevOps organizations
        </CardDescription>
      </CardHeader>

      <CardContent class="space-y-4">
        <!-- Configured orgs -->
        <div v-if="orgs.length > 0" class="space-y-2">
          <div
            v-for="(orgItem, index) in orgs"
            :key="orgItem.org"
            class="flex items-center justify-between px-3 py-2 rounded-md border border-border bg-muted/30"
          >
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-foreground truncate">{{ orgItem.org }}</p>
              <p class="text-xs text-muted-foreground truncate">{{ orgItem.projects.join(', ') }}</p>
            </div>
            <Button
              variant="ghost"
              size="icon"
              class="shrink-0 h-8 w-8 text-muted-foreground hover:text-destructive"
              @click="removeOrg(index)"
            >
              <Trash2 :size="14" />
            </Button>
          </div>
        </div>

        <Separator v-if="orgs.length > 0" />

        <!-- Add org form -->
        <div class="space-y-3">
          <div class="flex items-center gap-2">
            <Globe :size="14" class="text-muted-foreground" />
            <p class="text-xs font-medium text-muted-foreground">Add Organization</p>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <label class="text-xs text-muted-foreground">Organization</label>
              <Input
                v-model="newOrg"
                placeholder="e.g., xbox"
                @keydown.enter="addOrg"
              />
            </div>
            <div class="space-y-1">
              <label class="text-xs text-muted-foreground">Projects</label>
              <Input
                v-model="newProjects"
                placeholder="e.g., XES, XboxLive"
                @keydown.enter="addOrg"
              />
            </div>
          </div>
          <p class="text-[11px] text-muted-foreground">Separate multiple projects with commas</p>
          <Button
            variant="outline"
            size="sm"
            @click="addOrg"
            :disabled="!newOrg.trim() || !newProjects.trim()"
          >
            <Plus :size="14" class="mr-1.5" />
            Add
          </Button>
        </div>

        <Separator />

        <!-- Actions -->
        <div class="flex items-center justify-between">
          <Button variant="ghost" size="sm" @click="skip" class="text-muted-foreground">
            Skip for now
          </Button>
          <Button
            size="sm"
            @click="finish"
            :disabled="!canContinue || saving"
          >
            Continue
            <ArrowRight :size="14" class="ml-1.5" />
          </Button>
        </div>

        <p class="text-center text-[11px] text-muted-foreground">
          You can change this later in Settings
        </p>
      </CardContent>
    </Card>

    <!-- Step 2: Done -->
    <Card v-else class="w-full max-w-sm text-center">
      <CardContent class="py-8 space-y-3">
        <CheckCircle :size="40" class="mx-auto text-green-500" />
        <p class="text-sm font-medium text-foreground">You're all set!</p>
        <div class="flex flex-wrap justify-center gap-1.5">
          <Badge v-for="org in orgs" :key="org.org" variant="secondary">
            {{ org.org }} · {{ org.projects.length }} project{{ org.projects.length > 1 ? 's' : '' }}
          </Badge>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
