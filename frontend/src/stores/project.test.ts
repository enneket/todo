import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useProjectStore } from './project'
import axios from 'axios'

vi.mock('axios')

describe('Project Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches projects successfully', async () => {
    const store = useProjectStore()
    const mockProjects = [
      { id: 1, name: 'Work', description: '', color: '#000', created_at: '' }
    ]
    
    // @ts-ignore
    axios.get.mockResolvedValue({ data: mockProjects })

    await store.fetchProjects()

    expect(axios.get).toHaveBeenCalledWith('http://localhost:8081/api/projects')
    expect(store.projects).toEqual(mockProjects)
  })

  it('adds a project successfully', async () => {
    const store = useProjectStore()
    
    // @ts-ignore
    axios.post.mockResolvedValue({ data: { id: 1 } })
    // @ts-ignore
    axios.get.mockResolvedValue({ data: [{ id: 1, name: 'New Project' }] })

    await store.addProject('New Project')

    expect(axios.post).toHaveBeenCalled()
    expect(store.projects.length).toBe(1)
  })

  it('updates a project successfully', async () => {
    const store = useProjectStore()
    
    // @ts-ignore
    axios.put.mockResolvedValue({})
    // @ts-ignore
    axios.get.mockResolvedValue({ data: [] })

    await store.updateProject(1, { name: 'Updated' })

    expect(axios.put).toHaveBeenCalledWith('http://localhost:8081/api/projects/1', { name: 'Updated' })
  })

  it('deletes a project successfully', async () => {
    const store = useProjectStore()
    
    // @ts-ignore
    axios.delete.mockResolvedValue({})
    // @ts-ignore
    axios.get.mockResolvedValue({ data: [] })

    await store.deleteProject(1)

    expect(axios.delete).toHaveBeenCalledWith('http://localhost:8081/api/projects/1')
  })
})
