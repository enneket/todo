import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { useThemeStore } from './theme'

describe('Theme Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    document.documentElement.classList.remove('dark')
  })

  it('falls back to auto when no theme is stored', () => {
    const store = useThemeStore()
    expect(store.theme).toBe('auto')
  })

  it('falls back to auto when the stored value is invalid', () => {
    localStorage.setItem('theme', 'blue')
    const store = useThemeStore()
    expect(store.theme).toBe('auto')
  })

  it('applies and persists the dark theme', async () => {
    const store = useThemeStore()
    store.theme = 'dark'
    await nextTick()
    expect(document.documentElement.classList.contains('dark')).toBe(true)
    expect(localStorage.getItem('theme')).toBe('dark')
  })

  it('removes the dark class for the light theme', async () => {
    const store = useThemeStore()
    store.theme = 'dark'
    await nextTick()
    store.theme = 'light'
    await nextTick()
    expect(document.documentElement.classList.contains('dark')).toBe(false)
    expect(localStorage.getItem('theme')).toBe('light')
  })

  it('cycles light -> dark -> auto on toggle', () => {
    const store = useThemeStore()
    store.theme = 'light'
    store.toggleTheme()
    expect(store.theme).toBe('dark')
    store.toggleTheme()
    expect(store.theme).toBe('auto')
    store.toggleTheme()
    expect(store.theme).toBe('light')
  })

  it('exposes isDark for explicit themes', () => {
    const store = useThemeStore()
    store.theme = 'dark'
    expect(store.isDark).toBe(true)
    store.theme = 'light'
    expect(store.isDark).toBe(false)
  })
})
