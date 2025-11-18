package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleProtecteMiddleware(allowedRole ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "отсутсвует роль"})
			return
		}

		for _, r := range allowedRole {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "доступ запрещён"})
	}
}
