import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { ref } from 'vue'
import { useTodoFilter } from './useTodoFilter'
import type { Todo } from '../stores/todo'

describe('useTodoFilter', () => {
  const mockDate = new Date('2024-01-15T12:00:00Z')

  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(mockDate)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  const createTodo = (overrides: Partial<Todo> = {}): Todo => ({
    id: 1,
    title: 'Test Todo',
    description: '',
    completed: false,
    priority: 'medium',
    due_date: null,
    remind_at: null,
    repeat: '',
    tags: [],
    project_id: null,
    subtasks: [],
    created_at: new Date().toISOString(),
    ...overrides
  })

  it('filters by Smart View: Today', () => {
    const todos = ref([
      createTodo({ id: 1, title: 'Today Task', due_date: '2024-01-15T10:00:00Z' }),
      createTodo({ id: 2, title: 'Tomorrow Task', due_date: '2024-01-16T10:00:00Z' }),
      createTodo({ id: 3, title: 'No Date' })
    ])

    const filtered = useTodoFilter(
      todos,
      ref('today'),
      ref(null),
      ref(''),
      ref('created_desc'),
      ref('all')
    )

    expect(filtered.value).toHaveLength(1)
    expect(filtered.value[0].id).toBe(1)
  })

  it('filters by Smart View: Upcoming', () => {
    // Adjust test data to be in future for upcoming
    const futureTodos = ref([
        createTodo({ id: 5, title: 'Future Today', due_date: '2024-01-15T14:00:00Z' }),
        createTodo({ id: 6, title: 'Tomorrow', due_date: '2024-01-16T12:00:00Z' })
    ])
    
    const filtered = useTodoFilter(
      futureTodos,
      ref('upcoming'),
      ref(null),
      ref(''),
      ref('created_desc'),
      ref('all')
    )

    expect(filtered.value).toHaveLength(2)
  })

  it('filters by Smart View: Overdue', () => {
    const todos = ref([
      createTodo({ id: 1, title: 'Overdue', due_date: '2024-01-10T10:00:00Z', completed: false }),
      createTodo({ id: 2, title: 'Overdue Completed', due_date: '2024-01-10T10:00:00Z', completed: true }),
      createTodo({ id: 3, title: 'Future', due_date: '2024-01-20T10:00:00Z', completed: false })
    ])

    const filtered = useTodoFilter(
      todos,
      ref('overdue'),
      ref(null),
      ref(''),
      ref('created_desc'),
      ref('all')
    )

    expect(filtered.value).toHaveLength(1)
    expect(filtered.value[0].id).toBe(1)
  })

  it('sorts by Priority Desc', () => {
    const todos = ref([
      createTodo({ id: 1, priority: 'low' }),
      createTodo({ id: 2, priority: 'high' }),
      createTodo({ id: 3, priority: 'medium' })
    ])

    const filtered = useTodoFilter(
      todos,
      ref('all'),
      ref(null),
      ref(''),
      ref('priority_desc'),
      ref('all')
    )

    expect(filtered.value.map(t => t.id)).toEqual([2, 3, 1]) // High, Medium, Low
  })

  it('sorts by Due Date Asc', () => {
    const todos = ref([
      createTodo({ id: 1, due_date: '2024-01-20T10:00:00Z' }),
      createTodo({ id: 2, due_date: '2024-01-10T10:00:00Z' }),
      createTodo({ id: 3, due_date: null })
    ])

    const filtered = useTodoFilter(
      todos,
      ref('all'),
      ref(null),
      ref(''),
      ref('due_asc'),
      ref('all')
    )

    // Expected: Earliest first, null last
    expect(filtered.value.map(t => t.id)).toEqual([2, 1, 3])
  })
})
