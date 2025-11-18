package models

import (
	"auth-service/internal/server/token"
	"time"
)

type UserAuth struct {
	UserName string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,max=50"`
}

type User struct {
	UserName string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,max=50"`
	Role     string `json:"role" validate:"required, max=50"`
}

type UserRes struct {
	UserName string `json:"username" `
	Role     string `json:"role"`
}
type JwtToken struct {
	Token string `json:"token"`
}

type LoginUserRes struct {
	ID                   string    `json:"id_session"`
	AccessToken          string    `json:"access_token"`
	Refreshtoken         string    `json:"refresh_token"`
	AccessTokenExpireAt  time.Time `json:"access_token_expire_at"`
	RefreshTokenExpireAt time.Time `json:"refresh_token_expire_at"`
	User                 UserRes   `json:"user"`
}

type ReNewAccessTokenReq struct {
	RefreshRoken string `json:"refresh_token"`
}

type ReNewAccessTokenRes struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken *token.UserClaims `json:"refresh_token"`
}

type ResponseClaims struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}
