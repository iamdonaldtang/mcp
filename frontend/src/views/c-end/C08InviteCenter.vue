<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { cApi } from '../../api/client'
import type { InviteData, Referral } from '../../types/c-end'

const loading = ref(true)
const inviteData = ref<InviteData | null>(null)
const referrals = ref<Referral[]>([])

onMounted(async () => {
  try {
    const res = await cApi.get('/api/c/invite/link')
    inviteData.value = res.data.data
    referrals.value = inviteData.value?.referrals || []
  } finally { loading.value = false }
})

function copyLink() {
  if (inviteData.value?.referralUrl) {
    navigator.clipboard.writeText(inviteData.value.referralUrl)
  }
}

function shareTwitter() {
  if (!inviteData.value) return
  const text = encodeURIComponent(`Join our community! Complete tasks and earn rewards. ${inviteData.value.referralUrl}`)
  window.open(`https://twitter.com/intent/tweet?text=${text}`, '_blank')
}

function shareTelegram() {
  if (!inviteData.value) return
  window.open(`https://t.me/share/url?url=${encodeURIComponent(inviteData.value.referralUrl)}&text=${encodeURIComponent('Join our community!')}`, '_blank')
}

function truncAddr(addr: string) { return addr ? addr.slice(0, 6) + '...' + addr.slice(-4) : '—' }
</script>

<template>
  <div class="space-y-8">
    <!-- Hero -->
    <div class="bg-card-bg border border-border rounded-2xl p-8 text-center">
      <div class="w-16 h-16 rounded-2xl bg-c-accent/20 flex items-center justify-center mx-auto mb-4">
        <span class="material-symbols-rounded text-c-accent text-3xl">share</span>
      </div>
      <h1 class="text-2xl font-bold text-text-primary mb-2">Invite Friends, Earn Rewards</h1>
      <p class="text-sm text-text-secondary">Share your link and earn 30 points for every friend who signs up</p>
    </div>

    <!-- Referral Link -->
    <div class="bg-card-bg border border-border rounded-xl p-5">
      <div class="flex items-center gap-3 mb-4">
        <input :value="inviteData?.referralUrl || ''" readonly
          class="flex-1 px-4 py-2.5 bg-page-bg border border-border rounded-lg text-sm text-text-secondary" />
        <button class="px-5 py-2.5 bg-c-accent text-black text-sm font-medium rounded-lg hover:bg-c-accent/90" @click="copyLink">Copy Link</button>
      </div>
      <!-- Social share -->
      <div class="flex justify-center gap-3">
        <button class="px-5 py-2 bg-[#1DA1F2]/20 text-[#1DA1F2] text-sm font-medium rounded-lg hover:bg-[#1DA1F2]/30" @click="shareTwitter">Twitter</button>
        <button class="px-5 py-2 bg-[#5865F2]/20 text-[#5865F2] text-sm font-medium rounded-lg hover:bg-[#5865F2]/30">Discord</button>
        <button class="px-5 py-2 bg-[#0088CC]/20 text-[#0088CC] text-sm font-medium rounded-lg hover:bg-[#0088CC]/30" @click="shareTelegram">Telegram</button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-text-primary">{{ inviteData?.totalInvites || 0 }}</div>
        <div class="text-xs text-text-muted">Invites Sent</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-text-primary">{{ inviteData?.successfulJoins || 0 }}</div>
        <div class="text-xs text-text-muted">Successful Joins</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-c-accent">{{ inviteData?.pointsEarned || 0 }}</div>
        <div class="text-xs text-text-muted">Points Earned</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4 text-center">
        <div class="text-xl font-bold text-text-primary">{{ inviteData?.conversionRate || 0 }}%</div>
        <div class="text-xs text-text-muted">Conversion Rate</div>
      </div>
    </div>

    <!-- Referrals Table -->
    <div>
      <h2 class="text-base font-semibold text-text-primary mb-3">Your Referrals</h2>
      <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
        <table class="w-full">
          <thead>
            <tr class="border-b border-border">
              <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase">Wallet</th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase w-28">Status</th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase w-28">Points</th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase w-32">Date</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="referrals.length === 0">
              <td colspan="4" class="px-4 py-8 text-center text-sm text-text-muted">No referrals yet</td>
            </tr>
            <tr v-for="ref in referrals" :key="ref.address" class="border-b border-border last:border-b-0">
              <td class="px-4 py-3 text-sm text-text-primary">{{ truncAddr(ref.address) }}</td>
              <td class="px-4 py-3">
                <span class="px-2 py-0.5 text-xs rounded-full font-medium"
                  :class="ref.status === 'joined' ? 'bg-status-active-bg text-status-active' : 'bg-status-draft-bg text-status-draft'">
                  {{ ref.status === 'joined' ? 'Joined' : 'Pending' }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-c-accent">{{ ref.pointsEarned > 0 ? '+' + ref.pointsEarned : '—' }}</td>
              <td class="px-4 py-3 text-xs text-text-muted">{{ new Date(ref.date).toLocaleDateString() }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
