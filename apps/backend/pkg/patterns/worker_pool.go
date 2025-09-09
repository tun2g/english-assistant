package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Job represents a unit of work to be processed by the worker pool
type Job[T any, R any] struct {
	ID      string
	Data    T
	Process func(context.Context, T) (R, error)
}

// Result represents the result of processing a job
type Result[R any] struct {
	JobID  string
	Data   R
	Error  error
	Timing time.Duration
}

// WorkerPoolConfig holds configuration for the worker pool
type WorkerPoolConfig struct {
	WorkerCount    int           // Number of worker goroutines
	QueueSize      int           // Size of job queue buffer
	Timeout        time.Duration // Per-job timeout
	EnableMetrics  bool          // Whether to collect metrics
	Logger         *zap.Logger   // Logger instance
}

// WorkerPool implements a generic worker pool pattern
type WorkerPool[T any, R any] struct {
	config      WorkerPoolConfig
	jobs        chan Job[T, R]
	results     chan Result[R]
	workers     []Worker[T, R]
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	metrics     *WorkerPoolMetrics
	once        sync.Once
}

// Worker represents a single worker in the pool
type Worker[T any, R any] struct {
	ID      int
	pool    *WorkerPool[T, R]
	logger  *zap.Logger
}

// WorkerPoolMetrics holds metrics for the worker pool
type WorkerPoolMetrics struct {
	mu                sync.RWMutex
	JobsProcessed     int64
	JobsSucceeded     int64
	JobsFailed        int64
	AverageProcessingTime time.Duration
	totalProcessingTime   time.Duration
}

// NewWorkerPool creates a new worker pool with the given configuration
func NewWorkerPool[T any, R any](config WorkerPoolConfig) *WorkerPool[T, R] {
	if config.WorkerCount <= 0 {
		config.WorkerCount = 5
	}
	if config.QueueSize <= 0 {
		config.QueueSize = 100
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}

	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool[T, R]{
		config:  config,
		jobs:    make(chan Job[T, R], config.QueueSize),
		results: make(chan Result[R], config.QueueSize),
		ctx:     ctx,
		cancel:  cancel,
		metrics: &WorkerPoolMetrics{},
	}

	// Create workers
	wp.workers = make([]Worker[T, R], config.WorkerCount)
	for i := 0; i < config.WorkerCount; i++ {
		wp.workers[i] = Worker[T, R]{
			ID:     i,
			pool:   wp,
			logger: config.Logger.With(zap.Int("worker_id", i)),
		}
	}

	return wp
}

// Start starts all workers in the pool
func (wp *WorkerPool[T, R]) Start() {
	wp.once.Do(func() {
		wp.wg.Add(len(wp.workers))
		for i := range wp.workers {
			go wp.workers[i].run()
		}
		wp.config.Logger.Info("Worker pool started", 
			zap.Int("worker_count", len(wp.workers)),
			zap.Int("queue_size", wp.config.QueueSize))
	})
}

// Submit submits a job to the worker pool
func (wp *WorkerPool[T, R]) Submit(job Job[T, R]) error {
	select {
	case wp.jobs <- job:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	default:
		return fmt.Errorf("job queue is full")
	}
}

// SubmitAndWait submits a job and waits for the result
func (wp *WorkerPool[T, R]) SubmitAndWait(ctx context.Context, job Job[T, R]) (Result[R], error) {
	if err := wp.Submit(job); err != nil {
		return Result[R]{}, err
	}

	// Wait for result
	for {
		select {
		case result := <-wp.results:
			if result.JobID == job.ID {
				return result, nil
			}
			// Not our result, put it back (this is a limitation - in real use you'd need result routing)
			select {
			case wp.results <- result:
			case <-ctx.Done():
				return Result[R]{}, ctx.Err()
			}
		case <-ctx.Done():
			return Result[R]{}, ctx.Err()
		}
	}
}

// Results returns the results channel for consuming processed jobs
func (wp *WorkerPool[T, R]) Results() <-chan Result[R] {
	return wp.results
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool[T, R]) Stop() {
	wp.cancel()
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
	
	if wp.config.EnableMetrics {
		metrics := wp.GetMetrics()
		wp.config.Logger.Info("Worker pool stopped",
			zap.Int64("jobs_processed", metrics.JobsProcessed),
			zap.Int64("jobs_succeeded", metrics.JobsSucceeded),
			zap.Int64("jobs_failed", metrics.JobsFailed),
			zap.Duration("avg_processing_time", metrics.AverageProcessingTime))
	}
}

// GetMetrics returns current worker pool metrics
func (wp *WorkerPool[T, R]) GetMetrics() WorkerPoolMetrics {
	wp.metrics.mu.RLock()
	defer wp.metrics.mu.RUnlock()
	
	metrics := *wp.metrics
	if metrics.JobsProcessed > 0 {
		metrics.AverageProcessingTime = metrics.totalProcessingTime / time.Duration(metrics.JobsProcessed)
	}
	return metrics
}

// run starts the worker's processing loop
func (w *Worker[T, R]) run() {
	defer w.pool.wg.Done()
	
	w.logger.Debug("Worker started")
	defer w.logger.Debug("Worker stopped")

	for {
		select {
		case job, ok := <-w.pool.jobs:
			if !ok {
				return // Channel closed, worker should exit
			}
			w.processJob(job)
		case <-w.pool.ctx.Done():
			return
		}
	}
}

// processJob processes a single job
func (w *Worker[T, R]) processJob(job Job[T, R]) {
	start := time.Now()
	
	// Create timeout context for this job
	ctx, cancel := context.WithTimeout(w.pool.ctx, w.pool.config.Timeout)
	defer cancel()

	w.logger.Debug("Processing job", zap.String("job_id", job.ID))

	// Process the job
	data, err := job.Process(ctx, job.Data)
	
	processingTime := time.Since(start)
	
	result := Result[R]{
		JobID:  job.ID,
		Data:   data,
		Error:  err,
		Timing: processingTime,
	}

	// Update metrics
	if w.pool.config.EnableMetrics {
		w.pool.updateMetrics(result)
	}

	// Send result
	select {
	case w.pool.results <- result:
		if err != nil {
			w.logger.Error("Job failed", 
				zap.String("job_id", job.ID), 
				zap.Duration("processing_time", processingTime),
				zap.Error(err))
		} else {
			w.logger.Debug("Job completed successfully", 
				zap.String("job_id", job.ID),
				zap.Duration("processing_time", processingTime))
		}
	case <-w.pool.ctx.Done():
		w.logger.Warn("Failed to send job result, pool shutting down", zap.String("job_id", job.ID))
		return
	}
}

// updateMetrics updates the worker pool metrics
func (wp *WorkerPool[T, R]) updateMetrics(result Result[R]) {
	wp.metrics.mu.Lock()
	defer wp.metrics.mu.Unlock()
	
	wp.metrics.JobsProcessed++
	wp.metrics.totalProcessingTime += result.Timing
	
	if result.Error != nil {
		wp.metrics.JobsFailed++
	} else {
		wp.metrics.JobsSucceeded++
	}
}

// WorkerPoolMetrics getter methods
func (m *WorkerPoolMetrics) GetJobsProcessed() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.JobsProcessed
}

func (m *WorkerPoolMetrics) GetJobsSucceeded() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.JobsSucceeded
}

func (m *WorkerPoolMetrics) GetJobsFailed() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.JobsFailed
}

func (m *WorkerPoolMetrics) GetAverageProcessingTime() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.JobsProcessed == 0 {
		return 0
	}
	return m.totalProcessingTime / time.Duration(m.JobsProcessed)
}

