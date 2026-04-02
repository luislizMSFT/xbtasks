import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

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
  createdAt: string
  updatedAt: string
  completedAt: string | null
}

const MOCK_TASKS: Task[] = [
  {
    id: 1, title: 'Fix auth redirect loop', status: 'in_progress', priority: 'P0',
    category: 'auth', tags: 'bug,urgent', adoId: 'ADO-48291', parentId: null,
    personalPriority: '', description: 'The OAuth callback is redirecting in a loop when the token expires mid-session.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 2 * 3600000).toISOString(),
    updatedAt: new Date(Date.now() - 1800000).toISOString(), completedAt: null,
  },
  {
    id: 2, title: 'Build sidebar navigation', status: 'in_progress', priority: 'P1',
    category: 'frontend', tags: 'ui', adoId: '', parentId: null,
    personalPriority: '', description: 'Implement the main app shell sidebar with icon navigation.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 3600000).toISOString(), completedAt: null,
  },
  {
    id: 3, title: 'Review PR #234 — schema changes', status: 'in_review', priority: 'P1',
    category: '', tags: 'review', adoId: 'ADO-48350', parentId: null,
    personalPriority: '', description: 'Schema migration adding task_deps and pull_requests tables.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 172800000).toISOString(),
    updatedAt: new Date(Date.now() - 7200000).toISOString(), completedAt: null,
  },
  {
    id: 4, title: 'Write unit tests for DependencyService', status: 'todo', priority: 'P2',
    category: 'testing', tags: 'backend,testing', adoId: '', parentId: null,
    personalPriority: '', description: 'Cover Create, Delete, GetBlockers, and GetDependents methods.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 259200000).toISOString(),
    updatedAt: new Date(Date.now() - 259200000).toISOString(), completedAt: null,
  },
  {
    id: 5, title: 'Update API documentation', status: 'todo', priority: 'P3',
    category: 'docs', tags: 'docs', adoId: '', parentId: null,
    personalPriority: '', description: 'Document all TaskService and ProjectService endpoints in REFERENCE.md.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 345600000).toISOString(),
    updatedAt: new Date(Date.now() - 345600000).toISOString(), completedAt: null,
  },
  {
    id: 6, title: 'Waiting on schema review from team', status: 'blocked', priority: 'P1',
    category: 'backend', tags: 'blocked', adoId: 'ADO-47100', parentId: null,
    personalPriority: '', description: 'Cannot proceed with migration until the schema PR is approved.',
    projectId: null, area: '', dueDate: '', blockedReason: 'Needs team review on PR #234', blockedBy: '3',
    createdAt: new Date(Date.now() - 432000000).toISOString(),
    updatedAt: new Date(Date.now() - 86400000).toISOString(), completedAt: null,
  },
  {
    id: 7, title: 'Migrate user preferences to SQLite', status: 'done', priority: 'P2',
    category: 'backend', tags: 'migration', adoId: '', parentId: null,
    personalPriority: '', description: 'Moved user prefs from JSON file to SQLite users table.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 604800000).toISOString(),
    updatedAt: new Date(Date.now() - 518400000).toISOString(),
    completedAt: new Date(Date.now() - 518400000).toISOString(),
  },
  {
    id: 8, title: 'Dark mode color audit', status: 'done', priority: 'P2',
    category: 'design', tags: 'ui,theme', adoId: '', parentId: null,
    personalPriority: '', description: 'Audited all color tokens for WCAG AA contrast in dark mode.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 691200000).toISOString(),
    updatedAt: new Date(Date.now() - 604800000).toISOString(),
    completedAt: new Date(Date.now() - 604800000).toISOString(),
  },
  {
    id: 9, title: 'Set up CI pipeline', status: 'cancelled', priority: 'P3',
    category: 'infra', tags: 'ci', adoId: 'ADO-46500', parentId: null,
    personalPriority: '', description: 'Was going to set up GitHub Actions but decided to use ADO Pipelines instead.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 864000000).toISOString(),
    updatedAt: new Date(Date.now() - 777600000).toISOString(), completedAt: null,
  },
  {
    id: 10, title: 'Add keyboard shortcuts documentation', status: 'todo', priority: 'P3',
    category: 'docs', tags: 'docs,ux', adoId: '', parentId: null,
    personalPriority: '', description: 'Document all keyboard shortcuts available in the app.',
    projectId: null, area: '', dueDate: '', blockedReason: '', blockedBy: '',
    createdAt: new Date(Date.now() - 172800000).toISOString(),
    updatedAt: new Date(Date.now() - 172800000).toISOString(), completedAt: null,
  },
]

let useMock = true

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const selectedTaskId = ref<number | null>(null)
  const filterStatus = ref<string>('all')

  const selectedTask = computed(() =>
    tasks.value.find(t => t.id === selectedTaskId.value) ?? null
  )

  const filteredTasks = computed(() => {
    if (filterStatus.value === 'all') return tasks.value
    if (filterStatus.value === 'active')
      return tasks.value.filter(t => ['todo', 'in_progress', 'in_review'].includes(t.status))
    if (filterStatus.value === 'done')
      return tasks.value.filter(t => t.status === 'done')
    if (filterStatus.value === 'blocked')
      return tasks.value.filter(t => t.status === 'blocked')
    return tasks.value.filter(t => t.status === filterStatus.value)
  })

  const grouped = computed(() => {
    const groups: Record<string, Task[]> = {
      in_progress: [],
      in_review: [],
      todo: [],
      blocked: [],
      done: [],
      cancelled: [],
    }
    for (const t of filteredTasks.value) {
      if (groups[t.status]) groups[t.status].push(t)
    }
    return groups
  })

  const stats = computed(() => ({
    total: tasks.value.length,
    inProgress: tasks.value.filter(t => t.status === 'in_progress').length,
    inReview: tasks.value.filter(t => t.status === 'in_review').length,
    blocked: tasks.value.filter(t => t.status === 'blocked').length,
    done: tasks.value.filter(t => t.status === 'done').length,
    todo: tasks.value.filter(t => t.status === 'todo').length,
  }))

  async function fetchTasks(status = '') {
    loading.value = true
    try {
      if (!useMock) {
        const { List } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
        tasks.value = (await List(status)) as Task[]
      } else {
        await new Promise(r => setTimeout(r, 200)) // simulate network
        tasks.value = [...MOCK_TASKS]
      }
    } catch {
      useMock = true
      tasks.value = [...MOCK_TASKS]
    } finally {
      loading.value = false
    }
  }

  async function createTask(title: string, priority = 'P2') {
    if (useMock) {
      const newTask: Task = {
        id: Math.max(0, ...tasks.value.map(t => t.id)) + 1,
        title, description: '', status: 'todo', priority,
        category: '', projectId: null, area: '', dueDate: '',
        adoId: '', tags: '', blockedReason: '', blockedBy: '',
        parentId: null, personalPriority: '',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(), completedAt: null,
      }
      tasks.value.unshift(newTask)
      return newTask
    }
    const { Create } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    const t = await Create(title, '', priority, '', null, null) as Task
    tasks.value.unshift(t)
    return t
  }

  async function updateTask(task: Task) {
    if (useMock) {
      const idx = tasks.value.findIndex(t => t.id === task.id)
      if (idx !== -1) {
        task.updatedAt = new Date().toISOString()
        if (task.status === 'done' && !task.completedAt) {
          task.completedAt = new Date().toISOString()
        }
        tasks.value[idx] = { ...task }
      }
      return task
    }
    const { Update } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    const updated = await Update(
      task.id, task.title, task.description, task.status,
      task.priority, task.category, task.area, task.dueDate, task.tags
    ) as Task
    const idx = tasks.value.findIndex(t => t.id === task.id)
    if (idx !== -1) tasks.value[idx] = updated
    return updated
  }

  async function setStatus(id: number, status: string) {
    if (useMock) {
      const t = tasks.value.find(t => t.id === id)
      if (t) {
        t.status = status
        t.updatedAt = new Date().toISOString()
        if (status === 'done') t.completedAt = new Date().toISOString()
      }
      return t
    }
    const { SetStatus } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    const updated = await SetStatus(id, status) as Task
    const idx = tasks.value.findIndex(t => t.id === id)
    if (idx !== -1) tasks.value[idx] = updated
    return updated
  }

  async function deleteTask(id: number) {
    if (useMock) {
      tasks.value = tasks.value.filter(t => t.id !== id)
      if (selectedTaskId.value === id) selectedTaskId.value = null
      return
    }
    const { Delete } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    await Delete(id)
    tasks.value = tasks.value.filter(t => t.id !== id)
    if (selectedTaskId.value === id) selectedTaskId.value = null
  }

  function selectTask(id: number | null) {
    selectedTaskId.value = id
  }

  return {
    tasks, loading, selectedTaskId, filterStatus,
    selectedTask, filteredTasks, grouped, stats,
    fetchTasks, createTask, updateTask, setStatus, deleteTask, selectTask,
  }
})
