package ratelimiter

import (
	"testing"
	"time"
)

func TestSlidingWindowLog(t *testing.T) {
	limiter := NewSlidingWindowLog(1*time.Second, 3)

	for i := 0; i < 3; i++ {
		if !limiter.AllowRequest() {
			t.Errorf("Expected request %d to be allowed, but it was denied", i+1)
		}
	}

	if limiter.AllowRequest() {
		t.Errorf("Expected the 4th request to be denied, but it was allowed")
	}

	time.Sleep(1 * time.Second)

	if !limiter.AllowRequest() {
		t.Errorf("Expected the request after the window slide to be allowed, but it was denied")
	}
}
