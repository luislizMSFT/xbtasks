export async function manualSync() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
  return m.ManualSync()
}

export async function generateOutboundDiff(taskId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
  return m.GenerateOutboundDiff(taskId)
}

export async function pushChanges(taskId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
  return m.PushChanges(taskId)
}

export async function resolveConflict(taskId: number, resolutions: Record<string, string>) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/syncservice')
  return m.ResolveConflict(taskId, resolutions)
}
