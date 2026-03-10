import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/client'
import type { CommunityOverview, CommunityInsights, ChecklistItem } from '../types/b-community'

export const useCommunityStore = defineStore('community', () => {
  const overview = ref<CommunityOverview | null>(null)
  const insights = ref<CommunityInsights | null>(null)
  const checklist = ref<ChecklistItem[]>([])
  const loading = ref(false)

  async function fetchOverview() {
    loading.value = true
    try {
      const res = await api.get('/api/v1/community/overview')
      overview.value = res.data.data
    } finally {
      loading.value = false
    }
  }

  async function fetchInsights() {
    try {
      const res = await api.get('/api/v1/community/insights')
      insights.value = res.data.data
    } catch (e) {
      console.error('Failed to fetch insights:', e)
    }
  }

  async function fetchChecklist() {
    try {
      const res = await api.get('/api/v1/community/checklist')
      checklist.value = res.data.data
    } catch (e) {
      console.error('Failed to fetch checklist:', e)
    }
  }

  return { overview, insights, checklist, loading, fetchOverview, fetchInsights, fetchChecklist }
})
