export async function listMyPRs() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
  return m.ListMyPRs()
}

export async function listReviewPRs() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
  return m.ListReviewPRs()
}

export async function listTeamPRs() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
  return m.ListTeamPRs()
}

export async function dismissPR(prNumber: number, repo: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
  return m.DismissPR(prNumber, repo)
}

export async function undismissPR(prNumber: number, repo: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
  return m.UndismissPR(prNumber, repo)
}

export async function watchPR(prNumber: number, repo: string, watched: boolean) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/prservice')
  return m.WatchPR(prNumber, repo, watched)
}
