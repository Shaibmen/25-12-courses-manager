package handlers

import (
	"context"
	"online-courses/internal/domain/service"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	handler service.ReportService
}

func NewReportHandler(service service.ReportService) *ReportHandler {
	return &ReportHandler{handler: service}
}

func (r *ReportHandler) ReportPeriod(c *gin.Context) {

	startDate, err := time.Parse("2006-01-02", c.Query("start"))
	if err != nil {
		c.Error(err)
		return
	}
	endDate, err := time.Parse("2006-01-02", c.Query("end"))
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)

	defer cancel()

	filepath, err := r.handler.ReportPeriod(ctx, startDate, endDate)
	if err != nil {
		c.Error(err)
		return
	}

	c.File(filepath)
}

func (r *ReportHandler) MostExpensiveProgram(c *gin.Context) {
	startDate, err := time.Parse("2006-01-02", c.Query("start"))
	if err != nil {
		c.Error(err)
		return
	}
	endDate, err := time.Parse("2006-01-02", c.Query("end"))
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)

	defer cancel()

	filePath, err := r.handler.MostExpensiveProgram(ctx, startDate, endDate)
	if err != nil {
		c.Error(err)
		return
	}

	c.File(filePath)

}
