<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api/client'

const router = useRouter()
const loading = ref(true)
const hasCommunity = ref(false)
const selectedPath = ref<string | null>(null)

interface DeploymentPath {
  id: 'domain' | 'embed' | 'sdk'
  title: string
  icon: string
  description: string
  tags: string[]
  codeRequired: string
  recommended?: boolean
}

const deploymentPaths: DeploymentPath[] = [
  {
    id: 'domain',
    title: 'Host on Your Domain',
    icon: 'language',
    description: 'Zero code, custom domain. Your community at yourproject.com',
    tags: ['No Code', 'Custom Domain'],
    codeRequired: 'None',
  },
  {
    id: 'embed',
    title: 'Embed in Your App',
    icon: 'widgets',
    description: 'Widget Library + Page Builder + SSO. Integrate into your existing app',
    tags: ['Low Code', 'Flexible'],
    codeRequired: 'Low',
    recommended: true,
  },
  {
    id: 'sdk',
    title: 'Build with SDK',
    icon: 'code',
    description: 'Full programmatic control. Build custom experiences with our API',
    tags: ['Full Control', 'Developer'],
    codeRequired: 'Full',
  },
]

const highlights = [
  { icon: 'palette', label: 'Custom Branding', detail: 'Your logo, colors, and domain — users never see TaskOn' },
  { icon: 'lock', label: 'Data Ownership', detail: 'Full access to user data, analytics, and engagement metrics' },
  { icon: 'speed', label: '< 7 Days Integration', detail: 'From sign-up to live deployment in under a week' },
]

const resources = [
  { icon: 'integration_instructions', title: 'Integration Guide', description: 'Step-by-step setup for all deployment paths', href: '#' },
  { icon: 'play_circle', title: 'Watch Demo', description: 'See White Label in action — 3 minute walkthrough', href: '#' },
  { icon: 'api', title: 'API Docs', description: 'Full SDK and API reference documentation', href: '#', external: true },
]

onMounted(async () => {
  try {
    const res = await api.get('/api/v1/community/overview')
    hasCommunity.value = !!res.data?.data?.status && res.data.data.status !== 'none'
  } catch {
    hasCommunity.value = false
  } finally {
    loading.value = false
  }
})

function selectPath(path: DeploymentPath) {
  selectedPath.value = path.id
  router.push({ path: '/b/whitelabel/wizard/step1', query: { path: path.id } })
}
</script>

<template>
  <div class="max-w-5xl mx-auto space-y-10">
    <!-- Hero Section -->
    <div>
      <div class="w-16 h-16 rounded-2xl flex items-center justify-center mb-4" style="background: #1A1033">
        <span class="material-symbols-rounded text-3xl" style="color: #9B7EE0">auto_awesome</span>
      </div>
      <h1 class="text-2xl font-bold text-text-primary mb-2">Own Your Growth Experience</h1>
      <p class="text-base text-text-secondary max-w-2xl">
        Launch a fully branded community platform under your domain. Your users interact with your brand — TaskOn powers everything behind the scenes.
      </p>
    </div>

    <!-- Prerequisite Banner -->
    <div
      v-if="!loading && !hasCommunity"
      class="flex items-center gap-3 p-4 rounded-xl border"
      style="background: rgba(245, 158, 11, 0.08); border-color: #D97706"
    >
      <span class="material-symbols-rounded text-xl" style="color: #F59E0B">warning</span>
      <div class="flex-1">
        <span class="text-sm font-medium text-text-primary">White Label requires an active Community</span>
        <span class="text-sm text-text-secondary ml-1">— set up your community first to unlock White Label features.</span>
      </div>
      <router-link
        to="/b/community"
        class="text-sm font-medium hover:underline shrink-0"
        style="color: #F59E0B"
      >
        Set Up Community &rarr;
      </router-link>
    </div>

    <!-- Deployment Path Cards -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Choose Your Deployment Path</div>
      <div class="grid grid-cols-3 gap-6">
        <button
          v-for="path in deploymentPaths"
          :key="path.id"
          class="text-left bg-card-bg border rounded-xl p-5 transition-all duration-200 relative group"
          :class="[
            selectedPath === path.id
              ? 'border-2 shadow-[0_0_0_1px_#9B7EE0]'
              : 'border-border hover:bg-white/2',
          ]"
          :style="selectedPath === path.id ? { borderColor: '#9B7EE0' } : {}"
          :disabled="!hasCommunity && !loading"
          @click="selectPath(path)"
        >
          <!-- Recommended badge -->
          <div
            v-if="path.recommended"
            class="absolute -top-2.5 left-4 px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider rounded-full"
            style="background: #9B7EE0; color: #0A0F1A"
          >
            &#9733; Recommended
          </div>

          <div class="w-10 h-10 rounded-lg flex items-center justify-center mb-3" style="background: #1A1033">
            <span class="material-symbols-rounded text-xl" style="color: #9B7EE0">{{ path.icon }}</span>
          </div>
          <h3 class="text-base font-semibold text-text-primary mb-1">{{ path.title }}</h3>
          <p class="text-sm text-text-secondary mb-3 leading-relaxed">{{ path.description }}</p>

          <!-- Tags -->
          <div class="flex flex-wrap gap-1.5 mb-3">
            <span
              v-for="tag in path.tags"
              :key="tag"
              class="px-2 py-0.5 text-xs rounded-full"
              style="background: #1A1033; color: #9B7EE0"
            >{{ tag }}</span>
          </div>

          <!-- Code Required -->
          <div class="text-xs text-text-muted">
            Code Required: <span class="text-text-secondary font-medium">{{ path.codeRequired }}</span>
          </div>

          <!-- CTA -->
          <div class="mt-4 pt-3 border-t border-border">
            <span class="text-sm font-medium group-hover:underline" style="color: #9B7EE0">
              Get Started &rarr;
            </span>
          </div>
        </button>
      </div>
    </div>

    <!-- Highlight Strip -->
    <div class="bg-card-bg border border-border rounded-xl p-5">
      <div class="grid grid-cols-3 gap-6">
        <div v-for="h in highlights" :key="h.label" class="flex items-start gap-3">
          <span class="material-symbols-rounded text-lg mt-0.5 shrink-0" style="color: #9B7EE0">{{ h.icon }}</span>
          <div>
            <div class="text-sm font-medium text-text-primary">{{ h.label }}</div>
            <div class="text-xs text-text-secondary mt-0.5">{{ h.detail }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Divider -->
    <div class="border-t border-border"></div>

    <!-- Resources -->
    <div>
      <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-4">Resources</div>
      <div class="grid grid-cols-3 gap-6">
        <a
          v-for="r in resources"
          :key="r.title"
          :href="r.href"
          :target="r.external ? '_blank' : undefined"
          class="bg-card-bg border border-border rounded-xl p-5 hover:bg-white/2 transition-colors block"
        >
          <span class="material-symbols-rounded text-2xl mb-3 block" style="color: #9B7EE0">{{ r.icon }}</span>
          <h4 class="text-sm font-semibold text-text-primary mb-1">{{ r.title }}</h4>
          <p class="text-xs text-text-secondary">{{ r.description }}</p>
        </a>
      </div>
    </div>
  </div>
</template>
