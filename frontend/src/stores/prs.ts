import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { PRReviewer, PullRequest } from '@/types'
import * as prsApi from '@/api/prs'

export type { PRReviewer, PullRequest }

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

export { relativeTime } from '@/lib/date'

export const usePRStore = defineStore('prs', () => {
  const myPRs = ref<PullRequest[]>([])
  const reviewPRs = ref<PullRequest[]>([])
  const teamPRs = ref<PullRequest[]>([])
  const loading = ref(false)
  const connected = ref(false)
  const error = ref('')
  const showDismissed = ref(false)

  let fetchInFlight = false

  async function fetchMyPRs() {
    try {
      myPRs.value = (await prsApi.listMyPRs()) as PullRequest[]
      connected.value = true
    } catch (e: any) {
      console.warn('[PRStore] ListMyPRs failed:', e)
      error.value = e?.message || 'Failed to fetch PRs'
      myPRs.value = []
    }
  }

  async function fetchReviewPRs() {
    try {
      reviewPRs.value = (await prsApi.listReviewPRs()) as PullRequest[]
      connected.value = true
    } catch (e: any) {
      console.warn('[PRStore] ListReviewPRs failed:', e)
      reviewPRs.value = []
    }
  }

  async function fetchTeamPRs() {
    try {
      teamPRs.value = (await prsApi.listTeamPRs()) as PullRequest[]
      connected.value = true
    } catch (e: any) {
      console.warn('[PRStore] ListTeamPRs failed:', e)
      teamPRs.value = []
    }
  }

  async function fetchAll() {
    if (fetchInFlight) return
    fetchInFlight = true
    loading.value = true
    error.value = ''
    try {
      await Promise.all([fetchMyPRs(), fetchReviewPRs(), fetchTeamPRs()])
    } finally {
      loading.value = false
      fetchInFlight = false
    }
  }

  async function dismissPR(pr: PullRequest) {
    await prsApi.dismissPR(pr.prNumber, pr.repo)
    pr.dismissedAt = new Date().toISOString()
  }

  async function undismissPR(pr: PullRequest) {
    await prsApi.undismissPR(pr.prNumber, pr.repo)
    pr.dismissedAt = null
  }

  async function toggleWatchPR(pr: PullRequest) {
    const newVal = !pr.watched
    await prsApi.watchPR(pr.prNumber, pr.repo, newVal)
    pr.watched = newVal
  }

  return {
    myPRs, reviewPRs, teamPRs, loading, connected, error, showDismissed,
    fetchMyPRs, fetchReviewPRs, fetchTeamPRs, fetchAll,
    dismissPR, undismissPR, toggleWatchPR,
  }
})
