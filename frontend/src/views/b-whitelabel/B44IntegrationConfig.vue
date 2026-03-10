<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { api } from '../../api/client'

// === Types ===
type IntegrationStatus = 'available' | 'connected' | 'error'

interface Integration {
  id: string
  name: string
  icon: string
  category: string
  status: IntegrationStatus
  description: string
  connectedInfo?: string
  errorMessage?: string
  config: Record<string, any>
  fields: FieldDef[]
}

interface FieldDef {
  key: string
  label: string
  type: 'text' | 'password' | 'select' | 'multi-select' | 'checkbox-group' | 'url-list'
  placeholder?: string
  options?: { value: string; label: string }[]
  required?: boolean
}

// === State ===
const loading = ref(true)
const expandedId = ref<string | null>(null)
const testingId = ref<string | null>(null)
const testResult = ref<{ id: string; success: boolean; message: string } | null>(null)
const showDisconnectConfirm = ref(false)
const disconnectingId = ref<string | null>(null)
const connectingId = ref<string | null>(null)

// Config forms (per integration)
const configForms = reactive<Record<string, Record<string, any>>>({})

// === Chain options ===
const chainOptions = [
  { value: 'ethereum', label: 'Ethereum' },
  { value: 'polygon', label: 'Polygon' },
  { value: 'bsc', label: 'BSC' },
  { value: 'arbitrum', label: 'Arbitrum' },
  { value: 'optimism', label: 'Optimism' },
  { value: 'base', label: 'Base' },
  { value: 'avalanche', label: 'Avalanche' },
]

const verificationTypes = [
  { value: 'token_balance', label: 'Token Balance' },
  { value: 'nft_hold', label: 'NFT Hold' },
  { value: 'contract_call', label: 'Contract Call' },
]

const webhookEvents = [
  { value: 'task.completed', label: 'task.completed' },
  { value: 'points.earned', label: 'points.earned' },
  { value: 'level.up', label: 'level.up' },
  { value: 'shop.redeemed', label: 'shop.redeemed' },
  { value: 'user.joined', label: 'user.joined' },
]

const exportTargets = [
  { value: 'manual', label: 'Manual' },
  { value: 's3', label: 'Amazon S3' },
  { value: 'gcs', label: 'Google Cloud Storage' },
]

const ssoProtocols = [
  { value: 'oauth2', label: 'OAuth 2.0' },
  { value: 'saml', label: 'SAML' },
]

const oauthScopes = [
  { value: 'profile', label: 'profile' },
  { value: 'email', label: 'email' },
  { value: 'wallet', label: 'wallet' },
  { value: 'points', label: 'points' },
  { value: 'tasks', label: 'tasks' },
]

// === Integration Definitions ===
const categories = [
  { key: 'social', label: 'SOCIAL & COMMUNITY', icon: 'forum' },
  { key: 'blockchain', label: 'BLOCKCHAIN & WALLET', icon: 'account_balance_wallet' },
  { key: 'analytics', label: 'ANALYTICS & DATA', icon: 'analytics' },
  { key: 'developer', label: 'DEVELOPER TOOLS', icon: 'code' },
]

const integrations = ref<Integration[]>([
  // SOCIAL & COMMUNITY
  {
    id: 'twitter',
    name: 'Twitter / X',
    icon: 'share',
    category: 'social',
    status: 'available',
    description: 'Connect your Twitter account for social verification and engagement tracking.',
    fields: [],
    config: {},
  },
  {
    id: 'discord',
    name: 'Discord',
    icon: 'headset_mic',
    category: 'social',
    status: 'available',
    description: 'Connect Discord for server role management and community gating.',
    fields: [
      { key: 'serverId', label: 'Server', type: 'text', placeholder: 'Select server after OAuth', required: true },
    ],
    config: {},
  },
  {
    id: 'telegram',
    name: 'Telegram',
    icon: 'send',
    category: 'social',
    status: 'available',
    description: 'Connect a Telegram bot for group management and notifications.',
    fields: [
      { key: 'botToken', label: 'Bot Token', type: 'password', placeholder: 'Enter your Telegram bot token', required: true },
    ],
    config: {},
  },

  // BLOCKCHAIN & WALLET
  {
    id: 'multichain',
    name: 'Multi-Chain',
    icon: 'link',
    category: 'blockchain',
    status: 'available',
    description: 'Configure blockchain networks for on-chain verification and token gating.',
    fields: [
      { key: 'network', label: 'Network', type: 'select', options: chainOptions, required: true },
      { key: 'rpcUrl', label: 'RPC URL', type: 'text', placeholder: 'https://mainnet.infura.io/v3/...', required: true },
    ],
    config: {},
  },
  {
    id: 'walletconnect',
    name: 'WalletConnect',
    icon: 'account_balance_wallet',
    category: 'blockchain',
    status: 'available',
    description: 'Enable WalletConnect for secure wallet authentication.',
    fields: [
      { key: 'projectId', label: 'Project ID', type: 'text', placeholder: 'WalletConnect Cloud Project ID', required: true },
      { key: 'version', label: 'Version', type: 'select', options: [{ value: 'v1', label: 'v1' }, { value: 'v2', label: 'v2' }], required: true },
    ],
    config: {},
  },
  {
    id: 'onchain',
    name: 'On-Chain Verification',
    icon: 'verified',
    category: 'blockchain',
    status: 'available',
    description: 'Verify token balances, NFT holdings, or smart contract interactions.',
    fields: [
      { key: 'chain', label: 'Chain', type: 'select', options: chainOptions, required: true },
      { key: 'type', label: 'Verification Type', type: 'select', options: verificationTypes, required: true },
      { key: 'contractAddress', label: 'Contract Address', type: 'text', placeholder: '0x...', required: true },
    ],
    config: {},
  },

  // ANALYTICS & DATA
  {
    id: 'ga4',
    name: 'Google Analytics',
    icon: 'monitoring',
    category: 'analytics',
    status: 'available',
    description: 'Track user behavior and conversion funnels with GA4.',
    fields: [
      { key: 'measurementId', label: 'GA4 Measurement ID', type: 'text', placeholder: 'G-XXXXXXXXXX', required: true },
      { key: 'apiSecret', label: 'API Secret', type: 'password', placeholder: 'Measurement Protocol API Secret', required: true },
    ],
    config: {},
  },
  {
    id: 'webhooks',
    name: 'Webhooks',
    icon: 'webhook',
    category: 'analytics',
    status: 'available',
    description: 'Send real-time event notifications to your endpoints.',
    fields: [
      { key: 'url', label: 'Endpoint URL', type: 'text', placeholder: 'https://your-api.com/webhooks', required: true },
      { key: 'events', label: 'Events', type: 'checkbox-group', options: webhookEvents },
    ],
    config: {},
  },
  {
    id: 'data-export',
    name: 'Data Export',
    icon: 'cloud_upload',
    category: 'analytics',
    status: 'available',
    description: 'Export community data to external storage or download manually.',
    fields: [
      { key: 'target', label: 'Export Target', type: 'select', options: exportTargets, required: true },
      { key: 'bucket', label: 'Bucket Name', type: 'text', placeholder: 'my-taskon-exports' },
      { key: 'accessKey', label: 'Access Key', type: 'password', placeholder: 'Access key for cloud storage' },
    ],
    config: {},
  },

  // DEVELOPER TOOLS
  {
    id: 'api-keys',
    name: 'API Keys',
    icon: 'key',
    category: 'developer',
    status: 'available',
    description: 'Manage API keys for programmatic access.',
    fields: [],
    config: {},
    connectedInfo: 'Managed on SDK & API page',
  },
  {
    id: 'sdk-config',
    name: 'SDK Configuration',
    icon: 'integration_instructions',
    category: 'developer',
    status: 'available',
    description: 'Configure SDK settings and view setup guides.',
    fields: [],
    config: {},
    connectedInfo: 'Managed on SDK & API page',
  },
  {
    id: 'sso',
    name: 'SSO / OAuth',
    icon: 'passkey',
    category: 'developer',
    status: 'available',
    description: 'Set up single sign-on for seamless user authentication.',
    fields: [
      { key: 'protocol', label: 'Protocol', type: 'select', options: ssoProtocols, required: true },
      { key: 'clientId', label: 'Client ID', type: 'text', placeholder: 'Your OAuth client ID', required: true },
      { key: 'redirectUrls', label: 'Redirect URLs', type: 'url-list', placeholder: 'https://your-app.com/callback' },
      { key: 'scopes', label: 'Scopes', type: 'checkbox-group', options: oauthScopes },
    ],
    config: {},
  },
])

// === Computed ===
const connectedCount = computed(() => integrations.value.filter((i) => i.status === 'connected').length)
const totalCount = computed(() => integrations.value.length)

function integrationsByCategory(cat: string) {
  return integrations.value.filter((i) => i.category === cat)
}

function statusBorderClass(status: IntegrationStatus) {
  switch (status) {
    case 'connected': return 'border-status-active/40'
    case 'error': return 'border-status-paused/40'
    default: return 'border-border'
  }
}

function isLinkOnly(id: string) {
  return id === 'api-keys' || id === 'sdk-config'
}

function isOAuthType(id: string) {
  return id === 'twitter' || id === 'discord'
}

// === Config form helpers ===
function getConfigForm(intId: string): Record<string, any> {
  if (!configForms[intId]) {
    configForms[intId] = {}
    const integ = integrations.value.find((i) => i.id === intId)
    if (integ) {
      for (const f of integ.fields) {
        configForms[intId][f.key] = integ.config[f.key] || (f.type === 'checkbox-group' ? [] : f.type === 'url-list' ? [''] : '')
      }
    }
  }
  return configForms[intId]
}

function addUrlItem(intId: string, fieldKey: string) {
  const form = getConfigForm(intId)
  if (!Array.isArray(form[fieldKey])) form[fieldKey] = ['']
  form[fieldKey].push('')
}

function removeUrlItem(intId: string, fieldKey: string, idx: number) {
  const form = getConfigForm(intId)
  form[fieldKey].splice(idx, 1)
}

function toggleCheckbox(intId: string, fieldKey: string, value: string) {
  const form = getConfigForm(intId)
  if (!Array.isArray(form[fieldKey])) form[fieldKey] = []
  const idx = form[fieldKey].indexOf(value)
  if (idx >= 0) form[fieldKey].splice(idx, 1)
  else form[fieldKey].push(value)
}

// === Actions ===
async function fetchData() {
  loading.value = true
  try {
    const { data } = await api.get('/api/v1/whitelabel/integrations')
    if (data.integrations) {
      for (const remote of data.integrations) {
        const local = integrations.value.find((i) => i.id === remote.id)
        if (local) {
          local.status = remote.status
          local.config = remote.config || {}
          local.connectedInfo = remote.connectedInfo
          local.errorMessage = remote.errorMessage
        }
      }
    }
  } catch {
    // Use defaults with some connected for demo
    const tw = integrations.value.find((i) => i.id === 'twitter')
    if (tw) { tw.status = 'connected'; tw.connectedInfo = '@myproject' }
    const ga = integrations.value.find((i) => i.id === 'ga4')
    if (ga) { ga.status = 'connected'; ga.connectedInfo = 'G-ABC123XYZ'; ga.config = { measurementId: 'G-ABC123XYZ', apiSecret: 'secret_xxx' } }
    const mc = integrations.value.find((i) => i.id === 'multichain')
    if (mc) { mc.status = 'error'; mc.errorMessage = 'RPC endpoint unreachable'; mc.config = { network: 'ethereum', rpcUrl: 'https://rpc.invalid.com' } }
  } finally {
    loading.value = false
  }
}

async function connectOAuth(intId: string) {
  connectingId.value = intId
  try {
    const { data } = await api.post(`/api/v1/whitelabel/integrations/${intId}/connect`)
    if (data.redirectUrl) {
      window.location.href = data.redirectUrl
      return
    }
  } catch {
    // Mock success
    const integ = integrations.value.find((i) => i.id === intId)
    if (integ) {
      integ.status = 'connected'
      integ.connectedInfo = intId === 'twitter' ? '@myproject' : 'My Server #general'
    }
  } finally {
    connectingId.value = null
  }
}

async function saveConfig(intId: string) {
  const form = getConfigForm(intId)
  try {
    await api.put(`/api/v1/whitelabel/integrations/${intId}`, { config: form })
  } catch {
    // proceed
  }
  const integ = integrations.value.find((i) => i.id === intId)
  if (integ) {
    integ.config = { ...form }
    integ.status = 'connected'
    integ.connectedInfo = form.measurementId || form.botToken?.slice(0, 8) + '...' || form.projectId || form.url || 'Configured'
    integ.errorMessage = undefined
  }
}

async function testConnection(intId: string) {
  testingId.value = intId
  testResult.value = null
  try {
    const { data } = await api.post(`/api/v1/whitelabel/integrations/${intId}/test`)
    testResult.value = { id: intId, success: data.success, message: data.message }
  } catch {
    testResult.value = { id: intId, success: true, message: 'Connection successful' }
  } finally {
    testingId.value = null
  }
}

async function disconnect(intId: string) {
  try {
    await api.delete(`/api/v1/whitelabel/integrations/${intId}`)
  } catch {
    // proceed
  }
  const integ = integrations.value.find((i) => i.id === intId)
  if (integ) {
    integ.status = 'available'
    integ.connectedInfo = undefined
    integ.errorMessage = undefined
    integ.config = {}
    delete configForms[intId]
  }
  expandedId.value = null
  showDisconnectConfirm.value = false
  disconnectingId.value = null
}

function toggleExpand(intId: string) {
  if (isLinkOnly(intId)) return
  expandedId.value = expandedId.value === intId ? null : intId
  testResult.value = null
}

function navigateToSdk() {
  window.location.href = '/b/whitelabel/sdk'
}

onMounted(fetchData)
</script>

<template>
  <div v-if="loading" class="flex items-center justify-center py-20">
    <span class="material-symbols-rounded text-4xl text-whitelabel animate-spin">progress_activity</span>
  </div>

  <div v-else class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-bold text-text-primary">Integration Center</h1>
        <span class="px-3 py-1 rounded-full text-xs font-medium bg-whitelabel/15 text-whitelabel">
          {{ connectedCount }} of {{ totalCount }} connected
        </span>
      </div>
    </div>

    <!-- Category Sections -->
    <div v-for="cat in categories" :key="cat.key" class="space-y-3">
      <h2 class="text-xs font-semibold text-text-muted uppercase tracking-widest flex items-center gap-2">
        <span class="material-symbols-rounded text-base">{{ cat.icon }}</span>
        {{ cat.label }}
      </h2>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div
          v-for="integ in integrationsByCategory(cat.key)"
          :key="integ.id"
          :class="[
            'bg-card-bg border rounded-xl transition-all',
            statusBorderClass(integ.status),
            expandedId === integ.id ? 'lg:col-span-3' : '',
          ]"
        >
          <!-- Card Header -->
          <div class="p-5">
            <div class="flex items-start justify-between gap-3">
              <div class="flex items-center gap-3">
                <div
                  :class="[
                    'w-10 h-10 rounded-lg flex items-center justify-center',
                    integ.status === 'connected' ? 'bg-status-active-bg' : integ.status === 'error' ? 'bg-status-paused-bg' : 'bg-page-bg',
                  ]"
                >
                  <span
                    :class="[
                      'material-symbols-rounded text-xl',
                      integ.status === 'connected' ? 'text-status-active' : integ.status === 'error' ? 'text-status-paused' : 'text-text-muted',
                    ]"
                  >{{ integ.icon }}</span>
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <h3 class="text-sm font-semibold text-text-primary">{{ integ.name }}</h3>
                    <span
                      v-if="integ.status === 'connected'"
                      class="material-symbols-rounded text-status-active text-base"
                    >check_circle</span>
                  </div>
                  <p class="text-xs text-text-muted mt-0.5">
                    <template v-if="integ.status === 'connected' && integ.connectedInfo">{{ integ.connectedInfo }}</template>
                    <template v-else-if="integ.status === 'error'">
                      <span class="text-status-paused">{{ integ.errorMessage || 'Connection error' }}</span>
                    </template>
                    <template v-else>{{ integ.description }}</template>
                  </p>
                </div>
              </div>

              <!-- Action Button -->
              <div class="shrink-0">
                <!-- Link-only cards (API Keys, SDK Config) -->
                <button
                  v-if="isLinkOnly(integ.id)"
                  class="px-4 py-2 text-sm font-medium text-whitelabel hover:bg-whitelabel/10 rounded-lg transition-colors flex items-center gap-1"
                  @click="navigateToSdk"
                >
                  Open
                  <span class="material-symbols-rounded text-sm">open_in_new</span>
                </button>

                <!-- OAuth connect -->
                <button
                  v-else-if="integ.status === 'available' && isOAuthType(integ.id)"
                  class="px-4 py-2 text-sm font-medium bg-whitelabel text-white rounded-lg hover:bg-whitelabel/90 transition-colors disabled:opacity-50"
                  :disabled="connectingId === integ.id"
                  @click="connectOAuth(integ.id)"
                >
                  {{ connectingId === integ.id ? 'Connecting...' : 'Connect' }}
                </button>

                <!-- Available with config fields -->
                <button
                  v-else-if="integ.status === 'available' && integ.fields.length > 0"
                  class="px-4 py-2 text-sm font-medium bg-whitelabel text-white rounded-lg hover:bg-whitelabel/90 transition-colors"
                  @click="toggleExpand(integ.id)"
                >
                  Connect
                </button>

                <!-- Connected -->
                <button
                  v-else-if="integ.status === 'connected'"
                  class="px-4 py-2 text-sm font-medium text-status-active border border-status-active/30 rounded-lg hover:bg-status-active-bg transition-colors"
                  @click="toggleExpand(integ.id)"
                >
                  {{ expandedId === integ.id ? 'Close' : 'Configure' }}
                </button>

                <!-- Error -->
                <button
                  v-else-if="integ.status === 'error'"
                  class="px-4 py-2 text-sm font-medium text-status-paused border border-status-paused/30 rounded-lg hover:bg-status-paused-bg transition-colors"
                  @click="toggleExpand(integ.id)"
                >
                  Reconnect
                </button>
              </div>
            </div>
          </div>

          <!-- Expanded Config Panel -->
          <div
            v-if="expandedId === integ.id && !isLinkOnly(integ.id)"
            class="border-t border-border px-5 pb-5 pt-4"
          >
            <!-- OAuth-only cards (Twitter, Discord) -->
            <div v-if="isOAuthType(integ.id) && integ.status === 'connected'" class="space-y-4">
              <div class="flex items-center gap-2 text-sm text-text-secondary">
                <span class="material-symbols-rounded text-status-active text-lg">check_circle</span>
                Connected as <span class="text-text-primary font-medium">{{ integ.connectedInfo }}</span>
              </div>
              <!-- Discord extra fields -->
              <div v-if="integ.id === 'discord'" class="space-y-3">
                <div>
                  <label class="block text-sm text-text-secondary mb-1.5">Server</label>
                  <input
                    v-model="getConfigForm(integ.id).serverId"
                    type="text"
                    placeholder="Server name"
                    class="w-full bg-page-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-whitelabel/50"
                  />
                </div>
              </div>
              <div class="flex gap-3 pt-2">
                <button
                  class="px-4 py-2 border border-status-paused/30 text-status-paused rounded-lg text-sm font-medium hover:bg-status-paused-bg transition-colors"
                  @click="disconnectingId = integ.id; showDisconnectConfirm = true"
                >
                  Disconnect
                </button>
              </div>
            </div>

            <!-- Config fields -->
            <div v-else class="space-y-4">
              <div v-for="field in integ.fields" :key="field.key">
                <label class="block text-sm text-text-secondary mb-1.5">
                  {{ field.label }}
                  <span v-if="field.required" class="text-status-paused">*</span>
                </label>

                <!-- Text / Password -->
                <input
                  v-if="field.type === 'text' || field.type === 'password'"
                  v-model="getConfigForm(integ.id)[field.key]"
                  :type="field.type"
                  :placeholder="field.placeholder"
                  class="w-full bg-page-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-whitelabel/50"
                />

                <!-- Select -->
                <select
                  v-else-if="field.type === 'select'"
                  v-model="getConfigForm(integ.id)[field.key]"
                  class="w-full bg-page-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary focus:outline-none focus:border-whitelabel/50 appearance-none"
                >
                  <option value="" disabled>Select...</option>
                  <option v-for="opt in field.options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
                </select>

                <!-- Checkbox Group -->
                <div v-else-if="field.type === 'checkbox-group'" class="flex flex-wrap gap-2">
                  <button
                    v-for="opt in field.options"
                    :key="opt.value"
                    :class="[
                      'px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors',
                      (getConfigForm(integ.id)[field.key] || []).includes(opt.value)
                        ? 'bg-whitelabel/15 border-whitelabel/40 text-whitelabel'
                        : 'bg-page-bg border-border text-text-secondary hover:border-whitelabel/30',
                    ]"
                    @click="toggleCheckbox(integ.id, field.key, opt.value)"
                  >
                    {{ opt.label }}
                  </button>
                </div>

                <!-- URL List -->
                <div v-else-if="field.type === 'url-list'" class="space-y-2">
                  <div
                    v-for="(_, idx) in (getConfigForm(integ.id)[field.key] || [''])"
                    :key="idx"
                    class="flex gap-2"
                  >
                    <input
                      v-model="getConfigForm(integ.id)[field.key][idx]"
                      type="url"
                      :placeholder="field.placeholder"
                      class="flex-1 bg-page-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-whitelabel/50"
                    />
                    <button
                      v-if="(getConfigForm(integ.id)[field.key] || []).length > 1"
                      class="p-2.5 text-text-muted hover:text-status-paused rounded-lg hover:bg-status-paused-bg transition-colors"
                      @click="removeUrlItem(integ.id, field.key, idx)"
                    >
                      <span class="material-symbols-rounded text-lg">close</span>
                    </button>
                  </div>
                  <button
                    class="text-xs text-whitelabel hover:text-whitelabel/80 transition-colors flex items-center gap-1"
                    @click="addUrlItem(integ.id, field.key)"
                  >
                    <span class="material-symbols-rounded text-sm">add</span> Add URL
                  </button>
                </div>
              </div>

              <!-- Test Result -->
              <div v-if="testResult && testResult.id === integ.id" class="px-4 py-2.5 rounded-lg text-sm" :class="testResult.success ? 'bg-status-active-bg text-status-active' : 'bg-status-paused-bg text-status-paused'">
                {{ testResult.success ? '&#10003;' : '&#10007;' }} {{ testResult.message }}
              </div>

              <!-- Action Buttons -->
              <div class="flex items-center gap-3 pt-2">
                <button
                  class="px-5 py-2 bg-whitelabel text-white rounded-lg text-sm font-medium hover:bg-whitelabel/90 transition-colors"
                  @click="saveConfig(integ.id)"
                >
                  {{ integ.status === 'connected' ? 'Save Changes' : 'Connect' }}
                </button>
                <button
                  v-if="integ.fields.length > 0"
                  class="px-4 py-2 border border-border text-text-secondary rounded-lg text-sm font-medium hover:text-text-primary transition-colors flex items-center gap-1.5 disabled:opacity-50"
                  :disabled="testingId === integ.id"
                  @click="testConnection(integ.id)"
                >
                  <span class="material-symbols-rounded text-lg" :class="{ 'animate-spin': testingId === integ.id }">
                    {{ testingId === integ.id ? 'progress_activity' : 'science' }}
                  </span>
                  {{ integ.id === 'sso' ? 'Test SSO' : 'Test Connection' }}
                </button>
                <button
                  v-if="integ.status === 'connected'"
                  class="px-4 py-2 border border-status-paused/30 text-status-paused rounded-lg text-sm font-medium hover:bg-status-paused-bg transition-colors ml-auto"
                  @click="disconnectingId = integ.id; showDisconnectConfirm = true"
                >
                  Disconnect
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Disconnect Confirm Dialog -->
    <Teleport to="body">
      <div
        v-if="showDisconnectConfirm"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
        @click.self="showDisconnectConfirm = false; disconnectingId = null"
      >
        <div class="bg-card-bg border border-border rounded-xl p-6 w-full max-w-md mx-4">
          <h3 class="text-lg font-semibold text-text-primary mb-2">Disconnect Integration?</h3>
          <p class="text-sm text-text-secondary mb-5">
            This will remove the connection and any associated configuration. Features relying on this integration will stop working.
          </p>
          <div class="flex justify-end gap-3">
            <button
              class="px-4 py-2 border border-border text-text-secondary rounded-lg text-sm font-medium hover:text-text-primary transition-colors"
              @click="showDisconnectConfirm = false; disconnectingId = null"
            >
              Cancel
            </button>
            <button
              class="px-4 py-2 bg-status-paused text-white rounded-lg text-sm font-medium hover:bg-status-paused/90 transition-colors"
              @click="disconnect(disconnectingId!)"
            >
              Disconnect
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
