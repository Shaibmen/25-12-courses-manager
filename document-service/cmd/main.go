package main

import (
	"document-service/internal/config"
	"document-service/internal/server"
	"document-service/internal/server/handler"
	"document-service/internal/service"
	"document-service/internal/validate"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.MustInitConfig()

	validate.InitValid()

	r := gin.Default()

	s3Client := service.MustInitS3Client()
	personalCardService := service.NewPersonalCardService(s3Client)
	zayavlenieService := service.NewZayavlenieService(s3Client)
	dogovorService := service.NewDogovorService(s3Client)

	dogovorHandler := handler.NewDogovorHandler(dogovorService)
	zayavlenieHandler := handler.NewZayavlenieHandler(zayavlenieService)
	personalCardHandler := handler.NewPersonalCardHandler(personalCardService)

	if _, err := os.Stat("./personal_card"); os.IsNotExist(err) {
		os.MkdirAll("./personal_card", 0755)
	}

	server.MustServerInit(r, cfg.PORT, personalCardHandler, zayavlenieHandler, dogovorHandler)
}
