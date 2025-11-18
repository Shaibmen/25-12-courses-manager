package handler

import (
	"auth-service/internal/server/token"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authKey struct{}

type Middleware struct {
	token *token.JWTMaker
}

func NewMiddleware(secretKey string) *Middleware {
	return &Middleware{token: token.NewJwtMaker(secretKey)}
}

func (m *Middleware) GetAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, err := verifyClaimsFromHeader(c, m.token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), authKey{}, claims)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func verifyClaimsFromHeader(c *gin.Context, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {

	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, errors.New("invalid authorization header")
	}

	token := fields[1]

	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
