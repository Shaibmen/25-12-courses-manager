package handlers_utils

import (
	"errors"
	"net/http"
	"online-courses/internal/apperrors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  int
	Message string
}

func DefineError(err error) Response {

	if err == nil {
		return Response{Status: 200, Message: "запрос выполнен"}
	}

	var (
		validationErr validator.ValidationErrors
		strconvErr    *strconv.NumError
	)

	switch {

	case errors.Is(err, apperrors.ErrUnique):
		return Response{
			Status:  http.StatusConflict,
			Message: "запись c такими данными уже существует",
		}

	case errors.Is(err, apperrors.ErrForeignKey):
		return Response{
			Status:  http.StatusBadRequest,
			Message: "указанна недействительная ссылка",
		}
	case errors.Is(err, apperrors.ErrNotNull):
		return Response{
			Status:  http.StatusBadRequest,
			Message: "обязательное поле не указано",
		}
	case errors.Is(err, apperrors.ErrCheckViolation):
		return Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "данные не прошли проверку",
		}
	case errors.Is(err, apperrors.ErrNoTable):
		return Response{
			Status:  http.StatusNotFound,
			Message: "внутренняя ошибка сервера",
		}
	case errors.Is(err, apperrors.ErrNoRow):
		return Response{
			Status:  http.StatusNotFound,
			Message: "данные не найдены",
		}
	case errors.Is(err, apperrors.ErrEmptyData):
		return Response{
			Status:  http.StatusNotFound,
			Message: "нет данных для отображения",
		}

	case errors.As(err, &validationErr):
		return Response{
			Status:  http.StatusBadRequest,
			Message: "ошибка валидации данных",
		}

	case errors.As(err, &strconvErr):
		return Response{
			Status:  http.StatusBadRequest,
			Message: "неверные данные запроса",
		}

	case strings.Contains(err.Error(), "invalid UUID format"):
		return Response{
			Status:  http.StatusBadRequest,
			Message: "неверные данные запроса id",
		}
	case strings.Contains(err.Error(), "invalid UUID length"):
		return Response{
			Status:  http.StatusBadRequest,
			Message: "неверные данные запроса id",
		}
	case strings.Contains(err.Error(), "json: cannot unmarshal"):
		return Response{
			Status:  http.StatusBadRequest,
			Message: "неправильный тип данных",
		}
	case errors.Is(err, apperrors.ErrNoExists):
		return Response{
			Status:  http.StatusNotFound,
			Message: "нет данных для отображения",
		}
	}

	return Response{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

}
