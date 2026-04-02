<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useTheme } from '@/composables/useTheme'
import { Save, RotateCcw, FolderOpen, Database, Globe, RefreshCw, Palette } from 'lucide-vue-next'

const { mode, toggle: toggleTheme } = useTheme()

interface ConfigState {
  theme: string
  dbPath: string
  adoOrganization: string
  adoProject: string
  syncIntervalMinutes: number
  logLevel: string
  windowWidth: number
  windowHeight: number
}

const config = ref<ConfigState>({
  theme: 'system',
  dbPath: '',
  adoOrganization: '',
  adoProject: '',
  syncIntervalMinutes: 15,
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
      config.value.adoOrganization = all.ado?.organization || ''
      config.value.adoProject = all.ado?.project || ''
      config.value.syncIntervalMinutes = all.sync?.interval_minutes || 15
      config.value.logLevel = all.log?.level || 'info'
      config.value.windowWidth = all.window?.width || 1200
      config.value.windowHeight = all.window?.height || 800
    }
  } catch {
    // Fall back to defaults when bindings unavailable
  } finally {
    loading.value = false
  }
})

async function saveConfig() {
  try {
    const { Set } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
    await Set('theme', config.value.theme)
    await Set('db.path', config.value.dbPath)
    await Set('ado.organization', config.value.adoOrganization)
    await Set('ado.project', config.value.adoProject)
    await Set('sync.interval_minutes', config.value.syncIntervalMinutes)
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
    adoOrganization: '',
    adoProject: '',
    syncIntervalMinutes: 15,
    logLevel: 'info',
    windowWidth: 1200,
    windowHeight: 800,
  }
}
</script>

<template>
  <ScrollArea class="flex-1">
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

      <!-- Azure DevOps -->
      <Card>
        <CardHeader>
          <div class="flex items-center gap-2">
            <Globe :size="16" class="text-muted-foreground" />
            <CardTitle class="text-sm">Azure DevOps</CardTitle>
          </div>
          <CardDescription>Connect to your ADO organization</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">Organization</label>
            <Input v-model="config.adoOrganization" placeholder="e.g. xbox" />
          </div>
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">Default Project</label>
            <Input v-model="config.adoProject" placeholder="e.g. xb-tasks" />
          </div>
          <Separator />
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-foreground">Sync Interval</p>
              <p class="text-xs text-muted-foreground">How often to sync ADO items (minutes)</p>
            </div>
            <Input
              v-model.number="config.syncIntervalMinutes"
              type="number"
              :min="1"
              :max="120"
              class="w-20"
            />
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
