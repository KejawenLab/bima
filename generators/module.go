package generators

import (
	"fmt"
	"io/ioutil"
	"os"
	engine "text/template"

	configs "github.com/crowdeco/bima/configs"
	"github.com/crowdeco/bima/parsers"
	"gopkg.in/yaml.v2"
)

type Module struct {
	Config *parsers.Module
}

func (g *Module) Generate(template *configs.Template, modulePath string, packagePath string, templatePath string) {
	moduleTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/module.tpl", packagePath, templatePath))
	moduleFile, err := os.Create(fmt.Sprintf("%s/module.go", modulePath))
	if err != nil {
		panic(err)
	}

	g.Config.Parse()
	g.Config.Config = append(g.Config.Config, fmt.Sprintf("module:%s", template.ModuleLowercase))
	g.Config.Config = g.makeUnique(g.Config.Config)

	modules, err := yaml.Marshal(g.Config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(parsers.MODULES_FILE, modules, 0644)
	if err != nil {
		panic(err)
	}

	moduleTemplate.Execute(moduleFile, template)
}

func (g *Module) makeUnique(slices []string) []string {
	occured := make(map[string]bool)
	var result []string
	for e := range slices {
		if occured[slices[e]] != true {
			occured[slices[e]] = true

			result = append(result, slices[e])
		}
	}

	return result
}
