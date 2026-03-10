<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { cApi } from '../../api/client'
import type { Sprint } from '../../types/c-end'

const router = useRouter()
const loading = ref(true)
const currentSprint = ref<Sprint | null>(null)
const pastSprints = ref<any[]>([])
const hasActive = ref(false)

onMounted(async () => {
  try {
    const [currentRes, historyRes] = await Promise.allSettled([
      cApi.get('/api/c/lb-sprint/current'),
      cApi.get('/api/c/lb-sprint/history'),
    ])
    if (currentRes.status === 'fulfilled') {
      const d = currentRes.value.data.data
      if (d?.active) { currentSprint.value = d.sprint; hasActive.value = true }
    }
    if (historyRes.status === 'fulfilled') {
      pastSprints.value = historyRes.value.data.data || []
    }
  } finally { loading.value = false }
})

function countdown(endsAt: string) {
  const diff = new Date(endsAt).getTime() - Date.now()
  if (diff <= 0) return 'Ended'
  const d = Math.floor(diff / 86400000)
  const h = Math.floor((diff % 86400000) / 3600000)
  const m = Math.floor((diff % 3600000) / 60000)
  return `${d}d ${h}h ${m}m`
}
</script>

<template>
  <div class="space-y-8">
    <template v-if="loading">
      <div class="bg-card-bg border border-border rounded-xl h-40 animate-pulse"></div>
      <div class="space-y-3">
        <div v-for="i in 3" :key="i" class="bg-card-bg border border-border rounded-xl h-16 animate-pulse"></div>
      </div>
    </template>

    <template v-else-if="hasActive && currentSprint">
      <!-- Sprint Header -->
      <div class="bg-card-bg border border-border rounded-2xl p-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h1 class="text-2xl font-bold text-text-primary">{{ currentSprint.name }}</h1>
            <span class="text-sm text-text-secondary">{{ currentSprint.pointType }} · {{ currentSprint.totalParticipants.toLocaleString() }} participants</span>
          </div>
          <span class="px-3 py-1.5 bg-c-accent/20 text-c-accent text-sm font-semibold rounded-lg">
            Ends in {{ countdown(currentSprint.endsAt) }}
          </span>
        </div>
        <div class="text-3xl font-bold text-c-accent mb-2">{{ currentSprint.userPoints.toLocaleString() }}</div>
        <span class="text-xs text-text-muted">points earned this sprint</span>
      </div>

      <!-- Sprint Tasks -->
      <div>
        <h2 class="text-base font-semibold text-text-primary mb-3">Sprint Tasks</h2>
        <div class="space-y-2">
          <div v-for="task in currentSprint.tasks" :key="task.id"
            class="bg-card-bg border border-border rounded-xl p-4 flex items-center gap-4">
            <span class="material-symbols-rounded text-lg" :class="task.status === 'completed' ? 'text-status-active' : 'text-text-muted'">
              {{ task.status === 'completed' ? 'check_circle' : 'radio_button_unchecked' }}
            </span>
            <span class="flex-1 text-sm" :class="task.status === 'completed' ? 'text-text-secondary line-through' : 'text-text-primary'">{{ task.name }}</span>
            <span class="text-xs font-semibold text-c-accent">+{{ task.points }} pts</span>
          </div>
        </div>
      </div>

      <!-- Reward Tiers -->
      <div>
        <h2 class="text-base font-semibold text-text-primary mb-3">Reward Tiers</h2>
        <div class="space-y-2">
          <div v-for="tier in currentSprint.rewardTiers" :key="tier.threshold"
            class="bg-card-bg border rounded-xl p-4 flex items-center justify-between"
            :class="tier.status === 'earned' ? 'border-status-active' : tier.status === 'claimable' ? 'border-c-accent' : 'border-border'">
            <div class="flex items-center gap-3">
              <span class="material-symbols-rounded text-lg" :class="tier.status === 'earned' ? 'text-status-active' : tier.status === 'claimable' ? 'text-c-accent' : 'text-text-muted'">
                {{ tier.status === 'earned' ? 'check_circle' : tier.status === 'claimable' ? 'redeem' : 'lock' }}
              </span>
              <div>
                <div class="text-sm text-text-primary">{{ tier.reward }}</div>
                <div class="text-xs text-text-muted">{{ tier.threshold.toLocaleString() }} pts required</div>
              </div>
            </div>
            <span v-if="tier.status === 'earned'" class="text-xs text-status-active font-medium">Earned ✓</span>
            <button v-else-if="tier.status === 'claimable'" class="px-3 py-1 text-xs bg-c-accent text-black rounded-lg font-medium">Claim</button>
            <span v-else class="material-symbols-rounded text-text-muted text-base">lock</span>
          </div>
        </div>
      </div>

      <!-- Rankings -->
      <div v-if="currentSprint.rankings.length > 0">
        <h2 class="text-base font-semibold text-text-primary mb-3">Sprint Rankings</h2>
        <div class="space-y-1">
          <div v-for="entry in currentSprint.rankings.slice(0, 10)" :key="entry.rank"
            class="bg-card-bg border border-border rounded-lg px-4 py-2.5 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="text-sm font-medium text-text-muted w-6">#{{ entry.rank }}</span>
              <span class="text-sm text-text-primary">{{ entry.wallet_address?.slice(0, 6) }}...{{ entry.wallet_address?.slice(-4) }}</span>
            </div>
            <span class="text-sm font-medium text-c-accent">{{ entry.points?.toLocaleString() }} pts</span>
          </div>
        </div>
      </div>
    </template>

    <!-- No active sprint -->
    <template v-else>
      <div class="text-center py-16">
        <span class="material-symbols-rounded text-5xl text-text-muted block mb-3">sprint</span>
        <h2 class="text-lg font-semibold text-text-primary mb-1">No Active Sprint</h2>
        <p class="text-sm text-text-secondary">Check back soon for the next leaderboard sprint challenge</p>
      </div>
    </template>

    <!-- Past Sprints -->
    <div v-if="pastSprints.length > 0">
      <h2 class="text-base font-semibold text-text-primary mb-3">Past Sprints</h2>
      <div class="space-y-2">
        <div v-for="sprint in pastSprints" :key="sprint.id" class="bg-card-bg border border-border rounded-xl p-4 flex items-center justify-between">
          <div>
            <div class="text-sm font-medium text-text-primary">{{ sprint.name }}</div>
            <div class="text-xs text-text-muted">{{ sprint.participants }} participants</div>
          </div>
          <span class="px-2.5 py-0.5 text-xs rounded-full bg-status-completed-bg text-status-completed">Ended</span>
        </div>
      </div>
    </div>

    <!-- Cross-sell -->
    <div class="grid grid-cols-2 gap-4">
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/shop')">
        <span class="text-sm text-text-primary font-medium">Spend your points in the Shop →</span>
      </button>
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/leaderboard')">
        <span class="text-sm text-text-primary font-medium">Track your rank on the Leaderboard →</span>
      </button>
    </div>
  </div>
</template>
