<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTodoStore, type Todo } from './stores/todo'
import { PhPlus, PhTrash, PhCheckCircle, PhCircle, PhTranslate, PhPencil, PhClock, PhMagnifyingGlass } from '@phosphor-icons/vue'
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

const form = ref({
  id: undefined as number | undefined,
  title: '',
  description: '',
  priority: 'medium',
  due_date: getDefaultDueDate(),
  tags: ''
})

const showModal = ref(false)
const isEditing = computed(() => !!form.value.id)

const filter = ref<'all' | 'active' | 'completed'>('all')
const searchQuery = ref('')

onMounted(() => {
  todoStore.fetchTodos()
})

const openAddModal = () => {
  form.value = {
    id: undefined,
    title: '',
    description: '',
    priority: 'medium',
    due_date: getDefaultDueDate(),
    tags: ''
  }
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

  form.value = {
    id: todo.id,
    title: todo.title,
    description: todo.description || '',
    priority: todo.priority || 'medium',
    due_date: dateStr,
    tags: todo.tags ? todo.tags.join(', ') : ''
  }
  showModal.value = true
}

const handleSubmit = async () => {
  if (!form.value.title.trim()) return

  const dateToSend = form.value.due_date ? new Date(form.value.due_date).toISOString() : ''
  const tags = form.value.tags.split(',').map(t => t.trim()).filter(t => t)

  if (isEditing.value && form.value.id) {
      await todoStore.updateTodo(form.value.id, {
        title: form.value.title,
        description: form.value.description,
        priority: form.value.priority as 'high' | 'medium' | 'low',
        due_date: dateToSend,
        tags: tags
      })
  } else {
      await todoStore.addTodo(
        form.value.title, 
        form.value.priority, 
        dateToSend, 
        form.value.description, 
        tags
      )
  }
  showModal.value = false
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
    case 'high': return 'text-red-600 bg-red-50 border border-red-100'
    case 'medium': return 'text-amber-600 bg-amber-50 border border-amber-100'
    case 'low': return 'text-emerald-600 bg-emerald-50 border border-emerald-100'
    default: return 'text-slate-600 bg-slate-50 border border-slate-100'
  }
}

const formatDate = (dateStr: string | null) => {
    if (!dateStr) return ''
    return new Date(dateStr).toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
        month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    })
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 text-slate-800 p-6 md:p-12 font-sans selection:bg-primary-100 selection:text-primary-700">
    <div class="max-w-2xl mx-auto space-y-8">
      
      <!-- Minimal Header -->
      <header class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
            <h1 class="text-4xl font-bold font-display tracking-tight text-slate-900 flex items-center gap-3">
                <PhCheckCircle weight="fill" class="text-primary-500" />
                {{ t('title') }}
            </h1>
            <p class="text-slate-500 mt-1 font-medium">{{ new Date().toLocaleDateString(locale === 'zh' ? 'zh-CN' : 'en-US', { weekday: 'long', month: 'long', day: 'numeric' }) }}</p>
        </div>
        
        <div class="flex items-center gap-3">
            <div class="relative group">
                <PhMagnifyingGlass class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-primary-500 transition-colors" size="18" />
                <input
                    v-model="searchQuery"
                    type="text"
                    :placeholder="t('search')"
                    class="pl-10 pr-4 py-2 bg-white border border-slate-200 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-primary-100 focus:border-primary-500 w-40 focus:w-64 transition-all shadow-sm"
                />
            </div>
            <button
              @click="toggleLanguage"
              class="p-2 text-slate-400 hover:text-slate-700 hover:bg-white hover:shadow-sm rounded-full transition-all"
              :title="t('language')"
            >
              <PhTranslate size="24" />
            </button>
        </div>
      </header>

      <!-- Add Button (Hero) -->
      <button
        @click="openAddModal"
        class="w-full py-4 bg-white hover:bg-primary-50/50 border border-slate-200 hover:border-primary-200 rounded-2xl shadow-sm hover:shadow-md transition-all group flex items-center justify-center gap-3"
      >
        <div class="bg-primary-50 text-primary-600 p-1.5 rounded-lg group-hover:scale-110 transition-transform">
            <PhPlus weight="bold" size="20" />
        </div>
        <span class="text-slate-600 font-medium text-lg group-hover:text-primary-700 transition-colors">{{ t('add') }}</span>
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
              ? 'bg-slate-800 text-white shadow-md shadow-slate-200'
              : 'bg-white text-slate-500 hover:bg-slate-100 shadow-sm border border-slate-100'
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
            class="bg-white p-5 rounded-2xl shadow-sm border border-slate-100 hover:shadow-md hover:border-primary-100 transition-all group relative overflow-hidden"
            >
            <!-- Completion Check -->
            <div class="flex items-start gap-4">
                <button
                    @click="todoStore.updateTodo(todo.id, { completed: !todo.completed })"
                    class="mt-1 text-slate-300 hover:text-primary-500 transition-colors flex-shrink-0"
                >
                    <PhCheckCircle v-if="todo.completed" weight="fill" class="text-emerald-500 scale-110" size="24" />
                    <PhCircle v-else weight="bold" size="24" />
                </button>

                <div class="flex-1 min-w-0">
                    <div class="flex flex-wrap items-center gap-2 mb-1.5">
                        <span :class="['text-[10px] uppercase tracking-wider font-bold px-2 py-0.5 rounded-md', getPriorityColor(todo.priority)]">
                            {{ t('priority_' + (todo.priority || 'medium')) }}
                        </span>
                        <span v-if="todo.due_date" class="text-xs text-slate-400 flex items-center gap-1 font-medium bg-slate-50 px-2 py-0.5 rounded-md">
                            <PhClock size="12" weight="bold" /> {{ formatDate(todo.due_date) }}
                        </span>
                        <div class="flex gap-1">
                            <span v-for="tag in todo.tags" :key="tag" class="text-xs px-2 py-0.5 rounded-md bg-blue-50 text-blue-600 font-medium">
                                #{{ tag }}
                            </span>
                        </div>
                    </div>
                    
                    <h3
                        class="text-lg font-medium text-slate-800 transition-all leading-snug"
                        :class="{ 'line-through text-slate-400': todo.completed }"
                    >
                        {{ todo.title }}
                    </h3>
                    
                    <p v-if="todo.description" class="text-sm text-slate-500 mt-1.5 leading-relaxed">
                        {{ todo.description }}
                    </p>
                </div>

                <!-- Actions (Hover) -->
                <div class="flex flex-col gap-1 opacity-0 group-hover:opacity-100 transition-opacity absolute right-4 top-4 bg-white pl-2">
                    <button
                        @click="openEditModal(todo)"
                        class="p-2 text-slate-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-all"
                        :title="t('edit')"
                    >
                        <PhPencil size="18" weight="bold" />
                    </button>
                    <button
                        @click="todoStore.deleteTodo(todo.id)"
                        class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                    >
                        <PhTrash size="18" weight="bold" />
                    </button>
                </div>
            </div>
            </div>
        </transition-group>

        <div v-if="filteredTodos.length === 0" class="py-12 text-center">
            <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-100 text-slate-300 mb-4">
                <PhCheckCircle size="32" weight="duotone" />
            </div>
            <p class="text-slate-400 font-medium">{{ t('no_tasks') || 'No tasks found' }}</p>
        </div>
      </div>
    </div>
    
    <!-- Add/Edit Modal (Backdrop blur) -->
    <div v-if="showModal" class="fixed inset-0 bg-slate-900/20 backdrop-blur-sm flex items-center justify-center p-4 z-50 transition-all">
        <div class="bg-white rounded-2xl p-0 w-full max-w-lg shadow-2xl overflow-hidden border border-slate-100 transform transition-all scale-100">
            <!-- Modal Header -->
            <div class="px-6 py-4 border-b border-slate-100 bg-slate-50/50 flex justify-between items-center">
                <h2 class="text-lg font-bold font-display text-slate-800">{{ isEditing ? t('edit') : t('add') }}</h2>
                <button @click="showModal = false" class="text-slate-400 hover:text-slate-600">
                    <span class="text-2xl leading-none">&times;</span>
                </button>
            </div>

            <!-- Modal Body -->
            <div class="p-6 space-y-5">
                <!-- Title -->
                <div>
                    <label class="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1.5">{{ t('modal_title') }}</label>
                    <input 
                        v-model="form.title" 
                        type="text" 
                        class="w-full px-4 py-3 rounded-xl bg-slate-100 border-none focus:ring-2 focus:ring-primary-100 text-slate-800 placeholder-slate-400 font-medium transition-shadow" 
                        :placeholder="t('placeholder')" 
                        autofocus
                    />
                </div>
                
                <!-- Description -->
                <div>
                    <label class="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1.5">{{ t('description') }}</label>
                    <textarea 
                        v-model="form.description" 
                        class="w-full px-4 py-3 rounded-xl bg-slate-100 border-none focus:ring-2 focus:ring-primary-100 text-slate-800 placeholder-slate-400 resize-none h-24 transition-shadow" 
                        :placeholder="t('description')"
                    ></textarea>
                </div>
                
                <!-- Metadata Row -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                         <label class="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1.5">{{ t('priority') }}</label>
                         <BaseSelect v-model="form.priority" :options="priorityOptions" class="w-full" />
                    </div>
                    <div>
                         <label class="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1.5">{{ t('due_date') }}</label>
                         <input 
                            type="datetime-local" 
                            v-model="form.due_date" 
                            class="w-full px-4 py-2.5 rounded-xl bg-slate-100 border-none focus:ring-2 focus:ring-primary-100 text-slate-800 text-sm h-[42px]" 
                        />
                    </div>
                </div>
                
                <!-- Tags -->
                <div>
                     <label class="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1.5">{{ t('tags_label') }}</label>
                     <input 
                        v-model="form.tags" 
                        type="text" 
                        :placeholder="t('tags')" 
                        class="w-full px-4 py-3 rounded-xl bg-slate-100 border-none focus:ring-2 focus:ring-primary-100 text-slate-800 placeholder-slate-400 transition-shadow" 
                    />
                </div>
            </div>

            <!-- Modal Footer -->
            <div class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3">
                <button 
                    @click="showModal = false" 
                    class="px-5 py-2.5 text-slate-500 font-medium hover:bg-slate-200/50 rounded-xl transition-colors"
                >
                    {{ t('cancel') }}
                </button>
                <button 
                    @click="handleSubmit" 
                    class="px-5 py-2.5 bg-slate-900 text-white font-medium hover:bg-slate-800 rounded-xl shadow-lg shadow-slate-900/20 transition-all transform hover:scale-105 active:scale-95"
                >
                    {{ isEditing ? t('save') : t('add') }}
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
