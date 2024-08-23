package tokenbucket

import (
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_AddToken(t *testing.T) {
	type fiels struct {
		token          int
		maxToken       int
		refillRate     int
		lastRefillTime time.Time
		Mutex          sync.Mutex
	}

	tests := []struct {
		name          string
		fiels         fiels
		want          bool
		takedToken    int
		expectedToken int
	}{
		{
			"token available now",
			fiels{
				token:          10,
				maxToken:       10,
				refillRate:     1,
				lastRefillTime: time.Now(),
			},
			true,
			9,
			1,
		},
		{
			"no token available now",
			fiels{
				token:          0,
				maxToken:       10,
				refillRate:     1,
				lastRefillTime: time.Now(),
			},
			false,
			10,
			0,
		},
		{
			"tokens available after adjustment",
			fiels{
				token:          1,
				maxToken:       10,
				refillRate:     1,
				lastRefillTime: time.Now().Add(-1 * time.Second),
			},
			true,
			1,
			1,
		},
		{
			"tokens do not refresh above capacity",
			fiels{
				token:          10,
				maxToken:       10,
				refillRate:     1,
				lastRefillTime: time.Now().Add(-2 * time.Minute),
			},
			true,
			1,
			9,
		},
		{
			"refreshs 4 tokens",
			fiels{
				token:          4,
				maxToken:       10,
				refillRate:     1,
				lastRefillTime: time.Now().Add(-5 * time.Second),
			},
			true,
			4,
			6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &TokenBucket{
				tokens:         tt.fiels.token,
				maxTokens:      tt.fiels.maxToken,
				refillRate:     tt.fiels.refillRate,
				lastRefillTime: tt.fiels.lastRefillTime,
				Mutex:          tt.fiels.Mutex,
			}

			if got := b.AddToken(tt.takedToken); got != tt.want {
				t.Errorf("TokenBucket.Add(%v) = %v, want %v", tt.takedToken, got, tt.want)
			}

			if count := b.tokens; count != tt.expectedToken {
				t.Errorf("Token count incorrect.  Got %v, want %v", count, tt.expectedToken)
			}

			if b.tokens > b.maxTokens {
				t.Errorf("Max token is %v but current count is %v", b.maxTokens, b.tokens)
			}
		})
	}
}
