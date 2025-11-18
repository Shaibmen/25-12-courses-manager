package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustInitDB(ConnectionSting string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(ConnectionSting), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
