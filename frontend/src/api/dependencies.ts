export async function getDependencies(taskID: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/dependencyservice')
  return m.GetDependencies(taskID)
}

export async function addDependency(taskID: number, dependsOn: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/dependencyservice')
  return m.AddDependency(taskID, dependsOn)
}

export async function removeDependency(taskID: number, dependsOn: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/dependencyservice')
  return m.RemoveDependency(taskID, dependsOn)
}
