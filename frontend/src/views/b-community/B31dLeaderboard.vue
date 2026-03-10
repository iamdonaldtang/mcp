<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
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

interface LeaderboardItem {
  id: string
  name: string
  status: CampaignStatus
  point_type_id: string
  point_type_name: string
  period: 'weekly' | 'monthly' | 'all_time'
  display_top_n: number
  participants: number
  created_at: string
}

interface StatsData {
  participation_rate: number
  weekly_active_contestants: number
  top3_concentration: number
  avg_position_change: number
}

// --- State ---
const loading = ref(true)
const items = ref<LeaderboardItem[]>([])
const totalItems = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<StatsData>({ participation_rate: 0, weekly_active_contestants: 0, top3_concentration: 0, avg_position_change: 0 })
const pointTypes = ref<PointType[]>([])

// Filters
const filterStatus = ref<string>('all')
const filterPointType = ref<string>('all')
const searchQuery = ref('')
let searchTimeout: ReturnType<typeof setTimeout> | null = null

// Modal
const showCreate = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  point_type_id: '',
  period: 'weekly' as 'weekly' | 'monthly' | 'all_time',
  display_top_n: 100,
})

const formValid = computed(() =>
  form.value.name.trim().length >= 1 &&
  form.value.name.trim().length <= 50 &&
  form.value.point_type_id &&
  form.value.display_top_n >= 10 &&
  form.value.display_top_n <= 1000
)

const statusTabs = ['all', 'active', 'paused', 'draft'] as const

const periodLabels: Record<string, string> = {
  weekly: 'Weekly',
  monthly: 'Monthly',
  all_time: 'All Time',
}

// --- Fetch ---
onMounted(async () => {
  await Promise.all([fetchPointTypes(), fetchStats(), fetchItems()])
  loading.value = false
})

async function fetchPointTypes() {
  try {
    const res = await api.get('/api/v1/community/modules/points/types')
    pointTypes.value = res.data.data || []
  } catch { /* empty */ }
}

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/leaderboard/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

async function fetchItems() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    if (filterPointType.value !== 'all') params.point_type_id = filterPointType.value
    if (searchQuery.value) params.search = searchQuery.value

    const res = await api.get('/api/v1/community/modules/leaderboard', { params })
    items.value = res.data.data?.items || []
    totalItems.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

// Debounced search
function onSearchInput() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchItems()
  }, 300)
}

watch(filterStatus, () => { page.value = 1; fetchItems() })
watch(filterPointType, () => { page.value = 1; fetchItems() })

function onPageChange(p: number) {
  page.value = p
  fetchItems()
}

// --- Actions ---
function openCreate() {
  editingId.value = null
  form.value = { name: '', point_type_id: '', period: 'weekly', display_top_n: 100 }
  showCreate.value = true
}

function openEdit(item: LeaderboardItem) {
  editingId.value = item.id
  form.value = {
    name: item.name,
    point_type_id: item.point_type_id,
    period: item.period,
    display_top_n: item.display_top_n,
  }
  showCreate.value = true
}

async function saveForm() {
  if (!formValid.value || saving.value) return
  saving.value = true
  try {
    if (editingId.value) {
      await api.put(`/api/v1/community/modules/leaderboard/${editingId.value}`, form.value)
    } else {
      await api.post('/api/v1/community/modules/leaderboard', form.value)
    }
    showCreate.value = false
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function archiveItem(item: LeaderboardItem) {
  if (!confirm(`Archive "${item.name}"? Archived leaderboards cannot be edited.`)) return
  try {
    await api.put(`/api/v1/community/modules/leaderboard/${item.id}`, { status: 'completed' })
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

async function duplicateItem(item: LeaderboardItem) {
  try {
    await api.post('/api/v1/community/modules/leaderboard', {
      name: `${item.name} (Copy)`,
      point_type_id: item.point_type_id,
      period: item.period,
      display_top_n: item.display_top_n,
    })
    await fetchItems()
  } catch { /* empty */ }
}

async function deleteItem(id: string) {
  if (!confirm('Delete this leaderboard? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/modules/leaderboard/${id}`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

function isArchived(item: LeaderboardItem) {
  return item.status === 'completed'
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Leaderboard</h1>
        <p class="text-sm text-text-secondary">Manage recurring leaderboards based on custom point types</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-black text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New Leaderboard
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Participation Rate</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.participation_rate }}%</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Weekly Active Contestants</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.weekly_active_contestants.toLocaleString() }}</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Top 3 Concentration</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.top3_concentration }}%</p>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <p class="text-xs text-text-muted mb-1">Avg Position Change</p>
        <p class="text-2xl font-bold text-text-primary">{{ stats.avg_position_change >= 0 ? '+' : '' }}{{ stats.avg_position_change }}</p>
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

      <!-- Point Type Filter -->
      <select
        v-model="filterPointType"
        class="px-3 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary focus:border-community focus:outline-none"
      >
        <option value="all">All Point Types</option>
        <option v-for="pt in pointTypes" :key="pt.id" :value="pt.id">{{ pt.name }} ({{ pt.symbol }})</option>
      </select>

      <!-- Search -->
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search leaderboards..."
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
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Leaderboard Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Period</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Point Type</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Participants</th>
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
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">leaderboard</span>
              <p class="text-sm text-text-muted">No leaderboards found</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First Leaderboard</button>
            </td>
          </tr>
          <!-- Data rows -->
          <tr
            v-else
            v-for="item in items"
            :key="item.id"
            class="border-b border-border last:border-b-0 hover:bg-white/2 transition-colors"
            :class="{ 'opacity-50': isArchived(item) }"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-community/10 flex items-center justify-center">
                  <span class="material-symbols-rounded text-base text-community">leaderboard</span>
                </div>
                <span class="text-sm font-medium text-text-primary">{{ item.name }}</span>
              </div>
            </td>
            <td class="px-4 py-3">
              <StatusBadge :status="item.status" />
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ periodLabels[item.period] || item.period }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ item.point_type_name }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ item.participants.toLocaleString() }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ formatDate(item.created_at) }}</td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button
                  v-if="!isArchived(item)"
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Edit"
                  @click="openEdit(item)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button
                  v-if="!isArchived(item)"
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Archive"
                  @click="archiveItem(item)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">archive</span>
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

    <!-- Create/Edit Leaderboard Modal (D04) -->
    <Modal
      :open="showCreate"
      :title="editingId ? 'Edit Leaderboard' : 'Create Leaderboard'"
      @close="showCreate = false"
    >
      <div class="space-y-4">
        <!-- Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">
            Leaderboard Name
            <span class="text-text-muted text-xs ml-1">{{ form.name.length }}/50</span>
          </label>
          <input
            v-model="form.name"
            type="text"
            maxlength="50"
            placeholder="e.g. Weekly EXP Ranking"
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

        <!-- Period -->
        <div>
          <label class="block text-sm text-text-secondary mb-2">Period</label>
          <div class="flex gap-3">
            <label
              v-for="p in (['weekly', 'monthly', 'all_time'] as const)"
              :key="p"
              class="flex items-center gap-2 px-4 py-2.5 border rounded-lg cursor-pointer transition-colors text-sm"
              :class="form.period === p
                ? 'border-community bg-community/10 text-community'
                : 'border-border text-text-secondary hover:border-text-secondary'"
            >
              <input
                v-model="form.period"
                type="radio"
                :value="p"
                class="sr-only"
              />
              {{ periodLabels[p] }}
            </label>
          </div>
        </div>

        <!-- Display Top N -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Display Top N <span class="text-text-muted text-xs">(10-1000)</span></label>
          <input
            v-model.number="form.display_top_n"
            type="number"
            min="10"
            max="1000"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          />
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary transition-colors" @click="showCreate = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-black text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!formValid || saving"
          @click="saveForm"
        >
          {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Create Leaderboard') }}
        </button>
      </template>
    </Modal>
  </div>
</template>
