# XB-TASKS BACKEND INTEGRATION ANALYSIS

## 1. BACKEND SERVICE INVENTORY

### Services Registered with Wails (main.go lines 64-76)
All services are registered as Wails bindings, exposing all exported methods to the frontend via IPC.

#### Service Breakdown:

**ConfigService** - Configuration management
- Not detailed in provided files but handles app settings, org/project config, sync intervals

**TaskService** (tasks.go) - Core task CRUD operations
- Create(title, description, priority, category, projectID, parentID)  Task
- GetByID(id)  Task
- List(status)  []Task
- Update(id, title, description, status, priority, category, area, dueDate, tags)  Task
- Delete(id)  error
- SetStatus(id, status)  Task
- GetSubtasks(parentID)  []Task
- SetPersonalPriority(id, personalPriority)  Task
- CreateSubtask(parentID, title, description, priority)  Task
- ListFiltered(status, projectID, parentID, tag)  []Task
- GetADOLinks(taskID)  []TaskADOLink
- GetAllTags()  []string
- ReorderTasks(orderedIDs)  error

**ProjectService** (projects.go) - Project management
- Create(name, description)  Project
- GetByID(id)  Project
- List()  []Project
- Update(id, name, description, status)  Project
- Delete(id)  error
- LinkProjectToADO(projectID, adoID, direction)  error
- UnlinkProject(projectID, adoID, deleteLocal)  error
- GetProjectADOLink(projectID)  ProjectProviderLink
- PinProject(projectID, pinned)  error
- GetProjectProgress(projectID)  map[string]any

**DependencyService** (deps.go) - Task dependency tracking
- AddDependency(taskID, dependsOn)  error (with cycle detection)
- RemoveDependency(taskID, dependsOn)  error
- GetDependencies(taskID)  []Task
- GetBlockedBy(taskID)  []Task

**AuthService** - User authentication (not detailed but registered)

**ADOService** (adoservice.go) - Azure DevOps integration (read-heavy)
- CheckConnection()  string
- ListMyWorkItems()  []WorkItem (20s timeout, parallel org queries)
- GetWorkItem(adoID)  WorkItem
- SyncWorkItems()  error (fetches + upserts to cache)
- GetCachedWorkItems()  []WorkItem
- GetCachedWorkItem(adoID)  WorkItem
- GetSavedQueries()  []SavedQuery
- RunSavedQuery(queryID)  []WorkItem
- GetWorkItemTree()  []WorkItem (3-level parent hierarchy fetch)

**LinkService** (linkservice.go) - Task-ADO linking lifecycle
- LinkTask(taskID, adoID)  TaskADOLink
- PromoteTask(taskID, wiType)  TaskADOLink (creates new ADO work item)
- ImportWorkItem(adoID)  Task
- ImportWorkItemAsProject(adoID)  Project (imports parent + children)
- UnlinkTask(taskID, adoID, deleteLocal)  error
- IsPublic(taskID)  bool
- GetTaskLinks(taskID)  []TaskADOLink
- ListPublicTaskIDs()  []int
- ListLinkedAdoIDs()  []string

**PRService** (prservice.go) - Pull request queries (via az cli)
- ListMyPRs()  []PullRequest (top 50)
- ListReviewPRs()  []PullRequest (top 50, active)
- ListTeamPRs()  []PullRequest (top 100, active)
- SyncPRs()  error (deduplicates + upserts)
- GetCachedPRs()  []PullRequest

**PipelineService** (pipelineservice.go) - CI/CD pipeline runs (via az cli)
- ListRecentRuns()  []Pipeline (top 20)
- GetPipelineRun(runID)  Pipeline

**SyncService** (syncservice.go) - Bidirectional sync with ADO
- StartBackgroundSync() (auto-pulls on 15-min interval)
- StopSync()
- ManualSync()  []SyncDiff
- GenerateOutboundDiff(taskID)  SyncDiff
- PushChanges(taskID)  error
- ResolveConflict(taskID, resolutions)  error
- DetectConflicts(taskID)  []FieldDiff

**CommentService** (commentservice.go) - Task comments + ADO push
- AddComment(taskID, content)  TaskComment
- ListComments(taskID)  []TaskComment
- UpdateComment(id, content)  error
- DeleteComment(id)  error
- PushCommentToADO(commentID)  error
- FetchADOComments(taskID)  []ADOComment
- ReplyToADOComment(taskID, content)  *ADOComment

**ExternalLinksService** (externallinks.go) - URL link attachment
- AddLink(taskID, url, label)  TaskLink (auto-detects type: icm, grafana, ado, wiki, url)
- ListLinks(taskID)  []TaskLink
- DeleteLink(id)  error

**BrowserService** (browserservice.go) - System browser integration
- OpenURL(url)  error

---

## 2. DATA MODEL MAPPING

### Core Domain Types

**Task** (domain/task.go)
- id, title, description, status, priority, category
- projectId, area, dueDate
- adoId (single link; multi-link via task_ado_links table)
- tags (comma-separated string)
- blockedReason, blockedBy
- **parentId**  KEY: Enables task hierarchy/subtasks
- personalPriority (override)
- sortOrder (manual reordering)
- createdAt, updatedAt, completedAt

Status enum: todo, in_progress, in_review, done, blocked, cancelled
Priority enum: P0, P1, P2, P3

**WorkItem** (domain/workitem.go) - ADO representation
- adoId (string), title, state, type
- assignedTo, priority (int), areaPath, description, url
- org, project (denormalized)
- parentId (int)  Enables ADO tree fetch
- changedDate, syncedAt

**SyncState** (domain/sync.go) - Conflict detection baseline
- taskId, adoId
- lastSyncedAt
- localTitle, localStatus, localDesc (last known local)
- remoteTitle, remoteStatus, remoteDesc (last known remote)

**SyncDiff** (domain/sync.go) - Change summary for display
- taskId, adoId
- changes: []FieldDiff
- direction: "inbound", "outbound", "conflict"

**FieldDiff** - Per-field conflict view
- field, local, remote, proposed

### Mock Data to Real API Mapping

Mock playground structure  Real backend structure:

`
{
  id,
  title,
  status,
  priority,
  parentId,            Task.ParentID
  
  // ADO metadata (from WorkItem cache)
  ado: {
    type,              WorkItem.Type
    state,             WorkItem.State
    adoId,             WorkItem.AdoId (via task_ado_links)
  },
  
  // Subtasks
  subtasks: [...]      TaskService.GetSubtasks(parentID)
  
  // Links
  adoLinks: [...]      TaskService.GetADOLinks(taskID)
  
  // Sync status
  syncStatus,          SyncState + SyncDiff data
  
  // Comments (local + ADO)
  comments: [...]      CommentService.ListComments + FetchADOComments
  
  // PR links
  prs: [...]           PRService.GetCachedPRs (filtered by task context)
  
  // Notes
  notes: [...],        External links + Comments combined
  
  dirty: false         Local UI state (not persisted)
}
`

---

## 3. TASK TREE LOADING

### Current Implementation (ADO-style hierarchical)

**Local Task Hierarchy:**
- Uses parentId field in Task struct (nullable *int)
- Frontend reconstructs tree via:
  1. Load all root tasks (status filter, no parent)
  2. For each task, call TaskService.GetSubtasks(parentID)  **N+1 problem here!**

**ADO Work Item Tree:**
- ADOService.GetWorkItemTree() (lines 222-277):
  1. Fetch assigned items: ListMyWorkItems()
  2. Collect parent IDs not already in list
  3. Batch-fetch parents in 3 levels (multipart)
  4. Return flat list with ParentID relationships
  
  - Uses goroutines for parallel org/project queries
  - Tries all configured clients to find items

**Key Issue:** Frontend must query subtasks one level at a time
- Loading 30+ tasks = 30+ separate GetSubtasks() calls
- No single API returns a complete tree structure

---

## 4. ADO METADATA FLOW

### Data Path: ADO  Frontend

1. **Token acquisition** (via az cli or cached PAT):
   - uth.TokenProvider.GetToken() (5-min cache via CachedTokenProvider)
   
2. **ADO queries** (via REST client):
   - ADOService.ListMyWorkItems()  parallel queries to all configured orgs (20s timeout per call)
   - Results cached in SQLite via ADOService.SyncWorkItems()
   
3. **Cache layer** (SQLite ado_work_items table):
   - ADOService.GetCachedWorkItem(adoID)  Frontend reads from cache
   - Populated by sync operations, not on-demand
   
4. **Type/State mapping**:
   - ADO state  local status: do.MapADOToStatus(remote.State)
   - Local status  ADO state: do.MapStatusToADO(task.Status, remote.Type)
   - Type affects valid state transitions

5. **Frontend consumption**:
   - TaskService.GetADOLinks(taskID)  Links task to ADO
   - ADOService.GetCachedWorkItem(adoId)  Gets metadata (type, state, etc.)
   - Separate round-trip per task with ADO link

**Flow Diagram:**
`
Frontend     ADOService      ado.Client      SQLite Cache      ADO REST
   |            |                |                 |               |
   |--ListMyWorkItems()--------->|                 |               |
   |            |--Parallel queries (20s timeout)--|--GET /workitems|
   |            |                |                 |<--JSON---------|
   |            |                |                 |
   |            |--SyncWorkItems: Upsert----------->|
   |            |                |              (ado_work_items)
   |<--GetCachedWorkItem---------|                 |
   |            |                |              (read local copy)
`

### Performance Concerns:
- No eager loading of ADO metadata with task list
- Each task fetch requires separate ADO lookup if not cached
- 20-second timeout can block UI on slow networks
- Cache can become stale (depends on sync frequency)

---

## 5. PERFORMANCE CONCERNS

### Critical Issues:

#### 1. **N+1 Problem: Subtask Loading**
- TaskService.GetSubtasks() called once per parent task
- Loading 30 tasks  up to 30 DB queries
- **Fix needed:** Batch query in one call (GROUP BY parent_id)

#### 2. **N+1 Problem: ADO Metadata**
- Per-task ADO lookup if data not in cache
- TaskService.GetADOLinks(taskID) then ADOService.GetCachedWorkItem(adoId)
- 30 tasks = 30 link queries + 30 metadata queries
- **Fix needed:** Single join query with ado_work_items table

#### 3. **Large Payload Risk: GetWorkItemTree**
- Recursively fetches parent hierarchy (3 levels deep)
- Each level can spawn batch queries across all orgs
- **Impact:** 100+ API calls to ADO on cold sync
- **Fix:** Cap fetching to 2 levels or limit to immediate parents

#### 4. **Sync Service: Full Table Scan**
- pullChanges() queries ALL sync_states (line 143)
- For each state, fetches remote item individually
- 30 linked items = 30 ADO API calls (serialized per org)
- **Fix:** Use batch API for multiple items

#### 5. **Real-time Update Pattern**
- Background sync every 15 minutes (configurable)
- Polling-based, no event stream from ADO
- Large workloads will have stale data between cycles
- **Missing:** WebSocket or real-time event support

#### 6. **Comment Fetching**
- CommentService.FetchADOComments(taskID) queries ADO on-demand
- If task has 10 comments in ADO, all fetched into memory
- Not cached; every view = new API call

### Query Pattern Analysis:

**Eager Load Missing:**
- Tasks loaded without comments, links, or ADO metadata
- Frontend must issue additional requests

**Batch Operations:**
- ReorderTasks(orderedIDs) uses transaction 
- SyncPRs() deduplicates in-memory 
- But sync state pull is not batched 

**Index Recommendations:**
- 	ask_ado_links(task_id) for taskADO lookups
- do_work_items(parent_id) for tree traversal
- sync_state(task_id, ado_id) composite

---

## 6. MISSING ENDPOINTS / OPERATIONS

Based on playground requirements vs. backend capability:

### Exists ():
-  Task CRUD (Create, GetByID, Update, Delete)
-  Task status transitions (SetStatus)
-  Subtask creation (CreateSubtask)
-  Task linking to ADO (LinkTask, PromoteTask, ImportWorkItem)
-  Sync pull (ManualSync)
-  Sync push with diff (GenerateOutboundDiff, PushChanges)
-  Conflict resolution (ResolveConflict)
-  Comments local + push (AddComment, PushCommentToADO)
-  External links (AddLink)
-  Task reordering (ReorderTasks)
-  Tag management (GetAllTags)

### Missing ():

**Quick-add Task:**
-  TaskService.Create() exists
- But no "quick" variant that returns minimal data
- **Issue:** Full task fetch after create may be overhead

**Inline Subtask Creation:**
-  TaskService.CreateSubtask() exists
- Works as inline operation
- **OK**

**Bulk Task Operations:**
-  No bulk delete
-  No bulk status change
-  No bulk tag update
- **Impact:** Each operation = separate API call

**PR Linking:**
-  No endpoint to link PR to task
-  PRService only lists cached PRs
- **Fix needed:** PRService.LinkPRToTask(taskID, prNumber) + table for task_pr_links

**Note/Knowledge Base:**
-  No persistent note storage (comments work, but are time-based)
- Comments + External links partially cover this
- **Limitation:** No structured "note" concept distinct from comments

**Tree Operations:**
-  No bulk move (reparent multiple tasks)
-  No flatten/expand all
-  No export tree structure

**Search:**
-  No full-text search endpoint
-  No advanced filter (multi-field AND/OR)
- Workaround: Frontend implements client-side search

**Notification/Real-time:**
-  No event subscription
-  No WebSocket support
- Polling only (15-min background sync)

**Batch Sync:**
-  No batch conflict resolution
-  No batch push approval
- Each task = separate sync operation

---

## 7. WAILS BINDING PATTERNS

### Registration (main.go lines 64-76)

`go
wailsApp.RegisterService(application.NewService(taskService))
wailsApp.RegisterService(application.NewService(adoService))
// ... etc for each service
`

**How it works:**
1. Each service is passed to pplication.NewService()
2. Wails reflects on the service struct
3. **All exported methods** (PascalCase) become callable from frontend
4. IPC marshals arguments/results to JSON
5. Errors returned as JSON with message

### Frontend Call Pattern:

**Synchronous RPC (implied by service methods):**
`javascript
// Frontend (Vue)
const task = await TaskService.GetByID(123);
const subtasks = await TaskService.GetSubtasks(123);
`

**Event-based (for async operations):**
- SyncService.emitEvent() (line 37) sends events to frontend
- E.g., "sync:started", "sync:completed", "sync:conflict"
- Backend  Frontend real-time updates (limited by polling)

### Type Safety:
- No generated TypeScript types visible in provided files
- Frontend likely uses manual type definitions or untyped ny
- **Risk:** Type mismatches between Go structs and JS objects

### Authentication Flow:
- Token provider handles auth behind the scenes
- No explicit login/logout method visible
- uth.NewAuthService() likely has those methods (not detailed)

### Background Operations:
- syncService.StartBackgroundSync() (line 122 main.go)
- Runs in separate goroutine
- Emits events on completion/failure

---

## 8. CRITICAL BOTTLENECKS & RECOMMENDATIONS

### High Priority (Performance):

1. **Batch Subtask Fetch**
   - Replace NGetSubtasks with single query
   - Example: TaskService.GetTasksWithSubtasks(parentIDs []int)  map[int][]Task

2. **Eager Load ADO Metadata**
   - Include ADO join in task queries when available
   - TaskService.GetByID() could return WorkItem reference if linked

3. **Cache ADO Comments**
   - Store fetched comments locally with TTL
   - Invalidate on push, not on every view

4. **Batch Sync Pull**
   - Use GetWorkItemsByIDs() to fetch all remote items at once
   - Currently one-by-one in sync loop

### Medium Priority (Usability):

5. **PR Task Linking**
   - Add TaskService.LinkPRToTask() + persist link
   - Enable bi-directional reference

6. **Tree Navigation API**
   - TaskService.GetTaskTree(parentID)  NestedTask structure
   - Avoid multi-round-trip reconstruction

7. **Search API**
   - TaskService.Search(query string)  []Task
   - Enable efficient filtering on backend

### Low Priority (Future):

8. **Real-time Events**
   - WebSocket for instant sync updates
   - Subscription-based ADO change notifications

9. **Bulk Operations**
   - BulkDelete, BulkStatusChange, BulkReparent
   - Transactional safety

10. **Full-text Indexing**
    - SQLite FTS5 for task search
    - Faster than LIKE queries

---

## 9. QUICK REFERENCE: API CHEAT SHEET

### Task Operations
`
Create:       TaskService.Create(title, description, priority, category, projectID, parentID)  Task
Read:         TaskService.GetByID(id)  Task
Update:       TaskService.Update(id, title, ...)  Task
Delete:       TaskService.Delete(id)  error
GetSubtasks:  TaskService.GetSubtasks(parentID)  []Task
CreateSub:    TaskService.CreateSubtask(parentID, title, ...)  Task
Reorder:      TaskService.ReorderTasks(orderedIDs)  error
SetStatus:    TaskService.SetStatus(id, status)  Task
`

### ADO Linking
`
Link:         LinkService.LinkTask(taskID, adoID)  TaskADOLink
Promote:      LinkService.PromoteTask(taskID, wiType)  TaskADOLink
Import:       LinkService.ImportWorkItem(adoID)  Task
Unlink:       LinkService.UnlinkTask(taskID, adoID, deleteLocal)  error
GetLinks:     TaskService.GetADOLinks(taskID)  []TaskADOLink
`

### Sync
`
ManualSync:       SyncService.ManualSync()  []SyncDiff
OutboundDiff:     SyncService.GenerateOutboundDiff(taskID)  SyncDiff
Push:             SyncService.PushChanges(taskID)  error
Resolve:          SyncService.ResolveConflict(taskID, resolutions)  error
Detect:           SyncService.DetectConflicts(taskID)  []FieldDiff
`

### Comments
`
Add:          CommentService.AddComment(taskID, content)  TaskComment
List:         CommentService.ListComments(taskID)  []TaskComment
Update:       CommentService.UpdateComment(id, content)  error
Delete:       CommentService.DeleteComment(id)  error
Push:         CommentService.PushCommentToADO(commentID)  error
FetchADO:     CommentService.FetchADOComments(taskID)  []ADOComment
ReplyADO:     CommentService.ReplyToADOComment(taskID, content)  *ADOComment
`

### Projects
`
Create:       ProjectService.Create(name, description)  Project
GetByID:      ProjectService.GetByID(id)  Project
List:         ProjectService.List()  []Project
Update:       ProjectService.Update(id, name, ...)  Project
Delete:       ProjectService.Delete(id)  error
GetProgress:  ProjectService.GetProjectProgress(projectID)  map[string]any
`

### External Links
`
Add:          ExternalLinksService.AddLink(taskID, url, label)  TaskLink
List:         ExternalLinksService.ListLinks(taskID)  []TaskLink
Delete:       ExternalLinksService.DeleteLink(id)  error
`

### PRs
`
ListMine:     PRService.ListMyPRs()  []PullRequest
ListReview:   PRService.ListReviewPRs()  []PullRequest
ListTeam:     PRService.ListTeamPRs()  []PullRequest
GetCached:    PRService.GetCachedPRs()  []PullRequest
Sync:         PRService.SyncPRs()  error
`

---

## SUMMARY TABLE

| Category | Status | Key Points |
|----------|--------|-----------|
| **Services** | 13 registered | All methods auto-exposed to frontend |
| **Subtask Loading** |  N+1 | Each subtask = 1 DB query |
| **ADO Metadata** |  N+1 | Each ADO item = 1 cache lookup |
| **Sync Pull** |  Slow | 30 items = 30 ADO API calls |
| **Comments** | Not cached | Every view = fresh ADO fetch |
| **PRs** | Cached | Synced on-demand, readable offline |
| **Tree Navigation** | Manual | Frontend must reconstruct |
| **Real-time** | No events | Polling only (15-min cycle) |
| **PR-Task Link** | Missing | No persistent relationship |
| **Bulk Operations** | Limited | One-by-one only |

