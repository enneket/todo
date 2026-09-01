import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useTodoStore, type Todo } from './todo'
import axios from 'axios'

vi.mock('axios')

describe('Todo Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches todos successfully', async () => {
    const store = useTodoStore()
    const mockTodos = [
      { id: 1, title: 'Test Todo', completed: false, priority: 'medium', due_date: null, tags: [], project_id: null, subtasks: [], created_at: '' }
    ]
    
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: mockTodos })

    await store.fetchTodos()

    expect(axios.get).toHaveBeenCalledWith('http://localhost:8081/api/todos')
    expect(store.todos).toEqual(mockTodos)
  })

  it('adds a todo successfully', async () => {
    const store = useTodoStore()
    
    // @ts-expect-error -- Mocking axios
    axios.post.mockResolvedValue({ data: { id: 1 } })
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: [{ id: 1, title: 'New Todo' }] })

    await store.addTodo('New Todo')

    expect(axios.post).toHaveBeenCalled()
    expect(store.todos.length).toBe(1)
  })

  it('updates a todo successfully', async () => {
    const store = useTodoStore()

    // @ts-expect-error -- Mocking axios
    axios.put.mockResolvedValue({})
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: [] })

    await store.updateTodo(1, { completed: true })

    expect(axios.put).toHaveBeenCalledWith('http://localhost:8081/api/todos/1', { completed: true })
  })

  // Regression: the form sends '' for cleared date fields, but the backend
  // expects null on its `*time.Time` fields. If either slipped through, the
  // PUT would 400 and the user would see their edit (e.g. tag changes) get
  // silently dropped.
  it('normalizes empty date strings to null so the PUT body is valid', async () => {
    const store = useTodoStore()
    const existing: Todo = {
      id: 1, title: 'Test', description: '', completed: false, priority: 'medium',
      due_date: null, remind_at: null, notified_at: null, repeat: '', tags: ['a'],
      project_id: null, subtasks: [], created_at: '', deleted_at: null,
    }
    store.todos = [existing]

    // @ts-expect-error -- Mocking axios
    axios.put.mockResolvedValue({ data: existing })
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: [] })

    await store.updateTodo(1, {
      tags: ['b'],
      due_date: '',
      remind_at: '',
    })

    expect(axios.put).toHaveBeenCalledWith('http://localhost:8081/api/todos/1', {
      tags: ['b'],
      due_date: null,
      remind_at: null,
    })
  })

  it('deletes a todo successfully', async () => {
    const store = useTodoStore()
    
    // @ts-expect-error -- Mocking axios
    axios.delete.mockResolvedValue({})
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: [] })

    await store.deleteTodo(1)

    expect(axios.delete).toHaveBeenCalledWith('http://localhost:8081/api/todos/1')
  })
})

describe('Todo Store Trash', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches trashed todos', async () => {
    const store = useTodoStore()
    const mockTrash = [
      { id: 2, title: 'Trashed', description: '', completed: false, priority: 'medium', due_date: null, remind_at: null, notified_at: null, repeat: '', tags: [], project_id: null, subtasks: [], created_at: '', deleted_at: '2026-09-01T00:00:00Z' }
    ]

    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: mockTrash })

    await store.fetchTrashTodos()

    expect(axios.get).toHaveBeenCalledWith('http://localhost:8081/api/trash/todos')
    expect(store.trashTodos).toEqual(mockTrash)
  })

  it('restores a trashed todo', async () => {
    const store = useTodoStore()
    store.trashTodos = [
      { id: 1, title: 'Trashed', description: '', completed: false, priority: 'medium', due_date: null, remind_at: null, notified_at: null, repeat: '', tags: [], project_id: null, subtasks: [], created_at: '', deleted_at: '2026-09-01T00:00:00Z' }
    ]

    // @ts-expect-error -- Mocking axios
    axios.post.mockResolvedValue({})

    await store.restoreTrashTodo(1)

    expect(axios.post).toHaveBeenCalledWith('http://localhost:8081/api/trash/todos/1/restore')
    expect(store.trashTodos.length).toBe(0)
  })

  it('purges a trashed todo permanently', async () => {
    const store = useTodoStore()
    store.trashTodos = [
      { id: 1, title: 'Trashed', description: '', completed: false, priority: 'medium', due_date: null, remind_at: null, notified_at: null, repeat: '', tags: [], project_id: null, subtasks: [], created_at: '', deleted_at: '2026-09-01T00:00:00Z' }
    ]

    // @ts-expect-error -- Mocking axios
    axios.delete.mockResolvedValue({})

    await store.purgeTrashTodo(1)

    expect(axios.delete).toHaveBeenCalledWith('http://localhost:8081/api/trash/todos/1')
    expect(store.trashTodos.length).toBe(0)
  })
})
