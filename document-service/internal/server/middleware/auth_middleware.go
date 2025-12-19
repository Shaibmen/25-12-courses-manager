package middleware

import (
	"bytes"
	"document-service/internal/config"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	Message string
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, HttpResponse{Message: "authorization header is missing"})
			c.Abort()
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || fields[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, HttpResponse{Message: "invalid authorization header"})
			c.Abort()
			return
		}

		token := fields[1]

		requestBody, _ := json.Marshal(map[string]string{"token": token})

		resp, err := http.Post(cfg.SERVICE_AUTH+":8081/auth/v1/valid", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, HttpResponse{Message: "invalid serivce"})
			c.Abort()
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, HttpResponse{Message: "unauthorized"})
			c.Abort()
			return
		}

		c.Next()

	}
}
