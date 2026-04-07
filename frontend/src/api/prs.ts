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
