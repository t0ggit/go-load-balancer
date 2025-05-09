package postgres

import (
	"context"
	"fmt"
	"go-load-balancer/internal/rateLimiter/deciders/tokenBuckets/bucketSettings"
)

func (s *BucketSettingsStorage) Set(ctx context.Context, keyHash string, settings bucketSettings.BucketSettings) error {
	refillIntervalStr := settings.RefillInterval.String()

	query := `
        INSERT INTO bucket_settings (key_hash, capacity, refill, refill_interval)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (key_hash)
        DO UPDATE SET
            capacity = EXCLUDED.capacity,
            refill = EXCLUDED.refill,
            refill_interval = EXCLUDED.refill_interval
    `

	_, err := s.db.ExecContext(ctx, query, keyHash, settings.Capacity, settings.Refill, refillIntervalStr)
	if err != nil {
		return fmt.Errorf("failed to set bucket settings: %s", err.Error())
	}

	return nil
}
