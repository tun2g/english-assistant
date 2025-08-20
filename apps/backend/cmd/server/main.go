package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app-backend/internal/config"
	"app-backend/internal/container"
	"app-backend/internal/database"
	"app-backend/internal/logger"
	"app-backend/internal/middleware"
	"app-backend/internal/routes"
	_ "app-backend/docs" // Import generated swagger docs
	_ "app-backend/internal/docs" // Import docs for swagger generation

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	appLogger, err := logger.New(cfg.App.Environment)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting application",
		zap.String("environment", cfg.App.Environment),
		zap.String("port", cfg.App.Port),
		zap.String("log_level", cfg.App.LogLevel),
	)

	// Initialize database
	db, err := database.NewConnection(cfg.GetDatabaseURL())
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	appLogger.Info("Database connected successfully")

	// Auto-migrate database schemas
	if err := database.AutoMigrate(db); err != nil {
		appLogger.Fatal("Failed to migrate database", zap.Error(err))
	}
	appLogger.Info("Database migration completed")

	// Initialize dependency container
	appContainer := container.NewContainer(cfg, db, appLogger)
	appLogger.Info("Application dependencies initialized")

	// Setup Gin router
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware in order
	router.Use(middleware.RequestID())
	router.Use(middleware.LoggingMiddleware(appLogger.Slog()))
	router.Use(middleware.Recovery(appLogger))
	router.Use(middleware.ErrorHandler(appLogger))
	router.Use(middleware.CORS(cfg))

	// Setup all application routes
	routeConfig := &routes.RouteConfig{
		AuthHandler:        appContainer.AuthHandler,
		UserHandler:        appContainer.UserHandler,
		VideoHandler:       appContainer.VideoHandler,
		OAuthHandler:       appContainer.OAuthHandler,
		TranslationHandler: appContainer.TranslationHandler,
		AuthMiddleware:     appContainer.AuthMiddleware,
	}
	routes.SetupRoutes(router, routeConfig)
	appLogger.Info("Routes configured successfully")

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", zap.String("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited")
}