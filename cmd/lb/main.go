package main

import (
	"context"
	"errors"
	"fmt"
	"go-load-balancer/internal/config"
	"go-load-balancer/internal/loadBalancer"
	"go-load-balancer/internal/loadBalancer/iterators/roundRobin"
	"go-load-balancer/internal/rateLimiter"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/storage/postgres"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	// Config init
	cfg := config.MustLoad()

	// Logger init
	log := setupLogger(cfg.Env)

	// LoadBalancer init
	lb := loadBalancer.New(log, roundRobin.New())
	for _, backend := range cfg.Backends {
		lb.TryToRegisterNewBackend(backend.Url)
	}

	// BucketSettingsStorage init
	log.Info("initializing bucket settings storage")
	bss, err := postgres.New(cfg.RateLimiter.BucketSettingsDatabase)
	if err != nil {
		log.Error("failed to create bucket settings storage", slog.String("error", err.Error()))
	}
	// RateLimiter init
	log.Info("initializing rate limiter")
	rl := rateLimiter.New(log, tokenBuckets.New(
		bucketSettings.BucketSettings{
			Capacity:       cfg.RateLimiter.DefaultBucketCapacity,
			Refill:         cfg.RateLimiter.DefaultRefill,
			RefillInterval: cfg.RateLimiter.DefaultRefillInterval,
		},
		bss))
	apiServer := rl.Server(log, cfg.RateLimiter.API.Host)
	go func() {
		log.Info("starting api server", slog.String("addr", apiServer.Addr))
		_ = apiServer.ListenAndServe()
	}()

	// HTTP Server init
	server := &http.Server{
		Addr:    fmt.Sprintf(cfg.Host),
		Handler: rl.Middleware()(lb),
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		sig := <-sigChan
		log.Info("shutting down api server", slog.String("signal", sig.String()))
		_ = apiServer.Shutdown(context.TODO())
		log.Info("shutting down server", slog.String("signal", sig.String()))
		_ = server.Shutdown(context.TODO())
	}()

	// Start server
	log.Info("starting server", slog.String("addr", server.Addr))
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}

const (
	envLocal = "local"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
