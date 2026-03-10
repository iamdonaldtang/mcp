<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()

// === Types ===
type DeploymentPath = 'domain' | 'iframe' | 'widget' | 'pagebuilder' | 'sdk' | ''

interface PathOption {
  id: DeploymentPath
  icon: string
  title: string
  description: string
  badge?: string
  recommended?: boolean
}

// === Wizard State ===
const selectedPath = ref<DeploymentPath>('')
const saving = ref(false)
const loadingDraft = ref(true)
const lastSavedAt = ref<string | null>(null)
let autoSaveTimer: ReturnType<typeof setInterval> | null = null

// === Stepper ===
const steps = [
  { num: 1, label: 'Path' },
  { num: 2, label: 'Configure' },
  { num: 3, label: 'Brand' },
  { num: 4, label: 'Preview & Publish' },
]
const currentStep = 1

// === Path Options ===
const pathOptions: PathOption[] = [
  {
    id: 'domain',
    icon: 'language',
    title: 'Custom Domain',
    description: 'Host your community at yourproject.com. Zero code required.',
    badge: 'No Code',
  },
  {
    id: 'iframe',
    icon: 'web',
    title: 'Iframe Embed',
    description: 'Embed your community in an iframe. Quick setup with SSO support.',
  },
  {
    id: 'widget',
    icon: 'widgets',
    title: 'Widget Library',
    description: 'Drop-in React/Vue components. Customizable, lightweight.',
    recommended: true,
  },
  {
    id: 'pagebuilder',
    icon: 'dashboard',
    title: 'Page Builder',
    description: 'Build custom pages with drag-and-drop. No coding needed.',
    badge: 'Visual Editor',
  },
  {
    id: 'sdk',
    icon: 'code',
    title: 'SDK Integration',
    description: 'Full API access. Build completely custom experiences.',
    badge: 'Full Control',
  },
]

// === Validation ===
const isValid = computed(() => selectedPath.value !== '')

// === Embed paths that need step 1.5 ===
const embedPaths: DeploymentPath[] = ['iframe', 'widget', 'pagebuilder']

// === Draft Persistence ===
async function loadDraft() {
  loadingDraft.value = true
  try {
    const res = await api.get('/api/v1/whitelabel/wizard/draft')
    const draft = res.data?.data
    if (draft?.path) {
      selectedPath.value = draft.path
    }
  } catch {
    // No draft yet
  } finally {
    loadingDraft.value = false
  }
}

async function saveDraft(silent = true) {
  if (saving.value) return
  saving.value = true
  try {
    await api.post('/api/v1/whitelabel/wizard/draft', { path: selectedPath.value })
    lastSavedAt.value = new Date().toLocaleTimeString()
  } catch {
    if (!silent) {
      // Could show toast
    }
  } finally {
    saving.value = false
  }
}

// === Navigation ===
function goBack() {
  router.push('/b/whitelabel')
}

async function goNext() {
  if (!isValid.value) return
  await saveDraft(false)

  // Embed-related paths go to step 1.5 (embed options)
  if (embedPaths.includes(selectedPath.value)) {
    router.push('/b/whitelabel/wizard/step-1-5')
  } else {
    router.push('/b/whitelabel/wizard/step-2')
  }
}

// === Lifecycle ===
onMounted(async () => {
  await loadDraft()
  autoSaveTimer = setInterval(() => {
    if (selectedPath.value) saveDraft()
  }, 30000)
})

onUnmounted(() => {
  if (autoSaveTimer) clearInterval(autoSaveTimer)
})

watch(selectedPath, () => {
  // Triggers reactivity for template
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
      <div class="flex items-center gap-3">
        <span v-if="lastSavedAt" class="text-xs text-slate-500">
          Saved {{ lastSavedAt }}
        </span>
        <button
          class="px-4 py-2 text-sm font-medium text-slate-300 bg-border rounded-lg hover:bg-[#2D3B4E] transition-colors"
          :disabled="saving"
          @click="saveDraft(false)"
        >
          {{ saving ? 'Saving...' : 'Save Draft' }}
        </button>
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
              :class="step.num === currentStep ? 'text-white' : 'text-slate-500'"
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

    <!-- Loading State -->
    <div v-if="loadingDraft" class="flex items-center justify-center py-24">
      <div class="animate-spin w-8 h-8 border-2 border-whitelabel border-t-transparent rounded-full"></div>
    </div>

    <!-- Main Content -->
    <div v-else class="px-6 pb-28 max-w-4xl mx-auto">
      <div class="mb-8">
        <h2 class="text-xl font-semibold text-white mb-2">Choose Your Deployment Path</h2>
        <p class="text-sm text-slate-400">
          Select how you want to integrate your White Label community. You can change this later.
        </p>
      </div>

      <!-- Path Cards -->
      <div class="grid grid-cols-1 gap-4">
        <button
          v-for="option in pathOptions"
          :key="option.id"
          class="relative text-left bg-card-bg border rounded-xl p-6 transition-all hover:bg-[#151F2E] group"
          :class="selectedPath === option.id
            ? 'border-whitelabel ring-1 ring-whitelabel/30'
            : 'border-border hover:border-[#2D3B4E]'"
          @click="selectedPath = option.id"
        >
          <div class="flex items-start gap-4">
            <!-- Icon -->
            <div
              class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0 transition-colors"
              :class="selectedPath === option.id
                ? 'bg-whitelabel/15 text-whitelabel'
                : 'bg-border text-slate-400 group-hover:text-slate-300'"
            >
              <span class="material-symbols-rounded text-[24px]">{{ option.icon }}</span>
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <h3 class="text-base font-semibold text-white">{{ option.title }}</h3>
                <span
                  v-if="option.recommended"
                  class="px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider bg-whitelabel/15 text-whitelabel rounded-full"
                >
                  Recommended
                </span>
                <span
                  v-if="option.badge"
                  class="px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider bg-border text-slate-300 rounded-full"
                >
                  {{ option.badge }}
                </span>
              </div>
              <p class="text-sm text-slate-400 leading-relaxed">{{ option.description }}</p>
            </div>

            <!-- Selection indicator -->
            <div
              class="w-6 h-6 rounded-full border-2 flex items-center justify-center shrink-0 transition-colors mt-1"
              :class="selectedPath === option.id
                ? 'border-whitelabel bg-whitelabel'
                : 'border-[#334155]'"
            >
              <span
                v-if="selectedPath === option.id"
                class="material-symbols-rounded text-[16px] text-white"
              >check</span>
            </div>
          </div>
        </button>
      </div>
    </div>

    <!-- Bottom Action Bar -->
    <div class="fixed bottom-0 left-0 right-0 z-20 bg-card-bg border-t border-border px-6 py-4">
      <div class="max-w-4xl mx-auto flex items-center justify-end">
        <button
          class="px-6 py-2.5 text-sm font-semibold rounded-lg transition-all flex items-center gap-2"
          :class="isValid
            ? 'bg-whitelabel text-white hover:bg-[#8B6ED0]'
            : 'bg-border text-slate-500 cursor-not-allowed'"
          :disabled="!isValid"
          @click="goNext"
        >
          Next: Configure
          <span class="material-symbols-rounded text-[18px]">arrow_forward</span>
        </button>
      </div>
    </div>
  </div>
</template>
