# AGENTS.md

本仓库是 EUGL Go 公共基础能力库。AI / Agent 在本仓库工作时，必须遵守以下规则。

## 边界

允许放入：

- 配置加载。
- 结构化日志。
- request_id。
- HTTP 响应封装。
- Gin 中间件。
- PostgreSQL / GORM 连接封装。
- Redis client 和 TTL 基础方法。
- Resty HTTP client。
- 健康检查和可观测性基础结构。
- 跨领域事件 envelope 和 publisher interface。

禁止放入：

- 订单、餐品、门店、会员、营销、支付等业务实体。
- 领域服务、业务 workflow、业务 repository。
- 具体领域的状态机。
- 具体业务错误码号段。
- 为某个单一服务定制的逻辑。
- `common`、`utils`、`helper` 这类含义不清的包。

## 技术约束

- Go 版本使用 1.25。
- 配置加载使用 `github.com/caarlos0/env/v11`。
- 日志使用 Go 标准库 `log/slog`。
- HTTP client 使用 `github.com/go-resty/resty/v2`。
- Redis 使用 `github.com/redis/go-redis/v9`。
- GORM 只用于数据库连接封装。
- `db` 包必须保持 GORM 默认表名约定，不配置 `NamingStrategy.SingularTable = true`。
- `db` 包不得执行 `AutoMigrate`，生产表结构只能由领域项目 migration 管理。
- 数据库表结构、GORM model 和 migration 规则以 `data-schema-standards` 为准。
- 不提交真实 `.env`、token、secret、证书、私钥。

## 完成前检查

```text
go fmt ./...
go test ./...
go vet ./...
```

如果新增包，必须同步更新 README 的包列表。
