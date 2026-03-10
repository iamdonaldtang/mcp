import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const mockPush = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush, replace: vi.fn() }),
  useRoute: () => ({ params: {}, query: {} }),
}))

vi.mock('@/api/client', () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
  cApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import B09CommunityEmpty from '@/views/b-community/B09CommunityEmpty.vue'

describe('B09CommunityEmpty', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
  })

  it('renders welcome section with title', () => {
    const wrapper = mount(B09CommunityEmpty)
    expect(wrapper.text()).toContain('Welcome to Community')
    expect(wrapper.text()).toContain('Build a thriving community')
  })

  it('renders 3 strategy cards', () => {
    const wrapper = mount(B09CommunityEmpty)
    expect(wrapper.text()).toContain('Activate New Users')
    expect(wrapper.text()).toContain('Drive Daily Engagement')
    expect(wrapper.text()).toContain('Maximize Retention')
  })

  it('has activate strategy selected by default', () => {
    const wrapper = mount(B09CommunityEmpty)
    // The activate card should show the expanded "Includes" section
    expect(wrapper.text()).toContain('Includes')
    expect(wrapper.text()).toContain('Sectors & Tasks')
    expect(wrapper.text()).toContain('Points & Level')
    expect(wrapper.text()).toContain('TaskChain')
    expect(wrapper.text()).toContain('Activation Rate')
  })

  it('selects a different strategy on click', async () => {
    const wrapper = mount(B09CommunityEmpty)
    // Click the retention card (3rd button in the strategy grid)
    const strategyButtons = wrapper.findAll('.grid.grid-cols-3 button')
    expect(strategyButtons.length).toBe(3)

    await strategyButtons[2].trigger('click')

    // Retention card should now show expanded details
    expect(wrapper.text()).toContain('Benefits Shop')
    expect(wrapper.text()).toContain('Milestones')
    expect(wrapper.text()).toContain('30-Day Retention')
  })

  it('navigates to wizard with strategy template on CTA click', async () => {
    const wrapper = mount(B09CommunityEmpty)
    const ctaButton = wrapper.find('button.bg-community')
    expect(ctaButton.exists()).toBe(true)
    expect(ctaButton.text()).toContain('Create Community with This Strategy')

    await ctaButton.trigger('click')
    expect(mockPush).toHaveBeenCalledWith('/b/community/wizard/step-1?template=activate')
  })

  it('navigates to wizard with blank template on "Start from Scratch"', async () => {
    const wrapper = mount(B09CommunityEmpty)
    const scratchButton = wrapper.findAll('button').find(b => b.text().includes('start from scratch'))
    expect(scratchButton).toBeTruthy()

    await scratchButton!.trigger('click')
    expect(mockPush).toHaveBeenCalledWith('/b/community/wizard/step-1?template=blank')
  })

  it('renders How It Works engine steps', () => {
    const wrapper = mount(B09CommunityEmpty)
    expect(wrapper.text()).toContain('How It Works')
    expect(wrapper.text()).toContain('Quest')
    expect(wrapper.text()).toContain('Activate')
    expect(wrapper.text()).toContain('Engage')
    expect(wrapper.text()).toContain('Retain')
  })
})
