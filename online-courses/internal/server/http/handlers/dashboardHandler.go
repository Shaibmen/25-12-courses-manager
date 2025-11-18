package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/service"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	handler service.DashboardService
}

func NewDashboardHandler(service service.DashboardService) *DashboardHandler {
	return &DashboardHandler{handler: service}
}

func (d *DashboardHandler) UserDashboard(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()
	dto, err := d.handler.UserDashboard(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

func (d *DashboardHandler) AdminDashboard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()
	dto, err := d.handler.AdminDashboard(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto)
}
