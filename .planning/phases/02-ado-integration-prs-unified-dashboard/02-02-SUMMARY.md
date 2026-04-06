---
phase: 02-ado-integration-prs-unified-dashboard
plan: 02
subsystem: data-layer
tags: [schema, domain-types, config, multi-org, sqlite]
dependency_graph:
  requires: []
  provides: [task_links-table, task_comments-table, project_ado_links-table, sync_state-table, TaskLink-type, TaskComment-type, ProjectADOLink-type, SyncState-type, OrgProject-type, multi-org-config]
  affects: [internal/db, domain/types, internal/config, internal/app/adoservice]
tech_stack:
  added: []
  patterns: [idempotent-schema-migration, ON-CONFLICT-upsert, legacy-config-fallback]
key_files:
  created:
    - internal/db/comments.go
    - internal/db/links.go
  modified:
    - internal/db/db.go
    - internal/db/ado.go
    - domain/types.go
    - internal/config/config.go
    - internal/config/service.go
decisions:
  - "New domain types use same JSON tag pattern (camelCase) as existing types"
  - "Multi-org config falls back to legacy single-org when ado.orgs is empty"
  - "SyncState uses composite PK (task_id, ado_id) for per-link conflict tracking"
  - "is_public stored as INTEGER (0/1) in SQLite, mapped to bool in Go"
metrics:
  duration: 4m
  completed: 2026-04-06
---

# Phase 02 Plan 02: Schema & Domain Types Summary

Extended database schema with 4 new tables (task_links, task_comments, project_ado_links, sync_state), added 5 domain types, and multi-org ADO config with legacy fallback.

## What Was Done

### Task 1: Schema Migration — New Tables + DB Access Methods
**Commit:** `54aefcb`

Extended `internal/db/db.go` schema with 4 new tables:
- **task_links** — external URLs (ICM, Grafana, ADO, wiki) attached to tasks
- **task_comments** — local comments with is_public flag and ado_comment_id for ADO sync
- **project_ado_links** — project-to-ADO work item linking with direction (promoted/imported/linked)
- **sync_state** — per-field last-known state for conflict detection

Added 5 new columns to existing tables:
- `projects.is_pinned` (INTEGER DEFAULT 0)
- `ado_work_items.org`, `.project`, `.parent_id`, `.changed_date`

Created `internal/db/comments.go` with CRUD: CreateComment, ListComments, UpdateComment, DeleteComment, MarkCommentPublic

Created `internal/db/links.go` with CRUD: CreateLink, ListLinks, DeleteLink

Extended `internal/db/ado.go` with: UpsertSyncState, GetSyncState, ListSyncStates, CreateProjectADOLink, DeleteProjectADOLink, GetProjectADOLink, ListADOWorkItemsByOrg — plus updated existing methods to include new org/project/parent_id/changed_date fields.

### Task 2: Domain Types + Multi-Org Config Extension
**Commit:** `21c61d5`

Added to `domain/types.go`:
- **TaskLink** — external URL with type detection (url, icm, grafana, ado, wiki)
- **TaskComment** — comment with IsPublic bool + AdoCommentID
- **ProjectADOLink** — project ↔ ADO link with direction
- **SyncState** — per-field local/remote state snapshot for conflict detection
- **OrgProject** — org + projects list with mapstructure tags for Viper

Extended existing types:
- **ADOWorkItem** — added Org, Project, ParentID, ChangedDate fields
- **Project** — added IsPinned field

Extended `internal/config/config.go`:
- `GetOrgProjects()` — reads `ado.orgs`, falls back to legacy `ado.organization` + `ado.project`
- `SetOrgProjects()` — writes multi-org config to disk

Extended `internal/config/service.go`:
- GetOrgProjects, SetOrgProjects, GetSyncInterval, SetSyncInterval methods on ConfigService

## Deviations from Plan

None — plan executed exactly as written.

## Verification

- `go build ./domain/...` ✅
- `go build ./internal/config/...` ✅
- `go build ./internal/db/...` ✅
- `go build ./internal/app/...` ✅ (no regressions in existing ADOService/ProjectService)

## Self-Check: PASSED
