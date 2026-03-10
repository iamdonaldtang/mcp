// === API Response ===
export interface ApiResponse<T = unknown> {
  code: number
  data: T
  message: string
}

// === Pagination ===
export interface PaginationParams {
  page: number
  pageSize: number
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

// === Status Enums ===
export type CampaignStatus = 'active' | 'draft' | 'completed' | 'paused'

export type TaskType = 'social' | 'onchain' | 'verification' | 'custom' | 'recurring' | 'referral'

export type TaskStatus = 'available' | 'in_progress' | 'completed' | 'claimed' | 'locked' | 'expired' | 'cooldown'

export type QuestStatus = 'available' | 'completed' | 'in_progress' | 'ended'

export type SprintStatus = 'active' | 'ended' | 'upcoming'

export type MilestoneStatus = 'earned' | 'claimable' | 'locked'

export type ShopItemCategory = 'nft' | 'voucher' | 'merch' | 'whitelist'

export type ShopItemStatus = 'available' | 'sold_out'

export type RewardTierType = 'token' | 'nft' | 'whitelist' | 'points'

export type RewardTierStatus = 'earned' | 'claimable' | 'locked'

export type ActivityType = 'task' | 'points' | 'reward' | 'level_up' | 'invite' | 'wheel' | 'streak' | 'shop'

export type ReferralStatus = 'joined' | 'pending'

// === B-End Module Types ===
export type CommunityModule =
  | 'sectors_tasks'
  | 'task_chain'
  | 'day_chain'
  | 'points_level'
  | 'leaderboard'
  | 'badges'
  | 'lb_sprint'
  | 'milestone'
  | 'lucky_wheel'
  | 'benefits_shop'
  | 'access_rules'
  | 'homepage_editor'

export type ProductType = 'quest' | 'community' | 'whitelabel' | 'boost'

// === Shared UI Types ===
export interface StatusBadgeConfig {
  label: string
  bgColor: string
  textColor: string
}

export const STATUS_BADGE_MAP: Record<CampaignStatus, StatusBadgeConfig> = {
  active: { label: 'Active', bgColor: '#0A2E1A', textColor: '#16A34A' },
  draft: { label: 'Draft', bgColor: '#1F1A08', textColor: '#D97706' },
  completed: { label: 'Completed', bgColor: '#1E293B', textColor: '#64748B' },
  paused: { label: 'Paused', bgColor: '#2D1515', textColor: '#DC2626' },
}
