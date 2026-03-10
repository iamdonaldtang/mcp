<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCEndStore } from '../stores/c-end'

const route = useRoute()
const router = useRouter()
const store = useCEndStore()

const isConnected = computed(() => store.isWalletConnected)
const userStatus = computed(() => store.userStatus)
const tabVisibility = computed(() => store.tabVisibility)

interface NavTab {
  label: string
  route: string
  key: keyof typeof tabVisibility.value
}

const navTabs: NavTab[] = [
  { label: 'Home', route: '/c', key: 'home' },
  { label: 'Quests', route: '/c/quests', key: 'quests' },
  { label: 'Leaderboard', route: '/c/leaderboard', key: 'leaderboard' },
  { label: 'LB Sprint', route: '/c/lb-sprint', key: 'lbSprint' },
  { label: 'Milestone', route: '/c/milestones', key: 'milestone' },
  { label: 'Shop', route: '/c/shop', key: 'shop' },
]

const visibleTabs = computed(() =>
  navTabs.filter(tab => tabVisibility.value[tab.key])
)

function isActive(tabRoute: string) {
  return route.path === tabRoute
}

function connectWallet() {
  store.connectWallet()
}

function goProfile() {
  router.push('/c/profile')
}

function truncateAddress(addr: string) {
  return addr.slice(0, 6) + '...' + addr.slice(-4)
}
</script>

<template>
  <div class="min-h-screen bg-page-bg flex flex-col">
    <!-- Header -->
    <header class="bg-header-bg border-b border-border">
      <div class="max-w-7xl mx-auto px-6 h-14 flex items-center">
        <!-- Left: Logo + Project Name -->
        <div class="flex items-center gap-2 mr-8">
          <div class="w-8 h-8 rounded-lg bg-c-accent/20 flex items-center justify-center">
            <span class="material-symbols-rounded text-c-accent text-lg">bolt</span>
          </div>
          <span class="text-text-primary font-semibold">{{ store.communityName || 'Community' }}</span>
        </div>

        <!-- Center: Nav Tabs -->
        <nav class="flex-1 flex items-center gap-1">
          <button
            v-for="tab in visibleTabs"
            :key="tab.route"
            class="px-4 py-4 text-sm font-medium border-b-2 transition-colors"
            :class="isActive(tab.route)
              ? 'text-c-accent border-c-accent'
              : 'text-text-secondary border-transparent hover:text-text-primary'"
            @click="router.push(tab.route)"
          >
            {{ tab.label }}
          </button>
        </nav>

        <!-- Right: Wallet / User -->
        <div class="flex items-center gap-3">
          <button class="text-text-secondary hover:text-text-primary transition-colors">
            <span class="material-symbols-rounded text-xl">notifications</span>
          </button>
          <button
            v-if="!isConnected"
            class="px-4 py-1.5 text-sm font-medium border border-c-accent text-c-accent rounded-lg hover:bg-c-accent/10 transition-colors"
            @click="connectWallet"
          >
            Connect Wallet
          </button>
          <div v-else class="flex items-center gap-2 cursor-pointer" @click="goProfile">
            <div class="w-8 h-8 rounded-full bg-c-accent/20 flex items-center justify-center">
              <span class="material-symbols-rounded text-c-accent text-sm">person</span>
            </div>
            <span class="text-sm text-text-primary">{{ truncateAddress(userStatus?.walletAddress || '') }}</span>
            <button class="text-text-muted hover:text-text-primary" @click.stop="router.push('/c/profile')">
              <span class="material-symbols-rounded text-base">settings</span>
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- User Status Bar (logged in only) -->
    <div v-if="isConnected && userStatus" class="bg-card-bg border-b border-border">
      <div class="max-w-7xl mx-auto px-6 h-10 flex items-center gap-6 text-xs">
        <span class="text-text-primary">{{ truncateAddress(userStatus.walletAddress) }}</span>
        <span class="px-2 py-0.5 rounded bg-c-accent/20 text-c-accent font-semibold">Lv.{{ userStatus.level }}</span>
        <span class="text-text-secondary">{{ userStatus.xp.toLocaleString() }} XP</span>
        <span v-if="userStatus.dayStreak > 0" class="text-c-accent">🔥 {{ userStatus.dayStreak }} day streak</span>
      </div>
    </div>

    <!-- Page Content -->
    <main class="flex-1 max-w-7xl mx-auto w-full px-6 py-8">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="border-t border-border bg-page-bg">
      <div class="max-w-7xl mx-auto px-6 py-6">
        <div class="flex items-center justify-between text-xs text-text-muted">
          <span>TaskOn · Earn rewards by completing Web3 tasks</span>
          <div class="flex items-center gap-4">
            <a href="#" class="hover:text-text-secondary transition-colors">Help Center</a>
            <a href="#" class="hover:text-text-secondary transition-colors">Docs</a>
            <a href="#" class="hover:text-text-secondary transition-colors">Community</a>
            <a href="#" class="hover:text-text-secondary transition-colors">Terms</a>
            <a href="#" class="hover:text-text-secondary transition-colors">Privacy</a>
          </div>
          <a href="https://taskon.xyz" target="_blank" class="hover:text-text-secondary transition-colors">
            Powered by TaskOn
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>
