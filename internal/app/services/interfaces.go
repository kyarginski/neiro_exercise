package services

import (
	"log/slog"
	"time"

	"neiro/internal/app/health"
)

type IService interface {
	Logger() *slog.Logger
	Set(key, value string)
	Get(key string) (string, bool)
	Delete(key string)
	ClearCache(d time.Duration)
	Close() error

	health.LivenessChecker
	health.ReadinessChecker
}
