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
    
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: mockProjects })

    await store.fetchProjects()

    expect(axios.get).toHaveBeenCalledWith('http://localhost:8081/api/projects')
    expect(store.projects).toEqual(mockProjects)
  })

  it('adds a project successfully', async () => {
    const store = useProjectStore()
    
    // @ts-expect-error -- Mocking axios
    axios.post.mockResolvedValue({ data: { id: 1 } })
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: [{ id: 1, name: 'New Project' }] })

    await store.addProject('New Project')

    expect(axios.post).toHaveBeenCalled()
    expect(store.projects.length).toBe(1)
  })

  it('updates a project successfully', async () => {
    const store = useProjectStore()
    
    // @ts-expect-error -- Mocking axios
    axios.put.mockResolvedValue({})
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: [] })

    await store.updateProject(1, { name: 'Updated' })

    expect(axios.put).toHaveBeenCalledWith('http://localhost:8081/api/projects/1', { name: 'Updated' })
  })

  it('deletes a project successfully', async () => {
    const store = useProjectStore()
    
    // @ts-expect-error -- Mocking axios
    axios.delete.mockResolvedValue({})
    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: [] })

    await store.deleteProject(1)

    expect(axios.delete).toHaveBeenCalledWith('http://localhost:8081/api/projects/1')
  })
})

describe('Project Store Trash', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches trashed projects', async () => {
    const store = useProjectStore()
    const mockTrash = [
      { id: 2, name: 'Old Project', description: '', color: '#000', created_at: '', deleted_at: '2026-09-01T00:00:00Z' }
    ]

    // @ts-expect-error -- Mocking axios
    axios.get.mockResolvedValue({ data: mockTrash })

    await store.fetchTrashProjects()

    expect(axios.get).toHaveBeenCalledWith('http://localhost:8081/api/trash/projects')
    expect(store.trashProjects).toEqual(mockTrash)
  })

  it('restores a trashed project', async () => {
    const store = useProjectStore()
    store.trashProjects = [
      { id: 1, name: 'Old Project', description: '', color: '#000', created_at: '', deleted_at: '2026-09-01T00:00:00Z' }
    ]

    // @ts-expect-error -- Mocking axios
    axios.post.mockResolvedValue({})

    await store.restoreTrashProject(1)

    expect(axios.post).toHaveBeenCalledWith('http://localhost:8081/api/trash/projects/1/restore')
    expect(store.trashProjects.length).toBe(0)
  })

  it('purges a trashed project permanently', async () => {
    const store = useProjectStore()
    store.trashProjects = [
      { id: 1, name: 'Old Project', description: '', color: '#000', created_at: '', deleted_at: '2026-09-01T00:00:00Z' }
    ]

    // @ts-expect-error -- Mocking axios
    axios.delete.mockResolvedValue({})

    await store.purgeTrashProject(1)

    expect(axios.delete).toHaveBeenCalledWith('http://localhost:8081/api/trash/projects/1')
    expect(store.trashProjects.length).toBe(0)
  })
})
