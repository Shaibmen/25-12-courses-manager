package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"time"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	handler service.DocumentService
}

func NewDocumentHandler(handler service.DocumentService) *DocumentHandler {
	return &DocumentHandler{handler: handler}
}

func (d *DocumentHandler) DocumentDataDeliver(c *gin.Context) {

	var request request.DocumentsDataRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := d.handler.FormingDataDocument(ctx, request.ID_Listener, request.ID_Program, request.ID_Executor, dto.FrontDataDeliver(request.FrontData))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}
