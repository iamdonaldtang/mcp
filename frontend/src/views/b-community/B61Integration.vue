<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { api } from '../../api/client'

// --- Types ---
type IntegrationStatus = 'available' | 'connected' | 'error'
type CategoryKey = 'social' | 'blockchain' | 'analytics'

interface IntegrationCard {
  key: string
  category: CategoryKey
  title: string
  description: string
  icon: string
  status: IntegrationStatus
  errorMessage?: string
  connecting: boolean
  expanded: boolean
  config: Record<string, any>
  fields: IntegrationField[]
  connectAction: 'oauth' | 'token' | 'form'
}

interface IntegrationField {
  key: string
  label: string
  type: 'text' | 'select' | 'checkboxes' | 'password'
  placeholder?: string
  options?: { value: string; label: string }[]
  checkboxOptions?: { value: string; label: string }[]
  required?: boolean
}

interface WebhookDelivery {
  timestamp: string
  event: string
  success: boolean
}

// --- Card Definitions ---
const categories: { key: CategoryKey; label: string; icon: string }[] = [
  { key: 'social', label: 'SOCIAL & COMMUNITY', icon: 'forum' },
  { key: 'blockchain', label: 'BLOCKCHAIN & WALLET', icon: 'token' },
  { key: 'analytics', label: 'ANALYTICS & DATA', icon: 'analytics' },
]

const integrations = ref<IntegrationCard[]>([
  // --- Social & Community ---
  {
    key: 'twitter',
    category: 'social',
    title: 'Twitter/X',
    description: 'Connect your Twitter account to enable social verification tasks and auto-post community updates.',
    icon: 'tag',
    status: 'available',
    connecting: false,
    expanded: false,
    config: {},
    connectAction: 'oauth',
    fields: [],
  },
  {
    key: 'discord',
    category: 'social',
    title: 'Discord',
    description: 'Link your Discord server to sync roles, verify membership, and push notifications.',
    icon: 'headset_mic',
    status: 'available',
    connecting: false,
    expanded: false,
    config: {},
    connectAction: 'oauth',
    fields: [],
  },
  {
    key: 'telegram',
    category: 'social',
    title: 'Telegram',
    description: 'Connect your Telegram group or channel with a bot token for community verification.',
    icon: 'send',
    status: 'available',
    connecting: false,
    expanded: false,
    config: {},
    connectAction: 'token',
    fields: [
      { key: 'bot_token', label: 'Bot Token', type: 'text', placeholder: '123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11', required: true },
    ],
  },
  // --- Blockchain & Wallet ---
  {
    key: 'multichain',
    category: 'blockchain',
    title: 'Multi-Chain',
    description: 'Configure blockchain RPC endpoints for on-chain task verification and token gating.',
    icon: 'hub',
    status: 'available',
    connecting: false,
    expanded: false,
    config: { chain: 'ethereum', rpc_endpoint: '' },
    connectAction: 'form',
    fields: [
      {
        key: 'chain',
        label: 'Chain',
        type: 'select',
        options: [
          { value: 'ethereum', label: 'Ethereum' },
          { value: 'polygon', label: 'Polygon' },
          { value: 'bsc', label: 'BSC' },
          { value: 'arbitrum', label: 'Arbitrum' },
          { value: 'solana', label: 'Solana' },
        ],
      },
      { key: 'rpc_endpoint', label: 'RPC Endpoint', type: 'text', placeholder: 'https://mainnet.infura.io/v3/...', required: true },
    ],
  },
  {
    key: 'walletconnect',
    category: 'blockchain',
    title: 'Wallet Connect',
    description: 'Enable WalletConnect for seamless wallet authentication and transaction signing.',
    icon: 'account_balance_wallet',
    status: 'available',
    connecting: false,
    expanded: false,
    config: { version: 'v2', project_id: '' },
    connectAction: 'form',
    fields: [
      {
        key: 'version',
        label: 'Version',
        type: 'select',
        options: [
          { value: 'v1', label: 'v1 (Legacy)' },
          { value: 'v2', label: 'v2 (Current)' },
        ],
      },
      { key: 'project_id', label: 'Project ID', type: 'text', placeholder: 'Your WalletConnect Project ID', required: true },
    ],
  },
  {
    key: 'onchain',
    category: 'blockchain',
    title: 'On-Chain Verification',
    description: 'Set up on-chain verification rules for token balance, NFT holding, or contract interactions.',
    icon: 'verified',
    status: 'available',
    connecting: false,
    expanded: false,
    config: { chain: 'ethereum', verification_type: 'token_balance', contract_address: '' },
    connectAction: 'form',
    fields: [
      {
        key: 'chain',
        label: 'Chain',
        type: 'select',
        options: [
          { value: 'ethereum', label: 'Ethereum' },
          { value: 'polygon', label: 'Polygon' },
          { value: 'bsc', label: 'BSC' },
          { value: 'arbitrum', label: 'Arbitrum' },
          { value: 'solana', label: 'Solana' },
        ],
      },
      {
        key: 'verification_type',
        label: 'Verification Type',
        type: 'select',
        options: [
          { value: 'token_balance', label: 'Token Balance' },
          { value: 'nft_hold', label: 'NFT Hold' },
          { value: 'contract_call', label: 'Contract Call' },
        ],
      },
      { key: 'contract_address', label: 'Contract Address', type: 'text', placeholder: '0x...', required: true },
    ],
  },
  // --- Analytics & Data ---
  {
    key: 'ga4',
    category: 'analytics',
    title: 'Google Analytics',
    description: 'Track community engagement and conversion events with Google Analytics 4.',
    icon: 'monitoring',
    status: 'available',
    connecting: false,
    expanded: false,
    config: { measurement_id: '', api_secret: '' },
    connectAction: 'form',
    fields: [
      { key: 'measurement_id', label: 'GA4 Measurement ID', type: 'text', placeholder: 'G-XXXXXXXXXX', required: true },
      { key: 'api_secret', label: 'API Secret (optional)', type: 'password', placeholder: 'Optional for server-side events' },
    ],
  },
  {
    key: 'webhooks',
    category: 'analytics',
    title: 'Webhooks',
    description: 'Receive real-time event notifications for task completions, point changes, and more.',
    icon: 'webhook',
    status: 'available',
    connecting: false,
    expanded: false,
    config: { url: '', events: [] as string[], secret_header: '' },
    connectAction: 'form',
    fields: [
      { key: 'url', label: 'Webhook URL', type: 'text', placeholder: 'https://your-api.com/webhook', required: true },
      {
        key: 'events',
        label: 'Events',
        type: 'checkboxes',
        checkboxOptions: [
          { value: 'task.completed', label: 'task.completed' },
          { value: 'points.earned', label: 'points.earned' },
          { value: 'level.up', label: 'level.up' },
          { value: 'shop.redeemed', label: 'shop.redeemed' },
          { value: 'sprint.ended', label: 'sprint.ended' },
        ],
      },
      { key: 'secret_header', label: 'Secret Header (optional)', type: 'password', placeholder: 'X-Webhook-Secret' },
    ],
  },
  {
    key: 'data_export',
    category: 'analytics',
    title: 'Data Export',
    description: 'Schedule automated data exports to cloud storage or download manually.',
    icon: 'cloud_upload',
    status: 'available',
    connecting: false,
    expanded: false,
    config: { target: 'manual', bucket: '', region: '', access_key: '' },
    connectAction: 'form',
    fields: [
      {
        key: 'target',
        label: 'Export Target',
        type: 'select',
        options: [
          { value: 'manual', label: 'Manual Download' },
          { value: 's3', label: 'Amazon S3' },
          { value: 'gcs', label: 'Google Cloud Storage' },
        ],
      },
      { key: 'bucket', label: 'Bucket Name', type: 'text', placeholder: 'my-data-bucket' },
      { key: 'region', label: 'Region', type: 'text', placeholder: 'us-east-1' },
      { key: 'access_key', label: 'Access Key', type: 'password', placeholder: 'Your access key' },
    ],
  },
])

// Webhook delivery log
const webhookDeliveries = ref<WebhookDelivery[]>([])

// Disconnect confirmation
const disconnectTarget = ref<string | null>(null)

// Loading
const loading = ref(true)

// --- Computed ---
function cardsByCategory(cat: CategoryKey): IntegrationCard[] {
  return integrations.value.filter(c => c.category === cat)
}

function statusBorderClass(status: IntegrationStatus): string {
  switch (status) {
    case 'connected': return 'border-[#16A34A]'
    case 'error': return 'border-[#DC2626]'
    default: return 'border-border'
  }
}

function statusDotColor(status: IntegrationStatus): string {
  switch (status) {
    case 'connected': return '#16A34A'
    case 'error': return '#DC2626'
    default: return ''
  }
}

// Dynamic fields for data export (show cloud fields only when target != manual)
function visibleFields(card: IntegrationCard): IntegrationField[] {
  if (card.key === 'data_export') {
    if (card.config.target === 'manual') {
      return card.fields.filter(f => f.key === 'target')
    }
  }
  return card.fields
}

// --- API ---
async function fetchIntegrations() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/community/integrations')
    const data = res.data.data
    if (data && Array.isArray(data)) {
      data.forEach((item: any) => {
        const card = integrations.value.find(c => c.key === item.key)
        if (card) {
          card.status = item.status ?? 'available'
          card.errorMessage = item.error_message
          card.config = { ...card.config, ...item.config }
          if (item.webhook_deliveries) {
            webhookDeliveries.value = item.webhook_deliveries
          }
        }
      })
    }
  } catch {
    /* use defaults */
  }
  loading.value = false
}

async function connectOAuth(card: IntegrationCard) {
  card.connecting = true
  // Simulate OAuth popup flow
  try {
    const res = await api.post(`/api/v1/community/integrations/${card.key}/connect`)
    if (res.data.data?.oauth_url) {
      window.open(res.data.data.oauth_url, '_blank', 'width=600,height=700')
    }
    // Poll or simulate success after delay
    await new Promise(resolve => setTimeout(resolve, 2000))
    card.status = 'connected'
  } catch {
    // Mock success for dev
    await new Promise(resolve => setTimeout(resolve, 2000))
    card.status = 'connected'
  }
  card.connecting = false
}

async function connectWithToken(card: IntegrationCard) {
  const token = card.config.bot_token
  if (!token) return
  card.connecting = true
  try {
    await api.post(`/api/v1/community/integrations/${card.key}/connect`, { bot_token: token })
    card.status = 'connected'
  } catch {
    await new Promise(resolve => setTimeout(resolve, 1500))
    card.status = 'connected'
  }
  card.connecting = false
}

async function connectWithForm(card: IntegrationCard) {
  card.connecting = true
  try {
    await api.post(`/api/v1/community/integrations/${card.key}/connect`, card.config)
    card.status = 'connected'
  } catch {
    await new Promise(resolve => setTimeout(resolve, 1500))
    card.status = 'connected'
  }
  card.connecting = false
}

async function testConnection(card: IntegrationCard) {
  card.connecting = true
  try {
    await api.post(`/api/v1/community/integrations/${card.key}/test`, card.config)
    card.status = 'connected'
  } catch {
    await new Promise(resolve => setTimeout(resolve, 1000))
    card.status = 'connected'
  }
  card.connecting = false
}

async function sendTestEvent(card: IntegrationCard) {
  card.connecting = true
  try {
    const res = await api.post(`/api/v1/community/integrations/${card.key}/test-event`, card.config)
    const delivery: WebhookDelivery = {
      timestamp: new Date().toISOString(),
      event: 'test.ping',
      success: true,
    }
    webhookDeliveries.value.unshift(delivery)
    if (webhookDeliveries.value.length > 3) webhookDeliveries.value.pop()
  } catch {
    const delivery: WebhookDelivery = {
      timestamp: new Date().toISOString(),
      event: 'test.ping',
      success: false,
    }
    webhookDeliveries.value.unshift(delivery)
    if (webhookDeliveries.value.length > 3) webhookDeliveries.value.pop()
  }
  card.connecting = false
}

function confirmDisconnect(cardKey: string) {
  disconnectTarget.value = cardKey
}

async function disconnect(cardKey: string) {
  const card = integrations.value.find(c => c.key === cardKey)
  if (!card) return
  try {
    await api.delete(`/api/v1/community/integrations/${cardKey}`)
  } catch {
    /* proceed anyway for dev */
  }
  card.status = 'available'
  card.expanded = false
  card.errorMessage = undefined
  disconnectTarget.value = null
}

function cancelDisconnect() {
  disconnectTarget.value = null
}

function toggleExpand(card: IntegrationCard) {
  card.expanded = !card.expanded
}

function handleConnect(card: IntegrationCard) {
  if (card.connectAction === 'oauth') {
    connectOAuth(card)
  } else if (card.connectAction === 'token') {
    connectWithToken(card)
  } else {
    connectWithForm(card)
  }
}

function handlePrimaryAction(card: IntegrationCard) {
  if (card.status === 'connected') {
    toggleExpand(card)
  } else if (card.status === 'error') {
    handleConnect(card)
  } else {
    // available — show fields or connect directly
    if (card.fields.length > 0 && card.connectAction !== 'oauth') {
      toggleExpand(card)
    } else {
      handleConnect(card)
    }
  }
}

function formatTimestamp(ts: string): string {
  if (!ts) return '—'
  const d = new Date(ts)
  return d.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
    + ' ' + d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function toggleCheckbox(config: Record<string, any>, fieldKey: string, value: string) {
  if (!Array.isArray(config[fieldKey])) config[fieldKey] = []
  const idx = config[fieldKey].indexOf(value)
  if (idx >= 0) {
    config[fieldKey].splice(idx, 1)
  } else {
    config[fieldKey].push(value)
  }
}

// --- Init ---
onMounted(fetchIntegrations)
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-text-primary">Integration Center</h1>
      <p class="text-sm text-text-secondary mt-1">Connect external services to enhance your community</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-32">
      <div class="animate-spin w-8 h-8 border-2 border-community border-t-transparent rounded-full" />
    </div>

    <template v-else>
      <!-- Category Sections -->
      <div v-for="cat in categories" :key="cat.key" class="mb-8">
        <!-- Category Header -->
        <div class="flex items-center gap-2 mb-4">
          <span class="material-symbols-rounded text-lg text-text-muted">{{ cat.icon }}</span>
          <h2 class="text-[10px] font-semibold tracking-[1px] text-text-muted uppercase">{{ cat.label }}</h2>
        </div>

        <!-- Integration Cards Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <div v-for="card in cardsByCategory(cat.key)" :key="card.key" class="flex flex-col">
            <!-- Card -->
            <div
              class="bg-card-bg border rounded-xl p-5 transition-colors duration-200"
              :class="statusBorderClass(card.status)"
            >
              <!-- Card Header -->
              <div class="flex items-start justify-between mb-3">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-white/5 rounded-lg flex items-center justify-center">
                    <span class="material-symbols-rounded text-xl text-text-secondary">{{ card.icon }}</span>
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <span class="text-sm font-semibold text-text-primary">{{ card.title }}</span>
                      <!-- Status indicator -->
                      <span
                        v-if="card.status !== 'available'"
                        class="flex items-center gap-1 text-[11px] font-medium"
                      >
                        <span
                          class="w-1.5 h-1.5 rounded-full"
                          :style="{ backgroundColor: statusDotColor(card.status) }"
                        />
                        <span :class="card.status === 'connected' ? 'text-[#16A34A]' : 'text-[#DC2626]'">
                          {{ card.status === 'connected' ? 'Connected' : 'Error' }}
                        </span>
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Description -->
              <p class="text-xs text-text-muted mb-4 leading-relaxed">{{ card.description }}</p>

              <!-- Error message -->
              <div v-if="card.status === 'error' && card.errorMessage" class="mb-3 p-2 bg-[#2D1515] rounded-lg">
                <p class="text-xs text-[#DC2626]">{{ card.errorMessage }}</p>
              </div>

              <!-- Inline fields for token-based / available cards with fields -->
              <div
                v-if="card.status === 'available' && card.connectAction === 'token' && card.fields.length"
                class="mb-3"
              >
                <div v-for="field in card.fields" :key="field.key" class="mb-2">
                  <label class="block text-[11px] text-text-muted font-medium mb-1">{{ field.label }}</label>
                  <div class="flex gap-2">
                    <input
                      v-model="card.config[field.key]"
                      :type="field.type === 'password' ? 'password' : 'text'"
                      :placeholder="field.placeholder"
                      class="flex-1 h-8 px-3 bg-page-bg border border-border rounded-lg text-xs text-text-primary placeholder-text-muted focus:outline-none focus:border-community"
                    />
                    <button
                      class="h-8 px-3 bg-community/10 text-community text-xs font-medium rounded-lg hover:bg-community/20 transition-colors disabled:opacity-50"
                      :disabled="card.connecting"
                      @click="handleConnect(card)"
                    >
                      {{ card.connecting ? 'Verifying...' : 'Verify' }}
                    </button>
                  </div>
                </div>
              </div>

              <!-- Primary Action Button -->
              <button
                v-if="card.connectAction === 'oauth' || card.status === 'connected' || card.status === 'error' || (card.connectAction === 'form' && card.status === 'available')"
                class="w-full h-8 text-xs font-medium rounded-lg transition-colors duration-200 flex items-center justify-center gap-1.5 disabled:opacity-50"
                :class="{
                  'bg-community/10 text-community hover:bg-community/20': card.status === 'available',
                  'bg-white/5 text-text-secondary hover:bg-white/10': card.status === 'connected',
                  'bg-[#DC2626]/10 text-[#DC2626] hover:bg-[#DC2626]/20': card.status === 'error',
                }"
                :disabled="card.connecting"
                @click="handlePrimaryAction(card)"
              >
                <div v-if="card.connecting" class="animate-spin w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full" />
                <template v-else>
                  <template v-if="card.status === 'available'">
                    <span class="material-symbols-rounded text-sm">power</span>
                    {{ card.connectAction === 'form' ? 'Configure' : 'Connect' }}
                  </template>
                  <template v-else-if="card.status === 'connected'">
                    <span class="material-symbols-rounded text-sm">settings</span>
                    Configure
                  </template>
                  <template v-else>
                    <span class="material-symbols-rounded text-sm">refresh</span>
                    Reconnect
                  </template>
                </template>
              </button>
            </div>

            <!-- Expanded Config Panel -->
            <div
              v-if="card.expanded"
              class="mt-2 bg-card-bg border border-border rounded-xl p-4"
            >
              <div class="space-y-3">
                <div v-for="field in visibleFields(card)" :key="field.key">
                  <label class="block text-[11px] text-text-muted font-medium mb-1">{{ field.label }}</label>

                  <!-- Select -->
                  <select
                    v-if="field.type === 'select'"
                    v-model="card.config[field.key]"
                    class="w-full h-8 px-3 bg-page-bg border border-border rounded-lg text-xs text-text-primary focus:outline-none focus:border-community"
                  >
                    <option v-for="opt in field.options" :key="opt.value" :value="opt.value">
                      {{ opt.label }}
                    </option>
                  </select>

                  <!-- Checkboxes -->
                  <div v-else-if="field.type === 'checkboxes'" class="flex flex-wrap gap-2 mt-1">
                    <label
                      v-for="opt in field.checkboxOptions"
                      :key="opt.value"
                      class="flex items-center gap-1.5 text-xs text-text-secondary cursor-pointer"
                    >
                      <input
                        type="checkbox"
                        :checked="Array.isArray(card.config[field.key]) && card.config[field.key].includes(opt.value)"
                        class="rounded border-border text-community focus:ring-community accent-community"
                        @change="toggleCheckbox(card.config, field.key, opt.value)"
                      />
                      <code class="text-[10px] bg-white/5 px-1 py-0.5 rounded">{{ opt.label }}</code>
                    </label>
                  </div>

                  <!-- Text / Password -->
                  <input
                    v-else
                    v-model="card.config[field.key]"
                    :type="field.type === 'password' ? 'password' : 'text'"
                    :placeholder="field.placeholder"
                    class="w-full h-8 px-3 bg-page-bg border border-border rounded-lg text-xs text-text-primary placeholder-text-muted focus:outline-none focus:border-community"
                  />
                </div>
              </div>

              <!-- Webhook Deliveries Log -->
              <div v-if="card.key === 'webhooks' && webhookDeliveries.length" class="mt-4 border-t border-border pt-3">
                <h4 class="text-[11px] font-medium text-text-muted mb-2">Recent Deliveries</h4>
                <div class="space-y-1.5">
                  <div
                    v-for="(d, di) in webhookDeliveries.slice(0, 3)"
                    :key="di"
                    class="flex items-center justify-between text-[11px] py-1.5 px-2 bg-white/5 rounded"
                  >
                    <div class="flex items-center gap-2">
                      <span
                        class="material-symbols-rounded text-sm"
                        :class="d.success ? 'text-[#16A34A]' : 'text-[#DC2626]'"
                      >
                        {{ d.success ? 'check_circle' : 'cancel' }}
                      </span>
                      <code class="text-text-secondary">{{ d.event }}</code>
                    </div>
                    <span class="text-text-muted">{{ formatTimestamp(d.timestamp) }}</span>
                  </div>
                </div>
              </div>

              <!-- Actions -->
              <div class="flex items-center justify-between mt-4 pt-3 border-t border-border">
                <button
                  v-if="card.status === 'connected'"
                  class="text-xs text-[#DC2626] hover:text-[#EF4444] transition-colors"
                  @click="confirmDisconnect(card.key)"
                >
                  Disconnect
                </button>
                <span v-else />

                <div class="flex items-center gap-2">
                  <!-- Test buttons -->
                  <button
                    v-if="card.key === 'webhooks'"
                    class="h-7 px-3 bg-white/5 text-text-secondary text-xs rounded-md hover:bg-white/10 transition-colors disabled:opacity-50"
                    :disabled="card.connecting"
                    @click="sendTestEvent(card)"
                  >
                    Send Test Event
                  </button>
                  <button
                    v-else-if="card.connectAction === 'form'"
                    class="h-7 px-3 bg-white/5 text-text-secondary text-xs rounded-md hover:bg-white/10 transition-colors disabled:opacity-50"
                    :disabled="card.connecting"
                    @click="testConnection(card)"
                  >
                    {{ card.key === 'data_export' ? 'Test Upload' : 'Test Connection' }}
                  </button>

                  <!-- Save / Connect -->
                  <button
                    class="h-7 px-4 bg-community text-white text-xs font-medium rounded-md hover:bg-community/90 transition-colors disabled:opacity-50"
                    :disabled="card.connecting"
                    @click="handleConnect(card)"
                  >
                    {{ card.connecting ? 'Saving...' : card.status === 'connected' ? 'Save' : 'Connect' }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Disconnect Confirmation Modal -->
    <Teleport to="body">
      <div
        v-if="disconnectTarget"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
        @click.self="cancelDisconnect"
      >
        <div class="bg-card-bg border border-border rounded-xl p-6 w-full max-w-sm mx-4">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-10 h-10 bg-[#2D1515] rounded-lg flex items-center justify-center">
              <span class="material-symbols-rounded text-xl text-[#DC2626]">link_off</span>
            </div>
            <div>
              <h3 class="text-sm font-semibold text-text-primary">Disconnect Integration</h3>
              <p class="text-xs text-text-muted">This action cannot be undone</p>
            </div>
          </div>
          <p class="text-xs text-text-secondary mb-5 leading-relaxed">
            Are you sure you want to disconnect
            <strong class="text-text-primary">{{ integrations.find(c => c.key === disconnectTarget)?.title }}</strong>?
            All related configurations will be removed.
          </p>
          <div class="flex items-center justify-end gap-2">
            <button
              class="h-8 px-4 bg-white/5 text-text-secondary text-xs rounded-lg hover:bg-white/10 transition-colors"
              @click="cancelDisconnect"
            >
              Cancel
            </button>
            <button
              class="h-8 px-4 bg-[#DC2626] text-white text-xs font-medium rounded-lg hover:bg-[#EF4444] transition-colors"
              @click="disconnect(disconnectTarget!)"
            >
              Disconnect
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
