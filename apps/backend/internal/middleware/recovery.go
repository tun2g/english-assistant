package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"app-backend/internal/errors"
	"app-backend/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

// Recovery creates a custom recovery middleware with structured error handling
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		requestID := GetRequestID(c)
		
		// Create structured error with oops
		oopsErr := oops.
			In("panic_recovery").
			Tags("panic", "recovery").
			Code("PANIC_RECOVERED").
			Trace(requestID).
			With("request_method", c.Request.Method).
			With("request_path", c.Request.URL.Path).
			With("request_id", requestID).
			With("panic_value", err).
			With("stack_trace", string(debug.Stack())).
			Hint("Check server logs for detailed stack trace").
			Errorf("panic recovered: %v", err)

		// Log the error with full context
		log.WithRequest(requestID).Error(
			"Panic recovered",
			zap.Any("error", oopsErr),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Any("panic_value", err),
			zap.String("stack", string(debug.Stack())),
		)

		// Convert to app error and respond
		appErr := errors.FromOopsError(oopsErr).WithTraceID(requestID)
		appErr.Status = http.StatusInternalServerError
		
		c.JSON(appErr.Status, appErr)
		c.Abort()
	})
}

// ErrorHandler is a middleware to handle errors and convert them to standardized responses
func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			requestID := GetRequestID(c)
			err := c.Errors.Last().Err

			var appErr *errors.AppError

			// Handle different error types
			switch e := err.(type) {
			case *errors.AppError:
				appErr = e.WithTraceID(requestID)
			case oops.OopsError:
				appErr = errors.FromOopsError(e).WithTraceID(requestID)
			default:
				// Handle validation errors
				if validationErr := errors.HandleValidationError(err); validationErr != nil {
					appErr = validationErr.WithTraceID(requestID)
				} else {
					appErr = errors.NewInternalServerError(err.Error()).WithTraceID(requestID)
				}
			}

			// Log the error
			logLevel := log.Error
			if appErr.Status < 500 {
				logLevel = log.Warn
			}

			logLevel(
				"Request failed",
				zap.String("request_id", requestID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", appErr.Status),
				zap.String("error_code", appErr.Code),
				zap.String("error_id", appErr.ID),
				zap.Any("error", err),
			)

			// Return error response
			c.JSON(appErr.Status, appErr)
			c.Abort()
		}
	}
}

// HandleError is a helper function to set errors in the context
func HandleError(c *gin.Context, err error) {
	c.Error(err)
}

// HandleOopsError is a helper function to create and set oops errors
func HandleOopsError(c *gin.Context, domain, code, message string, attrs ...interface{}) {
	requestID := GetRequestID(c)
	
	builder := oops.
		In(domain).
		Code(code).
		Trace(requestID).
		With("request_method", c.Request.Method).
		With("request_path", c.Request.URL.Path).
		With("request_id", requestID)

	// Add custom attributes
	for i := 0; i < len(attrs); i += 2 {
		if i+1 < len(attrs) {
			key := fmt.Sprintf("%v", attrs[i])
			value := attrs[i+1]
			builder = builder.With(key, value)
		}
	}

	err := builder.Errorf(message)
	c.Error(err)
}