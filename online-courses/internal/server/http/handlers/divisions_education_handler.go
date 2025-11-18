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

type DivisionsEducationHandler struct {
	handler service.DivisionsEducationService
}

func NewDivisionsEducationHandler(divisions service.DivisionsEducationService) *DivisionsEducationHandler {
	return &DivisionsEducationHandler{handler: divisions}
}

func (d *DivisionsEducationHandler) CreateDivision(c *gin.Context) {

	var request request.DivisionsEducationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err := validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.DivisionsDTO{
		Divisions: request.Divisions,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := d.handler.CreateDivision(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "подразделение добавлено"})
}

func (d *DivisionsEducationHandler) ReadDivisions(c *gin.Context) {

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := d.handler.ReadDivisions(ctx, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

// func (d *DivisionsEducationHandler) ReadDivisionsByID(c *gin.Context) {

// 	idParam, err := utils.ParseUUID(c, "id")
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
// 	defer cancel()

// 	data, err := d.handler.ReadByID(ctx, idParam)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
// }

func (d *DivisionsEducationHandler) UpdateDivision(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var request request.DivisionsEducationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err = validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.DivisionsDTO{
		Divisions: request.Divisions,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := d.handler.UpdateDivision(ctx, id, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "подразделение обновлено"})
}

func (d *DivisionsEducationHandler) DeleteDivision(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := d.handler.DeleteDivision(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "подразделние удалено"})
}
