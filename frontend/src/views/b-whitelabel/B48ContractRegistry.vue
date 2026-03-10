<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'

// --- Types ---
interface ContractItem {
  id: string
  name: string
  network: string
  address: string
  abi: string
  monitored_events: string[]
  status: 'verified' | 'pending' | 'error'
  error_message?: string
  events_24h: number
  linked_rules: number
  created_at: string
  expanded?: boolean
  recent_events?: ContractEvent[]
}

interface ContractEvent {
  id: string
  event_name: string
  block_number: number
  tx_hash: string
  timestamp: string
  data: Record<string, unknown>
}

// --- State ---
const loading = ref(true)
const contracts = ref<ContractItem[]>([])
const totalContracts = ref(0)
const page = ref(1)
const pageSize = 20

const stats = ref({
  total: 0,
  verified: 0,
  events_24h: 0,
  linked_rules: 0,
})

// Filters
const filterTab = ref<string>('all')
const searchQuery = ref('')
let searchTimeout: ReturnType<typeof setTimeout> | null = null

// Register modal
const showRegister = ref(false)
const editingContract = ref<ContractItem | null>(null)
const saving = ref(false)
const verifying = ref<string | null>(null)
const duplicateWarning = ref('')

const contractForm = ref(getEmptyForm())

function getEmptyForm() {
  return {
    name: '',
    network: '',
    address: '',
    abi: '',
    monitored_events: [] as string[],
    manual_events: '',
  }
}

const networks = [
  { value: 'ethereum', label: 'Ethereum', color: '#627EEA' },
  { value: 'bsc', label: 'BSC', color: '#F3BA2F' },
  { value: 'polygon', label: 'Polygon', color: '#8247E5' },
  { value: 'arbitrum', label: 'Arbitrum', color: '#28A0F0' },
  { value: 'optimism', label: 'Optimism', color: '#FF0420' },
  { value: 'base', label: 'Base', color: '#0052FF' },
  { value: 'avalanche', label: 'Avalanche', color: '#E84142' },
]

const parsedAbiEvents = computed<string[]>(() => {
  if (!contractForm.value.abi) return []
  try {
    const parsed = JSON.parse(contractForm.value.abi)
    if (!Array.isArray(parsed)) return []
    return parsed
      .filter((item: { type?: string; name?: string }) => item.type === 'event' && item.name)
      .map((item: { name: string }) => item.name)
  } catch {
    return []
  }
})

const canSave = computed(() =>
  contractForm.value.name.trim().length > 0 &&
  contractForm.value.network !== '' &&
  /^0x[a-fA-F0-9]{40}$/.test(contractForm.value.address) &&
  !duplicateWarning.value
)

// --- Search debounce ---
watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchContracts()
  }, 300)
})

// --- Duplicate check ---
let dupTimeout: ReturnType<typeof setTimeout> | null = null
watch(
  () => [contractForm.value.address, contractForm.value.network],
  () => {
    duplicateWarning.value = ''
    if (dupTimeout) clearTimeout(dupTimeout)
    if (!contractForm.value.address || !contractForm.value.network) return
    if (!/^0x[a-fA-F0-9]{40}$/.test(contractForm.value.address)) return
    dupTimeout = setTimeout(async () => {
      try {
        const res = await api.get('/api/v1/whitelabel/contracts/check', {
          params: { address: contractForm.value.address, network: contractForm.value.network },
        })
        if (res.data.data?.exists && (!editingContract.value || editingContract.value.id !== res.data.data.id)) {
          duplicateWarning.value = 'This contract is already registered on this network.'
        }
      } catch { /* ignore */ }
    }, 500)
  },
)

// --- API ---
onMounted(async () => {
  await Promise.all([fetchStats(), fetchContracts()])
  loading.value = false
})

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/whitelabel/contracts/stats')
    if (res.data.data) stats.value = res.data.data
  } catch { /* defaults */ }
}

async function fetchContracts() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterTab.value !== 'all') params.status = filterTab.value
    if (searchQuery.value) params.search = searchQuery.value
    const res = await api.get('/api/v1/whitelabel/contracts', { params })
    contracts.value = (res.data.data?.items || []).map((c: ContractItem) => ({ ...c, expanded: false }))
    totalContracts.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

function openRegister() {
  editingContract.value = null
  contractForm.value = getEmptyForm()
  duplicateWarning.value = ''
  showRegister.value = true
}

function openEdit(contract: ContractItem) {
  editingContract.value = contract
  contractForm.value = {
    name: contract.name,
    network: contract.network,
    address: contract.address,
    abi: contract.abi || '',
    monitored_events: [...contract.monitored_events],
    manual_events: '',
  }
  duplicateWarning.value = ''
  showRegister.value = true
}

async function saveContract() {
  if (!canSave.value) return
  saving.value = true

  const events = contractForm.value.monitored_events.length > 0
    ? contractForm.value.monitored_events
    : contractForm.value.manual_events.split(',').map(e => e.trim()).filter(Boolean)

  const payload = {
    name: contractForm.value.name,
    network: contractForm.value.network,
    address: contractForm.value.address,
    abi: contractForm.value.abi || undefined,
    monitored_events: events,
  }

  try {
    if (editingContract.value) {
      await api.put(`/api/v1/whitelabel/contracts/${editingContract.value.id}`, {
        name: payload.name,
        monitored_events: payload.monitored_events,
      })
    } else {
      const res = await api.post('/api/v1/whitelabel/contracts', payload)
      // Auto-trigger verification after creation
      if (res.data.data?.id) {
        verifyContract(res.data.data.id)
      }
    }
    showRegister.value = false
    await Promise.all([fetchStats(), fetchContracts()])
  } catch { /* TODO: toast */ }
  finally { saving.value = false }
}

async function verifyContract(id: string) {
  verifying.value = id
  try {
    await api.post(`/api/v1/whitelabel/contracts/${id}/verify`)
    await Promise.all([fetchStats(), fetchContracts()])
  } catch { /* TODO: toast */ }
  finally { verifying.value = null }
}

async function deleteContract(contract: ContractItem) {
  const warn = contract.linked_rules > 0
    ? `This contract has ${contract.linked_rules} linked rule(s) that will be disabled. Delete anyway?`
    : `Delete contract "${contract.name}"?`
  if (!confirm(warn)) return
  try {
    await api.delete(`/api/v1/whitelabel/contracts/${contract.id}`)
    await Promise.all([fetchStats(), fetchContracts()])
  } catch { /* TODO: toast */ }
}

async function toggleExpand(contract: ContractItem) {
  contract.expanded = !contract.expanded
  if (contract.expanded && !contract.recent_events) {
    try {
      const res = await api.get(`/api/v1/whitelabel/contracts/${contract.id}/events`, {
        params: { limit: 10 },
      })
      contract.recent_events = res.data.data?.items || []
    } catch {
      contract.recent_events = []
    }
  }
}

function toggleEvent(eventName: string) {
  const idx = contractForm.value.monitored_events.indexOf(eventName)
  if (idx >= 0) {
    contractForm.value.monitored_events.splice(idx, 1)
  } else {
    contractForm.value.monitored_events.push(eventName)
  }
}

function getNetworkConfig(network: string) {
  return networks.find(n => n.value === network) || { label: network, color: '#94A3B8' }
}

function truncateAddress(addr: string): string {
  if (!addr || addr.length < 12) return addr
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`
}

function copyAddress(addr: string) {
  navigator.clipboard.writeText(addr)
}

function formatDate(d: string): string {
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatTxHash(hash: string): string {
  if (!hash || hash.length < 12) return hash
  return `${hash.slice(0, 10)}...${hash.slice(-6)}`
}

function onPageChange(p: number) {
  page.value = p
  fetchContracts()
}

function handleAbiUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = (e) => {
    contractForm.value.abi = (e.target?.result as string) || ''
  }
  reader.readAsText(file)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Contract Registry</h1>
        <p class="text-sm text-text-secondary">Register and monitor smart contracts for on-chain rule triggers</p>
      </div>
      <button
        class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors"
        @click="openRegister"
      >
        + Register Contract
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <StatsCard label="Total Contracts" :value="stats.total" icon="description" icon-color="#9B7EE0" />
      <StatsCard label="Verified" :value="stats.verified" icon="verified" icon-color="#16A34A" />
      <StatsCard label="Events Captured (24h)" :value="stats.events_24h" icon="bolt" icon-color="#F59E0B" />
      <StatsCard label="Linked Rules" :value="stats.linked_rules" icon="link" icon-color="#5D7EF1" />
    </div>

    <!-- Filter tabs + search -->
    <div class="flex items-center gap-4">
      <div class="flex gap-2">
        <button
          v-for="tab in ['all', 'verified', 'pending', 'error']"
          :key="tab"
          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
          :class="filterTab === tab ? 'bg-wl text-white' : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
          @click="filterTab = tab; page = 1; fetchContracts()"
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
            placeholder="Search contracts..."
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
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Contract Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Network</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-40">Address</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Status</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Events (24h)</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Linked Rules</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-32">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="7" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse" /></td>
          </tr>
          <tr v-else-if="contracts.length === 0">
            <td colspan="7" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">description</span>
              <p class="text-sm text-text-muted">No contracts registered</p>
              <button class="mt-2 text-xs text-wl hover:underline" @click="openRegister">+ Register First Contract</button>
            </td>
          </tr>
          <template v-else v-for="contract in contracts" :key="contract.id">
            <!-- Main row -->
            <tr
              class="border-b border-border last:border-b-0 hover:bg-white/[0.02] transition-colors"
              :class="contract.status === 'error' ? 'border-l-2 border-l-status-paused' : ''"
            >
              <td class="px-4 py-3">
                <div>
                  <span class="text-sm font-medium text-text-primary">{{ contract.name }}</span>
                  <p v-if="contract.status === 'error' && contract.error_message" class="text-xs text-status-paused mt-0.5">
                    {{ contract.error_message }}
                  </p>
                </div>
              </td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center gap-1.5 px-2 py-0.5 text-xs rounded-full font-medium"
                  :style="{ backgroundColor: getNetworkConfig(contract.network).color + '20', color: getNetworkConfig(contract.network).color }"
                >
                  <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: getNetworkConfig(contract.network).color }" />
                  {{ getNetworkConfig(contract.network).label }}
                </span>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-1.5">
                  <span class="text-sm font-mono text-text-secondary">{{ truncateAddress(contract.address) }}</span>
                  <button
                    class="p-0.5 rounded hover:bg-white/5 transition-colors"
                    title="Copy address"
                    @click="copyAddress(contract.address)"
                  >
                    <span class="material-symbols-rounded text-xs text-text-muted">content_copy</span>
                  </button>
                </div>
              </td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="{
                    'bg-status-active-bg text-status-active': contract.status === 'verified',
                    'bg-status-draft-bg text-status-draft': contract.status === 'pending',
                    'bg-[#2D1515] text-status-paused': contract.status === 'error',
                  }"
                >
                  {{ contract.status.charAt(0).toUpperCase() + contract.status.slice(1) }}
                </span>
              </td>
              <td class="px-4 py-3 text-right text-sm text-text-secondary">{{ contract.events_24h.toLocaleString() }}</td>
              <td class="px-4 py-3 text-right text-sm text-text-secondary">{{ contract.linked_rules }}</td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit" @click="openEdit(contract)">
                    <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                  </button>
                  <button
                    class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                    title="Verify"
                    :disabled="verifying === contract.id"
                    @click="verifyContract(contract.id)"
                  >
                    <span class="material-symbols-rounded text-base" :class="verifying === contract.id ? 'text-wl animate-spin' : 'text-text-muted'">
                      {{ verifying === contract.id ? 'progress_activity' : 'verified' }}
                    </span>
                  </button>
                  <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="View Events" @click="toggleExpand(contract)">
                    <span class="material-symbols-rounded text-base text-text-muted">
                      {{ contract.expanded ? 'expand_less' : 'expand_more' }}
                    </span>
                  </button>
                  <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteContract(contract)">
                    <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                  </button>
                </div>
              </td>
            </tr>
            <!-- Expanded events row -->
            <tr v-if="contract.expanded">
              <td colspan="7" class="bg-page-bg px-6 py-4">
                <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Recent Events</div>
                <div v-if="!contract.recent_events" class="text-sm text-text-muted">Loading...</div>
                <div v-else-if="contract.recent_events.length === 0" class="text-sm text-text-muted">No events captured yet</div>
                <table v-else class="w-full">
                  <thead>
                    <tr>
                      <th class="pb-2 text-left text-xs text-text-muted font-medium">Event</th>
                      <th class="pb-2 text-left text-xs text-text-muted font-medium">Block</th>
                      <th class="pb-2 text-left text-xs text-text-muted font-medium">Tx Hash</th>
                      <th class="pb-2 text-left text-xs text-text-muted font-medium">Time</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="evt in contract.recent_events" :key="evt.id" class="border-t border-border">
                      <td class="py-2 text-sm text-text-primary font-mono">{{ evt.event_name }}</td>
                      <td class="py-2 text-sm text-text-secondary">#{{ evt.block_number.toLocaleString() }}</td>
                      <td class="py-2 text-sm font-mono text-text-muted">{{ formatTxHash(evt.tx_hash) }}</td>
                      <td class="py-2 text-sm text-text-muted">{{ formatDate(evt.timestamp) }}</td>
                    </tr>
                  </tbody>
                </table>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
      <Pagination v-if="totalContracts > pageSize" :page="page" :page-size="pageSize" :total="totalContracts" @update:page="onPageChange" />
    </div>

    <!-- D12 Contract Register Modal -->
    <Modal :open="showRegister" :title="editingContract ? 'Edit Contract' : 'Register Contract'" max-width="640px" @close="showRegister = false">
      <div class="space-y-5">
        <!-- Contract Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1.5">Contract Name</label>
          <input
            v-model="contractForm.name"
            type="text"
            maxlength="50"
            placeholder="e.g. StakingPool V2"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-wl focus:outline-none"
          />
        </div>

        <!-- Network -->
        <div>
          <label class="block text-sm text-text-secondary mb-1.5">Network</label>
          <div class="grid grid-cols-4 gap-2">
            <button
              v-for="net in networks"
              :key="net.value"
              class="flex items-center gap-2 px-3 py-2 rounded-lg border text-xs font-medium transition-colors"
              :class="contractForm.network === net.value
                ? 'border-wl bg-wl/10 text-text-primary'
                : 'border-border text-text-secondary hover:text-text-primary'"
              :disabled="!!editingContract"
              @click="contractForm.network = net.value"
            >
              <span class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: net.color }" />
              {{ net.label }}
            </button>
          </div>
        </div>

        <!-- Contract Address -->
        <div>
          <label class="block text-sm text-text-secondary mb-1.5">Contract Address</label>
          <input
            v-model="contractForm.address"
            type="text"
            maxlength="42"
            placeholder="0x..."
            :disabled="!!editingContract"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm font-mono text-text-primary placeholder-text-muted focus:border-wl focus:outline-none disabled:opacity-50"
          />
          <p v-if="contractForm.address && !/^0x[a-fA-F0-9]{40}$/.test(contractForm.address)" class="text-xs text-status-paused mt-1">
            Invalid address format (must be 0x + 40 hex characters)
          </p>
          <p v-if="duplicateWarning" class="text-xs text-status-draft mt-1">
            {{ duplicateWarning }}
          </p>
        </div>

        <!-- ABI -->
        <div>
          <label class="block text-sm text-text-secondary mb-1.5">ABI (optional)</label>
          <div class="relative">
            <textarea
              v-model="contractForm.abi"
              rows="5"
              placeholder="Paste contract ABI JSON here..."
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm font-mono text-text-primary placeholder-text-muted focus:border-wl focus:outline-none resize-none"
            />
            <label class="absolute top-2 right-2 px-2 py-1 text-xs text-wl border border-wl/30 rounded hover:bg-wl/10 cursor-pointer transition-colors">
              Upload JSON
              <input type="file" accept=".json" class="hidden" @change="handleAbiUpload" />
            </label>
          </div>
          <p v-if="contractForm.abi && parsedAbiEvents.length > 0" class="text-xs text-status-active mt-1">
            {{ parsedAbiEvents.length }} events detected in ABI
          </p>
        </div>

        <!-- Monitored Events -->
        <div>
          <label class="block text-sm text-text-secondary mb-1.5">Monitored Events</label>
          <!-- From ABI -->
          <div v-if="parsedAbiEvents.length > 0" class="flex flex-wrap gap-2">
            <button
              v-for="evt in parsedAbiEvents"
              :key="evt"
              class="px-3 py-1.5 text-xs rounded-lg border transition-colors"
              :class="contractForm.monitored_events.includes(evt)
                ? 'border-wl bg-wl/10 text-wl'
                : 'border-border text-text-secondary hover:text-text-primary'"
              @click="toggleEvent(evt)"
            >
              {{ evt }}
            </button>
          </div>
          <!-- Manual input if no ABI -->
          <div v-else>
            <input
              v-model="contractForm.manual_events"
              type="text"
              placeholder="Comma-separated event names (e.g. Transfer, Approval, Stake)"
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-wl focus:outline-none"
            />
            <p class="text-xs text-text-muted mt-1">Provide ABI above to auto-detect events, or enter event names manually</p>
          </div>
        </div>
      </div>

      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showRegister = false">Cancel</button>
        <button
          class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 disabled:opacity-50 transition-colors"
          :disabled="!canSave || saving"
          @click="saveContract"
        >
          {{ saving ? 'Saving...' : editingContract ? 'Update Contract' : 'Register & Verify' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
