package handler

import (
	"auth-service/internal/domain/dto"
	"auth-service/internal/domain/entity"
	"auth-service/internal/server/models"
	"auth-service/internal/server/token"
	"auth-service/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthRepository interface {
	Login(ctx context.Context, username, passport string) (int, string, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, dto *dto.SessionDTO) (*dto.SessionDTO, error)
	ReadSession(ctx context.Context, id string) (*entity.Session, error)
	RevokeSession(ctx context.Context, id string) error
	DeleteSession(ctx context.Context, id string) error
	DeleteSessionByUsername(ctx context.Context, username string) error
}

type AuthHandler struct {
	auth       AuthRepository
	session    SessionRepository
	tokenMaker *token.JWTMaker
}

func NewAuthHandler(repo AuthRepository, session SessionRepository, secretKey string) *AuthHandler {
	return &AuthHandler{
		auth:       repo,
		session:    session,
		tokenMaker: token.NewJwtMaker(secretKey)}
}

func (a *AuthHandler) LoginHandler(c *gin.Context) {

	ctx := c.Request.Context()

	var request models.UserAuth

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := utils.Validate.Struct(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, role, err := a.auth.Login(ctx, request.UserName, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, accessClaims, err := a.tokenMaker.CreateJwt(id, request.UserName, role, 12*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, refreshClaims, err := a.tokenMaker.CreateJwt(id, request.UserName, role, 12*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := a.session.CreateSession(ctx, &dto.SessionDTO{
		ID:           refreshClaims.RegisteredClaims.ID,
		Username:     request.UserName,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	})

	if err != nil {
		if strings.Contains(err.Error(), "23505") {

			if err = a.session.DeleteSessionByUsername(ctx, request.UserName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "error creating session"})
				return
			}

			session, err = a.session.CreateSession(ctx, &dto.SessionDTO{
				ID:           refreshClaims.RegisteredClaims.ID,
				Username:     request.UserName,
				RefreshToken: refreshToken,
				IsRevoked:    false,
				ExpiresAt:    refreshClaims.ExpiresAt.Time,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "error creating session"})
				return
			}

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error creating session"})
			return
		}
	}

	response := models.LoginUserRes{
		ID:                   session.ID,
		Refreshtoken:         refreshToken,
		AccessTokenExpireAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpireAt: refreshClaims.ExpiresAt.Time,
		AccessToken:          accessToken,
		User: models.UserRes{
			UserName: request.UserName,
			Role:     role,
		},
	}

	c.JSON(200, gin.H{"token": response})
}

func (a *AuthHandler) ValidateHandler(c *gin.Context) {

	var request models.JwtToken

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := a.tokenMaker.VerifyToken(request.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	responseClaims := &models.ResponseClaims{
		ID:       claims.ID,
		UserName: claims.UserName,
		Role:     claims.Role,
	}

	c.JSON(http.StatusOK, responseClaims)
}

func (a *AuthHandler) LogoutUser(c *gin.Context) {

	ctx := c.Request.Context()

	claims := c.Request.Context().Value(authKey{}).(*token.UserClaims)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "empty authorization header"})
		return
	}
	fmt.Println(claims.RegisteredClaims.ID)

	if err := a.session.DeleteSession(ctx, claims.RegisteredClaims.ID); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(200, nil)
}

func (a *AuthHandler) ReNewAccessToken(c *gin.Context) {

	ctx := c.Request.Context()

	var request models.ReNewAccessTokenReq

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error decoding request body"})
		return
	}

	refreshClaims, err := a.tokenMaker.VerifyToken(request.RefreshRoken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	session, err := a.session.ReadSession(ctx, refreshClaims.RegisteredClaims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if session.IsRevoked {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "session revoked"})
		return
	}

	if session.Username != refreshClaims.UserName {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid session"})
		return
	}

	accessToken, refreshToken, err := a.tokenMaker.CreateJwt(refreshClaims.ID, refreshClaims.UserName, refreshClaims.Role, 2*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error creating token"})
		return
	}

	response := &models.ReNewAccessTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}

func (a *AuthHandler) RevokeSession(c *gin.Context) {

	ctx := c.Request.Context()

	claims := c.Request.Context().Value(authKey{}).(*token.UserClaims)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "empty authorization header"})
		return
	}

	if err := a.session.RevokeSession(ctx, claims.RegisteredClaims.ID); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "session revoke"})
}
