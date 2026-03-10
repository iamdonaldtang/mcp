<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { api } from '../../api/client'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'

// --- Types ---
type SectionType = 'banner' | 'quest_widget' | 'leaderboard_widget' | 'points_widget' | 'rich_text' | 'custom_html'
type VisibilityScope = 'all' | 'logged_in' | 'level_gated'

interface HomepageSection {
  id: string
  name: string
  type: SectionType
  visible: boolean
  sort_order: number
  visibility_scope: VisibilityScope
  min_level?: number
  config: Record<string, unknown>
  last_modified: string
}

interface StatsData {
  total_sections: number
  visible_sections: number
  page_views_24h: number
  avg_session_duration: string
}

// --- State ---
const loading = ref(true)
const items = ref<HomepageSection[]>([])
const totalItems = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<StatsData>({ total_sections: 0, visible_sections: 0, page_views_24h: 0, avg_session_duration: '0m 0s' })

// Filters
const filterVisibility = ref<string>('all')
const visibilityTabs = ['all', 'visible', 'hidden'] as const

// Drag state
const dragIndex = ref<number | null>(null)
const dropIndex = ref<number | null>(null)

// Modal
const showCreate = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  type: 'banner' as SectionType,
  visibility_scope: 'all' as VisibilityScope,
  min_level: 1,
  // Banner
  banner_image: '',
  banner_link: '',
  // Widget
  widget_instance_id: '',
  // Rich Text
  rich_text_content: '',
  // Custom HTML
  custom_html: '',
})

const sectionTypeLabels: Record<SectionType, string> = {
  banner: 'Banner',
  quest_widget: 'Quest Widget',
  leaderboard_widget: 'Leaderboard Widget',
  points_widget: 'Points Widget',
  rich_text: 'Rich Text',
  custom_html: 'Custom HTML',
}

const sectionTypeIcons: Record<SectionType, string> = {
  banner: 'image',
  quest_widget: 'task_alt',
  leaderboard_widget: 'leaderboard',
  points_widget: 'stars',
  rich_text: 'article',
  custom_html: 'code',
}

const formValid = computed(() => {
  if (form.value.name.trim().length < 1 || form.value.name.trim().length > 60) return false
  if (form.value.visibility_scope === 'level_gated' && form.value.min_level < 1) return false
  switch (form.value.type) {
    case 'banner':
      return form.value.banner_image.trim().length > 0
    case 'quest_widget':
    case 'leaderboard_widget':
    case 'points_widget':
      return form.value.widget_instance_id !== ''
    case 'rich_text':
      return form.value.rich_text_content.trim().length > 0
    case 'custom_html':
      return form.value.custom_html.trim().length > 0
  }
  return false
})

// Widget instances for dropdowns
const widgetInstances = ref<{ id: string; name: string; type: string }[]>([])

// --- Fetch ---
onMounted(async () => {
  await Promise.all([fetchStats(), fetchItems(), fetchWidgetInstances()])
  loading.value = false
})

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/settings/homepage/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

async function fetchItems() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterVisibility.value === 'visible') params.visible = 1
    if (filterVisibility.value === 'hidden') params.visible = 0
    const res = await api.get('/api/v1/community/settings/homepage/sections', { params })
    items.value = res.data.data?.items || []
    totalItems.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

async function fetchWidgetInstances() {
  try {
    const res = await api.get('/api/v1/community/settings/homepage/widget-instances')
    widgetInstances.value = res.data.data || []
  } catch { /* empty */ }
}

watch(filterVisibility, () => { page.value = 1; fetchItems() })

function onPageChange(p: number) {
  page.value = p
  fetchItems()
}

// --- Drag & Drop ---
function onDragStart(index: number) {
  dragIndex.value = index
}

function onDragOver(e: DragEvent, index: number) {
  e.preventDefault()
  dropIndex.value = index
}

function onDragEnd() {
  dragIndex.value = null
  dropIndex.value = null
}

async function onDrop(targetIndex: number) {
  if (dragIndex.value === null || dragIndex.value === targetIndex) {
    onDragEnd()
    return
  }
  const moved = items.value.splice(dragIndex.value, 1)[0]
  items.value.splice(targetIndex, 0, moved)
  items.value.forEach((item, i) => { item.sort_order = i + 1 })
  onDragEnd()
  try {
    await api.put('/api/v1/community/settings/homepage/reorder', {
      order: items.value.map(i => i.id),
    })
  } catch {
    await fetchItems()
  }
}

// --- Visibility Toggle ---
async function toggleVisibility(item: HomepageSection) {
  const newVisible = !item.visible
  item.visible = newVisible // optimistic
  try {
    await api.put(`/api/v1/community/settings/homepage/sections/${item.id}`, { visible: newVisible })
    await fetchStats()
  } catch {
    item.visible = !newVisible // revert
  }
}

// --- Actions ---
function openCreate() {
  editingId.value = null
  form.value = {
    name: '', type: 'banner', visibility_scope: 'all', min_level: 1,
    banner_image: '', banner_link: '', widget_instance_id: '',
    rich_text_content: '', custom_html: '',
  }
  showCreate.value = true
}

function openEdit(item: HomepageSection) {
  editingId.value = item.id
  form.value = {
    name: item.name,
    type: item.type,
    visibility_scope: item.visibility_scope,
    min_level: item.min_level || 1,
    banner_image: (item.config.banner_image as string) || '',
    banner_link: (item.config.banner_link as string) || '',
    widget_instance_id: (item.config.widget_instance_id as string) || '',
    rich_text_content: (item.config.rich_text_content as string) || '',
    custom_html: (item.config.custom_html as string) || '',
  }
  showCreate.value = true
}

function buildPayload() {
  const base = {
    name: form.value.name,
    type: form.value.type,
    visibility_scope: form.value.visibility_scope,
    min_level: form.value.visibility_scope === 'level_gated' ? form.value.min_level : undefined,
  }
  let config: Record<string, unknown> = {}
  switch (form.value.type) {
    case 'banner':
      config = { banner_image: form.value.banner_image, banner_link: form.value.banner_link }
      break
    case 'quest_widget':
    case 'leaderboard_widget':
    case 'points_widget':
      config = { widget_instance_id: form.value.widget_instance_id }
      break
    case 'rich_text':
      config = { rich_text_content: form.value.rich_text_content }
      break
    case 'custom_html':
      config = { custom_html: form.value.custom_html }
      break
  }
  return { ...base, config }
}

async function saveForm() {
  if (!formValid.value || saving.value) return
  saving.value = true
  try {
    const payload = buildPayload()
    if (editingId.value) {
      await api.put(`/api/v1/community/settings/homepage/sections/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/community/settings/homepage/sections', payload)
    }
    showCreate.value = false
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function deleteItem(item: HomepageSection) {
  if (!confirm(`Delete section "${item.name}"? This action cannot be undone.`)) return
  try {
    await api.delete(`/api/v1/community/settings/homepage/sections/${item.id}`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

function openPreview() {
  window.open('/b/community/preview', '_blank')
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function formatNumber(n: number) {
  return n.toLocaleString()
}

const filteredWidgets = computed(() => {
  if (form.value.type === 'quest_widget') return widgetInstances.value.filter(w => w.type === 'quest')
  if (form.value.type === 'leaderboard_widget') return widgetInstances.value.filter(w => w.type === 'leaderboard')
  if (form.value.type === 'points_widget') return widgetInstances.value.filter(w => w.type === 'points')
  return widgetInstances.value
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[#F1F5F9] mb-1">Homepage Editor</h1>
        <p class="text-sm text-[#94A3B8]">Customize the sections and layout of your community homepage</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          class="px-4 py-2 text-sm font-medium text-[#94A3B8] border border-[#1E293B] rounded-lg hover:text-[#F1F5F9] hover:border-[#94A3B8] transition-colors"
          @click="openPreview"
        >
          <span class="material-symbols-rounded text-base align-middle mr-1">open_in_new</span>
          Preview
        </button>
        <button
          class="px-4 py-2 bg-[#48BB78] text-black text-sm font-medium rounded-lg hover:bg-[#48BB78]/90 transition-colors"
          @click="openCreate"
        >
          + Add Section
        </button>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Total Sections</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.total_sections) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Visible Sections</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.visible_sections) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Page Views (24h)</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ formatNumber(stats.page_views_24h) }}</p>
      </div>
      <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4">
        <p class="text-xs text-[#64748B] mb-1">Avg Session Duration</p>
        <p class="text-2xl font-bold text-[#F1F5F9]">{{ stats.avg_session_duration }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="flex border-b border-[#1E293B]">
        <button
          v-for="tab in visibilityTabs"
          :key="tab"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="filterVisibility === tab
            ? 'text-[#48BB78]'
            : 'text-[#94A3B8] hover:text-[#F1F5F9]'"
          @click="filterVisibility = tab"
        >
          {{ tab === 'all' ? 'All' : tab.charAt(0).toUpperCase() + tab.slice(1) }}
          <span
            v-if="filterVisibility === tab"
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
            <th class="px-2 py-3 w-10"></th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider">Section Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-40">Type</th>
            <th class="px-4 py-3 text-center text-xs font-semibold text-[#64748B] uppercase tracking-wider w-28">Visibility</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-[#64748B] uppercase tracking-wider w-36">Last Modified</th>
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
              <span class="material-symbols-rounded text-4xl text-[#64748B] block mb-2">dashboard_customize</span>
              <p class="text-sm text-[#64748B]">No sections added yet</p>
              <p class="text-xs text-[#64748B] mt-1">Add sections to customize your community homepage</p>
              <button class="mt-3 text-xs text-[#48BB78] hover:underline" @click="openCreate">+ Add First Section</button>
            </td>
          </tr>
          <!-- Data rows -->
          <tr
            v-else
            v-for="(item, index) in items"
            :key="item.id"
            class="border-b border-[#1E293B] last:border-b-0 hover:bg-white/[0.02] transition-colors"
            :class="{
              'bg-[#48BB78]/5': dropIndex === index && dragIndex !== index,
              'opacity-50': !item.visible,
            }"
            draggable="true"
            @dragstart="onDragStart(index)"
            @dragover="onDragOver($event, index)"
            @drop="onDrop(index)"
            @dragend="onDragEnd"
          >
            <!-- Drag handle -->
            <td class="px-2 py-3 text-center cursor-grab active:cursor-grabbing">
              <span class="material-symbols-rounded text-base text-[#64748B]">drag_indicator</span>
            </td>
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-[#48BB78]/10 flex items-center justify-center">
                  <span class="material-symbols-rounded text-base text-[#48BB78]">{{ sectionTypeIcons[item.type] }}</span>
                </div>
                <div>
                  <span class="text-sm font-medium text-[#F1F5F9]">{{ item.name }}</span>
                  <p v-if="item.visibility_scope === 'level_gated'" class="text-xs text-[#64748B] mt-0.5">
                    Level {{ item.min_level }}+ only
                  </p>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <span class="inline-flex items-center gap-1.5 px-2 py-1 bg-[#1E293B] rounded text-xs text-[#94A3B8]">
                <span class="material-symbols-rounded text-xs">{{ sectionTypeIcons[item.type] }}</span>
                {{ sectionTypeLabels[item.type] }}
              </span>
            </td>
            <td class="px-4 py-3 text-center">
              <!-- Toggle switch -->
              <button
                class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors"
                :class="item.visible ? 'bg-[#48BB78]' : 'bg-[#1E293B]'"
                @click="toggleVisibility(item)"
              >
                <span
                  class="inline-block h-4 w-4 rounded-full bg-white transition-transform"
                  :class="item.visible ? 'translate-x-6' : 'translate-x-1'"
                />
              </button>
            </td>
            <td class="px-4 py-3 text-sm text-[#94A3B8]">{{ formatDate(item.last_modified) }}</td>
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
                  @click="deleteItem(item)"
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

    <!-- Create/Edit Section Modal (D11) -->
    <Modal
      :open="showCreate"
      :title="editingId ? 'Edit Section' : 'Add Section'"
      max-width="640px"
      @close="showCreate = false"
    >
      <div class="space-y-4">
        <!-- Section Type -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-2">Section Type <span class="text-[#DC2626]">*</span></label>
          <div class="grid grid-cols-3 gap-2">
            <button
              v-for="(label, key) in sectionTypeLabels"
              :key="key"
              class="flex flex-col items-center gap-1.5 px-3 py-3 border rounded-lg transition-colors text-center"
              :class="form.type === key
                ? 'border-[#48BB78] bg-[#48BB78]/10 text-[#48BB78]'
                : 'border-[#1E293B] text-[#94A3B8] hover:border-[#94A3B8]'"
              @click="form.type = key as SectionType"
            >
              <span class="material-symbols-rounded text-xl">{{ sectionTypeIcons[key as SectionType] }}</span>
              <span class="text-xs">{{ label }}</span>
            </button>
          </div>
        </div>

        <!-- Title -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-1">
            Title <span class="text-[#DC2626]">*</span>
            <span class="text-[#64748B] text-xs ml-1">{{ form.name.length }}/60</span>
          </label>
          <input
            v-model="form.name"
            type="text"
            maxlength="60"
            placeholder="Section title"
            class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm"
          />
        </div>

        <!-- Banner fields -->
        <template v-if="form.type === 'banner'">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Banner Image URL <span class="text-[#DC2626]">*</span></label>
            <input
              v-model="form.banner_image"
              type="text"
              placeholder="https://..."
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm"
            />
            <p class="text-xs text-[#64748B] mt-1">Recommended: 1440x400px, JPG/PNG/WebP</p>
          </div>
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Link URL <span class="text-[#64748B] text-xs">(optional)</span></label>
            <input
              v-model="form.banner_link"
              type="text"
              placeholder="https://..."
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm"
            />
          </div>
        </template>

        <!-- Widget fields -->
        <template v-if="['quest_widget', 'leaderboard_widget', 'points_widget'].includes(form.type)">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Instance <span class="text-[#DC2626]">*</span></label>
            <select
              v-model="form.widget_instance_id"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
            >
              <option value="" disabled>Select instance</option>
              <option v-for="w in filteredWidgets" :key="w.id" :value="w.id">{{ w.name }}</option>
            </select>
            <p v-if="filteredWidgets.length === 0" class="text-xs text-[#D97706] mt-1">
              No instances available. Create one in the corresponding module first.
            </p>
          </div>
        </template>

        <!-- Rich Text fields -->
        <template v-if="form.type === 'rich_text'">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">Content <span class="text-[#DC2626]">*</span></label>
            <textarea
              v-model="form.rich_text_content"
              rows="8"
              placeholder="Write your content here. Supports basic formatting."
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm resize-none"
            />
          </div>
        </template>

        <!-- Custom HTML fields -->
        <template v-if="form.type === 'custom_html'">
          <div>
            <label class="block text-sm text-[#94A3B8] mb-1">HTML Code <span class="text-[#DC2626]">*</span></label>
            <textarea
              v-model="form.custom_html"
              rows="10"
              placeholder="<div>Your custom HTML...</div>"
              class="w-full px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] placeholder-[#64748B] focus:border-[#48BB78] focus:outline-none text-sm font-mono resize-none"
            />
          </div>
          <div class="flex items-start gap-2 p-3 bg-[#1F1A08]/50 border border-[#D97706]/20 rounded-lg">
            <span class="material-symbols-rounded text-lg text-[#D97706] mt-0.5">warning</span>
            <div>
              <p class="text-sm text-[#F1F5F9] font-medium">XSS Sanitization</p>
              <p class="text-xs text-[#94A3B8] mt-0.5">Custom HTML is sanitized on save. Script tags, event handlers, and iframes from untrusted sources will be removed.</p>
            </div>
          </div>
        </template>

        <!-- Visibility -->
        <div>
          <label class="block text-sm text-[#94A3B8] mb-2">Visibility</label>
          <div class="flex gap-3">
            <label
              v-for="(label, key) in ({ all: 'All Users', logged_in: 'Logged-in Only', level_gated: 'Level Gated' } as Record<VisibilityScope, string>)"
              :key="key"
              class="flex items-center gap-2 px-3 py-2.5 border rounded-lg cursor-pointer transition-colors text-sm"
              :class="form.visibility_scope === key
                ? 'border-[#48BB78] bg-[#48BB78]/10 text-[#48BB78]'
                : 'border-[#1E293B] text-[#94A3B8] hover:border-[#94A3B8]'"
            >
              <input
                v-model="form.visibility_scope"
                type="radio"
                :value="key"
                class="sr-only"
              />
              {{ label }}
            </label>
          </div>
          <div v-if="form.visibility_scope === 'level_gated'" class="mt-3">
            <label class="block text-sm text-[#94A3B8] mb-1">Minimum Level</label>
            <input
              v-model.number="form.min_level"
              type="number"
              min="1"
              class="w-32 px-4 py-2.5 bg-[#0A0F1A] border border-[#1E293B] rounded-lg text-[#F1F5F9] text-sm focus:border-[#48BB78] focus:outline-none"
            />
          </div>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-[#64748B] hover:text-[#F1F5F9] transition-colors" @click="showCreate = false">Cancel</button>
        <button
          class="px-4 py-2 bg-[#48BB78] text-black text-sm font-medium rounded-lg hover:bg-[#48BB78]/90 disabled:opacity-50 transition-colors"
          :disabled="!formValid || saving"
          @click="saveForm"
        >
          {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Add Section') }}
        </button>
      </template>
    </Modal>
  </div>
</template>
