package commands

import (
	"sync"
	"time"
)

const (
	maxToolCallsPerMinute = 60
	burstSize             = 10
)

// TokenBucket implements a simple token bucket rate limiter.
type TokenBucket struct {
	mu         sync.Mutex
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
}

// NewTokenBucket creates a new token bucket rate limiter.
func NewTokenBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if the operation is allowed and consumes a token if so.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed / tb.refillRate)

	if tokensToAdd > 0 {
		tb.tokens = minInt(tb.maxTokens, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var globalRateLimiter = NewTokenBucket(burstSize, time.Minute/time.Duration(maxToolCallsPerMinute))
