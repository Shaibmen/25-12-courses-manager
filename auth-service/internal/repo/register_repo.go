package repo

import (
	"auth-service/internal/domain/entity"
	"auth-service/internal/utils"

	"gorm.io/gorm"
)

type RegisterRepo struct {
	db *gorm.DB
}

func NewRegisterRepo(db *gorm.DB) *RegisterRepo {
	return &RegisterRepo{db: db}
}

func (r *RegisterRepo) Register(username, password, role string) error {

	hash_password, err := utils.PasswordHash(password)
	if err != nil {
		return err
	}

	result := r.db.Create(&entity.User{Username: username, Password: string(hash_password), Role: role})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
