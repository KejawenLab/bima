package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func Validate(v interface{}) (string, error) {
	err := validator.New().Struct(v)
	if err == nil {
		return "", nil
	}

	var message string
	for _, ve := range err.(validator.ValidationErrors) {
		message = fmt.Sprintf("%s is %s", strcase.ToDelimited(ve.Field(), '_'), ve.Tag())
		break
	}

	return message, err
}
