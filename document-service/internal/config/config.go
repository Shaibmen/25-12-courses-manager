package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
}

func MustInitConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Printf("ошибка при загрузке конфига: %s", err.Error())
	}

	PORT := os.Getenv("PORT")

	return &Config{PORT: PORT}
}
