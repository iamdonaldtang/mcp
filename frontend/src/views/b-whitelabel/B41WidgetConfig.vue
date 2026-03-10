<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../../api/client'
import Modal from '../../components/common/Modal.vue'

type ModuleType = 'leaderboard' | 'sectors_tasks' | 'points_level' | 'benefits_shop' | 'day_chain'

interface WidgetConfig {
  id?: string
  name: string
  module_type: ModuleType | ''
  theme: 'light' | 'dark'
  primary_color: string
  border_radius: number
  padding: number
  refresh_interval: number
  display_options: Record<string, boolean>
}

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const saving = ref(false)
const showDeleteConfirm = ref(false)
const deleting = ref(false)

const isEdit = computed(() => route.params.id && route.params.id !== 'new')

const config = reactive<WidgetConfig>({
  name: '',
  module_type: '',
  theme: 'dark',
  primary_color: '#9B7EE0',
  border_radius: 12,
  padding: 16,
  refresh_interval: 30,
  display_options: {},
})

const MODULE_OPTIONS: Record<ModuleType, { label: string; icon: string; options: { key: string; label: string; default: boolean }[] }> = {
  leaderboard: {
    label: 'Leaderboard',
    icon: 'leaderboard',
    options: [
      { key: 'show_top_n', label: 'Show Top N (10-100)', default: true },
      { key: 'show_point_type', label: 'Show Point Type', default: true },
      { key: 'show_trend_arrows', label: 'Show Trend Arrows', default: true },
    ],
  },
  sectors_tasks: {
    label: 'Tasks / Quests',
    icon: 'task_alt',
    options: [
      { key: 'show_completed', label: 'Show Completed', default: false },
      { key: 'group_by_sector', label: 'Group by Sector', default: true },
      { key: 'show_points', label: 'Show Points', default: true },
    ],
  },
  points_level: {
    label: 'User Center',
    icon: 'person',
    options: [
      { key: 'show_level', label: 'Show Level', default: true },
      { key: 'show_achievements', label: 'Show Achievements', default: true },
      { key: 'show_activity', label: 'Show Activity', default: false },
    ],
  },
  benefits_shop: {
    label: 'Benefits Shop',
    icon: 'storefront',
    options: [
      { key: 'show_categories', label: 'Show Categories', default: true },
      { key: 'show_stock', label: 'Show Stock', default: true },
      { key: 'show_prices', label: 'Show Prices', default: true },
    ],
  },
  day_chain: {
    label: 'DayChain',
    icon: 'local_fire_department',
    options: [
      { key: 'show_streak', label: 'Show Streak Counter', default: true },
      { key: 'show_calendar', label: 'Show Calendar View', default: true },
      { key: 'show_rewards', label: 'Show Streak Rewards', default: true },
    ],
  },
}

const moduleTypes = Object.keys(MODULE_OPTIONS) as ModuleType[]
const refreshOptions = [
  { value: 15, label: '15 seconds' },
  { value: 30, label: '30 seconds' },
  { value: 60, label: '1 minute' },
  { value: 300, label: '5 minutes' },
]

const currentModuleConfig = computed(() => {
  if (!config.module_type) return null
  return MODULE_OPTIONS[config.module_type as ModuleType] || null
})

const embedCode = computed(() => {
  const id = config.id || 'YOUR_WIDGET_ID'
  return `<div id="taskon-widget" data-widget-id="${id}" data-theme="${config.theme}"></div>\n<script src="https://cdn.taskon.xyz/widgets/v1/loader.js"><\/script>`
})

const copiedEmbed = ref(false)

// Initialize display options when module type changes
watch(() => config.module_type, (newType) => {
  if (newType && MODULE_OPTIONS[newType as ModuleType]) {
    const opts: Record<string, boolean> = {}
    MODULE_OPTIONS[newType as ModuleType].options.forEach(o => {
      opts[o.key] = config.display_options[o.key] ?? o.default
    })
    config.display_options = opts
  }
})

onMounted(async () => {
  // Pre-fill module from query param
  const moduleParam = route.query.module as string
  if (moduleParam && MODULE_OPTIONS[moduleParam as ModuleType]) {
    config.module_type = moduleParam as ModuleType
  }

  if (isEdit.value) {
    await fetchWidget()
  }
  loading.value = false
})

async function fetchWidget() {
  try {
    const res = await api.get(`/api/v1/whitelabel/widgets/${route.params.id}`)
    const data = res.data.data
    if (data) {
      config.id = data.id
      config.name = data.name || ''
      config.module_type = data.module_type || ''
      config.theme = data.theme || 'dark'
      config.primary_color = data.primary_color || '#9B7EE0'
      config.border_radius = data.border_radius ?? 12
      config.padding = data.padding ?? 16
      config.refresh_interval = data.refresh_interval ?? 30
      config.display_options = data.display_options || {}
    }
  } catch { /* empty */ }
}

async function saveWidget() {
  if (!config.name.trim() || !config.module_type) return
  saving.value = true
  try {
    const payload = {
      name: config.name,
      module_type: config.module_type,
      theme: config.theme,
      primary_color: config.primary_color,
      border_radius: config.border_radius,
      padding: config.padding,
      refresh_interval: config.refresh_interval,
      display_options: config.display_options,
    }
    if (isEdit.value) {
      await api.put(`/api/v1/whitelabel/widgets/${route.params.id}`, payload)
    } else {
      await api.post('/api/v1/whitelabel/widgets', payload)
    }
    router.push('/b/whitelabel/widgets')
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function deleteWidget() {
  deleting.value = true
  try {
    await api.delete(`/api/v1/whitelabel/widgets/${route.params.id}`)
    router.push('/b/whitelabel/widgets')
  } catch { /* TODO: toast */ }
  deleting.value = false
}

function copyEmbedCode() {
  navigator.clipboard.writeText(embedCode.value).catch(() => {})
  copiedEmbed.value = true
  setTimeout(() => { copiedEmbed.value = false }, 2000)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
          @click="router.push('/b/whitelabel/widgets')"
        >
          <span class="material-symbols-rounded text-xl text-text-muted">arrow_back</span>
        </button>
        <div>
          <h1 class="text-2xl font-bold text-text-primary">Configure Widget</h1>
          <p class="text-sm text-text-secondary">
            {{ isEdit ? 'Edit widget settings' : 'Create a new embeddable widget' }}
          </p>
        </div>
        <span
          v-if="config.module_type && currentModuleConfig"
          class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-wl/15 text-wl"
        >
          <span class="material-symbols-rounded text-sm">{{ currentModuleConfig.icon }}</span>
          {{ currentModuleConfig.label }}
        </span>
      </div>
    </div>

    <!-- Loading -->
    <template v-if="loading">
      <div class="grid grid-cols-2 gap-6">
        <div class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
          <div class="h-5 w-32 bg-border rounded animate-pulse"></div>
          <div class="h-10 bg-border rounded animate-pulse"></div>
          <div class="h-10 bg-border rounded animate-pulse"></div>
          <div class="h-32 bg-border rounded animate-pulse"></div>
        </div>
        <div class="bg-card-bg border border-border rounded-xl p-6 h-80 animate-pulse"></div>
      </div>
    </template>

    <template v-else>
      <div class="grid grid-cols-2 gap-6">
        <!-- Left: Form -->
        <div class="space-y-6">
          <!-- Basic Settings -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-5">
            <h2 class="text-sm font-semibold text-text-muted uppercase tracking-wider">Basic Settings</h2>

            <!-- Widget Name -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Widget Name</label>
              <input
                v-model="config.name"
                type="text"
                maxlength="50"
                placeholder="My Leaderboard Widget"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted text-sm focus:border-wl focus:outline-none"
              />
              <p class="text-xs text-text-muted mt-1">{{ config.name.length }}/50 characters</p>
            </div>

            <!-- Module Type -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Module Type</label>
              <select
                v-model="config.module_type"
                :disabled="isEdit"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-wl focus:outline-none disabled:opacity-50"
              >
                <option value="" disabled>Select a module</option>
                <option v-for="mt in moduleTypes" :key="mt" :value="mt">
                  {{ MODULE_OPTIONS[mt].label }}
                </option>
              </select>
            </div>

            <!-- Theme -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Theme</label>
              <div class="flex gap-3">
                <button
                  class="flex-1 px-4 py-2.5 rounded-lg text-sm font-medium border transition-colors"
                  :class="config.theme === 'light'
                    ? 'bg-white text-gray-900 border-wl'
                    : 'bg-page-bg text-text-secondary border-border hover:border-text-muted'"
                  @click="config.theme = 'light'"
                >
                  <span class="material-symbols-rounded text-base align-middle mr-1">light_mode</span>
                  Light
                </button>
                <button
                  class="flex-1 px-4 py-2.5 rounded-lg text-sm font-medium border transition-colors"
                  :class="config.theme === 'dark'
                    ? 'bg-[#0A0F1A] text-white border-wl'
                    : 'bg-page-bg text-text-secondary border-border hover:border-text-muted'"
                  @click="config.theme = 'dark'"
                >
                  <span class="material-symbols-rounded text-base align-middle mr-1">dark_mode</span>
                  Dark
                </button>
              </div>
            </div>
          </div>

          <!-- Display Options -->
          <div v-if="currentModuleConfig" class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
            <h2 class="text-sm font-semibold text-text-muted uppercase tracking-wider">Display Options</h2>
            <div
              v-for="opt in currentModuleConfig.options"
              :key="opt.key"
              class="flex items-center justify-between"
            >
              <span class="text-sm text-text-primary">{{ opt.label }}</span>
              <button
                class="relative w-11 h-6 rounded-full transition-colors"
                :class="config.display_options[opt.key] ? 'bg-wl' : 'bg-border'"
                @click="config.display_options[opt.key] = !config.display_options[opt.key]"
              >
                <span
                  class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform"
                  :class="config.display_options[opt.key] ? 'translate-x-5' : ''"
                ></span>
              </button>
            </div>
          </div>

          <!-- Style Customization -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-5">
            <h2 class="text-sm font-semibold text-text-muted uppercase tracking-wider">Style Customization</h2>

            <!-- Primary Color -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Primary Color</label>
              <div class="flex items-center gap-3">
                <input
                  v-model="config.primary_color"
                  type="color"
                  class="w-10 h-10 rounded-lg border border-border cursor-pointer bg-transparent"
                />
                <input
                  v-model="config.primary_color"
                  type="text"
                  class="flex-1 px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm font-mono focus:border-wl focus:outline-none"
                />
              </div>
            </div>

            <!-- Border Radius -->
            <div>
              <div class="flex items-center justify-between mb-1.5">
                <label class="text-sm text-text-secondary">Border Radius</label>
                <span class="text-xs text-text-muted font-mono">{{ config.border_radius }}px</span>
              </div>
              <input
                v-model.number="config.border_radius"
                type="range"
                min="0"
                max="24"
                class="w-full accent-wl"
              />
            </div>

            <!-- Padding -->
            <div>
              <div class="flex items-center justify-between mb-1.5">
                <label class="text-sm text-text-secondary">Padding</label>
                <span class="text-xs text-text-muted font-mono">{{ config.padding }}px</span>
              </div>
              <input
                v-model.number="config.padding"
                type="range"
                min="8"
                max="32"
                class="w-full accent-wl"
              />
            </div>

            <!-- Refresh Interval -->
            <div>
              <label class="block text-sm text-text-secondary mb-1.5">Refresh Interval</label>
              <select
                v-model.number="config.refresh_interval"
                class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-wl focus:outline-none"
              >
                <option v-for="opt in refreshOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
            </div>
          </div>

          <!-- Embed Code -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-3">
            <div class="flex items-center justify-between">
              <h2 class="text-sm font-semibold text-text-muted uppercase tracking-wider">Embed Code</h2>
              <button
                class="text-xs font-medium flex items-center gap-1 transition-colors"
                :class="copiedEmbed ? 'text-status-active' : 'text-wl hover:text-wl/80'"
                @click="copyEmbedCode"
              >
                <span class="material-symbols-rounded text-sm">{{ copiedEmbed ? 'check' : 'content_copy' }}</span>
                {{ copiedEmbed ? 'Copied!' : 'Copy' }}
              </button>
            </div>
            <pre class="bg-page-bg border border-border rounded-lg p-4 text-xs text-text-secondary font-mono overflow-x-auto whitespace-pre-wrap">{{ embedCode }}</pre>
          </div>
        </div>

        <!-- Right: Live Preview -->
        <div class="space-y-6">
          <div class="bg-card-bg border border-border rounded-xl p-6 sticky top-6">
            <h2 class="text-sm font-semibold text-text-muted uppercase tracking-wider mb-4">Live Preview</h2>
            <div
              class="border transition-all overflow-hidden"
              :class="config.theme === 'dark' ? 'bg-[#0A0F1A] border-[#1E293B]' : 'bg-white border-gray-200'"
              :style="{
                borderRadius: config.border_radius + 'px',
                padding: config.padding + 'px',
              }"
            >
              <!-- Preview content based on module type -->
              <div v-if="!config.module_type" class="py-12 text-center">
                <span class="material-symbols-rounded text-3xl text-text-muted block mb-2">widgets</span>
                <p class="text-sm text-text-muted">Select a module type to preview</p>
              </div>

              <!-- Leaderboard preview -->
              <template v-else-if="config.module_type === 'leaderboard'">
                <div class="flex items-center justify-between mb-4">
                  <h3 class="text-sm font-semibold" :class="config.theme === 'dark' ? 'text-white' : 'text-gray-900'">Leaderboard</h3>
                  <span class="text-xs" :class="config.theme === 'dark' ? 'text-gray-400' : 'text-gray-500'">This Week</span>
                </div>
                <div v-for="i in 5" :key="i" class="flex items-center gap-3 py-2" :class="i < 5 ? 'border-b' : ''" :style="{ borderColor: config.theme === 'dark' ? '#1E293B' : '#E5E7EB' }">
                  <span class="w-6 text-center text-xs font-bold" :style="{ color: i <= 3 ? config.primary_color : (config.theme === 'dark' ? '#94A3B8' : '#6B7280') }">#{{ i }}</span>
                  <div class="w-7 h-7 rounded-full" :style="{ backgroundColor: config.primary_color + '30' }"></div>
                  <div class="flex-1">
                    <div class="h-3 rounded w-20" :class="config.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"></div>
                  </div>
                  <span class="text-xs font-medium" :style="{ color: config.primary_color }">{{ (600 - i * 80).toLocaleString() }} pts</span>
                </div>
              </template>

              <!-- Tasks preview -->
              <template v-else-if="config.module_type === 'sectors_tasks'">
                <h3 class="text-sm font-semibold mb-4" :class="config.theme === 'dark' ? 'text-white' : 'text-gray-900'">Tasks</h3>
                <div v-for="i in 3" :key="i" class="flex items-center gap-3 py-2.5" :class="i < 3 ? 'border-b' : ''" :style="{ borderColor: config.theme === 'dark' ? '#1E293B' : '#E5E7EB' }">
                  <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{ backgroundColor: config.primary_color + '20' }">
                    <span class="material-symbols-rounded text-base" :style="{ color: config.primary_color }">task_alt</span>
                  </div>
                  <div class="flex-1">
                    <div class="h-3 rounded w-28 mb-1.5" :class="config.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"></div>
                    <div class="h-2 rounded w-16" :class="config.theme === 'dark' ? 'bg-gray-800' : 'bg-gray-100'"></div>
                  </div>
                  <span class="text-xs font-medium px-2 py-0.5 rounded" :style="{ backgroundColor: config.primary_color + '20', color: config.primary_color }">+{{ i * 10 }} pts</span>
                </div>
              </template>

              <!-- Generic preview for other types -->
              <template v-else>
                <div class="flex items-center gap-2 mb-4">
                  <span class="material-symbols-rounded text-lg" :style="{ color: config.primary_color }">{{ currentModuleConfig?.icon }}</span>
                  <h3 class="text-sm font-semibold" :class="config.theme === 'dark' ? 'text-white' : 'text-gray-900'">{{ currentModuleConfig?.label }}</h3>
                </div>
                <div class="space-y-3">
                  <div v-for="i in 4" :key="i" class="h-8 rounded-lg" :class="config.theme === 'dark' ? 'bg-gray-800' : 'bg-gray-100'"></div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- Bottom Actions -->
      <div class="flex items-center justify-between pt-2">
        <button
          v-if="isEdit"
          class="px-4 py-2 text-sm font-medium text-status-paused hover:bg-status-paused/10 rounded-lg transition-colors"
          @click="showDeleteConfirm = true"
        >
          Delete Widget
        </button>
        <div v-else></div>
        <div class="flex gap-3">
          <button
            class="px-4 py-2 text-sm text-text-muted hover:text-text-primary transition-colors"
            @click="router.push('/b/whitelabel/widgets')"
          >
            Cancel
          </button>
          <button
            class="px-6 py-2.5 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors disabled:opacity-50 flex items-center gap-2"
            :disabled="saving || !config.name.trim() || !config.module_type"
            @click="saveWidget"
          >
            <span v-if="saving" class="material-symbols-rounded text-base animate-spin">progress_activity</span>
            {{ saving ? 'Saving...' : 'Save Widget' }}
          </button>
        </div>
      </div>
    </template>

    <!-- Delete Confirm Modal -->
    <Modal :open="showDeleteConfirm" title="Delete Widget" @close="showDeleteConfirm = false">
      <p class="text-sm text-text-secondary">
        Are you sure you want to delete this widget? Any pages or sites using this embed code will stop working. This action cannot be undone.
      </p>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showDeleteConfirm = false">Cancel</button>
        <button
          class="px-4 py-2 bg-status-paused text-white text-sm font-medium rounded-lg hover:bg-status-paused/90 disabled:opacity-50"
          :disabled="deleting"
          @click="deleteWidget"
        >
          {{ deleting ? 'Deleting...' : 'Delete Widget' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
