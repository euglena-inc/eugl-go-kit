package observability

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceInfo struct {
	ServiceName string
	Version     string
	Env         string
}

type Health struct {
	info  ServiceInfo
	db    *gorm.DB
	redis *goredis.Client
}

func NewHealth(info ServiceInfo, database *gorm.DB, redisClient *goredis.Client) *Health {
	return &Health{
		info:  info,
		db:    database,
		redis: redisClient,
	}
}

func (h *Health) Liveness() map[string]interface{} {
	return map[string]interface{}{
		"service_name": h.info.ServiceName,
		"version":      h.info.Version,
		"env":          h.info.Env,
		"status":       "alive",
	}
}

func (h *Health) Readiness(ctx context.Context) (map[string]interface{}, bool) {
	checks := map[string]interface{}{}
	ready := true

	if h.db == nil {
		checks["postgres"] = map[string]interface{}{"configured": false, "status": "skipped"}
	} else {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.PingContext(ctx) != nil {
			ready = false
			checks["postgres"] = map[string]interface{}{"configured": true, "status": "down"}
		} else {
			checks["postgres"] = map[string]interface{}{"configured": true, "status": "up"}
		}
	}

	if h.redis == nil {
		checks["redis"] = map[string]interface{}{"configured": false, "status": "skipped"}
	} else if err := h.redis.Ping(ctx).Err(); err != nil {
		ready = false
		checks["redis"] = map[string]interface{}{"configured": true, "status": "down"}
	} else {
		checks["redis"] = map[string]interface{}{"configured": true, "status": "up"}
	}

	status := "ready"
	if !ready {
		status = "not_ready"
	}

	return map[string]interface{}{
		"service_name": h.info.ServiceName,
		"version":      h.info.Version,
		"env":          h.info.Env,
		"status":       status,
		"checks":       checks,
	}, ready
}
