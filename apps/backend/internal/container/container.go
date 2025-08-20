package container

import (
	"app-backend/internal/config"
	"app-backend/internal/handlers/auth"
	"app-backend/internal/handlers/oauth"
	"app-backend/internal/handlers/translation"
	"app-backend/internal/handlers/user"
	"app-backend/internal/handlers/video"
	"app-backend/internal/logger"
	"app-backend/internal/middleware"
	"app-backend/internal/repositories"
	authService "app-backend/internal/services/auth"
	jwtService "app-backend/internal/services/jwt"
	oauthService "app-backend/internal/services/oauth"
	transcriptService "app-backend/internal/services/transcript"
	translationService "app-backend/internal/services/translation"
	userService "app-backend/internal/services/user"
	videoService "app-backend/internal/services/video"
	"app-backend/pkg/gemini"
	"app-backend/pkg/youtube"

	"gorm.io/gorm"
	"go.uber.org/zap"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	Config *config.Config

	// Database
	DB *gorm.DB

	// Logger
	Logger *logger.Logger

	// Repositories
	UserRepository    repositories.UserRepositoryInterface
	SessionRepository repositories.SessionRepositoryInterface

	// Services
	JWTService      jwtService.ServiceInterface
	UserService     userService.ServiceInterface
	AuthService     authService.ServiceInterface
	VideoService    videoService.ServiceInterface
	YouTubeOAuthService oauthService.ServiceInterface
	TranscriptService   transcriptService.ServiceInterface
	TranslationService  translationService.ServiceInterface

	// External Services
	YouTubeService *youtube.Service
	GeminiService  *gemini.Service

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// Handlers
	AuthHandler       auth.HandlerInterface
	UserHandler       user.HandlerInterface
	VideoHandler      video.HandlerInterface
	OAuthHandler      oauth.HandlerInterface
	TranslationHandler translation.HandlerInterface
}

// NewContainer creates and initializes all dependencies
func NewContainer(cfg *config.Config, db *gorm.DB, logger *logger.Logger) *Container {
	container := &Container{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}

	container.initRepositories()
	container.initExternalServices()
	container.initServices()
	container.initMiddleware()
	container.initHandlers()

	return container
}

// initRepositories initializes all repositories
func (c *Container) initRepositories() {
	c.UserRepository = repositories.NewUserRepository(c.DB)
	c.SessionRepository = repositories.NewSessionRepository(c.DB)
}

// initExternalServices initializes external API services
func (c *Container) initExternalServices() {
	youtubeKey := c.Config.ExternalAPIs.YouTube.APIKey
	geminiKey := c.Config.ExternalAPIs.Gemini.APIKey
	
	youtubePrefix := "empty"
	if len(youtubeKey) > 10 {
		youtubePrefix = youtubeKey[:10] + "..."
	} else if len(youtubeKey) > 0 {
		youtubePrefix = youtubeKey + "..."
	}
	
	geminiPrefix := "empty"
	if len(geminiKey) > 10 {
		geminiPrefix = geminiKey[:10] + "..."
	} else if len(geminiKey) > 0 {
		geminiPrefix = geminiKey + "..."
	}
	
	c.Logger.Zap().Info("Initializing external services", 
		zap.String("youtube_api_key_prefix", youtubePrefix),
		zap.String("gemini_api_key_prefix", geminiPrefix))
	c.GeminiService = gemini.NewService(geminiKey, c.Logger.Zap())
}

// initServices initializes all services
func (c *Container) initServices() {
	c.JWTService = jwtService.NewJWTService(c.Config)
	c.UserService = userService.NewUserService(c.UserRepository)
	c.AuthService = authService.NewAuthService(c.UserService, c.JWTService, c.SessionRepository)
	c.YouTubeOAuthService = oauthService.NewYouTubeOAuthService(c.Config, c.Logger)
	
	// Initialize YouTube service with OAuth support
	youtubeKey := c.Config.ExternalAPIs.YouTube.APIKey
	c.YouTubeService = youtube.NewServiceWithOAuth(youtubeKey, c.YouTubeOAuthService, c.Logger.Zap())
	
	// Initialize transcript service
	transcriptSvc, err := transcriptService.NewService(c.Config, c.Logger)
	if err != nil {
		c.Logger.Error("Failed to initialize transcript service", zap.Error(err))
	} else {
		c.TranscriptService = transcriptSvc
	}
	
	// Initialize translation service
	translationSvc, err := translationService.NewService(&translationService.Config{
		GeminiAPIKey: c.Config.ExternalAPIs.Gemini.APIKey,
		Logger:       c.Logger,
	})
	if err != nil {
		c.Logger.Error("Failed to initialize translation service", zap.Error(err))
	} else {
		c.TranslationService = translationSvc
	}
	
	c.VideoService = videoService.NewVideoService(c.YouTubeService, c.GeminiService, c.Logger.Zap())
}

// initMiddleware initializes all middleware
func (c *Container) initMiddleware() {
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.JWTService, c.AuthService, c.Logger)
}

// initHandlers initializes all handlers
func (c *Container) initHandlers() {
	c.AuthHandler = auth.NewAuthHandler(c.AuthService, c.Logger)
	c.UserHandler = user.NewUserHandler(c.UserService, c.Logger)
	c.VideoHandler = video.NewVideoHandler(c.VideoService, c.TranscriptService, c.Logger)
	c.OAuthHandler = oauth.NewOAuthHandler(c.YouTubeOAuthService, c.Logger)
	c.TranslationHandler = translation.NewTranslationHandler(c.TranslationService, c.Logger)
}