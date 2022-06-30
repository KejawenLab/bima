package {{.ModulePluralLowercase}}

import "github.com/sarulabs/dingo/v4"

var Dic = []dingo.Def{
	{
		Name:  "module:{{.ModuleLowercase}}:model",
		Build: (*{{.Module}})(nil),
	},
	{
		Name:  "module:{{.ModuleLowercase}}",
		Build: (*Module)(nil),
		Params: dingo.Params{
			"Module": dingo.Service("bima:module"),
		},
	},
	{
		Name:  "module:{{.ModuleLowercase}}:server",
		Build: (*Server)(nil),
		Params: dingo.Params{
			"Server": dingo.Service("bima:server"),
			"Module": dingo.Service("module:{{.ModuleLowercase}}"),
		},
	},
}
