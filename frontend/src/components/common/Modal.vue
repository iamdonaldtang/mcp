<script setup lang="ts">
defineProps<{
  open: boolean
  title: string
  maxWidth?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/60" @click="emit('close')" />
      <!-- Content -->
      <div
        class="relative bg-card-bg border border-border rounded-2xl shadow-xl"
        :style="{ maxWidth: maxWidth || '480px', width: '100%' }"
      >
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-4 border-b border-border">
          <h2 class="text-lg font-semibold text-text-primary">{{ title }}</h2>
          <button class="text-text-muted hover:text-text-primary transition-colors" @click="emit('close')">
            <span class="material-symbols-rounded text-xl">close</span>
          </button>
        </div>
        <!-- Body -->
        <div class="px-6 py-4">
          <slot />
        </div>
        <!-- Footer -->
        <div v-if="$slots.footer" class="px-6 py-4 border-t border-border flex justify-end gap-3">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </Teleport>
</template>
