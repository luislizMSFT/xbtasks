export async function listMyWorkItems() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
  return m.ListMyWorkItems()
}

export async function getCachedWorkItem(adoId: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
  return m.GetCachedWorkItem(adoId)
}

export async function getCachedWorkItems() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
  return m.GetCachedWorkItems()
}

export async function syncWorkItems() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
  return m.SyncWorkItems()
}

export async function getWorkItemTree() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
  return m.GetWorkItemTree()
}

export async function getSavedQueries() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
  return m.GetSavedQueries()
}

export async function runSavedQuery(queryId: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adoservice')
  return m.RunSavedQuery(queryId)
}
