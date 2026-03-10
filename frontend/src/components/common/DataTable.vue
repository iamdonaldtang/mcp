<script setup lang="ts">
defineProps<{
  columns: { key: string; label: string; width?: string }[]
  loading?: boolean
}>()
</script>

<template>
  <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
    <table class="w-full">
      <thead>
        <tr class="border-b border-border">
          <th
            v-for="col in columns"
            :key="col.key"
            class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider"
            :style="col.width ? { width: col.width } : {}"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td :colspan="columns.length" class="px-4 py-12 text-center text-text-muted">
            <span class="material-symbols-rounded animate-spin text-2xl">progress_activity</span>
          </td>
        </tr>
        <slot v-else />
      </tbody>
    </table>
  </div>
</template>
