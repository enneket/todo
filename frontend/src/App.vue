<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTodoStore } from './stores/todo'
import { PhPlus, PhTrash, PhCheckCircle, PhCircle, PhTranslate } from '@phosphor-icons/vue'

const { t, locale } = useI18n()
const todoStore = useTodoStore()
const newTodoTitle = ref('')
const filter = ref<'all' | 'active' | 'completed'>('all')

onMounted(() => {
  todoStore.fetchTodos()
})

const handleAddTodo = async () => {
  if (!newTodoTitle.value.trim()) return
  await todoStore.addTodo(newTodoTitle.value)
  newTodoTitle.value = ''
}

const filteredTodos = computed(() => {
  if (filter.value === 'active') return todoStore.todos.filter(t => !t.completed)
  if (filter.value === 'completed') return todoStore.todos.filter(t => t.completed)
  return todoStore.todos
})

const toggleLanguage = () => {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-800 dark:text-gray-200 p-8">
    <div class="max-w-2xl mx-auto bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
      <!-- Header -->
      <header class="p-6 bg-blue-600 text-white flex justify-between items-center">
        <h1 class="text-2xl font-bold flex items-center gap-2">
          <PhCheckCircle weight="bold" />
          {{ t('title') }}
        </h1>
        <button @click="toggleLanguage" class="p-2 hover:bg-blue-700 rounded-full transition-colors" :title="t('language')">
          <PhTranslate size="24" />
        </button>
      </header>

      <!-- Input -->
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex gap-2">
          <input
            v-model="newTodoTitle"
            @keyup.enter="handleAddTodo"
            type="text"
            :placeholder="t('placeholder')"
            class="flex-1 px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
          <button
            @click="handleAddTodo"
            class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors flex items-center gap-2"
          >
            <PhPlus weight="bold" />
            {{ t('add') }}
          </button>
        </div>
      </div>

      <!-- Filters -->
      <div class="flex border-b border-gray-200 dark:border-gray-700">
        <button
          v-for="f in ['all', 'active', 'completed']"
          :key="f"
          @click="filter = f as any"
          class="flex-1 py-3 font-medium transition-colors border-b-2"
          :class="filter === f ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700'"
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
          <button @click="todoStore.toggleTodo(todo)" class="text-gray-400 hover:text-blue-600 transition-colors">
            <PhCheckCircle v-if="todo.completed" weight="fill" class="text-green-500" size="24" />
            <PhCircle v-else size="24" />
          </button>
          
          <span
            class="flex-1 text-lg transition-all"
            :class="{ 'line-through text-gray-400': todo.completed }"
          >
            {{ todo.title }}
          </span>

          <button
            @click="todoStore.deleteTodo(todo.id)"
            class="text-gray-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-all p-2"
          >
            <PhTrash size="20" />
          </button>
        </div>
        
        <div v-if="filteredTodos.length === 0" class="p-8 text-center text-gray-400">
          No tasks found
        </div>
      </div>
    </div>
  </div>
</template>
