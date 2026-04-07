export async function signIn() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
  return m.SignIn()
}

export async function signInWithPAT(pat: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
  return m.SignInWithPAT(pat)
}

export async function signInWithAzCli() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
  return m.SignInWithAzCli()
}

export async function signOut() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
  return m.SignOut()
}

export async function tryRestoreSession() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
  return m.TryRestoreSession()
}

export async function getCurrentUser() {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/auth/authservice')
  return m.GetCurrentUser()
}
