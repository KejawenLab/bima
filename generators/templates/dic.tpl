package modules

import (
	{{.ModulePluralLowercase}} "{{.PackageName}}/{{.ModulePluralLowercase}}"
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	validations "{{.PackageName}}/{{.ModulePluralLowercase}}/validations"
	"github.com/sarulabs/dingo/v4"
)

var {{.Module}} = []dingo.Def{
	{
		Name:  "module:{{.ModuleLowercase}}:model",
		Build: (*models.{{.Module}})(nil),
	},
	{
		Name:  "module:{{.ModuleLowercase}}:validation",
		Build: (*validations.{{.Module}})(nil),
	},
	{
		Name:  "module:{{.ModuleLowercase}}",
		Build: (*{{.ModulePluralLowercase}}.Module)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("bima:context:background"),
			"Elasticsearch": dingo.Service("bima:connection:elasticsearch"),
			"Handler":       dingo.Service("bima:handler:handler"),
			"Logger":        dingo.Service("bima:handler:logger"),
			"Messenger":     dingo.Service("bima:handler:messager"),
			"Validator":     dingo.Service("module:{{.ModuleLowercase}}:validation"),
			"Cache":         dingo.Service("bima:cache:memory"),
			"Paginator":     dingo.Service("bima:pagination:paginator"),
		},
	},
	{
		Name:  "module:{{.ModuleLowercase}}:server",
		Build: (*{{.ModulePluralLowercase}}.Server)(nil),
		Params: dingo.Params{
			"Env":      dingo.Service("bima:config:env"),
			"Module":   dingo.Service("module:{{.ModuleLowercase}}"),
			"Database": dingo.Service("bima:connection:database"),
		},
	},
}
