package get

import (
	"context"
	"encoding/json"
	"go-load-balancer/internal/rateLimiter/api"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"log/slog"
	"net/http"
)

type BucketSettingsProvider interface {
	GetBucketSettings(ctx context.Context, log *slog.Logger, keyHash string) (settings bucketSettings.BucketSettings, err error)
}

func New(log *slog.Logger, bucketSettingsProvider BucketSettingsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rateLimiter.api.handlers.get.New"
		log := log.With(slog.String("op", op))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		keyHash := api.HashKey(req.Key)

		settings, err := bucketSettingsProvider.GetBucketSettings(r.Context(), log, keyHash)
		if err != nil {

			log.Error("failed to get bucket settings", slog.String("error", err.Error()))
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		resp := Response{
			BucketCapacity: settings.Capacity,
			Refill:         settings.Refill,
			RefillInterval: settings.RefillInterval.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", slog.String("error", err.Error()))
		}
	}
}
