export async function openURL(url: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/browserservice')
  return m.OpenURL(url)
}
