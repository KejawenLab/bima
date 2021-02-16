package models

import (
	bima "github.com/crowdeco/bima"
)

type {{.Module}} struct {
	*bima.Model
{{range .Columns}}
    {{.Name}} {{.GolangType}}
{{end}}
}

func (m *{{.Module}}) TableName() string {
	return "{{.ModuleLowercase}}"
}

func (m *{{.Module}}) IsSoftDelete() bool {
	return true
}
