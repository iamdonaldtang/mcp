<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { cApi } from '../../api/client'
import { useCEndStore } from '../../stores/c-end'
import TaskCard from '../../components/c-end/TaskCard.vue'
import type { DayChainState, Sector, Announcement, PulseStats } from '../../types/c-end'

const router = useRouter()
const store = useCEndStore()
const loading = ref(true)

const announcements = ref<Announcement[]>([])
const dayChain = ref<DayChainState | null>(null)
const sectors = ref<Sector[]>([])
const pulse = ref<PulseStats>({ totalMembers: 0, thisWeek: 0, liveActive: 0, tasksDone: 0 })

onMounted(async () => {
  try {
    const [homeRes, tasksRes, announceRes, dcRes] = await Promise.allSettled([
      cApi.get('/api/c/community/home'),
      cApi.get('/api/c/community/tasks'),
      cApi.get('/api/c/community/announcements'),
      store.isWalletConnected ? cApi.get('/api/c/community/daychain') : Promise.resolve(null),
    ])

    if (homeRes.status === 'fulfilled' && homeRes.value) {
      const d = homeRes.value.data.data
      pulse.value = d.pulse || pulse.value
    }
    if (tasksRes.status === 'fulfilled') {
      sectors.value = tasksRes.value.data.data || []
    }
    if (announceRes.status === 'fulfilled') {
      announcements.value = announceRes.value.data.data || []
    }
    if (dcRes.status === 'fulfilled' && dcRes.value) {
      const d = dcRes.value.data.data
      if (d?.enabled) dayChain.value = d
    }
  } finally {
    loading.value = false
  }
})

async function checkIn() {
  if (!dayChain.value || dayChain.value.checkedInToday) return
  try {
    await cApi.post('/api/c/community/daychain')
    dayChain.value.checkedInToday = true
    dayChain.value.currentStreak++
  } catch { /* TODO: toast */ }
}

function handleTaskAction(taskId: string) {
  // TODO: start/claim task via API
  console.log('Task action:', taskId)
}
</script>

<template>
  <div class="space-y-8">
    <!-- User Card (if connected) -->
    <div v-if="store.userStatus" class="bg-card-bg border border-border rounded-2xl p-6">
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 rounded-full bg-c-accent/20 flex items-center justify-center">
          <span class="material-symbols-rounded text-c-accent text-3xl">person</span>
        </div>
        <div class="flex-1">
          <div class="text-lg font-bold text-text-primary">{{ store.userStatus.walletAddress.slice(0, 6) }}...{{ store.userStatus.walletAddress.slice(-4) }}</div>
          <div class="flex items-center gap-4 mt-1 text-sm">
            <span class="text-c-accent font-semibold">{{ store.userStatus.xp.toLocaleString() }} XP</span>
            <span class="text-text-secondary">Level {{ store.userStatus.level }}</span>
            <span v-if="store.userStatus.dayStreak > 0" class="text-c-accent">🔥 {{ store.userStatus.dayStreak }} day streak</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Announcements -->
    <div v-if="announcements.length > 0" class="bg-card-bg border border-c-accent/30 rounded-2xl p-5 flex items-center gap-4">
      <span class="material-symbols-rounded text-c-accent text-2xl">campaign</span>
      <div class="flex-1">
        <div class="text-sm font-semibold text-text-primary">{{ announcements[0].title }}</div>
        <div class="text-xs text-text-secondary">{{ announcements[0].description }}</div>
      </div>
      <a v-if="announcements[0].link" :href="announcements[0].link" target="_blank" class="px-3 py-1.5 text-xs font-medium border border-c-accent text-c-accent rounded-lg hover:bg-c-accent/10">Learn More</a>
    </div>

    <!-- Daily Streak (DayChain) -->
    <div v-if="dayChain" class="bg-card-bg border border-border rounded-2xl p-6">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h2 class="text-base font-semibold text-text-primary">Daily Streak — Day {{ dayChain.currentStreak }} of {{ dayChain.targetDays }}</h2>
          <p class="text-xs text-text-secondary mt-0.5">Unlock streak: {{ dayChain.targetDays - dayChain.currentStreak }} days to go</p>
        </div>
        <button
          class="px-4 py-2 text-sm font-medium rounded-lg transition-colors"
          :class="dayChain.checkedInToday
            ? 'bg-text-muted/20 text-text-muted cursor-not-allowed'
            : 'bg-c-accent text-black hover:bg-c-accent/90'"
          :disabled="dayChain.checkedInToday"
          @click="checkIn"
        >
          {{ dayChain.checkedInToday ? 'Checked In ✓' : 'Continue Streak' }}
        </button>
      </div>
      <!-- Calendar strip (simplified) -->
      <div class="flex gap-1">
        <div
          v-for="i in Math.min(dayChain.targetDays, 30)"
          :key="i"
          class="w-3 h-3 rounded-sm"
          :class="i <= dayChain.currentStreak ? 'bg-status-active' : (i === dayChain.currentStreak + 1 ? 'bg-c-accent' : 'bg-border')"
        ></div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="flex gap-4 overflow-x-auto pb-2">
      <button class="shrink-0 flex items-center gap-2 px-4 py-3 bg-card-bg border border-border rounded-xl hover:border-c-accent/30 transition-colors" @click="router.push('/c/quests')">
        <div class="w-10 h-10 rounded-lg bg-c-accent/20 flex items-center justify-center">
          <span class="material-symbols-rounded text-c-accent">rocket_launch</span>
        </div>
        <span class="text-sm font-medium text-text-primary">Start Quest</span>
      </button>
      <button class="shrink-0 flex items-center gap-2 px-4 py-3 bg-card-bg border border-border rounded-xl hover:border-c-accent/30 transition-colors">
        <div class="w-10 h-10 rounded-lg bg-c-accent/20 flex items-center justify-center">
          <span class="material-symbols-rounded text-c-accent">casino</span>
        </div>
        <span class="text-sm font-medium text-text-primary">Lucky Wheel</span>
      </button>
      <button class="shrink-0 flex items-center gap-2 px-4 py-3 bg-card-bg border border-border rounded-xl hover:border-c-accent/30 transition-colors" @click="router.push('/c/invite')">
        <div class="w-10 h-10 rounded-lg bg-c-accent/20 flex items-center justify-center">
          <span class="material-symbols-rounded text-c-accent">group_add</span>
        </div>
        <span class="text-sm font-medium text-text-primary">Invite Friends</span>
      </button>
    </div>

    <!-- Task Sectors -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-20 bg-card-bg border border-border rounded-xl animate-pulse"></div>
    </div>
    <div v-else v-for="sector in sectors" :key="sector.id" class="space-y-3">
      <h2 class="text-base font-semibold text-text-primary">{{ sector.name }}</h2>
      <div class="space-y-2">
        <TaskCard
          v-for="task in sector.tasks"
          :key="task.id"
          :task="task"
          @action="handleTaskAction"
        />
      </div>
    </div>

    <!-- Community Pulse -->
    <div class="bg-card-bg border border-border rounded-2xl p-5">
      <div class="grid grid-cols-4 gap-4 text-center">
        <div>
          <div class="text-xl font-bold text-text-primary">{{ pulse.totalMembers.toLocaleString() }}</div>
          <div class="text-xs text-text-muted">Members</div>
        </div>
        <div>
          <div class="text-xl font-bold text-text-primary">{{ pulse.thisWeek.toLocaleString() }}</div>
          <div class="text-xs text-text-muted">This Week</div>
        </div>
        <div>
          <div class="text-xl font-bold text-text-primary">{{ pulse.liveActive.toLocaleString() }}</div>
          <div class="text-xs text-text-muted">Live Active</div>
        </div>
        <div>
          <div class="text-xl font-bold text-text-primary">{{ pulse.tasksDone.toLocaleString() }}</div>
          <div class="text-xs text-text-muted">Tasks Done</div>
        </div>
      </div>
    </div>

    <!-- Discover More -->
    <div>
      <h2 class="text-base font-semibold text-text-primary mb-3">Discover More</h2>
      <div class="grid grid-cols-2 gap-4">
        <button class="bg-card-bg border border-border rounded-xl p-5 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/lb-sprint')">
          <span class="material-symbols-rounded text-c-accent text-xl mb-2 block">sprint</span>
          <h3 class="text-sm font-semibold text-text-primary">Sprint Challenges</h3>
          <p class="text-xs text-text-secondary mt-1">Compete in time-limited leaderboard events</p>
        </button>
        <button class="bg-card-bg border border-border rounded-xl p-5 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/shop')">
          <span class="material-symbols-rounded text-c-accent text-xl mb-2 block">storefront</span>
          <h3 class="text-sm font-semibold text-text-primary">Rewards Shop</h3>
          <p class="text-xs text-text-secondary mt-1">Redeem your points for exclusive rewards</p>
        </button>
      </div>
    </div>
  </div>
</template>
