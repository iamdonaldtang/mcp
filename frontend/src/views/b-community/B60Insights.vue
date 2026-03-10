<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { api } from '../../api/client'
import { useRouter } from 'vue-router'

// --- Types ---
interface MetricCard {
  label: string
  value: string
  trend: number
  trendLabel: string
  icon: string
}

interface SegmentData {
  key: string
  label: string
  icon: string
  color: string
  colorBg: string
  count: number
  percentage: number
  users: SegmentUser[]
}

interface SegmentUser {
  wallet: string
  last_active: string
  points: number
  actions: number
}

interface RetentionModule {
  key: string
  label: string
  retention_pct: number
  color: string
  route: string
}

interface EconomyData {
  months: string[]
  earned: number[]
  burned: number[]
  net: number[]
}

interface ModuleOption {
  key: string
  label: string
}

// --- State ---
const router = useRouter()
const loading = ref(true)

// Header controls
const dateRange = ref<'7d' | '30d' | '90d' | 'custom'>('30d')
const dateRangeOptions = [
  { value: '7d', label: 'Last 7 Days' },
  { value: '30d', label: 'Last 30 Days' },
  { value: '90d', label: 'Last 90 Days' },
  { value: 'custom', label: 'Custom' },
] as const

const moduleFilters = ref<string[]>([])
const moduleOptions = ref<ModuleOption[]>([])
const showModuleDropdown = ref(false)

// Key metrics
const metrics = ref<MetricCard[]>([
  { label: 'D30 Retention Rate', value: '—', trend: 0, trendLabel: '', icon: 'trending_up' },
  { label: 'Daily Active Users', value: '—', trend: 0, trendLabel: '', icon: 'group' },
  { label: 'Points Burn Rate', value: '—', trend: 0, trendLabel: '', icon: 'local_fire_department' },
  { label: 'Economy Balance', value: '—', trend: 0, trendLabel: '', icon: 'account_balance' },
])

// Economy chart
const economy = ref<EconomyData>({
  months: [],
  earned: [],
  burned: [],
  net: [],
})

// User segments
const segments = ref<SegmentData[]>([
  { key: 'power', label: 'Power Users', icon: 'bolt', color: '#48BB78', colorBg: '#0A2E1A', count: 0, percentage: 0, users: [] },
  { key: 'active', label: 'Active', icon: 'person', color: '#5D7EF1', colorBg: '#0F1A3D', count: 0, percentage: 0, users: [] },
  { key: 'at_risk', label: 'At Risk', icon: 'warning', color: '#F59E0B', colorBg: '#1F1A08', count: 0, percentage: 0, users: [] },
  { key: 'dormant', label: 'Dormant', icon: 'person_off', color: '#64748B', colorBg: '#1E293B', count: 0, percentage: 0, users: [] },
])
const expandedSegment = ref<string | null>(null)

// Retention by module
const retentionModules = ref<RetentionModule[]>([])

// Module route map
const moduleRouteMap: Record<string, string> = {
  tasks: '/b/community/sectors',
  points: '/b/community/points-level',
  leaderboard: '/b/community/leaderboard',
  daychain: '/b/community/daychain',
  lb_sprint: '/b/community/lb-sprint',
  badges: '/b/community/badges',
  milestones: '/b/community/milestones',
  shop: '/b/community/benefits-shop',
  wheel: '/b/community/lucky-wheel',
}

const moduleColorMap: Record<string, string> = {
  tasks: '#48BB78',
  points: '#5D7EF1',
  leaderboard: '#9B7EE0',
  daychain: '#F59E0B',
  lb_sprint: '#ED8936',
  badges: '#EC4899',
  milestones: '#14B8A6',
  shop: '#EF4444',
  wheel: '#8B5CF6',
}

// Max retention value for bar scaling
const maxRetention = computed(() => {
  if (!retentionModules.value.length) return 100
  return Math.max(...retentionModules.value.map(m => m.retention_pct), 1)
})

// Economy chart max for bar scaling
const economyMax = computed(() => {
  if (!economy.value.earned.length) return 1
  return Math.max(...economy.value.earned, ...economy.value.burned, 1)
})

// --- Data Fetching ---
async function fetchOverview() {
  try {
    const res = await api.get('/api/v1/community/insights/overview', {
      params: { range: dateRange.value, modules: moduleFilters.value.join(',') },
    })
    const d = res.data.data
    if (d) {
      metrics.value = [
        {
          label: 'D30 Retention Rate',
          value: `${d.retention_rate ?? 0}%`,
          trend: d.retention_trend ?? 0,
          trendLabel: `${Math.abs(d.retention_trend ?? 0)}% vs prev period`,
          icon: 'trending_up',
        },
        {
          label: 'Daily Active Users',
          value: formatNumber(d.dau ?? 0),
          trend: d.dau_trend ?? 0,
          trendLabel: `${Math.abs(d.dau_trend ?? 0)}% vs prev period`,
          icon: 'group',
        },
        {
          label: 'Points Burn Rate',
          value: `${d.burn_rate ?? 0}%`,
          trend: d.burn_rate_trend ?? 0,
          trendLabel: `${Math.abs(d.burn_rate_trend ?? 0)}% vs prev period`,
          icon: 'local_fire_department',
        },
        {
          label: 'Economy Balance',
          value: formatNumber(d.economy_balance ?? 0),
          trend: d.economy_trend ?? 0,
          trendLabel: `${Math.abs(d.economy_trend ?? 0)}% vs prev period`,
          icon: 'account_balance',
        },
      ]
    }
  } catch {
    /* use defaults */
  }
}

async function fetchEconomy() {
  try {
    const res = await api.get('/api/v1/community/insights/economy', {
      params: { range: dateRange.value, modules: moduleFilters.value.join(',') },
    })
    const d = res.data.data
    if (d) {
      economy.value = {
        months: d.months ?? [],
        earned: d.earned ?? [],
        burned: d.burned ?? [],
        net: d.net ?? [],
      }
    }
  } catch {
    // Fallback mock data for development
    economy.value = {
      months: ['Oct', 'Nov', 'Dec', 'Jan', 'Feb', 'Mar'],
      earned: [45000, 52000, 48000, 61000, 58000, 67000],
      burned: [32000, 38000, 41000, 44000, 47000, 51000],
      net: [13000, 14000, 7000, 17000, 11000, 16000],
    }
  }
}

async function fetchSegments() {
  try {
    const res = await api.get('/api/v1/community/insights/segments', {
      params: { range: dateRange.value },
    })
    const d = res.data.data
    if (d && Array.isArray(d.segments)) {
      d.segments.forEach((s: any) => {
        const seg = segments.value.find(x => x.key === s.key)
        if (seg) {
          seg.count = s.count ?? 0
          seg.percentage = s.percentage ?? 0
          seg.users = s.users ?? []
        }
      })
    }
  } catch {
    // Fallback mock
    segments.value[0].count = 142; segments.value[0].percentage = 8
    segments.value[1].count = 891; segments.value[1].percentage = 52
    segments.value[2].count = 423; segments.value[2].percentage = 25
    segments.value[3].count = 256; segments.value[3].percentage = 15
  }
}

async function fetchRetention() {
  try {
    const res = await api.get('/api/v1/community/insights/retention', {
      params: { range: dateRange.value, modules: moduleFilters.value.join(',') },
    })
    const d = res.data.data
    if (d && Array.isArray(d.modules)) {
      retentionModules.value = d.modules.map((m: any) => ({
        key: m.key,
        label: m.label,
        retention_pct: m.retention_pct ?? 0,
        color: moduleColorMap[m.key] ?? '#48BB78',
        route: moduleRouteMap[m.key] ?? '#',
      }))
    }
  } catch {
    // Fallback mock
    retentionModules.value = [
      { key: 'tasks', label: 'Sectors & Tasks', retention_pct: 78, color: '#48BB78', route: moduleRouteMap.tasks },
      { key: 'points', label: 'Points & Level', retention_pct: 72, color: '#5D7EF1', route: moduleRouteMap.points },
      { key: 'leaderboard', label: 'Leaderboard', retention_pct: 65, color: '#9B7EE0', route: moduleRouteMap.leaderboard },
      { key: 'daychain', label: 'DayChain', retention_pct: 61, color: '#F59E0B', route: moduleRouteMap.daychain },
      { key: 'lb_sprint', label: 'LB Sprint', retention_pct: 58, color: '#ED8936', route: moduleRouteMap.lb_sprint },
      { key: 'badges', label: 'Badges', retention_pct: 54, color: '#EC4899', route: moduleRouteMap.badges },
      { key: 'milestones', label: 'Milestones', retention_pct: 49, color: '#14B8A6', route: moduleRouteMap.milestones },
      { key: 'shop', label: 'Benefits Shop', retention_pct: 43, color: '#EF4444', route: moduleRouteMap.shop },
      { key: 'wheel', label: 'Lucky Wheel', retention_pct: 38, color: '#8B5CF6', route: moduleRouteMap.wheel },
    ]
  }
}

async function fetchModules() {
  try {
    const res = await api.get('/api/v1/community/modules')
    const d = res.data.data
    if (d && Array.isArray(d)) {
      moduleOptions.value = d.map((m: any) => ({ key: m.key, label: m.label }))
    }
  } catch {
    moduleOptions.value = [
      { key: 'tasks', label: 'Tasks' },
      { key: 'points', label: 'Points' },
      { key: 'leaderboard', label: 'Leaderboard' },
      { key: 'daychain', label: 'DayChain' },
      { key: 'lb_sprint', label: 'LB Sprint' },
      { key: 'badges', label: 'Badges' },
      { key: 'milestones', label: 'Milestones' },
      { key: 'shop', label: 'Shop' },
      { key: 'wheel', label: 'Wheel' },
    ]
  }
}

async function loadAll() {
  loading.value = true
  await Promise.all([
    fetchOverview(),
    fetchEconomy(),
    fetchSegments(),
    fetchRetention(),
    fetchModules(),
  ])
  loading.value = false
}

// --- Actions ---
function toggleModuleFilter(key: string) {
  const idx = moduleFilters.value.indexOf(key)
  if (idx >= 0) {
    moduleFilters.value.splice(idx, 1)
  } else {
    moduleFilters.value.push(key)
  }
}

function applyFilters() {
  showModuleDropdown.value = false
  loadAll()
}

function onDateRangeChange() {
  loadAll()
}

function toggleSegment(key: string) {
  expandedSegment.value = expandedSegment.value === key ? null : key
}

function navigateToModule(route: string) {
  router.push(route)
}

function exportCSV() {
  // Trigger CSV download
  window.open(`/api/v1/community/insights/export?format=csv&range=${dateRange.value}`, '_blank')
}

function exportPDF() {
  window.open(`/api/v1/community/insights/export?format=pdf&range=${dateRange.value}`, '_blank')
}

function truncateWallet(wallet: string): string {
  if (!wallet || wallet.length < 12) return wallet
  return `${wallet.slice(0, 6)}...${wallet.slice(-4)}`
}

function formatNumber(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`
  return n.toLocaleString()
}

function formatDate(d: string): string {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

// --- Init ---
onMounted(loadAll)
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="flex flex-col gap-4 mb-6 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">Community Insights</h1>
        <p class="text-sm text-text-secondary mt-1">Cross-module analytics, economy health, and user segments</p>
      </div>
      <div class="flex items-center gap-3 flex-wrap">
        <!-- Date Range Selector -->
        <select
          v-model="dateRange"
          class="h-9 px-3 bg-card-bg border border-border rounded-lg text-sm text-text-primary focus:outline-none focus:border-community"
          @change="onDateRangeChange"
        >
          <option v-for="opt in dateRangeOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>

        <!-- Module Filter -->
        <div class="relative">
          <button
            class="h-9 px-3 bg-card-bg border border-border rounded-lg text-sm text-text-secondary hover:border-community transition-colors flex items-center gap-2"
            @click="showModuleDropdown = !showModuleDropdown"
          >
            <span class="material-symbols-rounded text-base">filter_list</span>
            Modules
            <span v-if="moduleFilters.length" class="ml-1 px-1.5 py-0.5 bg-community/20 text-community text-xs rounded-full font-medium">
              {{ moduleFilters.length }}
            </span>
          </button>
          <div
            v-if="showModuleDropdown"
            class="absolute right-0 mt-1 w-56 bg-card-bg border border-border rounded-lg shadow-xl z-20 py-2"
          >
            <label
              v-for="mod in moduleOptions"
              :key="mod.key"
              class="flex items-center gap-2 px-3 py-1.5 hover:bg-white/5 cursor-pointer text-sm text-text-secondary"
            >
              <input
                type="checkbox"
                :checked="moduleFilters.includes(mod.key)"
                class="rounded border-border text-community focus:ring-community accent-community"
                @change="toggleModuleFilter(mod.key)"
              />
              {{ mod.label }}
            </label>
            <div class="border-t border-border mt-2 pt-2 px-3">
              <button
                class="w-full h-8 bg-community/10 text-community text-sm rounded-md hover:bg-community/20 transition-colors"
                @click="applyFilters"
              >
                Apply
              </button>
            </div>
          </div>
        </div>

        <!-- Export -->
        <button
          class="h-9 px-3 bg-card-bg border border-border rounded-lg text-sm text-text-secondary hover:border-community transition-colors flex items-center gap-1.5"
          @click="exportCSV"
        >
          <span class="material-symbols-rounded text-base">download</span>
          CSV
        </button>
        <button
          class="h-9 px-3 bg-card-bg border border-border rounded-lg text-sm text-text-secondary hover:border-community transition-colors flex items-center gap-1.5"
          @click="exportPDF"
        >
          <span class="material-symbols-rounded text-base">picture_as_pdf</span>
          PDF
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-32">
      <div class="animate-spin w-8 h-8 border-2 border-community border-t-transparent rounded-full" />
    </div>

    <template v-else>
      <!-- Key Metrics Bar -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div
          v-for="(m, i) in metrics"
          :key="i"
          class="bg-card-bg border border-border rounded-xl p-5"
        >
          <div class="flex items-center justify-between mb-3">
            <span class="text-xs text-text-muted font-medium uppercase tracking-wider">{{ m.label }}</span>
            <span class="material-symbols-rounded text-xl text-text-muted">{{ m.icon }}</span>
          </div>
          <div class="text-2xl font-bold text-text-primary mb-1">{{ m.value }}</div>
          <div class="flex items-center gap-1 text-xs">
            <span
              v-if="m.trend !== 0"
              class="material-symbols-rounded text-sm"
              :class="m.trend > 0 ? 'text-green-400' : 'text-red-400'"
            >
              {{ m.trend > 0 ? 'arrow_upward' : 'arrow_downward' }}
            </span>
            <span :class="m.trend > 0 ? 'text-green-400' : m.trend < 0 ? 'text-red-400' : 'text-text-muted'">
              {{ m.trendLabel }}
            </span>
          </div>
        </div>
      </div>

      <!-- Points Economy Chart -->
      <div class="bg-card-bg border border-border rounded-xl p-6 mb-6">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-semibold text-text-primary">Points Economy — 6 Month</h2>
          <div class="flex items-center gap-4 text-xs">
            <span class="flex items-center gap-1.5">
              <span class="w-3 h-3 rounded-sm bg-green-500" />
              <span class="text-text-secondary">Earned</span>
            </span>
            <span class="flex items-center gap-1.5">
              <span class="w-3 h-3 rounded-sm bg-red-500" />
              <span class="text-text-secondary">Burned</span>
            </span>
            <span class="flex items-center gap-1.5">
              <span class="w-3 h-3 rounded-sm bg-blue-500" />
              <span class="text-text-secondary">Net</span>
            </span>
          </div>
        </div>

        <!-- Bar Chart Placeholder -->
        <div class="relative h-64 flex items-end gap-2">
          <!-- Y-axis labels -->
          <div class="absolute left-0 top-0 bottom-0 flex flex-col justify-between text-[10px] text-text-muted w-10">
            <span>{{ formatNumber(economyMax) }}</span>
            <span>{{ formatNumber(Math.round(economyMax * 0.75)) }}</span>
            <span>{{ formatNumber(Math.round(economyMax * 0.5)) }}</span>
            <span>{{ formatNumber(Math.round(economyMax * 0.25)) }}</span>
            <span>0</span>
          </div>

          <!-- Grid lines -->
          <div class="absolute left-12 right-0 top-0 bottom-6 flex flex-col justify-between pointer-events-none">
            <div v-for="n in 5" :key="n" class="border-t border-border/50" />
          </div>

          <!-- Bars -->
          <div class="flex-1 flex items-end justify-around ml-12 pb-6 gap-1">
            <div
              v-for="(month, idx) in economy.months"
              :key="month"
              class="flex-1 flex items-end justify-center gap-1"
            >
              <!-- Earned bar -->
              <div
                class="w-4 rounded-t bg-linear-to-t from-green-600 to-green-400 transition-all duration-500"
                :style="{ height: `${(economy.earned[idx] / economyMax) * 220}px` }"
                :title="`Earned: ${formatNumber(economy.earned[idx])}`"
              />
              <!-- Burned bar -->
              <div
                class="w-4 rounded-t bg-linear-to-t from-red-600 to-red-400 transition-all duration-500"
                :style="{ height: `${(economy.burned[idx] / economyMax) * 220}px` }"
                :title="`Burned: ${formatNumber(economy.burned[idx])}`"
              />
              <!-- Net bar -->
              <div
                class="w-4 rounded-t bg-linear-to-t from-blue-600 to-blue-400 transition-all duration-500"
                :style="{ height: `${(economy.net[idx] / economyMax) * 220}px` }"
                :title="`Net: ${formatNumber(economy.net[idx])}`"
              />
            </div>
          </div>

          <!-- X-axis labels -->
          <div class="absolute left-12 right-0 bottom-0 flex justify-around text-[10px] text-text-muted">
            <span v-for="month in economy.months" :key="month">{{ month }}</span>
          </div>
        </div>
      </div>

      <!-- User Segments -->
      <div class="mb-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">User Segments</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div v-for="seg in segments" :key="seg.key">
            <!-- Segment Card -->
            <button
              class="w-full text-left bg-card-bg border rounded-xl p-5 transition-all duration-200 hover:border-opacity-80"
              :class="expandedSegment === seg.key ? 'border-opacity-100' : 'border-border'"
              :style="expandedSegment === seg.key ? { borderColor: seg.color } : {}"
              @click="toggleSegment(seg.key)"
            >
              <div class="flex items-center justify-between mb-3">
                <div
                  class="w-10 h-10 rounded-lg flex items-center justify-center"
                  :style="{ backgroundColor: seg.colorBg }"
                >
                  <span class="material-symbols-rounded text-xl" :style="{ color: seg.color }">
                    {{ seg.icon }}
                  </span>
                </div>
                <span class="material-symbols-rounded text-text-muted text-lg transition-transform duration-200"
                  :class="{ 'rotate-180': expandedSegment === seg.key }"
                >
                  keyboard_arrow_down
                </span>
              </div>
              <div class="text-xl font-bold text-text-primary">{{ formatNumber(seg.count) }}</div>
              <div class="flex items-center justify-between mt-1">
                <span class="text-sm text-text-secondary">{{ seg.label }}</span>
                <span class="text-xs font-medium px-2 py-0.5 rounded-full" :style="{ backgroundColor: seg.colorBg, color: seg.color }">
                  {{ seg.percentage }}%
                </span>
              </div>
            </button>

            <!-- Expanded Segment Detail -->
            <div
              v-if="expandedSegment === seg.key"
              class="mt-2 bg-card-bg border border-border rounded-xl overflow-hidden"
            >
              <table class="w-full text-xs">
                <thead>
                  <tr class="border-b border-border text-text-muted">
                    <th class="text-left py-2.5 px-3 font-medium">Wallet</th>
                    <th class="text-left py-2.5 px-3 font-medium">Last Active</th>
                    <th class="text-right py-2.5 px-3 font-medium">Points</th>
                    <th class="text-right py-2.5 px-3 font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(u, ui) in seg.users.slice(0, 5)"
                    :key="ui"
                    class="border-b border-border/50 hover:bg-white/5"
                  >
                    <td class="py-2 px-3 text-text-primary font-mono">{{ truncateWallet(u.wallet) }}</td>
                    <td class="py-2 px-3 text-text-secondary">{{ formatDate(u.last_active) }}</td>
                    <td class="py-2 px-3 text-right text-text-primary">{{ formatNumber(u.points) }}</td>
                    <td class="py-2 px-3 text-right text-text-secondary">{{ u.actions }}</td>
                  </tr>
                  <tr v-if="!seg.users.length">
                    <td colspan="4" class="py-4 text-center text-text-muted">No user data available</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Retention by Module -->
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-5">Retention by Module</h2>
        <div class="space-y-3">
          <button
            v-for="mod in retentionModules"
            :key="mod.key"
            class="w-full flex items-center gap-4 group hover:bg-white/5 rounded-lg px-3 py-2 -mx-3 transition-colors"
            @click="navigateToModule(mod.route)"
          >
            <span class="text-sm text-text-secondary w-32 text-left shrink-0">{{ mod.label }}</span>
            <div class="flex-1 h-6 bg-border/30 rounded-full overflow-hidden relative">
              <div
                class="h-full rounded-full transition-all duration-700 ease-out"
                :style="{
                  width: `${(mod.retention_pct / maxRetention) * 100}%`,
                  backgroundColor: mod.color,
                  opacity: 0.8,
                }"
              />
            </div>
            <span class="text-sm font-semibold text-text-primary w-12 text-right shrink-0">
              {{ mod.retention_pct }}%
            </span>
            <span class="material-symbols-rounded text-text-muted text-sm opacity-0 group-hover:opacity-100 transition-opacity">
              arrow_forward
            </span>
          </button>
          <div v-if="!retentionModules.length" class="text-center text-text-muted py-8">
            No retention data available
          </div>
        </div>
      </div>
    </template>

    <!-- Click-away for module dropdown -->
    <div
      v-if="showModuleDropdown"
      class="fixed inset-0 z-10"
      @click="showModuleDropdown = false"
    />
  </div>
</template>
