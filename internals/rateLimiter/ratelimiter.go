package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	MaxTokens  int64
	Tokens     int64
	Rate       int64
	LastRefill time.Time
	Mutex      sync.Mutex
}

func NewTokenBucket(rate, maxTokens int64) *TokenBucket {
	newBucket := &TokenBucket{
		MaxTokens:  maxTokens,
		Tokens:     maxTokens,
		Rate:       rate,
		LastRefill: time.Now(),
	}

	return newBucket
}

func (tb *TokenBucket) Allow() bool {
	tb.Mutex.Lock()
	defer tb.Mutex.Unlock()

	timeElapsed := time.Since(tb.LastRefill).Seconds()
	tokensToAdd := timeElapsed * float64(tb.Rate)
	tb.Tokens = min(tb.Tokens+int64(tokensToAdd), tb.MaxTokens)
	tb.LastRefill = time.Now()

	if tb.Tokens > 0 {
		tb.Tokens--
		return true
	}

	return false
}
