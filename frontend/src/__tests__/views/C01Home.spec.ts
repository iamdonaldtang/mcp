import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const { mockPush, mockCGet, mockCPost } = vi.hoisted(() => ({
  mockPush: vi.fn(),
  mockCGet: vi.fn(),
  mockCPost: vi.fn(),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush, replace: vi.fn() }),
  useRoute: () => ({ params: {}, query: {} }),
}))

vi.mock('@/api/client', () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
  cApi: { get: mockCGet, post: mockCPost, put: vi.fn(), delete: vi.fn() },
}))

vi.mock('@/stores/c-end', () => ({
  useCEndStore: () => ({
    isWalletConnected: true,
    userStatus: {
      walletAddress: '0xABCDEF1234567890ABCDEF1234567890ABCDEF12',
      level: 3,
      xp: 1500,
      dayStreak: 5,
      totalPoints: 1200,
      rank: 7,
    },
    tabVisibility: { home: true, quests: true, leaderboard: true },
  }),
}))

import C01Home from '@/views/c-end/C01Home.vue'

describe('C01Home', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
    mockCGet.mockReset()
    mockCPost.mockReset()
  })

  function setupDefaultMocks() {
    mockCGet.mockImplementation((url: string) => {
      if (url.includes('/home')) {
        return Promise.resolve({
          data: { data: { pulse: { totalMembers: 2500, thisWeek: 180, liveActive: 42, tasksDone: 8900 } } },
        })
      }
      if (url.includes('/tasks')) {
        return Promise.resolve({
          data: {
            data: [
              {
                id: 's1',
                name: 'Getting Started',
                tasks: [
                  { id: 't1', title: 'Follow on Twitter', type: 'social', points: 50, status: 'available' },
                  { id: 't2', title: 'Join Discord', type: 'social', points: 30, status: 'completed' },
                ],
              },
            ],
          },
        })
      }
      if (url.includes('/announcements')) {
        return Promise.resolve({ data: { data: [] } })
      }
      if (url.includes('/daychain')) {
        return Promise.resolve({
          data: { data: { enabled: true, currentStreak: 5, targetDays: 30, checkedInToday: false } },
        })
      }
      return Promise.resolve({ data: { data: null } })
    })
  }

  it('renders loading state initially', () => {
    mockCGet.mockReturnValue(new Promise(() => {})) // never resolves
    const wrapper = mount(C01Home, {
      global: {
        stubs: {
          TaskCard: { template: '<div class="task-card-stub" />', props: ['task'] },
        },
      },
    })
    expect(wrapper.find('.animate-pulse').exists()).toBe(true)
  })

  it('renders sectors and tasks after loading', async () => {
    setupDefaultMocks()
    const wrapper = mount(C01Home, {
      global: {
        stubs: {
          TaskCard: { template: '<div class="task-card-stub" />', props: ['task'] },
        },
      },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('Getting Started')
    expect(wrapper.findAll('.task-card-stub').length).toBe(2)
  })

  it('shows user card when wallet is connected', async () => {
    setupDefaultMocks()
    const wrapper = mount(C01Home, {
      global: {
        stubs: {
          TaskCard: { template: '<div class="task-card-stub" />', props: ['task'] },
        },
      },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('0xABCD')
    expect(wrapper.text()).toContain('1,500 XP')
    expect(wrapper.text()).toContain('Level 3')
    expect(wrapper.text()).toContain('5 day streak')
  })

  it('shows DayChain check-in when enabled', async () => {
    setupDefaultMocks()
    const wrapper = mount(C01Home, {
      global: {
        stubs: {
          TaskCard: { template: '<div class="task-card-stub" />', props: ['task'] },
        },
      },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('Daily Streak')
    expect(wrapper.text()).toContain('Day 5 of 30')
    expect(wrapper.text()).toContain('Continue Streak')
  })

  it('handles check-in click', async () => {
    setupDefaultMocks()
    mockCPost.mockResolvedValue({ data: { code: 0 } })

    const wrapper = mount(C01Home, {
      global: {
        stubs: {
          TaskCard: { template: '<div class="task-card-stub" />', props: ['task'] },
        },
      },
    })
    await flushPromises()

    const checkInButton = wrapper.findAll('button').find(b => b.text().includes('Continue Streak'))
    expect(checkInButton).toBeTruthy()

    await checkInButton!.trigger('click')
    await flushPromises()

    expect(mockCPost).toHaveBeenCalledWith('/api/c/community/daychain')
    // After check-in, button should show checked state
    expect(wrapper.text()).toContain('Checked In')
  })

  it('renders community pulse stats', async () => {
    setupDefaultMocks()
    const wrapper = mount(C01Home, {
      global: {
        stubs: {
          TaskCard: { template: '<div class="task-card-stub" />', props: ['task'] },
        },
      },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('2,500')
    expect(wrapper.text()).toContain('Members')
    expect(wrapper.text()).toContain('180')
    expect(wrapper.text()).toContain('This Week')
    expect(wrapper.text()).toContain('8,900')
    expect(wrapper.text()).toContain('Tasks Done')
  })
})
