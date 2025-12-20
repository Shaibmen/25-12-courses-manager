package config

import (
	"os"

	"github.com/joho/godotenv"
)

const minSecretKeySize = 32

type Config struct {
	DB_STRING_CONN string
	DBUSER         string
	DBPASSWORD     string
	DBNAME         string
	JWTKey         string
	HOST           string
	SERVICE_CORE   string
	SERVICE_DOC    string
	SERVICE_AUTH   string
	SERVICE_FRONT  string
}

func MustInitConfig() *Config {
	_ = godotenv.Load("../.env")

	JWTKey := os.Getenv("JWT_SECRET_KEY")

	Host := os.Getenv("HOST")

	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")

	if len(JWTKey) < minSecretKeySize {
		panic("invalid key size")
	}

	if Host == "" {
		panic("хост не указан")
	}

	SERVICE_CORE, SERVICE_DOC, SERVICE_AUTH, SERVICE_FRONT, DB_HOST := "http://localhost", "http://localhost", "http://localhost", "http://localhost", "localhost"

	if Host != "LOCAL_HOST" {
		SERVICE_CORE, SERVICE_DOC, SERVICE_AUTH, SERVICE_FRONT, DB_HOST = "http://apicore", "http://apidoc", "http://apiauth", "http://frontend", "db"
	}

	dbStringConn := "host=" + DB_HOST + " " + "port=5432" + " " + "user=" + DBUser + " " + "password=" + DBPassword + " " + "dbname=" + DBName + " " + "sslmode=disable"

	return &Config{
		JWTKey:         JWTKey,
		DB_STRING_CONN: dbStringConn,
		HOST:           Host,
		SERVICE_CORE:   SERVICE_CORE,
		SERVICE_DOC:    SERVICE_DOC,
		SERVICE_AUTH:   SERVICE_AUTH,
		SERVICE_FRONT:  SERVICE_FRONT,
	}
}
