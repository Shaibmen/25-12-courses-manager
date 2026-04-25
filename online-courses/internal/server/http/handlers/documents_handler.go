package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"online-courses/internal/config"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	"online-courses/internal/server/http/request"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	handler service.DocumentService
	cfg     *config.Config
}

func NewDocumentHandler(handler service.DocumentService, cfg *config.Config) *DocumentHandler {
	return &DocumentHandler{handler: handler, cfg: cfg}
}

func (d *DocumentHandler) DocumentDataDeliver(c *gin.Context) {

	var request request.DocumentsDataRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("ошибка в маппинге структуры снова")
		c.Error(err)
		return
	}

	fmt.Println(request)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	data, err := d.handler.FormingDataDocument(ctx, request.ID_Listener, request.ID_Program, request.ID_Executor, request.FrontData)
	if err != nil {
		c.Error(err)
		return
	}

	log.Println("REQ id_program:", request.ID_Program, "id_listener:", request.ID_Listener)
	log.Println("FORM division:", data.DogovorRequest.ProgramEducation.DivisionEducation)
	log.Println("BRANCH company:", data.DogovorRequest.LegalEntity.CompanyName)

	if data.DogovorRequest.LegalEntity.CompanyName != "" {

		data.DogovorRequest.Enrollment.StartDate = request.FrontData.StartDate
		data.DogovorRequest.Enrollment.EndDate = request.FrontData.EndDate
		data.DogovorRequest.Enrollment.CurrentPrice = request.FrontData.CurrentPrice
		data.DogovorRequest.ProgramEducation.NameProfEducation = request.FrontData.NameProfEducation
		data.DogovorRequest.ProgramEducation.TimeEducation = request.FrontData.TimeEducation
		data.DogovorRequest.OptionNagruzka = request.FrontData.OptionNagruzka
		data.ZayavlenieCardInfo.Variant = request.FrontData.Variant

		responseDogovor, err := RequestToDoc(*data, "dogovor", c, d.cfg)
		if err != nil {
			c.Error(err)
			return
		}

		if responseDogovor != http.StatusOK {
			c.JSON(responseDogovor, nil)
			return
		}

	} else {
		responseCard, err := RequestToDoc(*data, "personal-card", c, d.cfg)
		if err != nil {
			c.Error(err)
			return
		}
		if responseCard != http.StatusOK {
			c.JSON(responseCard, nil)
			return
		}

		responseZayvlenie, err := RequestToDoc(*data, "zayavlenie", c, d.cfg)
		if err != nil {
			c.Error(err)
			return
		}
		if responseZayvlenie != http.StatusOK {
			c.JSON(responseZayvlenie, nil)
			return
		}

		responseDogovor, err := RequestToDoc(*data, "dogovor", c, d.cfg)
		if err != nil {
			c.Error(err)
			return
		}

		if responseDogovor != http.StatusOK {
			c.JSON(responseDogovor, nil)
			return
		}
	}

}

func RequestToDoc(data dto.FullDocumentInfoDTO, endpoint string, c *gin.Context, cfg *config.Config) (int, error) {
	requestBody, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", cfg.SERVICE_DOC+":8082/v1/doc/"+endpoint, bytes.NewBuffer(requestBody))
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
