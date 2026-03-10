<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../../api/client'
import Modal from '../../components/common/Modal.vue'

interface WidgetBlock {
  id: string
  type: string
  category: string
  label: string
  icon: string
  props: Record<string, unknown>
}

interface PageData {
  id?: string
  name: string
  slug: string
  status: 'published' | 'draft'
  theme: 'light' | 'dark'
  meta_title: string
  meta_description: string
  blocks: WidgetBlock[]
}

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const saving = ref(false)
const publishing = ref(false)
const showSettings = ref(false)

const isEdit = computed(() => route.params.id !== undefined)

const page = reactive<PageData>({
  name: 'Untitled Page',
  slug: '',
  status: 'draft',
  theme: 'dark',
  meta_title: '',
  meta_description: '',
  blocks: [],
})

const selectedBlockId = ref<string | null>(null)
const selectedBlock = computed(() => page.blocks.find(b => b.id === selectedBlockId.value) || null)
const dragType = ref<string | null>(null)

const editingName = ref(false)
const nameInput = ref<HTMLInputElement | null>(null)

// Widget palette definition
const palette = [
  {
    category: 'Layout',
    items: [
      { type: 'header', label: 'Header', icon: 'web_asset' },
      { type: 'footer', label: 'Footer', icon: 'call_to_action' },
      { type: 'divider', label: 'Divider', icon: 'horizontal_rule' },
      { type: 'spacer', label: 'Spacer', icon: 'expand' },
    ],
  },
  {
    category: 'Content',
    items: [
      { type: 'hero_banner', label: 'Hero Banner', icon: 'image' },
      { type: 'text_block', label: 'Text Block', icon: 'notes' },
      { type: 'image', label: 'Image', icon: 'photo' },
      { type: 'cta_button', label: 'CTA Button', icon: 'smart_button' },
    ],
  },
  {
    category: 'Community',
    items: [
      { type: 'leaderboard_widget', label: 'Leaderboard Widget', icon: 'leaderboard' },
      { type: 'task_list', label: 'Task List', icon: 'task_alt' },
      { type: 'points_display', label: 'Points Display', icon: 'stars' },
      { type: 'shop_grid', label: 'Shop Grid', icon: 'storefront' },
      { type: 'daychain_streak', label: 'DayChain Streak', icon: 'local_fire_department' },
    ],
  },
]

// Default properties per widget type
const defaultProps: Record<string, Record<string, unknown>> = {
  header: { title: 'Page Header', subtitle: '', showLogo: true },
  footer: { text: 'Powered by TaskOn', links: [] },
  divider: { thickness: 1, color: '#1E293B' },
  spacer: { height: 40 },
  hero_banner: { title: 'Welcome', subtitle: 'Your community awaits', bgColor: '#9B7EE0', textAlign: 'center' },
  text_block: { content: 'Add your content here...', fontSize: 16 },
  image: { src: '', alt: '', width: '100%' },
  cta_button: { label: 'Get Started', url: '#', variant: 'primary' },
  leaderboard_widget: { showTopN: 10, showTrend: true },
  task_list: { groupBySector: true, showPoints: true },
  points_display: { showLevel: true, showProgress: true },
  shop_grid: { columns: 3, showPrices: true },
  daychain_streak: { showCalendar: true, showRewards: true },
}

// Common property labels for display
const commonPropLabels: Record<string, string> = {
  title: 'Title',
  subtitle: 'Subtitle',
  content: 'Content',
  text: 'Text',
  label: 'Label',
  url: 'URL',
  src: 'Image Source',
  alt: 'Alt Text',
  bgColor: 'Background Color',
  color: 'Color',
  textAlign: 'Text Align',
  fontSize: 'Font Size',
  height: 'Height (px)',
  thickness: 'Thickness (px)',
  width: 'Width',
  columns: 'Columns',
  variant: 'Variant',
  showTopN: 'Show Top N',
  showTrend: 'Show Trend Arrows',
  showLogo: 'Show Logo',
  showLevel: 'Show Level',
  showProgress: 'Show Progress',
  showPoints: 'Show Points',
  showPrices: 'Show Prices',
  showCalendar: 'Show Calendar',
  showRewards: 'Show Rewards',
  groupBySector: 'Group by Sector',
}

let nextId = 1
function genId() {
  return `block_${Date.now()}_${nextId++}`
}

function autoSlug(name: string) {
  return name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')
}

onMounted(async () => {
  if (isEdit.value) {
    await fetchPage()
  } else {
    // Handle template from query
    const tpl = route.query.template as string
    if (tpl === 'community_hub') {
      page.name = 'Community Hub'
      page.slug = 'community-hub'
      page.blocks = [
        { id: genId(), type: 'header', category: 'Layout', label: 'Header', icon: 'web_asset', props: { ...defaultProps.header, title: 'Community Hub' } },
        { id: genId(), type: 'hero_banner', category: 'Content', label: 'Hero Banner', icon: 'image', props: { ...defaultProps.hero_banner, title: 'Welcome to Our Community' } },
        { id: genId(), type: 'leaderboard_widget', category: 'Community', label: 'Leaderboard Widget', icon: 'leaderboard', props: { ...defaultProps.leaderboard_widget } },
        { id: genId(), type: 'task_list', category: 'Community', label: 'Task List', icon: 'task_alt', props: { ...defaultProps.task_list } },
        { id: genId(), type: 'footer', category: 'Layout', label: 'Footer', icon: 'call_to_action', props: { ...defaultProps.footer } },
      ]
    } else if (tpl === 'rewards_portal') {
      page.name = 'Rewards Portal'
      page.slug = 'rewards-portal'
      page.blocks = [
        { id: genId(), type: 'header', category: 'Layout', label: 'Header', icon: 'web_asset', props: { ...defaultProps.header, title: 'Rewards Portal' } },
        { id: genId(), type: 'hero_banner', category: 'Content', label: 'Hero Banner', icon: 'image', props: { ...defaultProps.hero_banner, title: 'Earn & Redeem Rewards' } },
        { id: genId(), type: 'shop_grid', category: 'Community', label: 'Shop Grid', icon: 'storefront', props: { ...defaultProps.shop_grid } },
        { id: genId(), type: 'points_display', category: 'Community', label: 'Points Display', icon: 'stars', props: { ...defaultProps.points_display } },
        { id: genId(), type: 'footer', category: 'Layout', label: 'Footer', icon: 'call_to_action', props: { ...defaultProps.footer } },
      ]
    }
  }
  loading.value = false
})

async function fetchPage() {
  try {
    const res = await api.get(`/api/v1/whitelabel/pages/${route.params.id}`)
    const data = res.data.data
    if (data) {
      page.id = data.id
      page.name = data.name || 'Untitled Page'
      page.slug = data.slug || ''
      page.status = data.status || 'draft'
      page.theme = data.theme || 'dark'
      page.meta_title = data.meta_title || ''
      page.meta_description = data.meta_description || ''
      page.blocks = data.blocks || []
    }
  } catch { /* empty */ }
}

function onDragStart(item: { type: string; label: string; icon: string }, category: string) {
  dragType.value = item.type
}

function onDragEnd() {
  dragType.value = null
}

function onDrop(event: DragEvent) {
  event.preventDefault()
  if (!dragType.value) return

  const catItem = palette.flatMap(c => c.items.map(i => ({ ...i, category: c.category }))).find(i => i.type === dragType.value)
  if (!catItem) return

  const block: WidgetBlock = {
    id: genId(),
    type: catItem.type,
    category: catItem.category,
    label: catItem.label,
    icon: catItem.icon,
    props: { ...(defaultProps[catItem.type] || {}) },
  }
  page.blocks.push(block)
  selectedBlockId.value = block.id
  dragType.value = null
}

function onDragOver(event: DragEvent) {
  event.preventDefault()
}

function addBlock(item: { type: string; label: string; icon: string }, category: string) {
  const block: WidgetBlock = {
    id: genId(),
    type: item.type,
    category,
    label: item.label,
    icon: item.icon,
    props: { ...(defaultProps[item.type] || {}) },
  }
  page.blocks.push(block)
  selectedBlockId.value = block.id
}

function selectBlock(id: string) {
  selectedBlockId.value = selectedBlockId.value === id ? null : id
}

function removeBlock(id: string) {
  page.blocks = page.blocks.filter(b => b.id !== id)
  if (selectedBlockId.value === id) selectedBlockId.value = null
}

function moveBlock(id: string, dir: -1 | 1) {
  const idx = page.blocks.findIndex(b => b.id === id)
  if (idx < 0) return
  const newIdx = idx + dir
  if (newIdx < 0 || newIdx >= page.blocks.length) return
  const temp = page.blocks[idx]
  page.blocks[idx] = page.blocks[newIdx]
  page.blocks[newIdx] = temp
}

function startEditName() {
  editingName.value = true
  setTimeout(() => nameInput.value?.focus(), 50)
}

function finishEditName() {
  editingName.value = false
  if (!page.slug) page.slug = autoSlug(page.name)
}

async function saveDraft() {
  saving.value = true
  try {
    const payload = {
      name: page.name,
      slug: page.slug || autoSlug(page.name),
      status: 'draft' as const,
      theme: page.theme,
      meta_title: page.meta_title,
      meta_description: page.meta_description,
      blocks: page.blocks,
    }
    if (isEdit.value) {
      await api.put(`/api/v1/whitelabel/pages/${route.params.id}`, payload)
    } else {
      const res = await api.post('/api/v1/whitelabel/pages', payload)
      const id = res.data.data?.id
      if (id) router.replace(`/b/whitelabel/pages/${id}/edit`)
    }
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function publish() {
  publishing.value = true
  try {
    const payload = {
      name: page.name,
      slug: page.slug || autoSlug(page.name),
      status: 'published' as const,
      theme: page.theme,
      meta_title: page.meta_title,
      meta_description: page.meta_description,
      blocks: page.blocks,
    }
    if (isEdit.value) {
      await api.put(`/api/v1/whitelabel/pages/${route.params.id}`, payload)
    } else {
      await api.post('/api/v1/whitelabel/pages', payload)
    }
    page.status = 'published'
  } catch { /* TODO: toast */ }
  publishing.value = false
}

function previewPage() {
  const slug = page.slug || autoSlug(page.name)
  window.open(`/preview/pages/${slug}`, '_blank')
}
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-64px)] -m-6">
    <!-- Top Bar -->
    <div class="flex items-center justify-between px-4 py-3 bg-card-bg border-b border-border shrink-0">
      <div class="flex items-center gap-3">
        <button
          class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
          @click="router.push('/b/whitelabel/pages')"
        >
          <span class="material-symbols-rounded text-xl text-text-muted">arrow_back</span>
        </button>

        <!-- Editable page name -->
        <div class="flex items-center gap-2">
          <input
            v-if="editingName"
            ref="nameInput"
            v-model="page.name"
            type="text"
            class="px-2 py-1 bg-page-bg border border-wl rounded text-sm font-semibold text-text-primary focus:outline-none"
            @blur="finishEditName"
            @keyup.enter="finishEditName"
          />
          <button
            v-else
            class="text-sm font-semibold text-text-primary hover:text-wl transition-colors flex items-center gap-1"
            @click="startEditName"
          >
            {{ page.name }}
            <span class="material-symbols-rounded text-sm text-text-muted">edit</span>
          </button>
        </div>

        <!-- Status badge -->
        <span
          class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
          :class="page.status === 'published'
            ? 'bg-status-active-bg text-status-active'
            : 'bg-status-draft-bg text-status-draft'"
        >
          {{ page.status === 'published' ? 'Published' : 'Draft' }}
        </span>
      </div>

      <div class="flex items-center gap-2">
        <button
          class="p-2 rounded-lg hover:bg-white/5 transition-colors"
          title="Page Settings"
          @click="showSettings = true"
        >
          <span class="material-symbols-rounded text-lg text-text-muted">settings</span>
        </button>
        <button
          class="px-3 py-1.5 text-sm text-text-secondary border border-border rounded-lg hover:bg-white/5 transition-colors"
          @click="previewPage"
        >
          Preview
        </button>
        <button
          class="px-3 py-1.5 text-sm text-text-secondary border border-border rounded-lg hover:bg-white/5 transition-colors disabled:opacity-50 flex items-center gap-1.5"
          :disabled="saving"
          @click="saveDraft"
        >
          <span v-if="saving" class="material-symbols-rounded text-sm animate-spin">progress_activity</span>
          {{ saving ? 'Saving...' : 'Save Draft' }}
        </button>
        <button
          class="px-4 py-1.5 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors disabled:opacity-50 flex items-center gap-1.5"
          :disabled="publishing"
          @click="publish"
        >
          <span v-if="publishing" class="material-symbols-rounded text-sm animate-spin">progress_activity</span>
          {{ publishing ? 'Publishing...' : 'Publish' }}
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <span class="material-symbols-rounded text-3xl text-text-muted animate-spin">progress_activity</span>
    </div>

    <!-- Editor Layout -->
    <div v-else class="flex flex-1 overflow-hidden">
      <!-- Left Sidebar: Widget Palette -->
      <div class="w-50 bg-card-bg border-r border-border overflow-y-auto shrink-0">
        <div class="p-3">
          <p class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Widgets</p>
          <div v-for="group in palette" :key="group.category" class="mb-4">
            <p class="text-xs text-text-muted mb-2">{{ group.category }}</p>
            <div class="space-y-1">
              <div
                v-for="item in group.items"
                :key="item.type"
                draggable="true"
                class="flex items-center gap-2 px-2.5 py-2 rounded-lg cursor-grab hover:bg-white/5 transition-colors text-text-secondary hover:text-text-primary"
                @dragstart="onDragStart(item, group.category)"
                @dragend="onDragEnd"
                @click="addBlock(item, group.category)"
              >
                <span class="material-symbols-rounded text-base">{{ item.icon }}</span>
                <span class="text-xs font-medium">{{ item.label }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Center: Canvas -->
      <div
        class="flex-1 overflow-y-auto p-6"
        :class="page.theme === 'dark' ? 'bg-page-bg' : 'bg-gray-50'"
        @drop="onDrop"
        @dragover="onDragOver"
      >
        <div class="max-w-3xl mx-auto min-h-full">
          <!-- Empty canvas -->
          <div
            v-if="page.blocks.length === 0"
            class="border-2 border-dashed border-border rounded-xl py-20 text-center"
          >
            <span class="material-symbols-rounded text-4xl text-text-muted block mb-3">add_box</span>
            <p class="text-sm text-text-muted mb-1">Drag widgets here or click them in the palette</p>
            <p class="text-xs text-text-muted">Build your page by adding layout and content blocks</p>
          </div>

          <!-- Blocks -->
          <div v-else class="space-y-2">
            <div
              v-for="block in page.blocks"
              :key="block.id"
              class="relative group rounded-lg border-2 transition-all cursor-pointer"
              :class="selectedBlockId === block.id
                ? 'border-wl bg-wl/5'
                : 'border-transparent hover:border-border'"
              @click="selectBlock(block.id)"
            >
              <!-- Block toolbar (on hover / select) -->
              <div
                class="absolute -top-3 right-2 flex items-center gap-0.5 bg-card-bg border border-border rounded-lg px-1 py-0.5 z-10 opacity-0 group-hover:opacity-100 transition-opacity"
                :class="selectedBlockId === block.id ? 'opacity-100!' : ''"
              >
                <button
                  class="p-1 rounded hover:bg-white/10 transition-colors"
                  title="Move Up"
                  @click.stop="moveBlock(block.id, -1)"
                >
                  <span class="material-symbols-rounded text-sm text-text-muted">keyboard_arrow_up</span>
                </button>
                <button
                  class="p-1 rounded hover:bg-white/10 transition-colors"
                  title="Move Down"
                  @click.stop="moveBlock(block.id, 1)"
                >
                  <span class="material-symbols-rounded text-sm text-text-muted">keyboard_arrow_down</span>
                </button>
                <button
                  class="p-1 rounded hover:bg-white/10 transition-colors"
                  title="Remove"
                  @click.stop="removeBlock(block.id)"
                >
                  <span class="material-symbols-rounded text-sm text-status-paused">close</span>
                </button>
              </div>

              <!-- Block content preview -->
              <div class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <span class="material-symbols-rounded text-base text-text-muted cursor-grab">drag_indicator</span>
                  <span class="material-symbols-rounded text-base text-wl">{{ block.icon }}</span>
                  <span class="text-xs font-medium text-text-secondary">{{ block.label }}</span>
                  <span class="text-xs text-text-muted">{{ block.category }}</span>
                </div>

                <!-- Visual preview of block -->
                <div class="mt-2">
                  <!-- Header preview -->
                  <div v-if="block.type === 'header'" class="bg-page-bg rounded-lg p-4 border border-border">
                    <div class="text-base font-bold text-text-primary">{{ block.props.title || 'Header' }}</div>
                    <div v-if="block.props.subtitle" class="text-xs text-text-muted mt-1">{{ block.props.subtitle }}</div>
                  </div>
                  <!-- Hero banner preview -->
                  <div v-else-if="block.type === 'hero_banner'" class="rounded-lg p-6 text-center" :style="{ backgroundColor: (block.props.bgColor as string) || '#9B7EE0' }">
                    <div class="text-lg font-bold text-white">{{ block.props.title || 'Hero Banner' }}</div>
                    <div v-if="block.props.subtitle" class="text-sm text-white/70 mt-1">{{ block.props.subtitle }}</div>
                  </div>
                  <!-- Text block preview -->
                  <div v-else-if="block.type === 'text_block'" class="text-sm text-text-secondary">
                    {{ block.props.content || 'Text content...' }}
                  </div>
                  <!-- Divider preview -->
                  <div v-else-if="block.type === 'divider'">
                    <hr class="border-border" :style="{ borderWidth: (block.props.thickness as number || 1) + 'px' }" />
                  </div>
                  <!-- Spacer preview -->
                  <div v-else-if="block.type === 'spacer'" class="flex items-center justify-center" :style="{ height: (block.props.height as number || 40) + 'px' }">
                    <span class="text-xs text-text-muted">{{ block.props.height }}px spacer</span>
                  </div>
                  <!-- CTA button preview -->
                  <div v-else-if="block.type === 'cta_button'" class="flex justify-center">
                    <span class="px-6 py-2 bg-wl text-white text-sm font-medium rounded-lg">{{ block.props.label || 'Button' }}</span>
                  </div>
                  <!-- Community widget preview -->
                  <div v-else class="bg-page-bg rounded-lg p-4 border border-border">
                    <div class="flex items-center gap-2 mb-2">
                      <span class="material-symbols-rounded text-base text-wl">{{ block.icon }}</span>
                      <span class="text-xs font-medium text-text-primary">{{ block.label }}</span>
                    </div>
                    <div class="space-y-2">
                      <div v-for="i in 3" :key="i" class="h-4 bg-border rounded" :style="{ width: (100 - i * 15) + '%' }"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Drop hint when dragging -->
          <div
            v-if="dragType"
            class="mt-4 border-2 border-dashed border-wl/40 rounded-xl py-8 text-center"
          >
            <span class="material-symbols-rounded text-2xl text-wl/50 block mb-1">add_circle</span>
            <p class="text-xs text-wl/50">Drop here to add</p>
          </div>
        </div>
      </div>

      <!-- Right Sidebar: Properties -->
      <div
        class="w-70 bg-card-bg border-l border-border overflow-y-auto shrink-0 transition-all"
        :class="selectedBlock ? 'translate-x-0' : 'translate-x-full hidden'"
      >
        <div v-if="selectedBlock" class="p-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-sm font-semibold text-text-primary">Widget Properties</h3>
            <button
              class="p-1 rounded hover:bg-white/5 transition-colors"
              @click="selectedBlockId = null"
            >
              <span class="material-symbols-rounded text-base text-text-muted">close</span>
            </button>
          </div>

          <div class="space-y-4">
            <!-- Widget type info -->
            <div class="flex items-center gap-2 pb-3 border-b border-border">
              <span class="material-symbols-rounded text-base text-wl">{{ selectedBlock.icon }}</span>
              <span class="text-xs font-medium text-text-secondary">{{ selectedBlock.label }}</span>
            </div>

            <!-- Properties -->
            <div v-for="(value, key) in selectedBlock.props" :key="key" class="space-y-1">
              <label class="text-xs text-text-muted">{{ commonPropLabels[key as string] || key }}</label>

              <!-- Boolean toggle -->
              <div v-if="typeof value === 'boolean'" class="flex items-center justify-between">
                <span class="text-sm text-text-secondary">{{ commonPropLabels[key as string] || key }}</span>
                <button
                  class="relative w-9 h-5 rounded-full transition-colors"
                  :class="value ? 'bg-wl' : 'bg-border'"
                  @click="selectedBlock.props[key as string] = !value"
                >
                  <span
                    class="absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full transition-transform"
                    :class="value ? 'translate-x-4' : ''"
                  ></span>
                </button>
              </div>

              <!-- Number input -->
              <input
                v-else-if="typeof value === 'number'"
                :value="value"
                type="number"
                class="w-full px-3 py-1.5 bg-page-bg border border-border rounded-lg text-text-primary text-xs focus:border-wl focus:outline-none"
                @input="selectedBlock.props[key as string] = Number(($event.target as HTMLInputElement).value)"
              />

              <!-- Color input -->
              <div v-else-if="(key as string).toLowerCase().includes('color')" class="flex items-center gap-2">
                <input
                  :value="value"
                  type="color"
                  class="w-8 h-8 rounded border border-border cursor-pointer bg-transparent"
                  @input="selectedBlock.props[key as string] = ($event.target as HTMLInputElement).value"
                />
                <input
                  :value="value"
                  type="text"
                  class="flex-1 px-3 py-1.5 bg-page-bg border border-border rounded-lg text-text-primary text-xs font-mono focus:border-wl focus:outline-none"
                  @input="selectedBlock.props[key as string] = ($event.target as HTMLInputElement).value"
                />
              </div>

              <!-- Text input (default) -->
              <input
                v-else-if="typeof value === 'string'"
                :value="value"
                type="text"
                class="w-full px-3 py-1.5 bg-page-bg border border-border rounded-lg text-text-primary text-xs focus:border-wl focus:outline-none"
                @input="selectedBlock.props[key as string] = ($event.target as HTMLInputElement).value"
              />
            </div>

            <!-- Visibility toggle (common to all) -->
            <div class="pt-3 border-t border-border">
              <p class="text-xs text-text-muted mb-2">Visibility</p>
              <div class="flex gap-2">
                <button class="flex-1 px-3 py-1.5 text-xs font-medium bg-wl/15 text-wl rounded-lg">Visible</button>
                <button class="flex-1 px-3 py-1.5 text-xs font-medium text-text-muted border border-border rounded-lg hover:bg-white/5">Hidden</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Page Settings Slide-out -->
    <Teleport to="body">
      <div v-if="showSettings" class="fixed inset-0 z-50 flex justify-end">
        <div class="absolute inset-0 bg-black/60" @click="showSettings = false"></div>
        <div class="relative w-100 bg-card-bg border-l border-border h-full overflow-y-auto">
          <div class="flex items-center justify-between px-6 py-4 border-b border-border">
            <h2 class="text-lg font-semibold text-text-primary">Page Settings</h2>
            <button class="text-text-muted hover:text-text-primary transition-colors" @click="showSettings = false">
              <span class="material-symbols-rounded text-xl">close</span>
            </button>
          </div>
          <div class="px-6 py-5 space-y-5">
            <!-- Page Name -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Page Name</label>
              <input
                v-model="page.name"
                type="text"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-wl focus:outline-none"
              />
            </div>

            <!-- URL Slug -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">URL Slug</label>
              <div class="flex items-center">
                <span class="px-3 py-2.5 bg-page-bg border border-r-0 border-border rounded-l-lg text-xs text-text-muted">/pages/</span>
                <input
                  v-model="page.slug"
                  type="text"
                  :placeholder="autoSlug(page.name)"
                  class="flex-1 px-3 py-2.5 bg-page-bg border border-border rounded-r-lg text-text-primary text-sm font-mono focus:border-wl focus:outline-none"
                />
              </div>
            </div>

            <!-- Theme -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Theme</label>
              <div class="flex gap-3">
                <button
                  class="flex-1 px-4 py-2.5 rounded-lg text-sm font-medium border transition-colors"
                  :class="page.theme === 'light'
                    ? 'bg-white text-gray-900 border-wl'
                    : 'bg-page-bg text-text-secondary border-border hover:border-text-muted'"
                  @click="page.theme = 'light'"
                >
                  Light
                </button>
                <button
                  class="flex-1 px-4 py-2.5 rounded-lg text-sm font-medium border transition-colors"
                  :class="page.theme === 'dark'
                    ? 'bg-page-bg text-white border-wl'
                    : 'bg-page-bg text-text-secondary border-border hover:border-text-muted'"
                  @click="page.theme = 'dark'"
                >
                  Dark
                </button>
              </div>
            </div>

            <!-- Meta Title -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Meta Title</label>
              <input
                v-model="page.meta_title"
                type="text"
                placeholder="SEO page title"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-wl focus:outline-none"
              />
            </div>

            <!-- Meta Description -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Meta Description</label>
              <textarea
                v-model="page.meta_description"
                rows="3"
                placeholder="SEO page description"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-wl focus:outline-none resize-none"
              />
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
