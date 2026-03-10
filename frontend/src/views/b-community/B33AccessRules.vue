<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { api } from '../../api/client'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// --- Types ---
type RuleType = 'token_gate' | 'nft_hold' | 'level_requirement' | 'invite_only'

interface AccessRuleItem {
  id: string
  name: string
  type: RuleType
  status: CampaignStatus
  condition_summary: string
  priority: number
  config: {
    chain?: string
    contract_address?: string
    min_balance?: number
    collection_address?: string
    min_count?: number
    min_level?: number
    invite_codes?: string
    whitelist_file?: string
  }
  denial_message: string
  created_at: string
}

interface StatsData {
  total_rules: number
  active_rules: number
  users_validated: number
  access_denied: number
}

// --- State ---
const loading = ref(true)
const items = ref<AccessRuleItem[]>([])
const totalItems = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<StatsData>({ total_rules: 0, active_rules: 0, users_validated: 0, access_denied: 0 })

// Filters
const filterStatus = ref<string>('all')
const statusTabs = ['all', 'active', 'paused', 'draft'] as const

// Drag state
const dragIndex = ref<number | null>(null)
const dropIndex = ref<number | null>(null)

// Modal
const showCreate = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  type: 'token_gate' as RuleType,
  chain: '',
  contract_address: '',
  min_balance: 1,
  collection_address: '',
  min_count: 1,
  min_level: 1,
  invite_codes: '',
  denial_message: '',
})

// Preview rule
const showPreview = ref(false)
const previewAddress = ref('')
const previewResult = ref<{ passed: boolean; reason: string } | null>(null)
const previewLoading = ref(false)

const ruleTypeLabels: Record<RuleType, string> = {
  token_gate: 'Token Gate',
  nft_hold: 'NFT Hold',
  level_requirement: 'Level Requirement',
  invite_only: 'Invite Only',
}

const ruleTypeIcons: Record<RuleType, string> = {
  token_gate: 'token',
  nft_hold: 'image',
  level_requirement: 'military_tech',
  invite_only: 'vpn_key',
}

const chainOptions = ['Ethereum', 'Polygon', 'BNB Chain', 'Arbitrum', 'Optimism', 'Avalanche', 'Base', 'Solana'] as const

const formValid = computed(() => {
  if (form.value.name.trim().length < 1 || form.value.name.trim().length > 50) return false
  switch (form.value.type) {
    case 'token_gate':
      return form.value.chain !== '' && form.value.contract_address.trim().length > 0 && form.value.min_balance > 0
    case 'nft_hold':
      return form.value.chain !== '' && form.value.collection_address.trim().length > 0 && form.value.min_count > 0
    case 'level_requirement':
      return form.value.min_level >= 1
    case 'invite_only':
      return form.value.invite_codes.trim().length > 0
  }
  return false
})

// --- Fetch ---
onMounted(async () => {
  await Promise.all([fetchStats(), fetchItems()])
  loading.value = false
})

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/settings/access-rules/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

async function fetchItems() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    const res = await api.get('/api/v1/community/settings/access-rules', { params })
    items.value = res.data.data?.items || []
    totalItems.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

watch(filterStatus, () => { page.value = 1; fetchItems() })

function onPageChange(p: number) {
  page.value = p
  fetchItems()
}

// --- Drag & Drop ---
function onDragStart(index: number) {
  dragIndex.value = index
}

function onDragOver(e: DragEvent, index: number) {
  e.preventDefault()
  dropIndex.value = index
}

function onDragEnd() {
  dragIndex.value = null
  dropIndex.value = null
}

async function onDrop(targetIndex: number) {
  if (dragIndex.value === null || dragIndex.value === targetIndex) {
    onDragEnd()
    return
  }
  const moved = items.value.splice(dragIndex.value, 1)[0]
  items.value.splice(targetIndex, 0, moved)
  // Update priority numbers
  items.value.forEach((item, i) => { item.priority = i + 1 })
  onDragEnd()
  try {
    await api.put('/api/v1/community/settings/access-rules/reorder', {
      order: items.value.map(i => i.id),
    })
  } catch {
    await fetchItems()
  }
}

// --- Actions ---
function openCreate() {
  editingId.value = null
  form.value = {
    name: '', type: 'token_gate', chain: '', contract_address: '',
    min_balance: 1, collection_address: '', min_count: 1, min_level: 1,
    invite_codes: '', denial_message: '',
  }
  showCreate.value = true
}

function openEdit(item: AccessRuleItem) {
  editingId.value = item.id
  form.value = {
    name: item.name,
    type: item.type,
    chain: item.config.chain || '',
    contract_address: item.config.contract_address || '',
    min_balance: item.config.min_balance || 1,
    collection_address: item.config.collection_address || '',
    min_count: item.config.min_count || 1,
    min_level: item.config.min_level || 1,
    invite_codes: item.config.invite_codes || '',
    denial_message: item.denial_message || '',
  }
  showCreate.value = true
}

function buildPayload() {
  const base = { name: form.value.name, type: form.value.type, denial_message: form.value.denial_message }
  let config: Record<string, unknown> = {}
  switch (form.value.type) {
    case 'token_gate':
      config = { chain: form.value.chain, contract_address: form.value.contract_address, min_balance: form.value.min_balance }
      break
    case 'nft_hold':
      config = { chain: form.value.chain, collection_address: form.value.collection_address, min_count: form.value.min_count }
      break
    case 'level_requirement':
      config = { min_level: form.value.min_level }
      break
    case 'invite_only':
      config = { invite_codes: form.value.invite_codes }
      break
  }
  return { ...base, config }
}

async function saveForm() {
  if (!formValid.value || saving.value) return
  saving.value = true
  try {
    const payload = buildPayload()
    if (editingId.value) {
      await api.put(`/api/v1/community/settings/access-rules/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/community/settings/access-rules', payload)
    }
    showCreate.value = false
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function toggleStatus(item: AccessRuleItem) {
  const newStatus = item.status === 'active' ? 'paused' : 'active'
  try {
    await api.put(`/api/v1/community/settings/access-rules/${item.id}`, { status: newStatus })
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

async function deleteItem(item: AccessRuleItem) {
  const userCount = stats.value.users_validated
  const msg = userCount > 0
    ? `Delete "${item.name}"? ${userCount} users have been validated against this rule. This cannot be undone.`
    : `Delete "${item.name}"? This action cannot be undone.`
  if (!confirm(msg)) return
  try {
    await api.delete(`/api/v1/community/settings/access-rules/${item.id}`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

// --- Preview ---
function openPreviewRule() {
  previewAddress.value = ''
  previewResult.value = null
  showPreview.value = true
}

async function runPreview() {
  if (!previewAddress.value.trim()) return
  previewLoading.value = true
  previewResult.value = null
  try {
    const res = await api.post('/api/v1/community/settings/access-rules/preview', {
      wallet_address: previewAddress.value.trim(),
    })
    previewResult.value = res.data.data
  } catch {
    previewResult.value = { passed: false, reason: 'Preview failed. Check wallet address format.' }
  }
  previewLoading.value = false
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatNumber(n: number) {
  return n.toLocaleString()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[#F1F5F9] mb-1">Access Rules</h1>
        <p class="text-sm text-[#94A3B8]">Control who can join your community with layered access requirements</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          class="px-4 py-2 text-sm font-medium text-[#94A3B8] border border-[#1E293B] rounded-lg hover:text-[#F1F5F9] hover:border-[#94A3B8] transition-colors"
          @click="openPreviewRule"
        >
          <span class="material-symbols-rounded text-base align-middle mr-1">preview</span>
          Preview Rule
        </button>
        <button
          class="px-4 py-2 bg-[#48BB78] text-black text-sm font-medium rounded-lg hover:bg-[#48BB78]/90 transition-colors"
          @click="openCreate"
        >
          + Create Rule
        </button>
      </div>
    </div>

    <!-- AND relationship info -->
    <div class="flex items-start gap-3 p-3 bg-[#48BB78]/5 border border-[#48BB78]/20 rounded-lg">
      <span class="material-symbols-rounded text-lg text-[#48BB78] mt-0.5">info</span>
      <div>
        <p class="text-sm text-[#F1F5F9] font-medium">Rules use AND logic</p>
        <p class="text-xs text-[#94A3B8] mt-0.5">Users must pass <strong class="text-[#F1F5F9]">all active rules</strong> to gain access. Rules are evaluated in priority order (drag to reorder).</p>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Total Rules</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.total_rules) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Active Rules</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.active_rules) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Users Validated</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.users_validated) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Access Denied</p>
        <p class="text-2xl font-bold text-[#DC2626]">{{ formatNumber(stats.access_denied) }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="flex border-b border-[#1E293B]">
        <button
          v-for="tab in statusTabs"
          :key="tab"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="filterStatus === tab
            ? 'text-[#48BB78]'
            : 'text-[#94A3B8] hover:text-[#F1F5F9]'"
          @click="filterStatus = tab"
        >
          {{ tab === 'all' ? 'All' : tab.charAt(0).toUpperCase() + tab.slice(1) }}
          <span
            v-if="filterStatus === tab"
            class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#48BB78] rounded-t"
          />
        </button>
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-[#111B27] border border-[#1E293B] rounded-xl overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-[#1E293B]">
            <th class="px-2 py-3 w-10"></th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider">Rule Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-36">Type</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-28">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider">Condition</th>
            <th class="px-4 py-3 text-center text-xs font-semibold text-[#64748B] uppercase tracking-wider w-20">Priority</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-[#64748B] uppercase tracking-wider w-28">Actions</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading skeleton -->
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="7" class="px-4 py-4"><div class="h-4 bg-[#1E293B] rounded animate-pulse" /></td>
          </tr>
          <!-- Empty state -->
          <tr v-else-if="items.length === 0">
            <td colspan="7" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-[#64748B] block mb-2">shield</span>
              <p class="text-sm text-[#64748B]">No access rules configured</p>
              <p class="text-xs text-[#64748B] mt-1">Your community is open to everyone</p>
              <button class="mt-3 text-xs text-[#48BB78] hover:underline" @click="openCreate">+ Create First Rule</button>
            </td>
          </tr>
          <!-- Data rows -->
          <tr
            v-else
            v-for="(item, index) in items"
            :key="item.id"
            class="border-b border-[#1E293B] last:border-b-0 hover:bg-white/[0.02] transition-colors"
            :class="{ 'bg-[#48BB78]/5': dropIndex === index && dragIndex !== index }"
            draggable="true"
            @dragstart="onDragStart(index)"
            @dragover="onDragOver($event, index)"
            @drop="onDrop(index)"
            @dragend="onDragEnd"
          >
            <!-- Drag handle -->
            <td class="px-2 py-3 text-center cursor-grab active:cursor-grabbing">
              <span class="material-symbols-rounded text-base text-[#64748B]">drag_indicator</span>
            </td>
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-[#48BB78]/10 flex items-center justify-center">
                  <span class="material-symbols-rounded text-base text-[#48BB78]">{{ ruleTypeIcons[item.type] }}</span>
                </div>
                <span class="text-sm font-medium text-[#F1F5F9]">{{ item.name }}</span>
              </div>
            </td>
            <td class="px-4 py-3">
              <span class="inline-flex items-center gap-1.5 px-2 py-1 bg-[#1E293B] rounded text-xs text-[#94A3B8]">
                <span class="material-symbols-rounded text-xs">{{ ruleTypeIcons[item.type] }}</span>
                {{ ruleTypeLabels[item.type] }}
              </span>
            </td>
            <td class="px-4 py-3">
              <StatusBadge :status="item.status" />
            </td>
            <td class="px-4 py-3 text-sm text-[#94A3B8]">{{ item.condition_summary }}</td>
            <td class="px-4 py-3 text-center">
              <span class="inline-flex items-center justify-center w-6 h-6 rounded-full bg-[#1E293B] text-xs font-mono text-[#94A3B8]">{{ item.priority }}</span>
            </td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Edit"
                  @click="openEdit(item)"
                >
                  <span class="material-symbols-rounded text-base text-[#64748B]">edit</span>
                </button>
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  :title="item.status === 'active' ? 'Pause' : 'Activate'"
                  @click="toggleStatus(item)"
                >
                  <span class="material-symbols-rounded text-base text-[#64748B]">
                    {{ item.status === 'active' ? 'pause_circle' : 'play_circle' }}
                  </span>
                </button>
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Delete"
                  @click="deleteItem(item)"
                >
                  <span class="material-symbols-rounded text-base text-[#DC2626]">delete</span>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination
        v-if="totalItems > pageSize"
        :page="page"
        :page-size="pageSize"
        :total="totalItems"
        @update:page="onPageChange"
      />
    </div>

    <!-- Create/Edit Access Rule Modal (D10) -->
    <Modal
      :open="showCreate"
      :title="editingId ? 'Edit Access Rule' : 'Create Access Rule'"
      max-width="560px"
      @close="showCreate = false"
    >
      <div class="space-y-4">
        <!-- Rule Name -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Rule Name <span class="text-[#DC2626]">*</span>
            <span class="text-[#64748B] text-xs ml-1">{{ form.name.length }}/50</span>
          </label>
          <input
            v-model="form.name"
            type="text"
            maxlength="50"
            placeholder="e.g. Must hold 100 USDC"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm"
          />
        </div>

        <!-- Rule Type -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-2">Type <span class="text-[#DC2626]">*</span></label>
          <div class="grid grid-cols-2 gap-2">
            <button
              v-for="(label, key) in ruleTypeLabels"
              :key="key"
              class="flex items-center gap-2 px-3 py-2.5 border rounded-lg transition-colors text-sm text-left"
              :class="form.type === key
                ? 'border-[#48BB78] bg-[#48BB78]/10 text-[#48BB78]'
                : 'border-[#1E293B] text-[#94A3B8] hover:border-[#94A3B8]'"
              @click="form.type = key as RuleType"
            >
              <span class="material-symbols-rounded text-base">{{ ruleTypeIcons[key as RuleType] }}</span>
              {{ label }}
            </button>
          </div>
        </div>

        <!-- Token Gate fields -->
        <template v-if="form.type === 'token_gate'">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Chain <span class="text-[#DC2626]">*</span></label>
            <select
              v-model="form.chain"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
            >
              <option value="" disabled>Select chain</option>
              <option v-for="c in chainOptions" :key="c" :value="c">{{ c }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Contract Address <span class="text-[#DC2626]">*</span></label>
            <input
              v-model="form.contract_address"
              type="text"
              placeholder="0x..."
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm font-mono"
            />
          </div>
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Minimum Balance <span class="text-[#DC2626]">*</span></label>
            <input
              v-model.number="form.min_balance"
              type="number"
              min="0"
              step="any"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
            />
          </div>
        </template>

        <!-- NFT Hold fields -->
        <template v-if="form.type === 'nft_hold'">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Chain <span class="text-[#DC2626]">*</span></label>
            <select
              v-model="form.chain"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
            >
              <option value="" disabled>Select chain</option>
              <option v-for="c in chainOptions" :key="c" :value="c">{{ c }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Collection Address <span class="text-[#DC2626]">*</span></label>
            <input
              v-model="form.collection_address"
              type="text"
              placeholder="0x..."
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm font-mono"
            />
          </div>
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Minimum NFT Count <span class="text-[#DC2626]">*</span></label>
            <input
              v-model.number="form.min_count"
              type="number"
              min="1"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
            />
          </div>
        </template>

        <!-- Level Requirement fields -->
        <template v-if="form.type === 'level_requirement'">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Minimum Level <span class="text-[#DC2626]">*</span></label>
            <input
              v-model.number="form.min_level"
              type="number"
              min="1"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
            />
          </div>
        </template>

        <!-- Invite Only fields -->
        <template v-if="form.type === 'invite_only'">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">
              Invite Codes <span class="text-[#DC2626]">*</span>
              <span class="text-[#64748B] text-xs ml-1">One per line</span>
            </label>
            <textarea
              v-model="form.invite_codes"
              rows="5"
              placeholder="CODE001&#10;CODE002&#10;CODE003"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm font-mono resize-none"
            />
          </div>
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Or upload Whitelist CSV</label>
            <div class="flex items-center gap-3">
              <label class="flex-1 flex items-center justify-center gap-2 px-4 py-3 border border-dashed border-[#1E293B] rounded-lg cursor-pointer hover:border-[#94A3B8] transition-colors">
                <span class="material-symbols-rounded text-base text-[#64748B]">upload_file</span>
                <span class="text-sm text-[#64748B]">Choose CSV file</span>
                <input type="file" accept=".csv" class="sr-only" />
              </label>
            </div>
          </div>
        </template>

        <!-- Denial Message -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Denial Message
            <span class="text-[#64748B] text-xs ml-1">{{ form.denial_message.length }}/200</span>
          </label>
          <textarea
            v-model="form.denial_message"
            maxlength="200"
            rows="2"
            placeholder="Message shown to users who don't qualify..."
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm resize-none"
          />
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-[#64748B] hover:text-[#F1F5F9] transition-colors" @click="showCreate = false">Cancel</button>
        <button
          class="px-4 py-2 bg-[#48BB78] text-black text-sm font-medium rounded-lg hover:bg-[#48BB78]/90 disabled:opacity-50 transition-colors"
          :disabled="!formValid || saving"
          @click="saveForm"
        >
          {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Create Rule') }}
        </button>
      </template>
    </Modal>

    <!-- Preview Rule Modal -->
    <Modal
      :open="showPreview"
      title="Preview Rule Evaluation"
      @close="showPreview = false"
    >
      <div class="space-y-4">
        <p class="text-sm text-[#94A3B8]">Test how your current rules evaluate a wallet address.</p>
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">Wallet Address</label>
          <div class="flex gap-2">
            <input
              v-model="previewAddress"
              type="text"
              placeholder="0x..."
              class="flex-1 px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm font-mono"
              @keyup.enter="runPreview"
            />
            <button
              class="px-4 py-2.5 bg-[#48BB78] text-black text-sm font-medium rounded-lg hover:bg-[#48BB78]/90 disabled:opacity-50 transition-colors"
              :disabled="!previewAddress.trim() || previewLoading"
              @click="runPreview"
            >
              {{ previewLoading ? 'Checking...' : 'Evaluate' }}
            </button>
          </div>
        </div>
        <!-- Result -->
        <div v-if="previewResult" class="p-4 rounded-lg border" :class="previewResult.passed ? 'bg-[#0A2E1A]/50 border-[#16A34A]/30' : 'bg-[#2D1515]/50 border-[#DC2626]/30'">
          <div class="flex items-center gap-2 mb-1">
            <span class="material-symbols-rounded text-lg" :class="previewResult.passed ? 'text-[#16A34A]' : 'text-[#DC2626]'">
              {{ previewResult.passed ? 'check_circle' : 'cancel' }}
            </span>
            <span class="text-sm font-medium" :class="previewResult.passed ? 'text-[#16A34A]' : 'text-[#DC2626]'">
              {{ previewResult.passed ? 'Access Granted' : 'Access Denied' }}
            </span>
          </div>
          <p class="text-xs text-[#94A3B8] ml-7">{{ previewResult.reason }}</p>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-[#64748B] hover:text-[#F1F5F9] transition-colors" @click="showPreview = false">Close</button>
      </template>
    </Modal>
  </div>
</template>
