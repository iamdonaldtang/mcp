<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()
const loading = ref(true)

interface ChecklistStep {
  id: string
  section: 'wizard' | 'enrich' | 'go_live'
  label: string
  status: 'completed' | 'in_progress' | 'pending'
  expandable: boolean
  hint?: string
}

const steps = ref<ChecklistStep[]>([
  { id: 'community_created', section: 'wizard', label: 'Community created with strategy', status: 'completed', expandable: false },
  { id: 'starter_tasks', section: 'wizard', label: '3 starter tasks live', status: 'completed', expandable: false },
  { id: 'points_configured', section: 'wizard', label: 'Points & Levels configured', status: 'completed', expandable: false },
  { id: 'add_tasks', section: 'enrich', label: 'Add more tasks', status: 'pending', expandable: true, hint: 'Enrich your community with diverse tasks' },
  { id: 'setup_shop', section: 'enrich', label: 'Set up your Benefits Shop', status: 'pending', expandable: true, hint: 'Give members rewards to spend points on' },
  { id: 'customize_daychain', section: 'enrich', label: 'Customize DayChain rewards', status: 'pending', expandable: true, hint: 'Configure daily check-in bonuses' },
  { id: 'preview', section: 'go_live', label: 'Preview your community as a user', status: 'pending', expandable: true },
  { id: 'share', section: 'go_live', label: 'Share with your community', status: 'pending', expandable: true },
  { id: 'first_10', section: 'go_live', label: 'First 10 participants', status: 'pending', expandable: true },
])

const expandedStep = ref<string | null>(null)
const participantCount = ref(0)
const shareUrl = ref('https://community.example.com/join')

const completedCount = computed(() => steps.value.filter(s => s.status === 'completed').length)
const totalCount = computed(() => steps.value.length)
const progressPercent = computed(() => Math.round(completedCount.value / totalCount.value * 100))
const allCompleted = computed(() => completedCount.value === totalCount.value)

const wizardSteps = computed(() => steps.value.filter(s => s.section === 'wizard'))
const enrichSteps = computed(() => steps.value.filter(s => s.section === 'enrich'))
const goLiveSteps = computed(() => steps.value.filter(s => s.section === 'go_live'))

// Active modules and available modules
const activeModules = ref<{ type: string; name: string; icon: string; hasData: boolean }[]>([])
const addModules = ref<{ type: string; name: string; icon: string; brief: string }[]>([])

onMounted(async () => {
  loading.value = true
  try {
    const [progressRes, modulesRes] = await Promise.allSettled([
      api.get('/api/v1/community/onboarding/progress'),
      api.get('/api/v1/community/modules'),
    ])
    if (progressRes.status === 'fulfilled') {
      const data = progressRes.value.data.data
      if (data?.steps) {
        for (const s of data.steps) {
          const found = steps.value.find(x => x.id === s.id)
          if (found) found.status = s.status
        }
      }
    }
  } finally {
    loading.value = false
  }
})

function toggleExpand(id: string) {
  expandedStep.value = expandedStep.value === id ? null : id
}

function copyLink() {
  navigator.clipboard.writeText(shareUrl.value)
  // TODO: toast "Link copied!"
}

function shareTwitter() {
  const text = encodeURIComponent(`Join our community on @TaskOnXYZ! Complete tasks, earn points, and level up. ${shareUrl.value} #Web3 #TaskOn`)
  window.open(`https://twitter.com/intent/tweet?text=${text}`, '_blank')
}

function shareTelegram() {
  const text = encodeURIComponent('Join our community! Complete tasks and earn rewards.')
  window.open(`https://t.me/share/url?url=${encodeURIComponent(shareUrl.value)}&text=${text}`, '_blank')
}

async function enableModule(type: string) {
  try {
    await api.put(`/api/v1/community/modules/${type}/enable`)
    // Refresh modules
  } catch { /* TODO: toast */ }
}
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <div class="flex items-center gap-3 mb-1">
        <h1 class="text-2xl font-bold text-text-primary">Getting Started</h1>
        <span class="px-2.5 py-0.5 text-xs font-medium rounded-full" style="background: #1F1A08; color: #D97706">Setting Up</span>
      </div>
      <p class="text-sm text-text-secondary">Complete these steps to get the most out of your community</p>
    </div>

    <!-- Checklist Card -->
    <div class="bg-card-bg border border-border rounded-xl p-6">
      <!-- Progress Header -->
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm font-semibold text-text-primary">Getting Started</span>
        <span class="text-sm text-text-secondary">{{ completedCount }} of {{ totalCount }} complete</span>
      </div>
      <div class="h-2 bg-page-bg rounded-full mb-6 overflow-hidden">
        <div class="h-full bg-community rounded-full transition-all duration-300" :style="{ width: progressPercent + '%' }"></div>
      </div>

      <!-- Completed banner -->
      <div v-if="allCompleted" class="mb-6 p-4 rounded-lg" style="background: #0A2E1A; border: 1px solid #16A34A">
        <div class="flex items-center gap-3">
          <span class="material-symbols-rounded text-status-active text-xl">celebration</span>
          <div>
            <div class="text-sm font-semibold text-text-primary">Congratulations! Your community is fully set up.</div>
            <button class="text-xs text-community hover:underline mt-1" @click="router.push('/b/dashboard')">Go to Dashboard →</button>
          </div>
        </div>
      </div>

      <!-- Wizard section -->
      <div class="mb-4">
        <div class="text-[10px] font-semibold text-text-muted uppercase tracking-wider mb-2">Completed by Wizard</div>
        <div class="space-y-1">
          <div v-for="s in wizardSteps" :key="s.id" class="flex items-center gap-3 py-2 px-3 rounded-lg">
            <span class="material-symbols-rounded text-lg text-status-active">check_circle</span>
            <span class="text-sm text-text-secondary line-through">{{ s.label }}</span>
          </div>
        </div>
      </div>

      <!-- Enrich section -->
      <div class="mb-4">
        <div class="text-[10px] font-semibold text-text-muted uppercase tracking-wider mb-2">Enrich Your Community</div>
        <div class="space-y-1">
          <div v-for="s in enrichSteps" :key="s.id">
            <button
              class="flex items-center gap-3 py-2 px-3 rounded-lg w-full text-left hover:bg-white/2 transition-colors"
              @click="s.expandable && toggleExpand(s.id)"
            >
              <span v-if="s.status === 'completed'" class="material-symbols-rounded text-lg text-status-active">check_circle</span>
              <span v-else class="w-[18px] h-[18px] rounded-full border-2 border-border flex-shrink-0"></span>
              <span class="text-sm flex-1" :class="s.status === 'completed' ? 'text-text-secondary line-through' : 'text-text-primary'">{{ s.label }}</span>
              <span v-if="s.expandable" class="material-symbols-rounded text-base text-text-muted transition-transform" :class="expandedStep === s.id ? 'rotate-180' : ''">keyboard_arrow_down</span>
            </button>
            <!-- Expanded content -->
            <div v-if="expandedStep === s.id" class="ml-10 mt-1 mb-2 p-3 bg-page-bg rounded-lg text-sm">
              <template v-if="s.id === 'add_tasks'">
                <p class="text-text-muted mb-2">3 tasks configured</p>
                <button class="text-community text-xs font-medium hover:underline" @click="router.push('/b/community/modules/sectors')">Go to Sectors & Tasks →</button>
              </template>
              <template v-else-if="s.id === 'setup_shop'">
                <p class="text-text-muted mb-2">Not configured yet — Add rewards your community members can redeem with points</p>
                <button class="text-community text-xs font-medium hover:underline" @click="router.push('/b/community/modules/benefits-shop')">Open Benefits Shop →</button>
              </template>
              <template v-else-if="s.id === 'customize_daychain'">
                <p class="text-text-muted mb-2">Configure daily check-in bonuses and streak milestones</p>
                <button class="text-community text-xs font-medium hover:underline" @click="router.push('/b/community/modules/day-chain')">Configure DayChain →</button>
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- Go Live section -->
      <div>
        <div class="text-[10px] font-semibold text-text-muted uppercase tracking-wider mb-2">Go Live</div>
        <div class="space-y-1">
          <div v-for="s in goLiveSteps" :key="s.id">
            <button
              class="flex items-center gap-3 py-2 px-3 rounded-lg w-full text-left hover:bg-white/2 transition-colors"
              @click="s.expandable && toggleExpand(s.id)"
            >
              <span v-if="s.status === 'completed'" class="material-symbols-rounded text-lg text-status-active">check_circle</span>
              <span v-else class="w-[18px] h-[18px] rounded-full border-2 border-border flex-shrink-0"></span>
              <span class="text-sm flex-1" :class="s.status === 'completed' ? 'text-text-secondary line-through' : 'text-text-primary'">{{ s.label }}</span>
              <span v-if="s.expandable" class="material-symbols-rounded text-base text-text-muted transition-transform" :class="expandedStep === s.id ? 'rotate-180' : ''">keyboard_arrow_down</span>
            </button>
            <div v-if="expandedStep === s.id" class="ml-10 mt-1 mb-2 p-3 bg-page-bg rounded-lg text-sm">
              <template v-if="s.id === 'preview'">
                <p class="text-text-muted mb-2">See your community as a participant</p>
                <button class="text-community text-xs font-medium hover:underline" @click="router.push('/b/community/preview')">Open Full Preview →</button>
              </template>
              <template v-else-if="s.id === 'share'">
                <div class="space-y-3">
                  <div class="flex items-center gap-2">
                    <input type="text" :value="shareUrl" readonly class="flex-1 px-3 py-1.5 bg-card-bg border border-border rounded-lg text-xs text-text-secondary" />
                    <button class="px-3 py-1.5 bg-community text-white text-xs font-medium rounded-lg hover:bg-community/90" @click="copyLink">Copy</button>
                  </div>
                  <div class="flex gap-2">
                    <button class="px-3 py-1.5 bg-[#1DA1F2]/20 text-[#1DA1F2] text-xs font-medium rounded-lg hover:bg-[#1DA1F2]/30" @click="shareTwitter">Twitter</button>
                    <button class="px-3 py-1.5 bg-[#5865F2]/20 text-[#5865F2] text-xs font-medium rounded-lg hover:bg-[#5865F2]/30">Discord</button>
                    <button class="px-3 py-1.5 bg-[#0088CC]/20 text-[#0088CC] text-xs font-medium rounded-lg hover:bg-[#0088CC]/30" @click="shareTelegram">Telegram</button>
                  </div>
                </div>
              </template>
              <template v-else-if="s.id === 'first_10'">
                <div class="flex items-center gap-2">
                  <span class="text-text-muted">Participants:</span>
                  <span class="text-lg font-bold text-text-primary">{{ participantCount }}/10</span>
                </div>
                <p class="text-xs text-text-muted mt-1">This updates automatically as users join</p>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Active Modules -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Active Modules</div>
      <div class="grid grid-cols-3 gap-4">
        <div v-for="mod in activeModules" :key="mod.type" class="bg-card-bg border border-community/30 rounded-xl p-4">
          <span class="material-symbols-rounded text-community text-xl mb-2 block">{{ mod.icon }}</span>
          <h4 class="text-sm font-medium text-text-primary mb-2">{{ mod.name }}</h4>
          <button class="text-xs text-community font-medium hover:underline">
            {{ mod.hasData ? 'Manage' : 'Configure' }} →
          </button>
        </div>
        <!-- Fallback if empty -->
        <template v-if="activeModules.length === 0">
          <div v-for="mod in [{ icon: 'category', name: 'Sectors & Tasks' }, { icon: 'stars', name: 'Points & Level' }, { icon: 'leaderboard', name: 'Leaderboard' }]" :key="mod.name" class="bg-card-bg border border-community/30 rounded-xl p-4">
            <span class="material-symbols-rounded text-community text-xl mb-2 block">{{ mod.icon }}</span>
            <h4 class="text-sm font-medium text-text-primary mb-2">{{ mod.name }}</h4>
            <button class="text-xs text-community font-medium hover:underline" @click="router.push('/b/community/modules/sectors')">Configure →</button>
          </div>
        </template>
      </div>
    </div>

    <!-- Add More Modules -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Add More Modules</div>
      <div class="grid grid-cols-4 gap-4">
        <div v-for="mod in [
          { type: 'task_chain', icon: 'link', name: 'TaskChain', brief: 'Multi-step task sequences' },
          { type: 'day_chain', icon: 'local_fire_department', name: 'DayChain', brief: 'Daily check-in streaks' },
          { type: 'milestone', icon: 'flag', name: 'Milestones', brief: 'Achievement thresholds' },
          { type: 'benefits_shop', icon: 'storefront', name: 'Benefits Shop', brief: 'Points-based rewards' },
        ]" :key="mod.type" class="bg-card-bg border border-border rounded-xl p-4">
          <span class="material-symbols-rounded text-text-muted text-xl mb-2 block">{{ mod.icon }}</span>
          <h4 class="text-sm font-medium text-text-primary mb-1">{{ mod.name }}</h4>
          <p class="text-xs text-text-muted mb-3">{{ mod.brief }}</p>
          <button class="text-xs text-community font-medium hover:underline" @click="enableModule(mod.type)">+ Enable</button>
        </div>
      </div>
    </div>

    <div class="border-t border-border"></div>

    <!-- Resources -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Resources</div>
      <div class="grid grid-cols-3 gap-4">
        <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-xl mb-2 block">menu_book</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Community Playbook</h4>
          <p class="text-xs text-text-secondary">Proven retention strategies for Web3</p>
        </a>
        <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-xl mb-2 block">psychology</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Points Strategy</h4>
          <p class="text-xs text-text-secondary">Design your points economy</p>
        </a>
        <a href="#" target="_blank" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-xl mb-2 block">open_in_new</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Learn More</h4>
          <p class="text-xs text-text-secondary">Explore the full Community product</p>
        </a>
      </div>
    </div>
  </div>
</template>
