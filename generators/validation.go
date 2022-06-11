package generators

import (
	"fmt"
	"os"
	engine "text/template"
)

type Validation struct {
}

func (g *Validation) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	validationTemplate, err := engine.ParseFiles(fmt.Sprintf("%s/%s/validation.tpl", packagePath, templatePath))
	if err != nil {
		panic(err)
	}

	validationPath := fmt.Sprintf("%s/validations", modulePath)
	os.MkdirAll(validationPath, 0755)

	validationFile, err := os.Create(fmt.Sprintf("%s/%s.go", validationPath, template.ModuleLowercase))
	if err != nil {
		panic(err)
	}

	validationTemplate.Execute(validationFile, template)
}
