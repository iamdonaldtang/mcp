import type {
  TaskType, TaskStatus, QuestStatus, SprintStatus,
  RewardTierType, RewardTierStatus, ShopItemCategory,
  ShopItemStatus, ActivityType, ReferralStatus, MilestoneStatus,
} from './common'

// === C01 Home ===
export interface CTask {
  id: string
  name: string
  type: TaskType
  status: TaskStatus
  points: number
  icon: string
  iconColor: string
  sectorName: string
  completions: number
  requirement?: string
  cooldownEnd?: string
}

export interface Sector {
  id: string
  name: string
  tasks: CTask[]
}

export interface DayChainState {
  currentStreak: number
  targetDays: number
  checkedInToday: boolean
  calendar: { day: number; status: 'completed' | 'today' | 'future' | 'missed' }[]
  bonusMultiplier?: number
  bonusAtDay?: number
}

export interface Announcement {
  id: string
  title: string
  description: string
  image?: string
  link?: string
}

export interface FeaturedSlot {
  id: string
  title: string
  image: string
  link: string
}

export interface PulseStats {
  totalMembers: number
  thisWeek: number
  liveActive: number
  tasksDone: number
}

export interface UserStatus {
  walletAddress: string
  avatar?: string
  level: number
  xp: number
  dayStreak: number
  rank: number
  badges: string[]
}

export interface HomePageData {
  announcements: Announcement[]
  featured: FeaturedSlot[]
  dayChain: DayChainState
  sectors: Sector[]
  userStatus: UserStatus | null
  communityPulse: PulseStats
}

// === C02 Quests ===
export interface CQuest {
  id: string
  name: string
  description: string
  image: string
  status: QuestStatus
  taskCount: number
  completedTasks: number
  participants: number
  endsAt?: string
  points: number
}

// === C03 Leaderboard ===
export interface LeaderboardEntry {
  rank: number
  walletAddress: string
  avatar?: string
  level: number
  points: number
  change?: number
}

export interface LeaderboardData {
  pointType: string
  period: 'weekly' | 'monthly' | 'alltime'
  podium: LeaderboardEntry[]
  rankings: LeaderboardEntry[]
  userRank?: LeaderboardEntry
}

// === C04 LB Sprint ===
export interface SprintTask {
  id: string
  name: string
  points: number
  status: TaskStatus
}

export interface RewardTier {
  threshold: number
  reward: string
  type: RewardTierType
  status: RewardTierStatus
}

export interface Sprint {
  id: string
  name: string
  status: SprintStatus
  startsAt: string
  endsAt: string
  pointType: string
  userPoints: number
  totalParticipants: number
  tasks: SprintTask[]
  rewardTiers: RewardTier[]
  rankings: LeaderboardEntry[]
  userRank: number
}

// === C05 Milestone ===
export interface CMilestone {
  id: string
  name: string
  description: string
  requirement: string
  reward: string
  rewardType: RewardTierType
  status: MilestoneStatus
  progress?: number
  target?: number
}

// === C06 Shop ===
export interface ShopItem {
  id: string
  name: string
  description: string
  image: string
  category: ShopItemCategory
  price: number
  stock: number
  totalRedemptions: number
  isTimeLimited: boolean
  expiresAt?: string
  status: ShopItemStatus
}

// === C07 User Center ===
export interface Achievement {
  id: string
  name: string
  icon: string
  earnedAt?: string
  isEarned: boolean
}

export interface ActivityItem {
  id: string
  type: ActivityType
  title: string
  description?: string
  pointsDelta?: number
  timestamp: string
  iconColor: string
}

export interface UserProfile {
  walletAddress: string
  avatar?: string
  joinedAt: string
  level: number
  lifetimeXp: number
  totalPoints: number
  tasksDone: number
  dayStreak: number
  rank: number
}

export interface ReferralStats {
  totalInvites: number
  successfulJoins: number
  pointsEarned: number
  conversionRate: number
}

// === C08 Invite Center ===
export interface Referral {
  address: string
  status: ReferralStatus
  pointsEarned: number
  date: string
}

export interface InviteData {
  referralCode: string
  referralUrl: string
  totalInvites: number
  successfulJoins: number
  pointsEarned: number
  conversionRate: number
  referrals: Referral[]
}

// === Tab Visibility ===
export interface TabVisibility {
  home: boolean
  quests: boolean
  leaderboard: boolean
  lbSprint: boolean
  milestone: boolean
  shop: boolean
}
