package ratelimiter

import (
	"testing"
	"time"
)

func TestFixedWindowCounter(t *testing.T) {
	limiter := NewFixedWindowCounter(1, 5)

	for i := 0; i < 7; i++ {
		if i < 5 {
			if !limiter.AllowRequest() {
				t.Errorf("Expected request %d to be allowed, but it was denied", i+1)
			}
		} else {
			if limiter.AllowRequest() {
				t.Errorf("Expected request %d to be denied, but it was allowed", i+1)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)

	if !limiter.AllowRequest() {
		t.Errorf("Expected request after window reset to be allowed, but it was denied")
	}
}
