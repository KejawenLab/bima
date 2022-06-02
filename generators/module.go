package generators

import (
	"fmt"
	"io/ioutil"
	"os"
	engine "text/template"

	configs "github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/parsers"
	"gopkg.in/yaml.v2"
)

type Module struct {
	Config *parsers.Module
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

	g.Config.Parse(workDir)
	g.Config.Config = append(g.Config.Config, fmt.Sprintf("module:%s", template.ModuleLowercase))
	g.Config.Config = g.makeUnique(g.Config.Config)

	modules, err := yaml.Marshal(g.Config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", workDir, parsers.MODULES_FILE), modules, 0644)
	if err != nil {
		panic(err)
	}

	moduleTemplate.Execute(moduleFile, template)
}

func (g *Module) makeUnique(modules []string) []string {
	occured := make(map[string]bool)
	var result []string
	for e := range modules {
		if occured[modules[e]] != true {
			occured[modules[e]] = true

			result = append(result, modules[e])
		}
	}

	return result
}
