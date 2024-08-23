package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowLog struct {
	window_size time.Duration
	maxRequests int
	requestLog  []time.Time
	mu          sync.Mutex
}

func NewSlidingWindowLog(window_size time.Duration, max_requests int) *SlidingWindowLog {
	return &SlidingWindowLog{
		window_size: window_size,
		maxRequests: max_requests,
		requestLog:  make([]time.Time, 0),
	}
}

func (swl *SlidingWindowLog) AllowRequest() bool {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	now := time.Now()
	// Remove timestamps that are outside the current window
	for len(swl.requestLog) > 0 && now.Sub(swl.requestLog[0]) >= swl.window_size {
		swl.requestLog = swl.requestLog[1:]
	}

	// Check if we're still within the limit
	if len(swl.requestLog) < swl.maxRequests {
		swl.requestLog = append(swl.requestLog, now)
		return true
	}

	return false
}
