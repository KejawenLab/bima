package generators

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

const MODULES_FILE = "swaggers/modules.json"

type Swagger struct {
}

func (g *Swagger) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	workDir, _ := os.Getwd()
	modules, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, MODULES_FILE))
	if err != nil {
		panic(err)
	}

	modulesJson := []ModuleJson{}

	json.Unmarshal(modules, &modulesJson)
	modulesJson = append(modulesJson, ModuleJson{
		Name: template.Module,
		Url:  fmt.Sprintf("./%s.swagger.json", template.ModuleLowercase),
	})

	modulesJson = g.makeUnique(modulesJson)
	for k, m := range modulesJson {
		mUrl, _ := url.Parse(m.Url)
		query := mUrl.Query()

		query.Set("v", strconv.Itoa(int(time.Now().UnixMicro())))
		mUrl.RawQuery = query.Encode()
		m.Url = mUrl.String()

		modulesJson[k] = m
	}

	modules, _ = json.Marshal(modulesJson)

	err = os.WriteFile(fmt.Sprintf("%s/%s", workDir, MODULES_FILE), modules, 0644)
	if err != nil {
		panic(err)
	}
}

func (g *Swagger) makeUnique(modules []ModuleJson) []ModuleJson {
	occured := make(map[string]bool)
	var result []ModuleJson
	for _, e := range modules {
		if occured[e.Name] != true {
			occured[e.Name] = true

			result = append(result, e)
		}
	}

	return result
}
