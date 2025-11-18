package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey: secretKey}
}

func (maker *JWTMaker) CreateJwt(id int, username, role string, duration time.Duration) (string, *UserClaims, error) {

	claims, err := NewUserClaims(id, username, role, duration)
	if err != nil {
		return "", nil, errors.New("ошибка создания токена")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, errors.New("ошибка подписи токена")
	}

	return tokenString, claims, nil
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token signing method")
		}

		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
