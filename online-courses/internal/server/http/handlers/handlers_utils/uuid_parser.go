package handlers_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseUUID(c *gin.Context, param string) (uuid.UUID, error) {
	idParam := c.Param(param)
	return uuid.Parse(idParam)
}
