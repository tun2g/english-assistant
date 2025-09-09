package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Stage represents a single stage in a pipeline
type Stage[T any] interface {
	Process(ctx context.Context, input T) (T, error)
	Name() string
}

// Pipeline represents a chain of processing stages
type Pipeline[T any] struct {
	stages []Stage[T]
	logger *zap.Logger
}

// NewPipeline creates a new pipeline
func NewPipeline[T any](logger *zap.Logger) *Pipeline[T] {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Pipeline[T]{
		stages: make([]Stage[T], 0),
		logger: logger,
	}
}

// AddStage adds a stage to the pipeline
func (p *Pipeline[T]) AddStage(stage Stage[T]) *Pipeline[T] {
	p.stages = append(p.stages, stage)
	return p
}

// Execute executes the pipeline with the given input
func (p *Pipeline[T]) Execute(ctx context.Context, input T) (T, error) {
	current := input
	
	for i, stage := range p.stages {
		p.logger.Debug("Executing pipeline stage",
			zap.Int("stage_index", i),
			zap.String("stage_name", stage.Name()))
		
		result, err := stage.Process(ctx, current)
		if err != nil {
			p.logger.Error("Pipeline stage failed",
				zap.Int("stage_index", i),
				zap.String("stage_name", stage.Name()),
				zap.Error(err))
			return current, fmt.Errorf("stage %d (%s) failed: %w", i, stage.Name(), err)
		}
		
		current = result
		
		// Check for context cancellation between stages
		select {
		case <-ctx.Done():
			return current, ctx.Err()
		default:
		}
	}
	
	p.logger.Debug("Pipeline execution completed", zap.Int("stages", len(p.stages)))
	return current, nil
}

// ParallelPipeline executes multiple items through a pipeline concurrently
type ParallelPipeline[T any] struct {
	pipeline *Pipeline[T]
	semaphore *Semaphore
	logger   *zap.Logger
}

// NewParallelPipeline creates a new parallel pipeline
func NewParallelPipeline[T any](pipeline *Pipeline[T], maxConcurrency int, logger *zap.Logger) *ParallelPipeline[T] {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &ParallelPipeline[T]{
		pipeline:  pipeline,
		semaphore: NewSemaphore(maxConcurrency),
		logger:    logger,
	}
}

// ExecuteAll executes the pipeline for all inputs concurrently
func (pp *ParallelPipeline[T]) ExecuteAll(ctx context.Context, inputs []T) ([]T, []error) {
	if len(inputs) == 0 {
		return []T{}, []error{}
	}
	
	results := make([]T, len(inputs))
	errors := make([]error, len(inputs))
	var wg sync.WaitGroup
	
	for i, input := range inputs {
		wg.Add(1)
		go func(index int, item T) {
			defer wg.Done()
			
			err := pp.semaphore.Acquire(ctx)
			if err != nil {
				errors[index] = err
				return
			}
			defer pp.semaphore.Release()
			
			result, err := pp.pipeline.Execute(ctx, item)
			results[index] = result
			errors[index] = err
		}(i, input)
	}
	
	wg.Wait()
	return results, errors
}

// FunctionStage wraps a function as a pipeline stage
type FunctionStage[T any] struct {
	name string
	fn   func(context.Context, T) (T, error)
}

// NewFunctionStage creates a new function stage
func NewFunctionStage[T any](name string, fn func(context.Context, T) (T, error)) *FunctionStage[T] {
	return &FunctionStage[T]{
		name: name,
		fn:   fn,
	}
}

// Process implements Stage interface
func (fs *FunctionStage[T]) Process(ctx context.Context, input T) (T, error) {
	return fs.fn(ctx, input)
}

// Name implements Stage interface
func (fs *FunctionStage[T]) Name() string {
	return fs.name
}

// ConditionalStage executes a stage only if a condition is met
type ConditionalStage[T any] struct {
	name      string
	condition func(T) bool
	stage     Stage[T]
}

// NewConditionalStage creates a new conditional stage
func NewConditionalStage[T any](name string, condition func(T) bool, stage Stage[T]) *ConditionalStage[T] {
	return &ConditionalStage[T]{
		name:      name,
		condition: condition,
		stage:     stage,
	}
}

// Process implements Stage interface
func (cs *ConditionalStage[T]) Process(ctx context.Context, input T) (T, error) {
	if !cs.condition(input) {
		return input, nil // Skip processing
	}
	return cs.stage.Process(ctx, input)
}

// Name implements Stage interface
func (cs *ConditionalStage[T]) Name() string {
	return cs.name
}

// RetryStage wraps a stage with retry logic
type RetryStage[T any] struct {
	stage      Stage[T]
	maxRetries int
	logger     *zap.Logger
}

// NewRetryStage creates a new retry stage
func NewRetryStage[T any](stage Stage[T], maxRetries int, logger *zap.Logger) *RetryStage[T] {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &RetryStage[T]{
		stage:      stage,
		maxRetries: maxRetries,
		logger:     logger,
	}
}

// Process implements Stage interface
func (rs *RetryStage[T]) Process(ctx context.Context, input T) (T, error) {
	var lastErr error
	
	for attempt := 0; attempt <= rs.maxRetries; attempt++ {
		result, err := rs.stage.Process(ctx, input)
		if err == nil {
			if attempt > 0 {
				rs.logger.Info("Stage succeeded after retry",
					zap.String("stage", rs.stage.Name()),
					zap.Int("attempt", attempt))
			}
			return result, nil
		}
		
		lastErr = err
		
		if attempt < rs.maxRetries {
			rs.logger.Warn("Stage failed, retrying",
				zap.String("stage", rs.stage.Name()),
				zap.Int("attempt", attempt),
				zap.Error(err))
		}
		
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return input, ctx.Err()
		default:
		}
	}
	
	rs.logger.Error("Stage failed after all retries",
		zap.String("stage", rs.stage.Name()),
		zap.Int("max_retries", rs.maxRetries),
		zap.Error(lastErr))
	
	return input, fmt.Errorf("stage %s failed after %d retries: %w", rs.stage.Name(), rs.maxRetries, lastErr)
}

// Name implements Stage interface
func (rs *RetryStage[T]) Name() string {
	return fmt.Sprintf("retry-%s", rs.stage.Name())
}

// TimedStage wraps a stage with timing information
type TimedStage[T any] struct {
	stage  Stage[T]
	logger *zap.Logger
}

// NewTimedStage creates a new timed stage
func NewTimedStage[T any](stage Stage[T], logger *zap.Logger) *TimedStage[T] {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &TimedStage[T]{
		stage:  stage,
		logger: logger,
	}
}

// Process implements Stage interface
func (ts *TimedStage[T]) Process(ctx context.Context, input T) (T, error) {
	start := time.Now()
	result, err := ts.stage.Process(ctx, input)
	duration := time.Since(start)
	
	if err != nil {
		ts.logger.Error("Timed stage failed",
			zap.String("stage", ts.stage.Name()),
			zap.Duration("duration", duration),
			zap.Error(err))
	} else {
		ts.logger.Debug("Timed stage completed",
			zap.String("stage", ts.stage.Name()),
			zap.Duration("duration", duration))
	}
	
	return result, err
}

// Name implements Stage interface
func (ts *TimedStage[T]) Name() string {
	return fmt.Sprintf("timed-%s", ts.stage.Name())
}

