<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api } from '../../api/client'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

interface SectorItem {
  id: string
  name: string
  description: string
  sort_order: number
  status: string
  tasks?: TaskItem[]
}

interface TaskItem {
  id: string
  sector_id: string
  name: string
  type: string
  status: CampaignStatus
  points: number
  icon: string
  icon_color: string
  current_completions: number
  created_at: string
}

const loading = ref(true)
const sectors = ref<SectorItem[]>([])
const tasks = ref<TaskItem[]>([])
const totalTasks = ref(0)
const page = ref(1)
const pageSize = 20

// Filters
const filterStatus = ref<string>('all')
const filterSector = ref<string>('all')
const searchQuery = ref('')

// Modal states
const showCreateSector = ref(false)
const showCreateTask = ref(false)
const newSectorName = ref('')
const newSectorDesc = ref('')

const newTask = ref({
  sector_id: '',
  name: '',
  description: '',
  type: 'social',
  points: 10,
  icon: 'task_alt',
  icon_color: '#F59E0B',
})

const taskTypes = ['social', 'onchain', 'verification', 'custom', 'recurring', 'referral']

const filteredTasks = computed(() => {
  return tasks.value.filter(t => {
    if (filterStatus.value !== 'all' && t.status !== filterStatus.value) return false
    if (filterSector.value !== 'all' && t.sector_id !== filterSector.value) return false
    if (searchQuery.value && !t.name.toLowerCase().includes(searchQuery.value.toLowerCase())) return false
    return true
  })
})

onMounted(async () => {
  await Promise.all([fetchSectors(), fetchTasks()])
  loading.value = false
})

async function fetchSectors() {
  try {
    const res = await api.get('/api/v1/community/sectors')
    sectors.value = res.data.data || []
  } catch { /* empty */ }
}

async function fetchTasks() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterSector.value !== 'all') params.sector_id = filterSector.value
    if (filterStatus.value !== 'all') params.status = filterStatus.value

    const res = await api.get('/api/v1/community/tasks', { params })
    tasks.value = res.data.data?.items || []
    totalTasks.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

async function createSector() {
  if (!newSectorName.value.trim()) return
  try {
    await api.post('/api/v1/community/sectors', {
      name: newSectorName.value,
      description: newSectorDesc.value,
    })
    showCreateSector.value = false
    newSectorName.value = ''
    newSectorDesc.value = ''
    await fetchSectors()
  } catch { /* TODO: toast */ }
}

async function createTask() {
  if (!newTask.value.name.trim() || !newTask.value.sector_id) return
  try {
    await api.post('/api/v1/community/tasks', newTask.value)
    showCreateTask.value = false
    newTask.value = { sector_id: '', name: '', description: '', type: 'social', points: 10, icon: 'task_alt', icon_color: '#F59E0B' }
    await fetchTasks()
  } catch { /* TODO: toast */ }
}

async function toggleTaskStatus(task: TaskItem) {
  const newStatus = task.status === 'active' ? 'draft' : 'active'
  const oldStatus = task.status
  task.status = newStatus as CampaignStatus // optimistic
  try {
    await api.put(`/api/v1/community/tasks/${task.id}`, { status: newStatus })
  } catch {
    task.status = oldStatus as CampaignStatus // rollback
  }
}

async function deleteTask(id: string) {
  if (!confirm('Delete this task? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/tasks/${id}`)
    await fetchTasks()
  } catch { /* TODO: toast */ }
}

function onPageChange(p: number) {
  page.value = p
  fetchTasks()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Sectors & Tasks</h1>
        <p class="text-sm text-text-secondary">Manage task sectors and individual tasks for your community</p>
      </div>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 border border-community text-community text-sm font-medium rounded-lg hover:bg-community/10 transition-colors"
          @click="showCreateSector = true"
        >
          + New Sector
        </button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
          @click="showCreateTask = true"
        >
          + New Task
        </button>
      </div>
    </div>

    <!-- Sectors overview -->
    <div class="grid grid-cols-4 gap-4">
      <div
        v-for="sector in sectors"
        :key="sector.id"
        class="bg-card-bg border border-border rounded-xl p-4 cursor-pointer hover:border-community/30 transition-colors"
        :class="filterSector === sector.id ? 'border-community' : ''"
        @click="filterSector = filterSector === sector.id ? 'all' : sector.id; fetchTasks()"
      >
        <div class="flex items-center justify-between mb-1">
          <span class="text-sm font-semibold text-text-primary">{{ sector.name }}</span>
          <span class="text-xs text-text-muted">{{ sector.tasks?.length || 0 }} tasks</span>
        </div>
        <p class="text-xs text-text-muted truncate">{{ sector.description || 'No description' }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="flex gap-2">
        <button
          v-for="status in ['all', 'active', 'draft', 'paused']"
          :key="status"
          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
          :class="filterStatus === status ? 'bg-community text-white' : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
          @click="filterStatus = status; fetchTasks()"
        >
          {{ status === 'all' ? 'All' : status.charAt(0).toUpperCase() + status.slice(1) }}
        </button>
      </div>
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search tasks..."
            class="w-full pl-10 pr-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-community focus:outline-none"
          />
        </div>
      </div>
    </div>

    <!-- Tasks Table -->
    <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-border">
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Task</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Type</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Points</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Status</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Completions</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-24">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="6" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse"></div></td>
          </tr>
          <tr v-else-if="filteredTasks.length === 0">
            <td colspan="6" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">inbox</span>
              <p class="text-sm text-text-muted">No tasks yet</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="showCreateTask = true">+ Create First Task</button>
            </td>
          </tr>
          <tr
            v-else
            v-for="task in filteredTasks"
            :key="task.id"
            class="border-b border-border last:border-b-0 hover:bg-white/2 transition-colors"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{ background: task.icon_color + '20' }">
                  <span class="material-symbols-rounded text-base" :style="{ color: task.icon_color }">{{ task.icon }}</span>
                </div>
                <div>
                  <div class="text-sm font-medium text-text-primary">{{ task.name }}</div>
                  <div class="text-xs text-text-muted">{{ sectors.find(s => s.id === task.sector_id)?.name || '—' }}</div>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <span class="px-2 py-0.5 text-xs rounded bg-page-bg text-text-secondary capitalize">{{ task.type }}</span>
            </td>
            <td class="px-4 py-3 text-sm text-c-accent font-medium">+{{ task.points }}</td>
            <td class="px-4 py-3">
              <StatusBadge :status="task.status" />
            </td>
            <td class="px-4 py-3 text-sm text-text-secondary">{{ task.current_completions.toLocaleString() }}</td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button
                  class="p-1.5 rounded-lg hover:bg-white/5 transition-colors"
                  :title="task.status === 'active' ? 'Pause' : 'Activate'"
                  @click="toggleTaskStatus(task)"
                >
                  <span class="material-symbols-rounded text-base text-text-muted">
                    {{ task.status === 'active' ? 'pause_circle' : 'play_circle' }}
                  </span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit">
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteTask(task.id)">
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination v-if="totalTasks > pageSize" :page="page" :page-size="pageSize" :total="totalTasks" @update:page="onPageChange" />
    </div>

    <!-- Create Sector Modal -->
    <Modal :open="showCreateSector" title="Create Sector" @close="showCreateSector = false">
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-text-secondary mb-1">Sector Name</label>
          <input v-model="newSectorName" type="text" placeholder="e.g. Getting Started" class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm" />
        </div>
        <div>
          <label class="block text-sm text-text-secondary mb-1">Description</label>
          <textarea v-model="newSectorDesc" rows="3" placeholder="Brief description..." class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm resize-none" />
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showCreateSector = false">Cancel</button>
        <button class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50" :disabled="!newSectorName.trim()" @click="createSector">Create Sector</button>
      </template>
    </Modal>

    <!-- Create Task Modal -->
    <Modal :open="showCreateTask" title="Create Task" max-width="560px" @close="showCreateTask = false">
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-text-secondary mb-1">Sector</label>
          <select v-model="newTask.sector_id" class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none">
            <option value="" disabled>Select a sector</option>
            <option v-for="s in sectors" :key="s.id" :value="s.id">{{ s.name }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm text-text-secondary mb-1">Task Name</label>
          <input v-model="newTask.name" type="text" placeholder="e.g. Follow on Twitter" class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm" />
        </div>
        <div>
          <label class="block text-sm text-text-secondary mb-1">Description</label>
          <textarea v-model="newTask.description" rows="2" placeholder="What should the user do?" class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm resize-none" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-text-secondary mb-1">Type</label>
            <select v-model="newTask.type" class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none">
              <option v-for="t in taskTypes" :key="t" :value="t">{{ t.charAt(0).toUpperCase() + t.slice(1) }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1">Points Reward</label>
            <input v-model.number="newTask.points" type="number" min="0" class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none" />
          </div>
        </div>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showCreateTask = false">Cancel</button>
        <button class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50" :disabled="!newTask.name.trim() || !newTask.sector_id" @click="createTask">Create Task</button>
      </template>
    </Modal>
  </div>
</template>
