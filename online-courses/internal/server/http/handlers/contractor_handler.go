package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"strconv"
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

	passportSeria, err := strconv.Atoi(request.Passport.Seria)
	if err != nil {
		c.Error(err)
		return
	}

	passportNumber, err := strconv.Atoi(request.Passport.Number)
	if err != nil {
		c.Error(err)
		return
	}

	mainIndex, err := strconv.Atoi(request.RegAddress.MailIndex)
	if err != nil {
		c.Error(err)
		return
	}

	dto := &dto.ContractorCreateDTO{
		Contractor: dto.ContractorDTO{
			FirstName:     request.Contractor.FirstName,
			SecondName:    request.Contractor.SecondName,
			MiddleName:    request.Contractor.MiddleName,
			Contact_phone: request.Contractor.Contact_phone,
			Email:         request.Contractor.Email,
		},
		Passport: dto.PassportDTO{
			PlaceBirth:    request.Passport.PlaceBirth,
			Citizenship:   request.Passport.Citizenship,
			Gender:        request.Passport.Gender,
			Seria:         passportSeria,
			Number:        passportNumber,
			PassportGiven: request.Passport.PassportGiven,
			DateGiven:     request.Passport.DateGiven,
			Code:          request.Passport.Code,
		},
		RegAddress: dto.RegistrationAddressDTO{
			MailIndex: mainIndex,
			Region:    request.RegAddress.Region,
			City:      request.RegAddress.City,
			Street:    request.RegAddress.Street,
			House:     request.RegAddress.House,
			Building:  request.RegAddress.Building,
			Apartment: request.RegAddress.Apartment,
		},
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err = con.handler.Create(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "заказчик создан"})
}
