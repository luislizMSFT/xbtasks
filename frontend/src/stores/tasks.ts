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

const priorityOrder: Record<string, number> = { P0: 0, P1: 1, P2: 2, P3: 3 }

function startOfDay(d: Date): Date {
  return new Date(d.getFullYear(), d.getMonth(), d.getDate())
}

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const selectedTaskId = ref<number | null>(null)
  const filterStatus = ref<string>('all')

  // --- Public task IDs (personal→public model) ---
  const publicTaskIds = ref<Set<number>>(new Set())

  async function fetchPublicTaskIds() {
    try {
      const { ListPublicTaskIDs } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/linkservice')
      const ids = await ListPublicTaskIDs() as number[]
      publicTaskIds.value = new Set(ids)
    } catch {
      publicTaskIds.value = new Set()
    }
  }

  function isPublic(taskId: number): boolean {
    return publicTaskIds.value.has(taskId)
  }

  // --- Enhanced filter state ---
  const filterPriority = ref<string>('all')
  const filterProject = ref<number | null>(null)
  const filterDueDate = ref<string>('all')    // 'all', 'overdue', 'today', 'week', 'none'
  const filterAdoLink = ref<string>('all')    // 'all', 'linked', 'personal'
  const sortBy = ref<string>('priority')      // 'priority', 'dueDate', 'title', 'status'
  const groupBy = ref<string | null>(null)    // null, 'status', 'priority', 'project'

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

  // --- Enhanced filtered + sorted tasks ---
  const enhancedFilteredTasks = computed(() => {
    // Exclude subtasks — they render nested under their parent
    let result = tasks.value.filter(t => !t.parentId)

    // Status filter
    if (filterStatus.value !== 'all') {
      if (filterStatus.value === 'active') {
        result = result.filter(t => ['todo', 'in_progress', 'in_review'].includes(t.status))
      } else {
        result = result.filter(t => t.status === filterStatus.value)
      }
    }

    // Priority filter
    if (filterPriority.value !== 'all') {
      result = result.filter(t => t.priority === filterPriority.value)
    }

    // Project filter
    if (filterProject.value !== null) {
      result = result.filter(t => t.projectId === filterProject.value)
    }

    // Due date filter
    if (filterDueDate.value !== 'all') {
      const now = new Date()
      const today = startOfDay(now)
      const weekEnd = new Date(today)
      weekEnd.setDate(weekEnd.getDate() + 7)

      if (filterDueDate.value === 'overdue') {
        result = result.filter(t => t.dueDate && new Date(t.dueDate) < today)
      } else if (filterDueDate.value === 'today') {
        const tomorrow = new Date(today)
        tomorrow.setDate(tomorrow.getDate() + 1)
        result = result.filter(t => {
          if (!t.dueDate) return false
          const d = new Date(t.dueDate)
          return d >= today && d < tomorrow
        })
      } else if (filterDueDate.value === 'week') {
        result = result.filter(t => {
          if (!t.dueDate) return false
          const d = new Date(t.dueDate)
          return d >= today && d < weekEnd
        })
      } else if (filterDueDate.value === 'none') {
        result = result.filter(t => !t.dueDate)
      }
    }

    // ADO link filter
    if (filterAdoLink.value === 'linked') {
      result = result.filter(t => publicTaskIds.value.has(t.id))
    } else if (filterAdoLink.value === 'personal') {
      result = result.filter(t => !publicTaskIds.value.has(t.id))
    }

    // Sort — default: priority then due date
    result = [...result].sort((a, b) => {
      if (sortBy.value === 'priority') {
        const pa = priorityOrder[a.priority] ?? 2
        const pb = priorityOrder[b.priority] ?? 2
        if (pa !== pb) return pa - pb
        // Secondary sort by due date
        if (a.dueDate && b.dueDate) return new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime()
        if (a.dueDate) return -1
        if (b.dueDate) return 1
        return 0
      }
      if (sortBy.value === 'dueDate') {
        if (a.dueDate && b.dueDate) return new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime()
        if (a.dueDate) return -1
        if (b.dueDate) return 1
        return 0
      }
      if (sortBy.value === 'title') {
        return a.title.localeCompare(b.title)
      }
      if (sortBy.value === 'status') {
        const statusOrder: Record<string, number> = { in_progress: 0, in_review: 1, todo: 2, blocked: 3, done: 4, cancelled: 5 }
        return (statusOrder[a.status] ?? 3) - (statusOrder[b.status] ?? 3)
      }
      return 0
    })

    return result
  })

  // --- Grouped tasks computed (by groupBy dimension) ---
  const groupedEnhanced = computed(() => {
    const list = enhancedFilteredTasks.value
    if (!groupBy.value) return {} as Record<string, Task[]>

    const groups: Record<string, Task[]> = {}
    for (const t of list) {
      let key: string
      if (groupBy.value === 'status') {
        key = t.status
      } else if (groupBy.value === 'priority') {
        key = t.priority || 'None'
      } else if (groupBy.value === 'project') {
        key = t.projectId ? String(t.projectId) : 'No Project'
      } else {
        key = 'Other'
      }
      if (!groups[key]) groups[key] = []
      groups[key].push(t)
    }
    return groups
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
      const { List } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
      tasks.value = (await List(status)) as Task[]
      // Also fetch public task IDs alongside
      await fetchPublicTaskIds()
    } catch (e) {
      console.warn('[TaskStore] Wails binding unavailable:', e)
      tasks.value = []
    } finally {
      loading.value = false
    }
  }

  async function quickAdd(title: string) {
    return createTask(title, 'P2')
  }

  async function createTask(title: string, priority = 'P2') {
    const { Create } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    const t = await Create(title, '', priority, '', null, null) as Task
    tasks.value.unshift(t)
    return t
  }

  async function updateTask(task: Task) {
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
    const { SetStatus } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    const updated = await SetStatus(id, status) as Task
    const idx = tasks.value.findIndex(t => t.id === id)
    if (idx !== -1) tasks.value[idx] = updated
    return updated
  }

  async function deleteTask(id: number) {
    const { Delete } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    await Delete(id)
    tasks.value = tasks.value.filter(t => t.id !== id)
    if (selectedTaskId.value === id) selectedTaskId.value = null
  }

  async function setPersonalPriority(id: number, priority: string) {
    const { SetPersonalPriority } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    const updated = await SetPersonalPriority(id, priority) as Task
    const idx = tasks.value.findIndex(t => t.id === id)
    if (idx !== -1) tasks.value[idx] = updated
    if (selectedTaskId.value === id) { /* reactive via computed */ }
    return updated
  }

  async function getSubtasks(parentID: number): Promise<Task[]> {
    const { GetSubtasks } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    return (await GetSubtasks(parentID)) as Task[]
  }

  // DependencyService — task-to-task dependency management
  async function getDependencies(taskID: number): Promise<Task[]> {
    const { GetDependencies } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/dependencyservice')
    return (await GetDependencies(taskID)) as Task[]
  }

  async function addDependency(taskID: number, dependsOn: number) {
    const { AddDependency } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/dependencyservice')
    await AddDependency(taskID, dependsOn)
  }

  async function removeDependency(taskID: number, dependsOn: number) {
    const { RemoveDependency } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/dependencyservice')
    await RemoveDependency(taskID, dependsOn)
  }

  async function getAllTags(): Promise<string[]> {
    const { GetAllTags } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
    return (await GetAllTags()) as string[]
  }

  function selectTask(id: number | null) {
    selectedTaskId.value = id
  }

  return {
    tasks, loading, selectedTaskId, filterStatus,
    selectedTask, filteredTasks, grouped, stats,
    fetchTasks, createTask, updateTask, setStatus, deleteTask, selectTask,
    setPersonalPriority, getSubtasks,
    getDependencies, addDependency, removeDependency,
    getAllTags,
    // Phase 02 additions: personal→public model + enhanced filtering
    publicTaskIds, filterPriority, filterProject, filterDueDate, filterAdoLink,
    sortBy, groupBy, enhancedFilteredTasks, groupedEnhanced,
    fetchPublicTaskIds, isPublic, quickAdd,
  }
})
