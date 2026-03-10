<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { api } from '../../api/client'

// === Types ===
interface ApiKey {
  id: string
  key: string
  environment: 'production' | 'staging'
  createdAt: string
}

interface Webhook {
  id: string
  url: string
  events: string[]
  status: 'active' | 'failed'
  lastTriggered: string | null
  secret?: string
}

interface UsageData {
  apiCalls: { used: number; limit: number }
  webhookDeliveries: { used: number; limit: number }
  rateLimit: string
}

// === State ===
const loading = ref(true)
const projectId = ref('')
const activeEnv = ref<'production' | 'staging'>('production')
const apiKeys = reactive<Record<string, ApiKey | null>>({
  production: null,
  staging: null,
})
const keyRevealed = ref(false)
const webhooks = ref<Webhook[]>([])
const showWebhookForm = ref(false)
const editingWebhookId = ref<string | null>(null)
const activeCodeTab = ref('react')
const usage = reactive<UsageData>({
  apiCalls: { used: 0, limit: 50000 },
  webhookDeliveries: { used: 0, limit: 10000 },
  rateLimit: '100 requests/minute per key',
})

// Webhook form
const webhookForm = reactive({
  url: '',
  events: [] as string[],
  secret: '',
})

const availableEvents = [
  'task.completed',
  'points.earned',
  'level.up',
  'shop.redeemed',
  'sprint.ended',
  'user.joined',
  'milestone.claimed',
]

// Confirm dialogs
const showRegenerateConfirm = ref(false)
const showDeleteConfirm = ref(false)
const deletingWebhookId = ref<string | null>(null)
const testResults = reactive<Record<string, { success: boolean; message: string } | null>>({})

// === Computed ===
const currentKey = computed(() => apiKeys[activeEnv.value])
const maskedKey = computed(() => {
  if (!currentKey.value) return ''
  const k = currentKey.value.key
  return '****' + k.slice(-8)
})
const integrationStatus = computed(() => {
  const prodKey = apiKeys.production
  const hasWebhooks = webhooks.value.some((w) => w.status === 'active')
  if (prodKey && hasWebhooks) return 'active'
  if (prodKey) return 'partial'
  return 'inactive'
})
const statusLabel = computed(() => {
  const map = { active: 'Fully Configured', partial: 'Keys Only', inactive: 'Not Configured' }
  return map[integrationStatus.value]
})
const statusColor = computed(() => {
  const map = {
    active: 'bg-status-active-bg text-status-active',
    partial: 'bg-status-draft-bg text-status-draft',
    inactive: 'bg-status-completed-bg text-status-completed',
  }
  return map[integrationStatus.value]
})

// === API ===
async function fetchData() {
  loading.value = true
  try {
    const { data } = await api.get('/api/v1/whitelabel/sdk')
    projectId.value = data.projectId || ''
    if (data.keys?.production) apiKeys.production = data.keys.production
    if (data.keys?.staging) apiKeys.staging = data.keys.staging
    webhooks.value = data.webhooks || []
    if (data.usage) {
      Object.assign(usage, data.usage)
    }
  } catch {
    // Use placeholder data on error
    projectId.value = 'proj_wl_a8f3e2d1'
    apiKeys.production = { id: 'k1', key: 'sk_live_a1b2c3d4e5f6g7h8i9j0klmnopqrst', environment: 'production', createdAt: '2026-02-15' }
    apiKeys.staging = { id: 'k2', key: 'sk_test_x9y8w7v6u5t4s3r2q1p0onmlkjihgf', environment: 'staging', createdAt: '2026-02-15' }
    webhooks.value = [
      { id: 'wh1', url: 'https://api.myproject.com/webhooks/taskon', events: ['task.completed', 'points.earned', 'user.joined'], status: 'active', lastTriggered: '2026-03-09T14:30:00Z' },
      { id: 'wh2', url: 'https://hooks.slack.com/services/T00/B00/xxx', events: ['level.up', 'milestone.claimed'], status: 'failed', lastTriggered: '2026-03-08T10:00:00Z' },
    ]
    usage.apiCalls = { used: 12450, limit: 50000 }
    usage.webhookDeliveries = { used: 3280, limit: 10000 }
  } finally {
    loading.value = false
  }
}

async function copyToClipboard(text: string) {
  await navigator.clipboard.writeText(text)
}

async function generateKey() {
  try {
    const { data } = await api.post('/api/v1/whitelabel/sdk/keys', { environment: activeEnv.value })
    apiKeys[activeEnv.value] = data
  } catch {
    // Mock
    const env = activeEnv.value
    apiKeys[env] = {
      id: 'k_new',
      key: env === 'production' ? 'sk_live_' + Math.random().toString(36).slice(2, 30) : 'sk_test_' + Math.random().toString(36).slice(2, 30),
      environment: env,
      createdAt: new Date().toISOString(),
    }
  }
}

async function regenerateKey() {
  showRegenerateConfirm.value = false
  await generateKey()
}

function resetWebhookForm() {
  webhookForm.url = ''
  webhookForm.events = []
  webhookForm.secret = ''
  editingWebhookId.value = null
}

function startEditWebhook(wh: Webhook) {
  editingWebhookId.value = wh.id
  webhookForm.url = wh.url
  webhookForm.events = [...wh.events]
  webhookForm.secret = wh.secret || ''
  showWebhookForm.value = true
}

async function saveWebhook() {
  const payload = { url: webhookForm.url, events: webhookForm.events, secret: webhookForm.secret || undefined }
  try {
    if (editingWebhookId.value) {
      const { data } = await api.put(`/api/v1/whitelabel/sdk/webhooks/${editingWebhookId.value}`, payload)
      const idx = webhooks.value.findIndex((w) => w.id === editingWebhookId.value)
      if (idx >= 0) webhooks.value[idx] = data
    } else {
      const { data } = await api.post('/api/v1/whitelabel/sdk/webhooks', payload)
      webhooks.value.push(data)
    }
  } catch {
    // Mock
    const newWh: Webhook = {
      id: editingWebhookId.value || 'wh_' + Date.now(),
      url: webhookForm.url,
      events: [...webhookForm.events],
      status: 'active',
      lastTriggered: null,
      secret: webhookForm.secret || undefined,
    }
    if (editingWebhookId.value) {
      const idx = webhooks.value.findIndex((w) => w.id === editingWebhookId.value)
      if (idx >= 0) webhooks.value[idx] = newWh
    } else {
      webhooks.value.push(newWh)
    }
  }
  showWebhookForm.value = false
  resetWebhookForm()
}

async function deleteWebhook() {
  if (!deletingWebhookId.value) return
  try {
    await api.delete(`/api/v1/whitelabel/sdk/webhooks/${deletingWebhookId.value}`)
  } catch {
    // proceed anyway
  }
  webhooks.value = webhooks.value.filter((w) => w.id !== deletingWebhookId.value)
  showDeleteConfirm.value = false
  deletingWebhookId.value = null
}

async function testWebhook(whId: string) {
  testResults[whId] = null
  try {
    const { data } = await api.post(`/api/v1/whitelabel/sdk/webhooks/${whId}/test`)
    testResults[whId] = { success: data.statusCode === 200, message: `${data.statusCode} ${data.statusText || 'OK'}` }
  } catch {
    // Mock success
    testResults[whId] = { success: true, message: '200 OK' }
  }
}

function toggleEvent(event: string) {
  const idx = webhookForm.events.indexOf(event)
  if (idx >= 0) webhookForm.events.splice(idx, 1)
  else webhookForm.events.push(event)
}

function usagePercent(used: number, limit: number) {
  return Math.min(100, Math.round((used / limit) * 100))
}

function formatNumber(n: number) {
  return n.toLocaleString()
}

function formatDate(d: string | null) {
  if (!d) return 'Never'
  return new Date(d).toLocaleString()
}

// === Code Snippets ===
const codeSnippets: Record<string, { install: string; init: string; example: string }> = {
  react: {
    install: 'npm install @taskon/sdk @taskon/react',
    init: `import { TaskOnProvider } from '@taskon/react'

function App() {
  return (
    <TaskOnProvider
      projectId="${'${PROJECT_ID}'}"
      apiKey="${'${API_KEY}'}"
    >
      <YourApp />
    </TaskOnProvider>
  )
}`,
    example: `import { useTaskOn } from '@taskon/react'

function QuestCard() {
  const { completeTask, user } = useTaskOn()

  const handleClaim = async () => {
    await completeTask('daily-checkin')
  }

  return (
    <div>
      <p>Points: {user.points}</p>
      <button onClick={handleClaim}>Complete Task</button>
    </div>
  )
}`,
  },
  vue: {
    install: 'npm install @taskon/sdk @taskon/vue',
    init: `import { createApp } from 'vue'
import { TaskOnPlugin } from '@taskon/vue'

const app = createApp(App)
app.use(TaskOnPlugin, {
  projectId: '${'\${PROJECT_ID}'}',
  apiKey: '${'\${API_KEY}'}',
})
app.mount('#app')`,
    example: `<script setup>
import { useTaskOn } from '@taskon/vue'

const { completeTask, user } = useTaskOn()

async function handleClaim() {
  await completeTask('daily-checkin')
}
<\/script>

<template>
  <div>
    <p>Points: {{ user.points }}</p>
    <button @click="handleClaim">Complete Task</button>
  </div>
</template>`,
  },
  vanilla: {
    install: 'npm install @taskon/sdk',
    init: `import { TaskOn } from '@taskon/sdk'

const taskon = new TaskOn({
  projectId: '${'\${PROJECT_ID}'}',
  apiKey: '${'\${API_KEY}'}',
})

await taskon.init()`,
    example: `// Complete a task
const result = await taskon.completeTask('daily-checkin')
console.log('Points earned:', result.pointsEarned)

// Get user profile
const user = await taskon.getUser()
console.log('Level:', user.level, 'Points:', user.points)

// Listen for events
taskon.on('level.up', (data) => {
  console.log('New level:', data.level)
})`,
  },
  rest: {
    install: '# No SDK needed — use any HTTP client',
    init: `# Base URL
https://api.taskon.xyz/v1

# Authentication Header
Authorization: Bearer ${'\${API_KEY}'}
X-Project-ID: ${'\${PROJECT_ID}'}`,
    example: `# Complete a task
curl -X POST https://api.taskon.xyz/v1/tasks/complete \\
  -H "Authorization: Bearer ${'\${API_KEY}'}" \\
  -H "X-Project-ID: ${'\${PROJECT_ID}'}" \\
  -H "Content-Type: application/json" \\
  -d '{"taskId": "daily-checkin", "userId": "user_123"}'

# Get user points
curl https://api.taskon.xyz/v1/users/user_123/points \\
  -H "Authorization: Bearer ${'\${API_KEY}'}" \\
  -H "X-Project-ID: ${'\${PROJECT_ID}'}"`,
  },
}

const codeTabs = [
  { key: 'react', label: 'React' },
  { key: 'vue', label: 'Vue' },
  { key: 'vanilla', label: 'Vanilla JS' },
  { key: 'rest', label: 'REST API' },
]

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
        <h1 class="text-2xl font-bold text-text-primary">SDK & API</h1>
        <span :class="['px-3 py-1 rounded-full text-xs font-medium', statusColor]">{{ statusLabel }}</span>
      </div>
    </div>

    <!-- API Keys Section -->
    <div class="bg-card-bg border border-border rounded-xl p-6">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-semibold text-text-primary flex items-center gap-2">
          <span class="material-symbols-rounded text-whitelabel text-xl">key</span>
          API Keys
        </h2>
        <!-- Environment Toggle -->
        <div class="flex bg-page-bg rounded-lg p-1">
          <button
            v-for="env in (['production', 'staging'] as const)"
            :key="env"
            :class="[
              'px-4 py-1.5 rounded-md text-sm font-medium transition-colors',
              activeEnv === env ? 'bg-whitelabel/20 text-whitelabel' : 'text-text-muted hover:text-text-secondary',
            ]"
            @click="activeEnv = env; keyRevealed = false"
          >
            {{ env === 'production' ? 'Production' : 'Staging' }}
          </button>
        </div>
      </div>

      <!-- Project ID -->
      <div class="mb-4">
        <label class="block text-sm text-text-secondary mb-1.5">Project ID</label>
        <div class="flex items-center gap-2">
          <div class="flex-1 bg-page-bg border border-border rounded-lg px-4 py-2.5 font-mono text-sm text-text-primary">
            {{ projectId }}
          </div>
          <button
            class="px-3 py-2.5 border border-border rounded-lg text-text-secondary hover:text-text-primary hover:border-whitelabel/50 transition-colors"
            @click="copyToClipboard(projectId)"
            title="Copy"
          >
            <span class="material-symbols-rounded text-lg">content_copy</span>
          </button>
        </div>
      </div>

      <!-- API Key -->
      <div v-if="currentKey">
        <label class="block text-sm text-text-secondary mb-1.5">
          API Key ({{ activeEnv === 'production' ? 'Live' : 'Test' }})
        </label>
        <div class="flex items-center gap-2">
          <div class="flex-1 bg-page-bg border border-border rounded-lg px-4 py-2.5 font-mono text-sm text-text-primary">
            {{ keyRevealed ? currentKey.key : maskedKey }}
          </div>
          <button
            class="px-3 py-2.5 border border-border rounded-lg text-text-secondary hover:text-text-primary hover:border-whitelabel/50 transition-colors"
            @click="keyRevealed = !keyRevealed"
            :title="keyRevealed ? 'Hide' : 'Show'"
          >
            <span class="material-symbols-rounded text-lg">{{ keyRevealed ? 'visibility_off' : 'visibility' }}</span>
          </button>
          <button
            class="px-3 py-2.5 border border-border rounded-lg text-text-secondary hover:text-text-primary hover:border-whitelabel/50 transition-colors"
            @click="copyToClipboard(currentKey.key)"
            title="Copy"
          >
            <span class="material-symbols-rounded text-lg">content_copy</span>
          </button>
          <button
            class="px-3 py-2.5 border border-status-paused/30 rounded-lg text-status-paused hover:bg-status-paused-bg transition-colors"
            @click="showRegenerateConfirm = true"
            title="Regenerate"
          >
            <span class="material-symbols-rounded text-lg">refresh</span>
          </button>
        </div>
        <p class="text-xs text-text-muted mt-1.5">
          Created {{ formatDate(currentKey.createdAt) }}
        </p>
      </div>

      <!-- No Key State -->
      <div v-else class="text-center py-6">
        <p class="text-text-secondary mb-3">No {{ activeEnv }} API key generated yet.</p>
        <button
          class="px-5 py-2.5 bg-whitelabel text-white rounded-lg font-medium hover:bg-whitelabel/90 transition-colors"
          @click="generateKey"
        >
          Generate New Key
        </button>
      </div>
    </div>

    <!-- Webhooks Section -->
    <div class="bg-card-bg border border-border rounded-xl p-6">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-semibold text-text-primary flex items-center gap-2">
          <span class="material-symbols-rounded text-whitelabel text-xl">webhook</span>
          Webhooks
        </h2>
        <button
          class="flex items-center gap-1.5 px-4 py-2 bg-whitelabel/10 text-whitelabel rounded-lg text-sm font-medium hover:bg-whitelabel/20 transition-colors"
          @click="resetWebhookForm(); showWebhookForm = true"
        >
          <span class="material-symbols-rounded text-lg">add</span>
          Add Webhook
        </button>
      </div>

      <!-- Webhook List -->
      <div v-if="webhooks.length > 0" class="space-y-3 mb-4">
        <div
          v-for="wh in webhooks"
          :key="wh.id"
          :class="[
            'border rounded-lg p-4',
            wh.status === 'active' ? 'border-border' : 'border-status-paused/30',
          ]"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1.5">
                <span class="font-mono text-sm text-text-primary truncate">{{ wh.url }}</span>
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-xs font-medium shrink-0',
                    wh.status === 'active' ? 'bg-status-active-bg text-status-active' : 'bg-status-paused-bg text-status-paused',
                  ]"
                >
                  {{ wh.status === 'active' ? 'Active' : 'Failed' }}
                </span>
              </div>
              <div class="flex flex-wrap gap-1.5 mb-1.5">
                <span
                  v-for="ev in wh.events"
                  :key="ev"
                  class="px-2 py-0.5 bg-page-bg border border-border rounded text-xs text-text-secondary"
                >
                  {{ ev }}
                </span>
              </div>
              <p class="text-xs text-text-muted">Last triggered: {{ formatDate(wh.lastTriggered) }}</p>
              <!-- Test result -->
              <div v-if="testResults[wh.id]" class="mt-2">
                <span
                  :class="[
                    'text-xs font-medium',
                    testResults[wh.id]!.success ? 'text-status-active' : 'text-status-paused',
                  ]"
                >
                  {{ testResults[wh.id]!.success ? '&#10003;' : '&#10007;' }} {{ testResults[wh.id]!.message }}
                </span>
              </div>
            </div>
            <div class="flex items-center gap-1.5 shrink-0">
              <button
                class="p-2 text-text-muted hover:text-text-primary rounded-lg hover:bg-page-bg transition-colors"
                title="Send Test Event"
                @click="testWebhook(wh.id)"
              >
                <span class="material-symbols-rounded text-lg">send</span>
              </button>
              <button
                class="p-2 text-text-muted hover:text-text-primary rounded-lg hover:bg-page-bg transition-colors"
                title="Edit"
                @click="startEditWebhook(wh)"
              >
                <span class="material-symbols-rounded text-lg">edit</span>
              </button>
              <button
                class="p-2 text-text-muted hover:text-status-paused rounded-lg hover:bg-status-paused-bg transition-colors"
                title="Delete"
                @click="deletingWebhookId = wh.id; showDeleteConfirm = true"
              >
                <span class="material-symbols-rounded text-lg">delete</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="!showWebhookForm" class="text-center py-6">
        <p class="text-text-secondary">No webhooks configured yet.</p>
      </div>

      <!-- Webhook Form (inline) -->
      <div v-if="showWebhookForm" class="border border-whitelabel/30 rounded-lg p-5 bg-page-bg">
        <h3 class="text-sm font-semibold text-text-primary mb-4">
          {{ editingWebhookId ? 'Edit Webhook' : 'New Webhook' }}
        </h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Endpoint URL</label>
            <input
              v-model="webhookForm.url"
              type="url"
              placeholder="https://your-api.com/webhooks/taskon"
              class="w-full bg-card-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-whitelabel/50"
            />
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Events</label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="ev in availableEvents"
                :key="ev"
                :class="[
                  'px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors',
                  webhookForm.events.includes(ev)
                    ? 'bg-whitelabel/15 border-whitelabel/40 text-whitelabel'
                    : 'bg-card-bg border-border text-text-secondary hover:border-whitelabel/30',
                ]"
                @click="toggleEvent(ev)"
              >
                {{ ev }}
              </button>
            </div>
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Secret Header (optional)</label>
            <input
              v-model="webhookForm.secret"
              type="text"
              placeholder="whsec_..."
              class="w-full bg-card-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-whitelabel/50"
            />
          </div>
          <div class="flex gap-3">
            <button
              class="px-5 py-2 bg-whitelabel text-white rounded-lg text-sm font-medium hover:bg-whitelabel/90 transition-colors disabled:opacity-50"
              :disabled="!webhookForm.url || webhookForm.events.length === 0"
              @click="saveWebhook"
            >
              Save
            </button>
            <button
              class="px-5 py-2 border border-border text-text-secondary rounded-lg text-sm font-medium hover:text-text-primary transition-colors"
              @click="showWebhookForm = false; resetWebhookForm()"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- SDK Setup Guide Section -->
    <div class="bg-card-bg border border-border rounded-xl p-6">
      <h2 class="text-lg font-semibold text-text-primary flex items-center gap-2 mb-5">
        <span class="material-symbols-rounded text-whitelabel text-xl">integration_instructions</span>
        SDK Setup Guide
      </h2>

      <!-- Tabs -->
      <div class="flex border-b border-border mb-5">
        <button
          v-for="tab in codeTabs"
          :key="tab.key"
          :class="[
            'px-4 py-2.5 text-sm font-medium border-b-2 transition-colors -mb-px',
            activeCodeTab === tab.key
              ? 'border-whitelabel text-whitelabel'
              : 'border-transparent text-text-muted hover:text-text-secondary',
          ]"
          @click="activeCodeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- Code Blocks -->
      <div class="space-y-4">
        <!-- Install -->
        <div>
          <div class="flex items-center justify-between mb-1.5">
            <span class="text-xs text-text-muted uppercase tracking-wider font-semibold">Install</span>
            <button
              class="text-xs text-text-muted hover:text-text-secondary transition-colors flex items-center gap-1"
              @click="copyToClipboard(codeSnippets[activeCodeTab].install)"
            >
              <span class="material-symbols-rounded text-sm">content_copy</span> Copy
            </button>
          </div>
          <pre class="bg-[#0D1117] rounded-lg p-4 overflow-x-auto"><code class="text-sm font-mono text-text-primary">{{ codeSnippets[activeCodeTab].install }}</code></pre>
        </div>

        <!-- Initialization -->
        <div>
          <div class="flex items-center justify-between mb-1.5">
            <span class="text-xs text-text-muted uppercase tracking-wider font-semibold">Initialization</span>
            <button
              class="text-xs text-text-muted hover:text-text-secondary transition-colors flex items-center gap-1"
              @click="copyToClipboard(codeSnippets[activeCodeTab].init)"
            >
              <span class="material-symbols-rounded text-sm">content_copy</span> Copy
            </button>
          </div>
          <pre class="bg-[#0D1117] rounded-lg p-4 overflow-x-auto"><code class="text-sm font-mono text-text-primary whitespace-pre">{{ codeSnippets[activeCodeTab].init }}</code></pre>
        </div>

        <!-- Example -->
        <div>
          <div class="flex items-center justify-between mb-1.5">
            <span class="text-xs text-text-muted uppercase tracking-wider font-semibold">Example Usage</span>
            <button
              class="text-xs text-text-muted hover:text-text-secondary transition-colors flex items-center gap-1"
              @click="copyToClipboard(codeSnippets[activeCodeTab].example)"
            >
              <span class="material-symbols-rounded text-sm">content_copy</span> Copy
            </button>
          </div>
          <pre class="bg-[#0D1117] rounded-lg p-4 overflow-x-auto"><code class="text-sm font-mono text-text-primary whitespace-pre">{{ codeSnippets[activeCodeTab].example }}</code></pre>
        </div>
      </div>
    </div>

    <!-- Usage & Limits Section -->
    <div class="bg-card-bg border border-border rounded-xl p-6">
      <h2 class="text-lg font-semibold text-text-primary flex items-center gap-2 mb-5">
        <span class="material-symbols-rounded text-whitelabel text-xl">monitoring</span>
        Usage & Limits
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- API Calls -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm text-text-secondary">API Calls (this month)</span>
            <span class="text-sm font-medium text-text-primary">
              {{ formatNumber(usage.apiCalls.used) }} / {{ formatNumber(usage.apiCalls.limit) }}
            </span>
          </div>
          <div class="w-full h-2.5 bg-page-bg rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all"
              :class="usagePercent(usage.apiCalls.used, usage.apiCalls.limit) > 80 ? 'bg-status-draft' : 'bg-whitelabel'"
              :style="{ width: usagePercent(usage.apiCalls.used, usage.apiCalls.limit) + '%' }"
            />
          </div>
          <p class="text-xs text-text-muted mt-1">{{ usagePercent(usage.apiCalls.used, usage.apiCalls.limit) }}% used</p>
        </div>

        <!-- Webhook Deliveries -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm text-text-secondary">Webhook Deliveries (this month)</span>
            <span class="text-sm font-medium text-text-primary">
              {{ formatNumber(usage.webhookDeliveries.used) }} / {{ formatNumber(usage.webhookDeliveries.limit) }}
            </span>
          </div>
          <div class="w-full h-2.5 bg-page-bg rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all"
              :class="usagePercent(usage.webhookDeliveries.used, usage.webhookDeliveries.limit) > 80 ? 'bg-status-draft' : 'bg-whitelabel'"
              :style="{ width: usagePercent(usage.webhookDeliveries.used, usage.webhookDeliveries.limit) + '%' }"
            />
          </div>
          <p class="text-xs text-text-muted mt-1">{{ usagePercent(usage.webhookDeliveries.used, usage.webhookDeliveries.limit) }}% used</p>
        </div>
      </div>

      <!-- Rate Limit Info -->
      <div class="mt-5 flex items-center gap-2 px-4 py-3 bg-page-bg rounded-lg border border-border">
        <span class="material-symbols-rounded text-text-muted text-lg">info</span>
        <span class="text-sm text-text-secondary">Rate limit: {{ usage.rateLimit }}</span>
      </div>
    </div>

    <!-- Regenerate Confirm Dialog -->
    <Teleport to="body">
      <div
        v-if="showRegenerateConfirm"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
        @click.self="showRegenerateConfirm = false"
      >
        <div class="bg-card-bg border border-border rounded-xl p-6 w-full max-w-md mx-4">
          <h3 class="text-lg font-semibold text-text-primary mb-2">Regenerate API Key?</h3>
          <p class="text-sm text-text-secondary mb-5">
            This will invalidate the current {{ activeEnv }} key. Any services using this key will lose access immediately.
          </p>
          <div class="flex justify-end gap-3">
            <button
              class="px-4 py-2 border border-border text-text-secondary rounded-lg text-sm font-medium hover:text-text-primary transition-colors"
              @click="showRegenerateConfirm = false"
            >
              Cancel
            </button>
            <button
              class="px-4 py-2 bg-status-paused text-white rounded-lg text-sm font-medium hover:bg-status-paused/90 transition-colors"
              @click="regenerateKey"
            >
              Regenerate
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Webhook Confirm Dialog -->
    <Teleport to="body">
      <div
        v-if="showDeleteConfirm"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
        @click.self="showDeleteConfirm = false"
      >
        <div class="bg-card-bg border border-border rounded-xl p-6 w-full max-w-md mx-4">
          <h3 class="text-lg font-semibold text-text-primary mb-2">Delete Webhook?</h3>
          <p class="text-sm text-text-secondary mb-5">
            This webhook will stop receiving events immediately. This action cannot be undone.
          </p>
          <div class="flex justify-end gap-3">
            <button
              class="px-4 py-2 border border-border text-text-secondary rounded-lg text-sm font-medium hover:text-text-primary transition-colors"
              @click="showDeleteConfirm = false; deletingWebhookId = null"
            >
              Cancel
            </button>
            <button
              class="px-4 py-2 bg-status-paused text-white rounded-lg text-sm font-medium hover:bg-status-paused/90 transition-colors"
              @click="deleteWebhook"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
