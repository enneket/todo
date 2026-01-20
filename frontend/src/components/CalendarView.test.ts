import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, VueWrapper, DOMWrapper } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import CalendarView from './CalendarView.vue'
import { createI18n } from 'vue-i18n'

// Mock Phosphor Icons
vi.mock('@phosphor-icons/vue', () => ({
  PhCaretLeft: { template: '<span class="ph-caret-left"></span>' },
  PhCaretRight: { template: '<span class="ph-caret-right"></span>' },
  PhCheckCircle: { template: '<span class="ph-check-circle"></span>' },
  PhCircle: { template: '<span class="ph-circle"></span>' }
}))

// Mock translations
const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      today: 'Today',
      view_month: 'Month',
      view_week: 'Week'
    }
  }
})

describe('CalendarView', () => {
  let wrapper: VueWrapper

  const mockTodos = [
    {
      id: 1,
      title: 'Task 1',
      due_date: new Date().toISOString(), // Today
      completed: false,
      priority: 'high'
    },
    {
      id: 2,
      title: 'Task 2',
      due_date: new Date(new Date().setDate(new Date().getDate() + 1)).toISOString(), // Tomorrow
      completed: true,
      priority: 'medium'
    }
  ]

  beforeEach(() => {
    wrapper = mount(CalendarView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              todo: {
                todos: mockTodos
              }
            }
          }),
          i18n
        ]
      }
    })
  })

  it('renders calendar header with current month', () => {
    const currentMonth = new Date().toLocaleString('default', { month: 'long', year: 'numeric' })
    expect(wrapper.text()).toContain(currentMonth)
    expect(wrapper.text()).toContain('Today')
    expect(wrapper.text()).toContain('Month')
    expect(wrapper.text()).toContain('Week')
  })

  it('renders days grid correctly', () => {
    // Month view should have 42 cells (6 rows * 7 cols)
    const dayCells = wrapper.findAll('.grid-rows-6 > div')
    expect(dayCells).toHaveLength(42)
  })

  it('renders tasks on correct dates', () => {
    // We expect Task 1 to be visible since it's today
    expect(wrapper.text()).toContain('Task 1')
    
    // Depending on if tomorrow is in current view (it should be for month view usually)
    // Task 2 should also be there
    expect(wrapper.text()).toContain('Task 2')
  })

  it('toggles between month and week view', async () => {
    const viewButtons = wrapper.findAll('button')
    // Find Week button (based on text content or order)
    const weekBtn = viewButtons.find((b: DOMWrapper<Element>) => b.text().includes('Week'))
    if (weekBtn) {
      await weekBtn.trigger('click')
    }

    // In week view, grid should still be there but logic changes
    // The implementation keeps grid-rows-6 but changes days array length
    // Wait, let's check implementation of calendarDays in Week view
    // It pushes 7 days.
    
    // However, the template uses `v-for="(day, index) in calendarDays"`
    // So in week view, we expect 7 day cells.
    const dayCells = wrapper.findAll('.grid-rows-6 > div')
    expect(dayCells).toHaveLength(7)
  })

  it('emits edit-task event when clicking a task', async () => {
    const taskEl = wrapper.find('.cursor-pointer') // The task element has cursor-pointer class
    await taskEl.trigger('click')
    
    const emitted = wrapper.emitted('edit-task')
    expect(emitted).toBeTruthy()
    expect(emitted![0][0]).toEqual(mockTodos[0])
  })

  it('navigates to previous/next month', async () => {
    const initialMonth = wrapper.find('h2').text()
    
    // Find prev button (first button in the first group)
    const buttons = wrapper.findAll('button')
    const prevBtn = buttons[0] // Assuming structure: Prev, Today, Next
    
    await prevBtn.trigger('click')
    const newMonth = wrapper.find('h2').text()
    
    expect(newMonth).not.toBe(initialMonth)
  })
})
