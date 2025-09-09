package patterns_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"app-backend/pkg/patterns"
)

func TestSemaphore(t *testing.T) {
	sem := patterns.NewSemaphore(2)

	t.Run("basic acquire and release", func(t *testing.T) {
		ctx := context.Background()
		
		// Should be able to acquire
		err := sem.Acquire(ctx)
		if err != nil {
			t.Fatalf("Expected to acquire semaphore, got error: %v", err)
		}
		defer sem.Release()
		
		// Check available permits
		if sem.AvailablePermits() != 1 {
			t.Errorf("Expected 1 available permit, got %d", sem.AvailablePermits())
		}
	})

	t.Run("try acquire non-blocking", func(t *testing.T) {
		sem2 := patterns.NewSemaphore(1)
		
		// First acquire should succeed
		if !sem2.TryAcquire() {
			t.Error("Expected TryAcquire to succeed")
		}
		
		// Second should fail
		if sem2.TryAcquire() {
			sem2.Release() // cleanup if it somehow succeeded
			t.Error("Expected TryAcquire to fail when semaphore is full")
		}
		
		sem2.Release()
	})

	t.Run("with semaphore helper", func(t *testing.T) {
		sem3 := patterns.NewSemaphore(1)
		ctx := context.Background()
		
		executed := false
		err := sem3.WithSemaphore(ctx, func() error {
			executed = true
			return nil
		})
		
		if err != nil {
			t.Fatalf("WithSemaphore failed: %v", err)
		}
		
		if !executed {
			t.Error("Expected function to be executed")
		}
		
		if sem3.AvailablePermits() != 1 {
			t.Error("Expected semaphore to be released after WithSemaphore")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		sem4 := patterns.NewSemaphore(1)
		
		// Acquire the semaphore
		ctx := context.Background()
		err := sem4.Acquire(ctx)
		if err != nil {
			t.Fatalf("Failed to acquire semaphore: %v", err)
		}
		
		// Try to acquire with cancelled context
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		
		err = sem4.Acquire(cancelledCtx)
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
		
		sem4.Release()
	})

	t.Run("concurrent access", func(t *testing.T) {
		sem5 := patterns.NewSemaphore(3)
		var wg sync.WaitGroup
		ctx := context.Background()
		
		// Start 5 goroutines trying to acquire 3 permits
		acquired := make(chan bool, 5)
		
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				err := sem5.Acquire(ctx)
				if err != nil {
					acquired <- false
					return
				}
				
				acquired <- true
				time.Sleep(100 * time.Millisecond)
				sem5.Release()
			}()
		}
		
		wg.Wait()
		close(acquired)
		
		successCount := 0
		for success := range acquired {
			if success {
				successCount++
			}
		}
		
		if successCount != 5 {
			t.Errorf("Expected all 5 goroutines to eventually succeed, got %d", successCount)
		}
	})
}

func BenchmarkSemaphore(b *testing.B) {
	sem := patterns.NewSemaphore(100)
	ctx := context.Background()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = sem.Acquire(ctx)
			sem.Release()
		}
	})
}