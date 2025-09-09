package patterns_test

import (
	"context"
	"testing"
	"time"

	"app-backend/pkg/patterns"
	"go.uber.org/zap"
)

func TestTokenBucketLimiter(t *testing.T) {
	logger := zap.NewNop()

	t.Run("basic rate limiting", func(t *testing.T) {
		// 3 tokens, refill every 100ms
		limiter := patterns.NewTokenBucketLimiter(3, 100*time.Millisecond, logger)

		// First 3 requests should be allowed
		for i := 0; i < 3; i++ {
			if !limiter.Allow() {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}

		// Next request should be denied
		if limiter.Allow() {
			t.Error("Request 4 should be denied")
		}

		// Wait for refill and try again
		time.Sleep(150 * time.Millisecond)
		if !limiter.Allow() {
			t.Error("Request should be allowed after refill")
		}
	})

	t.Run("wait for token", func(t *testing.T) {
		limiter := patterns.NewTokenBucketLimiter(1, 100*time.Millisecond, logger)

		// Use the token
		if !limiter.Allow() {
			t.Error("First request should be allowed")
		}

		// Wait should succeed after refill
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		err := limiter.Wait(ctx)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("Wait should succeed: %v", err)
		}

		if elapsed < 90*time.Millisecond {
			t.Error("Wait should take at least ~100ms for refill")
		}
	})

	t.Run("wait timeout", func(t *testing.T) {
		limiter := patterns.NewTokenBucketLimiter(1, 1*time.Second, logger)

		// Use the token
		limiter.Allow()

		// Wait with short timeout should fail
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := limiter.Wait(ctx)
		if err != context.DeadlineExceeded {
			t.Errorf("Expected timeout, got %v", err)
		}
	})

	t.Run("statistics", func(t *testing.T) {
		limiter := patterns.NewTokenBucketLimiter(2, 100*time.Millisecond, logger)

		// Allow some requests
		limiter.Allow()
		limiter.Allow()

		// Deny one request
		limiter.Allow() // This should fail

		stats := limiter.GetStats()
		if stats.RequestsAllowed != 2 {
			t.Errorf("Expected 2 allowed requests, got %d", stats.RequestsAllowed)
		}
		if stats.RequestsDenied != 1 {
			t.Errorf("Expected 1 denied request, got %d", stats.RequestsDenied)
		}
	})

	t.Run("reset", func(t *testing.T) {
		limiter := patterns.NewTokenBucketLimiter(2, 100*time.Millisecond, logger)

		// Use all tokens
		limiter.Allow()
		limiter.Allow()

		// Should be no tokens left
		if limiter.Allow() {
			t.Error("Should not have tokens available")
		}

		// Reset should restore all tokens
		limiter.Reset()
		
		if !limiter.Allow() {
			t.Error("Should have tokens after reset")
		}
	})
}

func TestSlidingWindowLimiter(t *testing.T) {
	logger := zap.NewNop()

	t.Run("basic window limiting", func(t *testing.T) {
		// 3 requests per 200ms window
		limiter := patterns.NewSlidingWindowLimiter(3, 200*time.Millisecond, logger)

		// First 3 should be allowed
		for i := 0; i < 3; i++ {
			if !limiter.Allow() {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}

		// 4th should be denied
		if limiter.Allow() {
			t.Error("Request 4 should be denied")
		}

		// After window expires, should be allowed again
		time.Sleep(250 * time.Millisecond)
		if !limiter.Allow() {
			t.Error("Request should be allowed after window expiry")
		}
	})

	t.Run("sliding behavior", func(t *testing.T) {
		// 2 requests per 200ms
		limiter := patterns.NewSlidingWindowLimiter(2, 200*time.Millisecond, logger)

		// Use 2 requests
		limiter.Allow()
		limiter.Allow()

		// Wait half the window
		time.Sleep(100 * time.Millisecond)

		// Should still be limited
		if limiter.Allow() {
			t.Error("Should still be limited in middle of window")
		}

		// Wait for full window to pass
		time.Sleep(150 * time.Millisecond) // Total 250ms > 200ms window

		// Should be allowed now
		if !limiter.Allow() {
			t.Error("Should be allowed after full window")
		}
	})
}

func TestRateLimitedExecutor(t *testing.T) {
	logger := zap.NewNop()

	t.Run("execute with rate limiting", func(t *testing.T) {
		limiter := patterns.NewTokenBucketLimiter(2, 100*time.Millisecond, logger)
		executor := patterns.NewRateLimitedExecutor("test", limiter, logger)

		executed := 0
		
		// First two executions should succeed immediately
		for i := 0; i < 2; i++ {
			err := executor.TryExecute(func() error {
				executed++
				return nil
			})
			if err != nil {
				t.Errorf("Execution %d should succeed: %v", i+1, err)
			}
		}

		// Third should fail (rate limited)
		err := executor.TryExecute(func() error {
			executed++
			return nil
		})
		if err == nil {
			t.Error("Third execution should be rate limited")
		}

		if executed != 2 {
			t.Errorf("Expected 2 executions, got %d", executed)
		}
	})

	t.Run("execute with waiting", func(t *testing.T) {
		limiter := patterns.NewTokenBucketLimiter(1, 100*time.Millisecond, logger)
		executor := patterns.NewRateLimitedExecutor("test", limiter, logger)

		// Use up the token
		executor.TryExecute(func() error { return nil })

		// This should wait for refill
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		executed := false
		start := time.Now()
		
		err := executor.Execute(ctx, func() error {
			executed = true
			return nil
		})
		
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("Execute should succeed after waiting: %v", err)
		}

		if !executed {
			t.Error("Function should have been executed")
		}

		if elapsed < 90*time.Millisecond {
			t.Error("Should have waited for token refill")
		}
	})
}

func BenchmarkTokenBucketLimiter(b *testing.B) {
	logger := zap.NewNop()
	limiter := patterns.NewTokenBucketLimiter(1000, time.Microsecond, logger)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}

func BenchmarkSlidingWindowLimiter(b *testing.B) {
	logger := zap.NewNop()
	limiter := patterns.NewSlidingWindowLimiter(1000, time.Second, logger)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}