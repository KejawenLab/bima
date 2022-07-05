package {{.ModulePluralLowercase}}

import (
	"github.com/KejawenLab/bima/v3"
    "github.com/sarulabs/dingo/v4"
)

var Dic = []dingo.Def{
	{
		Name:  "module:{{.ModuleLowercase}}",
        Scope: bima.Application,
		Build: (*Module)(nil),
		Params: dingo.Params{
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
