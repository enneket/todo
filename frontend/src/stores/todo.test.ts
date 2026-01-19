import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useTodoStore } from './todo'
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
    
    // @ts-ignore
    axios.get.mockResolvedValue({ data: mockTodos })

    await store.fetchTodos()

    expect(axios.get).toHaveBeenCalledWith('http://localhost:8081/api/todos')
    expect(store.todos).toEqual(mockTodos)
  })

  it('adds a todo successfully', async () => {
    const store = useTodoStore()
    
    // @ts-ignore
    axios.post.mockResolvedValue({ data: { id: 1 } })
    // @ts-ignore
    axios.get.mockResolvedValue({ data: [{ id: 1, title: 'New Todo' }] })

    await store.addTodo('New Todo')

    expect(axios.post).toHaveBeenCalled()
    expect(store.todos.length).toBe(1)
  })

  it('updates a todo successfully', async () => {
    const store = useTodoStore()
    
    // @ts-ignore
    axios.put.mockResolvedValue({})
    // @ts-ignore
    axios.get.mockResolvedValue({ data: [] })

    await store.updateTodo(1, { completed: true })

    expect(axios.put).toHaveBeenCalledWith('http://localhost:8081/api/todos/1', { completed: true })
  })

  it('deletes a todo successfully', async () => {
    const store = useTodoStore()
    
    // @ts-ignore
    axios.delete.mockResolvedValue({})
    // @ts-ignore
    axios.get.mockResolvedValue({ data: [] })

    await store.deleteTodo(1)

    expect(axios.delete).toHaveBeenCalledWith('http://localhost:8081/api/todos/1')
  })
})
