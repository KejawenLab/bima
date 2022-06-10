package generators

import (
	"fmt"
	"os"
	engine "text/template"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/parsers"
	"gopkg.in/yaml.v2"
)

type Module struct {
	Config []string `yaml:"modules"`
}

func (g *Module) Generate(template *configs.Template, modulePath string, packagePath string, templatePath string) {
	workDir, _ := os.Getwd()
	moduleTemplate, err := engine.ParseFiles(fmt.Sprintf("%s/%s/module.tpl", packagePath, templatePath))
	if err != nil {
		panic(err)
	}

	moduleFile, err := os.Create(fmt.Sprintf("%s/module.go", modulePath))
	if err != nil {
		panic(err)
	}

	g.Config = parsers.ParseModule(workDir)
	g.Config = append(g.Config, fmt.Sprintf("module:%s", template.ModuleLowercase))
	g.Config = g.makeUnique(g.Config)

	modules, err := yaml.Marshal(g)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s", workDir, parsers.MODULES_FILE), modules, 0644)
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
