<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { cApi } from '../../api/client'
import type { CQuest } from '../../types/c-end'

const router = useRouter()
const loading = ref(true)
const quests = ref<CQuest[]>([])
const filter = ref<'all' | 'available' | 'completed'>('all')

const filteredQuests = computed(() => {
  if (filter.value === 'all') return quests.value
  if (filter.value === 'available') return quests.value.filter(q => q.status === 'available' || q.status === 'in_progress')
  return quests.value.filter(q => q.status === 'completed')
})

const counts = computed(() => ({
  all: quests.value.length,
  available: quests.value.filter(q => q.status === 'available' || q.status === 'in_progress').length,
  completed: quests.value.filter(q => q.status === 'completed').length,
}))

onMounted(async () => {
  try {
    const res = await cApi.get('/api/c/quests')
    quests.value = res.data.data || []
  } finally { loading.value = false }
})

function timeRemaining(endsAt?: string) {
  if (!endsAt) return null
  const diff = new Date(endsAt).getTime() - Date.now()
  if (diff <= 0) return null
  const d = Math.floor(diff / 86400000)
  return `Ends in ${d}d`
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-text-primary">All Quests</h1>
    </div>

    <!-- Filter Pills -->
    <div class="flex gap-2">
      <button v-for="f in (['all', 'available', 'completed'] as const)" :key="f"
        class="px-4 py-1.5 text-sm font-medium rounded-lg transition-colors"
        :class="filter === f ? 'bg-c-accent text-black' : 'bg-card-bg border border-border text-text-secondary'"
        @click="filter = f">
        {{ f === 'all' ? 'All' : f.charAt(0).toUpperCase() + f.slice(1) }} ({{ counts[f] }})
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-2 gap-4">
      <div v-for="i in 4" :key="i" class="bg-card-bg border border-border rounded-xl h-52 animate-pulse"></div>
    </div>

    <!-- Quest Grid -->
    <div v-else class="grid grid-cols-2 gap-4">
      <div v-for="quest in filteredQuests" :key="quest.id"
        class="bg-card-bg border border-border rounded-xl overflow-hidden hover:border-c-accent/30 transition-colors cursor-pointer"
        @click="router.push(`/c/quests/${quest.id}`)">
        <!-- Banner -->
        <div class="h-28 bg-page-bg flex items-center justify-center relative">
          <img v-if="quest.image" :src="quest.image" :alt="quest.name" class="w-full h-full object-cover" />
          <span v-else class="material-symbols-rounded text-3xl text-text-muted">image</span>
          <!-- Status badge -->
          <span v-if="quest.status === 'completed'" class="absolute top-2 right-2 px-2 py-0.5 text-[10px] font-semibold rounded-full bg-status-active-bg text-status-active">Completed</span>
          <span v-if="timeRemaining(quest.endsAt)" class="absolute top-2 right-2 px-2 py-0.5 text-[10px] font-semibold rounded-full bg-c-accent/20 text-c-accent">{{ timeRemaining(quest.endsAt) }}</span>
        </div>
        <!-- Content -->
        <div class="p-4">
          <h3 class="text-sm font-semibold text-text-primary mb-1">{{ quest.name }}</h3>
          <p class="text-xs text-text-muted mb-3 line-clamp-2">{{ quest.description }}</p>
          <div class="flex items-center justify-between">
            <span class="text-xs text-text-secondary">{{ quest.taskCount }} tasks · {{ quest.participants.toLocaleString() }} participants</span>
            <button v-if="quest.status !== 'completed'"
              class="px-3 py-1 text-xs font-medium bg-c-accent text-black rounded-lg hover:bg-c-accent/90">
              Start Quest
            </button>
            <span v-else class="text-xs text-status-active font-medium">All completed ✓</span>
          </div>
          <!-- Progress if in_progress -->
          <div v-if="quest.status === 'in_progress' && quest.taskCount > 0" class="mt-2">
            <div class="h-1.5 bg-page-bg rounded-full overflow-hidden">
              <div class="h-full bg-c-accent rounded-full" :style="{ width: (quest.completedTasks / quest.taskCount * 100) + '%' }"></div>
            </div>
            <span class="text-[10px] text-text-muted mt-0.5">{{ quest.completedTasks }}/{{ quest.taskCount }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Cross-sell -->
    <div class="grid grid-cols-2 gap-4 pt-4">
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/lb-sprint')">
        <span class="text-sm text-text-primary font-medium">Boost your rank in the Sprint →</span>
      </button>
      <button class="bg-card-bg border border-border rounded-xl p-4 text-left hover:border-c-accent/30 transition-colors" @click="router.push('/c/milestones')">
        <span class="text-sm text-text-primary font-medium">Unlock Milestones with your points →</span>
      </button>
    </div>
  </div>
</template>
