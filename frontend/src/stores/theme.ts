import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export type Theme = 'light' | 'dark' | 'auto'

export const useThemeStore = defineStore('theme', () => {
  // Validate the persisted value so corrupted entries fall back to 'auto'.
  const storedTheme = localStorage.getItem('theme')
  const theme = ref<Theme>(
    storedTheme === 'light' || storedTheme === 'dark' || storedTheme === 'auto'
      ? storedTheme
      : 'auto'
  )

  // Reactive system preference so `isDark` tracks OS-level changes.
  const systemDark = ref(window.matchMedia('(prefers-color-scheme: dark)').matches)

  const updateDocumentClass = (dark: boolean) => {
    if (dark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  const applyTheme = () => {
    if (theme.value === 'auto') {
      updateDocumentClass(systemDark.value)
    } else {
      updateDocumentClass(theme.value === 'dark')
    }
  }

  // Watch for changes in theme state
  watch(theme, (newTheme) => {
    localStorage.setItem('theme', newTheme)
    applyTheme()
  })

  // Listen for system preference changes
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    systemDark.value = e.matches
    if (theme.value === 'auto') {
      updateDocumentClass(e.matches)
    }
  })

  // Initial application
  applyTheme()

  // Reactive dark-state for JS-styled surfaces (e.g. chart.js colors).
  const isDark = computed(() => (theme.value === 'auto' ? systemDark.value : theme.value === 'dark'))

  const toggleTheme = () => {
    if (theme.value === 'light') theme.value = 'dark'
    else if (theme.value === 'dark') theme.value = 'auto'
    else theme.value = 'light'
  }

  return {
    theme,
    isDark,
    toggleTheme,
    applyTheme
  }
})
