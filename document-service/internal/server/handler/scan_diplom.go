package handler

import (
	"context"
	"document-service/internal/server/models"
	"document-service/internal/service"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ScanDiplomHandler struct {
	handler *service.ScanDiplomService
}

func NewScanDiplomHandler(handler *service.ScanDiplomService) *ScanDiplomHandler {
	return &ScanDiplomHandler{handler: handler}
}

func (s *ScanDiplomHandler) ScanDiplomCreate(c *gin.Context) {

	var request models.ScanDiplomRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "creation scan diplom error", "err": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	namePDF := request.Name + ".pdf"

	err := s.handler.CreateScanDiplom(ctx, request.File, namePDF)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "creation scan diplom error", "err": err.Error()})
		return
	}
}

func (s *ScanDiplomHandler) ScanDiplomDownload(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	param := c.Query("dogovor-name")

	doc, err := s.handler.DownloadScanDiplom(ctx, param)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusNotFound, gin.H{"message": "не удалось загрузить файл"})
		log.Println("ошибка при загрузке файла:", err)
		return
	}

	c.Data(200, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", doc)
}

func (s *ScanDiplomHandler) ScanDiplomExists(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	param := c.Query("dogovor-name")

	files, err := s.handler.ExistsScanDiplom(ctx, param)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "скана не существует"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (s *ScanDiplomHandler) DeleteScanDiplom(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	param := c.Query("dogovor-name")

	namePDF := param + ".pdf"
	_, err := s.handler.ExistsScanDiplom(ctx, namePDF)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "скана не существует"})
		return
	}

	err = s.handler.DeleteScanDiplom(namePDF)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверные данные скана дела"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
