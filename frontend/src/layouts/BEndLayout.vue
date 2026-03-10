<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const sidebarCollapsed = ref(false)
const product = computed(() => route.meta.product as string | undefined)

// Sidebar navigation structure
interface NavItem {
  id: string
  label: string
  icon: string
  route?: string
  children?: NavItem[]
  section?: string // non-clickable section header
  product?: 'community' | 'whitelabel'
}

const navItems: NavItem[] = [
  { id: 'dashboard', label: 'Home', icon: 'home', route: '/b/dashboard' },
  {
    id: 'community', label: 'Community', icon: 'groups', product: 'community',
    children: [
      { id: 'c-overview', label: 'Overview', icon: 'dashboard', route: '/b/community/overview' },
      // Module sections
      { id: 'c-sect-tasks', label: 'TASKS', icon: '', section: 'header' },
      { id: 'c-sectors', label: 'Sectors & Tasks', icon: 'category', route: '/b/community/modules/sectors' },
      { id: 'c-taskchain', label: 'TaskChain', icon: 'link', route: '/b/community/modules/task-chain' },
      { id: 'c-daychain', label: 'DayChain', icon: 'local_fire_department', route: '/b/community/modules/day-chain' },
      { id: 'c-sect-points', label: 'POINTS', icon: '', section: 'header' },
      { id: 'c-points', label: 'Points & Level', icon: 'stars', route: '/b/community/modules/points-level' },
      { id: 'c-leaderboard', label: 'Leaderboard', icon: 'leaderboard', route: '/b/community/modules/leaderboard' },
      { id: 'c-badges', label: 'Badges', icon: 'military_tech', route: '/b/community/modules/badges' },
      { id: 'c-sect-campaigns', label: 'CAMPAIGNS', icon: '', section: 'header' },
      { id: 'c-lbsprint', label: 'LB Sprint', icon: 'sprint', route: '/b/community/modules/lb-sprint' },
      { id: 'c-milestone', label: 'Milestone', icon: 'flag', route: '/b/community/modules/milestone' },
      { id: 'c-luckywheel', label: 'Lucky Wheel', icon: 'casino', route: '/b/community/modules/lucky-wheel' },
      { id: 'c-sect-rewards', label: 'REWARDS', icon: '', section: 'header' },
      { id: 'c-shop', label: 'Benefits Shop', icon: 'storefront', route: '/b/community/modules/benefits-shop' },
      { id: 'c-sect-settings', label: 'SETTINGS', icon: '', section: 'header' },
      { id: 'c-access', label: 'Access Rules', icon: 'lock', route: '/b/community/modules/access-rules' },
      { id: 'c-homepage', label: 'Homepage Editor', icon: 'web', route: '/b/community/modules/homepage-editor' },
    ],
  },
  {
    id: 'whitelabel', label: 'White Label', icon: 'palette', product: 'whitelabel',
    children: [
      { id: 'wl-overview', label: 'Overview', icon: 'dashboard', route: '/b/whitelabel/overview' },
      { id: 'wl-widgets', label: 'Widgets', icon: 'widgets', route: '/b/whitelabel/widgets' },
      { id: 'wl-pages', label: 'Pages', icon: 'article', route: '/b/whitelabel/pages' },
      { id: 'wl-smartrewards', label: 'Smart Rewards', icon: 'auto_awesome',
        children: [
          { id: 'wl-rules', label: 'Rule Builder', icon: 'rule', route: '/b/whitelabel/smart-rewards/rules' },
          { id: 'wl-privileges', label: 'Privilege Manager', icon: 'workspace_premium', route: '/b/whitelabel/smart-rewards/privileges' },
        ],
      },
    ],
  },
  { id: 'analytics', label: 'Analytics', icon: 'analytics', route: '/b/analytics' },
  { id: 'settings', label: 'Settings', icon: 'settings', route: '/b/settings' },
]

const expandedItems = ref<Set<string>>(new Set(['community', 'whitelabel']))

function toggleExpand(id: string) {
  if (expandedItems.value.has(id)) {
    expandedItems.value.delete(id)
  } else {
    expandedItems.value.add(id)
  }
}

function isActive(routePath?: string) {
  return routePath ? route.path === routePath : false
}

function navigate(routePath?: string) {
  if (routePath) router.push(routePath)
}

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="flex h-screen bg-page-bg">
    <!-- Sidebar -->
    <aside
      class="flex flex-col border-r border-border bg-sidebar-bg transition-all duration-200"
      :class="sidebarCollapsed ? 'w-16' : 'w-60'"
    >
      <!-- Logo -->
      <div class="flex items-center h-14 px-4 border-b border-border">
        <span class="material-symbols-rounded text-community text-2xl">rocket_launch</span>
        <span v-if="!sidebarCollapsed" class="ml-2 text-text-primary font-semibold text-base">TaskOn</span>
      </div>

      <!-- Nav -->
      <nav class="flex-1 overflow-y-auto py-2">
        <template v-for="item in navItems" :key="item.id">
          <!-- Top-level with children -->
          <div v-if="item.children">
            <button
              class="flex items-center w-full px-4 py-2.5 text-sm transition-colors hover:bg-white/5"
              :class="isActive(item.route) ? 'text-text-primary bg-white/10' : 'text-text-secondary'"
              @click="toggleExpand(item.id)"
            >
              <span class="material-symbols-rounded text-xl">{{ item.icon }}</span>
              <span v-if="!sidebarCollapsed" class="ml-3 flex-1 text-left">{{ item.label }}</span>
              <span v-if="!sidebarCollapsed" class="material-symbols-rounded text-base transition-transform"
                :class="expandedItems.has(item.id) ? '' : '-rotate-180'"
              >keyboard_arrow_up</span>
            </button>
            <!-- Sub items -->
            <div v-if="expandedItems.has(item.id) && !sidebarCollapsed" class="ml-0">
              <template v-for="sub in item.children" :key="sub.id">
                <!-- Section header -->
                <div v-if="sub.section === 'header'"
                  class="px-10 pt-4 pb-1 text-[10px] font-semibold tracking-wider text-text-muted uppercase">
                  {{ sub.label }}
                </div>
                <!-- Sub-nav item -->
                <button v-else
                  class="flex items-center w-full py-2 pl-10 pr-3 text-[13px] transition-colors hover:bg-white/5"
                  :class="isActive(sub.route)
                    ? (product === 'community' ? 'bg-[#ECFDF5] text-community font-semibold' : 'bg-[#1A1033] text-whitelabel font-semibold')
                    : 'text-text-secondary'"
                  @click="navigate(sub.route)"
                >
                  <span v-if="sub.icon" class="material-symbols-rounded text-base mr-2">{{ sub.icon }}</span>
                  <span>{{ sub.label }}</span>
                </button>
              </template>
            </div>
          </div>

          <!-- Top-level without children -->
          <button v-else
            class="flex items-center w-full px-4 py-2.5 text-sm transition-colors hover:bg-white/5"
            :class="isActive(item.route) ? 'text-text-primary bg-white/10' : 'text-text-secondary'"
            @click="navigate(item.route)"
          >
            <span class="material-symbols-rounded text-xl">{{ item.icon }}</span>
            <span v-if="!sidebarCollapsed" class="ml-3">{{ item.label }}</span>
          </button>
        </template>
      </nav>

      <!-- Collapse toggle -->
      <button
        class="flex items-center justify-center h-10 border-t border-border text-text-muted hover:text-text-primary transition-colors"
        @click="sidebarCollapsed = !sidebarCollapsed"
      >
        <span class="material-symbols-rounded text-xl">{{ sidebarCollapsed ? 'chevron_right' : 'chevron_left' }}</span>
      </button>
    </aside>

    <!-- Main content area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Top bar -->
      <header class="flex items-center h-14 px-6 border-b border-border bg-header-bg">
        <!-- Breadcrumb -->
        <div class="flex-1 text-sm text-text-secondary">
          <span class="text-text-muted">TaskOn</span>
          <span class="mx-2 text-text-muted">/</span>
          <span class="text-text-primary">{{ route.meta.title || route.name }}</span>
        </div>

        <!-- Right actions -->
        <div class="flex items-center gap-4">
          <button class="text-text-secondary hover:text-text-primary transition-colors">
            <span class="material-symbols-rounded text-xl">help_outline</span>
          </button>
          <button class="text-text-secondary hover:text-text-primary transition-colors">
            <span class="material-symbols-rounded text-xl">notifications</span>
          </button>
          <div class="flex items-center gap-2 cursor-pointer" @click="logout">
            <div class="w-8 h-8 rounded-full bg-community/20 flex items-center justify-center">
              <span class="material-symbols-rounded text-community text-sm">person</span>
            </div>
            <span v-if="auth.user" class="text-sm text-text-primary">{{ auth.user.email }}</span>
          </div>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>
