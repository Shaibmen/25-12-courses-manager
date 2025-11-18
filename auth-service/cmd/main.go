package main

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/repo"
	"auth-service/internal/server"
	"auth-service/internal/server/handler"
	"auth-service/internal/utils"
)

func main() {

	config := config.MustInitConfig()

	db := database.MustInitDB(config.DB)

	utils.InitdValidator()

	registerRepo := repo.NewRegisterRepo(db)
	authRepo := repo.NewAuthRepo(db)
	sessionRepo := repo.NewSessionRepo(db)

	registerHandler := handler.NewRegisterHandler(registerRepo, config.JWTKey)
	authHandler := handler.NewAuthHandler(authRepo, sessionRepo, config.JWTKey)

	middleware := handler.NewMiddleware(config.JWTKey)

	server.InitSever(authHandler, registerHandler, middleware)

}
