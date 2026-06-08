# EUGL Go Kit

`eugl-go-kit` 是 EUGL Go 后端服务的公共基础能力库。

它只放“脱离餐饮业务仍然成立”的基础设施能力，不放订单、餐品、门店、会员、支付等业务逻辑。

## 包列表

| 包 | 用途 |
| --- | --- |
| `config` | 基于 `github.com/caarlos0/env/v11` 的泛型 env 配置加载 |
| `logger` | 基于 Go 标准库 `log/slog` 的 JSON logger |
| `requestid` | `X-Request-Id` 生成、读取、写入 context |
| `errno` | 通用业务错误类型 |
| `response` | Gin 统一 JSON 响应 envelope |
| `middleware` | Gin request id、access log 中间件；access log 自动带常见业务 ID |
| `db` | PostgreSQL / GORM 连接初始化和关闭；不管理 schema 和 migration |
| `dbvalue` | 数据库 nullable 值的无业务转换 |
| `idgen` | 无业务含义的时间型 ID 生成器 |
| `redis` | go-redis client 初始化、关闭、TTL 写入和 Redis key 基础拼接 |
| `httpclient` | Resty HTTP client，自动透传 `X-Request-Id` |
| `httpserver` | HTTP server 启动、信号退出和优雅关闭封装 |
| `observability` | 健康检查数据结构和依赖探测；健康检查不打印日志 |
| `event` | 跨领域事件 envelope 和 publisher 接口 |

## 使用原则

- 只能放技术公共能力。
- 不能放任何餐饮业务逻辑。
- 不能放订单、餐品、门店、会员、营销、支付等领域对象。
- Redis 临时 key 必须通过带 TTL 的方法写入。
- Redis client 关闭、TTL 写入和无业务含义的 key 拼接统一使用 `redis` 包；业务服务只保留 key 前缀、TTL 策略和失效语义。
- 日志、HTTP client、事件发布都必须透传或打印 `request_id`。
- HTTP access log 必须打印 `request_id`，并尽量自动带上 `brand_id`、`store_id`、`order_no` 等常见业务定位字段；`/healthz`、`/readyz` 不打印 access log。
- HTTP 业务请求出现 4xx、5xx 或 Gin error 时必须额外打印 `http_error` 日志。
- `/readyz` 只返回依赖状态，不打印日志；额外依赖通过 `observability.Health.AddDependency` 注册。
- 服务启动入口优先使用 `httpserver.Run` 统一处理信号退出、HTTP graceful shutdown 和资源关闭。
- 对外 API 响应必须使用统一 envelope。
- `db` 包只负责连接初始化，不配置 `SingularTable`，不执行 `AutoMigrate`。
- 数据库表结构、GORM model 和 migration 以 [euglena-inc/data-schema-standards](https://github.com/euglena-inc/data-schema-standards) 为准。

## 本地检查

```bash
go fmt ./...
go test ./...
go vet ./...
```

或执行：

```bash
./scripts/check.sh
```

## 标准来源

公司级标准仓库：[euglena-inc/data-schema-standards](https://github.com/euglena-inc/data-schema-standards)
