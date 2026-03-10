<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'

interface ModuleStatus {
  key: string
  name: string
  icon: string
  configured: boolean
  community_route: string
}

interface WidgetItem {
  id: string
  module_key: string
  name: string
  status: 'active' | 'inactive'
  theme: 'light' | 'dark'
  impressions_24h: number
  avg_load_time_ms: number
  embed_code: string
  created_at: string
}

const router = useRouter()
const loading = ref(true)
const modules = ref<ModuleStatus[]>([])
const widgets = ref<WidgetItem[]>([])

const MODULE_DEFS: Record<string, { name: string; icon: string; route: string }> = {
  leaderboard: { name: 'Leaderboard', icon: 'leaderboard', route: '/b/community/modules/leaderboard' },
  sectors_tasks: { name: 'Tasks / Quests', icon: 'task_alt', route: '/b/community/modules/sectors' },
  points_level: { name: 'User Center', icon: 'person', route: '/b/community/modules/points' },
  benefits_shop: { name: 'Benefits Shop', icon: 'storefront', route: '/b/community/modules/shop' },
  day_chain: { name: 'DayChain', icon: 'local_fire_department', route: '/b/community/modules/daychain' },
}

const configuredModules = computed(() => modules.value.filter(m => m.configured))
const unconfiguredModules = computed(() => modules.value.filter(m => !m.configured))
const hasWidgets = computed(() => widgets.value.length > 0)
const hasCommunity = computed(() => modules.value.length > 0)

const totalWidgets = computed(() => widgets.value.length)
const activeWidgets = computed(() => widgets.value.filter(w => w.status === 'active').length)
const totalImpressions = computed(() => widgets.value.reduce((s, w) => s + w.impressions_24h, 0))
const avgLoadTime = computed(() => {
  if (!widgets.value.length) return '0ms'
  const avg = widgets.value.reduce((s, w) => s + w.avg_load_time_ms, 0) / widgets.value.length
  return `${Math.round(avg)}ms`
})

const copiedId = ref<string | null>(null)

onMounted(async () => {
  await Promise.all([fetchModules(), fetchWidgets()])
  loading.value = false
})

async function fetchModules() {
  try {
    const res = await api.get('/api/v1/community/modules/status')
    const data = res.data.data || []
    modules.value = data.map((m: { key: string; configured: boolean }) => ({
      key: m.key,
      name: MODULE_DEFS[m.key]?.name || m.key,
      icon: MODULE_DEFS[m.key]?.icon || 'widgets',
      configured: m.configured,
      community_route: MODULE_DEFS[m.key]?.route || '/b/community',
    }))
  } catch {
    // If community not set up, modules will be empty
  }
}

async function fetchWidgets() {
  try {
    const res = await api.get('/api/v1/whitelabel/widgets')
    widgets.value = res.data.data?.items || res.data.data || []
  } catch { /* empty */ }
}

async function toggleWidget(widget: WidgetItem) {
  const newStatus = widget.status === 'active' ? 'inactive' : 'active'
  const old = widget.status
  widget.status = newStatus
  try {
    await api.put(`/api/v1/whitelabel/widgets/${widget.id}`, { status: newStatus })
  } catch {
    widget.status = old
  }
}

function copyEmbed(widget: WidgetItem) {
  navigator.clipboard.writeText(widget.embed_code).catch(() => {})
  copiedId.value = widget.id
  setTimeout(() => { copiedId.value = null }, 2000)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Widget Library</h1>
        <p class="text-sm text-text-secondary">Embed community modules as widgets in your White Label site</p>
      </div>
      <button
        v-if="hasWidgets"
        class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors"
        @click="router.push('/b/whitelabel/widgets/new')"
      >
        + New Widget
      </button>
    </div>

    <!-- Loading -->
    <template v-if="loading">
      <div class="grid grid-cols-4 gap-4">
        <div v-for="i in 4" :key="i" class="bg-card-bg border border-border rounded-xl p-4">
          <div class="h-4 w-24 bg-border rounded animate-pulse mb-3"></div>
          <div class="h-8 w-16 bg-border rounded animate-pulse"></div>
        </div>
      </div>
      <div class="grid grid-cols-3 gap-4">
        <div v-for="i in 6" :key="i" class="bg-card-bg border border-border rounded-xl p-6 h-36 animate-pulse"></div>
      </div>
    </template>

    <!-- Empty: No Community at all -->
    <template v-else-if="!hasCommunity">
      <div class="flex flex-col items-center justify-center py-20 text-center">
        <span class="material-symbols-rounded text-5xl text-text-muted mb-4">widgets</span>
        <h3 class="text-lg font-semibold text-text-primary mb-2">Set Up Your Community First</h3>
        <p class="text-sm text-text-secondary max-w-md mb-6">
          Widgets are powered by your Community modules. Create a Community to start building widgets.
        </p>
        <router-link
          to="/b/community"
          class="px-6 py-2.5 bg-community text-white font-medium rounded-lg hover:bg-community/90 transition-colors"
        >
          Go to Community
        </router-link>
      </div>
    </template>

    <template v-else>
      <!-- Stats Row (if has widgets) -->
      <div v-if="hasWidgets" class="grid grid-cols-4 gap-4">
        <StatsCard label="Total Widgets" :value="totalWidgets" icon="widgets" icon-color="#9B7EE0" />
        <StatsCard label="Active" :value="activeWidgets" icon="toggle_on" icon-color="#16A34A" />
        <StatsCard label="Impressions (24h)" :value="totalImpressions" icon="visibility" icon-color="#3B82F6" />
        <StatsCard label="Avg Load Time" :value="avgLoadTime" icon="speed" icon-color="#F59E0B" />
      </div>

      <!-- CONFIGURED MODULES -->
      <div>
        <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3 flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-status-active"></span>
          Configured Modules
        </h2>
        <div v-if="configuredModules.length === 0" class="bg-card-bg border border-border rounded-xl p-6 text-center">
          <p class="text-sm text-text-muted">No modules configured yet. Set up modules in Community first.</p>
        </div>
        <div v-else class="grid grid-cols-3 gap-4">
          <div
            v-for="mod in configuredModules"
            :key="mod.key"
            class="bg-card-bg border border-border rounded-xl p-5 hover:border-wl/30 transition-colors"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="w-10 h-10 rounded-lg bg-status-active-bg flex items-center justify-center">
                <span class="material-symbols-rounded text-xl text-status-active">{{ mod.icon }}</span>
              </div>
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-status-active-bg text-status-active">
                Active
              </span>
            </div>
            <h3 class="text-sm font-semibold text-text-primary mb-3">{{ mod.name }}</h3>
            <button
              class="w-full px-3 py-2 text-sm font-medium text-wl border border-wl rounded-lg hover:bg-wl/10 transition-colors"
              @click="router.push(`/b/whitelabel/widgets/new?module=${mod.key}`)"
            >
              Configure Widget &rarr;
            </button>
          </div>
        </div>
      </div>

      <!-- NOT YET CONFIGURED -->
      <div v-if="unconfiguredModules.length > 0">
        <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3 flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-status-draft"></span>
          Not Yet Configured
        </h2>
        <div class="grid grid-cols-3 gap-4">
          <div
            v-for="mod in unconfiguredModules"
            :key="mod.key"
            class="bg-card-bg border border-border rounded-xl p-5 hover:border-amber-700/30 transition-colors"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="w-10 h-10 rounded-lg bg-status-draft-bg flex items-center justify-center">
                <span class="material-symbols-rounded text-xl text-status-draft">{{ mod.icon }}</span>
              </div>
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-status-draft-bg text-status-draft">
                Not Configured
              </span>
            </div>
            <h3 class="text-sm font-semibold text-text-primary mb-3">{{ mod.name }}</h3>
            <button
              class="w-full px-3 py-2 text-sm font-medium text-status-draft border border-status-draft rounded-lg hover:bg-amber-900/20 transition-colors"
              @click="router.push(mod.community_route)"
            >
              Set Up in Community &rarr;
            </button>
          </div>
        </div>
      </div>

      <!-- Existing Widgets (if any) -->
      <div v-if="hasWidgets">
        <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3 flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-wl"></span>
          Your Widgets
        </h2>
        <div class="grid grid-cols-2 gap-4">
          <div
            v-for="widget in widgets"
            :key="widget.id"
            class="bg-card-bg border border-border rounded-xl p-5 hover:border-wl/30 transition-colors"
          >
            <div class="flex items-start justify-between mb-3">
              <div>
                <h3 class="text-sm font-semibold text-text-primary">{{ widget.name }}</h3>
                <p class="text-xs text-text-muted mt-0.5 capitalize">{{ widget.module_key.replace(/_/g, ' ') }}</p>
              </div>
              <div class="flex items-center gap-2">
                <!-- Status toggle -->
                <button
                  class="relative w-9 h-5 rounded-full transition-colors"
                  :class="widget.status === 'active' ? 'bg-status-active' : 'bg-border'"
                  :title="widget.status === 'active' ? 'Deactivate' : 'Activate'"
                  @click="toggleWidget(widget)"
                >
                  <span
                    class="absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full transition-transform"
                    :class="widget.status === 'active' ? 'translate-x-4' : ''"
                  ></span>
                </button>
                <!-- Edit -->
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  title="Edit"
                  @click="router.push(`/b/whitelabel/widgets/${widget.id}/config`)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
              </div>
            </div>

            <!-- Embed code preview -->
            <div class="bg-page-bg border border-border rounded-lg p-3 mb-3">
              <code class="text-xs text-text-muted font-mono block truncate">{{ widget.embed_code }}</code>
            </div>

            <div class="flex items-center justify-between">
              <div class="flex items-center gap-4 text-xs text-text-muted">
                <span class="flex items-center gap-1">
                  <span class="material-symbols-rounded text-sm">visibility</span>
                  {{ widget.impressions_24h.toLocaleString() }}
                </span>
                <span class="flex items-center gap-1">
                  <span class="material-symbols-rounded text-sm">speed</span>
                  {{ widget.avg_load_time_ms }}ms
                </span>
              </div>
              <button
                class="text-xs font-medium transition-colors flex items-center gap-1"
                :class="copiedId === widget.id ? 'text-status-active' : 'text-wl hover:text-wl/80'"
                @click="copyEmbed(widget)"
              >
                <span class="material-symbols-rounded text-sm">{{ copiedId === widget.id ? 'check' : 'content_copy' }}</span>
                {{ copiedId === widget.id ? 'Copied!' : 'Copy Code' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
