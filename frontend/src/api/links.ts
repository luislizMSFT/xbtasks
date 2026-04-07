export async function linkTask(taskId: number, adoId: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
  return m.LinkTask(taskId, adoId)
}

export async function unlinkTask(taskId: number, adoId: string, deleteWorkItem: boolean) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
  return m.UnlinkTask(taskId, adoId, deleteWorkItem)
}

export async function promoteTask(taskId: number, wiType: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
  return m.PromoteTask(taskId, wiType)
}

export async function importWorkItem(adoId: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
  return m.ImportWorkItem(adoId)
}

export async function importWorkItemAsProject(adoId: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
  return m.ImportWorkItemAsProject(adoId)
}

export async function listPublicTaskIDs() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
  return m.ListPublicTaskIDs()
}

export async function listLinkedAdoIDs() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
  return m.ListLinkedAdoIDs()
}
