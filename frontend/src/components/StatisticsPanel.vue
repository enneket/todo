<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTodoStore } from '../stores/todo'
import { Doughnut, Bar } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title } from 'chart.js'
import { PhChartPie, PhCheckCircle, PhListBullets, PhTrendUp } from '@phosphor-icons/vue'

ChartJS.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title)

const { t } = useI18n()
const todoStore = useTodoStore()

// --- Stats Logic ---

const totalTodos = computed(() => todoStore.todos.length)
const completedTodos = computed(() => todoStore.todos.filter(t => t.completed).length)
const activeTodos = computed(() => totalTodos.value - completedTodos.value)
const completionRate = computed(() => totalTodos.value > 0 ? Math.round((completedTodos.value / totalTodos.value) * 100) : 0)

// --- Chart Data ---

const doughnutData = computed(() => ({
  labels: [t('active'), t('completed')],
  datasets: [
    {
      backgroundColor: ['#F59E0B', '#10B981'],
      data: [activeTodos.value, completedTodos.value]
    }
  ]
}))

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom' as const,
      labels: {
          usePointStyle: true,
          font: { family: 'Archivo' }
      }
    }
  }
}

// Priority Distribution Bar Chart
const barData = computed(() => {
    const high = todoStore.todos.filter(t => t.priority === 'high').length
    const medium = todoStore.todos.filter(t => t.priority === 'medium').length
    const low = todoStore.todos.filter(t => t.priority === 'low').length
    return {
        labels: [t('priority_high'), t('priority_medium'), t('priority_low')],
        datasets: [{
            label: t('count'),
            data: [high, medium, low],
            backgroundColor: ['#EF4444', '#F59E0B', '#10B981'],
            borderRadius: 6
        }]
    }
})

const barOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { display: false }
    },
    scales: {
        y: { beginAtZero: true, ticks: { stepSize: 1 } },
        x: { grid: { display: false } }
    }
}

</script>

<template>
  <div class="bg-white dark:bg-slate-900 p-6 rounded-2xl shadow-sm border border-slate-100 dark:border-slate-800">
    <h2 class="text-xl font-bold font-display text-slate-800 dark:text-white mb-6 flex items-center gap-2">
        <PhChartPie weight="duotone" class="text-primary-500" />
        {{ t('statistics') }}
    </h2>

    <!-- KPI Cards -->
    <div class="grid grid-cols-3 gap-4 mb-8">
        <div class="p-4 rounded-xl bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800/50 text-center">
            <div class="text-blue-500 mb-1 flex justify-center"><PhListBullets size="24" weight="duotone"/></div>
            <div class="text-2xl font-bold text-slate-800 dark:text-white">{{ totalTodos }}</div>
            <div class="text-xs font-bold text-blue-600 dark:text-blue-400 uppercase tracking-wider">{{ t('total') }}</div>
        </div>
        <div class="p-4 rounded-xl bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-100 dark:border-emerald-800/50 text-center">
            <div class="text-emerald-500 mb-1 flex justify-center"><PhCheckCircle size="24" weight="duotone"/></div>
            <div class="text-2xl font-bold text-slate-800 dark:text-white">{{ completedTodos }}</div>
            <div class="text-xs font-bold text-emerald-600 dark:text-emerald-400 uppercase tracking-wider">{{ t('completed') }}</div>
        </div>
        <div class="p-4 rounded-xl bg-indigo-50 dark:bg-indigo-900/20 border border-indigo-100 dark:border-indigo-800/50 text-center">
            <div class="text-indigo-500 mb-1 flex justify-center"><PhTrendUp size="24" weight="duotone"/></div>
            <div class="text-2xl font-bold text-slate-800 dark:text-white">{{ completionRate }}%</div>
            <div class="text-xs font-bold text-indigo-600 dark:text-indigo-400 uppercase tracking-wider">{{ t('rate') }}</div>
        </div>
    </div>

    <!-- Charts Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- Status Chart -->
        <div class="h-48 relative">
            <Doughnut :data="doughnutData" :options="doughnutOptions" />
        </div>
        
        <!-- Priority Chart -->
        <div class="h-48 relative">
            <Bar :data="barData" :options="barOptions" />
        </div>
    </div>
  </div>
</template>
