<script setup lang="ts">
defineProps<{
  label: string
  value: string | number
  icon: string
  iconColor?: string
  trend?: string
}>()
</script>

<template>
  <div class="bg-card-bg border border-border rounded-xl p-4">
    <div class="flex items-center justify-between mb-2">
      <span class="text-xs text-text-muted uppercase tracking-wider">{{ label }}</span>
      <span class="material-symbols-rounded text-lg" :style="{ color: iconColor || '#94A3B8' }">{{ icon }}</span>
    </div>
    <div class="text-2xl font-bold text-text-primary">{{ typeof value === 'number' ? value.toLocaleString() : value }}</div>
    <div v-if="trend" class="mt-1 text-xs text-status-active">{{ trend }}</div>
  </div>
</template>
