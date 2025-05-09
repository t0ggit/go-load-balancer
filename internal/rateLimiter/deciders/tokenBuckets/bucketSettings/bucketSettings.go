package bucketSettings

import "time"

type BucketSettings struct {
	Capacity       int
	Refill         int
	RefillInterval time.Duration
}
