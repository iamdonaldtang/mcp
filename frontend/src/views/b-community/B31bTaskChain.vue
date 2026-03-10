<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// === Types ===
interface TaskChainItem {
  id: string
  name: string
  status: CampaignStatus
  steps_count: number
  completions: number
  completion_rate: number
  created_at: string
}

interface StepFunnelEntry {
  step: number
  label: string
  count: number
  pct: number
}

interface ChainStep {
  task_id: string
  label: string
}

interface ChainForm {
  name: string
  steps: ChainStep[]
  reward_points: number
  reward_badge_id: string
}

// === State ===
const loading = ref(true)
const chains = ref<TaskChainItem[]>([])
const totalChains = ref(0)
const page = ref(1)
const pageSize = 20

// Stats
const statsTotal = ref(0)
const statsActive = ref(0)
const statsCompletions = ref(0)
const statsAvgRate = ref('0%')
const statsTrends = ref({ total: '', active: '', completions: '', rate: '' })

// Step funnel
const funnelData = ref<StepFunnelEntry[]>([])

// Filters
const filterStatus = ref<string>('all')
const searchQuery = ref('')
const searchDebounced = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null

// Modal
const showCreateModal = ref(false)
const editingChain = ref<TaskChainItem | null>(null)
const saving = ref(false)

const form = ref<ChainForm>({
  name: '',
  steps: [{ task_id: '', label: '' }, { task_id: '', label: '' }],
  reward_points: 100,
  reward_badge_id: '',
})

// Available tasks for step selector
const availableTasks = ref<{ id: string; name: string }[]>([])

// Drag state
const dragIndex = ref<number | null>(null)

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
    && form.value.steps.length >= 2
    && form.value.steps.every(s => s.task_id)
    && form.value.reward_points >= 1
})

const funnelMax = computed(() => Math.max(...funnelData.value.map(e => e.count), 1))

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
  await Promise.all([fetchChains(), fetchStats(), fetchFunnel(), fetchAvailableTasks()])
  loading.value = false
})

async function fetchChains() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    const res = await api.get('/api/v1/community/modules/taskchain', { params })
    chains.value = res.data.data?.items || []
    totalChains.value = res.data.data?.total || 0
  } catch { /* TODO: toast */ }
}

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/taskchain/stats')
    const s = res.data.data
    if (s) {
      statsTotal.value = s.total || 0
      statsActive.value = s.active || 0
      statsCompletions.value = s.completions || 0
      statsAvgRate.value = (s.avg_completion_rate ?? 0) + '%'
      statsTrends.value = {
        total: s.total_trend || '',
        active: s.active_trend || '',
        completions: s.completions_trend || '',
        rate: s.rate_trend || '',
      }
    }
  } catch { /* empty */ }
}

async function fetchFunnel() {
  try {
    const res = await api.get('/api/v1/community/modules/taskchain/funnel')
    funnelData.value = res.data.data || []
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
    const payload = {
      name: form.value.name.trim(),
      steps: form.value.steps,
      reward_points: form.value.reward_points,
      reward_badge_id: form.value.reward_badge_id || undefined,
    }
    if (editingChain.value) {
      await api.put(`/api/v1/community/modules/taskchain/${editingChain.value.id}`, payload)
    } else {
      await api.post('/api/v1/community/modules/taskchain', payload)
    }
    closeModal()
    await Promise.all([fetchChains(), fetchStats(), fetchFunnel()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function toggleStatus(chain: TaskChainItem) {
  const newStatus: CampaignStatus = chain.status === 'active' ? 'paused' : 'active'
  const oldStatus = chain.status
  chain.status = newStatus // optimistic
  try {
    await api.put(`/api/v1/community/modules/taskchain/${chain.id}`, { status: newStatus })
  } catch {
    chain.status = oldStatus // rollback
  }
}

async function duplicateChain(chain: TaskChainItem) {
  try {
    await api.post(`/api/v1/community/modules/taskchain/${chain.id}/duplicate`)
    await Promise.all([fetchChains(), fetchStats()])
  } catch { /* TODO: toast */ }
}

async function deleteChain(chain: TaskChainItem) {
  if (!confirm(`Delete "${chain.name}"? This action cannot be undone.`)) return
  try {
    await api.delete(`/api/v1/community/modules/taskchain/${chain.id}`)
    await Promise.all([fetchChains(), fetchStats(), fetchFunnel()])
  } catch { /* TODO: toast */ }
}

function openEdit(chain: TaskChainItem) {
  editingChain.value = chain
  form.value = {
    name: chain.name,
    steps: [{ task_id: '', label: '' }, { task_id: '', label: '' }], // will be loaded from detail endpoint
    reward_points: 100,
    reward_badge_id: '',
  }
  showCreateModal.value = true
  // Load full chain detail
  api.get(`/api/v1/community/modules/taskchain/${chain.id}`).then(res => {
    const d = res.data.data
    if (d) {
      form.value.steps = d.steps || form.value.steps
      form.value.reward_points = d.reward_points ?? 100
      form.value.reward_badge_id = d.reward_badge_id || ''
    }
  }).catch(() => { /* empty */ })
}

function openCreate() {
  editingChain.value = null
  form.value = {
    name: '',
    steps: [{ task_id: '', label: '' }, { task_id: '', label: '' }],
    reward_points: 100,
    reward_badge_id: '',
  }
  showCreateModal.value = true
}

function closeModal() {
  showCreateModal.value = false
  editingChain.value = null
}

// === Step management ===
function addStep() {
  if (form.value.steps.length >= 20) return
  form.value.steps.push({ task_id: '', label: '' })
}

function removeStep(index: number) {
  if (form.value.steps.length <= 2) return
  form.value.steps.splice(index, 1)
}

function onDragStart(index: number) {
  dragIndex.value = index
}

function onDragOver(e: DragEvent, index: number) {
  e.preventDefault()
  if (dragIndex.value === null || dragIndex.value === index) return
  const dragged = form.value.steps.splice(dragIndex.value, 1)[0]
  form.value.steps.splice(index, 0, dragged)
  dragIndex.value = index
}

function onDragEnd() {
  dragIndex.value = null
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
        <h1 class="text-2xl font-bold text-text-primary mb-1">TaskChain</h1>
        <p class="text-sm text-text-secondary">Multi-step sequential task chains that guide users through onboarding and engagement paths</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New TaskChain
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard label="Total Chains" :value="statsTotal" icon="conversion_path" icon-color="#48BB78" :trend="statsTrends.total" />
      <StatsCard label="Active Chains" :value="statsActive" icon="play_circle" icon-color="#60A5FA" :trend="statsTrends.active" />
      <StatsCard label="Completions" :value="statsCompletions" icon="check_circle" icon-color="#F59E0B" :trend="statsTrends.completions" />
      <StatsCard label="Avg Completion Rate" :value="statsAvgRate" icon="percent" icon-color="#A78BFA" :trend="statsTrends.rate" />
    </div>

    <!-- Step Funnel Visualization -->
    <div v-if="funnelData.length > 0" class="bg-card-bg border border-border rounded-xl p-6">
      <h3 class="text-sm font-semibold text-text-primary mb-4">Step Funnel (All Active Chains)</h3>
      <div class="space-y-2">
        <div v-for="entry in funnelData" :key="entry.step" class="flex items-center gap-3">
          <span class="text-xs text-text-muted w-16 shrink-0 text-right">Step {{ entry.step }}</span>
          <div class="flex-1 h-6 bg-page-bg rounded overflow-hidden relative">
            <div
              class="h-full rounded transition-all duration-500"
              :style="{
                width: (entry.count / funnelMax * 100) + '%',
                backgroundColor: entry.pct > 75 ? '#48BB78' : entry.pct > 50 ? '#60A5FA' : entry.pct > 25 ? '#F59E0B' : '#DC2626',
              }"
            />
            <span class="absolute inset-0 flex items-center px-2 text-xs font-medium text-text-primary">
              {{ entry.label || 'Step ' + entry.step }} &mdash; {{ entry.count.toLocaleString() }} ({{ entry.pct }}%)
            </span>
          </div>
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
            placeholder="Search chains..."
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
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-20">Steps</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Completions</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Rate</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Created</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Actions</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading -->
          <template v-if="loading">
            <tr v-for="i in 5" :key="i">
              <td colspan="7" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse" /></td>
            </tr>
          </template>
          <!-- Empty -->
          <tr v-else-if="filteredChains.length === 0">
            <td colspan="7" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">conversion_path</span>
              <p class="text-sm text-text-muted">No task chains yet</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First TaskChain</button>
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
                <div class="w-8 h-8 rounded-lg bg-community/10 flex items-center justify-center">
                  <span class="material-symbols-rounded text-base text-community">conversion_path</span>
                </div>
                <span class="text-sm font-medium text-text-primary">{{ chain.name }}</span>
              </div>
            </td>
            <td class="px-4 py-3">
              <StatusBadge :status="chain.status" />
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ chain.steps_count }}</td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ chain.completions.toLocaleString() }}</td>
            <td class="px-4 py-3">
              <span class="text-sm font-medium" :class="chain.completion_rate >= 50 ? 'text-community' : 'text-text-secondary'">
                {{ chain.completion_rate }}%
              </span>
            </td>
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

    <!-- D02 TaskChain Editor Modal -->
    <Modal :open="showCreateModal" :title="editingChain ? 'Edit TaskChain' : 'Create TaskChain'" max-width="600px" @close="closeModal">
      <div class="space-y-5 max-h-[60vh] overflow-y-auto pr-1">
        <!-- Chain Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Chain Name <span class="text-status-paused">*</span></label>
          <input
            v-model="form.name"
            type="text"
            maxlength="50"
            placeholder="e.g. Onboarding Journey"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
          <span class="text-xs text-text-muted mt-1 block">{{ form.name.length }}/50</span>
        </div>

        <!-- Steps -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="text-sm text-text-secondary">Steps <span class="text-status-paused">*</span> <span class="text-text-muted">({{ form.steps.length }}/20, min 2)</span></label>
            <button
              v-if="form.steps.length < 20"
              class="text-xs text-community hover:underline"
              @click="addStep"
            >
              + Add Step
            </button>
          </div>
          <div class="space-y-2">
            <div
              v-for="(step, idx) in form.steps"
              :key="idx"
              class="flex items-center gap-2 bg-page-bg border border-border rounded-lg p-2"
              draggable="true"
              @dragstart="onDragStart(idx)"
              @dragover="onDragOver($event, idx)"
              @dragend="onDragEnd"
            >
              <!-- Grip handle -->
              <span class="material-symbols-rounded text-base text-text-muted cursor-grab shrink-0">drag_indicator</span>
              <!-- Step number -->
              <span class="text-xs text-text-muted w-5 shrink-0 text-center">{{ idx + 1 }}</span>
              <!-- Task selector -->
              <select
                v-model="step.task_id"
                class="flex-1 px-3 py-1.5 bg-card-bg border border-border rounded text-sm text-text-primary focus:border-community focus:outline-none"
              >
                <option value="" disabled>Select task...</option>
                <option v-for="t in availableTasks" :key="t.id" :value="t.id">{{ t.name }}</option>
              </select>
              <!-- Step label -->
              <input
                v-model="step.label"
                type="text"
                placeholder="Label (optional)"
                class="w-32 px-3 py-1.5 bg-card-bg border border-border rounded text-sm text-text-primary placeholder-text-muted focus:border-community focus:outline-none"
              />
              <!-- Remove -->
              <button
                v-if="form.steps.length > 2"
                class="p-1 rounded hover:bg-white/5 transition-colors shrink-0"
                @click="removeStep(idx)"
              >
                <span class="material-symbols-rounded text-sm text-status-paused">close</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Completion Reward -->
        <div>
          <label class="block text-sm text-text-secondary mb-2">Completion Reward</label>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs text-text-muted mb-1">Points <span class="text-status-paused">*</span></label>
              <input
                v-model.number="form.reward_points"
                type="number"
                min="1"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
              />
            </div>
            <div>
              <label class="block text-xs text-text-muted mb-1">Badge (optional)</label>
              <input
                v-model="form.reward_badge_id"
                type="text"
                placeholder="Badge ID"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-community focus:outline-none"
              />
            </div>
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
          {{ saving ? 'Saving...' : editingChain ? 'Update Chain' : 'Create Chain' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
