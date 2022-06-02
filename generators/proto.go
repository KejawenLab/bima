package generators

import (
	"fmt"
	"os"
	engine "text/template"

	configs "github.com/KejawenLab/bima/v2/configs"
)

type Proto struct {
}

func (g *Proto) Generate(template *configs.Template, modulePath string, packagePath string, templatePath string) {
	workDir, _ := os.Getwd()
	protoTemplate, err := engine.ParseFiles(fmt.Sprintf("%s/%s/proto.tpl", packagePath, templatePath))
	if err != nil {
		panic(err)
	}

	protoFile, err := os.Create(fmt.Sprintf("%s/protos/%s.proto", workDir, template.ModuleLowercase))
	if err != nil {
		panic(err)
	}

	protoTemplate.Execute(protoFile, template)
}
