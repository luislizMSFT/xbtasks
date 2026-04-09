<script setup lang="ts">
/**
 * Playground: Wiki — Obsidian/Notion hybrid
 *
 * Dashboard → Areas → Pages flow.
 * Dashboard shows saved links + area cards. Click an area to browse its pages.
 * Full-width editor. Mock data — no backend.
 */
import { ref, computed, watch, nextTick } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import {
  FileText, BookOpen, Search, Plus, Trash2, ExternalLink,
  ChevronRight, ArrowLeft, Pencil, Globe, Star, Folder, X,
} from 'lucide-vue-next'

// ── Types ──
interface SavedLink {
  id: number
  title: string
  url: string
  description?: string
  icon?: string
  pinned?: boolean
}

interface WikiPage {
  id: number
  areaId: number
  title: string
  content: string
  updatedAt: string
}

interface Area {
  id: number
  name: string
  description: string
  icon: string
  color: string
  pageCount: number
}

type View = 'dashboard' | 'links' | 'area' | 'page'

// ── Mock Data ──
const mockLinks: SavedLink[] = [
  { id: 1, title: 'Grafana — Service Health', url: 'https://grafana.internal/d/service-health', description: 'CPU, memory, request rate, error rate', pinned: true },
  { id: 2, title: 'ADO Pipelines — xbox-services', url: 'https://dev.azure.com/xbox-services/_build', description: 'CI/CD pipelines for all services', pinned: true },
  { id: 3, title: 'Azure Portal — Resource Group', url: 'https://portal.azure.com/#view/HubsExtension/BrowseResourceGroupBlade', description: 'Xbox Services resource group', pinned: true },
  { id: 4, title: 'ICM — Incident Manager', url: 'https://icm.ad.msft.net', description: 'Create and track incidents', pinned: false },
  { id: 5, title: 'Kusto Explorer — Telemetry', url: 'https://dataexplorer.azure.com/clusters/xboxservices', description: 'Query telemetry and diagnostics logs', pinned: false },
  { id: 6, title: 'LaunchDarkly — Feature Flags', url: 'https://app.launchdarkly.com/xbox-services', description: 'Manage feature rollouts', pinned: false },
  { id: 7, title: 'Figma — Design System', url: 'https://figma.com/file/xbox-design-system', description: 'UI component library and design tokens', pinned: false },
  { id: 8, title: 'Confluence — Team Handbook', url: 'https://confluence.internal/xbox-services/handbook', description: 'Team processes, OKRs, meeting notes', pinned: false },
  { id: 9, title: 'GitHub — mono-repo', url: 'https://github.com/xbox-services/platform', description: 'Main service repository', pinned: false },
  { id: 10, title: 'PagerDuty — On-Call', url: 'https://pagerduty.com/schedules/xbox', description: 'On-call rotation schedule', pinned: false },
]

const mockAreas: Area[] = [
  { id: 1, name: 'On-Call & Runbooks', description: 'Service runbooks, alert procedures, escalation paths', icon: '🚨', color: 'text-red-500 bg-red-500/10 border-red-500/20', pageCount: 3 },
  { id: 2, name: 'Team & Onboarding', description: 'New hire guides, dev environment setup, architecture overview', icon: '👥', color: 'text-blue-500 bg-blue-500/10 border-blue-500/20', pageCount: 2 },
  { id: 3, name: 'Infrastructure', description: 'Deployment playbooks, monitoring, infra-as-code', icon: '⚙️', color: 'text-amber-500 bg-amber-500/10 border-amber-500/20', pageCount: 2 },
  { id: 4, name: 'ADO Wikis', description: 'Imported pages from Azure DevOps project wikis', icon: '📘', color: 'text-violet-500 bg-violet-500/10 border-violet-500/20', pageCount: 2 },
]

const mockPages: WikiPage[] = [
  // Area 1: On-Call & Runbooks
  {
    id: 1, areaId: 1, title: 'Service X Runbook',
    content: `# Service X Runbook\n\n## Overview\n\nService X handles authentication and token management for the Xbox Services platform.\n\n## On-Call Procedures\n\n### High CPU Alert\n\n1. Check the **Grafana dashboard** for CPU trends\n2. Look at recent deployments in ADO pipelines\n3. If caused by traffic spike → scale up via \`kubectl scale\`\n4. If caused by code regression → rollback latest deploy\n\n### Database Connection Failures\n\n- Check the connection pool metrics\n- Verify the SQL server is healthy via Azure Portal\n- Look for lock contention in \`sys.dm_exec_requests\`\n\n## Key Links\n\n- [[Team Onboarding]] — getting started guide\n- [[Monitoring Setup]] — Grafana dashboards and alerts\n\n## Contacts\n\n| Role | Person | Alias |\n|------|--------|-------|\n| Primary On-Call | Luis | luisliz |\n| Backup | Sarah | sarahc |\n| Escalation | Team Lead | teamlead |\n\n> **Note:** Always update the incident timeline in ICM before escalating.\n\n\`\`\`bash\n# Quick health check\ncurl -s https://service-x.internal/health | jq .\n\`\`\``,
    updatedAt: '2 hours ago',
  },
  {
    id: 2, areaId: 1, title: 'Monitoring Setup',
    content: `# Monitoring Setup\n\n## Grafana Dashboards\n\nWe maintain three primary dashboards:\n\n1. **Service Health** — CPU, memory, request rate, error rate\n2. **Database Performance** — query latency, connection pool, deadlocks\n3. **Pipeline Status** — deploy frequency, failure rate, MTTR\n\n## Alerting Rules\n\n| Alert | Threshold | Severity | Action |\n|-------|-----------|----------|--------|\n| High CPU | > 80% for 5m | P2 | Check \`kubectl top pods\` |\n| Error Rate | > 5% for 2m | P1 | Check logs, consider rollback |\n| DB Latency | > 500ms p95 | P2 | Check connection pool |\n| Deploy Failed | Any failure | P3 | Check pipeline logs |\n\n> See [[Service X Runbook]] for on-call procedures when alerts fire.`,
    updatedAt: '3 days ago',
  },
  {
    id: 3, areaId: 1, title: 'Escalation Matrix',
    content: `# Escalation Matrix\n\n## Severity Levels\n\n| Level | Response Time | Example |\n|-------|---------------|----------|\n| P0 | 15 min | Service down, data loss |\n| P1 | 30 min | Degraded performance, partial outage |\n| P2 | 4 hours | Non-critical bug, workaround exists |\n| P3 | Next business day | Feature request, minor UI issue |\n\n## Escalation Paths\n\n**P0/P1:**\n1. On-call engineer responds\n2. If unresolved in 30 min → page backup\n3. If unresolved in 1 hour → engage Team Lead\n4. If unresolved in 2 hours → engage Director`,
    updatedAt: '1 week ago',
  },
  // Area 2: Team & Onboarding
  {
    id: 4, areaId: 2, title: 'Team Onboarding',
    content: `# Team Onboarding\n\nWelcome to the Xbox Services team!\n\n## Day 1\n\n- [ ] Get access to ADO org \`xbox-services\`\n- [ ] Clone the mono-repo\n- [ ] Install Go 1.22+, Node 20+, and Wails v3\n- [ ] Run \`task setup\` to bootstrap local environment\n\n## Week 1\n\n- [ ] Read [[Service X Runbook]] and [[Monitoring Setup]]\n- [ ] Shadow an on-call rotation\n- [ ] Submit your first PR (any size)\n\n## Architecture\n\nThe app is a **Wails v3** desktop application:\n- **Backend:** Go with SQLite (local-first)\n- **Frontend:** Vue 3 + shadcn-vue\n- **Auth:** Azure CLI token provider`,
    updatedAt: '1 day ago',
  },
  {
    id: 5, areaId: 2, title: 'Dev Environment Setup',
    content: `# Dev Environment Setup\n\n\`\`\`bash\n# Install dependencies\ngo mod download\ncd frontend && npm install\n\n# Run the app\nwails3 dev\n\`\`\`\n\n## Required Tools\n\n- Go 1.22+\n- Node.js 20+\n- Wails v3 CLI\n- Azure CLI (for auth)\n- jj (Jujutsu) for version control`,
    updatedAt: '5 days ago',
  },
  // Area 3: Infrastructure
  {
    id: 6, areaId: 3, title: 'Deployment Playbook',
    content: `# Deployment Playbook\n\n## Standard Deploy Process\n\n1. Merge PR to \`main\`\n2. CI pipeline builds and runs tests\n3. CD pipeline deploys to staging\n4. Smoke tests run automatically\n5. Manual approval gate for production\n6. Rolling deploy to production (canary → 25% → 50% → 100%)\n\n## Rollback\n\n\`\`\`bash\n# Emergency rollback\naz pipelines run --name "rollback-prod" --branch main\n\`\`\`\n\n## Hotfix Process\n\nFor P0/P1 issues:\n1. Create branch from \`release/current\`\n2. Cherry-pick fix\n3. Fast-track through pipeline (skip staging soak)`,
    updatedAt: '2 weeks ago',
  },
  {
    id: 7, areaId: 3, title: 'Infrastructure Overview',
    content: `# Infrastructure Overview\n\n## Azure Resources\n\n- **3 regions:** West US 2, East US, West Europe\n- **AKS clusters:** 1 per region, auto-scaling 3-20 nodes\n- **CosmosDB:** Multi-region write, strong consistency\n- **Service Bus:** Premium tier, geo-DR enabled\n\n## Cost\n\nMonthly run rate: ~$45K\n- Compute: 40%\n- Database: 35%\n- Networking: 15%\n- Other: 10%`,
    updatedAt: '1 week ago',
  },
  // Area 4: ADO Wikis (imported)
  {
    id: 8, areaId: 4, title: 'Xbox Services Architecture',
    content: `# Xbox Services Architecture\n\n> **Source:** Azure DevOps Wiki — xbox-services / Platform Team\n\n## Service Topology\n\nThe platform consists of 12 microservices deployed across 3 Azure regions...\n\n## Authentication Flow\n\n1. Client sends request with Bearer token\n2. API Gateway validates token via Entra ID\n3. Service mesh handles mTLS between services\n4. Each service has its own managed identity`,
    updatedAt: '1 week ago',
  },
  {
    id: 9, areaId: 4, title: 'API Reference',
    content: `# API Reference\n\n> **Source:** Azure DevOps Wiki — xbox-services / Platform Team\n\n## Auth Service\n\n### POST /api/v1/token\n\nExchange credentials for access token.\n\n\`\`\`json\n{\n  "grant_type": "client_credentials",\n  "client_id": "...",\n  "scope": "xbox.services.read"\n}\n\`\`\`\n\n### GET /api/v1/health\n\nHealth check endpoint. Returns 200 if service is healthy.`,
    updatedAt: '2 weeks ago',
  },
]

// ── State ──
const currentView = ref<View>('dashboard')
const links = ref<SavedLink[]>([...mockLinks])
const areas = ref<Area[]>([...mockAreas])
const pages = ref<WikiPage[]>([...mockPages])
const selectedAreaId = ref<number | null>(null)
const selectedPageId = ref<number | null>(null)
const mode = ref<'edit' | 'preview'>('preview')
const searchQuery = ref('')
const editingTitle = ref(false)
const titleInput = ref('')
const showDeleteDialog = ref(false)
const showAddLink = ref(false)
const newLink = ref({ title: '', url: '', description: '' })

// ── Computed ──
const selectedArea = computed(() => areas.value.find(a => a.id === selectedAreaId.value) ?? null)
const selectedPage = computed(() => pages.value.find(p => p.id === selectedPageId.value) ?? null)
const areaPages = computed(() => selectedAreaId.value ? pages.value.filter(p => p.areaId === selectedAreaId.value) : [])
const pinnedLinks = computed(() => links.value.filter(l => l.pinned))
const otherLinks = computed(() => links.value.filter(l => !l.pinned))

const filteredLinks = computed(() => {
  if (!searchQuery.value.trim()) return links.value
  const q = searchQuery.value.toLowerCase()
  return links.value.filter(l =>
    l.title.toLowerCase().includes(q) ||
    l.url.toLowerCase().includes(q) ||
    (l.description?.toLowerCase().includes(q) ?? false),
  )
})

// ── Markdown rendering ──
function renderMarkdown(content: string): string {
  const withWikiLinks = content.replace(/\[\[([^\]]+)\]\]/g, (_match, pageName: string) => {
    const exists = pages.value.some(p => p.title === pageName)
    const cls = exists ? 'wiki-link' : 'wiki-link wiki-link-missing'
    return `<a href="#" class="${cls}" data-wiki-page="${pageName}">${pageName}</a>`
  })
  const html = marked.parse(withWikiLinks, { async: false }) as string
  return DOMPurify.sanitize(html)
}

const renderedContent = computed(() => {
  if (!selectedPage.value) return ''
  return renderMarkdown(selectedPage.value.content)
})

// ── Navigation ──
function goToDashboard() {
  currentView.value = 'dashboard'
  selectedAreaId.value = null
  selectedPageId.value = null
  searchQuery.value = ''
}

function goToLinks() {
  currentView.value = 'links'
  searchQuery.value = ''
}

function goToArea(areaId: number) {
  selectedAreaId.value = areaId
  selectedPageId.value = null
  currentView.value = 'area'
}

function goToPage(pageId: number) {
  selectedPageId.value = pageId
  mode.value = 'preview'
  currentView.value = 'page'
}

function handleWikiLinkClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (target.classList.contains('wiki-link')) {
    e.preventDefault()
    const pageName = target.getAttribute('data-wiki-page')
    if (!pageName) return
    const existing = pages.value.find(p => p.title === pageName)
    if (existing) {
      goToArea(existing.areaId)
      nextTick(() => goToPage(existing.id))
    }
  }
}

// ── Page Actions ──
function createPage() {
  if (!selectedAreaId.value) return
  const newId = Math.max(0, ...pages.value.map(p => p.id)) + 1
  const newPage: WikiPage = {
    id: newId,
    areaId: selectedAreaId.value,
    title: 'Untitled',
    content: '',
    updatedAt: 'just now',
  }
  pages.value.push(newPage)
  goToPage(newId)
  mode.value = 'edit'
  nextTick(() => {
    editingTitle.value = true
    titleInput.value = 'Untitled'
  })
}

function startEditTitle() {
  editingTitle.value = true
  titleInput.value = selectedPage.value?.title ?? ''
}

function saveTitle() {
  if (selectedPage.value && titleInput.value.trim()) {
    selectedPage.value.title = titleInput.value.trim()
  }
  editingTitle.value = false
}

function deletePage() {
  if (!selectedPage.value) return
  const areaId = selectedPage.value.areaId
  pages.value = pages.value.filter(p => p.id !== selectedPage.value!.id)
  showDeleteDialog.value = false
  goToArea(areaId)
}

// ── Link Actions ──
function addLink() {
  if (!newLink.value.title.trim() || !newLink.value.url.trim()) return
  const newId = Math.max(0, ...links.value.map(l => l.id)) + 1
  links.value.push({ id: newId, ...newLink.value, pinned: false })
  newLink.value = { title: '', url: '', description: '' }
  showAddLink.value = false
}

function togglePin(linkId: number) {
  const link = links.value.find(l => l.id === linkId)
  if (link) link.pinned = !link.pinned
}

function removeLink(linkId: number) {
  links.value = links.value.filter(l => l.id !== linkId)
}

// ── Breadcrumb ──
const breadcrumbs = computed(() => {
  const crumbs: { label: string; action?: () => void }[] = [{ label: 'Wiki', action: goToDashboard }]
  if (currentView.value === 'links') {
    crumbs.push({ label: 'Saved Links' })
  } else if (currentView.value === 'area' && selectedArea.value) {
    crumbs.push({ label: selectedArea.value.name })
  } else if (currentView.value === 'page' && selectedArea.value && selectedPage.value) {
    crumbs.push({ label: selectedArea.value.name, action: () => goToArea(selectedArea.value!.id) })
    crumbs.push({ label: selectedPage.value.title })
  }
  return crumbs
})
</script>

<template>
  <div class="flex flex-col h-screen w-full bg-background text-foreground overflow-hidden">
    <!-- ── Top Bar ── -->
    <header class="h-12 min-h-[48px] border-b flex items-center px-5 gap-3 bg-background">
      <button
        v-if="currentView !== 'dashboard'"
        class="text-muted-foreground hover:text-foreground transition-colors"
        @click="breadcrumbs.length > 1 ? breadcrumbs[breadcrumbs.length - 2].action?.() : goToDashboard()"
      >
        <ArrowLeft class="h-4 w-4" />
      </button>
      <nav class="flex items-center gap-1 text-sm">
        <template v-for="(crumb, i) in breadcrumbs" :key="i">
          <span v-if="i > 0" class="text-muted-foreground/40">/</span>
          <button
            v-if="crumb.action && i < breadcrumbs.length - 1"
            class="text-muted-foreground hover:text-foreground transition-colors"
            @click="crumb.action"
          >
            {{ crumb.label }}
          </button>
          <span v-else class="font-medium">{{ crumb.label }}</span>
        </template>
      </nav>
    </header>

    <!-- ── Dashboard View ── -->
    <ScrollArea v-if="currentView === 'dashboard'" class="flex-1">
      <div class="max-w-5xl mx-auto px-8 py-8 space-y-8">
        <!-- Pinned Links -->
        <section>
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold flex items-center gap-2">
              <Star class="h-4.5 w-4.5 text-amber-500" />
              Pinned Links
            </h2>
            <Button variant="ghost" size="sm" class="text-xs text-muted-foreground" @click="goToLinks">
              View all links
              <ChevronRight class="h-3 w-3 ml-1" />
            </Button>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <a
              v-for="link in pinnedLinks"
              :key="link.id"
              :href="link.url"
              target="_blank"
              class="group flex items-start gap-3 p-3.5 rounded-lg border bg-card hover:bg-muted/50 hover:border-primary/30 transition-all"
            >
              <div class="mt-0.5 h-8 w-8 rounded-md bg-primary/10 flex items-center justify-center shrink-0">
                <Globe class="h-4 w-4 text-primary" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium truncate group-hover:text-primary transition-colors">{{ link.title }}</p>
                <p class="text-xs text-muted-foreground mt-0.5 truncate">{{ link.description }}</p>
              </div>
              <ExternalLink class="h-3.5 w-3.5 text-muted-foreground/0 group-hover:text-muted-foreground/60 transition-all mt-1 shrink-0" />
            </a>
          </div>
        </section>

        <!-- Areas -->
        <section>
          <h2 class="text-lg font-semibold flex items-center gap-2 mb-4">
            <Folder class="h-4.5 w-4.5 text-primary" />
            Areas
          </h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <button
              v-for="area in areas"
              :key="area.id"
              class="group flex items-start gap-4 p-4 rounded-lg border bg-card hover:bg-muted/50 hover:border-primary/30 transition-all text-left"
              @click="goToArea(area.id)"
            >
              <div :class="['h-10 w-10 rounded-lg border flex items-center justify-center text-lg shrink-0', area.color]">
                {{ area.icon }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-semibold group-hover:text-primary transition-colors">{{ area.name }}</p>
                <p class="text-xs text-muted-foreground mt-0.5">{{ area.description }}</p>
                <p class="text-xs text-muted-foreground/60 mt-1.5">{{ area.pageCount }} pages</p>
              </div>
              <ChevronRight class="h-4 w-4 text-muted-foreground/30 group-hover:text-muted-foreground/60 transition-all mt-1 shrink-0" />
            </button>
          </div>
        </section>

        <!-- Recent Links -->
        <section>
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold flex items-center gap-2">
              <Globe class="h-4.5 w-4.5 text-muted-foreground" />
              All Links
            </h2>
            <Button variant="outline" size="sm" class="text-xs h-7" @click="showAddLink = true">
              <Plus class="h-3 w-3 mr-1" />
              Add Link
            </Button>
          </div>
          <div class="space-y-1">
            <a
              v-for="link in otherLinks.slice(0, 5)"
              :key="link.id"
              :href="link.url"
              target="_blank"
              class="group flex items-center gap-3 px-3 py-2 rounded-md hover:bg-muted/50 transition-colors"
            >
              <Globe class="h-3.5 w-3.5 text-muted-foreground shrink-0" />
              <span class="text-sm truncate flex-1">{{ link.title }}</span>
              <span class="text-xs text-muted-foreground/60 truncate max-w-[200px]">{{ link.description }}</span>
              <ExternalLink class="h-3 w-3 text-muted-foreground/0 group-hover:text-muted-foreground/60 shrink-0" />
            </a>
          </div>
          <Button v-if="otherLinks.length > 5" variant="ghost" size="sm" class="text-xs text-muted-foreground mt-2" @click="goToLinks">
            View all {{ links.length }} links →
          </Button>
        </section>
      </div>
    </ScrollArea>

    <!-- ── All Links View ── -->
    <ScrollArea v-else-if="currentView === 'links'" class="flex-1">
      <div class="max-w-4xl mx-auto px-8 py-6 space-y-4">
        <div class="flex items-center gap-3">
          <div class="relative flex-1">
            <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground" />
            <Input v-model="searchQuery" placeholder="Search links..." class="pl-8 h-9 text-sm" />
          </div>
          <Button variant="outline" size="sm" class="h-9" @click="showAddLink = true">
            <Plus class="h-3.5 w-3.5 mr-1.5" />
            Add Link
          </Button>
        </div>

        <div class="space-y-1">
          <div
            v-for="link in filteredLinks"
            :key="link.id"
            class="group flex items-center gap-3 px-3 py-2.5 rounded-md hover:bg-muted/50 transition-colors border border-transparent hover:border-border"
          >
            <div class="h-8 w-8 rounded-md bg-primary/10 flex items-center justify-center shrink-0">
              <Globe class="h-4 w-4 text-primary" />
            </div>
            <div class="min-w-0 flex-1">
              <a :href="link.url" target="_blank" class="text-sm font-medium hover:text-primary transition-colors">{{ link.title }}</a>
              <p class="text-xs text-muted-foreground truncate">{{ link.description }}</p>
            </div>
            <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
              <button
                :class="['h-7 w-7 flex items-center justify-center rounded hover:bg-muted transition-colors', link.pinned ? 'text-amber-500' : 'text-muted-foreground']"
                @click="togglePin(link.id)"
                :title="link.pinned ? 'Unpin' : 'Pin'"
              >
                <Star :class="['h-3.5 w-3.5', link.pinned ? 'fill-current' : '']" />
              </button>
              <button
                class="h-7 w-7 flex items-center justify-center rounded text-muted-foreground hover:text-destructive hover:bg-muted transition-colors"
                @click="removeLink(link.id)"
                title="Remove"
              >
                <X class="h-3.5 w-3.5" />
              </button>
            </div>
            <Badge v-if="link.pinned" variant="secondary" class="text-[10px] px-1.5 shrink-0">Pinned</Badge>
          </div>
        </div>

        <p v-if="filteredLinks.length === 0" class="text-sm text-muted-foreground text-center py-8">
          No links matching '{{ searchQuery }}'
        </p>
      </div>
    </ScrollArea>

    <!-- ── Area View (page list) ── -->
    <ScrollArea v-else-if="currentView === 'area' && selectedArea" class="flex-1">
      <div class="max-w-4xl mx-auto px-8 py-6 space-y-6">
        <!-- Area Header -->
        <div class="flex items-start gap-4">
          <div :class="['h-12 w-12 rounded-lg border flex items-center justify-center text-xl shrink-0', selectedArea.color]">
            {{ selectedArea.icon }}
          </div>
          <div>
            <h1 class="text-xl font-semibold">{{ selectedArea.name }}</h1>
            <p class="text-sm text-muted-foreground mt-0.5">{{ selectedArea.description }}</p>
          </div>
          <Button variant="outline" size="sm" class="ml-auto h-8 text-xs" @click="createPage">
            <Plus class="h-3 w-3 mr-1.5" />
            New Page
          </Button>
        </div>

        <!-- Pages List -->
        <div class="space-y-1">
          <button
            v-for="page in areaPages"
            :key="page.id"
            class="group w-full flex items-center gap-3 px-4 py-3 rounded-lg border bg-card hover:bg-muted/50 hover:border-primary/30 transition-all text-left"
            @click="goToPage(page.id)"
          >
            <FileText class="h-4 w-4 text-muted-foreground shrink-0" />
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium group-hover:text-primary transition-colors">{{ page.title }}</p>
              <p class="text-xs text-muted-foreground mt-0.5">Updated {{ page.updatedAt }}</p>
            </div>
            <ChevronRight class="h-4 w-4 text-muted-foreground/30 group-hover:text-muted-foreground/60 transition-all shrink-0" />
          </button>
        </div>

        <p v-if="areaPages.length === 0" class="text-sm text-muted-foreground text-center py-12">
          No pages yet. Create your first page in this area.
        </p>
      </div>
    </ScrollArea>

    <!-- ── Page View (full width editor/preview) ── -->
    <template v-else-if="currentView === 'page' && selectedPage">
      <!-- Page Header -->
      <div class="px-6 pt-3 pb-2 border-b flex items-center justify-between">
        <div class="flex-1 min-w-0">
          <input
            v-if="editingTitle"
            :value="titleInput"
            class="text-lg font-semibold bg-transparent border-b-2 border-primary outline-none w-full"
            @input="titleInput = ($event.target as HTMLInputElement).value"
            @keydown.enter="saveTitle"
            @keydown.escape="editingTitle = false"
            @blur="saveTitle"
            autofocus
          />
          <h1
            v-else
            class="text-lg font-semibold cursor-pointer hover:text-primary transition-colors"
            @dblclick="startEditTitle"
          >
            {{ selectedPage.title }}
            <Pencil class="inline h-3 w-3 ml-1 opacity-30" />
          </h1>
          <p class="text-xs text-muted-foreground">Updated {{ selectedPage.updatedAt }}</p>
        </div>
        <div class="flex items-center gap-2">
          <Tabs :model-value="mode" @update:model-value="(v: any) => mode = v as 'edit' | 'preview'">
            <TabsList class="h-8">
              <TabsTrigger value="edit" class="text-xs px-3">Edit</TabsTrigger>
              <TabsTrigger value="preview" class="text-xs px-3">Preview</TabsTrigger>
            </TabsList>
          </Tabs>
          <Button
            variant="ghost"
            size="icon"
            class="h-8 w-8 text-muted-foreground hover:text-destructive"
            aria-label="Delete page"
            @click="showDeleteDialog = true"
          >
            <Trash2 class="h-4 w-4" />
          </Button>
        </div>
      </div>

      <!-- Full-width content -->
      <ScrollArea class="flex-1">
        <div class="px-6 py-4">
          <textarea
            v-if="mode === 'edit'"
            :value="selectedPage.content"
            @input="selectedPage.content = ($event.target as HTMLTextAreaElement).value"
            class="w-full h-[calc(100vh-140px)] bg-muted/30 border rounded-md p-4 text-sm font-mono leading-relaxed resize-none outline-none focus:border-primary/50 focus:ring-1 focus:ring-primary/20"
            placeholder="Start writing in markdown..."
          />
          <div
            v-else
            class="wiki-content prose prose-zinc dark:prose-invert max-w-none"
            @click="handleWikiLinkClick"
            v-html="renderedContent"
          />
        </div>
      </ScrollArea>
    </template>

    <!-- ── Delete Confirmation ── -->
    <Dialog v-model:open="showDeleteDialog">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>Delete page?</DialogTitle>
          <DialogDescription>
            '{{ selectedPage?.title }}' will be permanently deleted.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="showDeleteDialog = false">Cancel</Button>
          <Button variant="destructive" @click="deletePage">Delete</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- ── Add Link Dialog ── -->
    <Dialog v-model:open="showAddLink">
      <DialogContent class="max-w-md">
        <DialogHeader>
          <DialogTitle>Add Link</DialogTitle>
          <DialogDescription>Save a bookmark for quick access</DialogDescription>
        </DialogHeader>
        <div class="space-y-3 py-2">
          <div>
            <label class="text-xs font-medium text-muted-foreground mb-1 block">Title</label>
            <Input v-model="newLink.title" placeholder="Grafana Dashboard" class="h-9 text-sm" />
          </div>
          <div>
            <label class="text-xs font-medium text-muted-foreground mb-1 block">URL</label>
            <Input v-model="newLink.url" placeholder="https://..." class="h-9 text-sm" />
          </div>
          <div>
            <label class="text-xs font-medium text-muted-foreground mb-1 block">Description (optional)</label>
            <Input v-model="newLink.description" placeholder="What is this link for?" class="h-9 text-sm" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showAddLink = false">Cancel</Button>
          <Button :disabled="!newLink.title.trim() || !newLink.url.trim()" @click="addLink">Save Link</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<style scoped>
.wiki-content :deep(h1) { font-size: 24px; font-weight: 600; line-height: 1.3; margin-top: 24px; margin-bottom: 8px; }
.wiki-content :deep(h2) { font-size: 20px; font-weight: 600; line-height: 1.3; margin-top: 20px; margin-bottom: 8px; }
.wiki-content :deep(h3) { font-size: 16px; font-weight: 600; line-height: 1.4; margin-top: 16px; margin-bottom: 4px; }
.wiki-content :deep(p) { font-size: 15px; line-height: 1.6; margin-bottom: 12px; }
.wiki-content :deep(li) { font-size: 15px; line-height: 1.6; margin-bottom: 4px; }
.wiki-content :deep(blockquote) { border-left: 3px solid hsl(var(--primary)); padding-left: 16px; font-style: italic; color: hsl(var(--muted-foreground)); }
.wiki-content :deep(code:not(pre code)) { font-size: 13px; background: hsl(var(--muted)); padding: 2px 6px; border-radius: 4px; font-family: "JetBrains Mono", "Cascadia Code", Consolas, monospace; }
.wiki-content :deep(pre) { background: hsl(var(--muted)); padding: 16px; border-radius: 8px; overflow-x: auto; margin-bottom: 16px; }
.wiki-content :deep(pre code) { font-size: 13px; line-height: 1.5; background: none; padding: 0; font-family: "JetBrains Mono", "Cascadia Code", Consolas, monospace; }
.wiki-content :deep(table) { width: 100%; border-collapse: collapse; margin-bottom: 16px; font-size: 14px; }
.wiki-content :deep(th), .wiki-content :deep(td) { border: 1px solid hsl(var(--border)); padding: 8px 12px; text-align: left; }
.wiki-content :deep(th) { background: hsl(var(--muted)); font-weight: 600; }
.wiki-content :deep(a.wiki-link) { color: hsl(var(--primary)); font-weight: 600; text-decoration: none; cursor: pointer; }
.wiki-content :deep(a.wiki-link:hover) { text-decoration: underline dotted; }
.wiki-content :deep(a.wiki-link-missing) { color: hsl(var(--muted-foreground)); font-weight: 400; text-decoration: underline dashed; }
.wiki-content :deep(a:not(.wiki-link)) { color: hsl(var(--primary)); text-decoration: none; }
.wiki-content :deep(a:not(.wiki-link):hover) { text-decoration: underline; }
.wiki-content :deep(ul), .wiki-content :deep(ol) { padding-left: 24px; margin-bottom: 12px; }
.wiki-content :deep(hr) { border: none; border-top: 1px solid hsl(var(--border)); margin: 24px 0; }
.wiki-content :deep(input[type="checkbox"]) { margin-right: 8px; }
</style>
