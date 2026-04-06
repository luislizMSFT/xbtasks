---
phase: 02-ado-integration-prs-unified-dashboard
plan: 05
subsystem: frontend-auth-settings
tags: [auth, login-ui, az-cli, multi-org, settings, config]
dependency_graph:
  requires: [02-03, 02-04]
  provides: [login-3-auth-methods, multi-org-settings, sync-interval-ui]
  affects: [frontend/src/views/LoginView.vue, frontend/src/stores/auth.ts, frontend/src/views/SettingsView.vue, internal/auth/auth.go]
tech_stack:
  added: []
  patterns: [dynamic-wails-imports, multi-org-config-ui, auth-method-tracking]
key_files:
  created: []
  modified:
    - frontend/src/views/LoginView.vue
    - frontend/src/stores/auth.ts
    - internal/auth/auth.go
    - frontend/src/views/SettingsView.vue
decisions:
  - "Az CLI auth fetches ADO profile for real user info, falls back to generic user"
  - "Auth method tracked via authMethod ref in store + keyring for session restore"
  - "Settings page replaced legacy single-org inputs with multi-org list management"
  - "Sync interval saved on change via ConfigService binding"
  - "Error message in login now displays actual error text from auth store"
metrics:
  duration: 4m
  completed: "2026-04-06T19:20:46Z"
  tasks_completed: 2
  tasks_total: 2
  files_modified: 4
---

# Phase 02 Plan 05: Login UI Az CLI Option & Multi-Org Settings Summary

Az CLI as third login auth method with ADO profile fetch; settings page rebuilt with multi-org management, sync interval, and auth section.

## What Was Done

### Task 1: Login UI — Add Az CLI Token Option + Wire Auth Store
- **LoginView.vue**: Added 3 auth options in order: Microsoft OAuth (primary), Az CLI Token (outline), PAT (ghost). Az CLI button uses Terminal icon from lucide-vue-next with helper text "Requires az login — run in terminal first".
- **auth.ts**: Added `signInWithAzCli()` and `signInWithPAT(pat)` methods with Wails binding dynamic imports and mock fallbacks. Added `authMethod` ref tracking which method was used (oauth/pat/azcli). All three methods exported in store return.
- **auth.go**: Added `SignInWithAzCli()` method that uses `AzCliTokenProvider` to get a token, then fetches user profile from ADO VSSPS API. Added `fetchADOProfile()` helper. Updated `TryRestoreSession()` to handle azcli auth method. Updated `SignOut()` to clean up auth_method keyring entry.
- **Error display**: Login error message now shows actual error text from store instead of generic string.

### Task 2: Settings Page — Multi-Org Configuration + Sync Interval
- **ADO Organizations section**: Displays configured org/project pairs as styled rows with org name, comma-separated projects, and trash button. Add form with org name + projects inputs, deduplication on add.
- **Sync Settings section**: Background sync interval with number input, saves via `SetSyncInterval` on change.
- **Authentication section**: Shows current auth method label, user avatar/initials, PAT update flow (for PAT users), and sign out button (destructive variant).
- **Removed**: Legacy `adoOrganization` and `adoProject` single-value inputs replaced by multi-org list management via `GetOrgProjects`/`SetOrgProjects`.
- **Kept**: Appearance (theme, window size), Data (db path), and About sections unchanged.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed error message display in LoginView**
- **Found during:** Task 1
- **Issue:** Login error showed hardcoded "Authentication failed" text instead of actual error
- **Fix:** Changed to `{{ authStore.error }}` to display actual error message
- **Files modified:** frontend/src/views/LoginView.vue
- **Commit:** 9cd685c

**2. [Rule 2 - Missing functionality] Added ADO profile fetch for Az CLI auth**
- **Found during:** Task 1
- **Issue:** Plan suggested generic "azcli-user" but the token should be used to get real user info
- **Fix:** Added `fetchADOProfile()` that calls ADO VSSPS profile API, with fallback to generic user
- **Files modified:** internal/auth/auth.go
- **Commit:** 9cd685c

**3. [Rule 2 - Missing functionality] Added auth method to session restore**
- **Found during:** Task 1
- **Issue:** `TryRestoreSession` didn't handle azcli method — would fail to restore
- **Fix:** Added auth_method keyring check at top of TryRestoreSession, delegates to SignInWithAzCli
- **Files modified:** internal/auth/auth.go
- **Commit:** 9cd685c

## Commits

| # | Hash | Message |
|---|------|---------|
| 1 | 9cd685c | feat(02-05): add Az CLI token as third auth option in login UI |
| 2 | 916df79 | feat(02-05): settings page with multi-org config, sync interval, auth section |

## Verification

- `go build ./internal/auth/...` exits 0
- LoginView.vue contains all three auth options with correct button variants
- auth.ts exports signIn, signInWithPAT, signInWithAzCli, authMethod
- SettingsView.vue has GetOrgProjects/SetOrgProjects/GetSyncInterval/SetSyncInterval bindings
- All acceptance criteria from both tasks verified

## Self-Check: PASSED

All files exist, all commits verified.
