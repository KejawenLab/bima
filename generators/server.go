package generators

import (
	"os"
	"strings"
	engine "text/template"
)

type Server struct {
}

func (g *Server) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	var path strings.Builder

	path.WriteString(packagePath)
	path.WriteString("/")
	path.WriteString(templatePath)
	path.WriteString("/server.tpl")

	serverTemplate, err := engine.ParseFiles(path.String())
	if err != nil {
		panic(err)
	}

	path.Reset()
	path.WriteString(modulePath)
	path.WriteString("/server.go")

	serverFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	serverTemplate.Execute(serverFile, template)
}
