<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'

const router = useRouter()
const loading = ref(true)

// === Types ===
interface Deployment {
  id: string
  type: 'domain' | 'widget' | 'page'
  label: string
  url: string
  status: 'live' | 'draft' | 'paused'
  last_activity: string
}

interface ConfiguredTool {
  key: string
  name: string
  icon: string
  status: 'active' | 'configured' | 'inactive'
  route: string
}

interface ActivityEvent {
  id: string
  icon: string
  description: string
  timestamp: string
  type: 'deployment' | 'config' | 'user' | 'system'
}

interface OverviewData {
  stats: {
    total_impressions: number
    impressions_trend: string
    active_deployments: number
    deployments_trend: string
    user_sessions: number
    sessions_trend: string
    revenue_impact: string
    revenue_trend: string
  }
  deployments: Deployment[]
  tools: ConfiguredTool[]
  activity: ActivityEvent[]
}

// === State ===
const data = ref<OverviewData>({
  stats: {
    total_impressions: 0,
    impressions_trend: '',
    active_deployments: 0,
    deployments_trend: '',
    user_sessions: 0,
    sessions_trend: '',
    revenue_impact: '$0',
    revenue_trend: '',
  },
  deployments: [],
  tools: [],
  activity: [],
})

// Fallback tools when API returns empty
const defaultTools: ConfiguredTool[] = [
  { key: 'domain', name: 'Custom Domain', icon: 'language', status: 'active', route: '/b/whitelabel/domain' },
  { key: 'widgets', name: 'Widget Library', icon: 'widgets', status: 'active', route: '/b/whitelabel/widgets' },
  { key: 'pages', name: 'Page Builder', icon: 'web', status: 'active', route: '/b/whitelabel/pages' },
  { key: 'sdk', name: 'SDK & API', icon: 'code', status: 'configured', route: '/b/whitelabel/sdk' },
  { key: 'brand', name: 'Brand Config', icon: 'palette', status: 'active', route: '/b/whitelabel/brand' },
  { key: 'integrations', name: 'Integrations', icon: 'integration_instructions', status: 'active', route: '/b/whitelabel/integrations' },
]

const tools = computed(() => data.value.tools.length > 0 ? data.value.tools : defaultTools)

// Deployment status styling
const deploymentStatusStyles: Record<string, { bg: string; text: string }> = {
  live: { bg: '#0A2E1A', text: '#16A34A' },
  draft: { bg: '#1F1A08', text: '#D97706' },
  paused: { bg: '#2D1515', text: '#DC2626' },
}

const deploymentTypeIcons: Record<string, string> = {
  domain: 'language',
  widget: 'widgets',
  page: 'web',
}

function toolStatusDot(status: string): string {
  switch (status) {
    case 'active': return '#16A34A'
    case 'configured': return '#D97706'
    default: return '#64748B'
  }
}

function activityIconColor(type: string): string {
  switch (type) {
    case 'deployment': return '#9B7EE0'
    case 'config': return '#60A5FA'
    case 'user': return '#48BB78'
    default: return '#94A3B8'
  }
}

function formatTimestamp(ts: string): string {
  if (!ts) return ''
  try {
    const date = new Date(ts)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffMins = Math.floor(diffMs / 60000)
    if (diffMins < 1) return 'Just now'
    if (diffMins < 60) return `${diffMins}m ago`
    const diffHours = Math.floor(diffMins / 60)
    if (diffHours < 24) return `${diffHours}h ago`
    const diffDays = Math.floor(diffHours / 24)
    if (diffDays < 7) return `${diffDays}d ago`
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  } catch { return ts }
}

// === API ===
onMounted(async () => {
  try {
    const res = await api.get('/api/v1/whitelabel/overview', { params: { view: 'management' } })
    if (res.data?.data) {
      const d = res.data.data
      if (d.stats) data.value.stats = d.stats
      if (d.deployments) data.value.deployments = d.deployments
      if (d.tools) data.value.tools = d.tools
      if (d.activity) data.value.activity = d.activity
    }
  } catch { /* TODO: error toast */ }
  finally { loading.value = false }
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-text-primary mb-1">White Label Management</h1>
      <p class="text-sm text-text-secondary">Full overview of your branded deployments, tools, and performance</p>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard
        label="Total Impressions"
        :value="data.stats.total_impressions"
        icon="visibility"
        icon-color="#9B7EE0"
        :trend="data.stats.impressions_trend"
      />
      <StatsCard
        label="Active Deployments"
        :value="data.stats.active_deployments"
        icon="rocket_launch"
        icon-color="#60A5FA"
        :trend="data.stats.deployments_trend"
      />
      <StatsCard
        label="User Sessions"
        :value="data.stats.user_sessions"
        icon="group"
        icon-color="#48BB78"
        :trend="data.stats.sessions_trend"
      />
      <StatsCard
        label="Revenue Impact"
        :value="data.stats.revenue_impact"
        icon="payments"
        icon-color="#F59E0B"
        :trend="data.stats.revenue_trend"
      />
    </div>

    <!-- Active Deployments -->
    <div>
      <h2 class="text-lg font-semibold text-text-primary mb-3">Active Deployments</h2>
      <div v-if="loading" class="grid grid-cols-3 gap-4">
        <div v-for="i in 3" :key="i" class="bg-card-bg border border-border rounded-xl p-5 h-32 animate-pulse" />
      </div>
      <div v-else-if="data.deployments.length === 0" class="bg-card-bg border border-border rounded-xl p-6 text-center">
        <span class="material-symbols-rounded text-3xl text-text-muted mb-2 block">cloud_off</span>
        <p class="text-sm text-text-muted">No active deployments yet.</p>
      </div>
      <div v-else class="grid grid-cols-3 gap-4">
        <div
          v-for="dep in data.deployments"
          :key="dep.id"
          class="bg-card-bg border border-border rounded-xl p-5 hover:border-whitelabel/30 transition-colors"
        >
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="material-symbols-rounded text-lg text-text-secondary">{{ deploymentTypeIcons[dep.type] || 'web' }}</span>
              <span class="text-xs font-medium uppercase tracking-wider text-text-muted">{{ dep.type }}</span>
            </div>
            <span
              class="text-xs font-medium px-2 py-0.5 rounded-full capitalize"
              :style="{
                background: deploymentStatusStyles[dep.status]?.bg || '#1E293B',
                color: deploymentStatusStyles[dep.status]?.text || '#64748B',
              }"
            >{{ dep.status }}</span>
          </div>
          <h4 class="text-sm font-semibold text-text-primary mb-1 truncate">{{ dep.label }}</h4>
          <p class="text-xs text-text-muted mb-3 truncate">{{ dep.url }}</p>
          <div class="flex items-center justify-between">
            <span class="text-[10px] text-text-muted">Last: {{ formatTimestamp(dep.last_activity) }}</span>
            <button class="text-xs font-medium hover:underline" style="color: #9B7EE0">Manage &rarr;</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Performance Overview -->
    <div class="grid grid-cols-2 gap-4">
      <!-- Chart Placeholder -->
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <h3 class="text-sm font-semibold text-text-primary mb-4">Impressions &amp; Sessions &mdash; 30 Days</h3>
        <div class="h-48 flex items-center justify-center border border-dashed border-border rounded-lg">
          <div class="text-center">
            <span class="material-symbols-rounded text-3xl text-text-muted mb-1 block">show_chart</span>
            <span class="text-sm text-text-muted">Performance Chart</span>
            <p class="text-xs text-text-muted mt-1">Chart integration pending</p>
          </div>
        </div>
      </div>

      <!-- Key Metrics Summary -->
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <h3 class="text-sm font-semibold text-text-primary mb-4">Key Metrics</h3>
        <div class="space-y-4">
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">Widget Adoption</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.stats.active_deployments > 0 ? '78%' : '—' }}</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full transition-all duration-500" style="background: #9B7EE0" :style="{ width: data.stats.active_deployments > 0 ? '78%' : '0%' }" />
            </div>
          </div>
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">Page Engagement</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.stats.active_deployments > 0 ? '64%' : '—' }}</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-[#60A5FA] transition-all duration-500" :style="{ width: data.stats.active_deployments > 0 ? '64%' : '0%' }" />
            </div>
          </div>
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">Session Duration (avg)</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.stats.active_deployments > 0 ? '4m 32s' : '—' }}</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-community transition-all duration-500" :style="{ width: data.stats.active_deployments > 0 ? '56%' : '0%' }" />
            </div>
          </div>
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-text-secondary">Conversion Rate</span>
              <span class="text-sm font-semibold text-text-primary">{{ data.stats.active_deployments > 0 ? '12.3%' : '—' }}</span>
            </div>
            <div class="h-2 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-c-accent transition-all duration-500" :style="{ width: data.stats.active_deployments > 0 ? '12.3%' : '0%' }" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Configured Tools -->
    <div>
      <h2 class="text-lg font-semibold text-text-primary mb-3">Configured Tools</h2>
      <div v-if="loading" class="grid grid-cols-3 gap-4">
        <div v-for="i in 6" :key="i" class="bg-card-bg border border-border rounded-xl p-4 h-24 animate-pulse" />
      </div>
      <div v-else class="grid grid-cols-3 gap-4">
        <div
          v-for="tool in tools"
          :key="tool.key"
          class="bg-card-bg border border-border rounded-xl p-4 hover:border-whitelabel/30 transition-colors cursor-pointer"
          @click="router.push(tool.route)"
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <span class="material-symbols-rounded text-lg text-text-secondary">{{ tool.icon }}</span>
              <span class="text-sm font-medium text-text-primary">{{ tool.name }}</span>
            </div>
            <span class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full" :style="{ background: toolStatusDot(tool.status) }"></span>
              <span class="text-[10px] font-medium capitalize" :style="{ color: toolStatusDot(tool.status) }">{{ tool.status }}</span>
            </span>
          </div>
          <div class="text-xs font-medium hover:underline" style="color: #9B7EE0">Manage &rarr;</div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div>
      <h2 class="text-lg font-semibold text-text-primary mb-3">Recent Activity</h2>
      <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
        <div v-if="loading" class="p-6">
          <div v-for="i in 5" :key="i" class="h-10 bg-page-bg rounded-lg mb-2 animate-pulse" />
        </div>
        <div v-else-if="data.activity.length === 0" class="p-6 text-center">
          <span class="material-symbols-rounded text-3xl text-text-muted mb-2 block">history</span>
          <p class="text-sm text-text-muted">No recent activity. Events will appear here as your deployment is used.</p>
        </div>
        <div v-else class="divide-y divide-border">
          <div
            v-for="event in data.activity.slice(0, 10)"
            :key="event.id"
            class="flex items-center gap-3 px-5 py-3 hover:bg-white/2 transition-colors"
          >
            <span class="material-symbols-rounded text-lg shrink-0" :style="{ color: activityIconColor(event.type) }">{{ event.icon }}</span>
            <span class="text-sm text-text-primary flex-1">{{ event.description }}</span>
            <span class="text-xs text-text-muted shrink-0">{{ formatTimestamp(event.timestamp) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Quick Actions</div>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 text-white text-sm font-medium rounded-lg hover:opacity-90 transition-opacity flex items-center gap-2"
          style="background: #9B7EE0"
          @click="router.push('/b/whitelabel/deploy')"
        >
          <span class="material-symbols-rounded text-base">rocket_launch</span>
          Deploy Update
        </button>
        <button
          class="px-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary hover:border-whitelabel/30 transition-colors flex items-center gap-2"
          @click="router.push('/b/analytics?product=whitelabel')"
        >
          <span class="material-symbols-rounded text-base" style="color: #9B7EE0">bar_chart</span>
          View Analytics
        </button>
        <button
          class="px-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary hover:border-whitelabel/30 transition-colors flex items-center gap-2"
          @click="router.push('/b/whitelabel/integrations')"
        >
          <span class="material-symbols-rounded text-base" style="color: #9B7EE0">integration_instructions</span>
          Manage Integrations
        </button>
      </div>
    </div>
  </div>
</template>
