package rateLimiter

import (
    "context"
    "fmt"
    "go-load-balancer/internal/rateLimiter/api/handlers/get"
    "go-load-balancer/internal/rateLimiter/api/handlers/set"
    "log/slog"
)

var (
    ErrBucketSettingsNotFound = fmt.Errorf("bucket settings not found")
)

type BucketSettingsStorage interface {
    get.BucketSettingsProvider
    set.BucketSettingsSetter
}

// Decider определяет разрешено ли обрабатывать запрос
type Decider interface {
    IsAllowed(ctx context.Context, log *slog.Logger, keyHash string) bool
    BucketSettingsStorage
}

type RateLimiter struct {
    decider Decider
    log     *slog.Logger
}

func New(log *slog.Logger, decider Decider) *RateLimiter {
    return &RateLimiter{
        decider: decider,
        log:     log,
    }
}
