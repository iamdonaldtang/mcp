<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()

// === Types ===
type EmbedSubPath = 'iframe' | 'widget' | 'pagebuilder' | ''

interface EmbedOption {
  id: EmbedSubPath
  icon: string
  title: string
  description: string
  features: string[]
}

// === State ===
const selectedSubPath = ref<EmbedSubPath>('')
const parentPath = ref('')
const loading = ref(true)
const saving = ref(false)

// === Stepper ===
const steps = [
  { num: 1, label: 'Path' },
  { num: 2, label: 'Configure' },
  { num: 3, label: 'Brand' },
  { num: 4, label: 'Preview & Publish' },
]

// === Embed Options ===
const embedOptions: EmbedOption[] = [
  {
    id: 'iframe',
    icon: 'web',
    title: 'Iframe Embed',
    description: 'Embed the entire community experience in an iframe. Simplest integration with optional SSO.',
    features: ['Quick 5-minute setup', 'SSO support', 'Responsive sizing', 'Minimal code changes'],
  },
  {
    id: 'widget',
    icon: 'widgets',
    title: 'Widget Library',
    description: 'Drop individual community modules into your app as React or Vue components.',
    features: ['Granular control', 'Tree-shakeable', 'Theme-aware', 'TypeScript support'],
  },
  {
    id: 'pagebuilder',
    icon: 'dashboard',
    title: 'Page Builder',
    description: 'Visually compose custom pages using community modules with drag-and-drop.',
    features: ['No coding needed', 'Visual editor', 'Custom layouts', 'Instant preview'],
  },
]

// === Validation ===
const isValid = computed(() => selectedSubPath.value !== '')

// === Load draft to check parent path ===
onMounted(async () => {
  loading.value = true
  try {
    const res = await api.get('/api/v1/whitelabel/wizard/draft')
    const draft = res.data?.data
    parentPath.value = draft?.path || ''

    // If parent path is not an embed path, redirect to step 2 directly
    const embedPaths = ['iframe', 'widget', 'pagebuilder']
    if (!embedPaths.includes(parentPath.value)) {
      router.replace('/b/whitelabel/wizard/step-2')
      return
    }

    // Pre-select the embed sub-path from draft if available
    if (draft?.embedSubPath) {
      selectedSubPath.value = draft.embedSubPath
    } else {
      selectedSubPath.value = parentPath.value as EmbedSubPath
    }
  } catch {
    router.replace('/b/whitelabel/wizard/step-1')
  } finally {
    loading.value = false
  }
})

// === Navigation ===
function goBack() {
  router.push('/b/whitelabel/wizard/step-1')
}

async function goNext() {
  if (!isValid.value) return
  saving.value = true
  try {
    await api.post('/api/v1/whitelabel/wizard/draft', {
      path: selectedSubPath.value,
      embedSubPath: selectedSubPath.value,
    })
    router.push('/b/whitelabel/wizard/step-2')
  } catch {
    // Error saving
  } finally {
    saving.value = false
  }
}
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
              :class="step.num === 1
                ? 'bg-whitelabel/20 text-whitelabel'
                : 'bg-border text-slate-500'"
            >
              <span v-if="step.num === 1" class="material-symbols-rounded text-[18px]">check</span>
              <span v-else>{{ step.num }}</span>
            </div>
            <span
              class="text-sm font-medium"
              :class="step.num === 1 ? 'text-whitelabel' : 'text-slate-500'"
            >
              {{ step.label }}
            </span>
          </div>
          <div
            v-if="idx < steps.length - 1"
            class="flex-1 h-px mx-3 bg-border"
          ></div>
        </template>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="animate-spin w-8 h-8 border-2 border-whitelabel border-t-transparent rounded-full"></div>
    </div>

    <!-- Main Content -->
    <div v-else class="px-6 pb-28 max-w-4xl mx-auto">
      <div class="mb-8">
        <div class="flex items-center gap-2 mb-2">
          <span class="px-2.5 py-1 text-xs font-medium bg-whitelabel/15 text-whitelabel rounded-lg">
            Step 1.5
          </span>
          <h2 class="text-xl font-semibold text-white">Choose Embed Method</h2>
        </div>
        <p class="text-sm text-slate-400">
          You selected an embed-based path. Choose the specific integration method that works best for your app.
        </p>
      </div>

      <!-- Embed Option Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <button
          v-for="option in embedOptions"
          :key="option.id"
          class="relative text-left bg-card-bg border rounded-xl p-6 transition-all hover:bg-[#151F2E] group flex flex-col"
          :class="selectedSubPath === option.id
            ? 'border-whitelabel ring-1 ring-whitelabel/30'
            : 'border-border hover:border-[#2D3B4E]'"
          @click="selectedSubPath = option.id"
        >
          <!-- Selection dot -->
          <div
            class="absolute top-4 right-4 w-5 h-5 rounded-full border-2 flex items-center justify-center transition-colors"
            :class="selectedSubPath === option.id
              ? 'border-whitelabel bg-whitelabel'
              : 'border-[#334155]'"
          >
            <span
              v-if="selectedSubPath === option.id"
              class="material-symbols-rounded text-[14px] text-white"
            >check</span>
          </div>

          <!-- Icon -->
          <div
            class="w-11 h-11 rounded-xl flex items-center justify-center mb-4 transition-colors"
            :class="selectedSubPath === option.id
              ? 'bg-whitelabel/15 text-whitelabel'
              : 'bg-border text-slate-400'"
          >
            <span class="material-symbols-rounded text-[22px]">{{ option.icon }}</span>
          </div>

          <!-- Title & Description -->
          <h3 class="text-base font-semibold text-white mb-1.5">{{ option.title }}</h3>
          <p class="text-sm text-slate-400 leading-relaxed mb-4">{{ option.description }}</p>

          <!-- Features -->
          <ul class="mt-auto space-y-1.5">
            <li
              v-for="feature in option.features"
              :key="feature"
              class="flex items-center gap-2 text-xs text-slate-500"
            >
              <span class="material-symbols-rounded text-[14px] text-whitelabel/60">check_circle</span>
              {{ feature }}
            </li>
          </ul>
        </button>
      </div>
    </div>

    <!-- Bottom Action Bar -->
    <div class="fixed bottom-0 left-0 right-0 z-20 bg-card-bg border-t border-border px-6 py-4">
      <div class="max-w-4xl mx-auto flex items-center justify-between">
        <button
          class="px-5 py-2.5 text-sm font-medium text-slate-300 bg-border rounded-lg hover:bg-[#2D3B4E] transition-colors flex items-center gap-2"
          @click="goBack"
        >
          <span class="material-symbols-rounded text-[18px]">arrow_back</span>
          Back
        </button>
        <button
          class="px-6 py-2.5 text-sm font-semibold rounded-lg transition-all flex items-center gap-2"
          :class="isValid
            ? 'bg-whitelabel text-white hover:bg-[#8B6ED0]'
            : 'bg-border text-slate-500 cursor-not-allowed'"
          :disabled="!isValid || saving"
          @click="goNext"
        >
          {{ saving ? 'Saving...' : 'Next: Configure' }}
          <span class="material-symbols-rounded text-[18px]">arrow_forward</span>
        </button>
      </div>
    </div>
  </div>
</template>
