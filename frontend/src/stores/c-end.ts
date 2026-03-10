import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { cApi } from '../api/client'
import type { UserStatus, TabVisibility } from '../types/c-end'

export const useCEndStore = defineStore('c-end', () => {
  const walletToken = ref<string | null>(localStorage.getItem('c_token'))
  const userStatus = ref<UserStatus | null>(null)
  const communityName = ref('Community')

  const tabVisibility = ref<TabVisibility>({
    home: true,
    quests: true,
    leaderboard: true,
    lbSprint: true,
    milestone: true,
    shop: true,
  })

  const isWalletConnected = computed(() => !!walletToken.value)

  async function connectWallet() {
    // TODO: integrate actual wallet provider (MetaMask, WalletConnect, etc.)
    // For now, mock the flow
    try {
      const mockAddress = '0x' + Math.random().toString(16).slice(2, 42)
      const res = await cApi.post('/api/c/wallet/connect', { address: mockAddress })
      walletToken.value = res.data.data.token
      localStorage.setItem('c_token', walletToken.value!)
      await fetchUserStatus()
    } catch (e) {
      console.error('Wallet connect failed:', e)
    }
  }

  async function fetchUserStatus() {
    if (!walletToken.value) return
    try {
      const res = await cApi.get('/api/c/user/status')
      userStatus.value = res.data.data
    } catch {
      disconnectWallet()
    }
  }

  async function fetchTabVisibility() {
    try {
      const res = await cApi.get('/api/c/community/home')
      if (res.data.data.tabVisibility) {
        tabVisibility.value = res.data.data.tabVisibility
      }
    } catch {
      // keep defaults
    }
  }

  function disconnectWallet() {
    walletToken.value = null
    userStatus.value = null
    localStorage.removeItem('c_token')
  }

  return {
    walletToken, userStatus, communityName,
    tabVisibility, isWalletConnected,
    connectWallet, fetchUserStatus, fetchTabVisibility, disconnectWallet,
  }
})
