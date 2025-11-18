package validate

import (
	"github.com/go-playground/validator/v10"
)

var (
	Validate *validator.Validate
)

func InitValid() {

	Validate = validator.New()

}
