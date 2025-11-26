package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/service"
	"online-courses/internal/mapping"
	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"time"

	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	handler service.ContractorService
}

func NewContractHandler(handler service.ContractorService) *ContractHandler {
	return &ContractHandler{handler: handler}
}

func (con *ContractHandler) CreateContract(c *gin.Context) {

	var request request.FullContractorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	dto, err := mapping.MapContcratorReqToDTO(request)
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err = con.handler.Create(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "заказчик создан"})
}

func (con *ContractHandler) UpdateContractor(c *gin.Context) {

	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var request request.FullContractorRequest

	dto, err := mapping.MapContcratorReqToDTO(request)
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err = con.handler.UpdateInTx(ctx, dto, id); err != nil {
		c.Error(err)
		return
	}
}

func (con *ContractHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err = con.handler.Delete(ctx, id); err != nil {
		c.Error(err)
		return
	}

}
