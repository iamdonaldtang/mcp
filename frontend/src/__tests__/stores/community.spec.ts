import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('@/api/client', () => import('../mocks/api'))

import { useCommunityStore } from '@/stores/community'
import { api } from '@/api/client'

describe('community store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('has null overview and insights, empty checklist', () => {
      const store = useCommunityStore()
      expect(store.overview).toBeNull()
      expect(store.insights).toBeNull()
      expect(store.checklist).toEqual([])
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchOverview', () => {
    const mockOverview = {
      totalMembers: 1200,
      activeMembers: 350,
      totalPoints: 50000,
      modulesEnabled: 5,
    }

    it('fetches and sets overview data', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockOverview, message: 'ok' },
      })

      const store = useCommunityStore()
      await store.fetchOverview()

      expect(api.get).toHaveBeenCalledWith('/api/v1/community/overview')
      expect(store.overview).toEqual(mockOverview)
    })

    it('sets loading true during fetch and false after', async () => {
      let resolveFn: (value: any) => void
      const pending = new Promise((resolve) => { resolveFn = resolve })
      vi.mocked(api.get).mockReturnValueOnce(pending as any)

      const store = useCommunityStore()
      const promise = store.fetchOverview()

      expect(store.loading).toBe(true)

      resolveFn!({ data: { code: 0, data: mockOverview, message: 'ok' } })
      await promise

      expect(store.loading).toBe(false)
    })

    it('sets loading false even on error', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('Server error'))

      const store = useCommunityStore()
      await expect(store.fetchOverview()).rejects.toThrow('Server error')

      expect(store.loading).toBe(false)
      expect(store.overview).toBeNull()
    })
  })

  describe('fetchInsights', () => {
    const mockInsights = {
      dailyActiveUsers: [100, 120, 115, 130],
      retentionRate: 0.72,
      topModules: ['tasks', 'leaderboard'],
    }

    it('fetches and sets insights data', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockInsights, message: 'ok' },
      })

      const store = useCommunityStore()
      await store.fetchInsights()

      expect(api.get).toHaveBeenCalledWith('/api/v1/community/insights')
      expect(store.insights).toEqual(mockInsights)
    })

    it('logs error and keeps insights null on failure', async () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      vi.mocked(api.get).mockRejectedValueOnce(new Error('Timeout'))

      const store = useCommunityStore()
      await store.fetchInsights()

      expect(store.insights).toBeNull()
      expect(consoleSpy).toHaveBeenCalledWith('Failed to fetch insights:', expect.any(Error))
    })

    it('does not throw on error (swallows it)', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => {})
      vi.mocked(api.get).mockRejectedValueOnce(new Error('fail'))

      const store = useCommunityStore()
      await expect(store.fetchInsights()).resolves.toBeUndefined()
    })
  })

  describe('fetchChecklist', () => {
    const mockChecklist = [
      { id: '1', title: 'Configure modules', completed: true },
      { id: '2', title: 'Add first task', completed: false },
      { id: '3', title: 'Share with community', completed: false },
    ]

    it('fetches and sets checklist items', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockChecklist, message: 'ok' },
      })

      const store = useCommunityStore()
      await store.fetchChecklist()

      expect(api.get).toHaveBeenCalledWith('/api/v1/community/checklist')
      expect(store.checklist).toEqual(mockChecklist)
      expect(store.checklist).toHaveLength(3)
    })

    it('logs error and keeps checklist empty on failure', async () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      vi.mocked(api.get).mockRejectedValueOnce(new Error('Network'))

      const store = useCommunityStore()
      await store.fetchChecklist()

      expect(store.checklist).toEqual([])
      expect(consoleSpy).toHaveBeenCalledWith('Failed to fetch checklist:', expect.any(Error))
    })

    it('replaces previous checklist data on re-fetch', async () => {
      const store = useCommunityStore()

      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: [{ id: '1', title: 'A', completed: false }], message: 'ok' },
      })
      await store.fetchChecklist()
      expect(store.checklist).toHaveLength(1)

      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockChecklist, message: 'ok' },
      })
      await store.fetchChecklist()
      expect(store.checklist).toHaveLength(3)
    })
  })
})
