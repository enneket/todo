import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseSelect from './BaseSelect.vue'

describe('BaseSelect', () => {
  const options = [
    { value: 'opt1', label: 'Option 1' },
    { value: 'opt2', label: 'Option 2' }
  ]

  it('renders correctly with initial value', () => {
    const wrapper = mount(BaseSelect, {
      props: {
        modelValue: 'opt1',
        options
      }
    })
    expect(wrapper.text()).toContain('Option 1')
  })

  it('toggles dropdown on click', async () => {
    const wrapper = mount(BaseSelect, {
      props: {
        modelValue: 'opt1',
        options
      }
    })
    
    // Initially closed (options not visible)
    expect(wrapper.find('div.absolute').exists()).toBe(false)

    // Click button
    await wrapper.find('button').trigger('click')
    
    // Now open
    expect(wrapper.find('div.absolute').exists()).toBe(true)
    expect(wrapper.text()).toContain('Option 2')
  })

  it('emits update event when option selected', async () => {
    const wrapper = mount(BaseSelect, {
      props: {
        modelValue: 'opt1',
        options
      }
    })

    await wrapper.find('button').trigger('click') // Open
    
    const optionsButtons = wrapper.findAll('div.absolute button')
    await optionsButtons[1].trigger('click') // Click Option 2

    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')![0]).toEqual(['opt2'])
    
    // Should close after selection
    expect(wrapper.find('div.absolute').exists()).toBe(false)
  })
})
