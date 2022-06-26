package generators

import (
	"bytes"
	"os"
	engine "text/template"
)

type Proto struct {
}

func (g *Proto) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	var path bytes.Buffer

	path.WriteString(packagePath)
	path.WriteString("/")
	path.WriteString(templatePath)
	path.WriteString("/proto.tpl")

	protoTemplate, err := engine.ParseFiles(path.String())
	if err != nil {
		panic(err)
	}

	workDir, _ := os.Getwd()

	path.Reset()
	path.WriteString(workDir)
	path.WriteString("/protos/")
	path.WriteString(template.ModuleLowercase)
	path.WriteString(".proto")

	protoFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	protoTemplate.Execute(protoFile, template)
}
