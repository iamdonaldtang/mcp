<script setup lang="ts">
import { ref, onMounted, computed, watch, onUnmounted } from 'vue'
import { api } from '../../api/client'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// --- Types ---
interface PointType {
  id: string
  name: string
  symbol: string
}

interface RewardTier {
  rank_from: number
  rank_to: number
  reward_type: 'token' | 'nft' | 'wl_spot'
  reward_value: string
}

interface SprintItem {
  id: string
  name: string
  status: CampaignStatus
  point_type_id: string
  point_type_name: string
  reward_summary: string
  start_date: string
  end_date: string
  participants: number
  reward_tiers: RewardTier[]
  distribution_mode: 'auto' | 'manual'
  created_at: string
}

interface StatsData {
  total_sprints: number
  active_participants: number
  tasks_completed: number
  rewards_given: number
}

// --- State ---
const loading = ref(true)
const items = ref<SprintItem[]>([])
const totalItems = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<StatsData>({ total_sprints: 0, active_participants: 0, tasks_completed: 0, rewards_given: 0 })
const pointTypes = ref<PointType[]>([])

// Filters
const filterStatus = ref<string>('all')
const searchQuery = ref('')
let searchTimeout: ReturnType<typeof setTimeout> | null = null

// Countdown timers
const now = ref(Date.now())
let countdownInterval: ReturnType<typeof setInterval> | null = null

// Modal
const showCreate = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  point_type_id: '',
  start_date: '',
  end_date: '',
  distribution_mode: 'auto' as 'auto' | 'manual',
  reward_tiers: [{ rank_from: 1, rank_to: 3, reward_type: 'token', reward_value: '' }] as RewardTier[],
})

const rewardTypeOptions = [
  { value: 'token', label: 'Token' },
  { value: 'nft', label: 'NFT' },
  { value: 'wl_spot', label: 'WL Spot' },
]

const statusTabs = ['all', 'active', 'scheduled', 'completed', 'draft'] as const

// Map scheduled to a CampaignStatus the badge can handle
const statusMap: Record<string, CampaignStatus> = {
  active: 'active',
  scheduled: 'draft',
  completed: 'completed',
  draft: 'draft',
}

// --- Validation ---
const todayStr = computed(() => {
  const d = new Date()
  return d.toISOString().split('T')[0]
})

const hasOverlappingRanks = computed(() => {
  const tiers = form.value.reward_tiers
  for (let i = 0; i < tiers.length; i++) {
    for (let j = i + 1; j < tiers.length; j++) {
      if (tiers[i].rank_from <= tiers[j].rank_to && tiers[j].rank_from <= tiers[i].rank_to) {
        return true
      }
    }
  }
  return false
})

const durationValid = computed(() => {
  if (!form.value.start_date || !form.value.end_date) return false
  const start = new Date(form.value.start_date)
  const end = new Date(form.value.end_date)
  const diffMs = end.getTime() - start.getTime()
  const diffDays = diffMs / (1000 * 60 * 60 * 24)
  return diffMs >= 3600000 && diffDays <= 90 // min 1 hour, max 90 days
})

const formValid = computed(() =>
  form.value.name.trim().length >= 1 &&
  form.value.name.trim().length <= 50 &&
  form.value.point_type_id &&
  form.value.start_date &&
  form.value.end_date &&
  durationValid.value &&
  form.value.reward_tiers.length >= 1 &&
  form.value.reward_tiers.every(t => t.rank_from > 0 && t.rank_to >= t.rank_from && t.reward_value.trim()) &&
  !hasOverlappingRanks.value
)

// --- Fetch ---
onMounted(async () => {
  await Promise.all([fetchPointTypes(), fetchStats(), fetchItems()])
  loading.value = false
  countdownInterval = setInterval(() => { now.value = Date.now() }, 1000)
})

onUnmounted(() => {
  if (countdownInterval) clearInterval(countdownInterval)
})

async function fetchPointTypes() {
  try {
    const res = await api.get('/api/v1/community/modules/points/types')
    pointTypes.value = res.data.data || []
  } catch { /* empty */ }
}

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/lb-sprint/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

async function fetchItems() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    if (searchQuery.value) params.search = searchQuery.value

    const res = await api.get('/api/v1/community/modules/lb-sprint', { params })
    items.value = res.data.data?.items || []
    totalItems.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

function onSearchInput() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchItems()
  }, 300)
}

watch(filterStatus, () => { page.value = 1; fetchItems() })

function onPageChange(p: number) {
  page.value = p
  fetchItems()
}

// --- Countdown ---
function getCountdown(startDate: string): string {
  const diff = new Date(startDate).getTime() - now.value
  if (diff <= 0) return 'Starting...'
  const days = Math.floor(diff / 86400000)
  const hours = Math.floor((diff % 86400000) / 3600000)
  const mins = Math.floor((diff % 3600000) / 60000)
  const secs = Math.floor((diff % 60000) / 1000)
  if (days > 0) return `${days}d ${hours}h ${mins}m`
  return `${hours}h ${mins}m ${secs}s`
}

// --- Actions ---
function openCreate() {
  editingId.value = null
  form.value = {
    name: '',
    point_type_id: '',
    start_date: '',
    end_date: '',
    distribution_mode: 'auto',
    reward_tiers: [{ rank_from: 1, rank_to: 3, reward_type: 'token', reward_value: '' }],
  }
  showCreate.value = true
}

function openEdit(item: SprintItem) {
  editingId.value = item.id
  form.value = {
    name: item.name,
    point_type_id: item.point_type_id,
    start_date: item.start_date.split('T')[0],
    end_date: item.end_date.split('T')[0],
    distribution_mode: item.distribution_mode,
    reward_tiers: item.reward_tiers.length > 0
      ? item.reward_tiers.map(t => ({ ...t }))
      : [{ rank_from: 1, rank_to: 3, reward_type: 'token', reward_value: '' }],
  }
  showCreate.value = true
}

function addTier() {
  const tiers = form.value.reward_tiers
  const lastTo = tiers.length > 0 ? tiers[tiers.length - 1].rank_to : 0
  tiers.push({ rank_from: lastTo + 1, rank_to: lastTo + 5, reward_type: 'token', reward_value: '' })
}

function removeTier(index: number) {
  form.value.reward_tiers.splice(index, 1)
}

async function saveForm() {
  if (!formValid.value || saving.value) return
  saving.value = true
  try {
    if (editingId.value) {
      await api.put(`/api/v1/community/modules/lb-sprint/${editingId.value}`, form.value)
    } else {
      await api.post('/api/v1/community/modules/lb-sprint', form.value)
    }
    showCreate.value = false
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function endEarly(item: SprintItem) {
  if (!confirm(`End "${item.name}" early? This will finalize results and distribute rewards.`)) return
  try {
    await api.put(`/api/v1/community/modules/lb-sprint/${item.id}/end`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

async function duplicateItem(item: SprintItem) {
  try {
    await api.post('/api/v1/community/modules/lb-sprint', {
      name: `${item.name} (Copy)`,
      point_type_id: item.point_type_id,
      start_date: item.start_date,
      end_date: item.end_date,
      distribution_mode: item.distribution_mode,
      reward_tiers: item.reward_tiers,
    })
    await fetchItems()
  } catch { /* empty */ }
}

async function deleteItem(id: string) {
  if (!confirm('Delete this LB Sprint? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/modules/lb-sprint/${id}`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

function formatDateRange(start: string, end: string) {
  const s = new Date(start).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  const e = new Date(end).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  return `${s} - ${e}`
}

// Status badge color for scheduled (which isn't in the standard map)
function getStatusBadge(status: string): { label: string; bgColor: string; textColor: string } {
  const map: Record<string, { label: string; bgColor: string; textColor: string }> = {
    active: { label: 'Active', bgColor: '#0A2E1A', textColor: '#16A34A' },
    scheduled: { label: 'Scheduled', bgColor: '#0C1A2E', textColor: '#3B82F6' },
    completed: { label: 'Completed', bgColor: '#1E293B', textColor: '#64748B' },
    draft: { label: 'Draft', bgColor: '#1F1A08', textColor: '#D97706' },
  }
  return map[status] || map.draft
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">LB Sprint</h1>
        <p class="text-sm text-text-secondary">Manage time-bounded competitions with reward incentives</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-black text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New LB Sprint
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Total LB Sprints</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.total_sprints }}</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Active Participants</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.active_participants.toLocaleString() }}</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Tasks Completed</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.tasks_completed.toLocaleString() }}</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">NFTs + Rewards Given</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.rewards_given.toLocaleString() }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <!-- Status Tabs -->
      <div class="flex border-b border-border">
        <button
          v-for="tab in statusTabs"
          :key="tab"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="filterStatus === tab
            ? 'text-community'
            : 'text-text-secondary hover:text-text-primary'"
          @click="filterStatus = tab"
        >
          {{ tab === 'all' ? 'All' : tab.charAt(0).toUpperCase() + tab.slice(1) }}
          <span
            v-if="filterStatus === tab"
            class="absolute bottom-0 left-0 right-0 h-0.5 bg-community rounded-t"
          />
        </button>
      </div>

      <!-- Search -->
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search sprints..."
            class="w-full pl-10 pr-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-community focus:outline-none"
            @input="onSearchInput"
          />
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-border">
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Sprint Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-40">Rewards</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-44">Duration</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Participants</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Actions</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading skeleton -->
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="6" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse" /></td>
          </tr>
          <!-- Empty state -->
          <tr v-else-if="items.length === 0">
            <td colspan="6" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">emoji_events</span>
              <p class="text-sm text-text-muted">No LB Sprints found</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First LB Sprint</button>
            </td>
          </tr>
          <!-- Data rows -->
          <tr
            v-else
            v-for="item in items"
            :key="item.id"
            class="border-b border-border last:border-b-0 hover:bg-white/2 transition-colors"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-community/10 flex items-center justify-center">
                  <span class="material-symbols-rounded text-base text-community">emoji_events</span>
                </div>
                <div>
                  <span class="text-sm font-medium text-text-primary">{{ item.name }}</span>
                  <div class="text-xs text-text-muted">{{ item.point_type_name }}</div>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :style="{ backgroundColor: getStatusBadge(item.status).bgColor, color: getStatusBadge(item.status).textColor }"
              >
                {{ getStatusBadge(item.status).label }}
              </span>
              <!-- Countdown for scheduled -->
              <div v-if="item.status === ('scheduled' as CampaignStatus)" class="text-xs text-text-muted mt-1">
                Starts in {{ getCountdown(item.start_date) }}
              </div>
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ item.reward_summary }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ formatDateRange(item.start_date, item.end_date) }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ item.participants.toLocaleString() }}</td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button
                  v-if="item.status !== 'completed'"
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Edit"
                  @click="openEdit(item)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button
                  v-if="item.status === 'active'"
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="End Early"
                  @click="endEarly(item)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">stop_circle</span>
                </button>
                <button
                  v-if="item.status === 'completed'"
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="View Results"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">bar_chart</span>
                </button>
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Duplicate"
                  @click="duplicateItem(item)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">content_copy</span>
                </button>
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Delete"
                  @click="deleteItem(item.id)"
                >
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
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

    <!-- Create/Edit LB Sprint Modal (D05) -->
    <Modal
      :open="showCreate"
      :title="editingId ? 'Edit LB Sprint' : 'Create LB Sprint'"
      max-width="640px"
      @close="showCreate = false"
    >
      <div class="space-y-4">
        <!-- Sprint Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">
            Sprint Name
            <span class="text-text-muted text-xs ml-1">{{ form.name.length }}/50</span>
          </label>
          <input
            v-model="form.name"
            type="text"
            maxlength="50"
            placeholder="e.g. Weekly EXP Championship"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
        </div>

        <!-- Point Type -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Point Type <span class="text-status-paused">*</span></label>
          <select
            v-model="form.point_type_id"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          >
            <option value="" disabled>Select point type</option>
            <option v-for="pt in pointTypes" :key="pt.id" :value="pt.id">{{ pt.name }} ({{ pt.symbol }})</option>
          </select>
        </div>

        <!-- Date Range -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-text-secondary mb-1">Start Date</label>
            <input
              v-model="form.start_date"
              type="date"
              :min="todayStr"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1">End Date</label>
            <input
              v-model="form.end_date"
              type="date"
              :min="form.start_date || todayStr"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            />
          </div>
        </div>
        <p v-if="form.start_date && form.end_date && !durationValid" class="text-xs text-status-paused">
          Duration must be at least 1 hour and no more than 90 days.
        </p>

        <!-- Reward Tiers -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="block text-sm text-text-secondary">Reward Tiers <span class="text-status-paused">*</span></label>
            <button
              class="text-xs text-community hover:underline"
              @click="addTier"
            >
              + Add Tier
            </button>
          </div>
          <div class="space-y-3">
            <div
              v-for="(tier, index) in form.reward_tiers"
              :key="index"
              class="flex items-start gap-2 p-3 bg-page-bg border border-border rounded-lg"
            >
              <div class="flex-1 grid grid-cols-4 gap-2">
                <!-- Rank From -->
                <div>
                  <label class="block text-xs text-text-muted mb-1">Rank From</label>
                  <input
                    v-model.number="tier.rank_from"
                    type="number"
                    min="1"
                    class="w-full px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
                  />
                </div>
                <!-- Rank To -->
                <div>
                  <label class="block text-xs text-text-muted mb-1">Rank To</label>
                  <input
                    v-model.number="tier.rank_to"
                    type="number"
                    :min="tier.rank_from"
                    class="w-full px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
                  />
                </div>
                <!-- Reward Type -->
                <div>
                  <label class="block text-xs text-text-muted mb-1">Reward Type</label>
                  <select
                    v-model="tier.reward_type"
                    class="w-full px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
                  >
                    <option v-for="opt in rewardTypeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
                  </select>
                </div>
                <!-- Reward Value -->
                <div>
                  <label class="block text-xs text-text-muted mb-1">Value</label>
                  <input
                    v-model="tier.reward_value"
                    type="text"
                    placeholder="e.g. 100 USDT"
                    class="w-full px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-community focus:outline-none"
                  />
                </div>
              </div>
              <button
                v-if="form.reward_tiers.length > 1"
                class="mt-5 p-1 rounded hover:bg-white/5 transition-colors"
                @click="removeTier(index)"
              >
                <span class="material-symbols-rounded text-base text-status-paused">close</span>
              </button>
            </div>
          </div>
          <p v-if="hasOverlappingRanks" class="text-xs text-status-paused mt-1">Rank ranges cannot overlap.</p>
        </div>

        <!-- Distribution Mode -->
        <div>
          <label class="block text-sm text-text-secondary mb-2">Distribution Mode</label>
          <div class="flex gap-3">
            <label
              v-for="mode in (['auto', 'manual'] as const)"
              :key="mode"
              class="flex items-center gap-2 px-4 py-2.5 border rounded-lg cursor-pointer transition-colors text-sm"
              :class="form.distribution_mode === mode
                ? 'border-community bg-community/10 text-community'
                : 'border-border text-text-secondary hover:border-text-secondary'"
            >
              <input
                v-model="form.distribution_mode"
                type="radio"
                :value="mode"
                class="sr-only"
              />
              {{ mode.charAt(0).toUpperCase() + mode.slice(1) }}
            </label>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary transition-colors" @click="showCreate = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-black text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!formValid || saving"
          @click="saveForm"
        >
          {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Create LB Sprint') }}
        </button>
      </template>
    </Modal>
  </div>
</template>
