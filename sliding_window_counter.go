package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	windowSize    time.Duration
	maxRequests   int
	currentWindow int64
	requestCount  int
	previousCount int
	mu            sync.Mutex
}

func NewSlidingWindowCounter(windowSize time.Duration, maxRequests int) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		windowSize:    windowSize,
		maxRequests:   maxRequests,
		currentWindow: time.Now().Unix() / int64(windowSize.Seconds()),
		requestCount:  0,
		previousCount: 0,
	}
}

func (swc *SlidingWindowCounter) AllowRequest() bool {
	swc.mu.Lock()
	defer swc.mu.Unlock()

	now := time.Now().Unix()
	window := now / int64(swc.windowSize.Seconds())

	// If we've moved to a new window, update the counts
	if window != swc.currentWindow {
		swc.previousCount = swc.requestCount
		swc.requestCount = 0
		swc.currentWindow = window
	}

	// Calculate the weighted request count
	windowElapsed := float64(now%int64(swc.windowSize.Seconds())) / float64(swc.windowSize.Seconds())
	threshold := float64(swc.previousCount)*(1-windowElapsed) + float64(swc.requestCount)

	// Check if we're within the limit
	if threshold < float64(swc.maxRequests) {
		swc.requestCount++
		return true
	}

	return false
}
