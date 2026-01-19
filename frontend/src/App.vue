<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTodoStore, type Todo } from './stores/todo'
import { PhPlus, PhTrash, PhCheckCircle, PhCircle, PhTranslate, PhPencil, PhClock } from '@phosphor-icons/vue'
import BaseSelect from './components/BaseSelect.vue'

const { t, locale } = useI18n()
const todoStore = useTodoStore()

const priorityOptions = computed(() => [
    { value: 'high', label: t('priority_high') },
    { value: 'medium', label: t('priority_medium') },
    { value: 'low', label: t('priority_low') }
])

const getDefaultDueDate = () => {
  const d = new Date()
  d.setDate(d.getDate() + 1) // Tomorrow
  d.setMinutes(d.getMinutes() - d.getTimezoneOffset()) // Adjust for local time ISO string
  return d.toISOString().slice(0, 16)
}

const newTodoTitle = ref('')
const newTodoDescription = ref('')
const newTodoPriority = ref('medium')
const newTodoDueDate = ref(getDefaultDueDate())
const newTodoTags = ref('')
const showAddModal = ref(false)

const filter = ref<'all' | 'active' | 'completed'>('all')
const searchQuery = ref('')

const editingTodo = ref<Todo | null>(null)
const editForm = ref({
  title: '',
  description: '',
  priority: 'medium',
  due_date: '',
  tags: ''
})

onMounted(() => {
  todoStore.fetchTodos()
})

const handleAddTodo = async () => {
  if (!newTodoTitle.value.trim()) return
  const dateToSend = newTodoDueDate.value ? new Date(newTodoDueDate.value).toISOString() : ''
  
  // Parse tags: comma separated
  const tags = newTodoTags.value.split(',').map(t => t.trim()).filter(t => t)

  await todoStore.addTodo(newTodoTitle.value, newTodoPriority.value, dateToSend, newTodoDescription.value, tags)
  newTodoTitle.value = ''
  newTodoDescription.value = ''
  newTodoPriority.value = 'medium'
  newTodoDueDate.value = getDefaultDueDate()
  newTodoTags.value = ''
  showAddModal.value = false
}

const startEdit = (todo: Todo) => {
  editingTodo.value = todo
  // ... (date conversion logic stays same)
  
  let dateStr = ''
  if (todo.due_date) {
      const d = new Date(todo.due_date)
      const offset = d.getTimezoneOffset() * 60000
      const localISOTime = (new Date(d.getTime() - offset)).toISOString().slice(0, 16)
      dateStr = localISOTime
  }

  editForm.value = {
    title: todo.title,
    description: todo.description || '',
    priority: todo.priority || 'medium',
    due_date: dateStr,
    tags: todo.tags ? todo.tags.join(', ') : ''
  }
}

const cancelEdit = () => {
  editingTodo.value = null
}

const saveEdit = async () => {
  if (!editingTodo.value) return
  if (!editForm.value.title.trim()) return

  const dateToSend = editForm.value.due_date ? new Date(editForm.value.due_date).toISOString() : ''
  const tags = editForm.value.tags.split(',').map(t => t.trim()).filter(t => t)
  
  await todoStore.updateTodo(editingTodo.value.id, {
    title: editForm.value.title,
    description: editForm.value.description,
    priority: editForm.value.priority as 'high' | 'medium' | 'low',
    due_date: dateToSend,
    tags: tags
  })
  editingTodo.value = null
}

const filteredTodos = computed(() => {
  let items = todoStore.todos
  
  // Search filter
  if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase()
      items = items.filter(t => 
        t.title.toLowerCase().includes(q) || 
        (t.description && t.description.toLowerCase().includes(q)) ||
        (t.tags && t.tags.some(tag => tag.toLowerCase().includes(q)))
      )
  }

  if (filter.value === 'active') return items.filter((t) => !t.completed)
  if (filter.value === 'completed') return items.filter((t) => t.completed)
  return items
})

const toggleLanguage = () => {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const getPriorityColor = (p: string) => {
  switch (p) {
    case 'high': return 'text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400'
    case 'medium': return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400'
    case 'low': return 'text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400'
    default: return 'text-gray-600 bg-gray-100'
  }
}

const formatDate = (dateStr: string | null) => {
    if (!dateStr) return ''
    return new Date(dateStr).toLocaleString()
}
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-800 dark:text-gray-200 p-8">
    <div class="max-w-3xl mx-auto bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
      <!-- Header -->
      <header class="p-6 bg-blue-600 text-white flex justify-between items-center">
        <div class="flex items-center gap-4 flex-1">
            <h1 class="text-2xl font-bold flex items-center gap-2 whitespace-nowrap">
            <PhCheckCircle weight="bold" />
            {{ t('title') }}
            </h1>
            <input
                v-model="searchQuery"
                type="text"
                :placeholder="t('search')"
                class="bg-blue-700/50 text-white placeholder-blue-200 border-none rounded-lg px-4 py-2 w-full max-w-xs focus:ring-2 focus:ring-white/50 transition-all"
            />
        </div>
        <button
          @click="toggleLanguage"
          class="p-2 hover:bg-blue-700 rounded-full transition-colors ml-4"
          :title="t('language')"
        >
          <PhTranslate size="24" />
        </button>
      </header>

      <!-- Add Button -->
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <button
            @click="showAddModal = true"
            class="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors flex items-center justify-center gap-2"
        >
            <PhPlus weight="bold" size="20" />
            {{ t('add') }}
        </button>
      </div>

      <!-- Filters -->
      <div class="flex border-b border-gray-200 dark:border-gray-700">
        <button
          v-for="f in ['all', 'active', 'completed']"
          :key="f"
          @click="filter = f as any"
          class="flex-1 py-3 font-medium transition-colors border-b-2"
          :class="
            filter === f
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700'
          "
        >
          {{ t(f) }}
        </button>
      </div>

      <!-- List -->
      <div class="divide-y divide-gray-200 dark:divide-gray-700 max-h-[60vh] overflow-y-auto">
        <div
          v-for="todo in filteredTodos"
          :key="todo.id"
          class="p-4 flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors group"
        >
          <button
            @click="todoStore.updateTodo(todo.id, { completed: !todo.completed })"
            class="text-gray-400 hover:text-blue-600 transition-colors"
          >
            <PhCheckCircle v-if="todo.completed" weight="fill" class="text-green-500" size="24" />
            <PhCircle v-else size="24" />
          </button>

          <div class="flex-1 min-w-0">
             <div class="flex items-center gap-2 mb-1">
                  <span :class="['text-xs px-2 py-0.5 rounded-full font-medium', getPriorityColor(todo.priority)]">
                      {{ t('priority_' + (todo.priority || 'medium')) }}
                  </span>
                  <span v-if="todo.due_date" class="text-xs text-gray-500 flex items-center gap-1">
                      <PhClock size="12" /> {{ formatDate(todo.due_date) }}
                  </span>
                  <div class="flex gap-1 ml-2">
                    <span v-for="tag in todo.tags" :key="tag" class="text-xs px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300">
                        #{{ tag }}
                    </span>
                  </div>
              </div>
            <span
              class="block text-lg transition-all truncate"
              :class="{ 'line-through text-gray-400': todo.completed }"
            >
              {{ todo.title }}
            </span>
            <p v-if="todo.description" class="text-sm text-gray-500 mt-1 line-clamp-2">
                {{ todo.description }}
            </p>
          </div>

          <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all">
               <button
                @click="startEdit(todo)"
                class="p-2 text-gray-400 hover:text-blue-500 transition-all"
                :title="t('edit')"
              >
                <PhPencil size="20" />
              </button>
              <button
                @click="todoStore.deleteTodo(todo.id)"
                class="p-2 text-gray-400 hover:text-red-500 transition-all"
              >
                <PhTrash size="20" />
              </button>
          </div>
        </div>

        <div v-if="filteredTodos.length === 0" class="p-8 text-center text-gray-400">
          No tasks found
        </div>
      </div>
    </div>
    
    <!-- Add Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
        <div class="bg-white dark:bg-gray-800 rounded-xl p-6 w-full max-w-md shadow-2xl">
            <h2 class="text-xl font-bold mb-4">{{ t('add') }}</h2>
            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium mb-1">Title</label>
                    <input v-model="newTodoTitle" type="text" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500" :placeholder="t('placeholder')" />
                </div>
                <div>
                    <label class="block text-sm font-medium mb-1">{{ t('description') }}</label>
                    <textarea v-model="newTodoDescription" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 h-24 resize-none" :placeholder="t('description')"></textarea>
                </div>
                <div class="flex gap-4">
                    <div class="flex-1">
                         <label class="block text-sm font-medium mb-1">{{ t('priority') }}</label>
                         <BaseSelect v-model="newTodoPriority" :options="priorityOptions" class="w-full" />
                    </div>
                    <div class="flex-1">
                         <label class="block text-sm font-medium mb-1">{{ t('due_date') }}</label>
                         <input type="datetime-local" v-model="newTodoDueDate" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700" />
                    </div>
                </div>
                <div>
                     <label class="block text-sm font-medium mb-1">{{ t('tags_label') }}</label>
                     <input v-model="newTodoTags" type="text" :placeholder="t('tags')" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>
                <div class="flex justify-end gap-2 mt-6">
                    <button @click="showAddModal = false" class="px-4 py-2 text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 rounded-lg transition-colors">{{ t('cancel') }}</button>
                    <button @click="handleAddTodo" class="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded-lg transition-colors">{{ t('add') }}</button>
                </div>
            </div>
        </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="editingTodo" class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
        <div class="bg-white dark:bg-gray-800 rounded-xl p-6 w-full max-w-md shadow-2xl">
            <h2 class="text-xl font-bold mb-4">{{ t('edit') }}</h2>
            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium mb-1">Title</label>
                    <input v-model="editForm.title" type="text" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>
                <div>
                    <label class="block text-sm font-medium mb-1">{{ t('description') }}</label>
                    <textarea v-model="editForm.description" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 h-24 resize-none"></textarea>
                </div>
                <div class="flex gap-4">
                    <div class="flex-1">
                         <label class="block text-sm font-medium mb-1">{{ t('priority') }}</label>
                         <BaseSelect v-model="editForm.priority" :options="priorityOptions" class="w-full" />
                    </div>
                    <div class="flex-1">
                         <label class="block text-sm font-medium mb-1">{{ t('due_date') }}</label>
                         <input type="datetime-local" v-model="editForm.due_date" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700" />
                    </div>
                </div>
                <div>
                     <label class="block text-sm font-medium mb-1">{{ t('tags_label') }}</label>
                     <input v-model="editForm.tags" type="text" :placeholder="t('tags')" class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>
                <div class="flex justify-end gap-2 mt-6">
                    <button @click="cancelEdit" class="px-4 py-2 text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 rounded-lg transition-colors">{{ t('cancel') }}</button>
                    <button @click="saveEdit" class="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded-lg transition-colors">{{ t('save') }}</button>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>
