package oauth

import (
	"context"
	"golang.org/x/oauth2"
)

// ServiceInterface defines the interface for OAuth operations
type ServiceInterface interface {
	// GenerateAuthURL creates an authorization URL for the user to visit
	GenerateAuthURL(state string) string
	
	// ExchangeCodeForTokens exchanges authorization code for access and refresh tokens
	ExchangeCodeForTokens(ctx context.Context, code string) (*oauth2.Token, error)
	
	// GetValidToken returns a valid access token, refreshing if necessary
	GetValidToken(ctx context.Context) (*oauth2.Token, error)
	
	// RefreshToken refreshes an expired access token using refresh token
	RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error)
	
	// SaveToken saves token to persistent storage
	SaveToken(token *oauth2.Token) error
	
	// LoadToken loads token from persistent storage
	LoadToken() (*oauth2.Token, error)
	
	// IsAuthenticated checks if user is currently authenticated
	IsAuthenticated() bool
	
	// RevokeToken revokes the current token
	RevokeToken(ctx context.Context) error
	
	// StoreState stores an OAuth state parameter for CSRF protection
	StoreState(state string)
	
	// ValidateAndClearState validates and removes an OAuth state parameter
	ValidateAndClearState(state string) bool
}