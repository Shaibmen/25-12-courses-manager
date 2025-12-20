package handler

import (
	"document-service/internal/mapper"
	"document-service/internal/server/models"
	"document-service/internal/service"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ZayavlenieHandler struct {
	service *service.ZayavlenieService
}

func NewZayavlenieHandler(service *service.ZayavlenieService) *ZayavlenieHandler {
	return &ZayavlenieHandler{service}
}

func (h *ZayavlenieHandler) CreateZayavlenie(c *gin.Context) {
	var fullRequest models.FullRequest

	if err := c.ShouldBindJSON(&fullRequest); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter request", "err": err.Error()})
		return
	}

	request := fullRequest.ZayavlenieData
	dto := mapper.ZayavlenieMapping(request)

	err := h.service.CreateZayavlenie(dto, request.DogovorAgeType)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "creation zayavlenie error", "err": err.Error()})
		return
	}
}

func (h *ZayavlenieHandler) ExistsZayavlenie(c *gin.Context) {
	param := c.Query("zayavlenie-name")

	files, err := h.service.ExistsZayavlenie(param)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "личного дела не существует"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *ZayavlenieHandler) DownloadZayavlenie(c *gin.Context) {
	// filePath := cardsPath + "/" + param
	//
	// c.File(filePath)
	param := c.Query("zayavlenie-name")

	doc, err := h.service.DownloadZayavlenie(param)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusNotFound, gin.H{"message": "не удалось загрузить файл"})
		log.Println("ошибка при загрузке файла:", err)
		return
	}

	// log.Println("Doc size:", len(doc))

	c.Data(200, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", doc)
}

func (h *ZayavlenieHandler) DeleteZayavlenie(c *gin.Context) {
	param := c.Query("zayavlenie-name")

	file := "Заявление-" + param + ".docx"
	filePath := cardsPath + "/" + file
	_, err := h.service.ExistsZayavlenie(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "личного дела не существует"})
		return
	}

	err = h.service.DeleteZayavlenie(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверные данные личного дела"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
