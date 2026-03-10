<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// === Types ===
interface DayChainItem {
  id: string
  name: string
  status: CampaignStatus
  daily_task: string
  daily_task_name: string
  base_reward: number
  active_streaks: number
  avg_streak_length: number
  created_at: string
}

interface StreakBar {
  day: number
  count: number
  pct: number
}

interface MilestoneBonus {
  day: number
  multiplier: number
}

interface DayChainForm {
  name: string
  daily_task_id: string
  base_reward: number
  milestones: MilestoneBonus[]
  grace_period_hours: number
}

// === State ===
const loading = ref(true)
const chains = ref<DayChainItem[]>([])
const totalChains = ref(0)
const page = ref(1)
const pageSize = 20

// Stats
const statsStreakRate = ref('0%')
const statsCompletionRate = ref('0%')
const statsDay7Pass = ref('0%')
const statsAvgStreak = ref('0')
const statsTrends = ref({ streak: '', completion: '', day7: '', avg: '' })

// Streak distribution
const streakData = ref<StreakBar[]>([])
const day7PassThrough = ref(100)

// Filters
const filterStatus = ref<string>('all')
const searchQuery = ref('')
const searchDebounced = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null

// Modal
const showCreateModal = ref(false)
const editingChain = ref<DayChainItem | null>(null)
const saving = ref(false)

const form = ref<DayChainForm>({
  name: '',
  daily_task_id: '',
  base_reward: 10,
  milestones: [{ day: 7, multiplier: 2 }],
  grace_period_hours: 4,
})

// Available tasks for selector
const availableTasks = ref<{ id: string; name: string }[]>([])

// === Computed ===
const filteredChains = computed(() => {
  return chains.value.filter(c => {
    if (filterStatus.value !== 'all' && c.status !== filterStatus.value) return false
    if (searchDebounced.value && !c.name.toLowerCase().includes(searchDebounced.value.toLowerCase())) return false
    return true
  })
})

const canSave = computed(() => {
  return form.value.name.trim().length >= 1
    && form.value.name.trim().length <= 50
    && form.value.daily_task_id
    && form.value.base_reward >= 1
    && form.value.grace_period_hours >= 0
    && form.value.grace_period_hours <= 24
})

const streakMax = computed(() => Math.max(...streakData.value.map(e => e.count), 1))

const showBottleneckBanner = computed(() => day7PassThrough.value < 60)

// Sorted milestones
const sortedMilestones = computed(() => [...form.value.milestones].sort((a, b) => a.day - b.day))

// === Search debounce ===
function onSearchInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    searchDebounced.value = searchQuery.value
  }, 300)
}

onUnmounted(() => {
  if (debounceTimer) clearTimeout(debounceTimer)
})

// === API ===
onMounted(async () => {
  await Promise.all([fetchChains(), fetchStats(), fetchStreakDistribution(), fetchAvailableTasks()])
  loading.value = false
})

async function fetchChains() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    const res = await api.get('/api/v1/community/modules/daychain', { params })
    chains.value = res.data.data?.items || []
    totalChains.value = res.data.data?.total || 0
  } catch { /* TODO: toast */ }
}

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/daychain/stats')
    const s = res.data.data
    if (s) {
      statsStreakRate.value = (s.active_streak_rate ?? 0) + '%'
      statsCompletionRate.value = (s.completion_rate ?? 0) + '%'
      statsDay7Pass.value = (s.day7_pass_through ?? 0) + '%'
      statsAvgStreak.value = (s.avg_streak_days ?? 0) + ' days'
      day7PassThrough.value = s.day7_pass_through ?? 100
      statsTrends.value = {
        streak: s.streak_trend || '',
        completion: s.completion_trend || '',
        day7: s.day7_trend || '',
        avg: s.avg_trend || '',
      }
    }
  } catch { /* empty */ }
}

async function fetchStreakDistribution() {
  try {
    const res = await api.get('/api/v1/community/modules/daychain/streak-distribution')
    streakData.value = res.data.data || []
  } catch { /* empty */ }
}

async function fetchAvailableTasks() {
  try {
    const res = await api.get('/api/v1/community/tasks', { params: { page_size: 200, status: 'active' } })
    availableTasks.value = (res.data.data?.items || []).map((t: { id: string; name: string }) => ({ id: t.id, name: t.name }))
  } catch { /* empty */ }
}

// === CRUD ===
async function saveChain() {
  if (!canSave.value || saving.value) return
  saving.value = true
  try {
    // Auto-sort milestones before sending
    const payload = {
      name: form.value.name.trim(),
      daily_task_id: form.value.daily_task_id,
      base_reward: form.value.base_reward,
      milestones: [...form.value.milestones].sort((a, b) => a.day - b.day),
      grace_period_hours: form.value.grace_period_hours,
    }
    if (editingChain.value) {
      await api.put(`/api/v1/community/modules/daychain/${editingChain.value.id}`, payload)
    } else {
      await api.post('/api/v1/community/modules/daychain', payload)
    }
    closeModal()
    await Promise.all([fetchChains(), fetchStats(), fetchStreakDistribution()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function toggleStatus(chain: DayChainItem) {
  const newStatus: CampaignStatus = chain.status === 'active' ? 'paused' : 'active'
  const oldStatus = chain.status
  chain.status = newStatus // optimistic
  try {
    await api.put(`/api/v1/community/modules/daychain/${chain.id}`, { status: newStatus })
  } catch {
    chain.status = oldStatus // rollback
  }
}

async function duplicateChain(chain: DayChainItem) {
  try {
    await api.post(`/api/v1/community/modules/daychain/${chain.id}/duplicate`)
    await Promise.all([fetchChains(), fetchStats()])
  } catch { /* TODO: toast */ }
}

async function deleteChain(chain: DayChainItem) {
  if (!confirm(`Delete "${chain.name}"? This action cannot be undone.`)) return
  try {
    await api.delete(`/api/v1/community/modules/daychain/${chain.id}`)
    await Promise.all([fetchChains(), fetchStats(), fetchStreakDistribution()])
  } catch { /* TODO: toast */ }
}

function openEdit(chain: DayChainItem) {
  editingChain.value = chain
  form.value = {
    name: chain.name,
    daily_task_id: chain.daily_task || '',
    base_reward: chain.base_reward,
    milestones: [{ day: 7, multiplier: 2 }],
    grace_period_hours: 4,
  }
  showCreateModal.value = true
  // Load full detail
  api.get(`/api/v1/community/modules/daychain/${chain.id}`).then(res => {
    const d = res.data.data
    if (d) {
      form.value.daily_task_id = d.daily_task_id || form.value.daily_task_id
      form.value.milestones = d.milestones || form.value.milestones
      form.value.grace_period_hours = d.grace_period_hours ?? 4
    }
  }).catch(() => { /* empty */ })
}

function openCreate() {
  editingChain.value = null
  form.value = {
    name: '',
    daily_task_id: '',
    base_reward: 10,
    milestones: [{ day: 7, multiplier: 2 }],
    grace_period_hours: 4,
  }
  showCreateModal.value = true
}

function closeModal() {
  showCreateModal.value = false
  editingChain.value = null
}

// === Milestone management ===
function addMilestone() {
  const maxDay = form.value.milestones.length > 0
    ? Math.max(...form.value.milestones.map(m => m.day))
    : 0
  form.value.milestones.push({ day: maxDay + 7, multiplier: 2 })
}

function removeMilestone(index: number) {
  form.value.milestones.splice(index, 1)
}

// === Pagination ===
function onPageChange(p: number) {
  page.value = p
  fetchChains()
}

function onFilterChange(status: string) {
  filterStatus.value = status
  page.value = 1
  fetchChains()
}

function formatDate(d: string) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">DayChain</h1>
        <p class="text-sm text-text-secondary">Daily streak challenges that drive habitual engagement through loss aversion</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New DayChain
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard label="Active Streak Rate" :value="statsStreakRate" icon="local_fire_department" icon-color="#F59E0B" :trend="statsTrends.streak" />
      <StatsCard label="Completion Rate" :value="statsCompletionRate" icon="check_circle" icon-color="#48BB78" :trend="statsTrends.completion" />
      <StatsCard label="Day 7 Pass-through" :value="statsDay7Pass" icon="filter_7" icon-color="#60A5FA" :trend="statsTrends.day7" />
      <StatsCard label="Avg Streak Days" :value="statsAvgStreak" icon="event_repeat" icon-color="#A78BFA" :trend="statsTrends.avg" />
    </div>

    <!-- Streak Bottleneck Banner -->
    <div
      v-if="showBottleneckBanner && !loading"
      class="flex items-start gap-3 bg-c-accent/10 border border-c-accent/30 rounded-xl p-4"
    >
      <span class="material-symbols-rounded text-xl text-c-accent shrink-0 mt-0.5">warning</span>
      <div>
        <h3 class="text-sm font-semibold text-text-primary mb-0.5">Streak Bottleneck Detected</h3>
        <p class="text-xs text-text-secondary">
          Day 7 pass-through is below 60% ({{ day7PassThrough }}%). Consider adding a Day 7 milestone bonus to incentivize users past the critical drop-off point.
        </p>
      </div>
    </div>

    <!-- Streak Distribution Chart -->
    <div v-if="streakData.length > 0" class="bg-card-bg border border-border rounded-xl p-6">
      <h3 class="text-sm font-semibold text-text-primary mb-4">Streak Distribution</h3>
      <div class="space-y-1.5">
        <div v-for="bar in streakData" :key="bar.day" class="flex items-center gap-3">
          <span
            class="text-xs w-12 shrink-0 text-right font-medium"
            :class="bar.day === 7 && bar.pct < 40 ? 'text-c-accent' : 'text-text-muted'"
          >
            Day {{ bar.day }}
          </span>
          <div class="flex-1 h-5 bg-page-bg rounded overflow-hidden relative">
            <div
              class="h-full rounded transition-all duration-500"
              :style="{
                width: (bar.count / streakMax * 100) + '%',
                backgroundColor: bar.day === 7 && bar.pct < 40 ? '#F59E0B' : '#48BB78',
              }"
            />
            <span class="absolute inset-0 flex items-center px-2 text-[11px] font-medium text-text-primary">
              {{ bar.count.toLocaleString() }} ({{ bar.pct }}%)
            </span>
          </div>
          <!-- Highlight Day 7 drop-off -->
          <span
            v-if="bar.day === 7 && bar.pct < 40"
            class="text-[10px] text-c-accent font-medium shrink-0"
          >
            &gt;40% drop-off
          </span>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="flex gap-2">
        <button
          v-for="status in ['all', 'active', 'paused', 'draft']"
          :key="status"
          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
          :class="filterStatus === status ? 'bg-community text-white' : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
          @click="onFilterChange(status)"
        >
          {{ status === 'all' ? 'All' : status.charAt(0).toUpperCase() + status.slice(1) }}
        </button>
      </div>
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search day chains..."
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
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Chain Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-36">Daily Task</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Base Reward</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Active Streaks</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Avg Length</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Created</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Actions</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading -->
          <template v-if="loading">
            <tr v-for="i in 5" :key="i">
              <td colspan="8" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse" /></td>
            </tr>
          </template>
          <!-- Empty -->
          <tr v-else-if="filteredChains.length === 0">
            <td colspan="8" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">local_fire_department</span>
              <p class="text-sm text-text-muted">No day chains yet</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First DayChain</button>
            </td>
          </tr>
          <!-- Rows -->
          <tr
            v-else
            v-for="chain in filteredChains"
            :key="chain.id"
            class="border-b border-border last:border-b-0 hover:bg-white/[0.02] transition-colors"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-c-accent/10 flex items-center justify-center">
                  <span class="material-symbols-rounded text-base text-c-accent">local_fire_department</span>
                </div>
                <span class="text-sm font-medium text-text-primary">{{ chain.name }}</span>
              </div>
            </td>
            <td class="px-4 py-3">
              <StatusBadge :status="chain.status" />
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary truncate max-w-[144px]">{{ chain.daily_task_name || '—' }}</td>
            <td class="px-4 py-3 text-sm text-c-accent font-medium">+{{ chain.base_reward }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ chain.active_streaks.toLocaleString() }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ chain.avg_streak_length }}d</td>
            <td class="px-4 py-3 text-sm text-text-muted">{{ formatDate(chain.created_at) }}</td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit" @click="openEdit(chain)">
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  :title="chain.status === 'active' ? 'Pause' : 'Activate'"
                  @click="toggleStatus(chain)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">
                    {{ chain.status === 'active' ? 'pause_circle' : 'play_circle' }}
                  </span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Duplicate" @click="duplicateChain(chain)">
                  <span class="material-symbols-rounded text-base text-text-muted">content_copy</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteChain(chain)">
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination v-if="totalChains > pageSize" :page="page" :page-size="pageSize" :total="totalChains" @update:page="onPageChange" />
    </div>

    <!-- D03 DayChain Config Modal -->
    <Modal :open="showCreateModal" :title="editingChain ? 'Edit DayChain' : 'Create DayChain'" max-width="560px" @close="closeModal">
      <div class="space-y-5 max-h-[60vh] overflow-y-auto pr-1">
        <!-- Chain Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Chain Name <span class="text-status-paused">*</span></label>
          <input
            v-model="form.name"
            type="text"
            maxlength="50"
            placeholder="e.g. Daily Login Streak"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
          <span class="text-xs text-text-muted mt-1 block">{{ form.name.length }}/50</span>
        </div>

        <!-- Daily Task -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Daily Task <span class="text-status-paused">*</span></label>
          <select
            v-model="form.daily_task_id"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          >
            <option value="" disabled>Select a task...</option>
            <option v-for="t in availableTasks" :key="t.id" :value="t.id">{{ t.name }}</option>
          </select>
        </div>

        <!-- Base Reward -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Base Reward (points per day) <span class="text-status-paused">*</span></label>
          <input
            v-model.number="form.base_reward"
            type="number"
            min="1"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          />
        </div>

        <!-- Milestone Bonuses -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="text-sm text-text-secondary">Milestone Bonuses</label>
            <button class="text-xs text-community hover:underline" @click="addMilestone">+ Add Milestone</button>
          </div>
          <div v-if="form.milestones.length === 0" class="text-xs text-text-muted bg-page-bg border border-border rounded-lg p-3 text-center">
            No milestones configured. Users earn base reward only.
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="(ms, idx) in sortedMilestones"
              :key="idx"
              class="flex items-center gap-3 bg-page-bg border border-border rounded-lg p-2.5"
            >
              <div class="flex-1">
                <label class="block text-[11px] text-text-muted mb-0.5">Day</label>
                <input
                  v-model.number="form.milestones[idx].day"
                  type="number"
                  min="1"
                  max="365"
                  class="w-full px-3 py-1.5 bg-card-bg border border-border rounded text-sm text-text-primary focus:border-community focus:outline-none"
                />
              </div>
              <div class="flex-1">
                <label class="block text-[11px] text-text-muted mb-0.5">Multiplier</label>
                <input
                  v-model.number="form.milestones[idx].multiplier"
                  type="number"
                  min="1"
                  max="100"
                  step="0.5"
                  class="w-full px-3 py-1.5 bg-card-bg border border-border rounded text-sm text-text-primary focus:border-community focus:outline-none"
                />
              </div>
              <div class="pt-4">
                <span class="text-xs text-text-muted">= {{ form.base_reward * (form.milestones[idx]?.multiplier || 1) }} pts</span>
              </div>
              <button class="pt-4 p-1 rounded hover:bg-white/5 transition-colors" @click="removeMilestone(idx)">
                <span class="material-symbols-rounded text-sm text-status-paused">close</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Grace Period -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">
            Grace Period: <span class="font-medium text-text-primary">{{ form.grace_period_hours }}h</span>
          </label>
          <input
            v-model.number="form.grace_period_hours"
            type="range"
            min="0"
            max="24"
            step="1"
            class="w-full h-2 bg-page-bg rounded-lg appearance-none cursor-pointer accent-community"
          />
          <div class="flex justify-between text-[11px] text-text-muted mt-1">
            <span>0h (strict)</span>
            <span>24h (lenient)</span>
          </div>
        </div>
      </div>

      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="closeModal">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!canSave || saving"
          @click="saveChain"
        >
          {{ saving ? 'Saving...' : editingChain ? 'Update DayChain' : 'Create DayChain' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
