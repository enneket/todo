<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTodoStore } from '../stores/todo'
import { storeToRefs } from 'pinia'
import { PhCaretLeft, PhCaretRight, PhCheckCircle, PhCircle } from '@phosphor-icons/vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const todoStore = useTodoStore()
const { todos } = storeToRefs(todoStore)

const viewMode = ref<'month' | 'week'>('month')
const currentDate = ref(new Date())

// Helper to get start of week (Sunday)
const getStartOfWeek = (date: Date) => {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day
  return new Date(d.setDate(diff))
}

const currentMonthName = computed(() => {
  return currentDate.value.toLocaleString('default', { month: 'long', year: 'numeric' })
})

const weekDays = computed(() => {
  const days = []
  const start = getStartOfWeek(new Date()) // Use current date just for names
  for (let i = 0; i < 7; i++) {
    const d = new Date(start)
    d.setDate(start.getDate() + i)
    days.push(d.toLocaleDateString('default', { weekday: 'short' }))
  }
  return days
})

const calendarDays = computed(() => {
  const days = []
  if (viewMode.value === 'month') {
    const year = currentDate.value.getFullYear()
    const month = currentDate.value.getMonth()
    const firstDayOfMonth = new Date(year, month, 1)
    
    // Start from the Sunday before the first day of the month
    const startDate = new Date(firstDayOfMonth)
    startDate.setDate(firstDayOfMonth.getDate() - firstDayOfMonth.getDay())
    
    // 6 weeks * 7 days = 42 days grid
    for (let i = 0; i < 42; i++) {
      const d = new Date(startDate)
      d.setDate(startDate.getDate() + i)
      days.push({
        date: d,
        isCurrentMonth: d.getMonth() === month,
        isToday: isSameDate(d, new Date())
      })
    }
  } else {
    // Week view
    const start = getStartOfWeek(currentDate.value)
    for (let i = 0; i < 7; i++) {
      const d = new Date(start)
      d.setDate(start.getDate() + i)
      days.push({
        date: d,
        isCurrentMonth: true, // Always relevant in week view
        isToday: isSameDate(d, new Date())
      })
    }
  }
  return days
})

const isSameDate = (d1: Date, d2: Date) => {
  return d1.getDate() === d2.getDate() &&
         d1.getMonth() === d2.getMonth() &&
         d1.getFullYear() === d2.getFullYear()
}

const getTasksForDate = (date: Date) => {
  return todos.value.filter(todo => {
    if (!todo.due_date) return false
    const d = new Date(todo.due_date)
    return isSameDate(d, date)
  })
}

const prev = () => {
  if (viewMode.value === 'month') {
    currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1, 1)
  } else {
    currentDate.value = new Date(currentDate.value.setDate(currentDate.value.getDate() - 7))
  }
}

const next = () => {
  if (viewMode.value === 'month') {
    currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 1)
  } else {
    currentDate.value = new Date(currentDate.value.setDate(currentDate.value.getDate() + 7))
  }
}

const setToday = () => {
  currentDate.value = new Date()
}

const getPriorityColor = (priority: string) => {
    switch (priority) {
        case 'high': return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
        case 'medium': return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
        case 'low': return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
        default: return 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-400'
    }
}
</script>

<template>
  <div class="h-full flex flex-col bg-white dark:bg-slate-900 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-800 overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200 dark:border-slate-800">
      <div class="flex items-center gap-4">
        <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100 capitalize">
          {{ currentMonthName }}
        </h2>
        <div class="flex items-center bg-slate-100 dark:bg-slate-800 rounded-lg p-1">
          <button 
            @click="prev"
            class="p-1 hover:bg-white dark:hover:bg-slate-700 rounded-md transition-colors text-slate-600 dark:text-slate-400"
          >
            <PhCaretLeft size="20" weight="bold" />
          </button>
          <button 
            @click="setToday"
            class="px-3 py-1 text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-primary-600 dark:hover:text-primary-400 transition-colors"
          >
            {{ t('today') }}
          </button>
          <button 
            @click="next"
            class="p-1 hover:bg-white dark:hover:bg-slate-700 rounded-md transition-colors text-slate-600 dark:text-slate-400"
          >
            <PhCaretRight size="20" weight="bold" />
          </button>
        </div>
      </div>
      
      <div class="flex bg-slate-100 dark:bg-slate-800 rounded-lg p-1">
        <button 
          @click="viewMode = 'month'"
          :class="['px-3 py-1.5 text-sm font-medium rounded-md transition-all', viewMode === 'month' ? 'bg-white dark:bg-slate-700 text-primary-600 dark:text-primary-400 shadow-sm' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300']"
        >
          {{ t('view_month') }}
        </button>
        <button 
          @click="viewMode = 'week'"
          :class="['px-3 py-1.5 text-sm font-medium rounded-md transition-all', viewMode === 'week' ? 'bg-white dark:bg-slate-700 text-primary-600 dark:text-primary-400 shadow-sm' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300']"
        >
          {{ t('view_week') }}
        </button>
      </div>
    </div>

    <!-- Weekday Headers -->
    <div class="grid grid-cols-7 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/50">
      <div 
        v-for="day in weekDays" 
        :key="day"
        class="py-2 text-center text-xs font-semibold text-slate-500 uppercase tracking-wider"
      >
        {{ day }}
      </div>
    </div>

    <!-- Calendar Grid -->
    <div class="flex-1 grid grid-cols-7 grid-rows-6 overflow-y-auto">
      <div 
        v-for="(day, index) in calendarDays" 
        :key="index"
        :class="[
          'border-b border-r border-slate-200 dark:border-slate-800 p-2 min-h-[100px] flex flex-col gap-1 transition-colors hover:bg-slate-50 dark:hover:bg-slate-800/30',
          !day.isCurrentMonth && viewMode === 'month' ? 'bg-slate-50/50 dark:bg-slate-900/50 text-slate-400 dark:text-slate-600' : 'bg-white dark:bg-slate-900',
          day.isToday ? 'bg-primary-50/30 dark:bg-primary-900/10' : ''
        ]"
      >
        <div class="flex justify-between items-start">
          <span 
            :class="[
              'text-sm font-medium w-7 h-7 flex items-center justify-center rounded-full',
              day.isToday ? 'bg-primary-600 text-white' : 'text-slate-700 dark:text-slate-300'
            ]"
          >
            {{ day.date.getDate() }}
          </span>
        </div>
        
        <div class="flex-1 flex flex-col gap-1 overflow-y-auto max-h-[100px] custom-scrollbar">
          <div 
            v-for="task in getTasksForDate(day.date)" 
            :key="task.id"
            :class="['text-[10px] px-1.5 py-1 rounded truncate cursor-pointer border border-transparent hover:border-slate-300 dark:hover:border-slate-600 flex items-center gap-1', task.completed ? 'opacity-50 line-through' : '', getPriorityColor(task.priority || 'medium')]"
            @click="$emit('edit-task', task)"
          >
             <PhCheckCircle v-if="task.completed" size="10" weight="fill" />
             <PhCircle v-else size="10" />
             <span class="truncate">{{ task.title }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 20px;
}
.dark .custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #475569;
}
</style>
