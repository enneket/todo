import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

export interface Project {
  id: number
  name: string
  description: string
  color: string
  created_at: string
  deleted_at: string | null
}

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const trashProjects = ref<Project[]>([])
  const API_URL = 'http://localhost:8081/api/projects'
  const TRASH_API_URL = 'http://localhost:8081/api/trash/projects'

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

  const fetchTrashProjects = async () => {
    try {
      const response = await axios.get<Project[]>(TRASH_API_URL)
      trashProjects.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch trashed projects:', error)
    }
  }

  const restoreTrashProject = async (id: number) => {
    try {
      await axios.post(`${TRASH_API_URL}/${id}/restore`)
      trashProjects.value = trashProjects.value.filter(p => p.id !== id)
    } catch (error) {
      console.error('Failed to restore project:', error)
    }
  }

  const purgeTrashProject = async (id: number) => {
    try {
      await axios.delete(`${TRASH_API_URL}/${id}`)
      trashProjects.value = trashProjects.value.filter(p => p.id !== id)
    } catch (error) {
      console.error('Failed to purge project:', error)
    }
  }

  return { projects, trashProjects, fetchProjects, addProject, updateProject, deleteProject, fetchTrashProjects, restoreTrashProject, purgeTrashProject }
})
