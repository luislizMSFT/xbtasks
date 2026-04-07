export async function listRecentRuns() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/pipelineservice')
  return m.ListRecentRuns()
}
