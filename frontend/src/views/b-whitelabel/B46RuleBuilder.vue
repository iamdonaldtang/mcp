<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// --- Types ---
interface RuleItem {
  id: string
  name: string
  status: CampaignStatus
  trigger_event: string
  trigger_category: 'onchain' | 'community'
  trigger_config: Record<string, unknown>
  action_type: 'award_points' | 'grant_badge' | 'upgrade_tier' | 'webhook'
  action_config: Record<string, unknown>
  frequency: 'once' | 'daily' | 'unlimited'
  triggered_today: number
  total_triggered: number
  created_at: string
}

interface ContractOption {
  id: string
  name: string
  address: string
  chain: string
}

// --- State ---
const loading = ref(true)
const rules = ref<RuleItem[]>([])
const totalRules = ref(0)
const page = ref(1)
const pageSize = 20

const stats = ref({
  total: 0,
  active: 0,
  triggered_today: 0,
  points_distributed_24h: 0,
})

// Anti-Sybil
const antiSybilExpanded = ref(false)
const antiSybil = ref({
  min_wallet_age: 30,
  min_transactions: 5,
  bot_detection: true,
})
const savingAntiSybil = ref(false)

// Filters
const filterTab = ref<string>('all')
const searchQuery = ref('')
let searchTimeout: ReturnType<typeof setTimeout> | null = null

// Modal
const showRuleEditor = ref(false)
const editingRule = ref<RuleItem | null>(null)
const saving = ref(false)

// Preset
const showPresetMenu = ref(false)
const presets = [
  { label: 'Welcome Bonus', trigger: 'daily_login', action: 'award_points', frequency: 'once' as const },
  { label: 'Daily Login Reward', trigger: 'daily_login', action: 'award_points', frequency: 'daily' as const },
  { label: 'Trading Reward', trigger: 'token_swap', action: 'award_points', frequency: 'unlimited' as const },
  { label: 'Referral Bonus', trigger: 'referral_success', action: 'award_points', frequency: 'unlimited' as const },
  { label: 'Milestone Achievement', trigger: 'milestone_reached', action: 'grant_badge', frequency: 'once' as const },
]

// Editor form
const ruleForm = ref(getEmptyForm())

function getEmptyForm() {
  return {
    name: '',
    trigger_category: 'community' as 'onchain' | 'community',
    trigger_event: '',
    chain: '',
    contract_id: '',
    condition_value: '',
    related_entity: '',
    action_type: 'award_points' as 'award_points' | 'grant_badge' | 'upgrade_tier' | 'webhook',
    point_type: 'EXP',
    point_value: 10,
    point_multiplier: 1,
    badge_id: '',
    tier_id: '',
    webhook_url: '',
    frequency: 'unlimited' as 'once' | 'daily' | 'unlimited',
  }
}

const onchainEvents = [
  { value: 'token_swap', label: 'Token Swap' },
  { value: 'add_liquidity', label: 'Add Liquidity' },
  { value: 'nft_mint', label: 'NFT Mint' },
  { value: 'nft_transfer', label: 'NFT Transfer' },
  { value: 'staking_deposit', label: 'Staking Deposit' },
  { value: 'contract_interaction', label: 'Contract Interaction' },
]

const communityEvents = [
  { value: 'task_completed', label: 'Task Completed' },
  { value: 'daily_login', label: 'Daily Login' },
  { value: 'referral_success', label: 'Referral Success' },
  { value: 'milestone_reached', label: 'Milestone Reached' },
  { value: 'level_up', label: 'Level Up' },
]

const availableEvents = computed(() =>
  ruleForm.value.trigger_category === 'onchain' ? onchainEvents : communityEvents
)

const contracts = ref<ContractOption[]>([])

const triggerSummary = computed(() => {
  const evt = availableEvents.value.find(e => e.value === ruleForm.value.trigger_event)
  return evt?.label || 'Select trigger...'
})

const actionSummary = computed(() => {
  switch (ruleForm.value.action_type) {
    case 'award_points': return `+${ruleForm.value.point_value} ${ruleForm.value.point_type}`
    case 'grant_badge': return 'Grant Badge'
    case 'upgrade_tier': return 'Upgrade Tier'
    case 'webhook': return 'Fire Webhook'
    default: return 'Select action...'
  }
})

const canSave = computed(() => ruleForm.value.name.trim().length > 0 && ruleForm.value.trigger_event !== '')

// --- Search debounce ---
watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchRules()
  }, 300)
})

// --- API ---
onMounted(async () => {
  await Promise.all([fetchStats(), fetchRules(), fetchContracts()])
  loading.value = false
})

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/whitelabel/rules/stats')
    if (res.data.data) stats.value = res.data.data
  } catch { /* defaults */ }
}

async function fetchRules() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterTab.value !== 'all') params.status = filterTab.value
    if (searchQuery.value) params.search = searchQuery.value
    const res = await api.get('/api/v1/whitelabel/rules', { params })
    rules.value = res.data.data?.items || []
    totalRules.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

async function fetchContracts() {
  try {
    const res = await api.get('/api/v1/whitelabel/contracts', { params: { page_size: 100, status: 'verified' } })
    contracts.value = (res.data.data?.items || []).map((c: Record<string, string>) => ({
      id: c.id, name: c.name, address: c.address, chain: c.chain,
    }))
  } catch { /* empty */ }
}

async function saveAntiSybil() {
  savingAntiSybil.value = true
  try {
    await api.put('/api/v1/whitelabel/rules/anti-sybil', antiSybil.value)
  } catch { /* TODO: toast */ }
  finally { savingAntiSybil.value = false }
}

function openCreate() {
  editingRule.value = null
  ruleForm.value = getEmptyForm()
  showRuleEditor.value = true
}

function openEdit(rule: RuleItem) {
  editingRule.value = rule
  ruleForm.value = {
    name: rule.name,
    trigger_category: rule.trigger_category,
    trigger_event: rule.trigger_event,
    chain: (rule.trigger_config.chain as string) || '',
    contract_id: (rule.trigger_config.contract_id as string) || '',
    condition_value: (rule.trigger_config.condition_value as string) || '',
    related_entity: (rule.trigger_config.related_entity as string) || '',
    action_type: rule.action_type,
    point_type: (rule.action_config.point_type as string) || 'EXP',
    point_value: (rule.action_config.point_value as number) || 10,
    point_multiplier: (rule.action_config.point_multiplier as number) || 1,
    badge_id: (rule.action_config.badge_id as string) || '',
    tier_id: (rule.action_config.tier_id as string) || '',
    webhook_url: (rule.action_config.webhook_url as string) || '',
    frequency: rule.frequency,
  }
  showRuleEditor.value = true
}

async function saveRule() {
  if (!canSave.value) return
  saving.value = true
  const payload = {
    name: ruleForm.value.name,
    trigger_category: ruleForm.value.trigger_category,
    trigger_event: ruleForm.value.trigger_event,
    trigger_config: ruleForm.value.trigger_category === 'onchain'
      ? { chain: ruleForm.value.chain, contract_id: ruleForm.value.contract_id, condition_value: ruleForm.value.condition_value }
      : { related_entity: ruleForm.value.related_entity },
    action_type: ruleForm.value.action_type,
    action_config: buildActionConfig(),
    frequency: ruleForm.value.frequency,
  }
  try {
    if (editingRule.value) {
      await api.put(`/api/v1/whitelabel/rules/${editingRule.value.id}`, payload)
    } else {
      await api.post('/api/v1/whitelabel/rules', payload)
    }
    showRuleEditor.value = false
    await Promise.all([fetchStats(), fetchRules()])
  } catch { /* TODO: toast */ }
  finally { saving.value = false }
}

function buildActionConfig(): Record<string, unknown> {
  switch (ruleForm.value.action_type) {
    case 'award_points': return { point_type: ruleForm.value.point_type, point_value: ruleForm.value.point_value, point_multiplier: ruleForm.value.point_multiplier }
    case 'grant_badge': return { badge_id: ruleForm.value.badge_id }
    case 'upgrade_tier': return { tier_id: ruleForm.value.tier_id }
    case 'webhook': return { webhook_url: ruleForm.value.webhook_url }
    default: return {}
  }
}

async function toggleRule(rule: RuleItem) {
  const newStatus = rule.status === 'active' ? 'paused' : 'active'
  const oldStatus = rule.status
  rule.status = newStatus as CampaignStatus
  try {
    await api.put(`/api/v1/whitelabel/rules/${rule.id}`, { status: newStatus })
    await fetchStats()
  } catch { rule.status = oldStatus as CampaignStatus }
}

async function duplicateRule(rule: RuleItem) {
  try {
    await api.post(`/api/v1/whitelabel/rules/${rule.id}/duplicate`)
    await Promise.all([fetchStats(), fetchRules()])
  } catch { /* TODO: toast */ }
}

async function deleteRule(rule: RuleItem) {
  if (!confirm(`Delete rule "${rule.name}"? It will be recoverable for 30 days.`)) return
  try {
    await api.delete(`/api/v1/whitelabel/rules/${rule.id}`)
    await Promise.all([fetchStats(), fetchRules()])
  } catch { /* TODO: toast */ }
}

function applyPreset(preset: typeof presets[0]) {
  ruleForm.value = getEmptyForm()
  ruleForm.value.name = preset.label
  ruleForm.value.trigger_category = ['token_swap', 'add_liquidity', 'nft_mint', 'nft_transfer', 'staking_deposit', 'contract_interaction'].includes(preset.trigger) ? 'onchain' : 'community'
  ruleForm.value.trigger_event = preset.trigger
  ruleForm.value.action_type = preset.action as typeof ruleForm.value.action_type
  ruleForm.value.frequency = preset.frequency
  showPresetMenu.value = false
  editingRule.value = null
  showRuleEditor.value = true
}

function formatEvent(event: string): string {
  return event.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')
}

function formatAction(type: string): string {
  return type.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')
}

function formatDate(d: string): string {
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function onPageChange(p: number) {
  page.value = p
  fetchRules()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Activity Rule Builder</h1>
        <p class="text-sm text-text-secondary">Create automated rules that trigger rewards based on user actions</p>
      </div>
      <div class="flex items-center gap-3">
        <!-- Preset dropdown -->
        <div class="relative">
          <button
            class="px-4 py-2 border border-wl text-wl text-sm font-medium rounded-lg hover:bg-wl/10 transition-colors"
            @click="showPresetMenu = !showPresetMenu"
          >
            <span class="material-symbols-rounded text-sm align-middle mr-1">auto_fix_high</span>
            Use Preset
          </button>
          <div
            v-if="showPresetMenu"
            class="absolute right-0 top-full mt-1 w-56 bg-[#111B27] border border-[#1E293B] rounded-xl shadow-xl z-20 py-1"
          >
            <button
              v-for="preset in presets"
              :key="preset.label"
              class="w-full text-left px-4 py-2.5 text-sm text-[#F1F5F9] hover:bg-white/5 transition-colors"
              @click="applyPreset(preset)"
            >
              {{ preset.label }}
            </button>
          </div>
        </div>
        <button
          class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors"
          @click="openCreate"
        >
          + Create Rule
        </button>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard label="Total Rules" :value="stats.total" icon="tune" icon-color="#9B7EE0" />
      <StatsCard label="Active Rules" :value="stats.active" icon="check_circle" icon-color="#16A34A" />
      <StatsCard label="Triggered Today" :value="stats.triggered_today" icon="bolt" icon-color="#F59E0B" />
      <StatsCard label="Points Distributed (24h)" :value="stats.points_distributed_24h" icon="toll" icon-color="#5D7EF1" />
    </div>

    <!-- Anti-Sybil Config -->
    <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
      <button
        class="w-full flex items-center justify-between px-5 py-4 hover:bg-white/2 transition-colors"
        @click="antiSybilExpanded = !antiSybilExpanded"
      >
        <div class="flex items-center gap-3">
          <span class="material-symbols-rounded text-lg text-[#F59E0B]">shield</span>
          <span class="text-sm font-semibold text-text-primary">Anti-Sybil Configuration</span>
          <span class="px-2 py-0.5 text-xs rounded bg-[#1F1A08] text-[#F59E0B]">Recommended</span>
        </div>
        <span class="material-symbols-rounded text-lg text-text-muted transition-transform" :class="antiSybilExpanded ? 'rotate-180' : ''">
          keyboard_arrow_down
        </span>
      </button>
      <div v-if="antiSybilExpanded" class="px-5 pb-5 border-t border-border pt-4">
        <div class="grid grid-cols-3 gap-6">
          <!-- Min Wallet Age -->
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Min Wallet Age (days)</label>
            <input
              v-model.number="antiSybil.min_wallet_age"
              type="number"
              min="0"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
            />
          </div>
          <!-- Min Transactions -->
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Min Transactions</label>
            <input
              v-model.number="antiSybil.min_transactions"
              type="number"
              min="0"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
            />
          </div>
          <!-- Bot Detection -->
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Bot Detection</label>
            <button
              class="flex items-center gap-2 px-4 py-2.5 rounded-lg border text-sm transition-colors"
              :class="antiSybil.bot_detection
                ? 'bg-[#0A2E1A] border-[#16A34A]/30 text-[#16A34A]'
                : 'bg-[#0A0F1A] border-[#1E293B] text-[#94A3B8]'"
              @click="antiSybil.bot_detection = !antiSybil.bot_detection"
            >
              <span class="material-symbols-rounded text-base">
                {{ antiSybil.bot_detection ? 'toggle_on' : 'toggle_off' }}
              </span>
              {{ antiSybil.bot_detection ? 'Enabled' : 'Disabled' }}
            </button>
          </div>
        </div>
        <div class="mt-4 flex justify-end">
          <button
            class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors disabled:opacity-50"
            :disabled="savingAntiSybil"
            @click="saveAntiSybil"
          >
            {{ savingAntiSybil ? 'Saving...' : 'Save Anti-Sybil Settings' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Filter tabs + search -->
    <div class="flex items-center gap-4">
      <div class="flex gap-2">
        <button
          v-for="tab in ['all', 'active', 'draft', 'paused']"
          :key="tab"
          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
          :class="filterTab === tab ? 'bg-wl text-white' : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
          @click="filterTab = tab; page = 1; fetchRules()"
        >
          {{ tab === 'all' ? 'All' : tab.charAt(0).toUpperCase() + tab.slice(1) }}
        </button>
      </div>
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search rules..."
            class="w-full pl-10 pr-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-wl focus:outline-none"
          />
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-border">
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Rule Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-36">Trigger Event</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Action Type</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Triggered Today</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Total</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Created</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="8" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse" /></td>
          </tr>
          <tr v-else-if="rules.length === 0">
            <td colspan="8" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">tune</span>
              <p class="text-sm text-text-muted">No rules yet</p>
              <button class="mt-2 text-xs text-wl hover:underline" @click="openCreate">+ Create First Rule</button>
            </td>
          </tr>
          <tr
            v-else
            v-for="rule in rules"
            :key="rule.id"
            class="border-b border-border last:border-b-0 hover:bg-white/[0.02] transition-colors"
          >
            <td class="px-4 py-3">
              <span class="text-sm font-medium text-text-primary">{{ rule.name }}</span>
            </td>
            <td class="px-4 py-3">
              <StatusBadge :status="rule.status" />
            </td>
            <td class="px-4 py-3">
              <span class="px-2 py-0.5 text-xs rounded bg-[#0A0F1A] text-text-secondary">{{ formatEvent(rule.trigger_event) }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-sm text-text-secondary">{{ formatAction(rule.action_type) }}</span>
            </td>
            <td class="px-4 py-3 text-right text-sm text-text-secondary">{{ rule.triggered_today.toLocaleString() }}</td>
            <td class="px-4 py-3 text-right text-sm text-text-secondary">{{ rule.total_triggered.toLocaleString() }}</td>
            <td class="px-4 py-3 text-sm text-text-muted">{{ formatDate(rule.created_at) }}</td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit" @click="openEdit(rule)">
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  :title="rule.status === 'active' ? 'Pause' : 'Activate'"
                  @click="toggleRule(rule)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">
                    {{ rule.status === 'active' ? 'pause_circle' : 'play_circle' }}
                  </span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Duplicate" @click="duplicateRule(rule)">
                  <span class="material-symbols-rounded text-base text-text-muted">content_copy</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteRule(rule)">
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination v-if="totalRules > pageSize" :page="page" :page-size="pageSize" :total="totalRules" @update:page="onPageChange" />
    </div>

    <!-- Click-away for preset menu -->
    <div v-if="showPresetMenu" class="fixed inset-0 z-10" @click="showPresetMenu = false" />

    <!-- D13 Activity Rule Editor Modal -->
    <Modal :open="showRuleEditor" :title="editingRule ? 'Edit Rule' : 'Create Rule'" max-width="800px" @close="showRuleEditor = false">
      <div class="flex gap-6">
        <!-- Left side: form -->
        <div class="flex-1 space-y-5 min-w-0">
          <!-- Rule Name -->
          <div>
            <label class="block text-sm text-text-secondary mb-1.5">Rule Name</label>
            <input
              v-model="ruleForm.name"
              type="text"
              maxlength="50"
              placeholder="e.g. Welcome Bonus"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
            />
            <span class="text-xs text-text-muted mt-1 block">{{ ruleForm.name.length }}/50</span>
          </div>

          <!-- TRIGGER Section -->
          <div class="border border-[#1E293B] rounded-lg p-4 space-y-4">
            <div class="flex items-center gap-2 mb-1">
              <span class="material-symbols-rounded text-base text-[#F59E0B]">bolt</span>
              <span class="text-sm font-semibold text-text-primary uppercase tracking-wider">Trigger</span>
            </div>

            <!-- Category toggle -->
            <div class="flex gap-2">
              <button
                v-for="cat in (['community', 'onchain'] as const)"
                :key="cat"
                class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
                :class="ruleForm.trigger_category === cat ? 'bg-wl text-white' : 'bg-[#0A0F1A] border border-[#1E293B] text-[#94A3B8]'"
                @click="ruleForm.trigger_category = cat; ruleForm.trigger_event = ''"
              >
                {{ cat === 'community' ? 'Community' : 'On-chain' }}
              </button>
            </div>

            <!-- Event dropdown -->
            <div>
              <label class="block text-xs text-text-muted mb-1">Event</label>
              <select
                v-model="ruleForm.trigger_event"
                class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
              >
                <option value="" disabled>Select event...</option>
                <option v-for="evt in availableEvents" :key="evt.value" :value="evt.value">{{ evt.label }}</option>
              </select>
            </div>

            <!-- On-chain specifics -->
            <template v-if="ruleForm.trigger_category === 'onchain'">
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs text-text-muted mb-1">Chain</label>
                  <select
                    v-model="ruleForm.chain"
                    class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
                  >
                    <option value="" disabled>Select chain</option>
                    <option value="ethereum">Ethereum</option>
                    <option value="bsc">BSC</option>
                    <option value="polygon">Polygon</option>
                    <option value="arbitrum">Arbitrum</option>
                    <option value="optimism">Optimism</option>
                    <option value="base">Base</option>
                    <option value="avalanche">Avalanche</option>
                  </select>
                </div>
                <div>
                  <label class="block text-xs text-text-muted mb-1">Contract</label>
                  <select
                    v-model="ruleForm.contract_id"
                    class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
                  >
                    <option value="" disabled>From Contract Registry</option>
                    <option v-for="c in contracts" :key="c.id" :value="c.id">{{ c.name }} ({{ c.address.slice(0, 6) }}...)</option>
                  </select>
                </div>
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1">Condition (e.g. min amount)</label>
                <input
                  v-model="ruleForm.condition_value"
                  type="text"
                  placeholder="e.g. min_amount >= 100"
                  class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
            </template>

            <!-- Community specifics -->
            <template v-if="ruleForm.trigger_category === 'community' && ['task_completed', 'milestone_reached', 'level_up'].includes(ruleForm.trigger_event)">
              <div>
                <label class="block text-xs text-text-muted mb-1">Related Entity (optional)</label>
                <input
                  v-model="ruleForm.related_entity"
                  type="text"
                  placeholder="Specific task/milestone/level ID"
                  class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
            </template>
          </div>

          <!-- ACTION Section -->
          <div class="border border-[#1E293B] rounded-lg p-4 space-y-4">
            <div class="flex items-center gap-2 mb-1">
              <span class="material-symbols-rounded text-base text-[#9B7EE0]">auto_awesome</span>
              <span class="text-sm font-semibold text-text-primary uppercase tracking-wider">Action</span>
            </div>

            <div>
              <label class="block text-xs text-text-muted mb-1">Type</label>
              <select
                v-model="ruleForm.action_type"
                class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
              >
                <option value="award_points">Award Points</option>
                <option value="grant_badge">Grant Badge</option>
                <option value="upgrade_tier">Upgrade Tier</option>
                <option value="webhook">Webhook</option>
              </select>
            </div>

            <!-- Award Points config -->
            <template v-if="ruleForm.action_type === 'award_points'">
              <div class="grid grid-cols-3 gap-3">
                <div>
                  <label class="block text-xs text-text-muted mb-1">Point Type</label>
                  <select
                    v-model="ruleForm.point_type"
                    class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
                  >
                    <option value="EXP">EXP</option>
                    <option value="GEM">GEM</option>
                    <option value="COIN">COIN</option>
                  </select>
                </div>
                <div>
                  <label class="block text-xs text-text-muted mb-1">Value</label>
                  <input
                    v-model.number="ruleForm.point_value"
                    type="number"
                    min="1"
                    class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
                  />
                </div>
                <div>
                  <label class="block text-xs text-text-muted mb-1">Multiplier</label>
                  <input
                    v-model.number="ruleForm.point_multiplier"
                    type="number"
                    min="0.1"
                    step="0.1"
                    class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] focus:border-wl focus:outline-none"
                  />
                </div>
              </div>
            </template>

            <!-- Grant Badge config -->
            <template v-if="ruleForm.action_type === 'grant_badge'">
              <div>
                <label class="block text-xs text-text-muted mb-1">Badge</label>
                <input
                  v-model="ruleForm.badge_id"
                  type="text"
                  placeholder="Badge ID"
                  class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
            </template>

            <!-- Upgrade Tier config -->
            <template v-if="ruleForm.action_type === 'upgrade_tier'">
              <div>
                <label class="block text-xs text-text-muted mb-1">Privilege Tier</label>
                <input
                  v-model="ruleForm.tier_id"
                  type="text"
                  placeholder="Tier ID"
                  class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
            </template>

            <!-- Webhook config -->
            <template v-if="ruleForm.action_type === 'webhook'">
              <div>
                <label class="block text-xs text-text-muted mb-1">Webhook URL</label>
                <input
                  v-model="ruleForm.webhook_url"
                  type="url"
                  placeholder="https://..."
                  class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-sm text-[#F1F5F9] placeholder-[#64748B] focus:border-wl focus:outline-none"
                />
              </div>
            </template>
          </div>

          <!-- Frequency -->
          <div>
            <label class="block text-sm text-text-secondary mb-2">Frequency</label>
            <div class="flex gap-3">
              <label
                v-for="freq in (['once', 'daily', 'unlimited'] as const)"
                :key="freq"
                class="flex items-center gap-2 px-3 py-2 rounded-lg border cursor-pointer transition-colors"
                :class="ruleForm.frequency === freq ? 'border-wl bg-wl/10 text-wl' : 'border-[#1E293B] text-[#94A3B8]'"
              >
                <input v-model="ruleForm.frequency" type="radio" :value="freq" class="sr-only" />
                <span class="text-sm">{{ freq === 'once' ? 'Once per user' : freq === 'daily' ? 'Daily' : 'Unlimited' }}</span>
              </label>
            </div>
          </div>
        </div>

        <!-- Right side: Rule Summary -->
        <div class="w-56 shrink-0">
          <div class="bg-[#0A0F1A] border border-[#1E293B] rounded-xl p-4 sticky top-0">
            <h3 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Rule Summary</h3>

            <!-- When -->
            <div class="mb-4">
              <div class="flex items-center gap-2 mb-2">
                <span class="w-6 h-6 rounded-full bg-[#1F1A08] flex items-center justify-center">
                  <span class="material-symbols-rounded text-xs text-[#F59E0B]">bolt</span>
                </span>
                <span class="text-xs text-text-muted uppercase">When</span>
              </div>
              <p class="text-sm text-text-primary pl-8">{{ triggerSummary }}</p>
            </div>

            <!-- Arrow -->
            <div class="flex justify-center my-2">
              <span class="material-symbols-rounded text-lg text-text-muted">arrow_downward</span>
            </div>

            <!-- Then -->
            <div class="mb-4">
              <div class="flex items-center gap-2 mb-2">
                <span class="w-6 h-6 rounded-full bg-[#1E1033] flex items-center justify-center">
                  <span class="material-symbols-rounded text-xs text-[#9B7EE0]">auto_awesome</span>
                </span>
                <span class="text-xs text-text-muted uppercase">Then</span>
              </div>
              <p class="text-sm text-text-primary pl-8">{{ actionSummary }}</p>
            </div>

            <!-- Frequency -->
            <div class="pt-3 border-t border-[#1E293B]">
              <div class="flex items-center gap-2">
                <span class="material-symbols-rounded text-sm text-text-muted">repeat</span>
                <span class="text-xs text-text-secondary capitalize">{{ ruleForm.frequency === 'once' ? 'Once per user' : ruleForm.frequency }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showRuleEditor = false">Cancel</button>
        <button
          class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 disabled:opacity-50 transition-colors"
          :disabled="!canSave || saving"
          @click="saveRule"
        >
          {{ saving ? 'Saving...' : editingRule ? 'Update Rule' : 'Create Rule' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
