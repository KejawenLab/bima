package generators

import (
	"fmt"
	"os"
	engine "text/template"
)

type Dic struct {
}

func (g *Dic) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	dicTemplate, err := engine.ParseFiles(fmt.Sprintf("%s/%s/dic.tpl", packagePath, templatePath))
	if err != nil {
		panic(err)
	}

	dicFile, err := os.Create(fmt.Sprintf("%s/dic.go", modulePath))
	if err != nil {
		panic(err)
	}

	dicTemplate.Execute(dicFile, template)
}
