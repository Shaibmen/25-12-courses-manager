package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	"online-courses/internal/validate"

	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProgramEducationHandler struct {
	handler service.ProgramEducationService
}

func NewProgramEducationHandler(service service.ProgramEducationService) *ProgramEducationHandler {
	return &ProgramEducationHandler{handler: service}
}

func (p *ProgramEducationHandler) CreateProgram(c *gin.Context) {

	var request request.ProgramEducationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err := validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.ProgramEducationDTO{
		NameProfEducation:     request.NameProfEducation,
		TimeEducation:         request.TimeEducation,
		IndividualPrice:       request.IndividualPrice,
		GroupPrice:            request.GroupPrice,
		CampusPrice:           request.CampusPrice,
		ID_EducationType:      request.ID_EducationType,
		ID_DivisionsEducation: request.ID_DivisionsEducation,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := p.handler.CreateProgram(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "программа обучения создана"})
}

func (p *ProgramEducationHandler) ReadProgram(c *gin.Context) {
	idParam := c.Query("page")
	page, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := p.handler.ReadProgram(ctx, page, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

func (p *ProgramEducationHandler) ReadByID(c *gin.Context) {
	idParam, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := p.handler.ReadByID(ctx, idParam)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

func (p *ProgramEducationHandler) UpdateProgram(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var request request.ProgramEducationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err = validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.ProgramEducationDTO{
		ID_ProgramEducation:   id,
		NameProfEducation:     request.NameProfEducation,
		TimeEducation:         request.TimeEducation,
		IndividualPrice:       request.IndividualPrice,
		GroupPrice:            request.GroupPrice,
		CampusPrice:           request.CampusPrice,
		ID_EducationType:      request.ID_EducationType,
		ID_DivisionsEducation: request.ID_DivisionsEducation,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := p.handler.UpdateProgram(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "программа обучения обновлена"})
}

func (p *ProgramEducationHandler) DeleteProgram(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := p.handler.DeleteProgram(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "программа обучения удалена"})
}
