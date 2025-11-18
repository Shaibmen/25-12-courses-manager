package apperrors

import "errors"

var (
	ErrUnique         = errors.New("unique violation error")      //409
	ErrForeignKey     = errors.New("foreign key violation error") //400
	ErrNotNull        = errors.New("not null violation error")    //400
	ErrCheckViolation = errors.New("check violation error")       //422
	ErrNoTable        = errors.New("table not found")             //404
	ErrNoRow          = errors.New("rows not found")              //404
	ErrEmptyData      = errors.New("data not found")              //404
	ErrValidate       = errors.New("validation error")            //400
	ErrNoExists       = errors.New("data no exists")              //400
)
