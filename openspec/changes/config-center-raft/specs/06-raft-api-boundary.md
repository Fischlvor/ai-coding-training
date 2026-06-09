### Requirement: Raft 与业务职责边界清晰
系统 SHALL 将共识能力与业务能力分离。Raft 层只负责 leader 选举、日志复制、提交与应用通知；配置管理业务层负责请求校验、业务语义、版本控制、幂等处理与发布策略。Raft 多节点部署 SHALL 采用同一代码、同一容器内端口、不同节点配置与 Docker service-name DNS 的方式组织 peer 地址，避免将宿主机映射端口作为节点间互联的唯一依据。

#### Scenario: 业务层提交配置写请求
- GIVEN 上层配置管理服务收到新增或修改配置请求
- WHEN 请求通过业务校验
- THEN 业务层 SHALL 将请求封装为 Raft 可复制的通用命令
- AND Raft 层 SHALL 仅负责复制、提交与应用通知

#### Scenario: Raft 不处理业务语义
- GIVEN 写入命令包含配置项、版本号、灰度标记等字段
- WHEN 命令进入 Raft 层
- THEN Raft SHALL 将其视为不透明命令
- AND Raft SHALL NOT 解析、校验或修改业务字段

### Requirement: Raft 对上层暴露最小必要接口
系统 SHALL 向上层暴露用于状态查询与命令复制的最小接口集，避免将成员管理、发布策略、配置校验等业务逻辑纳入 Raft API。

#### Scenario: 查询当前节点状态
- GIVEN 上层需要判断当前节点是否可写
- WHEN 调用状态查询接口
- THEN Raft SHALL 返回当前任期和是否为 leader
- AND 上层 SHALL 根据结果决定请求路由

#### Scenario: 发起命令复制
- GIVEN 上层准备提交一条业务命令
- WHEN 调用写入接口
- THEN Raft SHALL 返回日志索引、任期和是否为 leader
- AND 上层 SHALL 根据返回值决定是否等待提交或重试

#### Scenario: 应用结果通知
- GIVEN Raft 产生已提交条目
- WHEN Raft 向应用层发送 ApplyMsg
- THEN 应用层 SHALL 负责将该条目转换为业务动作
- AND Raft SHALL NOT 直接操作业务存储

### Requirement: 上层业务负责请求校验与幂等
系统 SHALL 在配置管理业务层完成参数校验、权限判断、幂等控制和版本控制，Raft 层不承担这些职责。

#### Scenario: 非法请求被拒绝
- GIVEN 客户端提交格式非法或缺少必要字段的请求
- WHEN 请求到达业务层
- THEN 业务层 SHALL 在进入 Raft 前直接拒绝
- AND Raft 层 SHALL 不接收该请求

#### Scenario: 重复请求处理
- GIVEN 同一业务请求被重复提交
- WHEN 业务层识别到重复请求
- THEN 业务层 SHALL 基于请求 ID 或版本号做幂等处理
- AND Raft 层 SHALL 仅处理去重后的有效命令

#### Scenario: 版本冲突
- GIVEN 客户端提交的版本号与当前生效版本不一致
- WHEN 业务层检测到冲突
- THEN 业务层 SHALL 拒绝该请求
- AND Raft 层 SHALL 不参与版本冲突判断

### Requirement: 配置写入必须经过多数派提交
系统 SHALL 保证配置写入在多数派提交后才可对外确认成功。

#### Scenario: 写入成功
- GIVEN 当前节点为 leader 且多数派可用
- WHEN 业务层提交配置变更命令
- THEN Raft SHALL 复制并提交该命令
- AND 业务层 SHALL 在收到应用通知后返回成功

#### Scenario: 多数派不可用
- GIVEN 集群无法形成多数派
- WHEN 业务层提交配置变更命令
- THEN Raft SHALL 不得返回已提交成功
- AND 业务层 SHALL 向客户端返回重试或失败

### Requirement: 节点恢复后继续参与一致性复制
系统 SHALL 在节点重启后恢复 Raft 与配置状态。

#### Scenario: 节点重启恢复
- GIVEN 节点已完成一次或多次配置提交
- WHEN 节点重启并重新加入集群
- THEN 系统 SHALL 从持久化数据中恢复必要状态
- AND 节点 SHALL 继续参与一致性复制

#### Scenario: 恢复期间禁止写入
- GIVEN 节点正在恢复但尚未完成同步
- WHEN 客户端提交配置写请求
- THEN 系统 SHALL 不将该节点视为可用写入节点
- AND 客户端 SHALL 继续路由到可用 leader

### Requirement: 接口易用且稳定
系统 SHOULD 提供便于上层集成的 Raft 接口抽象，避免上层直接依赖内部字段或日志结构。

#### Scenario: 上层仅依赖公开接口
- GIVEN 配置管理服务接入 Raft
- WHEN 上层实现写入与状态查询
- THEN 上层 SHALL 仅依赖公开方法和 ApplyMsg
- AND 上层 SHALL NOT 依赖 Raft 内部字段如日志数组、currentTerm、votedFor

#### Scenario: 后续扩展快照
- GIVEN 未来需要支持日志压缩与快照
- WHEN Raft 增加快照相关能力
- THEN 上层接口 SHALL 保持兼容
- AND 业务层 SHALL 仅补充快照恢复处理逻辑
