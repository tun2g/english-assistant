# Go Concurrency Patterns

This package provides comprehensive, reusable concurrency patterns for Go applications. All patterns
are designed to be thread-safe, performant, and easily integrated into any codebase.

## Available Patterns

### 1. Semaphore (`semaphore.go`)

Controls access to resources by limiting concurrent goroutines.

```go
// Create semaphore with capacity of 5
sem := patterns.NewSemaphore(5)

// Basic usage
ctx := context.Background()
if err := sem.Acquire(ctx); err != nil {
    log.Fatal("Failed to acquire semaphore:", err)
}
defer sem.Release()

// Do work...

// Helper method with automatic cleanup
err := sem.WithSemaphore(ctx, func() error {
    // Do protected work
    return someOperation()
})
```

**Use Cases:**

- Limiting concurrent database connections
- Controlling API request rate
- Resource pool management

### 2. Worker Pool (`worker_pool.go`)

Generic worker pool with metrics and timeout support.

```go
// Configuration
config := patterns.WorkerPoolConfig{
    WorkerCount:    5,
    QueueSize:      100,
    Timeout:        30 * time.Second,
    EnableMetrics:  true,
    Logger:         logger,
}

// Create and start pool
pool := patterns.NewWorkerPool[string, string](config)
pool.Start()
defer pool.Stop()

// Submit jobs
job := patterns.Job[string, string]{
    ID:   "job1",
    Data: "input data",
    Process: func(ctx context.Context, data string) (string, error) {
        // Process the data
        return "processed: " + data, nil
    },
}

err := pool.Submit(job)

// Consume results
for result := range pool.Results() {
    if result.Error != nil {
        log.Printf("Job %s failed: %v", result.JobID, result.Error)
    } else {
        log.Printf("Job %s completed: %s", result.JobID, result.Data)
    }
}
```

### 3. Circuit Breaker (`circuit_breaker.go`)

Prevents cascading failures by temporarily blocking calls to failing services.

```go
config := patterns.CircuitBreakerConfig{
    Name:                "external-api",
    FailureThreshold:    5,
    SuccessThreshold:    3,
    Timeout:            60 * time.Second,
    Interval:           120 * time.Second,
    Logger:             logger,
}

cb := patterns.NewCircuitBreaker(config)

// Execute with protection
err := cb.Execute(ctx, func() error {
    return callExternalAPI()
})

if patterns.IsCircuitBreakerError(err) {
    // Handle circuit breaker open
    log.Println("Service is temporarily unavailable")
}

// Execute with fallback
err = cb.ExecuteWithFallback(ctx,
    func() error {
        return callExternalAPI()
    },
    func() error {
        return useCachedData()
    },
)
```

### 4. Rate Limiter (`rate_limiter.go`)

Two implementations: Token Bucket and Sliding Window algorithms.

```go
// Token Bucket Rate Limiter
limiter := patterns.NewTokenBucketLimiter(10, time.Second, logger)

// Non-blocking check
if limiter.Allow() {
    // Process request
}

// Blocking wait
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := limiter.Wait(ctx); err != nil {
    log.Println("Rate limit exceeded")
} else {
    // Process request
}

// Sliding Window Rate Limiter
windowLimiter := patterns.NewSlidingWindowLimiter(100, time.Minute, logger)

// Rate-Limited Executor
executor := patterns.NewRateLimitedExecutor("api-calls", limiter, logger)
err := executor.Execute(ctx, func() error {
    return makeAPICall()
})
```

### 5. Batch Processor (`batch_processor.go`)

Processes items in batches for efficiency, with async support.

```go
config := patterns.BatchProcessorConfig{
    BatchSize:      10,
    FlushInterval:  5 * time.Second,
    MaxWorkers:     3,
    InputBuffer:    100,
    Logger:         logger,
}

// Processor function
processorFn := func(ctx context.Context, items []patterns.BatchItem[string]) ([]patterns.BatchResult[string], error) {
    results := make([]patterns.BatchResult[string], len(items))
    for i, item := range items {
        // Process each item
        results[i] = patterns.BatchResult[string]{
            ID:   item.ID,
            Data: "processed: " + item.Data,
        }
    }
    return results, nil
}

// Create and start processor
processor := patterns.NewBatchProcessor(config, processorFn)
processor.Start()
defer processor.Stop()

// Submit items
item := patterns.BatchItem[string]{
    ID:   "item1",
    Data: "raw data",
}
err := processor.Submit(item)

// Process results
for result := range processor.Results() {
    log.Printf("Processed %s: %s", result.ID, result.Data)
}
```

### 6. Concurrent Map (`concurrent_map.go`)

Thread-safe map with sharding for high performance.

```go
// Create concurrent map
cm := patterns.NewConcurrentMap[string, int]()

// Basic operations
cm.Set("key1", 100)
value, exists := cm.Get("key1")
cm.Delete("key1")

// Advanced operations
value, wasSet := cm.GetOrSet("key2", 200)
value = cm.GetOrCompute("key3", func() int {
    return expensiveComputation()
})

// Atomic operations
cm.Update("key1", func(oldValue int) int {
    return oldValue + 1
})

success := cm.CompareAndSwap("key1", 100, 101, func(current, expected int) bool {
    return current == expected
})

// Iteration
cm.ForEach(func(key string, value int) bool {
    log.Printf("%s: %d", key, value)
    return true // continue iteration
})

// Filter
filtered := cm.Filter(func(key string, value int) bool {
    return value > 50
})
```

### 7. Pipeline (`pipeline.go`)

Chain processing stages with parallel execution support.

```go
// Create pipeline
pipeline := patterns.NewPipeline[string](logger)

// Add stages
pipeline.AddStage(patterns.NewFunctionStage("validate", func(ctx context.Context, input string) (string, error) {
    if input == "" {
        return "", errors.New("empty input")
    }
    return input, nil
}))

pipeline.AddStage(patterns.NewFunctionStage("transform", func(ctx context.Context, input string) (string, error) {
    return strings.ToUpper(input), nil
}))

// Execute pipeline
result, err := pipeline.Execute(ctx, "hello world")

// Parallel pipeline for multiple inputs
parallelPipeline := patterns.NewParallelPipeline(pipeline, 5, logger)
results, errors := parallelPipeline.ExecuteAll(ctx, []string{"input1", "input2", "input3"})

// Advanced stages
retryStage := patterns.NewRetryStage(unstableStage, 3, logger)
timedStage := patterns.NewTimedStage(slowStage, logger)
conditionalStage := patterns.NewConditionalStage("optional", func(input string) bool {
    return len(input) > 10
}, expensiveStage)
```

## Best Practices

1. **Context Propagation**: Always pass context to support cancellation and timeouts
2. **Resource Cleanup**: Use defer statements or helper methods for proper cleanup
3. **Error Handling**: Check for specific pattern errors (circuit breaker, rate limit)
4. **Metrics**: Enable metrics in production for monitoring and debugging
5. **Configuration**: Use appropriate buffer sizes and timeouts for your use case
6. **Testing**: Patterns include comprehensive error handling and edge cases

## Performance Characteristics

- **Semaphore**: O(1) acquire/release with channel-based implementation
- **Worker Pool**: Configurable parallelism with bounded queues
- **Circuit Breaker**: O(1) state checks with time-based recovery
- **Rate Limiter**: Token bucket O(1), Sliding window O(n) where n is request count
- **Concurrent Map**: O(1) operations with configurable sharding (default 32 shards)
- **Batch Processor**: Configurable batch sizes and flush intervals for optimal throughput
- **Pipeline**: Sequential stages with optional parallel execution for multiple inputs

## Thread Safety

All patterns are designed to be thread-safe and can be safely used from multiple goroutines
simultaneously. Internal state is protected using appropriate synchronization primitives (mutexes,
channels, atomic operations).
