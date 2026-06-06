### Requirement: 配置版本历史
系统 SHALL 为每次配置变更保留版本历史，并区分草稿版本与已发布版本。

#### Scenario: 生成新草稿版本
- GIVEN 配置项已被修改
- WHEN 用户保存变更但未执行发布
- THEN 系统 SHALL 生成新的草稿版本
- AND 旧版本记录 SHALL 保持不变

#### Scenario: 重新保存草稿版本
- GIVEN 已存在未发布的草稿版本
- WHEN 用户再次修改并保存
- THEN 系统 SHALL 更新现有草稿版本或生成新的草稿版本号
- AND 系统 SHOULD 保持版本变更可追踪

#### Scenario: 查询历史版本
- GIVEN 某配置分组存在多个历史版本
- WHEN 用户查询版本列表
- THEN 系统 SHALL 返回按时间或版本号排序的历史记录
- AND 每条记录 SHALL 包含版本标识、发布时间与状态

#### Scenario: 查询空历史版本
- GIVEN 某配置分组尚未产生任何版本
- WHEN 用户查询版本列表
- THEN 系统 SHALL 返回空结果
- AND 不应返回错误

### Requirement: 配置发布
系统 SHALL 支持将指定版本发布为当前生效版本，并记录发布范围。系统 MUST 在发布前校验目标版本合法性。

#### Scenario: 发布指定版本
- GIVEN 目标版本存在且可用
- WHEN 用户提交发布请求
- THEN 系统 SHALL 将该版本设为当前生效版本
- AND 系统 SHALL 记录发布操作与发布范围

#### Scenario: 发布不存在版本
- GIVEN 目标版本不存在
- WHEN 用户提交发布请求
- THEN 系统 SHALL 拒绝该请求
- AND 返回明确的错误提示

#### Scenario: 重复发布当前版本
- GIVEN 目标版本已是当前生效版本
- WHEN 用户再次提交相同发布请求
- THEN 系统 SHALL 将其视为幂等操作或返回已生效提示
- AND 当前生效版本 SHALL 保持不变

#### Scenario: 发布前校验失败
- GIVEN 目标版本存在但内容不完整或不合法
- WHEN 用户提交发布请求
- THEN 系统 SHALL 拒绝该请求
- AND 当前生效版本 SHALL 保持不变

### Requirement: 配置回滚
系统 SHALL 支持将历史版本回滚为当前生效版本，并生成新的发布记录。系统 SHOULD 将回滚视为一次新的发布操作。

#### Scenario: 回滚到历史版本
- GIVEN 历史版本存在
- WHEN 用户提交回滚请求
- THEN 系统 SHALL 将目标历史版本设为当前生效版本
- AND 系统 SHALL 生成新的发布记录

#### Scenario: 回滚到当前版本
- GIVEN 目标版本已是当前生效版本
- WHEN 用户提交回滚请求
- THEN 系统 SHALL 返回已生效提示或将其视为幂等操作
- AND 当前生效版本 SHALL 保持不变

#### Scenario: 回滚到无效版本
- GIVEN 目标版本不存在或已失效
- WHEN 用户提交回滚请求
- THEN 系统 SHALL 拒绝该请求
- AND 系统 SHALL 保持当前生效版本不变