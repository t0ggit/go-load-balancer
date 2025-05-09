package set

type Request struct {
    Key            string `json:"key"`
    BucketCapacity int    `json:"bucket_capacity"`
    Refill         int    `json:"refill"`
    RefillInterval string `json:"refill_interval"`
}
