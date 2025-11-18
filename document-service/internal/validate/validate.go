package validate

import "github.com/go-playground/validator/v10"

var (
	Validator *validator.Validate
)

func InitValid() {
	valid := validator.New()

	Validator = valid
}
