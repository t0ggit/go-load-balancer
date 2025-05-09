package tokenBuckets

import (
	"context"
	"errors"
	"fmt"
	"go-load-balancer/internal/rateLimiter"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"log/slog"
	"reflect"
	"sync"
)

var (
	ErrNoSuchKeyHashInDatabase = fmt.Errorf("no such key hash in database")
)

type ExternalBucketSettingsStorage interface {
	ExternalBucketSettingsProvider
	ExternalBucketSettingsSetter
}

type TokenBuckets struct {
	buckets                       map[string]*Bucket
	mux                           *sync.RWMutex
	defaultBucketSettings         bucketSettings.BucketSettings
	externalBucketSettingsStorage ExternalBucketSettingsStorage
}

func New(defaultBucketSettings bucketSettings.BucketSettings,
	externalBucketSettingsStorage ExternalBucketSettingsStorage) *TokenBuckets {
	// Если передано значение интерфейса, но оно указывает на nil-значение, сбрасываем интерфейс в nil
	if externalBucketSettingsStorage == nil || reflect.ValueOf(externalBucketSettingsStorage).IsNil() {
		externalBucketSettingsStorage = nil
	}

	return &TokenBuckets{
		buckets:                       make(map[string]*Bucket, 32),
		mux:                           &sync.RWMutex{},
		defaultBucketSettings:         defaultBucketSettings,
		externalBucketSettingsStorage: externalBucketSettingsStorage,
	}
}

func (tb *TokenBuckets) IsAllowed(ctx context.Context, log *slog.Logger, keyHash string) bool {
	tb.mux.RLock()
	bucket, ok := tb.buckets[keyHash]
	tb.mux.RUnlock()

	if ok {
		return bucket.TakeToken()
	}

	var newBucket *Bucket
	// Если в нашей "локальной" таблице нет бакета, то создаем его.
	// Но чтобы его создать, нужны настройки бакета (либо из внешнего хранилища, либо дефолтные)

	settings, err := tb.GetBucketSettings(ctx, log, keyHash)
	if err != nil {
		switch {
		case errors.Is(err, rateLimiter.ErrBucketSettingsNotFound):
		default:
			log.Error("failed to get bucket settings", slog.String("error", err.Error()))
		}
		// Если нигде не нашли настроек, используем дефолтные
		newBucket = NewBucket(tb.defaultBucketSettings)
	} else {
		newBucket = NewBucket(settings)
	}

	tb.mux.Lock()
	tb.buckets[keyHash] = newBucket
	tb.mux.Unlock()
	return newBucket.TakeToken()
}
