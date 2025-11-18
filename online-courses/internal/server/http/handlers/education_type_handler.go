package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"online-courses/internal/validate"
	"time"

	"github.com/gin-gonic/gin"
)

type EducationTypeHandler struct {
	handler service.EducationTypeService
}

func NewEducationTypeHandler(educationType service.EducationTypeService) *EducationTypeHandler {
	return &EducationTypeHandler{handler: educationType}
}

func (e *EducationTypeHandler) CreateEducationType(c *gin.Context) {

	var request request.EducationTypeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err := validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.EducationTypeDTO{
		TypeName: request.TypeName,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := e.handler.CreateType(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "тип обучения создан"})
}

func (e *EducationTypeHandler) ReadEducationType(c *gin.Context) {

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := e.handler.ReadEducationType(ctx, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

// func (e *EducationTypeHandler) ReadEducationTypeByID(c *gin.Context) {

// 	idParam, err := utils.ParseUUID(c, "id")
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
// 	defer cancel()

// 	data, err := e.handler.ReadEducationTypeByID(ctx, idParam)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
// }

func (e *EducationTypeHandler) UpdateEducationType(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var request request.EducationTypeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err = validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.EducationTypeDTO{
		TypeName: request.TypeName,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := e.handler.UpdateEducationType(ctx, id, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "тип обучения обновлен"})
}

func (e *EducationTypeHandler) DeleteEducationtype(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := e.handler.DeleteEducationType(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "тип обучения удален"})
}
