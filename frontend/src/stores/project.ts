import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

export interface Project {
  id: number
  name: string
  description: string
  color: string
  created_at: string
}

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const API_URL = 'http://localhost:8081/api/projects'

  const fetchProjects = async () => {
    try {
      const response = await axios.get<Project[]>(API_URL)
      projects.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch projects:', error)
    }
  }

  const addProject = async (name: string, description: string = '', color: string = '#64748B') => {
    try {
      await axios.post(API_URL, { name, description, color })
      await fetchProjects()
    } catch (error) {
      console.error('Failed to add project:', error)
    }
  }

  const updateProject = async (id: number, updates: Partial<Project>) => {
    try {
      await axios.put(`${API_URL}/${id}`, updates)
      await fetchProjects()
    } catch (error) {
      console.error('Failed to update project:', error)
    }
  }

  const deleteProject = async (id: number) => {
    try {
      await axios.delete(`${API_URL}/${id}`)
      await fetchProjects()
    } catch (error) {
      console.error('Failed to delete project:', error)
    }
  }

  return { projects, fetchProjects, addProject, updateProject, deleteProject }
})
