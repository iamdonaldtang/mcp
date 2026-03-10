<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()
const loading = ref(true)
const publishing = ref(false)
const slugChecking = ref(false)
const slugAvailable = ref<boolean | null>(null)
const previewMode = ref<'desktop' | 'mobile'>('desktop')

// ── Draft data ──
interface DraftData {
  name: string
  description: string
  strategy: string
  enabled_modules: string[]
  content?: {
    sectors?: { id: string; name: string; tasks: { id: string; name: string; type: string; xp: number }[] }[]
    levels?: { id: string; name: string; threshold: number }[]
    taskchain?: { name: string; steps: { id: string; name: string; xp: number }[] }
    daychain?: { base_reward: number }
  }
}

const draft = ref<DraftData>({
  name: '',
  description: '',
  strategy: '',
  enabled_modules: [],
  content: {},
})

const slug = ref('')
let slugDebounce: ReturnType<typeof setTimeout> | null = null

// ── Readiness checklist ──
interface CheckItem {
  key: string
  label: string
  done: boolean
  icon: string
}

const checklist = computed<CheckItem[]>(() => [
  {
    key: 'info',
    label: 'Community info complete',
    done: !!(draft.value.name && draft.value.description),
    icon: draft.value.name && draft.value.description ? 'check_circle' : 'cancel',
  },
  {
    key: 'tasks',
    label: 'At least 1 task created',
    done: taskCount.value > 0,
    icon: taskCount.value > 0 ? 'check_circle' : 'cancel',
  },
  {
    key: 'points',
    label: 'Point system configured',
    done: levelCount.value >= 2,
    icon: levelCount.value >= 2 ? 'check_circle' : 'cancel',
  },
  {
    key: 'subscription',
    label: 'Subscription active (or trial)',
    done: subscriptionActive.value,
    icon: subscriptionActive.value ? 'check_circle' : 'cancel',
  },
])

const allReady = computed(() => checklist.value.every(c => c.done) && slugAvailable.value === true)

const taskCount = computed(() => {
  const sectors = draft.value.content?.sectors || []
  return sectors.reduce((sum, s) => sum + (s.tasks?.length || 0), 0)
})

const levelCount = computed(() => draft.value.content?.levels?.length || 0)

const subscriptionActive = ref(false)

// ── Module helpers ──
function hasModule(key: string) {
  return draft.value.enabled_modules.includes(key)
}

const previewTabs = computed(() => {
  const tabs: { key: string; label: string }[] = [{ key: 'home', label: 'Home' }]
  if (hasModule('sectors_tasks')) tabs.push({ key: 'quests', label: 'Quests' })
  if (hasModule('points_level') || hasModule('taskchain')) tabs.push({ key: 'leaderboard', label: 'Leaderboard' })
  return tabs
})

const activePreviewTab = ref('home')

// ── Slug validation ──
function generateSlug(name: string): string {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
    .slice(0, 48)
}

watch(slug, (val) => {
  if (slugDebounce) clearTimeout(slugDebounce)
  slugAvailable.value = null
  if (!val || val.length < 3) {
    slugAvailable.value = false
    return
  }
  slugChecking.value = true
  slugDebounce = setTimeout(async () => {
    try {
      const res = await api.get('/api/v1/community/slug/check', { params: { slug: val } })
      slugAvailable.value = res.data.data?.available ?? false
    } catch {
      slugAvailable.value = false
    } finally {
      slugChecking.value = false
    }
  }, 500)
})

// ── API ──
onMounted(async () => {
  try {
    const [draftRes, subRes] = await Promise.all([
      api.get('/api/v1/community/wizard/draft'),
      api.get('/api/v1/subscription/status'),
    ])

    draft.value = draftRes.data.data as DraftData
    slug.value = generateSlug(draft.value.name)

    const subData = subRes.data.data
    subscriptionActive.value = subData?.status === 'active' || subData?.status === 'trial'
  } catch {
    // Fallback for dev: assume draft exists with defaults
    draft.value.enabled_modules = ['sectors_tasks', 'points_level']
    slug.value = 'my-community'
  } finally {
    loading.value = false
  }
})

async function publishCommunity() {
  if (!allReady.value) return
  publishing.value = true
  try {
    await api.post('/api/v1/community/publish', {
      slug: slug.value,
    })
    router.push({ name: 'B10' })
  } catch {
    // TODO: toast error
    publishing.value = false
  }
}

async function goBack() {
  router.push({ name: 'B13-3' })
}

// ── Stepper ──
const steps = [
  { num: 1, label: 'Customize' },
  { num: 2, label: 'Modules' },
  { num: 3, label: 'Quick Setup' },
  { num: 4, label: 'Preview & Publish' },
]

// ── Mock preview data ──
const mockTasks = computed(() => {
  const sectors = draft.value.content?.sectors || []
  const tasks: { name: string; xp: number; type: string }[] = []
  for (const s of sectors) {
    for (const t of s.tasks || []) {
      tasks.push({ name: t.name, xp: t.xp, type: t.type })
    }
  }
  return tasks.slice(0, 5)
})

const mockLevels = computed(() => {
  return (draft.value.content?.levels || []).slice(0, 5)
})

const mockLeaderboard = [
  { rank: 1, name: 'CryptoWhale.eth', xp: 2450 },
  { rank: 2, name: '0xAlice...3f2a', xp: 1820 },
  { rank: 3, name: 'DeFiKing.eth', xp: 1540 },
  { rank: 4, name: '0xBob...9c1d', xp: 980 },
  { rank: 5, name: 'Web3Dev.eth', xp: 720 },
]

const taskTypeIcons: Record<string, string> = {
  social: 'share',
  onchain: 'account_balance_wallet',
  verification: 'verified_user',
  custom: 'extension',
  recurring: 'repeat',
  referral: 'group_add',
}
</script>

<template>
  <div class="min-h-screen bg-page-bg">
    <!-- Top Bar -->
    <div class="sticky top-0 z-30 bg-header-bg border-b border-border px-6 py-3 flex items-center justify-between">
      <button
        class="flex items-center gap-1.5 text-sm text-text-secondary hover:text-text-primary transition-colors"
        @click="goBack"
      >
        <span class="material-symbols-rounded text-lg">arrow_back</span>
        Back
      </button>
      <h1 class="text-base font-semibold text-text-primary">Create Community</h1>
      <div class="w-20" />
    </div>

    <!-- Stepper -->
    <div class="flex items-center justify-center gap-0 py-6 px-6 bg-page-bg">
      <template v-for="(step, idx) in steps" :key="step.num">
        <div class="flex items-center gap-2">
          <div
            class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-semibold transition-colors"
            :class="
              step.num < 4
                ? 'bg-community text-white'
                : 'bg-community text-white ring-2 ring-community/30'
            "
          >
            <span v-if="step.num < 4" class="material-symbols-rounded text-sm">check</span>
            <span v-else>{{ step.num }}</span>
          </div>
          <span
            class="text-sm font-medium"
            :class="step.num === 4 ? 'text-text-primary' : 'text-community'"
          >
            {{ step.label }}
          </span>
        </div>
        <div
          v-if="idx < steps.length - 1"
          class="w-16 h-px mx-3"
          :class="step.num < 4 ? 'bg-community' : 'bg-border'"
        />
      </template>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="flex items-center gap-3 text-text-muted">
        <span class="material-symbols-rounded animate-spin text-2xl">progress_activity</span>
        Loading preview...
      </div>
    </div>

    <!-- Main Content -->
    <div v-else class="max-w-6xl mx-auto px-6 pb-24 flex gap-6">
      <!-- Left: Preview Panel -->
      <div class="flex-1">
        <!-- Preview toolbar -->
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-sm font-semibold text-text-primary">Community Preview</h2>
          <div class="flex items-center gap-1 bg-card-bg border border-border rounded-lg p-0.5">
            <button
              class="px-3 py-1 text-xs font-medium rounded-md transition-colors"
              :class="previewMode === 'desktop' ? 'bg-community text-white' : 'text-text-muted hover:text-text-primary'"
              @click="previewMode = 'desktop'"
            >
              <span class="material-symbols-rounded text-sm align-middle mr-1">monitor</span>
              Desktop
            </button>
            <button
              class="px-3 py-1 text-xs font-medium rounded-md transition-colors"
              :class="previewMode === 'mobile' ? 'bg-community text-white' : 'text-text-muted hover:text-text-primary'"
              @click="previewMode = 'mobile'"
            >
              <span class="material-symbols-rounded text-sm align-middle mr-1">smartphone</span>
              Mobile
            </button>
          </div>
        </div>

        <!-- Preview frame -->
        <div class="flex justify-center">
          <div
            class="bg-page-bg border border-border overflow-hidden relative transition-all duration-300"
            :class="previewMode === 'mobile' ? 'w-[375px] rounded-3xl' : 'w-full rounded-xl'"
            :style="previewMode === 'mobile' ? 'min-height: 667px' : 'min-height: 500px'"
          >
            <!-- Preview Mode watermark -->
            <div class="absolute inset-0 flex items-center justify-center pointer-events-none z-10">
              <span class="text-text-muted/10 text-6xl font-bold rotate-[-15deg] select-none">Preview Mode</span>
            </div>

            <!-- C-end mock: Dark header -->
            <div class="bg-header-bg px-4 py-3 border-b border-border relative z-20">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm font-bold text-text-primary">{{ draft.name || 'My Community' }}</span>
                <div class="flex items-center gap-2">
                  <div class="w-6 h-6 rounded-full bg-community/20 flex items-center justify-center">
                    <span class="material-symbols-rounded text-xs text-community">person</span>
                  </div>
                </div>
              </div>
              <!-- Nav tabs -->
              <div class="flex gap-4">
                <button
                  v-for="tab in previewTabs"
                  :key="tab.key"
                  class="pb-2 text-xs font-medium border-b-2 transition-colors"
                  :class="activePreviewTab === tab.key
                    ? 'text-c-accent border-c-accent'
                    : 'text-text-muted border-transparent hover:text-text-secondary'"
                  @click="activePreviewTab = tab.key"
                >
                  {{ tab.label }}
                </button>
              </div>
            </div>

            <!-- C-end mock: Content area -->
            <div class="p-4 relative z-20">
              <!-- Home tab -->
              <div v-if="activePreviewTab === 'home'" class="space-y-3">
                <!-- Welcome banner -->
                <div class="bg-community/10 border border-community/20 rounded-xl p-4">
                  <h3 class="text-sm font-semibold text-text-primary mb-1">Welcome to {{ draft.name || 'the community' }}</h3>
                  <p class="text-xs text-text-secondary">{{ draft.description || 'Complete tasks, earn points, and unlock rewards.' }}</p>
                </div>

                <!-- Points display -->
                <div v-if="hasModule('points_level')" class="bg-card-bg border border-border rounded-lg p-3 flex items-center justify-between">
                  <div>
                    <span class="text-[10px] text-text-muted uppercase tracking-wider">Your Points</span>
                    <div class="text-lg font-bold text-c-accent">0 XP</div>
                  </div>
                  <div class="text-right">
                    <span class="text-[10px] text-text-muted uppercase tracking-wider">Level</span>
                    <div class="text-sm font-semibold text-text-primary">{{ mockLevels[0]?.name || 'Bronze' }}</div>
                  </div>
                </div>

                <!-- DayChain streak -->
                <div v-if="hasModule('daychain')" class="bg-card-bg border border-border rounded-lg p-3">
                  <div class="flex items-center gap-2 mb-2">
                    <span class="material-symbols-rounded text-sm text-boost">local_fire_department</span>
                    <span class="text-xs font-medium text-text-primary">Daily Check-in</span>
                    <span class="text-[10px] text-text-muted ml-auto">Day 0</span>
                  </div>
                  <div class="flex gap-1.5">
                    <div
                      v-for="d in 7"
                      :key="d"
                      class="flex-1 h-7 rounded bg-border flex items-center justify-center text-[10px] text-text-muted"
                    >
                      {{ d }}
                    </div>
                  </div>
                </div>

                <!-- Quick tasks -->
                <div v-if="mockTasks.length > 0">
                  <div class="flex items-center justify-between mb-2">
                    <span class="text-xs font-semibold text-text-primary">Quick Tasks</span>
                    <span class="text-[10px] text-community">View All</span>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="task in mockTasks.slice(0, 3)"
                      :key="task.name"
                      class="bg-card-bg border border-border rounded-lg p-3 flex items-center gap-3"
                    >
                      <div class="w-8 h-8 rounded-lg bg-community/10 flex items-center justify-center">
                        <span class="material-symbols-rounded text-sm text-community">{{ taskTypeIcons[task.type] || 'task_alt' }}</span>
                      </div>
                      <div class="flex-1 min-w-0">
                        <div class="text-xs font-medium text-text-primary truncate">{{ task.name }}</div>
                        <div class="text-[10px] text-c-accent">+{{ task.xp }} XP</div>
                      </div>
                      <button class="px-2.5 py-1 text-[10px] font-medium bg-community/20 text-community rounded-md">
                        Start
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Quests tab -->
              <div v-if="activePreviewTab === 'quests'" class="space-y-2">
                <div
                  v-for="task in mockTasks"
                  :key="task.name"
                  class="bg-card-bg border border-border rounded-lg p-3 flex items-center gap-3"
                >
                  <div class="w-8 h-8 rounded-lg bg-community/10 flex items-center justify-center">
                    <span class="material-symbols-rounded text-sm text-community">{{ taskTypeIcons[task.type] || 'task_alt' }}</span>
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="text-xs font-medium text-text-primary truncate">{{ task.name }}</div>
                    <div class="text-[10px] text-c-accent">+{{ task.xp }} XP</div>
                  </div>
                  <button class="px-2.5 py-1 text-[10px] font-medium bg-community/20 text-community rounded-md">
                    Start
                  </button>
                </div>
                <div v-if="mockTasks.length === 0" class="text-center py-8">
                  <span class="material-symbols-rounded text-3xl text-text-muted block mb-1">task_alt</span>
                  <p class="text-xs text-text-muted">No tasks configured yet</p>
                </div>
              </div>

              <!-- Leaderboard tab -->
              <div v-if="activePreviewTab === 'leaderboard'" class="space-y-2">
                <div
                  v-for="entry in mockLeaderboard"
                  :key="entry.rank"
                  class="bg-card-bg border border-border rounded-lg px-3 py-2.5 flex items-center gap-3"
                >
                  <div
                    class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold"
                    :class="
                      entry.rank === 1 ? 'bg-c-accent/20 text-c-accent' :
                      entry.rank === 2 ? 'bg-text-secondary/20 text-text-secondary' :
                      entry.rank === 3 ? 'bg-boost/20 text-boost' :
                      'bg-border text-text-muted'
                    "
                  >
                    {{ entry.rank }}
                  </div>
                  <span class="text-xs text-text-primary flex-1">{{ entry.name }}</span>
                  <span class="text-xs text-c-accent font-medium">{{ entry.xp.toLocaleString() }} XP</span>
                </div>
              </div>
            </div>

            <!-- C-end mock: Footer -->
            <div class="absolute bottom-0 left-0 right-0 border-t border-border bg-header-bg px-4 py-2 text-center z-20">
              <span class="text-[10px] text-text-muted">Powered by TaskOn</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Settings & Publish -->
      <div class="w-80 shrink-0 space-y-5">
        <!-- Community URL -->
        <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <div class="px-5 py-4 border-b border-border">
            <h3 class="text-sm font-semibold text-text-primary">Community URL</h3>
          </div>
          <div class="p-5">
            <div class="flex items-center gap-0 bg-page-bg border border-border rounded-lg overflow-hidden focus-within:border-community transition-colors">
              <span class="px-3 py-2.5 text-xs text-text-muted bg-page-bg border-r border-border whitespace-nowrap">
                taskon.xyz/c/
              </span>
              <input
                v-model="slug"
                type="text"
                placeholder="my-community"
                class="flex-1 px-3 py-2.5 bg-transparent text-sm text-text-primary placeholder-text-muted focus:outline-none"
                @input="slug = slug.toLowerCase().replace(/[^a-z0-9-]/g, '')"
              />
            </div>
            <!-- Slug status -->
            <div class="mt-2 flex items-center gap-1.5 h-5">
              <template v-if="slugChecking">
                <span class="material-symbols-rounded animate-spin text-sm text-text-muted">progress_activity</span>
                <span class="text-xs text-text-muted">Checking availability...</span>
              </template>
              <template v-else-if="slug.length > 0 && slug.length < 3">
                <span class="material-symbols-rounded text-sm text-status-paused">error</span>
                <span class="text-xs text-status-paused">Minimum 3 characters</span>
              </template>
              <template v-else-if="slugAvailable === true">
                <span class="material-symbols-rounded text-sm text-community">check_circle</span>
                <span class="text-xs text-community">Available</span>
              </template>
              <template v-else-if="slugAvailable === false && slug.length >= 3">
                <span class="material-symbols-rounded text-sm text-status-paused">cancel</span>
                <span class="text-xs text-status-paused">Not available — try another</span>
              </template>
            </div>
          </div>
        </div>

        <!-- Readiness Checklist -->
        <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <div class="px-5 py-4 border-b border-border">
            <h3 class="text-sm font-semibold text-text-primary">Readiness Checklist</h3>
            <p class="text-xs text-text-muted mt-0.5">
              {{ checklist.filter(c => c.done).length }} of {{ checklist.length }} ready
            </p>
          </div>

          <div class="p-4 space-y-3">
            <div
              v-for="item in checklist"
              :key="item.key"
              class="flex items-start gap-2.5"
            >
              <span
                class="material-symbols-rounded text-lg mt-px"
                :class="item.done ? 'text-community' : 'text-status-paused'"
              >
                {{ item.icon }}
              </span>
              <span
                class="text-sm leading-5"
                :class="item.done ? 'text-text-primary' : 'text-text-muted'"
              >
                {{ item.label }}
              </span>
            </div>
          </div>

          <!-- Progress -->
          <div class="px-4 pb-4">
            <div class="w-full h-1.5 bg-border rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500"
                :class="allReady ? 'bg-community' : 'bg-c-accent'"
                :style="{ width: `${(checklist.filter(c => c.done).length / checklist.length) * 100}%` }"
              />
            </div>
          </div>
        </div>

        <!-- Publish button -->
        <button
          class="w-full py-3 text-sm font-semibold rounded-xl transition-all flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          :class="allReady
            ? 'bg-community text-white hover:bg-community/90 shadow-lg shadow-community/20'
            : 'bg-card-bg border border-border text-text-muted'"
          :disabled="!allReady || publishing"
          @click="publishCommunity"
        >
          <span v-if="publishing" class="material-symbols-rounded animate-spin text-base">progress_activity</span>
          <span class="material-symbols-rounded text-base" v-else>rocket_launch</span>
          {{ publishing ? 'Publishing...' : 'Publish Community' }}
        </button>

        <p v-if="!allReady" class="text-xs text-text-muted text-center">
          Complete all checklist items to publish
        </p>
      </div>
    </div>

    <!-- Bottom Action Bar -->
    <div v-if="!loading" class="fixed bottom-0 left-0 right-0 bg-header-bg border-t border-border px-6 py-3 flex items-center justify-between z-20">
      <button
        class="px-5 py-2 text-sm text-text-secondary border border-border rounded-lg hover:text-text-primary hover:border-text-muted transition-colors"
        @click="goBack"
      >
        Back
      </button>
      <button
        class="px-5 py-2 text-sm font-semibold rounded-lg transition-all flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed"
        :class="allReady
          ? 'bg-community text-white hover:bg-community/90'
          : 'bg-card-bg border border-border text-text-muted'"
        :disabled="!allReady || publishing"
        @click="publishCommunity"
      >
        <span v-if="publishing" class="material-symbols-rounded animate-spin text-sm">progress_activity</span>
        <span class="material-symbols-rounded text-base" v-else>rocket_launch</span>
        {{ publishing ? 'Publishing...' : 'Publish Community' }}
      </button>
    </div>
  </div>
</template>
