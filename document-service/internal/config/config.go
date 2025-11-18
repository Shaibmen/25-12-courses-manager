package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
}

func MustInitConfig() *Config {
	_ = godotenv.Load()

	PORT := os.Getenv("PORT")

	if PORT == "" {
		panic("ошибка инициализации конфига")
	}

	return &Config{PORT: PORT}
}
