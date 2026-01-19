import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type Theme = 'light' | 'dark' | 'auto'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref<Theme>((localStorage.getItem('theme') as Theme) || 'auto')

  const updateDocumentClass = (isDark: boolean) => {
    if (isDark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  const applyTheme = () => {
    if (theme.value === 'auto') {
      const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      updateDocumentClass(systemDark)
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
    if (theme.value === 'auto') {
      updateDocumentClass(e.matches)
    }
  })

  // Initial application
  applyTheme()

  const toggleTheme = () => {
    if (theme.value === 'light') theme.value = 'dark'
    else if (theme.value === 'dark') theme.value = 'auto'
    else theme.value = 'light'
  }

  return {
    theme,
    toggleTheme,
    applyTheme
  }
})
