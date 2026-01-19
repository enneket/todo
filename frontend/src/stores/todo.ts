import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

export interface Subtask {
  id: number
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
  tags: string[]
  project_id: number | null
  subtasks: Subtask[]
  created_at: string
}

export const useTodoStore = defineStore('todo', () => {
  const todos = ref<Todo[]>([])
  // Assume API is running on localhost:8081 (from main.go)
  const API_URL = 'http://localhost:8081/api/todos'
  const SUBTASK_API_URL = 'http://localhost:8081/api/subtasks'

  const fetchTodos = async () => {
    try {
      const response = await axios.get<Todo[]>(API_URL)
      todos.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch todos:', error)
    }
  }

  const addTodo = async (title: string, priority: string = 'medium', dueDate: string | null = null, description: string = '', tags: string[] = [], projectId: number | null = null) => {
    try {
      // If dueDate is empty string, send null
      const dateToSend = dueDate === '' ? null : dueDate
      await axios.post(API_URL, { title, description, priority, due_date: dateToSend, tags, project_id: projectId })
      await fetchTodos()
    } catch (error) {
      console.error('Failed to add todo:', error)
    }
  }

  const updateTodo = async (id: number, updates: Partial<Todo>) => {
    try {
        // Ensure due_date is null if empty string
        if (updates.due_date === '') {
            updates.due_date = null
        }
      await axios.put(`${API_URL}/${id}`, updates)
      await fetchTodos()
    } catch (error) {
      console.error('Failed to update todo:', error)
    }
  }

  const deleteTodo = async (id: number) => {
    try {
      await axios.delete(`${API_URL}/${id}`)
      await fetchTodos()
    } catch (error) {
      console.error('Failed to delete todo:', error)
    }
  }

  const addSubtask = async (todoId: number, title: string) => {
    try {
      await axios.post(`${API_URL}/${todoId}/subtasks`, { title })
      await fetchTodos()
    } catch (error) {
      console.error('Failed to add subtask:', error)
    }
  }

  const updateSubtask = async (id: number, updates: Partial<Subtask>) => {
    try {
      await axios.put(`${SUBTASK_API_URL}/${id}`, updates)
      await fetchTodos()
    } catch (error) {
      console.error('Failed to update subtask:', error)
    }
  }

  const deleteSubtask = async (id: number) => {
    try {
      await axios.delete(`${SUBTASK_API_URL}/${id}`)
      await fetchTodos()
    } catch (error) {
      console.error('Failed to delete subtask:', error)
    }
  }

  return { todos, fetchTodos, addTodo, updateTodo, deleteTodo, addSubtask, updateSubtask, deleteSubtask }
})
