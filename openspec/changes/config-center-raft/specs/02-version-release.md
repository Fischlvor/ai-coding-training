### Requirement: 配置版本历史
系统 SHALL 为每次配置变更保留版本历史。

#### Scenario: 生成新版本
- GIVEN 配置项已被修改
- WHEN 用户执行发布操作
- THEN 系统 SHALL 生成新的版本号
- AND 系统 SHALL 保留旧版本记录

#### Scenario: 查询历史版本
- GIVEN 某配置分组存在多个历史版本
- WHEN 用户查询版本列表
- THEN 系统 SHALL 返回按时间或版本号排序的历史记录
- AND 每条记录 SHALL 包含版本标识与发布时间

### Requirement: 配置发布
系统 SHALL 支持将指定版本发布为当前生效版本。

#### Scenario: 发布指定版本
- GIVEN 目标版本存在且可用
- WHEN 用户提交发布请求
- THEN 系统 SHALL 将该版本设为当前生效版本
- AND 系统 SHALL 记录发布操作

#### Scenario: 发布不存在版本
- GIVEN 目标版本不存在
- WHEN 用户提交发布请求
- THEN 系统 SHALL 拒绝该请求
- AND 返回明确的错误提示

### Requirement: 配置回滚
系统 SHALL 支持将历史版本回滚为当前生效版本。

#### Scenario: 回滚到历史版本
- GIVEN 历史版本存在
- WHEN 用户提交回滚请求
- THEN 系统 SHALL 将目标历史版本设为当前生效版本
- AND 系统 SHALL 生成新的发布记录

#### Scenario: 回滚到无效版本
- GIVEN 目标版本不存在或已失效
- WHEN 用户提交回滚请求
- THEN 系统 SHALL 拒绝该请求
- AND 系统 SHALL 保持当前生效版本不变