import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const { mockPush, mockReplace, mockGet } = vi.hoisted(() => ({
  mockPush: vi.fn(),
  mockReplace: vi.fn(),
  mockGet: vi.fn(),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush, replace: mockReplace }),
  useRoute: () => ({ params: {}, query: {} }),
}))

vi.mock('@/api/client', () => ({
  api: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn() },
  cApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import B01Dashboard from '@/views/b-community/B01Dashboard.vue'

describe('B01Dashboard', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
    mockGet.mockReset()
  })

  it('renders loading state initially', () => {
    mockGet.mockReturnValue(new Promise(() => {})) // never resolves
    const wrapper = mount(B01Dashboard)
    expect(wrapper.find('.animate-spin').exists()).toBe(true)
  })

  it('renders new user state with welcome message and product cards', async () => {
    mockGet.mockRejectedValue(new Error('no dashboard'))
    const wrapper = mount(B01Dashboard)
    await flushPromises()

    expect(wrapper.text()).toContain('Welcome to TaskOn')
    expect(wrapper.text()).toContain('Get Users Fast')
    expect(wrapper.text()).toContain('Make Users Stay')
    expect(wrapper.text()).toContain('Own the Experience')
  })

  it('renders active user state with stats cards', async () => {
    mockGet.mockResolvedValue({
      data: {
        state: 'active',
        stats: {
          totalUsers: 12500,
          activeCampaigns: 3,
          pointsDistributed: 45000,
          conversionRate: 12.5,
          totalUsersTrend: 8,
          activeCampaignsTrend: -2,
          pointsDistributedTrend: 15,
          conversionRateTrend: 3,
        },
        campaigns: [
          { id: '1', name: 'Summer Quest', product: 'quest', status: 'active', metric: 'Participants', metricValue: 500, progress: 65 },
        ],
        activities: [],
        productBreakdown: {
          quest: { campaigns: 5 },
          community: { members: 1200 },
          whitelabel: { instances: 1 },
          boost: { conversions: 340 },
        },
      },
    })
    const wrapper = mount(B01Dashboard)
    await flushPromises()

    expect(wrapper.text()).toContain('Dashboard')
    expect(wrapper.text()).toContain('Total Users')
    expect(wrapper.text()).toContain('12.5K')
    expect(wrapper.text()).toContain('Active Campaigns')
    expect(wrapper.text()).toContain('Summer Quest')
  })

  it('renders power user state with advanced metrics', async () => {
    mockGet.mockResolvedValue({
      data: {
        state: 'power',
        stats: {
          totalUsers: 250000,
          activeCampaigns: 12,
          pointsDistributed: 1500000,
          conversionRate: 18.3,
          totalUsersTrend: 12,
          activeCampaignsTrend: 5,
          pointsDistributedTrend: 20,
          conversionRateTrend: 7,
        },
        campaigns: [],
        activities: [
          { id: 'a1', icon: 'person_add', description: '50 new users joined', timeAgo: '2h ago', product: 'community' },
        ],
        productBreakdown: {
          quest: { campaigns: 20 },
          community: { members: 15000 },
          whitelabel: { instances: 3 },
          boost: { conversions: 5000 },
        },
      },
    })
    const wrapper = mount(B01Dashboard)
    await flushPromises()

    expect(wrapper.text()).toContain('Your growth engine overview')
    expect(wrapper.text()).toContain('250.0K')
    expect(wrapper.text()).toContain('User Growth')
    expect(wrapper.text()).toContain('Product Breakdown')
    expect(wrapper.text()).toContain('50 new users joined')
  })

  it('navigates to community when clicking product card', async () => {
    mockGet.mockRejectedValue(new Error('no dashboard'))
    const wrapper = mount(B01Dashboard)
    await flushPromises()

    // Goal cards are buttons in the grid
    const goalCards = wrapper.findAll('.grid.grid-cols-3 button')
    expect(goalCards.length).toBe(3)

    // Click "Make Users Stay" (community card, index 1)
    await goalCards[1].trigger('click')
    expect(mockPush).toHaveBeenCalledWith('/b/community')
  })

  it('falls back to new state on API error', async () => {
    mockGet.mockRejectedValue(new Error('Network error'))
    const wrapper = mount(B01Dashboard)
    await flushPromises()

    // Should show new user state (default)
    expect(wrapper.text()).toContain('Welcome to TaskOn')
    expect(wrapper.find('.animate-spin').exists()).toBe(false)
  })
})
