package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"time"

	"github.com/gin-gonic/gin"
)

type ExecutorHandler struct {
	service service.ExecutorService
}

func NewExecutorHandler(service service.ExecutorService) *ExecutorHandler {
	return &ExecutorHandler{service: service}
}

func (e *ExecutorHandler) CreateExecutor(c *gin.Context) {

	var request request.ExecutorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	if request.Doverenost == "" {
		request.Doverenost = "-"
	}

	dto := dto.ExecutorDTO{
		Status:     request.Status,
		FirstName:  request.FirstName,
		SecondName: request.SecondName,
		MiddleName: request.MiddleName,
		Doverenost: request.Doverenost,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	if err := e.service.Create(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "исполнитель создан"})
}

func (e *ExecutorHandler) ReadExecutor(c *gin.Context) {

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := e.service.Read(ctx, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

func (e *ExecutorHandler) DeleteExecutor(c *gin.Context) {

	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err = e.service.Delete(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "исполнитель удалён"})
}
