<script setup lang="ts">
import { computed } from 'vue'
import type { CampaignStatus } from '../../types/common'
import { STATUS_BADGE_MAP } from '../../types/common'

const props = defineProps<{ status: CampaignStatus }>()

const config = computed(() => STATUS_BADGE_MAP[props.status])
</script>

<template>
  <span
    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
    :style="{ backgroundColor: config.bgColor, color: config.textColor }"
  >
    {{ config.label }}
  </span>
</template>
