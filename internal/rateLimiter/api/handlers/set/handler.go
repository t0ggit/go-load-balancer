package set

import (
	"context"
	"encoding/json"
	"go-load-balancer/internal/rateLimiter/api"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"log/slog"
	"net/http"
	"time"
)

type BucketSettingsSetter interface {
	SetBucketSettings(ctx context.Context, log *slog.Logger, keyHash string, settings bucketSettings.BucketSettings) error
}

func New(log *slog.Logger, bucketSettingsSetter BucketSettingsSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rateLimiter.api.handlers.set.New"
		log := log.With(slog.String("op", op))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info("failed to decode request", slog.String("error", err.Error()))
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		keyHash := api.HashKey(req.Key)

		parsedInterval, err := time.ParseDuration(req.RefillInterval)
		if err != nil {
			log.Info("invalid refill interval", slog.String("error", err.Error()))
			http.Error(w, "invalid refill interval", http.StatusBadRequest)
			return
		}

		err = bucketSettingsSetter.SetBucketSettings(r.Context(), log,
			keyHash, bucketSettings.BucketSettings{
				Capacity: req.BucketCapacity, Refill: req.Refill, RefillInterval: parsedInterval,
			})
		if err != nil {
			log.Error("failed to set bucket settings", slog.String("error", err.Error()))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		log.Info("bucket settings set", slog.String("key_hash", keyHash),
			slog.Int("capacity", req.BucketCapacity), slog.Int("refill", req.Refill),
			slog.String("refill_interval", req.RefillInterval))
	}
}
