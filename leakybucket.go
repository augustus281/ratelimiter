package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity int         // maximum number of requests in the bucket
	leakRate int         // rate at which requests leak (requests/second)
	bucket   []time.Time // hold request timestamps
	lastLeak time.Time   // last time we leaked from the bucket
	sync.Mutex
}

func NewLeakyBucket(capacity, leakRate int) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		bucket:   []time.Time{},
		lastLeak: time.Now(),
	}
}

func (lb *LeakyBucket) AllowRequest() bool {
	now := time.Now()
	leakedTime := now.Sub(lb.lastLeak).Seconds()

	leaked := int(leakedTime) * lb.leakRate
	if leaked > 0 {
		if leaked > len(lb.bucket) {
			lb.bucket = []time.Time{}
		} else {
			lb.bucket = lb.bucket[leaked:]
		}
		lb.lastLeak = now
	}

	if len(lb.bucket) < lb.capacity {
		lb.bucket = append(lb.bucket, now)
		return true
	}

	return false
}
