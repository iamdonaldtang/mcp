import type { CampaignStatus } from './common'

// === WL Hub ===
export interface WhiteLabelOverview {
  id: string
  communityId: string
  status: CampaignStatus
  deploymentPath: 'domain' | 'embed' | 'sdk'
  customDomain?: string
  brandConfig: WLBrandConfig
  createdAt: string
  isIntegrationVerified: boolean
  firstApiPingAt?: string
}

export interface WLBrandConfig {
  logo?: string
  primaryColor: string
  accentColor: string
  favicon?: string
  customCss?: string
}

// === Deployment Paths ===
export interface DomainSetup {
  domain: string
  dnsStatus: 'pending' | 'verified' | 'failed'
  sslStatus: 'pending' | 'active' | 'failed'
  cnameTarget: string
}

export interface EmbedConfig {
  mode: 'iframe' | 'widget_library' | 'page_builder'
  iframeUrl?: string
  widgetConfig?: WidgetConfig[]
  pageBuilderConfig?: PageBuilderConfig
}

// === Widget Library ===
export interface WidgetConfig {
  id: string
  moduleType: string
  moduleName: string
  isConfigured: boolean
  isActive: boolean
  embedCode?: string
  settings: Record<string, unknown>
}

// === Page Builder ===
export interface PageBuilderConfig {
  id: string
  pages: PageBuilderPage[]
}

export interface PageBuilderPage {
  id: string
  name: string
  slug: string
  layout: PageBuilderBlock[]
  isPublished: boolean
  createdAt: string
  updatedAt: string
}

export interface PageBuilderBlock {
  id: string
  type: 'hero' | 'tasks' | 'leaderboard' | 'shop' | 'custom_html' | 'milestones'
  config: Record<string, unknown>
  order: number
}

// === Smart Rewards ===
export interface RewardRule {
  id: string
  name: string
  description: string
  triggerType: 'task_complete' | 'level_up' | 'streak' | 'milestone' | 'referral' | 'custom'
  triggerConfig: Record<string, unknown>
  rewardType: 'points' | 'token' | 'nft' | 'whitelist' | 'badge'
  rewardValue: string
  isActive: boolean
  executionCount: number
  createdAt: string
}

export interface Privilege {
  id: string
  name: string
  description: string
  levelRequired: number
  pointCost?: number
  type: 'access' | 'discount' | 'exclusive_content' | 'custom'
  config: Record<string, unknown>
  isActive: boolean
  claims: number
  createdAt: string
}

// === SDK / Integration ===
export interface SDKConfig {
  projectId: string
  apiKey: string
  apiSecret: string
  webhookUrl?: string
  allowedOrigins: string[]
  ssoConfig?: SSOConfig
}

export interface SSOConfig {
  type: 'wallet' | 'oauth'
  oauthProvider?: string
  clientId?: string
  clientSecret?: string
  redirectUri?: string
}

export interface IntegrationStatus {
  isVerified: boolean
  firstPingAt?: string
  lastPingAt?: string
  totalApiCalls: number
  errors24h: number
}

// === Dev Kit ===
export interface DevKit {
  url: string
  projectId: string
  integrationCode: string
  ssoOptions: SSOConfig[]
  verificationTool: string
  estimatedTime: string
}

// === WL Wizard ===
export interface WLWizardData {
  step1: {
    path: 'domain' | 'embed' | 'sdk'
  }
  step1_5?: {
    embedMode: 'iframe' | 'widget_library' | 'page_builder'
  }
  step2: {
    domainSetup?: DomainSetup
    widgetSelection?: string[]
    sdkKeys?: SDKConfig
  }
  step3: {
    brand: WLBrandConfig
  }
  step4: {
    readinessChecks: { label: string; passed: boolean }[]
  }
}

// === WL Getting Started Checklist ===
export interface WLChecklistItem {
  id: string
  label: string
  completed: boolean
  autoDetected: boolean
  action?: string
  route?: string
}

// === Contract Registry ===
export interface ContractEntry {
  id: string
  name: string
  address: string
  chain: string
  abi?: string
  isVerified: boolean
  addedAt: string
}
