package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"app-backend/internal/config"
	"app-backend/internal/logger"
	
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Service implements OAuth operations for YouTube API
type Service struct {
	config      *oauth2.Config
	tokenPath   string
	logger      *logger.Logger
	stateStore  map[string]time.Time // In-memory state storage with expiration
	stateMutex  sync.RWMutex         // Mutex for thread-safe state operations
}

// NewYouTubeOAuthService creates a new OAuth service for YouTube API
func NewYouTubeOAuthService(cfg *config.Config, logger *logger.Logger) ServiceInterface {
	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ExternalAPIs.YouTube.OAuth.ClientID,
		ClientSecret: cfg.ExternalAPIs.YouTube.OAuth.ClientSecret,
		RedirectURL:  cfg.ExternalAPIs.YouTube.OAuth.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/youtube.force-ssl"},
		Endpoint:     google.Endpoint,
	}

	return &Service{
		config:     oauth2Config,
		tokenPath:  cfg.ExternalAPIs.YouTube.OAuth.TokenStorage,
		logger:     logger,
		stateStore: make(map[string]time.Time),
	}
}

// GenerateAuthURL creates an authorization URL for the user to visit
func (s *Service) GenerateAuthURL(state string) string {
	if state == "" {
		state = s.generateRandomState()
	}
	
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
}

// ExchangeCodeForTokens exchanges authorization code for access and refresh tokens
func (s *Service) ExchangeCodeForTokens(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", zap.Error(err))
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Save the token for future use
	if err := s.SaveToken(token); err != nil {
		s.logger.Warn("Failed to save token", zap.Error(err))
		// Don't return error here as the token exchange was successful
	}

	s.logger.Info("Successfully exchanged code for token")
	return token, nil
}

// GetValidToken returns a valid access token, refreshing if necessary
func (s *Service) GetValidToken(ctx context.Context) (*oauth2.Token, error) {
	token, err := s.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("no saved token found: %w", err)
	}

	// Check if token needs refresh
	if token.Expiry.Before(time.Now().Add(5 * time.Minute)) {
		s.logger.Info("Token is expired or will expire soon, refreshing...")
		refreshedToken, err := s.RefreshToken(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		return refreshedToken, nil
	}

	return token, nil
}

// RefreshToken refreshes an expired access token using refresh token
func (s *Service) RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	tokenSource := s.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		s.logger.Error("Failed to refresh token", zap.Error(err))
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Save the refreshed token
	if err := s.SaveToken(newToken); err != nil {
		s.logger.Warn("Failed to save refreshed token", zap.Error(err))
	}

	s.logger.Info("Successfully refreshed token")
	return newToken, nil
}

// SaveToken saves token to persistent storage
func (s *Service) SaveToken(token *oauth2.Token) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(s.tokenPath), 0700); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	// Marshal token to JSON
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Write token to file with restricted permissions
	if err := os.WriteFile(s.tokenPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	s.logger.Debug("Token saved successfully", zap.String("path", s.tokenPath))
	return nil
}

// LoadToken loads token from persistent storage
func (s *Service) LoadToken() (*oauth2.Token, error) {
	data, err := os.ReadFile(s.tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// IsAuthenticated checks if user is currently authenticated
func (s *Service) IsAuthenticated() bool {
	token, err := s.LoadToken()
	if err != nil {
		return false
	}

	// Check if token exists and is not expired (with 5 minute buffer)
	return token != nil && token.Valid() && token.Expiry.After(time.Now().Add(5*time.Minute))
}

// RevokeToken revokes the current token
func (s *Service) RevokeToken(ctx context.Context) error {
	token, err := s.LoadToken()
	if err != nil {
		return fmt.Errorf("no token to revoke: %w", err)
	}

	// Google OAuth2 revoke endpoint
	revokeURL := fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", token.AccessToken)
	
	// Make HTTP request to revoke token
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(revokeURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to revoke token, status: %d", resp.StatusCode)
	}

	// Remove token file
	if err := os.Remove(s.tokenPath); err != nil && !os.IsNotExist(err) {
		s.logger.Warn("Failed to remove token file", zap.Error(err))
	}

	s.logger.Info("Successfully revoked token")
	return nil
}

// generateRandomState generates a random state string for OAuth flow
func (s *Service) generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// StoreState stores an OAuth state parameter with expiration (10 minutes)
func (s *Service) StoreState(state string) {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()
	
	// Clean up expired states while we have the lock
	s.cleanupExpiredStates()
	
	// Store new state with expiration time
	s.stateStore[state] = time.Now().Add(10 * time.Minute)
	
	s.logger.Debug("Stored OAuth state", zap.String("state", state))
}

// ValidateAndClearState validates a state parameter and removes it from storage
func (s *Service) ValidateAndClearState(state string) bool {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()
	
	expiry, exists := s.stateStore[state]
	if !exists {
		s.logger.Warn("OAuth state not found", zap.String("state", state))
		return false
	}
	
	// Remove the state (use once)
	delete(s.stateStore, state)
	
	// Check if expired
	if time.Now().After(expiry) {
		s.logger.Warn("OAuth state expired", zap.String("state", state))
		return false
	}
	
	s.logger.Debug("OAuth state validated successfully", zap.String("state", state))
	return true
}

// cleanupExpiredStates removes expired states from storage (called with lock held)
func (s *Service) cleanupExpiredStates() {
	now := time.Now()
	for state, expiry := range s.stateStore {
		if now.After(expiry) {
			delete(s.stateStore, state)
			s.logger.Debug("Cleaned up expired OAuth state", zap.String("state", state))
		}
	}
}