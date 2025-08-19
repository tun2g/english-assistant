package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/samber/oops"
)

// AppError represents a standardized application error
type AppError struct {
	ID        string            `json:"id"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   string            `json:"details,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	TraceID   string            `json:"trace_id,omitempty"`
	Status    int               `json:"status"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(message string, err error, status int) *AppError {
	appErr := &AppError{
		ID:        uuid.New().String(),
		Message:   message,
		Timestamp: time.Now(),
		Status:    status,
	}

	if err != nil {
		appErr.Details = err.Error()
	}

	// Set appropriate error code based on status
	switch status {
	case http.StatusBadRequest:
		appErr.Code = ErrCodeValidation
	case http.StatusUnauthorized:
		appErr.Code = ErrCodeUnauthorized
	case http.StatusForbidden:
		appErr.Code = ErrCodeForbidden
	case http.StatusNotFound:
		appErr.Code = ErrCodeNotFound
	case http.StatusConflict:
		appErr.Code = ErrCodeConflict
	case http.StatusInternalServerError:
		appErr.Code = ErrCodeInternalServer
	default:
		appErr.Code = ErrCodeInternalServer
	}

	return appErr
}

// Common error codes
const (
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeConflict        = "CONFLICT"
	ErrCodeInternalServer  = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodeServiceUnavail  = "SERVICE_UNAVAILABLE"
)

// Error builder functions
func NewValidationError(details string, fields map[string]string) *AppError {
	return &AppError{
		ID:        uuid.New().String(),
		Code:      ErrCodeValidation,
		Message:   "Validation failed",
		Details:   details,
		Fields:    fields,
		Timestamp: time.Now(),
		Status:    http.StatusBadRequest,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		ID:        uuid.New().String(),
		Code:      ErrCodeNotFound,
		Message:   fmt.Sprintf("%s not found", resource),
		Timestamp: time.Now(),
		Status:    http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return &AppError{
		ID:        uuid.New().String(),
		Code:      ErrCodeUnauthorized,
		Message:   message,
		Timestamp: time.Now(),
		Status:    http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return &AppError{
		ID:        uuid.New().String(),
		Code:      ErrCodeForbidden,
		Message:   message,
		Timestamp: time.Now(),
		Status:    http.StatusForbidden,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		ID:        uuid.New().String(),
		Code:      ErrCodeConflict,
		Message:   message,
		Timestamp: time.Now(),
		Status:    http.StatusConflict,
	}
}

func NewInternalServerError(message string) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return &AppError{
		ID:        uuid.New().String(),
		Code:      ErrCodeInternalServer,
		Message:   message,
		Timestamp: time.Now(),
		Status:    http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		ID:        uuid.New().String(),
		Code:      ErrCodeBadRequest,
		Message:   message,
		Timestamp: time.Now(),
		Status:    http.StatusBadRequest,
	}
}

// Validation error helper
func HandleValidationError(err error) *AppError {
	var validationErrors validator.ValidationErrors
	if errors, ok := err.(validator.ValidationErrors); ok {
		validationErrors = errors
	} else {
		return NewBadRequestError("Invalid request format")
	}

	fields := make(map[string]string)
	var messages []string

	for _, err := range validationErrors {
		field := strings.ToLower(err.Field())
		tag := err.Tag()
		var message string

		switch tag {
		case "required":
			message = "This field is required"
		case "email":
			message = "Must be a valid email address"
		case "min":
			message = fmt.Sprintf("Must be at least %s characters", err.Param())
		case "max":
			message = fmt.Sprintf("Must be no more than %s characters", err.Param())
		case "len":
			message = fmt.Sprintf("Must be exactly %s characters", err.Param())
		case "oneof":
			message = fmt.Sprintf("Must be one of: %s", err.Param())
		default:
			message = fmt.Sprintf("Invalid value for %s", tag)
		}

		fields[field] = message
		messages = append(messages, fmt.Sprintf("%s: %s", field, message))
	}

	return NewValidationError(strings.Join(messages, "; "), fields)
}

// Oops error builder helper
func WithOops(domain string) oops.OopsErrorBuilder {
	return oops.
		In(domain).
		Trace(uuid.New().String()).
		Time(time.Now())
}

// Convert oops error to AppError
func FromOopsError(err error) *AppError {
	if oopsErr, ok := err.(oops.OopsError); ok {
		appErr := &AppError{
			ID:        uuid.New().String(),
			Code:      oopsErr.Code(),
			Message:   oopsErr.Error(),
			Timestamp: oopsErr.Time(),
			TraceID:   oopsErr.Trace(),
			Status:    http.StatusInternalServerError,
		}

		// Map domain to HTTP status
		switch oopsErr.Domain() {
		case "validation":
			appErr.Status = http.StatusBadRequest
			appErr.Code = ErrCodeValidation
		case "auth", "authentication":
			appErr.Status = http.StatusUnauthorized
			appErr.Code = ErrCodeUnauthorized
		case "authz", "authorization":
			appErr.Status = http.StatusForbidden
			appErr.Code = ErrCodeForbidden
		case "not_found":
			appErr.Status = http.StatusNotFound
			appErr.Code = ErrCodeNotFound
		}

		// Use the error message from oops
		if oopsErr.Error() != "" {
			appErr.Message = oopsErr.Error()
		}

		return appErr
	}

	// Fallback for regular errors
	return NewInternalServerError(err.Error())
}

// JSON serialization
func (e *AppError) JSON() []byte {
	data, _ := json.Marshal(e)
	return data
}

// Response helpers
func (e *AppError) WithTraceID(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}