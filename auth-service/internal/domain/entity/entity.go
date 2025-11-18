package entity

import (
	"time"
)

type User struct {
	ID       int    `gorm:"column:id_user;primaryKey;autoIncrement"`
	Username string `gorm:"column:user_name"`
	Password string `gorm:"column:password"`
	Role     string `gorm:"foreignKey:role column:role"`
}

func (User) TableName() string {
	return "user"
}

type Session struct {
	ID           string    `gorm:"column:id_session"`
	Username     string    `gorm:"column:user_name"`
	RefreshToken string    `gorm:"column:refresh_token"`
	IsRevoked    bool      `gorm:"column:is_revoked"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
}

func (Session) TableName() string {
	return "session"
}
