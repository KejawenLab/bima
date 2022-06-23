package generators

import (
	"fmt"
	"os"
	"strings"
)

const MODULE_IMPORT = "@modules:import"
const MODULE_REGISTER = "@modules:register"

type Provider struct {
}

func (p *Provider) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	workDir, _ := os.Getwd()
	path := fmt.Sprintf("%s/configs/provider.go", workDir)

	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	contents := strings.Split(string(file), "\n")
	importIdx := 0
	moduleIdx := 0
	skipImport := true

	for k, v := range contents {
		if strings.Contains(v, MODULE_IMPORT) {
			importIdx = k
			skipImport = false
			continue
		}

		if strings.Contains(v, MODULE_REGISTER) {
			moduleIdx = k
			break
		}
	}

	if !skipImport {
		contents[importIdx] = fmt.Sprintf(`    //%s
    %s %q`, MODULE_IMPORT, template.ModuleLowercase, fmt.Sprintf("%s/%s", template.PackageName, template.ModulePluralLowercase))
	}

	contents[moduleIdx] = fmt.Sprintf(`
    /*@module:%s*/if err := p.AddDefSlice(%s.Dic); err != nil {return err}
    //%s`, template.ModuleLowercase, template.ModuleLowercase, MODULE_REGISTER)

	body := strings.Join(contents, "\n")

	os.WriteFile(path, []byte(body), 0644)
}
