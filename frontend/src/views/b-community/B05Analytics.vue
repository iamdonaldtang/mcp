<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { api } from '../../api/client'

interface OverviewStats {
  totalUsersAcquired: number
  totalUsersAcquiredTrend: number
  activeUsers: number
  activeUsersTrend: number
  totalConversions: number
  totalConversionsTrend: number
  revenue: number
  revenueTrend: number
}

interface ProductPerformance {
  quest: { campaignsRun: number; usersAcquired: number; avgCPA: number }
  community: { members: number; dau: number; retentionRate: number }
  whitelabel: { instances: number; integrations: number; userSessions: number }
  boost: { conversions: number; cpsAvg: number; revenue: number }
}

interface TopCampaign {
  id: string
  name: string
  product: 'quest' | 'community' | 'whitelabel' | 'boost'
  status: 'active' | 'draft' | 'completed' | 'paused'
  users: number
  conversionRate: number
  created: string
}

interface FunnelStep {
  label: string
  value: number
  dropOff: number
}

interface AnalyticsData {
  stats: OverviewStats
  productPerformance: ProductPerformance
  topCampaigns: TopCampaign[]
  funnel: FunnelStep[]
}

const loading = ref(true)
const selectedPeriod = ref('30d')
const periodOptions = ['7d', '30d', '90d']

const data = ref<AnalyticsData>({
  stats: {
    totalUsersAcquired: 0, totalUsersAcquiredTrend: 0,
    activeUsers: 0, activeUsersTrend: 0,
    totalConversions: 0, totalConversionsTrend: 0,
    revenue: 0, revenueTrend: 0,
  },
  productPerformance: {
    quest: { campaignsRun: 0, usersAcquired: 0, avgCPA: 0 },
    community: { members: 0, dau: 0, retentionRate: 0 },
    whitelabel: { instances: 0, integrations: 0, userSessions: 0 },
    boost: { conversions: 0, cpsAvg: 0, revenue: 0 },
  },
  topCampaigns: [],
  funnel: [
    { label: 'Visitors', value: 0, dropOff: 0 },
    { label: 'Sign Ups', value: 0, dropOff: 0 },
    { label: 'Task Completed', value: 0, dropOff: 0 },
    { label: 'Retained', value: 0, dropOff: 0 },
  ],
})

const productColorMap: Record<string, string> = {
  quest: '#5D7EF1',
  community: '#48BB78',
  whitelabel: '#9B7EE0',
  boost: '#ED8936',
}

const productBgMap: Record<string, string> = {
  quest: '#0F1538',
  community: '#0A1F1A',
  whitelabel: '#1A1033',
  boost: '#1F1508',
}

const productLabelMap: Record<string, string> = {
  quest: 'Quest',
  community: 'Community',
  whitelabel: 'White Label',
  boost: 'Boost',
}

function formatNumber(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n.toLocaleString()
}

function formatCurrency(n: number): string {
  return '$' + formatNumber(n)
}

function trendIcon(val: number): string {
  return val >= 0 ? 'trending_up' : 'trending_down'
}

function trendColor(val: number): string {
  return val >= 0 ? '#16A34A' : '#DC2626'
}

function statusClass(status: string) {
  const map: Record<string, string> = {
    active: 'bg-status-active-bg text-status-active',
    draft: 'bg-status-draft-bg text-status-draft',
    completed: 'bg-status-completed-bg text-status-completed',
    paused: 'bg-status-paused-bg text-status-paused',
  }
  return map[status] || map.draft
}

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/analytics/overview?period=${selectedPeriod.value}`)
    data.value = res.data
  } catch {
    // Keep defaults
  } finally {
    loading.value = false
  }
}

watch(selectedPeriod, () => fetchData())
onMounted(() => fetchData())
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Analytics</h1>
        <p class="text-sm text-text-secondary">Cross-product performance overview.</p>
      </div>
      <div class="flex items-center bg-card-bg border border-border rounded-lg">
        <button
          v-for="period in periodOptions"
          :key="period"
          class="px-4 py-2 text-sm font-medium transition-colors rounded-lg"
          :class="selectedPeriod === period
            ? 'bg-quest text-white'
            : 'text-text-muted hover:text-text-secondary'"
          @click="selectedPeriod = period"
        >{{ period }}</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <span class="material-symbols-rounded text-4xl text-text-muted animate-spin">progress_activity</span>
    </div>

    <template v-else>
      <!-- Stats Row -->
      <div class="grid grid-cols-4 gap-6">
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Total Users Acquired</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatNumber(data.stats.totalUsersAcquired) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.totalUsersAcquiredTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.totalUsersAcquiredTrend) }}</span>
              {{ Math.abs(data.stats.totalUsersAcquiredTrend) }}%
            </span>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Active Users</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatNumber(data.stats.activeUsers) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.activeUsersTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.activeUsersTrend) }}</span>
              {{ Math.abs(data.stats.activeUsersTrend) }}%
            </span>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Total Conversions</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatNumber(data.stats.totalConversions) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.totalConversionsTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.totalConversionsTrend) }}</span>
              {{ Math.abs(data.stats.totalConversionsTrend) }}%
            </span>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Revenue</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatCurrency(data.stats.revenue) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.revenueTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.revenueTrend) }}</span>
              {{ Math.abs(data.stats.revenueTrend) }}%
            </span>
          </div>
        </div>
      </div>

      <!-- Overview Chart -->
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-text-primary">User Acquisition & Engagement — {{ selectedPeriod === '7d' ? '7 Days' : selectedPeriod === '30d' ? '30 Days' : '90 Days' }}</h2>
          <div class="flex items-center gap-4 text-xs text-text-muted">
            <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-full bg-quest inline-block"></span> Acquired</span>
            <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-full bg-community inline-block"></span> Active</span>
            <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-full bg-status-paused inline-block"></span> Churned</span>
          </div>
        </div>
        <div class="h-64 bg-page-bg rounded-lg flex items-center justify-center">
          <span class="text-text-muted text-sm">Chart placeholder — integrate charting library</span>
        </div>
      </div>

      <!-- Product Performance -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Product Performance</div>
        <div class="grid grid-cols-4 gap-6">
          <!-- Quest -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-4">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{ background: productBgMap.quest }">
                <span class="material-symbols-rounded text-lg text-quest">campaign</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">Quest</span>
            </div>
            <div class="space-y-3">
              <div>
                <div class="text-xs text-text-muted">Campaigns Run</div>
                <div class="text-lg font-bold text-text-primary">{{ data.productPerformance.quest.campaignsRun }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">Users Acquired</div>
                <div class="text-lg font-bold text-text-primary">{{ formatNumber(data.productPerformance.quest.usersAcquired) }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">Avg CPA</div>
                <div class="text-lg font-bold text-text-primary">{{ formatCurrency(data.productPerformance.quest.avgCPA) }}</div>
              </div>
            </div>
            <router-link to="/b/analytics/quest" class="block mt-4 text-xs font-medium text-quest hover:underline">View Details →</router-link>
          </div>

          <!-- Community -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-4">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{ background: productBgMap.community }">
                <span class="material-symbols-rounded text-lg text-community">groups</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">Community</span>
            </div>
            <div class="space-y-3">
              <div>
                <div class="text-xs text-text-muted">Members</div>
                <div class="text-lg font-bold text-text-primary">{{ formatNumber(data.productPerformance.community.members) }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">DAU</div>
                <div class="text-lg font-bold text-text-primary">{{ formatNumber(data.productPerformance.community.dau) }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">Retention Rate</div>
                <div class="text-lg font-bold text-text-primary">{{ data.productPerformance.community.retentionRate }}%</div>
              </div>
            </div>
            <router-link to="/b/analytics/community" class="block mt-4 text-xs font-medium text-community hover:underline">View Details →</router-link>
          </div>

          <!-- White Label -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-4">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{ background: productBgMap.whitelabel }">
                <span class="material-symbols-rounded text-lg text-whitelabel">deployed_code</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">White Label</span>
            </div>
            <div class="space-y-3">
              <div>
                <div class="text-xs text-text-muted">Instances</div>
                <div class="text-lg font-bold text-text-primary">{{ data.productPerformance.whitelabel.instances }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">Integrations</div>
                <div class="text-lg font-bold text-text-primary">{{ data.productPerformance.whitelabel.integrations }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">User Sessions</div>
                <div class="text-lg font-bold text-text-primary">{{ formatNumber(data.productPerformance.whitelabel.userSessions) }}</div>
              </div>
            </div>
            <router-link to="/b/analytics/whitelabel" class="block mt-4 text-xs font-medium text-whitelabel hover:underline">View Details →</router-link>
          </div>

          <!-- Boost -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-4">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{ background: productBgMap.boost }">
                <span class="material-symbols-rounded text-lg text-boost">bolt</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">Boost</span>
            </div>
            <div class="space-y-3">
              <div>
                <div class="text-xs text-text-muted">Conversions</div>
                <div class="text-lg font-bold text-text-primary">{{ formatNumber(data.productPerformance.boost.conversions) }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">CPS Avg</div>
                <div class="text-lg font-bold text-text-primary">{{ formatCurrency(data.productPerformance.boost.cpsAvg) }}</div>
              </div>
              <div>
                <div class="text-xs text-text-muted">Revenue</div>
                <div class="text-lg font-bold text-text-primary">{{ formatCurrency(data.productPerformance.boost.revenue) }}</div>
              </div>
            </div>
            <router-link to="/b/analytics/boost" class="block mt-4 text-xs font-medium text-boost hover:underline">View Details →</router-link>
          </div>
        </div>
      </div>

      <!-- Top Campaigns Table -->
      <div>
        <div class="flex items-center justify-between mb-4">
          <div class="text-xs font-semibold text-text-muted uppercase tracking-wider">Top Campaigns</div>
          <router-link to="/b/analytics/campaigns" class="text-xs text-quest hover:underline">View All →</router-link>
        </div>
        <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border">
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Campaign Name</th>
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Product</th>
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Status</th>
                <th class="text-right px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Users</th>
                <th class="text-right px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Conversion Rate</th>
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Created</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border">
              <tr
                v-for="campaign in data.topCampaigns.slice(0, 5)"
                :key="campaign.id"
                class="hover:bg-white/2 transition-colors"
              >
                <td class="px-5 py-3.5 text-sm font-medium text-text-primary">{{ campaign.name }}</td>
                <td class="px-5 py-3.5">
                  <span
                    class="text-xs font-medium px-2 py-0.5 rounded-full"
                    :style="{ background: productBgMap[campaign.product], color: productColorMap[campaign.product] }"
                  >{{ productLabelMap[campaign.product] }}</span>
                </td>
                <td class="px-5 py-3.5">
                  <span class="text-xs font-medium px-2 py-0.5 rounded-full capitalize" :class="statusClass(campaign.status)">{{ campaign.status }}</span>
                </td>
                <td class="px-5 py-3.5 text-sm text-text-secondary text-right">{{ formatNumber(campaign.users) }}</td>
                <td class="px-5 py-3.5 text-sm text-text-secondary text-right">{{ campaign.conversionRate }}%</td>
                <td class="px-5 py-3.5 text-sm text-text-muted">{{ campaign.created }}</td>
              </tr>
              <tr v-if="!data.topCampaigns.length">
                <td colspan="6" class="px-5 py-8 text-center text-sm text-text-muted">No campaign data available</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Funnel Visualization -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Conversion Funnel</div>
        <div class="bg-card-bg border border-border rounded-xl p-6">
          <div class="flex items-center">
            <div
              v-for="(step, i) in data.funnel"
              :key="step.label"
              class="flex items-center"
              :class="i < data.funnel.length - 1 ? 'flex-1' : ''"
            >
              <div class="text-center flex-1">
                <div class="text-2xl font-bold text-text-primary mb-1">{{ formatNumber(step.value) }}</div>
                <div class="text-sm font-medium text-text-secondary mb-1">{{ step.label }}</div>
                <div v-if="i > 0" class="text-xs text-status-paused">-{{ step.dropOff }}% drop-off</div>
              </div>
              <div v-if="i < data.funnel.length - 1" class="flex flex-col items-center mx-2 shrink-0">
                <span class="material-symbols-rounded text-text-muted text-2xl">arrow_forward</span>
              </div>
            </div>
          </div>
          <!-- Funnel bar visualization -->
          <div class="mt-6 space-y-2">
            <div
              v-for="(step, i) in data.funnel"
              :key="'bar-' + step.label"
              class="flex items-center gap-3"
            >
              <span class="text-xs text-text-muted w-28 text-right shrink-0">{{ step.label }}</span>
              <div class="flex-1 bg-page-bg rounded-full h-6 overflow-hidden">
                <div
                  class="h-full rounded-full transition-all flex items-center justify-end pr-2"
                  :style="{
                    width: data.funnel[0].value > 0 ? ((step.value / data.funnel[0].value) * 100) + '%' : '0%',
                    background: i === 0 ? '#5D7EF1' : i === 1 ? '#48BB78' : i === 2 ? '#9B7EE0' : '#ED8936',
                  }"
                >
                  <span v-if="step.value > 0" class="text-xs font-medium text-white">{{ formatNumber(step.value) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
