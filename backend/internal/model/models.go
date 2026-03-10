package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// === B-End Auth ===

type User struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	Role      string         `gorm:"default:admin" json:"role"`
	Avatar    string         `json:"avatar,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// === Community ===

type Community struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID         string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `json:"description"`
	Status         string         `gorm:"default:draft" json:"status"` // draft, active, paused, completed
	BrandColor     string         `gorm:"default:#48BB78" json:"brand_color"`
	Logo           string         `json:"logo,omitempty"`
	EnabledModules StringArray    `gorm:"type:text[]" json:"enabled_modules"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// StringArray is a custom type for PostgreSQL text arrays (also works with SQLite via JSON serialization)
type StringArray []string

// Scan implements the sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("unsupported type for StringArray: %T", value)
	}
}

// Value implements the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// === Sectors & Tasks ===

type Sector struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string         `gorm:"type:uuid;not null;index" json:"community_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      string         `gorm:"default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Tasks       []Task         `gorm:"foreignKey:SectorID" json:"tasks,omitempty"`
}

type Task struct {
	ID                string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SectorID          string         `gorm:"type:uuid;not null;index" json:"sector_id"`
	CommunityID       string         `gorm:"type:uuid;not null;index" json:"community_id"`
	Name              string         `gorm:"not null" json:"name"`
	Description       string         `json:"description"`
	Type              string         `gorm:"not null" json:"type"` // social, onchain, verification, custom, recurring, referral
	Status            string         `gorm:"default:draft" json:"status"`
	Points            int            `gorm:"default:0" json:"points"`
	Icon              string         `gorm:"default:task_alt" json:"icon"`
	IconColor         string         `gorm:"default:#F59E0B" json:"icon_color"`
	MaxCompletions    *int           `json:"max_completions,omitempty"`
	CooldownHours     *int           `json:"cooldown_hours,omitempty"`
	Requirement       string         `json:"requirement,omitempty"`
	CurrentCompletions int           `gorm:"default:0" json:"current_completions"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// === Points & Level ===

type PointType struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string   `gorm:"type:uuid;not null;index" json:"community_id"`
	Name        string    `gorm:"not null" json:"name"`
	Symbol      string    `gorm:"not null" json:"symbol"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
}

type LevelConfig struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string `gorm:"type:uuid;not null;index" json:"community_id"`
	Level       int    `gorm:"not null" json:"level"`
	Name        string `gorm:"not null" json:"name"`
	MinPoints   int    `gorm:"not null" json:"min_points"`
	Badge       string `json:"badge,omitempty"`
}

// === TaskChain ===

type TaskChain struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string         `gorm:"type:uuid;not null;index" json:"community_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:draft" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Steps       []TaskChainStep `gorm:"foreignKey:TaskChainID" json:"steps,omitempty"`
}

type TaskChainStep struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TaskChainID string `gorm:"type:uuid;not null;index" json:"task_chain_id"`
	Order       int    `gorm:"not null" json:"order"`
	TaskID      string `gorm:"type:uuid;not null" json:"task_id"`
	BonusPoints int    `gorm:"default:0" json:"bonus_points"`
}

// === DayChain ===

type DayChainConfig struct {
	ID             string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID    string `gorm:"type:uuid;uniqueIndex;not null" json:"community_id"`
	Enabled        bool   `gorm:"default:false" json:"enabled"`
	TargetDays     int    `gorm:"default:30" json:"target_days"`
	DailyPoints    int    `gorm:"default:10" json:"daily_points"`
	CatchUpEnabled bool   `gorm:"default:false" json:"catch_up_enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// === Leaderboard ===

type LeaderboardConfig struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string   `gorm:"type:uuid;not null;index" json:"community_id"`
	PointTypeID string    `gorm:"type:uuid;not null" json:"point_type_id"`
	Periods     StringArray `gorm:"type:text[]" json:"periods"` // weekly, monthly, alltime
	IsEnabled   bool      `gorm:"default:true" json:"is_enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

// === LB Sprint ===

type LBSprint struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID  string         `gorm:"type:uuid;not null;index" json:"community_id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	Status       string         `gorm:"default:draft" json:"status"`
	PointTypeID  string         `gorm:"type:uuid;not null" json:"point_type_id"`
	StartsAt     time.Time      `json:"starts_at"`
	EndsAt       time.Time      `json:"ends_at"`
	Participants int            `gorm:"default:0" json:"participants"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Rewards      []SprintReward `gorm:"foreignKey:SprintID" json:"rewards,omitempty"`
}

type SprintReward struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SprintID   string `gorm:"type:uuid;not null;index" json:"sprint_id"`
	Tier       int    `gorm:"not null" json:"tier"`
	Threshold  int    `gorm:"not null" json:"threshold"`
	RewardType string `gorm:"not null" json:"reward_type"` // token, nft, whitelist, points
	RewardValue string `gorm:"not null" json:"reward_value"`
	Quantity   *int   `json:"quantity,omitempty"`
}

// === Milestone ===

type Milestone struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID  string         `gorm:"type:uuid;not null;index" json:"community_id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	Requirement  string         `json:"requirement"`
	Threshold    int            `gorm:"not null" json:"threshold"`
	RewardType   string         `gorm:"not null" json:"reward_type"`
	RewardValue  string         `gorm:"not null" json:"reward_value"`
	Status       string         `gorm:"default:active" json:"status"`
	Claims       int            `gorm:"default:0" json:"claims"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// === Benefits Shop ===

type ShopItem struct {
	ID               string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID      string         `gorm:"type:uuid;not null;index" json:"community_id"`
	Name             string         `gorm:"not null" json:"name"`
	Description      string         `json:"description"`
	Image            string         `json:"image"`
	Category         string         `gorm:"not null" json:"category"` // nft, voucher, merch, whitelist
	Price            int            `gorm:"not null" json:"price"`
	Stock            int            `gorm:"default:-1" json:"stock"` // -1 = unlimited
	TotalRedemptions int            `gorm:"default:0" json:"total_redemptions"`
	IsTimeLimited    bool           `gorm:"default:false" json:"is_time_limited"`
	ExpiresAt        *time.Time     `json:"expires_at,omitempty"`
	Status           string         `gorm:"default:active" json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// === Lucky Wheel ===

type LuckyWheelConfig struct {
	ID             string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID    string `gorm:"type:uuid;uniqueIndex;not null" json:"community_id"`
	Enabled        bool   `gorm:"default:false" json:"enabled"`
	CostPoints     int    `gorm:"default:0" json:"cost_points"`
	CostPointTypeID string `gorm:"type:uuid" json:"cost_point_type_id"`
	DailySpinLimit int    `gorm:"default:1" json:"daily_spin_limit"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WheelSlot struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	WheelID     string  `gorm:"type:uuid;not null;index" json:"wheel_id"`
	Label       string  `gorm:"not null" json:"label"`
	RewardType  string  `gorm:"not null" json:"reward_type"` // points, token, nft, nothing
	RewardValue string  `json:"reward_value"`
	Probability float64 `gorm:"not null" json:"probability"`
	Color       string  `json:"color"`
}

// === Badges ===

type Badge struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string         `gorm:"type:uuid;not null;index" json:"community_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	Criteria    string         `json:"criteria"`
	AutoAward   bool           `gorm:"default:false" json:"auto_award"`
	AwardedCount int           `gorm:"default:0" json:"awarded_count"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// === Access Rules ===

type AccessRule struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string `gorm:"type:uuid;not null;index" json:"community_id"`
	Type        string `gorm:"not null" json:"type"` // token_gate, nft_gate, whitelist, open
	Config      string `gorm:"type:jsonb" json:"config"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

// === White Label ===

type WhiteLabel struct {
	ID                    string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID           string         `gorm:"type:uuid;uniqueIndex;not null" json:"community_id"`
	UserID                string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Status                string         `gorm:"default:draft" json:"status"`
	DeploymentPath        string         `json:"deployment_path"` // domain, embed, sdk
	CustomDomain          string         `json:"custom_domain,omitempty"`
	BrandLogo             string         `json:"brand_logo,omitempty"`
	BrandPrimaryColor     string         `gorm:"default:#F59E0B" json:"brand_primary_color"`
	BrandAccentColor      string         `gorm:"default:#F59E0B" json:"brand_accent_color"`
	BrandFavicon          string         `json:"brand_favicon,omitempty"`
	BrandCustomCSS        string         `json:"brand_custom_css,omitempty"`
	IsIntegrationVerified bool           `gorm:"default:false" json:"is_integration_verified"`
	FirstAPIPingAt        *time.Time     `json:"first_api_ping_at,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

// === Widget Config ===

type WidgetConfig struct {
	ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	WhiteLabelID string `gorm:"type:uuid;not null;index" json:"white_label_id"`
	ModuleType   string `gorm:"not null" json:"module_type"`
	ModuleName   string `gorm:"not null" json:"module_name"`
	IsConfigured bool   `gorm:"default:false" json:"is_configured"`
	IsActive     bool   `gorm:"default:false" json:"is_active"`
	EmbedCode    string `json:"embed_code,omitempty"`
	Settings     string `gorm:"type:jsonb;default:'{}'" json:"settings"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// === Page Builder ===

type PageBuilderPage struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	WhiteLabelID string         `gorm:"type:uuid;not null;index" json:"white_label_id"`
	Name         string         `gorm:"not null" json:"name"`
	Slug         string         `gorm:"not null" json:"slug"`
	Layout       string         `gorm:"type:jsonb;default:'[]'" json:"layout"` // JSON array of blocks
	IsPublished  bool           `gorm:"default:false" json:"is_published"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// === Smart Rewards ===

type RewardRule struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID    string    `gorm:"type:uuid;not null;index" json:"community_id"`
	Name           string    `gorm:"not null" json:"name"`
	Description    string    `json:"description"`
	TriggerType    string    `gorm:"not null" json:"trigger_type"` // task_complete, level_up, streak, milestone, referral, custom
	TriggerConfig  string    `gorm:"type:jsonb;default:'{}'" json:"trigger_config"`
	RewardType     string    `gorm:"not null" json:"reward_type"` // points, token, nft, whitelist, badge
	RewardValue    string    `gorm:"not null" json:"reward_value"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	ExecutionCount int       `gorm:"default:0" json:"execution_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Privilege struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID   string   `gorm:"type:uuid;not null;index" json:"community_id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `json:"description"`
	LevelRequired int       `gorm:"default:0" json:"level_required"`
	PointCost     *int      `json:"point_cost,omitempty"`
	Type          string    `gorm:"not null" json:"type"` // access, discount, exclusive_content, custom
	Config        string    `gorm:"type:jsonb;default:'{}'" json:"config"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	Claims        int       `gorm:"default:0" json:"claims"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// === SDK / Integration ===

type SDKConfig struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	WhiteLabelID   string    `gorm:"type:uuid;uniqueIndex;not null" json:"white_label_id"`
	ProjectID      string    `gorm:"uniqueIndex;not null" json:"project_id"`
	APIKey         string    `gorm:"not null" json:"api_key"`
	APISecret      string    `gorm:"not null" json:"-"`
	WebhookURL     string    `json:"webhook_url,omitempty"`
	AllowedOrigins StringArray `gorm:"type:text[]" json:"allowed_origins"`
	SSOType        string    `json:"sso_type,omitempty"` // wallet, oauth
	OAuthProvider  string    `json:"oauth_provider,omitempty"`
	OAuthClientID  string    `json:"oauth_client_id,omitempty"`
	OAuthSecret    string    `json:"-"`
	OAuthRedirect  string    `json:"oauth_redirect,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// === Contract Registry ===

type ContractEntry struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string   `gorm:"type:uuid;not null;index" json:"community_id"`
	Name        string    `gorm:"not null" json:"name"`
	Address     string    `gorm:"not null" json:"address"`
	Chain       string    `gorm:"not null" json:"chain"`
	ABI         string    `gorm:"type:text" json:"abi,omitempty"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	AddedAt     time.Time `json:"added_at"`
}

// === C-End: Community Members (wallet users) ===

type CommunityMember struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID   string   `gorm:"type:uuid;not null;index:idx_member_community" json:"community_id"`
	WalletAddress string    `gorm:"not null;index:idx_member_community" json:"wallet_address"`
	Level         int       `gorm:"default:1" json:"level"`
	TotalPoints   int       `gorm:"default:0" json:"total_points"`
	LifetimeXP    int       `gorm:"default:0" json:"lifetime_xp"`
	DayStreak     int       `gorm:"default:0" json:"day_streak"`
	LastCheckIn   *time.Time `json:"last_check_in,omitempty"`
	JoinedAt      time.Time `json:"joined_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// === C-End: Task Completions ===

type TaskCompletion struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID   string   `gorm:"type:uuid;not null;index" json:"community_id"`
	TaskID        string    `gorm:"type:uuid;not null;index" json:"task_id"`
	WalletAddress string    `gorm:"not null;index" json:"wallet_address"`
	Status        string    `gorm:"not null" json:"status"` // completed, claimed
	PointsEarned  int       `gorm:"default:0" json:"points_earned"`
	CompletedAt   time.Time `json:"completed_at"`
	ClaimedAt     *time.Time `json:"claimed_at,omitempty"`
}

// === C-End: Shop Redemptions ===

type ShopRedemption struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID   string   `gorm:"type:uuid;not null;index" json:"community_id"`
	ShopItemID    string    `gorm:"type:uuid;not null;index" json:"shop_item_id"`
	WalletAddress string    `gorm:"not null;index" json:"wallet_address"`
	PointsSpent   int       `gorm:"not null" json:"points_spent"`
	RedeemedAt    time.Time `json:"redeemed_at"`
}

// === C-End: Milestone Claims ===

type MilestoneClaim struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID   string   `gorm:"type:uuid;not null;index" json:"community_id"`
	MilestoneID   string    `gorm:"type:uuid;not null;index" json:"milestone_id"`
	WalletAddress string    `gorm:"not null;index" json:"wallet_address"`
	ClaimedAt     time.Time `json:"claimed_at"`
}

// === C-End: Activity Feed ===

type ActivityLog struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID   string   `gorm:"type:uuid;not null;index" json:"community_id"`
	WalletAddress string    `gorm:"not null;index" json:"wallet_address"`
	Type          string    `gorm:"not null" json:"type"` // task, points, reward, level_up, invite, wheel, streak, shop
	Title         string    `gorm:"not null" json:"title"`
	Description   string    `json:"description,omitempty"`
	PointsDelta   *int      `json:"points_delta,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// === C-End: Referrals ===

type Referral struct {
	ID              string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID     string     `gorm:"type:uuid;not null;index" json:"community_id"`
	ReferrerAddress string     `gorm:"not null;index" json:"referrer_address"`
	RefereeAddress  string     `gorm:"not null" json:"referee_address"`
	ReferralCode    string     `gorm:"not null" json:"referral_code"`
	Status          string     `gorm:"default:pending" json:"status"` // pending, joined
	PointsEarned    int        `gorm:"default:0" json:"points_earned"`
	CreatedAt       time.Time  `json:"created_at"`
	JoinedAt        *time.Time `json:"joined_at,omitempty"`
}

// === Domain Setup ===

type DomainSetup struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	WhiteLabelID string    `gorm:"type:uuid;uniqueIndex;not null" json:"white_label_id"`
	Domain       string    `gorm:"not null" json:"domain"`
	DNSStatus    string    `gorm:"default:pending" json:"dns_status"` // pending, verified, failed
	SSLStatus    string    `gorm:"default:pending" json:"ssl_status"` // pending, active, failed
	CNAMETarget  string    `gorm:"not null" json:"cname_target"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// === Announcements ===

type Announcement struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CommunityID string    `gorm:"type:uuid;not null;index" json:"community_id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Image       string     `json:"image,omitempty"`
	Link        string     `json:"link,omitempty"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// AllModels returns all model structs for GORM AutoMigrate
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&Community{},
		&Sector{},
		&Task{},
		&PointType{},
		&LevelConfig{},
		&TaskChain{},
		&TaskChainStep{},
		&DayChainConfig{},
		&LeaderboardConfig{},
		&LBSprint{},
		&SprintReward{},
		&Milestone{},
		&ShopItem{},
		&LuckyWheelConfig{},
		&WheelSlot{},
		&Badge{},
		&AccessRule{},
		&WhiteLabel{},
		&WidgetConfig{},
		&PageBuilderPage{},
		&RewardRule{},
		&Privilege{},
		&SDKConfig{},
		&ContractEntry{},
		&CommunityMember{},
		&TaskCompletion{},
		&ShopRedemption{},
		&MilestoneClaim{},
		&ActivityLog{},
		&Referral{},
		&DomainSetup{},
		&Announcement{},
	}
}
