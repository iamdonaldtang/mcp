<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const viewMode = ref<'desktop' | 'mobile'>('desktop')

// ─── Preview Data ───
interface PreviewData {
  community_name: string
  brand_color: string
  enabled_modules: string[]
  announcements: { title: string; content: string; type: string }[]
  mock_user: {
    name: string
    level: number
    points: number
    streak: number
  }
  tasks: {
    id: string
    name: string
    icon: string
    points: number
    status: 'available' | 'completed'
    sector: string
  }[]
}

const preview = ref<PreviewData>({
  community_name: 'My Community',
  brand_color: '#48BB78',
  enabled_modules: ['home', 'quests', 'leaderboard'],
  announcements: [],
  mock_user: { name: 'Preview User', level: 5, points: 1250, streak: 7 },
  tasks: [],
})

const activeTab = ref('home')
const toastMessage = ref('')
const toastTimeout = ref<ReturnType<typeof setTimeout> | null>(null)

// Nav tabs based on enabled modules
const navTabs = computed(() => {
  const tabs: { key: string; label: string }[] = []
  if (preview.value.enabled_modules.includes('home')) tabs.push({ key: 'home', label: 'Home' })
  if (preview.value.enabled_modules.includes('quests') || preview.value.enabled_modules.includes('sectors_tasks')) tabs.push({ key: 'quests', label: 'Quests' })
  if (preview.value.enabled_modules.includes('leaderboard')) tabs.push({ key: 'leaderboard', label: 'Leaderboard' })
  if (preview.value.enabled_modules.includes('lb_sprint')) tabs.push({ key: 'lb_sprint', label: 'LB Sprint' })
  if (preview.value.enabled_modules.includes('benefits_shop')) tabs.push({ key: 'shop', label: 'Shop' })
  // Always show at least Home
  if (tabs.length === 0) tabs.push({ key: 'home', label: 'Home' })
  return tabs
})

onMounted(async () => {
  loading.value = true
  try {
    const res = await api.get('/api/v1/community/preview')
    if (res.data.data) {
      preview.value = { ...preview.value, ...res.data.data }
    }
  } catch {
    // Use defaults — mock data for preview
    preview.value.tasks = [
      { id: '1', name: 'Follow us on Twitter', icon: 'share', points: 50, status: 'available', sector: 'Social' },
      { id: '2', name: 'Join Discord Server', icon: 'forum', points: 30, status: 'available', sector: 'Social' },
      { id: '3', name: 'Complete Profile Setup', icon: 'person', points: 20, status: 'completed', sector: 'Onboarding' },
    ]
    preview.value.announcements = [
      { title: 'Welcome to our community!', content: 'Complete tasks to earn points and unlock rewards.', type: 'general' },
    ]
  } finally {
    loading.value = false
  }
})

function interceptClick() {
  showToast('Links disabled in preview mode')
}

function showToast(message: string) {
  toastMessage.value = message
  if (toastTimeout.value) clearTimeout(toastTimeout.value)
  toastTimeout.value = setTimeout(() => { toastMessage.value = '' }, 2500)
}

function exitPreview() {
  const from = route.query.from as string
  if (from) {
    router.push(from)
  } else {
    router.back()
  }
}
</script>

<template>
  <div class="fixed inset-0 z-40 bg-[#0A0F1A] flex flex-col">
    <!-- ═══ Preview Banner ═══ -->
    <div class="h-12 flex-shrink-0 bg-[#1F1A08] border-b border-[#D97706] flex items-center justify-between px-6">
      <div class="flex items-center gap-2">
        <span class="material-symbols-rounded text-lg text-[#D97706]">warning</span>
        <span class="text-sm font-medium text-[#D97706]">Preview Mode — This is how your community looks to users</span>
      </div>
      <div class="flex items-center gap-3">
        <!-- Desktop / Mobile toggle -->
        <div class="flex bg-black/20 rounded-lg p-0.5">
          <button
            class="px-3 py-1 text-xs font-medium rounded-md transition-colors"
            :class="viewMode === 'desktop' ? 'bg-[#D97706] text-white' : 'text-[#D97706] hover:text-white'"
            @click="viewMode = 'desktop'"
          >
            <span class="material-symbols-rounded text-sm align-middle mr-1">desktop_windows</span>
            Desktop
          </button>
          <button
            class="px-3 py-1 text-xs font-medium rounded-md transition-colors"
            :class="viewMode === 'mobile' ? 'bg-[#D97706] text-white' : 'text-[#D97706] hover:text-white'"
            @click="viewMode = 'mobile'"
          >
            <span class="material-symbols-rounded text-sm align-middle mr-1">smartphone</span>
            Mobile
          </button>
        </div>
        <!-- Exit -->
        <button
          class="px-4 py-1.5 border border-[#D97706] text-[#D97706] text-xs font-medium rounded-lg hover:bg-[#D97706] hover:text-white transition-colors"
          @click="exitPreview"
        >
          Exit Preview
        </button>
      </div>
    </div>

    <!-- ═══ Preview Content Area ═══ -->
    <div class="flex-1 overflow-auto flex items-start justify-center" :class="viewMode === 'mobile' ? 'py-8' : ''">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center h-full">
        <div class="text-text-muted text-sm">Loading preview...</div>
      </div>

      <!-- C-End Simulation Frame -->
      <div
        v-else
        class="bg-[#0A0F1A] overflow-hidden flex flex-col"
        :class="viewMode === 'mobile' ? 'w-[375px] rounded-3xl border-2 border-[#1E293B] h-[720px]' : 'w-full h-full'"
      >
        <!-- C-End Header -->
        <div class="bg-[#0F172A] flex-shrink-0">
          <div class="px-6 py-4 flex items-center justify-between" :class="viewMode === 'mobile' ? 'px-4 py-3' : ''">
            <h1 class="text-lg font-bold text-white" :class="viewMode === 'mobile' ? 'text-base' : ''">{{ preview.community_name }}</h1>
            <!-- Mock user info -->
            <div class="flex items-center gap-2 cursor-pointer" @click="interceptClick">
              <div class="w-7 h-7 rounded-full bg-[#1E293B] flex items-center justify-center">
                <span class="material-symbols-rounded text-sm text-[#F59E0B]">person</span>
              </div>
              <span class="text-xs text-[#94A3B8]">{{ preview.mock_user.name }}</span>
            </div>
          </div>
          <!-- Nav tabs -->
          <div class="flex gap-0 px-6 border-b border-[#1E293B]" :class="viewMode === 'mobile' ? 'px-4 overflow-x-auto' : ''">
            <button
              v-for="tab in navTabs"
              :key="tab.key"
              class="px-4 py-2.5 text-sm font-medium transition-colors whitespace-nowrap"
              :class="activeTab === tab.key ? 'text-[#F59E0B] border-b-2 border-[#F59E0B]' : 'text-[#94A3B8] hover:text-white'"
              @click="activeTab = tab.key"
            >
              {{ tab.label }}
            </button>
          </div>
        </div>

        <!-- C-End Body -->
        <div class="flex-1 overflow-auto bg-[#0A0F1A]" :class="viewMode === 'mobile' ? 'px-4 py-4' : 'px-6 py-6'">
          <!-- User Status Bar -->
          <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4 mb-5 flex items-center justify-between" :class="viewMode === 'mobile' ? 'p-3 mb-4' : ''">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-[#1E293B] flex items-center justify-center">
                <span class="material-symbols-rounded text-lg text-[#F59E0B]">person</span>
              </div>
              <div>
                <div class="text-sm font-semibold text-white">{{ preview.mock_user.name }}</div>
                <div class="text-xs text-[#94A3B8]">Level {{ preview.mock_user.level }}</div>
              </div>
            </div>
            <div class="flex items-center gap-5" :class="viewMode === 'mobile' ? 'gap-3' : ''">
              <div class="text-center cursor-pointer" @click="interceptClick">
                <div class="text-sm font-bold text-[#F59E0B]">{{ preview.mock_user.points.toLocaleString() }}</div>
                <div class="text-[10px] text-[#64748B]">Points</div>
              </div>
              <div class="text-center cursor-pointer" @click="interceptClick">
                <div class="text-sm font-bold text-[#EF4444] flex items-center gap-0.5">
                  <span class="material-symbols-rounded text-sm">local_fire_department</span>
                  {{ preview.mock_user.streak }}
                </div>
                <div class="text-[10px] text-[#64748B]">Day Streak</div>
              </div>
            </div>
          </div>

          <!-- HOME TAB CONTENT -->
          <template v-if="activeTab === 'home'">
            <!-- Announcements -->
            <div v-if="preview.announcements.length > 0" class="mb-5">
              <div
                v-for="(ann, i) in preview.announcements"
                :key="i"
                class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4 cursor-pointer"
                @click="interceptClick"
              >
                <div class="flex items-center gap-2 mb-1">
                  <span class="material-symbols-rounded text-sm text-[#F59E0B]">campaign</span>
                  <span class="text-xs font-semibold text-white">{{ ann.title }}</span>
                </div>
                <p class="text-xs text-[#94A3B8]">{{ ann.content }}</p>
              </div>
            </div>

            <!-- Task Cards -->
            <div class="mb-4">
              <h3 class="text-xs font-semibold text-[#64748B] uppercase tracking-wider mb-3">Available Tasks</h3>
              <div class="space-y-3">
                <div
                  v-for="task in preview.tasks"
                  :key="task.id"
                  class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4 flex items-center justify-between cursor-pointer"
                  @click="interceptClick"
                >
                  <div class="flex items-center gap-3">
                    <div
                      class="w-10 h-10 rounded-lg flex items-center justify-center"
                      :style="{ background: task.status === 'completed' ? '#0A2E1A' : '#1E293B' }"
                    >
                      <span
                        class="material-symbols-rounded"
                        :style="{ color: task.status === 'completed' ? '#48BB78' : '#F59E0B' }"
                      >{{ task.status === 'completed' ? 'check_circle' : task.icon }}</span>
                    </div>
                    <div>
                      <div class="text-sm font-medium" :class="task.status === 'completed' ? 'text-[#64748B] line-through' : 'text-white'">{{ task.name }}</div>
                      <div class="text-[11px] text-[#64748B]">{{ task.sector }}</div>
                    </div>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-bold text-[#F59E0B]">+{{ task.points }}</span>
                    <button
                      v-if="task.status === 'available'"
                      class="px-3 py-1 bg-[#F59E0B] text-black text-xs font-semibold rounded-lg"
                      @click.stop="interceptClick"
                    >
                      Do
                    </button>
                    <span v-else class="material-symbols-rounded text-lg text-[#48BB78]">check</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Points Summary -->
            <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4 cursor-pointer" @click="interceptClick">
              <div class="flex items-center gap-2 mb-2">
                <span class="material-symbols-rounded text-sm text-[#F59E0B]">stars</span>
                <span class="text-xs font-semibold text-white">Your Progress</span>
              </div>
              <div class="flex items-center justify-between">
                <div>
                  <div class="text-2xl font-bold text-[#F59E0B]">{{ preview.mock_user.points.toLocaleString() }}</div>
                  <div class="text-[11px] text-[#64748B]">Total Points Earned</div>
                </div>
                <div class="w-20 h-20">
                  <!-- Simple circular progress placeholder -->
                  <div class="w-full h-full rounded-full border-4 border-[#1E293B] flex items-center justify-center relative">
                    <div class="absolute inset-0 rounded-full border-4 border-transparent border-t-[#F59E0B] border-r-[#F59E0B]" style="transform: rotate(45deg)"></div>
                    <span class="text-sm font-bold text-white z-10">Lv.{{ preview.mock_user.level }}</span>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- QUESTS TAB -->
          <template v-else-if="activeTab === 'quests'">
            <div class="space-y-3">
              <div
                v-for="task in preview.tasks"
                :key="task.id"
                class="bg-[#111B27] border border-[#1E293B] rounded-xl p-4 flex items-center justify-between cursor-pointer"
                @click="interceptClick"
              >
                <div class="flex items-center gap-3">
                  <div
                    class="w-10 h-10 rounded-lg flex items-center justify-center"
                    :style="{ background: task.status === 'completed' ? '#0A2E1A' : '#1E293B' }"
                  >
                    <span
                      class="material-symbols-rounded"
                      :style="{ color: task.status === 'completed' ? '#48BB78' : '#F59E0B' }"
                    >{{ task.status === 'completed' ? 'check_circle' : task.icon }}</span>
                  </div>
                  <div>
                    <div class="text-sm font-medium" :class="task.status === 'completed' ? 'text-[#64748B] line-through' : 'text-white'">{{ task.name }}</div>
                    <div class="text-[11px] text-[#64748B]">{{ task.sector }}</div>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs font-bold text-[#F59E0B]">+{{ task.points }}</span>
                  <button
                    v-if="task.status === 'available'"
                    class="px-3 py-1 bg-[#F59E0B] text-black text-xs font-semibold rounded-lg"
                    @click.stop="interceptClick"
                  >
                    Do
                  </button>
                  <span v-else class="material-symbols-rounded text-lg text-[#48BB78]">check</span>
                </div>
              </div>
            </div>
            <div v-if="preview.tasks.length === 0" class="text-center py-12">
              <span class="material-symbols-rounded text-4xl text-[#64748B] block mb-2">inbox</span>
              <p class="text-sm text-[#64748B]">No quests available</p>
            </div>
          </template>

          <!-- LEADERBOARD TAB -->
          <template v-else-if="activeTab === 'leaderboard'">
            <div class="bg-[#111B27] border border-[#1E293B] rounded-xl overflow-hidden">
              <div class="px-4 py-3 border-b border-[#1E293B] flex items-center justify-between">
                <span class="text-xs font-semibold text-white">Weekly Ranking</span>
                <span class="text-[11px] text-[#64748B] cursor-pointer" @click="interceptClick">All Time</span>
              </div>
              <!-- Mock leaderboard entries -->
              <div
                v-for="i in 5"
                :key="i"
                class="px-4 py-3 flex items-center justify-between border-b border-[#1E293B] last:border-b-0 cursor-pointer"
                @click="interceptClick"
              >
                <div class="flex items-center gap-3">
                  <span
                    class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold"
                    :class="i <= 3 ? 'bg-[#F59E0B]/20 text-[#F59E0B]' : 'bg-[#1E293B] text-[#64748B]'"
                  >{{ i }}</span>
                  <div class="w-7 h-7 rounded-full bg-[#1E293B]"></div>
                  <span class="text-sm text-white" :class="i === 3 ? 'font-semibold' : ''">
                    {{ i === 3 ? preview.mock_user.name + ' (You)' : `User ${i}` }}
                  </span>
                </div>
                <span class="text-xs font-bold text-[#F59E0B]">{{ (6000 - i * 800).toLocaleString() }} pts</span>
              </div>
            </div>
          </template>

          <!-- OTHER TABS (placeholder) -->
          <template v-else>
            <div class="text-center py-12">
              <span class="material-symbols-rounded text-4xl text-[#64748B] block mb-2">construction</span>
              <p class="text-sm text-[#64748B]">{{ activeTab.replace('_', ' ') }} preview coming soon</p>
            </div>
          </template>
        </div>

        <!-- C-End Footer -->
        <div class="flex-shrink-0 py-3 text-center border-t border-[#1E293B] bg-[#0F172A]">
          <span class="text-[11px] text-[#64748B]">Powered by <span class="text-[#94A3B8] font-medium">TaskOn</span></span>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="translate-y-4 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-4 opacity-0"
    >
      <div
        v-if="toastMessage"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 px-4 py-2 bg-[#1E293B] border border-[#334155] rounded-lg shadow-lg"
      >
        <span class="text-sm text-[#F1F5F9]">{{ toastMessage }}</span>
      </div>
    </Transition>
  </div>
</template>
