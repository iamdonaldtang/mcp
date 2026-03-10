import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('@/api/client', () => import('../mocks/api'))

import { useWhiteLabelStore } from '@/stores/whitelabel'
import { api } from '@/api/client'

describe('whitelabel store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('has null overview, empty checklist, loading false', () => {
      const store = useWhiteLabelStore()
      expect(store.overview).toBeNull()
      expect(store.checklist).toEqual([])
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchOverview', () => {
    const mockOverview = {
      status: 'active',
      domain: 'community.example.com',
      integrationMode: 'embed',
      widgetsDeployed: 3,
      totalUsers: 850,
    }

    it('fetches and sets overview data', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockOverview, message: 'ok' },
      })

      const store = useWhiteLabelStore()
      await store.fetchOverview()

      expect(api.get).toHaveBeenCalledWith('/api/v1/whitelabel/overview')
      expect(store.overview).toEqual(mockOverview)
    })

    it('manages loading state correctly during fetch', async () => {
      let resolveFn: (value: any) => void
      const pending = new Promise((resolve) => { resolveFn = resolve })
      vi.mocked(api.get).mockReturnValueOnce(pending as any)

      const store = useWhiteLabelStore()
      const promise = store.fetchOverview()

      expect(store.loading).toBe(true)

      resolveFn!({ data: { code: 0, data: mockOverview, message: 'ok' } })
      await promise

      expect(store.loading).toBe(false)
    })

    it('sets loading false on error', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('500'))

      const store = useWhiteLabelStore()
      await expect(store.fetchOverview()).rejects.toThrow('500')

      expect(store.loading).toBe(false)
      expect(store.overview).toBeNull()
    })

    it('overwrites previous overview on re-fetch', async () => {
      const store = useWhiteLabelStore()

      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: { status: 'draft' }, message: 'ok' },
      })
      await store.fetchOverview()
      expect(store.overview).toEqual({ status: 'draft' })

      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockOverview, message: 'ok' },
      })
      await store.fetchOverview()
      expect(store.overview).toEqual(mockOverview)
    })
  })

  describe('fetchChecklist', () => {
    const mockChecklist = [
      { id: '1', title: 'Set up domain', completed: true },
      { id: '2', title: 'Configure widgets', completed: false },
      { id: '3', title: 'Deploy integration', completed: false },
    ]

    it('fetches and sets checklist items', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockChecklist, message: 'ok' },
      })

      const store = useWhiteLabelStore()
      await store.fetchChecklist()

      expect(api.get).toHaveBeenCalledWith('/api/v1/whitelabel/checklist')
      expect(store.checklist).toEqual(mockChecklist)
      expect(store.checklist).toHaveLength(3)
    })

    it('logs error and keeps checklist empty on failure', async () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      vi.mocked(api.get).mockRejectedValueOnce(new Error('Network'))

      const store = useWhiteLabelStore()
      await store.fetchChecklist()

      expect(store.checklist).toEqual([])
      expect(consoleSpy).toHaveBeenCalledWith('Failed to fetch WL checklist:', expect.any(Error))
    })

    it('does not throw on API error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => {})
      vi.mocked(api.get).mockRejectedValueOnce(new Error('fail'))

      const store = useWhiteLabelStore()
      await expect(store.fetchChecklist()).resolves.toBeUndefined()
    })

    it('does not affect loading state (no loading wrapper)', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockChecklist, message: 'ok' },
      })

      const store = useWhiteLabelStore()
      // fetchChecklist does NOT toggle loading (only fetchOverview does)
      await store.fetchChecklist()
      expect(store.loading).toBe(false)
    })
  })
})
