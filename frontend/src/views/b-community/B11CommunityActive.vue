<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'

const router = useRouter()
const loading = ref(true)

const stats = ref({
  totalMembers: 0,
  activeThisWeek: 0,
  pointsDistributed: 0,
  tasksCompleted: 0,
  trends: { members: 0, active: 0, points: 0, tasks: 0 },
})

const checklistDismissed = ref(false)
const checklistRemaining = ref(0)

interface ModuleStat {
  type: string
  name: string
  icon: string
  metrics: { label: string; value: string | number }[]
  route: string
}

const moduleStats = ref<ModuleStat[]>([])

interface AvailableModule {
  type: string
  name: string
  icon: string
  description: string
}

const availableModules = ref<AvailableModule[]>([])

const moduleIcons: Record<string, string> = {
  sectors_tasks: 'category', task_chain: 'link', day_chain: 'local_fire_department',
  points_level: 'stars', leaderboard: 'leaderboard', badges: 'military_tech',
  lb_sprint: 'sprint', milestone: 'flag', lucky_wheel: 'casino', benefits_shop: 'storefront',
}

const moduleNames: Record<string, string> = {
  sectors_tasks: 'Sectors & Tasks', task_chain: 'TaskChain', day_chain: 'DayChain',
  points_level: 'Points & Level', leaderboard: 'Leaderboard', badges: 'Badges',
  lb_sprint: 'LB Sprint', milestone: 'Milestones', lucky_wheel: 'Lucky Wheel', benefits_shop: 'Benefits Shop',
}

const moduleRoutes: Record<string, string> = {
  sectors_tasks: '/b/community/modules/sectors', task_chain: '/b/community/modules/task-chain',
  day_chain: '/b/community/modules/day-chain', points_level: '/b/community/modules/points-level',
  leaderboard: '/b/community/modules/leaderboard', badges: '/b/community/modules/badges',
  lb_sprint: '/b/community/modules/lb-sprint', milestone: '/b/community/modules/milestone',
  lucky_wheel: '/b/community/modules/lucky-wheel', benefits_shop: '/b/community/modules/benefits-shop',
}

const allModuleTypes = Object.keys(moduleNames)
const showAddModules = computed(() => availableModules.value.length > 0)

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([fetchStats(), fetchModules(), fetchChecklist()])
  } finally {
    loading.value = false
  }
})

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/stats')
    if (res.data.data) stats.value = res.data.data
  } catch { /* defaults */ }
}

async function fetchModules() {
  try {
    const res = await api.get('/api/v1/community/overview')
    const overview = res.data.data
    if (!overview?.enabled_modules) return

    const enabled = overview.enabled_modules as string[]
    moduleStats.value = enabled.map((type: string) => ({
      type,
      name: moduleNames[type] || type,
      icon: moduleIcons[type] || 'extension',
      metrics: [{ label: 'Activity', value: '—' }, { label: 'Trend', value: '—' }],
      route: moduleRoutes[type] || '/b/community/overview',
    }))

    availableModules.value = allModuleTypes
      .filter(t => !enabled.includes(t))
      .map(type => ({
        type,
        name: moduleNames[type] || type,
        icon: moduleIcons[type] || 'extension',
        description: `Enable ${moduleNames[type]} to expand your community.`,
      }))
  } catch { /* defaults */ }
}

async function fetchChecklist() {
  try {
    const res = await api.get('/api/v1/community/checklist')
    const items = res.data.data || []
    checklistRemaining.value = items.filter((i: any) => !i.completed).length
    if (checklistRemaining.value === 0) checklistDismissed.value = true
  } catch {
    checklistDismissed.value = true
  }
}

async function enableModule(type: string) {
  try {
    await api.put(`/api/v1/community/modules/${type}/enable`)
    await fetchModules()
  } catch { /* TODO: toast */ }
}

function dismissChecklist() {
  checklistDismissed.value = true
  api.put('/api/v1/community/onboarding/dismiss').catch(() => {})
}

function formatTrend(val: number) {
  if (val > 0) return `↑${val}%`
  if (val < 0) return `↓${Math.abs(val)}%`
  return '—'
}
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <div class="flex items-center gap-3 mb-1">
        <h1 class="text-2xl font-bold text-text-primary">My Community</h1>
        <span class="px-2.5 py-0.5 text-xs font-medium rounded-full" style="background: #0A2E1A; color: #16A34A">Active</span>
      </div>
      <p class="text-sm text-text-secondary">
        {{ stats.totalMembers.toLocaleString() }} members · {{ moduleStats.length }} modules active
      </p>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard label="Total Members" :value="stats.totalMembers" icon="group" icon-color="#48BB78" :trend="formatTrend(stats.trends.members)" />
      <StatsCard label="Active This Week" :value="stats.activeThisWeek" icon="trending_up" icon-color="#3B82F6" :trend="formatTrend(stats.trends.active)" />
      <StatsCard label="Points Distributed" :value="stats.pointsDistributed" icon="stars" icon-color="#F59E0B" :trend="formatTrend(stats.trends.points)" />
      <StatsCard label="Tasks Completed" :value="stats.tasksCompleted" icon="task_alt" icon-color="#9B7EE0" :trend="formatTrend(stats.trends.tasks)" />
    </div>

    <!-- Checklist Banner -->
    <div
      v-if="!checklistDismissed && checklistRemaining > 0"
      class="bg-card-bg border border-border rounded-xl p-4 flex items-center justify-between cursor-pointer hover:bg-white/2"
      @click="router.push('/b/community/guided')"
    >
      <div class="flex items-center gap-3">
        <span class="material-symbols-rounded text-xl" style="color: #D97706">checklist</span>
        <span class="text-sm text-text-primary">{{ checklistRemaining }} steps remaining · Complete your community setup</span>
      </div>
      <button class="text-text-muted hover:text-text-primary transition-colors p-1" @click.stop="dismissChecklist">
        <span class="material-symbols-rounded text-lg">close</span>
      </button>
    </div>

    <!-- Module Performance -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Module Performance</div>
      <div class="grid grid-cols-3 gap-4">
        <template v-if="loading">
          <div v-for="i in 3" :key="i" class="bg-card-bg border border-border rounded-xl p-5 h-32 animate-pulse" />
        </template>
        <div
          v-else
          v-for="mod in moduleStats"
          :key="mod.type"
          class="bg-card-bg border border-border rounded-xl p-5 cursor-pointer hover:bg-[#161F2E] transition-colors"
          @click="router.push(mod.route)"
        >
          <div class="flex items-center gap-2 mb-3">
            <span class="material-symbols-rounded text-community">{{ mod.icon }}</span>
            <span class="text-sm font-semibold text-text-primary">{{ mod.name }}</span>
          </div>
          <div class="space-y-1">
            <div v-for="m in mod.metrics" :key="m.label" class="flex justify-between text-xs">
              <span class="text-text-muted">{{ m.label }}</span>
              <span class="text-text-secondary">{{ m.value }}</span>
            </div>
          </div>
          <button class="mt-3 text-xs text-community font-medium hover:underline" @click.stop="router.push(mod.route)">Manage →</button>
        </div>
      </div>
    </div>

    <!-- Add More Modules -->
    <div v-if="showAddModules">
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Add More Modules</div>
      <div class="grid grid-cols-4 gap-4">
        <div v-for="mod in availableModules" :key="mod.type" class="bg-card-bg border border-border rounded-xl p-4">
          <span class="material-symbols-rounded text-text-muted text-xl mb-2 block">{{ mod.icon }}</span>
          <h4 class="text-sm font-medium text-text-primary mb-1">{{ mod.name }}</h4>
          <p class="text-xs text-text-muted mb-3">{{ mod.description }}</p>
          <button class="text-xs text-community font-medium hover:underline" @click="enableModule(mod.type)">+ Enable</button>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Quick Actions</div>
      <div class="flex gap-4">
        <button class="flex items-center gap-2 px-5 py-3 bg-card-bg border border-border rounded-xl text-sm font-medium text-text-primary hover:bg-white/2 transition-colors" @click="router.push('/b/community/modules/sectors')">
          <span class="material-symbols-rounded text-community text-lg">add_circle</span> Create Task
        </button>
        <button class="flex items-center gap-2 px-5 py-3 bg-card-bg border border-border rounded-xl text-sm font-medium text-text-primary hover:bg-white/2 transition-colors" @click="router.push('/b/community/modules/benefits-shop')">
          <span class="material-symbols-rounded text-community text-lg">redeem</span> Add Reward
        </button>
        <button class="flex items-center gap-2 px-5 py-3 bg-card-bg border border-border rounded-xl text-sm font-medium text-text-primary hover:bg-white/2 transition-colors" @click="router.push('/b/community/insights')">
          <span class="material-symbols-rounded text-community text-lg">analytics</span> View Analytics
        </button>
      </div>
    </div>

    <div class="border-t border-border"></div>

    <!-- Resources -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Resources</div>
      <div class="grid grid-cols-2 gap-4">
        <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-xl mb-2 block">menu_book</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Community Playbook</h4>
          <p class="text-xs text-text-secondary">Best practices for growing your Web3 community</p>
        </a>
        <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-xl mb-2 block">psychology</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Points Strategy Guide</h4>
          <p class="text-xs text-text-secondary">Design an effective points economy</p>
        </a>
      </div>
    </div>
  </div>
</template>
