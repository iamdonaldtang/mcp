<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const error = ref('')

async function handleLogin() {
  error.value = ''
  try {
    await auth.login(email.value, password.value)
    router.push('/b/dashboard')
  } catch (e: any) {
    error.value = e.response?.data?.message || 'Login failed'
  }
}
</script>

<template>
  <div class="min-h-screen bg-page-bg flex items-center justify-center">
    <div class="w-full max-w-md">
      <div class="bg-card-bg border border-border rounded-2xl p-8">
        <div class="flex items-center justify-center mb-8">
          <span class="material-symbols-rounded text-community text-4xl mr-2">rocket_launch</span>
          <span class="text-2xl font-bold text-text-primary">TaskOn</span>
        </div>
        <h2 class="text-xl font-semibold text-text-primary text-center mb-6">Sign in to your account</h2>

        <div v-if="error" class="mb-4 p-3 bg-status-paused-bg border border-status-paused rounded-lg text-sm text-status-paused">
          {{ error }}
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm text-text-secondary mb-1">Email</label>
            <input v-model="email" type="email" required
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none" 
              placeholder="you@project.com" />
          </div>
          <div>
            <label class="block text-sm text-text-secondary mb-1">Password</label>
            <input v-model="password" type="password" required
              class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none"
              placeholder="••••••••" />
          </div>
          <button type="submit" :disabled="auth.loading"
            class="w-full py-2.5 bg-community text-white font-semibold rounded-lg hover:bg-community/90 transition-colors disabled:opacity-50">
            {{ auth.loading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
