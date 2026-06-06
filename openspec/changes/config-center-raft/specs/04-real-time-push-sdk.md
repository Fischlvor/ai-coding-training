### Requirement: 实时配置推送
系统 SHALL 在配置变更后向已订阅客户端推送通知。

#### Scenario: 配置变更推送
- GIVEN 客户端已订阅某配置分组
- WHEN 该分组配置发生发布或回滚
- THEN 系统 SHALL 向客户端发送变更通知
- AND 通知 SHALL 包含新的版本标识

#### Scenario: 客户端断开连接
- GIVEN 客户端订阅通道已断开
- WHEN 配置发生变化
- THEN 系统 SHALL 记录推送失败或断开状态
- AND 不应影响其他正常客户端

### Requirement: SDK 接入能力
系统 SHALL 提供统一的 SDK 接入方式，以便客户端拉取和订阅配置。

#### Scenario: SDK 拉取配置
- GIVEN SDK 已完成初始化
- WHEN 客户端请求某应用配置
- THEN SDK SHALL 调用服务端接口并返回当前生效配置
- AND 返回结果 SHALL 包含版本信息

#### Scenario: SDK 订阅配置变化
- GIVEN SDK 已完成订阅注册
- WHEN 服务端推送配置变化
- THEN SDK SHALL 接收并更新本地缓存
- AND 客户端 SHALL 能读取到最新版本配置

### Requirement: 可视化管理后台
系统 SHALL 提供可视化管理后台用于配置管理与集群状态查看。

#### Scenario: 后台查看配置列表
- GIVEN 管理员已登录后台
- WHEN 管理员打开配置管理页面
- THEN 后台 SHALL 展示应用、环境、配置分组与版本列表
- AND 用户可进行查询与筛选

#### Scenario: 后台执行发布操作
- GIVEN 管理员已选择目标配置版本
- WHEN 管理员在后台点击发布
- THEN 后台 SHALL 发起发布请求并展示执行结果
- AND 页面 SHALL 显示最新生效版本

#### Scenario: 后台查看集群状态
- GIVEN 集群节点正在运行
- WHEN 管理员打开集群状态页面
- THEN 后台 SHALL 展示节点列表、leader 标识与基础健康状态
- AND 页面 SHALL 能反映节点变化