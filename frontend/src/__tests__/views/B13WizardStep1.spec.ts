import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const { mockPush, mockGet, mockPost } = vi.hoisted(() => ({
  mockPush: vi.fn(),
  mockGet: vi.fn(),
  mockPost: vi.fn(),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush, replace: vi.fn() }),
  useRoute: () => ({ params: {}, query: {} }),
}))

vi.mock('@/api/client', () => ({
  api: { get: mockGet, post: mockPost, put: vi.fn(), delete: vi.fn() },
  cApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import B13WizardStep1 from '@/views/b-community/B13WizardStep1.vue'

describe('B13WizardStep1', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
    mockGet.mockReset()
    mockPost.mockReset()
    // Draft load returns nothing (fresh wizard)
    mockGet.mockRejectedValue(new Error('no draft'))
    mockPost.mockResolvedValue({ data: { code: 0 } })
  })

  it('renders step 1 form with name, description, color picker', async () => {
    const wrapper = mount(B13WizardStep1)
    await flushPromises()

    expect(wrapper.text()).toContain('Community Name')
    expect(wrapper.text()).toContain('Description')
    expect(wrapper.text()).toContain('Brand Color')
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('textarea').exists()).toBe(true)
  })

  it('shows stepper with 4 steps, step 1 active', async () => {
    const wrapper = mount(B13WizardStep1)
    await flushPromises()

    expect(wrapper.text()).toContain('Customize')
    expect(wrapper.text()).toContain('Modules')
    expect(wrapper.text()).toContain('Quick Setup')
    expect(wrapper.text()).toContain('Preview & Publish')

    // Step 1 dot should have the active class (bg-[#48BB78] text-white)
    const stepDots = wrapper.findAll('.w-8.h-8.rounded-full')
    expect(stepDots.length).toBe(4)
    expect(stepDots[0].classes()).toContain('bg-[#48BB78]')
  })

  it('validates name length (too short shows error)', async () => {
    const wrapper = mount(B13WizardStep1)
    await flushPromises()

    const nameInput = wrapper.find('input[type="text"]')
    await nameInput.setValue('AB')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Name must be at least 3 characters')
  })

  it('validates description length', async () => {
    const wrapper = mount(B13WizardStep1)
    await flushPromises()

    const textarea = wrapper.find('textarea')
    await textarea.setValue('Short')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Description must be at least 10 characters')
  })

  it('disables Next button when form is invalid', async () => {
    const wrapper = mount(B13WizardStep1)
    await flushPromises()

    // Find the Next button by its text content
    const nextButton = wrapper.findAll('button').find(b => b.text().includes('Next'))
    expect(nextButton).toBeTruthy()
    expect(nextButton!.attributes('disabled')).toBeDefined()
  })

  it('enables Next button when form is valid', async () => {
    const wrapper = mount(B13WizardStep1)
    await flushPromises()

    const nameInput = wrapper.find('input[type="text"]')
    const textarea = wrapper.find('textarea')
    await nameInput.setValue('My Test Community')
    await textarea.setValue('This is a valid description for the community that is long enough.')
    await wrapper.vm.$nextTick()

    const nextButton = wrapper.findAll('button').find(b => b.text().includes('Next'))
    expect(nextButton).toBeTruthy()
    expect(nextButton!.attributes('disabled')).toBeUndefined()
  })

  it('allows selecting preset color', async () => {
    const wrapper = mount(B13WizardStep1)
    await flushPromises()

    // Color preset buttons are w-9 h-9 rounded-lg
    const colorButtons = wrapper.findAll('button.w-9.h-9.rounded-lg')
    expect(colorButtons.length).toBeGreaterThanOrEqual(8) // 8 presets + custom trigger

    // Click the second preset (#5D7EF1)
    await colorButtons[1].trigger('click')

    // The live preview color bar should reflect the change
    const colorBar = wrapper.find('.h-2.transition-colors')
    expect(colorBar.attributes('style')).toContain('#5D7EF1')
  })
})
