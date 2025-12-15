package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/service"
	utils "online-courses/internal/server/http/handlers/handlers_utils"
	"online-courses/internal/server/http/models"
	"online-courses/internal/server/http/request"
	"online-courses/internal/validate"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EnrollmentListenerHandler struct {
	handler service.EnrollmentListenerService
}

func NewEnrollmentListenerHandler(service service.EnrollmentListenerService) *EnrollmentListenerHandler {
	return &EnrollmentListenerHandler{handler: service}
}

//	CREATE

func (e *EnrollmentListenerHandler) CreateEnrollment(c *gin.Context) {

	var request request.EnrollmentListenerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err := validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.EnrollmentListenerDTO{
		ID_Listener:      request.ID_Listener,
		ID_Program:       request.ID_ProgramEducation,
		StartDate:        request.StartDate,
		EndDate:          request.EndDate,
		CurrentPrice:     request.CurrentPrice,
		Is_active:        request.Is_active,
		Group:            request.Group,
		TypeOfRetraining: request.TypeOfRetraining,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := e.handler.Create(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	// info, err := e.handler.InfoToPersonalCard(ctx, request.ID_Listener, request.ID_ProgramEducation)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	// requestBody, _ := json.Marshal(info)

	// req, err := http.NewRequest("POST", "http://localhost:8082/v1/doc/personal-card", bytes.NewBuffer(requestBody))
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	// authHeader := c.Request.Header.Get("Authorization")
	// field := strings.Fields(authHeader)

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+field[1])

	// client := &http.Client{}
	// response, err := client.Do(req)

	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	// if response.StatusCode != http.StatusOK {

	// 	err := e.handler.Delete(ctx, request.ID_Listener, request.ID_ProgramEducation)
	// 	if err != nil {
	// 		c.Error(err)
	// 	}

	// 	c.JSON(http.StatusInternalServerError, models.HttpResponse{Message: "ошибка создания личного дела"})
	// 	return
	// }

	c.JSON(http.StatusOK, models.HttpResponse{Message: "слушатель записан на курс"})

}

//	READ

func (e *EnrollmentListenerHandler) ReadEnrollment(c *gin.Context) {

	pageParam := c.Query("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.Error(err)
		return
	}

	filter := c.Query("filter")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := e.handler.Read(ctx, page, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

//	UPDATE

func (e *EnrollmentListenerHandler) UpdateEnrollment(c *gin.Context) {

	idListener, err := utils.ParseUUID(c, "id_listener")
	if err != nil {
		c.Error(err)
		return
	}

	idProgram, err := utils.ParseUUID(c, "id_program")
	if err != nil {
		c.Error(err)
		return
	}

	var request request.EnrollemenUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	err = validate.Validate.Struct(&request)
	if err != nil {
		c.Error(err)
		return
	}

	dto := dto.EnrollmentListenerDTO{
		ID_Listener:      idListener,
		ID_Program:       request.ID_ProgramEducation,
		StartDate:        request.StartDate,
		EndDate:          request.EndDate,
		CurrentPrice:     request.CurrentPrice,
		Group:            request.Group,
		TypeOfRetraining: request.TypeOfRetraining,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := e.handler.Update(ctx, idListener, idProgram, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "запись на курс обновлена"})

}

//	DELETE

func (e *EnrollmentListenerHandler) DeleteEnrollment(c *gin.Context) {

	id_listener, err := utils.ParseUUID(c, "id_listener")
	if err != nil {
		c.Error(err)
		return
	}

	id_program, err := utils.ParseUUID(c, "id_program")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	// dto, err := e.handler.GetListenerFIO(ctx, id_listener)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	// fio := dto.SecondName + "_" + dto.FirstName + "_" + dto.MiddleName

	// url := "http://localhost:8082/v1/doc/delete?card-name=" + fio

	// req, err := http.NewRequest("DELETE", url, nil)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	// authHeader := c.Request.Header.Get("Authorization")
	// field := strings.Fields(authHeader)

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+field[1])

	// client := &http.Client{}
	// response, err := client.Do(req)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	// if response.StatusCode != http.StatusOK {
	// 	err = errors.New("личное дело не удалено ошибка сервиса документов")
	// 	c.Error(err)
	// 	return
	// }

	if err := e.handler.Delete(ctx, id_listener, id_program); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponse{Message: "запись на курс удалена"})
}

//	READ DETAIL Listener ENROLLMENT

func (e *EnrollmentListenerHandler) ReadDetailListener(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := e.handler.ReadDetailListener(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

// READ BY PROGRAM

func (e *EnrollmentListenerHandler) ReadByProgram(c *gin.Context) {
	id, err := utils.ParseUUID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	pageParam := c.Query("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	data, err := e.handler.ReadByProgram(ctx, id, page)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}

// func (e *EnrollmentListenerHandler) GetInfoToCreateCard(c *gin.Context) {
// 	var request request.CreateCardRequest

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
// 	defer cancel()

// 	data, err := e.handler.InfoToPersonalCard(ctx, request.ID_listener, request.ID_program)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, data)
// }
