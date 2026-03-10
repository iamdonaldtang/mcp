import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('@/api/client', () => import('../mocks/api'))

import { useAuthStore } from '@/stores/auth'
import { api } from '@/api/client'

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('initial state', () => {
    it('has null token and user by default', () => {
      const store = useAuthStore()
      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(store.loading).toBe(false)
    })

    it('isAuthenticated is false when no token', () => {
      const store = useAuthStore()
      expect(store.isAuthenticated).toBe(false)
    })

    it('reads token from localStorage on creation', () => {
      localStorage.setItem('b_token', 'existing-token')
      // Need a fresh pinia so the store re-initializes
      setActivePinia(createPinia())
      const store = useAuthStore()
      expect(store.token).toBe('existing-token')
      expect(store.isAuthenticated).toBe(true)
    })
  })

  describe('login', () => {
    const mockUser = { id: '1', email: 'test@example.com', name: 'Test User', role: 'admin' }
    const mockToken = 'jwt-token-123'

    it('sets token, user, and localStorage on successful login', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({
        data: { code: 0, data: { token: mockToken, user: mockUser }, message: 'success' },
      })

      const store = useAuthStore()
      await store.login('test@example.com', 'password123')

      expect(api.post).toHaveBeenCalledWith('/api/v1/auth/login', {
        email: 'test@example.com',
        password: 'password123',
      })
      expect(store.token).toBe(mockToken)
      expect(store.user).toEqual(mockUser)
      expect(store.isAuthenticated).toBe(true)
      expect(localStorage.getItem('b_token')).toBe(mockToken)
    })

    it('sets loading to true during request and false after', async () => {
      let resolveFn: (value: any) => void
      const pending = new Promise((resolve) => { resolveFn = resolve })
      vi.mocked(api.post).mockReturnValueOnce(pending as any)

      const store = useAuthStore()
      const loginPromise = store.login('a@b.com', 'pw')

      expect(store.loading).toBe(true)

      resolveFn!({
        data: { code: 0, data: { token: 'tok', user: mockUser }, message: 'ok' },
      })
      await loginPromise

      expect(store.loading).toBe(false)
    })

    it('sets loading to false even when login fails', async () => {
      vi.mocked(api.post).mockRejectedValueOnce(new Error('Network error'))

      const store = useAuthStore()
      await expect(store.login('a@b.com', 'pw')).rejects.toThrow('Network error')

      expect(store.loading).toBe(false)
      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })

    it('does not set localStorage when login fails', async () => {
      vi.mocked(api.post).mockRejectedValueOnce(new Error('401'))

      const store = useAuthStore()
      await expect(store.login('a@b.com', 'pw')).rejects.toThrow()

      expect(localStorage.getItem('b_token')).toBeNull()
    })
  })

  describe('fetchProfile', () => {
    const mockUser = { id: '2', email: 'user@test.com', name: 'Profile User', role: 'editor' }

    it('fetches and sets user profile when token exists', async () => {
      const store = useAuthStore()
      store.token = 'valid-token'

      vi.mocked(api.get).mockResolvedValueOnce({
        data: { code: 0, data: mockUser, message: 'ok' },
      })

      await store.fetchProfile()

      expect(api.get).toHaveBeenCalledWith('/api/v1/auth/profile')
      expect(store.user).toEqual(mockUser)
    })

    it('does nothing when no token', async () => {
      const store = useAuthStore()
      expect(store.token).toBeNull()

      await store.fetchProfile()

      expect(api.get).not.toHaveBeenCalled()
      expect(store.user).toBeNull()
    })

    it('calls logout on fetch error', async () => {
      localStorage.setItem('b_token', 'bad-token')
      setActivePinia(createPinia())
      const store = useAuthStore()

      vi.mocked(api.get).mockRejectedValueOnce(new Error('401 Unauthorized'))

      await store.fetchProfile()

      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(localStorage.getItem('b_token')).toBeNull()
    })
  })

  describe('logout', () => {
    it('clears token, user, and localStorage', () => {
      const store = useAuthStore()
      store.token = 'some-token'
      store.user = { id: '1', email: 'a@b.com', name: 'A', role: 'admin' }
      localStorage.setItem('b_token', 'some-token')

      store.logout()

      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
      expect(localStorage.getItem('b_token')).toBeNull()
    })

    it('is safe to call when already logged out', () => {
      const store = useAuthStore()
      expect(() => store.logout()).not.toThrow()
      expect(store.token).toBeNull()
    })
  })
})
