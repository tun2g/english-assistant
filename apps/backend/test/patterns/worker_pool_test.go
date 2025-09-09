package patterns_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"app-backend/pkg/patterns"
	"go.uber.org/zap"
)

func TestWorkerPool(t *testing.T) {
	logger := zap.NewNop()

	t.Run("basic job processing", func(t *testing.T) {
		config := patterns.WorkerPoolConfig{
			WorkerCount:   2,
			QueueSize:     10,
			Timeout:       5 * time.Second,
			EnableMetrics: true,
			Logger:        logger,
		}

		pool := patterns.NewWorkerPool[string, string](config)
		pool.Start()
		defer pool.Stop()

		// Submit jobs
		jobs := []patterns.Job[string, string]{
			{
				ID:   "job1",
				Data: "hello",
				Process: func(ctx context.Context, data string) (string, error) {
					return strings.ToUpper(data), nil
				},
			},
			{
				ID:   "job2", 
				Data: "world",
				Process: func(ctx context.Context, data string) (string, error) {
					return strings.ToUpper(data), nil
				},
			},
		}

		for _, job := range jobs {
			err := pool.Submit(job)
			if err != nil {
				t.Fatalf("Failed to submit job: %v", err)
			}
		}

		// Collect results
		results := make(map[string]string)
		timeout := time.After(10 * time.Second)
		
		for len(results) < len(jobs) {
			select {
			case result := <-pool.Results():
				if result.Error != nil {
					t.Errorf("Job %s failed: %v", result.JobID, result.Error)
				} else {
					results[result.JobID] = result.Data
				}
			case <-timeout:
				t.Fatal("Timeout waiting for results")
			}
		}

		// Verify results
		expected := map[string]string{
			"job1": "HELLO",
			"job2": "WORLD",
		}
		
		for jobID, expected := range expected {
			if actual := results[jobID]; actual != expected {
				t.Errorf("Job %s: expected %s, got %s", jobID, expected, actual)
			}
		}

		// Check metrics
		metrics := pool.GetMetrics()
		if metrics.JobsProcessed != 2 {
			t.Errorf("Expected 2 jobs processed, got %d", metrics.JobsProcessed)
		}
		if metrics.JobsSucceeded != 2 {
			t.Errorf("Expected 2 jobs succeeded, got %d", metrics.JobsSucceeded)
		}
	})

	t.Run("submit and wait", func(t *testing.T) {
		config := patterns.WorkerPoolConfig{
			WorkerCount: 1,
			QueueSize:   5,
			Timeout:     5 * time.Second,
			Logger:      logger,
		}

		pool := patterns.NewWorkerPool[int, int](config)
		pool.Start()
		defer pool.Stop()

		job := patterns.Job[int, int]{
			ID:   "test",
			Data: 42,
			Process: func(ctx context.Context, data int) (int, error) {
				return data * 2, nil
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := pool.SubmitAndWait(ctx, job)
		if err != nil {
			t.Fatalf("SubmitAndWait failed: %v", err)
		}

		if result.Error != nil {
			t.Fatalf("Job failed: %v", result.Error)
		}

		if result.Data != 84 {
			t.Errorf("Expected result 84, got %d", result.Data)
		}
	})

	t.Run("concurrent workers", func(t *testing.T) {
		config := patterns.WorkerPoolConfig{
			WorkerCount:   5,
			QueueSize:     100,
			Timeout:       5 * time.Second,
			EnableMetrics: true,
			Logger:        logger,
		}

		pool := patterns.NewWorkerPool[int, int](config)
		pool.Start()
		defer pool.Stop()

		// Submit many jobs
		numJobs := 50
		
		for i := 0; i < numJobs; i++ {
			job := patterns.Job[int, int]{
				ID:   fmt.Sprintf("job-%d", i),
				Data: i,
				Process: func(ctx context.Context, data int) (int, error) {
					// Simulate some work
					time.Sleep(10 * time.Millisecond)
					return data * 2, nil
				},
			}

			err := pool.Submit(job)
			if err != nil {
				t.Fatalf("Failed to submit job %d: %v", i, err)
			}
		}

		// Collect all results
		results := make([]patterns.Result[int], 0, numJobs)
		timeout := time.After(30 * time.Second)
		
		for len(results) < numJobs {
			select {
			case result := <-pool.Results():
				results = append(results, result)
			case <-timeout:
				t.Fatalf("Timeout waiting for results, got %d/%d", len(results), numJobs)
			}
		}

		// Verify all jobs completed
		metrics := pool.GetMetrics()
		if metrics.JobsProcessed != int64(numJobs) {
			t.Errorf("Expected %d jobs processed, got %d", numJobs, metrics.JobsProcessed)
		}
	})
}

func BenchmarkWorkerPool(b *testing.B) {
	logger := zap.NewNop()
	config := patterns.WorkerPoolConfig{
		WorkerCount: 4,
		QueueSize:   1000,
		Timeout:     10 * time.Second,
		Logger:      logger,
	}

	pool := patterns.NewWorkerPool[int, int](config)
	pool.Start()
	defer pool.Stop()

	// Start a goroutine to consume results
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < b.N; i++ {
			<-pool.Results()
		}
	}()

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		job := patterns.Job[int, int]{
			ID:   fmt.Sprintf("bench-%d", i),
			Data: i,
			Process: func(ctx context.Context, data int) (int, error) {
				return data * 2, nil
			},
		}
		
		_ = pool.Submit(job)
	}
	
	// Wait for all jobs to complete
	<-done
}