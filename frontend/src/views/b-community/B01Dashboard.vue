<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()

interface DashboardStats {
  totalUsers: number
  activeCampaigns: number
  pointsDistributed: number
  conversionRate: number
  totalUsersTrend: number
  activeCampaignsTrend: number
  pointsDistributedTrend: number
  conversionRateTrend: number
}

interface Campaign {
  id: string
  name: string
  product: 'quest' | 'community' | 'whitelabel' | 'boost'
  status: 'active' | 'draft' | 'completed' | 'paused'
  metric: string
  metricValue: number
  progress: number
}

interface Activity {
  id: string
  icon: string
  description: string
  timeAgo: string
  product: 'quest' | 'community' | 'whitelabel' | 'boost'
}

interface ProductBreakdown {
  quest: { campaigns: number }
  community: { members: number }
  whitelabel: { instances: number }
  boost: { conversions: number }
}

interface DashboardData {
  state: 'new' | 'active' | 'power'
  stats: DashboardStats
  campaigns: Campaign[]
  activities: Activity[]
  productBreakdown: ProductBreakdown
}

const loading = ref(true)
const data = ref<DashboardData>({
  state: 'new',
  stats: {
    totalUsers: 0,
    activeCampaigns: 0,
    pointsDistributed: 0,
    conversionRate: 0,
    totalUsersTrend: 0,
    activeCampaignsTrend: 0,
    pointsDistributedTrend: 0,
    conversionRateTrend: 0,
  },
  campaigns: [],
  activities: [],
  productBreakdown: {
    quest: { campaigns: 0 },
    community: { members: 0 },
    whitelabel: { instances: 0 },
    boost: { conversions: 0 },
  },
})

const userName = ref('there')

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

const goalCards = [
  {
    title: 'Get Users Fast',
    description: 'Launch quest campaigns to acquire new users with task-based incentives.',
    icon: 'campaign',
    product: 'quest',
    route: '/b/quest',
  },
  {
    title: 'Make Users Stay',
    description: 'Build a thriving community with retention mechanics and reward loops.',
    icon: 'groups',
    product: 'community',
    route: '/b/community',
  },
  {
    title: 'Own the Experience',
    description: 'Deploy a fully branded growth platform on your own domain.',
    icon: 'deployed_code',
    product: 'whitelabel',
    route: '/b/whitelabel',
  },
]

const quickStartSteps = [
  { number: 1, title: 'Choose a Product', description: 'Select Quest, Community, or White Label based on your growth goal.', icon: 'ads_click' },
  { number: 2, title: 'Create Your First Campaign', description: 'Follow the guided wizard to set up and launch in minutes.', icon: 'edit_note' },
  { number: 3, title: 'Share and Grow', description: 'Distribute your campaign link and watch your community grow.', icon: 'trending_up' },
]

const upsellMessages: Record<string, { text: string; route: string }> = {
  quest: { text: "You're using Quest. Try Community to retain your users", route: '/b/community' },
  community: { text: "You're using Community. Try White Label to own the full experience", route: '/b/whitelabel' },
  whitelabel: { text: "You're using White Label. Try Boost for pay-per-result campaigns", route: '/b/boost' },
  boost: { text: "You're using Boost. Try Quest to run your own acquisition campaigns", route: '/b/quest' },
}

const primaryProduct = computed(() => {
  if (!data.value.campaigns.length) return 'quest'
  const counts: Record<string, number> = {}
  data.value.campaigns.forEach((c) => {
    counts[c.product] = (counts[c.product] || 0) + 1
  })
  return Object.entries(counts).sort((a, b) => b[1] - a[1])[0][0]
})

const upsell = computed(() => upsellMessages[primaryProduct.value] || upsellMessages.quest)

function formatNumber(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n.toLocaleString()
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

onMounted(async () => {
  try {
    const res = await api.get('/api/v1/dashboard')
    data.value = res.data
    if (res.data.userName) userName.value = res.data.userName
  } catch {
    // Keep default new-user state on error
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="space-y-8">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <span class="material-symbols-rounded text-4xl text-text-muted animate-spin">progress_activity</span>
    </div>

    <!-- ===== NEW USER STATE ===== -->
    <template v-else-if="data.state === 'new'">
      <!-- Welcome Banner -->
      <div class="bg-card-bg border border-border rounded-xl p-8">
        <h1 class="text-3xl font-bold text-text-primary mb-2">Welcome to TaskOn, {{ userName }}!</h1>
        <p class="text-base text-text-secondary">Start growing your Web3 community. Choose a product to begin.</p>
      </div>

      <!-- Goal Cards -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Choose Your Goal</div>
        <div class="grid grid-cols-3 gap-6">
          <button
            v-for="card in goalCards"
            :key="card.product"
            class="text-left bg-card-bg border border-border rounded-xl p-6 hover:border-opacity-60 transition-all duration-200 group"
            @click="router.push(card.route)"
          >
            <div
              class="w-12 h-12 rounded-xl flex items-center justify-center mb-4"
              :style="{ background: productBgMap[card.product] }"
            >
              <span class="material-symbols-rounded text-2xl" :style="{ color: productColorMap[card.product] }">{{ card.icon }}</span>
            </div>
            <h3 class="text-lg font-semibold text-text-primary mb-1">{{ card.title }}</h3>
            <p class="text-sm text-text-secondary leading-relaxed mb-4">{{ card.description }}</p>
            <span class="text-sm font-medium group-hover:translate-x-1 transition-transform inline-block" :style="{ color: productColorMap[card.product] }">
              Get Started →
            </span>
          </button>
        </div>
      </div>

      <!-- Quick Start -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Quick Start</div>
        <div class="grid grid-cols-3 gap-6">
          <div
            v-for="step in quickStartSteps"
            :key="step.number"
            class="bg-card-bg border border-border rounded-xl p-5 flex items-start gap-4"
          >
            <div class="w-10 h-10 rounded-lg bg-page-bg border border-border flex items-center justify-center shrink-0">
              <span class="material-symbols-rounded text-text-secondary">{{ step.icon }}</span>
            </div>
            <div>
              <div class="text-sm font-semibold text-text-primary mb-0.5">{{ step.number }}. {{ step.title }}</div>
              <p class="text-xs text-text-secondary leading-relaxed">{{ step.description }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Resources -->
      <div class="border-t border-border pt-8">
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Resources</div>
        <div class="grid grid-cols-3 gap-6">
          <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/5 transition-colors block">
            <span class="material-symbols-rounded text-quest text-2xl mb-3 block">menu_book</span>
            <h4 class="text-sm font-semibold text-text-primary mb-1">Getting Started Guide</h4>
            <p class="text-xs text-text-secondary">Step-by-step guide to launching your first campaign</p>
          </a>
          <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/5 transition-colors block">
            <span class="material-symbols-rounded text-community text-2xl mb-3 block">play_circle</span>
            <h4 class="text-sm font-semibold text-text-primary mb-1">Watch Demo</h4>
            <p class="text-xs text-text-secondary">See how top projects use TaskOn for growth</p>
          </a>
          <a href="https://discord.gg/taskon" target="_blank" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/5 transition-colors block">
            <span class="material-symbols-rounded text-whitelabel text-2xl mb-3 block">forum</span>
            <h4 class="text-sm font-semibold text-text-primary mb-1">Join Discord</h4>
            <p class="text-xs text-text-secondary">Connect with other projects and get support</p>
          </a>
        </div>
      </div>
    </template>

    <!-- ===== ACTIVE USER STATE ===== -->
    <template v-else-if="data.state === 'active'">
      <!-- Header -->
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Dashboard</h1>
        <p class="text-sm text-text-secondary">Your growth overview at a glance.</p>
      </div>

      <!-- Stats Row -->
      <div class="grid grid-cols-4 gap-6">
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Total Users</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatNumber(data.stats.totalUsers) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.totalUsersTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.totalUsersTrend) }}</span>
              {{ Math.abs(data.stats.totalUsersTrend) }}%
            </span>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Active Campaigns</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ data.stats.activeCampaigns }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.activeCampaignsTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.activeCampaignsTrend) }}</span>
              {{ Math.abs(data.stats.activeCampaignsTrend) }}%
            </span>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Points Distributed</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatNumber(data.stats.pointsDistributed) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.pointsDistributedTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.pointsDistributedTrend) }}</span>
              {{ Math.abs(data.stats.pointsDistributedTrend) }}%
            </span>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Conversion Rate</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ data.stats.conversionRate }}%</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.conversionRateTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.conversionRateTrend) }}</span>
              {{ Math.abs(data.stats.conversionRateTrend) }}%
            </span>
          </div>
        </div>
      </div>

      <!-- Active Campaigns -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Active Campaigns</div>
        <div class="grid grid-cols-2 gap-6">
          <div
            v-for="campaign in data.campaigns"
            :key="campaign.id"
            class="bg-card-bg border border-border rounded-xl p-5"
          >
            <div class="flex items-start justify-between mb-3">
              <div>
                <h3 class="text-base font-semibold text-text-primary">{{ campaign.name }}</h3>
                <span
                  class="inline-block text-xs font-medium px-2 py-0.5 rounded-full mt-1"
                  :style="{ background: productBgMap[campaign.product], color: productColorMap[campaign.product] }"
                >{{ productLabelMap[campaign.product] }}</span>
              </div>
              <span
                class="text-xs font-medium px-2 py-0.5 rounded-full capitalize"
                :class="statusClass(campaign.status)"
              >{{ campaign.status }}</span>
            </div>
            <div class="text-sm text-text-secondary mb-3">{{ campaign.metric }}: <span class="text-text-primary font-medium">{{ formatNumber(campaign.metricValue) }}</span></div>
            <div class="w-full bg-page-bg rounded-full h-1.5">
              <div
                class="h-1.5 rounded-full transition-all"
                :style="{ width: campaign.progress + '%', background: productColorMap[campaign.product] }"
              ></div>
            </div>
            <div class="text-xs text-text-muted mt-1 text-right">{{ campaign.progress }}%</div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Quick Actions</div>
        <div class="flex gap-4">
          <button
            class="px-5 py-2.5 bg-quest text-white text-sm font-semibold rounded-xl hover:bg-quest/90 transition-colors"
            @click="router.push('/b/quest/wizard/step-1')"
          >Create Quest</button>
          <button
            class="px-5 py-2.5 bg-card-bg border border-border text-text-primary text-sm font-semibold rounded-xl hover:bg-white/5 transition-colors"
            @click="router.push('/b/community')"
          >Manage Community</button>
          <button
            class="px-5 py-2.5 bg-card-bg border border-border text-text-primary text-sm font-semibold rounded-xl hover:bg-white/5 transition-colors"
            @click="router.push('/b/analytics')"
          >View Analytics</button>
        </div>
      </div>

      <!-- Upsell Banner -->
      <div
        class="border rounded-xl p-5 flex items-center justify-between"
        :style="{ background: productBgMap[primaryProduct], borderColor: productColorMap[primaryProduct] + '33' }"
      >
        <div class="flex items-center gap-3">
          <span class="material-symbols-rounded text-2xl" :style="{ color: productColorMap[primaryProduct] }">auto_awesome</span>
          <span class="text-sm text-text-secondary">{{ upsell.text }}</span>
        </div>
        <button
          class="text-sm font-semibold px-4 py-2 rounded-lg transition-colors"
          :style="{ color: productColorMap[primaryProduct] }"
          @click="router.push(upsell.route)"
        >Learn More →</button>
      </div>
    </template>

    <!-- ===== POWER USER STATE ===== -->
    <template v-else>
      <!-- Header -->
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Dashboard</h1>
        <p class="text-sm text-text-secondary">Your growth engine overview.</p>
      </div>

      <!-- Enhanced Stats Row -->
      <div class="grid grid-cols-4 gap-6">
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Total Users</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatNumber(data.stats.totalUsers) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.totalUsersTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.totalUsersTrend) }}</span>
              {{ Math.abs(data.stats.totalUsersTrend) }}%
            </span>
          </div>
          <!-- Sparkline indicator -->
          <div class="mt-3 h-8 bg-page-bg rounded flex items-end gap-px px-1 pb-1">
            <div v-for="i in 12" :key="i" class="flex-1 rounded-t bg-quest/40" :style="{ height: (20 + Math.random() * 80) + '%' }"></div>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Active Campaigns</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ data.stats.activeCampaigns }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.activeCampaignsTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.activeCampaignsTrend) }}</span>
              {{ Math.abs(data.stats.activeCampaignsTrend) }}%
            </span>
          </div>
          <div class="mt-3 h-8 bg-page-bg rounded flex items-end gap-px px-1 pb-1">
            <div v-for="i in 12" :key="i" class="flex-1 rounded-t bg-community/40" :style="{ height: (20 + Math.random() * 80) + '%' }"></div>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Points Distributed</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ formatNumber(data.stats.pointsDistributed) }}</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.pointsDistributedTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.pointsDistributedTrend) }}</span>
              {{ Math.abs(data.stats.pointsDistributedTrend) }}%
            </span>
          </div>
          <div class="mt-3 h-8 bg-page-bg rounded flex items-end gap-px px-1 pb-1">
            <div v-for="i in 12" :key="i" class="flex-1 rounded-t bg-whitelabel/40" :style="{ height: (20 + Math.random() * 80) + '%' }"></div>
          </div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-2">Conversion Rate</div>
          <div class="flex items-end justify-between">
            <span class="text-2xl font-bold text-text-primary">{{ data.stats.conversionRate }}%</span>
            <span class="flex items-center gap-0.5 text-xs" :style="{ color: trendColor(data.stats.conversionRateTrend) }">
              <span class="material-symbols-rounded text-sm">{{ trendIcon(data.stats.conversionRateTrend) }}</span>
              {{ Math.abs(data.stats.conversionRateTrend) }}%
            </span>
          </div>
          <div class="mt-3 h-8 bg-page-bg rounded flex items-end gap-px px-1 pb-1">
            <div v-for="i in 12" :key="i" class="flex-1 rounded-t bg-boost/40" :style="{ height: (20 + Math.random() * 80) + '%' }"></div>
          </div>
        </div>
      </div>

      <!-- Growth Chart -->
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-text-primary">User Growth — Last 30 Days</h2>
          <div class="flex items-center gap-4 text-xs text-text-muted">
            <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-full bg-quest inline-block"></span> Acquired</span>
            <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-full bg-community inline-block"></span> Active</span>
            <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-full bg-status-paused inline-block"></span> Churned</span>
          </div>
        </div>
        <div class="h-52 bg-page-bg rounded-lg flex items-center justify-center">
          <span class="text-text-muted text-sm">Chart placeholder — integrate charting library</span>
        </div>
      </div>

      <!-- Product Breakdown -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Product Breakdown</div>
        <div class="grid grid-cols-4 gap-6">
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background: #0F1538">
                <span class="material-symbols-rounded text-lg text-quest">campaign</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">Quest</span>
            </div>
            <div class="text-xl font-bold text-text-primary">{{ data.productBreakdown.quest.campaigns }}</div>
            <div class="text-xs text-text-muted">campaigns</div>
          </div>
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background: #0A1F1A">
                <span class="material-symbols-rounded text-lg text-community">groups</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">Community</span>
            </div>
            <div class="text-xl font-bold text-text-primary">{{ formatNumber(data.productBreakdown.community.members) }}</div>
            <div class="text-xs text-text-muted">members</div>
          </div>
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background: #1A1033">
                <span class="material-symbols-rounded text-lg text-whitelabel">deployed_code</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">White Label</span>
            </div>
            <div class="text-xl font-bold text-text-primary">{{ data.productBreakdown.whitelabel.instances }}</div>
            <div class="text-xs text-text-muted">instances</div>
          </div>
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background: #1F1508">
                <span class="material-symbols-rounded text-lg text-boost">bolt</span>
              </div>
              <span class="text-sm font-semibold text-text-primary">Boost</span>
            </div>
            <div class="text-xl font-bold text-text-primary">{{ formatNumber(data.productBreakdown.boost.conversions) }}</div>
            <div class="text-xs text-text-muted">conversions</div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div>
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Recent Activity</div>
        <div class="bg-card-bg border border-border rounded-xl divide-y divide-border">
          <div
            v-for="activity in data.activities.slice(0, 10)"
            :key="activity.id"
            class="px-5 py-3.5 flex items-center gap-3"
          >
            <div
              class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0"
              :style="{ background: productBgMap[activity.product] }"
            >
              <span class="material-symbols-rounded text-base" :style="{ color: productColorMap[activity.product] }">{{ activity.icon }}</span>
            </div>
            <span class="text-sm text-text-secondary flex-1">{{ activity.description }}</span>
            <span class="text-xs text-text-muted shrink-0">{{ activity.timeAgo }}</span>
          </div>
          <div v-if="!data.activities.length" class="px-5 py-8 text-center text-sm text-text-muted">
            No recent activity
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
