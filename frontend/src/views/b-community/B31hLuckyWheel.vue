<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { api } from '../../api/client'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// === Types ===
interface Prize {
  name: string
  type: 'points' | 'badge' | 'shop_item' | 'nothing'
  value: string
  probability: number
}

interface WheelItem {
  id: string
  name: string
  status: CampaignStatus
  total_spins: number
  prizes: Prize[]
  win_rate: number
  spin_cost: number
  spin_limit: 'once' | 'daily' | 'unlimited'
  start_date: string | null
  end_date: string | null
  created_at: string
}

interface WheelStats {
  total_spins: number
  winners: number
  prizes_awarded: number
  avg_spins_per_day: number
}

// === State ===
const loading = ref(true)
const wheels = ref<WheelItem[]>([])
const totalWheels = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<WheelStats>({ total_spins: 0, winners: 0, prizes_awarded: 0, avg_spins_per_day: 0 })

// Filters
const filterStatus = ref<string>('all')
const searchQuery = ref('')

const statusTabs = [
  { key: 'all', label: 'All' },
  { key: 'active', label: 'Active' },
  { key: 'paused', label: 'Paused' },
  { key: 'draft', label: 'Draft' },
  { key: 'completed', label: 'Ended' },
]

// Modal
const showCreateModal = ref(false)
const editingWheel = ref<WheelItem | null>(null)
const saving = ref(false)

const defaultPrizes: Prize[] = [
  { name: '10 Points', type: 'points', value: '10', probability: 40 },
  { name: 'Nothing', type: 'nothing', value: '', probability: 60 },
]

const formDefaults = {
  name: '',
  prizes: [...defaultPrizes.map(p => ({ ...p }))] as Prize[],
  spin_cost: 0,
  spin_limit: 'daily' as 'once' | 'daily' | 'unlimited',
  has_duration: false,
  start_date: '',
  end_date: '',
}

const form = ref({ ...formDefaults })

const prizeTypes = [
  { key: 'points', label: 'Points' },
  { key: 'badge', label: 'Badge' },
  { key: 'shop_item', label: 'Shop Item' },
  { key: 'nothing', label: 'Nothing' },
]

// Wheel preview colors
const segmentColors = [
  '#48BB78', '#5D7EF1', '#ED8936', '#9B7EE0',
  '#F59E0B', '#EC4899', '#14B8A6', '#EF4444',
  '#06B6D4', '#8B5CF6', '#F97316', '#22C55E',
]

// === Computed ===
const filteredWheels = computed(() => {
  return wheels.value.filter(w => {
    if (filterStatus.value !== 'all' && w.status !== filterStatus.value) return false
    if (searchQuery.value && !w.name.toLowerCase().includes(searchQuery.value.toLowerCase())) return false
    return true
  })
})

const totalProbability = computed(() => {
  return form.value.prizes.reduce((sum, p) => sum + (p.probability || 0), 0)
})

const probabilityValid = computed(() => Math.abs(totalProbability.value - 100) < 0.01)

const canSave = computed(() => {
  return form.value.name.trim().length >= 1
    && form.value.name.trim().length <= 50
    && form.value.prizes.length >= 2
    && probabilityValid.value
    && form.value.prizes.every(p => p.name.trim().length > 0)
    && form.value.spin_cost >= 0
})

// === API ===
onMounted(async () => {
  await Promise.all([fetchWheels(), fetchStats()])
  loading.value = false
})

async function fetchWheels() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    const res = await api.get('/api/v1/community/modules/wheel', { params })
    wheels.value = res.data.data?.items || []
    totalWheels.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/wheel/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

function openCreate() {
  editingWheel.value = null
  form.value = {
    ...formDefaults,
    prizes: defaultPrizes.map(p => ({ ...p })),
  }
  showCreateModal.value = true
}

function openEdit(wheel: WheelItem) {
  editingWheel.value = wheel
  form.value = {
    name: wheel.name,
    prizes: wheel.prizes.map(p => ({ ...p })),
    spin_cost: wheel.spin_cost,
    spin_limit: wheel.spin_limit,
    has_duration: !!(wheel.start_date || wheel.end_date),
    start_date: wheel.start_date || '',
    end_date: wheel.end_date || '',
  }
  showCreateModal.value = true
}

async function saveWheel() {
  if (!canSave.value) return
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      prizes: form.value.prizes,
      spin_cost: form.value.spin_cost,
      spin_limit: form.value.spin_limit,
      start_date: form.value.has_duration && form.value.start_date ? form.value.start_date : null,
      end_date: form.value.has_duration && form.value.end_date ? form.value.end_date : null,
    }

    if (editingWheel.value) {
      await api.put(`/api/v1/community/modules/wheel/${editingWheel.value.id}`, payload)
    } else {
      await api.post('/api/v1/community/modules/wheel', payload)
    }

    showCreateModal.value = false
    await Promise.all([fetchWheels(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function toggleWheelStatus(wheel: WheelItem) {
  const newStatus = wheel.status === 'active' ? 'paused' : 'active'
  const oldStatus = wheel.status
  wheel.status = newStatus as CampaignStatus
  try {
    await api.put(`/api/v1/community/modules/wheel/${wheel.id}`, { status: newStatus })
  } catch {
    wheel.status = oldStatus as CampaignStatus
  }
}

async function duplicateWheel(wheel: WheelItem) {
  try {
    await api.post(`/api/v1/community/modules/wheel/${wheel.id}/duplicate`)
    await Promise.all([fetchWheels(), fetchStats()])
  } catch { /* empty */ }
}

async function deleteWheel(id: string) {
  if (!confirm('Delete this wheel? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/modules/wheel/${id}`)
    await Promise.all([fetchWheels(), fetchStats()])
  } catch { /* empty */ }
}

function addPrize() {
  form.value.prizes.push({ name: '', type: 'points', value: '', probability: 0 })
}

function removePrize(index: number) {
  if (form.value.prizes.length <= 2) return
  form.value.prizes.splice(index, 1)
}

function formatDate(dateStr: string) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function onPageChange(p: number) {
  page.value = p
  fetchWheels()
}

// Generate CSS conic-gradient for wheel preview
const wheelGradient = computed(() => {
  const prizes = form.value.prizes
  if (prizes.length === 0) return 'conic-gradient(#1E293B 0deg 360deg)'
  let segments: string[] = []
  let currentDeg = 0
  prizes.forEach((prize, i) => {
    const degrees = (prize.probability / 100) * 360
    const color = segmentColors[i % segmentColors.length]
    segments.push(`${color} ${currentDeg}deg ${currentDeg + degrees}deg`)
    currentDeg += degrees
  })
  return `conic-gradient(${segments.join(', ')})`
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Lucky Wheel</h1>
        <p class="text-sm text-text-secondary">Create and manage prize wheels for your community</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New Wheel
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Total Spins</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.total_spins.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Winners</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.winners.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Prizes Awarded</div>
        <div class="text-2xl font-bold text-community">{{ stats.prizes_awarded.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Avg Spins / Day</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.avg_spins_per_day.toLocaleString() }}</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="flex gap-2">
        <button
          v-for="tab in statusTabs"
          :key="tab.key"
          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
          :class="filterStatus === tab.key ? 'bg-community text-white' : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
          @click="filterStatus = tab.key; fetchWheels()"
        >
          {{ tab.label }}
        </button>
      </div>
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search wheels..."
            class="w-full pl-10 pr-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-community focus:outline-none"
          />
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-border">
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Wheel Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Total Spins</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Prize Pool</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Win Rate</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Created</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="7" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse"></div></td>
          </tr>
          <tr v-else-if="filteredWheels.length === 0">
            <td colspan="7" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">casino</span>
              <p class="text-sm text-text-muted">No lucky wheels yet</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First Wheel</button>
            </td>
          </tr>
          <tr
            v-else
            v-for="wheel in filteredWheels"
            :key="wheel.id"
            class="border-b border-border last:border-b-0 hover:bg-white/[0.02] transition-colors"
          >
            <td class="px-4 py-3">
              <span class="text-sm font-medium text-text-primary">{{ wheel.name }}</span>
            </td>
            <td class="px-4 py-3">
              <StatusBadge :status="wheel.status" />
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ wheel.total_spins.toLocaleString() }}</td>
            <td class="px-4 py-3">
              <span class="text-sm text-text-secondary">{{ wheel.prizes.length }} prizes</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-sm font-medium" :class="wheel.win_rate > 0 ? 'text-community' : 'text-text-muted'">
                {{ wheel.win_rate.toFixed(1) }}%
              </span>
            </td>
            <td class="px-4 py-3 text-sm text-text-muted">{{ formatDate(wheel.created_at) }}</td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  :title="wheel.status === 'active' ? 'Pause' : 'Activate'"
                  @click="toggleWheelStatus(wheel)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">
                    {{ wheel.status === 'active' ? 'pause_circle' : 'play_circle' }}
                  </span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit" @click="openEdit(wheel)">
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Duplicate" @click="duplicateWheel(wheel)">
                  <span class="material-symbols-rounded text-base text-text-muted">content_copy</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteWheel(wheel.id)">
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination v-if="totalWheels > pageSize" :page="page" :page-size="pageSize" :total="totalWheels" @update:page="onPageChange" />
    </div>

    <!-- D08 Lucky Wheel Config Modal -->
    <Modal
      :open="showCreateModal"
      :title="editingWheel ? 'Edit Lucky Wheel' : 'Create Lucky Wheel'"
      max-width="720px"
      @close="showCreateModal = false"
    >
      <div class="flex gap-6 max-h-[65vh]">
        <!-- Left: Form -->
        <div class="flex-1 space-y-4 overflow-y-auto pr-2">
          <!-- Wheel Name -->
          <div>
            <label class="block text-sm text-text-secondary mb-1">Wheel Name <span class="text-status-paused">*</span></label>
            <input
              v-model="form.name"
              type="text"
              maxlength="50"
              placeholder="e.g. Weekly Spin"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
            />
          </div>

          <!-- Prizes -->
          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="text-sm text-text-secondary">Prizes <span class="text-status-paused">*</span> (min 2)</label>
              <button
                class="text-xs text-community hover:underline"
                @click="addPrize"
              >
                + Add Prize
              </button>
            </div>
            <div class="space-y-2">
              <div
                v-for="(prize, idx) in form.prizes"
                :key="idx"
                class="flex items-start gap-2 bg-page-bg border border-border rounded-lg p-3"
              >
                <div
                  class="w-3 h-3 rounded-full mt-2.5 flex-shrink-0"
                  :style="{ backgroundColor: segmentColors[idx % segmentColors.length] }"
                ></div>
                <div class="flex-1 space-y-2">
                  <input
                    v-model="prize.name"
                    type="text"
                    placeholder="Prize name"
                    class="w-full px-3 py-1.5 bg-card-bg border border-border rounded text-text-primary text-sm focus:border-community focus:outline-none"
                  />
                  <div class="flex gap-2">
                    <select
                      v-model="prize.type"
                      class="w-28 px-2 py-1.5 bg-card-bg border border-border rounded text-text-primary text-xs focus:border-community focus:outline-none"
                    >
                      <option v-for="pt in prizeTypes" :key="pt.key" :value="pt.key">{{ pt.label }}</option>
                    </select>
                    <input
                      v-if="prize.type !== 'nothing'"
                      v-model="prize.value"
                      type="text"
                      placeholder="Value"
                      class="flex-1 px-2 py-1.5 bg-card-bg border border-border rounded text-text-primary text-xs focus:border-community focus:outline-none"
                    />
                    <div class="flex items-center gap-1">
                      <input
                        v-model.number="prize.probability"
                        type="number"
                        min="0"
                        max="100"
                        step="0.1"
                        class="w-16 px-2 py-1.5 bg-card-bg border border-border rounded text-text-primary text-xs text-right focus:border-community focus:outline-none"
                      />
                      <span class="text-xs text-text-muted">%</span>
                    </div>
                  </div>
                </div>
                <button
                  v-if="form.prizes.length > 2"
                  class="p-1 rounded hover:bg-white/5 mt-1"
                  @click="removePrize(idx)"
                >
                  <span class="material-symbols-rounded text-sm text-status-paused">close</span>
                </button>
              </div>
            </div>

            <!-- Probability total -->
            <div class="mt-2 flex items-center gap-2 text-sm">
              <span class="text-text-muted">Total:</span>
              <span
                class="font-semibold"
                :class="probabilityValid ? 'text-community' : 'text-status-paused'"
              >
                {{ totalProbability.toFixed(1) }}%
              </span>
              <span v-if="!probabilityValid" class="text-xs text-status-paused">Must equal 100%</span>
            </div>
          </div>

          <!-- Spin Cost -->
          <div>
            <label class="block text-sm text-text-secondary mb-1">Spin Cost</label>
            <div class="flex items-center gap-2">
              <input
                v-model.number="form.spin_cost"
                type="number"
                min="0"
                class="w-32 px-3 py-2 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
              />
              <span class="text-xs text-text-muted">points (0 = free)</span>
            </div>
          </div>

          <!-- Spin Limit -->
          <div>
            <label class="block text-sm text-text-secondary mb-1">Spin Limit</label>
            <div class="flex items-center gap-4">
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="form.spin_limit" type="radio" value="once" class="accent-community" />
                <span class="text-sm text-text-primary">Once</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="form.spin_limit" type="radio" value="daily" class="accent-community" />
                <span class="text-sm text-text-primary">Daily</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="form.spin_limit" type="radio" value="unlimited" class="accent-community" />
                <span class="text-sm text-text-primary">Unlimited</span>
              </label>
            </div>
          </div>

          <!-- Duration -->
          <div>
            <label class="flex items-center gap-2 cursor-pointer mb-2">
              <input v-model="form.has_duration" type="checkbox" class="accent-community" />
              <span class="text-sm text-text-secondary">Set Duration</span>
            </label>
            <div v-if="form.has_duration" class="flex gap-3">
              <div class="flex-1">
                <label class="block text-xs text-text-muted mb-1">Start</label>
                <input
                  v-model="form.start_date"
                  type="datetime-local"
                  class="w-full px-3 py-2 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
                />
              </div>
              <div class="flex-1">
                <label class="block text-xs text-text-muted mb-1">End</label>
                <input
                  v-model="form.end_date"
                  type="datetime-local"
                  class="w-full px-3 py-2 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Wheel Preview -->
        <div class="w-56 flex-shrink-0 flex flex-col items-center justify-start pt-4">
          <div class="text-xs text-text-muted uppercase tracking-wider mb-3">Preview</div>
          <div class="relative w-48 h-48">
            <!-- Wheel -->
            <div
              class="w-48 h-48 rounded-full border-4 border-border shadow-lg"
              :style="{ background: wheelGradient }"
            ></div>
            <!-- Center dot -->
            <div class="absolute inset-0 flex items-center justify-center">
              <div class="w-10 h-10 rounded-full bg-card-bg border-2 border-border flex items-center justify-center">
                <span class="material-symbols-rounded text-community text-lg">casino</span>
              </div>
            </div>
            <!-- Pointer -->
            <div class="absolute -top-2 left-1/2 -translate-x-1/2">
              <div class="w-0 h-0 border-l-[8px] border-r-[8px] border-t-[12px] border-l-transparent border-r-transparent border-t-community"></div>
            </div>
          </div>
          <!-- Legend -->
          <div class="mt-4 space-y-1 w-full">
            <div
              v-for="(prize, idx) in form.prizes"
              :key="idx"
              class="flex items-center gap-2 text-xs"
            >
              <div
                class="w-2.5 h-2.5 rounded-full flex-shrink-0"
                :style="{ backgroundColor: segmentColors[idx % segmentColors.length] }"
              ></div>
              <span class="text-text-secondary truncate flex-1">{{ prize.name || 'Unnamed' }}</span>
              <span class="text-text-muted">{{ prize.probability }}%</span>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showCreateModal = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!canSave || saving"
          @click="saveWheel"
        >
          {{ saving ? 'Saving...' : editingWheel ? 'Update Wheel' : 'Create Wheel' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
