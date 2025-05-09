package rateLimiter

import (
    "go-load-balancer/internal/rateLimiter/api/handlers/get"
    "go-load-balancer/internal/rateLimiter/api/handlers/set"
    "log/slog"
    "net/http"
)

func (rl *RateLimiter) Server(log *slog.Logger, host string) *http.Server {
    mux := http.NewServeMux()

    mux.HandleFunc("POST /get", get.New(log, rl.decider))
    mux.HandleFunc("POST /set", set.New(log, rl.decider))

    server := &http.Server{
        Addr:    host,
        Handler: mux,
    }

    return server
}
