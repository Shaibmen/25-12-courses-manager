package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"online-courses/internal/server/http/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestClaims struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.HttpResponse{Message: "заголовка авторизации нет"})
			c.Abort()
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || fields[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.HttpResponse{Message: "невалидный заголовок авторизации"})
			c.Abort()
			return
		}

		token := fields[1]

		requestBody, _ := json.Marshal(map[string]string{"token": token})

		resp, err := http.Post("http://apiauth:8081/auth/v1/valid", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.HttpResponse{Message: "invalid serivce"})
			c.Abort()
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, models.HttpResponse{Message: "unauthorized"})
			c.Abort()
			return
		}

		var claims RequestClaims

		if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("username", claims.UserName)
		c.Set("role", claims.Role)

		c.Next()

	}
}
