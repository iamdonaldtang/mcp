<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { cApi } from '../../api/client'
import { useCEndStore } from '../../stores/c-end'
import type { Achievement, ActivityItem, ReferralStats } from '../../types/c-end'

const router = useRouter()
const store = useCEndStore()
const loading = ref(true)

const profile = ref<any>(null)
const achievements = ref<Achievement[]>([])
const activities = ref<ActivityItem[]>([])
const referralStats = ref<ReferralStats | null>(null)

const activityColors: Record<string, string> = {
  task: '#48BB78', level_up: '#F59E0B', invite: '#9B7EE0',
  wheel: '#EF4444', milestone: '#3B82F6', shop: '#EC4899',
  streak: '#F59E0B', points: '#48BB78', reward: '#F59E0B',
}

onMounted(async () => {
  try {
    const [profileRes, achieveRes, actRes, refRes] = await Promise.allSettled([
      cApi.get('/api/c/user/profile'),
      cApi.get('/api/c/user/achievements'),
      cApi.get('/api/c/user/activity'),
      cApi.get('/api/c/user/referral-stats'),
    ])
    if (profileRes.status === 'fulfilled') profile.value = profileRes.value.data.data
    if (achieveRes.status === 'fulfilled') achievements.value = achieveRes.value.data.data || []
    if (actRes.status === 'fulfilled') activities.value = actRes.value.data.data || []
    if (refRes.status === 'fulfilled') referralStats.value = refRes.value.data.data
  } finally { loading.value = false }
})

function truncAddr(addr: string) { return addr ? addr.slice(0, 6) + '...' + addr.slice(-4) : '—' }

function timeAgo(ts: string) {
  const diff = Date.now() - new Date(ts).getTime()
  const h = Math.floor(diff / 3600000)
  if (h < 1) return 'just now'
  if (h < 24) return `${h}h ago`
  const d = Math.floor(h / 24)
  return `${d}d ago`
}
</script>

<template>
  <div class="space-y-8">
    <!-- Profile Card -->
    <div class="bg-card-bg border border-border rounded-2xl p-6">
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 rounded-full bg-c-accent/20 flex items-center justify-center">
          <span class="material-symbols-rounded text-c-accent text-3xl">person</span>
        </div>
        <div>
          <div class="text-lg font-bold text-text-primary">{{ truncAddr(store.userStatus?.walletAddress || '') }}</div>
          <div class="text-sm text-text-secondary">
            Joined {{ profile?.joined_at ? new Date(profile.joined_at).toLocaleDateString('en-US', { month: 'short', year: 'numeric' }) : '—' }}
            · Level {{ store.userStatus?.level || 1 }}
            · {{ (store.userStatus?.xp || 0).toLocaleString() }} XP Lifetime
          </div>
          <!-- Level progress -->
          <div class="mt-2 w-48 h-1.5 bg-page-bg rounded-full overflow-hidden">
            <div class="h-full bg-c-accent rounded-full" style="width: 65%"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-text-primary">{{ (store.userStatus?.xp || 0).toLocaleString() }}</div>
        <div class="text-xs text-text-muted">Total Points</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-text-primary">{{ profile?.tasks_done || 0 }}</div>
        <div class="text-xs text-text-muted">Tasks Done</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-text-primary">{{ store.userStatus?.dayStreak || 0 }}</div>
        <div class="text-xs text-text-muted">Day Streak</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-text-primary">#{{ store.userStatus?.rank || '—' }}</div>
        <div class="text-xs text-text-muted">Rank</div>
      </div>
    </div>

    <!-- Achievements -->
    <div>
      <h2 class="text-base font-semibold text-text-primary mb-3">Achievements</h2>
      <div class="grid grid-cols-6 gap-3">
        <div v-for="badge in achievements.slice(0, 6)" :key="badge.id"
          class="bg-card-bg border border-border rounded-xl p-3 text-center"
          :class="badge.isEarned ? '' : 'opacity-40'">
          <span class="material-symbols-rounded text-2xl block mb-1" :class="badge.isEarned ? 'text-c-accent' : 'text-text-muted'">
            {{ badge.icon || 'emoji_events' }}
          </span>
          <div class="text-[10px] font-medium" :class="badge.isEarned ? 'text-text-primary' : 'text-text-muted'">{{ badge.name }}</div>
          <span v-if="!badge.isEarned" class="material-symbols-rounded text-xs text-text-muted">lock</span>
        </div>
        <!-- Fill empty slots -->
        <template v-if="achievements.length < 6">
          <div v-for="i in (6 - achievements.length)" :key="'empty-' + i" class="bg-card-bg border border-border rounded-xl p-3 text-center opacity-30">
            <span class="material-symbols-rounded text-2xl text-text-muted block mb-1">help</span>
            <div class="text-[10px] text-text-muted">???</div>
          </div>
        </template>
      </div>
    </div>

    <!-- Recent Activity -->
    <div>
      <h2 class="text-base font-semibold text-text-primary mb-3">Recent Activity</h2>
      <div class="space-y-1">
        <div v-for="act in activities.slice(0, 6)" :key="act.id"
          class="bg-card-bg border border-border rounded-lg px-4 py-3 flex items-center gap-3">
          <div class="w-2 h-2 rounded-full shrink-0" :style="{ background: activityColors[act.type] || '#94A3B8' }"></div>
          <div class="flex-1 min-w-0">
            <div class="text-sm text-text-primary truncate">{{ act.title }}</div>
          </div>
          <span v-if="act.pointsDelta" class="text-xs font-semibold shrink-0" :class="(act.pointsDelta || 0) > 0 ? 'text-c-accent' : 'text-status-paused'">
            {{ (act.pointsDelta || 0) > 0 ? '+' : '' }}{{ act.pointsDelta }} pts
          </span>
          <span class="text-xs text-text-muted shrink-0">{{ timeAgo(act.timestamp) }}</span>
        </div>
        <div v-if="activities.length === 0" class="text-center py-8">
          <span class="material-symbols-rounded text-3xl text-text-muted block mb-1">history</span>
          <p class="text-sm text-text-muted">No activity yet</p>
        </div>
      </div>
    </div>

    <!-- Referral Program -->
    <div v-if="referralStats">
      <h2 class="text-base font-semibold text-text-primary mb-3">Referral Program</h2>
      <div class="bg-card-bg border border-border rounded-xl p-5">
        <div class="grid grid-cols-3 gap-4 mb-4">
          <div class="text-center">
            <div class="text-lg font-bold text-text-primary">{{ referralStats.totalInvites }}</div>
            <div class="text-xs text-text-muted">Total Referrals</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold text-text-primary">{{ referralStats.pointsEarned }}</div>
            <div class="text-xs text-text-muted">Bonus Points</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold text-text-primary">{{ referralStats.conversionRate }}%</div>
            <div class="text-xs text-text-muted">Conversion Rate</div>
          </div>
        </div>
        <button class="w-full py-2.5 bg-c-accent text-black text-sm font-medium rounded-lg hover:bg-c-accent/90" @click="router.push('/c/invite')">
          Invite More Friends
        </button>
      </div>
    </div>
  </div>
</template>
