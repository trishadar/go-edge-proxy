package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	tokens   int
	capacity int
	last     time.Time
	rate     int
}

// NewLimiter creates a token bucket
func NewLimiter(rate, capacity int) *Limiter {
	return &Limiter{
		tokens:   capacity,
		capacity: capacity,
		last:     time.Now(),
		rate:     rate,
	}
}

// Allow checks if request is allowed
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.last).Seconds()
	l.last = now

	// refill tokens
	l.tokens += int(elapsed * float64(l.rate))
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}

	if l.tokens > 0 {
		l.tokens--
		return true
	}

	return false
}
