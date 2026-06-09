# Design: 基于 Raft 的简化版配置管理中心

## 架构概览

系统采用 **Clean Architecture** 思路，将业务规则与技术细节分离，核心依赖方向始终保持从外层指向内层，避免业务逻辑直接依赖 HTTP、数据库或 Raft 的具体实现。

整体分层可概括为：

- **实体层（Entities / Domain）**：定义应用、环境、配置分组、配置项、版本、发布、灰度、订阅、事件等核心领域对象与领域规则。
- **用例层（Use Cases / Application）**：编排创建、更新、发布、回滚、灰度匹配、推送、订阅、恢复等业务流程，负责业务场景的输入输出转换。
- **接口适配层（Interface Adapters）**：负责 HTTP/API、WebSocket、Go SDK、后台页面、Raft command/event 之间的转换。
- **基础设施层（Frameworks & Drivers）**：负责 PostgreSQL、Raft 底座、Docker 部署、文件配置、日志等具体实现。

系统采用“管理后台 / Go SDK / HTTP API / 业务用例 / 领域模型 / Raft 一致性层 / 基础设施”的分层结构。

- 管理后台负责配置管理、发布、回滚与集群状态展示。
- Go SDK 负责客户端拉取配置、订阅变更、更新本地缓存。
- 业务用例负责应用、环境、配置分组、配置项、版本、发布与灰度逻辑。
- Raft 一致性层负责配置写操作在多节点间的顺序一致与多数派提交。
- 领域模型负责保持配置中心核心业务规则与状态约束。
- 基础设施层负责将已提交的命令落地到 PostgreSQL，并通过状态机将结果写回领域状态。

## 项目目录结构建议

参考业界常见的 Go 项目分层方式，并结合 Clean Architecture 的依赖倒置原则，本项目建议采用“实体层 + 用例层 + 接口适配层 + 基础设施层 + 独立 Raft 子模块”的结构，避免把 API、业务、持久化和共识实现混在一起。

建议目录结构如下：

```text
ai-coding-training/
├── openspec/
│   └── changes/config-center-raft/
│       ├── proposal.md
│       ├── specs/
│       ├── design.md
│       └── tasks.md
├── cmd/
│   ├── config-center-api/          # HTTP API 启动入口
│   ├── config-center-admin/        # 管理后台启动入口
│   └── config-center-sdk-demo/     # SDK 联调/示例入口（可选）
├── internal/
│   ├── entity/                     # 领域实体与领域规则（Entities）
│   ├── usecase/                    # 用例编排（Use Cases）
│   ├── adapter/
│   │   ├── http/                   # HTTP API handlers / middleware / router
│   │   ├── ws/                     # 订阅推送（如采用 WebSocket）
│   │   └── raft/                   # 业务命令与 Raft command/event 的适配
│   ├── infra/
│   │   ├── repository/             # PostgreSQL 仓储实现
│   │   ├── config/                 # 配置读取、环境变量、启动参数
│   │   ├── observability/          # 日志、指标、追踪（可按需启用）
│   │   └── auth/                   # 固定 token + hash 校验等基础鉴权实现
│   └── app/                        # 应用编排层，串联 usecase/adapter/infra
├── pkg/
│   └── sdk/                        # Go SDK 对外公共包（如需要独立发布）
├── migrations/                     # PostgreSQL DDL / migration 脚本
├── test/                           # 集成测试、端到端测试、测试数据
└── raft-stash/                     # 以子模块/子目录方式纳入的 Raft 源码
    └── src/
```

关于 `raft-stash` 的纳入方式，建议作为**项目内独立子模块或子目录**保留其原始工程边界，业务项目通过 `internal/raftadapter` 与其交互，而不是直接在业务代码中调用其内部实现。这样可以：

- 保持 Raft 课程工程的原始可运行性
- 让业务侧只依赖稳定接口（`Make` / `Start` / `GetState` / `ApplyMsg` / `Kill`）
- 方便后续替换共识实现而不影响业务层结构

## 接口边界与调用链

### Raft 层负责
- leader 选举
- 日志复制
- 多数派提交
- 持久化
- 崩溃恢复
- 通过 ApplyMsg 通知业务层已提交命令

### 业务层负责
- 请求鉴权
- 参数校验
- 业务语义解析
- 命令封装
- 幂等处理
- 版本控制
- 发布 / 回滚 / 灰度策略
- 对外响应生成

### 典型写入流程
1. 业务层接收配置变更请求。
2. 业务层完成校验和幂等检查。
3. 业务层将请求封装为不透明 command。
4. 调用 Raft Start(command)。
5. Raft 完成复制与提交。
6. Raft 通过 ApplyMsg 通知业务层。
7. 业务层应用命令并返回成功。

### 典型读取流程
1. 客户端请求配置。
2. 业务层读取本地状态机或缓存。
3. 必要时通过 GetState 判断当前节点是否 leader。
4. 若为 follower，则由上层做重定向或降级处理。

简化的数据流如下：

1. 管理后台或 SDK 发起配置相关请求。
2. 业务服务完成基础校验，并将写请求封装为 Raft command。
3. Leader 将 command 复制到多数节点并提交。
4. 状态机应用 command，更新配置版本、发布记录与灰度记录。
5. 配置变更生成标准化事件，Go SDK 收到事件后补偿拉取并更新本地缓存。
6. 管理后台查询集群状态、版本与发布结果。

## 模块划分

### Module 1: 域模型与配置管理
- 职责: 管理应用、环境、配置分组、配置项及其基础校验。
- 接口:
  - CreateApp / UpdateApp / DeleteApp / ListApp
  - CreateEnv / UpdateEnv / DeleteEnv / ListEnv
  - CreateGroup / UpdateGroup / DeleteGroup / ListGroup
  - CreateItem / UpdateItem / DeleteItem / GetItem / ListItem

### Module 2: 版本与发布管理
- 职责: 维护草稿版本、已发布版本、历史版本、发布记录与回滚记录。
- 接口:
  - SaveDraftVersion
  - ListVersions
  - PublishVersion
  - RollbackVersion
  - GetReleaseRecord

### Module 3: 灰度发布管理
- 职责: 维护灰度规则与命中结果，支持比例、标签和白名单三种简单策略。
- 接口:
  - CreateGrayRule
  - UpdateGrayRule
  - DeleteGrayRule
  - MatchGrayTarget
  - ListGrayRecords

### Module 4: 实时推送与客户端订阅
- 职责: 生成标准化变更事件，维护订阅关系，支持断线重连后的补偿拉取。
- 接口:
  - SubscribeConfig
  - UnsubscribeConfig
  - PublishConfigEvent
  - ReplayLatestEvent
  - GetChangeEvent

### Module 5: Go SDK
- 职责: 为客户端提供统一接入方式，完成拉取、订阅、缓存刷新和灰度结果判断。
- 接口:
  - NewClient
  - GetConfig
  - WatchConfig
  - RefreshConfig
  - IsGrayMatched

### Module 6: Raft 一致性层
- 职责: 处理 leader 选举、日志复制、多数派提交、提交结果回传、节点恢复与可选快照能力。
- 接口优先级:
  - 核心接口（实现优先级最高）:
    - SubmitCommand(command, timeout)
      - 用途: 提交写命令并等待结果或超时返回。
      - 说明: 仅 leader 接收写请求，非 leader 应返回重定向或重试提示。
    - GetState()
      - 用途: 查询当前节点的 term、leader 归属与角色状态。
    - GetLeader()
      - 用途: 查询当前集群 leader 标识或地址。
    - SubscribeApply()
      - 用途: 订阅已提交命令的 apply 结果，供业务层更新状态机。
    - Recover()
      - 用途: 从持久化日志和必要元数据恢复节点状态。
  - 次核心接口（实现时建议保留，先可占位）:
    - GetCommitIndex()
      - 用途: 查询当前已提交日志位置，用于同步进度展示。
  - 后置优化接口（本期可先占位，不作为主实现目标）:
    - Snapshot()
      - 用途: 对状态机生成快照，减少恢复成本。
    - Compact()
      - 用途: 清理已压缩日志，控制存储增长。
- 约束:
  - 写命令提交前必须形成多数派。
  - 提交成功后必须回传 apply 结果给业务层。
  - 非 leader 节点不得对外确认写成功。
  - 节点恢复期间不得对外提供写入能力。

### Module 7: 可视化管理后台
- 职责: 展示配置列表、版本历史、发布结果、灰度状态和节点状态。
- 接口:
  - Config list page
  - Version history page
  - Publish page
  - Rollback page
  - Cluster status page

## 接口草案

### 适用范围
接口草案覆盖本阶段会被业务调用、或对外提供稳定交互语义的模块，包括：
- 域模型与配置管理
- 版本与发布管理
- 灰度发布管理
- 实时推送与客户端订阅
- Go SDK
- Raft 一致性层

可视化管理后台属于页面层，主要定义页面入口与交互流程，不作为核心接口草案的重点。

### 说明
- 接口中的 `id` 字段优先对应数据模型主键。
- 接口中的 `name` 字段对应展示名或业务唯一名，由上层业务做唯一性校验。
- 涉及 `app_name/env_name/group_name` 的接口，建议在实现层映射到对应 `App/Environment/ConfigGroup` 主键。
- 涉及订阅、灰度命中记录、配置快照的接口，需要补充对应数据模型，见下方 `Subscription`、`GrayRecord`、`ConfigSnapshot`。

### Module 1: 域模型与配置管理接口草案
#### CreateApp
- 输入: `name`, `description`
- 输出: `App`
- 对应数据模型: `App`
- 约束: `name` 在租户或系统作用域内唯一，重复时返回业务错误

#### UpdateApp
- 输入: `app_id`, `name`, `description`
- 输出: `App`
- 对应数据模型: `App`
- 约束: `app_id` 必须存在，不允许修改为冲突名称

#### DeleteApp
- 输入: `app_id`
- 输出: `deleted`, `deleted_at`
- 对应数据模型: `App`
- 约束: 若仍有关联环境、分组或配置项，需先解除依赖

#### ListApp
- 输入: `page`, `page_size`, `filter`
- 输出: `items: []App`, `total`
- 对应数据模型: `App`

#### CreateEnv
- 输入: `name`, `description`
- 输出: `Environment`
- 对应数据模型: `Environment`
- 约束: `name` 保持唯一，重复时返回业务错误

#### UpdateEnv
- 输入: `env_id`, `name`, `description`
- 输出: `Environment`
- 对应数据模型: `Environment`

#### DeleteEnv
- 输入: `env_id`
- 输出: `deleted`, `deleted_at`
- 对应数据模型: `Environment`
- 约束: 删除需满足依赖检查

#### ListEnv
- 输入: `page`, `page_size`, `filter`
- 输出: `items: []Environment`, `total`
- 对应数据模型: `Environment`

#### CreateGroup
- 输入: `app_id`, `env_id`, `name`, `description`
- 输出: `ConfigGroup`
- 对应数据模型: `ConfigGroup`
- 约束: 同一作用域内名称唯一

#### UpdateGroup
- 输入: `group_id`, `name`, `description`
- 输出: `ConfigGroup`
- 对应数据模型: `ConfigGroup`

#### DeleteGroup
- 输入: `group_id`
- 输出: `deleted`, `deleted_at`
- 对应数据模型: `ConfigGroup`

#### ListGroup
- 输入: `app_id`, `env_id`, `page`, `page_size`
- 输出: `items: []ConfigGroup`, `total`
- 对应数据模型: `ConfigGroup`

#### CreateItem
- 输入: `group_id`, `key`, `value`, `status`
- 输出: `ConfigItem`
- 对应数据模型: `ConfigItem`
- 约束: `key` 在同一分组内唯一，写入前需校验格式与权限

#### UpdateItem
- 输入: `item_id`, `key`, `value`, `status`
- 输出: `ConfigItem`
- 对应数据模型: `ConfigItem`

#### DeleteItem
- 输入: `item_id`
- 输出: `deleted`, `deleted_at`
- 对应数据模型: `ConfigItem`

#### GetItem
- 输入: `item_id`
- 输出: `ConfigItem`
- 对应数据模型: `ConfigItem`

#### ListItem
- 输入: `group_id`, `page`, `page_size`, `filter`
- 输出: `items: []ConfigItem`, `total`
- 对应数据模型: `ConfigItem`

### Module 2: 版本与发布管理接口草案
#### SaveDraftVersion
- 输入: `group_id`, `draft_items: []ConfigItem`, `operator`
- 输出: `ConfigVersion`
- 对应数据模型: `ConfigVersion`
- 约束: 草稿保存不等于发布，不应改变当前生效版本

#### ListVersions
- 输入: `group_id`, `page`, `page_size`
- 输出: `items: []ConfigVersion`, `total`
- 对应数据模型: `ConfigVersion`

#### PublishVersion
- 输入: `version_id`, `publish_scope`, `operator`
- 输出: `ReleaseRecord`
- 对应数据模型: `ConfigVersion`, `ReleaseRecord`
- 约束: 必须经过业务校验和 Raft 提交后才可返回成功

#### RollbackVersion
- 输入: `version_id`, `target_version_id`, `operator`
- 输出: `ReleaseRecord`
- 对应数据模型: `ConfigVersion`, `ReleaseRecord`
- 约束: 回滚操作必须保持最终一致，不得产生部分生效

#### GetReleaseRecord
- 输入: `release_id`
- 输出: `ReleaseRecord`
- 对应数据模型: `ReleaseRecord`

### Module 3: 灰度发布管理接口草案
#### CreateGrayRule
- 输入: `group_id`, `rule_type`, `rule_value`, `target_scope`
- 输出: `GrayRule`
- 对应数据模型: `GrayRule`
- 约束: 仅支持比例、标签、白名单三类规则

#### UpdateGrayRule
- 输入: `gray_id`, `rule_type`, `rule_value`, `target_scope`, `status`
- 输出: `GrayRule`
- 对应数据模型: `GrayRule`

#### DeleteGrayRule
- 输入: `gray_id`
- 输出: `deleted`, `deleted_at`
- 对应数据模型: `GrayRule`

#### MatchGrayTarget
- 输入: `user_id`, `labels`, `env_context`
- 输出: `matched`, `gray_rule_id`, `matched_reason`
- 对应数据模型: `GrayRule`
- 约束: 仅做规则匹配，不负责发布动作

#### ListGrayRecords
- 输入: `group_id`, `page`, `page_size`
- 输出: `items: []GrayRecord`, `total`
- 对应数据模型: `GrayRecord`

### Module 4: 实时推送与客户端订阅接口草案
#### SubscribeConfig
- 输入: `app_id`, `env_id`, `group_id`, `client_id`
- 输出: `Subscription`
- 对应数据模型: `Subscription`
- 约束: 订阅建立后应能接收后续变更事件

#### UnsubscribeConfig
- 输入: `subscription_id`
- 输出: `deleted`, `deleted_at`
- 对应数据模型: `Subscription`

#### PublishConfigEvent
- 输入: `event_type`, `app_id`, `env_id`, `group_id`, `version_id`, `version_no`
- 输出: `ChangeEvent`
- 对应数据模型: `ChangeEvent`
- 约束: 事件应可幂等投递，重复事件不得改变最终状态

#### ReplayLatestEvent
- 输入: `subscription_id`, `checkpoint`
- 输出: `items: []ChangeEvent`, `next_checkpoint`
- 对应数据模型: `Subscription`, `ChangeEvent`

#### GetChangeEvent
- 输入: `event_id`
- 输出: `ChangeEvent`
- 对应数据模型: `ChangeEvent`

### Module 5: Go SDK 接口草案
#### NewClient
- 输入: `server_addr`, `timeout`, `cache_config`
- 输出: `Client`
- 对应数据模型: 无

#### GetConfig
- 输入: `app_id`, `env_id`, `group_id`, `key`
- 输出: `value`, `version_no`, `matched_gray_rule_id`, `matched`
- 对应数据模型: `ConfigItem`, `ConfigVersion`, `GrayRule`

#### WatchConfig
- 输入: `app_id`, `env_id`, `group_id`, `callback`
- 输出: `watch_handle`
- 对应数据模型: `Subscription`, `ChangeEvent`
- 约束: 支持变更通知与断线后的补偿拉取

#### RefreshConfig
- 输入: `app_id`, `env_id`, `group_id`
- 输出: `ConfigSnapshot`
- 对应数据模型: `ConfigSnapshot`

#### IsGrayMatched
- 输入: `user_context`, `gray_rule`
- 输出: `matched`
- 对应数据模型: `GrayRule`

### Module 6: Raft 一致性层接口草案
#### Start(command)
- 输入: 不透明业务命令
- 输出: `index`, `term`, `isLeader`
- 对应数据模型: 无
- 约束: 仅负责命令复制，不解析业务语义

#### GetState()
- 输入: 无
- 输出: `term`, `isLeader`
- 对应数据模型: 无
- 约束: 供上层路由与写入判断使用

#### ApplyMsg
- 输入: 已提交日志或快照信息
- 输出: 上层业务回调触发
- 对应数据模型: 无
- 约束: Raft 仅投递，不直接修改业务状态

#### RecoverFromLog
- 输入: 持久化日志与快照
- 输出: 恢复后的运行状态
- 对应数据模型: 无
- 约束: 恢复完成前不得对外宣布可写

## 数据模型

### 输入模型
- `CreateAppInput`
  - `name (text)`: 应用展示名/业务唯一名
  - `description (text)`: 应用描述
- `UpdateAppInput`
  - `app_id (uuid)`: 应用主键
  - `name (text)`: 应用展示名/业务唯一名
  - `description (text)`: 应用描述
- `ListAppQuery`
  - `page (int)`: 页码
  - `page_size (int)`: 每页数量
  - `filter (text/json)`: 过滤条件
- `CreateEnvInput`
  - `name (text)`: 环境展示名/业务唯一名
  - `description (text)`: 环境描述
- `UpdateEnvInput`
  - `env_id (uuid)`: 环境主键
  - `name (text)`: 环境展示名/业务唯一名
  - `description (text)`: 环境描述
- `ListEnvQuery`
  - `page (int)`: 页码
  - `page_size (int)`: 每页数量
  - `filter (text/json)`: 过滤条件
- `CreateGroupInput`
  - `app_id (uuid)`: 所属应用 ID
  - `env_id (uuid)`: 所属环境 ID
  - `name (text)`: 分组展示名/业务唯一名
  - `description (text)`: 分组描述
- `UpdateGroupInput`
  - `group_id (uuid)`: 分组主键
  - `name (text)`: 分组展示名/业务唯一名
  - `description (text)`: 分组描述
- `ListGroupQuery`
  - `app_id (uuid)`: 过滤的应用 ID
  - `env_id (uuid)`: 过滤的环境 ID
  - `page (int)`: 页码
  - `page_size (int)`: 每页数量
- `CreateItemInput`
  - `group_id (uuid)`: 所属分组 ID
  - `key (text)`: 配置键
  - `value (text)`: 配置值
  - `status (smallint)`: 配置状态
- `UpdateItemInput`
  - `item_id (uuid)`: 配置项主键
  - `key (text)`: 配置键
  - `value (text)`: 配置值
  - `status (smallint)`: 配置状态
- `ListItemQuery`
  - `group_id (uuid)`: 所属分组 ID
  - `page (int)`: 页码
  - `page_size (int)`: 每页数量
  - `filter (text/json)`: 过滤条件
- `SaveDraftVersionInput`
  - `group_id (uuid)`: 所属分组 ID
  - `draft_items ([]ConfigItemInput/json)`: 草稿配置集合
  - `operator (text)`: 操作人
- `ListVersionsQuery`
  - `group_id (uuid)`: 所属分组 ID
  - `page (int)`: 页码
  - `page_size (int)`: 每页数量
- `PublishVersionInput`
  - `version_id (uuid)`: 版本记录主键
  - `publish_scope (jsonb)`: 发布范围
  - `operator (text)`: 操作人
- `RollbackVersionInput`
  - `version_id (uuid)`: 当前版本主键
  - `target_version_id (uuid)`: 目标回滚版本主键
  - `operator (text)`: 操作人
- `GetReleaseRecordQuery`
  - `release_id (uuid)`: 发布记录主键
- `CreateGrayRuleInput`
  - `group_id (uuid)`: 所属分组 ID
  - `rule_type (text)`: 规则类型
  - `rule_value (jsonb)`: 规则值
  - `target_scope (jsonb)`: 命中目标范围
- `UpdateGrayRuleInput`
  - `gray_id (uuid)`: 灰度规则主键
  - `rule_type (text)`: 规则类型
  - `rule_value (jsonb)`: 规则值
  - `target_scope (jsonb)`: 命中目标范围
  - `status (smallint)`: 规则状态
- `ListGrayRecordsQuery`
  - `group_id (uuid)`: 所属分组 ID
  - `page (int)`: 页码
  - `page_size (int)`: 每页数量
- `SubscribeConfigInput`
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
  - `client_id (text)`: 客户端标识
- `UnsubscribeConfigInput`
  - `subscription_id (uuid)`: 订阅主键
- `PublishConfigEventInput`
  - `event_type (text)`: 事件类型
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
  - `version_id (uuid)`: 版本主键
  - `version_no (bigint)`: 业务版本号
- `ReplayLatestEventQuery`
  - `subscription_id (uuid)`: 订阅主键
  - `checkpoint (bigint)`: 补偿位点
- `GetChangeEventQuery`
  - `event_id (uuid)`: 事件主键
- `GetConfigQuery`
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
  - `key (text)`: 配置键
- `WatchConfigInput`
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
  - `callback (function/handler)`: 变更回调
- `RefreshConfigQuery`
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
- `IsGrayMatchedInput`
  - `user_context (json/text)`: 用户上下文
  - `gray_rule (GrayRule/DTO)`: 灰度规则

### 输出模型
- `AppOutput`
  - `app_id (uuid)`: 应用主键
  - `name (text)`: 应用展示名/业务唯一名
  - `description (text)`: 应用描述
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `EnvironmentOutput`
  - `env_id (uuid)`: 环境主键
  - `name (text)`: 环境展示名/业务唯一名
  - `description (text)`: 环境描述
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `ConfigGroupOutput`
  - `group_id (uuid)`: 分组主键
  - `app_id (uuid)`: 所属应用 ID
  - `env_id (uuid)`: 所属环境 ID
  - `name (text)`: 分组展示名/业务唯一名
  - `description (text)`: 分组描述
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `ConfigItemOutput`
  - `item_id (uuid)`: 配置项主键
  - `group_id (uuid)`: 所属分组 ID
  - `key (text)`: 配置键
  - `value (text)`: 配置值
  - `status (smallint)`: 配置状态
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `ConfigVersionOutput`
  - `version_id (uuid)`: 版本记录主键
  - `group_id (uuid)`: 所属分组 ID
  - `version_no (bigint)`: 业务版本号/展示版本号
  - `draft_flag (boolean)`: 是否草稿版本
  - `status (smallint)`: 版本状态
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `ReleaseRecordOutput`
  - `release_id (uuid)`: 发布记录主键
  - `group_id (uuid)`: 所属分组 ID
  - `version_id (uuid)`: 发布版本 ID
  - `release_type (text)`: 发布类型
  - `publish_scope (jsonb)`: 发布范围
  - `status (smallint)`: 发布状态
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `GrayRuleOutput`
  - `gray_id (uuid)`: 灰度规则主键
  - `group_id (uuid)`: 所属分组 ID
  - `rule_type (text)`: 规则类型
  - `rule_value (jsonb)`: 规则值
  - `target_scope (jsonb)`: 命中目标范围
  - `status (smallint)`: 规则状态
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `SubscriptionOutput`
  - `subscription_id (uuid)`: 订阅主键
  - `client_id (text)`: 客户端标识
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
  - `checkpoint (bigint)`: 订阅位点/游标
  - `status (smallint)`: 订阅状态
  - `created_at (timestamptz)`: 创建时间
  - `updated_at (timestamptz)`: 更新时间
- `GrayRecordOutput`
  - `record_id (uuid)`: 灰度命中记录主键
  - `group_id (uuid)`: 所属分组 ID
  - `user_id (text)`: 用户标识
  - `matched_rule_id (uuid)`: 命中的灰度规则 ID
  - `matched_at (timestamptz)`: 命中时间
  - `result (boolean)`: 命中结果
- `ConfigSnapshotOutput`
  - `snapshot_id (uuid)`: 快照主键
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
  - `items (jsonb)`: 快照中的配置项集合
  - `version_id (uuid)`: 对应版本 ID
  - `created_at (timestamptz)`: 创建时间
- `ChangeEventOutput`
  - `event_id (uuid)`: 事件主键
  - `app_id (uuid)`: 应用 ID
  - `env_id (uuid)`: 环境 ID
  - `group_id (uuid)`: 分组 ID
  - `version_id (uuid)`: 对应版本 ID
  - `version_no (bigint)`: 业务版本号
  - `event_type (text)`: 事件类型
  - `gray_flag (boolean)`: 是否灰度相关
  - `created_at (timestamptz)`: 创建时间
- `ClientOutput`
  - `server_addr (text)`: 服务端地址
  - `timeout (duration/int)`: 超时时间
  - `cache_config (json/text)`: 缓存配置
- `WatchHandleOutput`
  - `watch_handle (string/uuid)`: 监听句柄
- `ConfigValueOutput`
  - `value (text)`: 配置值
  - `version_no (bigint)`: 业务版本号
  - `matched_gray_rule_id (uuid)`: 命中的灰度规则 ID
  - `matched (boolean)`: 是否命中
- `ListOutput<T>`
  - `items ([]T)`: 列表数据
  - `total (int64)`: 总数
- `DeleteOutput`
  - `deleted (boolean)`: 是否删除成功
  - `deleted_at (timestamptz)`: 删除时间
- `BoolOutput`
  - `matched (boolean)`: 结果布尔值

### 数据库模型
- `app_row`: `app_id (uuid)`, `name (text)`, `description (text)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `environment_row`: `env_id (uuid)`, `name (text)`, `description (text)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `config_group_row`: `group_id (uuid)`, `app_id (uuid)`, `env_id (uuid)`, `name (text)`, `description (text)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `config_item_row`: `item_id (uuid)`, `group_id (uuid)`, `key (text)`, `value (text)`, `status (smallint)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `config_version_row`: `version_id (uuid)`, `group_id (uuid)`, `version_no (bigint)`, `draft_flag (boolean)`, `status (smallint)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `release_record_row`: `release_id (uuid)`, `group_id (uuid)`, `version_id (uuid)`, `release_type (text)`, `publish_scope (jsonb)`, `status (smallint)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `gray_rule_row`: `gray_id (uuid)`, `group_id (uuid)`, `rule_type (text)`, `rule_value (jsonb)`, `target_scope (jsonb)`, `status (smallint)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `subscription_row`: `subscription_id (uuid)`, `client_id (text)`, `app_id (uuid)`, `env_id (uuid)`, `group_id (uuid)`, `checkpoint (bigint)`, `status (smallint)`, `created_at (timestamptz)`, `updated_at (timestamptz)`
- `gray_record_row`: `record_id (uuid)`, `group_id (uuid)`, `user_id (text)`, `matched_rule_id (uuid)`, `matched_at (timestamptz)`, `result (boolean)`
- `config_snapshot_row`: `snapshot_id (uuid)`, `app_id (uuid)`, `env_id (uuid)`, `group_id (uuid)`, `items (jsonb)`, `version_id (uuid)`, `created_at (timestamptz)`
- `change_event_row`: `event_id (uuid)`, `app_id (uuid)`, `env_id (uuid)`, `group_id (uuid)`, `version_id (uuid)`, `version_no (bigint)`, `event_type (text)`, `gray_flag (boolean)`, `created_at (timestamptz)`

### 查询模型
- `AppQuery`: `page`, `page_size`, `name_like`, `status`, `created_after`, `created_before`
- `EnvironmentQuery`: `page`, `page_size`, `name_like`, `status`
- `ConfigGroupQuery`: `app_id`, `env_id`, `page`, `page_size`, `name_like`
- `ConfigItemQuery`: `group_id`, `page`, `page_size`, `key_like`, `status`
- `ConfigVersionQuery`: `group_id`, `page`, `page_size`, `status`, `version_no_from`, `version_no_to`
- `ReleaseRecordQuery`: `group_id`, `page`, `page_size`, `status`, `release_type`
- `GrayRuleQuery`: `group_id`, `page`, `page_size`, `status`, `rule_type`
- `GrayRecordQuery`: `group_id`, `page`, `page_size`, `user_id`, `matched_rule_id`
- `SubscriptionQuery`: `app_id`, `env_id`, `group_id`, `page`, `page_size`, `status`
- `ChangeEventQuery`: `app_id`, `env_id`, `group_id`, `page`, `page_size`, `version_no_from`, `version_no_to`
- `ConfigSnapshotQuery`: `app_id`, `env_id`, `group_id`, `version_id`

## 模型分层约束

- Input Model 仅用于接口入参绑定、参数校验与调用上下文承载，不直接用于持久化。
- Output Model 仅用于接口响应组装，不直接暴露数据库行结构。
- DB Model 仅用于 PostgreSQL 持久化映射，不直接作为 API 返回结构。
- Query Model 仅用于 repository/DAO 查询条件构造，不承担业务状态表达。
- 各层模型之间应通过显式转换函数或映射器完成转换，避免跨层直接复用。

## 命令与事件模型

- Command Model 表示业务请求封装后的 Raft 复制命令，进入 Raft 前必须保持不透明，不解析业务语义。
- Event Model 表示 Raft 提交后用于驱动状态机、推送与审计的业务事件。
- API Input 不应直接作为 Raft command；DB Model 不应直接作为 Event 输出；ApplyMsg 不应混入接口入参结构。
- 状态机应用与消息推送应基于 Event Model 进行统一处理，保证写路径语义一致。

## 状态枚举约定

- 各模型中的 `status` 字段应使用领域层统一枚举管理，禁止在业务代码中散落魔法数字。
- 枚举值应结合实体语义定义，例如草稿、已发布、已回滚、已删除、已禁用、已失效等。
- 若不同实体的状态集合不同，应在模型注释或领域定义中分别说明，不强制共享同一套枚举值。

## 版本号生成规则

- `version_id` 使用 UUID 作为版本记录主键，用于精确定位与关联。
- `version_no` 使用 `bigint`，并在 `group_id` 维度内单调递增，用于展示、排序与事件追踪。
- `version_no` 由业务层在事务内生成，禁止回退复用。
- 回滚或重发布不应修改历史 `version_no`，只应生成新的发布记录或事件。

## 数据库约束原则

- `app_row.name` 在同一租户或系统作用域内应唯一。
- `config_group_row (app_id, env_id, name)` 应唯一。
- `config_item_row (group_id, key)` 应唯一。
- `config_version_row (group_id, version_no)` 应唯一。
- `release_record_row.version_id` 应引用有效版本记录。
- `gray_rule_row.group_id`、`subscription_row (app_id, env_id, group_id)` 等外键关系应保持完整。
- 关键查询字段应建立必要索引，例如 `group_id`、`version_id`、`subscription_id`、`event_id`、`created_at`。

## 技术选型说明

- 选择 Go 作为统一实现语言，原因是其并发模型适合网络服务和 Raft 集群实现，且便于编写 Go SDK。
- 选择现有 Raft 课程工程作为一致性底座，原因是可以显著降低共识协议实现风险，将精力集中在配置中心业务。
- 选择 PostgreSQL 作为统一数据库，原因是其事务能力、约束能力、JSONB 与 UUID 支持都适合配置中心业务建模。
- 选择 `uuid-ossp` 或 `pgcrypto` 生成 UUID，原因是 PostgreSQL 原生支持 UUID 主键生成，适合分布式场景且避免中心化自增瓶颈。
- 选择 HTTP 作为管理接口，原因是简单直观，便于后台与 SDK 调用。
- 选择 WebSocket 或长连接式订阅作为推送实现，原因是能够满足实时变更通知需求，同时保持实现复杂度可控。
- 选择纯软件多节点模拟部署，原因是符合课题约束，无需额外硬件器件。

## 部署方案

系统采用 Docker + Docker Compose 进行本地与课题环境部署，所有服务运行在同一套虚拟网络中，避免宿主机端口与容器内部地址混淆。

### 镜像版本建议
- PostgreSQL：`postgres:15` 或 `pgvector/pgvector:pg15`
- Go 服务镜像：使用项目自定义构建镜像，统一以 `config-center:*` 命名
- 前端管理后台镜像：使用项目自定义构建镜像，统一以 `config-center-admin:*` 命名
- Raft 底座：作为源码子目录或子模块随项目一起构建，不单独依赖外部镜像

### Docker 网络建议
- 使用独立的 bridge 网络，例如 `config-center-network`
- 业务服务、后台、数据库与 Raft 节点应加入同一网络，便于容器间通过服务名互访
- 对外仅暴露必要端口，例如 API 端口、后台端口与数据库映射端口

### 数据库连接约定
- 数据库 host、port、user、password、dbname、sslmode 等信息应通过配置文件或环境变量注入。
- 本项目的 PostgreSQL 连接参数应在 `configs/config.yaml` 中集中配置，并在 Docker Compose 中通过环境变量或配置挂载方式传入。
- 配置文件采用“示例版 + 使用版”双文件模式：`configs/config.example.yaml` 用于提交到仓库，`configs/config.yaml` 用于本地或容器运行时使用。
- 使用版配置文件应加入 `.gitignore`，避免真实连接信息、token 或密钥被提交到版本库。
- 配置文件中涉及敏感信息的字段在仓库内应使用示例值或占位符，真实值通过本地环境注入，不应硬编码入文档或代码。

### 管理后台鉴权约定
- 管理后台采用固定管理员 token 作为最小鉴权方案。
- 服务端仅保存 token 的 hash 值，不保存明文 token。
- 请求进入后台时，服务端对输入 token 进行同样的 hash 计算后比对，通过后方可访问管理接口。
- 该方案仅用于本课题的最小可用鉴权，不引入多用户体系、RBAC 或复杂登录流程。

### 灰度与事件语义冻结
- 灰度规则优先级固定为：白名单 > 标签 > 比例；同一客户端同时命中多个规则时，仅采用最高优先级规则。
- 变更事件应包含唯一 `event_id`，客户端以 `event_id` 作为幂等去重主键，`checkpoint` 仅用于补偿拉取游标。
- `version_id` 使用 UUID 作为精确关联主键，`version_no` 使用 `bigint` 且在分组维度内单调递增，用于展示、排序与事件追踪。
- `ChangeEvent` 仅负责通知“发生了什么变化”，`ConfigSnapshot` 仅负责提供“当前完整配置视图”，两者职责不得混用。

## 设计约束

- Raft 只负责写入一致性，不负责灰度计算、SDK 缓存策略或前端展示逻辑。
- 灰度发布仅支持三种简单策略：比例、标签、白名单。
- Go SDK 为唯一正式实现的客户端 SDK。
- Prometheus 与更复杂的前端生态集成不纳入本阶段范围。
- 发布失败不得产生部分生效结果。
