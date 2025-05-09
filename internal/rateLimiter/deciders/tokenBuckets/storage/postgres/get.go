package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
	"time"
)

func (s *BucketSettingsStorage) Get(ctx context.Context, keyHash string) (bucketSettings.BucketSettings, error) {
	var settings bucketSettings.BucketSettings
	var refillIntervalStr string

	query := `SELECT capacity, refill, refill_interval FROM bucket_settings WHERE key_hash = $1`
	row := s.db.QueryRowContext(ctx, query, keyHash)
	err := row.Scan(&settings.Capacity, &settings.Refill, &refillIntervalStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return settings, tokenBuckets.ErrNoSuchKeyHashInDatabase
		}
		return bucketSettings.BucketSettings{}, fmt.Errorf("failed to fetch bucket settings: %s", err.Error())
	}

	refillInterval, err := time.ParseDuration(refillIntervalStr)
	if err != nil {
		return bucketSettings.BucketSettings{}, fmt.Errorf("invalid refill interval format fetched from database: %s", err.Error())
	}

	settings.RefillInterval = refillInterval
	return bucketSettings.BucketSettings{}, nil
}
