package generators

import (
	"os"
	"strings"
	engine "text/template"
)

type Model struct {
}

func (g *Model) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	var path strings.Builder

	path.WriteString(packagePath)
	path.WriteString("/")
	path.WriteString(templatePath)
	path.WriteString("/model.tpl")

	modelTemplate, err := engine.ParseFiles(path.String())
	if err != nil {
		panic(err)
	}

	path.Reset()
	path.WriteString(modulePath)
	path.WriteString("/model.go")

	modelFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	err = modelTemplate.Execute(modelFile, template)
	if err != nil {
		panic(err)
	}
}
