package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB         string
	Logger     string
	JWT        string
	DBUser     string
	DBPassword string
	DBName     string
	Host       string
}

func MustLoadConfig() *Config {

	_ = godotenv.Load()

	database := os.Getenv("DB")
	logger := os.Getenv("LEVEL_LOG")
	JWTKey := os.Getenv("JWT_SECRET_KEY")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")
	Host := os.Getenv("HOST")

	if database == "" || logger == "" || JWTKey == "" {
		log.Panic("конфиг не загружен")
	}

	log.Println("конфиг инициализирован")
	return &Config{
		DB:         database,
		Logger:     logger,
		JWT:        JWTKey,
		DBUser:     DBUser,
		DBPassword: DBPassword,
		DBName:     DBName,
		Host:       Host,
	}
}
