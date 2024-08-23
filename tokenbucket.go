package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	sync.Mutex
	tokens         int       // Current token count, start with a full bucket
	maxTokens      int       // Maximum number of tokens the bucket hold
	refillRate     int       // Rate at which tokens are added (tokens/second)
	lastRefillTime time.Time // Last time we checked the token count
}

func NewTokenBucket(maxTokens, refillRate int) *TokenBucket {
	return &TokenBucket{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) AddToken(tokens int) bool {
	tb.Lock()
	defer tb.Unlock()

	tb.refill()
	if tokens < tb.tokens {
		tb.tokens -= tokens
		return true
	}

	return false
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	duration := time.Since(tb.lastRefillTime)
	tokenAdd := tb.tokens * int(duration.Seconds())
	tb.tokens = tb.min(tb.maxTokens, tb.tokens+tokenAdd)
	tb.lastRefillTime = now
}

func (tb *TokenBucket) min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
