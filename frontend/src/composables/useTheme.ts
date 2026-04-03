import { useColorMode } from '@vueuse/core'
import { watch } from 'vue'

export function useTheme() {
  const mode = useColorMode({
    attribute: 'class',
    modes: { dark: 'dark', light: 'light' },
  })

  watch(mode, (val) => {
    document.documentElement.classList.toggle('dark', val === 'dark')
  }, { immediate: true })

  function toggle() {
    mode.value = mode.value === 'dark' ? 'light' : 'dark'
  }

  return { mode, toggle }
}
