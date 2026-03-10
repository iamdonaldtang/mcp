<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()

// === Types ===
type DeploymentPath = 'domain' | 'iframe' | 'widget' | 'pagebuilder' | 'sdk'
type PreviewMode = 'desktop' | 'mobile'

interface ChecklistItem {
  id: string
  label: string
  description: string
  passed: boolean
  loading: boolean
}

// === State ===
const loading = ref(true)
const publishing = ref(false)
const previewMode = ref<PreviewMode>('desktop')
const path = ref<DeploymentPath>('domain')

// === Draft data for summary ===
const summary = reactive({
  path: '' as string,
  pathLabel: '' as string,
  domain: '',
  brandColors: {
    primary: '#9B7EE0',
    secondary: '#5D7EF1',
    accent: '#F59E0B',
    background: '#0A0F1A',
  },
  modules: [] as string[],
  headingFont: 'Inter',
  bodyFont: 'Inter',
  buttonStyle: 'filled',
  logoUrl: '',
})

// === Stepper ===
const steps = [
  { num: 1, label: 'Path' },
  { num: 2, label: 'Configure' },
  { num: 3, label: 'Brand' },
  { num: 4, label: 'Preview & Publish' },
]
const currentStep = 4

// === Readiness Checklist ===
const checklist = ref<ChecklistItem[]>([
  { id: 'community', label: 'Community Active', description: 'An active community is required for White Label.', passed: false, loading: true },
  { id: 'path', label: 'Path Configured', description: 'Deployment path must be fully set up.', passed: false, loading: true },
  { id: 'brand', label: 'Brand Settings Saved', description: 'Brand colors and typography are configured.', passed: false, loading: true },
  { id: 'subscription', label: 'Subscription Active', description: 'A valid White Label subscription is required.', passed: false, loading: true },
  { id: 'twitter', label: 'Twitter Connected', description: 'Connect your Twitter account for distribution.', passed: false, loading: true },
])

const allPassed = computed(() => checklist.value.every(item => item.passed))
const passedCount = computed(() => checklist.value.filter(item => item.passed).length)

// === Path labels ===
const pathLabels: Record<string, string> = {
  domain: 'Custom Domain',
  iframe: 'Iframe Embed',
  widget: 'Widget Library',
  pagebuilder: 'Page Builder',
  sdk: 'SDK Integration',
}

// === Load all data ===
async function loadData() {
  loading.value = true
  try {
    // Load draft
    const draftRes = await api.get('/api/v1/whitelabel/wizard/draft')
    const draft = draftRes.data?.data
    if (draft) {
      const resolvedPath = draft.embedSubPath || draft.path || 'domain'
      path.value = resolvedPath
      summary.path = resolvedPath
      summary.pathLabel = pathLabels[resolvedPath] || resolvedPath
      summary.domain = draft.step2?.domain?.domain || ''

      // Modules from widget config
      if (draft.step2?.widget?.modules) {
        summary.modules = Object.entries(draft.step2.widget.modules)
          .filter(([, v]) => v)
          .map(([k]) => k)
      }
    }

    // Load brand
    const brandRes = await api.get('/api/v1/whitelabel/brand')
    const brand = brandRes.data?.data
    if (brand) {
      summary.brandColors.primary = brand.primaryColor || '#9B7EE0'
      summary.brandColors.secondary = brand.secondaryColor || '#5D7EF1'
      summary.brandColors.accent = brand.accentColor || '#F59E0B'
      summary.brandColors.background = brand.backgroundColor || '#0A0F1A'
      summary.headingFont = brand.headingFont || 'Inter'
      summary.bodyFont = brand.bodyFont || 'Inter'
      summary.buttonStyle = brand.buttonStyle || 'filled'
      summary.logoUrl = brand.logoUrl || ''
    }

    // Load readiness
    const readinessRes = await api.get('/api/v1/whitelabel/readiness')
    const readiness = readinessRes.data?.data
    if (readiness) {
      for (const item of checklist.value) {
        item.passed = readiness[item.id] ?? false
        item.loading = false
      }
    } else {
      checklist.value.forEach(item => { item.loading = false })
    }
  } catch {
    checklist.value.forEach(item => { item.loading = false })
  } finally {
    loading.value = false
  }
}

// === Publish ===
async function publish() {
  if (!allPassed.value || publishing.value) return
  publishing.value = true
  try {
    await api.post('/api/v1/whitelabel/publish')
    router.push('/b/whitelabel')
  } catch {
    // Error publishing
  } finally {
    publishing.value = false
  }
}

// === Navigation ===
function goBack() {
  router.push('/b/whitelabel/wizard/step-3')
}

// === Lifecycle ===
onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="min-h-screen bg-page-bg">
    <!-- Top Bar -->
    <div class="sticky top-0 z-20 bg-card-bg border-b border-border px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          class="flex items-center gap-1.5 text-sm text-slate-400 hover:text-white transition-colors"
          @click="goBack"
        >
          <span class="material-symbols-rounded text-[20px]">arrow_back</span>
          Back
        </button>
        <div class="w-px h-5 bg-border"></div>
        <h1 class="text-lg font-semibold text-white">Set Up White Label</h1>
      </div>
    </div>

    <!-- Stepper -->
    <div class="px-6 py-6 max-w-3xl mx-auto">
      <div class="flex items-center justify-between">
        <template v-for="(step, idx) in steps" :key="step.num">
          <div class="flex items-center gap-2">
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold transition-colors"
              :class="step.num === currentStep
                ? 'bg-whitelabel text-white'
                : step.num < currentStep
                  ? 'bg-whitelabel/20 text-whitelabel'
                  : 'bg-border text-slate-500'"
            >
              <span v-if="step.num < currentStep" class="material-symbols-rounded text-[18px]">check</span>
              <span v-else>{{ step.num }}</span>
            </div>
            <span
              class="text-sm font-medium"
              :class="step.num <= currentStep ? 'text-white' : 'text-slate-500'"
            >
              {{ step.label }}
            </span>
          </div>
          <div
            v-if="idx < steps.length - 1"
            class="flex-1 h-px mx-3"
            :class="step.num < currentStep ? 'bg-whitelabel/40' : 'bg-border'"
          ></div>
        </template>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="animate-spin w-8 h-8 border-2 border-whitelabel border-t-transparent rounded-full"></div>
    </div>

    <!-- Main Content -->
    <div v-else class="px-6 pb-28 max-w-6xl mx-auto">
      <div class="grid grid-cols-1 lg:grid-cols-5 gap-8">

        <!-- Left: Deployment Preview -->
        <div class="lg:col-span-3">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-white">Deployment Preview</h2>
            <!-- Desktop/Mobile Toggle -->
            <div class="flex items-center bg-card-bg border border-border rounded-lg p-0.5">
              <button
                class="px-3 py-1.5 text-xs font-medium rounded-md transition-colors flex items-center gap-1.5"
                :class="previewMode === 'desktop' ? 'bg-whitelabel text-white' : 'text-slate-400 hover:text-slate-300'"
                @click="previewMode = 'desktop'"
              >
                <span class="material-symbols-rounded text-[16px]">desktop_windows</span>
                Desktop
              </button>
              <button
                class="px-3 py-1.5 text-xs font-medium rounded-md transition-colors flex items-center gap-1.5"
                :class="previewMode === 'mobile' ? 'bg-whitelabel text-white' : 'text-slate-400 hover:text-slate-300'"
                @click="previewMode = 'mobile'"
              >
                <span class="material-symbols-rounded text-[16px]">smartphone</span>
                Mobile
              </button>
            </div>
          </div>

          <!-- Preview Frame -->
          <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
            <!-- Browser chrome -->
            <div class="bg-[#0D1420] px-4 py-2.5 flex items-center gap-3 border-b border-border">
              <div class="flex items-center gap-1.5">
                <div class="w-3 h-3 rounded-full bg-[#FF5F57]"></div>
                <div class="w-3 h-3 rounded-full bg-[#FEBC2E]"></div>
                <div class="w-3 h-3 rounded-full bg-[#28C840]"></div>
              </div>
              <div class="flex-1 px-3 py-1 text-xs text-slate-500 bg-page-bg rounded-md font-mono truncate">
                <template v-if="path === 'domain' && summary.domain">
                  https://{{ summary.domain }}
                </template>
                <template v-else-if="path === 'iframe'">
                  https://yourapp.com (iframe embed)
                </template>
                <template v-else-if="path === 'widget'">
                  https://yourapp.com (widget components)
                </template>
                <template v-else-if="path === 'pagebuilder'">
                  https://yourapp.com/community
                </template>
                <template v-else>
                  https://yourapp.com (SDK integration)
                </template>
              </div>
            </div>

            <!-- Preview content area -->
            <div
              class="flex justify-center p-6"
              :class="previewMode === 'mobile' ? 'py-8' : ''"
            >
              <div
                class="w-full rounded-lg overflow-hidden border transition-all"
                :class="previewMode === 'mobile' ? 'max-w-93.75' : ''"
                :style="{ borderColor: summary.brandColors.primary + '30' }"
              >
                <!-- Mock community header -->
                <div
                  class="px-5 py-4 flex items-center justify-between"
                  :style="{ backgroundColor: summary.brandColors.primary + '15' }"
                >
                  <div class="flex items-center gap-2.5">
                    <div
                      v-if="summary.logoUrl"
                      class="w-7 h-7 rounded-lg overflow-hidden"
                    >
                      <img :src="summary.logoUrl" alt="Logo" class="w-full h-full object-contain" />
                    </div>
                    <div
                      v-else
                      class="w-7 h-7 rounded-lg flex items-center justify-center text-white text-xs font-bold"
                      :style="{ backgroundColor: summary.brandColors.primary }"
                    >W</div>
                    <span
                      class="text-sm font-semibold text-slate-200"
                      :style="{ fontFamily: summary.headingFont }"
                    >Your Community</span>
                  </div>
                  <div class="flex items-center gap-3 text-xs text-slate-500">
                    <span>Quests</span>
                    <span>Leaderboard</span>
                    <span>Shop</span>
                  </div>
                </div>

                <!-- Mock content -->
                <div class="p-5 space-y-4" :style="{ backgroundColor: summary.brandColors.background }">
                  <div class="grid grid-cols-3 gap-2">
                    <div
                      v-for="n in 3"
                      :key="n"
                      class="h-16 rounded-lg"
                      :style="{ backgroundColor: summary.brandColors.primary + '10' }"
                    ></div>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="n in 2"
                      :key="n"
                      class="h-12 rounded-lg border"
                      :style="{ borderColor: summary.brandColors.primary + '20', backgroundColor: summary.brandColors.primary + '05' }"
                    ></div>
                  </div>
                  <button
                    class="w-full py-2.5 text-sm font-semibold text-white"
                    :style="{
                      backgroundColor: summary.brandColors.primary,
                      borderRadius: summary.buttonStyle === 'rounded' ? '9999px' : '8px',
                    }"
                  >
                    Get Started
                  </button>
                </div>

                <!-- Footer -->
                <div class="px-5 py-3 text-center border-t" :style="{ borderColor: summary.brandColors.primary + '15' }">
                  <span class="text-[10px] text-slate-600">Powered by TaskOn</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Summary & Checklist -->
        <div class="lg:col-span-2 space-y-6">

          <!-- Deployment Summary -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
            <h3 class="text-base font-semibold text-white">Deployment Summary</h3>

            <div class="space-y-3">
              <!-- Path -->
              <div class="flex items-center justify-between">
                <span class="text-sm text-slate-400">Path</span>
                <span class="text-sm font-medium text-white flex items-center gap-1.5">
                  <span class="material-symbols-rounded text-[16px] text-whitelabel">
                    {{ path === 'domain' ? 'language' : path === 'iframe' ? 'web' : path === 'widget' ? 'widgets' : path === 'pagebuilder' ? 'dashboard' : 'code' }}
                  </span>
                  {{ summary.pathLabel }}
                </span>
              </div>

              <!-- Domain (if applicable) -->
              <div v-if="path === 'domain' && summary.domain" class="flex items-center justify-between">
                <span class="text-sm text-slate-400">Domain</span>
                <span class="text-sm font-medium text-white font-mono">{{ summary.domain }}</span>
              </div>

              <!-- Brand Colors -->
              <div class="flex items-center justify-between">
                <span class="text-sm text-slate-400">Brand Colors</span>
                <div class="flex items-center gap-1.5">
                  <div class="w-5 h-5 rounded-full border border-border" :style="{ backgroundColor: summary.brandColors.primary }"></div>
                  <div class="w-5 h-5 rounded-full border border-border" :style="{ backgroundColor: summary.brandColors.secondary }"></div>
                  <div class="w-5 h-5 rounded-full border border-border" :style="{ backgroundColor: summary.brandColors.accent }"></div>
                </div>
              </div>

              <!-- Typography -->
              <div class="flex items-center justify-between">
                <span class="text-sm text-slate-400">Typography</span>
                <span class="text-sm font-medium text-white">{{ summary.headingFont }}</span>
              </div>

              <!-- Modules (for widget path) -->
              <div v-if="summary.modules.length > 0" class="flex items-start justify-between">
                <span class="text-sm text-slate-400">Modules</span>
                <div class="flex flex-wrap gap-1 justify-end max-w-40">
                  <span
                    v-for="mod in summary.modules"
                    :key="mod"
                    class="px-2 py-0.5 text-[10px] font-medium bg-whitelabel/10 text-whitelabel rounded-full capitalize"
                  >
                    {{ mod }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Readiness Checklist -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold text-white">Readiness Checklist</h3>
              <span class="text-xs font-medium" :class="allPassed ? 'text-status-active' : 'text-slate-500'">
                {{ passedCount }}/{{ checklist.length }}
              </span>
            </div>

            <!-- Progress bar -->
            <div class="w-full h-1.5 bg-border rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500"
                :class="allPassed ? 'bg-status-active' : 'bg-whitelabel'"
                :style="{ width: `${(passedCount / checklist.length) * 100}%` }"
              ></div>
            </div>

            <div class="space-y-2">
              <div
                v-for="item in checklist"
                :key="item.id"
                class="flex items-start gap-3 p-3 rounded-lg border transition-colors"
                :class="item.passed ? 'border-status-active-bg bg-status-active-bg/50' : 'border-border'"
              >
                <!-- Status icon -->
                <div class="shrink-0 mt-0.5">
                  <div v-if="item.loading" class="animate-spin w-5 h-5 border-2 border-slate-500 border-t-transparent rounded-full"></div>
                  <span
                    v-else
                    class="material-symbols-rounded text-[20px]"
                    :class="item.passed ? 'text-status-active' : 'text-slate-500'"
                  >
                    {{ item.passed ? 'check_circle' : 'radio_button_unchecked' }}
                  </span>
                </div>
                <!-- Content -->
                <div>
                  <p class="text-sm font-medium" :class="item.passed ? 'text-status-active' : 'text-slate-300'">
                    {{ item.label }}
                  </p>
                  <p class="text-xs text-slate-500 mt-0.5">{{ item.description }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Publish Button -->
          <button
            class="w-full py-3 text-sm font-semibold rounded-xl transition-all flex items-center justify-center gap-2"
            :class="allPassed
              ? 'bg-status-active text-white hover:bg-[#15803D]'
              : 'bg-border text-slate-500 cursor-not-allowed'"
            :disabled="!allPassed || publishing"
            @click="publish"
          >
            <span v-if="publishing" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
            <span v-else class="material-symbols-rounded text-[18px]">rocket_launch</span>
            {{ publishing ? 'Publishing...' : 'Publish White Label' }}
          </button>
          <p v-if="!allPassed" class="text-center text-xs text-slate-500 -mt-3">
            Complete all checklist items to publish.
          </p>
        </div>
      </div>
    </div>

    <!-- Bottom Action Bar -->
    <div class="fixed bottom-0 left-0 right-0 z-20 bg-card-bg border-t border-border px-6 py-4">
      <div class="max-w-6xl mx-auto flex items-center justify-between">
        <button
          class="px-5 py-2.5 text-sm font-medium text-slate-300 bg-border rounded-lg hover:bg-[#2D3B4E] transition-colors flex items-center gap-2"
          @click="goBack"
        >
          <span class="material-symbols-rounded text-[18px]">arrow_back</span>
          Back
        </button>
        <button
          class="px-6 py-2.5 text-sm font-semibold rounded-lg transition-all flex items-center gap-2"
          :class="allPassed
            ? 'bg-status-active text-white hover:bg-[#15803D]'
            : 'bg-border text-slate-500 cursor-not-allowed'"
          :disabled="!allPassed || publishing"
          @click="publish"
        >
          <span v-if="publishing" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
          <span v-else class="material-symbols-rounded text-[18px]">rocket_launch</span>
          {{ publishing ? 'Publishing...' : 'Publish' }}
        </button>
      </div>
    </div>
  </div>
</template>
