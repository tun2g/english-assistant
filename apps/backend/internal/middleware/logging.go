package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sloggin "github.com/samber/slog-gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// RequestID adds a unique request ID to each request and logs incoming/outgoing requests
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context and header
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)
		
		// Log incoming request with colorization
		methodColor := getMethodColor(c.Request.Method)
		
		fmt.Printf("%s [%s] %s %s %s - Request ID: %s\n",
			color.BlueString("====== INCOMING REQUEST"),
			time.Now().Format("2006-01-02 15:04:05"),
			methodColor.Sprint(c.Request.Method),
			color.YellowString(c.Request.URL.Path),
			color.MagentaString(c.ClientIP()),
			color.GreenString(requestID))
		
		c.Next()
		
		// Log outgoing response
		duration := time.Since(start)
		statusColor := getStatusColor(c.Writer.Status())
		
		fmt.Printf("%s [%s] %s %s %s %s %s - Request ID: %s\n",
			color.BlueString("====== OUTGOING REQUEST"),
			time.Now().Format("2006-01-02 15:04:05"),
			methodColor.Sprint(c.Request.Method),
			color.YellowString(c.Request.URL.Path),
			statusColor.Sprint(c.Writer.Status()),
			color.MagentaString(c.ClientIP()),
			color.CyanString(duration.String()),
			color.GreenString(requestID))
	}
}

// getMethodColor returns appropriate color for HTTP methods
func getMethodColor(method string) *color.Color {
	switch method {
	case "GET":
		return color.New(color.FgBlue)
	case "POST":
		return color.New(color.FgGreen)
	case "PUT":
		return color.New(color.FgYellow)
	case "DELETE":
		return color.New(color.FgRed)
	case "PATCH":
		return color.New(color.FgMagenta)
	case "HEAD":
		return color.New(color.FgCyan)
	case "OPTIONS":
		return color.New(color.FgWhite)
	default:
		return color.New(color.FgHiWhite)
	}
}

// getStatusColor returns appropriate color for HTTP status codes
func getStatusColor(status int) *color.Color {
	switch {
	case status >= 200 && status < 300:
		return color.New(color.FgGreen)
	case status >= 300 && status < 400:
		return color.New(color.FgYellow)
	case status >= 400 && status < 500:
		return color.New(color.FgRed)
	case status >= 500:
		return color.New(color.FgHiRed)
	default:
		return color.New(color.FgWhite)
	}
}

// LoggingMiddleware creates a structured logging middleware using slog
func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	config := sloggin.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithRequestID:      true,
		WithUserAgent:      true,
		WithRequestBody:    false, // Enable in development if needed
		WithRequestHeader:  false, // Enable in development if needed
		WithResponseBody:   false, // Enable in development if needed
		WithResponseHeader: false, // Enable in development if needed
		WithSpanID:         true,
		WithTraceID:        true,

		Filters: []sloggin.Filter{
			// Ignore health check endpoints
			sloggin.IgnorePath("/health"),
			sloggin.IgnorePath("/metrics"),
			// Ignore static assets
			sloggin.IgnorePathPrefix("/static/"),
			sloggin.IgnorePathPrefix("/assets/"),
		},
	}

	// Add custom attributes
	enrichedLogger := logger.With(
		slog.String("service", "app-backend"),
		slog.String("version", "1.0.0"),
		slog.Time("server_start_time", time.Now()),
	)

	return sloggin.NewWithConfig(enrichedLogger, config)
}

// Custom attributes helper for adding request-specific data
func AddLogAttributes(c *gin.Context, attrs ...slog.Attr) {
	for _, attr := range attrs {
		sloggin.AddCustomAttributes(c, attr)
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}