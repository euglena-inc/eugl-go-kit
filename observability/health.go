package observability

import (
	"context"
	"log/slog"
	"sync"

	"github.com/euglena-inc/eugl-go-kit/requestid"
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
	info             ServiceInfo
	db               *gorm.DB
	redis            *goredis.Client
	log              *slog.Logger
	mu               sync.Mutex
	dependencyStatus map[string]string
	extraChecks      []namedDependencyCheck
}

func NewHealth(info ServiceInfo, database *gorm.DB, redisClient *goredis.Client, logger ...*slog.Logger) *Health {
	log := slog.Default()
	if len(logger) > 0 && logger[0] != nil {
		log = logger[0]
	}
	return &Health{
		info:             info,
		db:               database,
		redis:            redisClient,
		log:              log,
		dependencyStatus: map[string]string{},
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
			h.logDependencyStatus(ctx, "postgres", "down", err)
		} else {
			checks["postgres"] = map[string]interface{}{"configured": true, "status": "up"}
			h.logDependencyStatus(ctx, "postgres", "up", nil)
		}
	}

	if h.redis == nil {
		checks["redis"] = map[string]interface{}{"configured": false, "status": "skipped"}
	} else if err := h.redis.Ping(ctx).Err(); err != nil {
		ready = false
		checks["redis"] = map[string]interface{}{"configured": true, "status": "down"}
		h.logDependencyStatus(ctx, "redis", "down", err)
	} else {
		checks["redis"] = map[string]interface{}{"configured": true, "status": "up"}
		h.logDependencyStatus(ctx, "redis", "up", nil)
	}

	for _, extra := range h.extraChecks {
		if err := extra.check(ctx); err != nil {
			ready = false
			checks[extra.name] = map[string]interface{}{"configured": true, "status": "down"}
			h.logDependencyStatus(ctx, extra.name, "down", err)
		} else {
			checks[extra.name] = map[string]interface{}{"configured": true, "status": "up"}
			h.logDependencyStatus(ctx, extra.name, "up", nil)
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

func (h *Health) logDependencyStatus(ctx context.Context, dependency string, status string, err error) {
	h.mu.Lock()
	previous := h.dependencyStatus[dependency]
	h.dependencyStatus[dependency] = status
	h.mu.Unlock()

	attrs := []slog.Attr{
		slog.String("request_id", requestid.FromContext(ctx)),
		slog.String("service_name", h.info.ServiceName),
		slog.String("dependency", dependency),
		slog.String("status", status),
	}
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	if status == "down" {
		h.log.LogAttrs(ctx, slog.LevelWarn, "dependency readiness down", attrs...)
		return
	}
	if status == "up" && previous == "down" {
		h.log.LogAttrs(ctx, slog.LevelInfo, "dependency readiness recovered", attrs...)
	}
}
