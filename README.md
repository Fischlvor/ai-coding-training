# config-center-raft

这是一个简化版配置中心的 Go 项目骨架，后续任务会逐步补齐基于 Raft 的一致性能力。

## 目录说明

- `cmd/config-center`：应用主入口
- `cmd/migrate`：数据库迁移命令入口
- `configs/config.example.yaml`：示例配置文件
- `configs/config.yaml`：本地运行配置文件，已加入忽略列表
- `migrations/`：PostgreSQL 迁移 SQL 文件
