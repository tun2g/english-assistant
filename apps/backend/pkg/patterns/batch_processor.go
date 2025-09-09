package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// BatchItem represents a single item to be processed in a batch
type BatchItem[T any] struct {
	ID   string
	Data T
}

// BatchResult represents the result of processing a batch item
type BatchResult[R any] struct {
	ID     string
	Data   R
	Error  error
}

// BatchProcessor processes items in batches for efficiency
type BatchProcessor[T any, R any] struct {
	config       BatchProcessorConfig
	inputChan    chan BatchItem[T]
	resultChan   chan BatchResult[R]
	processorFn  func(ctx context.Context, items []BatchItem[T]) ([]BatchResult[R], error)
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	metrics      *BatchProcessorMetrics
	logger       *zap.Logger
}

// BatchProcessorConfig holds configuration for the batch processor
type BatchProcessorConfig struct {
	BatchSize      int           // Maximum items per batch
	FlushInterval  time.Duration // Time to wait before processing partial batch
	MaxWorkers     int           // Number of worker goroutines
	InputBuffer    int           // Size of input channel buffer
	ResultBuffer   int           // Size of result channel buffer
	ProcessTimeout time.Duration // Timeout for processing each batch
	Logger         *zap.Logger   // Logger instance
}

// BatchProcessorMetrics holds metrics for the batch processor
type BatchProcessorMetrics struct {
	mu                    sync.RWMutex
	TotalItemsProcessed   int64
	TotalBatchesProcessed int64
	TotalItemsSucceeded   int64
	TotalItemsFailed      int64
	AverageBatchSize      float64
	AverageProcessingTime time.Duration
	totalProcessingTime   time.Duration
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor[T any, R any](
	config BatchProcessorConfig,
	processorFn func(ctx context.Context, items []BatchItem[T]) ([]BatchResult[R], error),
) *BatchProcessor[T, R] {
	// Set defaults
	if config.BatchSize <= 0 {
		config.BatchSize = 10
	}
	if config.FlushInterval <= 0 {
		config.FlushInterval = 5 * time.Second
	}
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = 3
	}
	if config.InputBuffer <= 0 {
		config.InputBuffer = 100
	}
	if config.ResultBuffer <= 0 {
		config.ResultBuffer = 100
	}
	if config.ProcessTimeout <= 0 {
		config.ProcessTimeout = 30 * time.Second
	}
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}

	ctx, cancel := context.WithCancel(context.Background())

	bp := &BatchProcessor[T, R]{
		config:      config,
		inputChan:   make(chan BatchItem[T], config.InputBuffer),
		resultChan:  make(chan BatchResult[R], config.ResultBuffer),
		processorFn: processorFn,
		ctx:         ctx,
		cancel:      cancel,
		metrics:     &BatchProcessorMetrics{},
		logger:      config.Logger,
	}

	return bp
}

// Start starts the batch processor workers
func (bp *BatchProcessor[T, R]) Start() {
	bp.logger.Info("Starting batch processor",
		zap.Int("batch_size", bp.config.BatchSize),
		zap.Duration("flush_interval", bp.config.FlushInterval),
		zap.Int("workers", bp.config.MaxWorkers))

	// Start worker goroutines
	for i := 0; i < bp.config.MaxWorkers; i++ {
		bp.wg.Add(1)
		go bp.worker(i)
	}
}

// Submit submits an item for batch processing
func (bp *BatchProcessor[T, R]) Submit(item BatchItem[T]) error {
	select {
	case bp.inputChan <- item:
		return nil
	case <-bp.ctx.Done():
		return fmt.Errorf("batch processor is shutting down")
	default:
		return fmt.Errorf("input buffer is full")
	}
}

// Results returns the channel for consuming processed results
func (bp *BatchProcessor[T, R]) Results() <-chan BatchResult[R] {
	return bp.resultChan
}

// Stop gracefully stops the batch processor
func (bp *BatchProcessor[T, R]) Stop() {
	bp.logger.Info("Stopping batch processor")
	bp.cancel()
	close(bp.inputChan)
	bp.wg.Wait()
	close(bp.resultChan)
	
	metrics := bp.GetMetrics()
	bp.logger.Info("Batch processor stopped",
		zap.Int64("total_items", metrics.TotalItemsProcessed),
		zap.Int64("total_batches", metrics.TotalBatchesProcessed),
		zap.Float64("avg_batch_size", metrics.AverageBatchSize),
		zap.Duration("avg_processing_time", metrics.AverageProcessingTime))
}

// GetMetrics returns current batch processor metrics
func (bp *BatchProcessor[T, R]) GetMetrics() BatchProcessorMetrics {
	bp.metrics.mu.RLock()
	defer bp.metrics.mu.RUnlock()

	metrics := *bp.metrics
	if metrics.TotalBatchesProcessed > 0 {
		metrics.AverageBatchSize = float64(metrics.TotalItemsProcessed) / float64(metrics.TotalBatchesProcessed)
		metrics.AverageProcessingTime = metrics.totalProcessingTime / time.Duration(metrics.TotalBatchesProcessed)
	}
	return metrics
}

// worker processes batches of items
func (bp *BatchProcessor[T, R]) worker(workerID int) {
	defer bp.wg.Done()
	
	workerLogger := bp.logger.With(zap.Int("worker_id", workerID))
	workerLogger.Debug("Batch processor worker started")
	defer workerLogger.Debug("Batch processor worker stopped")

	batch := make([]BatchItem[T], 0, bp.config.BatchSize)
	flushTimer := time.NewTimer(bp.config.FlushInterval)
	defer flushTimer.Stop()

	for {
		select {
		case item, ok := <-bp.inputChan:
			if !ok {
				// Input channel closed, process remaining batch
				if len(batch) > 0 {
					bp.processBatch(workerLogger, batch)
				}
				return
			}

			batch = append(batch, item)

			// Process batch if it's full
			if len(batch) >= bp.config.BatchSize {
				bp.processBatch(workerLogger, batch)
				batch = batch[:0] // Reset batch
				flushTimer.Reset(bp.config.FlushInterval)
			}

		case <-flushTimer.C:
			// Flush interval reached, process partial batch
			if len(batch) > 0 {
				bp.processBatch(workerLogger, batch)
				batch = batch[:0] // Reset batch
			}
			flushTimer.Reset(bp.config.FlushInterval)

		case <-bp.ctx.Done():
			// Context cancelled, process remaining batch
			if len(batch) > 0 {
				bp.processBatch(workerLogger, batch)
			}
			return
		}
	}
}

// processBatch processes a batch of items
func (bp *BatchProcessor[T, R]) processBatch(logger *zap.Logger, batch []BatchItem[T]) {
	if len(batch) == 0 {
		return
	}

	start := time.Now()
	batchSize := len(batch)

	logger.Debug("Processing batch", zap.Int("size", batchSize))

	// Create timeout context for batch processing
	ctx, cancel := context.WithTimeout(bp.ctx, bp.config.ProcessTimeout)
	defer cancel()

	// Process the batch
	results, err := bp.processorFn(ctx, batch)
	processingTime := time.Since(start)

	// Update metrics
	bp.updateMetrics(batchSize, processingTime, results, err)

	if err != nil {
		logger.Error("Batch processing failed", 
			zap.Int("batch_size", batchSize),
			zap.Duration("processing_time", processingTime),
			zap.Error(err))
		
		// Create error results for all items in batch
		for _, item := range batch {
			result := BatchResult[R]{
				ID:    item.ID,
				Error: err,
			}
			bp.sendResult(result)
		}
		return
	}

	// Send individual results
	for _, result := range results {
		bp.sendResult(result)
	}

	logger.Debug("Batch processed successfully",
		zap.Int("batch_size", batchSize),
		zap.Int("results", len(results)),
		zap.Duration("processing_time", processingTime))
}

// sendResult sends a result to the result channel
func (bp *BatchProcessor[T, R]) sendResult(result BatchResult[R]) {
	select {
	case bp.resultChan <- result:
	case <-bp.ctx.Done():
		bp.logger.Warn("Failed to send result, processor shutting down", zap.String("item_id", result.ID))
	}
}

// updateMetrics updates the batch processor metrics
func (bp *BatchProcessor[T, R]) updateMetrics(batchSize int, processingTime time.Duration, results []BatchResult[R], batchErr error) {
	bp.metrics.mu.Lock()
	defer bp.metrics.mu.Unlock()

	bp.metrics.TotalBatchesProcessed++
	bp.metrics.TotalItemsProcessed += int64(batchSize)
	bp.metrics.totalProcessingTime += processingTime

	if batchErr != nil {
		bp.metrics.TotalItemsFailed += int64(batchSize)
	} else {
		// Count individual successes/failures in results
		for _, result := range results {
			if result.Error != nil {
				bp.metrics.TotalItemsFailed++
			} else {
				bp.metrics.TotalItemsSucceeded++
			}
		}
	}
}

// AsyncBatchProcessor provides a higher-level async interface
type AsyncBatchProcessor[T any, R any] struct {
	*BatchProcessor[T, R]
	pendingResults map[string]chan BatchResult[R]
	resultsMu      sync.RWMutex
	resultProcessor *WorkerPool[BatchResult[R], struct{}]
}

// NewAsyncBatchProcessor creates a new async batch processor
func NewAsyncBatchProcessor[T any, R any](
	config BatchProcessorConfig,
	processorFn func(ctx context.Context, items []BatchItem[T]) ([]BatchResult[R], error),
) *AsyncBatchProcessor[T, R] {
	bp := NewBatchProcessor(config, processorFn)
	
	abp := &AsyncBatchProcessor[T, R]{
		BatchProcessor: bp,
		pendingResults: make(map[string]chan BatchResult[R]),
	}

	// Create worker pool for result processing
	resultConfig := WorkerPoolConfig{
		WorkerCount: 2,
		QueueSize:   100,
		Timeout:     5 * time.Second,
		Logger:      config.Logger,
	}
	
	abp.resultProcessor = NewWorkerPool[BatchResult[R], struct{}](resultConfig)
	
	return abp
}

// Start starts the async batch processor
func (abp *AsyncBatchProcessor[T, R]) Start() {
	abp.BatchProcessor.Start()
	abp.resultProcessor.Start()
	
	// Start result router
	go abp.routeResults()
}

// SubmitAsync submits an item and returns a channel for the result
func (abp *AsyncBatchProcessor[T, R]) SubmitAsync(item BatchItem[T]) (<-chan BatchResult[R], error) {
	resultChan := make(chan BatchResult[R], 1)
	
	abp.resultsMu.Lock()
	abp.pendingResults[item.ID] = resultChan
	abp.resultsMu.Unlock()
	
	err := abp.Submit(item)
	if err != nil {
		abp.resultsMu.Lock()
		delete(abp.pendingResults, item.ID)
		abp.resultsMu.Unlock()
		close(resultChan)
		return nil, err
	}
	
	return resultChan, nil
}

// SubmitAndWait submits an item and waits for the result
func (abp *AsyncBatchProcessor[T, R]) SubmitAndWait(ctx context.Context, item BatchItem[T]) (BatchResult[R], error) {
	resultChan, err := abp.SubmitAsync(item)
	if err != nil {
		return BatchResult[R]{}, err
	}
	
	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		// Clean up pending result
		abp.resultsMu.Lock()
		delete(abp.pendingResults, item.ID)
		abp.resultsMu.Unlock()
		return BatchResult[R]{}, ctx.Err()
	}
}

// Stop stops the async batch processor
func (abp *AsyncBatchProcessor[T, R]) Stop() {
	abp.BatchProcessor.Stop()
	abp.resultProcessor.Stop()
	
	// Close all pending result channels
	abp.resultsMu.Lock()
	for _, ch := range abp.pendingResults {
		close(ch)
	}
	abp.pendingResults = make(map[string]chan BatchResult[R])
	abp.resultsMu.Unlock()
}

// routeResults routes batch processing results to waiting callers
func (abp *AsyncBatchProcessor[T, R]) routeResults() {
	for result := range abp.Results() {
		abp.resultsMu.Lock()
		if ch, exists := abp.pendingResults[result.ID]; exists {
			select {
			case ch <- result:
			default:
				abp.logger.Warn("Failed to send result to waiting channel", zap.String("item_id", result.ID))
			}
			close(ch)
			delete(abp.pendingResults, result.ID)
		}
		abp.resultsMu.Unlock()
	}
}