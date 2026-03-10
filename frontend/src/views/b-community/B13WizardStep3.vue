<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()
const loading = ref(true)
const saving = ref(false)

// ── Draft data from previous steps ──
interface DraftData {
  name: string
  description: string
  strategy: string
  enabled_modules: string[]
  content?: {
    sectors?: SectorDraft[]
    levels?: LevelDraft[]
    taskchain?: TaskChainDraft
    daychain?: DayChainDraft
  }
}

interface SectorDraft {
  id: string
  name: string
  tasks: TaskDraft[]
}

interface TaskDraft {
  id: string
  name: string
  type: string
  xp: number
  editing: boolean
}

interface LevelDraft {
  id: string
  name: string
  threshold: number
  editing: boolean
}

interface TaskChainDraft {
  name: string
  steps: ChainStepDraft[]
}

interface ChainStepDraft {
  id: string
  name: string
  xp: number
  editing: boolean
}

interface DayChainDraft {
  base_reward: number
  editing: boolean
}

const draft = ref<DraftData>({
  name: '',
  description: '',
  strategy: '',
  enabled_modules: [],
  content: {},
})

// ── Module content state ──
const sectors = ref<SectorDraft[]>([])
const levels = ref<LevelDraft[]>([])
const taskchain = ref<TaskChainDraft>({ name: 'Getting Started Chain', steps: [] })
const daychain = ref<DayChainDraft>({ base_reward: 10, editing: false })

const MAX_TASKS_PER_SECTOR = 10

// ── Checklist ──
const checklistItems = computed(() => [
  { label: 'Community created', done: true, icon: 'check_circle' },
  {
    label: `Modules selected (${draft.value.enabled_modules.length} modules)`,
    done: draft.value.enabled_modules.length > 0,
    icon: 'check_circle',
  },
  {
    label: 'Content configured',
    done: isContentConfigured.value,
    icon: isContentConfigured.value ? 'check_circle' : 'radio_button_unchecked',
  },
  { label: 'Preview & Publish', done: false, icon: 'radio_button_unchecked' },
])

const isContentConfigured = computed(() => {
  const modules = draft.value.enabled_modules
  let configured = true
  if (modules.includes('sectors_tasks') && sectors.value.length > 0) {
    configured = configured && sectors.value.some(s => s.tasks.length > 0)
  }
  if (modules.includes('points_level') && levels.value.length > 0) {
    configured = configured && levels.value.length >= 2
  }
  return configured
})

// ── Helpers ──
let idCounter = 0
function genId() {
  return `tmp_${++idCounter}_${Date.now()}`
}

function hasModule(key: string) {
  return draft.value.enabled_modules.includes(key)
}

// ── Default content generators ──
function generateDefaultSectors() {
  sectors.value = [
    {
      id: genId(),
      name: 'General',
      tasks: [
        { id: genId(), name: 'Follow us on Twitter', type: 'social', xp: 20, editing: false },
        { id: genId(), name: 'Join Discord Server', type: 'social', xp: 15, editing: false },
        { id: genId(), name: 'Complete Profile', type: 'verification', xp: 25, editing: false },
      ],
    },
  ]
}

function generateDefaultLevels() {
  levels.value = [
    { id: genId(), name: 'Bronze', threshold: 0, editing: false },
    { id: genId(), name: 'Silver', threshold: 100, editing: false },
    { id: genId(), name: 'Gold', threshold: 500, editing: false },
    { id: genId(), name: 'Platinum', threshold: 2000, editing: false },
    { id: genId(), name: 'Diamond', threshold: 5000, editing: false },
  ]
}

function generateDefaultTaskChain() {
  taskchain.value = {
    name: 'Getting Started Chain',
    steps: [
      { id: genId(), name: 'Connect Wallet', xp: 10, editing: false },
      { id: genId(), name: 'Complete First Task', xp: 20, editing: false },
      { id: genId(), name: 'Earn 50 XP', xp: 30, editing: false },
    ],
  }
}

function generateDefaultDayChain() {
  daychain.value = { base_reward: 10, editing: false }
}

// ── Task CRUD ──
function addTask(sector: SectorDraft) {
  if (sector.tasks.length >= MAX_TASKS_PER_SECTOR) return
  sector.tasks.push({
    id: genId(),
    name: 'New Task',
    type: 'custom',
    xp: 10,
    editing: true,
  })
}

function removeTask(sector: SectorDraft, taskId: string) {
  sector.tasks = sector.tasks.filter(t => t.id !== taskId)
}

// ── Level CRUD ──
function addLevel() {
  const lastThreshold = levels.value.length > 0
    ? levels.value[levels.value.length - 1].threshold
    : 0
  levels.value.push({
    id: genId(),
    name: 'New Level',
    threshold: lastThreshold + 500,
    editing: true,
  })
}

function removeLevel(id: string) {
  if (levels.value.length <= 2) return
  levels.value = levels.value.filter(l => l.id !== id)
}

// ── Chain step CRUD ──
function addChainStep() {
  if (taskchain.value.steps.length >= 10) return
  taskchain.value.steps.push({
    id: genId(),
    name: 'New Step',
    xp: 15,
    editing: true,
  })
}

function removeChainStep(id: string) {
  taskchain.value.steps = taskchain.value.steps.filter(s => s.id !== id)
}

// ── Click to edit ──
function startEdit(item: { editing: boolean }) {
  item.editing = true
}

function finishEdit(item: { editing: boolean }) {
  item.editing = false
}

// ── API ──
onMounted(async () => {
  try {
    const res = await api.get('/api/v1/community/wizard/draft')
    const data = res.data.data as DraftData
    draft.value = data

    // Restore or generate default content per enabled module
    if (hasModule('sectors_tasks')) {
      sectors.value = data.content?.sectors?.length
        ? data.content.sectors
        : (generateDefaultSectors(), sectors.value)
    }
    if (hasModule('points_level')) {
      levels.value = data.content?.levels?.length
        ? data.content.levels
        : (generateDefaultLevels(), levels.value)
    }
    if (hasModule('taskchain')) {
      taskchain.value = data.content?.taskchain?.steps?.length
        ? data.content.taskchain
        : (generateDefaultTaskChain(), taskchain.value)
    }
    if (hasModule('daychain')) {
      daychain.value = data.content?.daychain
        ? data.content.daychain
        : (generateDefaultDayChain(), daychain.value)
    }
  } catch {
    // If no draft, generate defaults for all common modules
    draft.value.enabled_modules = ['sectors_tasks', 'points_level']
    generateDefaultSectors()
    generateDefaultLevels()
  } finally {
    loading.value = false
  }
})

async function saveDraft() {
  saving.value = true
  try {
    await api.put('/api/v1/community/wizard/draft', {
      ...draft.value,
      content: {
        sectors: sectors.value,
        levels: levels.value,
        taskchain: taskchain.value,
        daychain: daychain.value,
      },
    })
  } catch {
    // TODO: toast notification
  } finally {
    saving.value = false
  }
}

async function goBack() {
  await saveDraft()
  router.push({ name: 'B13-2' })
}

async function goNext() {
  await saveDraft()
  router.push({ name: 'B13-4' })
}

// ── Stepper ──
const steps = [
  { num: 1, label: 'Customize' },
  { num: 2, label: 'Modules' },
  { num: 3, label: 'Quick Setup' },
  { num: 4, label: 'Preview & Publish' },
]

const taskTypeIcons: Record<string, string> = {
  social: 'share',
  onchain: 'account_balance_wallet',
  verification: 'verified_user',
  custom: 'extension',
  recurring: 'repeat',
  referral: 'group_add',
}
</script>

<template>
  <div class="min-h-screen bg-page-bg">
    <!-- Top Bar -->
    <div class="sticky top-0 z-30 bg-header-bg border-b border-border px-6 py-3 flex items-center justify-between">
      <button
        class="flex items-center gap-1.5 text-sm text-text-secondary hover:text-text-primary transition-colors"
        @click="goBack"
      >
        <span class="material-symbols-rounded text-lg">arrow_back</span>
        Back
      </button>
      <h1 class="text-base font-semibold text-text-primary">Create Community</h1>
      <button
        class="px-4 py-1.5 text-sm text-text-secondary border border-border rounded-lg hover:text-text-primary hover:border-text-muted transition-colors disabled:opacity-50"
        :disabled="saving"
        @click="saveDraft"
      >
        {{ saving ? 'Saving...' : 'Save Draft' }}
      </button>
    </div>

    <!-- Stepper -->
    <div class="flex items-center justify-center gap-0 py-6 px-6 bg-page-bg">
      <template v-for="(step, idx) in steps" :key="step.num">
        <div class="flex items-center gap-2">
          <div
            class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-semibold transition-colors"
            :class="
              step.num < 3
                ? 'bg-community text-white'
                : step.num === 3
                  ? 'bg-community text-white ring-2 ring-community/30'
                  : 'bg-card-bg border border-border text-text-muted'
            "
          >
            <span v-if="step.num < 3" class="material-symbols-rounded text-sm">check</span>
            <span v-else>{{ step.num }}</span>
          </div>
          <span
            class="text-sm font-medium"
            :class="step.num === 3 ? 'text-text-primary' : step.num < 3 ? 'text-community' : 'text-text-muted'"
          >
            {{ step.label }}
          </span>
        </div>
        <div
          v-if="idx < steps.length - 1"
          class="w-16 h-px mx-3"
          :class="step.num < 3 ? 'bg-community' : 'bg-border'"
        />
      </template>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="flex items-center gap-3 text-text-muted">
        <span class="material-symbols-rounded animate-spin text-2xl">progress_activity</span>
        Loading draft...
      </div>
    </div>

    <!-- Main Content -->
    <div v-else class="max-w-6xl mx-auto px-6 pb-24 flex gap-6">
      <!-- Left: Content Editor -->
      <div class="flex-1 space-y-6">
        <!-- Sectors & Tasks -->
        <div v-if="hasModule('sectors_tasks')" class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <div class="px-5 py-4 border-b border-border flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="material-symbols-rounded text-lg text-community">folder</span>
              <h2 class="text-sm font-semibold text-text-primary">Sectors & Tasks</h2>
            </div>
            <span class="text-xs text-text-muted">Auto-generated starter content — edit inline</span>
          </div>

          <div v-for="sector in sectors" :key="sector.id" class="border-b border-border last:border-b-0">
            <!-- Sector header -->
            <div class="px-5 py-3 bg-page-bg/50 flex items-center justify-between">
              <span class="text-xs font-semibold text-text-secondary uppercase tracking-wider">{{ sector.name }}</span>
              <span class="text-xs text-text-muted">{{ sector.tasks.length }} / {{ MAX_TASKS_PER_SECTOR }} tasks</span>
            </div>

            <!-- Task rows -->
            <div v-for="task in sector.tasks" :key="task.id" class="px-5 py-3 flex items-center gap-3 hover:bg-white/2 transition-colors group">
              <div class="w-7 h-7 rounded-md flex items-center justify-center bg-community/10">
                <span class="material-symbols-rounded text-sm text-community">{{ taskTypeIcons[task.type] || 'task_alt' }}</span>
              </div>

              <!-- Name (click to edit) -->
              <div class="flex-1 min-w-0">
                <input
                  v-if="task.editing"
                  v-model="task.name"
                  type="text"
                  class="w-full bg-page-bg border border-community/40 rounded px-2 py-1 text-sm text-text-primary focus:outline-none focus:border-community"
                  @blur="finishEdit(task)"
                  @keyup.enter="finishEdit(task)"
                />
                <span
                  v-else
                  class="text-sm text-text-primary cursor-pointer hover:text-community transition-colors"
                  @click="startEdit(task)"
                >
                  {{ task.name }}
                  <span class="material-symbols-rounded text-xs text-text-muted ml-1 opacity-0 group-hover:opacity-100 transition-opacity">edit</span>
                </span>
              </div>

              <!-- Type badge -->
              <span class="px-2 py-0.5 text-[10px] rounded bg-page-bg text-text-muted uppercase tracking-wider">{{ task.type }}</span>

              <!-- XP (click to edit) -->
              <div class="w-20 text-right">
                <input
                  v-if="task.editing"
                  v-model.number="task.xp"
                  type="number"
                  min="1"
                  class="w-16 bg-page-bg border border-community/40 rounded px-2 py-1 text-sm text-c-accent text-right focus:outline-none focus:border-community"
                  @blur="finishEdit(task)"
                  @keyup.enter="finishEdit(task)"
                />
                <span
                  v-else
                  class="text-sm text-c-accent font-medium cursor-pointer hover:underline"
                  @click="startEdit(task)"
                >
                  +{{ task.xp }} XP
                </span>
              </div>

              <!-- Remove -->
              <button
                class="p-1 rounded hover:bg-status-paused-bg transition-colors opacity-0 group-hover:opacity-100"
                title="Remove task"
                @click="removeTask(sector, task.id)"
              >
                <span class="material-symbols-rounded text-sm text-status-paused">close</span>
              </button>
            </div>

            <!-- Add task -->
            <div class="px-5 py-2.5 border-t border-border/50">
              <button
                class="flex items-center gap-1.5 text-xs text-community hover:text-community/80 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                :disabled="sector.tasks.length >= MAX_TASKS_PER_SECTOR"
                @click="addTask(sector)"
              >
                <span class="material-symbols-rounded text-sm">add</span>
                Add Task
              </button>
            </div>
          </div>
        </div>

        <!-- Points & Levels -->
        <div v-if="hasModule('points_level')" class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <div class="px-5 py-4 border-b border-border flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="material-symbols-rounded text-lg text-c-accent">military_tech</span>
              <h2 class="text-sm font-semibold text-text-primary">Points & Levels</h2>
            </div>
            <span class="text-xs text-text-muted">{{ levels.length }} levels</span>
          </div>

          <div class="divide-y divide-border">
            <div
              v-for="(level, idx) in levels"
              :key="level.id"
              class="px-5 py-3 flex items-center gap-4 hover:bg-white/2 transition-colors group"
            >
              <!-- Level number -->
              <div class="w-7 h-7 rounded-full flex items-center justify-center bg-c-accent/10">
                <span class="text-xs font-bold text-c-accent">{{ idx + 1 }}</span>
              </div>

              <!-- Name (click to edit) -->
              <div class="w-32">
                <input
                  v-if="level.editing"
                  v-model="level.name"
                  type="text"
                  class="w-full bg-page-bg border border-community/40 rounded px-2 py-1 text-sm text-text-primary focus:outline-none focus:border-community"
                  @blur="finishEdit(level)"
                  @keyup.enter="finishEdit(level)"
                />
                <span
                  v-else
                  class="text-sm font-medium text-text-primary cursor-pointer hover:text-community transition-colors"
                  @click="startEdit(level)"
                >
                  {{ level.name }}
                  <span class="material-symbols-rounded text-xs text-text-muted ml-1 opacity-0 group-hover:opacity-100 transition-opacity">edit</span>
                </span>
              </div>

              <!-- Threshold (click to edit) -->
              <div class="flex-1">
                <input
                  v-if="level.editing"
                  v-model.number="level.threshold"
                  type="number"
                  min="0"
                  class="w-28 bg-page-bg border border-community/40 rounded px-2 py-1 text-sm text-text-secondary text-right focus:outline-none focus:border-community"
                  @blur="finishEdit(level)"
                  @keyup.enter="finishEdit(level)"
                />
                <span
                  v-else
                  class="text-sm text-text-secondary cursor-pointer hover:underline"
                  @click="startEdit(level)"
                >
                  {{ level.threshold.toLocaleString() }} XP required
                </span>
              </div>

              <!-- Progress bar visual -->
              <div class="w-24 h-1.5 bg-border rounded-full overflow-hidden">
                <div
                  class="h-full bg-c-accent rounded-full transition-all"
                  :style="{ width: `${Math.min(100, (level.threshold / 5000) * 100)}%` }"
                />
              </div>

              <!-- Remove -->
              <button
                v-if="levels.length > 2"
                class="p-1 rounded hover:bg-status-paused-bg transition-colors opacity-0 group-hover:opacity-100"
                title="Remove level"
                @click="removeLevel(level.id)"
              >
                <span class="material-symbols-rounded text-sm text-status-paused">close</span>
              </button>
            </div>
          </div>

          <div class="px-5 py-2.5 border-t border-border">
            <button
              class="flex items-center gap-1.5 text-xs text-community hover:text-community/80 transition-colors"
              @click="addLevel"
            >
              <span class="material-symbols-rounded text-sm">add</span>
              Add Level
            </button>
          </div>
        </div>

        <!-- TaskChain -->
        <div v-if="hasModule('taskchain')" class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <div class="px-5 py-4 border-b border-border flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="material-symbols-rounded text-lg text-quest">route</span>
              <h2 class="text-sm font-semibold text-text-primary">TaskChain</h2>
            </div>
            <span class="text-xs text-text-muted">{{ taskchain.steps.length }} steps</span>
          </div>

          <!-- Chain visualization -->
          <div class="px-5 py-4 space-y-0">
            <div
              v-for="(step, idx) in taskchain.steps"
              :key="step.id"
              class="flex items-start gap-3 group"
            >
              <!-- Connector -->
              <div class="flex flex-col items-center">
                <div class="w-6 h-6 rounded-full bg-quest/10 flex items-center justify-center">
                  <span class="text-[10px] font-bold text-quest">{{ idx + 1 }}</span>
                </div>
                <div v-if="idx < taskchain.steps.length - 1" class="w-px h-8 bg-border" />
              </div>

              <!-- Step content -->
              <div class="flex-1 flex items-center gap-3 pb-2">
                <div class="flex-1">
                  <input
                    v-if="step.editing"
                    v-model="step.name"
                    type="text"
                    class="w-full bg-page-bg border border-community/40 rounded px-2 py-1 text-sm text-text-primary focus:outline-none focus:border-community"
                    @blur="finishEdit(step)"
                    @keyup.enter="finishEdit(step)"
                  />
                  <span
                    v-else
                    class="text-sm text-text-primary cursor-pointer hover:text-community transition-colors"
                    @click="startEdit(step)"
                  >
                    {{ step.name }}
                    <span class="material-symbols-rounded text-xs text-text-muted ml-1 opacity-0 group-hover:opacity-100 transition-opacity">edit</span>
                  </span>
                </div>
                <div class="w-20 text-right">
                  <input
                    v-if="step.editing"
                    v-model.number="step.xp"
                    type="number"
                    min="1"
                    class="w-16 bg-page-bg border border-community/40 rounded px-2 py-1 text-sm text-c-accent text-right focus:outline-none focus:border-community"
                    @blur="finishEdit(step)"
                    @keyup.enter="finishEdit(step)"
                  />
                  <span
                    v-else
                    class="text-sm text-c-accent font-medium cursor-pointer hover:underline"
                    @click="startEdit(step)"
                  >
                    +{{ step.xp }} XP
                  </span>
                </div>
                <button
                  class="p-1 rounded hover:bg-status-paused-bg transition-colors opacity-0 group-hover:opacity-100"
                  title="Remove step"
                  @click="removeChainStep(step.id)"
                >
                  <span class="material-symbols-rounded text-sm text-status-paused">close</span>
                </button>
              </div>
            </div>
          </div>

          <div class="px-5 py-2.5 border-t border-border">
            <button
              class="flex items-center gap-1.5 text-xs text-community hover:text-community/80 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
              :disabled="taskchain.steps.length >= 10"
              @click="addChainStep"
            >
              <span class="material-symbols-rounded text-sm">add</span>
              Add Step
            </button>
          </div>
        </div>

        <!-- DayChain -->
        <div v-if="hasModule('daychain')" class="bg-card-bg border border-border rounded-xl overflow-hidden">
          <div class="px-5 py-4 border-b border-border flex items-center gap-2">
            <span class="material-symbols-rounded text-lg text-boost">local_fire_department</span>
            <h2 class="text-sm font-semibold text-text-primary">DayChain — Daily Check-in</h2>
          </div>

          <div class="px-5 py-4">
            <div class="flex items-center gap-4">
              <span class="text-sm text-text-secondary">Base daily reward:</span>
              <div class="flex items-center gap-2">
                <input
                  v-if="daychain.editing"
                  v-model.number="daychain.base_reward"
                  type="number"
                  min="1"
                  class="w-20 bg-page-bg border border-community/40 rounded px-2 py-1.5 text-sm text-c-accent text-center focus:outline-none focus:border-community"
                  @blur="finishEdit(daychain)"
                  @keyup.enter="finishEdit(daychain)"
                />
                <span
                  v-else
                  class="text-sm text-c-accent font-semibold cursor-pointer hover:underline"
                  @click="startEdit(daychain)"
                >
                  {{ daychain.base_reward }} XP
                </span>
                <span class="text-xs text-text-muted">(multiplied by streak day)</span>
              </div>
            </div>

            <!-- Streak preview -->
            <div class="mt-4 flex items-center gap-2">
              <div
                v-for="day in 7"
                :key="day"
                class="flex flex-col items-center gap-1"
              >
                <div
                  class="w-9 h-9 rounded-lg flex items-center justify-center text-xs font-semibold"
                  :class="day <= 3 ? 'bg-community/20 text-community' : 'bg-border text-text-muted'"
                >
                  {{ day }}x
                </div>
                <span class="text-[10px] text-text-muted">+{{ daychain.base_reward * day }}</span>
              </div>
              <div class="flex flex-col items-center gap-1">
                <div class="w-9 h-9 rounded-lg bg-border flex items-center justify-center">
                  <span class="material-symbols-rounded text-sm text-text-muted">more_horiz</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty state: no modules -->
        <div
          v-if="!hasModule('sectors_tasks') && !hasModule('points_level') && !hasModule('taskchain') && !hasModule('daychain')"
          class="bg-card-bg border border-border rounded-xl p-12 text-center"
        >
          <span class="material-symbols-rounded text-4xl text-text-muted block mb-3">widgets</span>
          <p class="text-sm text-text-muted mb-2">No configurable modules selected</p>
          <button
            class="text-xs text-community hover:underline"
            @click="router.push({ name: 'B13-2' })"
          >
            Go back to select modules
          </button>
        </div>
      </div>

      <!-- Right: Launch Checklist -->
      <div class="w-72 shrink-0">
        <div class="bg-card-bg border border-border rounded-xl overflow-hidden sticky top-20">
          <div class="px-5 py-4 border-b border-border">
            <h3 class="text-sm font-semibold text-text-primary">Launch Checklist</h3>
            <p class="text-xs text-text-muted mt-0.5">
              {{ checklistItems.filter(i => i.done).length }} of {{ checklistItems.length }} complete
            </p>
          </div>

          <div class="p-4 space-y-3">
            <div
              v-for="item in checklistItems"
              :key="item.label"
              class="flex items-start gap-2.5"
            >
              <span
                class="material-symbols-rounded text-lg mt-px"
                :class="item.done ? 'text-community' : 'text-text-muted'"
              >
                {{ item.icon }}
              </span>
              <span
                class="text-sm leading-5"
                :class="item.done ? 'text-text-primary' : 'text-text-muted'"
              >
                {{ item.label }}
              </span>
            </div>
          </div>

          <!-- Progress bar -->
          <div class="px-4 pb-4">
            <div class="w-full h-1.5 bg-border rounded-full overflow-hidden">
              <div
                class="h-full bg-community rounded-full transition-all duration-500"
                :style="{ width: `${(checklistItems.filter(i => i.done).length / checklistItems.length) * 100}%` }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Action Bar -->
    <div v-if="!loading" class="fixed bottom-0 left-0 right-0 bg-header-bg border-t border-border px-6 py-3 flex items-center justify-between z-20">
      <button
        class="px-5 py-2 text-sm text-text-secondary border border-border rounded-lg hover:text-text-primary hover:border-text-muted transition-colors"
        @click="goBack"
      >
        Back
      </button>
      <button
        class="px-5 py-2 bg-community text-white text-sm font-semibold rounded-lg hover:bg-community/90 transition-colors flex items-center gap-1.5"
        @click="goNext"
      >
        Next: Preview & Publish
        <span class="material-symbols-rounded text-base">arrow_forward</span>
      </button>
    </div>
  </div>
</template>
