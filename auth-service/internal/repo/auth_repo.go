package repo

import (
	"auth-service/internal/domain/entity"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (l *AuthRepo) Login(ctx context.Context, username, password string) (int, string, error) {
	query := `
       select u.id_user, u.user_name, u."password", r.role
		from "user" as u
		inner join role r on u.role = r.id_role
		where u.user_name = $1
    `
	//пароль -> хеш
	var user entity.User

	row := l.db.WithContext(ctx).Raw(query, username).Row()
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		return 0, "", err
	}

	password_byte := []byte(password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), password_byte)

	if err != nil {
		return 0, "", err
	}
	return user.ID, user.Role, nil
}
