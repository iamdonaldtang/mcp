<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'
import StatsCard from '../../components/common/StatsCard.vue'
import Modal from '../../components/common/Modal.vue'

interface PageItem {
  id: string
  name: string
  slug: string
  status: 'published' | 'draft' | 'unpublished'
  page_views: number
  last_modified: string
  template?: string
}

const router = useRouter()
const loading = ref(true)
const pages = ref<PageItem[]>([])
const filterTab = ref<string>('all')
const showDeleteConfirm = ref(false)
const deleteTarget = ref<PageItem | null>(null)
const deleting = ref(false)

const filterTabs = ['all', 'published', 'draft', 'unpublished']

const filteredPages = computed(() => {
  if (filterTab.value === 'all') return pages.value
  return pages.value.filter(p => p.status === filterTab.value)
})

const hasPages = computed(() => pages.value.length > 0)
const totalPages = computed(() => pages.value.length)
const publishedCount = computed(() => pages.value.filter(p => p.status === 'published').length)
const draftCount = computed(() => pages.value.filter(p => p.status === 'draft').length)
const totalViews = computed(() => pages.value.reduce((s, p) => s + p.page_views, 0))

const templates = [
  { key: 'blank', name: 'Blank Page', icon: 'article', description: 'Start from scratch with an empty canvas' },
  { key: 'community_hub', name: 'Community Hub', icon: 'groups', description: 'Pre-built layout with leaderboard, tasks, and rewards' },
  { key: 'rewards_portal', name: 'Rewards Portal', icon: 'redeem', description: 'Showcase your benefits shop and milestones' },
]

function statusBadge(status: string) {
  switch (status) {
    case 'published':
      return { label: 'Published', bg: 'bg-status-active-bg', text: 'text-status-active' }
    case 'draft':
      return { label: 'Draft', bg: 'bg-status-draft-bg', text: 'text-status-draft' }
    default:
      return { label: 'Unpublished', bg: 'bg-[#1E293B]', text: 'text-text-muted' }
  }
}

onMounted(async () => {
  await fetchPages()
  loading.value = false
})

async function fetchPages() {
  try {
    const res = await api.get('/api/v1/whitelabel/pages')
    pages.value = res.data.data?.items || res.data.data || []
  } catch { /* empty */ }
}

function confirmDelete(page: PageItem) {
  deleteTarget.value = page
  showDeleteConfirm.value = true
}

async function deletePage() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.delete(`/api/v1/whitelabel/pages/${deleteTarget.value.id}`)
    showDeleteConfirm.value = false
    deleteTarget.value = null
    await fetchPages()
  } catch { /* TODO: toast */ }
  deleting.value = false
}

function createFromTemplate(templateKey: string) {
  router.push(`/b/whitelabel/pages/new?template=${templateKey}`)
}

function formatDate(dateStr: string) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Page Builder</h1>
        <p class="text-sm text-text-secondary">Create and manage custom pages for your White Label site</p>
      </div>
      <button
        v-if="hasPages"
        class="px-4 py-2 bg-wl text-white text-sm font-medium rounded-lg hover:bg-wl/90 transition-colors"
        @click="router.push('/b/whitelabel/pages/new')"
      >
        + Create New Page
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
        <div v-for="i in 6" :key="i" class="bg-card-bg border border-border rounded-xl h-56 animate-pulse"></div>
      </div>
    </template>

    <!-- Empty State -->
    <template v-else-if="!hasPages">
      <div class="flex flex-col items-center justify-center py-16 text-center">
        <span class="material-symbols-rounded text-5xl text-text-muted mb-4">web</span>
        <h3 class="text-lg font-semibold text-text-primary mb-2">Build Custom Pages</h3>
        <p class="text-sm text-text-secondary max-w-md mb-8">
          Create pages with our visual editor. Combine community widgets, content blocks, and custom branding.
        </p>
      </div>

      <!-- Template Cards -->
      <div>
        <h2 class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Start with a Template</h2>
        <div class="grid grid-cols-3 gap-4">
          <button
            v-for="tpl in templates"
            :key="tpl.key"
            class="bg-card-bg border border-border rounded-xl p-6 text-left hover:border-wl/30 transition-colors group"
            @click="createFromTemplate(tpl.key)"
          >
            <div class="w-12 h-12 rounded-xl bg-wl/10 flex items-center justify-center mb-4 group-hover:bg-wl/20 transition-colors">
              <span class="material-symbols-rounded text-2xl text-wl">{{ tpl.icon }}</span>
            </div>
            <h3 class="text-sm font-semibold text-text-primary mb-1">{{ tpl.name }}</h3>
            <p class="text-xs text-text-muted">{{ tpl.description }}</p>
          </button>
        </div>
      </div>
    </template>

    <!-- Has Pages -->
    <template v-else>
      <!-- Stats Row -->
      <div class="grid grid-cols-4 gap-4">
        <StatsCard label="Total Pages" :value="totalPages" icon="web" icon-color="#9B7EE0" />
        <StatsCard label="Published" :value="publishedCount" icon="public" icon-color="#16A34A" />
        <StatsCard label="Draft" :value="draftCount" icon="edit_note" icon-color="#D97706" />
        <StatsCard label="Total Page Views" :value="totalViews" icon="visibility" icon-color="#3B82F6" />
      </div>

      <!-- Filter Tabs -->
      <div class="flex items-center gap-2">
        <button
          v-for="tab in filterTabs"
          :key="tab"
          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
          :class="filterTab === tab
            ? 'bg-wl text-white'
            : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
          @click="filterTab = tab"
        >
          {{ tab === 'all' ? 'All' : tab.charAt(0).toUpperCase() + tab.slice(1) }}
        </button>
      </div>

      <!-- Page Cards Grid -->
      <div class="grid grid-cols-3 gap-4">
        <!-- Create New Page card -->
        <button
          class="border-2 border-dashed border-border rounded-xl p-6 flex flex-col items-center justify-center text-center hover:border-wl/40 transition-colors min-h-55 group"
          @click="router.push('/b/whitelabel/pages/new')"
        >
          <span class="material-symbols-rounded text-3xl text-text-muted mb-2 group-hover:text-wl transition-colors">add_circle</span>
          <span class="text-sm font-medium text-text-secondary group-hover:text-text-primary transition-colors">Create New Page</span>
        </button>

        <!-- Page cards -->
        <div
          v-for="page in filteredPages"
          :key="page.id"
          class="bg-card-bg border border-border rounded-xl overflow-hidden hover:border-wl/30 transition-colors"
        >
          <!-- Thumbnail placeholder -->
          <div class="h-28 bg-page-bg border-b border-border flex items-center justify-center">
            <span class="material-symbols-rounded text-3xl text-text-muted">web</span>
          </div>

          <!-- Card body -->
          <div class="p-4">
            <div class="flex items-start justify-between mb-2">
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-semibold text-text-primary truncate">{{ page.name }}</h3>
                <p class="text-xs text-text-muted font-mono truncate mt-0.5">/{{ page.slug }}</p>
              </div>
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ml-2 shrink-0"
                :class="[statusBadge(page.status).bg, statusBadge(page.status).text]"
              >
                {{ statusBadge(page.status).label }}
              </span>
            </div>

            <div class="flex items-center justify-between text-xs text-text-muted mb-3">
              <span>{{ formatDate(page.last_modified) }}</span>
              <span class="flex items-center gap-1">
                <span class="material-symbols-rounded text-sm">visibility</span>
                {{ page.page_views.toLocaleString() }}
              </span>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-2">
              <button
                class="flex-1 px-3 py-1.5 text-xs font-medium text-wl border border-wl rounded-lg hover:bg-wl/10 transition-colors"
                @click="router.push(`/b/whitelabel/pages/${page.id}/edit`)"
              >
                Edit
              </button>
              <button
                class="px-3 py-1.5 text-xs font-medium text-text-secondary border border-border rounded-lg hover:bg-white/5 transition-colors"
                @click="router.push(`/b/whitelabel/pages/${page.id}/analytics`)"
              >
                Analytics
              </button>
              <button
                class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                title="Delete"
                @click="confirmDelete(page)"
              >
                <span class="material-symbols-rounded text-base text-status-paused">delete</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty filter results -->
      <div v-if="filteredPages.length === 0 && hasPages" class="text-center py-12">
        <span class="material-symbols-rounded text-3xl text-text-muted block mb-2">filter_list_off</span>
        <p class="text-sm text-text-muted">No pages match the current filter.</p>
      </div>
    </template>

    <!-- Delete Confirm Modal -->
    <Modal :open="showDeleteConfirm" title="Delete Page" @close="showDeleteConfirm = false">
      <p class="text-sm text-text-secondary">
        Are you sure you want to delete <strong class="text-text-primary">{{ deleteTarget?.name }}</strong>?
        This will permanently remove the page and its content. This action cannot be undone.
      </p>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showDeleteConfirm = false">Cancel</button>
        <button
          class="px-4 py-2 bg-status-paused text-white text-sm font-medium rounded-lg hover:bg-status-paused/90 disabled:opacity-50"
          :disabled="deleting"
          @click="deletePage"
        >
          {{ deleting ? 'Deleting...' : 'Delete Page' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
