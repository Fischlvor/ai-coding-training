### Requirement: 灰度发布策略
系统 SHALL 支持基础灰度发布策略，并使用统一的数据格式描述灰度规则。系统 MUST 对同一发布任务内的灰度规则冲突进行明确处理。

#### Scenario: 按比例灰度
- GIVEN 目标配置版本已存在
- WHEN 用户指定灰度比例
- THEN 系统 SHALL 仅让部分目标实例生效新版本
- AND 其余实例 SHALL 继续使用旧版本

#### Scenario: 按 0% 灰度
- GIVEN 目标配置版本已存在
- WHEN 用户指定灰度比例为 0%
- THEN 系统 SHALL 不向任何实例下发新版本
- AND 所有实例 SHALL 继续使用旧版本

#### Scenario: 按 100% 灰度
- GIVEN 目标配置版本已存在
- WHEN 用户指定灰度比例为 100%
- THEN 系统 SHALL 向全部目标实例下发新版本
- AND 所有实例 SHALL 视为命中新版本

#### Scenario: 按标签灰度
- GIVEN 目标实例已存在标签信息
- WHEN 用户指定标签条件
- THEN 系统 SHALL 仅对匹配标签的实例下发新版本
- AND 未匹配实例 SHALL 保持原版本

#### Scenario: 按实例白名单灰度
- GIVEN 目标实例列表已知
- WHEN 用户提交实例白名单
- THEN 系统 SHALL 仅对名单内实例下发新版本
- AND 名单外实例 SHALL 保持原版本

#### Scenario: 灰度规则冲突
- GIVEN 用户同时提交比例、标签和实例白名单
- WHEN 系统校验灰度请求
- THEN 系统 SHALL 按预定义优先级处理或拒绝该请求
- AND 返回规则冲突说明

#### Scenario: 灰度参数无效
- GIVEN 用户提交的比例、标签条件或实例名单非法
- WHEN 系统校验灰度请求
- THEN 系统 SHALL 拒绝该请求
- AND 返回校验失败原因

### Requirement: 灰度结果可追踪
系统 SHALL 记录灰度发布的执行结果与命中范围，并允许复核历史记录。

#### Scenario: 灰度发布成功
- GIVEN 灰度发布已执行完成
- WHEN 用户查询发布记录
- THEN 系统 SHALL 返回灰度范围、目标版本、命中规则与执行状态
- AND 记录 SHALL 可用于后续回滚或复核

#### Scenario: 灰度发布中断
- GIVEN 灰度发布过程中出现异常
- WHEN 系统停止当前发布流程
- THEN 系统 SHALL 记录失败状态
- AND 系统 SHALL 保留已完成的部分结果