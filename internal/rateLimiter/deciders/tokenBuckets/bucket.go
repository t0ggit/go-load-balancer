package tokenBuckets

import (
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"sync"
	"time"
)

type Bucket struct {
	tokens   int
	settings *bucketSettings.BucketSettings
	mux      *sync.RWMutex
}

func NewBucket(settings bucketSettings.BucketSettings) *Bucket {
	bucket := &Bucket{
		tokens:   settings.Capacity / 2,
		settings: &settings,
		mux:      &sync.RWMutex{},
	}

	// Start refilling goroutine
	go func() {
		interval := settings.RefillInterval
		for {
			time.Sleep(interval)
			bucket.mux.Lock()
			bucket.tokens += bucket.settings.Refill
			if bucket.tokens > bucket.settings.Capacity {
				bucket.tokens = bucket.settings.Capacity
			}
			interval = bucket.settings.RefillInterval
			bucket.mux.Unlock()
		}

	}()
	return bucket
}

func (b *Bucket) GetSettings() bucketSettings.BucketSettings {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return *b.settings
}

func (b *Bucket) SetSettings(settings bucketSettings.BucketSettings) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.settings = &settings
}

func (b *Bucket) TakeToken() bool {
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return b.tokens >= 0
}
