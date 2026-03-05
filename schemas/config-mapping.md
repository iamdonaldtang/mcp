# Quest配置格式映射文档

> 本文档说明如何将LLM生成的QuestConfig转换为API的SetCampaignParams格式

---

## 1. 格式概述

### LLM生成格式 (QuestConfig)
```json
{
  "basic_info": { ... },
  "tasks": [ ... ],
  "rewards": { ... },
  "eligibility": { ... }
}
```

### API要求格式 (SetCampaignParams)
```json
{
  "id": 0,
  "name": "",
  "desc": "",
  "campaign_start": 1234567890,
  "campaign_end": 1234567890,
  "campaign_image": "",
  "campaign_type": "Quest",
  "is_draft": false,
  "is_private": false,
  "eligibility_express": "and",
  "winner_rewards": [ ... ],
  "qualifier_rewards": [ ... ],
  "eligs": [ ... ],
  "tasks": [ ... ],
  "min_finished_optional_task_num": 0,
  "min_finished_optional_task_points": 0,
  "google_recaptcha": false
}
```

---

## 2. 字段映射

### 2.1 基础信息映射

| QuestConfig | SetCampaignParams | 转换说明 |
|-------------|-------------------|----------|
| `basic_info.name` | `name` | 直接复制 |
| `basic_info.description` | `desc` | 直接复制 |
| `basic_info.cover_image` | `campaign_image` | 直接复制 |
| `basic_info.start_time` | `campaign_start` | ISO 8601 → Unix时间戳 |
| `basic_info.end_time` | `campaign_end` | ISO 8601 → Unix时间戳 |
| `basic_info.is_private` | `is_private` | 直接复制 |
| - | `id` | 新建时为0 |
| - | `campaign_type` | 固定为"Quest" |
| - | `is_draft` | 发布时为false |

### 2.2 任务映射

**QuestConfig格式:**
```json
{
  "template_id": "FollowTwitter",
  "custom_name": "关注我们",
  "params": { "twitter_handle": "@SwapX" },
  "points": 100,
  "is_optional": false,
  "recurrence": "once"
}
```

**SetCampaignParams格式:**
```json
{
  "template_id": "FollowTwitter",
  "template_name": "FollowTwitter",
  "params": "{\"twitter_handle\":\"@SwapX\"}",
  "points": { "amount": 100 },
  "is_optional": false,
  "is_template_task": false,
  "is_ecosystem_task": false,
  "class_type": "OffChain",
  "is_hold": false,
  "platform": "Twitter",
  "custom_name": "关注我们",
  "recurrence": "Once"
}
```

**转换规则:**

| QuestConfig字段 | SetCampaignParams字段 | 转换说明 |
|-----------------|----------------------|----------|
| `template_id` | `template_id` | 直接复制 |
| `template_id` | `template_name` | 直接复制 |
| `params` | `params` | JSON对象 → JSON字符串 |
| `points` | `points.amount` | 封装为对象 |
| `is_optional` | `is_optional` | 直接复制 |
| - | `is_template_task` | 默认false |
| - | `is_ecosystem_task` | 默认false |
| - | `class_type` | 根据template_id推断 |
| - | `is_hold` | 根据template_id推断 |
| - | `platform` | 根据template_id推断 |
| `recurrence` | `recurrence` | 首字母大写: once→Once |

**class_type推断规则:**
```go
func inferClassType(templateId string) string {
    switch templateId {
    // OffChain类型
    case "FollowTwitter", "RetweetTwitter", "LikeATweet", "PostTweet",
         "JoinDiscord", "JoinTelegram", "DailyConnect", "Invite",
         "QuizChoose", "PowTask", "BindEmail":
        return "OffChain"
    // OnChain类型
    case "TokenBalance", "NftHolder", "SwapVolume", "ContractInteractive":
        return "OnChain"
    default:
        return "OffChain"
    }
}
```

**platform推断规则:**
```go
func inferPlatform(templateId string) string {
    switch templateId {
    case "FollowTwitter", "RetweetTwitter", "LikeATweet", "PostTweet":
        return "Twitter"
    case "JoinDiscord":
        return "Discord"
    case "JoinTelegram":
        return "Telegram"
    case "TokenBalance", "NftHolder", "SwapVolume", "ContractInteractive":
        return "Wallet"
    default:
        return ""
    }
}
```

### 2.3 奖励映射

**QuestConfig格式:**
```json
{
  "distribution_method": "lucky_draw",
  "layers": [
    {
      "max_winners": 100,
      "rewards": [
        { "type": "token", "amount": "10", "token_symbol": "USDC", "chain": "ethereum" }
      ]
    }
  ]
}
```

**SetCampaignParams格式:**
```json
{
  "winner_rewards": [
    {
      "winner_draw_type": "LuckyDraw",
      "automatically_winner_draw_type": "Automatic",
      "winner_layer_rewards": [
        {
          "winner_layer": { "max_winners": 100 },
          "max_winners": 100,
          "rewards": [
            {
              "reward_type": "Token",
              "chain": "ethereum",
              "symbol": "USDC",
              "amount": "10",
              "token_address": "0x..."
            }
          ]
        }
      ]
    }
  ]
}
```

**distribution_method映射:**

| QuestConfig | SetCampaignParams.winner_draw_type |
|-------------|-----------------------------------|
| `fcfs` | `FCFS` |
| `lucky_draw` | `LuckyDraw` |
| `ranking` | `Ranking` |
| `open_to_all` | `OpenToAll` |

**reward.type映射:**

| QuestConfig | SetCampaignParams.reward_type |
|-------------|-------------------------------|
| `token` | `Token` |
| `nft` | `NFT` |
| `points` | `Points` |
| `whitelist` | `Whitelist` |
| `discord_role` | `DiscordRole` |

### 2.4 资格条件映射

**QuestConfig格式:**
```json
{
  "express": "and",
  "conditions": [
    {
      "template_id": "TokenBalance",
      "params": { "network": "ethereum", "min_balance": "100" }
    }
  ]
}
```

**SetCampaignParams格式:**
```json
{
  "eligibility_express": "and",
  "eligs": [
    {
      "template_id": "TokenBalance",
      "params": "{\"network\":\"ethereum\",\"min_balance\":\"100\"}"
    }
  ]
}
```

**转换规则:**
- `express` → `eligibility_express`: 直接复制
- `conditions[].params` → `eligs[].params`: JSON对象 → JSON字符串

---

## 3. Go转换代码示例

```go
package quest

import (
    "encoding/json"
    "time"
)

// QuestConfig LLM生成的配置格式
type QuestConfig struct {
    BasicInfo   BasicInfo         `json:"basic_info"`
    Tasks       []TaskConfig      `json:"tasks"`
    Rewards     *RewardConfig     `json:"rewards,omitempty"`
    Eligibility *EligibilityConfig `json:"eligibility,omitempty"`
}

// SetCampaignParams API要求的格式
type SetCampaignParams struct {
    ID                          int                 `json:"id"`
    Name                        string              `json:"name"`
    Desc                        string              `json:"desc"`
    CampaignStart               int64               `json:"campaign_start"`
    CampaignEnd                 int64               `json:"campaign_end"`
    CampaignImage               string              `json:"campaign_image"`
    CampaignType                string              `json:"campaign_type"`
    IsDraft                     bool                `json:"is_draft"`
    IsPrivate                   bool                `json:"is_private"`
    EligibilityExpress          string              `json:"eligibility_express"`
    WinnerRewards               []WinnerRewards     `json:"winner_rewards"`
    QualifierRewards            []RewardInfo        `json:"qualifier_rewards"`
    Eligs                       []EligibilityItem   `json:"eligs"`
    Tasks                       []TaskParamItem     `json:"tasks"`
    MinFinishedOptionalTaskNum  int                 `json:"min_finished_optional_task_num"`
    MinFinishedOptionalTaskPoints int               `json:"min_finished_optional_task_points"`
    GoogleRecaptcha             bool                `json:"google_recaptcha"`
    IsPreview                   bool                `json:"is_preview"`
}

// ToSetCampaignParams 转换为API格式
func (q *QuestConfig) ToSetCampaignParams() (*SetCampaignParams, error) {
    params := &SetCampaignParams{
        ID:            0, // 0表示创建新Campaign
        Name:          q.BasicInfo.Name,
        Desc:          q.BasicInfo.Description,
        CampaignImage: q.BasicInfo.CoverImage,
        CampaignType:  "Quest",
        IsDraft:       false,
        IsPrivate:     q.BasicInfo.IsPrivate,
    }

    // 转换时间
    startTime, err := time.Parse(time.RFC3339, q.BasicInfo.StartTime)
    if err != nil {
        return nil, err
    }
    params.CampaignStart = startTime.Unix()

    endTime, err := time.Parse(time.RFC3339, q.BasicInfo.EndTime)
    if err != nil {
        return nil, err
    }
    params.CampaignEnd = endTime.Unix()

    // 转换任务
    for _, task := range q.Tasks {
        taskItem, err := convertTask(task)
        if err != nil {
            return nil, err
        }
        params.Tasks = append(params.Tasks, taskItem)
    }

    // 转换奖励
    if q.Rewards != nil {
        params.WinnerRewards = convertRewards(*q.Rewards)
    }

    // 转换资格条件
    if q.Eligibility != nil {
        params.EligibilityExpress = q.Eligibility.Express
        params.Eligs = convertEligibility(q.Eligibility.Conditions)
    }

    return params, nil
}

func convertTask(task TaskConfig) (TaskParamItem, error) {
    paramsJSON, err := json.Marshal(task.Params)
    if err != nil {
        return TaskParamItem{}, err
    }

    return TaskParamItem{
        TemplateID:      task.TemplateID,
        TemplateName:    task.TemplateID,
        Params:          string(paramsJSON),
        Points:          PointsInfo{Amount: task.Points},
        IsOptional:      task.IsOptional,
        IsTemplateTask:  false,
        IsEcosystemTask: false,
        ClassType:       inferClassType(task.TemplateID),
        IsHold:          false,
        Platform:        inferPlatform(task.TemplateID),
        CustomName:      task.CustomName,
        Recurrence:      capitalizeFirst(task.Recurrence),
    }, nil
}

func convertRewards(rewards RewardConfig) []WinnerRewards {
    // 实现奖励转换逻辑
    // ...
}

func convertEligibility(conditions []EligibilityCondition) []EligibilityItem {
    var items []EligibilityItem
    for _, cond := range conditions {
        paramsJSON, _ := json.Marshal(cond.Params)
        items = append(items, EligibilityItem{
            TemplateID: cond.TemplateID,
            Params:     string(paramsJSON),
        })
    }
    return items
}
```

---

## 4. 验证规则

### 4.1 必填字段检查

```go
func ValidateQuestConfig(config *QuestConfig) error {
    // 基础信息
    if config.BasicInfo.Name == "" {
        return errors.New("活动名称不能为空")
    }
    if config.BasicInfo.StartTime == "" {
        return errors.New("开始时间不能为空")
    }
    if config.BasicInfo.EndTime == "" {
        return errors.New("结束时间不能为空")
    }

    // 时间验证
    start, _ := time.Parse(time.RFC3339, config.BasicInfo.StartTime)
    end, _ := time.Parse(time.RFC3339, config.BasicInfo.EndTime)
    if !end.After(start) {
        return errors.New("结束时间必须晚于开始时间")
    }

    // 任务验证
    if len(config.Tasks) == 0 {
        return errors.New("至少需要一个任务")
    }

    for i, task := range config.Tasks {
        if task.TemplateID == "" {
            return fmt.Errorf("任务 %d 缺少template_id", i+1)
        }
        if err := validateTaskParams(task); err != nil {
            return fmt.Errorf("任务 %d: %v", i+1, err)
        }
    }

    return nil
}
```

### 4.2 任务参数验证

```go
func validateTaskParams(task TaskConfig) error {
    switch task.TemplateID {
    case "FollowTwitter":
        handle, ok := task.Params["twitter_handle"].(string)
        if !ok || handle == "" {
            return errors.New("twitter_handle不能为空")
        }
        if !strings.HasPrefix(handle, "@") {
            return errors.New("twitter_handle必须以@开头")
        }

    case "JoinDiscord":
        url, ok := task.Params["discord_server_url"].(string)
        if !ok || url == "" {
            return errors.New("discord_server_url不能为空")
        }
        if !strings.Contains(url, "discord") {
            return errors.New("无效的Discord链接")
        }

    case "TokenBalance":
        if _, ok := task.Params["network"]; !ok {
            return errors.New("network不能为空")
        }
        if _, ok := task.Params["token_contract_addr"]; !ok {
            return errors.New("token_contract_addr不能为空")
        }
        if _, ok := task.Params["min_balance"]; !ok {
            return errors.New("min_balance不能为空")
        }
    }

    return nil
}
```

---

## 5. 常见问题

### Q1: params字段为什么要转为JSON字符串？
API设计要求params为字符串类型，这是为了保持灵活性，不同任务类型的参数结构不同。

### Q2: 如何处理时区？
建议统一使用UTC时间，前端显示时转换为用户本地时区。

### Q3: reward的token_address从哪里获取？
需要调用API获取支持的Token列表，根据chain和symbol查找对应的address。

---

*文档版本: 1.0.0*
*更新时间: 2026-02-03*
