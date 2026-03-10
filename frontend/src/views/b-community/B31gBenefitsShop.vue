<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api } from '../../api/client'
import StatusBadge from '../../components/common/StatusBadge.vue'
import Pagination from '../../components/common/Pagination.vue'
import Modal from '../../components/common/Modal.vue'
import type { CampaignStatus } from '../../types/common'

// === Types ===
interface ShopItem {
  id: string
  name: string
  description: string
  category: string
  image_url: string
  price: number
  point_type: string
  stock_type: 'unlimited' | 'limited'
  stock_total: number
  stock_remaining: number
  redemptions: number
  status: CampaignStatus
  gate_type: 'all' | 'level' | 'badge'
  gate_value: string
  created_at: string
}

interface ShopStats {
  total_items: number
  total_redemptions: number
  points_spent: number
  sold_out_count: number
}

// === State ===
const loading = ref(true)
const items = ref<ShopItem[]>([])
const totalItems = ref(0)
const page = ref(1)
const pageSize = 20
const stats = ref<ShopStats>({ total_items: 0, total_redemptions: 0, points_spent: 0, sold_out_count: 0 })

// Filters
const filterStatus = ref<string>('all')
const filterCategory = ref<string>('all')
const searchQuery = ref('')

const statusTabs = [
  { key: 'all', label: 'All' },
  { key: 'active', label: 'Published' },
  { key: 'draft', label: 'Draft' },
  { key: 'sold_out', label: 'Sold Out' },
]

const categories = ['Merch', 'NFT', 'Whitelist', 'Coupon', 'Custom']

// Modal
const showCreateModal = ref(false)
const editingItem = ref<ShopItem | null>(null)
const saving = ref(false)

const formDefaults = {
  name: '',
  description: '',
  category: 'Merch',
  image: null as File | null,
  imagePreview: '',
  price: 100,
  point_type: 'EXP',
  stock_type: 'unlimited' as 'unlimited' | 'limited',
  stock_amount: 100,
  gate_type: 'all' as 'all' | 'level' | 'badge',
  gate_level: 1,
  gate_badge: '',
  publish_now: false,
}

const form = ref({ ...formDefaults })

// Replenish popup
const replenishItemId = ref<string | null>(null)
const replenishAmount = ref(50)

// === Computed ===
const filteredItems = computed(() => {
  return items.value.filter(item => {
    if (filterStatus.value === 'sold_out') {
      if (item.stock_type === 'limited' && item.stock_remaining > 0) return false
      if (item.stock_type === 'unlimited') return false
    } else if (filterStatus.value !== 'all' && item.status !== filterStatus.value) {
      return false
    }
    if (filterCategory.value !== 'all' && item.category !== filterCategory.value) return false
    if (searchQuery.value && !item.name.toLowerCase().includes(searchQuery.value.toLowerCase())) return false
    return true
  })
})

const isSoldOut = (item: ShopItem) => item.stock_type === 'limited' && item.stock_remaining <= 0

const canSave = computed(() => {
  return form.value.name.trim().length >= 1
    && form.value.name.trim().length <= 60
    && form.value.price >= 1
    && (form.value.stock_type === 'unlimited' || form.value.stock_amount > 0)
})

// === API ===
onMounted(async () => {
  await Promise.all([fetchItems(), fetchStats()])
  loading.value = false
})

async function fetchItems() {
  try {
    const params: Record<string, string | number> = { page: page.value, page_size: pageSize }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    if (filterCategory.value !== 'all') params.category = filterCategory.value
    const res = await api.get('/api/v1/community/modules/shop', { params })
    items.value = res.data.data?.items || []
    totalItems.value = res.data.data?.total || 0
  } catch { /* empty */ }
}

async function fetchStats() {
  try {
    const res = await api.get('/api/v1/community/modules/shop/stats')
    stats.value = res.data.data || stats.value
  } catch { /* empty */ }
}

function openCreate() {
  editingItem.value = null
  form.value = { ...formDefaults }
  showCreateModal.value = true
}

function openEdit(item: ShopItem) {
  editingItem.value = item
  form.value = {
    name: item.name,
    description: item.description,
    category: item.category,
    image: null,
    imagePreview: item.image_url,
    price: item.price,
    point_type: item.point_type || 'EXP',
    stock_type: item.stock_type,
    stock_amount: item.stock_total,
    gate_type: item.gate_type,
    gate_level: item.gate_type === 'level' ? parseInt(item.gate_value) || 1 : 1,
    gate_badge: item.gate_type === 'badge' ? item.gate_value : '',
    publish_now: item.status === 'active',
  }
  showCreateModal.value = true
}

async function saveItem() {
  if (!canSave.value) return
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      description: form.value.description.trim(),
      category: form.value.category,
      price: form.value.price,
      point_type: form.value.point_type,
      stock_type: form.value.stock_type,
      stock_total: form.value.stock_type === 'limited' ? form.value.stock_amount : null,
      gate_type: form.value.gate_type,
      gate_value: form.value.gate_type === 'level'
        ? String(form.value.gate_level)
        : form.value.gate_type === 'badge'
          ? form.value.gate_badge
          : null,
      status: form.value.publish_now ? 'active' : 'draft',
    }

    if (editingItem.value) {
      await api.put(`/api/v1/community/modules/shop/${editingItem.value.id}`, payload)
    } else {
      await api.post('/api/v1/community/modules/shop', payload)
    }

    // Upload image if present
    if (form.value.image) {
      const formData = new FormData()
      formData.append('image', form.value.image)
      const itemId = editingItem.value?.id || 'latest'
      await api.post(`/api/v1/community/modules/shop/${itemId}/image`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
    }

    showCreateModal.value = false
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* TODO: toast */ }
  saving.value = false
}

async function duplicateItem(item: ShopItem) {
  try {
    await api.post(`/api/v1/community/modules/shop/${item.id}/duplicate`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

async function deleteItem(id: string) {
  if (!confirm('Delete this shop item? This action cannot be undone.')) return
  try {
    await api.delete(`/api/v1/community/modules/shop/${id}`)
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

async function replenishStock() {
  if (!replenishItemId.value || replenishAmount.value < 1) return
  try {
    await api.post(`/api/v1/community/modules/shop/${replenishItemId.value}/replenish`, {
      amount: replenishAmount.value,
    })
    replenishItemId.value = null
    replenishAmount.value = 50
    await Promise.all([fetchItems(), fetchStats()])
  } catch { /* empty */ }
}

function handleImageDrop(e: DragEvent) {
  e.preventDefault()
  const file = e.dataTransfer?.files?.[0]
  if (file) handleImageFile(file)
}

function handleImageSelect(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) handleImageFile(file)
}

function handleImageFile(file: File) {
  if (!['image/png', 'image/jpeg'].includes(file.type)) return
  if (file.size > 2 * 1024 * 1024) return
  form.value.image = file
  form.value.imagePreview = URL.createObjectURL(file)
}

function onPageChange(p: number) {
  page.value = p
  fetchItems()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">Benefits Shop</h1>
        <p class="text-sm text-text-secondary">Manage redeemable items in your community shop</p>
      </div>
      <button
        class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 transition-colors"
        @click="openCreate"
      >
        + New Item
      </button>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-4 gap-4">
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Total Items</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.total_items.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Total Redemptions</div>
        <div class="text-2xl font-bold text-text-primary">{{ stats.total_redemptions.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Points Spent</div>
        <div class="text-2xl font-bold text-community">{{ stats.points_spent.toLocaleString() }}</div>
      </div>
      <div class="bg-card-bg border border-border rounded-xl p-4">
        <div class="text-xs text-text-muted uppercase tracking-wider mb-1">Items Sold Out</div>
        <div class="text-2xl font-bold" :class="stats.sold_out_count > 0 ? 'text-amber-400' : 'text-text-primary'">{{ stats.sold_out_count }}</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="flex gap-2">
        <button
          v-for="tab in statusTabs"
          :key="tab.key"
          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
          :class="filterStatus === tab.key ? 'bg-community text-white' : 'bg-card-bg border border-border text-text-secondary hover:text-text-primary'"
          @click="filterStatus = tab.key; fetchItems()"
        >
          {{ tab.label }}
        </button>
      </div>
      <select
        v-model="filterCategory"
        class="px-3 py-1.5 bg-card-bg border border-border rounded-lg text-xs text-text-secondary focus:border-community focus:outline-none"
        @change="fetchItems()"
      >
        <option value="all">All Categories</option>
        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
      </select>
      <div class="flex-1">
        <div class="relative">
          <span class="material-symbols-rounded absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-lg">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search items..."
            class="w-full pl-10 pr-4 py-2 bg-card-bg border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:border-community focus:outline-none"
          />
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-card-bg border border-border rounded-xl overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-border">
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Item Name</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Category</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Price</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-36">Redemptions / Stock</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Status</th>
            <th class="px-4 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider w-28">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" v-for="i in 5" :key="i">
            <td colspan="6" class="px-4 py-4"><div class="h-4 bg-border rounded animate-pulse"></div></td>
          </tr>
          <tr v-else-if="filteredItems.length === 0">
            <td colspan="6" class="px-4 py-12 text-center">
              <span class="material-symbols-rounded text-4xl text-text-muted block mb-2">storefront</span>
              <p class="text-sm text-text-muted">No shop items yet</p>
              <button class="mt-2 text-xs text-community hover:underline" @click="openCreate">+ Create First Item</button>
            </td>
          </tr>
          <tr
            v-else
            v-for="item in filteredItems"
            :key="item.id"
            class="border-b border-border last:border-b-0 hover:bg-white/[0.02] transition-colors"
            :class="isSoldOut(item) ? 'opacity-75' : ''"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-page-bg border border-border flex-shrink-0 overflow-hidden">
                  <img
                    v-if="item.image_url"
                    :src="item.image_url"
                    :alt="item.name"
                    class="w-full h-full object-cover"
                  />
                  <span v-else class="flex items-center justify-center w-full h-full">
                    <span class="material-symbols-rounded text-lg text-text-muted">image</span>
                  </span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-sm font-medium" :class="isSoldOut(item) ? 'text-amber-400' : 'text-text-primary'">{{ item.name }}</span>
                  <span
                    v-if="item.gate_type !== 'all'"
                    class="material-symbols-rounded text-sm text-amber-400"
                    :title="item.gate_type === 'level' ? `Level ${item.gate_value} required` : `Badge: ${item.gate_value}`"
                  >lock</span>
                </div>
              </div>
            </td>
            <td class="px-4 py-3">
              <span class="px-2 py-0.5 text-xs rounded bg-page-bg text-text-secondary">{{ item.category }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-sm font-medium text-community">{{ item.price.toLocaleString() }} {{ item.point_type || 'pts' }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-sm" :class="isSoldOut(item) ? 'text-amber-400' : 'text-text-secondary'">
                {{ item.redemptions.toLocaleString() }} /
                {{ item.stock_type === 'unlimited' ? '∞' : item.stock_total.toLocaleString() }}
              </span>
              <span v-if="item.stock_type === 'limited'" class="text-xs text-text-muted ml-1">
                ({{ item.stock_remaining }} left)
              </span>
            </td>
            <td class="px-4 py-3">
              <template v-if="isSoldOut(item)">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-amber-900/30 text-amber-400">
                  Sold Out
                </span>
              </template>
              <StatusBadge v-else :status="item.status" />
            </td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1 relative">
                <button
                  v-if="isSoldOut(item)"
                  class="px-2 py-1 text-xs font-medium rounded bg-amber-900/30 text-amber-400 hover:bg-amber-900/50 transition-colors"
                  @click="replenishItemId = item.id; replenishAmount = 50"
                >
                  Replenish
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Edit" @click="openEdit(item)">
                  <span class="material-symbols-rounded text-base text-text-muted">edit</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Duplicate" @click="duplicateItem(item)">
                  <span class="material-symbols-rounded text-base text-text-muted">content_copy</span>
                </button>
                <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" title="Delete" @click="deleteItem(item.id)">
                  <span class="material-symbols-rounded text-base text-status-paused">delete</span>
                </button>

                <!-- Replenish popup -->
                <div
                  v-if="replenishItemId === item.id"
                  class="absolute right-0 top-full mt-1 z-20 bg-card-bg border border-border rounded-lg p-3 shadow-xl w-52"
                >
                  <div class="text-xs text-text-secondary mb-2">Add stock quantity</div>
                  <input
                    v-model.number="replenishAmount"
                    type="number"
                    min="1"
                    class="w-full px-3 py-1.5 bg-page-bg border border-border rounded-lg text-sm text-text-primary focus:border-community focus:outline-none mb-2"
                  />
                  <div class="flex justify-end gap-2">
                    <button class="px-2 py-1 text-xs text-text-muted hover:text-text-primary" @click="replenishItemId = null">Cancel</button>
                    <button
                      class="px-3 py-1 text-xs bg-community text-white rounded-lg hover:bg-community/90"
                      @click="replenishStock"
                    >
                      Confirm
                    </button>
                  </div>
                </div>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination v-if="totalItems > pageSize" :page="page" :page-size="pageSize" :total="totalItems" @update:page="onPageChange" />
    </div>

    <!-- D07 Shop Item Editor Modal -->
    <Modal
      :open="showCreateModal"
      :title="editingItem ? 'Edit Shop Item' : 'Create Shop Item'"
      max-width="600px"
      @close="showCreateModal = false"
    >
      <div class="space-y-4 max-h-[60vh] overflow-y-auto pr-1">
        <!-- Item Name -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Item Name <span class="text-status-paused">*</span></label>
          <input
            v-model="form.name"
            type="text"
            maxlength="60"
            placeholder="e.g. Exclusive NFT Pass"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm"
          />
          <div class="text-xs text-text-muted mt-1 text-right">{{ form.name.length }}/60</div>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Description</label>
          <textarea
            v-model="form.description"
            rows="3"
            maxlength="300"
            placeholder="Describe this shop item..."
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary placeholder-text-muted focus:border-community focus:outline-none text-sm resize-none"
          />
          <div class="text-xs text-text-muted mt-1 text-right">{{ form.description.length }}/300</div>
        </div>

        <!-- Category -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Category</label>
          <select
            v-model="form.category"
            class="w-full px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
          >
            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>

        <!-- Image Upload -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Image</label>
          <div
            class="border-2 border-dashed border-border rounded-lg p-6 text-center cursor-pointer hover:border-community/50 transition-colors"
            @dragover.prevent
            @drop="handleImageDrop"
            @click="($refs.imageInput as HTMLInputElement)?.click()"
          >
            <template v-if="form.imagePreview">
              <img :src="form.imagePreview" class="w-20 h-20 object-cover rounded-lg mx-auto mb-2" />
              <p class="text-xs text-text-muted">Click or drop to replace</p>
            </template>
            <template v-else>
              <span class="material-symbols-rounded text-3xl text-text-muted block mb-1">cloud_upload</span>
              <p class="text-sm text-text-secondary">Drop image or click to upload</p>
              <p class="text-xs text-text-muted mt-1">PNG/JPG, max 2MB</p>
            </template>
          </div>
          <input ref="imageInput" type="file" accept="image/png,image/jpeg" class="hidden" @change="handleImageSelect" />
        </div>

        <!-- Price -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Price <span class="text-status-paused">*</span></label>
          <div class="flex gap-3">
            <input
              v-model.number="form.price"
              type="number"
              min="1"
              class="flex-1 px-4 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            />
            <select
              v-model="form.point_type"
              class="w-32 px-3 py-2.5 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            >
              <option value="EXP">EXP</option>
              <option value="GEM">GEM</option>
              <option value="POINTS">Points</option>
            </select>
          </div>
        </div>

        <!-- Stock -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Stock</label>
          <div class="flex items-center gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.stock_type" type="radio" value="unlimited" class="accent-community" />
              <span class="text-sm text-text-primary">Unlimited</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.stock_type" type="radio" value="limited" class="accent-community" />
              <span class="text-sm text-text-primary">Limited</span>
            </label>
            <input
              v-if="form.stock_type === 'limited'"
              v-model.number="form.stock_amount"
              type="number"
              min="1"
              placeholder="Quantity"
              class="w-32 px-3 py-2 bg-page-bg border border-border rounded-lg text-text-primary text-sm focus:border-community focus:outline-none"
            />
          </div>
        </div>

        <!-- Availability Gate -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Availability Gate</label>
          <div class="space-y-2">
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.gate_type" type="radio" value="all" class="accent-community" />
              <span class="text-sm text-text-primary">All Users</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.gate_type" type="radio" value="level" class="accent-community" />
              <span class="text-sm text-text-primary">Level Required</span>
              <input
                v-if="form.gate_type === 'level'"
                v-model.number="form.gate_level"
                type="number"
                min="1"
                class="w-20 px-2 py-1 bg-page-bg border border-border rounded text-text-primary text-sm focus:border-community focus:outline-none ml-2"
              />
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.gate_type" type="radio" value="badge" class="accent-community" />
              <span class="text-sm text-text-primary">Badge Required</span>
              <input
                v-if="form.gate_type === 'badge'"
                v-model="form.gate_badge"
                type="text"
                placeholder="Badge name"
                class="w-40 px-2 py-1 bg-page-bg border border-border rounded text-text-primary text-sm focus:border-community focus:outline-none ml-2"
              />
            </label>
          </div>
        </div>

        <!-- Status -->
        <div>
          <label class="block text-sm text-text-secondary mb-1">Status</label>
          <div class="flex items-center gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.publish_now" type="radio" :value="false" class="accent-community" />
              <span class="text-sm text-text-primary">Draft</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="form.publish_now" type="radio" :value="true" class="accent-community" />
              <span class="text-sm text-text-primary">Publish Now</span>
            </label>
          </div>
        </div>
      </div>

      <template #footer>
        <button class="px-4 py-2 text-sm text-text-muted hover:text-text-primary" @click="showCreateModal = false">Cancel</button>
        <button
          class="px-4 py-2 bg-community text-white text-sm font-medium rounded-lg hover:bg-community/90 disabled:opacity-50 transition-colors"
          :disabled="!canSave || saving"
          @click="saveItem"
        >
          {{ saving ? 'Saving...' : editingItem ? 'Update Item' : 'Create Item' }}
        </button>
      </template>
    </Modal>
  </div>
</template>
