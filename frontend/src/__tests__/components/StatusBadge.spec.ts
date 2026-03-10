import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusBadge from '@/components/common/StatusBadge.vue'

describe('StatusBadge', () => {
  it('renders active status with correct label and colors', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'active' } })
    expect(wrapper.text()).toBe('Active')
    const style = wrapper.find('span').attributes('style')
    expect(style).toContain('background-color: #0A2E1A')
    expect(style).toContain('color: #16A34A')
  })

  it('renders draft status with correct label and colors', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'draft' } })
    expect(wrapper.text()).toBe('Draft')
    const style = wrapper.find('span').attributes('style')
    expect(style).toContain('background-color: #1F1A08')
    expect(style).toContain('color: #D97706')
  })

  it('renders completed status with correct label and colors', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'completed' } })
    expect(wrapper.text()).toBe('Completed')
    const style = wrapper.find('span').attributes('style')
    expect(style).toContain('background-color: #1E293B')
    expect(style).toContain('color: #64748B')
  })

  it('renders paused status with correct label and colors', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'paused' } })
    expect(wrapper.text()).toBe('Paused')
    const style = wrapper.find('span').attributes('style')
    expect(style).toContain('background-color: #2D1515')
    expect(style).toContain('color: #DC2626')
  })

  it('renders as an inline span element', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'active' } })
    expect(wrapper.element.tagName).toBe('SPAN')
  })

  it('updates reactively when status prop changes', async () => {
    const wrapper = mount(StatusBadge, { props: { status: 'active' } })
    expect(wrapper.text()).toBe('Active')

    await wrapper.setProps({ status: 'paused' })
    expect(wrapper.text()).toBe('Paused')
    expect(wrapper.find('span').attributes('style')).toContain('color: #DC2626')
  })
})
