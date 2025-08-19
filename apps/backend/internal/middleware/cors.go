package middleware

import (
	"app-backend/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}

	// For development, allow all origins to support Chrome extensions
	if cfg.App.Environment == "development" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = cfg.CORS.AllowedOrigins
	}

	return cors.New(config)
}