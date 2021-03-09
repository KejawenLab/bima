package generators

import (
	"fmt"
	"os"
	engine "text/template"

	configs "github.com/crowdeco/bima/v2/configs"
)

type Model struct {
}

func (g *Model) Generate(template *configs.Template, modulePath string, packagePath string, templatePath string) {
	modelTemplate, _ := engine.ParseFiles(fmt.Sprintf("%s/%s/model.tpl", packagePath, templatePath))
	modelPath := fmt.Sprintf("%s/models", modulePath)
	os.MkdirAll(modelPath, 0755)

	modelFile, err := os.Create(fmt.Sprintf("%s/%s.go", modelPath, template.ModuleLowercase))
	if err != nil {
		panic(err)
	}

	err = modelTemplate.Execute(modelFile, template)
	if err != nil {
		panic(err)
	}
}
