package patterns

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	ErrCircuitBreakerOpen     = errors.New("circuit breaker is open")
	ErrTooManyRequests        = errors.New("too many requests")
	ErrCircuitBreakerTimeout  = errors.New("circuit breaker timeout")
)

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateHalfOpen
	StateOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF_OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	Name                   string        // Name for logging and metrics
	MaxRequests            uint32        // Max requests allowed when half-open
	Interval               time.Duration // Time window for failure counting
	Timeout                time.Duration // Time to wait before transitioning from open to half-open
	FailureThreshold       uint32        // Number of failures to trip the breaker
	SuccessThreshold       uint32        // Number of successes needed to close from half-open
	IsFailure              func(error) bool // Function to determine if error should count as failure
	OnStateChange          func(name string, from, to CircuitBreakerState) // Callback for state changes
	Logger                 *zap.Logger   // Logger instance
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config     CircuitBreakerConfig
	mutex      sync.RWMutex
	state      CircuitBreakerState
	counts     *Counts
	expiry     time.Time
	generation uint64
}

// Counts holds the statistics for the circuit breaker
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	// Set defaults
	if config.MaxRequests == 0 {
		config.MaxRequests = 1
	}
	if config.Interval <= 0 {
		config.Interval = 60 * time.Second
	}
	if config.Timeout <= 0 {
		config.Timeout = 60 * time.Second
	}
	if config.FailureThreshold == 0 {
		config.FailureThreshold = 5
	}
	if config.SuccessThreshold == 0 {
		config.SuccessThreshold = 1
	}
	if config.IsFailure == nil {
		config.IsFailure = func(err error) bool { return err != nil }
	}
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}
	if config.Name == "" {
		config.Name = "circuit-breaker"
	}

	cb := &CircuitBreaker{
		config:     config,
		state:      StateClosed,
		counts:     &Counts{},
		expiry:     time.Now().Add(config.Interval),
		generation: 0,
	}

	return cb
}

// Execute executes the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	generation, err := cb.beforeRequest()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			cb.afterRequest(generation, false)
			panic(r)
		}
	}()

	// Execute with timeout if context doesn't have one
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cb.config.Timeout)
		defer cancel()
	}

	// Execute the function
	err = fn()
	cb.afterRequest(generation, !cb.config.IsFailure(err))
	return err
}

// ExecuteWithFallback executes the function with a fallback if circuit breaker is open
func (cb *CircuitBreaker) ExecuteWithFallback(ctx context.Context, fn func() error, fallback func() error) error {
	err := cb.Execute(ctx, fn)
	if errors.Is(err, ErrCircuitBreakerOpen) || errors.Is(err, ErrTooManyRequests) {
		if fallback != nil {
			cb.config.Logger.Debug("Circuit breaker executing fallback",
				zap.String("name", cb.config.Name),
				zap.String("reason", err.Error()))
			return fallback()
		}
	}
	return err
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	now := time.Now()
	state, _ := cb.currentState(now)
	return state
}

// GetCounts returns a copy of the current counts
func (cb *CircuitBreaker) GetCounts() Counts {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return *cb.counts
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.changeState(StateClosed, time.Now())
	cb.config.Logger.Info("Circuit breaker reset", zap.String("name", cb.config.Name))
}

// beforeRequest checks if the request should be allowed
func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, ErrCircuitBreakerOpen
	}

	if state == StateHalfOpen && cb.counts.Requests >= cb.config.MaxRequests {
		return generation, ErrTooManyRequests
	}

	cb.counts.Requests++
	return generation, nil
}

// afterRequest updates the circuit breaker state after a request
func (cb *CircuitBreaker) afterRequest(before uint64, success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	
	// Ignore results from different generations
	if generation != before {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

// onSuccess handles a successful request
func (cb *CircuitBreaker) onSuccess(state CircuitBreakerState, now time.Time) {
	cb.counts.TotalSuccesses++
	cb.counts.ConsecutiveSuccesses++
	cb.counts.ConsecutiveFailures = 0

	if state == StateHalfOpen && cb.counts.ConsecutiveSuccesses >= cb.config.SuccessThreshold {
		cb.changeState(StateClosed, now)
	}
}

// onFailure handles a failed request
func (cb *CircuitBreaker) onFailure(state CircuitBreakerState, now time.Time) {
	cb.counts.TotalFailures++
	cb.counts.ConsecutiveFailures++
	cb.counts.ConsecutiveSuccesses = 0

	if state == StateClosed {
		if cb.counts.ConsecutiveFailures >= cb.config.FailureThreshold {
			cb.changeState(StateOpen, now)
		}
	} else if state == StateHalfOpen {
		cb.changeState(StateOpen, now)
	}
}

// currentState returns the current state and generation
func (cb *CircuitBreaker) currentState(now time.Time) (CircuitBreakerState, uint64) {
	switch cb.state {
	case StateClosed:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.changeState(StateClosed, now)
		}
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.changeState(StateHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

// changeState changes the state of the circuit breaker
func (cb *CircuitBreaker) changeState(state CircuitBreakerState, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state
	cb.generation++

	cb.counts = &Counts{}

	var expiry time.Time
	switch state {
	case StateClosed:
		expiry = now.Add(cb.config.Interval)
	case StateOpen:
		expiry = now.Add(cb.config.Timeout)
	default: // StateHalfOpen
		expiry = time.Time{} // No expiry for half-open
	}
	cb.expiry = expiry

	// Call state change callback
	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.config.Name, prev, state)
	}

	cb.config.Logger.Info("Circuit breaker state changed",
		zap.String("name", cb.config.Name),
		zap.String("from", prev.String()),
		zap.String("to", state.String()),
		zap.Time("expiry", expiry))
}

// IsCircuitBreakerError checks if an error is a circuit breaker error
func IsCircuitBreakerError(err error) bool {
	return errors.Is(err, ErrCircuitBreakerOpen) || 
		   errors.Is(err, ErrTooManyRequests) || 
		   errors.Is(err, ErrCircuitBreakerTimeout)
}

// CircuitBreakerMetrics provides metrics for monitoring
type CircuitBreakerMetrics struct {
	Name                   string
	State                  string
	TotalRequests          uint32
	TotalSuccesses         uint32
	TotalFailures          uint32
	ConsecutiveSuccesses   uint32
	ConsecutiveFailures    uint32
	FailureRate            float64
}

// GetMetrics returns current metrics for the circuit breaker
func (cb *CircuitBreaker) GetMetrics() CircuitBreakerMetrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	counts := *cb.counts
	var failureRate float64
	
	totalRequests := counts.TotalSuccesses + counts.TotalFailures
	if totalRequests > 0 {
		failureRate = float64(counts.TotalFailures) / float64(totalRequests)
	}

	return CircuitBreakerMetrics{
		Name:                 cb.config.Name,
		State:                cb.state.String(),
		TotalRequests:        counts.Requests,
		TotalSuccesses:       counts.TotalSuccesses,
		TotalFailures:        counts.TotalFailures,
		ConsecutiveSuccesses: counts.ConsecutiveSuccesses,
		ConsecutiveFailures:  counts.ConsecutiveFailures,
		FailureRate:         failureRate,
	}
}