package generators

import (
	"os"
	"strings"
	engine "text/template"
)

type Dic struct {
}

func (g *Dic) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	var path strings.Builder

	path.WriteString(packagePath)
	path.WriteString("/")
	path.WriteString(templatePath)
	path.WriteString("/dic.tpl")

	dicTemplate, err := engine.ParseFiles(path.String())
	if err != nil {
		panic(err)
	}

	path.Reset()
	path.WriteString(modulePath)
	path.WriteString("/dic.go")

	dicFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	dicTemplate.Execute(dicFile, template)
}
