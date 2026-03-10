import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TaskCard from '@/components/c-end/TaskCard.vue'
import type { CTask } from '@/types/c-end'

function makeTask(overrides: Partial<CTask> = {}): CTask {
  return {
    id: 'task-1',
    name: 'Follow on Twitter',
    type: 'social',
    status: 'available',
    points: 50,
    icon: 'share',
    iconColor: '#5D7EF1',
    sectorName: 'Social',
    completions: 1234,
    ...overrides,
  }
}

describe('TaskCard', () => {
  it('renders task name', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask() } })
    expect(wrapper.text()).toContain('Follow on Twitter')
  })

  it('renders task type and formatted completions', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask() } })
    expect(wrapper.text()).toContain('social')
    expect(wrapper.text()).toContain('1,234 completed')
  })

  it('renders points with XP suffix', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ points: 100 }) } })
    expect(wrapper.text()).toContain('+100 XP')
  })

  it('renders task icon with correct color', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask() } })
    const icons = wrapper.findAll('.material-symbols-rounded')
    const taskIcon = icons.find(i => i.text() === 'share')!
    expect(taskIcon.attributes('style')).toContain('color: #5D7EF1')
  })

  // --- Status-specific button labels ---

  it('shows "Start" button for available status', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'available' }) } })
    const btn = wrapper.find('button')
    expect(btn.text()).toBe('Start')
  })

  it('shows "Continue" button for in_progress status', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'in_progress' }) } })
    const btn = wrapper.find('button')
    expect(btn.text()).toBe('Continue')
  })

  it('shows "Claim" button for completed status', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'completed' }) } })
    const btn = wrapper.find('button')
    expect(btn.text()).toBe('Claim')
  })

  it('shows "Done ✓" button for claimed status', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'claimed' }) } })
    const btn = wrapper.find('button')
    expect(btn.text()).toContain('Done')
  })

  it('shows "Locked" button for locked status', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'locked' }) } })
    const btn = wrapper.find('button')
    expect(btn.text()).toBe('Locked')
  })

  it('shows "Expired" button for expired status', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'expired' }) } })
    const btn = wrapper.find('button')
    expect(btn.text()).toBe('Expired')
  })

  it('shows "Cooldown" button for cooldown status', () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'cooldown' }) } })
    const btn = wrapper.find('button')
    expect(btn.text()).toBe('Cooldown')
  })

  // --- Action emission ---

  it('emits action with task id when available task button is clicked', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ id: 'abc', status: 'available' }) } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toEqual([['abc']])
  })

  it('emits action for in_progress status', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'in_progress' }) } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toEqual([['task-1']])
  })

  it('emits action for completed status', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'completed' }) } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toEqual([['task-1']])
  })

  // --- Disabled states (no emission) ---

  it('does not emit action for claimed status', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'claimed' }) } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toBeUndefined()
  })

  it('does not emit action for locked status', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'locked' }) } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toBeUndefined()
  })

  it('does not emit action for expired status', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'expired' }) } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toBeUndefined()
  })

  it('does not emit action for cooldown status', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'cooldown' }) } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toBeUndefined()
  })

  // --- Disabled attribute ---

  it('sets disabled attribute on button for disabled statuses', () => {
    const disabledStatuses = ['claimed', 'locked', 'expired', 'cooldown'] as const
    for (const status of disabledStatuses) {
      const wrapper = mount(TaskCard, { props: { task: makeTask({ status }) } })
      expect(wrapper.find('button').attributes('disabled')).toBeDefined()
    }
  })

  it('does not set disabled attribute for active statuses', () => {
    const activeStatuses = ['available', 'in_progress', 'completed'] as const
    for (const status of activeStatuses) {
      const wrapper = mount(TaskCard, { props: { task: makeTask({ status }) } })
      expect(wrapper.find('button').attributes('disabled')).toBeUndefined()
    }
  })

  // --- Reactivity ---

  it('updates button label when task status changes', async () => {
    const wrapper = mount(TaskCard, { props: { task: makeTask({ status: 'available' }) } })
    expect(wrapper.find('button').text()).toBe('Start')

    await wrapper.setProps({ task: makeTask({ status: 'in_progress' }) })
    expect(wrapper.find('button').text()).toBe('Continue')
  })
})
