package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_STRING_CONN string
	Logger         string
	JWT            string
	DBUSER         string
	DBPASSWORD     string
	DBNAME         string
	DBHOST         string
	HOST           string
	SERVICE_CORE   string
	SERVICE_DOC    string
	SERVICE_AUTH   string
	SERVICE_FRONT  string
}

func MustLoadConfig() *Config {

	_ = godotenv.Load("../.env")

	logger := os.Getenv("LEVEL_LOG")
	JWTKey := os.Getenv("JWT_SECRET_KEY")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")

	if DBUser == "" || logger == "" || JWTKey == "" {
		log.Panic("конфиг не загружен")
	}

	Host := os.Getenv("HOST")

	if Host == "" {
		panic("хост не указан")
	}

	SERVICE_CORE, SERVICE_DOC, SERVICE_AUTH, SERVICE_FRONT, DB_HOST := "http://localhost", "http://localhost", "http://localhost", "http://localhost", "localhost"

	if Host != "LOCAL_HOST" {
		SERVICE_CORE, SERVICE_DOC, SERVICE_AUTH, SERVICE_FRONT, DB_HOST = "http://176.108.244.99", "http://176.108.244.99", "http://176.108.244.99", "http://176.108.244.99", "db"
	}

	dbStringConn := "host=" + DB_HOST + " " + "port=5432" + " " + "user=" + DBUser + " " + "password=" + DBPassword + " " + "dbname=" + DBName + " " + "sslmode=disable"

	log.Println("конфиг инициализирован")
	return &Config{
		DB_STRING_CONN: dbStringConn,
		Logger:         logger,
		JWT:            JWTKey,
		DBUSER:         DBUser,
		DBPASSWORD:     DBPassword,
		DBNAME:         DBName,
		DBHOST:         DB_HOST,
		HOST:           Host,
		SERVICE_CORE:   SERVICE_CORE,
		SERVICE_DOC:    SERVICE_DOC,
		SERVICE_AUTH:   SERVICE_AUTH,
		SERVICE_FRONT:  SERVICE_FRONT,
	}

}
