export async function listLinks(taskId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/externallinksservice')
  return m.ListLinks(taskId)
}

export async function addLink(taskId: number, url: string, label: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/externallinksservice')
  return m.AddLink(taskId, url, label)
}

export async function deleteLink(id: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/externallinksservice')
  return m.DeleteLink(id)
}
