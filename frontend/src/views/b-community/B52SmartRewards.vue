<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()
const loading = ref(true)

const overview = ref({
  rules: { total: 0, active: 0 },
  privileges: { total: 0, active: 0 },
})

onMounted(async () => {
  loading.value = true
  try {
    const res = await api.get('/api/v1/community/smart-rewards/overview')
    if (res.data.data) {
      overview.value = res.data.data
    }
  } catch { /* defaults */ }
  finally {
    loading.value = false
  }
})

const cards = [
  {
    key: 'rules',
    icon: 'tune',
    title: 'Reward Rules',
    description: 'Create automated rules that trigger rewards based on user actions',
    route: '/b/community/smart-rewards/rules',
    buttonLabel: 'Manage Rules →',
    color: '#48BB78',
    bgColor: '#0A2E1A',
  },
  {
    key: 'privileges',
    icon: 'shield',
    title: 'Privileges',
    description: 'Define special privileges and perks for different user tiers',
    route: '/b/community/smart-rewards/privileges',
    buttonLabel: 'Manage Privileges →',
    color: '#9B7EE0',
    bgColor: '#1E1033',
  },
]
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-text-primary mb-1">Smart Rewards</h1>
      <p class="text-sm text-text-secondary">Automate your reward distribution with rules and privileges</p>
    </div>

    <!-- Cards -->
    <div class="grid grid-cols-2 gap-6">
      <template v-if="loading">
        <div v-for="i in 2" :key="i" class="bg-card-bg border border-border rounded-2xl p-8 h-64 animate-pulse" />
      </template>
      <div
        v-else
        v-for="card in cards"
        :key="card.key"
        class="bg-card-bg border border-border rounded-2xl p-8 hover:border-community/30 transition-colors flex flex-col"
      >
        <!-- Icon -->
        <div
          class="w-14 h-14 rounded-xl flex items-center justify-center mb-5"
          :style="{ background: card.bgColor }"
        >
          <span class="material-symbols-rounded text-2xl" :style="{ color: card.color }">{{ card.icon }}</span>
        </div>

        <!-- Title & Description -->
        <h2 class="text-lg font-semibold text-text-primary mb-2">{{ card.title }}</h2>
        <p class="text-sm text-text-secondary mb-5 flex-1">{{ card.description }}</p>

        <!-- Stats -->
        <div class="flex items-center gap-6 mb-6">
          <div>
            <div class="text-2xl font-bold text-text-primary">
              {{ card.key === 'rules' ? overview.rules.total : overview.privileges.total }}
            </div>
            <div class="text-xs text-text-muted">Total {{ card.key === 'rules' ? 'Rules' : 'Privileges' }}</div>
          </div>
          <div class="w-px h-8 bg-border"></div>
          <div>
            <div class="text-2xl font-bold" :style="{ color: card.color }">
              {{ card.key === 'rules' ? overview.rules.active : overview.privileges.active }}
            </div>
            <div class="text-xs text-text-muted">Active</div>
          </div>
        </div>

        <!-- Button -->
        <button
          class="w-full py-3 text-sm font-medium rounded-xl transition-colors"
          :style="{
            background: card.bgColor,
            color: card.color,
          }"
          @mouseenter="($event.target as HTMLElement).style.opacity = '0.85'"
          @mouseleave="($event.target as HTMLElement).style.opacity = '1'"
          @click="router.push(card.route)"
        >
          {{ card.buttonLabel }}
        </button>
      </div>
    </div>

    <!-- Info Note -->
    <div class="bg-card-bg border border-border rounded-xl p-5 flex items-start gap-3">
      <span class="material-symbols-rounded text-lg text-text-muted mt-0.5">info</span>
      <div>
        <p class="text-sm text-text-secondary">
          Smart Rewards lets you create automated reward logic for your community.
          <strong class="text-text-primary">Rules</strong> define conditions and triggers (e.g., "award 50 points when a user completes 5 tasks"),
          while <strong class="text-text-primary">Privileges</strong> define tier-based perks (e.g., "Gold members get 2x points").
        </p>
      </div>
    </div>
  </div>
</template>
