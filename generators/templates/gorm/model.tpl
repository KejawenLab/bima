package models

import "github.com/KejawenLab/bima/v3"

type {{.Module}} struct {
	*bima.GormModel
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
