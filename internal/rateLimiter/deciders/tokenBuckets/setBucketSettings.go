package tokenBuckets

import (
	"context"
	"fmt"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"log/slog"
)

type ExternalBucketSettingsSetter interface {
	Set(ctx context.Context, keyHash string, settings bucketSettings.BucketSettings) error
}

func (tb *TokenBuckets) SetBucketSettings(ctx context.Context, log *slog.Logger, keyHash string, settings bucketSettings.BucketSettings) error {
	newBucket := NewBucket(settings)

	tb.mux.RLock()
	bucket, ok := tb.buckets[keyHash]
	tb.mux.RUnlock()
	if !ok {
		tb.mux.Lock()
		tb.buckets[keyHash] = newBucket
		tb.mux.Unlock()
	} else {
		bucket.SetSettings(settings)
	}

	if tb.externalBucketSettingsStorage != nil {
		log.Info("setting bucket settings to database", slog.String("key_hash", keyHash))
		err := tb.externalBucketSettingsStorage.Set(ctx, keyHash, settings)
		if err != nil {
			return fmt.Errorf("cannot set bucket settings to database: %s", err)
		}
	}

	return nil
}
