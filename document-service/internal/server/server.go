package server

import (
	"document-service/internal/server/handler"
	"document-service/internal/server/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MustServerInit(r *gin.Engine, port string, personalCard *handler.PersonalCardHandler, zayavlenie *handler.ZayavlenieHandler, dogovor *handler.DogovoreHandler) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.AuthMiddleware())
	api := r.Group("v1/doc")
	{

		api.POST("personal-card", personalCard.CreatePersonalCard)
		api.GET("exists", personalCard.ExistsPersonalCard)
		api.GET("download", personalCard.DownloadPersonalCard)
		api.DELETE("delete", personalCard.DeletePersonalCard)

		// api.POST("zayavlenie", zayavlenie.CreateZayavlenie)
		// api.GET("zayavlenie-exists", zayavlenie.ExistsZayavlenie)
		// api.GET("zayavlenie-download", zayavlenie.DownloadZayavlenie)
		// api.DELETE("zayavlenie-delete", zayavlenie.DeleteZayavlenie)

	}

	r.POST("dogovor", dogovor.CreateDogovor)
	r.GET("dogovor-exists", dogovor.ExistsDogovor)
	r.GET("dogovor-download", dogovor.DownloadDogovor)
	r.DELETE("dogovor-delete", dogovor.DeleteDogovor)

	r.POST("zayavlenie", zayavlenie.CreateZayavlenie)
	r.GET("zayavlenie-exists", zayavlenie.ExistsZayavlenie)
	r.GET("zayavlenie-download", zayavlenie.DownloadZayavlenie)
	r.DELETE("zayavlenie-delete", zayavlenie.DeleteZayavlenie)

	r.Run(port)

}
