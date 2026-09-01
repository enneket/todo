import { defineStore } from 'pinia'
import axios from 'axios'
import { ref, computed } from 'vue'

export interface Subtask {
  // id is `number` for server-side subtasks and `string` (UUID) for
  // temporary subtasks the user added to the form before the parent todo
  // was created.
  id: number | string
  todo_id: number
  title: string
  completed: boolean
  created_at: string
}

export interface Todo {
  id: number
  title: string
  description: string
  completed: boolean
  priority: 'high' | 'medium' | 'low'
  due_date: string | null
  remind_at: string | null
  notified_at: string | null
  repeat: string
  tags: string[]
  project_id: number | null
  subtasks: Subtask[]
  created_at: string
  deleted_at: string | null
}

export const useTodoStore = defineStore('todo', () => {
  const todos = ref<Todo[]>([])
  const trashTodos = ref<Todo[]>([])
  const lastError = ref<string | null>(null)

  const uniqueTags = computed(() => {
    const tags = new Set<string>()
    todos.value.forEach(todo => {
      if (todo.tags && Array.isArray(todo.tags)) {
        todo.tags.forEach(tag => {
          if (tag.trim()) tags.add(tag.trim())
        })
      }
    })
    return Array.from(tags).sort()
  })

  // Assume API is running on localhost:8081 (from main.go)
  const API_URL = 'http://localhost:8081/api/todos'
  const SUBTASK_API_URL = 'http://localhost:8081/api/subtasks'
  const TRASH_API_URL = 'http://localhost:8081/api/trash/todos'

  // Wrap a promise so any rejection is recorded on `lastError` instead of
  // swallowed into the console — the UI can react and show a toast.
  async function tracked<T>(label: string, p: Promise<T>): Promise<T | null> {
    try {
      return await p
    } catch (err: unknown) {
      let message = 'unknown error'
      if (axios.isAxiosError(err)) {
        message = String(err.response?.data || err.message)
      } else if (err instanceof Error) {
        message = err.message
      }
      lastError.value = `${label}: ${message}`
      console.error(lastError.value, err)
      return null
    }
  }

  function clearError() {
    lastError.value = null
  }

  const fetchTodos = async () => {
    const result = await tracked('fetchTodos', axios.get<Todo[]>(API_URL))
    if (result) todos.value = result.data || []
  }

  const addTodo = async (
    title: string,
    priority: string = 'medium',
    dueDate: string | null = null,
    description: string = '',
    tags: string[] = [],
    projectId: number | null = null,
    remindAt: string | null = null,
    repeat: string = ''
  ): Promise<number | null> => {
    const dateToSend = dueDate === '' ? null : dueDate
    const remindToSend = remindAt === '' ? null : remindAt
    const response = await tracked(
      'addTodo',
      axios.post<Todo>(API_URL, {
        title,
        description,
        priority,
        due_date: dateToSend,
        tags,
        project_id: projectId,
        remind_at: remindToSend,
        repeat,
      })
    )
    if (!response) return null
    // Prepend the freshly-created todo so the UI updates immediately without
    // re-fetching the whole list.
    todos.value.unshift(response.data)
    return response.data.id
  }

  const updateTodo = async (id: number, updates: Partial<Todo>) => {
    // The form sends '' for cleared date fields, but the backend expects a
    // `*time.Time` (null clears it). Leave the value untouched when it's
    // absent so a partial update doesn't wipe a field the caller didn't
    // intend to change.
    if (updates.due_date === '') {
      updates.due_date = null
    }
    if (updates.remind_at === '') {
      updates.remind_at = null
    }
    const response = await tracked(
      'updateTodo',
      axios.put<Todo>(`${API_URL}/${id}`, updates)
    )
    if (!response) return
    // Replace the row in place with whatever the server returned — that way
    // server-side defaults (e.g. notified_at reset, repeat reschedule) are
    // reflected immediately without an extra GET.
    const idx = todos.value.findIndex(t => t.id === id)
    if (idx !== -1) {
      todos.value[idx] = response.data
    }
  }

  const deleteTodo = async (id: number) => {
    const ok = await tracked('deleteTodo', axios.delete(`${API_URL}/${id}`))
    if (ok === null) return
    todos.value = todos.value.filter(t => t.id !== id)
  }

  const fetchTrashTodos = async () => {
    const result = await tracked('fetchTrashTodos', axios.get<Todo[]>(TRASH_API_URL))
    if (result) trashTodos.value = result.data || []
  }

  const restoreTrashTodo = async (id: number) => {
    const ok = await tracked('restoreTrashTodo', axios.post(`${TRASH_API_URL}/${id}/restore`))
    if (ok === null) return
    trashTodos.value = trashTodos.value.filter(t => t.id !== id)
  }

  const purgeTrashTodo = async (id: number) => {
    const ok = await tracked('purgeTrashTodo', axios.delete(`${TRASH_API_URL}/${id}`))
    if (ok === null) return
    trashTodos.value = trashTodos.value.filter(t => t.id !== id)
  }

  const addSubtask = async (todoId: number, title: string) => {
    const response = await tracked(
      'addSubtask',
      axios.post<Subtask>(`${API_URL}/${todoId}/subtasks`, { title })
    )
    if (!response) return
    const idx = todos.value.findIndex(t => t.id === todoId)
    if (idx !== -1) {
      todos.value[idx].subtasks = [...todos.value[idx].subtasks, response.data]
    }
  }

  const updateSubtask = async (id: number | string, updates: Partial<Subtask>) => {
    // Skip the server call for in-memory temp subtasks.
    if (typeof id !== 'number') return
    const ok = await tracked('updateSubtask', axios.put(`${SUBTASK_API_URL}/${id}`, updates))
    if (ok === null) return
    // Find which todo owns this subtask and patch in place.
    for (const t of todos.value) {
      const i = t.subtasks.findIndex(s => s.id === id)
      if (i !== -1) {
        t.subtasks[i] = { ...t.subtasks[i], ...updates }
        break
      }
    }
  }

  const deleteSubtask = async (id: number | string) => {
    // Skip the server call for in-memory temp subtasks; the caller has already
    // removed them from the form state.
    if (typeof id !== 'number') return
    const ok = await tracked('deleteSubtask', axios.delete(`${SUBTASK_API_URL}/${id}`))
    if (ok === null) return
    for (const t of todos.value) {
      const i = t.subtasks.findIndex(s => s.id === id)
      if (i !== -1) {
        t.subtasks.splice(i, 1)
        break
      }
    }
  }

  return {
    todos,
    trashTodos,
    uniqueTags,
    lastError,
    clearError,
    fetchTodos,
    addTodo,
    updateTodo,
    deleteTodo,
    addSubtask,
    updateSubtask,
    deleteSubtask,
    fetchTrashTodos,
    restoreTrashTodo,
    purgeTrashTodo,
  }
})
