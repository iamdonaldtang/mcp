import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/client'
import type { WhiteLabelOverview, WLChecklistItem } from '../types/b-whitelabel'

export const useWhiteLabelStore = defineStore('whitelabel', () => {
  const overview = ref<WhiteLabelOverview | null>(null)
  const checklist = ref<WLChecklistItem[]>([])
  const loading = ref(false)

  async function fetchOverview() {
    loading.value = true
    try {
      const res = await api.get('/api/v1/whitelabel/overview')
      overview.value = res.data.data
    } finally {
      loading.value = false
    }
  }

  async function fetchChecklist() {
    try {
      const res = await api.get('/api/v1/whitelabel/checklist')
      checklist.value = res.data.data
    } catch (e) {
      console.error('Failed to fetch WL checklist:', e)
    }
  }

  return { overview, checklist, loading, fetchOverview, fetchChecklist }
})
