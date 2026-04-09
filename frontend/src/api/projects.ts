export async function listProjects() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.List()
}

export async function createProject(name: string, description: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.Create(name, description)
}

export async function deleteProject(id: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.Delete(id)
}

export async function pinProject(id: number, pinned: boolean) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.PinProject(id, pinned)
}

export async function linkProjectToADO(projectId: number, adoId: string, direction: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.LinkProjectToADO(projectId, adoId, direction)
}

export async function unlinkProject(projectId: number, adoId: string, deleteWorkItem: boolean) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.UnlinkProject(projectId, adoId, deleteWorkItem)
}

export async function getProjectADOLink(projectId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.GetProjectADOLink(projectId)
}

export async function getProjectProgress(projectId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.GetProjectProgress(projectId)
}

export async function updateProject(id: number, name: string, description: string, status: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/projectservice')
  return m.Update(id, name, description, status)
}
