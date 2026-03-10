<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../../api/client'
import Modal from './Modal.vue'

const props = defineProps<{
  open: boolean
  entityType: string
  entityId: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'confirmed'): void
}>()

interface ReadinessItem {
  key: string
  label: string
  passed: boolean
  message?: string
}

const loading = ref(false)
const checks = ref<ReadinessItem[]>([])
const allPassed = ref(false)

const CACHE_TTL = 5 * 60 * 1000 // 5 minutes

function getCacheKey() {
  return `${props.entityType}_${props.entityId}_readiness`
}

function getCached(): ReadinessItem[] | null {
  try {
    const raw = sessionStorage.getItem(getCacheKey())
    if (!raw) return null
    const cached = JSON.parse(raw)
    if (Date.now() - cached.timestamp > CACHE_TTL) {
      sessionStorage.removeItem(getCacheKey())
      return null
    }
    return cached.data
  } catch {
    return null
  }
}

function setCache(data: ReadinessItem[]) {
  sessionStorage.setItem(getCacheKey(), JSON.stringify({ data, timestamp: Date.now() }))
}

async function fetchReadiness() {
  const cached = getCached()
  if (cached) {
    checks.value = cached
    allPassed.value = cached.every(c => c.passed)
    return
  }

  loading.value = true
  try {
    const res = await api.get(`/api/v1/readiness`, {
      params: { type: props.entityType, id: props.entityId },
    })
    checks.value = res.data.data
    allPassed.value = checks.value.every(c => c.passed)
    setCache(checks.value)
  } catch {
    checks.value = [
      { key: 'subscription', label: 'Active Subscription', passed: false, message: 'Unable to verify' },
      { key: 'twitter', label: 'Twitter Authorization', passed: false, message: 'Unable to verify' },
    ]
    allPassed.value = false
  } finally {
    loading.value = false
  }
}

function handleConfirm() {
  if (allPassed.value) {
    emit('confirmed')
    emit('close')
  }
}

onMounted(() => {
  if (props.open) fetchReadiness()
})
</script>

<template>
  <Modal :open="open" title="Publish Readiness Check" @close="emit('close')">
    <div v-if="loading" class="flex items-center justify-center py-8">
      <span class="material-symbols-rounded animate-spin text-2xl text-text-muted">progress_activity</span>
    </div>
    <div v-else class="space-y-3">
      <p class="text-sm text-text-secondary mb-4">
        All checks must pass before publishing.
      </p>
      <div
        v-for="check in checks"
        :key="check.key"
        class="flex items-center gap-3 p-3 rounded-lg border"
        :class="check.passed ? 'border-status-active/30 bg-status-active/5' : 'border-status-error/30 bg-status-error/5'"
      >
        <span
          class="material-symbols-rounded text-lg"
          :class="check.passed ? 'text-status-active' : 'text-status-error'"
        >
          {{ check.passed ? 'check_circle' : 'cancel' }}
        </span>
        <div class="flex-1">
          <div class="text-sm font-medium text-text-primary">{{ check.label }}</div>
          <div v-if="check.message" class="text-xs text-text-muted mt-0.5">{{ check.message }}</div>
        </div>
      </div>
    </div>
    <template #footer>
      <button
        class="px-4 py-2 text-sm font-medium text-text-secondary hover:text-text-primary transition-colors"
        @click="emit('close')"
      >
        Cancel
      </button>
      <button
        class="px-4 py-2 text-sm font-medium rounded-lg transition-colors"
        :class="allPassed ? 'bg-status-active text-white hover:bg-status-active/90' : 'bg-text-muted/20 text-text-muted cursor-not-allowed'"
        :disabled="!allPassed"
        @click="handleConfirm"
      >
        Confirm & Publish
      </button>
    </template>
  </Modal>
</template>
