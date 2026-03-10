<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

interface Strategy {
  id: string
  title: string
  description: string
  icon: string
  iconColor: string
  iconBg: string
  modules: string[]
  metric: string
}

const strategies: Strategy[] = [
  {
    id: 'activate',
    title: 'Activate New Users',
    description: 'Onboard and activate users with structured tasks, instant feedback, and clear progression paths.',
    icon: 'rocket_launch',
    iconColor: '#48BB78',
    iconBg: '#0A1F1A',
    modules: ['Sectors & Tasks', 'Points & Level', 'TaskChain'],
    metric: 'Activation Rate',
  },
  {
    id: 'engagement',
    title: 'Drive Daily Engagement',
    description: 'Build daily habits with streaks, leaderboards, and social competition that create fear of missing out.',
    icon: 'local_fire_department',
    iconColor: '#ED8936',
    iconBg: '#1F1508',
    modules: ['Sectors & Tasks', 'Points & Level', 'DayChain', 'Leaderboard'],
    metric: 'DAU',
  },
  {
    id: 'retention',
    title: 'Maximize Retention',
    description: 'Create leaving costs through sunk cost accumulation, variable rewards, and exclusive milestone benefits.',
    icon: 'shield',
    iconColor: '#9B7EE0',
    iconBg: '#1A1033',
    modules: ['Sectors & Tasks', 'Points & Level', 'DayChain', 'Milestones', 'Benefits Shop'],
    metric: '30-Day Retention',
  },
]

const selectedStrategy = ref<string>('activate')

const engineSteps = [
  { label: 'Quest', desc: 'Acquire users', icon: 'campaign' },
  { label: 'Activate', desc: 'Complete onboarding', icon: 'rocket_launch' },
  { label: 'Engage', desc: 'Drive daily habits', icon: 'local_fire_department' },
  { label: 'Retain', desc: 'Create leaving costs', icon: 'shield' },
]

function createWithStrategy() {
  router.push(`/b/community/wizard/step-1?template=${selectedStrategy.value}`)
}

function createFromScratch() {
  router.push('/b/community/wizard/step-1?template=blank')
}
</script>

<template>
  <div class="max-w-5xl mx-auto space-y-10">
    <!-- Welcome Section -->
    <div>
      <div class="w-16 h-16 rounded-2xl flex items-center justify-center mb-4" style="background: #0A1F1A">
        <span class="material-symbols-rounded text-3xl text-community">groups</span>
      </div>
      <h1 class="text-2xl font-bold text-text-primary mb-2">Welcome to Community</h1>
      <p class="text-base text-text-secondary max-w-2xl">
        Build a thriving community that turns one-time visitors into loyal, engaged members. Choose a retention strategy to get started.
      </p>
    </div>

    <!-- Retention Strategies -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Retention Strategies</div>
      <div class="grid grid-cols-3 gap-6">
        <button
          v-for="s in strategies"
          :key="s.id"
          class="text-left bg-card-bg border rounded-xl p-5 transition-all duration-200 hover:border-opacity-80"
          :class="selectedStrategy === s.id
            ? 'border-community border-2 shadow-[0_0_0_1px_#48BB78]'
            : 'border-border hover:bg-white/2'"
          @click="selectedStrategy = s.id"
        >
          <div class="w-10 h-10 rounded-lg flex items-center justify-center mb-3" :style="{ background: s.iconBg }">
            <span class="material-symbols-rounded text-xl" :style="{ color: s.iconColor }">{{ s.icon }}</span>
          </div>
          <h3 class="text-base font-semibold text-text-primary mb-1">{{ s.title }}</h3>
          <p class="text-sm text-text-secondary mb-3 leading-relaxed">{{ s.description }}</p>

          <!-- Expanded details when selected -->
          <div v-if="selectedStrategy === s.id" class="mt-3 pt-3 border-t border-border">
            <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2">Includes</div>
            <div class="flex flex-wrap gap-1.5">
              <span
                v-for="mod in s.modules"
                :key="mod"
                class="px-2 py-0.5 text-xs rounded-full"
                :style="{ background: s.iconBg, color: s.iconColor }"
              >{{ mod }}</span>
            </div>
            <div class="mt-2 text-xs text-text-muted">Key metric: <span class="text-text-secondary">{{ s.metric }}</span></div>
          </div>
        </button>
      </div>
    </div>

    <!-- How It Works -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">How It Works</div>
      <div class="bg-card-bg border border-border rounded-xl p-6">
        <div class="flex items-center justify-between">
          <div
            v-for="(step, i) in engineSteps"
            :key="step.label"
            class="flex items-center"
          >
            <div class="text-center">
              <div class="w-12 h-12 rounded-xl bg-page-bg border border-border flex items-center justify-center mb-2 mx-auto">
                <span class="material-symbols-rounded text-text-secondary">{{ step.icon }}</span>
              </div>
              <div class="text-sm font-medium text-text-primary">{{ step.label }}</div>
              <div class="text-xs text-text-muted">{{ step.desc }}</div>
            </div>
            <span v-if="i < engineSteps.length - 1" class="material-symbols-rounded text-text-muted mx-6">arrow_forward</span>
          </div>
        </div>
      </div>
    </div>

    <!-- CTA -->
    <div class="flex items-center gap-4">
      <button
        class="px-6 py-3 bg-community text-white font-semibold rounded-xl hover:bg-community/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="!selectedStrategy"
        @click="createWithStrategy"
      >
        Create Community with This Strategy
      </button>
      <button
        class="text-sm text-text-muted hover:text-text-secondary transition-colors"
        @click="createFromScratch"
      >
        Or start from scratch →
      </button>
    </div>

    <!-- Divider -->
    <div class="border-t border-border"></div>

    <!-- Resources -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Resources</div>
      <div class="grid grid-cols-3 gap-6">
        <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-2xl mb-3 block">play_circle</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Video Tutorial</h4>
          <p class="text-xs text-text-secondary">Watch a 5-minute walkthrough of Community setup</p>
        </a>
        <a href="#" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-2xl mb-3 block">menu_book</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Retention Playbook</h4>
          <p class="text-xs text-text-secondary">Learn proven Web3 retention strategies</p>
        </a>
        <a href="#" target="_blank" class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block">
          <span class="material-symbols-rounded text-community text-2xl mb-3 block">open_in_new</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">Learn More</h4>
          <p class="text-xs text-text-secondary">Explore Community product page</p>
        </a>
      </div>
    </div>
  </div>
</template>
