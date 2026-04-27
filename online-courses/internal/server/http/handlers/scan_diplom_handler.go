package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"online-courses/internal/config"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"strings"

	"github.com/gin-gonic/gin"
)

type ScanDiplom struct {
	cfg *config.Config
}

func NewScanDiplom(cfg *config.Config) *ScanDiplom {
	return &ScanDiplom{cfg: cfg}
}

func (s *ScanDiplom) ScanDiplomHandler(c *gin.Context) {
	Param := c.Query("name-scan")

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		c.Error(err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.Error(err)
		return
	}
	defer file.Close()

	filebytes, err := io.ReadAll(file)
	if err != nil {
		c.Error(err)
		return
	}

	data := request.ScanDiplomRequest{
		File: filebytes,
		Name: Param,
	}

	requestBody, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", s.cfg.SERVICE_DOC+":8082/v1/doc/"+"scan-diplom", bytes.NewBuffer(requestBody))
	if err != nil {
		c.Error(err)
		return
	}

	authHeader := c.Request.Header.Get("Authorization")
	field := strings.Fields(authHeader)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+field[1])

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(response.StatusCode, models.HttpResponse{Message: "Скан загружен"})
}
