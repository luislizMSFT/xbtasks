<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useTheme } from '@/composables/useTheme'
import { useAuthStore } from '@/stores/auth'
import { Save, RotateCcw, FolderOpen, Database, Globe, Palette, Trash2, Plus, RefreshCw, LogOut, KeyRound, Shield } from 'lucide-vue-next'

const router = useRouter()
const { mode, toggle: toggleTheme } = useTheme()
const authStore = useAuthStore()

// ---- Org/Project Management ----

interface OrgProject {
  org: string
  projects: string[]
}

const orgs = ref<OrgProject[]>([])
const newOrg = ref('')
const newProjects = ref('')
const orgsLoading = ref(false)
const orgsSaved = ref(false)

async function loadOrgs() {
  try {
    const { GetOrgProjects } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
    orgs.value = (await GetOrgProjects() as OrgProject[]) || []
  } catch {
    orgs.value = []
  }
}

async function saveOrgs() {
  orgsLoading.value = true
  try {
    const { SetOrgProjects } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
    await SetOrgProjects(orgs.value)
    orgsSaved.value = true
    setTimeout(() => { orgsSaved.value = false }, 2000)
  } catch {
    // Bindings unavailable
  } finally {
    orgsLoading.value = false
  }
}

function addOrg() {
  const org = newOrg.value.trim()
  if (!org) return
  const projects = newProjects.value
    .split(',')
    .map(p => p.trim())
    .filter(p => p.length > 0)
  if (projects.length === 0) return

  // Check if org already exists — merge projects
  const existing = orgs.value.find(o => o.org.toLowerCase() === org.toLowerCase())
  if (existing) {
    const newPs = projects.filter(p => !existing.projects.includes(p))
    existing.projects.push(...newPs)
  } else {
    orgs.value.push({ org, projects })
  }

  newOrg.value = ''
  newProjects.value = ''
  saveOrgs()
}

function removeOrg(index: number) {
  orgs.value.splice(index, 1)
  saveOrgs()
}

// ---- Sync Interval ----

const syncInterval = ref(15)

async function loadSyncInterval() {
  try {
    const { GetSyncInterval } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
    syncInterval.value = (await GetSyncInterval()) || 15
  } catch {
    syncInterval.value = 15
  }
}

async function saveSyncInterval() {
  try {
    const { SetSyncInterval } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
    await SetSyncInterval(syncInterval.value)
  } catch {
    // Bindings unavailable
  }
}

// ---- General Config ----

interface ConfigState {
  theme: string
  dbPath: string
  logLevel: string
  windowWidth: number
  windowHeight: number
}

const config = ref<ConfigState>({
  theme: 'system',
  dbPath: '',
  logLevel: 'info',
  windowWidth: 1200,
  windowHeight: 800,
})

const saved = ref(false)
const loading = ref(true)

onMounted(async () => {
  try {
    const { GetAll } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
    const all = await GetAll()
    if (all) {
      config.value.theme = all.theme || 'system'
      config.value.dbPath = all.db?.path || ''
      config.value.logLevel = all.log?.level || 'info'
      config.value.windowWidth = all.window?.width || 1200
      config.value.windowHeight = all.window?.height || 800
    }
  } catch {
    // Fall back to defaults when bindings unavailable
  } finally {
    loading.value = false
  }

  await loadOrgs()
  await loadSyncInterval()
})

async function saveConfig() {
  try {
    const { Set } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
    await Set('theme', config.value.theme)
    await Set('db.path', config.value.dbPath)
    await Set('log.level', config.value.logLevel)
    await Set('window.width', config.value.windowWidth)
    await Set('window.height', config.value.windowHeight)
    saved.value = true
    setTimeout(() => { saved.value = false }, 2000)
  } catch {
    // Bindings unavailable
  }
}

function resetDefaults() {
  config.value = {
    theme: 'system',
    dbPath: '',
    logLevel: 'info',
    windowWidth: 1200,
    windowHeight: 800,
  }
  syncInterval.value = 15
}

// ---- Auth ----

const patUpdateInput = ref('')
const showPatUpdate = ref(false)

async function handleSignOut() {
  await authStore.signOut()
  router.push('/')
}

async function updatePAT() {
  const token = patUpdateInput.value.trim()
  if (!token) return
  await authStore.signInWithPAT(token)
  patUpdateInput.value = ''
  showPatUpdate.value = false
}

const authMethodLabel = {
  oauth: 'Microsoft OAuth',
  pat: 'Personal Access Token',
  azcli: 'Az CLI Token',
} as const
</script>

<template>
  <ScrollArea class="flex-1 h-full">
    <div class="max-w-2xl mx-auto px-6 py-6 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold text-foreground">Settings</h1>
          <p class="text-sm text-muted-foreground mt-0.5">Configure your Team ADO Tool</p>
        </div>
        <div class="flex items-center gap-2">
          <Button variant="outline" size="sm" @click="resetDefaults">
            <RotateCcw :size="14" class="mr-1.5" />
            Reset
          </Button>
          <Button size="sm" @click="saveConfig">
            <Save :size="14" class="mr-1.5" />
            {{ saved ? 'Saved!' : 'Save' }}
          </Button>
        </div>
      </div>

      <!-- Azure DevOps Organizations -->
      <Card>
        <CardHeader>
          <div class="flex items-center gap-2">
            <Globe :size="16" class="text-muted-foreground" />
            <CardTitle class="text-sm">Azure DevOps Organizations</CardTitle>
          </div>
          <CardDescription>Configure which ADO orgs and projects to sync</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <!-- Org list -->
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

          <div v-else class="px-3 py-4 rounded-md border border-dashed border-border text-center">
            <p class="text-sm text-muted-foreground">No organizations configured</p>
            <p class="text-xs text-muted-foreground mt-0.5">Add an org below to start syncing</p>
          </div>

          <Separator />

          <!-- Add org form -->
          <div class="space-y-3">
            <p class="text-xs font-medium text-muted-foreground">Add Organization</p>
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
              Add Organization
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Sync Settings -->
      <Card>
        <CardHeader>
          <div class="flex items-center gap-2">
            <RefreshCw :size="16" class="text-muted-foreground" />
            <CardTitle class="text-sm">Sync Settings</CardTitle>
          </div>
          <CardDescription>Background sync frequency and behavior</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-foreground">Background sync interval</p>
              <p class="text-xs text-muted-foreground">How often to pull updates from ADO (minutes)</p>
            </div>
            <Input
              v-model.number="syncInterval"
              type="number"
              :min="1"
              :max="120"
              class="w-20"
              @change="saveSyncInterval"
            />
          </div>
        </CardContent>
      </Card>

      <!-- Authentication -->
      <Card>
        <CardHeader>
          <div class="flex items-center gap-2">
            <Shield :size="16" class="text-muted-foreground" />
            <CardTitle class="text-sm">Authentication</CardTitle>
          </div>
          <CardDescription>Manage your sign-in method</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-foreground">Current method</p>
              <p class="text-xs text-muted-foreground">
                {{ authStore.authMethod ? authMethodLabel[authStore.authMethod] : 'Not signed in' }}
              </p>
            </div>
            <Badge variant="secondary">
              {{ authStore.isAuthenticated ? 'Active' : 'Inactive' }}
            </Badge>
          </div>

          <div v-if="authStore.user" class="flex items-center gap-3 px-3 py-2 rounded-md border border-border bg-muted/30">
            <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-semibold text-primary">
              {{ authStore.initials }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-foreground truncate">{{ authStore.user.displayName }}</p>
              <p v-if="authStore.user.email" class="text-xs text-muted-foreground truncate">{{ authStore.user.email }}</p>
            </div>
          </div>

          <!-- PAT update -->
          <div v-if="authStore.authMethod === 'pat'" class="space-y-2">
            <Button
              v-if="!showPatUpdate"
              variant="outline"
              size="sm"
              @click="showPatUpdate = true"
            >
              <KeyRound :size="14" class="mr-1.5" />
              Update PAT
            </Button>
            <div v-if="showPatUpdate" class="flex gap-2">
              <Input
                v-model="patUpdateInput"
                type="password"
                placeholder="New personal access token"
                class="flex-1"
                @keydown.enter="updatePAT"
              />
              <Button size="sm" @click="updatePAT" :disabled="!patUpdateInput.trim()">Save</Button>
              <Button variant="outline" size="sm" @click="showPatUpdate = false">Cancel</Button>
            </div>
          </div>

          <Separator />

          <Button
            variant="destructive"
            size="sm"
            @click="handleSignOut"
          >
            <LogOut :size="14" class="mr-1.5" />
            Sign Out
          </Button>
        </CardContent>
      </Card>

      <!-- Appearance -->
      <Card>
        <CardHeader>
          <div class="flex items-center gap-2">
            <Palette :size="16" class="text-muted-foreground" />
            <CardTitle class="text-sm">Appearance</CardTitle>
          </div>
          <CardDescription>Theme and window preferences</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-muted-foreground">Theme</label>
              <Select v-model="config.theme">
                <SelectTrigger>
                  <SelectValue placeholder="Select theme" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="system">System</SelectItem>
                  <SelectItem value="light">Light</SelectItem>
                  <SelectItem value="dark">Dark</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-muted-foreground">Log Level</label>
              <Select v-model="config.logLevel">
                <SelectTrigger>
                  <SelectValue placeholder="Select level" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="debug">Debug</SelectItem>
                  <SelectItem value="info">Info</SelectItem>
                  <SelectItem value="warn">Warn</SelectItem>
                  <SelectItem value="error">Error</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <Separator />
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-muted-foreground">Window Width</label>
              <Input v-model.number="config.windowWidth" type="number" :min="800" :max="3840" />
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-muted-foreground">Window Height</label>
              <Input v-model.number="config.windowHeight" type="number" :min="600" :max="2160" />
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Data -->
      <Card>
        <CardHeader>
          <div class="flex items-center gap-2">
            <Database :size="16" class="text-muted-foreground" />
            <CardTitle class="text-sm">Data</CardTitle>
          </div>
          <CardDescription>Database location and storage</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">Database Path</label>
            <div class="flex items-center gap-2">
              <Input v-model="config.dbPath" placeholder="Auto-detected" class="flex-1" readonly />
              <Button variant="outline" size="icon" class="shrink-0">
                <FolderOpen :size="14" />
              </Button>
            </div>
            <p class="text-[11px] text-muted-foreground">
              Leave empty for default OS location
            </p>
          </div>
        </CardContent>
      </Card>

      <!-- About -->
      <Card>
        <CardHeader>
          <CardTitle class="text-sm">About</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-foreground">Team ADO Tool</p>
              <p class="text-xs text-muted-foreground">Xbox Services team productivity dashboard</p>
            </div>
            <Badge variant="secondary">v0.1.0</Badge>
          </div>
        </CardContent>
      </Card>
    </div>
  </ScrollArea>
</template>
