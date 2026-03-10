<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { cApi } from '../../api/client'
import type { ActivityItem } from '../../types/c-end'

const loading = ref(true)
const activities = ref<ActivityItem[]>([])
const filter = ref('all')
const filters = ['all', 'task', 'points', 'reward', 'level_up', 'invite']
let pollInterval: number | null = null

const activityColors: Record<string, string> = {
  task: '#48BB78', points: '#48BB78', level_up: '#F59E0B',
  invite: '#9B7EE0', wheel: '#EF4444', streak: '#F59E0B',
  shop: '#EC4899', reward: '#F59E0B',
}

const filteredActivities = computed(() => {
  if (filter.value === 'all') return activities.value
  return activities.value.filter(a => a.type === filter.value)
})

// Group by day
const groupedActivities = computed(() => {
  const groups: { label: string; items: ActivityItem[] }[] = []
  const today = new Date().toDateString()
  const yesterday = new Date(Date.now() - 86400000).toDateString()
  let currentGroup = ''

  for (const act of filteredActivities.value) {
    const dayStr = new Date(act.timestamp).toDateString()
    const label = dayStr === today ? 'Today' : dayStr === yesterday ? 'Yesterday' : new Date(act.timestamp).toLocaleDateString()
    if (label !== currentGroup) {
      currentGroup = label
      groups.push({ label, items: [] })
    }
    groups[groups.length - 1].items.push(act)
  }
  return groups
})

onMounted(async () => {
  await fetchFeed()
  // Auto-refresh every 15s
  pollInterval = window.setInterval(fetchFeed, 15000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})

async function fetchFeed() {
  try {
    const res = await cApi.get('/api/c/activity/feed')
    activities.value = res.data.data || []
  } finally { loading.value = false }
}

function timeAgo(ts: string) {
  const diff = Date.now() - new Date(ts).getTime()
  const m = Math.floor(diff / 60000)
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  return `${Math.floor(h / 24)}d ago`
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-text-primary">Activity Feed</h1>
      <div class="flex items-center gap-2">
        <div class="w-2 h-2 rounded-full bg-status-active animate-pulse"></div>
        <span class="text-xs text-status-active font-medium">Live</span>
      </div>
    </div>

    <!-- Filter Chips -->
    <div class="flex gap-2 flex-wrap">
      <button v-for="f in filters" :key="f"
        class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors capitalize"
        :class="filter === f ? 'bg-c-accent text-black' : 'bg-card-bg border border-border text-text-secondary'"
        @click="filter = f">
        {{ f === 'all' ? 'All' : f === 'level_up' ? 'Level Up' : f.charAt(0).toUpperCase() + f.slice(1) }}
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="space-y-3">
      <div v-for="i in 5" :key="i" class="bg-card-bg border border-border rounded-xl h-14 animate-pulse"></div>
    </div>

    <!-- Feed grouped by day -->
    <template v-else>
      <div v-for="group in groupedActivities" :key="group.label" class="space-y-2">
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider">{{ group.label }}</div>
        <div v-for="act in group.items" :key="act.id"
          class="bg-card-bg border border-border rounded-xl px-4 py-3 flex items-center gap-3">
          <div class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ background: activityColors[act.type] || '#94A3B8' }"></div>
          <div class="flex-1 min-w-0">
            <div class="text-sm text-text-primary">{{ act.title }}</div>
            <div v-if="act.description" class="text-xs text-text-muted truncate">{{ act.description }}</div>
          </div>
          <span v-if="act.pointsDelta" class="text-xs font-semibold shrink-0"
            :class="(act.pointsDelta || 0) > 0 ? 'text-c-accent' : 'text-status-paused'">
            {{ (act.pointsDelta || 0) > 0 ? '+' : '' }}{{ act.pointsDelta }} pts
          </span>
          <span class="text-xs text-text-muted shrink-0">{{ timeAgo(act.timestamp) }}</span>
        </div>
      </div>

      <div v-if="filteredActivities.length === 0" class="text-center py-12">
        <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">history</span>
        <p class="text-sm text-text-muted">No activity to show</p>
      </div>
    </template>
  </div>
</template>
