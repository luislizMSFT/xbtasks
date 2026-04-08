<script setup lang="ts">
/**
 * Playground: Work Item Lifecycle Tracking — PRs, Pipelines & Task Traceability.
 *
 * Two views:
 *  1. Traceability Table: rows with ADO type icons, parent chain, PR/pipeline status.
 *     Click any row → detail modal.
 *  2. Impact Graph: structured DAG showing ADO hierarchy + task relationships.
 */
import { ref, computed } from 'vue'
import type { Task, PullRequest, WorkItem } from '@/types'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import AzureDevOpsIcon from '@/components/icons/AzureDevOpsIcon.vue'
import LifecycleModal from './LifecycleModal.vue'
import ImpactGraph from './ImpactGraph.vue'
import {
  CheckCircle2, CircleDot, GitPullRequest, GitMerge,
  Rocket, AlertTriangle, Loader2, Check, X, Clock, XCircle,
  Activity, ChevronRight, User,
} from 'lucide-vue-next'
import {
  statusColor, adoTypeIcon, adoTypeColor, adoTypeHex,
  adoStateClasses, adoPriorityClasses,
} from '@/lib/styles'

// ── Types ──

interface PipelineRun {
  id: number
  name: string
  status: 'running' | 'succeeded' | 'failed' | 'queued' | 'cancelled'
  prId: number | null
  branch: string
  startedAt: string
  finishedAt: string | null
  duration: string | null
}

interface Deployment {
  id: number
  environment: string
  status: 'succeeded' | 'failed' | 'in_progress' | 'pending'
  pipelineRunId: number
  deployedAt: string
}

interface LifecycleItem {
  task: Task
  ado?: WorkItem          // ADO work item data (type, state, hierarchy)
  prs: PullRequest[]
  pipelines: PipelineRun[]
  deployments: Deployment[]
}

// ── Mock ADO hierarchy (parent items that aren't tracked as personal tasks) ──

const adoHierarchy: WorkItem[] = [
  {
    id: 54000, adoId: '54000', title: 'Sprint 12 — Auth & Backend Hardening',
    state: 'Active', type: 'Epic', assignedTo: 'Luis Lizama', priority: 1,
    areaPath: 'XboxServices\\Auth', description: '', url: '', org: 'msazure', project: 'Xbox',
    parentId: 0, changedDate: '2026-04-01T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
  },
  {
    id: 54300, adoId: '54300', title: 'Authentication Overhaul',
    state: 'Active', type: 'Feature', assignedTo: 'Luis Lizama', priority: 2,
    areaPath: 'XboxServices\\Auth', description: '', url: '', org: 'msazure', project: 'Xbox',
    parentId: 54000, changedDate: '2026-04-01T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
  },
  {
    id: 54320, adoId: '54320', title: 'Token refresh for seamless auth',
    state: 'Active', type: 'User Story', assignedTo: 'Luis Lizama', priority: 2,
    areaPath: 'XboxServices\\Auth', description: '', url: '', org: 'msazure', project: 'Xbox',
    parentId: 54300, changedDate: '2026-04-05T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
  },
  {
    id: 54390, adoId: '54390', title: 'Harden batch API security',
    state: 'Active', type: 'User Story', assignedTo: 'Alice Chen', priority: 1,
    areaPath: 'XboxServices\\Backend', description: '', url: '', org: 'msazure', project: 'Xbox',
    parentId: 54300, changedDate: '2026-04-02T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
  },
  {
    id: 54450, adoId: '54450', title: 'Observability Improvements',
    state: 'Resolved', type: 'Feature', assignedTo: 'Luis Lizama', priority: 2,
    areaPath: 'XboxServices\\Observability', description: '', url: '', org: 'msazure', project: 'Xbox',
    parentId: 54000, changedDate: '2026-04-05T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
  },
  {
    id: 54550, adoId: '54550', title: 'Infrastructure & Tooling',
    state: 'Active', type: 'Feature', assignedTo: 'Luis Lizama', priority: 2,
    areaPath: 'XboxServices\\Infra', description: '', url: '', org: 'msazure', project: 'Xbox',
    parentId: 54000, changedDate: '2026-04-01T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
  },
]

// All ADO items indexed by adoId for parent lookups
const adoById = computed(() => {
  const map = new Map<string, WorkItem>()
  for (const item of adoHierarchy) map.set(item.adoId, item)
  for (const item of mockItems) {
    if (item.ado) map.set(item.ado.adoId, item.ado)
  }
  return map
})

function getParentChain(adoId: string): WorkItem[] {
  const chain: WorkItem[] = []
  let current = adoById.value.get(adoId)
  while (current && current.parentId) {
    const parent = adoById.value.get(String(current.parentId))
    if (!parent || chain.includes(parent)) break
    chain.unshift(parent)
    current = parent
  }
  return chain
}

// ── Mock lifecycle items ──

const mockItems: LifecycleItem[] = [
  // ── ADO work items (with hierarchy) ──
  {
    task: {
      id: 1, title: 'Implement auth token refresh', description: 'Auto-refresh tokens before expiry',
      status: 'in_progress', priority: 'P0', category: '', projectId: 1, area: 'Auth',
      dueDate: '2026-04-10', adoId: '54321', tags: 'auth,backend', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 0,
      createdAt: '2026-04-01T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
    },
    ado: {
      id: 54321, adoId: '54321', title: 'Implement auth token refresh',
      state: 'Active', type: 'Task', assignedTo: 'Luis Lizama', priority: 1,
      areaPath: 'XboxServices\\Auth', description: '', url: '', org: 'msazure', project: 'Xbox',
      parentId: 54320, changedDate: '2026-04-07T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
    },
    prs: [{
      id: 1, title: 'feat: add token refresh logic', prUrl: '',
      prNumber: 101, repo: 'xbox-services', taskId: 1, adoId: '54321',
      status: 'active', reviewers: JSON.stringify([
        { displayName: 'Alice Chen', uniqueName: 'alice@ms.com', vote: 10 },
        { displayName: 'Bob Kim', uniqueName: 'bob@ms.com', vote: 0 },
      ]),
      sourceBranch: 'feature/token-refresh', targetBranch: 'main',
      votes: 10, createdBy: 'luisliz', createdAt: '2026-04-05T14:00:00Z',
      updatedAt: '2026-04-07T09:00:00Z', mergedAt: null,
    }],
    pipelines: [{
      id: 10, name: 'xbox-services-ci', status: 'succeeded', prId: 1,
      branch: 'feature/token-refresh', startedAt: '2026-04-07T09:01:00Z',
      finishedAt: '2026-04-07T09:08:00Z', duration: '7m 12s',
    }],
    deployments: [],
  },
  {
    task: {
      id: 2, title: 'Fix rate-limiting bypass on batch endpoints', description: 'Batch API can exceed rate limits',
      status: 'in_review', priority: 'P1', category: '', projectId: 1, area: 'Backend',
      dueDate: '2026-04-08', adoId: '54400', tags: 'security', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 1,
      createdAt: '2026-04-02T10:00:00Z', updatedAt: '2026-04-07T16:00:00Z', completedAt: null,
    },
    ado: {
      id: 54400, adoId: '54400', title: 'Fix rate-limiting bypass on batch endpoints',
      state: 'Active', type: 'Bug', assignedTo: 'Luis Lizama', priority: 1,
      areaPath: 'XboxServices\\Backend', description: '', url: '', org: 'msazure', project: 'Xbox',
      parentId: 54390, changedDate: '2026-04-07T16:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
    },
    prs: [{
      id: 2, title: 'fix: enforce rate limit on batch routes', prUrl: '',
      prNumber: 102, repo: 'xbox-services', taskId: 2, adoId: '54400',
      status: 'active', reviewers: JSON.stringify([
        { displayName: 'Alice Chen', uniqueName: 'alice@ms.com', vote: 10 },
        { displayName: 'Carol Wu', uniqueName: 'carol@ms.com', vote: 5 },
        { displayName: 'Dan Lee', uniqueName: 'dan@ms.com', vote: -5 },
      ]),
      sourceBranch: 'fix/rate-limit-batch', targetBranch: 'main',
      votes: 5, createdBy: 'luisliz', createdAt: '2026-04-06T11:00:00Z',
      updatedAt: '2026-04-07T15:00:00Z', mergedAt: null,
    }],
    pipelines: [
      { id: 20, name: 'xbox-services-ci', status: 'failed', prId: 2, branch: 'fix/rate-limit-batch', startedAt: '2026-04-07T15:01:00Z', finishedAt: '2026-04-07T15:06:00Z', duration: '5m 32s' },
      { id: 21, name: 'xbox-services-ci', status: 'running', prId: 2, branch: 'fix/rate-limit-batch', startedAt: '2026-04-07T16:00:00Z', finishedAt: null, duration: null },
    ],
    deployments: [],
  },
  {
    task: {
      id: 3, title: 'Add telemetry for sync operations', description: 'Track sync latency and failure rates',
      status: 'done', priority: 'P2', category: '', projectId: 1, area: 'Observability',
      dueDate: '', adoId: '54500', tags: 'telemetry', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 2,
      createdAt: '2026-03-28T10:00:00Z', updatedAt: '2026-04-05T10:00:00Z', completedAt: '2026-04-05T10:00:00Z',
    },
    ado: {
      id: 54500, adoId: '54500', title: 'Add telemetry for sync operations',
      state: 'Closed', type: 'Task', assignedTo: 'Luis Lizama', priority: 2,
      areaPath: 'XboxServices\\Observability', description: '', url: '', org: 'msazure', project: 'Xbox',
      parentId: 54450, changedDate: '2026-04-05T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
    },
    prs: [{
      id: 3, title: 'feat: add sync telemetry with OpenTelemetry', prUrl: '',
      prNumber: 98, repo: 'xbox-services', taskId: 3, adoId: '54500',
      status: 'completed', reviewers: JSON.stringify([{ displayName: 'Alice Chen', uniqueName: 'alice@ms.com', vote: 10 }]),
      sourceBranch: 'feature/sync-telemetry', targetBranch: 'main',
      votes: 10, createdBy: 'luisliz', createdAt: '2026-04-03T09:00:00Z',
      updatedAt: '2026-04-04T14:00:00Z', mergedAt: '2026-04-04T14:30:00Z',
    }],
    pipelines: [{
      id: 30, name: 'xbox-services-ci', status: 'succeeded', prId: 3, branch: 'feature/sync-telemetry',
      startedAt: '2026-04-04T14:31:00Z', finishedAt: '2026-04-04T14:38:00Z', duration: '6m 45s',
    }],
    deployments: [
      { id: 1, environment: 'staging', status: 'succeeded', pipelineRunId: 30, deployedAt: '2026-04-04T15:00:00Z' },
      { id: 2, environment: 'production', status: 'succeeded', pipelineRunId: 30, deployedAt: '2026-04-05T09:00:00Z' },
    ],
  },
  {
    task: {
      id: 5, title: 'Upgrade Go to 1.23', description: 'Module dependency update',
      status: 'blocked', priority: 'P1', category: '', projectId: 1, area: 'Infrastructure',
      dueDate: '2026-04-12', adoId: '54600', tags: 'infra', blockedReason: 'Waiting for upstream fix', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 4,
      createdAt: '2026-04-01T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
    },
    ado: {
      id: 54600, adoId: '54600', title: 'Upgrade Go to 1.23',
      state: 'Active', type: 'Task', assignedTo: 'Luis Lizama', priority: 2,
      areaPath: 'XboxServices\\Infra', description: '', url: '', org: 'msazure', project: 'Xbox',
      parentId: 54550, changedDate: '2026-04-07T10:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
    },
    prs: [{
      id: 5, title: 'chore: upgrade go 1.23 + update deps', prUrl: '',
      prNumber: 99, repo: 'xbox-services', taskId: 5, adoId: '54600',
      status: 'draft', reviewers: '[]',
      sourceBranch: 'chore/go-1.23', targetBranch: 'main',
      votes: 0, createdBy: 'luisliz', createdAt: '2026-04-03T10:00:00Z',
      updatedAt: '2026-04-03T10:00:00Z', mergedAt: null,
    }],
    pipelines: [{
      id: 50, name: 'xbox-services-ci', status: 'cancelled', prId: 5, branch: 'chore/go-1.23',
      startedAt: '2026-04-03T10:05:00Z', finishedAt: '2026-04-03T10:06:00Z', duration: '1m 02s',
    }],
    deployments: [],
  },
  {
    task: {
      id: 6, title: 'Migrate user settings to new schema', description: 'DB schema migration for user prefs',
      status: 'in_progress', priority: 'P1', category: '', projectId: 1, area: 'Backend',
      dueDate: '2026-04-15', adoId: '54700', tags: 'migration,database', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 5,
      createdAt: '2026-04-05T10:00:00Z', updatedAt: '2026-04-08T12:00:00Z', completedAt: null,
    },
    ado: {
      id: 54700, adoId: '54700', title: 'Migrate user settings to new schema',
      state: 'Active', type: 'Task', assignedTo: 'Luis Lizama', priority: 2,
      areaPath: 'XboxServices\\Backend', description: '', url: '', org: 'msazure', project: 'Xbox',
      parentId: 54550, changedDate: '2026-04-08T12:00:00Z', syncedAt: '2026-04-08T10:00:00Z',
    },
    prs: [{
      id: 6, title: 'feat: user settings schema v2 migration', prUrl: '',
      prNumber: 105, repo: 'xbox-services', taskId: 6, adoId: '54700',
      status: 'active', reviewers: JSON.stringify([
        { displayName: 'Bob Kim', uniqueName: 'bob@ms.com', vote: 0 },
      ]),
      sourceBranch: 'feature/settings-migration', targetBranch: 'main',
      votes: 0, createdBy: 'luisliz', createdAt: '2026-04-07T10:00:00Z',
      updatedAt: '2026-04-08T11:00:00Z', mergedAt: null,
    }],
    pipelines: [{
      id: 60, name: 'xbox-services-ci', status: 'succeeded', prId: 6, branch: 'feature/settings-migration',
      startedAt: '2026-04-08T11:01:00Z', finishedAt: '2026-04-08T11:09:00Z', duration: '8m 03s',
    }],
    deployments: [],
  },
  // ── Personal tasks (no ADO) ──
  {
    task: {
      id: 4, title: 'Refactor project settings page', description: 'Better org for settings layout',
      status: 'in_progress', priority: 'P2', category: '', projectId: 1, area: 'Frontend',
      dueDate: '', adoId: '', tags: 'ui', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 3,
      createdAt: '2026-04-06T10:00:00Z', updatedAt: '2026-04-07T10:00:00Z', completedAt: null,
    },
    prs: [], pipelines: [], deployments: [],
  },
  {
    task: {
      id: 7, title: 'Write onboarding runbook', description: 'Document local dev setup for new hires',
      status: 'todo', priority: 'P3', category: '', projectId: 1, area: 'Docs',
      dueDate: '', adoId: '', tags: 'docs', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 6,
      createdAt: '2026-04-06T08:00:00Z', updatedAt: '2026-04-06T08:00:00Z', completedAt: null,
    },
    prs: [], pipelines: [], deployments: [],
  },
  {
    task: {
      id: 8, title: 'Investigate flaky test in sync_test.go', description: 'TestSyncConcurrency fails ~10% of runs',
      status: 'in_progress', priority: 'P2', category: '', projectId: 1, area: 'Quality',
      dueDate: '2026-04-09', adoId: '', tags: 'testing,flaky', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 7,
      createdAt: '2026-04-07T09:00:00Z', updatedAt: '2026-04-08T14:00:00Z', completedAt: null,
    },
    prs: [{
      id: 8, title: 'fix: resolve sync test race condition', prUrl: '',
      prNumber: 107, repo: 'xbox-services', taskId: 8, adoId: '',
      status: 'active', reviewers: JSON.stringify([
        { displayName: 'Alice Chen', uniqueName: 'alice@ms.com', vote: 5 },
      ]),
      sourceBranch: 'fix/sync-flaky-test', targetBranch: 'main',
      votes: 5, createdBy: 'luisliz', createdAt: '2026-04-08T10:00:00Z',
      updatedAt: '2026-04-08T13:00:00Z', mergedAt: null,
    }],
    pipelines: [{
      id: 80, name: 'xbox-services-ci', status: 'succeeded', prId: 8, branch: 'fix/sync-flaky-test',
      startedAt: '2026-04-08T13:01:00Z', finishedAt: '2026-04-08T13:07:00Z', duration: '5m 55s',
    }],
    deployments: [],
  },
  {
    task: {
      id: 9, title: 'Update team dashboard widgets', description: 'Add new sprint velocity chart',
      status: 'done', priority: 'P3', category: '', projectId: 1, area: 'Frontend',
      dueDate: '', adoId: '', tags: 'ui,dashboard', blockedReason: '', blockedBy: '',
      parentId: null, personalPriority: '', sortOrder: 8,
      createdAt: '2026-04-03T10:00:00Z', updatedAt: '2026-04-06T16:00:00Z', completedAt: '2026-04-06T16:00:00Z',
    },
    prs: [{
      id: 9, title: 'feat: sprint velocity widget', prUrl: '',
      prNumber: 103, repo: 'xb-tasks', taskId: 9, adoId: '',
      status: 'completed', reviewers: JSON.stringify([{ displayName: 'Bob Kim', uniqueName: 'bob@ms.com', vote: 10 }]),
      sourceBranch: 'feature/velocity-widget', targetBranch: 'main',
      votes: 10, createdBy: 'luisliz', createdAt: '2026-04-05T14:00:00Z',
      updatedAt: '2026-04-06T15:30:00Z', mergedAt: '2026-04-06T15:45:00Z',
    }],
    pipelines: [{
      id: 90, name: 'xb-tasks-ci', status: 'succeeded', prId: 9, branch: 'feature/velocity-widget',
      startedAt: '2026-04-06T15:46:00Z', finishedAt: '2026-04-06T15:50:00Z', duration: '4m 12s',
    }],
    deployments: [
      { id: 3, environment: 'staging', status: 'succeeded', pipelineRunId: 90, deployedAt: '2026-04-06T16:00:00Z' },
    ],
  },
]

// ── State ──

const activeVariant = ref<'table' | 'graph'>('table')
const variants = [
  { id: 'table' as const, label: 'Traceability Table' },
  { id: 'graph' as const, label: 'Impact Graph' },
]

const modalOpen = ref(false)
const modalItem = ref<LifecycleItem | null>(null)
function openModal(item: LifecycleItem) {
  modalItem.value = item
  modalOpen.value = true
}

const selectedGraphNode = ref<number | null>(null)

// ── Impact graph: ADO hierarchy + tasks + PRs as nodes ──

const graphNodes = computed(() => [
  // ADO hierarchy nodes (Epic → Feature → User Story)
  ...adoHierarchy.map(wi => ({
    id: wi.id,
    title: wi.title,
    status: wi.state === 'Closed' || wi.state === 'Resolved' ? 'done' : wi.state === 'Active' ? 'in_progress' : 'todo',
    priority: `P${wi.priority}`,
    type: wi.type.toLowerCase(),
    adoId: wi.adoId,
  })),
  // Task nodes (both ADO and personal)
  ...mockItems.map(item => ({
    id: item.ado ? item.ado.id : item.task.id,
    title: item.task.title,
    status: item.task.status,
    priority: item.task.priority,
    type: item.ado?.type.toLowerCase() ?? 'personal',
    adoId: item.task.adoId || undefined,
  })),
])

const graphEdges = computed(() => [
  // ADO hierarchy: parent → child
  ...adoHierarchy
    .filter(wi => wi.parentId && adoHierarchy.some(p => p.id === wi.parentId))
    .map(wi => ({ source: wi.parentId, target: wi.id, kind: 'parent' as const })),
  // ADO task → parent (User Story / Feature)
  ...mockItems
    .filter(item => item.ado?.parentId)
    .map(item => ({ source: item.ado!.parentId, target: item.ado!.id, kind: 'parent' as const })),
  // Task → PR links
  ...mockItems.flatMap(item => item.prs.map(pr => ({
    source: item.ado ? item.ado.id : item.task.id,
    target: 2000 + pr.id,
    kind: 'pr_link' as const,
  }))),
  // Cross-item dependencies
  { source: 54321, target: 54400, kind: 'dependency' as const },   // auth task → rate-limit bug
  { source: 54600, target: 54700, kind: 'triggers' as const },     // Go upgrade triggers settings migration
])

// PR nodes (only for items that have PRs)
const graphNodesWithPRs = computed(() => [
  ...graphNodes.value,
  ...mockItems.flatMap(item => item.prs.map(pr => ({
    id: 2000 + pr.id,
    title: `!${pr.prNumber} ${pr.title}`,
    status: pr.status === 'completed' ? 'done' : pr.status === 'draft' ? 'todo' : 'in_progress',
    priority: 'P2',
    type: 'pr',
  }))),
])

// ── Table helpers ──

function pipelineStatusIcon(status: string) {
  switch (status) {
    case 'succeeded': return Check
    case 'failed': return X
    case 'running': return Loader2
    case 'queued': return Clock
    case 'cancelled': return XCircle
    default: return Clock
  }
}

function pipelineStatusColor(status: string): string {
  switch (status) {
    case 'succeeded': return 'text-emerald-500'
    case 'failed': return 'text-red-500'
    case 'running': return 'text-blue-500'
    default: return 'text-zinc-400'
  }
}

function deployStatusColor(status: string): string {
  switch (status) {
    case 'succeeded': return 'text-emerald-500'
    case 'failed': return 'text-red-500'
    case 'in_progress': return 'text-blue-500'
    default: return 'text-zinc-400'
  }
}

function parseReviewers(json: string): Array<{ displayName: string; uniqueName: string; vote: number }> {
  try { return JSON.parse(json) } catch { return [] }
}

function relativeTime(iso: string): string {
  const ms = Date.now() - new Date(iso).getTime()
  const min = Math.floor(ms / 60000)
  if (min < 1) return 'just now'
  if (min < 60) return `${min}m ago`
  const hrs = Math.floor(min / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
}

function lifecycleStageLabel(item: LifecycleItem): string {
  if (item.deployments.some(d => d.status === 'succeeded' && d.environment === 'production')) return 'Shipped'
  if (item.deployments.some(d => d.status === 'in_progress' || d.status === 'pending')) return 'Deploying'
  if (item.pipelines.some(p => p.status === 'running' || p.status === 'queued')) return 'Building'
  if (item.prs.length > 0) return 'Code Review'
  return 'Development'
}

function stageBgClasses(label: string): string {
  switch (label) {
    case 'Development': return 'bg-zinc-500/15 text-zinc-500 border-zinc-500/25'
    case 'Code Review': return 'bg-blue-500/15 text-blue-600 dark:text-blue-400 border-blue-500/25'
    case 'Building': return 'bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/25'
    case 'Deploying': return 'bg-violet-500/15 text-violet-600 dark:text-violet-400 border-violet-500/25'
    case 'Shipped': return 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/25'
    default: return 'bg-muted text-muted-foreground border-border'
  }
}

const stats = computed(() => ({
  total: mockItems.length,
  ado: mockItems.filter(i => i.ado).length,
  personal: mockItems.filter(i => !i.ado).length,
  withPR: mockItems.filter(i => i.prs.length > 0).length,
  building: mockItems.filter(i => i.pipelines.some(p => p.status === 'running')).length,
  deployed: mockItems.filter(i => i.deployments.some(d => d.status === 'succeeded')).length,
  failing: mockItems.filter(i => i.pipelines.some(p => p.status === 'failed')).length,
}))
</script>

<template>
  <div class="h-screen flex flex-col bg-background text-foreground">
    <!-- Header -->
    <div class="h-10 shrink-0 border-b border-border flex items-center px-4 gap-3">
      <Activity :size="14" class="text-blue-500" />
      <span class="text-sm font-semibold">Playground: Work Item Lifecycle</span>
      <Badge variant="secondary" class="text-[9px] h-4">Phase 4</Badge>
      <div class="flex-1" />

      <!-- Summary pills -->
      <div class="flex items-center gap-2 text-[10px] mr-4">
        <span class="text-muted-foreground tabular-nums">{{ stats.total }} items</span>
        <Badge variant="outline" class="text-[9px] h-4 px-1.5 text-blue-500 border-blue-500/30">
          <AzureDevOpsIcon :size="8" class="mr-0.5" />{{ stats.ado }} ADO
        </Badge>
        <Badge variant="outline" class="text-[9px] h-4 px-1.5 text-zinc-400 border-zinc-400/30">
          <User :size="8" class="mr-0.5" />{{ stats.personal }} personal
        </Badge>
        <Badge variant="outline" class="text-[9px] h-4 px-1.5 text-blue-500 border-blue-500/30">
          <GitPullRequest :size="9" class="mr-0.5" />{{ stats.withPR }} PRs
        </Badge>
        <Badge v-if="stats.building" variant="outline" class="text-[9px] h-4 px-1.5 text-amber-500 border-amber-500/30">
          <Loader2 :size="9" class="mr-0.5 animate-spin" />{{ stats.building }} building
        </Badge>
        <Badge v-if="stats.failing" variant="outline" class="text-[9px] h-4 px-1.5 text-red-500 border-red-500/30">
          <AlertTriangle :size="9" class="mr-0.5" />{{ stats.failing }} failing
        </Badge>
        <Badge v-if="stats.deployed" variant="outline" class="text-[9px] h-4 px-1.5 text-emerald-500 border-emerald-500/30">
          <Rocket :size="9" class="mr-0.5" />{{ stats.deployed }} shipped
        </Badge>
      </div>

      <!-- Variant switcher -->
      <div class="flex gap-1">
        <Button
          v-for="v in variants" :key="v.id"
          :variant="activeVariant === v.id ? 'default' : 'outline'"
          size="sm" class="h-7 text-[10px]"
          @click="activeVariant = v.id"
        >{{ v.label }}</Button>
      </div>
    </div>

    <!-- ═══ TRACEABILITY TABLE ═══ -->
    <div v-if="activeVariant === 'table'" class="flex-1 overflow-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-border text-left">
            <th class="px-4 py-2 text-[10px] font-semibold text-muted-foreground uppercase tracking-wider w-[340px]">Work Item</th>
            <th class="px-4 py-2 text-[10px] font-semibold text-muted-foreground uppercase tracking-wider">Stage</th>
            <th class="px-4 py-2 text-[10px] font-semibold text-muted-foreground uppercase tracking-wider">PR</th>
            <th class="px-4 py-2 text-[10px] font-semibold text-muted-foreground uppercase tracking-wider">Build</th>
            <th class="px-4 py-2 text-[10px] font-semibold text-muted-foreground uppercase tracking-wider">Deploy</th>
            <th class="px-4 py-2 text-[10px] font-semibold text-muted-foreground uppercase tracking-wider">Updated</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr
            v-for="item in mockItems" :key="item.task.id"
            class="hover:bg-muted/50 transition-colors cursor-pointer"
            @click="openModal(item)"
          >
            <!-- Work Item cell: ADO type icon + title + parent breadcrumb -->
            <td class="px-4 py-3">
              <div class="flex items-start gap-2">
                <!-- ADO type icon or generic status icon -->
                <div class="mt-0.5 shrink-0">
                  <component
                    v-if="item.ado"
                    :is="adoTypeIcon(item.ado.type)"
                    :size="14" :class="adoTypeColor(item.ado.type)"
                  />
                  <component
                    v-else
                    :is="item.task.status === 'done' ? CheckCircle2 : item.task.status === 'blocked' ? AlertTriangle : CircleDot"
                    :size="14" :class="statusColor(item.task.status)"
                  />
                </div>
                <div class="min-w-0 flex-1">
                  <!-- Parent breadcrumb for ADO items -->
                  <div v-if="item.ado && item.task.adoId" class="flex items-center gap-0.5 flex-wrap mb-0.5">
                    <template v-for="(parent, i) in getParentChain(item.task.adoId)" :key="parent.adoId">
                      <component :is="adoTypeIcon(parent.type)" :size="9" :class="adoTypeColor(parent.type)" />
                      <span class="text-[9px] text-muted-foreground truncate max-w-[100px]">{{ parent.title }}</span>
                      <ChevronRight :size="8" class="text-muted-foreground/40" />
                    </template>
                  </div>
                  <!-- Title row -->
                  <div class="text-sm font-medium truncate">{{ item.task.title }}</div>
                  <!-- Meta row -->
                  <div class="flex items-center gap-1.5 mt-0.5">
                    <span class="text-[10px] text-muted-foreground">{{ item.task.priority }}</span>
                    <Badge v-if="item.ado" variant="outline" class="text-[9px] h-3.5 px-1 gap-0.5" :class="adoTypeColor(item.ado.type)">
                      {{ item.ado.type }}
                    </Badge>
                    <Badge v-if="item.ado" :class="['text-[9px] h-3.5 px-1', adoStateClasses(item.ado.state)]">
                      {{ item.ado.state }}
                    </Badge>
                    <Badge v-if="item.task.adoId" variant="outline" class="text-[9px] h-3.5 px-1 gap-0.5">
                      <AzureDevOpsIcon :size="8" />#{{ item.task.adoId }}
                    </Badge>
                    <Badge v-if="!item.ado" variant="outline" class="text-[9px] h-3.5 px-1 text-zinc-400 border-zinc-400/30">
                      <User :size="8" class="mr-0.5" />Personal
                    </Badge>
                  </div>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <Badge :class="['text-[9px] h-4 px-1.5', stageBgClasses(lifecycleStageLabel(item))]">
                {{ lifecycleStageLabel(item) }}
              </Badge>
            </td>
            <td class="px-4 py-3">
              <div v-if="item.prs.length" class="space-y-1">
                <div v-for="pr in item.prs" :key="pr.id" class="flex items-center gap-1.5">
                  <component :is="pr.status === 'completed' ? GitMerge : GitPullRequest" :size="12"
                    :class="pr.status === 'completed' ? 'text-violet-500' : pr.status === 'draft' ? 'text-zinc-400' : 'text-emerald-500'" />
                  <span class="text-[11px]">!{{ pr.prNumber }}</span>
                  <div class="flex items-center -space-x-0.5">
                    <span v-for="r in parseReviewers(pr.reviewers)" :key="r.uniqueName"
                      class="w-4 h-4 rounded-full flex items-center justify-center text-[8px] border border-background"
                      :class="r.vote >= 10 ? 'bg-emerald-500/20 text-emerald-600' : r.vote === 5 ? 'bg-emerald-400/20 text-emerald-500' : r.vote === -5 ? 'bg-amber-500/20 text-amber-600' : r.vote <= -10 ? 'bg-red-500/20 text-red-600' : 'bg-muted text-muted-foreground'"
                      :title="r.displayName + ': ' + r.vote"
                    >{{ r.displayName[0] }}</span>
                  </div>
                </div>
              </div>
              <span v-else class="text-[10px] text-muted-foreground/60">&mdash;</span>
            </td>
            <td class="px-4 py-3">
              <div v-if="item.pipelines.length" class="space-y-1">
                <div v-for="run in item.pipelines" :key="run.id" class="flex items-center gap-1.5">
                  <component :is="pipelineStatusIcon(run.status)" :size="12"
                    :class="[pipelineStatusColor(run.status), run.status === 'running' && 'animate-spin']" />
                  <span class="text-[11px]" :class="pipelineStatusColor(run.status)">{{ run.status }}</span>
                  <span v-if="run.duration" class="text-[9px] text-muted-foreground">{{ run.duration }}</span>
                </div>
              </div>
              <span v-else class="text-[10px] text-muted-foreground/60">&mdash;</span>
            </td>
            <td class="px-4 py-3">
              <div v-if="item.deployments.length" class="space-y-1">
                <div v-for="dep in item.deployments" :key="dep.id" class="flex items-center gap-1.5">
                  <Rocket :size="12" :class="deployStatusColor(dep.status)" />
                  <span class="text-[10px] capitalize">{{ dep.environment }}</span>
                </div>
              </div>
              <span v-else class="text-[10px] text-muted-foreground/60">&mdash;</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-[10px] text-muted-foreground">{{ relativeTime(item.task.updatedAt) }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ═══ IMPACT GRAPH (ADO hierarchy + tasks + PRs) ═══ -->
    <div v-else-if="activeVariant === 'graph'" class="flex-1 min-h-0 relative">
      <ImpactGraph
        :nodes="graphNodesWithPRs"
        :edges="graphEdges"
        :selected-id="selectedGraphNode"
        @select="(id: number) => selectedGraphNode = id"
      />
    </div>

    <!-- Detail modal -->
    <LifecycleModal v-model:open="modalOpen" :item="modalItem" />

    <!-- Footer -->
    <div class="border-t border-border px-4 py-3 text-[10px] text-muted-foreground">
      <strong>Phase 4 Playground:</strong>
      Traceability Table = ADO hierarchy (Epic→Feature→Story→Task/Bug) with parent breadcrumbs, type icons, state badges. Click rows for full detail.
      &middot; Impact Graph = structured DAG showing full ADO hierarchy + task dependencies + PR links.
    </div>
  </div>
</template>
