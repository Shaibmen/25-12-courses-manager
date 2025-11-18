package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/service"
	"online-courses/internal/server/http/models"
	"time"

	"github.com/gin-gonic/gin"
)

type LevelEducationHandler struct {
	handler service.LevelEducationService
}

func NewLevelEducationHandler(service service.LevelEducationService) *LevelEducationHandler {
	return &LevelEducationHandler{handler: service}
}

func (lvl *LevelEducationHandler) GetLevelEducations(c *gin.Context) {
	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := lvl.handler.ReadAll(ctx, filter)
	if err != nil {
		// errorResponse := utils.DefineError(err)
		// c.JSON(errorResponse.Status, models.HttpResponse{Message: errorResponse.Message})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}
