package config

import (
	"os"
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
}

func MustInitConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		panic(fmt.Sprintf("ошибка при загрузке конфига: %s", err.Error()))
	}

	PORT := os.Getenv("PORT")

	return &Config{PORT: PORT}
}
