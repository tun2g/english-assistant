package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwaggerRoutes configures Swagger documentation routes
func SetupSwaggerRoutes(router *gin.Engine) {
	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Redirect root to swagger docs in development
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
}