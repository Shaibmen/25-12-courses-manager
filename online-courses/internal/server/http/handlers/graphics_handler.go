package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/service"
	"online-courses/internal/server/http/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GraphicsHandler struct {
	handler service.GraphicsService
}

func NewGraphicsHandler(handler service.GraphicsService) *GraphicsHandler {
	return &GraphicsHandler{handler: handler}
}

func (g *GraphicsHandler) CountListenersOnProgram(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.CountListenersOnProgram(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}

func (g *GraphicsHandler) PopularProgramType(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.PopularProgramType(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}

func (g *GraphicsHandler) CountListenersOnProgramAccurate(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.CountListenersOnProgramAccurate(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}

func (g *GraphicsHandler) WorthProgramAccurate(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.WorthProgramAccurate(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}

func (g *GraphicsHandler) AgeDiff(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.AgeDiff(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}

func (g *GraphicsHandler) WhoEnrolled(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.WhoEnrolled(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}

func (g *GraphicsHandler) GroupMembers(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.GroupMembers(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}

func (g *GraphicsHandler) DivisionMember(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	data, err := g.handler.DivisionMember(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})

}
