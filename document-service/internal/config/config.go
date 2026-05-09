package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT          string
	HOST          string
	SERVICE_CORE  string
	SERVICE_DOC   string
	SERVICE_AUTH  string
	SERVICE_FRONT string
}

func MustInitConfig() *Config {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Printf("ошибка при загрузке конфига: %s", err.Error())
	}

	PORT := os.Getenv("PORT")

	Host := os.Getenv("HOST")

	if Host == "" {
		panic("хост не указан")
	}

	SERVICE_CORE, SERVICE_DOC, SERVICE_AUTH, SERVICE_FRONT := "http://localhost", "http://localhost", "http://localhost", "http://localhost"

	if Host != "LOCAL_HOST" {
		SERVICE_CORE, SERVICE_DOC, SERVICE_AUTH, SERVICE_FRONT = "http://176.108.244.99", "http://176.108.244.99", "http://176.108.244.99", "http://176.108.244.99"
	}

	return &Config{
		PORT:          PORT,
		HOST:          Host,
		SERVICE_CORE:  SERVICE_CORE,
		SERVICE_DOC:   SERVICE_DOC,
		SERVICE_AUTH:  SERVICE_AUTH,
		SERVICE_FRONT: SERVICE_FRONT,
	}
}
