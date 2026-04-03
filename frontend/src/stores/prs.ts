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
  createdBy: string
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

export const usePRStore = defineStore('prs', () => {
  const myPRs = ref<PullRequest[]>([])
  const reviewPRs = ref<PullRequest[]>([])
  const teamPRs = ref<PullRequest[]>([])
  const loading = ref(false)
  const connected = ref(false)
  const error = ref('')

  async function fetchMyPRs() {
    try {
      const { ListMyPRs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
      myPRs.value = (await ListMyPRs()) as PullRequest[]
      connected.value = true
    } catch (e: any) {
      console.warn('[PRStore] ListMyPRs failed:', e)
      error.value = e?.message || 'Failed to fetch PRs'
      myPRs.value = []
    }
  }

  async function fetchReviewPRs() {
    try {
      const { ListReviewPRs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
      reviewPRs.value = (await ListReviewPRs()) as PullRequest[]
      connected.value = true
    } catch (e: any) {
      console.warn('[PRStore] ListReviewPRs failed:', e)
      reviewPRs.value = []
    }
  }

  async function fetchTeamPRs() {
    try {
      const { ListTeamPRs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
      teamPRs.value = (await ListTeamPRs()) as PullRequest[]
      connected.value = true
    } catch (e: any) {
      console.warn('[PRStore] ListTeamPRs failed:', e)
      teamPRs.value = []
    }
  }

  async function fetchAll() {
    loading.value = true
    error.value = ''
    try {
      await Promise.all([fetchMyPRs(), fetchReviewPRs(), fetchTeamPRs()])
    } finally {
      loading.value = false
    }
  }

  return {
    myPRs, reviewPRs, teamPRs, loading, connected, error,
    fetchMyPRs, fetchReviewPRs, fetchTeamPRs, fetchAll,
  }
})
