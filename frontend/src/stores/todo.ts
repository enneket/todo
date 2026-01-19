import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

export interface Todo {
  id: number
  title: string
  description: string
  completed: boolean
  priority: 'high' | 'medium' | 'low'
  due_date: string | null
  tags: string[]
  created_at: string
}

export const useTodoStore = defineStore('todo', () => {
  const todos = ref<Todo[]>([])
  // Assume API is running on localhost:8081 (from main.go)
  const API_URL = 'http://localhost:8081/api/todos'

  const fetchTodos = async () => {
    try {
      const response = await axios.get<Todo[]>(API_URL)
      todos.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch todos:', error)
    }
  }

  const addTodo = async (title: string, priority: string = 'medium', dueDate: string | null = null, description: string = '', tags: string[] = []) => {
    try {
      // If dueDate is empty string, send null
      const dateToSend = dueDate === '' ? null : dueDate
      await axios.post(API_URL, { title, description, priority, due_date: dateToSend, tags })
      await fetchTodos()
    } catch (error) {
      console.error('Failed to add todo:', error)
    }
  }

  const updateTodo = async (id: number, updates: Partial<Todo>) => {
    try {
        // Ensure due_date is null if empty string
        if (updates.due_date === '') {
            updates.due_date = null as any
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

  return { todos, fetchTodos, addTodo, updateTodo, deleteTodo }
})
