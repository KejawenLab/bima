package models

import "github.com/KejawenLab/bima/v2"

type {{.Module}} struct {
	mgm.DefaultModel `bson:",inline"`
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
