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
      const { List } = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/taskservice')
      tasks.value = (await List(status)) as Task[]
    } catch (e) {
      console.warn('[TaskStore] Wails binding unavailable:', e)
      tasks.value = []
    } finally {
      loading.value = false
    }
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
  }
})
