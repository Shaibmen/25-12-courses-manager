package handler

import (
	"document-service/internal/mapper"
	"document-service/internal/server/models"
	"document-service/internal/service"
	"document-service/internal/validate"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	cardsPath = "./personal_card"
)

type PersonalCardHandler struct {
	service *service.PersonalCardService
}

func NewPersonalCardHandler(service *service.PersonalCardService) *PersonalCardHandler {
	return &PersonalCardHandler{service: service}
}

func (p *PersonalCardHandler) CreatePersonalCard(c *gin.Context) {
	var request models.FullListenerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter request", "err": err.Error()})
		return
	}

	err := validate.Validator.Struct(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter validation", "err": err.Error()})
		return
	}

	dto, err := mapper.FullListenerMapping(request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter request mapping", "err": err.Error()})
		return
	}

	err = p.service.CreatePersonalCard(dto)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "creation listener card error", "err": err.Error()})
		return
	}

}

func (p *PersonalCardHandler) ExistsPersonalCard(c *gin.Context) {
	param := c.Query("card-name")

	files, err := p.service.ExistsPersonalCard(param)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "личного дела не существует"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (p *PersonalCardHandler) DownloadPersonalCard(c *gin.Context) {
	// filePath := cardsPath + "/" + param
	//
	// c.File(filePath)
	param := c.Query("card-name")

	doc, err := p.service.DownloadPersonalCard(param)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusNotFound, gin.H{"message": "не удалось загрузить файл"})
		log.Println("ошибка при загрузке файла:", err)
		return
	}

	// log.Println("Doc size:", len(doc))

	c.Data(200, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", doc)
}

func (p *PersonalCardHandler) DeletePersonalCard(c *gin.Context) {
	param := c.Query("card-name")

	file := "Личное-дело-" + param + ".docx"
	filePath := cardsPath + "/" + file
	_, err := p.service.ExistsPersonalCard(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "личного дела не существует"})
		return
	}

	err = p.service.DeletePersonalCard(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверные данные личного дела"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
