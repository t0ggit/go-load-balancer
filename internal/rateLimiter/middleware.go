package rateLimiter

import (
    "go-load-balancer/internal/rateLimiter/api"
    "log/slog"
    "net/http"
)

func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := r.RemoteAddr
            apiKey := r.Header.Get("Authorization")
            if apiKey != "" {
                key = apiKey
            }

            keyHash := api.HashKey(key)

            if rl.decider.IsAllowed(r.Context(), rl.log, keyHash) {
                next.ServeHTTP(w, r)
            } else {
                rl.log.Info("request rejected", slog.String("client", keyHash))
                http.Error(w, "too many requests, try again later", http.StatusTooManyRequests)
            }
        })
    }
}
