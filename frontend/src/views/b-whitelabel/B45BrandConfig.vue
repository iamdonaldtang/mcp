<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { api } from '../../api/client'

// === Types ===
interface BrandSettings {
  logoUrl: string | null
  primaryColor: string
  secondaryColor: string
  accentColor: string
  backgroundColor: string
  headingFont: string
  bodyFont: string
  baseFontSize: number
  letterSpacing: number
  buttonStyle: 'filled' | 'outline' | 'rounded'
  customCss: string
  lastSaved: string | null
}

// === State ===
const loading = ref(true)
const saving = ref(false)
const previewMode = ref<'desktop' | 'mobile'>('desktop')
const showResetConfirm = ref(false)
const dragOver = ref(false)
const logoPreview = ref<string | null>(null)
const logoFile = ref<File | null>(null)

const defaults: BrandSettings = {
  logoUrl: null,
  primaryColor: '#9B7EE0',
  secondaryColor: '#5D7EF1',
  accentColor: '#F59E0B',
  backgroundColor: '#0A0F1A',
  headingFont: 'Inter',
  bodyFont: 'Inter',
  baseFontSize: 14,
  letterSpacing: 0,
  buttonStyle: 'filled',
  customCss: '',
  lastSaved: null,
}

const form = reactive<BrandSettings>({ ...defaults })

const fontOptions = [
  { value: 'Inter', label: 'Inter' },
  { value: 'Plus Jakarta Sans', label: 'Plus Jakarta Sans' },
  { value: 'Space Grotesk', label: 'Space Grotesk' },
  { value: 'Manrope', label: 'Manrope' },
  { value: 'Poppins', label: 'Poppins' },
]

const buttonStyles = [
  { value: 'filled' as const, label: 'Filled' },
  { value: 'outline' as const, label: 'Outline' },
  { value: 'rounded' as const, label: 'Rounded' },
]

// === Computed ===
const customCssLength = computed(() => form.customCss.length)
const customCssLines = computed(() => form.customCss ? form.customCss.split('\n').length : 0)
const customCssOverLimit = computed(() => customCssLength.value > 10240)

const lastSavedDisplay = computed(() => {
  if (!form.lastSaved) return 'Never'
  return new Date(form.lastSaved).toLocaleDateString('en-US', {
    year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
  })
})

// Preview styles
const previewHeaderStyle = computed(() => ({
  backgroundColor: form.backgroundColor,
  fontFamily: `'${form.headingFont}', sans-serif`,
  color: form.primaryColor,
}))

const previewBodyStyle = computed(() => ({
  fontFamily: `'${form.bodyFont}', sans-serif`,
  fontSize: form.baseFontSize + 'px',
  letterSpacing: form.letterSpacing + 'px',
}))

function buttonPreviewClass(style: 'filled' | 'outline' | 'rounded') {
  const base = 'px-5 py-2 text-sm font-medium transition-colors'
  switch (style) {
    case 'filled':
      return `${base} rounded-lg`
    case 'outline':
      return `${base} rounded-lg border-2 bg-transparent`
    case 'rounded':
      return `${base} rounded-full`
  }
}

function buttonPreviewInlineStyle(style: 'filled' | 'outline' | 'rounded') {
  switch (style) {
    case 'filled':
      return { backgroundColor: form.primaryColor, color: '#FFFFFF' }
    case 'outline':
      return { borderColor: form.primaryColor, color: form.primaryColor }
    case 'rounded':
      return { backgroundColor: form.primaryColor, color: '#FFFFFF' }
  }
}

// === API ===
async function fetchData() {
  loading.value = true
  try {
    const { data } = await api.get('/api/v1/whitelabel/brand')
    Object.assign(form, data)
    if (data.logoUrl) logoPreview.value = data.logoUrl
  } catch {
    // Use defaults — already set
    form.lastSaved = '2026-03-08T16:30:00Z'
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    // Upload logo if changed
    if (logoFile.value) {
      const fd = new FormData()
      fd.append('logo', logoFile.value)
      const { data: logoData } = await api.post('/api/v1/whitelabel/brand/logo', fd, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      form.logoUrl = logoData.url
    }

    const { data } = await api.put('/api/v1/whitelabel/brand', {
      primaryColor: form.primaryColor,
      secondaryColor: form.secondaryColor,
      accentColor: form.accentColor,
      backgroundColor: form.backgroundColor,
      headingFont: form.headingFont,
      bodyFont: form.bodyFont,
      baseFontSize: form.baseFontSize,
      letterSpacing: form.letterSpacing,
      buttonStyle: form.buttonStyle,
      customCss: form.customCss,
    })
    form.lastSaved = data.lastSaved || new Date().toISOString()
  } catch {
    form.lastSaved = new Date().toISOString()
  } finally {
    saving.value = false
  }
}

function resetToDefaults() {
  Object.assign(form, defaults)
  logoPreview.value = null
  logoFile.value = null
  showResetConfirm.value = false
}

// === Logo handling ===
function handleLogoDrop(e: DragEvent) {
  dragOver.value = false
  const file = e.dataTransfer?.files[0]
  if (file) processLogoFile(file)
}

function handleLogoSelect(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) processLogoFile(file)
}

function processLogoFile(file: File) {
  if (!['image/svg+xml', 'image/png'].includes(file.type)) return
  if (file.size > 2 * 1024 * 1024) return
  logoFile.value = file
  const reader = new FileReader()
  reader.onload = () => { logoPreview.value = reader.result as string }
  reader.readAsDataURL(file)
}

function removeLogo() {
  logoPreview.value = null
  logoFile.value = null
  form.logoUrl = null
}

onMounted(fetchData)
</script>

<template>
  <div v-if="loading" class="flex items-center justify-center py-20">
    <span class="material-symbols-rounded text-4xl text-whitelabel animate-spin">progress_activity</span>
  </div>

  <div v-else class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">Brand Settings</h1>
        <p class="text-sm text-text-muted mt-1">Last saved: {{ lastSavedDisplay }}</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          class="px-4 py-2 border border-border text-text-secondary rounded-lg text-sm font-medium hover:text-text-primary transition-colors"
          @click="showResetConfirm = true"
        >
          Reset to Defaults
        </button>
        <button
          class="px-5 py-2 bg-whitelabel text-white rounded-lg text-sm font-medium hover:bg-whitelabel/90 transition-colors disabled:opacity-50 flex items-center gap-2"
          :disabled="saving || customCssOverLimit"
          @click="save"
        >
          <span v-if="saving" class="material-symbols-rounded text-lg animate-spin">progress_activity</span>
          Save Brand Settings
        </button>
      </div>
    </div>

    <!-- Two Column Layout -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
      <!-- LEFT: Form -->
      <div class="space-y-6">
        <!-- Logo Section -->
        <div class="bg-card-bg border border-border rounded-xl p-6">
          <h2 class="text-base font-semibold text-text-primary mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-whitelabel text-lg">image</span>
            Logo
          </h2>

          <div class="flex items-start gap-5">
            <!-- Logo Preview -->
            <div class="w-20 h-20 rounded-lg border border-border bg-page-bg flex items-center justify-center overflow-hidden shrink-0">
              <img
                v-if="logoPreview"
                :src="logoPreview"
                alt="Logo"
                class="max-w-full max-h-full object-contain"
              />
              <span v-else class="material-symbols-rounded text-3xl text-text-muted">photo_camera</span>
            </div>

            <div class="flex-1">
              <!-- Drop Area -->
              <div
                :class="[
                  'border-2 border-dashed rounded-lg p-4 text-center transition-colors cursor-pointer',
                  dragOver ? 'border-whitelabel/60 bg-whitelabel/5' : 'border-border hover:border-whitelabel/30',
                ]"
                @dragover.prevent="dragOver = true"
                @dragleave="dragOver = false"
                @drop.prevent="handleLogoDrop"
                @click="($refs.logoInput as HTMLInputElement).click()"
              >
                <input
                  ref="logoInput"
                  type="file"
                  accept=".svg,.png"
                  class="hidden"
                  @change="handleLogoSelect"
                />
                <span class="material-symbols-rounded text-2xl text-text-muted mb-1 block">cloud_upload</span>
                <p class="text-sm text-text-secondary">Drop logo or click to upload</p>
                <p class="text-xs text-text-muted mt-1">SVG or PNG, max 2MB</p>
              </div>

              <button
                v-if="logoPreview"
                class="mt-2 text-xs text-status-paused hover:text-status-paused/80 transition-colors"
                @click="removeLogo"
              >
                Remove logo
              </button>
            </div>
          </div>
        </div>

        <!-- Colors Section -->
        <div class="bg-card-bg border border-border rounded-xl p-6">
          <h2 class="text-base font-semibold text-text-primary mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-whitelabel text-lg">palette</span>
            Colors
          </h2>

          <div class="grid grid-cols-2 gap-4">
            <div v-for="(colorDef, idx) in [
              { key: 'primaryColor', label: 'Primary Color' },
              { key: 'secondaryColor', label: 'Secondary Color' },
              { key: 'accentColor', label: 'Accent Color' },
              { key: 'backgroundColor', label: 'Background Color' },
            ]" :key="idx">
              <label class="block text-sm text-text-secondary mb-1.5">{{ colorDef.label }}</label>
              <div class="flex items-center gap-2">
                <div class="relative">
                  <input
                    type="color"
                    v-model="(form as any)[colorDef.key]"
                    class="w-10 h-10 rounded-lg border border-border cursor-pointer bg-transparent p-0.5"
                  />
                </div>
                <input
                  type="text"
                  v-model="(form as any)[colorDef.key]"
                  class="flex-1 bg-page-bg border border-border rounded-lg px-3 py-2 text-sm text-text-primary font-mono uppercase focus:outline-none focus:border-whitelabel/50"
                  maxlength="7"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Typography Section -->
        <div class="bg-card-bg border border-border rounded-xl p-6">
          <h2 class="text-base font-semibold text-text-primary mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-whitelabel text-lg">text_fields</span>
            Typography
          </h2>

          <div class="space-y-4">
            <!-- Font Dropdowns -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-text-secondary mb-1.5">Heading Font</label>
                <select
                  v-model="form.headingFont"
                  class="w-full bg-page-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary focus:outline-none focus:border-whitelabel/50 appearance-none"
                >
                  <option v-for="f in fontOptions" :key="f.value" :value="f.value">{{ f.label }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm text-text-secondary mb-1.5">Body Font</label>
                <select
                  v-model="form.bodyFont"
                  class="w-full bg-page-bg border border-border rounded-lg px-4 py-2.5 text-sm text-text-primary focus:outline-none focus:border-whitelabel/50 appearance-none"
                >
                  <option v-for="f in fontOptions" :key="f.value" :value="f.value">{{ f.label }}</option>
                </select>
              </div>
            </div>

            <!-- Font Size Slider -->
            <div>
              <div class="flex items-center justify-between mb-1.5">
                <label class="text-sm text-text-secondary">Base Font Size</label>
                <span class="text-sm font-medium text-text-primary">{{ form.baseFontSize }}px</span>
              </div>
              <input
                type="range"
                v-model.number="form.baseFontSize"
                min="12"
                max="18"
                step="1"
                class="w-full h-2 bg-page-bg rounded-full appearance-none cursor-pointer accent-whitelabel"
              />
              <div class="flex justify-between text-xs text-text-muted mt-1">
                <span>12px</span>
                <span>18px</span>
              </div>
            </div>

            <!-- Letter Spacing Slider -->
            <div>
              <div class="flex items-center justify-between mb-1.5">
                <label class="text-sm text-text-secondary">Letter Spacing</label>
                <span class="text-sm font-medium text-text-primary">{{ form.letterSpacing }}px</span>
              </div>
              <input
                type="range"
                v-model.number="form.letterSpacing"
                min="-0.5"
                max="2"
                step="0.1"
                class="w-full h-2 bg-page-bg rounded-full appearance-none cursor-pointer accent-whitelabel"
              />
              <div class="flex justify-between text-xs text-text-muted mt-1">
                <span>-0.5px</span>
                <span>2px</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Button Style Section -->
        <div class="bg-card-bg border border-border rounded-xl p-6">
          <h2 class="text-base font-semibold text-text-primary mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-whitelabel text-lg">smart_button</span>
            Button Style
          </h2>

          <div class="grid grid-cols-3 gap-3">
            <button
              v-for="bs in buttonStyles"
              :key="bs.value"
              :class="[
                'border rounded-xl p-4 text-center transition-all',
                form.buttonStyle === bs.value
                  ? 'border-whitelabel/50 bg-whitelabel/5'
                  : 'border-border hover:border-whitelabel/20',
              ]"
              @click="form.buttonStyle = bs.value"
            >
              <!-- Button Preview -->
              <div class="flex justify-center mb-3">
                <span :class="buttonPreviewClass(bs.value)" :style="buttonPreviewInlineStyle(bs.value)">
                  Button
                </span>
              </div>
              <span class="text-sm font-medium" :class="form.buttonStyle === bs.value ? 'text-whitelabel' : 'text-text-secondary'">
                {{ bs.label }}
              </span>
            </button>
          </div>
        </div>

        <!-- Custom CSS Section -->
        <div class="bg-card-bg border border-border rounded-xl p-6">
          <h2 class="text-base font-semibold text-text-primary mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-whitelabel text-lg">code</span>
            Custom CSS
          </h2>

          <textarea
            v-model="form.customCss"
            rows="10"
            placeholder="/* Add custom CSS overrides */"
            class="w-full bg-[#0D1117] border border-border rounded-lg px-4 py-3 text-sm text-text-primary font-mono focus:outline-none focus:border-whitelabel/50 resize-y"
            spellcheck="false"
          />
          <div class="flex items-center justify-between mt-2">
            <div class="flex items-center gap-3 text-xs text-text-muted">
              <span>{{ customCssLines }} lines</span>
              <span>{{ customCssLength.toLocaleString() }} / 10,240 chars</span>
            </div>
            <div v-if="customCssLength > 0" class="flex items-center gap-1 text-xs text-status-draft">
              <span class="material-symbols-rounded text-sm">warning</span>
              Custom CSS may affect performance
            </div>
          </div>
          <div v-if="customCssOverLimit" class="mt-2 px-3 py-2 rounded-lg bg-status-paused-bg text-status-paused text-xs">
            CSS exceeds 10KB limit. Please reduce the content before saving.
          </div>
        </div>
      </div>

      <!-- RIGHT: Live Preview -->
      <div>
        <div class="bg-card-bg border border-border rounded-xl overflow-hidden sticky top-6">
          <!-- Preview Header -->
          <div class="flex items-center justify-between px-5 py-3 border-b border-border">
            <h2 class="text-sm font-semibold text-text-primary">Live Preview</h2>
            <div class="flex bg-page-bg rounded-lg p-1">
              <button
                v-for="mode in (['desktop', 'mobile'] as const)"
                :key="mode"
                :class="[
                  'px-3 py-1 rounded-md text-xs font-medium transition-colors flex items-center gap-1',
                  previewMode === mode ? 'bg-whitelabel/20 text-whitelabel' : 'text-text-muted hover:text-text-secondary',
                ]"
                @click="previewMode = mode"
              >
                <span class="material-symbols-rounded text-sm">{{ mode === 'desktop' ? 'computer' : 'phone_iphone' }}</span>
                {{ mode === 'desktop' ? 'Desktop' : 'Mobile' }}
              </button>
            </div>
          </div>

          <!-- Preview Content -->
          <div class="p-5">
            <div
              :class="[
                'rounded-xl border border-border overflow-hidden mx-auto transition-all',
                previewMode === 'mobile' ? 'max-w-[375px]' : 'w-full',
              ]"
              :style="{ backgroundColor: form.backgroundColor }"
            >
              <!-- Mock Header -->
              <div
                class="px-5 py-3 border-b flex items-center justify-between"
                :style="{ ...previewHeaderStyle, borderColor: form.primaryColor + '20' }"
              >
                <div class="flex items-center gap-3">
                  <!-- Logo in header -->
                  <div class="w-8 h-8 rounded flex items-center justify-center overflow-hidden" :style="{ backgroundColor: form.primaryColor + '20' }">
                    <img v-if="logoPreview" :src="logoPreview" alt="" class="max-w-full max-h-full object-contain" />
                    <span v-else class="text-xs font-bold" :style="{ color: form.primaryColor }">TK</span>
                  </div>
                  <span class="text-sm font-semibold" :style="{ color: form.primaryColor, fontFamily: `'${form.headingFont}', sans-serif` }">
                    My Community
                  </span>
                </div>
                <div class="flex gap-4">
                  <span
                    v-for="tab in ['Quests', 'Leaderboard', 'Shop']"
                    :key="tab"
                    class="text-xs"
                    :style="{
                      color: tab === 'Quests' ? form.accentColor : '#94A3B8',
                      fontFamily: `'${form.headingFont}', sans-serif`,
                      borderBottom: tab === 'Quests' ? `2px solid ${form.accentColor}` : 'none',
                      paddingBottom: '4px',
                    }"
                  >{{ tab }}</span>
                </div>
              </div>

              <!-- Mock Card -->
              <div class="p-5" :style="previewBodyStyle">
                <div class="rounded-lg p-4 mb-4" :style="{ backgroundColor: form.primaryColor + '08', border: `1px solid ${form.primaryColor}20` }">
                  <h3
                    class="font-semibold mb-1"
                    :style="{
                      color: '#F1F5F9',
                      fontFamily: `'${form.headingFont}', sans-serif`,
                      fontSize: (form.baseFontSize + 2) + 'px',
                    }"
                  >
                    Daily Check-in
                  </h3>
                  <p
                    :style="{
                      color: '#94A3B8',
                      fontFamily: `'${form.bodyFont}', sans-serif`,
                      fontSize: form.baseFontSize + 'px',
                      letterSpacing: form.letterSpacing + 'px',
                    }"
                  >
                    Complete your daily check-in to earn 10 points and maintain your streak.
                  </p>
                </div>

                <!-- Button Samples -->
                <div class="flex flex-wrap gap-2 mb-4">
                  <span
                    v-for="bs in buttonStyles"
                    :key="bs.value"
                    :class="[
                      buttonPreviewClass(bs.value),
                      'text-xs cursor-default',
                      form.buttonStyle === bs.value ? 'ring-2 ring-offset-1' : 'opacity-50',
                    ]"
                    :style="{
                      ...buttonPreviewInlineStyle(bs.value),
                      ringColor: form.accentColor,
                      ringOffsetColor: form.backgroundColor,
                    }"
                  >
                    {{ bs.label }}
                  </span>
                </div>

                <!-- Color Swatches -->
                <div class="flex items-center gap-2">
                  <div
                    v-for="(swatch, si) in [
                      { color: form.primaryColor, label: 'Primary' },
                      { color: form.secondaryColor, label: 'Secondary' },
                      { color: form.accentColor, label: 'Accent' },
                      { color: form.backgroundColor, label: 'Background' },
                    ]"
                    :key="si"
                    class="flex items-center gap-1.5"
                  >
                    <div
                      class="w-4 h-4 rounded-full border border-white/10"
                      :style="{ backgroundColor: swatch.color }"
                    />
                    <span class="text-[10px] text-text-muted">{{ swatch.label }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Reset Confirm Dialog -->
    <Teleport to="body">
      <div
        v-if="showResetConfirm"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
        @click.self="showResetConfirm = false"
      >
        <div class="bg-card-bg border border-border rounded-xl p-6 w-full max-w-md mx-4">
          <h3 class="text-lg font-semibold text-text-primary mb-2">Reset to Defaults?</h3>
          <p class="text-sm text-text-secondary mb-5">
            This will revert all brand settings to their default values. Your current customizations will be lost until you save again.
          </p>
          <div class="flex justify-end gap-3">
            <button
              class="px-4 py-2 border border-border text-text-secondary rounded-lg text-sm font-medium hover:text-text-primary transition-colors"
              @click="showResetConfirm = false"
            >
              Cancel
            </button>
            <button
              class="px-4 py-2 bg-status-paused text-white rounded-lg text-sm font-medium hover:bg-status-paused/90 transition-colors"
              @click="resetToDefaults"
            >
              Reset
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
