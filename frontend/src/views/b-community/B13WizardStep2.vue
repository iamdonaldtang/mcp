<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'
import type { CommunityWizardData } from '../../types/b-community'
import type { CommunityModule } from '../../types/common'

const router = useRouter()

// === Stepper ===
const steps = [
  { num: 1, label: 'Customize' },
  { num: 2, label: 'Modules' },
  { num: 3, label: 'Quick Setup' },
  { num: 4, label: 'Preview & Publish' },
]
const currentStep = 2

// === Module Definitions ===
interface ModuleDef {
  id: CommunityModule
  name: string
  description: string
  icon: string
  required: boolean
  system: 'task_engine' | 'points_recognition' | 'incentive_campaigns' | 'rewards_economy'
}

const allModules: ModuleDef[] = [
  // Task Engine
  { id: 'sectors_tasks', name: 'Sectors & Tasks', description: 'Core task infrastructure — create sectors and assign tasks for users to complete', icon: 'task_alt', required: true, system: 'task_engine' },
  { id: 'task_chain', name: 'TaskChain', description: 'Multi-step task sequences that guide users through onboarding or key actions', icon: 'link', required: false, system: 'task_engine' },
  { id: 'day_chain', name: 'DayChain', description: 'Daily check-in streaks with escalating multipliers to drive habitual engagement', icon: 'local_fire_department', required: false, system: 'task_engine' },
  // Points & Recognition
  { id: 'points_level', name: 'Points & Level', description: 'Core point system with custom point types and level progression', icon: 'stars', required: true, system: 'points_recognition' },
  { id: 'leaderboard', name: 'Leaderboard', description: 'Recurring rankings (weekly/monthly/all-time) by point type to fuel competition', icon: 'leaderboard', required: false, system: 'points_recognition' },
  { id: 'badges', name: 'Badges', description: 'Achievement badges awarded for milestones, streaks, or special actions', icon: 'military_tech', required: false, system: 'points_recognition' },
  // Incentive Campaigns
  { id: 'lb_sprint', name: 'LB Sprint', description: 'Time-bounded competitions with NFT, token, or whitelist spot rewards', icon: 'sprint', required: false, system: 'incentive_campaigns' },
  { id: 'milestone', name: 'Milestones', description: 'One-time achievement rewards when users hit key thresholds', icon: 'flag', required: false, system: 'incentive_campaigns' },
  // Rewards Economy
  { id: 'benefits_shop', name: 'Benefits Shop', description: 'Redemption store where users spend earned points on real rewards', icon: 'storefront', required: false, system: 'rewards_economy' },
  { id: 'lucky_wheel', name: 'Lucky Wheel', description: 'Gamified spin-to-win mechanic that adds variable reward excitement', icon: 'casino', required: false, system: 'rewards_economy' },
]

const systemLabels: Record<string, { label: string; icon: string }> = {
  task_engine: { label: 'Task Engine', icon: 'engineering' },
  points_recognition: { label: 'Points & Recognition', icon: 'emoji_events' },
  incentive_campaigns: { label: 'Incentive Campaigns', icon: 'campaign' },
  rewards_economy: { label: 'Rewards Economy', icon: 'paid' },
}

const systemOrder = ['task_engine', 'points_recognition', 'incentive_campaigns', 'rewards_economy'] as const

// === Strategy Definitions ===
interface StrategyDef {
  id: CommunityWizardData['step2']['strategy'] | 'blank'
  name: string
  description: string
  color: string
  icon: string
  modules: CommunityModule[]
}

const strategies: StrategyDef[] = [
  {
    id: 'activate',
    name: 'Activate New Users',
    description: 'TaskChain + Points + Levels — Goal gradient + instant feedback',
    color: '#48BB78',
    icon: 'rocket_launch',
    modules: ['sectors_tasks', 'points_level', 'task_chain'],
  },
  {
    id: 'engage',
    name: 'Drive Daily Engagement',
    description: 'DayChain + Leaderboard + Sprint — Loss aversion + social comparison',
    color: '#ED8936',
    icon: 'trending_up',
    modules: ['sectors_tasks', 'points_level', 'day_chain', 'leaderboard'],
  },
  {
    id: 'retain',
    name: 'Maximize Retention',
    description: 'Benefits Shop + Milestones + DayChain — Sunk cost + variable reinforcement',
    color: '#9B7EE0',
    icon: 'loyalty',
    modules: ['sectors_tasks', 'points_level', 'day_chain', 'milestone', 'benefits_shop'],
  },
  {
    id: 'blank',
    name: 'Start Blank',
    description: 'Only required modules — build your own configuration from scratch',
    color: '#64748B',
    icon: 'edit_note',
    modules: ['sectors_tasks', 'points_level'],
  },
]

// === State ===
const selectedStrategy = ref<string | null>(null)
const enabledModules = ref<Set<CommunityModule>>(new Set(['sectors_tasks', 'points_level']))
const loading = ref(true)
const saving = ref(false)

// === Computed ===
const enabledCount = computed(() => enabledModules.value.size)
const optionalEnabled = computed(() => [...enabledModules.value].filter(m => !allModules.find(d => d.id === m)?.required).length)

const estimatedTime = computed(() => {
  const base = 5 // minutes for required modules
  const extra = optionalEnabled.value * 3
  return base + extra
})

const modulesBySystem = computed(() => {
  const map: Record<string, ModuleDef[]> = {}
  for (const sys of systemOrder) {
    map[sys] = allModules.filter(m => m.system === sys)
  }
  return map
})

// === Strategy Selection ===
function selectStrategy(strategy: StrategyDef) {
  selectedStrategy.value = strategy.id
  // Reset to only required modules first
  enabledModules.value = new Set(['sectors_tasks', 'points_level'])
  // Then add strategy modules
  for (const mod of strategy.modules) {
    enabledModules.value.add(mod)
  }
}

// === Module Toggle ===
function toggleModule(mod: ModuleDef) {
  if (mod.required) return
  if (enabledModules.value.has(mod.id)) {
    enabledModules.value.delete(mod.id)
  } else {
    enabledModules.value.add(mod.id)
  }
  // Clear strategy selection when manually toggling
  selectedStrategy.value = null
}

function isModuleEnabled(id: CommunityModule): boolean {
  return enabledModules.value.has(id)
}

// === Draft ===
async function loadDraft() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/community/wizard/draft')
    const draft = res.data?.data
    if (draft?.step2) {
      selectedStrategy.value = draft.step2.strategy || null
      if (draft.step2.enabledModules?.length) {
        enabledModules.value = new Set(draft.step2.enabledModules)
        // Always include required modules
        enabledModules.value.add('sectors_tasks')
        enabledModules.value.add('points_level')
      }
    }
  } catch {
    // No draft — defaults are fine
  } finally {
    loading.value = false
  }
}

async function saveDraft() {
  saving.value = true
  try {
    const strategyValue = selectedStrategy.value === 'blank' ? null : selectedStrategy.value
    await api.post('/api/v1/community/wizard/draft', {
      step2: {
        strategy: strategyValue as CommunityWizardData['step2']['strategy'],
        enabledModules: [...enabledModules.value],
      },
    })
  } catch {
    // Silent fail
  } finally {
    saving.value = false
  }
}

// === Navigation ===
function goBack() {
  router.push('/b/community/wizard/step1')
}

async function goNext() {
  await saveDraft()
  router.push('/b/community/wizard/step3')
}

function goToCommunity() {
  router.push('/b/community')
}

// === Lifecycle ===
onMounted(() => {
  loadDraft()
})
</script>

<template>
  <div class="min-h-screen bg-[#0A0F1A]">
    <!-- Top Bar -->
    <div class="sticky top-0 z-20 bg-[#111B27] border-b border-[#1E293B] px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          class="flex items-center gap-1.5 text-sm text-slate-400 hover:text-white transition-colors"
          @click="goToCommunity"
        >
          <span class="material-symbols-rounded text-[20px]">arrow_back</span>
          Back
        </button>
        <div class="w-px h-5 bg-[#1E293B]"></div>
        <h1 class="text-lg font-semibold text-white">Create Community</h1>
      </div>
      <button
        class="px-4 py-2 text-sm font-medium text-slate-300 bg-[#1E293B] rounded-lg hover:bg-[#2D3B4E] transition-colors"
        :disabled="saving"
        @click="saveDraft"
      >
        {{ saving ? 'Saving...' : 'Save Draft' }}
      </button>
    </div>

    <!-- Stepper -->
    <div class="px-6 py-6 max-w-3xl mx-auto">
      <div class="flex items-center justify-between">
        <template v-for="(step, idx) in steps" :key="step.num">
          <div class="flex items-center gap-2">
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold transition-colors"
              :class="step.num === currentStep
                ? 'bg-[#48BB78] text-white'
                : step.num < currentStep
                  ? 'bg-[#48BB78]/20 text-[#48BB78]'
                  : 'bg-[#1E293B] text-slate-500'"
            >
              <span v-if="step.num < currentStep" class="material-symbols-rounded text-[18px]">check</span>
              <span v-else>{{ step.num }}</span>
            </div>
            <span
              class="text-sm font-medium"
              :class="step.num === currentStep ? 'text-white' : 'text-slate-500'"
            >
              {{ step.label }}
            </span>
          </div>
          <div
            v-if="idx < steps.length - 1"
            class="flex-1 h-px mx-3"
            :class="step.num < currentStep ? 'bg-[#48BB78]/40' : 'bg-[#1E293B]'"
          ></div>
        </template>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="animate-spin w-8 h-8 border-2 border-[#48BB78] border-t-transparent rounded-full"></div>
    </div>

    <!-- Main Content -->
    <div v-else class="px-6 pb-28 max-w-6xl mx-auto">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Left: Strategy + Modules (3 cols) -->
        <div class="lg:col-span-3 space-y-8">
          <!-- Strategy Selector -->
          <div>
            <h2 class="text-base font-semibold text-white mb-1">Choose a Strategy</h2>
            <p class="text-sm text-slate-400 mb-4">
              Select a recommended strategy to auto-configure modules, or start blank.
            </p>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <button
                v-for="strategy in strategies"
                :key="strategy.id"
                class="relative text-left p-4 rounded-xl border-2 transition-all group"
                :class="selectedStrategy === strategy.id
                  ? 'border-current bg-current/5'
                  : 'border-[#1E293B] bg-[#111B27] hover:border-[#2D3B4E]'"
                :style="selectedStrategy === strategy.id ? { borderColor: strategy.color, color: strategy.color } : {}"
                @click="selectStrategy(strategy)"
              >
                <!-- Selected indicator -->
                <div
                  v-if="selectedStrategy === strategy.id"
                  class="absolute top-3 right-3 w-5 h-5 rounded-full flex items-center justify-center"
                  :style="{ backgroundColor: strategy.color }"
                >
                  <span class="material-symbols-rounded text-[14px] text-white">check</span>
                </div>

                <div class="flex items-start gap-3">
                  <div
                    class="w-10 h-10 rounded-lg flex items-center justify-center shrink-0 transition-colors"
                    :style="{ backgroundColor: `${strategy.color}20`, color: strategy.color }"
                  >
                    <span class="material-symbols-rounded text-[22px]">{{ strategy.icon }}</span>
                  </div>
                  <div>
                    <p class="text-sm font-semibold text-white mb-0.5">{{ strategy.name }}</p>
                    <p class="text-xs text-slate-400 leading-relaxed">{{ strategy.description }}</p>
                    <div v-if="strategy.id !== 'blank'" class="flex flex-wrap gap-1 mt-2">
                      <span
                        v-for="modId in strategy.modules.filter(m => m !== 'sectors_tasks' && m !== 'points_level')"
                        :key="modId"
                        class="px-2 py-0.5 text-[10px] font-medium rounded-full"
                        :style="{ backgroundColor: `${strategy.color}15`, color: strategy.color }"
                      >
                        {{ allModules.find(m => m.id === modId)?.name }}
                      </span>
                    </div>
                  </div>
                </div>
              </button>
            </div>
          </div>

          <!-- Module Grid by System -->
          <div>
            <h2 class="text-base font-semibold text-white mb-1">Configure Modules</h2>
            <p class="text-sm text-slate-400 mb-5">
              Toggle individual modules on or off. Required modules cannot be disabled.
            </p>

            <div class="space-y-6">
              <div v-for="sys in systemOrder" :key="sys">
                <!-- System Header -->
                <div class="flex items-center gap-2 mb-3">
                  <span class="material-symbols-rounded text-[18px] text-slate-500">
                    {{ systemLabels[sys].icon }}
                  </span>
                  <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider">
                    {{ systemLabels[sys].label }}
                  </h3>
                </div>

                <!-- Module Cards -->
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                  <div
                    v-for="mod in modulesBySystem[sys]"
                    :key="mod.id"
                    class="bg-[#111B27] border rounded-xl p-4 transition-all"
                    :class="isModuleEnabled(mod.id)
                      ? 'border-[#48BB78]/30'
                      : 'border-[#1E293B]'"
                  >
                    <div class="flex items-start justify-between mb-2">
                      <div class="flex items-center gap-2">
                        <span
                          class="material-symbols-rounded text-[20px]"
                          :class="isModuleEnabled(mod.id) ? 'text-[#48BB78]' : 'text-slate-500'"
                        >
                          {{ mod.icon }}
                        </span>
                        <span class="text-sm font-semibold text-white">{{ mod.name }}</span>
                        <span
                          v-if="mod.required"
                          class="px-1.5 py-0.5 text-[10px] font-semibold uppercase tracking-wider rounded bg-[#48BB78]/10 text-[#48BB78]"
                        >
                          Core
                        </span>
                      </div>

                      <!-- Toggle -->
                      <button
                        class="relative w-10 h-5 rounded-full transition-colors shrink-0"
                        :class="isModuleEnabled(mod.id) ? 'bg-[#48BB78]' : 'bg-[#1E293B]'"
                        :disabled="mod.required"
                        @click="toggleModule(mod)"
                      >
                        <span
                          class="absolute top-0.5 w-4 h-4 rounded-full bg-white shadow transition-transform"
                          :class="isModuleEnabled(mod.id) ? 'translate-x-[22px]' : 'translate-x-0.5'"
                        ></span>
                      </button>
                    </div>
                    <p class="text-xs text-slate-400 leading-relaxed">{{ mod.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Summary Sidebar -->
        <div class="lg:col-span-1">
          <div class="sticky top-24">
            <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-5 space-y-5">
              <h3 class="text-sm font-semibold text-white">Configuration Summary</h3>

              <!-- Enabled count -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <span class="text-xs text-slate-400">Enabled Modules</span>
                  <span class="text-sm font-bold text-[#48BB78]">{{ enabledCount }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-xs text-slate-400">Optional Modules</span>
                  <span class="text-sm font-semibold text-white">{{ optionalEnabled }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-xs text-slate-400">Est. Setup Time</span>
                  <span class="text-sm font-semibold text-white">~{{ estimatedTime }} min</span>
                </div>
              </div>

              <!-- Divider -->
              <div class="border-t border-[#1E293B]"></div>

              <!-- Module List -->
              <div class="space-y-2">
                <p class="text-xs text-slate-500 uppercase tracking-wider font-semibold">Active Modules</p>
                <div
                  v-for="mod in allModules.filter(m => isModuleEnabled(m.id))"
                  :key="mod.id"
                  class="flex items-center gap-2"
                >
                  <span class="material-symbols-rounded text-[16px] text-[#48BB78]">check_circle</span>
                  <span class="text-xs text-slate-300">{{ mod.name }}</span>
                  <span
                    v-if="mod.required"
                    class="text-[9px] text-slate-500"
                  >(core)</span>
                </div>
              </div>

              <!-- Strategy badge -->
              <div v-if="selectedStrategy && selectedStrategy !== 'blank'" class="pt-2 border-t border-[#1E293B]">
                <div class="flex items-center gap-2">
                  <span class="material-symbols-rounded text-[14px] text-slate-500">auto_awesome</span>
                  <span class="text-xs text-slate-400">
                    Strategy: <span class="text-white font-medium">{{ strategies.find(s => s.id === selectedStrategy)?.name }}</span>
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Action Bar -->
    <div class="fixed bottom-0 left-0 right-0 z-20 bg-[#111B27] border-t border-[#1E293B] px-6 py-4">
      <div class="max-w-6xl mx-auto flex items-center justify-between">
        <button
          class="px-5 py-2.5 text-sm font-medium text-slate-300 bg-[#1E293B] rounded-lg hover:bg-[#2D3B4E] transition-colors flex items-center gap-2"
          @click="goBack"
        >
          <span class="material-symbols-rounded text-[18px]">arrow_back</span>
          Back
        </button>
        <button
          class="px-6 py-2.5 text-sm font-semibold bg-[#48BB78] text-white rounded-lg hover:bg-[#38A169] transition-colors flex items-center gap-2"
          @click="goNext"
        >
          Next: Quick Setup
          <span class="material-symbols-rounded text-[18px]">arrow_forward</span>
        </button>
      </div>
    </div>
  </div>
</template>
