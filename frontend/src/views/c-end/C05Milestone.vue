<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { cApi } from '../../api/client'
import { useCEndStore } from '../../stores/c-end'

const router = useRouter()
const store = useCEndStore()
const loading = ref(true)

interface MilestoneItem {
  id: string
  name: string
  description: string
  requirement: string
  threshold: number
  reward: string
  reward_type: string
  status: 'earned' | 'claimable' | 'locked'
}

const milestones = ref<MilestoneItem[]>([])

const earned = ref<MilestoneItem[]>([])
const claimable = ref<MilestoneItem[]>([])
const locked = ref<MilestoneItem[]>([])

onMounted(async () => {
  try {
    const res = await cApi.get('/api/c/milestones')
    milestones.value = res.data.data || []
    earned.value = milestones.value.filter(m => m.status === 'earned')
    claimable.value = milestones.value.filter(m => m.status === 'claimable')
    locked.value = milestones.value.filter(m => m.status === 'locked')
  } finally { loading.value = false }
})

async function claimMilestone(id: string) {
  try {
    await cApi.post(`/api/c/milestones/${id}/claim`)
    const idx = milestones.value.findIndex(m => m.id === id)
    if (idx !== -1) milestones.value[idx].status = 'earned'
    claimable.value = milestones.value.filter(m => m.status === 'claimable')
    earned.value = milestones.value.filter(m => m.status === 'earned')
  } catch { /* TODO: toast */ }
}

const userPoints = store.userStatus?.xp || 0
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-text-primary mb-2">Milestones</h1>
      <div class="flex items-center gap-3">
        <span class="text-sm text-text-secondary">{{ userPoints.toLocaleString() }} points</span>
        <div class="flex-1 max-w-xs h-2 bg-page-bg rounded-full overflow-hidden">
          <div class="h-full bg-status-active rounded-full" style="width: 60%"></div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 4" :key="i" class="bg-card-bg border border-border rounded-xl h-20 animate-pulse"></div>
    </div>

    <template v-else>
      <!-- Earned -->
      <div v-if="earned.length > 0">
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Earned Milestones</div>
        <div class="space-y-2">
          <div v-for="m in earned" :key="m.id" class="bg-card-bg border border-status-active rounded-xl p-4 flex items-center gap-4">
            <span class="material-symbols-rounded text-xl text-status-active">check_circle</span>
            <div class="flex-1">
              <div class="text-sm font-medium text-text-primary line-through">{{ m.name }}</div>
              <div class="text-xs text-text-muted">Claimed: {{ m.reward }}</div>
            </div>
            <span class="text-xs text-status-active font-medium">Claimed ✓</span>
          </div>
        </div>
      </div>

      <!-- Claimable -->
      <div v-if="claimable.length > 0">
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Ready to Claim</div>
        <div class="space-y-2">
          <div v-for="m in claimable" :key="m.id" class="bg-card-bg border border-c-accent rounded-xl p-4 flex items-center gap-4 animate-pulse-slow">
            <span class="material-symbols-rounded text-xl text-c-accent">redeem</span>
            <div class="flex-1">
              <div class="text-sm font-semibold text-text-primary">{{ m.name }}</div>
              <div class="text-xs text-text-secondary">Reward: {{ m.reward }}</div>
            </div>
            <button class="px-4 py-2 bg-c-accent text-black text-sm font-medium rounded-lg hover:bg-c-accent/90" @click="claimMilestone(m.id)">
              Claim Reward
            </button>
          </div>
        </div>
      </div>

      <!-- Locked -->
      <div v-if="locked.length > 0">
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Locked</div>
        <div class="space-y-2">
          <div v-for="m in locked" :key="m.id" class="bg-card-bg border border-border rounded-xl p-4 flex items-center gap-4 opacity-60">
            <span class="material-symbols-rounded text-xl text-text-muted">lock</span>
            <div class="flex-1">
              <div class="text-sm font-medium text-text-primary">{{ m.name }}</div>
              <div class="text-xs text-text-muted">{{ m.requirement }}</div>
              <div class="text-xs text-text-secondary mt-0.5">Reward: {{ m.reward }}</div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Cross-sell -->
    <div class="grid grid-cols-2 gap-4">
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/lb-sprint')">
        <span class="text-sm text-text-primary font-medium">Join this week's Sprint to earn faster →</span>
      </button>
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/quests')">
        <span class="text-sm text-text-primary font-medium">Complete Quests for milestone points →</span>
      </button>
    </div>
  </div>
</template>
