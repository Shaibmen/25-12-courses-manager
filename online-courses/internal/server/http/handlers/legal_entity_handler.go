package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/service"
	"online-courses/internal/mapping"
	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LegalEntityHandler struct {
	handler service.LegalEntityService
}

func NewLegalEntityHandler(service service.LegalEntityService) *LegalEntityHandler {
	return &LegalEntityHandler{handler: service}
}

func (l *LegalEntityHandler) CreateLegalEntity(c *gin.Context) {

	var request request.FullLegalEntityRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dto, err := mapping.LegalEntityFullMappping(request)
	if err != nil {
		c.Error(err)
		return
	}

	if err = l.handler.Create(ctx, *dto); err != nil {
		c.Error(err)
		return
	}
}

func (l *LegalEntityHandler) ReadLegalEntity(c *gin.Context) {

	idParam := c.Query("page")
	page, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := l.handler.Read(ctx, page, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

func (l *LegalEntityHandler) UpdateLegalEntity(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	var request request.FullLegalEntityRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	data, err := mapping.LegalEntityFullMappping(request)
	if err != nil {
		c.Error(err)
		return
	}

	if err = l.handler.Update(ctx, data, id); err != nil {
		c.Error(err)
		return
	}

}

func (l *LegalEntityHandler) DeleteLegalEntity(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err = l.handler.Delete(ctx, id); err != nil {
		c.Error(err)
		return
	}
}
