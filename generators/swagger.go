package generators

import (
	"bytes"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"
)

type Swagger struct {
}

func (g *Swagger) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	workDir, _ := os.Getwd()
	var path bytes.Buffer

	path.WriteString(workDir)
	path.WriteString("/")
	path.WriteString("swaggers/modules.json")

	modules, err := os.ReadFile(path.String())
	if err != nil {
		panic(err)
	}

	modulesJson := []ModuleJson{}

	json.Unmarshal(modules, &modulesJson)

	path.Reset()
	path.WriteString("./")
	path.WriteString(template.ModuleLowercase)
	path.WriteString(".swagger.json")

	modulesJson = append(modulesJson, ModuleJson{
		Name: template.Module,
		Url:  path.String(),
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

	path.Reset()
	path.WriteString(workDir)
	path.WriteString("/")
	path.WriteString("swaggers/modules.json")

	err = os.WriteFile(path.String(), modules, 0644)
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
