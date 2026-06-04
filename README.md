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
| `middleware` | Gin request id、access log 中间件 |
| `db` | PostgreSQL / GORM 连接初始化和关闭 |
| `redis` | go-redis client 初始化和 TTL 写入封装 |
| `httpclient` | Resty HTTP client，自动透传 `X-Request-Id` |
| `observability` | 健康检查数据结构和依赖探测 |
| `event` | 跨领域事件 envelope 和 publisher 接口 |

## 使用原则

- 只能放技术公共能力。
- 不能放任何餐饮业务逻辑。
- 不能放订单、餐品、门店、会员、营销、支付等领域对象。
- Redis 临时 key 必须通过带 TTL 的方法写入。
- 日志、HTTP client、事件发布都必须透传或打印 `request_id`。
- 对外 API 响应必须使用统一 envelope。

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

