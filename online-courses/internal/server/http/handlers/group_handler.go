package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	handler service.GroupService
}

func NewGroupHandler(handler service.GroupService) *GroupHandler {
	return &GroupHandler{handler: handler}
}

func (g *GroupHandler) CreateGroup(c *gin.Context) {

	var request request.GroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dto := dto.GroupDTO{
		NameGroup:  request.NameGroup,
		Raspisanie: request.Raspisanie,
	}

	if err := g.handler.Create(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "Группа создана"})
}

func (g *GroupHandler) UpdateGroup(c *gin.Context) {

	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var request request.GroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dto := dto.GroupDTO{
		ID_Group:   id,
		NameGroup:  request.NameGroup,
		Raspisanie: request.Raspisanie,
	}

	if err := g.handler.Update(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "Группа обновлена"})
}

func (g *GroupHandler) DeleteGroup(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := g.handler.Delete(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "Группа удалена"})
}

func (g *GroupHandler) ReadGroup(c *gin.Context) {

	idParam := c.Query("page")
	page, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := g.handler.Read(ctx, page, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

func (g *GroupHandler) ReadByIDGroup(c *gin.Context) {

	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()
	data, err := g.handler.ReadByID(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}
