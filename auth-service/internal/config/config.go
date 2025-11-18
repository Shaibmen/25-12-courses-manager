package config

import (
	"os"

	"github.com/joho/godotenv"
)

const minSecretKeySize = 32

type Config struct {
	JWTKey string
	DB     string
}

func MustInitConfig() *Config {
	_ = godotenv.Load()

	JWTKey := os.Getenv("JWT_SECRET_KEY")

	DBConnection := os.Getenv("DB")

	if len(JWTKey) < minSecretKeySize {
		panic("invalid key size")
	}

	if DBConnection == "" {
		panic("ошибка инициализации конфига")
	}

	return &Config{
		JWTKey: JWTKey,
		DB:     DBConnection,
	}
}
