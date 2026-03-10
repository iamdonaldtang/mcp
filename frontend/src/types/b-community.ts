import type { CampaignStatus, CommunityModule, PaginatedData } from './common'

// === Community Hub ===
export interface CommunityOverview {
  id: string
  name: string
  description: string
  status: CampaignStatus
  brandColor: string
  logo?: string
  createdAt: string
  totalMembers: number
  activeMembers: number
  totalTasks: number
  totalPoints: number
  enabledModules: CommunityModule[]
}

// === Sectors & Tasks (B-end) ===
export interface BSector {
  id: string
  name: string
  description: string
  sortOrder: number
  taskCount: number
  status: CampaignStatus
}

export interface BTask {
  id: string
  sectorId: string
  name: string
  description: string
  type: 'social' | 'onchain' | 'verification' | 'custom' | 'recurring' | 'referral'
  status: CampaignStatus
  points: number
  icon: string
  iconColor: string
  maxCompletions?: number
  currentCompletions: number
  cooldownHours?: number
  requirement?: string
  createdAt: string
  updatedAt: string
}

// === Points & Level ===
export interface PointType {
  id: string
  name: string
  symbol: string
  icon: string
  color: string
  isDefault: boolean
}

export interface LevelConfig {
  level: number
  name: string
  minPoints: number
  badge?: string
}

// === TaskChain ===
export interface TaskChain {
  id: string
  name: string
  description: string
  status: CampaignStatus
  steps: TaskChainStep[]
  totalCompletions: number
  createdAt: string
}

export interface TaskChainStep {
  id: string
  order: number
  taskId: string
  taskName: string
  bonusPoints: number
}

// === DayChain ===
export interface DayChainConfig {
  id: string
  enabled: boolean
  targetDays: number
  dailyPoints: number
  bonusDays: { day: number; multiplier: number }[]
  catchUpEnabled: boolean
}

// === Leaderboard Config ===
export interface LeaderboardConfig {
  id: string
  pointTypeId: string
  periods: ('weekly' | 'monthly' | 'alltime')[]
  isEnabled: boolean
}

// === LB Sprint ===
export interface BLBSprint {
  id: string
  name: string
  description: string
  status: CampaignStatus
  pointTypeId: string
  startsAt: string
  endsAt: string
  rewards: SprintReward[]
  participants: number
  createdAt: string
}

export interface SprintReward {
  tier: number
  threshold: number
  rewardType: 'token' | 'nft' | 'whitelist' | 'points'
  rewardValue: string
  quantity?: number
}

// === Milestone ===
export interface BMilestone {
  id: string
  name: string
  description: string
  requirement: string
  threshold: number
  rewardType: 'token' | 'nft' | 'whitelist' | 'points'
  rewardValue: string
  status: CampaignStatus
  claims: number
  createdAt: string
}

// === Benefits Shop ===
export interface BShopItem {
  id: string
  name: string
  description: string
  image: string
  category: 'nft' | 'voucher' | 'merch' | 'whitelist'
  price: number
  stock: number
  totalRedemptions: number
  isTimeLimited: boolean
  expiresAt?: string
  status: CampaignStatus
  createdAt: string
}

// === Lucky Wheel ===
export interface LuckyWheelConfig {
  id: string
  enabled: boolean
  costPoints: number
  costPointTypeId: string
  slots: WheelSlot[]
  dailySpinLimit: number
}

export interface WheelSlot {
  id: string
  label: string
  rewardType: 'points' | 'token' | 'nft' | 'nothing'
  rewardValue: string
  probability: number
  color: string
}

// === Badges ===
export interface BBadge {
  id: string
  name: string
  description: string
  icon: string
  criteria: string
  autoAward: boolean
  awardedCount: number
  createdAt: string
}

// === Access Rules ===
export interface AccessRule {
  id: string
  type: 'token_gate' | 'nft_gate' | 'whitelist' | 'open'
  config: Record<string, unknown>
  enabled: boolean
}

// === Wizard ===
export interface CommunityWizardData {
  step1: {
    name: string
    description: string
    brandColor: string
    logo?: string
  }
  step2: {
    strategy: 'activate' | 'engage' | 'retain'
    enabledModules: CommunityModule[]
  }
  step3: {
    sectors: { name: string; tasks: Partial<BTask>[] }[]
    pointConfig: Partial<PointType>
  }
  step4: {
    previewMode: 'desktop' | 'mobile'
    readinessChecks: { label: string; passed: boolean }[]
  }
}

// === Community Insights ===
export interface CommunityInsights {
  memberGrowth: { date: string; count: number }[]
  taskCompletionRate: number
  avgDailyActive: number
  topModules: { module: CommunityModule; engagement: number }[]
  retentionByModule: { module: CommunityModule; rate: number }[]
  economyHealth: {
    totalPointsIssued: number
    totalPointsSpent: number
    circulatingSupply: number
  }
}

// === Getting Started Checklist ===
export interface ChecklistItem {
  id: string
  label: string
  completed: boolean
  autoDetected: boolean
  action?: string
  route?: string
}
