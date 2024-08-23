package ratelimiter

import (
	"testing"
	"time"
)

func TestLeakyBucket(t *testing.T) {
	tests := []struct {
		name           string
		capacity       int
		leakRate       int
		initRequests   int
		sleepDuration  time.Duration
		expectedResult bool
	}{
		{
			name:           "Allow request when bucket is not full",
			capacity:       5,
			leakRate:       1,
			initRequests:   4,
			sleepDuration:  0,
			expectedResult: true,
		},
		{
			name:           "Deny request when bucket is full",
			capacity:       5,
			leakRate:       1,
			initRequests:   5,
			sleepDuration:  0,
			expectedResult: false,
		},
		{
			name:           "Allow request after bucket leaks",
			capacity:       5,
			leakRate:       1,
			initRequests:   5,
			sleepDuration:  2 * time.Second,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := NewLeakyBucket(tt.capacity, tt.leakRate)

			for i := 0; i < tt.initRequests; i++ {
				lb.AllowRequest()
			}

			if tt.sleepDuration > 0 {
				time.Sleep(tt.sleepDuration)
			}

			if lb.AllowRequest() != tt.expectedResult {
				t.Errorf("AllowRequest() = %v, want %v", lb.AllowRequest(), tt.expectedResult)
			}
		})
	}
}
