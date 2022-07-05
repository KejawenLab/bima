package generators

import (
	"os"
	"strings"
	engine "text/template"

	"github.com/KejawenLab/bima/v4/parsers"
	"gopkg.in/yaml.v2"
)

type Module struct {
	Config []string `yaml:"modules"`
}

func (g *Module) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	var str strings.Builder
	str.WriteString(packagePath)
	str.WriteString("/")
	str.WriteString(templatePath)
	str.WriteString("/module.tpl")

	moduleTemplate, err := engine.ParseFiles(str.String())
	if err != nil {
		panic(err)
	}

	str.Reset()
	str.WriteString(modulePath)
	str.WriteString("/module.go")

	moduleFile, err := os.Create(str.String())
	if err != nil {
		panic(err)
	}

	str.Reset()
	str.WriteString("module:")
	str.WriteString(template.ModuleLowercase)

	workDir, _ := os.Getwd()
	g.Config = parsers.ParseModule(workDir)
	g.Config = append(g.Config, str.String())
	g.Config = g.makeUnique(g.Config)

	modules, err := yaml.Marshal(g)
	if err != nil {
		panic(err)
	}

	str.Reset()
	str.WriteString(workDir)
	str.WriteString("/")
	str.WriteString(parsers.ModulePath)

	err = os.WriteFile(str.String(), modules, 0644)
	if err != nil {
		panic(err)
	}

	moduleTemplate.Execute(moduleFile, template)
}

func (g *Module) makeUnique(modules []string) []string {
	exists := make(map[string]bool)
	var result []string
	for e := range modules {
		if exists[modules[e]] != true {
			exists[modules[e]] = true

			result = append(result, modules[e])
		}
	}

	return result
}
