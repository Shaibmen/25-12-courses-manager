package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/service"
	"online-courses/internal/mapping"
	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"online-courses/internal/validate"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListenerHandler struct {
	handler service.ListenerService
}

func NewListenerHandler(handler service.ListenerService) *ListenerHandler {
	return &ListenerHandler{handler: handler}
}

func (h *ListenerHandler) CreateListenerHandler(c *gin.Context) {

	var request request.FullListenerRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.Error(err)
		return
	}

	err := validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto, err := mapping.MapListenerReqToDto(request)
	if err != nil {
		c.Error(err)
		return
	}

	var idLegal *uuid.UUID = nil
	if request.Listener.ID_LegalEntity != uuid.Nil {
		idLegal = &request.Listener.ID_LegalEntity
	}

	var idContractor *uuid.UUID = nil
	if request.Listener.ID_Contractor != uuid.Nil {
		idContractor = &request.Listener.ID_Contractor
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := h.handler.CreateFullListener(ctx, dto, idLegal, idContractor); err != nil {
		c.Error(err)
		return

	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "слушатель создан"})
}

func (h *ListenerHandler) GetListener(c *gin.Context) {

	idParam := c.Query("page")
	page, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := h.handler.ReadListener(ctx, page, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

func (h *ListenerHandler) GetFullListener(c *gin.Context) {

	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := h.handler.ReadFullListener(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

func (h *ListenerHandler) UpdateListener(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var request request.FullListenerRequest

	if err = c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err = validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto, err := mapping.MapListenerReqToDto(request)
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err = h.handler.UpdateListener(ctx, dto, id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "слушатель успешно обновлён"})
}

func (h *ListenerHandler) DeleteListenerHandler(c *gin.Context) {

	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err = h.handler.DeleteListener(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "слушатель удалён"})

}
