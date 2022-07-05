package {{.ModulePluralLowercase}}

import (
	"context"
	"time"

    "github.com/KejawenLab/bima/v4/configs"
)

type {{.Module}} struct {
	configs.MongoBase `bson:",inline"`
{{range .Columns}}
    {{.Name}} {{.GolangType}} `bson:"{{.NameUnderScore}}" {{if .IsRequired}}validate:"required"{{end}}`
{{end}}
}

func (m *{{.Module}}) CollectionName() string {
	return "{{.ModuleLowercase}}"
}

func (m *{{.Module}}) Creating(context.Context) error {
	m.CreatedAt = time.Now()

	return nil
}

func (m *{{.Module}}) Updating(context.Context) error {
	m.UpdatedAt = time.Now()

	return nil
}
