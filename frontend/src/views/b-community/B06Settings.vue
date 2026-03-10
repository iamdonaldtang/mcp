<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../../api/client'

const route = useRoute()

const tabs = ['Profile', 'Organization', 'Billing', 'API Keys', 'Notifications'] as const
type Tab = (typeof tabs)[number]

const activeTab = ref<Tab>((route.query.tab as Tab) || 'Profile')
const loading = ref(true)
const saving = ref(false)
const saveMessage = ref('')

// --- Profile ---
interface Profile {
  avatar: string
  initials: string
  email: string
  emailVerified: boolean
  displayName: string
  timezone: string
}

const profile = ref<Profile>({
  avatar: '',
  initials: 'U',
  email: '',
  emailVerified: false,
  displayName: '',
  timezone: 'UTC',
})

const timezones = [
  'UTC', 'America/New_York', 'America/Chicago', 'America/Denver', 'America/Los_Angeles',
  'Europe/London', 'Europe/Berlin', 'Europe/Paris', 'Asia/Tokyo', 'Asia/Shanghai',
  'Asia/Singapore', 'Asia/Dubai', 'Australia/Sydney', 'Pacific/Auckland',
]

// --- Organization ---
interface TeamMember {
  id: string
  name: string
  email: string
  role: 'owner' | 'admin' | 'member'
  status: 'active' | 'invited'
}

interface Organization {
  name: string
  logo: string
  website: string
  industry: string
  teamMembers: TeamMember[]
}

const org = ref<Organization>({
  name: '',
  logo: '',
  website: '',
  industry: '',
  teamMembers: [],
})

const industries = ['DeFi', 'NFT', 'Gaming', 'DAO', 'Infrastructure', 'Other']
const showInviteForm = ref(false)
const inviteEmail = ref('')
const inviteRole = ref<'admin' | 'member'>('member')

// --- Billing ---
interface BillingData {
  planName: string
  planPrice: string
  renewalDate: string
  paymentMethod: string
  cardLast4: string
  history: { id: string; date: string; description: string; amount: string; status: 'paid' | 'pending' }[]
  usage: {
    apiCalls: { used: number; limit: number }
    storage: { used: number; limit: number }
    teamMembers: { used: number; limit: number }
  }
}

const billing = ref<BillingData>({
  planName: '',
  planPrice: '',
  renewalDate: '',
  paymentMethod: '',
  cardLast4: '',
  history: [],
  usage: {
    apiCalls: { used: 0, limit: 1 },
    storage: { used: 0, limit: 1 },
    teamMembers: { used: 0, limit: 1 },
  },
})

// --- API Keys ---
interface ApiKey {
  id: string
  name: string
  key: string
  created: string
  lastUsed: string
  visible: boolean
}

const apiKeys = ref<ApiKey[]>([])
const showCreateKey = ref(false)
const newKeyName = ref('')
const newKeyPermissions = ref<string[]>([])
const permissionOptions = ['Read', 'Write', 'Admin', 'Analytics', 'Campaigns']

// --- Notifications ---
interface NotificationPrefs {
  email: {
    campaignCompleted: boolean
    newSignUp: boolean
    weeklyDigest: boolean
    billingAlerts: boolean
  }
  inApp: {
    campaignCompleted: boolean
    newSignUp: boolean
    weeklyDigest: boolean
    billingAlerts: boolean
  }
}

const notifications = ref<NotificationPrefs>({
  email: { campaignCompleted: true, newSignUp: true, weeklyDigest: true, billingAlerts: true },
  inApp: { campaignCompleted: true, newSignUp: true, weeklyDigest: false, billingAlerts: true },
})

// --- Helpers ---
function showSaveSuccess(msg = 'Changes saved') {
  saveMessage.value = msg
  setTimeout(() => { saveMessage.value = '' }, 3000)
}

async function saveSection(section: string, payload: unknown) {
  saving.value = true
  try {
    await api.put(`/api/v1/settings/${section}`, payload)
    showSaveSuccess()
  } catch {
    saveMessage.value = 'Failed to save. Please try again.'
    setTimeout(() => { saveMessage.value = '' }, 3000)
  } finally {
    saving.value = false
  }
}

function formatUsage(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n.toLocaleString()
}

function usagePercent(used: number, limit: number): number {
  return limit > 0 ? Math.min(Math.round((used / limit) * 100), 100) : 0
}

function usageBarColor(used: number, limit: number): string {
  const pct = usagePercent(used, limit)
  if (pct >= 90) return '#DC2626'
  if (pct >= 70) return '#D97706'
  return '#48BB78'
}

function maskKey(key: string): string {
  if (key.length <= 8) return '****' + key.slice(-4)
  return '****' + key.slice(-8)
}

async function inviteMember() {
  if (!inviteEmail.value) return
  saving.value = true
  try {
    await api.post('/api/v1/settings/organization/invite', { email: inviteEmail.value, role: inviteRole.value })
    org.value.teamMembers.push({
      id: Date.now().toString(),
      name: inviteEmail.value.split('@')[0],
      email: inviteEmail.value,
      role: inviteRole.value,
      status: 'invited',
    })
    inviteEmail.value = ''
    inviteRole.value = 'member'
    showInviteForm.value = false
    showSaveSuccess('Invitation sent')
  } catch {
    saveMessage.value = 'Failed to send invitation.'
    setTimeout(() => { saveMessage.value = '' }, 3000)
  } finally {
    saving.value = false
  }
}

async function removeMember(id: string) {
  try {
    await api.delete(`/api/v1/settings/organization/members/${id}`)
    org.value.teamMembers = org.value.teamMembers.filter(m => m.id !== id)
    showSaveSuccess('Member removed')
  } catch {
    saveMessage.value = 'Failed to remove member.'
    setTimeout(() => { saveMessage.value = '' }, 3000)
  }
}

async function createApiKey() {
  if (!newKeyName.value) return
  saving.value = true
  try {
    const res = await api.post('/api/v1/settings/api-keys', {
      name: newKeyName.value,
      permissions: newKeyPermissions.value,
    })
    apiKeys.value.push({ ...res.data, visible: false })
    newKeyName.value = ''
    newKeyPermissions.value = []
    showCreateKey.value = false
    showSaveSuccess('API key created')
  } catch {
    saveMessage.value = 'Failed to create API key.'
    setTimeout(() => { saveMessage.value = '' }, 3000)
  } finally {
    saving.value = false
  }
}

async function revokeApiKey(id: string) {
  try {
    await api.delete(`/api/v1/settings/api-keys/${id}`)
    apiKeys.value = apiKeys.value.filter(k => k.id !== id)
    showSaveSuccess('API key revoked')
  } catch {
    saveMessage.value = 'Failed to revoke key.'
    setTimeout(() => { saveMessage.value = '' }, 3000)
  }
}

onMounted(async () => {
  try {
    const res = await api.get('/api/v1/settings')
    if (res.data.profile) profile.value = res.data.profile
    if (res.data.organization) org.value = res.data.organization
    if (res.data.billing) billing.value = res.data.billing
    if (res.data.apiKeys) apiKeys.value = res.data.apiKeys.map((k: ApiKey) => ({ ...k, visible: false }))
    if (res.data.notifications) notifications.value = res.data.notifications
  } catch {
    // Keep defaults
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-text-primary mb-1">Settings</h1>
      <p class="text-sm text-text-secondary">Manage your account, organization, and preferences.</p>
    </div>

    <!-- Save Message Toast -->
    <Transition name="fade">
      <div v-if="saveMessage" class="fixed top-6 right-6 z-50 px-4 py-2.5 rounded-lg text-sm font-medium shadow-lg"
        :class="saveMessage.includes('Failed') ? 'bg-status-paused-bg text-status-paused' : 'bg-status-active-bg text-status-active'">
        {{ saveMessage }}
      </div>
    </Transition>

    <!-- Tabs -->
    <div class="border-b border-border">
      <div class="flex gap-0">
        <button
          v-for="tab in tabs"
          :key="tab"
          class="px-5 py-3 text-sm font-medium border-b-2 transition-colors -mb-px"
          :class="activeTab === tab
            ? 'border-quest text-text-primary'
            : 'border-transparent text-text-muted hover:text-text-secondary'"
          @click="activeTab = tab"
        >{{ tab }}</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <span class="material-symbols-rounded text-4xl text-text-muted animate-spin">progress_activity</span>
    </div>

    <template v-else>
      <!-- ===== PROFILE TAB ===== -->
      <div v-if="activeTab === 'Profile'" class="max-w-2xl space-y-6">
        <!-- Avatar -->
        <div class="flex items-center gap-5">
          <div class="w-20 h-20 rounded-full bg-card-bg border border-border flex items-center justify-center overflow-hidden shrink-0">
            <img v-if="profile.avatar" :src="profile.avatar" alt="Avatar" class="w-full h-full object-cover" />
            <span v-else class="text-2xl font-bold text-text-muted">{{ profile.initials }}</span>
          </div>
          <button class="px-4 py-2 text-sm font-medium bg-card-bg border border-border rounded-lg text-text-primary hover:bg-white/5 transition-colors">
            Change Avatar
          </button>
        </div>

        <!-- Email -->
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-1.5">Email</label>
          <div class="flex items-center gap-2">
            <input
              type="email"
              :value="profile.email"
              readonly
              class="flex-1 px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-muted cursor-not-allowed"
            />
            <span v-if="profile.emailVerified" class="flex items-center gap-1 text-xs text-status-active">
              <span class="material-symbols-rounded text-sm">verified</span> Verified
            </span>
          </div>
        </div>

        <!-- Display Name -->
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-1.5">Display Name</label>
          <input
            v-model="profile.displayName"
            type="text"
            placeholder="Enter your display name"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-quest transition-colors"
          />
        </div>

        <!-- Timezone -->
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-1.5">Timezone</label>
          <select
            v-model="profile.timezone"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary focus:outline-none focus:border-quest transition-colors appearance-none"
          >
            <option v-for="tz in timezones" :key="tz" :value="tz">{{ tz }}</option>
          </select>
        </div>

        <!-- Save -->
        <button
          class="px-6 py-2.5 bg-quest text-white text-sm font-semibold rounded-lg hover:bg-quest/90 transition-colors disabled:opacity-50"
          :disabled="saving"
          @click="saveSection('profile', profile)"
        >{{ saving ? 'Saving...' : 'Save Changes' }}</button>
      </div>

      <!-- ===== ORGANIZATION TAB ===== -->
      <div v-if="activeTab === 'Organization'" class="space-y-8">
        <div class="max-w-2xl space-y-6">
          <!-- Org Name -->
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1.5">Organization Name</label>
            <input
              v-model="org.name"
              type="text"
              placeholder="Your organization name"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-quest transition-colors"
            />
          </div>

          <!-- Org Logo -->
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1.5">Organization Logo</label>
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 rounded-xl bg-page-bg border border-border flex items-center justify-center overflow-hidden">
                <img v-if="org.logo" :src="org.logo" alt="Logo" class="w-full h-full object-cover" />
                <span v-else class="material-symbols-rounded text-2xl text-text-muted">business</span>
              </div>
              <button class="px-4 py-2 text-sm font-medium bg-card-bg border border-border rounded-lg text-text-primary hover:bg-white/5 transition-colors">
                Upload Logo
              </button>
            </div>
          </div>

          <!-- Website URL -->
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1.5">Website URL</label>
            <input
              v-model="org.website"
              type="url"
              placeholder="https://your-project.com"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-quest transition-colors"
            />
          </div>

          <!-- Industry -->
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1.5">Industry</label>
            <select
              v-model="org.industry"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary focus:outline-none focus:border-quest transition-colors appearance-none"
            >
              <option value="" disabled>Select industry</option>
              <option v-for="ind in industries" :key="ind" :value="ind">{{ ind }}</option>
            </select>
          </div>

          <!-- Save Org -->
          <button
            class="px-6 py-2.5 bg-quest text-white text-sm font-semibold rounded-lg hover:bg-quest/90 transition-colors disabled:opacity-50"
            :disabled="saving"
            @click="saveSection('organization', { name: org.name, logo: org.logo, website: org.website, industry: org.industry })"
          >{{ saving ? 'Saving...' : 'Save Changes' }}</button>
        </div>

        <!-- Team Members -->
        <div>
          <div class="flex items-center justify-between mb-4">
            <div class="text-xs font-semibold text-text-muted uppercase tracking-wider">Team Members</div>
            <button
              class="flex items-center gap-1 text-sm font-medium text-quest hover:text-quest/80 transition-colors"
              @click="showInviteForm = !showInviteForm"
            >
              <span class="material-symbols-rounded text-lg">add</span> Invite Member
            </button>
          </div>

          <!-- Invite Form -->
          <div v-if="showInviteForm" class="bg-card-bg border border-border rounded-xl p-5 mb-4">
            <div class="flex items-end gap-4">
              <div class="flex-1">
                <label class="block text-xs text-text-muted mb-1">Email</label>
                <input
                  v-model="inviteEmail"
                  type="email"
                  placeholder="team@example.com"
                  class="w-full px-3 py-2 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-quest transition-colors"
                />
              </div>
              <div class="w-36">
                <label class="block text-xs text-text-muted mb-1">Role</label>
                <select
                  v-model="inviteRole"
                  class="w-full px-3 py-2 bg-page-bg border border-border rounded-lg text-sm text-text-primary focus:outline-none focus:border-quest transition-colors appearance-none"
                >
                  <option value="admin">Admin</option>
                  <option value="member">Member</option>
                </select>
              </div>
              <button
                class="px-4 py-2 bg-quest text-white text-sm font-semibold rounded-lg hover:bg-quest/90 transition-colors disabled:opacity-50"
                :disabled="saving || !inviteEmail"
                @click="inviteMember"
              >Send Invite</button>
              <button
                class="px-4 py-2 text-sm text-text-muted hover:text-text-secondary transition-colors"
                @click="showInviteForm = false"
              >Cancel</button>
            </div>
          </div>

          <!-- Members Table -->
          <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
            <table class="w-full">
              <thead>
                <tr class="border-b border-border">
                  <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Name</th>
                  <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Email</th>
                  <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Role</th>
                  <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Status</th>
                  <th class="text-right px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border">
                <tr v-for="member in org.teamMembers" :key="member.id" class="hover:bg-white/2 transition-colors">
                  <td class="px-5 py-3.5 text-sm font-medium text-text-primary">{{ member.name }}</td>
                  <td class="px-5 py-3.5 text-sm text-text-secondary">{{ member.email }}</td>
                  <td class="px-5 py-3.5">
                    <select
                      v-if="member.role !== 'owner'"
                      v-model="member.role"
                      class="px-2 py-1 bg-page-bg border border-border rounded text-xs text-text-primary focus:outline-none focus:border-quest transition-colors appearance-none"
                    >
                      <option value="admin">Admin</option>
                      <option value="member">Member</option>
                    </select>
                    <span v-else class="text-xs font-medium text-quest">Owner</span>
                  </td>
                  <td class="px-5 py-3.5">
                    <span
                      class="text-xs font-medium px-2 py-0.5 rounded-full capitalize"
                      :class="member.status === 'active' ? 'bg-status-active-bg text-status-active' : 'bg-status-draft-bg text-status-draft'"
                    >{{ member.status }}</span>
                  </td>
                  <td class="px-5 py-3.5 text-right">
                    <button
                      v-if="member.role !== 'owner'"
                      class="text-xs text-status-paused hover:underline"
                      @click="removeMember(member.id)"
                    >Remove</button>
                  </td>
                </tr>
                <tr v-if="!org.teamMembers.length">
                  <td colspan="5" class="px-5 py-8 text-center text-sm text-text-muted">No team members yet</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- ===== BILLING TAB ===== -->
      <div v-if="activeTab === 'Billing'" class="space-y-8">
        <!-- Current Plan -->
        <div class="bg-card-bg border border-border rounded-xl p-6">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Current Plan</div>
              <div class="text-xl font-bold text-text-primary">{{ billing.planName || 'Free' }}</div>
              <div class="text-sm text-text-secondary mt-1">{{ billing.planPrice || '$0' }} / month</div>
              <div v-if="billing.renewalDate" class="text-xs text-text-muted mt-1">Renews on {{ billing.renewalDate }}</div>
            </div>
            <div class="flex gap-3">
              <button class="px-5 py-2.5 bg-quest text-white text-sm font-semibold rounded-lg hover:bg-quest/90 transition-colors">Upgrade</button>
              <button class="px-5 py-2.5 bg-card-bg border border-border text-text-primary text-sm font-medium rounded-lg hover:bg-white/5 transition-colors">Manage</button>
            </div>
          </div>
        </div>

        <!-- Payment Method -->
        <div class="bg-card-bg border border-border rounded-xl p-5">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="material-symbols-rounded text-xl text-text-secondary">credit_card</span>
              <span class="text-sm text-text-primary">Card ending in {{ billing.cardLast4 || '****' }}</span>
            </div>
            <button class="text-sm text-quest hover:underline">Update</button>
          </div>
        </div>

        <!-- Billing History -->
        <div>
          <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Billing History</div>
          <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
            <table class="w-full">
              <thead>
                <tr class="border-b border-border">
                  <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Date</th>
                  <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Description</th>
                  <th class="text-right px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Amount</th>
                  <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Status</th>
                  <th class="text-right px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Invoice</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border">
                <tr v-for="item in billing.history" :key="item.id" class="hover:bg-white/2 transition-colors">
                  <td class="px-5 py-3.5 text-sm text-text-muted">{{ item.date }}</td>
                  <td class="px-5 py-3.5 text-sm text-text-primary">{{ item.description }}</td>
                  <td class="px-5 py-3.5 text-sm text-text-primary text-right">{{ item.amount }}</td>
                  <td class="px-5 py-3.5">
                    <span
                      class="text-xs font-medium px-2 py-0.5 rounded-full capitalize"
                      :class="item.status === 'paid' ? 'bg-status-active-bg text-status-active' : 'bg-status-draft-bg text-status-draft'"
                    >{{ item.status }}</span>
                  </td>
                  <td class="px-5 py-3.5 text-right">
                    <button class="text-xs text-quest hover:underline">Download</button>
                  </td>
                </tr>
                <tr v-if="!billing.history.length">
                  <td colspan="5" class="px-5 py-8 text-center text-sm text-text-muted">No billing history</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Usage Meters -->
        <div>
          <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Usage</div>
          <div class="grid grid-cols-3 gap-6">
            <div class="bg-card-bg border border-border rounded-xl p-5">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm text-text-secondary">API Calls</span>
                <span class="text-xs text-text-muted">{{ formatUsage(billing.usage.apiCalls.used) }} / {{ formatUsage(billing.usage.apiCalls.limit) }}</span>
              </div>
              <div class="w-full bg-page-bg rounded-full h-2">
                <div class="h-2 rounded-full transition-all" :style="{ width: usagePercent(billing.usage.apiCalls.used, billing.usage.apiCalls.limit) + '%', background: usageBarColor(billing.usage.apiCalls.used, billing.usage.apiCalls.limit) }"></div>
              </div>
            </div>
            <div class="bg-card-bg border border-border rounded-xl p-5">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm text-text-secondary">Storage</span>
                <span class="text-xs text-text-muted">{{ formatUsage(billing.usage.storage.used) }} / {{ formatUsage(billing.usage.storage.limit) }}</span>
              </div>
              <div class="w-full bg-page-bg rounded-full h-2">
                <div class="h-2 rounded-full transition-all" :style="{ width: usagePercent(billing.usage.storage.used, billing.usage.storage.limit) + '%', background: usageBarColor(billing.usage.storage.used, billing.usage.storage.limit) }"></div>
              </div>
            </div>
            <div class="bg-card-bg border border-border rounded-xl p-5">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm text-text-secondary">Team Members</span>
                <span class="text-xs text-text-muted">{{ billing.usage.teamMembers.used }} of {{ billing.usage.teamMembers.limit }}</span>
              </div>
              <div class="w-full bg-page-bg rounded-full h-2">
                <div class="h-2 rounded-full transition-all" :style="{ width: usagePercent(billing.usage.teamMembers.used, billing.usage.teamMembers.limit) + '%', background: usageBarColor(billing.usage.teamMembers.used, billing.usage.teamMembers.limit) }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== API KEYS TAB ===== -->
      <div v-if="activeTab === 'API Keys'" class="space-y-6">
        <!-- Warning -->
        <div class="flex items-start gap-3 bg-status-draft-bg border border-status-draft/20 rounded-xl p-4">
          <span class="material-symbols-rounded text-status-draft text-xl shrink-0 mt-0.5">warning</span>
          <p class="text-sm text-text-secondary">API keys provide full access to your account. Keep them secure and never share them publicly.</p>
        </div>

        <!-- Create Key Button -->
        <div class="flex justify-end">
          <button
            class="flex items-center gap-1 text-sm font-medium text-quest hover:text-quest/80 transition-colors"
            @click="showCreateKey = !showCreateKey"
          >
            <span class="material-symbols-rounded text-lg">add</span> Create New Key
          </button>
        </div>

        <!-- Create Key Form -->
        <div v-if="showCreateKey" class="bg-card-bg border border-border rounded-xl p-5">
          <div class="space-y-4">
            <div>
              <label class="block text-xs text-text-muted mb-1">Key Name</label>
              <input
                v-model="newKeyName"
                type="text"
                placeholder="e.g., Production API Key"
                class="w-full px-3 py-2 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-quest transition-colors"
              />
            </div>
            <div>
              <label class="block text-xs text-text-muted mb-2">Permissions</label>
              <div class="flex flex-wrap gap-3">
                <label v-for="perm in permissionOptions" :key="perm" class="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    :value="perm"
                    v-model="newKeyPermissions"
                    class="w-4 h-4 rounded border-border bg-page-bg text-quest focus:ring-quest"
                  />
                  <span class="text-sm text-text-secondary">{{ perm }}</span>
                </label>
              </div>
            </div>
            <div class="flex gap-3">
              <button
                class="px-4 py-2 bg-quest text-white text-sm font-semibold rounded-lg hover:bg-quest/90 transition-colors disabled:opacity-50"
                :disabled="saving || !newKeyName"
                @click="createApiKey"
              >Create Key</button>
              <button
                class="px-4 py-2 text-sm text-text-muted hover:text-text-secondary transition-colors"
                @click="showCreateKey = false"
              >Cancel</button>
            </div>
          </div>
        </div>

        <!-- Keys List -->
        <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border">
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Name</th>
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Key</th>
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Created</th>
                <th class="text-left px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Last Used</th>
                <th class="text-right px-5 py-3 text-xs font-semibold text-text-muted uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border">
              <tr v-for="key in apiKeys" :key="key.id" class="hover:bg-white/2 transition-colors">
                <td class="px-5 py-3.5 text-sm font-medium text-text-primary">{{ key.name }}</td>
                <td class="px-5 py-3.5">
                  <div class="flex items-center gap-2">
                    <code class="text-xs bg-page-bg px-2 py-1 rounded text-text-secondary font-mono">{{ key.visible ? key.key : maskKey(key.key) }}</code>
                    <button class="text-xs text-quest hover:underline" @click="key.visible = !key.visible">{{ key.visible ? 'Hide' : 'Show' }}</button>
                  </div>
                </td>
                <td class="px-5 py-3.5 text-sm text-text-muted">{{ key.created }}</td>
                <td class="px-5 py-3.5 text-sm text-text-muted">{{ key.lastUsed || 'Never' }}</td>
                <td class="px-5 py-3.5 text-right">
                  <button class="text-xs text-status-paused hover:underline" @click="revokeApiKey(key.id)">Revoke</button>
                </td>
              </tr>
              <tr v-if="!apiKeys.length">
                <td colspan="5" class="px-5 py-8 text-center text-sm text-text-muted">No API keys created yet</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ===== NOTIFICATIONS TAB ===== -->
      <div v-if="activeTab === 'Notifications'" class="max-w-2xl space-y-8">
        <!-- Email Notifications -->
        <div>
          <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Email Notifications</div>
          <div class="bg-card-bg border border-border rounded-xl divide-y divide-border">
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">Campaign Completed</div>
                <div class="text-xs text-text-muted">Get notified when a campaign finishes</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.email.campaignCompleted" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">New Sign-up</div>
                <div class="text-xs text-text-muted">Get notified when a new user joins</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.email.newSignUp" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">Weekly Digest</div>
                <div class="text-xs text-text-muted">Receive a weekly summary of your activity</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.email.weeklyDigest" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">Billing Alerts</div>
                <div class="text-xs text-text-muted">Important billing and payment updates</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.email.billingAlerts" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
          </div>
        </div>

        <!-- In-App Notifications -->
        <div>
          <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">In-App Notifications</div>
          <div class="bg-card-bg border border-border rounded-xl divide-y divide-border">
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">Campaign Completed</div>
                <div class="text-xs text-text-muted">Show in-app notification when a campaign finishes</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.inApp.campaignCompleted" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">New Sign-up</div>
                <div class="text-xs text-text-muted">Show in-app notification for new sign-ups</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.inApp.newSignUp" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">Weekly Digest</div>
                <div class="text-xs text-text-muted">Show in-app weekly summary notification</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.inApp.weeklyDigest" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
            <div class="px-5 py-4 flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-text-primary">Billing Alerts</div>
                <div class="text-xs text-text-muted">Show in-app billing and payment alerts</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notifications.inApp.billingAlerts" class="sr-only peer" />
                <div class="w-10 h-5 bg-border rounded-full peer peer-checked:bg-quest transition-colors after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-5"></div>
              </label>
            </div>
          </div>
        </div>

        <!-- Save Prefs -->
        <button
          class="px-6 py-2.5 bg-quest text-white text-sm font-semibold rounded-lg hover:bg-quest/90 transition-colors disabled:opacity-50"
          :disabled="saving"
          @click="saveSection('notifications', notifications)"
        >{{ saving ? 'Saving...' : 'Save Preferences' }}</button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
