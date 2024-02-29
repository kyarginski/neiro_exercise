package services

import (
	"log/slog"
	"time"

	"neiro/internal/app/repository"
)

type MyService struct {
	log     *slog.Logger
	storage *repository.KeyValueStore
	ttl     int
}

func (s *MyService) Set(key, value string) {
	s.storage.Set(key, value, s.ttl)
}

func (s *MyService) Get(key string) (string, bool) {
	return s.storage.Get(key)
}

func (s *MyService) Delete(key string) {
	s.storage.Delete(key)
}

func NewService(log *slog.Logger, ttl int) (IService, error) {
	storage := repository.NewKeyValueStore(log)

	return &MyService{
		log:     log,
		storage: storage,
		ttl:     ttl,
	}, nil
}

// ClearCache runs in a separate goroutine to remove expired keys.
func (s *MyService) ClearCache(d time.Duration) {
	ticker := time.NewTicker(d * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go s.storage.ClearExpired()
	}
}

func (s *MyService) LivenessCheck() bool {
	// Implement liveness check logic
	return true
}

func (s *MyService) ReadinessCheck() bool {
	// Implement readiness check logic
	return true
}

func (s *MyService) Logger() *slog.Logger {
	return s.log
}

func (s *MyService) Close() error {
	return nil
}
