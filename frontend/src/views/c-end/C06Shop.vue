<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { cApi } from '../../api/client'
import { useCEndStore } from '../../stores/c-end'
import Modal from '../../components/common/Modal.vue'
import type { ShopItem } from '../../types/c-end'

const store = useCEndStore()
const loading = ref(true)
const items = ref<ShopItem[]>([])
const filterCategory = ref('all')
const categories = ['all', 'nft', 'voucher', 'merch', 'whitelist']

const userPoints = computed(() => store.userStatus?.xp || 0)

const filteredItems = computed(() =>
  filterCategory.value === 'all'
    ? items.value
    : items.value.filter(i => i.category === filterCategory.value)
)

// Redeem modal
const showRedeem = ref(false)
const redeemItem = ref<ShopItem | null>(null)
const redeeming = ref(false)

onMounted(async () => {
  try {
    const res = await cApi.get('/api/c/shop/items')
    items.value = res.data.data || []
  } finally {
    loading.value = false
  }
})

function getItemState(item: ShopItem) {
  if (item.stock === 0) return 'sold_out'
  if (userPoints.value < item.price) return 'not_enough'
  return 'affordable'
}

function getButtonConfig(item: ShopItem) {
  const state = getItemState(item)
  switch (state) {
    case 'affordable': return { label: 'Redeem', class: 'bg-c-accent text-black hover:bg-c-accent/90', disabled: false }
    case 'not_enough': return { label: 'Not Enough', class: 'bg-text-muted/20 text-text-muted', disabled: true }
    case 'sold_out': return { label: 'Sold Out', class: 'bg-text-muted/20 text-text-muted line-through', disabled: true }
    default: return { label: 'Redeem', class: 'bg-c-accent text-black', disabled: false }
  }
}

function openRedeem(item: ShopItem) {
  if (getItemState(item) !== 'affordable') return
  redeemItem.value = item
  showRedeem.value = true
}

async function confirmRedeem() {
  if (!redeemItem.value) return
  redeeming.value = true
  try {
    await cApi.post('/api/c/shop/redeem', { item_id: redeemItem.value.id })
    // Update local state
    const idx = items.value.findIndex(i => i.id === redeemItem.value!.id)
    if (idx !== -1 && items.value[idx].stock > 0) {
      items.value[idx].stock--
      items.value[idx].totalRedemptions++
      if (items.value[idx].stock === 0) items.value[idx].status = 'sold_out'
    }
    showRedeem.value = false
    redeemItem.value = null
    // Refresh user status
    await store.fetchUserStatus()
  } catch {
    // TODO: toast error
  } finally {
    redeeming.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-text-primary">Rewards Shop</h1>
      <div class="px-4 py-2 bg-c-accent/20 rounded-xl">
        <span class="text-sm font-semibold text-c-accent">{{ userPoints.toLocaleString() }} Points Available</span>
      </div>
    </div>

    <!-- Category Filter -->
    <div class="flex gap-2">
      <button
        v-for="cat in categories"
        :key="cat"
        class="px-4 py-1.5 text-sm font-medium rounded-lg transition-colors capitalize"
        :class="filterCategory === cat
          ? 'bg-c-accent text-black'
          : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
        @click="filterCategory = cat"
      >
        {{ cat === 'all' ? 'All Items' : cat + 's' }}
      </button>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="grid grid-cols-3 gap-4">
      <div v-for="i in 6" :key="i" class="bg-card-bg border border-border rounded-xl h-64 animate-pulse"></div>
    </div>

    <!-- Items Grid -->
    <div v-else class="grid grid-cols-3 gap-4">
      <div
        v-for="item in filteredItems"
        :key="item.id"
        class="bg-card-bg border border-border rounded-xl overflow-hidden hover:border-c-accent/30 transition-colors"
      >
        <!-- Image -->
        <div class="h-36 bg-page-bg flex items-center justify-center">
          <img v-if="item.image" :src="item.image" :alt="item.name" class="h-full w-full object-cover" />
          <span v-else class="material-symbols-rounded text-4xl text-text-muted">image</span>
        </div>

        <!-- Info -->
        <div class="p-4">
          <div class="flex items-center gap-2 mb-1">
            <span class="px-2 py-0.5 text-[10px] uppercase font-semibold rounded text-text-muted bg-page-bg">{{ item.category }}</span>
            <span v-if="item.isTimeLimited" class="text-[10px] text-c-accent">⏰ Limited</span>
          </div>
          <h3 class="text-sm font-semibold text-text-primary mb-1">{{ item.name }}</h3>
          <div class="flex items-center justify-between mb-3">
            <span class="text-sm font-bold text-c-accent">{{ item.price.toLocaleString() }} pts</span>
            <span v-if="item.stock > 0 && item.stock <= 50" class="text-xs text-text-muted">{{ item.stock }} remaining</span>
          </div>

          <!-- Not enough hint -->
          <p v-if="getItemState(item) === 'not_enough'" class="text-xs text-text-muted mb-2">
            Need {{ (item.price - userPoints).toLocaleString() }} more
          </p>

          <button
            class="w-full py-2 text-sm font-medium rounded-lg transition-colors"
            :class="getButtonConfig(item).class"
            :disabled="getButtonConfig(item).disabled"
            @click="openRedeem(item)"
          >
            {{ getButtonConfig(item).label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-if="!loading && filteredItems.length === 0" class="text-center py-16">
      <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">storefront</span>
      <p class="text-sm text-text-muted">No items in this category</p>
    </div>

    <!-- Redeem Confirmation Modal -->
    <Modal :open="showRedeem" title="Confirm Redemption" @close="showRedeem = false">
      <div v-if="redeemItem" class="text-center py-2">
        <p class="text-text-primary text-base mb-2">Redeem <strong>{{ redeemItem.name }}</strong>?</p>
        <p class="text-text-secondary text-sm">This will cost <span class="text-c-accent font-semibold">{{ redeemItem.price.toLocaleString() }} points</span></p>
        <p class="text-xs text-text-muted mt-2">Remaining balance: {{ (userPoints - redeemItem.price).toLocaleString() }} pts</p>
      </div>
      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showRedeem = false">Cancel</button>
        <button
          class="px-4 py-2 bg-c-accent text-black text-sm font-medium rounded-lg hover:bg-c-accent/90 disabled:opacity-50"
          :disabled="redeeming"
          @click="confirmRedeem"
        >
          {{ redeeming ? 'Redeeming...' : 'Confirm' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
