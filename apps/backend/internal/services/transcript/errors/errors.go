package errors

import (
	"fmt"
	"net/http"

	"app-backend/internal/errors"
)

var (
	// Provider-specific errors
	ErrTranscriptNotFound      = errors.NewAppError("Transcript not found for this video", nil, http.StatusNotFound)
	ErrTranscriptDisabled      = errors.NewAppError("Transcripts are disabled for this video", nil, http.StatusForbidden)
	ErrInvalidVideoID          = errors.NewAppError("Invalid YouTube video ID", nil, http.StatusBadRequest)
	ErrProviderNotAvailable    = errors.NewAppError("Transcript provider is not available", nil, http.StatusServiceUnavailable)
	ErrAllProvidersFailed      = errors.NewAppError("All transcript providers failed", nil, http.StatusServiceUnavailable)
	ErrInvalidLanguage         = errors.NewAppError("Invalid or unsupported language code", nil, http.StatusBadRequest)
	ErrRateLimitExceeded       = errors.NewAppError("Rate limit exceeded for transcript provider", nil, http.StatusTooManyRequests)
	ErrAuthenticationFailed    = errors.NewAppError("Authentication failed with transcript provider", nil, http.StatusUnauthorized)
)

// NewProviderError creates a new provider-specific error
func NewProviderError(provider string, err error) *errors.AppError {
	return errors.NewAppError(
		fmt.Sprintf("Provider %s failed: %v", provider, err),
		err,
		http.StatusServiceUnavailable,
	)
}

// NewVideoIDExtractionError creates an error for video ID extraction failures
func NewVideoIDExtractionError(url string, err error) *errors.AppError {
	return errors.NewAppError(
		fmt.Sprintf("Failed to extract video ID from URL: %s", url),
		err,
		http.StatusBadRequest,
	)
}