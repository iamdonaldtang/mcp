import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

// === B-End Community Routes ===
const bCommunityRoutes: RouteRecordRaw[] = [
  {
    path: '/b/community',
    component: () => import('../layouts/BEndLayout.vue'),
    meta: { product: 'community' },
    children: [
      // Hub states
      { path: '', redirect: '/b/community/overview' },
      { path: 'empty', name: 'B09', component: () => import('../views/b-community/B09CommunityEmpty.vue') },
      { path: 'guided', name: 'B10', component: () => import('../views/b-community/B10CommunityGuided.vue') },
      { path: 'overview', name: 'B11', component: () => import('../views/b-community/B11CommunityActive.vue') },
      { path: 'management', name: 'B12', component: () => import('../views/b-community/B12CommunityManagement.vue') },

      // Wizard
      { path: 'wizard/step-1', name: 'B13-1', component: () => import('../views/b-community/B13WizardStep1.vue') },
      { path: 'wizard/step-2', name: 'B13-2', component: () => import('../views/b-community/B13WizardStep2.vue') },
      { path: 'wizard/step-3', name: 'B13-3', component: () => import('../views/b-community/B13WizardStep3.vue') },
      { path: 'wizard/step-4', name: 'B13-4', component: () => import('../views/b-community/B13WizardStep4.vue') },

      // Module Management pages
      { path: 'modules/sectors', name: 'B31a', component: () => import('../views/b-community/B31aSectors.vue') },
      { path: 'modules/task-chain', name: 'B31b', component: () => import('../views/b-community/B31bTaskChain.vue') },
      { path: 'modules/day-chain', name: 'B31c', component: () => import('../views/b-community/B31cDayChain.vue') },
      { path: 'modules/leaderboard', name: 'B31d', component: () => import('../views/b-community/B31dLeaderboard.vue') },
      { path: 'modules/lb-sprint', name: 'B31e', component: () => import('../views/b-community/B31eLBSprint.vue') },
      { path: 'modules/milestone', name: 'B31f', component: () => import('../views/b-community/B31fMilestone.vue') },
      { path: 'modules/benefits-shop', name: 'B31g', component: () => import('../views/b-community/B31gBenefitsShop.vue') },
      { path: 'modules/lucky-wheel', name: 'B31h', component: () => import('../views/b-community/B31hLuckyWheel.vue') },
      { path: 'modules/badges', name: 'B31i', component: () => import('../views/b-community/B31iBadges.vue') },
      { path: 'modules/points-level', name: 'B32', component: () => import('../views/b-community/B32PointsLevel.vue') },
      { path: 'modules/access-rules', name: 'B33', component: () => import('../views/b-community/B33AccessRules.vue') },
      { path: 'modules/homepage-editor', name: 'B34', component: () => import('../views/b-community/B34HomepageEditor.vue') },

      // Operations
      { path: 'content', name: 'B50', component: () => import('../views/b-community/B50Content.vue') },
      { path: 'preview', name: 'B51', component: () => import('../views/b-community/B51Preview.vue') },

      // Smart Rewards (shared with WL)
      { path: 'smart-rewards', name: 'B52', component: () => import('../views/b-community/B52SmartRewards.vue') },

      // Insights
      { path: 'insights', name: 'B60', component: () => import('../views/b-community/B60Insights.vue') },

      // Integration Center
      { path: 'integration', name: 'B61', component: () => import('../views/b-community/B61Integration.vue') },
    ],
  },
]

// === B-End White Label Routes ===
const bWhiteLabelRoutes: RouteRecordRaw[] = [
  {
    path: '/b/whitelabel',
    component: () => import('../layouts/BEndLayout.vue'),
    meta: { product: 'whitelabel' },
    children: [
      { path: '', redirect: '/b/whitelabel/overview' },
      { path: 'empty', name: 'B14', component: () => import('../views/b-whitelabel/B14WLEmpty.vue') },
      { path: 'active', name: 'B15', component: () => import('../views/b-whitelabel/B15WLActive.vue') },
      { path: 'overview', name: 'B16', component: () => import('../views/b-whitelabel/B16WLManagement.vue') },

      // Wizard
      { path: 'wizard/step-1', name: 'B17-1', component: () => import('../views/b-whitelabel/B17WizardStep1.vue') },
      { path: 'wizard/step-1-5', name: 'B17-1.5', component: () => import('../views/b-whitelabel/B17WizardStep1_5.vue') },
      { path: 'wizard/step-2', name: 'B17-2', component: () => import('../views/b-whitelabel/B17WizardStep2.vue') },
      { path: 'wizard/step-3', name: 'B17-3', component: () => import('../views/b-whitelabel/B17WizardStep3.vue') },
      { path: 'wizard/step-4', name: 'B17-4', component: () => import('../views/b-whitelabel/B17WizardStep4.vue') },

      // Domain
      { path: 'domain', name: 'B40', component: () => import('../views/b-whitelabel/B40DomainSetup.vue') },

      // Widget Library
      { path: 'widgets', name: 'B41', component: () => import('../views/b-whitelabel/B41WidgetLibrary.vue') },
      { path: 'widgets/:id/config', name: 'B41-config', component: () => import('../views/b-whitelabel/B41WidgetConfig.vue') },

      // Page Builder
      { path: 'pages', name: 'B42', component: () => import('../views/b-whitelabel/B42PageBuilder.vue') },
      { path: 'pages/:id/editor', name: 'B42-editor', component: () => import('../views/b-whitelabel/B42PageEditor.vue') },

      // SDK / Integration
      { path: 'sdk', name: 'B43', component: () => import('../views/b-whitelabel/B43SDKConfig.vue') },
      { path: 'integration', name: 'B44', component: () => import('../views/b-whitelabel/B44IntegrationConfig.vue') },

      // Brand
      { path: 'brand', name: 'B45', component: () => import('../views/b-whitelabel/B45BrandConfig.vue') },

      // Smart Rewards
      { path: 'smart-rewards/rules', name: 'B46', component: () => import('../views/b-whitelabel/B46RuleBuilder.vue') },
      { path: 'smart-rewards/privileges', name: 'B47', component: () => import('../views/b-whitelabel/B47PrivilegeManager.vue') },

      // Contract Registry
      { path: 'contracts', name: 'B48', component: () => import('../views/b-whitelabel/B48ContractRegistry.vue') },

      // Dev Kit
      { path: 'devkit', name: 'B49', component: () => import('../views/b-whitelabel/B49DevKit.vue') },
    ],
  },
]

// === B-End Shared Routes (Dashboard, Analytics, Settings) ===
const bSharedRoutes: RouteRecordRaw[] = [
  {
    path: '/b',
    component: () => import('../layouts/BEndLayout.vue'),
    children: [
      { path: '', redirect: '/b/dashboard' },
      { path: 'dashboard', name: 'B01', component: () => import('../views/b-community/B01Dashboard.vue') },
      { path: 'analytics', name: 'B05', component: () => import('../views/b-community/B05Analytics.vue') },
      { path: 'settings', name: 'B06', component: () => import('../views/b-community/B06Settings.vue') },
    ],
  },
]

// === C-End Routes ===
const cEndRoutes: RouteRecordRaw[] = [
  {
    path: '/c',
    component: () => import('../layouts/CEndLayout.vue'),
    children: [
      { path: '', name: 'C01', component: () => import('../views/c-end/C01Home.vue') },
      { path: 'quests', name: 'C02', component: () => import('../views/c-end/C02Quests.vue') },
      { path: 'leaderboard', name: 'C03', component: () => import('../views/c-end/C03Leaderboard.vue') },
      { path: 'lb-sprint', name: 'C04', component: () => import('../views/c-end/C04LBSprint.vue') },
      { path: 'milestones', name: 'C05', component: () => import('../views/c-end/C05Milestone.vue') },
      { path: 'shop', name: 'C06', component: () => import('../views/c-end/C06Shop.vue') },
      { path: 'profile', name: 'C07', component: () => import('../views/c-end/C07UserCenter.vue') },
      { path: 'invite', name: 'C08', component: () => import('../views/c-end/C08InviteCenter.vue') },
      { path: 'activity', name: 'C09', component: () => import('../views/c-end/C09ActivityFeed.vue') },
    ],
  },
]

// === Auth Routes ===
const authRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/LoginPage.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    ...authRoutes,
    ...bSharedRoutes,
    ...bCommunityRoutes,
    ...bWhiteLabelRoutes,
    ...cEndRoutes,
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('../views/NotFound.vue') },
  ],
})

// Navigation guard for B-end auth
router.beforeEach((to) => {
  const token = localStorage.getItem('b_token')
  if (to.path.startsWith('/b') && !token && to.name !== 'Login') {
    return { name: 'Login' }
  }
})

export default router
