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
  color: string
  is_default: boolean
}

interface LevelItem {
  id: string
  name: string
  status: CampaignStatus
  threshold: number
  point_type_id: string
  point_type_name: string
  badge_icon: string
  perks: string
  members_at_level: number
  created_at: string
}

interface StatsData {
  total_points_issued: number
  active_earners: number
  avg_points_per_user: number
  level_up_events: number
}

// --- State ---
const loading = ref(true)
const items = ref<LevelItem[]>([])
const totalItems = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<StatsData>({ total_points_issued: 0, active_earners: 0, avg_points_per_user: 0, level_up_events: 0 })
const pointTypes = ref<PointType[]>([])

// Filters
const filterStatus = ref<string>('all')
const statusTabs = ['all', 'active', 'draft'] as const

// Point Type Modal
const showPointTypeModal = ref(false)
const savingPointType = ref(false)
const editingPointTypeId = ref<string | null>(null)
const pointTypeForm = ref({
  name: '',
  symbol: '',
  color: '#48BB78',
})
const pointTypeFormValid = computed(() =>
  pointTypeForm.value.name.trim().length >= 1 &&
  pointTypeForm.value.name.trim().length <= 30 &&
  pointTypeForm.value.symbol.trim().length >= 1 &&
  pointTypeForm.value.symbol.trim().length <= 6
)

// Level Modal
const showCreate = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  threshold: 0,
  point_type_id: '',
  badge_icon: 'military_tech',
  perks: '',
})

const badgeIcons = [
  'military_tech', 'star', 'diamond', 'workspace_premium', 'emoji_events',
  'shield', 'verified', 'local_fire_department', 'bolt', 'rocket_launch',
  'auto_awesome', 'grade', 'favorite', 'crown',
] as const

const formValid = computed(() =>
  form.value.name.trim().length >= 1 &&
  form.value.name.trim().length <= 30 &&
  form.value.point_type_id !== '' &&
  form.value.threshold >= 0 &&
  form.value.badge_icon !== ''
)

// Reassignment dialog
const showReassign = ref(false)
const deletingLevel = ref<LevelItem | null>(null)
const reassignTargetId = ref<string>('')

// Predefined colors for point types
const presetColors = ['#48BB78', '#5D7EF1', '#ED8936', '#9B7EE0', '#F59E0B', '#DC2626', '#06B6D4', '#EC4899']

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
    const res = await api.get('/api/v1/community/modules/points/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

async function fetchItems() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    const res = await api.get('/api/v1/community/modules/points/levels', { params })
    items.value = res.data.data?.items || []
    totalItems.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

watch(filterStatus, () => { page.value = 1; fetchItems() })

function onPageChange(p: number) {
  page.value = p
  fetchItems()
}

// --- Point Type Actions ---
function openCreatePointType() {
  editingPointTypeId.value = null
  pointTypeForm.value = { name: '', symbol: '', color: '#48BB78' }
  showPointTypeModal.value = true
}

function openEditPointType(pt: PointType) {
  editingPointTypeId.value = pt.id
  pointTypeForm.value = { name: pt.name, symbol: pt.symbol, color: pt.color }
  showPointTypeModal.value = true
}

async function savePointType() {
  if (!pointTypeFormValid.value || savingPointType.value) return
  savingPointType.value = true
  try {
    if (editingPointTypeId.value) {
      await api.put(`/api/v1/community/modules/points/types/${editingPointTypeId.value}`, pointTypeForm.value)
    } else {
      await api.post('/api/v1/community/modules/points/types', pointTypeForm.value)
    }
    showPointTypeModal.value = false
    await fetchPointTypes()
  } catch { /* TODO: toast */ }
  savingPointType.value = false
}

async function deletePointType(pt: PointType) {
  if (pt.is_default) return
  if (!confirm(`Delete point type "${pt.name}"? All associated levels and points will be affected.`)) return
  try {
    await api.delete(`/api/v1/community/modules/points/types/${pt.id}`)
    await fetchPointTypes()
  } catch { /* empty */ }
}

// --- Level Actions ---
function openCreate() {
  editingId.value = null
  form.value = { name: '', threshold: 0, point_type_id: pointTypes.value[0]?.id || '', badge_icon: 'military_tech', perks: '' }
  showCreate.value = true
}

function openEdit(item: LevelItem) {
  editingId.value = item.id
  form.value = {
    name: item.name,
    threshold: item.threshold,
    point_type_id: item.point_type_id,
    badge_icon: item.badge_icon || 'military_tech',
    perks: item.perks || '',
  }
  showCreate.value = true
}

async function saveForm() {
  if (!formValid.value || saving.value) return
  saving.value = true
  try {
    if (editingId.value) {
      await api.put(`/api/v1/community/modules/points/levels/${editingId.value}`, form.value)
    } else {
      await api.post('/api/v1/community/modules/points/levels', form.value)
    }
    showCreate.value = false
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

function confirmDelete(item: LevelItem) {
  if (item.members_at_level > 0) {
    deletingLevel.value = item
    reassignTargetId.value = ''
    showReassign.value = true
  } else {
    directDelete(item.id)
  }
}

async function directDelete(id: string) {
  if (!confirm('Delete this level? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/modules/points/levels/${id}`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

async function deleteWithReassign() {
  if (!deletingLevel.value || !reassignTargetId.value) return
  try {
    await api.delete(`/api/v1/community/modules/points/levels/${deletingLevel.value.id}`, {
      data: { reassign_to: reassignTargetId.value }
    })
    showReassign.value = false
    deletingLevel.value = null
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

const reassignOptions = computed(() =>
  items.value.filter(i => i.id !== deletingLevel.value?.id)
)

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
        <h1 class="text-2xl font-bold text-[#F1F5F9] mb-1">Points & Level</h1>
        <p class="text-sm text-[#94A3B8]">Manage point types, issuance, and level progression</p>
      </div>
      <button
        class="px-4 py-2 bg-[#48BB78] text-black text-sm font-medium rounded-lg hover:bg-[#48BB78]/90 transition-colors"
        @click="openCreate"
      >
        + Create Level
      </button>
    </div>

    <!-- Point Types Section -->
    <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-5">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2">
          <span class="material-symbols-rounded text-lg text-[#48BB78]">token</span>
          <h2 class="text-base font-semibold text-[#F1F5F9]">Point Types</h2>
          <span class="text-xs text-[#64748B]">({{ pointTypes.length }}/10)</span>
        </div>
        <button
          v-if="pointTypes.length < 10"
          class="px-3 py-1.5 text-xs font-medium text-[#48BB78] border border-[#48BB78]/30 rounded-lg hover:bg-[#48BB78]/10 transition-colors"
          @click="openCreatePointType"
        >
          + Add Point Type
        </button>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
        <div
          v-for="pt in pointTypes"
          :key="pt.id"
          class="flex items-center gap-3 p-3 bg-[#0A0F1A] border border-[#1E293B] rounded-lg group hover:border-[#48BB78]/30 transition-colors"
        >
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold text-white shrink-0"
            :style="{ backgroundColor: pt.color }"
          >
            {{ pt.symbol }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-[#F1F5F9] truncate">{{ pt.name }}</p>
            <p class="text-xs text-[#64748B]">{{ pt.symbol }}</p>
          </div>
          <div class="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              class="p-1 rounded hover:bg-white/5 transition-colors"
              title="Edit"
              @click="openEditPointType(pt)"
            >
              <span class="material-symbols-rounded text-sm text-[#64748B]">edit</span>
            </button>
            <button
              v-if="!pt.is_default"
              class="p-1 rounded hover:bg-white/5 transition-colors"
              title="Delete"
              @click="deletePointType(pt)"
            >
              <span class="material-symbols-rounded text-sm text-[#DC2626]">delete</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Total Points Issued</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.total_points_issued) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Active Earners</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.active_earners) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Avg Points / User</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.avg_points_per_user) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Level-Up Events</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.level_up_events) }}</p>
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
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider">Level Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-28">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-32">Threshold</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-36">Members at Level</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-32">Created Date</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-[#64748B] uppercase tracking-wider w-24">Actions</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading skeleton -->
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="6" class="px-4 py-4"><div class="h-4 bg-[#1E293B] rounded animate-pulse" /></td>
          </tr>
          <!-- Empty state -->
          <tr v-else-if="items.length === 0">
            <td colspan="6" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-[#64748B] block mb-2">military_tech</span>
              <p class="text-sm text-[#64748B]">No levels configured yet</p>
              <button class="mt-2 text-xs text-[#48BB78] hover:underline" @click="openCreate">+ Create First Level</button>
            </td>
          </tr>
          <!-- Data rows -->
          <tr
            v-else
            v-for="item in items"
            :key="item.id"
            class="border-b border-[#1E293B] last:border-b-0 hover:bg-white/[0.02] transition-colors"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-[#48BB78]/10 flex items-center justify-center">
                  <span class="material-symbols-rounded text-base text-[#48BB78]">{{ item.badge_icon || 'military_tech' }}</span>
                </div>
                <div>
                  <span class="text-sm font-medium text-[#F1F5F9]">{{ item.name }}</span>
                  <p v-if="item.perks" class="text-xs text-[#64748B] mt-0.5 truncate max-w-[240px]">{{ item.perks }}</p>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <StatusBadge :status="item.status" />
            </td>
            <td class="px-4 py-3">
              <span class="text-sm text-[#94A3B8]">{{ formatNumber(item.threshold) }}</span>
              <span class="text-xs text-[#64748B] ml-1">{{ item.point_type_name }}</span>
            </td>
            <td class="px-4 py-3 text-sm text-[#94A3B8]">{{ formatNumber(item.members_at_level) }}</td>
            <td class="px-4 py-3 text-sm text-[#94A3B8]">{{ formatDate(item.created_at) }}</td>
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
                  title="Delete"
                  @click="confirmDelete(item)"
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

    <!-- Create/Edit Point Type Modal -->
    <Modal
      :open="showPointTypeModal"
      :title="editingPointTypeId ? 'Edit Point Type' : 'Add Point Type'"
      @close="showPointTypeModal = false"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Name
            <span class="text-[#64748B] text-xs ml-1">{{ pointTypeForm.name.length }}/30</span>
          </label>
          <input
            v-model="pointTypeForm.name"
            type="text"
            maxlength="30"
            placeholder="e.g. Experience Points"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm"
          />
        </div>
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Symbol
            <span class="text-[#64748B] text-xs ml-1">{{ pointTypeForm.symbol.length }}/6</span>
          </label>
          <input
            v-model="pointTypeForm.symbol"
            type="text"
            maxlength="6"
            placeholder="e.g. XP, GEM"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm uppercase"
          />
        </div>
        <div>
          <label class="block text-sm text-[#94A3B8] mb-2">Color</label>
          <div class="flex items-center gap-2">
            <button
              v-for="c in presetColors"
              :key="c"
              class="w-7 h-7 rounded-full border-2 transition-all"
              :class="pointTypeForm.color === c ? 'border-white scale-110' : 'border-transparent'"
              :style="{ backgroundColor: c }"
              @click="pointTypeForm.color = c"
            />
            <input
              v-model="pointTypeForm.color"
              type="color"
              class="w-7 h-7 rounded-full border-0 cursor-pointer bg-transparent"
            />
          </div>
        </div>
        <!-- Preview -->
        <div class="flex items-center gap-3 p-3 bg-[#0A0F1A] border border-[#1E293B] rounded-lg">
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold text-white"
            :style="{ backgroundColor: pointTypeForm.color }"
          >
            {{ pointTypeForm.symbol || '?' }}
          </div>
          <div>
            <p class="text-sm text-[#F1F5F9]">{{ pointTypeForm.name || 'Point Type' }}</p>
            <p class="text-xs text-[#64748B]">Preview</p>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-[#64748B] hover:text-[#F1F5F9] transition-colors" @click="showPointTypeModal = false">Cancel</button>
        <button
          class="px-4 py-2 bg-[#48BB78] text-black text-sm font-medium rounded-lg hover:bg-[#48BB78]/90 disabled:opacity-50 transition-colors"
          :disabled="!pointTypeFormValid || savingPointType"
          @click="savePointType"
        >
          {{ savingPointType ? 'Saving...' : (editingPointTypeId ? 'Save Changes' : 'Add Point Type') }}
        </button>
      </template>
    </Modal>

    <!-- Create/Edit Level Modal (D01) -->
    <Modal
      :open="showCreate"
      :title="editingId ? 'Edit Level' : 'Create Level'"
      @close="showCreate = false"
    >
      <div class="space-y-4">
        <!-- Level Name -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Level Name <span class="text-[#DC2626]">*</span>
            <span class="text-[#64748B] text-xs ml-1">{{ form.name.length }}/30</span>
          </label>
          <input
            v-model="form.name"
            type="text"
            maxlength="30"
            placeholder="e.g. Bronze, Silver, Gold"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm"
          />
        </div>

        <!-- Point Type -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">Point Type <span class="text-[#DC2626]">*</span></label>
          <select
            v-model="form.point_type_id"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
          >
            <option value="" disabled>Select point type</option>
            <option v-for="pt in pointTypes" :key="pt.id" :value="pt.id">{{ pt.name }} ({{ pt.symbol }})</option>
          </select>
        </div>

        <!-- Threshold -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Threshold <span class="text-[#DC2626]">*</span>
            <span class="text-[#64748B] text-xs ml-1">Must be monotonically increasing</span>
          </label>
          <input
            v-model.number="form.threshold"
            type="number"
            min="0"
            placeholder="Points required for this level"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm"
          />
        </div>

        <!-- Level Badge -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-2">Level Badge</label>
          <div class="grid grid-cols-7 gap-2">
            <button
              v-for="icon in badgeIcons"
              :key="icon"
              class="w-10 h-10 rounded-lg flex items-center justify-center transition-all"
              :class="form.badge_icon === icon
                ? 'bg-[#48BB78]/20 border-2 border-[#48BB78]'
                : 'bg-[#0A0F1A] border border-[#1E293B] hover:border-[#94A3B8]'"
              @click="form.badge_icon = icon"
            >
              <span
                class="material-symbols-rounded text-xl"
                :class="form.badge_icon === icon ? 'text-[#48BB78]' : 'text-[#94A3B8]'"
              >{{ icon }}</span>
            </button>
          </div>
        </div>

        <!-- Level Perks -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Level Perks
            <span class="text-[#64748B] text-xs ml-1">{{ form.perks.length }}/200</span>
          </label>
          <textarea
            v-model="form.perks"
            maxlength="200"
            rows="3"
            placeholder="Describe perks for this level..."
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
          {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Create Level') }}
        </button>
      </template>
    </Modal>

    <!-- Reassignment Dialog -->
    <Modal
      :open="showReassign"
      title="Reassign Members"
      @close="showReassign = false"
    >
      <div class="space-y-4">
        <div class="flex items-start gap-3 p-3 bg-[#2D1515]/50 border border-[#DC2626]/20 rounded-lg">
          <span class="material-symbols-rounded text-lg text-[#DC2626] mt-0.5">warning</span>
          <div>
            <p class="text-sm text-[#F1F5F9] font-medium">This level has {{ deletingLevel?.members_at_level.toLocaleString() }} members</p>
            <p class="text-xs text-[#94A3B8] mt-1">Choose a level to reassign these members to before deleting.</p>
          </div>
        </div>
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">Reassign to <span class="text-[#DC2626]">*</span></label>
          <select
            v-model="reassignTargetId"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
          >
            <option value="" disabled>Select target level</option>
            <option v-for="opt in reassignOptions" :key="opt.id" :value="opt.id">
              {{ opt.name }} ({{ formatNumber(opt.threshold) }} pts)
            </option>
          </select>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-[#64748B] hover:text-[#F1F5F9] transition-colors" @click="showReassign = false">Cancel</button>
        <button
          class="px-4 py-2 bg-[#DC2626] text-white text-sm font-medium rounded-lg hover:bg-[#DC2626]/90 disabled:opacity-50 transition-colors"
          :disabled="!reassignTargetId"
          @click="deleteWithReassign"
        >
          Delete & Reassign
        </button>
      </template>
    </Modal>
  </div>
</template>
