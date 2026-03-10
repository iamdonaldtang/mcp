<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api } from '../../api/client'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// === Types ===
interface BadgeItem {
  id: string
  name: string
  description: string
  icon: string
  category: 'achievement' | 'engagement' | 'special'
  condition_type: 'auto' | 'manual'
  auto_condition: string | null
  auto_value: number | null
  is_rare: boolean
  earned_count: number
  status: CampaignStatus
  created_at: string
}

interface BadgeStats {
  total_badges: number
  badges_earned: number
  unique_holders: number
  earn_rate: number
}

// === State ===
const loading = ref(true)
const badges = ref<BadgeItem[]>([])
const totalBadges = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<BadgeStats>({ total_badges: 0, badges_earned: 0, unique_holders: 0, earn_rate: 0 })

// Filters
const filterStatus = ref<string>('all')
const filterCategory = ref<string>('all')
const searchQuery = ref('')

const statusTabs = [
  { key: 'all', label: 'All' },
  { key: 'active', label: 'Active' },
  { key: 'draft', label: 'Draft' },
]

const categoryOptions = ['Achievement', 'Engagement', 'Special']

// Modal
const showCreateModal = ref(false)
const editingBadge = ref<BadgeItem | null>(null)
const saving = ref(false)

// Grant modal
const showGrantModal = ref(false)
const grantBadgeId = ref<string | null>(null)
const grantBadgeName = ref('')
const grantAddress = ref('')
const granting = ref(false)

// Icon picker
const presetIcons = [
  'emoji_events', 'military_tech', 'star', 'bolt', 'diamond',
  'favorite', 'rocket_launch', 'local_fire_department', 'psychology', 'diversity_3',
  'workspace_premium', 'eco', 'auto_awesome', 'celebration', 'shield',
  'verified', 'trophy', 'whatshot', 'grade', 'flash_on',
  'pets', 'cruelty_free', 'savings', 'palette', 'music_note',
  'sports_esports', 'code', 'terminal', 'token', 'currency_bitcoin',
  'public', 'language', 'explore', 'travel_explore', 'flight_takeoff',
  'handshake', 'volunteer_activism', 'loyalty', 'redeem', 'card_giftcard',
  'thumb_up', 'sentiment_very_satisfied', 'mood', 'self_improvement', 'spa',
  'hiking', 'sailing', 'surfing', 'skateboarding', 'snowboarding',
  'fitness_center', 'sports_martial_arts', 'sports_kabaddi', 'directions_run', 'sprint',
]

const autoConditions = [
  { key: 'complete_tasks', label: 'Complete N tasks' },
  { key: 'reach_level', label: 'Reach Level N' },
  { key: 'maintain_streak', label: 'Maintain N-day streak' },
  { key: 'earn_points', label: 'Earn N points' },
]

const formDefaults = {
  name: '',
  description: '',
  icon: 'emoji_events',
  category: 'achievement' as 'achievement' | 'engagement' | 'special',
  condition_type: 'auto' as 'auto' | 'manual',
  auto_condition: 'complete_tasks',
  auto_value: 10,
  is_rare: false,
}

const form = ref({ ...formDefaults })

// === Computed ===
const filteredBadges = computed(() => {
  return badges.value.filter(b => {
    if (filterStatus.value !== 'all' && b.status !== filterStatus.value) return false
    if (filterCategory.value !== 'all' && b.category !== filterCategory.value.toLowerCase()) return false
    if (searchQuery.value && !b.name.toLowerCase().includes(searchQuery.value.toLowerCase())) return false
    return true
  })
})

const canSave = computed(() => {
  return form.value.name.trim().length >= 1
    && form.value.name.trim().length <= 30
    && form.value.icon
    && (form.value.condition_type === 'manual' || (form.value.auto_value && form.value.auto_value > 0))
})

// === API ===
onMounted(async () => {
  await Promise.all([fetchBadges(), fetchStats()])
  loading.value = false
})

async function fetchBadges() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    if (filterCategory.value !== 'all') params.category = filterCategory.value.toLowerCase()
    const res = await api.get('/api/v1/community/modules/badges', { params })
    badges.value = res.data.data?.items || []
    totalBadges.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/badges/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

function openCreate() {
  editingBadge.value = null
  form.value = { ...formDefaults }
  showCreateModal.value = true
}

function openEdit(badge: BadgeItem) {
  editingBadge.value = badge
  form.value = {
    name: badge.name,
    description: badge.description,
    icon: badge.icon,
    category: badge.category,
    condition_type: badge.condition_type,
    auto_condition: badge.auto_condition || 'complete_tasks',
    auto_value: badge.auto_value || 10,
    is_rare: badge.is_rare,
  }
  showCreateModal.value = true
}

async function saveBadge() {
  if (!canSave.value) return
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      description: form.value.description.trim(),
      icon: form.value.icon,
      category: form.value.category,
      condition_type: form.value.condition_type,
      auto_condition: form.value.condition_type === 'auto' ? form.value.auto_condition : null,
      auto_value: form.value.condition_type === 'auto' ? form.value.auto_value : null,
      is_rare: form.value.is_rare,
    }

    if (editingBadge.value) {
      await api.put(`/api/v1/community/modules/badges/${editingBadge.value.id}`, payload)
    } else {
      await api.post('/api/v1/community/modules/badges', payload)
    }

    showCreateModal.value = false
    await Promise.all([fetchBadges(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

function openGrant(badge: BadgeItem) {
  grantBadgeId.value = badge.id
  grantBadgeName.value = badge.name
  grantAddress.value = ''
  showGrantModal.value = true
}

async function grantBadge() {
  if (!grantBadgeId.value || !grantAddress.value.trim()) return
  granting.value = true
  try {
    await api.post(`/api/v1/community/modules/badges/${grantBadgeId.value}/grant`, {
      wallet_address: grantAddress.value.trim(),
    })
    showGrantModal.value = false
    await Promise.all([fetchBadges(), fetchStats()])
  } catch { /* TODO: toast */ }
  granting.value = false
}

async function duplicateBadge(badge: BadgeItem) {
  try {
    await api.post(`/api/v1/community/modules/badges/${badge.id}/duplicate`)
    await Promise.all([fetchBadges(), fetchStats()])
  } catch { /* empty */ }
}

async function deleteBadge(id: string) {
  if (!confirm('Delete this badge? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/modules/badges/${id}`)
    await Promise.all([fetchBadges(), fetchStats()])
  } catch { /* empty */ }
}

function onPageChange(p: number) {
  page.value = p
  fetchBadges()
}

function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Badges</h1>
        <p class="text-sm text-text-secondary">Create and manage achievement badges for your community</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New Badge
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Total Badges</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.total_badges }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Badges Earned</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.badges_earned.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Unique Holders</div>
        <div class="text-2xl font-bold text-community">{{ stats.unique_holders.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Earn Rate</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.earn_rate.toFixed(1) }}%</div>
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
          @click="filterStatus = tab.key; fetchBadges()"
        >
          {{ tab.label }}
        </button>
      </div>
      <select
        v-model="filterCategory"
        class="px-3 py-1.5 bg-card-bg border border-border rounded-lg text-xs text-text-secondary focus:border-community focus:outline-none"
        @change="fetchBadges()"
      >
        <option value="all">All Categories</option>
        <option v-for="cat in categoryOptions" :key="cat" :value="cat">{{ cat }}</option>
      </select>
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search badges..."
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
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Badge Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Category</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Condition</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Earned</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Status</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-36">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="6" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse"></div></td>
          </tr>
          <tr v-else-if="filteredBadges.length === 0">
            <td colspan="6" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">military_tech</span>
              <p class="text-sm text-text-muted">No badges yet</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First Badge</button>
            </td>
          </tr>
          <tr
            v-else
            v-for="badge in filteredBadges"
            :key="badge.id"
            class="border-b border-border last:border-b-0 hover:bg-white/[0.02] transition-colors"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div
                  class="w-8 h-8 rounded-lg flex items-center justify-center"
                  :class="badge.is_rare ? 'bg-amber-900/30 ring-1 ring-amber-400/50' : 'bg-community/10'"
                >
                  <span
                    class="material-symbols-rounded text-base"
                    :class="badge.is_rare ? 'text-amber-400' : 'text-community'"
                  >{{ badge.icon }}</span>
                </div>
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="text-sm font-medium text-text-primary">{{ badge.name }}</span>
                    <span
                      v-if="badge.is_rare"
                      class="px-1.5 py-0.5 text-[10px] font-semibold uppercase rounded bg-amber-900/30 text-amber-400"
                    >Rare</span>
                  </div>
                  <div v-if="badge.description" class="text-xs text-text-muted truncate max-w-xs">{{ badge.description }}</div>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <span class="px-2 py-0.5 text-xs rounded bg-page-bg text-text-secondary capitalize">{{ badge.category }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-sm text-text-secondary capitalize">{{ badge.condition_type }}</span>
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ badge.earned_count.toLocaleString() }}</td>
            <td class="px-4 py-3">
              <StatusBadge :status="badge.status" />
            </td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit" @click="openEdit(badge)">
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button
                  v-if="badge.condition_type === 'manual'"
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Grant Badge"
                  @click="openGrant(badge)"
                >
                  <span class="material-symbols-rounded text-base text-community">card_giftcard</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Duplicate" @click="duplicateBadge(badge)">
                  <span class="material-symbols-rounded text-base text-text-muted">content_copy</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteBadge(badge.id)">
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination v-if="totalBadges > pageSize" :page="page" :page-size="pageSize" :total="totalBadges" @update:page="onPageChange" />
    </div>

    <!-- D09 Badge Editor Modal -->
    <Modal
      :open="showCreateModal"
      :title="editingBadge ? 'Edit Badge' : 'Create Badge'"
      max-width="560px"
      @close="showCreateModal = false"
    >
      <div class="space-y-4 max-h-[65vh] overflow-y-auto pr-1">
        <!-- Badge Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Badge Name <span class="text-status-paused">*</span></label>
          <input
            v-model="form.name"
            type="text"
            maxlength="30"
            placeholder="e.g. Early Adopter"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
          <div class="text-xs text-text-muted mt-1 text-right">{{ form.name.length }}/30</div>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Description</label>
          <textarea
            v-model="form.description"
            rows="2"
            maxlength="200"
            placeholder="What does this badge represent?"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm resize-none"
          />
          <div class="text-xs text-text-muted mt-1 text-right">{{ form.description.length }}/200</div>
        </div>

        <!-- Icon Picker -->
        <div>
          <label class="block text-sm text-text-secondary mb-2">Icon <span class="text-status-paused">*</span></label>
          <div class="bg-page-bg border border-border rounded-lg p-3">
            <!-- Selected preview -->
            <div class="flex items-center gap-3 mb-3 pb-3 border-b border-border">
              <div
                class="w-12 h-12 rounded-xl flex items-center justify-center"
                :class="form.is_rare ? 'bg-amber-900/30 ring-2 ring-amber-400/50' : 'bg-community/15'"
              >
                <span
                  class="material-symbols-rounded text-2xl"
                  :class="form.is_rare ? 'text-amber-400' : 'text-community'"
                >{{ form.icon }}</span>
              </div>
              <div>
                <div class="text-sm text-text-primary font-medium">{{ form.icon }}</div>
                <div class="text-xs text-text-muted">Selected icon</div>
              </div>
            </div>
            <!-- Grid -->
            <div class="grid grid-cols-10 gap-1.5 max-h-36 overflow-y-auto">
              <button
                v-for="icon in presetIcons"
                :key="icon"
                class="w-8 h-8 rounded-lg flex items-center justify-center transition-colors"
                :class="form.icon === icon
                  ? 'bg-community/20 ring-1 ring-community'
                  : 'hover:bg-white/5'"
                @click="form.icon = icon"
              >
                <span
                  class="material-symbols-rounded text-lg"
                  :class="form.icon === icon ? 'text-community' : 'text-text-muted'"
                >{{ icon }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Category -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Category</label>
          <select
            v-model="form.category"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          >
            <option value="achievement">Achievement</option>
            <option value="engagement">Engagement</option>
            <option value="special">Special</option>
          </select>
        </div>

        <!-- Earn Condition -->
        <div>
          <label class="block text-sm text-text-secondary mb-2">Earn Condition</label>
          <div class="space-y-2">
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.condition_type" type="radio" value="auto" class="accent-community" />
              <span class="text-sm text-text-primary">Auto-trigger</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.condition_type" type="radio" value="manual" class="accent-community" />
              <span class="text-sm text-text-primary">Manual Only</span>
            </label>
          </div>

          <!-- Auto condition config -->
          <div v-if="form.condition_type === 'auto'" class="mt-3 bg-page-bg border border-border rounded-lg p-3 space-y-3">
            <select
              v-model="form.auto_condition"
              class="w-full px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            >
              <option v-for="cond in autoConditions" :key="cond.key" :value="cond.key">{{ cond.label }}</option>
            </select>
            <div class="flex items-center gap-2">
              <span class="text-sm text-text-muted">N =</span>
              <input
                v-model.number="form.auto_value"
                type="number"
                min="1"
                class="w-24 px-3 py-2 bg-card-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
              />
              <span class="text-xs text-text-muted">
                {{ form.auto_condition === 'complete_tasks' ? 'tasks' :
                   form.auto_condition === 'reach_level' ? 'level' :
                   form.auto_condition === 'maintain_streak' ? 'days' : 'points' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Is Rare -->
        <div>
          <label class="flex items-center gap-3 cursor-pointer">
            <div class="relative">
              <input v-model="form.is_rare" type="checkbox" class="sr-only peer" />
              <div class="w-9 h-5 bg-border rounded-full peer-checked:bg-amber-500 transition-colors"></div>
              <div class="absolute left-0.5 top-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform peer-checked:translate-x-4"></div>
            </div>
            <div>
              <span class="text-sm text-text-primary">Rare Badge</span>
              <span class="text-xs text-text-muted ml-2">Adds gold glow visual indicator</span>
            </div>
          </label>
        </div>
      </div>

      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showCreateModal = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!canSave || saving"
          @click="saveBadge"
        >
          {{ saving ? 'Saving...' : editingBadge ? 'Update Badge' : 'Create Badge' }}
        </button>
      </template>
    </Modal>

    <!-- Grant Badge Modal -->
    <Modal
      :open="showGrantModal"
      title="Grant Badge"
      max-width="420px"
      @close="showGrantModal = false"
    >
      <div class="space-y-4">
        <p class="text-sm text-text-secondary">
          Manually grant <span class="font-medium text-community">{{ grantBadgeName }}</span> to a user.
        </p>
        <div>
          <label class="block text-sm text-text-secondary mb-1">Wallet Address <span class="text-status-paused">*</span></label>
          <input
            v-model="grantAddress"
            type="text"
            placeholder="0x..."
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm font-mono"
          />
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showGrantModal = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!grantAddress.trim() || granting"
          @click="grantBadge"
        >
          {{ granting ? 'Granting...' : 'Grant Badge' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
