<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { PhCaretDown, PhCheck } from '@phosphor-icons/vue'

interface Option {
  value: string
  label: string
}

const props = defineProps<{
  modelValue: string
  options: Option[]
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const isOpen = ref(false)
const containerRef = ref<HTMLElement | null>(null)

const toggle = () => {
  isOpen.value = !isOpen.value
}

const select = (value: string) => {
  emit('update:modelValue', value)
  isOpen.value = false
}

const handleClickOutside = (event: MouseEvent) => {
  if (containerRef.value && !containerRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

const getLabel = (val: string) => {
    return props.options.find(o => o.value === val)?.label || val
}
</script>

<template>
  <div ref="containerRef" class="relative inline-block text-left min-w-[120px]">
    <button
      @click="toggle"
      type="button"
      class="inline-flex w-full justify-between items-center gap-x-1.5 rounded-xl bg-slate-100 px-4 py-2.5 text-sm font-semibold text-slate-800 shadow-sm hover:bg-slate-200 transition-colors h-[42px] focus:outline-none focus:ring-2 focus:ring-primary-100"
    >
      {{ getLabel(modelValue) }}
      <PhCaretDown class="-mr-1 h-4 w-4 text-slate-400" aria-hidden="true" />
    </button>

    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute right-0 z-50 mt-2 w-full origin-top-right rounded-xl bg-white shadow-xl ring-1 ring-black ring-opacity-5 focus:outline-none overflow-hidden"
      >
        <div class="py-1">
          <button
            v-for="opt in options"
            :key="opt.value"
            @click="select(opt.value)"
            class="group flex w-full items-center justify-between px-4 py-2.5 text-sm text-slate-600 hover:bg-slate-50 hover:text-primary-600 transition-colors"
          >
            {{ opt.label }}
            <PhCheck v-if="modelValue === opt.value" class="h-4 w-4 text-primary-500" />
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>
