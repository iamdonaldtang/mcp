<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'
import Modal from '../../components/common/Modal.vue'

// --- Types ---
interface PrivilegeTier {
  id: string
  name: string
  icon: string
  color: string
  rank_order: number
  status: 'active' | 'draft'
  holders_count: number
  qualification_mode: 'status' | 'redemption' | 'achievement'
  qualification_config: Record<string, unknown>
  privileges: PrivilegeConfig
  duration_type: 'permanent' | 'time_limited'
  duration_days?: number
  created_at: string
}

interface PrivilegeConfig {
  early_access: boolean
  exclusive_shop: boolean
  point_multiplier: boolean
  point_multiplier_value: number
  custom_badge: boolean
  custom_badge_id: string
  fee_discount: boolean
  fee_discount_pct: number
  gas_rebate: boolean
  gas_rebate_pct: number
  yield_boost: boolean
  yield_boost_pct: number
  priority_support: boolean
  custom: boolean
  custom_label: string
  custom_value: string
}

interface TierMember {
  id: string
  wallet_address: string
  joined_at: string
  status: 'active' | 'expired'
}

// --- State ---
const loading = ref(true)
const tiers = ref<PrivilegeTier[]>([])

const stats = ref({
  total_tiers: 0,
  active_holders: 0,
  total_value_distributed: 0,
  last_granted: '',
})

// Editor modal
const showEditor = ref(false)
const editingTier = ref<PrivilegeTier | null>(null)
const saving = ref(false)
const tierForm = ref(getEmptyForm())

function getEmptyForm() {
  return {
    name: '',
    icon: 'star',
    color: '#9B7EE0',
    rank_order: 1,
    qualification_mode: 'status' as 'status' | 'redemption' | 'achievement',
    // Status-based
    token_gate_enabled: false,
    token_gate_chain: '',
    token_gate_contract: '',
    token_gate_min_balance: 0,
    level_requirement: 0,
    points_threshold: 0,
    manual_assignment: false,
    // Redemption-based
    shop_item_id: '',
    // Achievement-based
    badge_selector: '',
    milestone_selector: '',
    // Privileges
    early_access: false,
    exclusive_shop: false,
    point_multiplier: false,
    point_multiplier_value: 1.5,
    custom_badge: false,
    custom_badge_id: '',
    fee_discount: false,
    fee_discount_pct: 10,
    gas_rebate: false,
    gas_rebate_pct: 5,
    yield_boost: false,
    yield_boost_pct: 10,
    priority_support: false,
    custom: false,
    custom_label: '',
    custom_value: '',
    // Duration
    duration_type: 'permanent' as 'permanent' | 'time_limited',
    duration_days: 30,
  }
}

const iconOptions = ['star', 'diamond', 'shield', 'workspace_premium', 'military_tech', 'emoji_events', 'local_fire_department', 'rocket_launch']
const colorOptions = ['#9B7EE0', '#5D7EF1', '#48BB78', '#F59E0B', '#ED8936', '#DC2626', '#EC4899', '#06B6D4']

const canSaveTier = computed(() => tierForm.value.name.trim().length > 0 && tierForm.value.name.length <= 30)

// Members panel
const showMembers = ref(false)
const membersTier = ref<PrivilegeTier | null>(null)
const members = ref<TierMember[]>([])
const membersTotal = ref(0)
const membersPage = ref(1)
const membersSearch = ref('')
const membersLoading = ref(false)
const newMemberAddress = ref('')
const addingMember = ref(false)

const sortedTiers = computed(() => [...tiers.value].sort((a, b) => a.rank_order - b.rank_order))

// --- API ---
onMounted(async () => {
  await Promise.all([fetchStats(), fetchTiers()])
  loading.value = false
})

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/whitelabel/privileges/stats')
    if (res.data.data) stats.value = res.data.data
  } catch { /* defaults */ }
}

async function fetchTiers() {
  try {
    const res = await api.get('/api/v1/whitelabel/privileges')
    tiers.value = res.data.data?.items || []
  } catch { /* empty */ }
}

function openCreate() {
  editingTier.value = null
  tierForm.value = getEmptyForm()
  tierForm.value.rank_order = tiers.value.length + 1
  showEditor.value = true
}

function openEdit(tier: PrivilegeTier) {
  editingTier.value = tier
  tierForm.value = {
    name: tier.name,
    icon: tier.icon,
    color: tier.color,
    rank_order: tier.rank_order,
    qualification_mode: tier.qualification_mode,
    token_gate_enabled: !!(tier.qualification_config.token_gate_chain),
    token_gate_chain: (tier.qualification_config.token_gate_chain as string) || '',
    token_gate_contract: (tier.qualification_config.token_gate_contract as string) || '',
    token_gate_min_balance: (tier.qualification_config.token_gate_min_balance as number) || 0,
    level_requirement: (tier.qualification_config.level_requirement as number) || 0,
    points_threshold: (tier.qualification_config.points_threshold as number) || 0,
    manual_assignment: (tier.qualification_config.manual_assignment as boolean) || false,
    shop_item_id: (tier.qualification_config.shop_item_id as string) || '',
    badge_selector: (tier.qualification_config.badge_selector as string) || '',
    milestone_selector: (tier.qualification_config.milestone_selector as string) || '',
    early_access: tier.privileges.early_access,
    exclusive_shop: tier.privileges.exclusive_shop,
    point_multiplier: tier.privileges.point_multiplier,
    point_multiplier_value: tier.privileges.point_multiplier_value || 1.5,
    custom_badge: tier.privileges.custom_badge,
    custom_badge_id: tier.privileges.custom_badge_id || '',
    fee_discount: tier.privileges.fee_discount,
    fee_discount_pct: tier.privileges.fee_discount_pct || 10,
    gas_rebate: tier.privileges.gas_rebate,
    gas_rebate_pct: tier.privileges.gas_rebate_pct || 5,
    yield_boost: tier.privileges.yield_boost,
    yield_boost_pct: tier.privileges.yield_boost_pct || 10,
    priority_support: tier.privileges.priority_support,
    custom: tier.privileges.custom,
    custom_label: tier.privileges.custom_label || '',
    custom_value: tier.privileges.custom_value || '',
    duration_type: tier.duration_type,
    duration_days: tier.duration_days || 30,
  }
  showEditor.value = true
}

async function saveTier() {
  if (!canSaveTier.value) return
  saving.value = true
  const payload = {
    name: tierForm.value.name,
    icon: tierForm.value.icon,
    color: tierForm.value.color,
    rank_order: tierForm.value.rank_order,
    qualification_mode: tierForm.value.qualification_mode,
    qualification_config: buildQualConfig(),
    privileges: {
      early_access: tierForm.value.early_access,
      exclusive_shop: tierForm.value.exclusive_shop,
      point_multiplier: tierForm.value.point_multiplier,
      point_multiplier_value: tierForm.value.point_multiplier_value,
      custom_badge: tierForm.value.custom_badge,
      custom_badge_id: tierForm.value.custom_badge_id,
      fee_discount: tierForm.value.fee_discount,
      fee_discount_pct: tierForm.value.fee_discount_pct,
      gas_rebate: tierForm.value.gas_rebate,
      gas_rebate_pct: tierForm.value.gas_rebate_pct,
      yield_boost: tierForm.value.yield_boost,
      yield_boost_pct: tierForm.value.yield_boost_pct,
      priority_support: tierForm.value.priority_support,
      custom: tierForm.value.custom,
      custom_label: tierForm.value.custom_label,
      custom_value: tierForm.value.custom_value,
    },
    duration_type: tierForm.value.duration_type,
    duration_days: tierForm.value.duration_type === 'time_limited' ? tierForm.value.duration_days : undefined,
  }
  try {
    if (editingTier.value) {
      await api.put(`/api/v1/whitelabel/privileges/${editingTier.value.id}`, payload)
    } else {
      await api.post('/api/v1/whitelabel/privileges', payload)
    }
    showEditor.value = false
    await Promise.all([fetchStats(), fetchTiers()])
  } catch { /* TODO: toast */ }
  finally { saving.value = false }
}

function buildQualConfig(): Record<string, unknown> {
  switch (tierForm.value.qualification_mode) {
    case 'status': return {
      token_gate_chain: tierForm.value.token_gate_enabled ? tierForm.value.token_gate_chain : undefined,
      token_gate_contract: tierForm.value.token_gate_enabled ? tierForm.value.token_gate_contract : undefined,
      token_gate_min_balance: tierForm.value.token_gate_enabled ? tierForm.value.token_gate_min_balance : undefined,
      level_requirement: tierForm.value.level_requirement || undefined,
      points_threshold: tierForm.value.points_threshold || undefined,
      manual_assignment: tierForm.value.manual_assignment,
    }
    case 'redemption': return { shop_item_id: tierForm.value.shop_item_id }
    case 'achievement': return { badge_selector: tierForm.value.badge_selector, milestone_selector: tierForm.value.milestone_selector }
    default: return {}
  }
}

async function deleteTier(tier: PrivilegeTier) {
  if (!confirm(`Delete tier "${tier.name}"? Members will lose their privileges.`)) return
  try {
    await api.delete(`/api/v1/whitelabel/privileges/${tier.id}`)
    await Promise.all([fetchStats(), fetchTiers()])
  } catch { /* TODO: toast */ }
}

// --- Members panel ---
async function openMembers(tier: PrivilegeTier) {
  membersTier.value = tier
  membersPage.value = 1
  membersSearch.value = ''
  showMembers.value = true
  await fetchMembers()
}

async function fetchMembers() {
  if (!membersTier.value) return
  membersLoading.value = true
  try {
    const params: Record<string, string | number> = { page: membersPage.value, page_size: 20 }
    if (membersSearch.value) params.search = membersSearch.value
    const res = await api.get(`/api/v1/whitelabel/privileges/${membersTier.value.id}/members`, { params })
    members.value = res.data.data?.items || []
    membersTotal.value = res.data.data?.total || 0
  } catch { /* empty */ }
  finally { membersLoading.value = false }
}

async function addMember() {
  if (!newMemberAddress.value.trim() || !membersTier.value) return
  addingMember.value = true
  try {
    await api.post(`/api/v1/whitelabel/privileges/${membersTier.value.id}/members`, {
      wallet_address: newMemberAddress.value,
    })
    newMemberAddress.value = ''
    await fetchMembers()
  } catch { /* TODO: toast */ }
  finally { addingMember.value = false }
}

async function removeMember(memberId: string) {
  if (!membersTier.value) return
  try {
    await api.delete(`/api/v1/whitelabel/privileges/${membersTier.value.id}/members/${memberId}`)
    await fetchMembers()
  } catch { /* TODO: toast */ }
}

function handleCsvUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file || !membersTier.value) return
  const formData = new FormData()
  formData.append('file', file)
  api.post(`/api/v1/whitelabel/privileges/${membersTier.value.id}/members/import`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }).then(() => fetchMembers()).catch(() => { /* TODO: toast */ })
}

async function exportMembers() {
  if (!membersTier.value) return
  try {
    const res = await api.get(`/api/v1/whitelabel/privileges/${membersTier.value.id}/members/export`, {
      responseType: 'blob',
    })
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const a = document.createElement('a')
    a.href = url
    a.download = `${membersTier.value.name}_members.csv`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch { /* TODO: toast */ }
}

function truncateAddress(addr: string): string {
  if (!addr || addr.length < 12) return addr
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`
}

function formatDate(d: string): string {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function getPrivilegeTags(tier: PrivilegeTier): string[] {
  const tags: string[] = []
  const p = tier.privileges
  if (p.early_access) tags.push('Early Access')
  if (p.exclusive_shop) tags.push('Exclusive Shop')
  if (p.point_multiplier) tags.push(`${p.point_multiplier_value}x Points`)
  if (p.custom_badge) tags.push('Custom Badge')
  if (p.fee_discount) tags.push(`${p.fee_discount_pct}% Fee Discount`)
  if (p.gas_rebate) tags.push(`${p.gas_rebate_pct}% Gas Rebate`)
  if (p.yield_boost) tags.push(`${p.yield_boost_pct}% Yield Boost`)
  if (p.priority_support) tags.push('Priority Support')
  if (p.custom) tags.push(p.custom_label || 'Custom')
  return tags
}

function getQualSummary(tier: PrivilegeTier): string {
  switch (tier.qualification_mode) {
    case 'status': {
      const parts: string[] = []
      if (tier.qualification_config.token_gate_chain) parts.push('Token Gate')
      if (tier.qualification_config.level_requirement) parts.push(`Level ${tier.qualification_config.level_requirement}+`)
      if (tier.qualification_config.points_threshold) parts.push(`${tier.qualification_config.points_threshold}+ pts`)
      if (tier.qualification_config.manual_assignment) parts.push('Manual')
      return parts.length ? parts.join(' + ') : 'Status-based'
    }
    case 'redemption': return 'Redemption-based'
    case 'achievement': return 'Achievement-based'
    default: return '—'
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Privilege Manager</h1>
        <p class="text-sm text-text-secondary">Define tiered privileges and perks for your community members</p>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard label="Total Tiers" :value="stats.total_tiers" icon="layers" icon-color="#9B7EE0" />
      <StatsCard label="Active Holders" :value="stats.active_holders" icon="group" icon-color="#48BB78" />
      <StatsCard label="Total Value Distributed" :value="stats.total_value_distributed" icon="payments" icon-color="#5D7EF1" />
      <StatsCard label="Last Granted" :value="stats.last_granted ? formatDate(stats.last_granted) : '—'" icon="schedule" icon-color="#F59E0B" />
    </div>

    <!-- Tier Ladder -->
    <div class="relative">
      <!-- Loading -->
      <div v-if="loading" class="space-y-4">
        <div v-for="i in 3" :key="i" class="bg-card-bg border border-border rounded-xl p-6 h-32 animate-pulse" />
      </div>

      <!-- Empty -->
      <div v-else-if="sortedTiers.length === 0" class="bg-card-bg border border-border rounded-xl p-12 text-center">
        <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">layers</span>
        <p class="text-sm text-text-muted mb-3">No privilege tiers yet</p>
        <button class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors" @click="openCreate">
          + Add First Tier
        </button>
      </div>

      <!-- Tier cards with connecting line -->
      <div v-else class="relative pl-8">
        <!-- Vertical connector line -->
        <div
          v-if="sortedTiers.length > 1"
          class="absolute left-[15px] top-6 bottom-6 w-0.5 bg-[#1E293B]"
        />

        <div class="space-y-4">
          <div
            v-for="(tier, idx) in sortedTiers"
            :key="tier.id"
            class="relative bg-card-bg border rounded-xl p-5 transition-all"
            :class="tier.status === 'draft' ? 'opacity-50 border-dashed border-[#1E293B]' : 'border-[#1E293B]'"
            :style="tier.status === 'active' ? { borderLeftWidth: '3px', borderLeftColor: tier.color } : {}"
          >
            <!-- Rank dot on connector line -->
            <div
              class="absolute -left-8 top-1/2 -translate-y-1/2 w-[14px] h-[14px] rounded-full border-2 border-[#111B27]"
              :style="{ backgroundColor: tier.color }"
            />

            <div class="flex items-start gap-4">
              <!-- Icon circle -->
              <div
                class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0"
                :style="{ backgroundColor: tier.color + '20' }"
              >
                <span class="material-symbols-rounded text-xl" :style="{ color: tier.color }">{{ tier.icon }}</span>
              </div>

              <!-- Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-3 mb-1">
                  <span class="text-xs font-bold text-text-muted uppercase">Rank #{{ tier.rank_order }}</span>
                  <span class="text-base font-semibold text-text-primary">{{ tier.name }}</span>
                  <span
                    class="px-2 py-0.5 text-xs rounded-full font-medium"
                    :class="tier.status === 'active' ? 'bg-[#0A2E1A] text-[#16A34A]' : 'bg-[#1F1A08] text-[#D97706]'"
                  >
                    {{ tier.status === 'active' ? 'Active' : 'Draft' }}
                  </span>
                  <span class="text-xs text-text-muted ml-auto">{{ tier.holders_count.toLocaleString() }} holders</span>
                </div>

                <!-- Qualification summary -->
                <p class="text-xs text-text-muted mb-3">{{ getQualSummary(tier) }}</p>

                <!-- Privilege tags -->
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="tag in getPrivilegeTags(tier)"
                    :key="tag"
                    class="px-2 py-0.5 text-xs rounded bg-[#0A0F1A] border border-[#1E293B] text-text-secondary"
                  >
                    {{ tag }}
                  </span>
                  <span v-if="getPrivilegeTags(tier).length === 0" class="text-xs text-text-muted">No privileges configured</span>
                </div>
              </div>

              <!-- Actions -->
              <div class="flex items-center gap-1 shrink-0">
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit" @click="openEdit(tier)">
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="View Members" @click="openMembers(tier)">
                  <span class="material-symbols-rounded text-base text-text-muted">group</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteTier(tier)">
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Add Tier button at bottom -->
        <button
          class="relative mt-4 w-full py-3 border border-dashed border-[#1E293B] rounded-xl text-sm text-text-muted hover:text-wl hover:border-wl/30 transition-colors"
          @click="openCreate"
        >
          + Add Tier
        </button>
      </div>
    </div>

    <!-- D14 Privilege Tier Editor Modal -->
    <Modal :open="showEditor" :title="editingTier ? 'Edit Tier' : 'Add Tier'" max-width="640px" @close="showEditor = false">
      <div class="space-y-5 max-h-[65vh] overflow-y-auto pr-1">
        <!-- Basic Info -->
        <div class="grid grid-cols-2 gap-4">
          <div class="col-span-2">
            <label class="block text-sm text-text-secondary mb-1.5">Tier Name</label>
            <input
              v-model="tierForm.name"
              type="text"
              maxlength="30"
              placeholder="e.g. Gold"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
            />
            <span class="text-xs text-text-muted mt-1 block">{{ tierForm.name.length }}/30</span>
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Icon</label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="ic in iconOptions"
                :key="ic"
                class="w-9 h-9 rounded-lg flex items-center justify-center border transition-colors"
                :class="tierForm.icon === ic ? 'border-wl bg-wl/10' : 'border-[#1E293B] hover:border-[#94A3B8]'"
                @click="tierForm.icon = ic"
              >
                <span class="material-symbols-rounded text-base" :style="{ color: tierForm.color }">{{ ic }}</span>
              </button>
            </div>
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Color</label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="c in colorOptions"
                :key="c"
                class="w-9 h-9 rounded-full border-2 transition-all"
                :class="tierForm.color === c ? 'border-white scale-110' : 'border-transparent'"
                :style="{ backgroundColor: c }"
                @click="tierForm.color = c"
              />
            </div>
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Rank Order</label>
            <input
              v-model.number="tierForm.rank_order"
              type="number"
              min="1"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
            />
            <span class="text-xs text-text-muted mt-1 block">Lower = higher rank</span>
          </div>
        </div>

        <!-- Qualification -->
        <div class="border border-[#1E293B] rounded-lg p-4 space-y-4">
          <span class="text-sm font-semibold text-text-primary uppercase tracking-wider">Qualification</span>

          <div class="flex gap-2">
            <label
              v-for="mode in (['status', 'redemption', 'achievement'] as const)"
              :key="mode"
              class="flex items-center gap-2 px-3 py-2 rounded-lg border cursor-pointer transition-colors"
              :class="tierForm.qualification_mode === mode ? 'border-wl bg-wl/10 text-wl' : 'border-[#1E293B] text-[#94A3B8]'"
            >
              <input v-model="tierForm.qualification_mode" type="radio" :value="mode" class="sr-only" />
              <span class="text-xs capitalize">{{ mode }}-Based</span>
            </label>
          </div>

          <!-- Status-Based -->
          <template v-if="tierForm.qualification_mode === 'status'">
            <div class="space-y-3">
              <!-- Token Gate -->
              <div class="flex items-center gap-2">
                <input id="token_gate" v-model="tierForm.token_gate_enabled" type="checkbox" class="accent-wl" />
                <label for="token_gate" class="text-sm text-text-secondary">Token Gate</label>
              </div>
              <div v-if="tierForm.token_gate_enabled" class="grid grid-cols-3 gap-3 pl-6">
                <select
                  v-model="tierForm.token_gate_chain"
                  class="px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] focus:border-wl focus:outline-none"
                >
                  <option value="">Chain</option>
                  <option value="ethereum">Ethereum</option>
                  <option value="bsc">BSC</option>
                  <option value="polygon">Polygon</option>
                </select>
                <input
                  v-model="tierForm.token_gate_contract"
                  type="text"
                  placeholder="Contract"
                  class="px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
                <input
                  v-model.number="tierForm.token_gate_min_balance"
                  type="number"
                  placeholder="Min balance"
                  class="px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
              <!-- Level -->
              <div>
                <label class="block text-xs text-text-muted mb-1">Level Requirement (0 = none)</label>
                <input
                  v-model.number="tierForm.level_requirement"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] focus:border-wl focus:outline-none"
                />
              </div>
              <!-- Points -->
              <div>
                <label class="block text-xs text-text-muted mb-1">Points Threshold (0 = none)</label>
                <input
                  v-model.number="tierForm.points_threshold"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] focus:border-wl focus:outline-none"
                />
              </div>
              <!-- Manual -->
              <div class="flex items-center gap-2">
                <input id="manual_assign" v-model="tierForm.manual_assignment" type="checkbox" class="accent-wl" />
                <label for="manual_assign" class="text-sm text-text-secondary">Allow Manual Assignment</label>
              </div>
            </div>
          </template>

          <!-- Redemption-Based -->
          <template v-if="tierForm.qualification_mode === 'redemption'">
            <div>
              <label class="block text-xs text-text-muted mb-1">Shop Item</label>
              <input
                v-model="tierForm.shop_item_id"
                type="text"
                placeholder="Shop Item ID from Benefits Shop"
                class="w-full px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
              />
            </div>
          </template>

          <!-- Achievement-Based -->
          <template v-if="tierForm.qualification_mode === 'achievement'">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs text-text-muted mb-1">Badge</label>
                <input
                  v-model="tierForm.badge_selector"
                  type="text"
                  placeholder="Badge ID"
                  class="w-full px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1">Milestone</label>
                <input
                  v-model="tierForm.milestone_selector"
                  type="text"
                  placeholder="Milestone ID"
                  class="w-full px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
            </div>
          </template>
        </div>

        <!-- Privileges -->
        <div class="border border-[#1E293B] rounded-lg p-4 space-y-3">
          <span class="text-sm font-semibold text-text-primary uppercase tracking-wider">Privileges</span>

          <!-- Simple toggles -->
          <div class="flex items-center justify-between">
            <label class="text-sm text-text-secondary">Early Access</label>
            <input v-model="tierForm.early_access" type="checkbox" class="accent-wl" />
          </div>
          <div class="flex items-center justify-between">
            <label class="text-sm text-text-secondary">Exclusive Shop</label>
            <input v-model="tierForm.exclusive_shop" type="checkbox" class="accent-wl" />
          </div>
          <div class="flex items-center justify-between">
            <label class="text-sm text-text-secondary">Priority Support</label>
            <input v-model="tierForm.priority_support" type="checkbox" class="accent-wl" />
          </div>

          <!-- Toggle + value inputs -->
          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <input v-model="tierForm.point_multiplier" type="checkbox" class="accent-wl" />
              <label class="text-sm text-text-secondary">Point Multiplier</label>
            </div>
            <input
              v-if="tierForm.point_multiplier"
              v-model.number="tierForm.point_multiplier_value"
              type="number"
              min="1"
              step="0.1"
              class="w-20 px-2 py-1 bg-[#0A0F1A] border border-[#1E293B] rounded text-xs text-[#F1F5F9] text-right focus:border-wl focus:outline-none"
            />
            <span v-if="tierForm.point_multiplier" class="text-xs text-text-muted">x</span>
          </div>

          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <input v-model="tierForm.custom_badge" type="checkbox" class="accent-wl" />
              <label class="text-sm text-text-secondary">Custom Badge</label>
            </div>
            <input
              v-if="tierForm.custom_badge"
              v-model="tierForm.custom_badge_id"
              type="text"
              placeholder="Badge ID"
              class="w-32 px-2 py-1 bg-[#0A0F1A] border border-[#1E293B] rounded text-xs text-[#F1F5F9] placeholder-[#64748B] text-right focus:border-wl focus:outline-none"
            />
          </div>

          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <input v-model="tierForm.fee_discount" type="checkbox" class="accent-wl" />
              <label class="text-sm text-text-secondary">Fee Discount</label>
            </div>
            <div v-if="tierForm.fee_discount" class="flex items-center gap-1">
              <input
                v-model.number="tierForm.fee_discount_pct"
                type="number"
                min="1"
                max="100"
                class="w-16 px-2 py-1 bg-[#0A0F1A] border border-[#1E293B] rounded text-xs text-[#F1F5F9] text-right focus:border-wl focus:outline-none"
              />
              <span class="text-xs text-text-muted">%</span>
            </div>
          </div>

          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <input v-model="tierForm.gas_rebate" type="checkbox" class="accent-wl" />
              <label class="text-sm text-text-secondary">Gas Rebate</label>
            </div>
            <div v-if="tierForm.gas_rebate" class="flex items-center gap-1">
              <input
                v-model.number="tierForm.gas_rebate_pct"
                type="number"
                min="1"
                max="100"
                class="w-16 px-2 py-1 bg-[#0A0F1A] border border-[#1E293B] rounded text-xs text-[#F1F5F9] text-right focus:border-wl focus:outline-none"
              />
              <span class="text-xs text-text-muted">%</span>
            </div>
          </div>

          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <input v-model="tierForm.yield_boost" type="checkbox" class="accent-wl" />
              <label class="text-sm text-text-secondary">Yield Boost</label>
            </div>
            <div v-if="tierForm.yield_boost" class="flex items-center gap-1">
              <input
                v-model.number="tierForm.yield_boost_pct"
                type="number"
                min="1"
                max="100"
                class="w-16 px-2 py-1 bg-[#0A0F1A] border border-[#1E293B] rounded text-xs text-[#F1F5F9] text-right focus:border-wl focus:outline-none"
              />
              <span class="text-xs text-text-muted">%</span>
            </div>
          </div>

          <!-- Custom privilege -->
          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <input v-model="tierForm.custom" type="checkbox" class="accent-wl" />
              <label class="text-sm text-text-secondary">Custom</label>
            </div>
          </div>
          <div v-if="tierForm.custom" class="grid grid-cols-2 gap-3 pl-6">
            <input
              v-model="tierForm.custom_label"
              type="text"
              placeholder="Label"
              class="px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
            />
            <input
              v-model="tierForm.custom_value"
              type="text"
              placeholder="Value"
              class="px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-xs text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
            />
          </div>
        </div>

        <!-- Duration -->
        <div>
          <label class="block text-sm text-text-secondary mb-2">Duration</label>
          <div class="flex gap-3 items-center">
            <label
              v-for="dt in (['permanent', 'time_limited'] as const)"
              :key="dt"
              class="flex items-center gap-2 px-3 py-2 rounded-lg border cursor-pointer transition-colors"
              :class="tierForm.duration_type === dt ? 'border-wl bg-wl/10 text-wl' : 'border-[#1E293B] text-[#94A3B8]'"
            >
              <input v-model="tierForm.duration_type" type="radio" :value="dt" class="sr-only" />
              <span class="text-sm">{{ dt === 'permanent' ? 'Permanent' : 'Time-Limited' }}</span>
            </label>
            <input
              v-if="tierForm.duration_type === 'time_limited'"
              v-model.number="tierForm.duration_days"
              type="number"
              min="1"
              class="w-24 px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
            />
            <span v-if="tierForm.duration_type === 'time_limited'" class="text-sm text-text-muted">days</span>
          </div>
        </div>
      </div>

      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showEditor = false">Cancel</button>
        <button
          class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 disabled:opacity-50 transition-colors"
          :disabled="!canSaveTier || saving"
          @click="saveTier"
        >
          {{ saving ? 'Saving...' : editingTier ? 'Update Tier' : 'Create Tier' }}
        </button>
      </template>
    </Modal>

    <!-- D15 Members Panel (slide-in) -->
    <Teleport to="body">
      <div v-if="showMembers" class="fixed inset-0 z-50 flex justify-end">
        <div class="absolute inset-0 bg-black/60" @click="showMembers = false" />
        <div class="relative w-[480px] bg-[#111B27] border-l border-[#1E293B] h-full flex flex-col">
          <!-- Panel Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-[#1E293B]">
            <div>
              <h2 class="text-lg font-semibold text-[#F1F5F9]">{{ membersTier?.name }} Members</h2>
              <p class="text-xs text-[#94A3B8]">{{ membersTotal }} total members</p>
            </div>
            <button class="text-[#94A3B8] hover:text-[#F1F5F9] transition-colors" @click="showMembers = false">
              <span class="material-symbols-rounded text-xl">close</span>
            </button>
          </div>

          <!-- Search + Add -->
          <div class="px-6 py-3 border-b border-[#1E293B] space-y-3">
            <div class="relative">
              <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-[#64748B] text-lg">search</span>
              <input
                v-model="membersSearch"
                type="text"
                placeholder="Search by wallet address..."
                class="w-full pl-10 pr-4 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                @input="membersPage = 1; fetchMembers()"
              />
            </div>
            <div class="flex gap-2">
              <input
                v-model="newMemberAddress"
                type="text"
                placeholder="0x..."
                class="flex-1 px-3 py-2 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
              />
              <button
                class="px-3 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 disabled:opacity-50 transition-colors"
                :disabled="!newMemberAddress.trim() || addingMember"
                @click="addMember"
              >
                Add
              </button>
            </div>
            <div class="flex gap-2">
              <label class="px-3 py-1.5 text-xs text-wl border border-wl/30 rounded-lg hover:bg-wl/10 transition-colors cursor-pointer">
                Import CSV
                <input type="file" accept=".csv" class="hidden" @change="handleCsvUpload" />
              </label>
              <button
                class="px-3 py-1.5 text-xs text-text-secondary border border-[#1E293B] rounded-lg hover:text-text-primary transition-colors"
                @click="exportMembers"
              >
                Export CSV
              </button>
            </div>
          </div>

          <!-- Member list -->
          <div class="flex-1 overflow-y-auto">
            <div v-if="membersLoading" class="p-6 space-y-3">
              <div v-for="i in 5" :key="i" class="h-12 bg-[#1E293B] rounded animate-pulse" />
            </div>
            <div v-else-if="members.length === 0" class="p-12 text-center">
              <span class="material-symbols-rounded text-3xl text-[#64748B] block mb-2">group_off</span>
              <p class="text-sm text-[#64748B]">No members found</p>
            </div>
            <div v-else>
              <div
                v-for="member in members"
                :key="member.id"
                class="flex items-center justify-between px-6 py-3 border-b border-[#1E293B] hover:bg-white/[0.02] transition-colors"
              >
                <div>
                  <div class="text-sm font-mono text-[#F1F5F9]">{{ truncateAddress(member.wallet_address) }}</div>
                  <div class="text-xs text-[#64748B]">Joined {{ formatDate(member.joined_at) }}</div>
                </div>
                <div class="flex items-center gap-2">
                  <span
                    class="px-2 py-0.5 text-xs rounded-full"
                    :class="member.status === 'active' ? 'bg-[#0A2E1A] text-[#16A34A]' : 'bg-[#1E293B] text-[#64748B]'"
                  >
                    {{ member.status }}
                  </span>
                  <button
                    class="p-1 rounded hover:bg-white/5 transition-colors"
                    title="Remove"
                    @click="removeMember(member.id)"
                  >
                    <span class="material-symbols-rounded text-sm text-status-paused">close</span>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Pagination -->
          <div v-if="membersTotal > 20" class="px-6 py-3 border-t border-[#1E293B] flex items-center justify-between">
            <span class="text-xs text-[#64748B]">Page {{ membersPage }} of {{ Math.ceil(membersTotal / 20) }}</span>
            <div class="flex gap-1">
              <button
                class="px-2 py-1 text-xs rounded text-[#94A3B8] hover:bg-white/5 disabled:opacity-30"
                :disabled="membersPage <= 1"
                @click="membersPage--; fetchMembers()"
              >
                Prev
              </button>
              <button
                class="px-2 py-1 text-xs rounded text-[#94A3B8] hover:bg-white/5 disabled:opacity-30"
                :disabled="membersPage >= Math.ceil(membersTotal / 20)"
                @click="membersPage++; fetchMembers()"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
