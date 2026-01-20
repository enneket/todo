<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTodoStore, type Todo, type Subtask } from './stores/todo'
import { useThemeStore } from './stores/theme'
import { useProjectStore } from './stores/project'
import { 
  PhPlus, PhTrash, PhCheckCircle, PhCircle, PhTranslate, PhPencil, 
  PhClock, PhMagnifyingGlass, PhWarning, PhChartBar, PhSun, PhMoon, 
  PhDesktop, PhList, PhFolder, PhCheckSquare, PhX, PhBell, PhArrowsClockwise,
  PhCalendar, PhCalendarCheck, PhWarningCircle, PhSortAscending, PhSortDescending
} from '@phosphor-icons/vue'
import BaseSelect from './components/BaseSelect.vue'
import StatisticsPanel from './components/StatisticsPanel.vue'
import CalendarView from './components/CalendarView.vue'
import { useTodoFilter, type ViewType } from './composables/useTodoFilter'

const { t, locale } = useI18n()
const todoStore = useTodoStore()
const themeStore = useThemeStore()
const projectStore = useProjectStore()

const showStats = ref(false)
const showSidebar = ref(true) // For mobile/responsive toggle if needed

const priorityOptions = computed(() => [
    { value: 'high', label: t('priority_high') },
    { value: 'medium', label: t('priority_medium') },
    { value: 'low', label: t('priority_low') }
])

const projectOptions = computed(() => [
  { value: -1, label: t('inbox') }, // -1 as placeholder for null/Inbox
  ...projectStore.projects.map(p => ({ value: p.id, label: p.name }))
])

const getDefaultDueDate = () => {
  const d = new Date()
  d.setDate(d.getDate() + 1) // Tomorrow
  d.setMinutes(d.getMinutes() - d.getTimezoneOffset()) // Adjust for local time ISO string
  return d.toISOString().slice(0, 16)
}

// Todo Form
const form = ref({
  id: undefined as number | undefined,
  title: '',
  description: '',
  priority: 'medium',
  due_date: getDefaultDueDate(),
  remind_at: '',
  repeat: '',
  tags: '',
  project_id: -1 as number // -1 for Inbox (null)
})

const repeatOptions = computed(() => [
    { value: '', label: t('repeat_none') || 'No Repeat' },
    { value: 'daily', label: t('repeat_daily') || 'Daily' },
    { value: 'weekly', label: t('repeat_weekly') || 'Weekly' },
    { value: 'monthly', label: t('repeat_monthly') || 'Monthly' },
    { value: 'weekdays', label: t('repeat_weekdays') || 'Weekdays' },
])

// Project Form
const showProjectModal = ref(false)
const projectForm = ref({
  name: '',
  description: '',
  color: '#3B82F6'
})

const showModal = ref(false)
const isEditing = computed(() => !!form.value.id)

const filter = ref<'all' | 'active' | 'completed'>('all')
const searchQuery = ref('')
const currentView = ref<ViewType>('all')
const currentProjectId = ref<number | null>(null)
const sortOption = ref<'created_desc' | 'due_asc' | 'due_desc' | 'priority_desc'>('created_desc')
const showSortMenu = ref(false)

const currentProjectName = computed(() => {
  if (currentView.value === 'all') return t('all_tasks')
  if (currentView.value === 'inbox') return t('inbox')
  if (currentView.value === 'today') return t('today')
  if (currentView.value === 'upcoming') return t('upcoming')
  if (currentView.value === 'overdue') return t('overdue')
  if (currentView.value === 'calendar') return t('calendar')
  
  if (currentView.value === 'project' && currentProjectId.value) {
      const p = projectStore.projects.find(p => p.id === currentProjectId.value)
      return p ? p.name : 'Unknown'
  }
  return t('all_tasks')
})

// Subtasks Logic
const newSubtaskTitle = ref('')

onMounted(() => {
  todoStore.fetchTodos()
  projectStore.fetchProjects()
})

const openAddModal = () => {
  form.value = {
    id: undefined,
    title: '',
    description: '',
    priority: 'medium',
    due_date: getDefaultDueDate(),
    remind_at: '',
    repeat: '',
    tags: '',
    project_id: currentProjectId.value && currentProjectId.value > 0 ? currentProjectId.value : -1
  }
  newSubtaskTitle.value = ''
  showModal.value = true
}

const openEditModal = (todo: Todo) => {
  let dateStr = ''
  if (todo.due_date) {
      const d = new Date(todo.due_date)
      const offset = d.getTimezoneOffset() * 60000
      const localISOTime = (new Date(d.getTime() - offset)).toISOString().slice(0, 16)
      dateStr = localISOTime
  }

  let remindStr = ''
  if (todo.remind_at) {
      const d = new Date(todo.remind_at)
      const offset = d.getTimezoneOffset() * 60000
      const localISOTime = (new Date(d.getTime() - offset)).toISOString().slice(0, 16)
      remindStr = localISOTime
  }

  form.value = {
    id: todo.id,
    title: todo.title,
    description: todo.description || '',
    priority: todo.priority || 'medium',
    due_date: dateStr,
    remind_at: remindStr,
    repeat: todo.repeat || '',
    tags: todo.tags ? todo.tags.join(', ') : '',
    project_id: todo.project_id || -1
  }
  showModal.value = true
}

const handleSubmit = async () => {
  if (!form.value.title.trim()) return

  const dateToSend = form.value.due_date ? new Date(form.value.due_date).toISOString() : ''
  const remindToSend = form.value.remind_at ? new Date(form.value.remind_at).toISOString() : ''
  const tags = form.value.tags.split(',').map(t => t.trim()).filter(t => t)
  const projectId = form.value.project_id === -1 ? null : form.value.project_id

  if (isEditing.value && form.value.id) {
      await todoStore.updateTodo(form.value.id, {
        title: form.value.title,
        description: form.value.description,
        priority: form.value.priority as 'high' | 'medium' | 'low',
        due_date: dateToSend,
        remind_at: remindToSend,
        repeat: form.value.repeat,
        tags: tags,
        project_id: projectId
      })
  } else {
      await todoStore.addTodo(
        form.value.title, 
        form.value.priority, 
        dateToSend, 
        form.value.description, 
        tags,
        projectId,
        remindToSend,
        form.value.repeat
      )
  }
  showModal.value = false
}

// Project Submit
const handleProjectSubmit = async () => {
  if (!projectForm.value.name.trim()) return
  await projectStore.addProject(projectForm.value.name, projectForm.value.description, projectForm.value.color)
  showProjectModal.value = false
  projectForm.value = { name: '', description: '', color: '#3B82F6' }
}

const deleteProject = async (id: number) => {
  if (confirm(t('confirm_delete_project') || 'Delete this project?')) {
    await projectStore.deleteProject(id)
    if (currentProjectId.value === id) {
      currentProjectId.value = null
      currentView.value = 'all'
    }
  }
}

// Subtask Handlers
const addSubtask = async () => {
  if (!newSubtaskTitle.value.trim() || !form.value.id) return
  await todoStore.addSubtask(form.value.id, newSubtaskTitle.value)
  newSubtaskTitle.value = ''
}

const toggleSubtask = async (subtask: Subtask) => {
  await todoStore.updateSubtask(subtask.id, { completed: !subtask.completed })
}

const deleteSubtask = async (id: number) => {
  await todoStore.deleteSubtask(id)
}

const filteredTodos = useTodoFilter(
  computed(() => todoStore.todos),
  currentView,
  currentProjectId,
  searchQuery,
  sortOption,
  filter
)

const toggleLanguage = () => {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const getPriorityColor = (p: string) => {
  switch (p) {
    case 'high': return 'text-red-700 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-900/50'
    case 'medium': return 'text-amber-700 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-100 dark:border-amber-900/50'
    case 'low': return 'text-emerald-700 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-100 dark:border-emerald-900/50'
    default: return 'text-slate-600 dark:text-slate-400 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700'
  }
}

const formatDate = (dateStr: string | null) => {
    if (!dateStr) return ''
    return new Date(dateStr).toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
        month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    })
}

const getDueDateStatus = (dateStr: string | null, completed: boolean) => {
    if (!dateStr || completed) return { status: '', label: '', color: '' }
    
    const now = new Date()
    const due = new Date(dateStr)
    const diff = due.getTime() - now.getTime()
    const hours = diff / (1000 * 60 * 60)
    
    if (diff < 0) {
        return { status: 'overdue', label: 'overdue', color: 'text-red-600 dark:text-red-400 font-bold' }
    } else if (hours <= 24) {
        return { status: 'soon', label: 'due_soon', color: 'text-orange-600 dark:text-orange-400 font-medium' }
    }
    return { status: 'future', label: '', color: 'text-slate-500 dark:text-slate-500' }
}

const getRepeatLabel = (value: string) => {
    const opt = repeatOptions.value.find(o => o.value === value)
    return opt ? opt.label : value
}
// Ignore unused warning for now as it's used in template
void getRepeatLabel


const currentTodoSubtasks = computed(() => {
  if (!form.value.id) return []
  const todo = todoStore.todos.find(t => t.id === form.value.id)
  return todo ? todo.subtasks : []
})
</script>

<template>
  <div class="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-800 dark:text-slate-200 font-sans selection:bg-primary-100 selection:text-primary-700 flex">
    
    <!-- Sidebar -->
    <aside 
      class="w-64 bg-white dark:bg-slate-900 border-r border-slate-100 dark:border-slate-800 flex-shrink-0 flex flex-col transition-all"
      :class="{ '-ml-64': !showSidebar }"
    >
      <div class="p-6">
        <h1 class="text-xl font-bold font-display tracking-tight text-slate-900 dark:text-white flex items-center gap-2 mb-8">
            <PhCheckCircle weight="fill" class="text-primary-500" />
            Todo App
        </h1>

        <nav class="space-y-1">
          <button 
            @click="currentView = 'all'; currentProjectId = null"
            class="w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
            :class="currentView === 'all' ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-400' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
          >
            <PhList size="18" />
            {{ t('all_tasks') || 'All Tasks' }}
          </button>
          <button 
            @click="currentView = 'inbox'; currentProjectId = null"
            class="w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
            :class="currentView === 'inbox' ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-400' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
          >
            <PhFolder size="18" />
            {{ t('inbox') || 'Inbox' }}
          </button>
          
          <div class="pt-2 pb-1">
            <div class="h-px bg-slate-100 dark:bg-slate-800 mx-3"></div>
          </div>

          <button 
            @click="currentView = 'today'; currentProjectId = null"
            class="w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
            :class="currentView === 'today' ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-400' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
          >
            <PhCalendar size="18" />
            {{ t('today') || 'Today' }}
          </button>
          <button 
            @click="currentView = 'upcoming'; currentProjectId = null"
            class="w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
            :class="currentView === 'upcoming' ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-400' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
          >
            <PhCalendarCheck size="18" />
            {{ t('upcoming') || 'Upcoming' }}
          </button>
          <button 
            @click="currentView = 'overdue'; currentProjectId = null"
            class="w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
            :class="currentView === 'overdue' ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-400' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
          >
            <PhWarningCircle size="18" />
            {{ t('overdue') || 'Overdue' }}
          </button>
          <button 
            @click="currentView = 'calendar'; currentProjectId = null"
            class="w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
            :class="currentView === 'calendar' ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-400' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
          >
            <PhCalendar size="18" />
            {{ t('calendar') || 'Calendar' }}
          </button>
        </nav>

        <div class="mt-8">
          <div class="flex items-center justify-between px-3 mb-2">
            <h3 class="text-xs font-bold text-slate-500 uppercase tracking-wider">{{ t('projects') || 'Projects' }}</h3>
            <button @click="showProjectModal = true" class="text-slate-500 hover:text-primary-500 transition-colors" :aria-label="t('add_project') || 'Add Project'">
              <PhPlus size="14" weight="bold" />
            </button>
          </div>
          <nav class="space-y-1">
            <div 
              v-for="project in projectStore.projects" 
              :key="project.id"
              class="group flex items-center justify-between px-3 py-2 rounded-lg text-sm font-medium transition-colors cursor-pointer"
              :class="currentView === 'project' && currentProjectId === project.id ? 'bg-slate-100 dark:bg-slate-800 text-slate-900 dark:text-white' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
              @click="currentView = 'project'; currentProjectId = project.id"
            >
              <div class="flex items-center gap-2">
                <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: project.color }"></span>
                {{ project.name }}
              </div>
              <button @click.stop="deleteProject(project.id)" class="opacity-0 group-hover:opacity-100 text-slate-400 hover:text-red-500" :aria-label="t('delete') || 'Delete'">
                <PhX size="14" />
              </button>
            </div>
          </nav>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 min-w-0 h-screen overflow-y-auto p-6 md:p-12">
      <div :class="['mx-auto space-y-8 transition-all', currentView === 'calendar' ? 'max-w-6xl h-[calc(100vh-6rem)] flex flex-col' : 'max-w-3xl']">
        
        <!-- Header -->
        <header class="flex flex-col md:flex-row md:items-end justify-between gap-4">
          <div>
              <h2 class="text-3xl font-bold font-display tracking-tight text-slate-900 dark:text-white flex items-center gap-3">
                  {{ currentProjectName }}
              </h2>
              <p class="text-slate-500 dark:text-slate-400 mt-1 font-medium">{{ new Date().toLocaleDateString(locale === 'zh' ? 'zh-CN' : 'en-US', { weekday: 'long', month: 'long', day: 'numeric' }) }}</p>
          </div>
          
          <div class="flex items-center gap-3">
              <div class="relative group">
                  <PhMagnifyingGlass class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-primary-600 transition-colors" size="18" />
                  <input
                      v-model="searchQuery"
                      type="text"
                      :placeholder="t('search')"
                      :aria-label="t('search') || 'Search'"
                      class="pl-10 pr-4 py-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-primary-100 dark:focus:ring-primary-900 focus:border-primary-500 w-40 focus:w-64 transition-all shadow-sm dark:text-white dark:placeholder-slate-500 text-slate-700"
                  />
              </div>
              <button
                @click="showSortMenu = !showSortMenu"
                class="relative p-2 text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 hover:bg-white dark:hover:bg-slate-900 hover:shadow-sm rounded-full transition-all"
                :aria-label="t('sort_by')"
              >
                <PhSortAscending v-if="sortOption === 'due_asc'" size="24" />
                <PhSortDescending v-else size="24" />
                
                <!-- Sort Menu -->
                <div v-if="showSortMenu" class="absolute right-0 top-full mt-2 w-48 bg-white dark:bg-slate-900 rounded-xl shadow-xl border border-slate-100 dark:border-slate-800 z-50 py-1 overflow-hidden">
                    <button @click="sortOption = 'created_desc'; showSortMenu = false" class="w-full text-left px-4 py-2 text-sm hover:bg-slate-50 dark:hover:bg-slate-800 flex items-center justify-between" :class="sortOption === 'created_desc' ? 'text-primary-600 font-medium' : 'text-slate-600 dark:text-slate-400'">
                        {{ t('sort_created') }}
                        <PhCheckCircle v-if="sortOption === 'created_desc'" weight="fill" size="16" />
                    </button>
                    <button @click="sortOption = 'due_asc'; showSortMenu = false" class="w-full text-left px-4 py-2 text-sm hover:bg-slate-50 dark:hover:bg-slate-800 flex items-center justify-between" :class="sortOption === 'due_asc' ? 'text-primary-600 font-medium' : 'text-slate-600 dark:text-slate-400'">
                        {{ t('sort_due_date') }} ({{ t('ascending') }})
                        <PhCheckCircle v-if="sortOption === 'due_asc'" weight="fill" size="16" />
                    </button>
                    <button @click="sortOption = 'due_desc'; showSortMenu = false" class="w-full text-left px-4 py-2 text-sm hover:bg-slate-50 dark:hover:bg-slate-800 flex items-center justify-between" :class="sortOption === 'due_desc' ? 'text-primary-600 font-medium' : 'text-slate-600 dark:text-slate-400'">
                        {{ t('sort_due_date') }} ({{ t('descending') }})
                        <PhCheckCircle v-if="sortOption === 'due_desc'" weight="fill" size="16" />
                    </button>
                    <button @click="sortOption = 'priority_desc'; showSortMenu = false" class="w-full text-left px-4 py-2 text-sm hover:bg-slate-50 dark:hover:bg-slate-800 flex items-center justify-between" :class="sortOption === 'priority_desc' ? 'text-primary-600 font-medium' : 'text-slate-600 dark:text-slate-400'">
                        {{ t('sort_priority') }}
                        <PhCheckCircle v-if="sortOption === 'priority_desc'" weight="fill" size="16" />
                    </button>
                </div>
              </button>
              <button
                @click="themeStore.toggleTheme()"
                class="p-2 text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 hover:bg-white dark:hover:bg-slate-900 hover:shadow-sm rounded-full transition-all"
                :title="themeStore.theme"
                :aria-label="t('theme_toggle') || 'Toggle Theme'"
              >
                <PhSun v-if="themeStore.theme === 'light'" size="24" />
                <PhMoon v-else-if="themeStore.theme === 'dark'" size="24" />
                <PhDesktop v-else size="24" />
              </button>
              <button
                @click="showStats = !showStats"
                class="p-2 text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 hover:bg-white dark:hover:bg-slate-900 hover:shadow-sm rounded-full transition-all"
                :title="t('statistics')"
                :aria-label="t('statistics') || 'Statistics'"
                :class="{ 'text-primary-600 bg-white dark:bg-slate-900 shadow-sm': showStats }"
              >
                <PhChartBar size="24" />
              </button>
              <button
                @click="toggleLanguage"
                class="p-2 text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 hover:bg-white dark:hover:bg-slate-900 hover:shadow-sm rounded-full transition-all"
                :title="t('language')"
                :aria-label="t('language') || 'Language'"
              >
                <PhTranslate size="24" />
              </button>
          </div>
        </header>

        <!-- Stats Panel -->
        <transition
          enter-active-class="transition ease-out duration-200"
          enter-from-class="opacity-0 -translate-y-4"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition ease-in duration-150"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 -translate-y-4"
        >
          <StatisticsPanel v-if="showStats && currentView !== 'calendar'" />
        </transition>

        <CalendarView 
          v-if="currentView === 'calendar'" 
          class="flex-1 min-h-0 shadow-sm"
          @edit-task="openEditModal"
        />

        <div v-else class="space-y-8">
        <!-- Add Button (Hero) -->
        <button
          @click="openAddModal"
          class="w-full py-4 bg-white dark:bg-slate-900 hover:bg-primary-50/50 dark:hover:bg-primary-900/10 border border-slate-200 dark:border-slate-800 hover:border-primary-200 dark:hover:border-primary-800 rounded-2xl shadow-sm hover:shadow-md transition-all group flex items-center justify-center gap-3"
        >
          <div class="bg-primary-50 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 p-1.5 rounded-lg group-hover:scale-110 transition-transform">
              <PhPlus weight="bold" size="20" />
          </div>
          <span class="text-slate-600 dark:text-slate-300 font-medium text-lg group-hover:text-primary-700 dark:group-hover:text-primary-400 transition-colors">{{ t('add') }}</span>
        </button>

        <!-- Filters -->
        <div class="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-hide">
          <button
            v-for="f in ['all', 'active', 'completed']"
            :key="f"
            @click="filter = f as any"
            class="px-4 py-1.5 rounded-full text-sm font-medium transition-all whitespace-nowrap"
            :class="
              filter === f
                ? 'bg-slate-800 dark:bg-primary-600 text-white shadow-md shadow-slate-200 dark:shadow-none'
                : 'bg-white dark:bg-slate-900 text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 shadow-sm border border-slate-100 dark:border-slate-800'
            "
          >
            {{ t(f) }}
          </button>
        </div>

        <!-- List -->
        <div class="space-y-3">
          <transition-group name="list">
              <div
              v-for="todo in filteredTodos"
              :key="todo.id"
              class="bg-white dark:bg-slate-900 p-5 rounded-2xl shadow-sm border border-slate-100 dark:border-slate-800 hover:shadow-md hover:border-primary-100 dark:hover:border-primary-900 transition-all group relative overflow-hidden"
              >
              <!-- Completion Check -->
              <div class="flex items-start gap-4">
                  <button
                          @click="todoStore.updateTodo(todo.id, { completed: !todo.completed })"
                          class="mt-1 text-slate-400 dark:text-slate-600 hover:text-primary-500 transition-colors flex-shrink-0"
                          :aria-label="todo.completed ? t('mark_incomplete') : t('mark_complete')"
                      >
                          <PhCheckCircle v-if="todo.completed" weight="fill" class="text-emerald-500 scale-110" size="24" />
                          <PhCircle v-else weight="bold" size="24" />
                      </button>

                  <div class="flex-1 min-w-0">
                      <div class="flex flex-wrap items-center gap-2 mb-1.5">
                          <span :class="['text-[10px] uppercase tracking-wider font-bold px-2 py-0.5 rounded-md', getPriorityColor(todo.priority || 'medium')]">
                              {{ t('priority_' + (todo.priority || 'medium')) }}
                          </span>
                          <span v-if="todo.project_id" class="text-[10px] uppercase tracking-wider font-bold px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-400 flex items-center gap-1">
                             <span class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: projectStore.projects.find(p => p.id === todo.project_id)?.color }"></span>
                             {{ projectStore.projects.find(p => p.id === todo.project_id)?.name }}
                          </span>
                          <span v-if="todo.due_date" :class="['text-xs flex items-center gap-1 font-medium bg-slate-50 dark:bg-slate-800 px-2 py-0.5 rounded-md', getDueDateStatus(todo.due_date, todo.completed).color]">
                              <PhClock v-if="getDueDateStatus(todo.due_date, todo.completed).status !== 'overdue'" size="12" weight="bold" />
                              <PhWarning v-else size="12" weight="fill" />
                              {{ formatDate(todo.due_date) }}
                              <span v-if="getDueDateStatus(todo.due_date, todo.completed).label">
                                  ({{ t(getDueDateStatus(todo.due_date, todo.completed).label) }})
                              </span>
                          </span>
                          <span v-if="todo.remind_at" class="text-xs px-2 py-0.5 rounded-md bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 font-medium flex items-center gap-1" :title="formatDate(todo.remind_at)">
                              <PhBell size="12" weight="bold" />
                          </span>
                          <span v-if="todo.repeat" class="text-xs px-2 py-0.5 rounded-md bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 font-medium flex items-center gap-1">
                              <PhArrowsClockwise size="12" weight="bold" />
                              {{ getRepeatLabel(todo.repeat) }}
                          </span>
                          <div class="flex gap-1">
                              <span v-for="tag in todo.tags" :key="tag" class="text-xs px-2 py-0.5 rounded-md bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 font-medium">
                                  #{{ tag }}
                              </span>
                          </div>
                      </div>
                      
                      <h3
                          class="text-lg font-medium text-slate-800 dark:text-slate-200 transition-all leading-snug"
                          :class="{ 'line-through text-slate-400 dark:text-slate-600': todo.completed }"
                      >
                          {{ todo.title }}
                      </h3>
                      
                      <p v-if="todo.description" class="text-sm text-slate-600 dark:text-slate-400 mt-1.5 leading-relaxed">
                          {{ todo.description }}
                      </p>
                      
                      <!-- Subtasks Preview -->
                      <div v-if="todo.subtasks && todo.subtasks.length > 0" class="mt-3 space-y-1">
                        <div class="flex items-center gap-2 text-xs text-slate-500">
                          <PhCheckSquare weight="bold" />
                          <span>{{ todo.subtasks.filter(s => s.completed).length }}/{{ todo.subtasks.length }}</span>
                        </div>
                        <div class="w-full bg-slate-100 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
                          <div class="h-full bg-primary-500 transition-all duration-300" :style="{ width: `${(todo.subtasks.filter(s => s.completed).length / todo.subtasks.length) * 100}%` }"></div>
                        </div>
                      </div>
                  </div>

                  <!-- Actions (Hover) -->
                  <div class="flex flex-col gap-1 opacity-0 group-hover:opacity-100 transition-opacity absolute right-4 top-4 bg-white dark:bg-slate-900 pl-2">
                      <button
                          @click="openEditModal(todo)"
                          class="p-2 text-slate-400 hover:text-primary-600 dark:hover:text-primary-400 hover:bg-primary-50 dark:hover:bg-primary-900/30 rounded-lg transition-all"
                          :title="t('edit')"
                          :aria-label="t('edit') || 'Edit'"
                      >
                          <PhPencil size="18" weight="bold" />
                      </button>
                      <button
                          @click="todoStore.deleteTodo(todo.id)"
                          class="p-2 text-slate-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-all"
                          :aria-label="t('delete') || 'Delete'"
                      >
                          <PhTrash size="18" weight="bold" />
                      </button>
                  </div>
              </div>
              </div>
          </transition-group>

          <div v-if="filteredTodos.length === 0" class="py-12 text-center">
              <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-300 dark:text-slate-600 mb-4">
                  <PhCheckCircle size="32" weight="duotone" />
              </div>
              <p class="text-slate-400 dark:text-slate-500 font-medium">{{ t('no_tasks') || 'No tasks found' }}</p>
          </div>
        </div>
        </div>
      </div>
    </main>
    
    <!-- Add/Edit Modal (Backdrop blur) -->
    <div v-if="showModal" class="fixed inset-0 bg-slate-900/20 dark:bg-slate-900/50 backdrop-blur-sm flex items-center justify-center p-4 z-50 transition-all">
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-0 w-full max-w-lg shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800 transform transition-all scale-100 max-h-[90vh] flex flex-col">
            <!-- Modal Header -->
            <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/50 flex justify-between items-center flex-shrink-0">
                <h2 class="text-lg font-bold font-display text-slate-800 dark:text-white">{{ isEditing ? t('edit') : t('add') }}</h2>
                <button @click="showModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200">
                    <span class="text-2xl leading-none">&times;</span>
                </button>
            </div>

            <!-- Modal Body -->
            <div class="p-6 space-y-5 overflow-y-auto">
                <!-- Title -->
                <div>
                    <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('modal_title') }}</label>
                    <input 
                        v-model="form.title" 
                        type="text" 
                        class="w-full px-4 py-3 rounded-xl bg-slate-100 dark:bg-slate-800 border-none focus:ring-2 focus:ring-primary-100 dark:focus:ring-primary-900 text-slate-800 dark:text-white placeholder-slate-400 font-medium transition-shadow" 
                        :placeholder="t('placeholder')" 
                        autofocus
                    />
                </div>
                
                <!-- Description -->
                <div>
                    <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('description') }}</label>
                    <textarea 
                        v-model="form.description" 
                        class="w-full px-4 py-3 rounded-xl bg-slate-100 dark:bg-slate-800 border-none focus:ring-2 focus:ring-primary-100 dark:focus:ring-primary-900 text-slate-800 dark:text-white placeholder-slate-400 resize-none h-24 transition-shadow" 
                        :placeholder="t('description')"
                    ></textarea>
                </div>
                
                <!-- Metadata Row -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                         <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('priority') }}</label>
                         <BaseSelect v-model="form.priority" :options="priorityOptions" class="w-full" />
                    </div>
                    <div>
                         <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('due_date') }}</label>
                         <input 
                            type="datetime-local" 
                            v-model="form.due_date" 
                            class="w-full px-4 py-2.5 rounded-xl bg-slate-100 dark:bg-slate-800 border-none focus:ring-2 focus:ring-primary-100 dark:focus:ring-primary-900 text-slate-800 dark:text-white text-sm h-[42px] dark:[color-scheme:dark]" 
                        />
                    </div>
                </div>

                <!-- Remind & Repeat -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                         <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('remind_at') }}</label>
                         <input 
                            type="datetime-local" 
                            v-model="form.remind_at" 
                            class="w-full px-4 py-2.5 rounded-xl bg-slate-100 dark:bg-slate-800 border-none focus:ring-2 focus:ring-primary-100 dark:focus:ring-primary-900 text-slate-800 dark:text-white text-sm h-[42px] dark:[color-scheme:dark]" 
                        />
                    </div>
                    <div>
                         <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('repeat') }}</label>
                         <BaseSelect v-model="form.repeat" :options="repeatOptions" class="w-full" />
                    </div>
                </div>

                <!-- Project & Tags -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                         <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('project') }}</label>
                         <BaseSelect v-model="form.project_id" :options="projectOptions" class="w-full" />
                    </div>
                    <div>
                         <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('tags_label') }}</label>
                         <input 
                            v-model="form.tags" 
                            type="text" 
                            :placeholder="t('tags')" 
                            class="w-full px-4 py-2.5 rounded-xl bg-slate-100 dark:bg-slate-800 border-none focus:ring-2 focus:ring-primary-100 dark:focus:ring-primary-900 text-slate-800 dark:text-white placeholder-slate-400 transition-shadow text-sm h-[42px]" 
                        />
                    </div>
                </div>

                <!-- Subtasks (Only in Edit Mode) -->
                <div v-if="isEditing">
                    <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">{{ t('subtasks') }}</label>
                    <div class="space-y-2 mb-3">
                      <div 
                        v-for="subtask in currentTodoSubtasks" 
                        :key="subtask.id"
                        class="flex items-center gap-3 p-2 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800 group"
                      >
                        <button 
                          @click="toggleSubtask(subtask)"
                          class="text-slate-400 hover:text-primary-500 transition-colors"
                        >
                          <PhCheckSquare v-if="subtask.completed" weight="fill" class="text-primary-500" size="20" />
                          <PhCheckSquare v-else size="20" />
                        </button>
                        <span 
                          class="flex-1 text-sm text-slate-700 dark:text-slate-300"
                          :class="{ 'line-through text-slate-400': subtask.completed }"
                        >
                          {{ subtask.title }}
                        </span>
                        <button @click="deleteSubtask(subtask.id)" class="opacity-0 group-hover:opacity-100 text-slate-400 hover:text-red-500">
                          <PhX size="16" />
                        </button>
                      </div>
                    </div>
                    <div class="flex items-center gap-2">
                      <PhPlus class="text-slate-400" />
                      <input 
                        v-model="newSubtaskTitle"
                        @keydown.enter.prevent="addSubtask"
                        type="text" 
                        :placeholder="t('add_subtask')" 
                        class="flex-1 bg-transparent border-none focus:ring-0 p-0 text-sm text-slate-800 dark:text-white placeholder-slate-400" 
                      />
                    </div>
                </div>
            </div>

            <!-- Modal Footer -->
            <div class="px-6 py-4 bg-slate-50 dark:bg-slate-800/50 border-t border-slate-100 dark:border-slate-800 flex justify-end gap-3 flex-shrink-0">
                <button 
                    @click="showModal = false" 
                    class="px-5 py-2.5 text-slate-500 dark:text-slate-400 font-medium hover:bg-slate-200/50 dark:hover:bg-slate-700/50 rounded-xl transition-colors"
                >
                    {{ t('cancel') }}
                </button>
                <button 
                    @click="handleSubmit" 
                    class="px-5 py-2.5 bg-slate-900 dark:bg-primary-600 text-white font-medium hover:bg-slate-800 dark:hover:bg-primary-500 rounded-xl shadow-lg shadow-slate-900/20 dark:shadow-primary-900/20 transition-all transform hover:scale-105 active:scale-95"
                >
                    {{ isEditing ? t('save') : t('add') }}
                </button>
            </div>
        </div>
    </div>

    <!-- Add Project Modal -->
    <div v-if="showProjectModal" class="fixed inset-0 bg-slate-900/20 dark:bg-slate-900/50 backdrop-blur-sm flex items-center justify-center p-4 z-50 transition-all">
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-0 w-full max-w-sm shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800">
            <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/50 flex justify-between items-center">
                <h2 class="text-lg font-bold font-display text-slate-800 dark:text-white">{{ t('add_project') }}</h2>
                <button @click="showProjectModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200">
                    <span class="text-2xl leading-none">&times;</span>
                </button>
            </div>
            <div class="p-6 space-y-4">
                <div>
                    <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('name') }}</label>
                    <input 
                        v-model="projectForm.name" 
                        type="text" 
                        class="w-full px-4 py-3 rounded-xl bg-slate-100 dark:bg-slate-800 border-none focus:ring-2 focus:ring-primary-100 dark:focus:ring-primary-900 text-slate-800 dark:text-white placeholder-slate-400 font-medium" 
                        autofocus
                    />
                </div>
                <div>
                    <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">{{ t('color') }}</label>
                    <div class="flex gap-2 flex-wrap">
                      <button 
                        v-for="color in ['#EF4444', '#F59E0B', '#10B981', '#3B82F6', '#6366F1', '#8B5CF6', '#EC4899', '#64748B']"
                        :key="color"
                        @click="projectForm.color = color"
                        class="w-8 h-8 rounded-full border-2 transition-transform hover:scale-110"
                        :class="projectForm.color === color ? 'border-slate-900 dark:border-white' : 'border-transparent'"
                        :style="{ backgroundColor: color }"
                      ></button>
                    </div>
                </div>
            </div>
            <div class="px-6 py-4 bg-slate-50 dark:bg-slate-800/50 border-t border-slate-100 dark:border-slate-800 flex justify-end gap-3">
                <button @click="showProjectModal = false" class="px-4 py-2 text-slate-500 dark:text-slate-400 font-medium hover:bg-slate-200/50 rounded-lg">
                    {{ t('cancel') }}
                </button>
                <button @click="handleProjectSubmit" class="px-4 py-2 bg-primary-600 text-white font-medium hover:bg-primary-700 rounded-lg">
                    {{ t('create') }}
                </button>
            </div>
        </div>
    </div>

  </div>
</template>

<style scoped>
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.4s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.list-leave-active {
  position: absolute;
  width: 100%;
}
</style>
