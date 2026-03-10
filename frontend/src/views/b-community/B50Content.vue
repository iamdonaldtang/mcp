<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../../api/client'
import Modal from '../../components/common/Modal.vue'

// ─── Types ───
interface Announcement {
  id: string
  title: string
  content: string
  type: 'general' | 'event' | 'alert'
  image_url?: string
  cta_text?: string
  cta_url?: string
  pinned: boolean
  published_at: string
  scheduled_at?: string
}

interface FeaturedSlot {
  id: string
  slot_index: number
  content_type: 'quest' | 'lb_sprint' | 'milestone' | 'external_url'
  content_id?: string
  content_name?: string
  custom_title?: string
  custom_image_url?: string
  external_url?: string
  thumbnail_url?: string
}

interface ContentOption {
  id: string
  name: string
  status: string
}

interface ModuleStatus {
  key: string
  name: string
  icon: string
  active: boolean
  route: string
}

// ─── State ───
const loading = ref(true)
const announcements = ref<Announcement[]>([])
const featuredSlots = ref<(FeaturedSlot | null)[]>([null, null, null, null, null, null])
const moduleStatuses = ref<ModuleStatus[]>([])

// Announcement Modal (D16)
const showAnnouncementModal = ref(false)
const editingAnnouncement = ref<Announcement | null>(null)
const announcementForm = ref({
  title: '',
  content: '',
  type: 'general' as 'general' | 'event' | 'alert',
  image: null as File | null,
  cta_text: '',
  cta_url: '',
  schedule_mode: 'now' as 'now' | 'schedule',
  scheduled_at: '',
  pinned: false,
})

// Featured Slot Modal (D17)
const showFeaturedModal = ref(false)
const editingSlotIndex = ref<number | null>(null)
const featuredForm = ref({
  content_type: 'quest' as 'quest' | 'lb_sprint' | 'milestone' | 'external_url',
  content_id: '',
  custom_title: '',
  custom_image: null as File | null,
  external_url: '',
})
const contentOptions = ref<ContentOption[]>([])
const loadingContentOptions = ref(false)

// Confirm delete
const confirmDeleteAnnouncement = ref<string | null>(null)
const confirmRemoveSlot = ref<number | null>(null)

// ─── Data Fetching ───
onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([fetchAnnouncements(), fetchFeatured(), fetchModuleStatuses()])
  } finally {
    loading.value = false
  }
})

async function fetchAnnouncements() {
  try {
    const res = await api.get('/api/v1/community/content/announcements')
    announcements.value = res.data.data || []
  } catch { /* defaults */ }
}

async function fetchFeatured() {
  try {
    const res = await api.get('/api/v1/community/content/featured')
    const slots: FeaturedSlot[] = res.data.data || []
    const arr: (FeaturedSlot | null)[] = [null, null, null, null, null, null]
    slots.forEach(s => {
      if (s.slot_index >= 0 && s.slot_index < 6) arr[s.slot_index] = s
    })
    featuredSlots.value = arr
  } catch { /* defaults */ }
}

async function fetchModuleStatuses() {
  try {
    const res = await api.get('/api/v1/community/overview')
    const enabled: string[] = res.data.data?.enabled_modules || []
    const modules: { key: string; name: string; icon: string; route: string }[] = [
      { key: 'points_level', name: 'Points & Level', icon: 'stars', route: '/b/community/modules/points-level' },
      { key: 'leaderboard', name: 'Leaderboard', icon: 'leaderboard', route: '/b/community/modules/leaderboard' },
      { key: 'task_chain', name: 'TaskChain', icon: 'link', route: '/b/community/modules/task-chain' },
      { key: 'day_chain', name: 'DayChain', icon: 'local_fire_department', route: '/b/community/modules/day-chain' },
      { key: 'benefits_shop', name: 'Shop', icon: 'storefront', route: '/b/community/modules/benefits-shop' },
      { key: 'lucky_wheel', name: 'Lucky Wheel', icon: 'casino', route: '/b/community/modules/lucky-wheel' },
    ]
    moduleStatuses.value = modules.map(m => ({
      ...m,
      active: enabled.includes(m.key),
    }))
  } catch { /* defaults */ }
}

// ─── Announcement Actions ───
function openNewAnnouncement() {
  editingAnnouncement.value = null
  announcementForm.value = {
    title: '', content: '', type: 'general', image: null,
    cta_text: '', cta_url: '', schedule_mode: 'now', scheduled_at: '', pinned: false,
  }
  showAnnouncementModal.value = true
}

function openEditAnnouncement(a: Announcement) {
  editingAnnouncement.value = a
  announcementForm.value = {
    title: a.title,
    content: a.content,
    type: a.type,
    image: null,
    cta_text: a.cta_text || '',
    cta_url: a.cta_url || '',
    schedule_mode: a.scheduled_at ? 'schedule' : 'now',
    scheduled_at: a.scheduled_at || '',
    pinned: a.pinned,
  }
  showAnnouncementModal.value = true
}

async function saveAnnouncement() {
  if (!announcementForm.value.title.trim()) return
  const payload: Record<string, unknown> = {
    title: announcementForm.value.title,
    content: announcementForm.value.content,
    type: announcementForm.value.type,
    pinned: announcementForm.value.pinned,
  }
  if (announcementForm.value.cta_text) {
    payload.cta_text = announcementForm.value.cta_text
    payload.cta_url = announcementForm.value.cta_url
  }
  if (announcementForm.value.schedule_mode === 'schedule' && announcementForm.value.scheduled_at) {
    payload.scheduled_at = announcementForm.value.scheduled_at
  }
  // Image upload would use FormData in production
  if (announcementForm.value.image) {
    payload.has_image = true
  }

  try {
    if (editingAnnouncement.value) {
      await api.put(`/api/v1/community/content/announcements/${editingAnnouncement.value.id}`, payload)
    } else {
      await api.post('/api/v1/community/content/announcements', payload)
    }
    showAnnouncementModal.value = false
    await fetchAnnouncements()
  } catch { /* TODO: toast */ }
}

async function togglePin(a: Announcement) {
  try {
    await api.put(`/api/v1/community/content/announcements/${a.id}`, { pinned: !a.pinned })
    await fetchAnnouncements()
  } catch { /* TODO: toast */ }
}

async function deleteAnnouncement(id: string) {
  try {
    await api.delete(`/api/v1/community/content/announcements/${id}`)
    confirmDeleteAnnouncement.value = null
    await fetchAnnouncements()
  } catch { /* TODO: toast */ }
}

function onImageSelected(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files?.[0]) {
    announcementForm.value.image = input.files[0]
  }
}

// ─── Featured Slot Actions ───
function openFeaturedModal(index: number) {
  editingSlotIndex.value = index
  const existing = featuredSlots.value[index]
  if (existing) {
    featuredForm.value = {
      content_type: existing.content_type,
      content_id: existing.content_id || '',
      custom_title: existing.custom_title || '',
      custom_image: null,
      external_url: existing.external_url || '',
    }
  } else {
    featuredForm.value = {
      content_type: 'quest', content_id: '', custom_title: '', custom_image: null, external_url: '',
    }
  }
  showFeaturedModal.value = true
  fetchContentOptions(featuredForm.value.content_type)
}

async function fetchContentOptions(type: string) {
  if (type === 'external_url') {
    contentOptions.value = []
    return
  }
  loadingContentOptions.value = true
  try {
    const endpoint = type === 'quest'
      ? '/api/v1/community/content/options/quests'
      : type === 'lb_sprint'
        ? '/api/v1/community/content/options/lb-sprints'
        : '/api/v1/community/content/options/milestones'
    const res = await api.get(endpoint)
    contentOptions.value = res.data.data || []
  } catch {
    contentOptions.value = []
  } finally {
    loadingContentOptions.value = false
  }
}

function onContentTypeChange() {
  featuredForm.value.content_id = ''
  featuredForm.value.external_url = ''
  fetchContentOptions(featuredForm.value.content_type)
}

function onFeaturedImageSelected(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files?.[0]) {
    featuredForm.value.custom_image = input.files[0]
  }
}

async function saveFeatured() {
  if (editingSlotIndex.value === null) return
  const f = featuredForm.value
  if (f.content_type === 'external_url' && !f.external_url.trim()) return
  if (f.content_type !== 'external_url' && !f.content_id) return

  const payload: Record<string, unknown> = {
    slot_index: editingSlotIndex.value,
    content_type: f.content_type,
  }
  if (f.content_type === 'external_url') {
    payload.external_url = f.external_url
  } else {
    payload.content_id = f.content_id
  }
  if (f.custom_title) payload.custom_title = f.custom_title

  try {
    const existing = featuredSlots.value[editingSlotIndex.value]
    if (existing) {
      await api.put(`/api/v1/community/content/featured/${existing.id}`, payload)
    } else {
      await api.post('/api/v1/community/content/featured', payload)
    }
    showFeaturedModal.value = false
    await fetchFeatured()
  } catch { /* TODO: toast */ }
}

async function removeFeaturedSlot(index: number) {
  const slot = featuredSlots.value[index]
  if (!slot) return
  try {
    await api.delete(`/api/v1/community/content/featured/${slot.id}`)
    confirmRemoveSlot.value = null
    await fetchFeatured()
  } catch { /* TODO: toast */ }
}

// ─── Helpers ───
const typeBadgeStyles: Record<string, { bg: string; text: string }> = {
  general: { bg: '#172554', text: '#60A5FA' },
  event: { bg: '#1E1033', text: '#A78BFA' },
  alert: { bg: '#2D1515', text: '#F87171' },
}

const contentTypeIcons: Record<string, string> = {
  quest: 'rocket_launch',
  lb_sprint: 'sprint',
  milestone: 'flag',
  external_url: 'link',
}

function formatDate(dateStr: string) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div>
      <h1 class="text-2xl font-bold text-text-primary mb-1">Content Management</h1>
      <p class="text-sm text-text-secondary">Manage announcements, featured content, and module visibility for your community</p>
    </div>

    <!-- ═══ ANNOUNCEMENTS ═══ -->
    <div>
      <div class="flex items-center justify-between mb-4">
        <div class="text-xs font-semibold text-text-muted uppercase tracking-wider">Announcements</div>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
          @click="openNewAnnouncement"
        >
          + New Announcement
        </button>
      </div>

      <!-- Loading skeleton -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 2" :key="i" class="bg-card-bg border border-border rounded-xl p-5 h-24 animate-pulse" />
      </div>

      <!-- Empty state -->
      <div v-else-if="announcements.length === 0" class="bg-card-bg border border-border rounded-xl p-8 text-center">
        <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">campaign</span>
        <p class="text-sm text-text-muted mb-2">No announcements yet</p>
        <button class="text-xs text-community hover:underline" @click="openNewAnnouncement">+ Create your first announcement</button>
      </div>

      <!-- Announcement cards -->
      <div v-else class="space-y-3">
        <div
          v-for="a in announcements"
          :key="a.id"
          class="bg-card-bg border border-border rounded-xl p-5 hover:bg-[#161F2E] transition-colors"
          :class="a.pinned ? 'border-community/30' : ''"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1.5">
                <span v-if="a.pinned" class="material-symbols-rounded text-base text-community" title="Pinned">push_pin</span>
                <h3 class="text-sm font-semibold text-text-primary truncate">{{ a.title }}</h3>
                <span
                  class="px-2 py-0.5 text-[10px] font-medium rounded-full uppercase tracking-wide"
                  :style="{ background: typeBadgeStyles[a.type]?.bg, color: typeBadgeStyles[a.type]?.text }"
                >
                  {{ a.type }}
                </span>
              </div>
              <p class="text-xs text-text-secondary line-clamp-2 mb-2">{{ a.content }}</p>
              <div class="flex items-center gap-3 text-[11px] text-text-muted">
                <span>{{ formatDate(a.published_at) }}</span>
                <span v-if="a.cta_text" class="flex items-center gap-1">
                  <span class="material-symbols-rounded text-xs">link</span>
                  {{ a.cta_text }}
                </span>
              </div>
            </div>
            <!-- Thumbnail -->
            <img v-if="a.image_url" :src="a.image_url" alt="" class="w-16 h-16 rounded-lg object-cover flex-shrink-0" />
          </div>
          <!-- Actions -->
          <div class="flex items-center gap-1 mt-3 pt-3 border-t border-border">
            <button
              class="flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs text-text-muted hover:text-text-primary hover:bg-white/5 transition-colors"
              @click="togglePin(a)"
            >
              <span class="material-symbols-rounded text-sm">{{ a.pinned ? 'push_pin' : 'push_pin' }}</span>
              {{ a.pinned ? 'Unpin' : 'Pin' }}
            </button>
            <button
              class="flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs text-text-muted hover:text-text-primary hover:bg-white/5 transition-colors"
              @click="openEditAnnouncement(a)"
            >
              <span class="material-symbols-rounded text-sm">edit</span>
              Edit
            </button>
            <button
              v-if="confirmDeleteAnnouncement !== a.id"
              class="flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs text-status-paused hover:bg-white/5 transition-colors"
              @click="confirmDeleteAnnouncement = a.id"
            >
              <span class="material-symbols-rounded text-sm">delete</span>
              Delete
            </button>
            <div v-else class="flex items-center gap-1">
              <span class="text-xs text-status-paused">Confirm?</span>
              <button class="px-2 py-0.5 text-xs text-status-paused hover:bg-white/5 rounded" @click="deleteAnnouncement(a.id)">Yes</button>
              <button class="px-2 py-0.5 text-xs text-text-muted hover:bg-white/5 rounded" @click="confirmDeleteAnnouncement = null">No</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ═══ FEATURED SLOTS ═══ -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Featured Slots</div>
      <div class="grid grid-cols-3 gap-4">
        <div
          v-for="(slot, idx) in featuredSlots"
          :key="idx"
          class="bg-card-bg border border-border rounded-xl overflow-hidden cursor-pointer hover:border-community/30 transition-colors"
          :class="slot ? '' : 'border-dashed'"
          @click="openFeaturedModal(idx)"
        >
          <!-- Filled slot -->
          <template v-if="slot">
            <div class="relative">
              <img v-if="slot.thumbnail_url || slot.custom_image_url" :src="slot.custom_image_url || slot.thumbnail_url" alt="" class="w-full h-28 object-cover" />
              <div v-else class="w-full h-28 bg-page-bg flex items-center justify-center">
                <span class="material-symbols-rounded text-3xl text-text-muted">{{ contentTypeIcons[slot.content_type] || 'article' }}</span>
              </div>
              <!-- Remove button -->
              <button
                v-if="confirmRemoveSlot !== idx"
                class="absolute top-2 right-2 w-6 h-6 rounded-full bg-black/60 flex items-center justify-center hover:bg-black/80 transition-colors"
                @click.stop="confirmRemoveSlot = idx"
              >
                <span class="material-symbols-rounded text-sm text-white">close</span>
              </button>
              <div v-else class="absolute top-2 right-2 flex gap-1" @click.stop>
                <button class="px-2 py-0.5 bg-status-paused text-white text-[10px] rounded" @click="removeFeaturedSlot(idx)">Remove</button>
                <button class="px-2 py-0.5 bg-black/60 text-white text-[10px] rounded" @click="confirmRemoveSlot = null">Cancel</button>
              </div>
            </div>
            <div class="p-3">
              <div class="flex items-center gap-1.5 mb-1">
                <span class="material-symbols-rounded text-xs text-text-muted">{{ contentTypeIcons[slot.content_type] || 'article' }}</span>
                <span class="text-[10px] text-text-muted uppercase">{{ slot.content_type.replace('_', ' ') }}</span>
              </div>
              <p class="text-sm font-medium text-text-primary truncate">{{ slot.custom_title || slot.content_name || 'Untitled' }}</p>
            </div>
          </template>

          <!-- Empty slot -->
          <template v-else>
            <div class="h-[172px] flex flex-col items-center justify-center gap-2">
              <span class="material-symbols-rounded text-2xl text-text-muted">add</span>
              <span class="text-xs text-text-muted">Add Featured</span>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- ═══ MODULE STATUS ═══ -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Module Status</div>
      <div class="grid grid-cols-3 gap-4">
        <template v-if="loading">
          <div v-for="i in 6" :key="i" class="bg-card-bg border border-border rounded-xl p-4 h-20 animate-pulse" />
        </template>
        <div
          v-else
          v-for="mod in moduleStatuses"
          :key="mod.key"
          class="bg-card-bg border border-border rounded-xl p-4 flex items-center justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-9 h-9 rounded-lg flex items-center justify-center"
              :style="{ background: mod.active ? '#0A2E1A' : '#1E293B' }"
            >
              <span
                class="material-symbols-rounded text-lg"
                :style="{ color: mod.active ? '#48BB78' : '#64748B' }"
              >{{ mod.icon }}</span>
            </div>
            <div>
              <span class="text-sm font-medium text-text-primary block">{{ mod.name }}</span>
              <span
                class="text-[11px] font-medium"
                :style="{ color: mod.active ? '#16A34A' : '#64748B' }"
              >
                {{ mod.active ? 'Active' : 'Not Configured' }}
              </span>
            </div>
          </div>
          <a
            :href="mod.route"
            class="text-xs text-community font-medium hover:underline"
            @click.prevent="$router.push(mod.route)"
          >
            Configure →
          </a>
        </div>
      </div>
    </div>

    <!-- ═══ D16 — Announcement Modal ═══ -->
    <Modal
      :open="showAnnouncementModal"
      :title="editingAnnouncement ? 'Edit Announcement' : 'New Announcement'"
      max-width="560px"
      @close="showAnnouncementModal = false"
    >
      <div class="space-y-4">
        <!-- Title -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Title</label>
          <input
            v-model="announcementForm.title"
            type="text"
            maxlength="80"
            placeholder="Announcement title"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
          <span class="text-[11px] text-text-muted mt-0.5 block text-right">{{ announcementForm.title.length }}/80</span>
        </div>

        <!-- Content -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Content</label>
          <textarea
            v-model="announcementForm.content"
            rows="4"
            maxlength="500"
            placeholder="Announcement content... URLs will be auto-linked."
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm resize-none"
          />
          <span class="text-[11px] text-text-muted mt-0.5 block text-right">{{ announcementForm.content.length }}/500</span>
        </div>

        <!-- Type -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Type</label>
          <select
            v-model="announcementForm.type"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          >
            <option value="general">General</option>
            <option value="event">Event</option>
            <option value="alert">Alert</option>
          </select>
        </div>

        <!-- Image (optional) -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Image <span class="text-text-muted">(optional)</span></label>
          <div class="flex items-center gap-3">
            <label class="px-4 py-2 bg-page-bg border border-border rounded-lg text-sm text-text-secondary cursor-pointer hover:border-community/30 transition-colors">
              <span class="material-symbols-rounded text-sm align-middle mr-1">upload</span>
              {{ announcementForm.image ? announcementForm.image.name : 'Choose file' }}
              <input type="file" accept="image/*" class="hidden" @change="onImageSelected" />
            </label>
            <button v-if="announcementForm.image" class="text-xs text-text-muted hover:text-text-primary" @click="announcementForm.image = null">Remove</button>
          </div>
        </div>

        <!-- CTA Button (optional) -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">CTA Button <span class="text-text-muted">(optional)</span></label>
          <div class="grid grid-cols-2 gap-3">
            <input
              v-model="announcementForm.cta_text"
              type="text"
              placeholder="Button text"
              class="px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
            />
            <input
              v-model="announcementForm.cta_url"
              type="url"
              placeholder="https://..."
              class="px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
            />
          </div>
        </div>

        <!-- Schedule -->
        <div>
          <label class="block text-sm text-text-secondary mb-2">Schedule</label>
          <div class="flex gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" value="now" v-model="announcementForm.schedule_mode" class="accent-[#48BB78]" />
              <span class="text-sm text-text-primary">Publish Now</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" value="schedule" v-model="announcementForm.schedule_mode" class="accent-[#48BB78]" />
              <span class="text-sm text-text-primary">Schedule</span>
            </label>
          </div>
          <input
            v-if="announcementForm.schedule_mode === 'schedule'"
            v-model="announcementForm.scheduled_at"
            type="datetime-local"
            class="mt-2 w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          />
        </div>

        <!-- Pin toggle -->
        <div class="flex items-center justify-between">
          <label class="text-sm text-text-secondary">Pin this announcement</label>
          <button
            class="relative w-10 h-5 rounded-full transition-colors"
            :class="announcementForm.pinned ? 'bg-community' : 'bg-border'"
            @click="announcementForm.pinned = !announcementForm.pinned"
          >
            <span
              class="absolute top-0.5 w-4 h-4 bg-white rounded-full transition-transform"
              :class="announcementForm.pinned ? 'left-[22px]' : 'left-0.5'"
            />
          </button>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showAnnouncementModal = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50"
          :disabled="!announcementForm.title.trim()"
          @click="saveAnnouncement"
        >
          {{ editingAnnouncement ? 'Save Changes' : 'Publish' }}
        </button>
      </template>
    </Modal>

    <!-- ═══ D17 — Featured Slot Modal ═══ -->
    <Modal
      :open="showFeaturedModal"
      title="Featured Content"
      max-width="480px"
      @close="showFeaturedModal = false"
    >
      <div class="space-y-4">
        <!-- Content Type -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Content Type</label>
          <select
            v-model="featuredForm.content_type"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            @change="onContentTypeChange"
          >
            <option value="quest">Quest</option>
            <option value="lb_sprint">LB Sprint</option>
            <option value="milestone">Milestone</option>
            <option value="external_url">External URL</option>
          </select>
        </div>

        <!-- Content Select (for non-URL types) -->
        <div v-if="featuredForm.content_type !== 'external_url'">
          <label class="block text-sm text-text-secondary mb-1">Select Content</label>
          <select
            v-model="featuredForm.content_id"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            :disabled="loadingContentOptions"
          >
            <option value="" disabled>{{ loadingContentOptions ? 'Loading...' : 'Select content' }}</option>
            <option v-for="opt in contentOptions" :key="opt.id" :value="opt.id">
              {{ opt.name }} · {{ opt.status }}
            </option>
          </select>
        </div>

        <!-- External URL (for URL type) -->
        <div v-else>
          <label class="block text-sm text-text-secondary mb-1">URL</label>
          <input
            v-model="featuredForm.external_url"
            type="url"
            placeholder="https://..."
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
        </div>

        <!-- Custom Title (optional) -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Custom Title <span class="text-text-muted">(optional, 60 chars)</span></label>
          <input
            v-model="featuredForm.custom_title"
            type="text"
            maxlength="60"
            placeholder="Override display title"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
        </div>

        <!-- Custom Image (optional) -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Custom Image <span class="text-text-muted">(optional)</span></label>
          <label class="px-4 py-2 bg-page-bg border border-border rounded-lg text-sm text-text-secondary cursor-pointer hover:border-community/30 transition-colors inline-block">
            <span class="material-symbols-rounded text-sm align-middle mr-1">upload</span>
            {{ featuredForm.custom_image ? featuredForm.custom_image.name : 'Choose file' }}
            <input type="file" accept="image/*" class="hidden" @change="onFeaturedImageSelected" />
          </label>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showFeaturedModal = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50"
          :disabled="featuredForm.content_type === 'external_url' ? !featuredForm.external_url.trim() : !featuredForm.content_id"
          @click="saveFeatured"
        >
          Save
        </button>
      </template>
    </Modal>
  </div>
</template>
