export async function listTasks(status = '') {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.List(status)
}
export async function createTask(title: string, desc: string, priority: string, category: string, projectId: number | null, parentId: number | null) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.Create(title, desc, priority, category, projectId, parentId)
}
export async function updateTask(id: number, title: string, description: string, status: string, priority: string, category: string, area: string, dueDate: string, tags: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.Update(id, title, description, status, priority, category, area, dueDate, tags)
}
export async function setStatus(id: number, status: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.SetStatus(id, status)
}
export async function deleteTask(id: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.Delete(id)
}
export async function reorderTasks(orderedIds: number[]) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.ReorderTasks(orderedIds)
}
export async function setPersonalPriority(id: number, priority: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.SetPersonalPriority(id, priority)
}
export async function getSubtasks(parentID: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.GetSubtasks(parentID)
}
export async function createSubtask(parentId: number, title: string, desc: string, priority: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.CreateSubtask(parentId, title, desc, priority)
}
export async function getAllTags() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
  return m.GetAllTags()
}
