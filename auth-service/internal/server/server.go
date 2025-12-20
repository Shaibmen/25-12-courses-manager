package server

import (
	"auth-service/internal/config"
	"auth-service/internal/server/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitSever(AuthHandler *handler.AuthHandler, RegisterHandler *handler.RegisterHandler, Middleware *handler.Middleware, config *config.Config) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.SERVICE_FRONT + ":5173", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	version := "auth/v1"

	api := r.Group(version)
	{

		api.POST("/login", AuthHandler.LoginHandler)
		api.POST("/valid", AuthHandler.ValidateHandler)
		api.POST("/register", RegisterHandler.Registration)
		api.POST("/new-access-token", AuthHandler.ReNewAccessToken)

		protected := r.Group(version + "/protected")
		{
			protected.Use(Middleware.GetAuthMiddleware())
			protected.POST("/logout", AuthHandler.LogoutUser)
			protected.POST("/revoke", AuthHandler.RevokeSession)
		}
	}

	r.Run(":8081")
}
