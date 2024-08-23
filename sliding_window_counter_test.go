package ratelimiter

import (
	"testing"
	"time"
)

func newSlidingWindowCounter(windowSize time.Duration, maxRequests int) *SlidingWindowCounter {
	return NewSlidingWindowCounter(windowSize, maxRequests)
}

func TestSlidingWindowCounter_AllowRequest(t *testing.T) {
	tests := []struct {
		windowSize    time.Duration
		maxRequests   int
		requests      int
		expectAllowed bool
	}{
		{time.Second * 10, 5, 5, true},   // within limit
		{time.Second * 10, 5, 6, false},  // exceeding limit
		{time.Second * 10, 5, 10, false}, // far exceeding limit
		{time.Second * 5, 2, 2, true},    // within smaller window
		{time.Second * 5, 2, 3, false},   // exceeding limit in smaller window
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			swc := newSlidingWindowCounter(tt.windowSize, tt.maxRequests)

			for i := 0; i < tt.requests; i++ {
				allowed := swc.AllowRequest()
				if i < tt.maxRequests && !allowed {
					t.Errorf("Request %d was not allowed, but it should be", i)
				}
				if i >= tt.maxRequests && allowed {
					t.Errorf("Request %d was allowed, but it should not be", i)
				}
			}
		})
	}
}

func TestSlidingWindowCounter_WindowExpiration(t *testing.T) {
	windowSize := time.Second * 2
	maxRequests := 2

	swc := newSlidingWindowCounter(windowSize, maxRequests)

	if !swc.AllowRequest() {
		t.Errorf("First request should be allowed")
	}

	if !swc.AllowRequest() {
		t.Errorf("Second request should be allowed")
	}

	time.Sleep(windowSize + time.Second)

	if swc.AllowRequest() {
		t.Errorf("Request after window expiration should not be allowed")
	}
}
