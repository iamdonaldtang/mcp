<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { cApi } from '../../api/client'
import { useCEndStore } from '../../stores/c-end'
import type { LeaderboardEntry } from '../../types/c-end'

const router = useRouter()
const store = useCEndStore()
const loading = ref(true)

const period = ref<'weekly' | 'monthly' | 'alltime'>('alltime')
const pointType = ref('Community Points')
const podium = ref<LeaderboardEntry[]>([])
const rankings = ref<LeaderboardEntry[]>([])
const userRank = ref<LeaderboardEntry | null>(null)

onMounted(() => fetchLeaderboard())

async function fetchLeaderboard() {
  loading.value = true
  try {
    const res = await cApi.get('/api/c/leaderboard', { params: { period: period.value } })
    const data = res.data.data?.rankings || []
    podium.value = data.slice(0, 3)
    rankings.value = data.slice(3)
    // Find user rank
    if (store.userStatus) {
      userRank.value = data.find((e: LeaderboardEntry) =>
        e.wallet_address === store.userStatus?.walletAddress
      ) || null
    }
  } finally {
    loading.value = false
  }
}

function truncAddr(addr: string) {
  return addr ? addr.slice(0, 6) + '...' + addr.slice(-4) : '—'
}

function changePeriod(p: 'weekly' | 'monthly' | 'alltime') {
  period.value = p
  fetchLeaderboard()
}

const podiumColors = ['#F59E0B', '#CBD5E1', '#CD7F32']
const podiumHeights = ['h-28', 'h-20', 'h-16']
const podiumOrder = [1, 0, 2] // Display #2, #1, #3
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-text-primary mb-2">Leaderboard</h1>
      <div class="flex items-center gap-4">
        <span class="text-sm text-text-secondary">{{ pointType }}</span>
        <div class="flex gap-2">
          <button
            v-for="p in (['weekly', 'monthly', 'alltime'] as const)"
            :key="p"
            class="px-3 py-1 text-xs font-medium rounded-lg transition-colors capitalize"
            :class="period === p ? 'bg-c-accent text-black' : 'bg-card-bg border border-border text-text-secondary'"
            @click="changePeriod(p)"
          >
            {{ p === 'alltime' ? 'All Time' : p === 'weekly' ? 'This Week' : 'This Month' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Podium (Top 3) -->
    <div v-if="podium.length >= 3" class="flex items-end justify-center gap-6 py-8">
      <template v-for="idx in podiumOrder" :key="idx">
        <div class="flex flex-col items-center">
          <div class="w-14 h-14 rounded-full border-2 flex items-center justify-center mb-2" :style="{ borderColor: podiumColors[idx] }">
            <span class="material-symbols-rounded text-2xl" :style="{ color: podiumColors[idx] }">person</span>
          </div>
          <span v-if="idx === 0" class="material-symbols-rounded text-lg mb-1" style="color: #F59E0B">emoji_events</span>
          <div class="text-xs text-text-primary font-medium">{{ truncAddr(podium[idx]?.wallet_address || '') }}</div>
          <div class="text-sm font-bold" :style="{ color: podiumColors[idx] }">{{ podium[idx]?.points?.toLocaleString() }} pts</div>
          <div class="mt-2 w-20 rounded-t-lg flex items-end justify-center" :class="podiumHeights[idx]" :style="{ background: podiumColors[idx] + '30' }">
            <span class="text-2xl font-bold pb-2" :style="{ color: podiumColors[idx] }">#{{ idx + 1 }}</span>
          </div>
        </div>
      </template>
    </div>

    <!-- Your Rank -->
    <div v-if="userRank" class="bg-c-accent/10 border border-c-accent/30 rounded-xl p-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <span class="text-lg font-bold text-c-accent">#{{ userRank.rank }}</span>
        <div class="w-8 h-8 rounded-full bg-c-accent/20 flex items-center justify-center">
          <span class="material-symbols-rounded text-c-accent text-sm">person</span>
        </div>
        <span class="text-sm text-text-primary">You · Lv.{{ userRank.level }}</span>
      </div>
      <span class="text-sm font-bold text-c-accent">{{ userRank.points?.toLocaleString() }} pts</span>
    </div>

    <!-- Rankings Table -->
    <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-border">
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase w-16">Rank</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase">User</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase w-20">Level</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase w-28">Points</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="4" class="px-4 py-3"><div class="h-4 bg-border rounded animate-pulse"></div></td>
          </tr>
          <tr
            v-else
            v-for="entry in rankings"
            :key="entry.rank"
            class="border-b border-border last:border-b-0 hover:bg-white/2"
          >
            <td class="px-4 py-3 text-sm font-medium text-text-secondary">#{{ entry.rank }}</td>
            <td class="px-4 py-3">
              <div class="flex items-center gap-2">
                <div class="w-7 h-7 rounded-full bg-border flex items-center justify-center">
                  <span class="material-symbols-rounded text-text-muted text-xs">person</span>
                </div>
                <span class="text-sm text-text-primary">{{ truncAddr(entry.wallet_address) }}</span>
              </div>
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary">Lv.{{ entry.level }}</td>
            <td class="px-4 py-3 text-sm text-text-primary text-right font-medium">{{ entry.points?.toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Cross-sell -->
    <div class="grid grid-cols-2 gap-4">
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/lb-sprint')">
        <span class="text-sm text-text-primary font-medium">Earn more points in the Sprint →</span>
      </button>
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/shop')">
        <span class="text-sm text-text-primary font-medium">Redeem rewards in the Shop →</span>
      </button>
    </div>
  </div>
</template>
