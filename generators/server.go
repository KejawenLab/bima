package generators

import (
	"fmt"
	"os"
	engine "text/template"

	"github.com/KejawenLab/bima/v2/configs"
)

type Server struct {
}

func (g *Server) Generate(template *configs.Template, modulePath string, packagePath string, templatePath string) {
	serverTemplate, err := engine.ParseFiles(fmt.Sprintf("%s/%s/server.tpl", packagePath, templatePath))
	if err != nil {
		panic(err)
	}

	serverFile, err := os.Create(fmt.Sprintf("%s/server.go", modulePath))
	if err != nil {
		panic(err)
	}

	serverTemplate.Execute(serverFile, template)
}
