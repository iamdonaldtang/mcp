<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'

interface DnsRecord {
  type: string
  host: string
  value: string
  status: 'verified' | 'pending' | 'error'
  ttl: string
}

interface DomainConfig {
  domain: string
  ssl_status: 'active' | 'provisioning' | 'error' | 'none'
  ssl_expiry?: string
  dns_records: DnsRecord[]
  verified: boolean
  verified_at?: string
  redirect_https: boolean
  custom_404_url: string
  force_www: boolean
}

const loading = ref(true)
const saving = ref(false)
const verifying = ref(false)
const verifyTimer = ref<ReturnType<typeof setInterval> | null>(null)
const verifyElapsed = ref(0)

const config = ref<DomainConfig>({
  domain: '',
  ssl_status: 'none',
  dns_records: [],
  verified: false,
  redirect_https: true,
  custom_404_url: '',
  force_www: false,
})

const form = ref({
  domain: '',
  redirect_https: true,
  custom_404_url: '',
  force_www: false,
})

const advancedOpen = ref(false)

const domainStatus = computed(() => {
  if (config.value.verified) return 'verified'
  if (config.value.domain) return 'pending'
  return 'not_configured'
})

const domainStatusConfig = computed(() => {
  switch (domainStatus.value) {
    case 'verified':
      return { label: 'Verified', bg: '#0A2E1A', text: '#16A34A' }
    case 'pending':
      return { label: 'Pending Verification', bg: '#1F1A08', text: '#D97706' }
    default:
      return { label: 'Not Configured', bg: '#1E293B', text: '#64748B' }
  }
})

const sslStatusConfig = computed(() => {
  switch (config.value.ssl_status) {
    case 'active':
      return { label: 'Active', icon: 'verified_user', color: '#16A34A' }
    case 'provisioning':
      return { label: 'Provisioning...', icon: 'hourglass_top', color: '#D97706' }
    case 'error':
      return { label: 'Error', icon: 'error', color: '#DC2626' }
    default:
      return { label: 'Not Provisioned', icon: 'shield', color: '#64748B' }
  }
})

onMounted(async () => {
  await fetchDomain()
  loading.value = false
})

onUnmounted(() => {
  if (verifyTimer.value) clearInterval(verifyTimer.value)
})

async function fetchDomain() {
  try {
    const res = await api.get('/api/v1/whitelabel/domain')
    const data = res.data.data
    if (data) {
      config.value = data
      form.value = {
        domain: data.domain || '',
        redirect_https: data.redirect_https ?? true,
        custom_404_url: data.custom_404_url || '',
        force_www: data.force_www ?? false,
      }
    }
  } catch { /* empty */ }
}

async function saveDomain() {
  if (!form.value.domain.trim()) return
  saving.value = true
  try {
    await api.put('/api/v1/whitelabel/domain', {
      domain: form.value.domain,
      redirect_https: form.value.redirect_https,
      custom_404_url: form.value.custom_404_url,
      force_www: form.value.force_www,
    })
    await fetchDomain()
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function verifyDns() {
  verifying.value = true
  verifyElapsed.value = 0

  try {
    await api.post('/api/v1/whitelabel/domain/verify')
  } catch { /* empty */ }

  // Poll every 10s, max 5 min (300s)
  verifyTimer.value = setInterval(async () => {
    verifyElapsed.value += 10
    try {
      const res = await api.get('/api/v1/whitelabel/domain')
      const data = res.data.data
      if (data) {
        config.value = data
        if (data.verified || verifyElapsed.value >= 300) {
          stopVerify()
        }
      }
    } catch {
      stopVerify()
    }
  }, 10000)
}

function stopVerify() {
  verifying.value = false
  if (verifyTimer.value) {
    clearInterval(verifyTimer.value)
    verifyTimer.value = null
  }
}

function dnsStatusIcon(status: string) {
  return status === 'verified' ? 'check_circle' : 'hourglass_top'
}

function dnsStatusColor(status: string) {
  return status === 'verified' ? '#16A34A' : '#D97706'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-bold text-text-primary">Domain Setup</h1>
        <span
          class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
          :style="{ backgroundColor: domainStatusConfig.bg, color: domainStatusConfig.text }"
        >
          {{ domainStatusConfig.label }}
        </span>
      </div>
    </div>

    <!-- Loading skeleton -->
    <template v-if="loading">
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <div class="space-y-4">
          <div class="h-5 w-48 bg-border rounded animate-pulse"></div>
          <div class="h-10 bg-border rounded animate-pulse"></div>
          <div class="h-32 bg-border rounded animate-pulse"></div>
        </div>
      </div>
    </template>

    <template v-else>
      <!-- Current Domain Display (if configured) -->
      <div v-if="config.domain" class="bg-card-bg border border-border rounded-xl p-6">
        <h2 class="text-sm font-semibold text-text-muted uppercase tracking-wider mb-4">Current Domain</h2>
        <div class="grid grid-cols-3 gap-6">
          <div>
            <p class="text-xs text-text-muted mb-1">Domain URL</p>
            <p class="text-sm font-medium text-text-primary flex items-center gap-2">
              <span class="material-symbols-rounded text-base text-wl">language</span>
              {{ config.domain }}
            </p>
          </div>
          <div>
            <p class="text-xs text-text-muted mb-1">SSL Certificate</p>
            <p class="text-sm font-medium flex items-center gap-2">
              <span class="material-symbols-rounded text-base" :style="{ color: sslStatusConfig.color }">{{ sslStatusConfig.icon }}</span>
              <span :style="{ color: sslStatusConfig.color }">{{ sslStatusConfig.label }}</span>
            </p>
          </div>
          <div>
            <p class="text-xs text-text-muted mb-1">Last Verified</p>
            <p class="text-sm text-text-secondary">
              {{ config.verified_at ? new Date(config.verified_at).toLocaleDateString() : 'Never' }}
            </p>
          </div>
        </div>
      </div>

      <!-- Configuration Form -->
      <div class="bg-card-bg border border-border rounded-xl p-6 space-y-6">
        <h2 class="text-sm font-semibold text-text-muted uppercase tracking-wider">Configuration</h2>

        <!-- Custom Domain input -->
        <div>
          <label class="block text-sm text-text-secondary mb-1.5">Custom Domain</label>
          <input
            v-model="form.domain"
            type="text"
            placeholder="community.yourproject.com"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-wl focus:outline-none"
          />
          <p class="text-xs text-text-muted mt-1.5">Enter the subdomain or domain you want to use for your White Label community.</p>
        </div>

        <!-- DNS Records Table -->
        <div>
          <div class="flex items-center justify-between mb-3">
            <label class="text-sm font-medium text-text-secondary">DNS Records</label>
            <button
              class="px-4 py-2 text-sm font-medium rounded-lg transition-colors flex items-center gap-2 disabled:opacity-50"
              :class="verifying
                ? 'bg-amber-900/30 text-amber-400 cursor-wait'
                : 'bg-wl text-white hover:bg-wl/90'"
              :disabled="verifying || !config.domain"
              @click="verifyDns"
            >
              <span v-if="verifying" class="material-symbols-rounded text-base animate-spin">progress_activity</span>
              <span v-else class="material-symbols-rounded text-base">dns</span>
              {{ verifying ? `Verifying... (${verifyElapsed}s)` : 'Verify DNS' }}
            </button>
          </div>

          <div class="border border-border rounded-xl overflow-hidden">
            <table class="w-full">
              <thead>
                <tr class="border-b border-border bg-page-bg">
                  <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Type</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Host</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Value</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Status</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-20">TTL</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="config.dns_records.length === 0">
                  <td colspan="5" class="px-4 py-8 text-center text-sm text-text-muted">
                    Enter a domain and save to generate DNS records.
                  </td>
                </tr>
                <tr
                  v-for="(record, idx) in config.dns_records"
                  :key="idx"
                  class="border-b border-border last:border-b-0"
                >
                  <td class="px-4 py-3">
                    <span class="px-2 py-0.5 text-xs font-mono rounded bg-page-bg text-text-secondary">{{ record.type }}</span>
                  </td>
                  <td class="px-4 py-3 text-sm text-text-primary font-mono">{{ record.host }}</td>
                  <td class="px-4 py-3 text-sm text-text-secondary font-mono">{{ record.value }}</td>
                  <td class="px-4 py-3">
                    <span class="flex items-center gap-1.5">
                      <span class="material-symbols-rounded text-base" :style="{ color: dnsStatusColor(record.status) }">
                        {{ dnsStatusIcon(record.status) }}
                      </span>
                      <span class="text-xs capitalize" :style="{ color: dnsStatusColor(record.status) }">{{ record.status }}</span>
                    </span>
                  </td>
                  <td class="px-4 py-3 text-sm text-text-muted font-mono">{{ record.ttl }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p class="text-xs text-text-muted mt-2">
            Add these records in your DNS provider's dashboard. Propagation may take up to 48 hours.
          </p>
        </div>

        <!-- SSL Certificate -->
        <div class="bg-page-bg border border-border rounded-xl p-4">
          <div class="flex items-center gap-3">
            <span class="material-symbols-rounded text-xl" :style="{ color: sslStatusConfig.color }">{{ sslStatusConfig.icon }}</span>
            <div>
              <p class="text-sm font-medium text-text-primary">SSL Certificate</p>
              <p class="text-xs text-text-muted">
                Auto-provisioned via Let's Encrypt after DNS verification.
                <span :style="{ color: sslStatusConfig.color }">Status: {{ sslStatusConfig.label }}</span>
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Advanced Settings -->
      <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
        <button
          class="w-full flex items-center justify-between px-6 py-4 hover:bg-white/2 transition-colors"
          @click="advancedOpen = !advancedOpen"
        >
          <span class="text-sm font-semibold text-text-secondary">Advanced Settings</span>
          <span class="material-symbols-rounded text-lg text-text-muted transition-transform" :class="advancedOpen ? 'rotate-180' : ''">
            keyboard_arrow_down
          </span>
        </button>

        <div v-if="advancedOpen" class="px-6 pb-6 space-y-5 border-t border-border pt-5">
          <!-- Redirect HTTP to HTTPS -->
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-primary">Redirect HTTP to HTTPS</p>
              <p class="text-xs text-text-muted">Automatically redirect insecure requests to HTTPS</p>
            </div>
            <button
              class="relative w-11 h-6 rounded-full transition-colors"
              :class="form.redirect_https ? 'bg-wl' : 'bg-border'"
              @click="form.redirect_https = !form.redirect_https"
            >
              <span
                class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform"
                :class="form.redirect_https ? 'translate-x-5' : ''"
              ></span>
            </button>
          </div>

          <!-- Custom 404 Page URL -->
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Custom 404 Page URL</label>
            <input
              v-model="form.custom_404_url"
              type="text"
              placeholder="https://yourproject.com/404"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-wl focus:outline-none"
            />
            <p class="text-xs text-text-muted mt-1">Leave blank to use the default TaskOn 404 page.</p>
          </div>

          <!-- Force WWW -->
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-primary">Force WWW Prefix</p>
              <p class="text-xs text-text-muted">Redirect non-www requests to the www version of your domain</p>
            </div>
            <button
              class="relative w-11 h-6 rounded-full transition-colors"
              :class="form.force_www ? 'bg-wl' : 'bg-border'"
              @click="form.force_www = !form.force_www"
            >
              <span
                class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform"
                :class="form.force_www ? 'translate-x-5' : ''"
              ></span>
            </button>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <div class="flex justify-end">
        <button
          class="px-6 py-2.5 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors disabled:opacity-50 flex items-center gap-2"
          :disabled="saving || !form.domain.trim()"
          @click="saveDomain"
        >
          <span v-if="saving" class="material-symbols-rounded text-base animate-spin">progress_activity</span>
          {{ saving ? 'Saving...' : 'Save Domain Settings' }}
        </button>
      </div>
    </template>
  </div>
</template>
