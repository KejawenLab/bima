package models

import (
	"context"
	"time"

	"github.com/kamva/mgm/v3"
)

type {{.Module}} struct {
	mgm.DefaultModel `bson:",inline"`
{{range .Columns}}
    {{.Name}} {{.GolangType}} `bson:"{{.NameUnderScore}}"`
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
