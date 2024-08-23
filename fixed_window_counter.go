package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowCounter struct {
	window_size    int64 // size of the window in seconds
	max_requests   int   // maximum number of requests per window
	current_window int64 // current window time
	request_count  int   // count of requests in the current window
	mu             sync.Mutex
}

func NewFixedWindowCounter(window_size int64, max_requests int) *FixedWindowCounter {
	return &FixedWindowCounter{
		window_size:    window_size,
		max_requests:   max_requests,
		current_window: time.Now().Unix() / window_size,
		request_count:  0,
	}
}

func (fxc *FixedWindowCounter) AllowRequest() bool {
	fxc.mu.Lock()
	defer fxc.mu.Unlock()

	current_time := time.Now()
	window := current_time.Unix() / fxc.window_size

	// If we've moved to a new window, reset the counter
	if window != fxc.current_window {
		fxc.current_window = window
		fxc.request_count = 0
	}

	// Check if we're still within the limit for this window
	if fxc.request_count < fxc.max_requests {
		fxc.request_count++
		return true
	}

	return false
}
