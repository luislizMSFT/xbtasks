import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface PRReviewer {
  displayName: string
  uniqueName: string
  vote: number // 10=approve, 5=approve-with-suggestions, -5=wait, -10=reject, 0=none
}

export interface PullRequest {
  id: number
  title: string
  prUrl: string
  prNumber: number
  repo: string
  taskId: number | null
  adoId: string
  status: string // draft, active, completed, abandoned
  reviewers: string // JSON array
  sourceBranch: string
  targetBranch: string
  votes: number
  createdAt: string
  updatedAt: string
  mergedAt: string | null
}

export function parseReviewers(reviewersJson: string): PRReviewer[] {
  try { return JSON.parse(reviewersJson) } catch { return [] }
}

export function branchName(ref: string): string {
  return ref.replace('refs/heads/', '')
}

export function voteIcon(vote: number): string {
  switch (vote) {
    case 10: return '👍'
    case 5: return '👍?'
    case -5: return '⏳'
    case -10: return '👎'
    default: return ''
  }
}

export function relativeTime(dateStr: string): string {
  if (!dateStr) return ''
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diffMs = now - then
  if (diffMs < 0) return 'just now'
  const minutes = Math.floor(diffMs / 60000)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}d ago`
  const months = Math.floor(days / 30)
  return `${months}mo ago`
}

const MOCK_REVIEW_PRS: PullRequest[] = [
  {
    id: 1001, title: 'Add rate limiting to gateway', prUrl: 'https://dev.azure.com/xbox/xb-tasks/_git/xb-gateway/pullrequest/1001',
    prNumber: 1001, repo: 'xb-gateway', taskId: null, adoId: 'ADO-48400', status: 'active',
    reviewers: JSON.stringify([
      { displayName: 'Luis', uniqueName: 'luis@xbox.com', vote: 0 },
      { displayName: 'Jordan', uniqueName: 'jordan@xbox.com', vote: 10 },
    ]),
    sourceBranch: 'refs/heads/feature/rate-limiting', targetBranch: 'refs/heads/main',
    votes: 1, createdAt: new Date(Date.now() - 7200000).toISOString(),
    updatedAt: new Date(Date.now() - 3600000).toISOString(), mergedAt: null,
  },
  {
    id: 1002, title: 'Schema migration v2', prUrl: 'https://dev.azure.com/xbox/xb-tasks/_git/xb-services/pullrequest/1002',
    prNumber: 1002, repo: 'xb-services', taskId: 3, adoId: 'ADO-48350', status: 'active',
    reviewers: JSON.stringify([
      { displayName: 'Luis', uniqueName: 'luis@xbox.com', vote: 0 },
    ]),
    sourceBranch: 'refs/heads/feature/schema-v2', targetBranch: 'refs/heads/main',
    votes: 0, createdAt: new Date(Date.now() - 14400000).toISOString(),
    updatedAt: new Date(Date.now() - 7200000).toISOString(), mergedAt: null,
  },
]

const MOCK_MY_PRS: PullRequest[] = [
  {
    id: 2001, title: 'Fix auth redirect loop', prUrl: 'https://dev.azure.com/xbox/xb-tasks/_git/xb-services/pullrequest/2001',
    prNumber: 2001, repo: 'xb-services', taskId: 1, adoId: 'ADO-48291', status: 'active',
    reviewers: JSON.stringify([
      { displayName: 'Alex', uniqueName: 'alex@xbox.com', vote: 10 },
      { displayName: 'Sam', uniqueName: 'sam@xbox.com', vote: 0 },
    ]),
    sourceBranch: 'refs/heads/fix/auth-redirect', targetBranch: 'refs/heads/main',
    votes: 1, createdAt: new Date(Date.now() - 10800000).toISOString(),
    updatedAt: new Date(Date.now() - 1800000).toISOString(), mergedAt: null,
  },
  {
    id: 2002, title: 'Update CI pipeline', prUrl: 'https://dev.azure.com/xbox/xb-tasks/_git/xb-infra/pullrequest/2002',
    prNumber: 2002, repo: 'xb-infra', taskId: null, adoId: '', status: 'completed',
    reviewers: JSON.stringify([
      { displayName: 'Sam', uniqueName: 'sam@xbox.com', vote: 10 },
      { displayName: 'Alex', uniqueName: 'alex@xbox.com', vote: 10 },
    ]),
    sourceBranch: 'refs/heads/ci/pipeline-update', targetBranch: 'refs/heads/main',
    votes: 2, createdAt: new Date(Date.now() - 172800000).toISOString(),
    updatedAt: new Date(Date.now() - 86400000).toISOString(), mergedAt: new Date(Date.now() - 86400000).toISOString(),
  },
  {
    id: 2003, title: 'Add dark mode token audit', prUrl: 'https://dev.azure.com/xbox/xb-tasks/_git/xb-tasks/pullrequest/2003',
    prNumber: 2003, repo: 'xb-tasks', taskId: 8, adoId: '', status: 'draft',
    reviewers: JSON.stringify([]),
    sourceBranch: 'refs/heads/feature/dark-mode-audit', targetBranch: 'refs/heads/main',
    votes: 0, createdAt: new Date(Date.now() - 259200000).toISOString(),
    updatedAt: new Date(Date.now() - 172800000).toISOString(), mergedAt: null,
  },
]

const MOCK_TEAM_PRS: PullRequest[] = [
  {
    id: 3001, title: 'Refactor notification service', prUrl: 'https://dev.azure.com/xbox/xb-tasks/_git/xb-services/pullrequest/3001',
    prNumber: 3001, repo: 'xb-services', taskId: null, adoId: 'ADO-48500', status: 'active',
    reviewers: JSON.stringify([
      { displayName: 'Jordan', uniqueName: 'jordan@xbox.com', vote: 5 },
      { displayName: 'Sam', uniqueName: 'sam@xbox.com', vote: 0 },
    ]),
    sourceBranch: 'refs/heads/refactor/notifications', targetBranch: 'refs/heads/main',
    votes: 1, createdAt: new Date(Date.now() - 43200000).toISOString(),
    updatedAt: new Date(Date.now() - 21600000).toISOString(), mergedAt: null,
  },
  {
    id: 3002, title: 'Add telemetry hooks', prUrl: 'https://dev.azure.com/xbox/xb-tasks/_git/xb-infra/pullrequest/3002',
    prNumber: 3002, repo: 'xb-infra', taskId: null, adoId: 'ADO-48510', status: 'active',
    reviewers: JSON.stringify([
      { displayName: 'Alex', uniqueName: 'alex@xbox.com', vote: 0 },
    ]),
    sourceBranch: 'refs/heads/feature/telemetry', targetBranch: 'refs/heads/main',
    votes: 0, createdAt: new Date(Date.now() - 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 43200000).toISOString(), mergedAt: null,
  },
]

let useMock = false

export const usePRStore = defineStore('prs', () => {
  const myPRs = ref<PullRequest[]>([])
  const reviewPRs = ref<PullRequest[]>([])
  const teamPRs = ref<PullRequest[]>([])
  const loading = ref(false)
  const connected = ref(false)

  async function fetchMyPRs() {
    try {
      if (!useMock) {
        const { ListMyPRs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
        myPRs.value = (await ListMyPRs()) as PullRequest[]
        connected.value = true
      } else {
        await new Promise(r => setTimeout(r, 150))
        myPRs.value = [...MOCK_MY_PRS]
      }
    } catch (e) {
      console.warn('[PRStore] Wails binding unavailable for ListMyPRs, using mock data:', e)
      useMock = true
      myPRs.value = [...MOCK_MY_PRS]
    }
  }

  async function fetchReviewPRs() {
    try {
      if (!useMock) {
        const { ListReviewPRs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
        reviewPRs.value = (await ListReviewPRs()) as PullRequest[]
        connected.value = true
      } else {
        await new Promise(r => setTimeout(r, 150))
        reviewPRs.value = [...MOCK_REVIEW_PRS]
      }
    } catch (e) {
      console.warn('[PRStore] Wails binding unavailable for ListReviewPRs, using mock data:', e)
      useMock = true
      reviewPRs.value = [...MOCK_REVIEW_PRS]
    }
  }

  async function fetchTeamPRs() {
    try {
      if (!useMock) {
        const { ListTeamPRs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
        teamPRs.value = (await ListTeamPRs()) as PullRequest[]
        connected.value = true
      } else {
        await new Promise(r => setTimeout(r, 150))
        teamPRs.value = [...MOCK_TEAM_PRS]
      }
    } catch (e) {
      console.warn('[PRStore] Wails binding unavailable for ListTeamPRs, using mock data:', e)
      useMock = true
      teamPRs.value = [...MOCK_TEAM_PRS]
    }
  }

  async function fetchAll() {
    loading.value = true
    try {
      await Promise.all([fetchMyPRs(), fetchReviewPRs(), fetchTeamPRs()])
    } finally {
      loading.value = false
    }
  }

  return {
    myPRs, reviewPRs, teamPRs, loading, connected,
    fetchMyPRs, fetchReviewPRs, fetchTeamPRs, fetchAll,
  }
})
