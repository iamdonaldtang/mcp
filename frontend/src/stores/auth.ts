import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../api/client'

export interface BUser {
  id: string
  email: string
  name: string
  role: string
  avatar?: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('b_token'))
  const user = ref<BUser | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  async function login(email: string, password: string) {
    loading.value = true
    try {
      const res = await api.post('/api/v1/auth/login', { email, password })
      token.value = res.data.data.token
      user.value = res.data.data.user
      localStorage.setItem('b_token', token.value!)
    } finally {
      loading.value = false
    }
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res = await api.get('/api/v1/auth/profile')
      user.value = res.data.data
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('b_token')
  }

  return { token, user, loading, isAuthenticated, login, fetchProfile, logout }
})
