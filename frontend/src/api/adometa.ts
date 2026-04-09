// ADO metadata cache API wrapper
// Maps to internal/app/adometa.go ADOMetaCacheService Wails bindings

export interface AdoMeta {
  type: string
  state: string
}

export async function getAllADOMeta(): Promise<Record<number, AdoMeta>> {
  try {
    const m = await import(
      '../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adometacacheservice'
    )
    return (await m.GetAll()) as Record<number, AdoMeta>
  } catch {
    console.warn('[adometa] Bindings not available, returning empty cache')
    return {}
  }
}

export async function refreshADOMeta(): Promise<void> {
  try {
    const m = await import(
      '../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/adometacacheservice'
    )
    await m.Refresh()
  } catch {
    console.warn('[adometa] Refresh binding not available')
  }
}
