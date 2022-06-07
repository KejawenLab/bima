package validations

import (
	"{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	"github.com/go-ozzo/ozzo-validation/v4"
)

type {{.Module}} struct{}

func (v *{{.Module}}) Validate(m *models.{{.Module}}) (bool, error) {
	err := validation.ValidateStruct(m,
    {{range .Columns}}
        {{if .IsRequired}}
        validation.Field(&m.{{.Name}}, validation.Required),
        {{end}}
    {{end}}
	)

	if err != nil {
		return false, err
	}

	return true, nil
}
