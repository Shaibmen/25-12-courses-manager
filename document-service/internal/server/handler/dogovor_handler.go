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

type DogovoreHandler struct {
	service *service.DogovorService
}

func NewDogovorHandler(service *service.DogovorService) *DogovoreHandler {
	return &DogovoreHandler{service}
}

func (h *DogovoreHandler) CreateDogovor(c *gin.Context) {
	var fullRequest models.FullRequest

	if err := c.ShouldBindJSON(&fullRequest); err != nil {
		log.Println("пиздец при шудбинджсон!", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter request", "err": err.Error()})
		return
	}

	request := fullRequest.DogovorData

	log.Println("request:", request)

	dto, err := mapper.DogovorMapping(request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter request mapping", "err": err.Error()})
		return
	}

	err = h.service.CreateDogovor(dto, request.DogovorType)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "creation dogovor error", "err": err.Error()})
		return
	}
}

func (h *DogovoreHandler) ExistsDogovor(c *gin.Context) {
	param := c.Query("dogovor-name")

	files, err := h.service.ExistsDogovor(param)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "личного дела не существует"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *DogovoreHandler) DownloadDogovor(c *gin.Context) {
	// filePath := cardsPath + "/" + param
	//
	// c.File(filePath)
	param := c.Query("dogovor-name")

	doc, err := h.service.DownloadDogovor(param)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusNotFound, gin.H{"message": "не удалось загрузить файл"})
		log.Println("ошибка при загрузке файла:", err)
		return
	}

	// log.Println("Doc size:", len(doc))

	c.Data(200, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", doc)
}

func (h *DogovoreHandler) DeleteDogovor(c *gin.Context) {
	param := c.Query("dogovor-name")

	file := "Заявление-" + param + ".docx"
	filePath := cardsPath + "/" + file
	_, err := h.service.ExistsDogovor(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "личного дела не существует"})
		return
	}

	err = h.service.DeleteDogovor(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверные данные личного дела"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
