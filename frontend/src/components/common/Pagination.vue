<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{
  (e: 'update:page', value: number): void
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

function goTo(p: number) {
  if (p >= 1 && p <= totalPages.value) {
    emit('update:page', p)
  }
}
</script>

<template>
  <div class="flex items-center justify-between px-4 py-3">
    <span class="text-sm text-text-muted">
      Showing {{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, total) }} of {{ total }}
    </span>
    <div class="flex items-center gap-1">
      <button
        class="px-3 py-1.5 text-sm rounded-lg transition-colors"
        :class="page <= 1 ? 'text-text-muted cursor-not-allowed' : 'text-text-secondary hover:bg-white/5'"
        :disabled="page <= 1"
        @click="goTo(page - 1)"
      >
        Prev
      </button>
      <template v-for="p in totalPages" :key="p">
        <button
          v-if="p <= 5 || p === totalPages || Math.abs(p - page) <= 1"
          class="w-8 h-8 text-sm rounded-lg transition-colors"
          :class="p === page ? 'bg-community text-white' : 'text-text-secondary hover:bg-white/5'"
          @click="goTo(p)"
        >
          {{ p }}
        </button>
        <span v-else-if="p === 6 || p === totalPages - 1" class="text-text-muted px-1">...</span>
      </template>
      <button
        class="px-3 py-1.5 text-sm rounded-lg transition-colors"
        :class="page >= totalPages ? 'text-text-muted cursor-not-allowed' : 'text-text-secondary hover:bg-white/5'"
        :disabled="page >= totalPages"
        @click="goTo(page + 1)"
      >
        Next
      </button>
    </div>
  </div>
</template>
