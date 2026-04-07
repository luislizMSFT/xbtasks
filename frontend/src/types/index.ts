// Shared type definitions for xb-tasks frontend.
// Mirrors Go domain/ entities with camelCase fields (Wails JSON serialization).

// ─── Tasks ──────────────────────────────────────────────────────────────────

export interface Task {
  id: number
  title: string
  description: string
  status: string
  priority: string
  category: string
  projectId: number | null
  area: string
  dueDate: string
  adoId: string
  tags: string
  blockedReason: string
  blockedBy: string
  parentId: number | null
  personalPriority: string
  sortOrder: number
  createdAt: string
  updatedAt: string
  completedAt: string | null
}

export interface TaskComment {
  id: number
  taskId: number
  content: string
  isPublic: boolean
  adoCommentId: string
  createdAt: string
  updatedAt: string
}

export interface TaskLink {
  id: number
  taskId: number
  url: string
  label: string
  type: string // 'url', 'icm', 'grafana', 'ado', 'wiki'
  createdAt: string
}

// ─── Azure DevOps ───────────────────────────────────────────────────────────

export interface WorkItem {
  id: number
  adoId: string
  title: string
  state: string
  type: string
  assignedTo: string
  priority: number
  areaPath: string
  description: string
  url: string
  org: string
  project: string
  parentId: number
  changedDate: string
  syncedAt: string
}

/** @deprecated Use WorkItem instead */
export type ADOWorkItem = WorkItem

export interface Pipeline {
  id: number
  name: string
  status: string
  result: string
  url: string
  sourceBranch: string
  queueTime: string
  finishTime: string | null
}

/** @deprecated Use Pipeline instead */
export type ADOPipeline = Pipeline

export interface SavedQuery {
  id: string
  name: string
  path: string
  isFolder: boolean
}

export interface ADOComment {
  id: number
  text: string
  createdBy: string
  createdDate: string
}

// ─── Pull Requests ──────────────────────────────────────────────────────────

export interface PRReviewer {
  displayName: string
  uniqueName: string
  vote: number // 10=approve, 5=approve-with-suggestions, -5=wait, -10=reject, 0=none
}

export interface PullRequest {
  id: number
  title: string
  prUrl: string
  prNumber: number
  repo: string
  taskId: number | null
  adoId: string
  status: string // draft, active, completed, abandoned
  reviewers: string // JSON array
  sourceBranch: string
  targetBranch: string
  votes: number
  createdBy: string
  createdAt: string
  updatedAt: string
  mergedAt: string | null
}

// ─── Projects ───────────────────────────────────────────────────────────────

export interface Project {
  id: number
  name: string
  description: string
  status: string
  isPinned: boolean
  createdAt: string
  updatedAt: string
  taskCount?: number
}

export interface ProjectADOLink {
  projectId: number
  adoId: string
  direction: string
  createdAt: string
}

export interface ProjectProgress {
  localDone: number
  localTotal: number
  adoDone: number
  adoTotal: number
}

// ─── Sync ───────────────────────────────────────────────────────────────────

export interface FieldDiff {
  field: string    // 'title', 'status', 'description'
  local: string    // current local value
  remote: string   // current ADO value
  proposed: string // what will be written
}

export interface SyncDiff {
  taskId: number
  adoId: string
  changes: FieldDiff[]
  direction: string // 'push', 'pull', or 'conflict'
}

export interface Conflict {
  taskId: number
  adoId: string
  fields: FieldDiff[]
}

// ─── Auth ───────────────────────────────────────────────────────────────────

export interface User {
  id: string
  displayName: string
  email: string
  avatarUrl: string
}

// ─── Graph / Dependency Visualization ───────────────────────────────────────

export interface SubtaskItem {
  id: number
  title: string
  done: boolean
}

export interface GraphNode {
  id: number
  title: string
  status: string
  priority: string
  description: string
  dependsOn: number[]
  children: GraphNode[]
  subtasks?: SubtaskItem[]
}

export interface ForceGraphNode {
  id: number
  title: string
  status: string
  priority: string
  projectId: number | null
  source: 'local' | 'ado'
  adoId?: string
  adoUrl?: string
  type?: string // ADO work item type (e.g., 'Epic', 'Feature', 'User Story')
  x?: number
  y?: number
}

export type ForceGraphEdgeKind = 'dependency' | 'parent'

export interface ForceGraphEdge {
  source: number
  target: number
  kind: ForceGraphEdgeKind
}
