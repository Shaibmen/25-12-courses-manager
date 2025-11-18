package server

import (
	"document-service/internal/server/handler"
	"document-service/internal/server/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MustServerInit(r *gin.Engine, port string, PersonalCard *handler.PersonalCardHandler) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("v1/doc")
	{
		api.Use(middleware.AuthMiddleware())

		api.POST("personal-card", PersonalCard.CreatePersonalCard)

		api.GET("exists", PersonalCard.ExistsPersonalCard)
		api.GET("download", PersonalCard.DownloadPersonalCard)
		api.DELETE("delete", PersonalCard.DeletePersonalCard)
	}

	r.Run(port)

}
