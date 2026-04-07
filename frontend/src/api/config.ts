export async function getAll() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
  return m.GetAll()
}

export async function get(key: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
  return m.Get(key)
}

export async function set(key: string, value: any) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
  return m.Set(key, value)
}

export async function getOrgProjects() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
  return m.GetOrgProjects()
}

export async function setOrgProjects(orgs: any[]) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
  return m.SetOrgProjects(orgs)
}

export async function getSyncInterval() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
  return m.GetSyncInterval()
}

export async function setSyncInterval(minutes: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/config/configservice')
  return m.SetSyncInterval(minutes)
}
