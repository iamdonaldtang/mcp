<script setup lang="ts">
defineProps<{
  icon: string
  title: string
  description: string
  actionLabel?: string
  actionRoute?: string
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-16 text-center">
    <span class="material-symbols-rounded text-5xl text-text-muted mb-4">{{ icon }}</span>
    <h3 class="text-lg font-semibold text-text-primary mb-2">{{ title }}</h3>
    <p class="text-sm text-text-secondary max-w-md mb-6">{{ description }}</p>
    <router-link
      v-if="actionLabel && actionRoute"
      :to="actionRoute"
      class="px-6 py-2.5 bg-community text-white font-medium rounded-lg hover:bg-community/90 transition-colors"
    >
      {{ actionLabel }}
    </router-link>
    <slot />
  </div>
</template>
