package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// RateLimiter interface defines the contract for rate limiting implementations
type RateLimiter interface {
	Allow() bool
	Wait(ctx context.Context) error
	Reset()
	GetStats() RateLimiterStats
}

// RateLimiterStats provides statistics about the rate limiter
type RateLimiterStats struct {
	RequestsAllowed  int64
	RequestsDenied   int64
	CurrentTokens    int
	RefillRate       float64
	LastRefill       time.Time
}

// TokenBucketLimiter implements rate limiting using the token bucket algorithm
type TokenBucketLimiter struct {
	mu           sync.Mutex
	capacity     int           // Maximum tokens in bucket
	tokens       int           // Current tokens available
	refillRate   time.Duration // Time between token refills
	lastRefill   time.Time     // Last time tokens were added
	allowed      int64         // Number of requests allowed
	denied       int64         // Number of requests denied
	logger       *zap.Logger
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(capacity int, refillRate time.Duration, logger *zap.Logger) *TokenBucketLimiter {
	if capacity <= 0 {
		panic("capacity must be positive")
	}
	if refillRate <= 0 {
		panic("refillRate must be positive")
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	return &TokenBucketLimiter{
		capacity:   capacity,
		tokens:     capacity, // Start with full bucket
		refillRate: refillRate,
		lastRefill: time.Now(),
		logger:     logger,
	}
}

// Allow checks if a request is allowed (non-blocking)
func (tbl *TokenBucketLimiter) Allow() bool {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	tbl.refillTokens()

	if tbl.tokens > 0 {
		tbl.tokens--
		tbl.allowed++
		return true
	}

	tbl.denied++
	return false
}

// Wait blocks until a token is available or context is cancelled
func (tbl *TokenBucketLimiter) Wait(ctx context.Context) error {
	for {
		if tbl.Allow() {
			return nil
		}

		// Calculate time to next token
		tbl.mu.Lock()
		nextRefill := tbl.lastRefill.Add(tbl.refillRate)
		tbl.mu.Unlock()

		waitTime := time.Until(nextRefill)
		if waitTime <= 0 {
			continue // Try again immediately
		}

		select {
		case <-time.After(waitTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Reset resets the rate limiter to initial state
func (tbl *TokenBucketLimiter) Reset() {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	tbl.tokens = tbl.capacity
	tbl.lastRefill = time.Now()
	tbl.allowed = 0
	tbl.denied = 0

	tbl.logger.Debug("Token bucket rate limiter reset")
}

// GetStats returns current statistics
func (tbl *TokenBucketLimiter) GetStats() RateLimiterStats {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	tbl.refillTokens()

	return RateLimiterStats{
		RequestsAllowed: tbl.allowed,
		RequestsDenied:  tbl.denied,
		CurrentTokens:   tbl.tokens,
		RefillRate:      float64(time.Second) / float64(tbl.refillRate),
		LastRefill:      tbl.lastRefill,
	}
}

// refillTokens adds tokens based on elapsed time (must be called with mutex held)
func (tbl *TokenBucketLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(tbl.lastRefill)
	
	if elapsed >= tbl.refillRate {
		tokensToAdd := int(elapsed / tbl.refillRate)
		tbl.tokens = min(tbl.capacity, tbl.tokens+tokensToAdd)
		tbl.lastRefill = tbl.lastRefill.Add(time.Duration(tokensToAdd) * tbl.refillRate)
	}
}

// SlidingWindowLimiter implements rate limiting using sliding window algorithm
type SlidingWindowLimiter struct {
	mu         sync.Mutex
	requests   []time.Time   // Request timestamps
	window     time.Duration // Time window
	limit      int           // Max requests in window
	allowed    int64         // Number of requests allowed
	denied     int64         // Number of requests denied
	logger     *zap.Logger
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter
func NewSlidingWindowLimiter(limit int, window time.Duration, logger *zap.Logger) *SlidingWindowLimiter {
	if limit <= 0 {
		panic("limit must be positive")
	}
	if window <= 0 {
		panic("window must be positive")
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	return &SlidingWindowLimiter{
		requests: make([]time.Time, 0),
		window:   window,
		limit:    limit,
		logger:   logger,
	}
}

// Allow checks if a request is allowed
func (swl *SlidingWindowLimiter) Allow() bool {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	now := time.Now()
	swl.cleanOldRequests(now)

	if len(swl.requests) < swl.limit {
		swl.requests = append(swl.requests, now)
		swl.allowed++
		return true
	}

	swl.denied++
	return false
}

// Wait blocks until a request is allowed
func (swl *SlidingWindowLimiter) Wait(ctx context.Context) error {
	for {
		if swl.Allow() {
			return nil
		}

		swl.mu.Lock()
		if len(swl.requests) == 0 {
			swl.mu.Unlock()
			continue
		}
		
		// Wait until the oldest request expires
		oldestRequest := swl.requests[0]
		waitTime := swl.window - time.Since(oldestRequest)
		swl.mu.Unlock()

		if waitTime <= 0 {
			continue
		}

		select {
		case <-time.After(waitTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Reset resets the rate limiter
func (swl *SlidingWindowLimiter) Reset() {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	swl.requests = swl.requests[:0]
	swl.allowed = 0
	swl.denied = 0

	swl.logger.Debug("Sliding window rate limiter reset")
}

// GetStats returns current statistics
func (swl *SlidingWindowLimiter) GetStats() RateLimiterStats {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	now := time.Now()
	swl.cleanOldRequests(now)

	return RateLimiterStats{
		RequestsAllowed: swl.allowed,
		RequestsDenied:  swl.denied,
		CurrentTokens:   swl.limit - len(swl.requests), // Available "slots"
		RefillRate:      float64(swl.limit) / swl.window.Seconds(),
		LastRefill:      time.Now(), // Not applicable, but provide current time
	}
}

// cleanOldRequests removes requests outside the time window
func (swl *SlidingWindowLimiter) cleanOldRequests(now time.Time) {
	cutoff := now.Add(-swl.window)
	
	// Find first request within window
	i := 0
	for i < len(swl.requests) && swl.requests[i].Before(cutoff) {
		i++
	}
	
	// Remove old requests
	if i > 0 {
		copy(swl.requests, swl.requests[i:])
		swl.requests = swl.requests[:len(swl.requests)-i]
	}
}

// RateLimitedExecutor provides rate-limited execution of functions
type RateLimitedExecutor struct {
	limiter RateLimiter
	name    string
	logger  *zap.Logger
}

// NewRateLimitedExecutor creates a new rate-limited executor
func NewRateLimitedExecutor(name string, limiter RateLimiter, logger *zap.Logger) *RateLimitedExecutor {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &RateLimitedExecutor{
		limiter: limiter,
		name:    name,
		logger:  logger,
	}
}

// Execute executes a function with rate limiting
func (rle *RateLimitedExecutor) Execute(ctx context.Context, fn func() error) error {
	if err := rle.limiter.Wait(ctx); err != nil {
		rle.logger.Warn("Rate limiter wait cancelled",
			zap.String("executor", rle.name),
			zap.Error(err))
		return fmt.Errorf("rate limiter wait failed: %w", err)
	}

	rle.logger.Debug("Executing rate-limited function", zap.String("executor", rle.name))
	return fn()
}

// TryExecute tries to execute a function without blocking
func (rle *RateLimitedExecutor) TryExecute(fn func() error) error {
	if !rle.limiter.Allow() {
		rle.logger.Debug("Rate limit exceeded, request denied", zap.String("executor", rle.name))
		return fmt.Errorf("rate limit exceeded for executor: %s", rle.name)
	}

	rle.logger.Debug("Executing rate-limited function", zap.String("executor", rle.name))
	return fn()
}

// GetStats returns rate limiter statistics
func (rle *RateLimitedExecutor) GetStats() RateLimiterStats {
	return rle.limiter.GetStats()
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}