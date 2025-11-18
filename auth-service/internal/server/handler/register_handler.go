package handler

import (
	"auth-service/internal/server/models"
	"auth-service/internal/server/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegisterRepository interface {
	Register(username, password, role string) error
}

type RegisterHandler struct {
	handler    RegisterRepository
	tokenMaker *token.JWTMaker
}

func NewRegisterHandler(repo RegisterRepository, secretKey string) *RegisterHandler {
	return &RegisterHandler{handler: repo, tokenMaker: token.NewJwtMaker(secretKey)}
}

func (r *RegisterHandler) Registration(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "нет заголовка авторизации"})
		c.Abort()
		return
	}

	fields := strings.Fields(authHeader)
	if len(fields) < 2 || fields[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "евалидный заголовок авторизации"})
		c.Abort()
		return
	}

	token := fields[1]
	claims, err := r.tokenMaker.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "неверный токен"})
		return
	}

	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "доступ запрещён"})
		return
	}

	var request models.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := r.handler.Register(request.UserName, request.Password, request.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"message": "регистрация успешна"})
}
