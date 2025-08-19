package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App          AppConfig           `mapstructure:"app"`
	Database     DatabaseConfig     `mapstructure:"database"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	CORS         CORSConfig         `mapstructure:"cors"`
	Security     SecurityConfig     `mapstructure:"security"`
	ExternalAPIs ExternalAPIsConfig `mapstructure:"external_apis"`
	Transcript   TranscriptConfig   `mapstructure:"transcript"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`
	LogLevel    string `mapstructure:"log_level"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	AccessTTLMinutes   int    `mapstructure:"access_ttl_minutes"`
	RefreshTTLHours    int    `mapstructure:"refresh_ttl_hours"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	ExposeHeaders  []string `mapstructure:"expose_headers"`
	MaxAge         int      `mapstructure:"max_age"`
}

type SecurityConfig struct {
	BcryptCost int           `mapstructure:"bcrypt_cost"`
	RateLimit  RateLimitConfig `mapstructure:"rate_limit"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	Burst             int `mapstructure:"burst"`
}

type ExternalAPIsConfig struct {
	YouTube YouTubeConfig `mapstructure:"youtube"`
	Gemini  GeminiConfig  `mapstructure:"gemini"`
}

type YouTubeConfig struct {
	APIKey       string      `mapstructure:"api_key"`
	APIURL       string      `mapstructure:"api_url"`
	RateLimit    int         `mapstructure:"rate_limit"`
	OAuth        OAuthConfig `mapstructure:"oauth"`
}

type OAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	TokenStorage string `mapstructure:"token_storage"`
}

type GeminiConfig struct {
	APIKey    string `mapstructure:"api_key"`
	APIURL    string `mapstructure:"api_url"`
	RateLimit int    `mapstructure:"rate_limit"`
}

type TranscriptConfig struct {
	Providers []TranscriptProviderConfig `mapstructure:"providers"`
}

type TranscriptProviderConfig struct {
	Type     string                 `mapstructure:"type"`
	Enabled  bool                   `mapstructure:"enabled"`
	Priority int                    `mapstructure:"priority"`
	Config   map[string]interface{} `mapstructure:"config"`
}

// GetDatabaseURL returns the formatted database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// New creates and initializes a new configuration
func New() (*Config, error) {
	// Set configuration file settings
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	// Set environment variable settings
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("APP")

	// Set default values
	setDefaults()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is okay, we'll use env vars and defaults
	}

	// Check for environment-specific config file
	env := viper.GetString("app.environment")
	if env != "" {
		viper.SetConfigName(fmt.Sprintf("app.%s", env))
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read environment config file: %w", err)
			}
		}
	}

	// Unmarshal configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "app-backend")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.log_level", "info")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "app_backend_dev")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "1h")

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-in-production")
	viper.SetDefault("jwt.access_ttl_minutes", 15)
	viper.SetDefault("jwt.refresh_ttl_hours", 168)

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"})
	viper.SetDefault("cors.expose_headers", []string{"Content-Length", "Content-Type"})
	viper.SetDefault("cors.max_age", 300)

	// Security defaults
	viper.SetDefault("security.bcrypt_cost", 12)
	viper.SetDefault("security.rate_limit.requests_per_minute", 60)
	viper.SetDefault("security.rate_limit.burst", 10)

	// External APIs defaults
	viper.SetDefault("external_apis.youtube.api_key", "")
	viper.SetDefault("external_apis.youtube.api_url", "https://www.googleapis.com/youtube/v3")
	viper.SetDefault("external_apis.youtube.rate_limit", 100)
	
	// YouTube OAuth defaults
	viper.SetDefault("external_apis.youtube.oauth.client_id", "")
	viper.SetDefault("external_apis.youtube.oauth.client_secret", "")
	viper.SetDefault("external_apis.youtube.oauth.redirect_url", "http://localhost:8000/api/v1/oauth/youtube/callback")
	viper.SetDefault("external_apis.youtube.oauth.token_storage", "./.oauth_tokens")
	
	viper.SetDefault("external_apis.gemini.api_key", "")
	viper.SetDefault("external_apis.gemini.api_url", "https://generativelanguage.googleapis.com")
	viper.SetDefault("external_apis.gemini.rate_limit", 60)
	
	// Transcript service defaults
	viper.SetDefault("transcript.providers", []map[string]interface{}{
		{
			"type":     "youtube_api",
			"enabled":  false,
			"priority": 1,
			"config": map[string]interface{}{
				"api_key": "",
			},
		},
		{
			"type":     "yt_transcript",
			"enabled":  true,
			"priority": 2,
			"config":   map[string]interface{}{},
		},
		{
			"type":     "kkdai_youtube",
			"enabled":  true,
			"priority": 3,
			"config":   map[string]interface{}{},
		},
		{
			"type":     "innertube",
			"enabled":  true,
			"priority": 4,
			"config": map[string]interface{}{
				"timeout": 30,
			},
		},
	})
}