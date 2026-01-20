import { computed, type Ref } from 'vue'
import type { Todo } from '../stores/todo'

export type ViewType = 'all' | 'inbox' | 'today' | 'upcoming' | 'overdue' | 'project' | 'calendar'
export type SortOption = 'created_desc' | 'due_asc' | 'due_desc' | 'priority_desc'
export type FilterType = 'all' | 'active' | 'completed'

export function useTodoFilter(
  todos: Ref<Todo[]>,
  currentView: Ref<ViewType>,
  currentProjectId: Ref<number | null>,
  searchQuery: Ref<string>,
  sortOption: Ref<SortOption>,
  filter: Ref<FilterType>
) {
  return computed(() => {
    let items = todos.value
    
    // View filtering
    if (currentView.value === 'inbox') {
      items = items.filter(t => t.project_id === null)
    } else if (currentView.value === 'project' && currentProjectId.value) {
      items = items.filter(t => t.project_id === currentProjectId.value)
    } else if (currentView.value === 'today') {
        const now = new Date()
        items = items.filter(t => {
            if (!t.due_date) return false
            const d = new Date(t.due_date)
            return d.getDate() === now.getDate() && 
                   d.getMonth() === now.getMonth() && 
                   d.getFullYear() === now.getFullYear()
        })
    } else if (currentView.value === 'upcoming') {
        const now = new Date()
        const nextWeek = new Date()
        nextWeek.setDate(now.getDate() + 7)
        items = items.filter(t => {
            if (!t.due_date) return false
            const d = new Date(t.due_date)
            // Reset time part for accurate date comparison if needed, 
            // but for "upcoming" generally we want future dates.
            // Let's keep existing logic: d >= now
            return d >= now && d <= nextWeek
        })
    } else if (currentView.value === 'overdue') {
        const now = new Date()
        items = items.filter(t => !t.completed && t.due_date && new Date(t.due_date) < now)
    }
  
    // Search filter
    if (searchQuery.value.trim()) {
        const q = searchQuery.value.toLowerCase()
        items = items.filter(t => 
          t.title.toLowerCase().includes(q) || 
          (t.description && t.description.toLowerCase().includes(q)) ||
          (t.tags && t.tags.some(tag => tag.toLowerCase().includes(q)))
        )
    }
  
    // Sorting
    // Create a copy to avoid mutating the original array in store if it was a direct reference
    // (though filter() creates new array, if no filter ran, it might be ref to store array)
    items = [...items]
    
    items.sort((a, b) => {
        switch (sortOption.value) {
            case 'due_asc':
                if (!a.due_date && !b.due_date) return 0
                if (!a.due_date) return 1 // No due date at bottom
                if (!b.due_date) return -1
                return new Date(a.due_date).getTime() - new Date(b.due_date).getTime()
                
            case 'due_desc':
                if (!a.due_date && !b.due_date) return 0
                if (!a.due_date) return 1 // No due date at bottom
                if (!b.due_date) return -1
                return new Date(b.due_date).getTime() - new Date(a.due_date).getTime()
                
            case 'priority_desc':
                const pMap: Record<string, number> = { high: 3, medium: 2, low: 1 }
                const pA = pMap[a.priority || 'medium'] || 0
                const pB = pMap[b.priority || 'medium'] || 0
                if (pA !== pB) return pB - pA
                // Secondary sort by created_at desc
                return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
                
            case 'created_desc':
            default:
                return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        }
    })
  
    if (filter.value === 'active') return items.filter((t) => !t.completed)
    if (filter.value === 'completed') return items.filter((t) => t.completed)
    return items
  })
}
