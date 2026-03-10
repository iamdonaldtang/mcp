<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const projectId = computed(() => route.params.projectId as string)

const loading = ref(true)
const project = ref({
  name: '',
  email: '',
  sdk_version: '1.0.0',
  api_base: 'https://api.taskon.xyz',
})

// Verification state
const verificationDomain = ref('')
const verificationStatus = ref<'idle' | 'testing' | 'success' | 'error'>('idle')
const verificationChecks = ref([
  { key: 'sdk', label: 'SDK loaded', passed: false },
  { key: 'api', label: 'API responding', passed: false },
  { key: 'widget', label: 'Widget rendering', passed: false },
  { key: 'sso', label: 'SSO working', passed: false },
])

// Code tab
const codeTab = ref<'react' | 'vue' | 'vanilla'>('react')

const allPassed = computed(() => verificationChecks.value.every(c => c.passed))

// --- API ---
onMounted(async () => {
  try {
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
    const res = await axios.get(`${baseUrl}/api/devkit/${projectId.value}`)
    if (res.data.data) {
      project.value = { ...project.value, ...res.data.data }
    }
  } catch { /* use defaults */ }
  finally { loading.value = false }
})

async function runVerification() {
  verificationStatus.value = 'testing'
  verificationChecks.value.forEach(c => { c.passed = false })

  try {
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
    const res = await axios.post(`${baseUrl}/api/devkit/${projectId.value}/verify`, {
      domain: verificationDomain.value,
    })
    const results = res.data.data?.checks || {}
    verificationChecks.value.forEach(c => {
      c.passed = !!results[c.key]
    })
    verificationStatus.value = allPassed.value ? 'success' : 'error'
  } catch {
    verificationStatus.value = 'error'
  }
}

function copyText(text: string) {
  navigator.clipboard.writeText(text)
}

// --- Code snippets ---
const installCmd = 'npm install @taskon/sdk'

const initCode = computed(() =>
`import { TaskOnSDK } from '@taskon/sdk'

const taskon = new TaskOnSDK({
  projectId: '${projectId.value}',
  apiBase: '${project.value.api_base}',
})

taskon.init()`)

const widgetCode = computed(() =>
`<!-- Add this where you want the widget to appear -->
<div id="taskon-widget"></div>

<script>
  taskon.mount('#taskon-widget', {
    modules: ['tasks', 'leaderboard', 'shop'],
    theme: 'dark',
  })
<\/script>`)

const reactCode = computed(() =>
`import { TaskOnProvider, TaskWidget } from '@taskon/react'

function App() {
  return (
    <TaskOnProvider projectId="${projectId.value}">
      <TaskWidget
        modules={['tasks', 'leaderboard', 'shop']}
        theme="dark"
        onUserAction={(event) => console.log(event)}
      />
    </TaskOnProvider>
  )
}

export default App`)

const vueCode = computed(() => {
  const pkg = '@taskon/vue'
  return `<template>
  <TaskOnProvider :project-id="'${projectId.value}'">
    <TaskWidget
      :modules="['tasks', 'leaderboard', 'shop']"
      theme="dark"
      @user-action="onUserAction"
    />
  </TaskOnProvider>
</template>

<script setup>
import { TaskOnProvider, TaskWidget } from '${pkg}'

function onUserAction(event) {
  console.log(event)
}
<\/script>`
})

const vanillaCode = computed(() =>
`<script src="https://cdn.taskon.xyz/sdk/${project.value.sdk_version}/taskon.min.js"><\/script>

<div id="taskon-widget"></div>

<script>
  const taskon = new TaskOnSDK({
    projectId: '${projectId.value}',
  })

  taskon.mount('#taskon-widget', {
    modules: ['tasks', 'leaderboard', 'shop'],
    theme: 'dark',
  })
<\/script>`)

const walletSsoCode = computed(() =>
`// Wallet-Based SSO (Recommended)
taskon.auth.connectWallet({
  providers: ['metamask', 'walletconnect', 'coinbase'],
  onSuccess: (user) => {
    console.log('Connected:', user.address)
  },
  onError: (err) => {
    console.error('Auth failed:', err)
  },
})`)

const oauthSsoCode = computed(() =>
`// OAuth 2.0 SSO
taskon.auth.configureOAuth({
  clientId: 'YOUR_CLIENT_ID',
  clientSecret: 'YOUR_CLIENT_SECRET',
  redirectUri: 'https://yourapp.com/callback',
  scopes: ['profile', 'wallet'],
})

// Trigger login
taskon.auth.login('oauth')`)
</script>

<template>
  <!-- Full-page standalone layout (no sidebar) -->
  <div class="min-h-screen bg-page-bg">
    <!-- Header -->
    <header class="border-b border-border bg-card-bg">
      <div class="max-w-5xl mx-auto px-8 py-4 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <!-- TaskOn Logo -->
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 bg-wl rounded-lg flex items-center justify-center">
              <span class="text-white text-sm font-bold">T</span>
            </div>
            <span class="text-lg font-bold text-text-primary">TaskOn</span>
          </div>
          <span class="px-2.5 py-0.5 text-xs font-medium rounded-full bg-wl/10 text-wl border border-wl/20">
            Developer Kit
          </span>
        </div>
        <div v-if="project.name" class="text-sm text-text-secondary">
          <span class="text-text-muted">Project:</span> {{ project.name }}
        </div>
      </div>
    </header>

    <!-- Loading -->
    <div v-if="loading" class="max-w-5xl mx-auto px-8 py-20 text-center">
      <span class="material-symbols-rounded text-4xl text-wl animate-spin block mb-4">progress_activity</span>
      <p class="text-text-muted">Loading developer kit...</p>
    </div>

    <!-- Content -->
    <main v-else class="max-w-5xl mx-auto px-8 py-10 space-y-12">
      <!-- Hero -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-text-primary mb-3">
          Integration Guide for {{ project.name || 'Your Project' }}
        </h1>
        <p class="text-lg text-text-secondary mb-2">
          Everything you need to integrate TaskOn White Label
        </p>
        <div class="inline-flex items-center gap-2 px-3 py-1.5 bg-card-bg border border-border rounded-full">
          <span class="material-symbols-rounded text-sm text-text-muted">schedule</span>
          <span class="text-sm text-text-secondary">Estimated time: ~30 minutes</span>
        </div>
      </div>

      <!-- Section 1: Quick Start -->
      <section>
        <div class="flex items-center gap-3 mb-6">
          <div class="w-8 h-8 rounded-full bg-wl/10 flex items-center justify-center">
            <span class="text-wl font-bold text-sm">1</span>
          </div>
          <h2 class="text-xl font-bold text-text-primary">Quick Start</h2>
        </div>

        <div class="space-y-6">
          <!-- Step 1: Install -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <span class="w-6 h-6 rounded-full bg-wl text-white text-xs font-bold flex items-center justify-center">1</span>
              <h3 class="text-sm font-semibold text-text-primary">Install SDK</h3>
            </div>
            <div class="relative bg-page-bg rounded-lg p-4">
              <code class="text-sm text-text-primary font-mono">{{ installCmd }}</code>
              <button
                class="absolute top-2 right-2 p-1.5 rounded hover:bg-white/5 transition-colors"
                title="Copy"
                @click="copyText(installCmd)"
              >
                <span class="material-symbols-rounded text-sm text-text-muted">content_copy</span>
              </button>
            </div>
          </div>

          <!-- Step 2: Initialize -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <span class="w-6 h-6 rounded-full bg-wl text-white text-xs font-bold flex items-center justify-center">2</span>
              <h3 class="text-sm font-semibold text-text-primary">Initialize</h3>
            </div>
            <div class="relative bg-page-bg rounded-lg p-4 overflow-x-auto">
              <pre class="text-sm text-text-primary font-mono whitespace-pre">{{ initCode }}</pre>
              <button
                class="absolute top-2 right-2 p-1.5 rounded hover:bg-white/5 transition-colors"
                title="Copy"
                @click="copyText(initCode)"
              >
                <span class="material-symbols-rounded text-sm text-text-muted">content_copy</span>
              </button>
            </div>
          </div>

          <!-- Step 3: Add Widget -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <span class="w-6 h-6 rounded-full bg-wl text-white text-xs font-bold flex items-center justify-center">3</span>
              <h3 class="text-sm font-semibold text-text-primary">Add Widget</h3>
            </div>
            <div class="relative bg-page-bg rounded-lg p-4 overflow-x-auto">
              <pre class="text-sm text-text-primary font-mono whitespace-pre">{{ widgetCode }}</pre>
              <button
                class="absolute top-2 right-2 p-1.5 rounded hover:bg-white/5 transition-colors"
                title="Copy"
                @click="copyText(widgetCode)"
              >
                <span class="material-symbols-rounded text-sm text-text-muted">content_copy</span>
              </button>
            </div>
          </div>

          <!-- Step 4: Verify -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-3">
              <span class="w-6 h-6 rounded-full bg-wl text-white text-xs font-bold flex items-center justify-center">4</span>
              <h3 class="text-sm font-semibold text-text-primary">Verify Integration</h3>
            </div>
            <p class="text-sm text-text-secondary mb-3">Run verification after deploying to check everything works.</p>
            <button
              class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors disabled:opacity-50"
              :disabled="verificationStatus === 'testing'"
              @click="runVerification"
            >
              {{ verificationStatus === 'testing' ? 'Verifying...' : 'Run Verification' }}
            </button>
          </div>
        </div>
      </section>

      <!-- Section 2: Integration Code -->
      <section>
        <div class="flex items-center gap-3 mb-6">
          <div class="w-8 h-8 rounded-full bg-wl/10 flex items-center justify-center">
            <span class="text-wl font-bold text-sm">2</span>
          </div>
          <h2 class="text-xl font-bold text-text-primary">Integration Code</h2>
        </div>

        <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <!-- Tabs -->
          <div class="flex border-b border-border">
            <button
              v-for="tab in (['react', 'vue', 'vanilla'] as const)"
              :key="tab"
              class="px-5 py-3 text-sm font-medium transition-colors border-b-2"
              :class="codeTab === tab
                ? 'border-wl text-wl'
                : 'border-transparent text-text-muted hover:text-text-primary'"
              @click="codeTab = tab"
            >
              {{ tab === 'react' ? 'React' : tab === 'vue' ? 'Vue' : 'Vanilla JS' }}
            </button>
          </div>

          <!-- Code blocks -->
          <div class="relative p-5">
            <div class="bg-page-bg rounded-lg p-4 overflow-x-auto">
              <pre class="text-sm text-text-primary font-mono whitespace-pre">{{ codeTab === 'react' ? reactCode : codeTab === 'vue' ? vueCode : vanillaCode }}</pre>
            </div>
            <button
              class="absolute top-7 right-7 p-1.5 rounded hover:bg-white/5 transition-colors"
              title="Copy"
              @click="copyText(codeTab === 'react' ? reactCode : codeTab === 'vue' ? vueCode : vanillaCode)"
            >
              <span class="material-symbols-rounded text-sm text-text-muted">content_copy</span>
            </button>
          </div>
        </div>
      </section>

      <!-- Section 3: SSO Setup -->
      <section>
        <div class="flex items-center gap-3 mb-6">
          <div class="w-8 h-8 rounded-full bg-wl/10 flex items-center justify-center">
            <span class="text-wl font-bold text-sm">3</span>
          </div>
          <h2 class="text-xl font-bold text-text-primary">SSO Setup</h2>
        </div>

        <div class="grid grid-cols-2 gap-6">
          <!-- Wallet-Based -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-1">
              <span class="material-symbols-rounded text-lg text-wl">account_balance_wallet</span>
              <h3 class="text-sm font-semibold text-text-primary">Wallet-Based</h3>
              <span class="px-2 py-0.5 text-xs rounded bg-status-active-bg text-status-active">Recommended</span>
            </div>
            <p class="text-xs text-text-secondary mb-4">Users connect their Web3 wallet. Simplest integration with zero backend changes.</p>

            <div class="space-y-3 mb-4">
              <div class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-wl/10 text-wl text-xs font-bold flex items-center justify-center mt-0.5 shrink-0">1</span>
                <p class="text-xs text-text-secondary">SDK auto-detects wallet providers (MetaMask, WalletConnect, Coinbase)</p>
              </div>
              <div class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-wl/10 text-wl text-xs font-bold flex items-center justify-center mt-0.5 shrink-0">2</span>
                <p class="text-xs text-text-secondary">User signs a message to verify ownership</p>
              </div>
              <div class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-wl/10 text-wl text-xs font-bold flex items-center justify-center mt-0.5 shrink-0">3</span>
                <p class="text-xs text-text-secondary">TaskOn creates/links user account automatically</p>
              </div>
            </div>

            <div class="relative bg-page-bg rounded-lg p-3 overflow-x-auto">
              <pre class="text-xs text-text-primary font-mono whitespace-pre">{{ walletSsoCode }}</pre>
              <button
                class="absolute top-1 right-1 p-1 rounded hover:bg-white/5"
                @click="copyText(walletSsoCode)"
              >
                <span class="material-symbols-rounded text-xs text-text-muted">content_copy</span>
              </button>
            </div>
          </div>

          <!-- OAuth 2.0 -->
          <div class="bg-card-bg border border-border rounded-xl p-5">
            <div class="flex items-center gap-2 mb-1">
              <span class="material-symbols-rounded text-lg text-text-secondary">key</span>
              <h3 class="text-sm font-semibold text-text-primary">OAuth 2.0</h3>
            </div>
            <p class="text-xs text-text-secondary mb-4">Traditional OAuth flow for apps with existing user systems. Requires backend setup.</p>

            <div class="space-y-3 mb-4">
              <div class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-card-bg border border-border text-text-muted text-xs font-bold flex items-center justify-center mt-0.5 shrink-0">1</span>
                <p class="text-xs text-text-secondary">Register OAuth app in TaskOn dashboard to get client credentials</p>
              </div>
              <div class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-card-bg border border-border text-text-muted text-xs font-bold flex items-center justify-center mt-0.5 shrink-0">2</span>
                <p class="text-xs text-text-secondary">Configure redirect URI and scopes</p>
              </div>
              <div class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-card-bg border border-border text-text-muted text-xs font-bold flex items-center justify-center mt-0.5 shrink-0">3</span>
                <p class="text-xs text-text-secondary">Handle authorization callback on your server</p>
              </div>
            </div>

            <div class="relative bg-page-bg rounded-lg p-3 overflow-x-auto">
              <pre class="text-xs text-text-primary font-mono whitespace-pre">{{ oauthSsoCode }}</pre>
              <button
                class="absolute top-1 right-1 p-1 rounded hover:bg-white/5"
                @click="copyText(oauthSsoCode)"
              >
                <span class="material-symbols-rounded text-xs text-text-muted">content_copy</span>
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- Section 4: Self-Service Verification -->
      <section>
        <div class="flex items-center gap-3 mb-6">
          <div class="w-8 h-8 rounded-full bg-wl/10 flex items-center justify-center">
            <span class="text-wl font-bold text-sm">4</span>
          </div>
          <h2 class="text-xl font-bold text-text-primary">Self-Service Verification</h2>
        </div>

        <div class="bg-card-bg border border-border rounded-xl p-6">
          <!-- Domain input -->
          <div class="mb-6">
            <label class="block text-sm text-text-secondary mb-1.5">Deployment Domain</label>
            <div class="flex gap-3">
              <input
                v-model="verificationDomain"
                type="text"
                placeholder="https://yourapp.com"
                class="flex-1 px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-wl focus:outline-none"
              />
              <button
                class="px-5 py-2.5 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors disabled:opacity-50"
                :disabled="!verificationDomain.trim() || verificationStatus === 'testing'"
                @click="runVerification"
              >
                <span v-if="verificationStatus === 'testing'" class="flex items-center gap-2">
                  <span class="material-symbols-rounded text-sm animate-spin">progress_activity</span>
                  Testing...
                </span>
                <span v-else>Test Connection</span>
              </button>
            </div>
          </div>

          <!-- Checklist -->
          <div class="space-y-3">
            <div
              v-for="check in verificationChecks"
              :key="check.key"
              class="flex items-center gap-3 px-4 py-3 rounded-lg border transition-colors"
              :class="check.passed ? 'border-status-active/20 bg-status-active-bg' : 'border-border bg-page-bg'"
            >
              <span
                class="material-symbols-rounded text-lg"
                :class="check.passed ? 'text-status-active' : 'text-text-muted'"
              >
                {{ check.passed ? 'check_circle' : 'radio_button_unchecked' }}
              </span>
              <span class="text-sm" :class="check.passed ? 'text-status-active' : 'text-text-secondary'">
                {{ check.label }}
              </span>
              <span v-if="check.passed" class="text-xs text-status-active ml-auto">Passed</span>
            </div>
          </div>

          <!-- Success banner -->
          <div
            v-if="allPassed && verificationStatus === 'success'"
            class="mt-6 flex items-center gap-3 p-4 rounded-xl bg-status-active-bg border border-status-active/20"
          >
            <span class="material-symbols-rounded text-2xl text-status-active">celebration</span>
            <div>
              <h3 class="text-base font-semibold text-status-active">Integration Complete!</h3>
              <p class="text-sm text-status-active/80">All checks passed. Your integration is ready for production.</p>
            </div>
          </div>

          <!-- Error banner -->
          <div
            v-if="verificationStatus === 'error' && !allPassed"
            class="mt-6 flex items-center gap-3 p-4 rounded-xl bg-[#2D1515] border border-status-paused/20"
          >
            <span class="material-symbols-rounded text-2xl text-status-paused">error</span>
            <div>
              <h3 class="text-base font-semibold text-status-paused">Verification Failed</h3>
              <p class="text-sm text-status-paused/80">Some checks did not pass. Please review and fix the issues above.</p>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="border-t border-border bg-card-bg mt-12">
      <div class="max-w-5xl mx-auto px-8 py-6 flex items-center justify-between">
        <div class="text-sm text-text-muted">
          Need help? Contact
          <span v-if="project.email" class="text-wl">{{ project.email }}</span>
          <span v-else class="text-wl">support@taskon.xyz</span>
        </div>
        <div class="flex items-center gap-2">
          <div class="w-6 h-6 bg-wl rounded flex items-center justify-center">
            <span class="text-white text-xs font-bold">T</span>
          </div>
          <span class="text-xs text-text-muted">Powered by TaskOn</span>
        </div>
      </div>
    </footer>
  </div>
</template>
