import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const { mockPush, mockGet } = vi.hoisted(() => ({
  mockPush: vi.fn(),
  mockGet: vi.fn(),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush, replace: vi.fn() }),
  useRoute: () => ({ params: {}, query: {} }),
}))

vi.mock('@/api/client', () => ({
  api: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn() },
  cApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import B14WLEmpty from '@/views/b-whitelabel/B14WLEmpty.vue'

describe('B14WLEmpty', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
    mockGet.mockReset()
  })

  it('renders welcome section with title', async () => {
    mockGet.mockRejectedValue(new Error('no community'))
    const wrapper = mount(B14WLEmpty)
    await flushPromises()

    expect(wrapper.text()).toContain('Own Your Growth Experience')
    expect(wrapper.text()).toContain('fully branded community platform')
  })

  it('renders 3 deployment path cards', async () => {
    mockGet.mockResolvedValue({ data: { data: { status: 'active' } } })
    const wrapper = mount(B14WLEmpty)
    await flushPromises()

    expect(wrapper.text()).toContain('Host on Your Domain')
    expect(wrapper.text()).toContain('Embed in Your App')
    expect(wrapper.text()).toContain('Build with SDK')
  })

  it('shows recommended badge on Embed path', async () => {
    mockGet.mockResolvedValue({ data: { data: { status: 'active' } } })
    const wrapper = mount(B14WLEmpty)
    await flushPromises()

    expect(wrapper.text()).toContain('Recommended')
  })

  it('shows prerequisite banner when no community exists', async () => {
    mockGet.mockRejectedValue(new Error('no community'))
    const wrapper = mount(B14WLEmpty)
    await flushPromises()

    expect(wrapper.text()).toContain('White Label requires an active Community')
    expect(wrapper.text()).toContain('Set Up Community')
  })

  it('hides prerequisite banner when community exists', async () => {
    mockGet.mockResolvedValue({ data: { data: { status: 'active' } } })
    const wrapper = mount(B14WLEmpty)
    await flushPromises()

    expect(wrapper.text()).not.toContain('White Label requires an active Community')
  })

  it('navigates to wizard on path card click', async () => {
    mockGet.mockResolvedValue({ data: { data: { status: 'active' } } })
    const wrapper = mount(B14WLEmpty)
    await flushPromises()

    // Click the first deployment path card (domain)
    const pathCards = wrapper.findAll('.grid.grid-cols-3 button')
    expect(pathCards.length).toBe(3)

    await pathCards[0].trigger('click')
    expect(mockPush).toHaveBeenCalledWith({
      path: '/b/whitelabel/wizard/step1',
      query: { path: 'domain' },
    })
  })

  it('renders highlight strip with key features', async () => {
    mockGet.mockResolvedValue({ data: { data: { status: 'active' } } })
    const wrapper = mount(B14WLEmpty)
    await flushPromises()

    expect(wrapper.text()).toContain('Custom Branding')
    expect(wrapper.text()).toContain('Data Ownership')
    expect(wrapper.text()).toContain('< 7 Days Integration')
  })
})
