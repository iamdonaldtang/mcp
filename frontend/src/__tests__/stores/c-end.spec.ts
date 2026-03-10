import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('@/api/client', () => import('../mocks/api'))

import { useCEndStore } from '@/stores/c-end'
import { cApi } from '@/api/client'

describe('c-end store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('initial state', () => {
    it('has null walletToken and userStatus by default', () => {
      const store = useCEndStore()
      expect(store.walletToken).toBeNull()
      expect(store.userStatus).toBeNull()
    })

    it('has default communityName', () => {
      const store = useCEndStore()
      expect(store.communityName).toBe('Community')
    })

    it('has all tabs visible by default', () => {
      const store = useCEndStore()
      expect(store.tabVisibility).toEqual({
        home: true,
        quests: true,
        leaderboard: true,
        lbSprint: true,
        milestone: true,
        shop: true,
      })
    })

    it('isWalletConnected is false when no token', () => {
      const store = useCEndStore()
      expect(store.isWalletConnected).toBe(false)
    })

    it('reads walletToken from localStorage on creation', () => {
      localStorage.setItem('c_token', 'saved-wallet-token')
      setActivePinia(createPinia())
      const store = useCEndStore()
      expect(store.walletToken).toBe('saved-wallet-token')
      expect(store.isWalletConnected).toBe(true)
    })
  })

  describe('connectWallet', () => {
    const mockToken = 'c-jwt-token-456'
    const mockUserStatus = {
      address: '0xabc123',
      points: 500,
      level: 3,
      joinedAt: '2026-01-15',
    }

    it('calls wallet connect API and sets token + localStorage', async () => {
      vi.mocked(cApi.post).mockResolvedValueOnce({
        data: { code: 0, data: { token: mockToken }, message: 'ok' },
      })
      vi.mocked(cApi.get).mockResolvedValueOnce({
        data: { code: 0, data: mockUserStatus, message: 'ok' },
      })

      const store = useCEndStore()
      await store.connectWallet()

      expect(cApi.post).toHaveBeenCalledWith('/api/c/wallet/connect', expect.objectContaining({
        address: expect.stringMatching(/^0x/),
      }))
      expect(store.walletToken).toBe(mockToken)
      expect(store.isWalletConnected).toBe(true)
      expect(localStorage.getItem('c_token')).toBe(mockToken)
    })

    it('calls fetchUserStatus after successful connect', async () => {
      vi.mocked(cApi.post).mockResolvedValueOnce({
        data: { code: 0, data: { token: mockToken }, message: 'ok' },
      })
      vi.mocked(cApi.get).mockResolvedValueOnce({
        data: { code: 0, data: mockUserStatus, message: 'ok' },
      })

      const store = useCEndStore()
      await store.connectWallet()

      expect(cApi.get).toHaveBeenCalledWith('/api/c/user/status')
      expect(store.userStatus).toEqual(mockUserStatus)
    })

    it('logs error and does not set token on connect failure', async () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      vi.mocked(cApi.post).mockRejectedValueOnce(new Error('Wallet rejected'))

      const store = useCEndStore()
      await store.connectWallet()

      expect(store.walletToken).toBeNull()
      expect(store.isWalletConnected).toBe(false)
      expect(localStorage.getItem('c_token')).toBeNull()
      expect(consoleSpy).toHaveBeenCalledWith('Wallet connect failed:', expect.any(Error))
    })

    it('does not throw on failure (swallows error)', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => {})
      vi.mocked(cApi.post).mockRejectedValueOnce(new Error('fail'))

      const store = useCEndStore()
      await expect(store.connectWallet()).resolves.toBeUndefined()
    })
  })

  describe('fetchUserStatus', () => {
    const mockUserStatus = {
      address: '0xdef456',
      points: 1200,
      level: 5,
      joinedAt: '2025-12-01',
    }

    it('fetches and sets user status when wallet is connected', async () => {
      const store = useCEndStore()
      store.walletToken = 'valid-c-token'

      vi.mocked(cApi.get).mockResolvedValueOnce({
        data: { code: 0, data: mockUserStatus, message: 'ok' },
      })

      await store.fetchUserStatus()

      expect(cApi.get).toHaveBeenCalledWith('/api/c/user/status')
      expect(store.userStatus).toEqual(mockUserStatus)
    })

    it('does nothing when no walletToken', async () => {
      const store = useCEndStore()
      expect(store.walletToken).toBeNull()

      await store.fetchUserStatus()

      expect(cApi.get).not.toHaveBeenCalled()
      expect(store.userStatus).toBeNull()
    })

    it('calls disconnectWallet on fetch error', async () => {
      localStorage.setItem('c_token', 'expired-token')
      setActivePinia(createPinia())
      const store = useCEndStore()

      vi.mocked(cApi.get).mockRejectedValueOnce(new Error('401'))

      await store.fetchUserStatus()

      expect(store.walletToken).toBeNull()
      expect(store.userStatus).toBeNull()
      expect(localStorage.getItem('c_token')).toBeNull()
    })
  })

  describe('fetchTabVisibility', () => {
    it('updates tabVisibility from API response', async () => {
      const customVisibility = {
        home: true,
        quests: true,
        leaderboard: false,
        lbSprint: false,
        milestone: true,
        shop: false,
      }

      vi.mocked(cApi.get).mockResolvedValueOnce({
        data: { code: 0, data: { tabVisibility: customVisibility }, message: 'ok' },
      })

      const store = useCEndStore()
      await store.fetchTabVisibility()

      expect(cApi.get).toHaveBeenCalledWith('/api/c/community/home')
      expect(store.tabVisibility).toEqual(customVisibility)
    })

    it('keeps default tabVisibility when response has no tabVisibility field', async () => {
      vi.mocked(cApi.get).mockResolvedValueOnce({
        data: { code: 0, data: { communityName: 'My Project' }, message: 'ok' },
      })

      const store = useCEndStore()
      await store.fetchTabVisibility()

      // All should remain true (defaults)
      expect(store.tabVisibility.home).toBe(true)
      expect(store.tabVisibility.quests).toBe(true)
      expect(store.tabVisibility.leaderboard).toBe(true)
      expect(store.tabVisibility.lbSprint).toBe(true)
      expect(store.tabVisibility.milestone).toBe(true)
      expect(store.tabVisibility.shop).toBe(true)
    })

    it('keeps defaults on API error', async () => {
      vi.mocked(cApi.get).mockRejectedValueOnce(new Error('503'))

      const store = useCEndStore()
      await store.fetchTabVisibility()

      expect(store.tabVisibility).toEqual({
        home: true,
        quests: true,
        leaderboard: true,
        lbSprint: true,
        milestone: true,
        shop: true,
      })
    })

    it('does not throw on error', async () => {
      vi.mocked(cApi.get).mockRejectedValueOnce(new Error('fail'))

      const store = useCEndStore()
      await expect(store.fetchTabVisibility()).resolves.toBeUndefined()
    })
  })

  describe('disconnectWallet', () => {
    it('clears walletToken, userStatus, and localStorage', () => {
      const store = useCEndStore()
      store.walletToken = 'some-token'
      store.userStatus = { address: '0x123', points: 100, level: 1, joinedAt: '2026-01-01' } as any
      localStorage.setItem('c_token', 'some-token')

      store.disconnectWallet()

      expect(store.walletToken).toBeNull()
      expect(store.userStatus).toBeNull()
      expect(store.isWalletConnected).toBe(false)
      expect(localStorage.getItem('c_token')).toBeNull()
    })

    it('is safe to call when already disconnected', () => {
      const store = useCEndStore()
      expect(() => store.disconnectWallet()).not.toThrow()
      expect(store.walletToken).toBeNull()
    })

    it('does not affect tabVisibility or communityName', () => {
      const store = useCEndStore()
      store.walletToken = 'tok'
      store.communityName = 'My Community'

      store.disconnectWallet()

      expect(store.communityName).toBe('My Community')
      expect(store.tabVisibility.home).toBe(true)
    })
  })
})
