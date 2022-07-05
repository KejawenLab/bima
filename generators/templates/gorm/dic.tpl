package {{.ModulePluralLowercase}}

import (
	"github.com/KejawenLab/bima/v4"
    "github.com/sarulabs/dingo/v4"
)

var Dic = []dingo.Def{
	{
		Name:  "module:{{.ModuleLowercase}}:model",
        Scope: bima.Application,
		Build: (*{{.Module}})(nil),
        Params: dingo.Params{
			"GormModel": dingo.Service("bima:model"),
		},
	},
	{
		Name:  "module:{{.ModuleLowercase}}",
        Scope: bima.Application,
		Build: (*Module)(nil),
		Params: dingo.Params{
            "Model":  dingo.Service("module:{{.ModuleLowercase}}:model"),
			"Module": dingo.Service("bima:module"),
		},
	},
	{
		Name:  "module:{{.ModuleLowercase}}:server",
        Scope: bima.Application,
		Build: (*Server)(nil),
		Params: dingo.Params{
			"Server": dingo.Service("bima:server"),
			"Module": dingo.Service("module:{{.ModuleLowercase}}"),
		},
	},
}
