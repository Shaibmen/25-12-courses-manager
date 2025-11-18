package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewUserClaims(id int, username, role string, duration time.Duration) (*UserClaims, error) {

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &UserClaims{
		ID:       id,
		UserName: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
