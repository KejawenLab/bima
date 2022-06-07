package generators

import (
	"fmt"
	"os"
	engine "text/template"

	"github.com/KejawenLab/bima/v2/configs"
)

type Dic struct {
}

func (g *Dic) Generate(template *configs.Template, modulePath string, packagePath string, templatePath string) {
	dicTemplate, err := engine.ParseFiles(fmt.Sprintf("%s/%s/dic.tpl", packagePath, templatePath))
	if err != nil {
		panic(err)
	}

	dicPath := fmt.Sprintf("%s/dics", modulePath)
	os.MkdirAll(dicPath, 0755)

	dicFile, err := os.Create(fmt.Sprintf("%s/%s.go", dicPath, template.ModuleLowercase))
	if err != nil {
		panic(err)
	}

	dicTemplate.Execute(dicFile, template)
}
