### Requirement: Raft 写入一致性
系统 SHALL 使用 Raft 保证配置写操作在多节点间的一致性。

#### Scenario: 提交配置写请求
- GIVEN 集群中存在可用 leader
- WHEN 客户端提交配置新增或修改请求
- THEN leader SHALL 将命令复制给多数节点
- AND 命令提交后 SHALL 被状态机应用

#### Scenario: Leader 不可用
- GIVEN 当前 leader 已失联
- WHEN 客户端提交配置写请求
- THEN 系统 SHALL 触发新的 leader 选举
- AND 在新的 leader 可用前拒绝写入或返回重试提示

### Requirement: 节点重启恢复
系统 SHALL 在节点重启后恢复 Raft 与配置状态。

#### Scenario: 节点重启后恢复
- GIVEN 节点已完成一次或多次配置提交
- WHEN 节点重启并重新加入集群
- THEN 系统 SHALL 从持久化数据中恢复必要状态
- AND 节点 SHALL 继续参与一致性复制

#### Scenario: 恢复失败
- GIVEN 节点持久化数据损坏或缺失
- WHEN 节点尝试重启恢复
- THEN 系统 SHALL 返回恢复失败状态
- AND 不应将该节点视为可用写入节点