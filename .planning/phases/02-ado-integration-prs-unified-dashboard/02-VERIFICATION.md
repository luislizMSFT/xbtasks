---
phase: 02-ado-integration-prs-unified-dashboard
verified: 2026-04-06T19:33:05Z
status: gaps_found
score: 13/14 must-haves verified
gaps:
  - truth: "ConflictResolver component is mounted in the component tree and renders when conflicts are detected"
    status: failed
    reason: "ConflictResolver.vue exists and is fully implemented, but is not imported or mounted in any parent component (App.vue, TaskDetail.vue, TasksView.vue, etc). The sync store correctly sets showConflictResolver=true on conflict, but no Dialog renders because the component is orphaned."
    artifacts:
      - path: "frontend/src/components/ConflictResolver.vue"
        issue: "Orphaned — not imported/mounted in any parent component"
    missing:
      - "Import ConflictResolver in TaskDetail.vue (or App.vue) and add <ConflictResolver /> to the template alongside the existing <SyncConfirmDialog />"
---

# Phase 02: ADO Integration & Sync Workflow Verification Report

**Phase Goal:** User can authenticate via az cli, browse ADO items, link/promote/import tasks with the personal→public model, and bidirectionally sync with safe confirmation before any outbound changes
**Verified:** 2026-04-06T19:33:05Z
**Status:** gaps_found
**Re-verification:** No — initial verification

## Goal Achievement

### Observable Truths

| #  | Truth | Status | Evidence |
|----|-------|--------|----------|
| 1  | TokenProvider interface exists with az cli, PAT, and cached wrapper implementations | ✓ VERIFIED | `internal/auth/token.go` L10 interface, `azcli.go` L27 NewAzCliTokenProvider, `pat.go` L22 NewPATTokenProvider, `token.go` L28 NewCachedTokenProvider |
| 2  | Login UI shows three auth options (OAuth, Az CLI, PAT) | ✓ VERIFIED | `LoginView.vue` L83-98 Az CLI button with Terminal icon, L124-139 PAT input, existing OAuth button |
| 3  | ADO REST client uses direct HTTP with Bearer token (no az cli subprocess) | ✓ VERIFIED | `pkg/ado/client.go` L60 Authorization Bearer header, `adoservice.go` has NO exec.Command/runAzCli references |
| 4  | ADO browser shows tree view with filter/search/area/saved queries | ✓ VERIFIED | `AdoView.vue` imports AdoTreeBrowser, has filterArea/filterType/filterState/searchQuery, saved query picker L365-373 |
| 5  | Link/Promote/Import/Unlink flows all implemented | ✓ VERIFIED | `linkservice.go` L70 LinkTask, L114 PromoteTask, L169 ImportWorkItem, L215 UnlinkTask; dialogs wired L89/L44 in Vue |
| 6  | isPublic computed from task_ado_links table presence | ✓ VERIFIED | `linkservice.go` L245 IsPublic, `tasks.ts` L38-51 publicTaskIds + isPublic, `TaskRow.vue` L20 isPublic prop |
| 7  | Background sync pulls silently on configurable timer | ✓ VERIFIED | `syncservice.go` L36 StartBackgroundSync, L43 time.NewTicker, `main.go` L123 StartBackgroundSync() |
| 8  | Outbound push requires confirmation via preview diff dialog | ✓ VERIFIED | `SyncConfirmDialog.vue` mounted in TaskDetail.vue L740, shows pendingDiff field-by-field, confirmPush/cancelPush wired |
| 9  | Per-field conflict resolution UI shows local vs ADO with pick buttons | ✗ FAILED | `ConflictResolver.vue` fully implemented (L74 local, L88 remote, L112 resolve button) BUT orphaned — not imported/mounted in any parent component |
| 10 | Comments local-first with selective ADO push | ✓ VERIFIED | `CommentsSection.vue` L120 Push to ADO for private comments on public tasks, `commentservice.go` L56 PushCommentToADO |
| 11 | External links with auto-detected types open in real browser | ✓ VERIFIED | `ExternalLinks.vue` L84-85 window.open, `externallinks.go` L26-28 ICM/Grafana/ADO/Wiki patterns |
| 12 | Projects page with card grid, pin, ADO linking, dual progress | ✓ VERIFIED | `ProjectCard.vue` L42 pin, L51 ADO badge, L72-93 dual progress bars; `ProjectDetailView.vue` L173 PROJ-06 dual bars |
| 13 | Multi-org configuration via settings | ✓ VERIFIED | `SettingsView.vue` L34-45 GetOrgProjects/SetOrgProjects, `config.go` L121 GetOrgProjects with legacy fallback |
| 14 | Unified task list with filter bar (status, priority, project, due date, ADO link) | ✓ VERIFIED | `FilterBar.vue` all 5 filter dimensions + sort + groupBy, `TasksView.vue` L212 FilterBar integrated, L247 quick-add |

**Score:** 13/14 truths verified

### Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `internal/auth/token.go` | TokenProvider interface + CachedTokenProvider | ✓ VERIFIED | 1861 bytes, interface L10, cached L18 |
| `internal/auth/azcli.go` | Az CLI token provider | ✓ VERIFIED | 2134 bytes, adoResourceID L12, exec.Command L36 |
| `internal/auth/pat.go` | PAT token provider | ✓ VERIFIED | 1105 bytes, keyring.Get L33 |
| `pkg/ado/client.go` | ADO HTTP client with auth | ✓ VERIFIED | 3488 bytes, Bearer auth L60, json-patch+json L89 |
| `pkg/ado/models.go` | ADO domain types | ✓ VERIFIED | 2284 bytes, WorkItem/PatchOp/SyncDiff/FieldDiff/OrgProject |
| `pkg/ado/query.go` | WIQL + batch fetch | ✓ VERIFIED | 7741 bytes, QueryMyWorkItems L12, @Me query, batch fetch L41 |
| `pkg/ado/push.go` | Create/update work items | ✓ VERIFIED | 2998 bytes, CreateWorkItem L11, UpdateWorkItemFields L84 |
| `pkg/ado/comments.go` | Comments API | ✓ VERIFIED | 1322 bytes, GetComments L11, AddComment L28 |
| `pkg/ado/sync.go` | State mapping + diff generation | ✓ VERIFIED | 2556 bytes, StatusToADO/ADOToStatus maps, GenerateSyncDiff L59 |
| `pkg/ado/urls.go` | URL builders | ✓ VERIFIED | 572 bytes, WorkItemWebURL L6, dev.azure.com L7 |
| `internal/db/db.go` | Extended schema | ✓ VERIFIED | 6769 bytes, 4 new tables (task_links/comments/project_ado_links/sync_state), ALTER TABLE migrations |
| `internal/db/comments.go` | Comment CRUD | ✓ VERIFIED | 2195 bytes, CreateComment/ListComments/MarkCommentPublic |
| `internal/db/links.go` | External links CRUD | ✓ VERIFIED | 1370 bytes, CreateLink/ListLinks/DeleteLink |
| `internal/db/ado.go` | Extended ADO DB ops | ✓ VERIFIED | 6055 bytes, UpsertADOWorkItem/UpsertSyncState/CreateProjectADOLink/ListSyncStates |
| `internal/config/config.go` | Multi-org config | ✓ VERIFIED | 4036 bytes, GetOrgProjects L121 with legacy fallback, SetOrgProjects L135 |
| `internal/config/service.go` | Config service bindings | ✓ VERIFIED | 1414 bytes, GetOrgProjects/SetOrgProjects/GetSyncInterval/SetSyncInterval |
| `domain/types.go` | New domain types | ✓ VERIFIED | 6228 bytes, TaskLink/TaskComment/ProjectADOLink/SyncState/OrgProject + IsPinned on Project |
| `internal/app/adoservice.go` | Refactored ADO service (REST) | ✓ VERIFIED | 7330 bytes, tokenProv L19, ado.QueryMyWorkItems L84, NO exec.Command |
| `internal/app/linkservice.go` | Link/Promote/Import/Unlink | ✓ VERIFIED | 12290 bytes, all 4 flows + IsPublic + ListPublicTaskIDs |
| `internal/app/syncservice.go` | Background sync + diff + conflicts | ✓ VERIFIED | 17976 bytes, StartBackgroundSync/ManualSync/GenerateOutboundDiff/PushChanges/ResolveConflict/DetectConflicts |
| `internal/app/commentservice.go` | Comment CRUD + ADO push | ✓ VERIFIED | 4714 bytes, AddComment/PushCommentToADO with ado.AddComment |
| `internal/app/externallinks.go` | Links + URL detection | ✓ VERIFIED | 1836 bytes, DetectLinkType with ICM/Grafana/ADO/Wiki patterns |
| `internal/app/projects.go` | Extended projects | ✓ VERIFIED | 5441 bytes, LinkProjectToADO/UnlinkProject/PinProject/GetProjectProgress |
| `main.go` | Service wiring | ✓ VERIFIED | 4225 bytes, token chain L35-36, all services registered L63-74, StartBackgroundSync L123 |
| `frontend/src/views/LoginView.vue` | 3 auth options | ✓ VERIFIED | 5679 bytes, OAuth/AzCLI/PAT buttons |
| `frontend/src/stores/auth.ts` | Extended auth store | ✓ VERIFIED | 3798 bytes, signInWithAzCli/signInWithPAT/authMethod |
| `frontend/src/views/SettingsView.vue` | Multi-org config UI | ✓ VERIFIED | 16473 bytes, org management + sync interval + auth section |
| `frontend/src/views/TasksView.vue` | Full task list | ✓ VERIFIED | 15634 bytes, FilterBar + quickAdd + enhancedFilteredTasks + isPublic |
| `frontend/src/components/TaskRow.vue` | Enhanced task row | ✓ VERIFIED | 6455 bytes, isPublic prop, AzureDevOpsIcon, PriorityBadge |
| `frontend/src/components/FilterBar.vue` | Reusable filter component | ✓ VERIFIED | 7479 bytes, 5 filter dimensions + sort + groupBy + sync button |
| `frontend/src/stores/tasks.ts` | Extended task store | ✓ VERIFIED | 11273 bytes, publicTaskIds/isPublic/enhancedFilteredTasks/quickAdd/groupBy |
| `frontend/src/views/AdoView.vue` | ADO browser | ✓ VERIFIED | 32184 bytes, tree browser + filters + saved queries + import/link |
| `frontend/src/components/AdoTreeBrowser.vue` | Tree view component | ✓ VERIFIED | 8213 bytes, recursive AdoTreeNode, expand/collapse, linked indicator |
| `frontend/src/components/LinkDialog.vue` | Link to ADO dialog | ✓ VERIFIED | 6784 bytes, search by ID/title, LinkTask binding |
| `frontend/src/components/PromoteDialog.vue` | Promote to ADO dialog | ✓ VERIFIED | 5713 bytes, wiType selector, preview, PromoteTask binding |
| `frontend/src/stores/ado.ts` | Extended ADO store | ✓ VERIFIED | 7546 bytes, workItemTree/GetWorkItemTree/savedQueries |
| `frontend/src/components/SyncConfirmDialog.vue` | Push preview dialog | ✓ VERIFIED | 2678 bytes, field-by-field diff, confirm/cancel |
| `frontend/src/components/ConflictResolver.vue` | Conflict resolution UI | ⚠️ ORPHANED | 4122 bytes, fully implemented but not mounted in any parent component |
| `frontend/src/stores/sync.ts` | Sync state management | ✓ VERIFIED | 4031 bytes, manualSync/generateOutboundDiff/confirmPush/resolveConflict/conflicts |
| `frontend/src/components/ExternalLinks.vue` | Links with type icons | ✓ VERIFIED | 4031 bytes, AddLink/ListLinks/DeleteLink + openExternal(window.open) |
| `frontend/src/components/CommentsSection.vue` | Comments with public/private | ✓ VERIFIED | 5198 bytes, PushCommentToADO + Private/Public badges |
| `frontend/src/components/TaskDetail.vue` | Extended detail sidebar | ✓ VERIFIED | 30024 bytes, ExternalLinks L495, CommentsSection L502, SyncConfirmDialog L740, Push to ADO L509 |
| `frontend/src/views/ProjectsView.vue` | Project card grid | ✓ VERIFIED | 4877 bytes, ProjectCard, pinned/unpinned sections, grid layout |
| `frontend/src/views/ProjectDetailView.vue` | Project dashboard | ✓ VERIFIED | 10855 bytes, dual progress PROJ-06, ADO link/unlink, filtered tasks |
| `frontend/src/components/ProjectCard.vue` | Project card | ✓ VERIFIED | 3586 bytes, isPinned, ADO badge, dual progress bars |
| `frontend/src/stores/projects.ts` | Extended project store | ✓ VERIFIED | 4055 bytes, pinProject/linkProjectToADO/unlinkProject/GetProjectProgress |
| `frontend/src/router/index.ts` | Route config | ✓ VERIFIED | 1238 bytes, /projects and /projects/:id routes |

### Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `internal/auth/azcli.go` | az cli binary | exec.Command("az"...) | ✓ WIRED | L36 exec.Command with adoResourceID |
| `pkg/ado/client.go` | ADO REST API | Authorization: Bearer | ✓ WIRED | L60 Bearer token header, L89 json-patch+json |
| `internal/app/adoservice.go` | `pkg/ado/client.go` | ado.NewClient/ado.QueryMyWorkItems | ✓ WIRED | L84 ado.QueryMyWorkItems, config.GetOrgProjects L33 |
| `internal/app/adoservice.go` | `internal/auth/token.go` | tokenProv.GetToken() | ✓ WIRED | L29 tokenProv.GetToken() |
| `internal/app/linkservice.go` | `internal/db/db.go` | task_ado_links CRUD | ✓ WIRED | L79 INSERT task_ado_links, L245 SELECT COUNT task_ado_links |
| `internal/app/syncservice.go` | `pkg/ado/client.go` | REST calls for pull/push | ✓ WIRED | Uses ado package for queries and updates |
| `internal/app/syncservice.go` | `internal/db/ado.go` | UpsertSyncState | ✓ WIRED | Calls db.UpsertSyncState, db.GetSyncState, db.ListSyncStates |
| `internal/app/commentservice.go` | `pkg/ado/comments.go` | ado.AddComment | ✓ WIRED | L81 ado.AddComment(client, adoIDInt, comment.Content) |
| `main.go` | All services | RegisterService | ✓ WIRED | L63-74 all services registered, L35-36 token chain, L123 StartBackgroundSync |
| `frontend/src/stores/auth.ts` | Wails bindings | authservice imports | ✓ WIRED | L56 SignInWithPAT, L82 SignInWithAzCli dynamic imports |
| `frontend/src/stores/sync.ts` | Wails bindings | syncservice imports | ✓ WIRED | L41 ManualSync, L61 GenerateOutboundDiff, L71 PushChanges, L82 ResolveConflict |
| `frontend/src/stores/tasks.ts` | Wails bindings | LinkService.ListPublicTaskIDs | ✓ WIRED | L42 ListPublicTaskIDs from linkservice |
| `frontend/src/views/TasksView.vue` | FilterBar | filter state binding | ✓ WIRED | L212-229 FilterBar with all v-model bindings |
| `frontend/src/stores/ado.ts` | Wails bindings | ADOService.GetWorkItemTree | ✓ WIRED | L166 GetWorkItemTree dynamic import |
| `frontend/src/components/LinkDialog.vue` | Wails bindings | LinkService.LinkTask | ✓ WIRED | L89 LinkTask dynamic import |
| `frontend/src/components/SyncConfirmDialog.vue` | sync store | useSyncStore | ✓ WIRED | L11 useSyncStore, mounted in TaskDetail L740 |
| `frontend/src/components/ConflictResolver.vue` | sync store | useSyncStore | ✗ NOT WIRED | Component uses syncStore correctly but is NOT mounted in any parent |
| `frontend/src/components/ExternalLinks.vue` | Wails bindings | ExternalLinksService | ✓ WIRED | L51 ListLinks, L64 AddLink, L76 DeleteLink |
| `frontend/src/components/CommentsSection.vue` | Wails bindings | CommentService | ✓ WIRED | L32 ListComments, L45 AddComment, L57 PushCommentToADO |
| `frontend/src/components/TaskDetail.vue` | ExternalLinks + CommentsSection | Component imports | ✓ WIRED | L22-23 imports, L495/L502 usage, L509 Push to ADO button |
| `frontend/src/views/SettingsView.vue` | Wails bindings | ConfigService | ✓ WIRED | L34-45 GetOrgProjects/SetOrgProjects, L89-99 Sync interval |
| `frontend/src/stores/projects.ts` | Wails bindings | ProjectService | ✓ WIRED | L70-71 PinProject, L80-81 LinkProjectToADO, L108-109 GetProjectProgress |

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|-----------|-------------|--------|----------|
| AUTH-01 | 02-01, 02-05 | Token provider + az cli auth | ✓ SATISFIED | TokenProvider interface, AzCliTokenProvider, LoginView 3 options |
| AUTH-02 | 02-01 | Token auto-refresh | ✓ SATISFIED | CachedTokenProvider with TTL and buffer duration |
| AUTH-03 | 02-01 | Swappable token provider | ✓ SATISFIED | TokenProvider interface allows any implementation |
| TASK-08 | 02-03, 02-06 | Personal→public task model | ✓ SATISFIED | IsPublic from task_ado_links, badges in TaskRow |
| TASK-09 | 02-06 | Quick-add with just title | ✓ SATISFIED | TasksView L166-170 handleQuickAdd, tasks.ts quickAdd |
| ADO-01 | 02-03, 02-07 | View assigned ADO items | ✓ SATISFIED | QueryMyWorkItems with @Me, AdoView tree browser |
| ADO-02 | 02-07 | View work item details | ✓ SATISFIED | AdoView detail panel with fields |
| ADO-03 | 02-03 | Link personal task to ADO | ✓ SATISFIED | LinkTask in linkservice.go + LinkDialog.vue |
| ADO-04 | 02-03 | Promote task to ADO | ✓ SATISFIED | PromoteTask in linkservice.go + PromoteDialog.vue |
| ADO-05 | 02-03 | Import ADO item as local task | ✓ SATISFIED | ImportWorkItem in linkservice.go + AdoView import button |
| ADO-06 | 02-07 | ADO browser with linked status | ✓ SATISFIED | AdoTreeBrowser with isLinked, hideLinked toggle |
| ADO-07 | 02-03 | Unlink task from ADO | ✓ SATISFIED | UnlinkTask with deleteLocal option |
| ADO-08 | 02-01 | Direct REST API (no az cli per query) | ✓ SATISFIED | pkg/ado client, adoservice has NO exec.Command |
| ADO-09 | 02-02, 02-05 | Multiple org/project config | ✓ SATISFIED | GetOrgProjects with legacy fallback, SettingsView org manager |
| ADO-10 | 02-03, 02-07 | Multi-org unified list | ✓ SATISFIED | getClients iterates all orgs, AdoView tree shows all items |
| SYNC-01 | 02-04 | Background auto-sync on timer | ✓ SATISFIED | StartBackgroundSync with time.NewTicker, configurable interval |
| SYNC-02 | 02-04, 02-08 | Outbound push requires confirmation | ✓ SATISFIED | SyncConfirmDialog previews diff, confirmPush explicit action |
| SYNC-03 | 02-04, 02-08 | Per-field conflict resolution | ⚠️ PARTIAL | Backend ResolveConflict works, ConflictResolver.vue complete but ORPHANED (not mounted) |
| SYNC-04 | 02-04 | Sync only title/status/description | ✓ SATISFIED | GenerateSyncDiff compares title/state/description only |
| LINK-01 | 02-02 | External URL links on tasks | ✓ SATISFIED | task_links table, CreateLink, ExternalLinks.vue |
| LINK-02 | 02-04 | Auto-detect URL patterns | ✓ SATISFIED | DetectLinkType with ICM/Grafana/ADO/Wiki patterns |
| LINK-03 | 02-09 | Links in task detail with icons | ✓ SATISFIED | ExternalLinks.vue type icons, mounted in TaskDetail L495 |
| CMT-01 | 02-02, 02-09 | Local comments (private by default) | ✓ SATISFIED | CommentsSection "Comments are private by default", Private badge |
| CMT-02 | 02-04, 02-09 | Selective push to ADO | ✓ SATISFIED | PushCommentToADO, "Push to ADO" button for private comments on public tasks |
| CMT-03 | 02-09 | Description push preview | ✓ SATISFIED | GenerateOutboundDiff previews changes, SyncConfirmDialog shows diff |
| DASH-01 | 02-06 | Unified list personal + public | ✓ SATISFIED | enhancedFilteredTasks, isPublic badges, AzureDevOpsIcon |
| DASH-02 | 02-06 | Filterable by 5 dimensions | ✓ SATISFIED | FilterBar: status/priority/project/dueDate/adoLink |
| DASH-03 | 02-06 | Linked items show connection | ✓ SATISFIED | TaskRow ADO badge, priority badge, project tag |
| TL-01 | 02-06 | Global task list across projects | ✓ SATISFIED | TasksView shows all tasks, project filter optional |
| TL-02 | 02-06 | Medium-density rows with indicators | ✓ SATISFIED | TaskRow: checkbox + title + priority + project + ADO badge + due date |
| TL-03 | 02-06 | Sort + group-by toggle | ✓ SATISFIED | sortBy default priority, groupBy status/priority/project |
| TL-04 | 02-06 | Click → slide-out detail | ✓ SATISFIED | TaskDetail sidebar opens on task selection |
| TL-05 | 02-06 | Tasks without project (orphan) | ✓ SATISFIED | No project requirement enforced, filterProject optional |
| PROJ-01 | 02-10 | Card grid layout | ✓ SATISFIED | ProjectsView grid cols-1/2/3, ProjectCard component |
| PROJ-02 | 02-10 | Pin/star favorites | ✓ SATISFIED | PinProject backend, pinnedProjects/unpinnedProjects in store |
| PROJ-03 | 02-10 | Flat projects (no sub-projects) | ✓ SATISFIED | No project nesting in schema or UI |
| PROJ-04 | 02-04, 02-10 | Projects linked to ADO scenarios | ✓ SATISFIED | LinkProjectToADO, project_ado_links table |
| PROJ-05 | 02-10 | Project dashboard on click | ✓ SATISFIED | ProjectDetailView with stats, ADO context, filtered tasks |
| PROJ-06 | 02-10 | Dual progress (local + ADO) | ✓ SATISFIED | ProjectCard/ProjectDetailView dual progress bars L173 |
| PROJ-07 | 02-04, 02-10 | Link/unlink project to ADO | ✓ SATISFIED | LinkProjectToADO/UnlinkProject + UI in ProjectDetailView |
| UX-01 | 02-07 | In-app ADO detail + Open in ADO | ✓ SATISFIED | AdoView detail panel L468, "Open in ADO" button L494 |
| UX-02 | 02-09, 02-10 | External links open real browser | ✓ SATISFIED | ExternalLinks.vue window.open(url, '_blank') |
| UX-03 | 02-06, 02-08 | Compact toolbar with all controls | ✓ SATISFIED | FilterBar single row with chips + sync button |
| UX-04 | 02-07 | ADO tree view (Scenario→Deliverable→Task) | ✓ SATISFIED | AdoTreeBrowser recursive tree, expand/collapse |
| UX-05 | 02-07 | Filter chips + search + saved queries | ✓ SATISFIED | AdoView filterType/filterState/filterArea/searchQuery + saved query picker |

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| `frontend/src/components/TaskDetail.vue` | 219 | TODO: Wire to real comment backend | ℹ️ Info | Mock comments in Discussion tab; real CommentsSection properly wired at L502 |
| `frontend/src/components/TaskDetail.vue` | 245 | TODO: Wire to real activity backend | ℹ️ Info | Mock activity timeline — activity feed is a future feature |
| `build/ios/` | - | Missing main function | ℹ️ Info | Pre-existing issue (initial commits), unrelated to Phase 2, causes `go build ./...` to fail |

No blocker or warning anti-patterns found. No placeholder/stub implementations. All `return nil` statements are proper error-handling returns in Go functions.

### Human Verification Required

### 1. Login Flow — All Three Auth Methods
**Test:** Try signing in with Microsoft OAuth, Az CLI token, and PAT
**Expected:** Each method authenticates and routes to /tasks; Az CLI shows "Requires az login" hint
**Why human:** Requires real Microsoft/ADO credentials and az CLI session

### 2. ADO Work Item Tree Browser
**Test:** Navigate to ADO browser, verify tree hierarchy loads from configured orgs
**Expected:** Tree shows Scenario → Deliverable → Task/Bug/Story hierarchy; filters and search work
**Why human:** Requires live ADO connection with real work items

### 3. Link/Promote/Import End-to-End
**Test:** Link a local task to ADO item, promote a task to new ADO item, import an ADO item
**Expected:** task_ado_links entries created, ADO work items created/fetched, isPublic badge updates
**Why human:** Requires ADO write access

### 4. Sync Push Confirmation Dialog
**Test:** Modify a linked task locally, click "Push to ADO", verify preview diff appears
**Expected:** SyncConfirmDialog shows field-by-field changes with local/ADO comparison before confirming
**Why human:** Requires linked task with pending changes

### 5. Background Sync + Conflict Detection
**Test:** Modify same item in both app and ADO, trigger manual sync
**Expected:** Conflict detected; ConflictResolver dialog should appear (CURRENTLY BROKEN — see gap)
**Why human:** Requires simultaneous modification of linked items

### 6. External Links Type Detection
**Test:** Paste ICM/Grafana/ADO/Wiki URLs as links on a task
**Expected:** Each auto-detects correct type with appropriate icon
**Why human:** Visual verification of icon rendering

### 7. Comment Push to ADO
**Test:** Add private comment, click "Push to ADO" button
**Expected:** Comment pushed, marked as Public with globe icon, ADO comment created
**Why human:** Requires ADO write access

### Gaps Summary

**1 gap found — ConflictResolver component not mounted (SYNC-03 partial):**

The `ConflictResolver.vue` component is fully implemented with per-field local/remote selection UI, but it is never imported or rendered in any parent component. The `SyncConfirmDialog.vue` IS correctly mounted in `TaskDetail.vue` at line 740, but the companion `ConflictResolver.vue` was not given the same treatment.

**Root cause:** The ConflictResolver was built (Plan 08) but never wired into the component tree. The sync store correctly manages `showConflictResolver` state and sets it to `true` when conflicts are detected, but without the component being mounted, the Dialog never renders.

**Fix (minimal):** Add `<ConflictResolver />` to `TaskDetail.vue` or `App.vue` alongside the existing `<SyncConfirmDialog />`. This is a ~3 line fix (import + template usage).

**Impact:** Per-field conflict resolution (SYNC-03) will not work at runtime despite full backend support and a complete UI component. This is the only gap preventing SYNC-03 from being fully satisfied.

---

_Verified: 2026-04-06T19:33:05Z_
_Verifier: Claude (gsd-verifier)_
