package validator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("notblank", notBlank)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("notblank", notBlank)
	}
}

func V() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

func FormatErrors(err error) string {
	var messages []string

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			messages = append(messages, fmt.Sprintf("%s failed on '%s'", e.Field(), e.Tag()))
		}
		return strings.Join(messages, "; ")
	}

	return err.Error()
}

func notBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}
