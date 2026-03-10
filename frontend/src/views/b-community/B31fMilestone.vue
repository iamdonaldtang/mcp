<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { api } from '../../api/client'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// --- Types ---
interface MilestoneTier {
  threshold: number
  reward_type: 'badge' | 'shop_item' | 'custom'
  reward_value: string
  completion_rate?: number
}

interface MilestoneItem {
  id: string
  name: string
  status: CampaignStatus
  tiers: MilestoneTier[]
  completions: number
  created_at: string
}

interface StatsData {
  total_sets: number
  completions: number
  rewards_claimed: number
  claim_rate: number
}

// --- State ---
const loading = ref(true)
const items = ref<MilestoneItem[]>([])
const totalItems = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<StatsData>({ total_sets: 0, completions: 0, rewards_claimed: 0, claim_rate: 0 })

// Filters
const filterStatus = ref<string>('all')
const searchQuery = ref('')
let searchTimeout: ReturnType<typeof setTimeout> | null = null

// Expandable rows
const expandedIds = ref<Set<string>>(new Set())

// Modal
const showCreate = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  tiers: [{ threshold: 100, reward_type: 'badge', reward_value: '' }] as MilestoneTier[],
})

const rewardTypeOptions = [
  { value: 'badge', label: 'Badge' },
  { value: 'shop_item', label: 'Shop Item' },
  { value: 'custom', label: 'Custom' },
]

const statusTabs = ['all', 'active', 'paused', 'draft'] as const

// --- Validation ---
const tiersAscending = computed(() => {
  const tiers = form.value.tiers
  for (let i = 1; i < tiers.length; i++) {
    if (tiers[i].threshold <= tiers[i - 1].threshold) return false
  }
  return true
})

const formValid = computed(() =>
  form.value.name.trim().length >= 1 &&
  form.value.name.trim().length <= 50 &&
  form.value.tiers.length >= 1 &&
  form.value.tiers.length <= 10 &&
  form.value.tiers.every(t => t.threshold > 0 && t.reward_value.trim()) &&
  tiersAscending.value
)

// --- Fetch ---
onMounted(async () => {
  await Promise.all([fetchStats(), fetchItems()])
  loading.value = false
})

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/milestone/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

async function fetchItems() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    if (searchQuery.value) params.search = searchQuery.value

    const res = await api.get('/api/v1/community/modules/milestone', { params })
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

// --- Expandable rows ---
function toggleExpand(id: string) {
  if (expandedIds.value.has(id)) {
    expandedIds.value.delete(id)
  } else {
    expandedIds.value.add(id)
  }
}

// --- Actions ---
function openCreate() {
  editingId.value = null
  form.value = {
    name: '',
    tiers: [{ threshold: 100, reward_type: 'badge', reward_value: '' }],
  }
  showCreate.value = true
}

function openEdit(item: MilestoneItem) {
  editingId.value = item.id
  form.value = {
    name: item.name,
    tiers: item.tiers.map(t => ({
      threshold: t.threshold,
      reward_type: t.reward_type,
      reward_value: t.reward_value,
    })),
  }
  showCreate.value = true
}

function addTier() {
  if (form.value.tiers.length >= 10) return
  const lastThreshold = form.value.tiers.length > 0
    ? form.value.tiers[form.value.tiers.length - 1].threshold
    : 0
  form.value.tiers.push({ threshold: lastThreshold + 100, reward_type: 'badge', reward_value: '' })
}

function removeTier(index: number) {
  form.value.tiers.splice(index, 1)
}

// Auto-sort tiers by threshold on blur
function sortTiers() {
  form.value.tiers.sort((a, b) => a.threshold - b.threshold)
}

async function saveForm() {
  if (!formValid.value || saving.value) return
  saving.value = true
  try {
    if (editingId.value) {
      await api.put(`/api/v1/community/modules/milestone/${editingId.value}`, form.value)
    } else {
      await api.post('/api/v1/community/modules/milestone', form.value)
    }
    showCreate.value = false
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function toggleStatus(item: MilestoneItem) {
  const newStatus = item.status === 'active' ? 'paused' : 'active'
  const oldStatus = item.status
  item.status = newStatus as CampaignStatus
  try {
    await api.put(`/api/v1/community/modules/milestone/${item.id}`, { status: newStatus })
    await fetchStats()
  } catch {
    item.status = oldStatus
  }
}

async function duplicateItem(item: MilestoneItem) {
  try {
    await api.post('/api/v1/community/modules/milestone', {
      name: `${item.name} (Copy)`,
      tiers: item.tiers.map(t => ({ threshold: t.threshold, reward_type: t.reward_type, reward_value: t.reward_value })),
    })
    await fetchItems()
  } catch { /* empty */ }
}

async function deleteItem(id: string) {
  if (!confirm('Delete this milestone set? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/modules/milestone/${id}`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function rewardTypeLabel(type: string) {
  return rewardTypeOptions.find(o => o.value === type)?.label || type
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Milestones</h1>
        <p class="text-sm text-text-secondary">Manage milestone sets with tiered rewards for user progression</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-black text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New Milestone Set
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Total Sets</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.total_sets }}</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Completions</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.completions.toLocaleString() }}</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Rewards Claimed</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.rewards_claimed.toLocaleString() }}</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Claim Rate</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.claim_rate }}%</p>
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
            placeholder="Search milestones..."
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
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-10"></th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Milestone Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-20">Tiers</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Completions</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Created</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Actions</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading skeleton -->
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="7" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse" /></td>
          </tr>
          <!-- Empty state -->
          <tr v-else-if="items.length === 0">
            <td colspan="7" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">flag</span>
              <p class="text-sm text-text-muted">No milestone sets found</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First Milestone Set</button>
            </td>
          </tr>
          <!-- Data rows -->
          <template v-else v-for="item in items" :key="item.id">
            <tr
              class="border-b border-border hover:bg-white/2 transition-colors cursor-pointer"
              @click="toggleExpand(item.id)"
            >
              <td class="px-4 py-3">
                <span class="material-symbols-rounded text-base text-text-muted transition-transform" :class="{ 'rotate-90': expandedIds.has(item.id) }">
                  chevron_right
                </span>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-lg bg-community/10 flex items-center justify-center">
                    <span class="material-symbols-rounded text-base text-community">flag</span>
                  </div>
                  <span class="text-sm font-medium text-text-primary">{{ item.name }}</span>
                </div>
              </td>
              <td class="px-4 py-3">
                <StatusBadge :status="item.status" />
              </td>
              <td class="px-4 py-3 text-sm text-text-secondary">{{ item.tiers.length }}</td>
              <td class="px-4 py-3 text-sm text-text-secondary">{{ item.completions.toLocaleString() }}</td>
              <td class="px-4 py-3 text-sm text-text-secondary">{{ formatDate(item.created_at) }}</td>
              <td class="px-4 py-3 text-right" @click.stop>
                <div class="flex items-center justify-end gap-1">
                  <button
                    class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                    title="Edit"
                    @click="openEdit(item)"
                  >
                    <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                  </button>
                  <button
                    class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                    :title="item.status === 'active' ? 'Pause' : 'Activate'"
                    @click="toggleStatus(item)"
                  >
                    <span class="material-symbols-rounded text-base text-text-muted">
                      {{ item.status === 'active' ? 'pause_circle' : 'play_circle' }}
                    </span>
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
            <!-- Expanded tier details -->
            <tr v-if="expandedIds.has(item.id)">
              <td colspan="7" class="px-4 py-0">
                <div class="py-3 pl-14 pr-4">
                  <div class="grid grid-cols-4 gap-2 text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 px-3">
                    <span>Tier</span>
                    <span>Threshold</span>
                    <span>Reward</span>
                    <span>Completion Rate</span>
                  </div>
                  <div
                    v-for="(tier, idx) in item.tiers"
                    :key="idx"
                    class="grid grid-cols-4 gap-2 items-center px-3 py-2 rounded-lg odd:bg-page-bg text-sm"
                  >
                    <span class="text-text-secondary font-medium">Tier {{ idx + 1 }}</span>
                    <span class="text-text-primary font-medium">{{ tier.threshold.toLocaleString() }}</span>
                    <span class="text-text-secondary">
                      <span class="px-2 py-0.5 text-xs rounded bg-page-bg text-text-muted border border-border">{{ rewardTypeLabel(tier.reward_type) }}</span>
                      {{ tier.reward_value }}
                    </span>
                    <div class="flex items-center gap-2">
                      <div class="flex-1 h-1.5 bg-border rounded-full overflow-hidden">
                        <div class="h-full bg-community rounded-full transition-all" :style="{ width: (tier.completion_rate || 0) + '%' }" />
                      </div>
                      <span class="text-xs text-text-muted w-10 text-right">{{ tier.completion_rate || 0 }}%</span>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </template>
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

    <!-- Create/Edit Milestone Modal (D06) -->
    <Modal
      :open="showCreate"
      :title="editingId ? 'Edit Milestone Set' : 'Create Milestone Set'"
      max-width="600px"
      @close="showCreate = false"
    >
      <div class="space-y-4">
        <!-- Milestone Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">
            Milestone Name
            <span class="text-text-muted text-xs ml-1">{{ form.name.length }}/50</span>
          </label>
          <input
            v-model="form.name"
            type="text"
            maxlength="50"
            placeholder="e.g. Community Champion"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
        </div>

        <!-- Tiers -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="block text-sm text-text-secondary">
              Tiers <span class="text-status-paused">*</span>
              <span class="text-text-muted text-xs ml-1">{{ form.tiers.length }}/10</span>
            </label>
            <button
              v-if="form.tiers.length < 10"
              class="text-xs text-community hover:underline"
              @click="addTier"
            >
              + Add Tier
            </button>
          </div>
          <div class="space-y-3">
            <div
              v-for="(tier, index) in form.tiers"
              :key="index"
              class="flex items-start gap-2 p-3 bg-page-bg border border-border rounded-lg"
            >
              <!-- Tier number -->
              <div class="w-7 h-7 rounded-full bg-community/10 flex items-center justify-center mt-5 shrink-0">
                <span class="text-xs font-bold text-community">{{ index + 1 }}</span>
              </div>
              <div class="flex-1 grid grid-cols-3 gap-2">
                <!-- Threshold -->
                <div>
                  <label class="block text-xs text-text-muted mb-1">Threshold</label>
                  <input
                    v-model.number="tier.threshold"
                    type="number"
                    min="1"
                    class="w-full px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
                    @blur="sortTiers"
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
                  <label class="block text-xs text-text-muted mb-1">Reward Value</label>
                  <input
                    v-model="tier.reward_value"
                    type="text"
                    placeholder="e.g. Gold Badge"
                    class="w-full px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-community focus:outline-none"
                  />
                </div>
              </div>
              <button
                v-if="form.tiers.length > 1"
                class="mt-5 p-1 rounded hover:bg-white/5 transition-colors shrink-0"
                @click="removeTier(index)"
              >
                <span class="material-symbols-rounded text-base text-status-paused">close</span>
              </button>
            </div>
          </div>
          <p v-if="!tiersAscending && form.tiers.length > 1" class="text-xs text-status-paused mt-1">
            Thresholds will be auto-sorted in ascending order on blur.
          </p>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary transition-colors" @click="showCreate = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-black text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!formValid || saving"
          @click="saveForm"
        >
          {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Create Milestone Set') }}
        </button>
      </template>
    </Modal>
  </div>
</template>
