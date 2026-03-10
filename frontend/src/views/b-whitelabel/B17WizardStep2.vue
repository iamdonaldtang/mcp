<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()

// === Types ===
type DeploymentPath = 'domain' | 'iframe' | 'widget' | 'pagebuilder' | 'sdk'

interface DnsRecord {
  type: string
  host: string
  value: string
  status: 'pending' | 'verified' | 'error'
}

// === State ===
const path = ref<DeploymentPath>('domain')
const loading = ref(true)
const saving = ref(false)
const lastSavedAt = ref<string | null>(null)
let autoSaveTimer: ReturnType<typeof setInterval> | null = null

// === Stepper ===
const steps = [
  { num: 1, label: 'Path' },
  { num: 2, label: 'Configure' },
  { num: 3, label: 'Brand' },
  { num: 4, label: 'Preview & Publish' },
]
const currentStep = 2

// === Domain Config ===
const domainForm = reactive({
  domain: '',
  dnsRecords: [
    { type: 'CNAME', host: '', value: 'taskon-wl.xyz', status: 'pending' as DnsRecord['status'] },
  ] as DnsRecord[],
  sslStatus: 'pending' as 'pending' | 'active' | 'error',
})
const verifyingDns = ref(false)
let dnsPollingTimer: ReturnType<typeof setInterval> | null = null

const domainError = computed(() => {
  if (!domainForm.domain) return ''
  const pattern = /^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)*\.[a-z]{2,}$/
  return pattern.test(domainForm.domain) ? '' : 'Enter a valid domain (e.g., community.yourproject.com)'
})

function updateDnsHost() {
  if (domainForm.domain) {
    domainForm.dnsRecords[0].host = domainForm.domain
  }
}

async function verifyDns() {
  verifyingDns.value = true
  try {
    const res = await api.post('/api/v1/whitelabel/domain/verify', { domain: domainForm.domain })
    const data = res.data?.data
    if (data?.verified) {
      domainForm.dnsRecords[0].status = 'verified'
      domainForm.sslStatus = 'active'
      verifyingDns.value = false
    } else {
      // Start polling
      startDnsPolling()
    }
  } catch {
    domainForm.dnsRecords[0].status = 'error'
    verifyingDns.value = false
  }
}

function startDnsPolling() {
  if (dnsPollingTimer) clearInterval(dnsPollingTimer)
  let attempts = 0
  dnsPollingTimer = setInterval(async () => {
    attempts++
    if (attempts > 18) {
      // Stop after 3 minutes
      clearInterval(dnsPollingTimer!)
      dnsPollingTimer = null
      verifyingDns.value = false
      return
    }
    try {
      const res = await api.get(`/api/v1/whitelabel/domain/status?domain=${domainForm.domain}`)
      if (res.data?.data?.verified) {
        domainForm.dnsRecords[0].status = 'verified'
        domainForm.sslStatus = 'active'
        clearInterval(dnsPollingTimer!)
        dnsPollingTimer = null
        verifyingDns.value = false
      }
    } catch {
      // Keep polling
    }
  }, 10000)
}

// === Iframe Config ===
const iframeForm = reactive({
  sourceUrl: '',
  allowedOrigins: '',
  height: 800,
  ssoEnabled: false,
  ssoCallbackUrl: '',
  ssoSharedSecret: '',
})
const testingSso = ref(false)
const ssoTestResult = ref<'success' | 'error' | null>(null)

const iframeEmbedCode = computed(() => {
  const url = iframeForm.sourceUrl || 'https://your-community.taskon.xyz'
  return `<iframe
  src="${url}"
  width="100%"
  height="${iframeForm.height}px"
  frameborder="0"
  allow="clipboard-write"
  style="border: none; border-radius: 8px;"
></iframe>`
})

async function testSso() {
  testingSso.value = true
  ssoTestResult.value = null
  try {
    const res = await api.post('/api/v1/whitelabel/sso/test', {
      callbackUrl: iframeForm.ssoCallbackUrl,
    })
    ssoTestResult.value = res.data?.data?.success ? 'success' : 'error'
  } catch {
    ssoTestResult.value = 'error'
  } finally {
    testingSso.value = false
  }
}

// === Widget Config ===
const widgetForm = reactive({
  modules: {
    leaderboard: true,
    tasks: true,
    userCenter: true,
    shop: false,
    dayChain: false,
  } as Record<string, boolean>,
  theme: 'dark' as 'light' | 'dark',
  primaryColor: '#9B7EE0',
})

const moduleOptions = [
  { key: 'leaderboard', label: 'Leaderboard', icon: 'leaderboard' },
  { key: 'tasks', label: 'Tasks', icon: 'task_alt' },
  { key: 'userCenter', label: 'User Center', icon: 'person' },
  { key: 'shop', label: 'Benefits Shop', icon: 'storefront' },
  { key: 'dayChain', label: 'DayChain', icon: 'local_fire_department' },
]

// === Page Builder Config ===
const pageBuilderForm = reactive({
  pageName: '',
  urlSlug: '',
  template: 'blank' as 'blank' | 'community-hub' | 'rewards-portal',
  theme: 'dark' as 'light' | 'dark',
})

const templateOptions = [
  { id: 'blank', label: 'Blank', description: 'Start from scratch', icon: 'note_add' },
  { id: 'community-hub', label: 'Community Hub', description: 'Pre-built community layout', icon: 'hub' },
  { id: 'rewards-portal', label: 'Rewards Portal', description: 'Rewards-focused layout', icon: 'card_giftcard' },
]

function generateSlug() {
  pageBuilderForm.urlSlug = pageBuilderForm.pageName
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

// === SDK Config ===
const sdkForm = reactive({
  projectId: '',
  apiKey: '',
  apiKeyMasked: '',
  environment: 'production' as 'production' | 'staging',
  webhookUrl: '',
  allowedOrigins: '',
})
const generatingKey = ref(false)
const keyCopied = ref(false)

async function generateApiKey() {
  generatingKey.value = true
  try {
    const res = await api.post('/api/v1/whitelabel/sdk/generate-key', {
      environment: sdkForm.environment,
    })
    const data = res.data?.data
    sdkForm.apiKey = data?.apiKey || ''
    sdkForm.apiKeyMasked = data?.apiKey
      ? data.apiKey.substring(0, 8) + '...' + data.apiKey.substring(data.apiKey.length - 4)
      : ''
  } catch {
    // Error generating key
  } finally {
    generatingKey.value = false
  }
}

async function copyApiKey() {
  if (!sdkForm.apiKey) return
  try {
    await navigator.clipboard.writeText(sdkForm.apiKey)
    keyCopied.value = true
    setTimeout(() => { keyCopied.value = false }, 2000)
  } catch {
    // Clipboard API not available
  }
}

// === Validation per path ===
const isValid = computed(() => {
  switch (path.value) {
    case 'domain':
      return domainForm.domain.length > 0 && !domainError.value
    case 'iframe':
      return true // Source URL is auto-generated
    case 'widget':
      return Object.values(widgetForm.modules).some(Boolean)
    case 'pagebuilder':
      return pageBuilderForm.pageName.length > 0
    case 'sdk':
      return true // Project ID is auto-generated
    default:
      return false
  }
})

// === Draft Persistence ===
async function loadDraft() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/whitelabel/wizard/draft')
    const draft = res.data?.data
    if (draft?.path) {
      path.value = draft.embedSubPath || draft.path
    }
    if (draft?.step2) {
      const s2 = draft.step2
      // Restore based on path
      if (s2.domain) Object.assign(domainForm, s2.domain)
      if (s2.iframe) Object.assign(iframeForm, s2.iframe)
      if (s2.widget) Object.assign(widgetForm, s2.widget)
      if (s2.pagebuilder) Object.assign(pageBuilderForm, s2.pagebuilder)
      if (s2.sdk) Object.assign(sdkForm, s2.sdk)
    }
    // Auto-generate some fields
    if (!iframeForm.sourceUrl && draft?.projectId) {
      iframeForm.sourceUrl = `https://${draft.projectId}.taskon.xyz`
    }
    if (!sdkForm.projectId && draft?.projectId) {
      sdkForm.projectId = draft.projectId
    }
  } catch {
    router.replace('/b/whitelabel/wizard/step-1')
  } finally {
    loading.value = false
  }
}

function getStep2Data() {
  return {
    domain: { ...domainForm },
    iframe: { ...iframeForm },
    widget: { ...widgetForm },
    pagebuilder: { ...pageBuilderForm },
    sdk: { ...sdkForm },
  }
}

async function saveDraft(silent = true) {
  if (saving.value) return
  saving.value = true
  try {
    await api.post('/api/v1/whitelabel/wizard/draft', { step2: getStep2Data() })
    lastSavedAt.value = new Date().toLocaleTimeString()
  } catch {
    if (!silent) { /* toast */ }
  } finally {
    saving.value = false
  }
}

// === Navigation ===
function goBack() {
  const embedPaths = ['iframe', 'widget', 'pagebuilder']
  if (embedPaths.includes(path.value)) {
    router.push('/b/whitelabel/wizard/step-1-5')
  } else {
    router.push('/b/whitelabel/wizard/step-1')
  }
}

async function goNext() {
  if (!isValid.value) return
  await saveDraft(false)
  router.push('/b/whitelabel/wizard/step-3')
}

// === Lifecycle ===
onMounted(async () => {
  await loadDraft()
  autoSaveTimer = setInterval(() => { saveDraft() }, 30000)
})

onUnmounted(() => {
  if (autoSaveTimer) clearInterval(autoSaveTimer)
  if (dnsPollingTimer) clearInterval(dnsPollingTimer)
})
</script>

<template>
  <div class="min-h-screen bg-page-bg">
    <!-- Top Bar -->
    <div class="sticky top-0 z-20 bg-card-bg border-b border-border px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          class="flex items-center gap-1.5 text-sm text-slate-400 hover:text-white transition-colors"
          @click="goBack"
        >
          <span class="material-symbols-rounded text-[20px]">arrow_back</span>
          Back
        </button>
        <div class="w-px h-5 bg-border"></div>
        <h1 class="text-lg font-semibold text-white">Set Up White Label</h1>
      </div>
      <div class="flex items-center gap-3">
        <span v-if="lastSavedAt" class="text-xs text-slate-500">Saved {{ lastSavedAt }}</span>
        <button
          class="px-4 py-2 text-sm font-medium text-slate-300 bg-border rounded-lg hover:bg-[#2D3B4E] transition-colors"
          :disabled="saving"
          @click="saveDraft(false)"
        >
          {{ saving ? 'Saving...' : 'Save Draft' }}
        </button>
      </div>
    </div>

    <!-- Stepper -->
    <div class="px-6 py-6 max-w-3xl mx-auto">
      <div class="flex items-center justify-between">
        <template v-for="(step, idx) in steps" :key="step.num">
          <div class="flex items-center gap-2">
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold transition-colors"
              :class="step.num === currentStep
                ? 'bg-whitelabel text-white'
                : step.num < currentStep
                  ? 'bg-whitelabel/20 text-whitelabel'
                  : 'bg-border text-slate-500'"
            >
              <span v-if="step.num < currentStep" class="material-symbols-rounded text-[18px]">check</span>
              <span v-else>{{ step.num }}</span>
            </div>
            <span
              class="text-sm font-medium"
              :class="step.num <= currentStep ? 'text-white' : 'text-slate-500'"
            >
              {{ step.label }}
            </span>
          </div>
          <div
            v-if="idx < steps.length - 1"
            class="flex-1 h-px mx-3"
            :class="step.num < currentStep ? 'bg-whitelabel/40' : 'bg-border'"
          ></div>
        </template>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="animate-spin w-8 h-8 border-2 border-whitelabel border-t-transparent rounded-full"></div>
    </div>

    <!-- Main Content -->
    <div v-else class="px-6 pb-28 max-w-5xl mx-auto">

      <!-- ========== DOMAIN PATH ========== -->
      <template v-if="path === 'domain'">
        <div class="mb-6">
          <h2 class="text-xl font-semibold text-white mb-1">Domain Configuration</h2>
          <p class="text-sm text-slate-400">Point your custom domain to TaskOn's white label infrastructure.</p>
        </div>

        <div class="bg-card-bg border border-border rounded-xl p-6 space-y-6">
          <!-- Domain Input -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">
              Custom Domain <span class="text-red-400">*</span>
            </label>
            <input
              v-model="domainForm.domain"
              type="text"
              placeholder="community.yourproject.com"
              class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border rounded-lg outline-none transition-colors placeholder:text-slate-600"
              :class="domainError ? 'border-red-500 focus:border-red-500' : 'border-border focus:border-whitelabel'"
              @input="updateDnsHost"
            />
            <p v-if="domainError" class="text-xs text-red-400">{{ domainError }}</p>
          </div>

          <!-- DNS Records Table -->
          <div class="space-y-3">
            <label class="block text-sm font-medium text-slate-300">DNS Records</label>
            <div class="overflow-hidden rounded-lg border border-border">
              <table class="w-full text-sm">
                <thead>
                  <tr class="bg-page-bg">
                    <th class="px-4 py-2.5 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Type</th>
                    <th class="px-4 py-2.5 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Host</th>
                    <th class="px-4 py-2.5 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Value</th>
                    <th class="px-4 py-2.5 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="record in domainForm.dnsRecords" :key="record.type" class="border-t border-border">
                    <td class="px-4 py-3 text-slate-300 font-mono text-xs">{{ record.type }}</td>
                    <td class="px-4 py-3 text-slate-300 font-mono text-xs">{{ record.host || domainForm.domain || '—' }}</td>
                    <td class="px-4 py-3 text-slate-300 font-mono text-xs">{{ record.value }}</td>
                    <td class="px-4 py-3">
                      <span
                        class="inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full"
                        :class="{
                          'bg-status-draft-bg text-status-draft': record.status === 'pending',
                          'bg-status-active-bg text-status-active': record.status === 'verified',
                          'bg-status-paused-bg text-status-paused': record.status === 'error',
                        }"
                      >
                        <span class="material-symbols-rounded text-[14px]">
                          {{ record.status === 'verified' ? 'check_circle' : record.status === 'error' ? 'error' : 'schedule' }}
                        </span>
                        {{ record.status === 'verified' ? 'Verified' : record.status === 'error' ? 'Error' : 'Pending' }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Verify Button -->
          <div class="flex items-center gap-4">
            <button
              class="px-5 py-2.5 text-sm font-medium bg-whitelabel text-white rounded-lg hover:bg-[#8B6ED0] transition-colors flex items-center gap-2 disabled:opacity-50"
              :disabled="!domainForm.domain || !!domainError || verifyingDns"
              @click="verifyDns"
            >
              <span v-if="verifyingDns" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
              <span v-else class="material-symbols-rounded text-[18px]">dns</span>
              {{ verifyingDns ? 'Verifying...' : 'Verify DNS' }}
            </button>
            <!-- SSL Status -->
            <div class="flex items-center gap-2 text-sm">
              <span class="material-symbols-rounded text-[18px]" :class="domainForm.sslStatus === 'active' ? 'text-status-active' : 'text-slate-500'">
                {{ domainForm.sslStatus === 'active' ? 'lock' : 'lock_open' }}
              </span>
              <span :class="domainForm.sslStatus === 'active' ? 'text-status-active' : 'text-slate-500'">
                SSL {{ domainForm.sslStatus === 'active' ? 'Active' : 'Pending' }}
              </span>
            </div>
          </div>
        </div>
      </template>

      <!-- ========== IFRAME PATH ========== -->
      <template v-if="path === 'iframe'">
        <div class="mb-6">
          <h2 class="text-xl font-semibold text-white mb-1">Iframe Configuration</h2>
          <p class="text-sm text-slate-400">Configure iframe embed settings and optional SSO.</p>
        </div>

        <div class="bg-card-bg border border-border rounded-xl p-6 space-y-6">
          <!-- Source URL -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">Source URL</label>
            <input
              v-model="iframeForm.sourceUrl"
              type="text"
              readonly
              class="w-full px-4 py-2.5 text-sm text-slate-400 bg-page-bg border border-border rounded-lg cursor-default"
            />
            <p class="text-xs text-slate-500">Auto-generated based on your project.</p>
          </div>

          <!-- Allowed Origins -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">Allowed Origins</label>
            <textarea
              v-model="iframeForm.allowedOrigins"
              rows="3"
              placeholder="https://yourapp.com&#10;https://staging.yourapp.com"
              class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none resize-none focus:border-whitelabel placeholder:text-slate-600"
            ></textarea>
            <p class="text-xs text-slate-500">One origin per line. Required for cross-origin communication.</p>
          </div>

          <!-- Height -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">Height</label>
            <div class="flex items-center gap-2">
              <input
                v-model.number="iframeForm.height"
                type="number"
                min="300"
                max="2000"
                class="w-32 px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel"
              />
              <span class="text-sm text-slate-500">px</span>
            </div>
          </div>

          <!-- SSO Toggle -->
          <div class="border-t border-border pt-6 space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-slate-300">Single Sign-On (SSO)</p>
                <p class="text-xs text-slate-500 mt-0.5">Let your users authenticate through your app.</p>
              </div>
              <button
                class="relative w-11 h-6 rounded-full transition-colors"
                :class="iframeForm.ssoEnabled ? 'bg-whitelabel' : 'bg-border'"
                @click="iframeForm.ssoEnabled = !iframeForm.ssoEnabled"
              >
                <div
                  class="absolute top-0.5 w-5 h-5 bg-white rounded-full shadow transition-transform"
                  :class="iframeForm.ssoEnabled ? 'translate-x-[22px]' : 'translate-x-0.5'"
                ></div>
              </button>
            </div>

            <template v-if="iframeForm.ssoEnabled">
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Callback URL</label>
                <input
                  v-model="iframeForm.ssoCallbackUrl"
                  type="text"
                  placeholder="https://yourapp.com/auth/callback"
                  class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel placeholder:text-slate-600"
                />
              </div>
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Shared Secret</label>
                <input
                  v-model="iframeForm.ssoSharedSecret"
                  type="password"
                  placeholder="Your shared secret key"
                  class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel placeholder:text-slate-600"
                />
              </div>
              <button
                class="px-4 py-2 text-sm font-medium bg-whitelabel/15 text-whitelabel rounded-lg hover:bg-whitelabel/25 transition-colors flex items-center gap-2 disabled:opacity-50"
                :disabled="!iframeForm.ssoCallbackUrl || testingSso"
                @click="testSso"
              >
                <span v-if="testingSso" class="animate-spin w-4 h-4 border-2 border-whitelabel border-t-transparent rounded-full"></span>
                <span v-else class="material-symbols-rounded text-[18px]">verified_user</span>
                {{ testingSso ? 'Testing...' : 'Test SSO' }}
              </button>
              <div v-if="ssoTestResult" class="flex items-center gap-2 text-sm">
                <span
                  class="material-symbols-rounded text-[18px]"
                  :class="ssoTestResult === 'success' ? 'text-status-active' : 'text-status-paused'"
                >
                  {{ ssoTestResult === 'success' ? 'check_circle' : 'error' }}
                </span>
                <span :class="ssoTestResult === 'success' ? 'text-status-active' : 'text-status-paused'">
                  {{ ssoTestResult === 'success' ? 'SSO test passed' : 'SSO test failed — check callback URL and secret' }}
                </span>
              </div>
            </template>
          </div>

          <!-- Embed Code Preview -->
          <div class="border-t border-border pt-6 space-y-2">
            <label class="block text-sm font-medium text-slate-300">Embed Code</label>
            <pre class="w-full px-4 py-3 text-xs text-slate-300 bg-page-bg border border-border rounded-lg overflow-x-auto font-mono leading-relaxed">{{ iframeEmbedCode }}</pre>
          </div>
        </div>
      </template>

      <!-- ========== WIDGET PATH ========== -->
      <template v-if="path === 'widget'">
        <div class="mb-6">
          <h2 class="text-xl font-semibold text-white mb-1">Widget Configuration</h2>
          <p class="text-sm text-slate-400">Select which community modules to include and customize the appearance.</p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
          <!-- Left: Config -->
          <div class="lg:col-span-3 space-y-6">
            <!-- Module Selector -->
            <div class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
              <h3 class="text-base font-semibold text-white">Community Modules</h3>
              <div class="space-y-3">
                <label
                  v-for="mod in moduleOptions"
                  :key="mod.key"
                  class="flex items-center gap-3 p-3 rounded-lg border transition-colors cursor-pointer"
                  :class="widgetForm.modules[mod.key]
                    ? 'border-whitelabel/40 bg-whitelabel/5'
                    : 'border-border hover:border-[#2D3B4E]'"
                >
                  <input
                    v-model="widgetForm.modules[mod.key]"
                    type="checkbox"
                    class="sr-only"
                  />
                  <div
                    class="w-5 h-5 rounded border-2 flex items-center justify-center shrink-0 transition-colors"
                    :class="widgetForm.modules[mod.key]
                      ? 'border-whitelabel bg-whitelabel'
                      : 'border-[#334155]'"
                  >
                    <span v-if="widgetForm.modules[mod.key]" class="material-symbols-rounded text-[14px] text-white">check</span>
                  </div>
                  <span class="material-symbols-rounded text-[20px] text-slate-400">{{ mod.icon }}</span>
                  <span class="text-sm font-medium text-slate-200">{{ mod.label }}</span>
                </label>
              </div>
            </div>

            <!-- Theme & Color -->
            <div class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
              <h3 class="text-base font-semibold text-white">Appearance</h3>
              <!-- Theme -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Theme</label>
                <div class="flex gap-3">
                  <button
                    v-for="theme in (['light', 'dark'] as const)"
                    :key="theme"
                    class="flex-1 px-4 py-2.5 text-sm font-medium rounded-lg border transition-colors capitalize"
                    :class="widgetForm.theme === theme
                      ? 'border-whitelabel bg-whitelabel/10 text-whitelabel'
                      : 'border-border text-slate-400 hover:border-[#2D3B4E]'"
                    @click="widgetForm.theme = theme"
                  >
                    {{ theme }}
                  </button>
                </div>
              </div>
              <!-- Primary Color -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Primary Color</label>
                <div class="flex items-center gap-3">
                  <input
                    v-model="widgetForm.primaryColor"
                    type="color"
                    class="w-10 h-10 rounded-lg border border-border cursor-pointer bg-transparent"
                  />
                  <input
                    v-model="widgetForm.primaryColor"
                    type="text"
                    maxlength="7"
                    class="w-28 px-3 py-2 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel font-mono"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Right: Preview -->
          <div class="lg:col-span-2">
            <div class="sticky top-24">
              <h3 class="text-sm font-medium text-slate-400 mb-3">Widget Preview</h3>
              <div
                class="rounded-xl border border-border overflow-hidden"
                :class="widgetForm.theme === 'dark' ? 'bg-card-bg' : 'bg-white'"
              >
                <div class="h-1.5" :style="{ backgroundColor: widgetForm.primaryColor }"></div>
                <div class="p-4 space-y-3">
                  <template v-for="mod in moduleOptions" :key="mod.key">
                    <div
                      v-if="widgetForm.modules[mod.key]"
                      class="flex items-center gap-2 px-3 py-2.5 rounded-lg border text-sm"
                      :class="widgetForm.theme === 'dark'
                        ? 'border-border bg-page-bg text-slate-300'
                        : 'border-slate-200 bg-slate-50 text-slate-600'"
                    >
                      <span
                        class="material-symbols-rounded text-[18px]"
                        :style="{ color: widgetForm.primaryColor }"
                      >{{ mod.icon }}</span>
                      {{ mod.label }}
                    </div>
                  </template>
                  <div v-if="!Object.values(widgetForm.modules).some(Boolean)" class="text-center py-6">
                    <span class="material-symbols-rounded text-[32px] text-slate-600">widgets</span>
                    <p class="text-xs text-slate-500 mt-2">Select modules to preview</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- ========== PAGE BUILDER PATH ========== -->
      <template v-if="path === 'pagebuilder'">
        <div class="mb-6">
          <h2 class="text-xl font-semibold text-white mb-1">Page Setup</h2>
          <p class="text-sm text-slate-400">Create the first page for your White Label community.</p>
        </div>

        <div class="bg-card-bg border border-border rounded-xl p-6 space-y-6">
          <!-- Page Name -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">
              Page Name <span class="text-red-400">*</span>
            </label>
            <input
              v-model="pageBuilderForm.pageName"
              type="text"
              placeholder="e.g. Community Hub"
              class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel placeholder:text-slate-600"
              @input="generateSlug"
            />
          </div>

          <!-- URL Slug -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">URL Slug</label>
            <div class="flex items-center gap-0 overflow-hidden rounded-lg border border-border">
              <span class="px-3 py-2.5 text-sm text-slate-500 bg-page-bg border-r border-border shrink-0">yourproject.com/</span>
              <input
                v-model="pageBuilderForm.urlSlug"
                type="text"
                class="flex-1 px-3 py-2.5 text-sm text-white bg-page-bg outline-none"
              />
            </div>
          </div>

          <!-- Template Selector -->
          <div class="space-y-3">
            <label class="block text-sm font-medium text-slate-300">Template</label>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
              <button
                v-for="tpl in templateOptions"
                :key="tpl.id"
                class="text-left p-4 rounded-lg border transition-colors"
                :class="pageBuilderForm.template === tpl.id
                  ? 'border-whitelabel bg-whitelabel/5'
                  : 'border-border hover:border-[#2D3B4E]'"
                @click="pageBuilderForm.template = tpl.id as typeof pageBuilderForm.template"
              >
                <span
                  class="material-symbols-rounded text-[24px] mb-2 block"
                  :class="pageBuilderForm.template === tpl.id ? 'text-whitelabel' : 'text-slate-500'"
                >{{ tpl.icon }}</span>
                <p class="text-sm font-medium text-white">{{ tpl.label }}</p>
                <p class="text-xs text-slate-500 mt-0.5">{{ tpl.description }}</p>
              </button>
            </div>
          </div>

          <!-- Theme -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">Theme</label>
            <div class="flex gap-3">
              <button
                v-for="theme in (['light', 'dark'] as const)"
                :key="theme"
                class="px-5 py-2.5 text-sm font-medium rounded-lg border transition-colors capitalize"
                :class="pageBuilderForm.theme === theme
                  ? 'border-whitelabel bg-whitelabel/10 text-whitelabel'
                  : 'border-border text-slate-400 hover:border-[#2D3B4E]'"
                @click="pageBuilderForm.theme = theme"
              >
                {{ theme }}
              </button>
            </div>
          </div>
        </div>
      </template>

      <!-- ========== SDK PATH ========== -->
      <template v-if="path === 'sdk'">
        <div class="mb-6">
          <h2 class="text-xl font-semibold text-white mb-1">SDK Configuration</h2>
          <p class="text-sm text-slate-400">Generate API keys and configure SDK access for your integration.</p>
        </div>

        <div class="bg-card-bg border border-border rounded-xl p-6 space-y-6">
          <!-- Project ID -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">Project ID</label>
            <div class="flex items-center gap-2">
              <input
                :value="sdkForm.projectId || 'Loading...'"
                type="text"
                readonly
                class="flex-1 px-4 py-2.5 text-sm text-slate-400 bg-page-bg border border-border rounded-lg cursor-default font-mono"
              />
            </div>
            <p class="text-xs text-slate-500">Auto-generated. Use this in your SDK initialization.</p>
          </div>

          <!-- API Key -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">API Key</label>
            <div class="flex items-center gap-2">
              <input
                :value="sdkForm.apiKeyMasked || 'Not yet generated'"
                type="text"
                readonly
                class="flex-1 px-4 py-2.5 text-sm text-slate-400 bg-page-bg border border-border rounded-lg cursor-default font-mono"
              />
              <button
                v-if="sdkForm.apiKey"
                class="px-3 py-2.5 text-sm font-medium bg-border rounded-lg hover:bg-[#2D3B4E] transition-colors flex items-center gap-1"
                :class="keyCopied ? 'text-status-active' : 'text-slate-300'"
                @click="copyApiKey"
              >
                <span class="material-symbols-rounded text-[16px]">{{ keyCopied ? 'check' : 'content_copy' }}</span>
                {{ keyCopied ? 'Copied' : 'Copy' }}
              </button>
              <button
                class="px-4 py-2.5 text-sm font-medium bg-whitelabel text-white rounded-lg hover:bg-[#8B6ED0] transition-colors flex items-center gap-1.5 disabled:opacity-50"
                :disabled="generatingKey"
                @click="generateApiKey"
              >
                <span v-if="generatingKey" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
                <span v-else class="material-symbols-rounded text-[16px]">key</span>
                {{ generatingKey ? 'Generating...' : sdkForm.apiKey ? 'Regenerate' : 'Generate Key' }}
              </button>
            </div>
          </div>

          <!-- Environment -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">Environment</label>
            <div class="flex gap-3">
              <button
                v-for="env in (['production', 'staging'] as const)"
                :key="env"
                class="px-5 py-2.5 text-sm font-medium rounded-lg border transition-colors capitalize"
                :class="sdkForm.environment === env
                  ? 'border-whitelabel bg-whitelabel/10 text-whitelabel'
                  : 'border-border text-slate-400 hover:border-[#2D3B4E]'"
                @click="sdkForm.environment = env"
              >
                {{ env }}
              </button>
            </div>
          </div>

          <!-- Webhook URL -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">
              Webhook URL <span class="text-xs text-slate-500 font-normal">(optional)</span>
            </label>
            <input
              v-model="sdkForm.webhookUrl"
              type="text"
              placeholder="https://yourapp.com/webhooks/taskon"
              class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel placeholder:text-slate-600"
            />
          </div>

          <!-- Allowed Origins -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">Allowed Origins</label>
            <textarea
              v-model="sdkForm.allowedOrigins"
              rows="3"
              placeholder="https://yourapp.com&#10;https://staging.yourapp.com"
              class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none resize-none focus:border-whitelabel placeholder:text-slate-600"
            ></textarea>
            <p class="text-xs text-slate-500">One origin per line.</p>
          </div>
        </div>
      </template>
    </div>

    <!-- Bottom Action Bar -->
    <div class="fixed bottom-0 left-0 right-0 z-20 bg-card-bg border-t border-border px-6 py-4">
      <div class="max-w-5xl mx-auto flex items-center justify-between">
        <button
          class="px-5 py-2.5 text-sm font-medium text-slate-300 bg-border rounded-lg hover:bg-[#2D3B4E] transition-colors flex items-center gap-2"
          @click="goBack"
        >
          <span class="material-symbols-rounded text-[18px]">arrow_back</span>
          Back
        </button>
        <button
          class="px-6 py-2.5 text-sm font-semibold rounded-lg transition-all flex items-center gap-2"
          :class="isValid
            ? 'bg-whitelabel text-white hover:bg-[#8B6ED0]'
            : 'bg-border text-slate-500 cursor-not-allowed'"
          :disabled="!isValid"
          @click="goNext"
        >
          Next: Brand Settings
          <span class="material-symbols-rounded text-[18px]">arrow_forward</span>
        </button>
      </div>
    </div>
  </div>
</template>
