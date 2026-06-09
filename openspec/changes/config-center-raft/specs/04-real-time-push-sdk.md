### Requirement: 实时配置推送
系统 SHALL 在配置变更后向已订阅客户端推送标准化通知事件，并支持客户端断线重连后的补偿拉取。系统 MUST 保证同一事件的重复投递不影响客户端最终状态。

#### Scenario: 配置变更推送
- GIVEN 客户端已订阅某配置分组
- WHEN 该分组配置发生发布或回滚
- THEN 系统 SHALL 向客户端发送变更通知
- AND 通知 SHALL 包含应用、环境、分组、版本号与事件类型

#### Scenario: 重复事件投递
- GIVEN 客户端已经处理过某一配置变更事件
- WHEN 系统再次投递相同事件
- THEN 客户端 SHALL 识别该事件为重复事件
- AND 客户端最终状态 SHALL 保持不变

#### Scenario: 客户端断开连接
- GIVEN 客户端订阅通道已断开
- WHEN 配置发生变化
- THEN 系统 SHALL 记录推送失败或断开状态
- AND 不应影响其他正常客户端

#### Scenario: 客户端重连补偿拉取
- GIVEN 客户端曾经断开并重新连接
- WHEN 客户端携带上次已知版本号重新订阅
- THEN 系统 SHALL 返回缺失期间的最新版本信息
- AND 客户端 SHALL 重新拉取当前生效配置

### Requirement: SDK 接入能力
系统 SHALL 提供 Go SDK 接入方式，以便客户端拉取、订阅和解析配置。系统 SHOULD 支持基础重试与本地缓存刷新。系统 SHOULD 支持按环境区分加载不同配置源。

#### Scenario: SDK 初始化
- GIVEN SDK 已获取服务端地址与客户端身份信息
- WHEN 客户端创建 SDK 实例
- THEN SDK SHALL 完成初始化
- AND SDK SHALL 具备拉取与订阅能力

#### Scenario: SDK 拉取配置
- GIVEN SDK 已完成初始化
- WHEN 客户端请求某应用配置
- THEN SDK SHALL 调用服务端接口并返回当前生效配置
- AND 返回结果 SHALL 包含版本信息与灰度状态

#### Scenario: SDK 订阅配置变化
- GIVEN SDK 已完成订阅注册
- WHEN 服务端推送配置变化
- THEN SDK SHALL 接收标准化事件并更新本地缓存
- AND 客户端 SHALL 能读取到最新版本配置

#### Scenario: SDK 解析灰度结果
- GIVEN 服务端返回灰度发布结果
- WHEN SDK 根据客户端身份信息计算命中情况
- THEN SDK SHALL 返回是否命中灰度的结果
- AND 客户端 SHALL 据此选择生效版本

#### Scenario: SDK 拉取失败
- GIVEN 服务端暂时不可达
- WHEN 客户端请求配置
- THEN SDK SHALL 返回可识别的错误或重试结果
- AND 客户端 SHOULD 继续使用本地缓存

### Requirement: 可视化管理后台
系统 SHALL 提供可视化管理后台用于配置管理与集群状态查看，并展示发布结果。系统 SHOULD 保留核心操作路径，避免引入复杂审批流。系统 SHOULD 支持按环境切换不同配置视图。

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

#### Scenario: 后台查询发布失败原因
- GIVEN 一次发布操作失败
- WHEN 管理员查看发布详情
- THEN 后台 SHALL 展示失败原因与当前状态
- AND 页面 SHOULD 保持最近一次有效状态可见