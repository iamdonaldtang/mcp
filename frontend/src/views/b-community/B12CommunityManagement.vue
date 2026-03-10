<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'

// === Types ===
interface ModuleInfo {
  id: string
  name: string
  icon: string
  status: 'active' | 'inactive'
  metric_label: string
  metric_value: string | number
}

interface AiInsight {
  id: string
  type: 'warning' | 'insight' | 'success'
  icon: string
  title: string
  description: string
  action_label: string
  action_route: string
}

interface IntegrationCard {
  key: string
  name: string
  icon: string
  connected: boolean
}

interface RetentionMetrics {
  d1: number
  d7: number
  d30: number
  mau: number
}

interface OverviewData {
  stats: {
    total_members: number
    members_trend: string
    active_this_week: number
    active_trend: string
    points_distributed: number
    points_trend: string
    tasks_completed: number
    tasks_trend: string
  }
  modules: ModuleInfo[]
  insights: AiInsight[]
  integrations: IntegrationCard[]
  retention: RetentionMetrics
}

// === State ===
const loading = ref(true)
const data = ref<OverviewData>({
  stats: {
    total_members: 0,
    members_trend: '',
    active_this_week: 0,
    active_trend: '',
    points_distributed: 0,
    points_trend: '',
    tasks_completed: 0,
    tasks_trend: '',
  },
  modules: [],
  insights: [],
  integrations: [],
  retention: { d1: 0, d7: 0, d30: 0, mau: 0 },
})

// === Insight styling ===
const insightStyles: Record<string, { bg: string; border: string; icon_color: string }> = {
  warning: { bg: 'rgba(245,158,11,0.08)', border: '#D97706', icon_color: '#F59E0B' },
  insight: { bg: 'rgba(59,130,246,0.08)', border: '#3B82F6', icon_color: '#60A5FA' },
  success: { bg: 'rgba(72,187,120,0.08)', border: '#48BB78', icon_color: '#48BB78' },
}

// === API ===
onMounted(async () => {
  await fetchOverview()
  loading.value = false
})

async function fetchOverview() {
  try {
    const res = await api.get('/api/v1/community/overview', { params: { view: 'deep' } })
    if (res.data.data) {
      data.value = res.data.data
    }
  } catch { /* TODO: error toast */ }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-text-primary mb-1">Community Management</h1>
      <p class="text-sm text-text-secondary">Advanced overview for your community — modules, insights, integrations, and engagement</p>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard
        label="Total Members"
        :value="data.stats.total_members"
        icon="group"
        icon-color="#48BB78"
        :trend="data.stats.members_trend"
      />
      <StatsCard
        label="Active This Week"
        :value="data.stats.active_this_week"
        icon="trending_up"
        icon-color="#60A5FA"
        :trend="data.stats.active_trend"
      />
      <StatsCard
        label="Points Distributed"
        :value="data.stats.points_distributed"
        icon="stars"
        icon-color="#F59E0B"
        :trend="data.stats.points_trend"
      />
      <StatsCard
        label="Tasks Completed"
        :value="data.stats.tasks_completed"
        icon="task_alt"
        icon-color="#A78BFA"
        :trend="data.stats.tasks_trend"
      />
    </div>

    <!-- Active Modules -->
    <div>
      <h2 class="text-lg font-semibold text-text-primary mb-3">Active Modules</h2>
      <div v-if="loading" class="grid grid-cols-4 gap-4">
        <div v-for="i in 4" :key="i" class="bg-card-bg border border-border rounded-xl p-4 h-24 animate-pulse" />
      </div>
      <div v-else-if="data.modules.length === 0" class="bg-card-bg border border-border rounded-xl p-6 text-center">
        <p class="text-sm text-text-muted">No modules activated yet.</p>
      </div>
      <div v-else class="grid grid-cols-4 gap-4">
        <div
          v-for="mod in data.modules"
          :key="mod.id"
          class="bg-card-bg border border-border rounded-xl p-4 hover:border-community/30 transition-colors"
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <span class="material-symbols-rounded text-lg text-text-secondary">{{ mod.icon }}</span>
              <span class="text-sm font-medium text-text-primary">{{ mod.name }}</span>
            </div>
            <span
              class="w-2 h-2 rounded-full"
              :class="mod.status === 'active' ? 'bg-status-active' : 'bg-[#64748B]'"
            />
          </div>
          <div class="text-xs text-text-muted">{{ mod.metric_label }}</div>
          <div class="text-lg font-bold text-text-primary mt-0.5">
            {{ typeof mod.metric_value === 'number' ? mod.metric_value.toLocaleString() : mod.metric_value }}
          </div>
        </div>
      </div>
    </div>

    <!-- AI Insights -->
    <div>
      <h2 class="text-lg font-semibold text-text-primary mb-3">AI Insights</h2>
      <div v-if="loading" class="grid grid-cols-3 gap-4">
        <div v-for="i in 3" :key="i" class="bg-card-bg border border-border rounded-xl p-5 h-32 animate-pulse" />
      </div>
      <div v-else-if="data.insights.length === 0" class="bg-card-bg border border-border rounded-xl p-6 text-center">
        <p class="text-sm text-text-muted">No insights available right now. Check back later.</p>
      </div>
      <div v-else class="grid grid-cols-3 gap-4">
        <div
          v-for="insight in data.insights"
          :key="insight.id"
          class="rounded-xl p-5 border"
          :style="{
            backgroundColor: insightStyles[insight.type]?.bg || insightStyles.insight.bg,
            borderColor: insightStyles[insight.type]?.border || insightStyles.insight.border,
          }"
        >
          <div class="flex items-start gap-3">
            <span
              class="material-symbols-rounded text-xl mt-0.5 shrink-0"
              :style="{ color: insightStyles[insight.type]?.icon_color || '#60A5FA' }"
            >
              {{ insight.icon }}
            </span>
            <div class="min-w-0">
              <h3 class="text-sm font-semibold text-text-primary mb-1">{{ insight.title }}</h3>
              <p class="text-xs text-text-secondary leading-relaxed mb-3">{{ insight.description }}</p>
              <router-link
                :to="insight.action_route"
                class="text-xs font-medium hover:underline"
                :style="{ color: insightStyles[insight.type]?.icon_color || '#60A5FA' }"
              >
                {{ insight.action_label }} &rarr;
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Integrations -->
    <div>
      <h2 class="text-lg font-semibold text-text-primary mb-3">Integrations</h2>
      <div class="grid grid-cols-4 gap-4">
        <template v-if="loading">
          <div v-for="i in 4" :key="i" class="bg-card-bg border border-border rounded-xl p-4 h-20 animate-pulse" />
        </template>
        <template v-else>
          <div
            v-for="intg in data.integrations"
            :key="intg.key"
            class="bg-card-bg border border-border rounded-xl p-4 flex items-center justify-between hover:border-community/30 transition-colors cursor-pointer"
          >
            <div class="flex items-center gap-3">
              <span class="material-symbols-rounded text-xl text-text-secondary">{{ intg.icon }}</span>
              <span class="text-sm font-medium text-text-primary">{{ intg.name }}</span>
            </div>
            <span
              class="text-xs font-medium px-2 py-0.5 rounded-full"
              :class="intg.connected
                ? 'bg-status-active-bg text-status-active'
                : 'bg-[#1E293B] text-[#64748B]'"
            >
              {{ intg.connected ? 'Connected' : 'Available' }}
            </span>
          </div>
          <!-- All Integrations link -->
          <router-link
            to="/community/integrations"
            class="bg-card-bg border border-border border-dashed rounded-xl p-4 flex items-center justify-center gap-2 hover:border-community/30 transition-colors"
          >
            <span class="material-symbols-rounded text-lg text-text-muted">integration_instructions</span>
            <span class="text-sm text-text-secondary">All Integrations &rarr;</span>
          </router-link>
        </template>
      </div>
    </div>

    <!-- Engagement Overview -->
    <div class="grid grid-cols-2 gap-4">
      <!-- WAU Chart Placeholder -->
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <h3 class="text-sm font-semibold text-text-primary mb-4">Weekly Active Users</h3>
        <div class="h-48 flex items-center justify-center border border-dashed border-border rounded-lg">
          <div class="text-center">
            <span class="material-symbols-rounded text-3xl text-text-muted mb-1 block">bar_chart</span>
            <span class="text-sm text-text-muted">WAU Chart</span>
          </div>
        </div>
      </div>

      <!-- Retention Metrics -->
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <h3 class="text-sm font-semibold text-text-primary mb-4">Retention Metrics</h3>
        <div class="space-y-4">
          <!-- D1 Retention -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">D1 Retention</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.retention.d1 }}%</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-community transition-all duration-500" :style="{ width: data.retention.d1 + '%' }" />
            </div>
          </div>
          <!-- 7-Day Retention -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">7-Day Retention</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.retention.d7 }}%</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-[#60A5FA] transition-all duration-500" :style="{ width: data.retention.d7 + '%' }" />
            </div>
          </div>
          <!-- 30-Day Retention -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">30-Day Retention</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.retention.d30 }}%</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-c-accent transition-all duration-500" :style="{ width: data.retention.d30 + '%' }" />
            </div>
          </div>
          <!-- MAU -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">MAU</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.retention.mau }}%</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-[#A78BFA] transition-all duration-500" :style="{ width: data.retention.mau + '%' }" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
