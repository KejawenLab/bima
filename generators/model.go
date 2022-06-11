package generators

import (
	"fmt"
	"os"
	engine "text/template"
)

type Model struct {
}

func (g *Model) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	modelTemplate, err := engine.ParseFiles(fmt.Sprintf("%s/%s/model.tpl", packagePath, templatePath))
	if err != nil {
		panic(err)
	}

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
