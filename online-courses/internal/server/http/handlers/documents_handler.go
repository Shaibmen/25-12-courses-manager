package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	"online-courses/internal/server/http/request"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	handler service.DocumentService
}

func NewDocumentHandler(handler service.DocumentService) *DocumentHandler {
	return &DocumentHandler{handler: handler}
}

func (d *DocumentHandler) DocumentDataDeliver(c *gin.Context) {

	var request request.DocumentsDataRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := d.handler.FormingDataDocument(ctx, request.ID_Listener, request.ID_Program, request.ID_Executor, dto.FrontDataDeliver(request.FrontData))
	if err != nil {
		c.Error(err)
		return
	}

	responseCard, err := RequestToDoc(*data, "personal-card", c)
	if err != nil {
		c.Error(err)
		return
	}
	if responseCard != http.StatusOK {
		c.JSON(responseCard, nil)
		return
	}

	responseZayvlenie, err := RequestToDoc(*data, "zayavlenie", c)
	if err != nil {
		c.Error(err)
		return
	}
	if responseZayvlenie != http.StatusOK {
		c.JSON(responseZayvlenie, nil)
		return
	}

	responseDogovor, err := RequestToDoc(*data, "dogovor", c)
	if err != nil {
		c.Error(err)
		return
	}

	if responseDogovor != http.StatusOK {
		c.JSON(responseDogovor, nil)
		return
	}

}

func RequestToDoc(data dto.FullDocumentInfoDTO, endpoint string, c *gin.Context) (int, error) {
	requestBody, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://apidoc:8082/v1/doc/"+endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, err
	}

	authHeader := c.Request.Header.Get("Authorization")
	field := strings.Fields(authHeader)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+field[1])

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}
