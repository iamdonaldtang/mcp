<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'

const router = useRouter()
const loading = ref(true)

// === Types ===
interface ChecklistStep {
  id: string
  section: 'setup' | 'configure' | 'go_live'
  label: string
  status: 'completed' | 'in_progress' | 'pending'
  expandable: boolean
  hint?: string
  pathDependent?: ('domain' | 'embed' | 'sdk')[]
}

interface ToolkitCard {
  key: string
  name: string
  icon: string
  status: 'active' | 'configured' | 'not_configured'
  route: string
}

interface OverviewStats {
  total_impressions: number
  impressions_trend: string
  active_widgets: number
  widgets_trend: string
  page_views: number
  page_views_trend: string
  integration_status: 'connected' | 'pending' | 'not_started'
}

// === State ===
const deploymentPath = ref<'domain' | 'embed' | 'sdk'>('embed')
const expandedStep = ref<string | null>(null)
const devKitEmail = ref('')
const devKitSending = ref(false)
const integrationVerified = ref(false)
const promoKitOpen = ref(false)
const shareUrl = ref('https://app.yourproject.com/community')

const stats = ref<OverviewStats>({
  total_impressions: 0,
  impressions_trend: '',
  active_widgets: 0,
  widgets_trend: '',
  page_views: 0,
  page_views_trend: '',
  integration_status: 'pending',
})

const steps = ref<ChecklistStep[]>([
  // SETUP (auto-completed by wizard)
  { id: 'community_created', section: 'setup', label: 'Community created', status: 'completed', expandable: false },
  { id: 'path_selected', section: 'setup', label: 'Deployment path selected', status: 'completed', expandable: false },
  { id: 'brand_configured', section: 'setup', label: 'Brand configured', status: 'completed', expandable: false },
  // CONFIGURE (path-dependent)
  { id: 'configure_widgets', section: 'configure', label: 'Configure widgets', status: 'pending', expandable: true, hint: 'Add community widgets to your app', pathDependent: ['embed'] },
  { id: 'build_page', section: 'configure', label: 'Build a page', status: 'pending', expandable: true, hint: 'Create a custom community page with Page Builder', pathDependent: ['embed'] },
  { id: 'setup_domain', section: 'configure', label: 'Set up custom domain', status: 'pending', expandable: true, hint: 'Point your domain to your White Label community', pathDependent: ['domain'] },
  { id: 'setup_sdk', section: 'configure', label: 'Set up SDK keys', status: 'pending', expandable: true, hint: 'Generate API keys and configure your SDK environment', pathDependent: ['sdk'] },
  // GO LIVE
  { id: 'preview', section: 'go_live', label: 'Preview deployment', status: 'pending', expandable: true, hint: 'See your branded community before going live' },
  { id: 'send_devkit', section: 'go_live', label: 'Send Dev Kit', status: 'pending', expandable: true, hint: 'Share integration instructions with your developer' },
  { id: 'integration_verified', section: 'go_live', label: 'Integration verified', status: 'pending', expandable: true, hint: 'Waiting for first API ping from your domain' },
  { id: 'announce', section: 'go_live', label: 'Announce your community', status: 'pending', expandable: true, hint: 'Generate and share promo materials' },
  { id: 'first_interaction', section: 'go_live', label: 'First user interaction', status: 'pending', expandable: true, hint: 'Real-time — updates when a user completes an action' },
])

// Filter steps by current deployment path
const filteredSteps = computed(() =>
  steps.value.filter(s => !s.pathDependent || s.pathDependent.includes(deploymentPath.value))
)

const setupSteps = computed(() => filteredSteps.value.filter(s => s.section === 'setup'))
const configureSteps = computed(() => filteredSteps.value.filter(s => s.section === 'configure'))
const goLiveSteps = computed(() => filteredSteps.value.filter(s => s.section === 'go_live'))

const completedCount = computed(() => filteredSteps.value.filter(s => s.status === 'completed').length)
const totalCount = computed(() => filteredSteps.value.length)
const progressPercent = computed(() => Math.round(completedCount.value / totalCount.value * 100))
const allCompleted = computed(() => completedCount.value === totalCount.value)

const toolkitCards = ref<ToolkitCard[]>([
  { key: 'widgets', name: 'Widget Library', icon: 'widgets', status: 'not_configured', route: '/b/whitelabel/widgets' },
  { key: 'pages', name: 'Page Builder', icon: 'web', status: 'not_configured', route: '/b/whitelabel/pages' },
  { key: 'domain', name: 'Custom Domain', icon: 'language', status: 'not_configured', route: '/b/whitelabel/domain' },
  { key: 'sdk', name: 'SDK & API', icon: 'code', status: 'not_configured', route: '/b/whitelabel/sdk' },
])

const integrationStatusLabel = computed(() => {
  switch (stats.value.integration_status) {
    case 'connected': return 'Connected'
    case 'pending': return 'Pending'
    default: return 'Not Started'
  }
})

const integrationStatusColor = computed(() => {
  switch (stats.value.integration_status) {
    case 'connected': return '#16A34A'
    case 'pending': return '#D97706'
    default: return '#64748B'
  }
})

// WebSocket for integration verification
let ws: WebSocket | null = null

function startIntegrationListener() {
  try {
    const wsUrl = (import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080') + '/ws/integration/ping'
    ws = new WebSocket(wsUrl)
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (data.type === 'integration_verified') {
          integrationVerified.value = true
          stats.value.integration_status = 'connected'
          const step = steps.value.find(s => s.id === 'integration_verified')
          if (step) step.status = 'completed'
        }
      } catch { /* ignore parse errors */ }
    }
  } catch { /* WebSocket not available */ }
}

// === API ===
onMounted(async () => {
  try {
    const [overviewRes, onboardingRes] = await Promise.allSettled([
      api.get('/api/v1/whitelabel/overview'),
      api.get('/api/v1/whitelabel/onboarding'),
    ])

    if (overviewRes.status === 'fulfilled' && overviewRes.value.data?.data) {
      const d = overviewRes.value.data.data
      stats.value = {
        total_impressions: d.total_impressions ?? 0,
        impressions_trend: d.impressions_trend ?? '',
        active_widgets: d.active_widgets ?? 0,
        widgets_trend: d.widgets_trend ?? '',
        page_views: d.page_views ?? 0,
        page_views_trend: d.page_views_trend ?? '',
        integration_status: d.integration_status ?? 'pending',
      }
      if (d.deployment_path) deploymentPath.value = d.deployment_path
      if (d.toolkit) {
        for (const tool of d.toolkit) {
          const found = toolkitCards.value.find(t => t.key === tool.key)
          if (found) found.status = tool.status
        }
      }
    }

    if (onboardingRes.status === 'fulfilled' && onboardingRes.value.data?.data?.steps) {
      for (const s of onboardingRes.value.data.data.steps) {
        const found = steps.value.find(x => x.id === s.id)
        if (found) found.status = s.status
      }
    }

    startIntegrationListener()
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (ws) {
    ws.close()
    ws = null
  }
})

function toggleExpand(id: string) {
  expandedStep.value = expandedStep.value === id ? null : id
}

function toolStatusColor(status: string): string {
  switch (status) {
    case 'active': return '#16A34A'
    case 'configured': return '#D97706'
    default: return '#64748B'
  }
}

function toolStatusLabel(status: string): string {
  switch (status) {
    case 'active': return 'Active'
    case 'configured': return 'Configured'
    default: return 'Not Configured'
  }
}

async function sendDevKit() {
  if (!devKitEmail.value.trim()) return
  devKitSending.value = true
  try {
    await api.post('/api/v1/whitelabel/devkit/send', { email: devKitEmail.value.trim() })
    const step = steps.value.find(s => s.id === 'send_devkit')
    if (step) step.status = 'completed'
    devKitEmail.value = ''
  } catch { /* TODO: error toast */ }
  finally { devKitSending.value = false }
}

function shareTwitter() {
  const text = encodeURIComponent(`We just launched our branded community powered by @TaskOnXYZ! Join us: ${shareUrl.value} #Web3`)
  window.open(`https://twitter.com/intent/tweet?text=${text}`, '_blank')
}

function shareTelegram() {
  const text = encodeURIComponent('Check out our new branded community!')
  window.open(`https://t.me/share/url?url=${encodeURIComponent(shareUrl.value)}&text=${text}`, '_blank')
}

function copyLink() {
  navigator.clipboard.writeText(shareUrl.value)
}
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <div class="flex items-center gap-3 mb-1">
        <h1 class="text-2xl font-bold text-text-primary">White Label</h1>
        <span class="px-2.5 py-0.5 text-xs font-medium rounded-full" style="background: #1F1A08; color: #D97706">Setting Up</span>
      </div>
      <p class="text-sm text-text-secondary">Complete the onboarding checklist to go live with your branded experience</p>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard
        label="Total Impressions"
        :value="stats.total_impressions"
        icon="visibility"
        icon-color="#9B7EE0"
        :trend="stats.impressions_trend"
      />
      <StatsCard
        label="Active Widgets"
        :value="stats.active_widgets"
        icon="widgets"
        icon-color="#60A5FA"
        :trend="stats.widgets_trend"
      />
      <StatsCard
        label="Page Views"
        :value="stats.page_views"
        icon="web"
        icon-color="#F59E0B"
        :trend="stats.page_views_trend"
      />
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="flex items-center justify-between mb-2">
          <span class="text-xs text-text-muted uppercase tracking-wider">Integration</span>
          <span class="material-symbols-rounded text-lg" :style="{ color: integrationStatusColor }">lan</span>
        </div>
        <div class="text-2xl font-bold text-text-primary">{{ integrationStatusLabel }}</div>
        <div class="mt-1 flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full" :style="{ background: integrationStatusColor }"></span>
          <span class="text-xs" :style="{ color: integrationStatusColor }">{{ stats.integration_status === 'connected' ? 'Live' : 'Awaiting ping' }}</span>
        </div>
      </div>
    </div>

    <!-- Onboarding Checklist -->
    <div class="bg-card-bg border border-border rounded-xl p-6">
      <!-- Progress Header -->
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm font-semibold text-text-primary">Getting Started</span>
        <span class="text-sm text-text-secondary">{{ completedCount }} of {{ totalCount }} complete</span>
      </div>
      <div class="h-2 bg-page-bg rounded-full mb-6 overflow-hidden">
        <div class="h-full rounded-full transition-all duration-300" style="background: #9B7EE0" :style="{ width: progressPercent + '%' }"></div>
      </div>

      <!-- Completed banner -->
      <div v-if="allCompleted" class="mb-6 p-4 rounded-lg" style="background: #0A2E1A; border: 1px solid #16A34A">
        <div class="flex items-center gap-3">
          <span class="material-symbols-rounded text-status-active text-xl">celebration</span>
          <div>
            <div class="text-sm font-semibold text-text-primary">Your White Label deployment is live!</div>
            <button class="text-xs hover:underline mt-1" style="color: #9B7EE0" @click="router.push('/b/whitelabel/management')">Go to Management &rarr;</button>
          </div>
        </div>
      </div>

      <!-- SETUP section -->
      <div class="mb-4">
        <div class="text-[10px] font-semibold text-text-muted uppercase tracking-wider mb-2">Setup</div>
        <div class="space-y-1">
          <div v-for="s in setupSteps" :key="s.id" class="flex items-center gap-3 py-2 px-3 rounded-lg">
            <span class="material-symbols-rounded text-lg text-status-active">check_circle</span>
            <span class="text-sm text-text-secondary line-through">{{ s.label }}</span>
          </div>
        </div>
      </div>

      <!-- CONFIGURE section -->
      <div class="mb-4">
        <div class="text-[10px] font-semibold text-text-muted uppercase tracking-wider mb-2">Configure</div>
        <div class="space-y-1">
          <div v-for="s in configureSteps" :key="s.id">
            <button
              class="flex items-center gap-3 py-2 px-3 rounded-lg w-full text-left hover:bg-white/2 transition-colors"
              @click="s.expandable && toggleExpand(s.id)"
            >
              <span v-if="s.status === 'completed'" class="material-symbols-rounded text-lg text-status-active">check_circle</span>
              <span v-else class="w-[18px] h-[18px] rounded-full border-2 border-border flex-shrink-0"></span>
              <span class="text-sm flex-1" :class="s.status === 'completed' ? 'text-text-secondary line-through' : 'text-text-primary'">{{ s.label }}</span>
              <span v-if="s.expandable" class="material-symbols-rounded text-base text-text-muted transition-transform" :class="expandedStep === s.id ? 'rotate-180' : ''">keyboard_arrow_down</span>
            </button>
            <div v-if="expandedStep === s.id" class="ml-10 mt-1 mb-2 p-3 bg-page-bg rounded-lg text-sm">
              <p class="text-text-muted mb-2">{{ s.hint }}</p>
              <button
                v-if="s.id === 'configure_widgets'"
                class="text-xs font-medium hover:underline"
                style="color: #9B7EE0"
                @click="router.push('/b/whitelabel/widgets')"
              >Open Widget Library &rarr;</button>
              <button
                v-else-if="s.id === 'build_page'"
                class="text-xs font-medium hover:underline"
                style="color: #9B7EE0"
                @click="router.push('/b/whitelabel/pages')"
              >Open Page Builder &rarr;</button>
              <button
                v-else-if="s.id === 'setup_domain'"
                class="text-xs font-medium hover:underline"
                style="color: #9B7EE0"
                @click="router.push('/b/whitelabel/domain')"
              >Configure Domain &rarr;</button>
              <button
                v-else-if="s.id === 'setup_sdk'"
                class="text-xs font-medium hover:underline"
                style="color: #9B7EE0"
                @click="router.push('/b/whitelabel/sdk')"
              >Set Up SDK &rarr;</button>
            </div>
          </div>
        </div>
      </div>

      <!-- GO LIVE section -->
      <div>
        <div class="text-[10px] font-semibold text-text-muted uppercase tracking-wider mb-2">Go Live</div>
        <div class="space-y-1">
          <div v-for="s in goLiveSteps" :key="s.id">
            <button
              class="flex items-center gap-3 py-2 px-3 rounded-lg w-full text-left hover:bg-white/2 transition-colors"
              @click="s.expandable && toggleExpand(s.id)"
            >
              <span v-if="s.status === 'completed'" class="material-symbols-rounded text-lg text-status-active">check_circle</span>
              <span v-else class="w-[18px] h-[18px] rounded-full border-2 border-border flex-shrink-0"></span>
              <span class="text-sm flex-1" :class="s.status === 'completed' ? 'text-text-secondary line-through' : 'text-text-primary'">{{ s.label }}</span>
              <span v-if="s.expandable" class="material-symbols-rounded text-base text-text-muted transition-transform" :class="expandedStep === s.id ? 'rotate-180' : ''">keyboard_arrow_down</span>
            </button>
            <div v-if="expandedStep === s.id" class="ml-10 mt-1 mb-2 p-3 bg-page-bg rounded-lg text-sm">
              <!-- Preview -->
              <template v-if="s.id === 'preview'">
                <p class="text-text-muted mb-2">See your branded community exactly as your users will</p>
                <button class="text-xs font-medium hover:underline" style="color: #9B7EE0" @click="router.push('/b/whitelabel/preview')">Open Preview &rarr;</button>
              </template>

              <!-- Send Dev Kit -->
              <template v-else-if="s.id === 'send_devkit'">
                <p class="text-text-muted mb-2">Send a standalone integration page to your developer — no TaskOn account required</p>
                <div class="flex items-center gap-2 mt-2">
                  <input
                    v-model="devKitEmail"
                    type="email"
                    placeholder="developer@yourproject.com"
                    class="flex-1 px-3 py-1.5 bg-card-bg border border-border rounded-lg text-xs text-text-primary placeholder:text-text-muted focus:outline-none focus:border-[#9B7EE0]"
                  />
                  <button
                    class="px-4 py-1.5 text-white text-xs font-medium rounded-lg hover:opacity-90 disabled:opacity-50"
                    style="background: #9B7EE0"
                    :disabled="!devKitEmail.trim() || devKitSending"
                    @click="sendDevKit"
                  >
                    {{ devKitSending ? 'Sending...' : 'Send' }}
                  </button>
                </div>
              </template>

              <!-- Integration Verified -->
              <template v-else-if="s.id === 'integration_verified'">
                <div v-if="integrationVerified" class="flex items-center gap-2">
                  <span class="material-symbols-rounded text-status-active text-lg">check_circle</span>
                  <span class="text-text-primary text-sm">First API ping detected! Integration is live.</span>
                </div>
                <div v-else class="flex items-center gap-2">
                  <span class="material-symbols-rounded text-text-muted text-lg animate-pulse">sensors</span>
                  <span class="text-text-muted">Listening for first API ping from your domain...</span>
                </div>
              </template>

              <!-- Announce (Promo Kit) -->
              <template v-else-if="s.id === 'announce'">
                <p class="text-text-muted mb-2">Share your branded community with your audience</p>
                <div class="space-y-3">
                  <div class="flex items-center gap-2">
                    <input type="text" :value="shareUrl" readonly class="flex-1 px-3 py-1.5 bg-card-bg border border-border rounded-lg text-xs text-text-secondary" />
                    <button class="px-3 py-1.5 text-white text-xs font-medium rounded-lg hover:opacity-90" style="background: #9B7EE0" @click="copyLink">Copy</button>
                  </div>
                  <div class="flex gap-2">
                    <button class="px-3 py-1.5 bg-[#1DA1F2]/20 text-[#1DA1F2] text-xs font-medium rounded-lg hover:bg-[#1DA1F2]/30" @click="shareTwitter">Twitter</button>
                    <button class="px-3 py-1.5 bg-[#5865F2]/20 text-[#5865F2] text-xs font-medium rounded-lg hover:bg-[#5865F2]/30">Discord</button>
                    <button class="px-3 py-1.5 bg-[#0088CC]/20 text-[#0088CC] text-xs font-medium rounded-lg hover:bg-[#0088CC]/30" @click="shareTelegram">Telegram</button>
                  </div>
                </div>
              </template>

              <!-- First User Interaction -->
              <template v-else-if="s.id === 'first_interaction'">
                <div class="flex items-center gap-2">
                  <span class="material-symbols-rounded text-text-muted text-lg">person</span>
                  <span class="text-text-muted">Waiting for first user action — this updates automatically</span>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Toolkit Cards -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Your Toolkit</div>
      <div class="grid grid-cols-4 gap-4">
        <div
          v-for="tool in toolkitCards"
          :key="tool.key"
          class="bg-card-bg border border-border rounded-xl p-4 hover:border-whitelabel/30 transition-colors"
        >
          <div class="flex items-center justify-between mb-3">
            <span class="material-symbols-rounded text-xl text-text-secondary">{{ tool.icon }}</span>
            <span class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full" :style="{ background: toolStatusColor(tool.status) }"></span>
              <span class="text-[10px] font-medium" :style="{ color: toolStatusColor(tool.status) }">{{ toolStatusLabel(tool.status) }}</span>
            </span>
          </div>
          <h4 class="text-sm font-medium text-text-primary mb-2">{{ tool.name }}</h4>
          <button
            class="text-xs font-medium hover:underline"
            style="color: #9B7EE0"
            @click="router.push(tool.route)"
          >
            {{ tool.status === 'not_configured' ? 'Configure' : 'Manage' }} &rarr;
          </button>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Quick Actions</div>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary hover:border-whitelabel/30 transition-colors flex items-center gap-2"
          @click="router.push('/b/whitelabel/widgets')"
        >
          <span class="material-symbols-rounded text-base" style="color: #9B7EE0">add</span>
          Add Widget
        </button>
        <button
          class="px-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary hover:border-whitelabel/30 transition-colors flex items-center gap-2"
          @click="router.push('/b/whitelabel/pages/new')"
        >
          <span class="material-symbols-rounded text-base" style="color: #9B7EE0">add</span>
          Create Page
        </button>
        <button
          class="px-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary hover:border-whitelabel/30 transition-colors flex items-center gap-2"
          @click="router.push('/b/analytics?product=whitelabel')"
        >
          <span class="material-symbols-rounded text-base" style="color: #9B7EE0">bar_chart</span>
          View Analytics
        </button>
      </div>
    </div>
  </div>
</template>
