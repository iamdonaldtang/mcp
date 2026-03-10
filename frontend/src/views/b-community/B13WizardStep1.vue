<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'
import type { CommunityWizardData } from '../../types/b-community'

const router = useRouter()

// === Wizard State ===
const form = ref<CommunityWizardData['step1']>({
  name: '',
  description: '',
  brandColor: '#48BB78',
})

const saving = ref(false)
const loadingDraft = ref(true)
const lastSavedAt = ref<string | null>(null)
const formTouched = ref(false)
let autoSaveTimer: ReturnType<typeof setInterval> | null = null

// === Stepper ===
const steps = [
  { num: 1, label: 'Customize' },
  { num: 2, label: 'Modules' },
  { num: 3, label: 'Quick Setup' },
  { num: 4, label: 'Preview & Publish' },
]
const currentStep = 1

// === Brand Color Presets ===
const colorPresets = [
  '#48BB78', '#5D7EF1', '#9B7EE0', '#ED8936',
  '#EF4444', '#EC4899', '#14B8A6', '#F59E0B',
]
const customHexInput = ref('')
const showCustomInput = ref(false)

// === Validation ===
const nameError = computed(() => {
  if (!formTouched.value && form.value.name === '') return ''
  if (form.value.name.length < 3) return 'Name must be at least 3 characters'
  if (form.value.name.length > 50) return 'Name must be 50 characters or less'
  return ''
})

const descError = computed(() => {
  if (!formTouched.value && form.value.description === '') return ''
  if (form.value.description.length > 0 && form.value.description.length < 10) return 'Description must be at least 10 characters'
  if (form.value.description.length > 500) return 'Description must be 500 characters or less'
  return ''
})

const isValid = computed(() => {
  return (
    form.value.name.length >= 3 &&
    form.value.name.length <= 50 &&
    form.value.description.length >= 10 &&
    form.value.description.length <= 500 &&
    /^#[0-9A-Fa-f]{6}$/.test(form.value.brandColor)
  )
})

// === Color Selection ===
function selectColor(color: string) {
  form.value.brandColor = color
  showCustomInput.value = false
  customHexInput.value = ''
}

function enableCustomColor() {
  showCustomInput.value = true
  customHexInput.value = form.value.brandColor
}

function applyCustomColor() {
  const hex = customHexInput.value.startsWith('#') ? customHexInput.value : `#${customHexInput.value}`
  if (/^#[0-9A-Fa-f]{6}$/.test(hex)) {
    form.value.brandColor = hex
  }
}

// === Draft Persistence ===
async function loadDraft() {
  loadingDraft.value = true
  try {
    const res = await api.get('/api/v1/community/wizard/draft')
    const draft = res.data?.data
    if (draft?.step1) {
      form.value.name = draft.step1.name || ''
      form.value.description = draft.step1.description || ''
      form.value.brandColor = draft.step1.brandColor || '#48BB78'
      if (draft.step1.brandColor && !colorPresets.includes(draft.step1.brandColor)) {
        showCustomInput.value = true
        customHexInput.value = draft.step1.brandColor
      }
    }
  } catch {
    // No draft yet — start fresh
  } finally {
    loadingDraft.value = false
  }
}

async function saveDraft(silent = true) {
  if (saving.value) return
  saving.value = true
  try {
    await api.post('/api/v1/community/wizard/draft', { step1: form.value })
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
  router.push('/b/community')
}

async function goNext() {
  formTouched.value = true
  if (!isValid.value) return
  await saveDraft(false)
  router.push('/b/community/wizard/step2')
}

// === Lifecycle ===
onMounted(async () => {
  await loadDraft()
  // Auto-save every 30 seconds
  autoSaveTimer = setInterval(() => {
    if (formTouched.value || form.value.name || form.value.description) {
      saveDraft()
    }
  }, 30000)
})

onUnmounted(() => {
  if (autoSaveTimer) clearInterval(autoSaveTimer)
})

// Mark form as touched on any change
watch(form, () => {
  formTouched.value = true
}, { deep: true })
</script>

<template>
  <div class="min-h-screen bg-[#0A0F1A]">
    <!-- Top Bar -->
    <div class="sticky top-0 z-20 bg-[#111B27] border-b border-[#1E293B] px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          class="flex items-center gap-1.5 text-sm text-slate-400 hover:text-white transition-colors"
          @click="goBack"
        >
          <span class="material-symbols-rounded text-[20px]">arrow_back</span>
          Back
        </button>
        <div class="w-px h-5 bg-[#1E293B]"></div>
        <h1 class="text-lg font-semibold text-white">Create Community</h1>
      </div>
      <div class="flex items-center gap-3">
        <span v-if="lastSavedAt" class="text-xs text-slate-500">
          Saved {{ lastSavedAt }}
        </span>
        <button
          class="px-4 py-2 text-sm font-medium text-slate-300 bg-[#1E293B] rounded-lg hover:bg-[#2D3B4E] transition-colors"
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
    <div v-if="loadingDraft" class="flex items-center justify-center py-24">
      <div class="animate-spin w-8 h-8 border-2 border-[#48BB78] border-t-transparent rounded-full"></div>
    </div>

    <!-- Main Content -->
    <div v-else class="px-6 pb-28 max-w-6xl mx-auto">
      <div class="grid grid-cols-1 lg:grid-cols-5 gap-8">
        <!-- Left: Form -->
        <div class="lg:col-span-3 space-y-6">
          <div class="bg-[#111B27] border border-[#1E293B] rounded-xl p-6 space-y-6">
            <h2 class="text-base font-semibold text-white">Community Details</h2>

            <!-- Name -->
            <div class="space-y-2">
              <label class="block text-sm font-medium text-slate-300">
                Community Name <span class="text-red-400">*</span>
              </label>
              <input
                v-model="form.name"
                type="text"
                maxlength="50"
                placeholder="e.g. MyCrypto Community"
                class="w-full px-4 py-2.5 text-sm text-white bg-[#0A0F1A] border rounded-lg outline-none transition-colors placeholder:text-slate-600"
                :class="nameError ? 'border-red-500 focus:border-red-500' : 'border-[#1E293B] focus:border-[#48BB78]'"
              />
              <div class="flex items-center justify-between">
                <span v-if="nameError" class="text-xs text-red-400">{{ nameError }}</span>
                <span v-else class="text-xs text-transparent">.</span>
                <span class="text-xs text-slate-500">{{ form.name.length }}/50</span>
              </div>
            </div>

            <!-- Description -->
            <div class="space-y-2">
              <label class="block text-sm font-medium text-slate-300">
                Description <span class="text-red-400">*</span>
              </label>
              <textarea
                v-model="form.description"
                maxlength="500"
                rows="4"
                placeholder="Describe what your community is about and what members can expect..."
                class="w-full px-4 py-2.5 text-sm text-white bg-[#0A0F1A] border rounded-lg outline-none resize-none transition-colors placeholder:text-slate-600"
                :class="descError ? 'border-red-500 focus:border-red-500' : 'border-[#1E293B] focus:border-[#48BB78]'"
              ></textarea>
              <div class="flex items-center justify-between">
                <span v-if="descError" class="text-xs text-red-400">{{ descError }}</span>
                <span v-else class="text-xs text-transparent">.</span>
                <span class="text-xs text-slate-500">{{ form.description.length }}/500</span>
              </div>
            </div>

            <!-- Brand Color -->
            <div class="space-y-3">
              <label class="block text-sm font-medium text-slate-300">Brand Color</label>
              <div class="flex flex-wrap items-center gap-3">
                <button
                  v-for="color in colorPresets"
                  :key="color"
                  class="w-9 h-9 rounded-lg transition-all"
                  :class="form.brandColor === color && !showCustomInput
                    ? 'ring-2 ring-white ring-offset-2 ring-offset-[#111B27] scale-110'
                    : 'hover:scale-105'"
                  :style="{ backgroundColor: color }"
                  :title="color"
                  @click="selectColor(color)"
                ></button>
                <!-- Custom color trigger -->
                <button
                  class="w-9 h-9 rounded-lg border-2 border-dashed flex items-center justify-center transition-colors"
                  :class="showCustomInput
                    ? 'border-[#48BB78] text-[#48BB78]'
                    : 'border-[#1E293B] text-slate-500 hover:border-slate-400 hover:text-slate-400'"
                  title="Custom color"
                  @click="enableCustomColor"
                >
                  <span class="material-symbols-rounded text-[18px]">palette</span>
                </button>
              </div>

              <!-- Custom hex input -->
              <div v-if="showCustomInput" class="flex items-center gap-2 mt-2">
                <div
                  class="w-9 h-9 rounded-lg border border-[#1E293B] shrink-0"
                  :style="{ backgroundColor: /^#[0-9A-Fa-f]{6}$/.test(customHexInput.startsWith('#') ? customHexInput : `#${customHexInput}`) ? (customHexInput.startsWith('#') ? customHexInput : `#${customHexInput}`) : '#333' }"
                ></div>
                <input
                  v-model="customHexInput"
                  type="text"
                  placeholder="#48BB78"
                  maxlength="7"
                  class="w-28 px-3 py-2 text-sm text-white bg-[#0A0F1A] border border-[#1E293B] rounded-lg outline-none focus:border-[#48BB78]"
                  @keyup.enter="applyCustomColor"
                  @blur="applyCustomColor"
                />
                <button
                  class="px-3 py-2 text-xs font-medium text-white bg-[#48BB78] rounded-lg hover:bg-[#38A169] transition-colors"
                  @click="applyCustomColor"
                >
                  Apply
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Live Preview -->
        <div class="lg:col-span-2">
          <div class="sticky top-24">
            <h3 class="text-sm font-medium text-slate-400 mb-3">Live Preview</h3>
            <div class="bg-[#111B27] border border-[#1E293B] rounded-xl overflow-hidden">
              <!-- Brand color bar -->
              <div
                class="h-2 transition-colors duration-300"
                :style="{ backgroundColor: form.brandColor }"
              ></div>

              <div class="p-5 space-y-4">
                <!-- Community name -->
                <div class="flex items-center gap-3">
                  <div
                    class="w-10 h-10 rounded-xl flex items-center justify-center text-white font-bold text-lg transition-colors duration-300"
                    :style="{ backgroundColor: form.brandColor }"
                  >
                    {{ form.name ? form.name[0]?.toUpperCase() : 'C' }}
                  </div>
                  <div>
                    <p class="text-sm font-semibold text-white">
                      {{ form.name || 'Community Name' }}
                    </p>
                    <p class="text-xs text-slate-500">your-project.taskon.xyz</p>
                  </div>
                </div>

                <!-- Description -->
                <p class="text-xs text-slate-400 leading-relaxed line-clamp-3">
                  {{ form.description || 'Your community description will appear here...' }}
                </p>

                <!-- Placeholder stats -->
                <div class="grid grid-cols-3 gap-3 pt-2 border-t border-[#1E293B]">
                  <div class="text-center">
                    <p class="text-lg font-bold text-white">0</p>
                    <p class="text-[10px] text-slate-500 uppercase tracking-wider">Members</p>
                  </div>
                  <div class="text-center">
                    <p class="text-lg font-bold text-white">0</p>
                    <p class="text-[10px] text-slate-500 uppercase tracking-wider">Tasks</p>
                  </div>
                  <div class="text-center">
                    <p class="text-lg font-bold text-white">0</p>
                    <p class="text-[10px] text-slate-500 uppercase tracking-wider">Points</p>
                  </div>
                </div>

                <!-- Preview modules hint -->
                <div class="pt-3 border-t border-[#1E293B]">
                  <div class="flex items-center gap-2 text-xs text-slate-500">
                    <span class="material-symbols-rounded text-[14px]">info</span>
                    Modules will be configured in the next step
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Action Bar -->
    <div class="fixed bottom-0 left-0 right-0 z-20 bg-[#111B27] border-t border-[#1E293B] px-6 py-4">
      <div class="max-w-6xl mx-auto flex items-center justify-end">
        <button
          class="px-6 py-2.5 text-sm font-semibold rounded-lg transition-all flex items-center gap-2"
          :class="isValid
            ? 'bg-[#48BB78] text-white hover:bg-[#38A169]'
            : 'bg-[#1E293B] text-slate-500 cursor-not-allowed'"
          :disabled="!isValid"
          @click="goNext"
        >
          Next: Configure Modules
          <span class="material-symbols-rounded text-[18px]">arrow_forward</span>
        </button>
      </div>
    </div>
  </div>
</template>
