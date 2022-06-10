package generators

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
)

const MODULES_FILE = "swaggers/modules.json"

type Swagger struct {
}

func (g *Swagger) Generate(template *configs.Template, modulePath string, packagePath string, templatePath string) {
	workDir, _ := os.Getwd()
	modules, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, MODULES_FILE))
	if err != nil {
		panic(err)
	}

	modulesJson := []configs.ModuleJson{}

	json.Unmarshal(modules, &modulesJson)
	modulesJson = append(modulesJson, configs.ModuleJson{
		Name: template.Module,
		Url:  fmt.Sprintf("./%s.swagger.json", template.ModuleLowercase),
	})

	modulesJson = g.makeUnique(modulesJson)
	for k, m := range modulesJson {
		m.Url = fmt.Sprintf("%s?v=%s", m.Url, time.Now().Format(time.RFC3339Nano))
		modulesJson[k] = m
	}

	modules, _ = json.Marshal(modulesJson)

	err = os.WriteFile(fmt.Sprintf("%s/%s", workDir, MODULES_FILE), modules, 0644)
	if err != nil {
		panic(err)
	}
}

func (g *Swagger) makeUnique(modules []configs.ModuleJson) []configs.ModuleJson {
	occured := make(map[string]bool)
	var result []configs.ModuleJson
	for _, e := range modules {
		if occured[e.Name] != true {
			occured[e.Name] = true

			result = append(result, e)
		}
	}

	return result
}
