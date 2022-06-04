package models

import "github.com/KejawenLab/bima/v2/configs"

type {{.Module}} struct {
	configs.GormBase
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
