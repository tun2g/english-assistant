package patterns_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"app-backend/pkg/patterns"
	"go.uber.org/zap"
)

func TestCircuitBreaker(t *testing.T) {
	logger := zap.NewNop()

	t.Run("closed to open transition", func(t *testing.T) {
		config := patterns.CircuitBreakerConfig{
			Name:             "test",
			FailureThreshold: 3,
			SuccessThreshold: 2,
			Timeout:          1 * time.Second,
			Interval:         5 * time.Second,
			Logger:           logger,
		}

		cb := patterns.NewCircuitBreaker(config)
		ctx := context.Background()

		// Circuit should start in closed state
		if cb.GetState() != patterns.StateClosed {
			t.Errorf("Expected initial state to be closed, got %v", cb.GetState())
		}

		// Simulate failures to trip the circuit
		for i := 0; i < 3; i++ {
			err := cb.Execute(ctx, func() error {
				return fmt.Errorf("failure %d", i+1)
			})
			if err == nil {
				t.Error("Expected error from failing function")
			}
		}

		// Circuit should now be open
		if cb.GetState() != patterns.StateOpen {
			t.Errorf("Expected state to be open after failures, got %v", cb.GetState())
		}

		// Calls should be rejected immediately
		err := cb.Execute(ctx, func() error {
			return nil
		})
		if !patterns.IsCircuitBreakerError(err) {
			t.Errorf("Expected circuit breaker error, got %v", err)
		}
	})

	t.Run("open to half-open transition", func(t *testing.T) {
		config := patterns.CircuitBreakerConfig{
			Name:             "test-halfopen",
			FailureThreshold: 2,
			SuccessThreshold: 1,
			Timeout:          100 * time.Millisecond, // Short timeout for testing
			Interval:         1 * time.Second,
			Logger:           logger,
		}

		cb := patterns.NewCircuitBreaker(config)
		ctx := context.Background()

		// Trip the circuit
		for i := 0; i < 2; i++ {
			cb.Execute(ctx, func() error {
				return fmt.Errorf("failure")
			})
		}

		// Should be open
		if cb.GetState() != patterns.StateOpen {
			t.Error("Expected circuit to be open")
		}

		// Wait for timeout to transition to half-open
		time.Sleep(150 * time.Millisecond)

		// Next call should transition to half-open
		err := cb.Execute(ctx, func() error {
			return nil // Success
		})

		if err != nil {
			t.Errorf("Expected successful call in half-open state, got %v", err)
		}

		// Should transition back to closed
		if cb.GetState() != patterns.StateClosed {
			t.Errorf("Expected state to be closed after success, got %v", cb.GetState())
		}
	})

	t.Run("execute with fallback", func(t *testing.T) {
		config := patterns.CircuitBreakerConfig{
			Name:             "test-fallback",
			FailureThreshold: 1,
			Timeout:          100 * time.Millisecond,
			Logger:           logger,
		}

		cb := patterns.NewCircuitBreaker(config)
		ctx := context.Background()

		// Trip the circuit
		cb.Execute(ctx, func() error {
			return fmt.Errorf("failure")
		})

		fallbackCalled := false
		err := cb.ExecuteWithFallback(ctx, 
			func() error {
				return fmt.Errorf("main function")
			},
			func() error {
				fallbackCalled = true
				return nil
			},
		)

		if err != nil {
			t.Errorf("Expected fallback to succeed, got %v", err)
		}

		if !fallbackCalled {
			t.Error("Expected fallback to be called")
		}
	})

	t.Run("metrics collection", func(t *testing.T) {
		config := patterns.CircuitBreakerConfig{
			Name:             "test-metrics",
			FailureThreshold: 5,
			Logger:           logger,
		}

		cb := patterns.NewCircuitBreaker(config)
		ctx := context.Background()

		// Execute some successful calls
		for i := 0; i < 3; i++ {
			cb.Execute(ctx, func() error {
				return nil
			})
		}

		// Execute some failed calls
		for i := 0; i < 2; i++ {
			cb.Execute(ctx, func() error {
				return fmt.Errorf("failure")
			})
		}

		metrics := cb.GetMetrics()
		
		if metrics.TotalSuccesses != 3 {
			t.Errorf("Expected 3 successes, got %d", metrics.TotalSuccesses)
		}
		
		if metrics.TotalFailures != 2 {
			t.Errorf("Expected 2 failures, got %d", metrics.TotalFailures)
		}
		
		expectedRate := 2.0 / 5.0 // 2 failures out of 5 total
		if metrics.FailureRate != expectedRate {
			t.Errorf("Expected failure rate %.2f, got %.2f", expectedRate, metrics.FailureRate)
		}
	})

	t.Run("reset functionality", func(t *testing.T) {
		config := patterns.CircuitBreakerConfig{
			Name:             "test-reset",
			FailureThreshold: 2,
			Logger:           logger,
		}

		cb := patterns.NewCircuitBreaker(config)
		ctx := context.Background()

		// Trip the circuit
		for i := 0; i < 2; i++ {
			cb.Execute(ctx, func() error {
				return fmt.Errorf("failure")
			})
		}

		// Should be open
		if cb.GetState() != patterns.StateOpen {
			t.Error("Expected circuit to be open")
		}

		// Reset the circuit
		cb.Reset()

		// Should be closed
		if cb.GetState() != patterns.StateClosed {
			t.Error("Expected circuit to be closed after reset")
		}

		// Should accept calls again
		err := cb.Execute(ctx, func() error {
			return nil
		})

		if err != nil {
			t.Errorf("Expected call to succeed after reset, got %v", err)
		}
	})
}

func BenchmarkCircuitBreaker(b *testing.B) {
	logger := zap.NewNop()
	config := patterns.CircuitBreakerConfig{
		Name:             "bench",
		FailureThreshold: 1000, // High threshold so it doesn't trip
		Logger:           logger,
	}

	cb := patterns.NewCircuitBreaker(config)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb.Execute(ctx, func() error {
				return nil
			})
		}
	})
}