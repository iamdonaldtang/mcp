<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()

// === State ===
const loading = ref(true)
const saving = ref(false)
const lastSavedAt = ref<string | null>(null)
let autoSaveTimer: ReturnType<typeof setInterval> | null = null

// === Stepper ===
const steps = [
  { num: 1, label: 'Path' },
  { num: 2, label: 'Configure' },
  { num: 3, label: 'Brand' },
  { num: 4, label: 'Preview & Publish' },
]
const currentStep = 3

// === Brand Form ===
const form = reactive({
  logoFile: null as File | null,
  logoPreview: '',
  primaryColor: '#9B7EE0',
  secondaryColor: '#5D7EF1',
  accentColor: '#F59E0B',
  backgroundColor: '#0A0F1A',
  headingFont: 'Inter',
  bodyFont: 'Inter',
  buttonStyle: 'filled' as 'filled' | 'outline' | 'rounded',
})

const fontOptions = [
  { value: 'Inter', label: 'Inter' },
  { value: 'Plus Jakarta Sans', label: 'Plus Jakarta Sans' },
  { value: 'Space Grotesk', label: 'Space Grotesk' },
  { value: 'Manrope', label: 'Manrope' },
]

const buttonStyles = [
  { value: 'filled', label: 'Filled', desc: 'Solid background color' },
  { value: 'outline', label: 'Outline', desc: 'Border only, transparent fill' },
  { value: 'rounded', label: 'Rounded', desc: 'Pill-shaped with full radius' },
]

// === Logo Upload ===
const dragOver = ref(false)
const logoError = ref('')

function handleLogoDrop(e: DragEvent) {
  dragOver.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file) validateAndSetLogo(file)
}

function handleLogoSelect(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) validateAndSetLogo(file)
}

function validateAndSetLogo(file: File) {
  logoError.value = ''
  const validTypes = ['image/svg+xml', 'image/png']
  if (!validTypes.includes(file.type)) {
    logoError.value = 'Only SVG and PNG files are supported.'
    return
  }
  if (file.size > 2 * 1024 * 1024) {
    logoError.value = 'File must be under 2MB.'
    return
  }
  form.logoFile = file
  form.logoPreview = URL.createObjectURL(file)
}

function removeLogo() {
  if (form.logoPreview) URL.revokeObjectURL(form.logoPreview)
  form.logoFile = null
  form.logoPreview = ''
}

// === Color inputs (hex validation) ===
function isValidHex(color: string): boolean {
  return /^#[0-9A-Fa-f]{6}$/.test(color)
}

// === Validation ===
const isValid = computed(() => {
  return (
    isValidHex(form.primaryColor) &&
    isValidHex(form.secondaryColor) &&
    isValidHex(form.accentColor) &&
    isValidHex(form.backgroundColor)
  )
})

// === Preview computed styles ===
const previewButtonClass = computed(() => {
  switch (form.buttonStyle) {
    case 'outline': return 'bg-transparent border-2'
    case 'rounded': return 'rounded-full'
    default: return 'rounded-lg'
  }
})

// === Draft Persistence ===
async function loadDraft() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/whitelabel/wizard/draft')
    const draft = res.data?.data
    if (draft?.brand) {
      const b = draft.brand
      if (b.primaryColor) form.primaryColor = b.primaryColor
      if (b.secondaryColor) form.secondaryColor = b.secondaryColor
      if (b.accentColor) form.accentColor = b.accentColor
      if (b.backgroundColor) form.backgroundColor = b.backgroundColor
      if (b.headingFont) form.headingFont = b.headingFont
      if (b.bodyFont) form.bodyFont = b.bodyFont
      if (b.buttonStyle) form.buttonStyle = b.buttonStyle
      if (b.logoUrl) form.logoPreview = b.logoUrl
    }
  } catch {
    // No draft
  } finally {
    loading.value = false
  }
}

async function saveDraft(silent = true) {
  if (saving.value) return
  saving.value = true
  try {
    await api.put('/api/v1/whitelabel/brand', {
      primaryColor: form.primaryColor,
      secondaryColor: form.secondaryColor,
      accentColor: form.accentColor,
      backgroundColor: form.backgroundColor,
      headingFont: form.headingFont,
      bodyFont: form.bodyFont,
      buttonStyle: form.buttonStyle,
    })

    // Upload logo if changed
    if (form.logoFile) {
      const formData = new FormData()
      formData.append('logo', form.logoFile)
      await api.post('/api/v1/whitelabel/brand/logo', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      form.logoFile = null // Reset after upload
    }

    lastSavedAt.value = new Date().toLocaleTimeString()
  } catch {
    if (!silent) { /* toast */ }
  } finally {
    saving.value = false
  }
}

// === Navigation ===
function goBack() {
  router.push('/b/whitelabel/wizard/step-2')
}

async function goNext() {
  if (!isValid.value) return
  await saveDraft(false)
  router.push('/b/whitelabel/wizard/step-4')
}

// === Lifecycle ===
onMounted(async () => {
  await loadDraft()
  autoSaveTimer = setInterval(() => { saveDraft() }, 30000)
})

onUnmounted(() => {
  if (autoSaveTimer) clearInterval(autoSaveTimer)
  if (form.logoPreview && form.logoFile) URL.revokeObjectURL(form.logoPreview)
})

watch(
  () => [form.primaryColor, form.secondaryColor, form.accentColor, form.backgroundColor],
  () => { /* triggers reactivity for live preview */ },
)
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
        <span v-if="lastSavedAt" class="text-xs text-slate-500">Saved {{ lastSavedAt }}</span>
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

        <!-- Left: Form -->
        <div class="lg:col-span-3 space-y-6">

          <!-- Logo Upload -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
            <h3 class="text-base font-semibold text-white">Logo</h3>
            <div
              v-if="!form.logoPreview"
              class="border-2 border-dashed rounded-xl p-8 text-center transition-colors cursor-pointer"
              :class="dragOver ? 'border-whitelabel bg-whitelabel/5' : 'border-border hover:border-[#2D3B4E]'"
              @dragover.prevent="dragOver = true"
              @dragleave="dragOver = false"
              @drop.prevent="handleLogoDrop"
              @click="($refs.logoInput as HTMLInputElement).click()"
            >
              <span class="material-symbols-rounded text-[36px] text-slate-500 mb-2">cloud_upload</span>
              <p class="text-sm text-slate-400">Drag and drop or click to upload</p>
              <p class="text-xs text-slate-600 mt-1">SVG or PNG, max 2MB</p>
              <input
                ref="logoInput"
                type="file"
                accept=".svg,.png,image/svg+xml,image/png"
                class="hidden"
                @change="handleLogoSelect"
              />
            </div>
            <div v-else class="flex items-center gap-4">
              <div class="w-16 h-16 rounded-xl border border-border bg-page-bg flex items-center justify-center overflow-hidden">
                <img :src="form.logoPreview" alt="Logo preview" class="max-w-full max-h-full object-contain" />
              </div>
              <button
                class="text-sm text-red-400 hover:text-red-300 transition-colors flex items-center gap-1"
                @click="removeLogo"
              >
                <span class="material-symbols-rounded text-[16px]">delete</span>
                Remove
              </button>
            </div>
            <p v-if="logoError" class="text-xs text-red-400">{{ logoError }}</p>
          </div>

          <!-- Colors -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-5">
            <h3 class="text-base font-semibold text-white">Colors</h3>

            <div class="grid grid-cols-2 gap-4">
              <!-- Primary Color -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Primary Color</label>
                <div class="flex items-center gap-2">
                  <input v-model="form.primaryColor" type="color" class="w-9 h-9 rounded-lg border border-border cursor-pointer bg-transparent" />
                  <input
                    v-model="form.primaryColor"
                    type="text"
                    maxlength="7"
                    class="flex-1 px-3 py-2 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel font-mono"
                  />
                </div>
              </div>

              <!-- Secondary Color -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Secondary Color</label>
                <div class="flex items-center gap-2">
                  <input v-model="form.secondaryColor" type="color" class="w-9 h-9 rounded-lg border border-border cursor-pointer bg-transparent" />
                  <input
                    v-model="form.secondaryColor"
                    type="text"
                    maxlength="7"
                    class="flex-1 px-3 py-2 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel font-mono"
                  />
                </div>
              </div>

              <!-- Accent Color -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Accent Color</label>
                <div class="flex items-center gap-2">
                  <input v-model="form.accentColor" type="color" class="w-9 h-9 rounded-lg border border-border cursor-pointer bg-transparent" />
                  <input
                    v-model="form.accentColor"
                    type="text"
                    maxlength="7"
                    class="flex-1 px-3 py-2 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel font-mono"
                  />
                </div>
              </div>

              <!-- Background Color -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Background Color</label>
                <div class="flex items-center gap-2">
                  <input v-model="form.backgroundColor" type="color" class="w-9 h-9 rounded-lg border border-border cursor-pointer bg-transparent" />
                  <input
                    v-model="form.backgroundColor"
                    type="text"
                    maxlength="7"
                    class="flex-1 px-3 py-2 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel font-mono"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Typography -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-5">
            <h3 class="text-base font-semibold text-white">Typography</h3>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Heading Font</label>
                <select
                  v-model="form.headingFont"
                  class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel appearance-none cursor-pointer"
                >
                  <option v-for="font in fontOptions" :key="font.value" :value="font.value">{{ font.label }}</option>
                </select>
              </div>
              <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Body Font</label>
                <select
                  v-model="form.bodyFont"
                  class="w-full px-4 py-2.5 text-sm text-white bg-page-bg border border-border rounded-lg outline-none focus:border-whitelabel appearance-none cursor-pointer"
                >
                  <option v-for="font in fontOptions" :key="font.value" :value="font.value">{{ font.label }}</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Button Style -->
          <div class="bg-card-bg border border-border rounded-xl p-6 space-y-4">
            <h3 class="text-base font-semibold text-white">Button Style</h3>
            <div class="grid grid-cols-3 gap-3">
              <button
                v-for="style in buttonStyles"
                :key="style.value"
                class="p-4 rounded-lg border transition-colors text-left"
                :class="form.buttonStyle === style.value
                  ? 'border-whitelabel bg-whitelabel/5'
                  : 'border-border hover:border-[#2D3B4E]'"
                @click="form.buttonStyle = style.value as typeof form.buttonStyle"
              >
                <p class="text-sm font-medium text-white">{{ style.label }}</p>
                <p class="text-xs text-slate-500 mt-0.5">{{ style.desc }}</p>
              </button>
            </div>
          </div>
        </div>

        <!-- Right: Live Preview -->
        <div class="lg:col-span-2">
          <div class="sticky top-24">
            <h3 class="text-sm font-medium text-slate-400 mb-3">Live Preview</h3>
            <div
              class="rounded-xl border border-border overflow-hidden transition-colors"
              :style="{ backgroundColor: form.backgroundColor }"
            >
              <!-- Header bar -->
              <div class="px-5 py-4 flex items-center justify-between" :style="{ backgroundColor: form.primaryColor + '15' }">
                <div class="flex items-center gap-3">
                  <div
                    v-if="form.logoPreview"
                    class="w-8 h-8 rounded-lg overflow-hidden"
                  >
                    <img :src="form.logoPreview" alt="Logo" class="w-full h-full object-contain" />
                  </div>
                  <div
                    v-else
                    class="w-8 h-8 rounded-lg flex items-center justify-center text-white text-sm font-bold"
                    :style="{ backgroundColor: form.primaryColor }"
                  >
                    W
                  </div>
                  <span
                    class="text-sm font-semibold"
                    :style="{ fontFamily: form.headingFont, color: '#F1F5F9' }"
                  >
                    Your Community
                  </span>
                </div>
                <div class="flex items-center gap-2">
                  <div class="w-6 h-6 rounded-full bg-white/10"></div>
                </div>
              </div>

              <!-- Content -->
              <div class="p-5 space-y-4">
                <!-- Heading -->
                <h4
                  class="text-lg font-bold"
                  :style="{ fontFamily: form.headingFont, color: '#F1F5F9' }"
                >
                  Welcome Back
                </h4>
                <p
                  class="text-sm leading-relaxed"
                  :style="{ fontFamily: form.bodyFont, color: '#94A3B8' }"
                >
                  Complete tasks, earn points, and climb the leaderboard.
                </p>

                <!-- Stats row -->
                <div class="grid grid-cols-3 gap-3">
                  <div
                    v-for="(stat, i) in [{ label: 'Points', val: '2,450' }, { label: 'Rank', val: '#12' }, { label: 'Tasks', val: '8/15' }]"
                    :key="i"
                    class="text-center p-3 rounded-lg"
                    :style="{ backgroundColor: form.primaryColor + '10' }"
                  >
                    <p class="text-base font-bold" :style="{ color: form.primaryColor }">{{ stat.val }}</p>
                    <p class="text-[10px] mt-0.5" :style="{ fontFamily: form.bodyFont, color: '#64748B' }">{{ stat.label }}</p>
                  </div>
                </div>

                <!-- Card -->
                <div class="p-4 rounded-lg border" :style="{ borderColor: form.primaryColor + '30' }">
                  <div class="flex items-center justify-between mb-2">
                    <span class="text-sm font-medium" :style="{ fontFamily: form.headingFont, color: '#F1F5F9' }">Daily Task</span>
                    <span
                      class="px-2 py-0.5 text-[10px] font-bold rounded-full"
                      :style="{ backgroundColor: form.accentColor + '20', color: form.accentColor }"
                    >
                      +50 PTS
                    </span>
                  </div>
                  <p class="text-xs mb-3" :style="{ fontFamily: form.bodyFont, color: '#94A3B8' }">
                    Follow our Twitter account
                  </p>
                  <button
                    class="w-full py-2 text-sm font-semibold text-white transition-colors"
                    :class="previewButtonClass"
                    :style="{
                      backgroundColor: form.buttonStyle === 'outline' ? 'transparent' : form.primaryColor,
                      borderColor: form.primaryColor,
                      color: form.buttonStyle === 'outline' ? form.primaryColor : '#FFFFFF',
                      borderRadius: form.buttonStyle === 'rounded' ? '9999px' : form.buttonStyle === 'outline' ? '8px' : '8px',
                    }"
                  >
                    Start Task
                  </button>
                </div>

                <!-- Secondary button -->
                <button
                  class="w-full py-2 text-sm font-medium border rounded-lg transition-colors"
                  :style="{ borderColor: form.secondaryColor + '40', color: form.secondaryColor }"
                >
                  View All Tasks
                </button>
              </div>
            </div>
          </div>
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
          :class="isValid
            ? 'bg-whitelabel text-white hover:bg-[#8B6ED0]'
            : 'bg-border text-slate-500 cursor-not-allowed'"
          :disabled="!isValid"
          @click="goNext"
        >
          Next: Preview & Publish
          <span class="material-symbols-rounded text-[18px]">arrow_forward</span>
        </button>
      </div>
    </div>
  </div>
</template>
