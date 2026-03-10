<script setup lang="ts">
import { computed } from 'vue'
import type { CTask } from '../../types/c-end'

const props = defineProps<{ task: CTask }>()

const emit = defineEmits<{
  (e: 'action', taskId: string): void
}>()

const buttonConfig = computed(() => {
  switch (props.task.status) {
    case 'available': return { label: 'Start', class: 'bg-c-accent text-black' }
    case 'in_progress': return { label: 'Continue', class: 'border border-c-accent text-c-accent' }
    case 'completed': return { label: 'Claim', class: 'bg-status-active text-white' }
    case 'claimed': return { label: 'Done ✓', class: 'bg-text-muted/20 text-text-muted cursor-not-allowed' }
    case 'locked': return { label: 'Locked', class: 'bg-text-muted/20 text-text-muted cursor-not-allowed' }
    case 'expired': return { label: 'Expired', class: 'bg-text-muted/20 text-text-muted cursor-not-allowed' }
    case 'cooldown': return { label: 'Cooldown', class: 'bg-text-muted/20 text-text-muted cursor-not-allowed' }
    default: return { label: 'Start', class: 'bg-c-accent text-black' }
  }
})

const isDisabled = computed(() =>
  ['claimed', 'locked', 'expired', 'cooldown'].includes(props.task.status)
)
</script>

<template>
  <div class="bg-card-bg border border-border rounded-xl p-4 flex items-center gap-4">
    <!-- Icon -->
    <div
      class="w-10 h-10 rounded-[10px] flex items-center justify-center flex-shrink-0"
      :style="{ backgroundColor: task.iconColor + '20' }"
    >
      <span class="material-symbols-rounded text-xl" :style="{ color: task.iconColor }">{{ task.icon }}</span>
    </div>

    <!-- Content -->
    <div class="flex-1 min-w-0">
      <div class="text-sm font-semibold text-text-primary truncate">{{ task.name }}</div>
      <div class="text-xs text-text-secondary">{{ task.type }} · {{ task.completions.toLocaleString() }} completed</div>
    </div>

    <!-- Points + Action -->
    <div class="flex items-center gap-3 flex-shrink-0">
      <span class="text-xs font-semibold text-c-accent">+{{ task.points }} XP</span>
      <button
        class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
        :class="buttonConfig.class"
        :disabled="isDisabled"
        @click="!isDisabled && emit('action', task.id)"
      >
        {{ buttonConfig.label }}
      </button>
    </div>
  </div>
</template>
