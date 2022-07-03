package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func Validate(v interface{}) (string, error) {
	err := validator.New().Struct(v)
	if err == nil {
		return "", nil
	}

	var message strings.Builder
	for _, ve := range err.(validator.ValidationErrors) {
		message.WriteString(strcase.ToDelimited(ve.Field(), '_'))
		message.WriteString(" is ")
		message.WriteString(ve.Tag())

		break
	}

	return message.String(), err
}
