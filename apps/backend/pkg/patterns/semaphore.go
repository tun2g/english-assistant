package patterns

import (
	"context"
	"fmt"
)

// Semaphore controls access to a resource by limiting concurrent goroutines
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a new semaphore with the given capacity
func NewSemaphore(capacity int) *Semaphore {
	if capacity <= 0 {
		panic("semaphore capacity must be positive")
	}
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

// Acquire tries to acquire the semaphore. Blocks if no permits available.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire tries to acquire the semaphore without blocking
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release releases a permit back to the semaphore
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
		panic("semaphore: release called without acquire")
	}
}

// AvailablePermits returns the number of available permits
func (s *Semaphore) AvailablePermits() int {
	return cap(s.ch) - len(s.ch)
}

// WithSemaphore is a helper function that automatically handles acquire/release
func (s *Semaphore) WithSemaphore(ctx context.Context, fn func() error) error {
	if err := s.Acquire(ctx); err != nil {
		return fmt.Errorf("semaphore acquire failed: %w", err)
	}
	defer s.Release()
	return fn()
}