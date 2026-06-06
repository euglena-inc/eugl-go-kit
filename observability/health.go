package observability

import (
	"context"
	"log/slog"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceInfo struct {
	ServiceName string
	Version     string
	Env         string
}

type DependencyCheck func(ctx context.Context) error

type namedDependencyCheck struct {
	name  string
	check DependencyCheck
}

type Health struct {
	info        ServiceInfo
	db          *gorm.DB
	redis       *goredis.Client
	extraChecks []namedDependencyCheck
}

func NewHealth(info ServiceInfo, database *gorm.DB, redisClient *goredis.Client, logger ...*slog.Logger) *Health {
	return &Health{
		info:  info,
		db:    database,
		redis: redisClient,
	}
}

func (h *Health) AddDependency(name string, check DependencyCheck) {
	if name == "" || check == nil {
		return
	}
	h.extraChecks = append(h.extraChecks, namedDependencyCheck{name: name, check: check})
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
		if err == nil {
			err = sqlDB.PingContext(ctx)
		}
		if err != nil {
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

	for _, extra := range h.extraChecks {
		if err := extra.check(ctx); err != nil {
			ready = false
			checks[extra.name] = map[string]interface{}{"configured": true, "status": "down"}
		} else {
			checks[extra.name] = map[string]interface{}{"configured": true, "status": "up"}
		}
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
