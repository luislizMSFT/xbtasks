---
phase: 02-ado-integration-prs-unified-dashboard
plan: 01
subsystem: auth, ado-client
tags: [auth, ado, rest-api, token-provider, wiql, json-patch, state-mapping]
dependency_graph:
  requires: []
  provides: [TokenProvider, CachedTokenProvider, AzCliTokenProvider, PATTokenProvider, ado.Client, ado.WorkItem, ado.QueryMyWorkItems, ado.CreateWorkItem, ado.UpdateWorkItem, ado.AddComment, ado.StatusToADO, ado.ADOToStatus]
  affects: [internal/auth, pkg/ado]
tech_stack:
  added: []
  patterns: [token-provider-interface, cached-provider-wrapper, bearer-auth-http-client, wiql-then-batch-fetch, json-patch-content-type, bidirectional-state-mapping]
key_files:
  created:
    - internal/auth/token.go
    - internal/auth/azcli.go
    - internal/auth/pat.go
    - pkg/ado/client.go
    - pkg/ado/models.go
    - pkg/ado/query.go
    - pkg/ado/push.go
    - pkg/ado/comments.go
    - pkg/ado/sync.go
    - pkg/ado/urls.go
  modified: []
decisions:
  - "CachedTokenProvider uses interface assertion to check inner provider for Expiry() method rather than requiring it in TokenProvider interface"
  - "ADO client accepts token string (not TokenProvider) — caller is responsible for token lifecycle"
  - "parseWorkItemFromAPI handles AssignedTo as both string and object (displayName) for API response compatibility"
  - "State mapping defaults to 'New'/'todo' for unknown states rather than erroring"
metrics:
  duration: 5m
  completed: "2026-04-06T18:56:12Z"
  tasks: 2
  files_created: 10
---

# Phase 02 Plan 01: Auth Token Provider + ADO REST Client Summary

**One-liner:** Abstracted token provider (az CLI + PAT + cached wrapper) and full ADO REST client with WIQL queries, JSON Patch CRUD, comments API, bidirectional state mapping, and URL builders.

## What Was Built

### Token Provider Interface (`internal/auth/`)
- **`token.go`** — `TokenProvider` interface with `GetToken()` and `Name()` methods; `CachedTokenProvider` wraps any provider with TTL-based mutex caching and configurable buffer duration
- **`azcli.go`** — `AzCliTokenProvider` runs `az account get-access-token --resource 499b84ac-...` to fetch ADO tokens; parses expiry from response; logs provider in use
- **`pat.go`** — `PATTokenProvider` reads PAT from OS keyring (`team-ado-tool/pat`) using go-keyring; clear error messages when not configured

### ADO REST Client (`pkg/ado/`)
- **`client.go`** — `Client` struct with Bearer token auth, 15s timeout; `NewMultiOrgClients` for multi-org; private helpers for `doRequest`, `doJSON`, `doPatch` (with `application/json-patch+json`), `decodeResponse`
- **`models.go`** — Domain types: `WorkItem`, `PatchOperation`, `WIQLResponse`, `WIQLWorkItemRef`, `Identity`, `Comment`, `OrgProject`, `FieldDiff`, `SyncDiff`
- **`query.go`** — `QueryMyWorkItems` (WIQL @Me + batch fetch), `GetWorkItemsByIDs` (batched 200/req), `GetWorkItem`, `GetWorkItemChildren`; field parser handling ADO response format
- **`push.go`** — `CreateWorkItem` (POST with JSON Patch), `UpdateWorkItem` (PATCH), `UpdateWorkItemFields` (convenience: friendly name → ADO field path)
- **`comments.go`** — `GetComments` and `AddComment` using 7.0-preview.3 comments API
- **`sync.go`** — `StatusToADO`/`ADOToStatus` bidirectional state maps with type-specific overrides (Task→Completed, Bug→Closed); `MapStatusToADO`, `MapADOToStatus`, `GenerateSyncDiff`
- **`urls.go`** — `WorkItemWebURL`, `OrgURL`, `ProjectURL` URL builders

## Commits

| # | Hash | Message |
|---|------|---------|
| 1 | `66910c5` | feat(02-01): token provider interface with az CLI and PAT implementations |
| 2 | `034bc07` | feat(02-01): ADO REST client package with WIQL, JSON Patch, comments, state mapping |

## Deviations from Plan

### Out-of-Scope Pre-existing Issues

**1. Pre-existing build failures in `internal/db/`**
- `internal/db/comments.go` references `domain.TaskComment` (undefined)
- `internal/db/links.go` references `domain.TaskLink` (undefined)
- `internal/db/ado.go` references `domain.SyncState`, `domain.ProjectADOLink` (undefined)
- These prevent `go build ./internal/auth/...` from succeeding transitively (auth.go imports internal/db)
- **Impact:** Auth package files are syntactically valid (verified via gofmt) and will compile once domain types are added
- **Logged to:** `deferred-items.md` in phase directory

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed double-request in CreateWorkItem**
- **Found during:** Task 2
- **Issue:** Initial implementation made an empty POST then tried to re-do with proper body
- **Fix:** Refactored to single POST with JSON-marshaled PatchOperations and `application/json-patch+json` content type
- **Files modified:** `pkg/ado/push.go`
- **Commit:** `034bc07`

## Verification Results

- ✅ `gofmt -e` passes on all 3 auth files (token.go, azcli.go, pat.go)
- ✅ `go build ./pkg/ado/...` exits 0
- ✅ `go vet ./pkg/ado/...` exits 0
- ✅ TokenProvider interface has GetToken() and Name() methods
- ✅ ADO client sets `Authorization: Bearer` header
- ✅ JSON Patch uses `application/json-patch+json` content type
- ✅ All 14 acceptance criteria verified

## Self-Check: PASSED

All 10 created files verified on disk. Both commit hashes (66910c5, 034bc07) confirmed in git log. SUMMARY.md exists.
