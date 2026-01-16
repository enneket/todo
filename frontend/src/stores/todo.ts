import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

export interface Todo {
  id: number
  title: string
  completed: boolean
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

  const addTodo = async (title: string) => {
    try {
      await axios.post(API_URL, { title })
      await fetchTodos()
    } catch (error) {
      console.error('Failed to add todo:', error)
    }
  }

  const toggleTodo = async (todo: Todo) => {
    try {
      await axios.put(`${API_URL}/${todo.id}`, { completed: !todo.completed })
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

  return { todos, fetchTodos, addTodo, toggleTodo, deleteTodo }
})
