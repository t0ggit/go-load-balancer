package tokenBuckets

import (
	"context"
	"errors"
	"go-load-balancer/internal/rateLimiter"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"log/slog"
)

type ExternalBucketSettingsProvider interface {
	Get(ctx context.Context, keyHash string) (settings bucketSettings.BucketSettings, err error)
}

// GetBucketSettings сначала ищет бакет в локальной таблице,
// если не находит, то ищет во внешнем хранилище настройки бакета для указанного keyHash,
// если такого нет, то возвращает ошибку ErrBucketSettingsNotFound
func (tb *TokenBuckets) GetBucketSettings(ctx context.Context, log *slog.Logger, keyHash string) (settings bucketSettings.BucketSettings, err error) {
	// Сначала смотрим в "локальной" таблице
	tb.mux.RLock()
	bucket, ok := tb.buckets[keyHash]
	tb.mux.RUnlock()
	if ok {
		return bucket.GetSettings(), nil
	}

	// Если не нашли, то идем во внешнее хранилище
	if tb.externalBucketSettingsStorage != nil {
		settings, err = tb.externalBucketSettingsStorage.Get(ctx, keyHash)
		if err != nil {
			switch {
			case errors.Is(ErrNoSuchKeyHashInDatabase, err):
				log.Info("no such key hash in database", slog.String("key_hash", keyHash))
				return bucketSettings.BucketSettings{}, rateLimiter.ErrBucketSettingsNotFound
			default:
				return bucketSettings.BucketSettings{}, err
			}
		}
		return settings, nil
	}

	// Если внешнего хранилища нет, то возвращаем ошибку
	return bucketSettings.BucketSettings{}, rateLimiter.ErrBucketSettingsNotFound
}
