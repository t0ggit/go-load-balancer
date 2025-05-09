package get

type Response struct {
    BucketCapacity int    `json:"bucket_capacity"`
    Refill         int    `json:"refill"`
    RefillInterval string `json:"refill_interval"`
}
